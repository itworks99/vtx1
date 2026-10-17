package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempAsm(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestPreprocessPassthrough(t *testing.T) {
	dir := t.TempDir()
	src := "main:\n        ADD T0, T1, T2\n        WFI\n"
	path := writeTempAsm(t, dir, "plain.asm", src)
	out, sourceMap, err := Preprocess(path)
	if err != nil {
		t.Fatal(err)
	}
	if out != src {
		t.Errorf("passthrough changed source:\n%q\nwant:\n%q", out, src)
	}
	if _, ok := sourceMap[path]; !ok {
		t.Error("source map is missing the input file")
	}
}

func TestPreprocessInclude(t *testing.T) {
	dir := t.TempDir()
	writeTempAsm(t, dir, "defs.asm", "        .EQU IOBASE, 0x2000\n")
	main := writeTempAsm(t, dir, "main.asm",
		"        .INCLUDE \"defs.asm\"\n        LD T0, IOBASE\n")
	out, _, err := Preprocess(main)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, ".EQU IOBASE, 0x2000") {
		t.Errorf("include not inlined:\n%s", out)
	}
	if strings.Contains(out, ".INCLUDE") {
		t.Errorf(".INCLUDE directive not removed:\n%s", out)
	}
}

func TestPreprocessIncludeIsRelativeToIncluder(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "lib")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	writeTempAsm(t, sub, "inner.asm", "        NOP\n")
	writeTempAsm(t, sub, "outer.asm", "        .INCLUDE \"inner.asm\"\n")
	main := writeTempAsm(t, dir, "main.asm", "        .INCLUDE \"lib/outer.asm\"\n")
	out, _, err := Preprocess(main)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "NOP") {
		t.Errorf("nested include not resolved:\n%s", out)
	}
}

func TestPreprocessCircularInclude(t *testing.T) {
	dir := t.TempDir()
	writeTempAsm(t, dir, "a.asm", "        .INCLUDE \"b.asm\"\n")
	writeTempAsm(t, dir, "b.asm", "        .INCLUDE \"a.asm\"\n")
	_, _, err := Preprocess(filepath.Join(dir, "a.asm"))
	if err == nil || !strings.Contains(err.Error(), "circular") {
		t.Errorf("circular include should fail, got: %v", err)
	}
}

func TestPreprocessMacroExpansion(t *testing.T) {
	dir := t.TempDir()
	main := writeTempAsm(t, dir, "main.asm", strings.Join([]string{
		".MACRO STORE2 src addr",
		"        ST src, addr",
		"        ST src, addr",
		".ENDM",
		"start:",
		"        STORE2 T1, 0x2000",
		"        WFI",
		"",
	}, "\n"))
	out, _, err := Preprocess(main)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, ".MACRO") || strings.Contains(out, ".ENDM") {
		t.Errorf("macro definition not stripped:\n%s", out)
	}
	if got := strings.Count(out, "ST T1, 0x2000"); got != 2 {
		t.Errorf("macro body expanded %d times, want 2:\n%s", got, out)
	}
	if !strings.Contains(out, "start:") {
		t.Errorf("surrounding code lost:\n%s", out)
	}
}

func TestPreprocessMacroWithLabelOnInvocation(t *testing.T) {
	dir := t.TempDir()
	main := writeTempAsm(t, dir, "main.asm", strings.Join([]string{
		".MACRO ZERO r",
		"        LD r, 0",
		".ENDM",
		"init: ZERO T3",
		"",
	}, "\n"))
	out, _, err := Preprocess(main)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "init:") || !strings.Contains(out, "LD T3, 0") {
		t.Errorf("labelled invocation mis-expanded:\n%s", out)
	}
}

func TestPreprocessMacroInvokesMacro(t *testing.T) {
	dir := t.TempDir()
	main := writeTempAsm(t, dir, "main.asm", strings.Join([]string{
		".MACRO ZERO r",
		"        LD r, 0",
		".ENDM",
		".MACRO ZERO2 a b",
		"        ZERO a",
		"        ZERO b",
		".ENDM",
		"        ZERO2 T1, T2",
		"",
	}, "\n"))
	out, _, err := Preprocess(main)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "LD T1, 0") || !strings.Contains(out, "LD T2, 0") {
		t.Errorf("nested macro invocation failed:\n%s", out)
	}
}

func TestPreprocessMacroArgumentCountMismatch(t *testing.T) {
	dir := t.TempDir()
	main := writeTempAsm(t, dir, "main.asm", strings.Join([]string{
		".MACRO ZERO r",
		"        LD r, 0",
		".ENDM",
		"        ZERO T1, T2",
		"",
	}, "\n"))
	_, _, err := Preprocess(main)
	if err == nil || !strings.Contains(err.Error(), "argument") {
		t.Errorf("argument count mismatch should fail, got: %v", err)
	}
}

func TestPreprocessRecursiveMacro(t *testing.T) {
	dir := t.TempDir()
	main := writeTempAsm(t, dir, "main.asm", strings.Join([]string{
		".MACRO LOOPY",
		"        LOOPY",
		".ENDM",
		"        LOOPY",
		"",
	}, "\n"))
	_, _, err := Preprocess(main)
	if err == nil || !strings.Contains(err.Error(), "depth") {
		t.Errorf("recursive macro should hit the depth limit, got: %v", err)
	}
}

func TestPreprocessMissingEndm(t *testing.T) {
	dir := t.TempDir()
	main := writeTempAsm(t, dir, "main.asm", ".MACRO BROKEN\n        NOP\n")
	_, _, err := Preprocess(main)
	if err == nil || !strings.Contains(err.Error(), ".ENDM") {
		t.Errorf("unterminated macro should fail, got: %v", err)
	}
}
