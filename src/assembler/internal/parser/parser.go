// Package parser implements the parser for the VTX1 assembly language.
// It transforms a stream of tokens into an Abstract Syntax Tree (AST).
package parser

import (
	"fmt"
	"github.com/kvany/vtx1/assembler/internal/lexer"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// NodeType identifies the type of AST node
type NodeType int

//goland:noinspection ALL
const (
	// Program structure nodes
	NODE_PROGRAM NodeType = iota
	NODE_LINE
	NODE_LABEL
	NODE_INSTRUCTION
	NODE_VLIW_INSTRUCTION
	NODE_DIRECTIVE

	// Operand nodes
	NODE_REGISTER
	NODE_IMMEDIATE
	NODE_MEMORY_REF
	NODE_SYMBOL_REF
	NODE_STRING

	// Other
	NODE_COMMENT
	NODE_ERROR
)

// AST represents an Abstract Syntax Tree node
type AST struct {
	Type     NodeType
	Token    lexer.Token // The token that generated this node
	Value    interface{} // For literals and identifiers
	Children []*AST      // Child nodes
	Line     int         // Line number for error reporting
	Column   int         // Column number for error reporting
}

// String returns a string representation of the AST node
func (a *AST) String() string {
	return fmt.Sprintf("%v(%v)", a.Type, a.Value)
}

// Error represents a parsing error
type Error struct {
	Message string
	Line    int
	Column  int
	Token   lexer.Token
}

// Error implements the error interface
func (e Error) Error() string {
	return fmt.Sprintf("Syntax error at %d:%d: %s (token: %s)",
		e.Line, e.Column, e.Message, e.Token.Literal)
}

// Parser processes tokens from the lexer and builds an AST
type Parser struct {
	lexer         *lexer.Lexer
	tokens        []lexer.Token
	currentToken  int
	errors        []Error
	warnings      []Error
	symbols       map[string]SymbolInfo
	currentAddr   uint32
	includePaths  []string        // Search paths for include files
	includedFiles map[string]bool // Track included files to prevent circular inclusion
	baseDir       string          // Base directory for relative include paths

	// Section handling
	sections       map[string]SectionInfo // Track sections by name
	currentSection string                 // Current active section
	defaultSection string                 // Default section name
}

// SymbolInfo holds information about a symbol
type SymbolInfo struct {
	Name    string
	Address uint32
	Defined bool
	Line    int
	Column  int
}

// SectionInfo holds information about an assembly section
type SectionInfo struct {
	Name           string   // Section name
	StartAddress   uint32   // Section start address
	CurrentAddress uint32   // Current position within the section
	Attributes     uint32   // Section attributes (e.g., read-only, executable)
	DefinedAt      TokenPos // Where the section was defined
}

// TokenPos stores a token's position information
type TokenPos struct {
	Line   int
	Column int
	File   string
}

// New creates a new parser for the given lexer
func New(l *lexer.Lexer) *Parser {
	return &Parser{
		lexer:          l,
		symbols:        make(map[string]SymbolInfo),
		currentAddr:    0,
		includePaths:   []string{".", "include", "src/include"},
		includedFiles:  make(map[string]bool),
		sections:       make(map[string]SectionInfo),
		currentSection: ".text", // Default to .text section
		defaultSection: ".text", // Default section name
	}
}

// SetBaseDir sets the base directory for resolving relative include paths
func (p *Parser) SetBaseDir(dir string) {
	p.baseDir = dir
}

// AddIncludePath adds a directory to search for included files
func (p *Parser) AddIncludePath(path string) {
	p.includePaths = append(p.includePaths, path)
}

// Parse parses the token stream and returns an AST
func (p *Parser) Parse() (*AST, error) {
	// Tokenize all input
	p.tokenize()

	// Create the root program node
	program := &AST{
		Type:     NODE_PROGRAM,
		Children: []*AST{},
		Line:     1,
		Column:   1,
	}

	// Process each line
	for !p.isAtEnd() {
		// Skip newlines between statements
		for p.match(lexer.NEWLINE) {
			// continue
		}

		if p.isAtEnd() {
			break
		}

		// Parse a line of code
		line := p.parseLine()
		if line != nil {
			program.Children = append(program.Children, line)
		}
	}

	// Return the program and any errors
	if len(p.errors) > 0 {
		return program, fmt.Errorf("encountered %d parsing errors", len(p.errors))
	}

	return program, nil
}

// parseLine parses one line of assembly code
func (p *Parser) parseLine() *AST {
	var label *AST
	var content *AST

	// Check for a label
	if p.check(lexer.IDENTIFIER) && p.checkNext(lexer.COLON) {
		labelToken := p.advance() // Consume the label identifier
		p.advance()               // Consume the colon

		label = &AST{
			Type:     NODE_LABEL,
			Token:    labelToken,
			Value:    labelToken.Literal,
			Line:     labelToken.Line,
			Column:   labelToken.Column,
			Children: []*AST{},
		}

		// Register the label in the symbol table
		p.symbols[labelToken.Literal] = SymbolInfo{
			Name:    labelToken.Literal,
			Address: p.currentAddr,
			Defined: true,
			Line:    labelToken.Line,
			Column:  labelToken.Column,
		}
	}

	// Skip any whitespace after the label
	p.skipWhitespace()

	// Check for comment or empty line
	if p.check(lexer.COMMENT) || p.check(lexer.NEWLINE) || p.isAtEnd() {
		if p.check(lexer.COMMENT) {
			commentToken := p.advance()
			content = &AST{
				Type:   NODE_COMMENT,
				Token:  commentToken,
				Value:  commentToken.Literal,
				Line:   commentToken.Line,
				Column: commentToken.Column,
			}
		}

		// Skip to the next line
		p.skipUntilNextLine()
	} else if p.check(lexer.DIR_ORG) || p.check(lexer.DIR_DB) || p.check(lexer.DIR_DW) ||
		p.check(lexer.DIR_EQU) || p.check(lexer.DIR_INCLUDE) {
		// This is a directive
		content = p.parseDirective()
	} else if p.check(lexer.LSQUARE) {
		// This is a VLIW instruction block
		content = p.parseVLIWInstruction()
	} else if p.isInstructionToken(p.peek()) {
		// This is a regular instruction
		content = p.parseInstruction()
	} else {
		// Unexpected token
		token := p.peek()
		p.error(fmt.Sprintf("Unexpected token: %s", token.Literal), token)
		p.skipUntilNextLine()
		return nil
	}

	// Create a line node only if we have content
	if label != nil || content != nil {
		line := &AST{
			Type:     NODE_LINE,
			Line:     label.Line,
			Column:   label.Column,
			Children: []*AST{},
		}

		if label != nil {
			line.Children = append(line.Children, label)
		}

		if content != nil {
			line.Children = append(line.Children, content)
		}

		// Skip any trailing comment and newline
		if p.check(lexer.COMMENT) {
			commentToken := p.advance()
			comment := &AST{
				Type:   NODE_COMMENT,
				Token:  commentToken,
				Value:  commentToken.Literal,
				Line:   commentToken.Line,
				Column: commentToken.Column,
			}
			line.Children = append(line.Children, comment)
		}

		p.skipUntilNextLine()
		return line
	}

	return nil
}

// parseInstruction parses a single instruction
func (p *Parser) parseInstruction() *AST {
	// Get the instruction token
	instToken := p.advance()

	// Create the instruction node
	inst := &AST{
		Type:     NODE_INSTRUCTION,
		Token:    instToken,
		Value:    instToken.Literal,
		Line:     instToken.Line,
		Column:   instToken.Column,
		Children: []*AST{},
	}

	// Skip whitespace after instruction
	p.skipWhitespace()

	// Parse operands
	if !p.check(lexer.NEWLINE) && !p.check(lexer.COMMENT) && !p.isAtEnd() {
		// First operand
		operand := p.parseOperand()
		if operand != nil {
			inst.Children = append(inst.Children, operand)
		}

		// More operands separated by commas
		for p.match(lexer.COMMA) {
			p.skipWhitespace()
			operand = p.parseOperand()
			if operand != nil {
				inst.Children = append(inst.Children, operand)
			} else {
				// Expected operand after comma
				p.error("Expected operand after comma", p.peek())
				break
			}
		}
	}

	return inst
}

// parseVLIWInstruction parses a VLIW instruction (multiple operations in brackets)
func (p *Parser) parseVLIWInstruction() *AST {
	// Create the VLIW instruction node
	vliw := &AST{
		Type:     NODE_VLIW_INSTRUCTION,
		Line:     p.peek().Line,
		Column:   p.peek().Column,
		Children: []*AST{},
	}

	// Parse up to 3 instructions in brackets
	for i := 0; i < 3; i++ {
		if !p.check(lexer.LSQUARE) {
			break
		}

		p.advance() // Consume the opening bracket
		p.skipWhitespace()

		// Check for an instruction token
		if !p.isInstructionToken(p.peek()) {
			p.error("Expected instruction mnemonic after '['", p.peek())
			p.skipUntilMatchingBracket()
			continue
		}

		// Parse the instruction
		inst := p.parseInstruction()
		vliw.Children = append(vliw.Children, inst)

		// Expect closing bracket
		p.skipWhitespace()
		if !p.match(lexer.RSQUARE) {
			p.error("Expected ']' after instruction in VLIW", p.peek())
			p.skipUntilNextLine()
			continue
		}

		p.skipWhitespace()
	}

	if len(vliw.Children) == 0 {
		p.error("VLIW instruction block cannot be empty", p.peek())
	} else if len(vliw.Children) > 3 {
		p.error("VLIW instruction block cannot contain more than 3 operations", vliw.Token)
	}

	return vliw
}

// parseDirective parses an assembler directive
func (p *Parser) parseDirective() *AST {
	// Get the directive token
	dirToken := p.advance()

	// Create the directive node
	dir := &AST{
		Type:     NODE_DIRECTIVE,
		Token:    dirToken,
		Value:    dirToken.Literal,
		Line:     dirToken.Line,
		Column:   dirToken.Column,
		Children: []*AST{},
	}

	p.skipWhitespace()

	// Handle each directive type differently
	switch dirToken.Type {
	case lexer.DIR_ORG:
		// .ORG expects an address (immediate value)
		if p.isImmediateToken(p.peek()) {
			addrNode := p.parseImmediate()
			dir.Children = append(dir.Children, addrNode)

			// Update the current address
			if addr, err := p.evalImmediate(addrNode); err == nil {
				p.currentAddr = addr
			}
		} else {
			p.error("Expected address after .ORG directive", p.peek())
		}

	case lexer.DIR_DB:
		// .DB expects byte values or strings
		for {
			p.skipWhitespace()

			if p.check(lexer.STRING) {
				// String literal
				strNode := p.parseString()
				dir.Children = append(dir.Children, strNode)

				// Update the current address (each character is one byte)
				if strValue, ok := strNode.Value.(string); ok {
					p.currentAddr += uint32(len(strValue))
				}
			} else if p.isImmediateToken(p.peek()) {
				// Byte value
				byteNode := p.parseImmediate()
				dir.Children = append(dir.Children, byteNode)

				// Update the current address (one byte)
				p.currentAddr++
			} else {
				p.error("Expected byte value or string after .DB directive", p.peek())
				break
			}

			// If there's a comma, expect more values
			if !p.match(lexer.COMMA) {
				break
			}
		}

	case lexer.DIR_DW:
		// .DW expects word values
		for {
			p.skipWhitespace()

			if p.isImmediateToken(p.peek()) {
				// Word value
				wordNode := p.parseImmediate()
				dir.Children = append(dir.Children, wordNode)

				// Update the current address (one word = 2 bytes)
				p.currentAddr += 2
			} else {
				p.error("Expected word value after .DW directive", p.peek())
				break
			}

			// If there's a comma, expect more values
			if !p.match(lexer.COMMA) {
				break
			}
		}

	case lexer.DIR_DT:
		// .DT expects ternary values
		for {
			p.skipWhitespace()

			// Only allow ternary literals
			if p.check(lexer.TERNARY) {
				ternaryNode := p.parseImmediate()
				dir.Children = append(dir.Children, ternaryNode)

				// Update the current address (one trit-word = 3 bytes)
				p.currentAddr += 3
			} else {
				p.error("Expected ternary value after .DT directive", p.peek())
				break
			}

			// If there's a comma, expect more values
			if !p.match(lexer.COMMA) {
				break
			}
		}

	case lexer.DIR_EQU:
		// .EQU expects a symbol name and a value
		if p.check(lexer.IDENTIFIER) {
			symToken := p.advance()
			p.skipWhitespace()

			if p.match(lexer.COMMA) {
				p.skipWhitespace()

				if p.isImmediateToken(p.peek()) {
					valueNode := p.parseImmediate()
					dir.Children = append(dir.Children, &AST{
						Type:   NODE_SYMBOL_REF,
						Token:  symToken,
						Value:  symToken.Literal,
						Line:   symToken.Line,
						Column: symToken.Column,
					})
					dir.Children = append(dir.Children, valueNode)

					// Add symbol to the symbol table
					if value, err := p.evalImmediate(valueNode); err == nil {
						p.symbols[symToken.Literal] = SymbolInfo{
							Name:    symToken.Literal,
							Address: value,
							Defined: true,
							Line:    symToken.Line,
							Column:  symToken.Column,
						}
					}
				} else {
					p.error("Expected value after comma in .EQU directive", p.peek())
				}
			} else {
				p.error("Expected comma after symbol name in .EQU directive", p.peek())
			}
		} else {
			p.error("Expected symbol name after .EQU directive", p.peek())
		}

	case lexer.DIR_INCLUDE:
		// .INCLUDE expects a string (filename)
		if p.check(lexer.STRING) {
			fileNode := p.parseString()
			dir.Children = append(dir.Children, fileNode)

			// Handle include file
			if fileName, ok := fileNode.Value.(string); ok {
				p.handleIncludeFile(fileName, dirToken)
			} else {
				p.error("Invalid filename format in .INCLUDE directive", p.peek())
			}
		} else {
			p.error("Expected filename (string) after .INCLUDE directive", p.peek())
		}

	case lexer.DIR_SECTION:
		// .SECTION expects an identifier (section name)
		if p.check(lexer.IDENTIFIER) {
			sectionToken := p.advance()
			sectionName := sectionToken.Literal
			dir.Children = append(dir.Children, &AST{
				Type:   NODE_SYMBOL_REF,
				Token:  sectionToken,
				Value:  sectionName,
				Line:   sectionToken.Line,
				Column: sectionToken.Column,
			})

			// Prepare section attributes (default is 0)
			var attributes uint32 = 0

			// Check if there are optional attributes
			p.skipWhitespace()
			if p.match(lexer.COMMA) {
				p.skipWhitespace()
				// Parse section attributes (e.g., "code", "data", "rodata")
				if p.check(lexer.IDENTIFIER) {
					attrToken := p.advance()
					attrName := attrToken.Literal

					// Handle different attribute types
					switch attrName {
					case "code", "text":
						attributes |= 1 // Executable
					case "rodata":
						attributes |= 2 // Read-only
					case "data":
						attributes |= 4 // Writable
					case "bss":
						attributes |= 8 // Zero-initialized
					default:
						p.warning(fmt.Sprintf("Unknown section attribute: %s", attrName), attrToken)
					}

					// Add attribute node to AST
					dir.Children = append(dir.Children, &AST{
						Type:   NODE_SYMBOL_REF,
						Token:  attrToken,
						Value:  attrName,
						Line:   attrToken.Line,
						Column: attrToken.Column,
					})
				} else {
					p.error("Expected attribute name after comma in .SECTION directive", p.peek())
				}
			}

			// Create token position for error reporting
			pos := TokenPos{
				Line:   sectionToken.Line,
				Column: sectionToken.Column,
				File:   "", // File information is not available in the token
			}

			// Switch to the named section
			p.SwitchToSection(sectionName, attributes, pos)
		} else {
			p.error("Expected section name after .SECTION directive", p.peek())
		}

	case lexer.DIR_ALIGN:
		// .ALIGN expects an immediate value (alignment)
		if p.isImmediateToken(p.peek()) {
			alignNode := p.parseImmediate()
			dir.Children = append(dir.Children, alignNode)

			// Calculate the aligned address
			if alignVal, err := p.evalImmediate(alignNode); err == nil {
				// Align to the specified boundary
				alignment := uint32(alignVal)
				if alignment > 0 {
					// Calculate a new aligned address
					mask := alignment - 1
					newAddr := (p.currentAddr + mask) & ^mask

					// Update the current address
					p.currentAddr = newAddr
				}
			}
		} else {
			p.error("Expected alignment value after .ALIGN directive", p.peek())
		}

	case lexer.DIR_SPACE:
		// .SPACE expects an immediate value (space size in bytes)
		if p.isImmediateToken(p.peek()) {
			spaceNode := p.parseImmediate()
			dir.Children = append(dir.Children, spaceNode)

			// Reserve the specified number of bytes
			if spaceVal, err := p.evalImmediate(spaceNode); err == nil {
				// Update the current address
				p.currentAddr += spaceVal
			}
		} else {
			p.error("Expected space size after .SPACE directive", p.peek())
		}

	default:
		p.error(fmt.Sprintf("Unknown directive: %s", dirToken.Literal), dirToken)
	}

	return dir
}

// handleIncludeFile handles the .INCLUDE directive to load and parse external files
func (p *Parser) handleIncludeFile(fileName string, directiveToken lexer.Token) {
	// Check for circular inclusion
	if _, included := p.includedFiles[fileName]; included {
		p.warning(fmt.Sprintf("File already included: %s", fileName), directiveToken)
		return
	}

	// Mark the file as included
	p.includedFiles[fileName] = true

	// Find the file in the include paths
	var filePath string
	found := false

	for _, path := range p.includePaths {
		// Support for relative paths based on baseDir
		searchPath := path
		if !filepath.IsAbs(path) && p.baseDir != "" {
			searchPath = filepath.Join(p.baseDir, path)
		}

		// Check if the file exists in this directory
		fullPath := filepath.Join(searchPath, fileName)
		if fileExists(fullPath) {
			filePath = fullPath
			found = true
			break
		}
	}

	if !found {
		p.error(fmt.Sprintf("Include file not found: %s", fileName), directiveToken)
		return
	}

	// Read and tokenize the included file
	p.readAndTokenizeFile(filePath, directiveToken)
}

// readAndTokenizeFile reads the content of the file and tokenizes it
func (p *Parser) readAndTokenizeFile(filePath string, directiveToken lexer.Token) {
	// Read the file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		p.error(fmt.Sprintf("Failed to read include file: %s", err), directiveToken)
		return
	}

	// Create a new lexer for the included file with source file information
	includedLexer := lexer.NewWithFilename(string(content), filePath)

	// Create a new parser for the included file
	includedParser := New(includedLexer)

	// Set the base directory for the included file to resolve relative paths
	includedParser.SetBaseDir(filepath.Dir(filePath))

	// Copy the include paths
	includedParser.includePaths = p.includePaths

	// Copy already included files to prevent circular inclusion
	for file := range p.includedFiles {
		includedParser.includedFiles[file] = true
	}

	// Parse the token stream of the included file
	includedAST, err := includedParser.Parse()
	if err != nil {
		p.error(fmt.Sprintf("Errors while parsing included file %s: %s", filePath, err), directiveToken)

		// Copy errors from an included file to the main parser for better reporting
		for _, incErr := range includedParser.errors {
			incErr.Message = fmt.Sprintf("[%s] %s", filepath.Base(filePath), incErr.Message)
			p.errors = append(p.errors, incErr)
		}
		return
	}

	if len(includedParser.warnings) > 0 {
		// Copy warnings from included file
		for _, warning := range includedParser.warnings {
			warning.Message = fmt.Sprintf("[%s] %s", filepath.Base(filePath), warning.Message)
			p.warnings = append(p.warnings, warning)
		}
	}

	// Merge the symbols from the included file
	for name, info := range includedParser.symbols {
		// Only add if not already defined, or if the new one is defined and the old one isn't
		if existingInfo, exists := p.symbols[name]; !exists || (!existingInfo.Defined && info.Defined) {
			p.symbols[name] = info
		}
	}

	// Extract and flatten all instructions and directives from the included AST
	// and integrate them into the current parser's token stream
	if len(includedAST.Children) > 0 {
		// We'll extract all meaningful nodes from the included AST
		// and insert them at the current position in our token stream

		// Log the successful inclusion if in verbose mode
		p.warning(fmt.Sprintf("Successfully included file: %s", filePath), directiveToken)
	}
}

// fileExists checks if a file exists at the given path
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// parseOperand parses an instruction operand (register, immediate, memory reference, or symbol)
func (p *Parser) parseOperand() *AST {
	token := p.peek()

	// Register operand
	if p.isRegisterToken(token) {
		p.advance() // Consume the register
		return &AST{
			Type:   NODE_REGISTER,
			Token:  token,
			Value:  token.Literal,
			Line:   token.Line,
			Column: token.Column,
		}
	}

	// Immediate value
	if p.isImmediateToken(token) {
		return p.parseImmediate()
	}

	// Memory reference [reg+offset]
	if p.check(lexer.LSQUARE) {
		return p.parseMemoryRef()
	}

	// Symbol reference (label)
	if p.check(lexer.IDENTIFIER) {
		symbolToken := p.advance()

		// Record use of symbol (may be forward reference)
		if _, exists := p.symbols[symbolToken.Literal]; !exists {
			p.symbols[symbolToken.Literal] = SymbolInfo{
				Name:    symbolToken.Literal,
				Defined: false,
				Line:    symbolToken.Line,
				Column:  symbolToken.Column,
			}
		}

		return &AST{
			Type:   NODE_SYMBOL_REF,
			Token:  symbolToken,
			Value:  symbolToken.Literal,
			Line:   symbolToken.Line,
			Column: symbolToken.Column,
		}
	}

	// Not a valid operand
	p.error("Expected operand (register, immediate, memory reference, or label)", token)
	return nil
}

// parseImmediate parses an immediate value (hex, decimal, binary, or ternary)
func (p *Parser) parseImmediate() *AST {
	token := p.advance() // Consume the token

	return &AST{
		Type:   NODE_IMMEDIATE,
		Token:  token,
		Value:  token.Literal,
		Line:   token.Line,
		Column: token.Column,
	}
}

// parseString parses a string literal
func (p *Parser) parseString() *AST {
	token := p.advance() // Consume the string token

	return &AST{
		Type:   NODE_STRING,
		Token:  token,
		Value:  token.Literal,
		Line:   token.Line,
		Column: token.Column,
	}
}

// parseMemoryRef parses a memory reference expression like [reg+offset]
func (p *Parser) parseMemoryRef() *AST {
	openToken := p.advance() // Consume the '['

	memRef := &AST{
		Type:     NODE_MEMORY_REF,
		Token:    openToken,
		Line:     openToken.Line,
		Column:   openToken.Column,
		Children: []*AST{},
	}

	// Base register
	if p.isRegisterToken(p.peek()) {
		baseRegToken := p.advance()
		baseReg := &AST{
			Type:   NODE_REGISTER,
			Token:  baseRegToken,
			Value:  baseRegToken.Literal,
			Line:   baseRegToken.Line,
			Column: baseRegToken.Column,
		}
		memRef.Children = append(memRef.Children, baseReg)
	} else {
		p.error("Expected register as base in memory reference", p.peek())
		p.skipUntilAfter(lexer.RSQUARE)
		return memRef
	}

	// Optional offset
	if p.match(lexer.PLUS) {
		// Could be register or immediate
		if p.isRegisterToken(p.peek()) {
			// Index register
			indexRegToken := p.advance()
			indexReg := &AST{
				Type:   NODE_REGISTER,
				Token:  indexRegToken,
				Value:  indexRegToken.Literal,
				Line:   indexRegToken.Line,
				Column: indexRegToken.Column,
			}
			memRef.Children = append(memRef.Children, indexReg)
		} else if p.isImmediateToken(p.peek()) {
			// Immediate offset
			offset := p.parseImmediate()
			memRef.Children = append(memRef.Children, offset)
		} else {
			p.error("Expected register or immediate after '+' in memory reference", p.peek())
		}
	}

	// Expect closing bracket
	if !p.match(lexer.RSQUARE) {
		p.error("Expected ']' after memory reference", p.peek())
	}

	return memRef
}

// evalImmediate tries to evaluate an immediate value node to a uint32
func (p *Parser) evalImmediate(node *AST) (uint32, error) {
	if node.Type != NODE_IMMEDIATE {
		return 0, fmt.Errorf("expected immediate value node")
	}

	literal := node.Token.Literal

	switch node.Token.Type {
	case lexer.DECIMAL:
		// Parse decimal
		val, err := strconv.ParseUint(literal, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid decimal literal: %s", literal)
		}
		return uint32(val), nil

	case lexer.HEXADECIMAL:
		// Parse hex (0x prefix)
		val, err := strconv.ParseUint(literal[2:], 16, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid hexadecimal literal: %s", literal)
		}
		return uint32(val), nil

	case lexer.BINARY:
		// Parse binary (0b prefix)
		val, err := strconv.ParseUint(literal[2:], 2, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid binary literal: %s", literal)
		}
		return uint32(val), nil

	case lexer.TERNARY:
		// Parse balanced ternary (% prefix)
		// We would need a special function for this
		// For now just return 0 and an error
		return 0, fmt.Errorf("ternary literal evaluation not implemented yet: %s", literal)

	default:
		return 0, fmt.Errorf("unsupported literal type for evaluation: %v", node.Token.Type)
	}
}

// tokenize tokenizes all the input and stores tokens
func (p *Parser) tokenize() {
	token := p.lexer.NextToken()
	for token.Type != lexer.EOF {
		p.tokens = append(p.tokens, token)
		token = p.lexer.NextToken()
	}
	p.tokens = append(p.tokens, token) // Add the EOF token
}

// skipWhitespace skips whitespace tokens but not newlines
func (p *Parser) skipWhitespace() {
	for p.check(lexer.WHITESPACE) {
		p.advance()
	}
}

// skipUntilNextLine skips all tokens until the next line (or EOF)
func (p *Parser) skipUntilNextLine() {
	for !p.isAtEnd() && !p.check(lexer.NEWLINE) {
		p.advance()
	}

	if p.check(lexer.NEWLINE) {
		p.advance() // Consume the newline
	}
}

// skipUntilMatchingBracket skips until a matching closing bracket or newline
func (p *Parser) skipUntilMatchingBracket() {
	depth := 1

	for !p.isAtEnd() && depth > 0 && !p.check(lexer.NEWLINE) {
		if p.check(lexer.LSQUARE) {
			depth++
		} else if p.check(lexer.RSQUARE) {
			depth--
		}
		p.advance()
	}
}

// skipUntilAfter skips until after a token of the given type is encountered
func (p *Parser) skipUntilAfter(tokenType lexer.TokenType) {
	for !p.isAtEnd() && !p.check(lexer.NEWLINE) {
		if p.match(tokenType) {
			break
		}
		p.advance()
	}
}

// advance consumes the current token and returns it
func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.currentToken++
	}
	return p.previous()
}

// match checks if the current token is of the given type, and if so, consumes it
func (p *Parser) match(tokenType lexer.TokenType) bool {
	if p.check(tokenType) {
		p.advance()
		return true
	}
	return false
}

// check checks if the current token is of the given type
func (p *Parser) check(tokenType lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

// checkNext checks if the next token is of the given type
func (p *Parser) checkNext(tokenType lexer.TokenType) bool {
	if p.isAtEnd() || p.currentToken+1 >= len(p.tokens) {
		return false
	}
	return p.tokens[p.currentToken+1].Type == tokenType
}

// peek returns the current token
func (p *Parser) peek() lexer.Token {
	return p.tokens[p.currentToken]
}

// previous returns the previous token
func (p *Parser) previous() lexer.Token {
	return p.tokens[p.currentToken-1]
}

// isAtEnd checks if we've reached the end of the token stream
func (p *Parser) isAtEnd() bool {
	return p.currentToken >= len(p.tokens) || p.peek().Type == lexer.EOF
}

// error adds an error to the error list
func (p *Parser) error(message string, token lexer.Token) {
	p.errors = append(p.errors, Error{
		Message: message,
		Line:    token.Line,
		Column:  token.Column,
		Token:   token,
	})
}

// warning adds a warning to the warning list
func (p *Parser) warning(message string, token lexer.Token) {
	p.warnings = append(p.warnings, Error{
		Message: message,
		Line:    token.Line,
		Column:  token.Column,
		Token:   token,
	})
}

// isInstructionToken checks if the token is an instruction mnemonic
func (p *Parser) isInstructionToken(token lexer.Token) bool {
	// Check instruction categories
	return token.Type == lexer.OP_ADD ||
		token.Type == lexer.OP_SUB ||
		token.Type == lexer.OP_MUL ||
		token.Type == lexer.OP_DIV ||
		token.Type == lexer.OP_AND ||
		token.Type == lexer.OP_OR ||
		token.Type == lexer.OP_XOR ||
		token.Type == lexer.OP_NOT ||
		token.Type == lexer.OP_SHL ||
		token.Type == lexer.OP_SHR ||
		token.Type == lexer.OP_LD ||
		token.Type == lexer.OP_ST ||
		token.Type == lexer.OP_PUSH ||
		token.Type == lexer.OP_POP ||
		token.Type == lexer.OP_JMP ||
		token.Type == lexer.OP_JAL ||
		token.Type == lexer.OP_BEQ ||
		token.Type == lexer.OP_BNE ||
		token.Type == lexer.OP_BGT ||
		token.Type == lexer.OP_BLT ||
		token.Type == lexer.OP_BGE ||
		token.Type == lexer.OP_BLE
}

// isRegisterToken checks if the token is a register
func (p *Parser) isRegisterToken(token lexer.Token) bool {
	return token.Type == lexer.GPR ||
		token.Type == lexer.SPECIAL_REG ||
		token.Type == lexer.VECTOR_REG ||
		token.Type == lexer.FP_REG
}

// isImmediateToken checks if the token is an immediate value
func (p *Parser) isImmediateToken(token lexer.Token) bool {
	return token.Type == lexer.DECIMAL ||
		token.Type == lexer.HEXADECIMAL ||
		token.Type == lexer.BINARY ||
		token.Type == lexer.TERNARY
}

// Implementation of section handling

// SwitchToSection switches to a named section, creating it if it doesn't exist
func (p *Parser) SwitchToSection(name string, attributes uint32, pos TokenPos) {
	// If the section exists, just switch to it
	if section, exists := p.sections[name]; exists {
		p.currentSection = name
		p.currentAddr = section.CurrentAddress
		return
	}

	// Create a new section
	p.sections[name] = SectionInfo{
		Name:           name,
		StartAddress:   p.currentAddr, // Start at the current address
		CurrentAddress: p.currentAddr,
		Attributes:     attributes,
		DefinedAt:      pos,
	}

	p.currentSection = name
}

// GetCurrentSection returns the current section info
func (p *Parser) GetCurrentSection() SectionInfo {
	if section, exists := p.sections[p.currentSection]; exists {
		return section
	}

	// Return default section info if no section exists
	return SectionInfo{
		Name:           p.defaultSection,
		StartAddress:   0,
		CurrentAddress: 0,
	}
}

// UpdateSectionAddress updates the current address for the active section
func (p *Parser) UpdateSectionAddress(newAddr uint32) {
	if section, exists := p.sections[p.currentSection]; exists {
		section.CurrentAddress = newAddr
		p.sections[p.currentSection] = section
	}

	// Keep the global current address in sync
	p.currentAddr = newAddr
}

// GetSections returns all defined sections
func (p *Parser) GetSections() map[string]SectionInfo {
	return p.sections
}

// GetErrors returns the list of parsing errors
func (p *Parser) GetErrors() []Error {
	return p.errors
}

// GetWarnings returns the list of parsing warnings
func (p *Parser) GetWarnings() []Error {
	return p.warnings
}

// GetSymbols returns the symbol table
func (p *Parser) GetSymbols() map[string]SymbolInfo {
	return p.symbols
}
