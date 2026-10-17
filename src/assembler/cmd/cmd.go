package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/antlr4-go/antlr/v4"
	parser "github.com/kvany/vtx1/assembler/grammar"
)

const (
	ExitSuccess = 0
	ExitError   = 1
	Version     = "0.1.0"
)

// ShowUsage prints usage information for the assembler
func ShowUsage() {
	fmt.Printf("VTX1 Assembler v%s\n\n", Version)
	fmt.Println("Usage:")
	fmt.Println("  vtx1asm [options] input.asm")
	fmt.Println("\nOptions:")
	fmt.Println("  --wordsize=36|108   Output unit: 18-trit words as 36 bits (default) or 108-bit VLIW bundles")
	// The actual flag.PrintDefaults() should be called from main
}

// RunAssembler is the main entry point for assembling a file
func RunAssembler(inputFile, outputFile, listingFile, format string, verbose bool, errorsFile string, wordSize string) error {
	// Default output file is input file with .bin extension
	if outputFile == "" {
		baseName := filepath.Base(inputFile)
		ext := filepath.Ext(baseName)
		nameWithoutExt := baseName[:len(baseName)-len(ext)]
		outputFile = nameWithoutExt + ".bin"
	}
	return assembleFile(inputFile, outputFile, listingFile, format, verbose, errorsFile, wordSize)
}

// CompilationContext holds state and outputs from each compilation stage
type CompilationContext struct {
	SourceFile   string // Input file path
	SourceCode   string // Preprocessed source code content
	Verbose      bool   // Verbose output enabled
	OutputFormat string // Output format

	// Error handling
	ErrorManager *ErrorManager     // Centralized error management system
	SourceMap    map[string]string // Maps filenames to source content for error reporting

	// ANTLR-generated lexer and parser
	Tree antlr.ParseTree // Parse tree from ANTLR

	// Custom AST
	AST *AST // Abstract Syntax Tree

	// Symbol Table
	SymbolTable *SymbolTable

	// Code generation outputs
	MachineCode  []byte            // Generated machine code (flat word stream)
	MachineWords []EmittedWord     // Every word with its address (preserves .ORG layout)
	Symbols      map[string]uint32 // Symbol table for debugging
}

// Minimal stub for ErrorManager
// Replace with real error handling as needed
type ErrorManager struct {
	Errors   []error // List of errors encountered
	Warnings []error // List of warnings encountered
}

func NewErrorManager() *ErrorManager {
	return &ErrorManager{Errors: []error{}, Warnings: []error{}}
}

func (em *ErrorManager) Error() string {
	if len(em.Errors) > 0 {
		return em.Errors[0].Error()
	}
	return ""
}

func (em *ErrorManager) Warning() string {
	if len(em.Warnings) > 0 {
		return em.Warnings[0].Error()
	}
	return ""
}

func (em *ErrorManager) Summary() string {
	return fmt.Sprintf("%d error(s), %d warning(s)", len(em.Errors), len(em.Warnings))
}

func (em *ErrorManager) HasErrors() bool {
	return len(em.Errors) > 0
}

func (em *ErrorManager) HasWarnings() bool {
	return len(em.Warnings) > 0
}

// PrintErrorWithSource prints an error message with the offending source line and a caret under the error column, if possible.
func PrintErrorWithSource(err error, ctx *CompilationContext) {
	errStr := err.Error()
	// Try to extract line and column from error string (format: ... at line X:Y ...)
	var line, col int
	found := false
	_, _ = fmt.Sscanf(errStr, "%*[^l]line %d:%d", &line, &col)
	if line > 0 {
		found = true
	}
	if found && ctx != nil && ctx.SourceCode != "" {
		lines := splitLines(ctx.SourceCode)
		if line-1 >= 0 && line-1 < len(lines) {
			fmt.Fprintf(os.Stderr, "%s\n", errStr)
			fmt.Fprintf(os.Stderr, "%4d | %s\n", line, lines[line-1])
			caret := make([]rune, col+6) // 4 for number, 2 for ' | '
			for i := 0; i < col+5; i++ {
				caret[i] = ' '
			}
			caret[col+5] = '^'
			fmt.Fprintf(os.Stderr, "%s\n", string(caret))
			return
		}
	}
	// Fallback: just print the error
	fmt.Fprintf(os.Stderr, "%s\n", errStr)
}

// splitLines splits a string into lines (handles \r\n, \n, \r)
func splitLines(s string) []string {
	lines := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		} else if s[i] == '\r' {
			if i+1 < len(s) && s[i+1] == '\n' {
				lines = append(lines, s[start:i])
				start = i + 2
				i++
			} else {
				lines = append(lines, s[start:i])
				start = i + 1
			}
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// assembleFile processes the input file and generates the output binary
func assembleFile(inputFile, outputFile, listingFile, format string, verbose bool, errorsFile string, wordSize string) error {
	// Default output format to binary if empty
	if format == "" {
		format = "binary"
	}

	// Create the error manager
	errorManager := NewErrorManager()

	// Stage 0: Preprocessing (.INCLUDE resolution and macro expansion).
	// Note: line numbers in later diagnostics refer to the preprocessed source.
	if verbose {
		fmt.Println("Stage 0: Preprocessing")
	}
	source, sourceMap, err := Preprocess(inputFile)
	if err != nil {
		return fmt.Errorf("preprocessing failed: %v", err)
	}

	// Create a compilation context to pass information between stages
	ctx := &CompilationContext{
		SourceFile:   inputFile,
		SourceCode:   source,
		Verbose:      verbose,
		OutputFormat: format,
		ErrorManager: errorManager,
		SourceMap:    sourceMap,
		SymbolTable:  NewSymbolTable(),
	}
	ctx.SymbolTable.AttachWarnings(&errorManager.Warnings)

	// Attach codegen warnings to ErrorManager before code generation
	AttachCodegenWarnings(&errorManager.Warnings)

	// Stage 1: Lexical Analysis
	if verbose {
		fmt.Println("Stage 1: Lexical Analysis")
	}
	if err := runLexicalAnalysis(ctx); err != nil {
		PrintErrorWithSource(err, ctx)
		return fmt.Errorf("lexical analysis failed: %v", err)
	}

	// Stage 2: Parsing
	if verbose {
		fmt.Println("Stage 2: Parsing")
	}
	if err := runParsing(ctx); err != nil {
		PrintErrorWithSource(err, ctx)
		return fmt.Errorf("parsing failed: %v", err)
	}

	// Stage 2.5: Symbol Table Population (First Pass)
	if verbose {
		fmt.Println("Stage 2.5: Symbol Table Population")
	}
	if err := runSymbolPass(ctx); err != nil {
		PrintErrorWithSource(err, ctx)
		return fmt.Errorf("symbol resolution failed: %v", err)
	}

	// After all passes, check for unused labels and add warnings
	for _, sym := range ctx.SymbolTable.UnusedLabels() {
		msg := fmt.Sprintf("warning: label '%s' defined at %s:%d:%d is never used", sym.Name, sym.File, sym.Line, sym.Column)
		errorManager.Warnings = append(errorManager.Warnings, errors.New(msg))
	}

	// Stage 3: Code Generation
	if verbose {
		fmt.Println("Stage 3: Code Generation")
	}
	if err := runCodeGeneration(ctx); err != nil {
		PrintErrorWithSource(err, ctx)
		return fmt.Errorf("code generation failed: %v", err)
	}

	// Print warnings with source lines if any
	if errorManager.HasWarnings() {
		for _, warn := range errorManager.Warnings {
			PrintErrorWithSource(warn, ctx)
		}
	}

	// Write output based on format
	if err := writeOutputCtx(ctx, outputFile, format, wordSize); err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}
	if verbose {
		fmt.Printf("Output written to %s (%d bytes of machine code, wordsize=%s, format=%s)\n",
			outputFile, len(ctx.MachineCode), wordSize, format)
	}

	// Generate a listing file if requested
	if listingFile != "" {
		if err := generateListing([]byte(source), ctx.Tree, ctx.MachineCode, listingFile); err != nil {
			return fmt.Errorf("failed to generate listing: %v", err)
		}

		if verbose {
			fmt.Printf("Assembly listing written to %s\n", listingFile)
		}
	}

	return nil
}

// runLexicalAnalysis performs lexical analysis on the source code
func runLexicalAnalysis(ctx *CompilationContext) error {
	// ANTLR handles lexing internally; nothing to do here
	return nil
}

// runParsing parses the source code and builds the parse tree
func runParsing(ctx *CompilationContext) error {
	input := antlr.NewInputStream(ctx.SourceCode)
	lexer := parser.Newvtx1_grammarLexer(input)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.Newvtx1_grammarParser(tokens)

	// Remove default error listener and add our custom one
	p.RemoveErrorListeners()
	p.AddErrorListener(NewCustomErrorListener(ctx.ErrorManager))

	ctx.Tree = p.Program()

	if ctx.ErrorManager.HasErrors() {
		for _, err := range ctx.ErrorManager.Errors {
			PrintErrorWithSource(err, ctx)
		}
		return errors.New("syntax errors found during parsing")
	}

	// Build custom AST from parse tree
	ctx.AST = BuildAST(ctx.Tree)
	if ctx.AST == nil {
		err := errors.New("parsing failed: could not build AST due to syntax errors")
		ctx.ErrorManager.Errors = append(ctx.ErrorManager.Errors, err)
		return err
	}
	return nil
}

// runSymbolPass performs the first pass of assembly: populating the symbol table.
func runSymbolPass(ctx *CompilationContext) error {
	var currentAddress uint32 = 0 // Start at address 0, can be changed by .ORG

	for _, line := range ctx.AST.Program.Lines {
		// If there's a label, define it in the symbol table with the current address.
		if line.Label != nil {
			_, err := ctx.SymbolTable.Define(line.Label.Name, currentAddress, ctx.SourceFile, line.Label.Line, line.Label.Column)
			if err != nil {
				// Add the detailed error to the ErrorManager
				ctx.ErrorManager.Errors = append(ctx.ErrorManager.Errors, err)
			}
		}

		// Advance the word address based on the statement type.
		if line.Statement != nil {
			if s, ok := line.Statement.(*DirectiveNode); ok && s.Name == ".ORG" {
				if len(s.Params) > 0 {
					if val, err := parseImmediateOperand(s.Params[0]); err == nil {
						currentAddress = uint32(val)
					}
				}
			} else {
				currentAddress += statementWords(line.Statement)
			}
		}
	}

	// After the pass, check for any undefined symbols that were referenced.
	// Note: References aren't tracked yet, this is just a check for definition.
	undefinedSymbols := ctx.SymbolTable.AllUndefined()
	if len(undefinedSymbols) > 0 {
		for _, s := range undefinedSymbols {
			// This error should be more specific, with line numbers of references.
			err := fmt.Errorf("undefined symbol: %s", s.Name)
			ctx.ErrorManager.Errors = append(ctx.ErrorManager.Errors, err)
		}
	}

	if ctx.ErrorManager.HasErrors() {
		return errors.New("errors found during symbol pass")
	}

	return nil
}

// runCodeGeneration generates machine code from the AST
func runCodeGeneration(ctx *CompilationContext) error {
	if ctx.AST == nil {
		return errors.New("cannot run code generation without a valid AST")
	}

	cg := NewCodeGenerator(ctx.SymbolTable)
	if err := cg.Generate(ctx.AST); err != nil {
		ctx.ErrorManager.Errors = append(ctx.ErrorManager.Errors, err)
		return fmt.Errorf("code generation failed: %v", err)
	}
	ctx.MachineCode = cg.Output
	ctx.MachineWords = cg.Words

	if ctx.Verbose {
		fmt.Printf("Generated %d bytes of machine code.\n", len(ctx.MachineCode))
	}

	return nil
}

// TODO: Add CLI flag for symbol/debug output (e.g., --symbols)
// Stub for writing symbol table/debug info
func writeSymbols(symbols map[string]uint32, outputFile string) error {
	// TODO: Implement symbol table/debug info output
	return errors.New("symbol table/debug info output not implemented yet")
}

// generateListing creates an assembly listing file
func generateListing(source []byte, tree antlr.ParseTree, binary []byte, listingFile string) error {
	// Dummy implementation
	return ioutil.WriteFile(listingFile, []byte("Listing not implemented\n"), 0644)
}
