package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAssembleExamples assembles every shipped example program end to end
// through the full pipeline (preprocess, parse, symbols, codegen, output)
// in every word size and output format.
func TestAssembleExamples(t *testing.T) {
	examples, err := filepath.Glob("../examples/*.asm")
	if err != nil || len(examples) == 0 {
		t.Fatalf("no examples found: %v", err)
	}

	// Silence the pipeline's stdout chatter during the test.
	old := os.Stdout
	devnull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	os.Stdout = devnull
	defer func() { os.Stdout = old; devnull.Close() }()

	for _, example := range examples {
		name := filepath.Base(example)
		if name == "io_defs.asm" {
			continue // include file, not a standalone program
		}
		for _, wordSize := range []string{"36", "108"} {
			for _, format := range []string{"binary", "hex", "objdump"} {
				t.Run(name+"/"+wordSize+"/"+format, func(t *testing.T) {
					out := filepath.Join(t.TempDir(), "out.bin")
					if err := RunAssembler(example, out, "", format, false, "", wordSize); err != nil {
						t.Fatalf("assembling %s (wordsize=%s, format=%s): %v", name, wordSize, format, err)
					}
					data, err := ioutil.ReadFile(out)
					if err != nil {
						t.Fatal(err)
					}
					if len(data) == 0 {
						t.Fatalf("assembling %s produced an empty output file", name)
					}
					if format == "hex" && !strings.HasSuffix(strings.TrimSpace(string(data)), ":00000001FF") {
						t.Errorf("hex output for %s is missing the EOF record", name)
					}
				})
			}
		}
	}
}
