package codegen

import (
	"bytes"
	"os"
	"testing"

	"github.com/kvany/vtx1/assembler/internal/lexer"
	"github.com/kvany/vtx1/assembler/internal/parser"
)

func TestCodeGenerator(t *testing.T) {
	// Test cases for different instruction formats
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name: "Basic ALU R-format instruction",
			input: `
				ADD r1, r2, r3
			`,
			// Encoding for ADD r1, r2, r3
			// Type: OpTypeALU (000), Opcode: 000001
			// Format: | Opcode(6) | rs2(5) | rs1(5) | Type(3) | rd(5) | Opcode(7) |
			// Values: | 000001    | 00011  | 00010  | 000     | 00001 | 0000001   |
			expected: []byte{0x01, 0x09, 0x43, 0x04}, // Little-endian encoding
		},
		{
			name: "Basic MEM I-format instruction",
			input: `
				LD r1, 16(r2)
			`,
			// Encoding for LD r1, 16(r2)
			// Type: OpTypeMEM (001), Opcode: 000000
			// Format: | imm(12) | rs1(5) | Type(3) | rd(5) | Opcode(7) |
			// Values: | 000000010000 | 00010 | 001 | 00001 | 0000000 |
			expected: []byte{0x00, 0x12, 0x10, 0x10}, // Little-endian encoding
		},
		{
			name: "Jump J-format instruction",
			input: `
				start:
				JMP start
			`,
			// Encoding for JMP start (relative jump to same location, offset = 0)
			// Opcode: 0x0C (assuming from initial code), rd: 0
			// Format: | target(20) | rd(5) | opcode(7) |
			// Values: | 00000000000000000000 | 00000 | 0001100 |
			expected: []byte{0x0C, 0x00, 0x00, 0x00}, // Little-endian encoding
		},
		{
			name: "VLIW instruction with 2 operations",
			input: `
				{
					ADD r1, r2, r3;
					SUB r4, r5, r6
				}
			`,
			// Encoding for a VLIW with 2 operations (ADD and SUB)
			// Format bits: 10 (2 operations)
			// op1 = ADD encoded in 15 bits, op2 = SUB encoded in 15 bits
			// This is a complex encoding that depends on the actual VLIW implementation
			// For testing we're using a simplified expected value
			expected: []byte{0x02, 0x00, 0x00, 0x00}, // Little-endian with format bits 10
		},
		{
			name: "Data directives",
			input: `
				.word 0x12345678
				.byte 0xAA, 0xBB, 0xCC, 0xDD
			`,
			// One 32-bit word followed by four bytes
			expected: []byte{0x78, 0x56, 0x34, 0x12, 0xAA, 0xBB, 0xCC, 0xDD},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a lexer from the input
			lex := lexer.New(tc.input)

			// Parse the input
			p := parser.New(lex)
			ast, err := p.Parse()
			if err != nil {
				t.Fatalf("Failed to parse input: %v", err)
			}

			// Generate code
			cg := New(ast, BinaryFormatRaw)
			err = cg.Generate()
			if err != nil {
				t.Fatalf("Failed to generate code: %v", err)
			}

			// Check for errors
			if len(cg.Errors()) > 0 {
				t.Fatalf("Unexpected errors: %v", cg.Errors())
			}

			// Write output to a buffer
			var buf bytes.Buffer
			err = cg.WriteOutput(&buf)
			if err != nil {
				t.Fatalf("Failed to write output: %v", err)
			}

			// Compare with expected output
			if !bytes.Equal(buf.Bytes(), tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, buf.Bytes())
			}
		})
	}
}

func TestBinaryOutputFormats(t *testing.T) {
	// Simple test program
	const input = `
		start:
			ADD r1, r2, r3
			LD  r4, 8(r5)
			JMP start
	`

	// Create the AST
	lex := lexer.New(input)
	p := parser.New(lex)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse input: %v", err)
	}

	// Test each output format
	formats := []struct {
		name   string
		format BinaryFormat
	}{
		{"Raw", BinaryFormatRaw},
		{"Intel HEX", BinaryFormatHEX},
		{"ELF", BinaryFormatELF},
	}

	for _, fmt := range formats {
		t.Run(fmt.name, func(t *testing.T) {
			// Generate code with the specific format
			cg := New(ast, fmt.format)
			err = cg.Generate()
			if err != nil {
				t.Fatalf("Failed to generate code: %v", err)
			}

			// Check for errors
			if len(cg.Errors()) > 0 {
				t.Fatalf("Unexpected errors: %v", cg.Errors())
			}

			// Write output to a buffer
			var buf bytes.Buffer
			err = cg.WriteOutput(&buf)
			if err != nil {
				t.Fatalf("Failed to write output: %v", err)
			}

			// Just verify we got some output
			if buf.Len() == 0 {
				t.Errorf("No output generated for format %s", fmt.name)
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name: "Unknown instruction",
			input: `
				UNKNOWN r1, r2, r3
			`,
			expectedError: "unknown instruction",
		},
		{
			name: "Missing operands",
			input: `
				ADD r1
			`,
			expectedError: "not enough operands",
		},
		{
			name: "Invalid register",
			input: `
				ADD r1, r99, r3
			`,
			expectedError: "invalid register",
		},
		{
			name: "Undefined symbol",
			input: `
				JMP nonexistent_label
			`,
			expectedError: "undefined symbol",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a lexer from the input
			lex := lexer.New(tc.input)

			// Parse the input
			p := parser.New(lex)
			ast, err := p.Parse()
			if err != nil {
				t.Logf("Parser error (expected): %v", err)
				return
			}

			// Generate code
			cg := New(ast, BinaryFormatRaw)
			_ = cg.Generate() // Ignore the returned error, we check the internal errors

			// Check for expected error
			foundExpectedError := false
			for _, err := range cg.Errors() {
				if err.Error() != "" && bytes.Contains([]byte(err.Error()), []byte(tc.expectedError)) {
					foundExpectedError = true
					break
				}
			}

			if !foundExpectedError {
				t.Errorf("Expected error containing '%s', got errors: %v", tc.expectedError, cg.Errors())
			}
		})
	}
}

func TestRealAssemblyFile(t *testing.T) {
	// Create a temporary test assembly file
	tempFile, err := os.CreateTemp("", "testasm*.asm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}(tempFile.Name())

	// Write a simple but complete assembly program to the file
	const program = `
		; VTX1 Assembly Test Program
		
		.section .text
		
		start:
			; Initialize registers
			LD r1, 0(r0)     ; Load r1 from address 0
			LD r2, 4(r0)     ; Load r2 from address 4
			
		loop:
			ADD r3, r1, r2   ; r3 = r1 + r2
			ST r3, 8(r0)     ; Store r3 to address 8
			
			; VLIW instruction with parallel operations
			{
				ADD r1, r3, r0;  ; r1 = r3 (copy)
				ADD r2, r1, r0   ; r2 = r1 (copy)
			}
			
			CMP r3, r0       ; Compare r3 with zero
			JMP loop         ; Jump back to loop
			
		.section .data
			.word 0x12345678  ; Initial value for r1
			.word 0xABCDEF01  ; Initial value for r2
	`

	_, err = tempFile.Write([]byte(program))
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Load the file and process it
	sourceBytes, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	// Create a lexer from the source
	lex := lexer.New(string(sourceBytes))

	// Parse the input
	p := parser.New(lex)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	// Try all output formats
	formats := []BinaryFormat{
		BinaryFormatRaw,
		BinaryFormatHEX,
		BinaryFormatELF,
	}

	for _, format := range formats {
		// Generate code
		cg := New(ast, format)
		err = cg.Generate()
		if err != nil {
			t.Fatalf("Failed to generate code for format %d: %v", format, err)
		}

		// Check for errors
		if len(cg.Errors()) > 0 {
			t.Fatalf("Unexpected errors for format %d: %v", format, cg.Errors())
		}

		// Write output to a buffer
		var buf bytes.Buffer
		err = cg.WriteOutput(&buf)
		if err != nil {
			t.Fatalf("Failed to write output for format %d: %v", format, err)
		}

		// Just verify we got some output
		if buf.Len() == 0 {
			t.Errorf("No output generated for format %d", format)
		}
	}
}
