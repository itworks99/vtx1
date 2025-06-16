package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

	// Check for input file
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

	// TODO: Implement the assembler logic
	fmt.Println("VTX1 Assembler - Implementation in progress")
	fmt.Println("This is a placeholder for the actual assembler logic")

	// For now, we'll just exit successfully
	os.Exit(ExitSuccess)
}

func showUsage() {
	fmt.Printf("VTX1 Assembler v%s\n\n", Version)
	fmt.Println("Usage:")
	fmt.Println("  vtx1asm [options] input.asm")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}
