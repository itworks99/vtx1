package grammar

// TokenType defines the type of token in the grammar
type TokenType int

const (
	// Token types
	TokenInvalid TokenType = iota
	TokenIdentifier
	TokenString
	TokenOperator
	TokenLeftParen
	TokenRightParen
	TokenLeftBracket
	TokenRightBracket
	TokenLeftBrace
	TokenRightBrace
	TokenAssignment  // ::=
	TokenAlternative // |
	TokenRepetition  // *
	TokenOneOrMore   // +
	TokenOptional    // ?
	TokenEOF
)

// Token represents a lexical token from the grammar file
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// IsTerminal returns true if the token represents a terminal symbol
func (t Token) IsTerminal() bool {
	return t.Type == TokenString
}
