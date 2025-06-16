package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
	// Read a source file
	source, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	if verbose {
		fmt.Println("Stage 1: Lexical Analysis")
	}

	// Initialize lexer
	lex := lexer.New(string(source))

	// Parse tokens directly from the lexer
	if verbose {
		fmt.Println("Stage 2: Parsing")
	}

	// Parse using the lexer instance
	p := parser.New(lex)
	ast, err := p.Parse()
	if err != nil {
		return fmt.Errorf("parsing failed: %v", err)
	}

	if verbose {
		fmt.Println("Parsing completed successfully")
	}

	if verbose {
		fmt.Println("Stage 3: Code Generation")
	}

	// Determine output format
	var outputFormat codegen.BinaryFormat
	switch format {
	case "binary":
		outputFormat = codegen.BinaryFormatRaw
	case "hex":
		outputFormat = codegen.BinaryFormatHEX
	case "elf":
		outputFormat = codegen.BinaryFormatELF
	default:
		outputFormat = codegen.BinaryFormatRaw
	}

	// Generate machine code
	gen := codegen.New(ast, outputFormat)
	err = gen.Generate()
	if err != nil {
		return fmt.Errorf("code generation failed: %v", err)
	}

	if verbose {
		fmt.Printf("Code generation completed: %d bytes of machine code\n", len(gen.MachineCode()))
	}

	// Write output based on format
	if err := writeOutput(gen.MachineCode(), outputFile, format); err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}

	// Generate listing file if requested
	if listingFile != "" {
		// We'll pass an empty tokens slice, our generateListing function handles this case
		var tokens []lexer.Token
		if err := generateListing(source, tokens, ast, gen.MachineCode(), listingFile); err != nil {
			return fmt.Errorf("failed to generate listing: %v", err)
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
	listing := fmt.Sprintf("VTX1 Assembler Listing\n\n")

	// If tokens weren't pre-collected, we can tokenize here for the listing
	if len(tokens) == 0 {
		// Initialize a new lexer for tokenization
		tokenLexer := lexer.New(string(source))

		// Collect tokens for the listing
		var collectedTokens []lexer.Token
		for {
			token := tokenLexer.NextToken()
			collectedTokens = append(collectedTokens, token)
			if token.Type == lexer.EOF {
				break
			}
		}
		tokens = collectedTokens
	}

	// TODO: Implement proper listing format with address, machine code, and source code alignment
	// Using source, tokens, ast, and binary data to create a well-formatted listing

	return ioutil.WriteFile(listingFile, []byte(listing), 0644)
}
