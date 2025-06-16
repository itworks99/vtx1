package grammar

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// Lexer converts an input stream into tokens
type Lexer struct {
	reader    *bufio.Reader
	line      int
	column    int
	lastChar  rune
	lastToken Token
}

// NewLexer creates a new lexer from an io.Reader
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(r),
		line:   1,
		column: 0,
	}
}

// NextToken reads and returns the next token from the input
func (l *Lexer) NextToken() (Token, error) {
	// Skip whitespace
	for {
		ch, _, err := l.readChar()
		if err != nil {
			if err == io.EOF {
				return Token{Type: TokenEOF, Line: l.line, Column: l.column}, nil
			}
			return Token{}, err
		}

		if !unicode.IsSpace(ch) {
			l.unreadChar()
			break
		}

		if ch == '\n' {
			l.line++
			l.column = 0
		}
	}

	// Read the next character
	ch, _, err := l.readChar()
	if err != nil {
		if err == io.EOF {
			return Token{Type: TokenEOF, Line: l.line, Column: l.column}, nil
		}
		return Token{}, err
	}

	startColumn := l.column

	// Check for comments
	if ch == '/' {
		next, _, err := l.readChar()
		if err == nil {
			if next == '/' {
				// Line comment, read until end of line
				for {
					ch, _, err = l.readChar()
					if err != nil || ch == '\n' {
						break
					}
				}
				if ch == '\n' {
					l.line++
					l.column = 0
				}
				return l.NextToken() // Recursively get the next token after the comment
			} else if next == '*' {
				// Block comment, read until */
				for {
					ch, _, err = l.readChar()
					if err != nil {
						break
					}

					if ch == '*' {
						ch, _, err = l.readChar()
						if err != nil {
							break
						}
						if ch == '/' {
							return l.NextToken() // Recursively get the next token after the comment
						}
					}

					if ch == '\n' {
						l.line++
						l.column = 0
					}
				}
			} else {
				l.unreadChar() // Put back the character that isn't part of a comment
			}
		}
		l.unreadChar() // Put back the '/' character
	}

	// Check for specific tokens
	switch ch {
	case '(':
		return Token{Type: TokenLeftParen, Value: "(", Line: l.line, Column: startColumn}, nil
	case ')':
		return Token{Type: TokenRightParen, Value: ")", Line: l.line, Column: startColumn}, nil
	case '[':
		return Token{Type: TokenLeftBracket, Value: "[", Line: l.line, Column: startColumn}, nil
	case ']':
		return Token{Type: TokenRightBracket, Value: "]", Line: l.line, Column: startColumn}, nil
	case '{':
		return Token{Type: TokenLeftBrace, Value: "{", Line: l.line, Column: startColumn}, nil
	case '}':
		return Token{Type: TokenRightBrace, Value: "}", Line: l.line, Column: startColumn}, nil
	case '|':
		return Token{Type: TokenAlternative, Value: "|", Line: l.line, Column: startColumn}, nil
	case '*':
		return Token{Type: TokenRepetition, Value: "*", Line: l.line, Column: startColumn}, nil
	case '+':
		return Token{Type: TokenOneOrMore, Value: "+", Line: l.line, Column: startColumn}, nil
	case '?':
		return Token{Type: TokenOptional, Value: "?", Line: l.line, Column: startColumn}, nil
	case ':':
		// Check for ::= assignment
		nextCh, _, err := l.readChar()
		if err == nil && nextCh == ':' {
			nextCh, _, err = l.readChar()
			if err == nil && nextCh == '=' {
				return Token{Type: TokenAssignment, Value: "::=", Line: l.line, Column: startColumn}, nil
			}
			l.unreadChar() // Put back the character after ':'
		}
		l.unreadChar() // Put back the first character after ':'
	case '"', '\'':
		// String literal
		quote := ch
		var value strings.Builder

		for {
			ch, _, err = l.readChar()
			if err != nil || ch == quote {
				break
			}

			// Handle escaped characters
			if ch == '\\' {
				escape, _, err := l.readChar()
				if err != nil {
					break
				}
				switch escape {
				case 'n':
					value.WriteRune('\n')
				case 't':
					value.WriteRune('\t')
				case 'r':
					value.WriteRune('\r')
				case '\\':
					value.WriteRune('\\')
				case quote:
					value.WriteRune(quote)
				default:
					value.WriteRune('\\')
					value.WriteRune(escape)
				}
			} else {
				value.WriteRune(ch)
			}
		}

		return Token{Type: TokenString, Value: value.String(), Line: l.line, Column: startColumn}, nil
	}

	// Identifier (rule name)
	if unicode.IsLetter(ch) || ch == '_' {
		var value strings.Builder
		value.WriteRune(ch)

		for {
			ch, _, err = l.readChar()
			if err != nil {
				break
			}

			if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
				value.WriteRune(ch)
			} else {
				l.unreadChar()
				break
			}
		}

		return Token{Type: TokenIdentifier, Value: value.String(), Line: l.line, Column: startColumn}, nil
	}

	// Unrecognized token
	return Token{Type: TokenInvalid, Value: string(ch), Line: l.line, Column: startColumn}, nil
}

// Helper method to read a character
func (l *Lexer) readChar() (rune, int, error) {
	ch, size, err := l.reader.ReadRune()
	if err != nil {
		return 0, 0, err
	}

	l.lastChar = ch
	l.column++

	return ch, size, nil
}

// Helper method to unread a character
func (l *Lexer) unreadChar() error {
	err := l.reader.UnreadRune()
	if err != nil {
		return err
	}

	l.column--
	return nil
}
