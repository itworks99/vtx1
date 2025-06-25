package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kvany/vtx1/assembler/cmd"
)

func main() {
	// Set up command line flags
	outputFile := flag.String("o", "", "Output binary file (default: input.bin)")
	flag.StringVar(outputFile, "output", "", "Output binary file (default: input.bin)")

	listingFile := flag.String("l", "", "Generate assembly listing file")
	flag.StringVar(listingFile, "listing", "", "Generate assembly listing file")

	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.BoolVar(verbose, "verbose", false, "Enable verbose output")

	format := flag.String("f", "binary", "Output format: binary, hex, or objdump")
	flag.StringVar(format, "format", "binary", "Output format: binary, hex, or objdump")

	showVersion := flag.Bool("version", false, "Show version information")
	showHelp := flag.Bool("h", false, "Show help information")
	flag.BoolVar(showHelp, "help", false, "Show help information")

	flag.Parse()

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

	err := cmd.RunAssembler(inputFile, *outputFile, *listingFile, *format, *verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(cmd.ExitError)
	}

	if *verbose {
		fmt.Println("Assembly completed successfully.")
	}

	os.Exit(cmd.ExitSuccess)
}
