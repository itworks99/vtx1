package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kvany/vtx1/assembler/cmd"
)

func main() {
	// Debug: print raw os.Args and flag.Args()
	fmt.Printf("[DEBUG] os.Args: %v\n", os.Args)

	// Set up command line flags
	outputFile := flag.String("o", "", "Output binary file (default: input.bin)")
	flag.StringVar(outputFile, "output", "", "Output binary file (default: input.bin)")

	listingFile := flag.String("l", "", "Generate assembly listing file")
	flag.StringVar(listingFile, "listing", "", "Generate assembly listing file")

	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.BoolVar(verbose, "verbose", false, "Enable verbose output")

	format := flag.String("f", "binary", "Output format: binary, hex, or objdump")
	flag.StringVar(format, "format", "binary", "Output format: binary, hex, or objdump")

	// Add errors flag
	errorsFile := flag.String("errors", "", "Write errors and warnings to this file")

	showVersion := flag.Bool("version", false, "Show version information")
	showHelp := flag.Bool("h", false, "Show help information")
	flag.BoolVar(showHelp, "help", false, "Show help information")

	wordSize := flag.String("wordsize", "8", "Output word size/format: 8, 36, 108, ternary")
	flag.StringVar(wordSize, "w", "8", "Output word size/format: 8, 36, 108, ternary")

	flag.Parse()
	fmt.Printf("[DEBUG] flag.Args(): %v\n", flag.Args())

	if *showVersion {
		fmt.Printf("VTX1 Assembler v%s\n", cmd.Version)
		os.Exit(cmd.ExitSuccess)
	}

	if *showHelp {
		cmd.ShowUsage()
		os.Exit(cmd.ExitSuccess)
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Missing input file\n")
		cmd.ShowUsage()
		os.Exit(cmd.ExitError)
	}

	inputFile := args[0]

	err := cmd.RunAssembler(inputFile, *outputFile, *listingFile, *format, *verbose, *errorsFile, *wordSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(cmd.ExitError)
	}

	if *verbose {
		fmt.Println("Assembly completed successfully.")
	}

	os.Exit(cmd.ExitSuccess)
}
