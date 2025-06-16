// Package lexer implements the lexical analyzer for the VTX1 assembly language.
// It converts source code into a stream of tokens that can be processed by the parser.
package lexer

import (
	"fmt"
	"strings"
)

// TokenType represents the type of a token in the VTX1 assembly language
type TokenType int

// Token types for the VTX1 assembly language
const (
	EOF TokenType = iota
	ILLEGAL
	COMMENT
	WHITESPACE
	NEWLINE

	// Identifiers and literals
	IDENTIFIER
	DECIMAL
	HEXADECIMAL
	BINARY
	TERNARY
	STRING

	// Registers
	GPR         // General Purpose Register (T0-T6)
	SPECIAL_REG // Special Register (TA, TB, TC, TS, TI)
	VECTOR_REG  // Vector Register (VA, VB, VC)
	FP_REG      // Floating Point Register (FA, FB, FC)

	// Operators
	PLUS
	MINUS
	ASTERISK
	SLASH
	PERCENT
	COMMA
	COLON
	LBRACKET // [
	RBRACKET // ]
	LSQUARE  // [
	RSQUARE  // ]
	LPAREN   // (
	RPAREN   // )

	// Instructions - ALU operations (from EBNF grammar)
	OP_ADD
	OP_SUB
	OP_MUL
	OP_AND
	OP_OR
	OP_NOT
	OP_XOR
	OP_SHL
	OP_SHR
	OP_ROL
	OP_ROR
	OP_CMP
	OP_TEST
	OP_INC
	OP_DEC
	OP_NEG
	OP_DIV
	OP_MOD
	OP_UDIV
	OP_UMOD
	OP_SQRT
	OP_ABS

	// Instructions - Memory operations (from EBNF grammar)
	OP_LD
	OP_ST
	OP_VLD
	OP_VST
	OP_FLD
	OP_FST
	OP_LEA
	OP_PUSH
	OP_POP // Added missing POP opcode
	OP_CACHE
	OP_FLUSH
	OP_MEMBAR

	// Instructions - Control flow (from EBNF grammar)
	OP_JMP
	OP_JAL
	OP_JR
	OP_JALR
	OP_BEQ
	OP_BNE
	OP_BLT
	OP_BGE
	OP_BLTU
	OP_BGEU
	OP_BGT // Added missing opcode
	OP_BLE // Added missing opcode
	OP_CALL
	OP_RET
	OP_SYSCALL
	OP_BREAK
	OP_HALT

	// Instructions - Vector operations (from EBNF grammar)
	OP_VADD
	OP_VSUB
	OP_VMUL
	OP_VAND
	OP_VOR
	OP_VNOT
	OP_VSHL
	OP_VSHR
	OP_VDOT
	OP_VREDUCE
	OP_VMAX
	OP_VMIN
	OP_VSUM
	OP_VPERM

	// Instructions - Floating point operations (from EBNF grammar)
	OP_FADD
	OP_FSUB
	OP_FMUL
	OP_FCMP
	OP_FMOV
	OP_FNEG
	OP_SIN
	OP_COS
	OP_TAN
	OP_ASIN
	OP_ACOS
	OP_ATAN
	OP_EXP
	OP_LOG

	// Instructions - System operations (from EBNF grammar)
	OP_NOP
	OP_WFI

	// Directives (from EBNF grammar)
	DIR_ORG
	DIR_DB
	DIR_DW
	DIR_DT
	DIR_EQU
	DIR_INCLUDE
	DIR_SECTION
	DIR_ALIGN
	DIR_SPACE
)

// Token represents a lexical token in the VTX1 assembly language
type Token struct {
	Type    TokenType // The type of this token
	Literal string    // The literal text of this token
	Line    int       // The line number where this token appears
	Column  int       // The column number where this token starts
}

// String returns a string representation of the token
func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %v, Literal: %q, Line: %d, Column: %d}", t.Type, t.Literal, t.Line, t.Column)
}

// Position returns a string representation of the token's position
func (t Token) Position() string {
	return fmt.Sprintf("%d:%d", t.Line, t.Column)
}

// Lexer performs lexical analysis of VTX1 assembly code
type Lexer struct {
	input        string // The input source code
	position     int    // Current position in input (points to current char)
	readPosition int    // Current reading position in input (after current char)
	char         byte   // Current character under examination
	line         int    // Current line number
	column       int    // Current column number
}

// New creates a new Lexer for the given input string
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

// readChar advances the lexer to the next character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0 // ASCII NUL (end of input)
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	l.column++
	// Track newlines to maintain line/column information
	if l.char == '\n' {
		l.line++
		l.column = 0
	}
}

// peekChar returns the next character without advancing the position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	// Skip whitespace except newlines (which we want to track)
	l.skipWhitespace()

	// Note the start position of this token
	tok.Line = l.line
	tok.Column = l.column

	switch l.char {
	case ';':
		tok.Type = COMMENT
		tok.Literal = l.readComment()
	case ',':
		tok = l.newToken(COMMA, string(l.char))
	case ':':
		tok = l.newToken(COLON, string(l.char))
	case '+':
		tok = l.newToken(PLUS, string(l.char))
	case '-':
		tok = l.newToken(MINUS, string(l.char))
	case '*':
		tok = l.newToken(ASTERISK, string(l.char))
	case '/':
		tok = l.newToken(SLASH, string(l.char))
	case '%':
		// Check if this is a ternary literal (should have +, -, 0 next)
		if l.peekChar() == '+' || l.peekChar() == '-' || l.peekChar() == '0' {
			// This is a ternary literal
			return l.readTernary()
		}
		tok = l.newToken(PERCENT, string(l.char))
	case '[':
		// Could be either LBRACKET (VLIW) or LSQUARE (memory reference)
		// For now, we'll treat both as LSQUARE and let the parser distinguish
		tok = l.newToken(LSQUARE, string(l.char))
	case ']':
		// Could be either RBRACKET (VLIW) or RSQUARE (memory reference)
		// For now, we'll treat both as RSQUARE and let the parser distinguish
		tok = l.newToken(RSQUARE, string(l.char))
	case '(':
		tok = l.newToken(LPAREN, string(l.char))
	case ')':
		tok = l.newToken(RPAREN, string(l.char))
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case '\n', '\r':
		if l.char == '\r' && l.peekChar() == '\n' {
			l.readChar() // Consume the \n in \r\n
		}
		tok.Type = NEWLINE
		tok.Literal = "\\n"
	case 0:
		tok.Type = EOF
		tok.Literal = ""
	case '.':
		// This is likely a directive
		directive := l.readDirective()
		tok.Literal = directive
		tok.Type = l.lookupDirective(directive)
	default:
		if isLetter(l.char) {
			// Handle identifiers, keywords, and directives
			tok.Literal = l.readIdentifier()
			tok.Type = l.lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			// Handle numeric literals
			return l.readNumber()
		} else {
			// Unknown character
			tok = l.newToken(ILLEGAL, string(l.char))
		}
	}

	l.readChar()
	return tok
}

// newToken creates a new token with the given type and literal
func (l *Lexer) newToken(tokenType TokenType, literal string) Token {
	return Token{Type: tokenType, Literal: literal, Line: l.line, Column: l.column}
}

// skipWhitespace skips any whitespace characters (except newlines)
func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' {
		l.readChar()
	}
}

// readComment reads a comment until the end of the line
func (l *Lexer) readComment() string {
	position := l.position
	for l.char != '\n' && l.char != '\r' && l.char != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString reads a string literal enclosed in double quotes
func (l *Lexer) readString() string {
	position := l.position + 1 // Skip the opening quote
	for {
		l.readChar()
		if l.char == '"' || l.char == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// readIdentifier reads an identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.char) || isDigit(l.char) || l.char == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readDirective reads a directive starting with '.'
func (l *Lexer) readDirective() string {
	position := l.position
	// Read the dot
	l.readChar()

	// Read the directive name (letters only)
	for isLetter(l.char) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// readNumber reads a numeric literal
func (l *Lexer) readNumber() Token {
	// Check for hexadecimal (0x...)
	if l.char == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		return l.readHexadecimal()
	}

	// Check for binary (0b...)
	if l.char == '0' && (l.peekChar() == 'b' || l.peekChar() == 'B') {
		return l.readBinary()
	}

	// Check for ternary (0t...)
	if l.char == '0' && (l.peekChar() == 't' || l.peekChar() == 'T') {
		return l.readTernary()
	}

	// Otherwise it's a decimal number
	return l.readDecimal()
}

// readDecimal reads a decimal number
func (l *Lexer) readDecimal() Token {
	tok := Token{Type: DECIMAL, Line: l.line, Column: l.column}
	position := l.position

	for isDigit(l.char) {
		l.readChar()
	}

	tok.Literal = l.input[position:l.position]
	return tok
}

// readHexadecimal reads a hexadecimal number (0x...)
func (l *Lexer) readHexadecimal() Token {
	tok := Token{Type: HEXADECIMAL, Line: l.line, Column: l.column}
	position := l.position

	// Skip the '0x' prefix
	l.readChar() // Skip 0
	l.readChar() // Skip x

	for isHexDigit(l.char) {
		l.readChar()
	}

	tok.Literal = l.input[position:l.position]
	return tok
}

// readBinary reads a binary number (0b...)
func (l *Lexer) readBinary() Token {
	tok := Token{Type: BINARY, Line: l.line, Column: l.column}
	position := l.position

	// Skip the '0b' prefix
	l.readChar() // Skip 0
	l.readChar() // Skip b

	for l.char == '0' || l.char == '1' {
		l.readChar()
	}

	tok.Literal = l.input[position:l.position]
	return tok
}

// readTernary reads a balanced ternary number (0t...) with +, -, 0
func (l *Lexer) readTernary() Token {
	tok := Token{Type: TERNARY, Line: l.line, Column: l.column}
	position := l.position

	// Skip the '0t' prefix
	l.readChar() // Skip 0
	l.readChar() // Skip t

	// Read ternary digits (+, -, 0)
	for l.char == '+' || l.char == '-' || l.char == '0' {
		l.readChar()
	}

	tok.Literal = l.input[position:l.position]
	return tok
}

// lookupDirective determines the token type for a directive
func (l *Lexer) lookupDirective(directive string) TokenType {
	switch strings.ToUpper(directive) {
	case ".ORG":
		return DIR_ORG
	case ".DB":
		return DIR_DB
	case ".DW":
		return DIR_DW
	case ".DT":
		return DIR_DT
	case ".EQU":
		return DIR_EQU
	case ".INCLUDE":
		return DIR_INCLUDE
	case ".SECTION":
		return DIR_SECTION
	case ".ALIGN":
		return DIR_ALIGN
	case ".SPACE":
		return DIR_SPACE
	default:
		return IDENTIFIER // Unknown directive
	}
}

// lookupIdent checks if the identifier is a keyword or instruction
func (l *Lexer) lookupIdent(ident string) TokenType {
	// Convert to uppercase for case-insensitive matching
	upperIdent := strings.ToUpper(ident)

	// Check for instructions
	switch upperIdent {
	case "ADD":
		return OP_ADD
	case "SUB":
		return OP_SUB
	case "MUL":
		return OP_MUL
	case "DIV":
		return OP_DIV
	case "AND":
		return OP_AND
	case "OR":
		return OP_OR
	case "XOR":
		return OP_XOR
	case "NOT":
		return OP_NOT
	case "SHL":
		return OP_SHL
	case "SHR":
		return OP_SHR
	case "ROL":
		return OP_ROL
	case "ROR":
		return OP_ROR
	case "CMP":
		return OP_CMP
	case "TEST":
		return OP_TEST
	case "INC":
		return OP_INC
	case "DEC":
		return OP_DEC
	case "NEG":
		return OP_NEG
	case "MOD":
		return OP_MOD
	case "UDIV":
		return OP_UDIV
	case "UMOD":
		return OP_UMOD
	case "SQRT":
		return OP_SQRT
	case "ABS":
		return OP_ABS
	case "LD":
		return OP_LD
	case "ST":
		return OP_ST
	case "VLD":
		return OP_VLD
	case "VST":
		return OP_VST
	case "FLD":
		return OP_FLD
	case "FST":
		return OP_FST
	case "LEA":
		return OP_LEA
	case "PUSH":
		return OP_PUSH
	case "POP": // Added missing opcode mapping
		return OP_POP
	case "CACHE":
		return OP_CACHE
	case "FLUSH":
		return OP_FLUSH
	case "MEMBAR":
		return OP_MEMBAR
	case "JMP":
		return OP_JMP
	case "JAL":
		return OP_JAL
	case "JR":
		return OP_JR
	case "JALR":
		return OP_JALR
	case "BEQ":
		return OP_BEQ
	case "BNE":
		return OP_BNE
	case "BLT":
		return OP_BLT
	case "BGE":
		return OP_BGE
	case "BLTU":
		return OP_BLTU
	case "BGEU":
		return OP_BGEU
	case "BGT": // Added missing opcode
		return OP_BGT
	case "BLE": // Added missing opcode
		return OP_BLE
	case "CALL":
		return OP_CALL
	case "RET":
		return OP_RET
	case "SYSCALL":
		return OP_SYSCALL
	case "BREAK":
		return OP_BREAK
	case "HALT":
		return OP_HALT
	}

	// Check for registers
	if len(upperIdent) >= 2 {
		prefix := upperIdent[0:1]
		if prefix == "T" && len(upperIdent) == 2 && upperIdent[1] >= '0' && upperIdent[1] <= '6' {
			return GPR // T0-T6
		} else if prefix == "T" && (upperIdent == "TA" || upperIdent == "TB" || upperIdent == "TC" || upperIdent == "TS" || upperIdent == "TI") {
			return SPECIAL_REG
		} else if prefix == "V" && (upperIdent == "VA" || upperIdent == "VB" || upperIdent == "VC") {
			return VECTOR_REG
		} else if prefix == "F" && (upperIdent == "FA" || upperIdent == "FB" || upperIdent == "FC") {
			return FP_REG
		}
	}

	// Otherwise, it's just an identifier
	return IDENTIFIER
}

// isLetter returns true if the character is a letter or underscore
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit returns true if the character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isHexDigit returns true if the character is a hexadecimal digit
func isHexDigit(ch byte) bool {
	return isDigit(ch) || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}
