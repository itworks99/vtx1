package grammar

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Parser is responsible for parsing an EBNF grammar from tokens
type Parser struct {
	lexer        *Lexer
	current      Token
	errors       []string
	lookahead    Token
	hasLookahead bool
}

// NewParser creates a new parser from an io.Reader
func NewParser(r io.Reader) *Parser {
	return &Parser{
		lexer:  NewLexer(r),
		errors: []string{},
	}
}

// ParseGrammar parses a full EBNF grammar
func ParseGrammar(filename string) (*Grammar, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open grammar file: %w", err)
	}
	defer file.Close()

	parser := NewParser(file)
	return parser.Parse()
}

// Parse parses the grammar from the input
func (p *Parser) Parse() (*Grammar, error) {
	grammar := &Grammar{
		Rules: make(map[string]GrammarRule),
	}

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	// Parse each rule in the grammar
	for p.current.Type != TokenEOF {
		rule, err := p.parseRule()
		if err != nil {
			return nil, err
		}
		grammar.Rules[rule.Name] = *rule

		// Check for errors
		if len(p.errors) > 0 {
			return nil, fmt.Errorf("parsing errors: %s", strings.Join(p.errors, "; "))
		}
	}

	return grammar, nil
}

// parseRule parses a single EBNF rule
func (p *Parser) parseRule() (*GrammarRule, error) {
	// A rule should start with an identifier
	if p.current.Type != TokenIdentifier {
		p.addError(fmt.Sprintf("expected rule name, got %s", p.current.Value))
		return nil, fmt.Errorf("expected rule name at line %d, column %d",
			p.current.Line, p.current.Column)
	}

	rule := &GrammarRule{
		Name: p.current.Value,
	}

	// Next should be the assignment operator ::=
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	if p.current.Type != TokenAssignment {
		p.addError(fmt.Sprintf("expected ::=, got %s", p.current.Value))
		return nil, fmt.Errorf("expected ::= at line %d, column %d",
			p.current.Line, p.current.Column)
	}

	// Next comes the rule expression
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	rule.Definition = expr
	rule.Alternatives = strings.Split(expr, "|")

	// Trim spaces from alternatives
	for i, alt := range rule.Alternatives {
		rule.Alternatives[i] = strings.TrimSpace(alt)
	}

	return rule, nil
}

// parseExpression parses an EBNF expression
func (p *Parser) parseExpression() (string, error) {
	var result strings.Builder

	// Parse the first term
	term, err := p.parseTerm()
	if err != nil {
		return "", err
	}
	result.WriteString(term)

	// Continue parsing terms separated by |
	for {
		if err := p.peek(); err != nil {
			return "", err
		}

		if p.lookahead.Type == TokenAlternative {
			if err := p.nextToken(); err != nil { // consume the |
				return "", err
			}

			if err := p.nextToken(); err != nil { // move to the next term
				return "", err
			}

			term, err := p.parseTerm()
			if err != nil {
				return "", err
			}

			result.WriteString(" | " + term)
		} else {
			break
		}
	}

	return result.String(), nil
}

// parseTerm parses a term in an EBNF expression
func (p *Parser) parseTerm() (string, error) {
	var result strings.Builder

	// Parse token based on its type
	switch p.current.Type {
	case TokenIdentifier:
		// Non-terminal symbol
		result.WriteString(p.current.Value)

	case TokenString:
		// Terminal symbol
		result.WriteString("'" + p.current.Value + "'")

	case TokenLeftParen:
		// Group ( expr )
		if err := p.nextToken(); err != nil {
			return "", err
		}

		expr, err := p.parseExpression()
		if err != nil {
			return "", err
		}

		if p.current.Type != TokenRightParen {
			p.addError(fmt.Sprintf("expected ), got %s", p.current.Value))
			return "", fmt.Errorf("expected ) at line %d, column %d",
				p.current.Line, p.current.Column)
		}

		result.WriteString("(" + expr + ")")

	case TokenLeftBracket:
		// Optional [ expr ]
		if err := p.nextToken(); err != nil {
			return "", err
		}

		expr, err := p.parseExpression()
		if err != nil {
			return "", err
		}

		if p.current.Type != TokenRightBracket {
			p.addError(fmt.Sprintf("expected ], got %s", p.current.Value))
			return "", fmt.Errorf("expected ] at line %d, column %d",
				p.current.Line, p.current.Column)
		}

		result.WriteString("[" + expr + "]")

	case TokenLeftBrace:
		// Repetition { expr }
		if err := p.nextToken(); err != nil {
			return "", err
		}

		expr, err := p.parseExpression()
		if err != nil {
			return "", err
		}

		if p.current.Type != TokenRightBrace {
			p.addError(fmt.Sprintf("expected }, got %s", p.current.Value))
			return "", fmt.Errorf("expected } at line %d, column %d",
				p.current.Line, p.current.Column)
		}

		result.WriteString("{" + expr + "}")

	default:
		p.addError(fmt.Sprintf("unexpected token: %s", p.current.Value))
		return "", fmt.Errorf("unexpected token at line %d, column %d: %s",
			p.current.Line, p.current.Column, p.current.Value)
	}

	// Check for modifiers like * (zero or more), + (one or more), ? (optional)
	if err := p.peek(); err == nil {
		switch p.lookahead.Type {
		case TokenRepetition, TokenOneOrMore, TokenOptional:
			if err := p.nextToken(); err != nil {
				return "", err
			}
			result.WriteString(p.current.Value)
		}
	}

	// Move to the next token for the next call
	if err := p.nextToken(); err != nil {
		if err != io.EOF {
			return "", err
		}
	}

	return result.String(), nil
}

// nextToken advances to the next token
func (p *Parser) nextToken() error {
	var err error
	if p.hasLookahead {
		p.current = p.lookahead
		p.hasLookahead = false
	} else {
		p.current, err = p.lexer.NextToken()
	}
	return err
}

// peek looks ahead at the next token without consuming it
func (p *Parser) peek() error {
	var err error
	if !p.hasLookahead {
		p.lookahead, err = p.lexer.NextToken()
		if err != nil {
			return err
		}
		p.hasLookahead = true
	}
	return nil
}

// addError adds a parsing error
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("line %d, column %d: %s",
		p.current.Line, p.current.Column, msg))
}
