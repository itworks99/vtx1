package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

// Preprocessor
// ============
// Runs before lexing/parsing and handles the two source-rewriting features:
//
//   .INCLUDE "file"            textually inlines the file (paths are resolved
//                              relative to the including file; include cycles
//                              are detected and reported)
//
//   .MACRO name [p1 p2 ...]    records a macro; an invocation line
//   ...body...                 `name arg1, arg2` expands to the body with
//   .ENDM                      each parameter replaced by its argument
//
// Macro names are case-insensitive. Expansion is recursive (macros may invoke
// macros and included files may define or invoke them) with a depth limit to
// catch accidental self-recursion.

const maxMacroDepth = 64

type macroDef struct {
	Name   string
	Params []string
	Body   []string
	File   string
	Line   int
}

type preprocessor struct {
	macros       map[string]*macroDef // keyed by upper-case name
	includeStack []string             // absolute paths of files being processed
	sourceMap    map[string]string    // file -> original content, for error context
	out          []string
}

// Preprocess reads inputFile, resolves .INCLUDE directives and expands macros,
// and returns the flattened source plus a map of every file read.
func Preprocess(inputFile string) (string, map[string]string, error) {
	p := &preprocessor{
		macros:    make(map[string]*macroDef),
		sourceMap: make(map[string]string),
	}
	if err := p.processFile(inputFile, 0); err != nil {
		return "", p.sourceMap, err
	}
	src := strings.Join(p.out, "\n")
	if len(src) > 0 {
		src += "\n"
	}
	return src, p.sourceMap, nil
}

func (p *preprocessor) processFile(path string, depth int) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	for _, active := range p.includeStack {
		if active == abs {
			return fmt.Errorf("circular include detected: %s is already being included (chain: %s)",
				path, strings.Join(p.includeStack, " -> "))
		}
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read include file %s: %v", path, err)
	}
	p.sourceMap[path] = string(content)
	p.includeStack = append(p.includeStack, abs)
	defer func() { p.includeStack = p.includeStack[:len(p.includeStack)-1] }()

	return p.processLines(splitLines(string(content)), path, depth)
}

var includeRe = regexp.MustCompile(`(?i)^\s*\.INCLUDE\s+"([^"]*)"\s*(;.*)?$`)
var macroStartRe = regexp.MustCompile(`(?i)^\s*\.MACRO\s+([A-Za-z_][A-Za-z0-9_]*)\s*(.*?)\s*(;.*)?$`)
var endmRe = regexp.MustCompile(`(?i)^\s*\.ENDM\s*(;.*)?$`)
var invocationRe = regexp.MustCompile(`^\s*(?:([A-Za-z_][A-Za-z0-9_]*)\s*:\s*)?([A-Za-z_][A-Za-z0-9_]*)\s*(.*?)\s*(;.*)?$`)

func (p *preprocessor) processLines(lines []string, file string, depth int) error {
	if depth > maxMacroDepth {
		return fmt.Errorf("macro expansion exceeds maximum depth of %d in %s (recursive macro?)", maxMacroDepth, file)
	}
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// .INCLUDE "file"
		if m := includeRe.FindStringSubmatch(line); m != nil {
			incPath := m[1]
			if !filepath.IsAbs(incPath) {
				incPath = filepath.Join(filepath.Dir(file), incPath)
			}
			if err := p.processFile(incPath, depth); err != nil {
				return fmt.Errorf("%s:%d: %v", file, i+1, err)
			}
			continue
		}

		// .MACRO name [params...] ... .ENDM
		if m := macroStartRe.FindStringSubmatch(line); m != nil {
			def := &macroDef{Name: m[1], File: file, Line: i + 1}
			if params := strings.Fields(m[2]); len(params) > 0 {
				def.Params = params
			}
			end := -1
			for j := i + 1; j < len(lines); j++ {
				if endmRe.MatchString(lines[j]) {
					end = j
					break
				}
				if macroStartRe.MatchString(lines[j]) {
					return fmt.Errorf("%s:%d: nested .MACRO definitions are not supported", file, j+1)
				}
				def.Body = append(def.Body, lines[j])
			}
			if end < 0 {
				return fmt.Errorf("%s:%d: .MACRO '%s' has no matching .ENDM", file, i+1, def.Name)
			}
			key := strings.ToUpper(def.Name)
			if prev, exists := p.macros[key]; exists {
				return fmt.Errorf("%s:%d: macro '%s' already defined at %s:%d", file, i+1, def.Name, prev.File, prev.Line)
			}
			p.macros[key] = def
			i = end
			continue
		}

		// Macro invocation: `[label:] NAME [args]` where NAME is a known macro.
		if m := invocationRe.FindStringSubmatch(line); m != nil {
			if def, ok := p.macros[strings.ToUpper(m[2])]; ok {
				if m[1] != "" {
					p.out = append(p.out, m[1]+":")
				}
				expanded, err := expandMacro(def, m[3], file, i+1)
				if err != nil {
					return err
				}
				if err := p.processLines(expanded, file, depth+1); err != nil {
					return err
				}
				continue
			}
		}

		p.out = append(p.out, line)
	}
	return nil
}

// expandMacro substitutes macro parameters with the invocation arguments.
func expandMacro(def *macroDef, argText, file string, line int) ([]string, error) {
	var args []string
	if strings.TrimSpace(argText) != "" {
		for _, a := range strings.Split(argText, ",") {
			args = append(args, strings.TrimSpace(a))
		}
	}
	if len(args) != len(def.Params) {
		return nil, fmt.Errorf("%s:%d: macro '%s' expects %d argument(s), got %d",
			file, line, def.Name, len(def.Params), len(args))
	}

	expanded := make([]string, len(def.Body))
	copy(expanded, def.Body)
	for idx, param := range def.Params {
		re, err := regexp.Compile(`\b` + regexp.QuoteMeta(param) + `\b`)
		if err != nil {
			return nil, fmt.Errorf("macro '%s': invalid parameter name %q", def.Name, param)
		}
		for li := range expanded {
			expanded[li] = re.ReplaceAllString(expanded[li], args[idx])
		}
	}
	return expanded, nil
}
