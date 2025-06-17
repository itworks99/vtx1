package grammar

import (
	"fmt"
	"os"
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
	content, err := ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading grammar file: %w", err)
	}

	grammar := &Grammar{
		Rules: make(map[string]GrammarRule),
	}

	lines := strings.Split(content, "\n")
	var currentRule *GrammarRule
	var ruleDefinitionBuilder strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue // Skip empty lines and comments
		}

		if strings.Contains(line, "::=") {
			// Save the previous rule if there was one
			if currentRule != nil {
				definition := strings.TrimSpace(ruleDefinitionBuilder.String())
				currentRule.Definition = definition
				currentRule.Alternatives = parseAlternatives(definition)
				grammar.Rules[currentRule.Name] = *currentRule
			}

			// Start a new rule
			parts := strings.SplitN(line, "::=", 2)
			if len(parts) != 2 {
				continue
			}

			ruleName := strings.TrimSpace(parts[0])
			ruleDefinitionBuilder.Reset()
			ruleDefinitionBuilder.WriteString(strings.TrimSpace(parts[1]))

			currentRule = &GrammarRule{
				Name: ruleName,
			}
		} else if currentRule != nil {
			// Continue with the current rule definition
			ruleDefinitionBuilder.WriteString(" ")
			ruleDefinitionBuilder.WriteString(line)
		}
	}

	// Add the last rule
	if currentRule != nil {
		definition := strings.TrimSpace(ruleDefinitionBuilder.String())
		currentRule.Definition = definition
		currentRule.Alternatives = parseAlternatives(definition)
		grammar.Rules[currentRule.Name] = *currentRule
	}

	return grammar, nil
}

// tokenizeForGeneration breaks an alternative into tokens for test generation
func tokenizeForGeneration(alt string) []string {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)

	for _, char := range alt {
		if inQuotes {
			current.WriteRune(char)
			if char == quoteChar {
				inQuotes = false
				tokens = append(tokens, current.String())
				current.Reset()
			}
		} else if char == '\'' || char == '"' {
			inQuotes = true
			quoteChar = char
			current.WriteRune(char)
		} else if char == ' ' || char == '\t' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		} else if char == '*' || char == '+' || char == '?' || char == '|' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(char))
		} else {
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

// parseAlternatives splits a rule definition into its alternative productions
func parseAlternatives(definition string) []string {
	var alternatives []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)
	inGroup := 0

	for _, char := range definition {
		if inQuotes {
			current.WriteRune(char)
			if char == quoteChar {
				inQuotes = false
			}
		} else if char == '\'' || char == '"' {
			inQuotes = true
			quoteChar = char
			current.WriteRune(char)
		} else if char == '(' {
			inGroup++
			current.WriteRune(char)
		} else if char == ')' {
			inGroup--
			current.WriteRune(char)
		} else if char == '|' && inGroup == 0 {
			alternatives = append(alternatives, strings.TrimSpace(current.String()))
			current.Reset()
		} else {
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		alternatives = append(alternatives, strings.TrimSpace(current.String()))
	}

	return alternatives
}

// ReadFile reads the content of a file and returns it as a string
func ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GenerateTestCases generates test cases based on the grammar
func (g *Grammar) GenerateTestCases() []string {
	var testCases []string

	// Start with the Program rule as the entry point
	if _, ok := g.Rules["Program"]; ok {
		testCases = append(testCases, g.generateFromRule("Program", 0, make(map[string]bool))...)
	}

	return testCases
}

// generateFromRule recursively generates examples for a rule
func (g *Grammar) generateFromRule(ruleName string, depth int, visited map[string]bool) []string {
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
	for _, alt := range rule.Alternatives {
		if example := g.generateFromAlternative(alt, depth+1, visited); example != "" {
			examples = append(examples, example)
		}
		if len(examples) >= 3 {
			break
		}
	}

	return examples
}

// generateFromAlternative generates an example for a single alternative
func (g *Grammar) generateFromAlternative(alt string, depth int, visited map[string]bool) string {
	var result strings.Builder
	tokens := tokenizeForGeneration(alt)

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if i < len(tokens)-1 {
			nextToken := tokens[i+1]
			if nextToken == "*" || nextToken == "+" || nextToken == "?" {
				i++
				if examples := g.generateFromRule(token, depth+1, visited); len(examples) > 0 {
					result.WriteString(examples[0] + " ")
				}
				continue
			}
		}

		if !strings.HasPrefix(token, "'") && !strings.HasPrefix(token, "\"") {
			if examples := g.generateFromRule(token, depth+1, visited); len(examples) > 0 {
				result.WriteString(examples[0] + " ")
			}
		} else {
			token = strings.Trim(token, "'\"")
			result.WriteString(token + " ")
		}
	}

	return strings.TrimSpace(result.String())
}

// ExpandEBNFGrammar expands any rules in the grammar to handle special notation
func (g *Grammar) ExpandEBNFGrammar() {
	expanded := make(map[string]GrammarRule)

	// Copy existing rules first
	for name, rule := range g.Rules {
		expanded[name] = rule
	}

	// Expand rules with special notation (*, +, ?)
	for name, rule := range g.Rules {
		for i, alt := range rule.Alternatives {
			newAlt, newRules := g.expandAlternative(alt, name)
			rule.Alternatives[i] = newAlt

			// Add newly created rules
			for newName, newRule := range newRules {
				expanded[newName] = newRule
			}
		}

		expanded[name] = rule
	}

	g.Rules = expanded
}

// expandAlternative expands an alternative to handle repetition operators (*, +, ?)
func (g *Grammar) expandAlternative(alt string, parentName string) (string, map[string]GrammarRule) {
	newRules := make(map[string]GrammarRule)
	tokens := tokenizeForGeneration(alt)
	var result strings.Builder

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if i < len(tokens)-1 {
			nextToken := tokens[i+1]

			if nextToken == "*" {
				// Zero or more repetitions
				newRuleName := fmt.Sprintf("%s_%s_Star", parentName, token)
				newRules[newRuleName] = GrammarRule{
					Name:         newRuleName,
					Definition:   fmt.Sprintf("| %s %s", token, newRuleName),
					Alternatives: []string{"", fmt.Sprintf("%s %s", token, newRuleName)},
				}
				result.WriteString(newRuleName + " ")
				i++ // Skip the * token
				continue
			} else if nextToken == "+" {
				// One or more repetitions
				newRuleName := fmt.Sprintf("%s_%s_Plus", parentName, token)
				newRules[newRuleName] = GrammarRule{
					Name:         newRuleName,
					Definition:   fmt.Sprintf("%s | %s %s", token, token, newRuleName),
					Alternatives: []string{token, fmt.Sprintf("%s %s", token, newRuleName)},
				}
				result.WriteString(newRuleName + " ")
				i++ // Skip the + token
				continue
			} else if nextToken == "?" {
				// Optional occurrence
				newRuleName := fmt.Sprintf("%s_%s_Optional", parentName, token)
				newRules[newRuleName] = GrammarRule{
					Name:         newRuleName,
					Definition:   fmt.Sprintf("| %s", token),
					Alternatives: []string{"", token},
				}
				result.WriteString(newRuleName + " ")
				i++ // Skip the ? token
				continue
			}
		}

		result.WriteString(token + " ")
	}

	return strings.TrimSpace(result.String()), newRules
}

// PrintGrammar prints the grammar rules in a readable format
func (g *Grammar) PrintGrammar() string {
	var result strings.Builder

	for name, rule := range g.Rules {
		result.WriteString(name)
		result.WriteString(" ::= ")
		result.WriteString(rule.Definition)
		result.WriteString("\n")
	}

	return result.String()
}

// ValidateSyntax checks if a string matches the grammar
func (g *Grammar) ValidateSyntax(input string, startRule string) bool {
	// Placeholder for a real validation implementation
	// A proper validation would require a parsing algorithm like LL(1) or recursive descent
	return true
}
