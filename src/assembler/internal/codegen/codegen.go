// Package codegen implements code generation for the VTX1 assembly language.
// It transforms an AST from the parser into machine code.
package codegen

import (
	"encoding/binary"
	"fmt"
	"github.com/kvany/vtx1/assembler/internal/lexer"
	"io"
	"strconv"
	"strings"

	"github.com/kvany/vtx1/assembler/internal/parser"
)

// InstructionFormat defines the format of a VTX1 instruction
type InstructionFormat uint8

// Instruction formats
const (
	FormatR InstructionFormat = iota // Register format
	FormatI                          // Immediate format
	FormatJ                          // Jump format
	FormatV                          // VLIW format
)

// OperationType defines the VTX1 instruction operation type (3 bits)
type OperationType uint8

// VTX1 operation types from architecture specification
const (
	OpTypeALU    OperationType = 0b000
	OpTypeMEM    OperationType = 0b001
	OpTypeCTRL   OperationType = 0b010
	OpTypeVEC    OperationType = 0b011
	OpTypeFPU    OperationType = 0b100
	OpTypeSYS    OperationType = 0b101
	OpTypeMICRO  OperationType = 0b110
	OpTypeRESERV OperationType = 0b111 // Reserved
)

// InstructionEncoding contains the encoding info for an instruction
type InstructionEncoding struct {
	Mnemonic    string            // Instruction mnemonic (e.g., "ADD")
	Format      InstructionFormat // Instruction format (R, I, J, V)
	Type        OperationType     // Operation type (3 bits)
	Opcode      uint8             // Operation code (6 bits)
	CycleCount  uint8             // Number of cycles to execute
	Description string            // Human-readable description
}

// Encoding tables for all VTX1 instructions based on architecture specification
var encodingTable = map[string]InstructionEncoding{
	// ALU operations (Type: 000)
	"NEG":  {Mnemonic: "NEG", Format: FormatR, Type: OpTypeALU, Opcode: 0b000000, CycleCount: 1, Description: "Negate register"},
	"ADD":  {Mnemonic: "ADD", Format: FormatR, Type: OpTypeALU, Opcode: 0b000001, CycleCount: 1, Description: "Add two registers"},
	"SUB":  {Mnemonic: "SUB", Format: FormatR, Type: OpTypeALU, Opcode: 0b000010, CycleCount: 1, Description: "Subtract two registers"},
	"MUL":  {Mnemonic: "MUL", Format: FormatR, Type: OpTypeALU, Opcode: 0b000011, CycleCount: 2, Description: "Multiply two registers"},
	"AND":  {Mnemonic: "AND", Format: FormatR, Type: OpTypeALU, Opcode: 0b000100, CycleCount: 1, Description: "Bitwise AND"},
	"OR":   {Mnemonic: "OR", Format: FormatR, Type: OpTypeALU, Opcode: 0b000101, CycleCount: 1, Description: "Bitwise OR"},
	"NOT":  {Mnemonic: "NOT", Format: FormatR, Type: OpTypeALU, Opcode: 0b000110, CycleCount: 1, Description: "Bitwise NOT"},
	"XOR":  {Mnemonic: "XOR", Format: FormatR, Type: OpTypeALU, Opcode: 0b000111, CycleCount: 1, Description: "Bitwise XOR"},
	"SHL":  {Mnemonic: "SHL", Format: FormatR, Type: OpTypeALU, Opcode: 0b001000, CycleCount: 1, Description: "Shift left logical"},
	"SHR":  {Mnemonic: "SHR", Format: FormatR, Type: OpTypeALU, Opcode: 0b001001, CycleCount: 1, Description: "Shift right logical"},
	"ROL":  {Mnemonic: "ROL", Format: FormatR, Type: OpTypeALU, Opcode: 0b001010, CycleCount: 1, Description: "Rotate left"},
	"ROR":  {Mnemonic: "ROR", Format: FormatR, Type: OpTypeALU, Opcode: 0b001011, CycleCount: 1, Description: "Rotate right"},
	"CMP":  {Mnemonic: "CMP", Format: FormatR, Type: OpTypeALU, Opcode: 0b001100, CycleCount: 1, Description: "Compare registers"},
	"TEST": {Mnemonic: "TEST", Format: FormatR, Type: OpTypeALU, Opcode: 0b001101, CycleCount: 1, Description: "Test register"},
	"INC":  {Mnemonic: "INC", Format: FormatR, Type: OpTypeALU, Opcode: 0b001110, CycleCount: 1, Description: "Increment register"},
	"DEC":  {Mnemonic: "DEC", Format: FormatR, Type: OpTypeALU, Opcode: 0b001111, CycleCount: 1, Description: "Decrement register"},

	// Memory operations (Type: 001)
	"LD":   {Mnemonic: "LD", Format: FormatI, Type: OpTypeMEM, Opcode: 0b000000, CycleCount: 2, Description: "Load from memory"},
	"ST":   {Mnemonic: "ST", Format: FormatI, Type: OpTypeMEM, Opcode: 0b000001, CycleCount: 2, Description: "Store to memory"},
	"VLD":  {Mnemonic: "VLD", Format: FormatI, Type: OpTypeMEM, Opcode: 0b000010, CycleCount: 3, Description: "Vector load from memory"},
	"VST":  {Mnemonic: "VST", Format: FormatI, Type: OpTypeMEM, Opcode: 0b000011, CycleCount: 3, Description: "Vector store to memory"},
	"FLD":  {Mnemonic: "FLD", Format: FormatI, Type: OpTypeMEM, Opcode: 0b000100, CycleCount: 2, Description: "Floating-point load"},
	"FST":  {Mnemonic: "FST", Format: FormatI, Type: OpTypeMEM, Opcode: 0b000101, CycleCount: 2, Description: "Floating-point store"},
	"LEA":  {Mnemonic: "LEA", Format: FormatI, Type: OpTypeMEM, Opcode: 0b000110, CycleCount: 1, Description: "Load effective address"},
	"PUSH": {Mnemonic: "PUSH", Format: FormatI, Type: OpTypeMEM, Opcode: 0b000111, CycleCount: 2, Description: "Push to stack"},
}

// VLIWInstruction represents a VLIW instruction containing multiple operations
type VLIWInstruction struct {
	Operations []uint32 // Array of encoded operations
	LineInfo   string   // Source line information for error reporting
}

// BinaryFormat specifies the output format for generated code
type BinaryFormat int

// Binary output formats
const (
	BinaryFormatRaw BinaryFormat = iota // Raw binary output
	BinaryFormatELF                     // ELF executable format
	BinaryFormatHEX                     // Intel HEX format
)

// CodeGenerator handles the generation of machine code from an AST
type CodeGenerator struct {
	ast          *parser.AST      // The abstract syntax tree to process
	machineCode  []byte           // Generated machine code
	symbolTable  map[string]int   // Symbol table for labels
	currentAddr  int              // Current address during code generation
	errors       []CodegenError   // Errors encountered during code generation
	sectionStart int              // Start address of the current section
	outputFormat BinaryFormat     // Output format for the generated code
	vliwBuffer   *VLIWInstruction // Buffer for building VLIW instructions
}

// CodegenError represents a code generation error
type CodegenError struct {
	Message string
	Line    int
	Column  int
}

// Error returns the string representation of the error
func (e CodegenError) Error() string {
	return fmt.Sprintf("codegen error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// New creates a new code generator for the given AST with optional output format
func New(ast *parser.AST, format BinaryFormat) *CodeGenerator {
	return &CodeGenerator{
		ast:          ast,
		machineCode:  make([]byte, 0),
		symbolTable:  make(map[string]int),
		errors:       make([]CodegenError, 0),
		outputFormat: format,
	}
}

// Generate processes the AST and generates machine code
func (cg *CodeGenerator) Generate() error {
	// First pass: collect all symbols and their addresses
	cg.collectSymbols()

	// Reset address counter for second pass
	cg.currentAddr = 0

	// Second pass: generate machine code
	return cg.generateCode(cg.ast)
}

// addError adds a code generation error with line and column information
func (cg *CodeGenerator) addError(msg string, line, column int) error {
	cg.errors = append(cg.errors, CodegenError{
		Message: msg,
		Line:    line,
		Column:  column,
	})
	return fmt.Errorf(msg) // Return an error for convenience
}

// WriteOutput writes the generated machine code to the given writer in the specified format
func (cg *CodeGenerator) WriteOutput(w io.Writer) error {
	switch cg.outputFormat {
	case BinaryFormatRaw:
		_, err := w.Write(cg.machineCode)
		return err
	case BinaryFormatHEX:
		return cg.writeIntelHex(w)
	case BinaryFormatELF:
		return cg.writeELF(w)
	default:
		return fmt.Errorf("unsupported binary format: %d", cg.outputFormat)
	}
}

// writeIntelHex writes the machine code in Intel HEX format
func (cg *CodeGenerator) writeIntelHex(w io.Writer) error {
	const bytesPerRecord = 16

	for addr := 0; addr < len(cg.machineCode); addr += bytesPerRecord {
		// Determine how many bytes to write in this record
		remainingBytes := len(cg.machineCode) - addr
		bytesToWrite := bytesPerRecord
		if remainingBytes < bytesPerRecord {
			bytesToWrite = remainingBytes
		}

		// Start of record
		recordType := byte(0x00) // Data record

		// Calculate checksum (sum of all bytes + record type + length + address)
		var checksum byte = byte(bytesToWrite) + byte(addr>>8) + byte(addr&0xFF) + recordType

		// Write record start
		fmt.Fprintf(w, ":%02X%04X%02X", bytesToWrite, addr, recordType)

		// Write data bytes
		for i := 0; i < bytesToWrite; i++ {
			dataByte := cg.machineCode[addr+i]
			fmt.Fprintf(w, "%02X", dataByte)
			checksum += dataByte
		}

		// Write checksum (two's complement of the LSB of the sum)
		checksum = byte(0x100 - int(checksum))
		fmt.Fprintf(w, "%02X\n", checksum)
	}

	// End of file record
	fmt.Fprintln(w, ":00000001FF")

	return nil
}

// writeELF writes the machine code in ELF format (simplified)
func (cg *CodeGenerator) writeELF(w io.Writer) error {
	// This is a simplified ELF writer - for a production implementation
	// you would want to use a dedicated ELF library

	// ELF header constants
	const (
		ELFCLASS32    = 1      // 32-bit architecture
		ELFDATA2LSB   = 1      // Little endian
		EV_CURRENT    = 1      // Current version
		ELFOSABI_SYSV = 0      // System V ABI
		ET_EXEC       = 2      // Executable file
		EM_CUSTOM     = 0xBEEF // Custom machine type for VTX1
		EI_NIDENT     = 16     // Size of e_ident[]
	)

	// Write ELF Header
	// Magic number
	_, err := w.Write([]byte{0x7F, 'E', 'L', 'F'})
	if err != nil {
		return err
	}

	// File class (32-bit)
	_, err = w.Write([]byte{ELFCLASS32})
	if err != nil {
		return err
	}

	// Data encoding (little endian)
	_, err = w.Write([]byte{ELFDATA2LSB})
	if err != nil {
		return err
	}

	// File version
	_, err = w.Write([]byte{EV_CURRENT})
	if err != nil {
		return err
	}

	// OS ABI
	_, err = w.Write([]byte{ELFOSABI_SYSV})
	if err != nil {
		return err
	}

	// ABI Version
	_, err = w.Write([]byte{0})
	if err != nil {
		return err
	}

	// Padding
	_, err = w.Write(make([]byte, EI_NIDENT-8))
	if err != nil {
		return err
	}

	// File type (ET_EXEC)
	err = binary.Write(w, binary.LittleEndian, uint16(ET_EXEC))
	if err != nil {
		return err
	}

	// Machine type (custom VTX1)
	err = binary.Write(w, binary.LittleEndian, uint16(EM_CUSTOM))
	if err != nil {
		return err
	}

	// ELF version
	err = binary.Write(w, binary.LittleEndian, uint32(EV_CURRENT))
	if err != nil {
		return err
	}

	// Entry point address (start at 0 for simplicity)
	err = binary.Write(w, binary.LittleEndian, uint32(0))
	if err != nil {
		return err
	}

	// Program header offset (fixed at 52 bytes from the start for 32-bit ELF)
	err = binary.Write(w, binary.LittleEndian, uint32(52))
	if err != nil {
		return err
	}

	// Section header offset (none for simplicity)
	err = binary.Write(w, binary.LittleEndian, uint32(0))
	if err != nil {
		return err
	}

	// Flags (architecture specific)
	err = binary.Write(w, binary.LittleEndian, uint32(0))
	if err != nil {
		return err
	}

	// ELF header size
	err = binary.Write(w, binary.LittleEndian, uint16(52))
	if err != nil {
		return err
	}

	// Program header entry size
	err = binary.Write(w, binary.LittleEndian, uint16(32))
	if err != nil {
		return err
	}

	// Program header entries count (1 for simplicity)
	err = binary.Write(w, binary.LittleEndian, uint16(1))
	if err != nil {
		return err
	}

	// Section header entry size
	err = binary.Write(w, binary.LittleEndian, uint16(0))
	if err != nil {
		return err
	}

	// Section header entries count
	err = binary.Write(w, binary.LittleEndian, uint16(0))
	if err != nil {
		return err
	}

	// Section header string table index
	err = binary.Write(w, binary.LittleEndian, uint16(0))
	if err != nil {
		return err
	}

	// Program header (simplified - one loadable segment)
	const (
		PT_LOAD = 1 // Loadable segment
		PF_R    = 4 // Read permission
		PF_X    = 1 // Execute permission
	)

	// Segment type (PT_LOAD)
	err = binary.Write(w, binary.LittleEndian, uint32(PT_LOAD))
	if err != nil {
		return err
	}

	// Segment file offset (immediately after headers)
	err = binary.Write(w, binary.LittleEndian, uint32(52+32))
	if err != nil {
		return err
	}

	// Segment virtual address
	err = binary.Write(w, binary.LittleEndian, uint32(0))
	if err != nil {
		return err
	}

	// Segment physical address
	err = binary.Write(w, binary.LittleEndian, uint32(0))
	if err != nil {
		return err
	}

	// Segment size in a file
	err = binary.Write(w, binary.LittleEndian, uint32(len(cg.machineCode)))
	if err != nil {
		return err
	}

	// Segment size in memory
	err = binary.Write(w, binary.LittleEndian, uint32(len(cg.machineCode)))
	if err != nil {
		return err
	}

	// Segment flags (read + execute)
	err = binary.Write(w, binary.LittleEndian, uint32(PF_R|PF_X))
	if err != nil {
		return err
	}

	// Segment alignment
	err = binary.Write(w, binary.LittleEndian, uint32(4)) // 4-byte alignment
	if err != nil {
		return err
	}

	// Write the actual machine code
	_, err = w.Write(cg.machineCode)
	return err
}

// SymbolTable returns a copy of the symbol table
func (cg *CodeGenerator) SymbolTable() map[string]int {
	symbolTableCopy := make(map[string]int)
	for k, v := range cg.symbolTable {
		symbolTableCopy[k] = v
	}
	return symbolTableCopy
}

// Errors returns any errors encountered during code generation
func (cg *CodeGenerator) Errors() []CodegenError {
	return cg.errors
}

// MachineCode returns the generated machine code
func (cg *CodeGenerator) MachineCode() []byte {
	return cg.machineCode
}

// collectSymbols performs the first pass over the AST to gather symbol addresses
func (cg *CodeGenerator) collectSymbols() {
	// Start with the program node
	for _, node := range cg.ast.Children {
		// Each child is a line
		if node.Type == parser.NODE_LINE {
			for _, lineItem := range node.Children {
				if lineItem.Type == parser.NODE_LABEL {
					labelName := lineItem.Value.(string)
					cg.symbolTable[labelName] = cg.currentAddr
				} else if lineItem.Type == parser.NODE_INSTRUCTION {
					// Each instruction is 4 bytes
					cg.currentAddr += 4
				} else if lineItem.Type == parser.NODE_DIRECTIVE {
					// Handle directives that affect the address counter
					cg.processDirectiveForSymbols(lineItem)
				}
			}
		}
	}
}

// processDirectiveForSymbols handles directives during symbol collection phase
func (cg *CodeGenerator) processDirectiveForSymbols(node *parser.AST) {
	directive := node.Token.Literal
	switch directive {
	case "org":
		if len(node.Children) > 0 && node.Children[0].Type == parser.NODE_IMMEDIATE {
			// Set the address counter to the specified value
			addr, ok := node.Children[0].Value.(int)
			if ok {
				cg.currentAddr = addr
			}
		}
	case "word", "dw":
		// Each word takes 4 bytes
		cg.currentAddr += 4 * len(node.Children)
	case "byte", "db":
		// Each byte takes 1 byte
		cg.currentAddr += len(node.Children)
	case "section":
		if len(node.Children) > 0 && node.Children[0].Type == parser.NODE_STRING {
			// Record the section start
			cg.sectionStart = cg.currentAddr
		}
	}
}

// generateCode performs the second pass over the AST to generate machine code
func (cg *CodeGenerator) generateCode(node *parser.AST) error {
	switch node.Type {
	case parser.NODE_PROGRAM:
		for _, child := range node.Children {
			if err := cg.generateCode(child); err != nil {
				return err
			}
		}

	case parser.NODE_LINE:
		for _, child := range node.Children {
			if child.Type == parser.NODE_INSTRUCTION {
				if err := cg.generateInstruction(child); err != nil {
					return err
				}
			} else if child.Type == parser.NODE_DIRECTIVE {
				if err := cg.processDirective(child); err != nil {
					return err
				}
			}
			// Skip labels in the second pass since we already processed them
		}

	default:
		// Skip other node types
	}

	return nil
}

// generateInstruction generates machine code for an instruction node
func (cg *CodeGenerator) generateInstruction(node *parser.AST) error {
	// Extract the instruction mnemonic
	mnemonic := node.Token.Literal
	encoding, exists := encodingTable[mnemonic]

	if !exists {
		cg.errors = append(cg.errors, CodegenError{
			Message: fmt.Sprintf("unknown instruction: %s", mnemonic),
			Line:    node.Line,
			Column:  node.Column,
		})
		// Add a placeholder to maintain alignment
		cg.emitUint32(0)
		return nil
	}

	// Generate code based on instruction format
	switch encoding.Format {
	case FormatR:
		return cg.generateRFormat(node, encoding)
	case FormatI:
		return cg.generateIFormat(node, encoding)
	case FormatJ:
		return cg.generateJFormat(node, encoding)
	case FormatV:
		return cg.generateVFormat(node, encoding)
	default:
		cg.errors = append(cg.errors, CodegenError{
			Message: fmt.Sprintf("unsupported instruction format for %s", mnemonic),
			Line:    node.Line,
			Column:  node.Column,
		})
		// Add a placeholder to maintain alignment
		cg.emitUint32(0)
		return nil
	}
}

// generateRFormat generates R-format instructions (register-register operations)
func (cg *CodeGenerator) generateRFormat(node *parser.AST, encoding InstructionEncoding) error {
	// R-format: | opcode(6) | rs2(5) | rs1(5) | type(3) | rd(5) | opcode(7) |
	if len(node.Children) < 3 && encoding.Mnemonic != "NOT" && encoding.Mnemonic != "NEG" && encoding.Mnemonic != "INC" && encoding.Mnemonic != "DEC" {
		cg.addError(fmt.Sprintf("Not enough operands for %s instruction (expected 3, got %d)",
			encoding.Mnemonic, len(node.Children)), node.Line, node.Column)
		cg.emitUint32(0)
		return nil
	}

	// Extract registers
	var rd, rs1, rs2 uint8

	// Get the register value from the AST node
	rd = cg.getRegisterFromNode(node.Children[0])

	// For unary operations (NOT, NEG, INC, DEC), rs2 is unused
	if encoding.Mnemonic == "NOT" || encoding.Mnemonic == "NEG" || encoding.Mnemonic == "INC" || encoding.Mnemonic == "DEC" {
		if len(node.Children) >= 2 {
			rs1 = cg.getRegisterFromNode(node.Children[1])
		} else {
			cg.addError(fmt.Sprintf("Missing source register for %s instruction", encoding.Mnemonic),
				node.Line, node.Column)
		}
		rs2 = 0 // Not used for unary operations
	} else {
		// Binary operations (ADD, SUB, MUL, etc.)
		rs1 = cg.getRegisterFromNode(node.Children[1])
		rs2 = cg.getRegisterFromNode(node.Children[2])
	}

	// Validate instruction encoding
	if encoding.Type != OpTypeALU {
		cg.addError(fmt.Sprintf("Instruction %s has incorrect type field: expected ALU (000), got %03b",
			encoding.Mnemonic, encoding.Type), node.Line, node.Column)
	}

	// Construct the instruction with correct VTX1 encoding
	instruction := uint32(encoding.Opcode) << 26  // opcode (6 bits) in upper bits
	instruction |= uint32(rs2) << 20              // rs2 (5 bits)
	instruction |= uint32(rs1) << 15              // rs1 (5 bits)
	instruction |= uint32(encoding.Type) << 12    // type (3 bits)
	instruction |= uint32(rd) << 7                // rd (5 bits)
	instruction |= uint32(encoding.Opcode & 0x7F) // Lower 7 bits of opcode

	// Emit the instruction
	cg.emitUint32(instruction)

	return nil
}

// getRegisterFromNode extracts a register value from an AST node
func (cg *CodeGenerator) getRegisterFromNode(node *parser.AST) uint8 {
	if node == nil {
		cg.addError("Nil node passed to getRegisterFromNode", 0, 0)
		return 0
	}

	if node.Type != parser.NODE_REGISTER {
		cg.addError(fmt.Sprintf("Expected register node, got %v", node.Type), node.Line, node.Column)
		return 0
	}

	regName := node.Token.Literal
	reg, err := cg.getRegisterNumber(regName)
	if err != nil {
		cg.addError(fmt.Sprintf("Invalid register: %s", err.Error()), node.Line, node.Column)
		return 0
	}

	return reg
}

// generateIFormat generates I-format instructions (register-immediate operations)
func (cg *CodeGenerator) generateIFormat(node *parser.AST, encoding InstructionEncoding) error {
	// I-format for VTX1: | imm(12) | rs1(5) | type(3) | rd(5) | opcode(7) |
	if len(node.Children) < 2 {
		cg.addError(fmt.Sprintf("Not enough operands for %s instruction", encoding.Mnemonic),
			node.Line, node.Column)
		cg.emitUint32(0)
		return nil
	}

	// Extract registers and immediate
	var rd, rs1 uint8
	imm := 0
	var err error

	// Get destination register
	rd = cg.getRegisterFromNode(node.Children[0])

	// Handle different I-format instruction patterns
	if encoding.Type == OpTypeMEM {
		// Memory operations like LD, ST: LD rd, imm(rs1)
		if node.Children[1].Type == parser.NODE_MEMORY_REF {
			memRef := node.Children[1]
			if len(memRef.Children) >= 2 {
				imm, err = cg.getImmediateValue(memRef.Children[0])
				if err != nil {
					cg.addError(fmt.Sprintf("Invalid offset in memory reference: %s", err.Error()),
						memRef.Children[0].Line, memRef.Children[0].Column)
					imm = 0
				}

				rs1 = cg.getRegisterFromNode(memRef.Children[1])
			} else {
				cg.addError("Incomplete memory reference, expected format: imm(rs)",
					node.Children[1].Line, node.Children[1].Column)
			}
		} else {
			cg.addError("Invalid memory reference format for load/store instruction",
				node.Children[1].Line, node.Children[1].Column)
		}
	} else {
		// ALU immediate operations like ADDI: ADDI rd, rs1, imm
		rs1 = cg.getRegisterFromNode(node.Children[1])

		if len(node.Children) >= 3 {
			imm, err = cg.getImmediateValue(node.Children[2])
			if err != nil {
				cg.addError(fmt.Sprintf("Invalid immediate value: %s", err.Error()),
					node.Children[2].Line, node.Children[2].Column)
				imm = 0
			}
		} else {
			cg.addError(fmt.Sprintf("Missing immediate value for %s instruction", encoding.Mnemonic),
				node.Line, node.Column)
		}
	}

	// Check if immediate value fits in 12 bits with sign extension
	if imm > 2047 || imm < -2048 {
		cg.addError(fmt.Sprintf("Immediate value %d out of range for 12-bit signed value (-2048 to 2047)",
			imm), node.Line, node.Column)
		imm &= 0xFFF // Truncate to 12 bits
	} else {
		// Ensure 12-bit 2's complement representation
		imm &= 0xFFF
	}

	// Construct the instruction with correct encoding for VTX1
	instruction := uint32(imm) << 20              // imm (12 bits)
	instruction |= uint32(rs1) << 15              // rs1 (5 bits)
	instruction |= uint32(encoding.Type) << 12    // type (3 bits)
	instruction |= uint32(rd) << 7                // rd (5 bits)
	instruction |= uint32(encoding.Opcode & 0x7F) // opcode (7 bits)

	// Emit the instruction
	cg.emitUint32(instruction)

	return nil
}

// generateJFormat generates J-format instructions (jumps)
func (cg *CodeGenerator) generateJFormat(node *parser.AST, encoding InstructionEncoding) error {
	// J-format: | imm(20) | rd(5) | opcode(7) |
	if len(node.Children) < 1 {
		cg.addError(fmt.Sprintf("Not enough operands for %s instruction", encoding.Mnemonic),
			node.Line, node.Column)
		cg.emitUint32(0)
		return nil
	}

	// Extract destination register and target
	var rd uint8
	target := 0
	var err error

	// For CALL, the first operand is the destination register
	if encoding.Mnemonic == "CALL" {
		rd = cg.getRegisterFromNode(node.Children[0])
	}

	// Get the jump target
	targetIdx := 0
	if encoding.Mnemonic == "CALL" {
		targetIdx = 1
	}

	if len(node.Children) <= targetIdx {
		cg.addError(fmt.Sprintf("Missing jump target for %s instruction", encoding.Mnemonic),
			node.Line, node.Column)
		target = 0
	} else if node.Children[targetIdx].Type == parser.NODE_SYMBOL_REF {
		symbolName := node.Children[targetIdx].Token.Literal
		if addr, exists := cg.symbolTable[symbolName]; exists {
			// Calculate relative offset
			target = addr - cg.currentAddr
		} else {
			cg.addError(fmt.Sprintf("Undefined symbol: %s", symbolName),
				node.Children[targetIdx].Line, node.Children[targetIdx].Column)
			target = 0
		}
	} else {
		target, err = cg.getImmediateValue(node.Children[targetIdx])
		if err != nil {
			cg.addError(fmt.Sprintf("Invalid jump target: %s", err.Error()),
				node.Children[targetIdx].Line, node.Children[targetIdx].Column)
			target = 0
		}
	}

	// Ensure target fits in 20 bits
	target &= 0xFFFFF

	// Construct the instruction
	instruction := uint32(target) << 12
	instruction |= uint32(rd) << 7
	instruction |= uint32(encoding.Opcode)

	// Emit the instruction
	cg.emitUint32(instruction)

	return nil
}

// generateVFormat generates VLIW-format instructions
func (cg *CodeGenerator) generateVFormat(node *parser.AST, encoding InstructionEncoding) error {
	// VLIW instructions in VTX1 combine multiple operations in a single instruction
	// Format: | op1(10) | op2(10) | op3(10) | format(2) |

	// Check if this is a VLIW_INSTRUCTION node
	if node.Type != parser.NODE_VLIW_INSTRUCTION {
		cg.addError(fmt.Sprintf("Invalid VLIW instruction format for %s", encoding.Mnemonic), node.Line, node.Column)
		cg.emitUint32(0) // Placeholder
		return nil
	}

	// Initialize a new VLIW instruction
	vliw := VLIWInstruction{
		Operations: make([]uint32, 0, 3), // VTX1 can bundle up to 3 operations
		LineInfo:   fmt.Sprintf("%d:%d", node.Line, node.Column),
	}

	// Process each operation in the VLIW instruction
	for _, opNode := range node.Children {
		if opNode.Type != parser.NODE_INSTRUCTION {
			cg.addError("VLIW element must be a valid instruction", opNode.Line, opNode.Column)
			continue
		}

		// Get the encoding for this operation
		opMnemonic := opNode.Token.Literal
		opEncoding, exists := encodingTable[opMnemonic]
		if !exists {
			cg.addError(fmt.Sprintf("Unknown instruction '%s' in VLIW bundle", opMnemonic), opNode.Line, opNode.Column)
			continue
		}

		// Generate the operation code based on its format
		var opCode uint32
		var err error

		switch opEncoding.Format {
		case FormatR:
			// For R-format operations in VLIW, we encode them in a compressed format
			// Format: | type(3) | opcode(6) | rd(5) | rs1(5) | rs2(5) | reserved(8) |

			if len(opNode.Children) < 3 {
				cg.addError(fmt.Sprintf("Not enough operands for %s in VLIW", opEncoding.Mnemonic), opNode.Line, opNode.Column)
				continue
			}

			rd, rs1, rs2 := 0, 0, 0

			rd = int(cg.getRegisterFromNode(opNode.Children[0]))
			rs1 = int(cg.getRegisterFromNode(opNode.Children[1]))
			rs2 = int(cg.getRegisterFromNode(opNode.Children[2]))

			// Encode in VLIW-compressed R-format
			opCode = uint32(opEncoding.Type) << 29    // Type (3 bits)
			opCode |= uint32(opEncoding.Opcode) << 23 // Opcode (6 bits)
			opCode |= uint32(rd) << 18                // rd (5 bits)
			opCode |= uint32(rs1) << 13               // rs1 (5 bits)
			opCode |= uint32(rs2) << 8                // rs2 (5 bits)
			// bits 0-7 are reserved (set to 0)

		case FormatI:
			// For I-format operations in VLIW, we encode them in a compressed format
			// Format: | type(3) | opcode(6) | rd(5) | rs1(5) | imm(13) |

			if len(opNode.Children) < 2 {
				cg.addError(fmt.Sprintf("Not enough operands for %s in VLIW", opEncoding.Mnemonic), opNode.Line, opNode.Column)
				continue
			}

			rd, rs1 := 0, 0
			imm := 0

			rd = int(cg.getRegisterFromNode(opNode.Children[0]))

			// Handle memory references specially
			if opNode.Children[1].Type == parser.NODE_MEMORY_REF {
				memRef := opNode.Children[1]
				if len(memRef.Children) >= 2 {
					imm, err = cg.getImmediateValue(memRef.Children[0])
					if err != nil {
						cg.addError(err.Error(), memRef.Line, memRef.Column)
						imm = 0
					}

					rs1 = int(cg.getRegisterFromNode(memRef.Children[1]))
				}
			} else {
				// Regular immediate format
				rs1 = int(cg.getRegisterFromNode(opNode.Children[1]))

				if len(opNode.Children) >= 3 {
					imm, err = cg.getImmediateValue(opNode.Children[2])
					if err != nil {
						cg.addError(err.Error(), opNode.Line, opNode.Column)
						imm = 0
					}
				}
			}

			// Ensure immediate fits in 13 bits
			imm &= 0x1FFF

			// Encode in VLIW-compressed I-format
			opCode = uint32(opEncoding.Type) << 29    // Type (3 bits)
			opCode |= uint32(opEncoding.Opcode) << 23 // Opcode (6 bits)
			opCode |= uint32(rd) << 18                // rd (5 bits)
			opCode |= uint32(rs1) << 13               // rs1 (5 bits)
			opCode |= uint32(imm)                     // imm (13 bits)

		default:
			cg.addError(fmt.Sprintf("Instruction '%s' has unsupported format for VLIW", opEncoding.Mnemonic), opNode.Line, opNode.Column)
			continue
		}

		// Add the encoded operation to the VLIW instruction
		vliw.Operations = append(vliw.Operations, opCode)
	}

	// Check if we have any operations in the VLIW bundle
	if len(vliw.Operations) == 0 {
		cg.addError("Empty VLIW instruction bundle", node.Line, node.Column)
		cg.emitUint32(0) // Placeholder
		return nil
	}

	// Encode the final VLIW instruction
	// In VTX1, a VLIW instruction is 32-bits with:
	// - For 1 operation: | op1(30) | format(2)=01 |
	// - For 2 operations: | op1(15) | op2(15) | format(2)=10 |
	// - For 3 operations: | op1(10) | op2(10) | op3(10) | format(2)=11 |

	var vliwInstruction uint32

	switch len(vliw.Operations) {
	case 1:
		// Single operation VLIW
		vliwInstruction = vliw.Operations[0] << 2
		vliwInstruction |= 0x1 // Format 01

	case 2:
		// Two operations VLIW
		op1 := vliw.Operations[0] & 0x7FFF // Take lower 15 bits
		op2 := vliw.Operations[1] & 0x7FFF // Take lower 15 bits

		vliwInstruction = op1 << 17
		vliwInstruction |= op2 << 2
		vliwInstruction |= 0x2 // Format 10

	case 3:
		// Three operations VLIW
		op1 := vliw.Operations[0] & 0x3FF // Take lower 10 bits
		op2 := vliw.Operations[1] & 0x3FF // Take lower 10 bits
		op3 := vliw.Operations[2] & 0x3FF // Take lower 10 bits

		vliwInstruction = op1 << 22
		vliwInstruction |= op2 << 12
		vliwInstruction |= op3 << 2
		vliwInstruction |= 0x3 // Format 11

	default:
		cg.addError(fmt.Sprintf("VLIW instruction with %d operations not supported (max 3)", len(vliw.Operations)), node.Line, node.Column)
		cg.emitUint32(0) // Placeholder
		return nil
	}

	// Emit the VLIW instruction
	cg.emitUint32(vliwInstruction)
	cg.currentAddr += 4

	return nil
}

// startVLIW begins a new VLIW instruction bundle
func (cg *CodeGenerator) startVLIW(line, column int) {
	cg.vliwBuffer = &VLIWInstruction{
		Operations: make([]uint32, 0, 3),
		LineInfo:   fmt.Sprintf("%d:%d", line, column),
	}
}

// endVLIW finalizes and emits the current VLIW instruction bundle
func (cg *CodeGenerator) endVLIW() {
	if cg.vliwBuffer == nil || len(cg.vliwBuffer.Operations) == 0 {
		return
	}

	// Encode the VLIW instruction
	var vliwInstruction uint32

	switch len(cg.vliwBuffer.Operations) {
	case 1:
		vliwInstruction = cg.vliwBuffer.Operations[0] << 2
		vliwInstruction |= 0x1 // Format 01

	case 2:
		op1 := cg.vliwBuffer.Operations[0] & 0x7FFF
		op2 := cg.vliwBuffer.Operations[1] & 0x7FFF

		vliwInstruction = op1 << 17
		vliwInstruction |= op2 << 2
		vliwInstruction |= 0x2 // Format 10

	case 3:
		op1 := cg.vliwBuffer.Operations[0] & 0x3FF
		op2 := cg.vliwBuffer.Operations[1] & 0x3FF
		op3 := cg.vliwBuffer.Operations[2] & 0x3FF

		vliwInstruction = op1 << 22
		vliwInstruction |= op2 << 12
		vliwInstruction |= op3 << 2
		vliwInstruction |= 0x3 // Format 11
	}

	// Emit the instruction
	cg.emitUint32(vliwInstruction)
	cg.currentAddr += 4

	// Reset the VLIW buffer
	cg.vliwBuffer = nil
}

// emitUint32 writes a 32-bit value to the binary output
func (cg *CodeGenerator) emitUint32(value uint32) {
	// Create a 4-byte buffer
	buf := make([]byte, 4)

	// Write the value in little-endian format
	binary.LittleEndian.PutUint32(buf, value)

	// Append to the binary output
	cg.machineCode = append(cg.machineCode, buf...)

	// Increment the current position
	cg.currentAddr += 4
}

// getRegisterNumber converts a register token to its numeric value
func (cg *CodeGenerator) getRegisterNumber(reg string) (uint8, error) {
	// Handle general purpose registers
	if len(reg) == 2 && reg[0] == 'T' && reg[1] >= '0' && reg[1] <= '6' {
		return uint8(reg[1] - '0'), nil
	}

	// Handle special registers
	switch reg {
	case "TA":
		return 7, nil
	case "TB":
		return 8, nil
	case "TC":
		return 9, nil
	case "TS":
		return 10, nil
	case "TI":
		return 11, nil
	case "VA":
		return 12, nil
	case "VB":
		return 13, nil
	case "VC":
		return 14, nil
	case "FA":
		return 15, nil
	case "FB":
		return 16, nil
	case "FC":
		return 17, nil
	}

	return 0, fmt.Errorf("invalid register: %s", reg)
}

// processDirective handles assembly directives
func (cg *CodeGenerator) processDirective(node *parser.AST) error {
	directive := node.Token.Literal

	switch directive {
	case ".ORG":
		// Set the current address
		if len(node.Children) != 1 {
			return cg.addError("ORG directive requires one argument", node.Line, node.Column)
		}

		// Parse the address
		address, err := cg.evaluateImmediate(node.Children[0])
		if err != nil {
			return err
		}

		// Set the current address
		cg.currentAddr = int(address)

	case ".DB", ".DW", ".DT":
		// Define data bytes, words, or ternary values
		for _, child := range node.Children {
			value, err := cg.evaluateImmediate(child)
			if err != nil {
				return err
			}

			if directive == ".DB" {
				// Emit 8-bit value
				cg.machineCode = append(cg.machineCode, byte(value))
				cg.currentAddr++
			} else if directive == ".DW" {
				// Emit 16-bit value
				buf := make([]byte, 2)
				binary.LittleEndian.PutUint16(buf, uint16(value))
				cg.machineCode = append(cg.machineCode, buf...)
				cg.currentAddr += 2
			} else if directive == ".DT" {
				// Emit 32-bit value for ternary
				cg.emitUint32(uint32(value))
			}
		}

	case ".ALIGN":
		// Align the current address
		if len(node.Children) != 1 {
			return cg.addError("ALIGN directive requires one argument", node.Line, node.Column)
		}

		// Parse the alignment value
		alignment, err := cg.evaluateImmediate(node.Children[0])
		if err != nil {
			return err
		}

		// Calculate padding needed
		alignmentVal := int(alignment)
		remainder := cg.currentAddr % alignmentVal
		if remainder > 0 {
			padding := alignmentVal - remainder
			// Add padding bytes
			for i := 0; i < padding; i++ {
				cg.machineCode = append(cg.machineCode, 0)
				cg.currentAddr++
			}
		}

	case ".SPACE":
		// Reserve space
		if len(node.Children) != 1 {
			return cg.addError("SPACE directive requires one argument", node.Line, node.Column)
		}

		// Parse the size
		size, err := cg.evaluateImmediate(node.Children[0])
		if err != nil {
			return err
		}

		// Add the space
		sizeVal := int(size)
		for i := 0; i < sizeVal; i++ {
			cg.machineCode = append(cg.machineCode, 0)
			cg.currentAddr++
		}

	case ".EQU":
		// Define a constant
		if len(node.Children) != 2 {
			return cg.addError("EQU directive requires two arguments", node.Line, node.Column)
		}

		// Get the identifier
		if node.Children[0].Type != parser.NODE_SYMBOL_REF {
			return cg.addError("First argument of EQU must be an identifier", node.Line, node.Column)
		}

		// Get the value
		value, err := cg.evaluateImmediate(node.Children[1])
		if err != nil {
			return err
		}

		// Add to symbol table
		cg.symbolTable[node.Children[0].Token.Literal] = int(value)

	case ".INCLUDE":
		// Include another file
		// This would normally be handled by the parser, but we'll add a placeholder here
		return cg.addError("INCLUDE directive not implemented in this version", node.Line, node.Column)

	case ".SECTION":
		// Define a section
		// This would normally change the current section, but we'll add a placeholder here
		return cg.addError("SECTION directive not implemented in this version", node.Line, node.Column)

	default:
		return cg.addError(fmt.Sprintf("Unknown directive: %s", directive), node.Line, node.Column)
	}

	return nil
}

// evaluateImmediate evaluates an immediate value node
func (cg *CodeGenerator) evaluateImmediate(node *parser.AST) (int64, error) {
	if node == nil {
		return 0, fmt.Errorf("nil node passed to evaluateImmediate")
	}

	switch node.Type {
	case parser.NODE_IMMEDIATE:
		// Direct immediate value
		literal := node.Token.Literal

		// Check the type from the token
		switch node.Token.Type {
		case lexer.DECIMAL:
			// Parse decimal
			value, err := strconv.ParseInt(literal, 10, 64)
			if err != nil {
				return 0, cg.addError(fmt.Sprintf("Invalid decimal: %s", literal), node.Line, node.Column)
			}
			return value, nil

		case lexer.HEXADECIMAL:
			// Parse hexadecimal (remove 0x prefix)
			value, err := strconv.ParseInt(strings.TrimPrefix(literal, "0x"), 16, 64)
			if err != nil {
				return 0, cg.addError(fmt.Sprintf("Invalid hexadecimal: %s", literal), node.Line, node.Column)
			}
			return value, nil

		case lexer.BINARY:
			// Parse binary (remove 0b prefix)
			value, err := strconv.ParseInt(strings.TrimPrefix(literal, "0b"), 2, 64)
			if err != nil {
				return 0, cg.addError(fmt.Sprintf("Invalid binary: %s", literal), node.Line, node.Column)
			}
			return value, nil

		case lexer.TERNARY:
			// Parse balanced ternary (remove 0t prefix)
			ternaryStr := strings.TrimPrefix(literal, "0t")
			value := int64(0)

			// Process each ternary digit
			for i := 0; i < len(ternaryStr); i++ {
				// Shift existing value
				value *= 3

				// Add new digit
				switch ternaryStr[i] {
				case '+':
					value += 1
				case '-':
					value -= 1
				case '0':
					// No change
				default:
					return 0, cg.addError(fmt.Sprintf("Invalid ternary digit: %c", ternaryStr[i]), node.Line, node.Column)
				}
			}

			return value, nil
		default:
			// No modifier, do nothing
		}

	case parser.NODE_SYMBOL_REF:
		// Symbol reference
		symbol := node.Token.Literal

		// Look up in symbol table
		value, exists := cg.symbolTable[symbol]
		if !exists {
			return 0, cg.addError(fmt.Sprintf("Undefined symbol: %s", symbol), node.Line, node.Column)
		}

		return int64(value), nil
	default:
		// No modifier, do nothing
	}

	return 0, cg.addError(fmt.Sprintf("Cannot evaluate node as immediate: %v", node.Type), node.Line, node.Column)
}

// getImmediateValue evaluates an immediate value from a node and returns as integer
func (cg *CodeGenerator) getImmediateValue(node *parser.AST) (int, error) {
	if node == nil {
		return 0, fmt.Errorf("nil node passed to getImmediateValue")
	}

	// Evaluate the immediate value
	value, err := cg.evaluateImmediate(node)
	if err != nil {
		return 0, err
	}

	// Convert to integer (possibly truncating)
	return int(value), nil
}
