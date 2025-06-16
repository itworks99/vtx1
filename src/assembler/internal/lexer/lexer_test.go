package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
; This is a test for the VTX1 lexer
.ORG 0x1000    ; Set origin to 0x1000

main:           ; Main program entry
        LD T0, 0x1234        ; Load immediate value
        LD T1, [T0]          ; Load from memory
        ADD T2, T0, T1       ; Add registers
        
        ; VLIW instruction (parallel execution)
        [ADD T3, T0, T1] [SUB T4, T0, T1] [MUL T5, T0, T1]
        
        ; Test balanced ternary
        LD T0, %+-0          ; Load ternary value
        
        ; Test string
        DB "Hello, VTX1!"    ; Store string
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{NEWLINE, "\\n"},
		{COMMENT, "; This is a test for the VTX1 lexer"},
		{NEWLINE, "\\n"},
		{DIR_ORG, ".ORG"},
		{HEXADECIMAL, "0x1000"},
		{COMMENT, "; Set origin to 0x1000"},
		{NEWLINE, "\\n"},
		{NEWLINE, "\\n"},
		{IDENTIFIER, "main"},
		{COLON, ":"},
		{COMMENT, "; Main program entry"},
		{NEWLINE, "\\n"},
		{OP_LD, "LD"},
		{GPR, "T0"},
		{COMMA, ","},
		{HEXADECIMAL, "0x1234"},
		{COMMENT, "; Load immediate value"},
		{NEWLINE, "\\n"},
		{OP_LD, "LD"},
		{GPR, "T1"},
		{COMMA, ","},
		{LSQUARE, "["},
		{GPR, "T0"},
		{RSQUARE, "]"},
		{COMMENT, "; Load from memory"},
		{NEWLINE, "\\n"},
		{OP_ADD, "ADD"},
		{GPR, "T2"},
		{COMMA, ","},
		{GPR, "T0"},
		{COMMA, ","},
		{GPR, "T1"},
		{COMMENT, "; Add registers"},
		{NEWLINE, "\\n"},
		{NEWLINE, "\\n"},
		{COMMENT, "; VLIW instruction (parallel execution)"},
		{NEWLINE, "\\n"},
		{LSQUARE, "["},
		{OP_ADD, "ADD"},
		{GPR, "T3"},
		{COMMA, ","},
		{GPR, "T0"},
		{COMMA, ","},
		{GPR, "T1"},
		{RSQUARE, "]"},
		{LSQUARE, "["},
		{OP_SUB, "SUB"},
		{GPR, "T4"},
		{COMMA, ","},
		{GPR, "T0"},
		{COMMA, ","},
		{GPR, "T1"},
		{RSQUARE, "]"},
		{LSQUARE, "["},
		{OP_MUL, "MUL"},
		{GPR, "T5"},
		{COMMA, ","},
		{GPR, "T0"},
		{COMMA, ","},
		{GPR, "T1"},
		{RSQUARE, "]"},
		{NEWLINE, "\\n"},
		{NEWLINE, "\\n"},
		{COMMENT, "; Test balanced ternary"},
		{NEWLINE, "\\n"},
		{OP_LD, "LD"},
		{GPR, "T0"},
		{COMMA, ","},
		{TERNARY, "%+-0"},
		{COMMENT, "; Load ternary value"},
		{NEWLINE, "\\n"},
		{NEWLINE, "\\n"},
		{COMMENT, "; Test string"},
		{NEWLINE, "\\n"},
		{DIR_DB, ".DB"},
		{STRING, "Hello, VTX1!"},
		{COMMENT, "; Store string"},
		{NEWLINE, "\\n"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%d, got=%d (%q)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIndentationHandling(t *testing.T) {
	input := `
label:
        ADD T0, T1, T2  ; Indented instruction
	LD T3, 0x1234   ; Tab-indented instruction
    ST T4, [T5]     ; Different indentation
`

	l := New(input)

	// Skip the newline
	_ = l.NextToken()

	// Check for label token
	tok := l.NextToken()
	if tok.Type != IDENTIFIER || tok.Literal != "label" {
		t.Fatalf("Expected label identifier, got %q", tok.Literal)
	}

	// Check for colon
	tok = l.NextToken()
	if tok.Type != COLON {
		t.Fatalf("Expected colon, got %q", tok.Literal)
	}

	// Skip newline
	_ = l.NextToken()

	// Check that we correctly handle the indented instruction
	tok = l.NextToken()
	if tok.Type != OP_ADD {
		t.Fatalf("Expected ADD instruction, got token type %d", tok.Type)
	}

	// Advance tokens until the tab-indented instruction
	for i := 0; i < 8; i++ { // Skip to the next line
		_ = l.NextToken()
	}

	// Check the tab-indented instruction
	tok = l.NextToken()
	if tok.Type != OP_LD {
		t.Fatalf("Expected LD instruction, got token type %d", tok.Type)
	}
}

func TestBalancedTernary(t *testing.T) {
	input := `
	%+0- ; Simple ternary
	%+-+-+0-0-+  ; Complex ternary
	%+           ; Single trit
	`

	l := New(input)

	// Skip the newline and whitespace (implicit in NextToken for the first token)
	_ = l.NextToken()

	// Check first ternary
	tok := l.NextToken()
	if tok.Type != TERNARY || tok.Literal != "%+0-" {
		t.Fatalf("Expected ternary literal %%+0-, got %q", tok.Literal)
	}

	// Skip comment and newline
	_ = l.NextToken() // comment
	_ = l.NextToken() // newline

	// Check second ternary
	tok = l.NextToken()
	if tok.Type != TERNARY || tok.Literal != "%+-+-+0-0-+" {
		t.Fatalf("Expected ternary literal %%+-+-+0-0-+, got %q", tok.Literal)
	}

	// Skip comment and newline
	_ = l.NextToken() // comment
	_ = l.NextToken() // newline

	// Check third ternary
	tok = l.NextToken()
	if tok.Type != TERNARY || tok.Literal != "%+" {
		t.Fatalf("Expected ternary literal %%+, got %q", tok.Literal)
	}
}

func TestVLIWInstructions(t *testing.T) {
	input := `
	[ADD T0, T1, T2] [SUB T3, T4, T5] [MUL T6, T0, T1]
	[LD T0, 0x1000]  ; Single operation in VLIW
	`

	l := New(input)

	// Skip newline
	_ = l.NextToken()

	// First VLIW instruction - check the brackets and operations
	tok := l.NextToken() // '['
	if tok.Type != LSQUARE {
		t.Fatalf("Expected LSQUARE, got %d", tok.Type)
	}

	tok = l.NextToken() // 'ADD'
	if tok.Type != OP_ADD {
		t.Fatalf("Expected OP_ADD, got %d", tok.Type)
	}

	// Skip ahead to check the next VLIW operation
	for i := 0; i < 7; i++ {
		_ = l.NextToken()
	}

	tok = l.NextToken() // '['
	if tok.Type != LSQUARE {
		t.Fatalf("Expected LSQUARE, got %d", tok.Type)
	}

	tok = l.NextToken() // 'SUB'
	if tok.Type != OP_SUB {
		t.Fatalf("Expected OP_SUB, got %d", tok.Type)
	}
}
