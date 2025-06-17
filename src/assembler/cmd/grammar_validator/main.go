package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kvany/vtx1/assembler/internal/grammar"
	"github.com/kvany/vtx1/assembler/internal/lexer"
	"github.com/kvany/vtx1/assembler/internal/parser"
)

func main() {
	// Get the grammar file path
	grammarFile := filepath.Join("prototype", "grammar", "vtx1_grammar.ebnf")

	// Parse the EBNF grammar
	g, err := grammar.ParseEBNF(grammarFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing grammar: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully parsed EBNF grammar\n")

	// Extract language constructs
	instructions := g.ListInstructions()
	fmt.Printf("Found %d instruction types\n", len(instructions))

	// Generate a simple test case
	testCases := g.GenerateTestCases()
	if len(testCases) > 0 {
		fmt.Printf("Generated example: %s\n", testCases[0])
	}

	// Validate that our lexer/parser implementation matches the grammar
	if len(os.Args) > 1 {
		inputFile := os.Args[1]

		// Read the input file
		source, err := os.ReadFile(inputFile)
		if err != nil {
			fprintf, err := fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}

		// Check if it conforms to our grammar
		isValid, err := g.ValidateSyntax(string(source))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating syntax: %v\n", err)
			os.Exit(1)
		}
		if isValid {
			fmt.Println("Input file conforms to the EBNF grammar")
		} else {
			fmt.Println("Input file has syntax errors according to the EBNF grammar")
		}

		// Process with our implementation
		l := lexer.New(string(source))
		p := parser.New(l)

		ast, err := p.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing with our implementation: %v\n", err)
			for _, parseErr := range p.GetErrors() {
				fmt.Fprintf(os.Stderr, "  %s\n", parseErr.Error())
			}
			os.Exit(1)
		}

		fmt.Println("Successfully parsed input with our implementation")
		fmt.Printf("AST has %d nodes\n", countNodes(ast))
	}
}

func countNodes(ast *parser.AST) int {
	if ast == nil {
		return 0
	}

	count := 1 // Count this node
	for _, child := range ast.Children {
		count += countNodes(child)
	}

	return count
}
