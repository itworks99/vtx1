package cmd

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// CodeGenerator is responsible for traversing the AST and emitting machine code.
type CodeGenerator struct {
	Output      []byte
	SymbolTable *SymbolTable
	CurrentAddr uint32
	Labels      map[string]uint32
	Equs        map[string]uint32
}

// NewCodeGenerator creates a new code generator with the given symbol table.
func NewCodeGenerator(symbolTable *SymbolTable) *CodeGenerator {
	return &CodeGenerator{
		Output:      make([]byte, 0),
		SymbolTable: symbolTable,
		CurrentAddr: 0,
		Labels:      make(map[string]uint32),
		Equs:        make(map[string]uint32),
	}
}

// Pass 1: Collect labels and .EQUs, handle .ORG/.SPACE for address tracking
func (cg *CodeGenerator) collectSymbols(ast *AST) error {
	addr := uint32(0)
	for _, line := range ast.Program.Lines {
		if line.Label != nil {
			cg.Labels[line.Label.Name] = addr
		}
		if line.Statement == nil {
			continue
		}
		switch stmt := line.Statement.(type) {
		case *InstructionNode:
			addr += 4
		case *VLIWInstructionNode:
			addr += 12
		case *DirectiveNode:
			name := strings.ToUpper(stmt.Name)
			switch name {
			case ".ORG":
				if len(stmt.Params) > 0 {
					imm, _ := parseImmediateOperand(stmt.Params[0])
					addr = uint32(imm)
				}
			case ".SPACE":
				if len(stmt.Params) > 0 {
					imm, _ := parseImmediateOperand(stmt.Params[0])
					addr += uint32(imm)
				}
			case ".DW":
				addr += 2 * uint32(len(stmt.Params))
			case ".DB":
				addr += uint32(len(stmt.Params))
			case ".EQU":
				if len(stmt.Params) == 2 {
					if id, ok := stmt.Params[0].(*IdentifierNode); ok {
						imm, _ := parseImmediateOperand(stmt.Params[1])
						cg.Equs[id.Name] = uint32(imm)
					}
				}
			}
		}
	}
	return nil
}

// Pass 2: Emit code/data, resolving symbols
func (cg *CodeGenerator) Generate(ast *AST) error {
	fmt.Println("[DEBUG] CodeGenerator.Generate called")
	cg.Output = make([]byte, 0)
	cg.CurrentAddr = 0
	cg.Labels = make(map[string]uint32)
	cg.Equs = make(map[string]uint32)
	if err := cg.collectSymbols(ast); err != nil {
		return err
	}
	cg.CurrentAddr = 0
	for i, line := range ast.Program.Lines {
		if line.Statement == nil {
			continue
		}
		typeName := reflect.TypeOf(line.Statement)
		fmt.Printf("[DEBUG] Generate: line %d, reflect.TypeOf=%v, type=%T, label=%v, statement=%#v\n", i, typeName, line.Statement, line.Label, line.Statement)
		switch stmt := line.Statement.(type) {
		case *InstructionNode:
			fmt.Printf("[DEBUG] emitInstruction called: %+v\n", stmt)
			if err := cg.emitInstruction(stmt); err != nil {
				return err
			}
		case *VLIWInstructionNode:
			fmt.Printf("[DEBUG] emitVLIWInstruction called: %+v\n", stmt)
			if err := cg.emitVLIWInstruction(stmt); err != nil {
				return err
			}
		case *DirectiveNode:
			fmt.Printf("[DEBUG] emitDirective called: %+v\n", stmt)
			if err := cg.emitDirective(stmt); err != nil {
				return err
			}
		default:
			fmt.Printf("[DEBUG] Generate: unhandled node type %T\n", stmt)
		}
	}
	fmt.Printf("[DEBUG] CodeGenerator output length: %d bytes\n", len(cg.Output))
	return nil
}

// --- Instruction Encoding ---
var opcodeMap = map[string]byte{
	// ALU
	"NOP": 0x00, "ADD": 0x01, "SUB": 0x02, "MUL": 0x03, "AND": 0x04, "OR": 0x05, "NOT": 0x06, "XOR": 0x07,
	"SHL": 0x08, "SHR": 0x09, "ROL": 0x0A, "ROR": 0x0B, "CMP": 0x0C, "TEST": 0x0D, "INC": 0x0E, "DEC": 0x0F, "NEG": 0x10,
	// MEM
	"LD": 0x20, "ST": 0x21, "VLD": 0x22, "VST": 0x23, "FLD": 0x24, "FST": 0x25, "LEA": 0x26, "PUSH": 0x27, "POP": 0x28,
	// CTRL
	"JMP": 0x30, "JAL": 0x31, "JR": 0x32, "JALR": 0x33, "BEQ": 0x34, "BNE": 0x35, "BLT": 0x36, "BGE": 0x37, "BLTU": 0x38, "BGEU": 0x39, "BGT": 0x3A, "BLE": 0x3B, "CALL": 0x3C, "RET": 0x3D,
	// VEC
	"VADD": 0x40, "VSUB": 0x41, "VMUL": 0x42, "VAND": 0x43, "VOR": 0x44, "VNOT": 0x45, "VSHL": 0x46, "VSHR": 0x47,
	// FP
	"FADD": 0x50, "FSUB": 0x51, "FMUL": 0x52, "FCMP": 0x53, "FMOV": 0x54, "FNEG": 0x55,
	// SYS
	"WFI": 0x60,
	// COMPLEX
	"DIV": 0x70, "MOD": 0x71, "UDIV": 0x72, "UMOD": 0x73, "SQRT": 0x74, "ABS": 0x75, "SIN": 0x76, "COS": 0x77, "TAN": 0x78, "ASIN": 0x79, "ACOS": 0x7A, "ATAN": 0x7B, "EXP": 0x7C, "LOG": 0x7D,
	// COMPLEX_VEC
	"VDOT": 0x80, "VREDUCE": 0x81, "VMAX": 0x82, "VMIN": 0x83, "VSUM": 0x84, "VPERM": 0x85,
	// COMPLEX_MEM
	"CACHE": 0x90, "FLUSH": 0x91, "MEMBAR": 0x92,
	// COMPLEX_SYS
	"SYSCALL": 0xA0, "BREAK": 0xA1, "HALT": 0xA2,
}

func regNum(name string) (byte, error) {
	if name == "TB" {
		return 7, nil
	}
	if name == "TA" {
		return 8, nil
	}
	if name == "TC" {
		return 9, nil
	}
	if name == "TS" {
		return 10, nil
	}
	if name == "TI" {
		return 11, nil
	}
	if name == "T0" || name == "T1" || name == "T2" || name == "T3" || name == "T4" || name == "T5" || name == "T6" {
		return name[1] - '0', nil
	}
	return 0, fmt.Errorf("unknown register: %s", name)
}

// Add a global set of deprecated mnemonics and directives
var deprecatedMnemonics = map[string]bool{
	"OLDOP": true, // Example: replace with real deprecated mnemonics
}
var deprecatedDirectives = map[string]bool{
	".OLD": true, // Example: replace with real deprecated directives
}

// Add a pointer to ErrorManager.Warnings for warning reporting
var codegenWarnings *[]error

func AttachCodegenWarnings(w *[]error) {
	codegenWarnings = w
}

func (cg *CodeGenerator) emitInstruction(instr *InstructionNode) error {
	fmt.Printf("[DEBUG] In emitInstruction: %+v\n", instr)
	if deprecatedMnemonics[strings.ToUpper(instr.Mnemonic)] && codegenWarnings != nil {
		*codegenWarnings = append(*codegenWarnings, fmt.Errorf("warning: instruction '%s' at line %d is deprecated", instr.Mnemonic, instr.Line))
	}
	opc, ok := opcodeMap[strings.ToUpper(instr.Mnemonic)]
	if !ok {
		return fmt.Errorf("unsupported instruction: %s at line %d", instr.Mnemonic, instr.Line)
	}
	var out [4]byte
	out[0] = opc
	// Encoding: opcode | dst | src1 | src2/imm
	switch opc {
	case 0x00: // NOP
		// nothing more
	case 0x01, 0x02, 0x03, 0x04, 0x05, 0x07: // ADD, SUB, MUL, AND, OR, XOR
		if len(instr.Operands) != 3 {
			return fmt.Errorf("%s requires 3 operands at line %d", instr.Mnemonic, instr.Line)
		}
		for i := 0; i < 3; i++ {
			reg, ok := instr.Operands[i].(*RegisterNode)
			if !ok {
				return fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			rn, err := regNum(reg.Name)
			if err != nil {
				return err
			}
			out[i+1] = rn
		}
	case 0x06: // NOT
		if len(instr.Operands) != 2 {
			return fmt.Errorf("NOT requires 2 operands at line %d", instr.Line)
		}
		for i := 0; i < 2; i++ {
			reg, ok := instr.Operands[i].(*RegisterNode)
			if !ok {
				return fmt.Errorf("NOT operand %d is not a register at line %d", i+1, instr.Line)
			}
			rn, err := regNum(reg.Name)
			if err != nil {
				return err
			}
			out[i+1] = rn
		}
	case 0x10: // LD
		if len(instr.Operands) != 2 {
			return fmt.Errorf("LD requires 2 operands at line %d", instr.Line)
		}
		dst, ok := instr.Operands[0].(*RegisterNode)
		if !ok {
			return fmt.Errorf("LD dst must be register")
		}
		out[1], _ = regNum(dst.Name)
		// src can be ImmediateNode, IdentifierNode, or MemoryOperandNode
		switch src := instr.Operands[1].(type) {
		case *ImmediateNode:
			imm, _ := parseImmediateOperand(src)
			binary.BigEndian.PutUint16(out[2:], uint16(imm))
		case *IdentifierNode:
			addr := cg.resolveSymbol(src.Name)
			binary.BigEndian.PutUint16(out[2:], uint16(addr))
		case *MemoryOperandNode:
			addr := cg.resolveMemOperand(src)
			binary.BigEndian.PutUint16(out[2:], uint16(addr))
		default:
			return fmt.Errorf("LD src must be immediate, label, or memory")
		}
	case 0x11: // ST
		if len(instr.Operands) != 2 {
			return fmt.Errorf("ST requires 2 operands at line %d", instr.Line)
		}
		src, ok := instr.Operands[0].(*RegisterNode)
		if !ok {
			return fmt.Errorf("ST src must be register")
		}
		out[1], _ = regNum(src.Name)
		switch dst := instr.Operands[1].(type) {
		case *ImmediateNode:
			imm, _ := parseImmediateOperand(dst)
			binary.BigEndian.PutUint16(out[2:], uint16(imm))
		case *IdentifierNode:
			addr := cg.resolveSymbol(dst.Name)
			binary.BigEndian.PutUint16(out[2:], uint16(addr))
		case *MemoryOperandNode:
			addr := cg.resolveMemOperand(dst)
			binary.BigEndian.PutUint16(out[2:], uint16(addr))
		default:
			return fmt.Errorf("ST dst must be immediate, label, or memory")
		}
	case 0x20: // JMP
		if len(instr.Operands) != 1 {
			return fmt.Errorf("JMP requires 1 operand at line %d", instr.Line)
		}
		addr := cg.resolveOperandAddr(instr.Operands[0])
		binary.BigEndian.PutUint16(out[2:], uint16(addr))
	case 0x21, 0x22, 0x23, 0x24: // BEQ, BNE, BLT, BGT
		if len(instr.Operands) != 3 {
			return fmt.Errorf("%s requires 3 operands at line %d", instr.Mnemonic, instr.Line)
		}
		for i := 0; i < 2; i++ {
			reg, ok := instr.Operands[i].(*RegisterNode)
			if !ok {
				return fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			rn, err := regNum(reg.Name)
			if err != nil {
				return err
			}
			out[i+1] = rn
		}
		addr := cg.resolveOperandAddr(instr.Operands[2])
		binary.BigEndian.PutUint16(out[2:], uint16(addr))
	case 0x30: // WFI
		// nothing more
	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47: // VEC
		if len(instr.Operands) != 3 {
			return fmt.Errorf("%s requires 3 operands at line %d", instr.Mnemonic, instr.Line)
		}
		vecRegs := map[string]byte{"VA": 0, "VT": 1, "VB": 2}
		for i := 0; i < 3; i++ {
			reg, ok := instr.Operands[i].(*RegisterNode)
			if !ok {
				return fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			vn, ok := vecRegs[strings.ToUpper(reg.Name)]
			if !ok {
				return fmt.Errorf("%s operand %d is not a vector register (VA/VT/VB) at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			out[i+1] = vn
		}
	case 0x50, 0x51, 0x52, 0x53, 0x54, 0x55: // FP
		if len(instr.Operands) != 3 {
			return fmt.Errorf("%s requires 3 operands at line %d", instr.Mnemonic, instr.Line)
		}
		fpRegs := map[string]byte{"FA": 0, "FT": 1, "FB": 2}
		for i := 0; i < 3; i++ {
			reg, ok := instr.Operands[i].(*RegisterNode)
			if !ok {
				return fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			fn, ok := fpRegs[strings.ToUpper(reg.Name)]
			if !ok {
				return fmt.Errorf("%s operand %d is not a floating point register (FA/FT/FB) at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			out[i+1] = fn
		}
	case 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D: // COMPLEX
		// Most are binary, some unary (e.g., SQRT, ABS, SIN, COS, TAN, ASIN, ACOS, ATAN, EXP, LOG)
		unaryOps := map[string]bool{"SQRT": true, "ABS": true, "SIN": true, "COS": true, "TAN": true, "ASIN": true, "ACOS": true, "ATAN": true, "EXP": true, "LOG": true}
		if unaryOps[strings.ToUpper(instr.Mnemonic)] {
			if len(instr.Operands) != 2 {
				return fmt.Errorf("%s requires 2 operands (dst, src) at line %d", instr.Mnemonic, instr.Line)
			}
			for i := 0; i < 2; i++ {
				reg, ok := instr.Operands[i].(*RegisterNode)
				if !ok {
					return fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
				}
				rn, err := regNum(reg.Name)
				if err != nil {
					return err
				}
				out[i+1] = rn
			}
			out[3] = 0 // unused
		} else {
			if len(instr.Operands) != 3 {
				return fmt.Errorf("%s requires 3 operands (dst, src1, src2) at line %d", instr.Mnemonic, instr.Line)
			}
			for i := 0; i < 3; i++ {
				reg, ok := instr.Operands[i].(*RegisterNode)
				if !ok {
					return fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
				}
				rn, err := regNum(reg.Name)
				if err != nil {
					return err
				}
				out[i+1] = rn
			}
		}
	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85: // COMPLEX_VEC
		if len(instr.Operands) != 3 {
			return fmt.Errorf("%s requires 3 operands at line %d", instr.Mnemonic, instr.Line)
		}
		vecRegs := map[string]byte{"VA": 0, "VT": 1, "VB": 2}
		for i := 0; i < 3; i++ {
			reg, ok := instr.Operands[i].(*RegisterNode)
			if !ok {
				return fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			vn, ok := vecRegs[strings.ToUpper(reg.Name)]
			if !ok {
				return fmt.Errorf("%s operand %d is not a vector register (VA/VT/VB) at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			out[i+1] = vn
		}
	case 0x90, 0x91, 0x92: // COMPLEX_MEM
		if len(instr.Operands) != 0 {
			return fmt.Errorf("%s does not take any operands at line %d", instr.Mnemonic, instr.Line)
		}
		// All bytes except opcode are zero
		// No further action needed
	case 0xA0, 0xA1, 0xA2: // COMPLEX_SYS
		if len(instr.Operands) != 0 {
			return fmt.Errorf("%s does not take any operands at line %d", instr.Mnemonic, instr.Line)
		}
		// All bytes except opcode are zero
		// No further action needed
	default:
		return fmt.Errorf("unsupported opcode: %02X", opc)
	}
	cg.Output = append(cg.Output, out[:]...)
	cg.CurrentAddr += 4
	fmt.Printf("[DEBUG] emitInstruction: appended %d bytes, opcode=%02X, addr=0x%X\n", len(out), out[0], cg.CurrentAddr)
	return nil
}

func (cg *CodeGenerator) emitVLIWInstruction(vliw *VLIWInstructionNode) error {
	fmt.Printf("[DEBUG] In emitVLIWInstruction: %+v\n", vliw)
	const vliwWordSize = 12
	var word [vliwWordSize]byte
	usedDestRegs := make(map[byte]bool)
	for i := 0; i < 3; i++ {
		var instr *InstructionNode
		if i < len(vliw.Instructions) {
			instr = vliw.Instructions[i]
		} else {
			instr = &InstructionNode{Mnemonic: "NOP", Operands: nil, Line: vliw.Line}
		}
		enc, destReg, err := cg.encodeVLIWSubInstr(instr)
		if err != nil {
			return fmt.Errorf("VLIW error at line %d: %v", instr.Line, err)
		}
		if destReg != 0xFF {
			if usedDestRegs[destReg] {
				return fmt.Errorf("VLIW resource conflict: destination register T%d written by more than one instruction in the same VLIW word (line %d)", destReg, instr.Line)
			}
			usedDestRegs[destReg] = true
		}
		copy(word[i*4:(i+1)*4], enc[:])
	}
	cg.Output = append(cg.Output, word[:]...)
	cg.CurrentAddr += vliwWordSize
	return nil
}

func (cg *CodeGenerator) encodeVLIWSubInstr(instr *InstructionNode) ([4]byte, byte, error) {
	var out [4]byte
	opc, ok := opcodeMap[strings.ToUpper(instr.Mnemonic)]
	if !ok {
		return out, 0xFF, fmt.Errorf("unsupported instruction: %s at line %d", instr.Mnemonic, instr.Line)
	}
	out[0] = opc
	var destReg byte = 0xFF
	switch opc {
	case 0x00: // NOP
		// nothing more
	case 0x01, 0x02, 0x03, 0x04, 0x05, 0x07: // ADD, SUB, MUL, AND, OR, XOR
		if len(instr.Operands) != 3 {
			return out, 0xFF, fmt.Errorf("%s requires 3 operands at line %d", instr.Mnemonic, instr.Line)
		}
		for i := 0; i < 3; i++ {
			reg, ok := instr.Operands[i].(*RegisterNode)
			if !ok {
				return out, 0xFF, fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			rn, err := regNum(reg.Name)
			if err != nil {
				return out, 0xFF, err
			}
			out[i+1] = rn
		}
		destReg = out[1]
	case 0x10: // LD
		if len(instr.Operands) != 2 {
			return out, 0xFF, fmt.Errorf("LD requires 2 operands at line %d", instr.Line)
		}
		dst, ok := instr.Operands[0].(*RegisterNode)
		if !ok {
			return out, 0xFF, fmt.Errorf("LD dst must be register")
		}
		out[1], _ = regNum(dst.Name)
		destReg = out[1]
		// src can be ImmediateNode, IdentifierNode, or MemoryOperandNode
		switch src := instr.Operands[1].(type) {
		case *ImmediateNode:
			imm, _ := parseImmediateOperand(src)
			binary.BigEndian.PutUint16(out[2:], uint16(imm))
		case *IdentifierNode:
			addr := cg.resolveSymbol(src.Name)
			binary.BigEndian.PutUint16(out[2:], uint16(addr))
		case *MemoryOperandNode:
			addr := cg.resolveMemOperand(src)
			binary.BigEndian.PutUint16(out[2:], uint16(addr))
		default:
			return out, 0xFF, fmt.Errorf("LD src must be immediate, label, or memory")
		}
	case 0x11: // ST
		if len(instr.Operands) != 2 {
			return out, 0xFF, fmt.Errorf("ST requires 2 operands at line %d", instr.Line)
		}
		src, ok := instr.Operands[0].(*RegisterNode)
		if !ok {
			return out, 0xFF, fmt.Errorf("ST src must be register")
		}
		out[1], _ = regNum(src.Name)
		// dst can be ImmediateNode, IdentifierNode, or MemoryOperandNode
		switch dst := instr.Operands[1].(type) {
		case *ImmediateNode:
			imm, _ := parseImmediateOperand(dst)
			binary.BigEndian.PutUint16(out[2:], uint16(imm))
		case *IdentifierNode:
			addr := cg.resolveSymbol(dst.Name)
			binary.BigEndian.PutUint16(out[2:], uint16(addr))
		case *MemoryOperandNode:
			addr := cg.resolveMemOperand(dst)
			binary.BigEndian.PutUint16(out[2:], uint16(addr))
		default:
			return out, 0xFF, fmt.Errorf("ST dst must be immediate, label, or memory")
		}
	case 0x20: // JMP
		if len(instr.Operands) != 1 {
			return out, 0xFF, fmt.Errorf("JMP requires 1 operand at line %d", instr.Line)
		}
		addr := cg.resolveOperandAddr(instr.Operands[0])
		binary.BigEndian.PutUint16(out[2:], uint16(addr))
	case 0x21, 0x22, 0x23, 0x24: // BEQ, BNE, BLT, BGT
		if len(instr.Operands) != 3 {
			return out, 0xFF, fmt.Errorf("%s requires 3 operands at line %d", instr.Mnemonic, instr.Line)
		}
		for i := 0; i < 2; i++ {
			reg, ok := instr.Operands[i].(*RegisterNode)
			if !ok {
				return out, 0xFF, fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			rn, err := regNum(reg.Name)
			if err != nil {
				return out, 0xFF, err
			}
			out[i+1] = rn
		}
		addr := cg.resolveOperandAddr(instr.Operands[2])
		binary.BigEndian.PutUint16(out[2:], uint16(addr))
	case 0x30: // WFI
		// nothing more
	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47: // VEC
		if len(instr.Operands) != 3 {
			return out, 0xFF, fmt.Errorf("%s requires 3 operands at line %d", instr.Mnemonic, instr.Line)
		}
		vecRegs := map[string]byte{"VA": 0, "VT": 1, "VB": 2}
		for i := 0; i < 3; i++ {
			reg, ok := instr.Operands[i].(*RegisterNode)
			if !ok {
				return out, 0xFF, fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			vn, ok := vecRegs[strings.ToUpper(reg.Name)]
			if !ok {
				return out, 0xFF, fmt.Errorf("%s operand %d is not a vector register (VA/VT/VB) at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			out[i+1] = vn
		}
	case 0x50, 0x51, 0x52, 0x53, 0x54, 0x55: // FP
		if len(instr.Operands) != 3 {
			return out, 0xFF, fmt.Errorf("%s requires 3 operands at line %d", instr.Mnemonic, instr.Line)
		}
		fpRegs := map[string]byte{"FA": 0, "FT": 1, "FB": 2}
		for i := 0; i < 3; i++ {
			reg, ok := instr.Operands[i].(*RegisterNode)
			if !ok {
				return out, 0xFF, fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			fn, ok := fpRegs[strings.ToUpper(reg.Name)]
			if !ok {
				return out, 0xFF, fmt.Errorf("%s operand %d is not a floating point register (FA/FT/FB) at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			out[i+1] = fn
		}
	case 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D: // COMPLEX
		// Most are binary, some unary (e.g., SQRT, ABS, SIN, COS, TAN, ASIN, ACOS, ATAN, EXP, LOG)
		unaryOps := map[string]bool{"SQRT": true, "ABS": true, "SIN": true, "COS": true, "TAN": true, "ASIN": true, "ACOS": true, "ATAN": true, "EXP": true, "LOG": true}
		if unaryOps[strings.ToUpper(instr.Mnemonic)] {
			if len(instr.Operands) != 2 {
				return out, 0xFF, fmt.Errorf("%s requires 2 operands (dst, src) at line %d", instr.Mnemonic, instr.Line)
			}
			for i := 0; i < 2; i++ {
				reg, ok := instr.Operands[i].(*RegisterNode)
				if !ok {
					return out, 0xFF, fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
				}
				rn, err := regNum(reg.Name)
				if err != nil {
					return out, 0xFF, err
				}
				out[i+1] = rn
			}
			out[3] = 0 // unused
		} else {
			if len(instr.Operands) != 3 {
				return out, 0xFF, fmt.Errorf("%s requires 3 operands (dst, src1, src2) at line %d", instr.Mnemonic, instr.Line)
			}
			for i := 0; i < 3; i++ {
				reg, ok := instr.Operands[i].(*RegisterNode)
				if !ok {
					return out, 0xFF, fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
				}
				rn, err := regNum(reg.Name)
				if err != nil {
					return out, 0xFF, err
				}
				out[i+1] = rn
			}
		}
	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85: // COMPLEX_VEC
		if len(instr.Operands) != 3 {
			return out, 0xFF, fmt.Errorf("%s requires 3 operands at line %d", instr.Mnemonic, instr.Line)
		}
		vecRegs := map[string]byte{"VA": 0, "VT": 1, "VB": 2}
		for i := 0; i < 3; i++ {
			reg, ok := instr.Operands[i].(*RegisterNode)
			if !ok {
				return out, 0xFF, fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			vn, ok := vecRegs[strings.ToUpper(reg.Name)]
			if !ok {
				return out, 0xFF, fmt.Errorf("%s operand %d is not a vector register (VA/VT/VB) at line %d", instr.Mnemonic, i+1, instr.Line)
			}
			out[i+1] = vn
		}
	case 0x90, 0x91, 0x92: // COMPLEX_MEM
		// All bytes except opcode are zero
		// No further action needed
	case 0xA0, 0xA1, 0xA2: // COMPLEX_SYS
		// All bytes except opcode are zero
		// No further action needed
	default:
		return out, 0xFF, fmt.Errorf("unsupported opcode: %02X", opc)
	}
	return out, destReg, nil
}

// --- Directive/Data Emission ---
func (cg *CodeGenerator) emitDirective(dir *DirectiveNode) error {
	fmt.Printf("[DEBUG] In emitDirective: %+v\n", dir)
	if deprecatedDirectives[strings.ToUpper(dir.Name)] && codegenWarnings != nil {
		*codegenWarnings = append(*codegenWarnings, fmt.Errorf("warning: directive '%s' at line %d is deprecated", dir.Name, dir.Line))
	}
	name := strings.ToUpper(dir.Name)
	switch name {
	case ".ORG":
		if len(dir.Params) > 0 {
			imm, _ := parseImmediateOperand(dir.Params[0])
			cg.CurrentAddr = uint32(imm)
		}
	case ".SPACE":
		if len(dir.Params) > 0 {
			imm, _ := parseImmediateOperand(dir.Params[0])
			for i := 0; i < int(imm); i++ {
				cg.Output = append(cg.Output, 0)
				cg.CurrentAddr++
			}
		}
	case ".DW":
		for _, op := range dir.Params {
			switch v := op.(type) {
			case *ImmediateNode:
				if isQuotedString(v.Value) {
					for _, c := range unquoteString(v.Value) {
						var buf [2]byte
						binary.BigEndian.PutUint16(buf[:], uint16(c))
						cg.Output = append(cg.Output, buf[:]...)
						cg.CurrentAddr += 2
					}
				} else {
					imm, _ := parseImmediateOperand(v)
					var buf [2]byte
					binary.BigEndian.PutUint16(buf[:], uint16(imm))
					cg.Output = append(cg.Output, buf[:]...)
					cg.CurrentAddr += 2
				}
			default:
				imm, _ := parseImmediateOperand(op)
				var buf [2]byte
				binary.BigEndian.PutUint16(buf[:], uint16(imm))
				cg.Output = append(cg.Output, buf[:]...)
				cg.CurrentAddr += 2
			}
		}
	case ".DB":
		for _, op := range dir.Params {
			switch v := op.(type) {
			case *ImmediateNode:
				if isQuotedString(v.Value) {
					for _, c := range unquoteString(v.Value) {
						cg.Output = append(cg.Output, byte(c))
						cg.CurrentAddr++
					}
				} else {
					imm, _ := parseImmediateOperand(v)
					cg.Output = append(cg.Output, byte(imm))
					cg.CurrentAddr++
				}
			default:
				imm, _ := parseImmediateOperand(op)
				cg.Output = append(cg.Output, byte(imm))
				cg.CurrentAddr++
			}
		}
	case ".EQU":
		// Already handled in pass 1
		return nil
	case ".INCLUDE":
		return fmt.Errorf(".INCLUDE directive not implemented yet")
	// TODO: Add support for other assembler directives as needed
	default:
		// Ignore other directives for now
	}
	return nil
}

// Helper: check if a string is quoted (e.g., '"Hello"')
func isQuotedString(s string) bool {
	return len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"'
}

// Helper: remove quotes and unescape (basic)
func unquoteString(s string) string {
	if isQuotedString(s) {
		return s[1 : len(s)-1]
	}
	return s
}

// --- Helpers ---
func parseImmediateOperand(op OperandNode) (int64, error) {
	switch v := op.(type) {
	case *ImmediateNode:
		// Support decimal, hex, binary, ternary
		val := v.Value
		if strings.HasPrefix(val, "0x") {
			return strconv.ParseInt(val[2:], 16, 32)
		} else if strings.HasPrefix(val, "0b") {
			return strconv.ParseInt(val[2:], 2, 32)
		} else if strings.HasPrefix(val, "0t") {
			// Ternary: +, -, 0
			// For now, treat as 0
			return 0, nil
		} else {
			return strconv.ParseInt(val, 10, 32)
		}
	case *IdentifierNode:
		// Symbol/label reference
		return 0, nil // Will be resolved in pass 2
	}
	return 0, fmt.Errorf("unsupported immediate operand")
}

func (cg *CodeGenerator) resolveSymbol(name string) uint32 {
	if v, ok := cg.Labels[name]; ok {
		return v
	}
	if v, ok := cg.Equs[name]; ok {
		return v
	}
	return 0
}

func (cg *CodeGenerator) resolveOperandAddr(op OperandNode) uint32 {
	switch v := op.(type) {
	case *ImmediateNode:
		imm, _ := parseImmediateOperand(v)
		return uint32(imm)
	case *IdentifierNode:
		return cg.resolveSymbol(v.Name)
	}
	return 0
}

func (cg *CodeGenerator) resolveMemOperand(mem *MemoryOperandNode) uint32 {
	// For now, just use base register as address (stub)
	if mem.Base != "" {
		if addr, ok := cg.Labels[mem.Base]; ok {
			return addr
		}
	}
	return 0
}
