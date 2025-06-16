package grammar

import (
	"fmt"
	"strings"
)

// GrammarRule represents a single rule in the EBNF grammar
type GrammarRule struct {
	Name         string
	Definition   string
	Alternatives []string
}

// Grammar represents the entire EBNF grammar
type Grammar struct {
	Rules map[string]GrammarRule
}

// ParseEBNF parses an EBNF grammar file and returns a Grammar structure
func ParseEBNF(filename string) (*Grammar, error) {
	// Use the new token-based parser
	return ParseGrammar(filename)
}

// parseAlternatives splits a rule definition into alternatives
func parseAlternatives(definition string) []string {
	// This is enhanced to handle parentheses and quoted strings properly
	if !strings.Contains(definition, "|") {
		return []string{strings.TrimSpace(definition)}
	}

	// We need to handle cases where | appears inside parentheses or quotes
	var alternatives []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)
	parenLevel := 0

	for _, c := range definition {
		switch c {
		case '\'', '"':
			if inQuotes && c == quoteChar {
				inQuotes = false
			} else if !inQuotes {
				inQuotes = true
				quoteChar = c
			}
			current.WriteRune(c)
		case '(':
			if !inQuotes {
				parenLevel++
			}
			current.WriteRune(c)
		case ')':
			if !inQuotes && parenLevel > 0 {
				parenLevel--
			}
			current.WriteRune(c)
		case '|':
			if !inQuotes && parenLevel == 0 {
				// This is a top-level alternative separator
				alternatives = append(alternatives, strings.TrimSpace(current.String()))
				current.Reset()
			} else {
				current.WriteRune(c)
			}
		default:
			current.WriteRune(c)
		}
	}

	// Don't forget the last part
	if current.Len() > 0 {
		alternatives = append(alternatives, strings.TrimSpace(current.String()))
	}

	return alternatives
}

// GenerateTestCases generates test cases based on the grammar
func (g *Grammar) GenerateTestCases() []string {
	// Generate simple examples of each rule
	var testCases []string

	// Start with the Program rule as the entry point
	if programRule, ok := g.Rules["Program"]; ok {
		testCases = append(testCases, g.generateFromRule("Program", 0, make(map[string]bool))...)
	}

	return testCases
}

// generateFromRule recursively generates examples for a rule
func (g *Grammar) generateFromRule(ruleName string, depth int, visited map[string]bool) []string {
	// Prevent infinite recursion
	if depth > 3 || visited[ruleName] {
		return nil
	}
	visited[ruleName] = true
	defer func() { visited[ruleName] = false }()

	rule, ok := g.Rules[ruleName]
	if !ok {
		return nil
	}

	var examples []string

	// For simple rules, generate an example for each alternative
	for _, alt := range rule.Alternatives {
		if example := g.generateFromAlternative(alt, depth+1, visited); example != "" {
			examples = append(examples, example)
		}
		if len(examples) >= 3 {
			break // Limit to 3 examples per rule
		}
	}

	return examples
}

// generateFromAlternative generates an example for a single alternative
func (g *Grammar) generateFromAlternative(alt string, depth int, visited map[string]bool) string {
	var result strings.Builder

	// This implementation is enhanced to handle EBNF operators
	tokens := tokenizeForGeneration(alt)

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// Handle operators and modifiers
		if i < len(tokens)-1 {
			nextToken := tokens[i+1]
			if nextToken == "*" || nextToken == "+" || nextToken == "?" {
				// Skip this token and its modifier for now - simplified approach
				i++
				continue
			}
		}

		// If it's a reference to another rule
		if !strings.HasPrefix(token, "'") && !strings.HasPrefix(token, "\"") {
			if examples := g.generateFromRule(token, depth+1, visited); len(examples) > 0 {
				result.WriteString(examples[0] + " ")
			}
		} else {
			// It's a literal, strip quotes and add it
			token = strings.Trim(token, "'\"")
			result.WriteString(token + " ")
		}
	}

	return strings.TrimSpace(result.String())
}

// tokenizeForGeneration breaks an alternative into tokens for test generation
func tokenizeForGeneration(alt string) []string {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)

	for _, c := range alt {
		switch {
		case c == '\'' || c == '"':
			if inQuotes && c == quoteChar {
				inQuotes = false
				current.WriteRune(c)
				tokens = append(tokens, current.String())
				current.Reset()
			} else if !inQuotes {
				if current.Len() > 0 {
					tokens = append(tokens, strings.TrimSpace(current.String()))
					current.Reset()
				}
				inQuotes = true
				quoteChar = c
				current.WriteRune(c)
			} else {
				current.WriteRune(c)
			}
		case c == ' ' && !inQuotes:
			if current.Len() > 0 {
				tokens = append(tokens, strings.TrimSpace(current.String()))
				current.Reset()
			}
		case c == '(' || c == ')' || c == '[' || c == ']' || c == '{' || c == '}' || c == '*' || c == '+' || c == '?':
			if !inQuotes {
				if current.Len() > 0 {
					tokens = append(tokens, strings.TrimSpace(current.String()))
					current.Reset()
				}
				tokens = append(tokens, string(c))
			} else {
				current.WriteRune(c)
			}
		default:
			current.WriteRune(c)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, strings.TrimSpace(current.String()))
	}

	return tokens
}

// ValidateSyntax checks if a given source code conforms to the grammar
func (g *Grammar) ValidateSyntax(source string) (bool, error) {
	// This would require implementing a full parser based on the grammar
	// For now, we'll return a placeholder
	return true, fmt.Errorf("syntax validation not yet implemented")
}

// ListInstructions returns a list of all instruction types defined in the grammar
func (g *Grammar) ListInstructions() []string {
	var instructions []string

	// Look for the Mnemonic rule and its components
	if mnemonicRule, ok := g.Rules["Mnemonic"]; ok {
		for _, alt := range mnemonicRule.Alternatives {
			if strings.Contains(alt, "ALU_Op") {
				if aluOpRule, ok := g.Rules["ALU_Op"]; ok {
					for _, op := range parseInstructionList(aluOpRule.Definition) {
						instructions = append(instructions, op)
					}
				}
			}
			if strings.Contains(alt, "Memory_Op") {
				if memOpRule, ok := g.Rules["Memory_Op"]; ok {
					for _, op := range parseInstructionList(memOpRule.Definition) {
						instructions = append(instructions, op)
					}
				}
			}
			if strings.Contains(alt, "Control_Op") {
				if ctrlOpRule, ok := g.Rules["Control_Op"]; ok {
					for _, op := range parseInstructionList(ctrlOpRule.Definition) {
						instructions = append(instructions, op)
					}
				}
			}
		}
	}

	return instructions
}

// Helper to parse a list of instructions from a rule like "ADD | SUB | MUL"
func parseInstructionList(definition string) []string {
	var instructions []string

	// Extract literals by finding quoted strings
	inQuote := false
	quoteChar := rune(0)
	var currentInst strings.Builder

	for _, c := range definition {
		switch c {
		case '\'', '"':
			if inQuote && c == quoteChar {
				inQuote = false
				instructions = append(instructions, strings.TrimSpace(currentInst.String()))
				currentInst.Reset()
			} else if !inQuote {
				inQuote = true
				quoteChar = c
			} else {
				currentInst.WriteRune(c)
			}
		case '|':
			if !inQuote && currentInst.Len() > 0 {
				instructions = append(instructions, strings.TrimSpace(currentInst.String()))
				currentInst.Reset()
			} else if inQuote {
				currentInst.WriteRune(c)
			}
		default:
			if inQuote {
				currentInst.WriteRune(c)
			}
		}
	}

	if currentInst.Len() > 0 && inQuote {
		instructions = append(instructions, strings.TrimSpace(currentInst.String()))
	}

	return instructions
}
