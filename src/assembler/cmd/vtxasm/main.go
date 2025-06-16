package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kvany/vtx1/assembler/internal/codegen"
	"github.com/kvany/vtx1/assembler/internal/lexer"
	"github.com/kvany/vtx1/assembler/internal/parser"
)

const (
	ExitSuccess = 0
	ExitError   = 1
	Version     = "0.1.0"
)

// Command line flags
var (
	outputFile  string
	listingFile string
	verbose     bool
	format      string
	showVersion bool
	showHelp    bool
)

func init() {
	// Set up command line flags
	flag.StringVar(&outputFile, "o", "", "Output binary file (default: input.bin)")
	flag.StringVar(&outputFile, "output", "", "Output binary file (default: input.bin)")

	flag.StringVar(&listingFile, "l", "", "Generate assembly listing file")
	flag.StringVar(&listingFile, "listing", "", "Generate assembly listing file")

	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")

	flag.StringVar(&format, "f", "binary", "Output format: binary, hex, or objdump")
	flag.StringVar(&format, "format", "binary", "Output format: binary, hex, or objdump")

	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showHelp, "h", false, "Show help information")
	flag.BoolVar(&showHelp, "help", false, "Show help information")
}

func main() {
	// Parse command line flags
	flag.Parse()

	// Show version and exit
	if showVersion {
		fmt.Printf("VTX1 Assembler v%s\n", Version)
		os.Exit(ExitSuccess)
	}

	// Show help and exit
	if showHelp {
		showUsage()
		os.Exit(ExitSuccess)
	}

	// Check for an input file
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Missing input file\n")
		showUsage()
		os.Exit(ExitError)
	}

	inputFile := args[0]

	// Default output file is input file with .bin extension
	if outputFile == "" {
		baseName := filepath.Base(inputFile)
		ext := filepath.Ext(baseName)
		nameWithoutExt := baseName[:len(baseName)-len(ext)]
		outputFile = nameWithoutExt + ".bin"
	}

	if verbose {
		fmt.Printf("VTX1 Assembler v%s\n", Version)
		fmt.Printf("Input file: %s\n", inputFile)
		fmt.Printf("Output file: %s\n", outputFile)
		if listingFile != "" {
			fmt.Printf("Listing file: %s\n", listingFile)
		}
		fmt.Printf("Format: %s\n", format)
	}

	// Implement the assembler logic
	err := assembleFile(inputFile, outputFile, listingFile, format, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ExitError)
	}

	if verbose {
		fmt.Println("Assembly completed successfully.")
	}

	os.Exit(ExitSuccess)
}

func showUsage() {
	fmt.Printf("VTX1 Assembler v%s\n\n", Version)
	fmt.Println("Usage:")
	fmt.Println("  vtx1asm [options] input.asm")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}

// assembleFile processes the input file and generates the output binary
func assembleFile(inputFile, outputFile, listingFile, format string, verbose bool) error {
	// Read the source file
	source, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	// Create a compilation context to pass information between stages
	ctx := &CompilationContext{
		SourceFile:   inputFile,
		SourceCode:   string(source),
		Verbose:      verbose,
		OutputFormat: format,
	}

	// Stage 1: Lexical Analysis
	if verbose {
		fmt.Println("Stage 1: Lexical Analysis")
	}

	if err := runLexicalAnalysis(ctx); err != nil {
		return fmt.Errorf("lexical analysis failed: %v", err)
	}

	// Stage 2: Parsing
	if verbose {
		fmt.Println("Stage 2: Parsing")
	}

	if err := runParsing(ctx); err != nil {
		return fmt.Errorf("parsing failed: %v", err)
	}

	// Stage 3: Code Generation
	if verbose {
		fmt.Println("Stage 3: Code Generation")
	}

	if err := runCodeGeneration(ctx); err != nil {
		return fmt.Errorf("code generation failed: %v", err)
	}

	// Write output based on format
	if err := writeOutput(ctx.MachineCode, outputFile, format); err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}

	// Generate a listing file if requested
	if listingFile != "" {
		if err := generateListing([]byte(ctx.SourceCode), ctx.Tokens, ctx.AST, ctx.MachineCode, listingFile); err != nil {
			return fmt.Errorf("failed to generate listing: %v", err)
		}

		if verbose {
			fmt.Printf("Assembly listing written to %s\n", listingFile)
		}
	}

	if verbose {
		fmt.Printf("Successfully assembled %s (%d bytes of machine code)\n",
			filepath.Base(inputFile), len(ctx.MachineCode))
	}

	return nil
}

// CompilationContext holds state and outputs from each compilation stage
type CompilationContext struct {
	SourceFile   string // Input file path
	SourceCode   string // Source code content
	Verbose      bool   // Verbose output enabled
	OutputFormat string // Output format

	// Lexical analysis outputs
	Lexer  *lexer.Lexer  // Lexer instance
	Tokens []lexer.Token // Collected tokens (for diagnostics)

	// Parsing outputs
	Parser *parser.Parser // Parser instance
	AST    *parser.AST    // Abstract Syntax Tree

	// Code generation outputs
	CodeGen     *codegen.CodeGenerator // Code generator instance (if available)
	MachineCode []byte                 // Generated machine code
	Symbols     map[string]int         // Symbol table (for debugging)
	Diagnostics []string               // Warnings and information
}

// runLexicalAnalysis performs lexical analysis on the source code
func runLexicalAnalysis(ctx *CompilationContext) error {
	// Initialize lexer
	ctx.Lexer = lexer.New(ctx.SourceCode)

	// If verbose, collect tokens for diagnostics
	if ctx.Verbose {
		// Collect tokens for diagnostics while not disturbing the lexer state
		tempLexer := lexer.New(ctx.SourceCode)
		ctx.Tokens = []lexer.Token{}

		for {
			token := tempLexer.NextToken()
			ctx.Tokens = append(ctx.Tokens, token)

			if token.Type == lexer.EOF {
				break
			}
		}

		fmt.Printf("Tokenization completed: %d tokens generated\n", len(ctx.Tokens))

		// Print sample tokens if very verbose
		if len(os.Getenv("VTX1_VERY_VERBOSE")) > 0 {
			fmt.Println("First 10 tokens:")
			for i, t := range ctx.Tokens {
				if i >= 10 {
					break
				}
				fmt.Printf("  %s\n", t.String())
			}
		}
	}

	return nil
}

// runParsing parses the tokens into an AST
func runParsing(ctx *CompilationContext) error {
	// Initialize parser with the lexer
	ctx.Parser = parser.New(ctx.Lexer)

	// Parse the tokens into an AST
	ast, err := ctx.Parser.Parse()
	if err != nil {
		return err
	}

	ctx.AST = ast

	if ctx.Verbose {
		fmt.Println("Parsing completed successfully")

		// Print AST node count if very verbose
		if len(os.Getenv("VTX1_VERY_VERBOSE")) > 0 {
			nodeCount := countASTNodes(ast)
			fmt.Printf("AST contains %d nodes\n", nodeCount)
		}
	}

	return nil
}

// countASTNodes counts the number of nodes in the AST (helper for diagnostics)
func countASTNodes(node *parser.AST) int {
	if node == nil {
		return 0
	}

	count := 1 // Count this node

	for _, child := range node.Children {
		count += countASTNodes(child)
	}

	return count
}

// runCodeGeneration generates machine code from the AST
func runCodeGeneration(ctx *CompilationContext) error {
	// Determine output format
	var outputFormat codegen.BinaryFormat
	switch ctx.OutputFormat {
	case "binary":
		outputFormat = codegen.BinaryFormatRaw
	case "hex":
		outputFormat = codegen.BinaryFormatHEX
	case "elf":
		outputFormat = codegen.BinaryFormatELF
	default:
		outputFormat = codegen.BinaryFormatRaw
	}

	// Create code generator
	generator := codegen.New(ctx.AST, outputFormat)
	ctx.CodeGen = generator

	// Generate machine code
	err := generator.Generate()
	if err != nil {
		return err
	}

	// Get the generated machine code
	ctx.MachineCode = generator.MachineCode()

	// Extract symbol table if available
	if table := generator.SymbolTable(); table != nil {
		ctx.Symbols = table
	}

	if ctx.Verbose {
		fmt.Printf("Code generation completed: %d bytes of machine code\n", len(ctx.MachineCode))

		// Show an address map of main sections if very verbose
		if len(os.Getenv("VTX1_VERY_VERBOSE")) > 0 && len(ctx.Symbols) > 0 {
			fmt.Println("Symbol addresses:")
			for sym, addr := range ctx.Symbols {
				fmt.Printf("  %-20s 0x%06X\n", sym, addr)
			}
		}
	}

	return nil
}

// writeOutput writes the binary data to the specified output file in the requested format
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
	var result string
	addr := 0

	for i := 0; i < len(data); i += 16 {
		end := i + 16
		if end > len(data) {
			end = len(data)
		}

		chunk := data[i:end]
		length := len(chunk)

		// Calculate checksum (2's complement of sum of all bytes in record)
		checksum := byte(length)
		checksum += byte(addr >> 8)
		checksum += byte(addr & 0xFF)

		for _, b := range chunk {
			checksum += b
		}
		checksum = byte(-(int(checksum)))

		// Format record
		record := fmt.Sprintf(":%02X%04X00", length, addr)

		for _, b := range chunk {
			record += fmt.Sprintf("%02X", b)
		}

		record += fmt.Sprintf("%02X\n", checksum)
		result += record

		addr += length
	}

	// Add EOF record
	result += ":00000001FF\n"

	return result
}

// formatAsObjDump formats binary data similar to objdump output
func formatAsObjDump(data []byte) string {
	var result string
	result += fmt.Sprintf("VTX1 binary file - %d bytes\n\n", len(data))
	result += "Contents of binary file:\n\n"
	result += "Addr    Bytes                                   ASCII\n"
	result += "------------------------------------------------------------\n"

	for i := 0; i < len(data); i += 16 {
		// Address
		result += fmt.Sprintf("%06x  ", i)

		end := i + 16
		if end > len(data) {
			end = len(data)
		}

		// Hex bytes
		for j := i; j < i+16; j++ {
			if j < end {
				result += fmt.Sprintf("%02x ", data[j])
			} else {
				result += "   "
			}

			// Extra space after 8 bytes
			if j == i+7 {
				result += " "
			}
		}

		// ASCII representation
		result += " |"
		for j := i; j < end; j++ {
			b := data[j]
			if b >= 32 && b <= 126 {
				result += string(b)
			} else {
				result += "."
			}
		}
		result += "|\n"
	}

	return result
}

// generateListing creates a listing file with source, tokens, and binary representation
func generateListing(source []byte, tokens []lexer.Token, ast *parser.AST, binary []byte, listingFile string) error {
	// Create a map to track which machine code addresses correspond to which source lines
	addressToSourceLineMap := buildAddressSourceMap(ast, binary)

	// Split source code into lines for display
	sourceLines := strings.Split(string(source), "\n")

	// Create the header for the listing
	listing := strings.Builder{}
	listing.WriteString(fmt.Sprintf("VTX1 Assembler Listing - Generated %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	listing.WriteString(fmt.Sprintf("%-10s %-30s %s\n", "Address", "Machine Code", "Source Code"))
	listing.WriteString(fmt.Sprintf("%-10s %-30s %s\n", "-------", "---------------", "--------------------------------"))

	// Process by source line
	for lineNum, sourceLine := range sourceLines {
		// Ensure source line is not excessively long
		displaySource := sourceLine
		if len(displaySource) > 60 {
			displaySource = displaySource[:57] + "..."
		}

		// Check if this source line generated any machine code
		if addresses, found := getAddressesForSourceLine(lineNum+1, addressToSourceLineMap); found {
			// This line generated machine code - display with address and bytes
			for i, addr := range addresses {
				// Display the memory address
				addrStr := fmt.Sprintf("0x%06X", addr)

				// Get machine code bytes (typically 4 or 8 bytes per instruction for VTX1)
				mcodeBytes := getInstructionBytes(binary, addr)

				// Format machine code as hex values
				mcodeHex := formatBytesAsHex(mcodeBytes)

				if i == 0 {
					// First or only instruction for this source line - show with source
					listing.WriteString(fmt.Sprintf("%-10s %-30s %s\n", addrStr, mcodeHex, displaySource))
				} else {
					// Additional instructions for this source line - just show the code
					listing.WriteString(fmt.Sprintf("%-10s %-30s\n", addrStr, mcodeHex))
				}
			}
		} else {
			// This line generated no machine code (comment, directive, etc.)
			listing.WriteString(fmt.Sprintf("%-10s %-30s %s\n", "", "", displaySource))
		}
	}

	// Add symbol table if we have access to symbols
	if ast != nil {
		// Try to get symbols from the code generator
		// Note: In a full implementation, these would be passed from the compiler context
		var symbols map[string]struct {
			Address uint32
			Defined bool
		}

		// Create a simple symbol table for demonstration
		symbols = make(map[string]struct {
			Address uint32
			Defined bool
		})

		// Extract some basic symbols from the AST - in a real implementation,
		// this would use the actual symbol table from the compiler
		extractSymbolsFromAST(ast, symbols)

		// Only display the section if we found symbols
		if len(symbols) > 0 {
			listing.WriteString("\n\nSymbol Table:\n")
			listing.WriteString(fmt.Sprintf("%-20s %-10s %s\n", "Name", "Address", "Defined"))
			listing.WriteString(fmt.Sprintf("%-20s %-10s %s\n", "--------------------", "----------", "-------"))

			// Sort symbols by address
			var sortedSymbols []struct {
				Name    string
				Address uint32
				Defined bool
			}

			for name, info := range symbols {
				sortedSymbols = append(sortedSymbols, struct {
					Name    string
					Address uint32
					Defined bool
				}{
					Name:    name,
					Address: info.Address,
					Defined: info.Defined,
				})
			}

			sort.Slice(sortedSymbols, func(i, j int) bool {
				return sortedSymbols[i].Address < sortedSymbols[j].Address
			})

			// Print sorted symbols
			for _, sym := range sortedSymbols {
				listing.WriteString(fmt.Sprintf("%-20s 0x%08X %v\n", sym.Name, sym.Address, sym.Defined))
			}
		}
	}

	return ioutil.WriteFile(listingFile, []byte(listing.String()), 0644)
}

// extractSymbolsFromAST extracts symbol information from the AST
// This is a simplified implementation for demonstration purposes
func extractSymbolsFromAST(ast *parser.AST, symbols map[string]struct {
	Address uint32
	Defined bool
}) {
	if ast == nil {
		return
	}

	// Check if this node is a label
	if ast.Type == parser.NODE_LABEL && ast.Value != nil {
		if labelName, ok := ast.Value.(string); ok {
			// In a real implementation, you'd get the actual address from the symbol table
			// Here we're just assigning sequential addresses for demonstration
			addr := uint32(len(symbols) * 4) // Just for demonstration

			symbols[labelName] = struct {
				Address uint32
				Defined bool
			}{
				Address: addr,
				Defined: true,
			}
		}
	}

	// Process child nodes
	for _, child := range ast.Children {
		extractSymbolsFromAST(child, symbols)
	}
}

// buildAddressSourceMap creates a mapping from machine code addresses to source line numbers
func buildAddressSourceMap(ast *parser.AST, binary []byte) map[uint32]int {
	addressMap := make(map[uint32]int)

	// Start at address 0 and process the AST
	processASTForAddressMap(ast, addressMap, 0)

	return addressMap
}

// processASTForAddressMap recursively processes AST nodes to build address-to-line mapping
func processASTForAddressMap(node *parser.AST, addressMap map[uint32]int, currentAddr uint32) uint32 {
	if node == nil {
		return currentAddr
	}

	// If this is an instruction node, map its address to its line number
	if node.Type == parser.NODE_INSTRUCTION {
		addressMap[currentAddr] = node.Line

		// For VTX1, assume each instruction is 4 bytes (can be customized based on actual architecture)
		currentAddr += 4
	}

	// Process all child nodes
	for _, child := range node.Children {
		currentAddr = processASTForAddressMap(child, addressMap, currentAddr)
	}

	return currentAddr
}

// getAddressesForSourceLine finds all instruction addresses for a given source line
func getAddressesForSourceLine(lineNum int, addressMap map[uint32]int) ([]uint32, bool) {
	var addresses []uint32

	for addr, line := range addressMap {
		if line == lineNum {
			addresses = append(addresses, addr)
		}
	}

	// Sort addresses for consistent output
	if len(addresses) > 0 {
		sort.Slice(addresses, func(i, j int) bool {
			return addresses[i] < addresses[j]
		})
		return addresses, true
	}

	return nil, false
}

// getInstructionBytes extracts the machine code bytes for an instruction at the given address
func getInstructionBytes(binary []byte, addr uint32) []byte {
	if int(addr) >= len(binary) {
		return []byte{}
	}

	// Define instruction size (4 bytes for VTX1)
	instructionSize := 4
	if int(addr)+instructionSize > len(binary) {
		instructionSize = len(binary) - int(addr)
	}

	return binary[addr : addr+uint32(instructionSize)]
}

// formatBytesAsHex formats a byte slice as a hex string with spaces between bytes
func formatBytesAsHex(bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, b := range bytes {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprintf("%02X", b))
	}

	return sb.String()
}
