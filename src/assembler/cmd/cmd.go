package cmd

import (
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
	// The actual flag.PrintDefaults() should be called from main
}

// RunAssembler is the main entry point for assembling a file
func RunAssembler(inputFile, outputFile, listingFile, format string, verbose bool) error {
	// Default output file is input file with .bin extension
	if outputFile == "" {
		baseName := filepath.Base(inputFile)
		ext := filepath.Ext(baseName)
		nameWithoutExt := baseName[:len(baseName)-len(ext)]
		outputFile = nameWithoutExt + ".bin"
	}
	return assembleFile(inputFile, outputFile, listingFile, format, verbose)
}

// CompilationContext holds state and outputs from each compilation stage
type CompilationContext struct {
	SourceFile   string // Input file path
	SourceCode   string // Source code content
	Verbose      bool   // Verbose output enabled
	OutputFormat string // Output format

	// Error handling
	ErrorManager *ErrorManager     // Centralized error management system
	SourceMap    map[string]string // Maps filenames to source content for error reporting

	// ANTLR-generated lexer and parser
	Tree antlr.ParseTree // Parse tree from ANTLR

	// Code generation outputs
	MachineCode []byte            // Generated machine code
	Symbols     map[string]uint32 // Symbol table for debugging
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

// assembleFile processes the input file and generates the output binary
func assembleFile(inputFile, outputFile, listingFile, format string, verbose bool) error {
	// Read the source file
	source, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	// Create the error manager
	errorManager := NewErrorManager()

	// Initialize source map for enhanced error reporting
	sourceMap := make(map[string]string)
	sourceMap[inputFile] = string(source)

	// Create a compilation context to pass information between stages
	ctx := &CompilationContext{
		SourceFile:   inputFile,
		SourceCode:   string(source),
		Verbose:      verbose,
		OutputFormat: format,
		ErrorManager: errorManager,
		SourceMap:    sourceMap,
	}

	// Stage 1: Lexical Analysis
	if verbose {
		fmt.Println("Stage 1: Lexical Analysis")
	}

	if err := runLexicalAnalysis(ctx); err != nil {
		// Print detailed error information with source context
		fmt.Fprintln(os.Stderr, errorManager.Summary())
		return fmt.Errorf("lexical analysis failed: %v", err)
	}

	// Stage 2: Parsing
	if verbose {
		fmt.Println("Stage 2: Parsing")
	}

	if err := runParsing(ctx); err != nil {
		// Print detailed error information with source context
		fmt.Fprintln(os.Stderr, errorManager.Summary())
		return fmt.Errorf("parsing failed: %v", err)
	}

	// Check for warnings even if no errors
	if errorManager.HasWarnings() && verbose {
		fmt.Fprintln(os.Stderr, "Warnings:")
		fmt.Fprintln(os.Stderr, errorManager.Summary())
	}

	// Stage 3: Code Generation
	if verbose {
		fmt.Println("Stage 3: Code Generation")
	}

	if err := runCodeGeneration(ctx); err != nil {
		// Print detailed error information with source context
		fmt.Fprintln(os.Stderr, errorManager.Summary())
		return fmt.Errorf("code generation failed: %v", err)
	}

	// Write output based on format
	if err := writeOutput(ctx.MachineCode, outputFile, format); err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}

	// Generate a listing file if requested
	if listingFile != "" {
		if err := generateListing(source, ctx.Tree, ctx.MachineCode, listingFile); err != nil {
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
	ctx.Tree = p.Program()
	return nil
}

// runCodeGeneration generates machine code from the AST
func runCodeGeneration(ctx *CompilationContext) error {
	if ctx.Tree == nil {
		return fmt.Errorf("cannot run code generation without a valid AST")
	}

	// For now, just print the parse tree since codegen expects the old AST format
	if ctx.Verbose {
		fmt.Println("Parse tree:")
		fmt.Println(ctx.Tree.ToStringTree(nil, nil))
	}

	// TODO: Convert ANTLR parse tree to the expected AST format for codegen
	// For now, create empty machine code
	ctx.MachineCode = []byte{}

	if ctx.Verbose {
		fmt.Printf("Generated %d bytes of machine code.\n", len(ctx.MachineCode))
	}

	return nil
}

// writeOutput writes the generated binary to the specified file
func writeOutput(binary []byte, outputFile, format string) error {
	switch format {
	case "binary":
		return ioutil.WriteFile(outputFile, binary, 0644)

	case "hex":
		// Format as Intel HEX file
		hexData := formatAsHex(binary)
		return ioutil.WriteFile(outputFile, []byte(hexData), 0644)

	case "objdump":
		// Format as objdump-like output
		objDump := formatAsObjDump(binary)
		return ioutil.WriteFile(outputFile, []byte(objDump), 0644)

	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// formatAsHex formats binary data as Intel HEX format
func formatAsHex(data []byte) string {
	// Dummy implementation
	return fmt.Sprintf(":%X\n", data)
}

// formatAsObjDump formats binary data as a simple objdump-like output
func formatAsObjDump(data []byte) string {
	// Dummy implementation
	return fmt.Sprintf("OBJ DUMP: %X\n", data)
}

// generateListing creates an assembly listing file
func generateListing(source []byte, tree antlr.ParseTree, binary []byte, listingFile string) error {
	// Dummy implementation
	return ioutil.WriteFile(listingFile, []byte("Listing not implemented\n"), 0644)
}
