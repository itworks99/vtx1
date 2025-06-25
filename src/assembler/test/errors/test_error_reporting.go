// This program demonstrates the enhanced error reporting system
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kvany/vtx1/assembler/internal/errors"
	"github.com/kvany/vtx1/assembler/internal/lexer"
	"github.com/kvany/vtx1/assembler/internal/parser"
)

func main() {
	// Get test file path from command line or use default
	filename := "multiple_errors.asm"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	// Resolve full path
	testFilePath := filepath.Join(".", filename)

	// Read the test file
	source, err := ioutil.ReadFile(testFilePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Create error manager that will collect all errors
	errorManager := errors.NewErrorManager()

	// Set a higher limit to see more errors at once
	errorManager.MaxErrors = 50

	fmt.Printf("Testing error reporting on file: %s\n\n", testFilePath)
	fmt.Println("Starting lexical analysis and parsing...")

	// Create lexer with source filename for better error reporting
	lex := lexer.NewWithFilename(string(source), testFilePath)
	lex.SetErrorManager(errorManager)

	// Create parser with the lexer
	p := parser.New(lex)
	p.SetErrorManager(errorManager)

	// Add source file contents for better error reporting
	p.AddSourceFile(testFilePath, string(source))

	// Parse the input and collect errors
	ast, err := p.Parse()

	// Check if there were any errors
	if errorManager.HasErrors() {
		fmt.Println("\nFound errors in the assembly code:")
		fmt.Println(errorManager.Summary())
	} else if errorManager.HasWarnings() {
		fmt.Println("\nNo errors, but found warnings:")
		fmt.Println(errorManager.Summary())
	} else {
		fmt.Println("\nNo errors or warnings detected.")
	}

	if ast != nil {
		fmt.Printf("\nSuccessfully created AST with %d nodes.\n", countASTNodes(ast))
	}
}

// Helper function to count nodes in the AST
func countASTNodes(node *parser.ASTNode) int {
	if node == nil {
		return 0
	}

	count := 1 // Count this node
	for _, child := range node.Children {
		count += countASTNodes(child)
	}

	return count
}
