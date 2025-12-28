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
	fmt.Println("  --wordsize=8|36|108|ternary   Output word size/format (default: 8-bit bytes)")
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
	SourceCode   string // Source code content
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
	if found && ctx != nil && ctx.SourceMap != nil {
		src, ok := ctx.SourceMap[ctx.SourceFile]
		if ok {
			lines := splitLines(src)
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
	fmt.Println("[DEBUG] Entered assembleFile")
	// Read the source file
	source, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	// Default output format to binary if empty
	if format == "" {
		format = "binary"
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
		SymbolTable:  NewSymbolTable(),
	}
	ctx.SymbolTable.AttachWarnings(&errorManager.Warnings)

	// Attach codegen warnings to ErrorManager before code generation
	AttachCodegenWarnings(&errorManager.Warnings)

	// Stage 1: Lexical Analysis
	if verbose {
		fmt.Println("Stage 1: Lexical Analysis")
	}
	fmt.Println("[DEBUG] Starting lexical analysis...")
	if err := runLexicalAnalysis(ctx); err != nil {
		PrintErrorWithSource(err, ctx)
		fmt.Println("[DEBUG] Lexical analysis failed.")
		return fmt.Errorf("lexical analysis failed: %v", err)
	}
	fmt.Println("[DEBUG] Lexical analysis complete.")

	// Stage 2: Parsing
	if verbose {
		fmt.Println("Stage 2: Parsing")
	}
	fmt.Println("[DEBUG] Starting parsing...")
	if err := runParsing(ctx); err != nil {
		PrintErrorWithSource(err, ctx)
		fmt.Println("[DEBUG] Parsing failed.")
		return fmt.Errorf("parsing failed: %v", err)
	}
	fmt.Println("[DEBUG] Parsing complete.")

	// Stage 2.5: Symbol Table Population (First Pass)
	if verbose {
		fmt.Println("Stage 2.5: Symbol Table Population")
	}
	fmt.Println("[DEBUG] Starting symbol pass...")
	if err := runSymbolPass(ctx); err != nil {
		PrintErrorWithSource(err, ctx)
		fmt.Println("[DEBUG] Symbol pass failed.")
		return fmt.Errorf("symbol resolution failed: %v", err)
	}
	fmt.Println("[DEBUG] Symbol pass complete.")

	// After all passes, check for unused labels and add warnings
	for _, sym := range ctx.SymbolTable.UnusedLabels() {
		msg := fmt.Sprintf("warning: label '%s' defined at %s:%d:%d is never used", sym.Name, sym.File, sym.Line, sym.Column)
		errorManager.Warnings = append(errorManager.Warnings, fmt.Errorf(msg))
	}

	// Print warnings with source lines if any
	if errorManager.HasWarnings() {
		for _, warn := range errorManager.Warnings {
			PrintErrorWithSource(warn, ctx)
		}
	}

	// Stage 3: Code Generation
	if verbose {
		fmt.Println("Stage 3: Code Generation")
	}
	fmt.Println("[DEBUG] Starting code generation...")
	if err := runCodeGeneration(ctx); err != nil {
		PrintErrorWithSource(err, ctx)
		fmt.Println("[DEBUG] Code generation failed.")
		return fmt.Errorf("code generation failed: %v", err)
	}
	fmt.Println("[DEBUG] Code generation complete.")

	// Write output based on format
	fmt.Printf("[DEBUG] MachineCode length before writeOutput: %d bytes\n", len(ctx.MachineCode))
	if err := writeOutput(ctx.MachineCode, outputFile, format, wordSize); err != nil {
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

	// Remove default error listener and add our custom one
	p.RemoveErrorListeners()
	p.AddErrorListener(NewCustomErrorListener(ctx.ErrorManager))

	ctx.Tree = p.Program()

	// Debug: confirm parse tree is built
	if ctx.Tree == nil {
		fmt.Println("[DEBUG] ctx.Tree is nil after parsing!")
	} else {
		fmt.Println("[DEBUG] ctx.Tree is not nil after parsing.")
	}

	// Debug: print any errors collected by the custom error listener
	if ctx.ErrorManager.HasErrors() {
		fmt.Println("[DEBUG] Errors collected by custom error listener:")
		for _, err := range ctx.ErrorManager.Errors {
			fmt.Println("  ", err)
		}
	}

	// Debug: print the parse tree as a string
	fmt.Println("[DEBUG] Parse tree:")
	fmt.Println(ctx.Tree.ToStringTree(nil, p))

	// Build custom AST from parse tree
	ctx.AST = BuildAST(ctx.Tree)
	if ctx.AST == nil {
		err := fmt.Errorf("Parsing failed: could not build AST due to syntax errors.")
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

		// Advance the address based on the statement type.
		if line.Statement != nil {
			switch s := line.Statement.(type) {
			case *InstructionNode:
				// Standard instructions are 1 word (e.g., 4 bytes in a 32-bit model)
				// This is a placeholder size.
				currentAddress += 4
			case *VLIWInstructionNode:
				// VLIW are larger, e.g., 12 bytes. Placeholder size.
				currentAddress += 12
			case *DirectiveNode:
				// Handle directives that affect the address counter.
				switch s.Name {
				case ".ORG":
					// Placeholder: need to parse immediate value
					// val, _ := parseImmediate(s.Params[0])
					// currentAddress = val
				case ".DB", ".DW", ".DT":
					// Placeholder: need to calculate size of data
					// currentAddress += calculateDataSize(s)
				case ".SPACE":
					// Placeholder: need to parse immediate value
					// val, _ := parseImmediate(s.Params[0])
					// currentAddress += val
				}
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
		return fmt.Errorf("errors found during symbol pass")
	}

	return nil
}

// runCodeGeneration generates machine code from the AST
func runCodeGeneration(ctx *CompilationContext) error {
	if ctx.AST == nil {
		return fmt.Errorf("cannot run code generation without a valid AST")
	}

	fmt.Printf("[DEBUG] AST before code generation: %+v\n", ctx.AST)

	cg := NewCodeGenerator(ctx.SymbolTable)
	if err := cg.Generate(ctx.AST); err != nil {
		ctx.ErrorManager.Errors = append(ctx.ErrorManager.Errors, err)
		return fmt.Errorf("code generation failed: %v", err)
	}
	ctx.MachineCode = cg.Output

	if ctx.Verbose {
		fmt.Printf("Generated %d bytes of machine code.\n", len(ctx.MachineCode))
	}

	return nil
}

// writeOutput writes the generated binary to the specified file
func writeOutput(binary []byte, outputFile, format, wordSize string) error {
	switch wordSize {
	case "8":
		fmt.Printf("[DEBUG] writeOutput: writing %d bytes (8-bit)\n", len(binary))
		return writeOutput8(binary, outputFile, format)
	case "36":
		packed := pack36BitWords(binary)
		fmt.Printf("[DEBUG] writeOutput: writing %d bytes (36-bit packed)\n", len(packed))
		return writeOutput8(packed, outputFile, format)
	case "108":
		packed := pack108BitWords(binary)
		fmt.Printf("[DEBUG] writeOutput: writing %d bytes (108-bit packed)\n", len(packed))
		return writeOutput8(packed, outputFile, format)
	case "ternary":
		packed := packTernary(binary)
		fmt.Printf("[DEBUG] writeOutput: writing %d bytes (ternary packed)\n", len(packed))
		return writeOutput8(packed, outputFile, format)
	default:
		return fmt.Errorf("unsupported word size: %s", wordSize)
	}
}

// writeOutput8 is the existing logic for 8-bit output
func writeOutput8(binary []byte, outputFile, format string) error {
	switch format {
	case "binary":
		return ioutil.WriteFile(outputFile, binary, 0644)
	case "hex":
		// TODO: Implement real Intel HEX output. Current implementation is a placeholder.
		hexData := formatAsHex(binary)
		return ioutil.WriteFile(outputFile, []byte(hexData), 0644)
	case "objdump":
		// TODO: Implement real objdump output. Current implementation is a placeholder.
		objDump := formatAsObjDump(binary)
		return ioutil.WriteFile(outputFile, []byte(objDump), 0644)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// pack36BitWords packs every 36 bits into 5 bytes (last 4 bits zeroed if needed)
func pack36BitWords(data []byte) []byte {
	var out []byte
	for i := 0; i+4 < len(data); i += 4 {
		// Take 4 bytes (32 bits)
		b0 := data[i]
		b1 := data[i+1]
		b2 := data[i+2]
		b3 := data[i+3]
		var b4 byte = 0
		if i+4 < len(data) {
			b4 = data[i+4] & 0x0F // Only lower 4 bits for 36th-39th bits
		}
		out = append(out, b0, b1, b2, b3, b4)
	}
	// Handle trailing bytes (if any)
	rest := len(data) % 4
	if rest > 0 {
		start := len(data) - rest
		chunk := make([]byte, 5)
		copy(chunk, data[start:])
		out = append(out, chunk...)
	}
	return out
}

// pack108BitWords packs every 108 bits into 14 bytes
func pack108BitWords(data []byte) []byte {
	var out []byte
	for i := 0; i+12 < len(data); i += 12 {
		// 12 bytes = 96 bits, need 12 more bits (1.5 bytes)
		chunk := make([]byte, 14)
		copy(chunk, data[i:i+12])
		if i+13 < len(data) {
			chunk[12] = data[i+12]
		}
		if i+14 < len(data) {
			chunk[13] = data[i+13] & 0x0F // Only lower 4 bits for last nibble
		}
		out = append(out, chunk...)
	}
	// Handle trailing bytes (if any)
	rest := len(data) % 12
	if rest > 0 {
		start := len(data) - rest
		chunk := make([]byte, 14)
		copy(chunk, data[start:])
		out = append(out, chunk...)
	}
	return out
}

// packTernary packs trits as 2 bits each using VTX1 encoding
// Each 36-bit word = 18 trits, each trit encoded as 2 bits:
// 00 = -1, 01 = 0, 10 = +1, 11 = undefined
func packTernary(data []byte) []byte {
	var out []byte
	for i := 0; i+4 < len(data); i += 4 {
		// Treat 4 bytes as a 32-bit integer, pad with 0 if needed
		var word uint32 = 0
		for j := 0; j < 4 && i+j < len(data); j++ {
			word |= uint32(data[i+j]) << (8 * (3 - j))
		}
		// Convert to 18 trits (balanced ternary)
		trits := intToBalancedTrits(word, 18)
		// Encode trits as 2 bits each
		packed := encodeTritsToBytes(trits)
		out = append(out, packed...)
	}
	return out
}

// intToBalancedTrits converts an integer to a slice of n balanced trits (-1,0,+1)
func intToBalancedTrits(val uint32, n int) []int {
	trits := make([]int, n)
	v := int(val)
	for i := n - 1; i >= 0; i-- {
		rem := v % 3
		v /= 3
		if rem == 2 {
			trits[i] = -1
			v += 1
		} else {
			trits[i] = rem
		}
	}
	return trits
}

// encodeTritsToBytes encodes 18 trits (-1,0,+1) as 36 bits (5 bytes)
func encodeTritsToBytes(trits []int) []byte {
	var bits uint64 = 0
	for i, t := range trits {
		var enc uint64
		switch t {
		case -1:
			enc = 0b00
		case 0:
			enc = 0b01
		case 1:
			enc = 0b10
		default:
			enc = 0b11 // undefined
		}
		bits |= (enc << (2 * (17 - i)))
	}
	out := make([]byte, 5)
	for i := 0; i < 5; i++ {
		out[i] = byte((bits >> (8 * (4 - i))) & 0xFF)
	}
	return out
}

// TODO: Add CLI flag for symbol/debug output (e.g., --symbols)
// Stub for writing symbol table/debug info
func writeSymbols(symbols map[string]uint32, outputFile string) error {
	// TODO: Implement symbol table/debug info output
	return fmt.Errorf("symbol table/debug info output not implemented yet")
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
