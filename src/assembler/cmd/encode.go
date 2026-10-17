package cmd

import (
	"fmt"
	"strings"

	"github.com/kvany/vtx1/assembler/isa"
)

// Trit-native operation encoding
// ==============================
// Every operation is encoded into one 18-trit machine word (stored as 5
// bytes: 36 bits, 2 bits per trit, 4 padding bits). A VLIW bundle is three
// consecutive words. The field layout is defined in package isa, which is
// the single source of truth shared with the generated Verilog constants.

// Slot is one instruction word as trits (index 0 = most significant).
type Slot [isa.SlotTrits]int

// Bytes returns the stored form of the slot (5 bytes, 36 bits + 4 zero bits).
func (s *Slot) Bytes() []byte {
	return isa.EncodeTrits(s[:])
}

// setField writes the balanced-ternary encoding of v into trits [pos, pos+n).
func (s *Slot) setField(pos, n int, v int64) error {
	trits, err := isa.BalancedTrits(v, n)
	if err != nil {
		return err
	}
	copy(s[pos:pos+n], trits)
	return nil
}

func (s *Slot) setOpcode(value int) {
	s.setField(0, isa.OpcodeTrits, int64(value-isa.OpcodeBias))
}

func (s *Slot) setReg(pos int, index int) {
	s.setField(pos, isa.RegTrits, int64(index-isa.RegBias))
}

// opInfo describes an encoded operation for VLIW packing validation.
type opInfo struct {
	Mnemonic string
	Category isa.Category
	Writes   []string // canonical register names written by this op
	Reads    []string // canonical register names read by this op
}

// lookupReg resolves a register name and enforces its class.
func lookupReg(name string, cls isa.RegClass, instr *InstructionNode, operand int) (*isa.Register, error) {
	r, ok := isa.RegByName[strings.ToUpper(name)]
	if !ok {
		return nil, fmt.Errorf("%s operand %d: unknown register %s at line %d", instr.Mnemonic, operand, name, instr.Line)
	}
	if r.Class != cls && r.Index != isa.ZeroRegIndex {
		want := map[isa.RegClass]string{
			isa.ClassGPR: "a general purpose register",
			isa.ClassVec: "a vector register (VA/VT/VB)",
			isa.ClassFP:  "a floating point register (FA/FT/FB)",
		}[cls]
		return nil, fmt.Errorf("%s operand %d: %s is not %s at line %d", instr.Mnemonic, operand, name, want, instr.Line)
	}
	return r, nil
}

// operandReg extracts a register operand of the given class.
func operandReg(instr *InstructionNode, idx int, cls isa.RegClass) (*isa.Register, error) {
	reg, ok := instr.Operands[idx].(*RegisterNode)
	if !ok {
		return nil, fmt.Errorf("%s operand %d is not a register at line %d", instr.Mnemonic, idx+1, instr.Line)
	}
	return lookupReg(reg.Name, cls, instr, idx+1)
}

func requireOperands(instr *InstructionNode, n int) error {
	if len(instr.Operands) != n {
		return fmt.Errorf("%s requires %d operand(s), got %d at line %d",
			instr.Mnemonic, n, len(instr.Operands), instr.Line)
	}
	return nil
}

// encodeOperation encodes one operation into an 18-trit slot. wordAddr is the
// word address the operation will execute at (the bundle base for VLIW slots),
// used for PC-relative branch offsets. It also reports the registers the
// operation reads and writes for VLIW conflict detection.
func (cg *CodeGenerator) encodeOperation(instr *InstructionNode, wordAddr uint32) (Slot, opInfo, error) {
	var slot Slot
	mnem := strings.ToUpper(instr.Mnemonic)
	info := opInfo{Mnemonic: mnem}

	ins, ok := isa.ByName[mnem]
	if !ok {
		return slot, info, fmt.Errorf("unsupported instruction: %s at line %d", instr.Mnemonic, instr.Line)
	}
	slot.setOpcode(ins.Value)
	info.Category = ins.Category

	switch ins.Format {
	case isa.FormatNone:
		if err := requireOperands(instr, 0); err != nil {
			return slot, info, err
		}

	case isa.FormatR1:
		if err := requireOperands(instr, 1); err != nil {
			return slot, info, err
		}
		r, err := operandReg(instr, 0, ins.Class)
		if err != nil {
			return slot, info, err
		}
		slot.setReg(4, r.Index)
		switch mnem {
		case "INC", "DEC":
			info.Writes = append(info.Writes, r.Name)
			info.Reads = append(info.Reads, r.Name)
		case "POP":
			info.Writes = append(info.Writes, r.Name)
		default: // PUSH, JR, JALR
			info.Reads = append(info.Reads, r.Name)
		}

	case isa.FormatR2:
		if err := requireOperands(instr, 2); err != nil {
			return slot, info, err
		}
		rd, err := operandReg(instr, 0, ins.Class)
		if err != nil {
			return slot, info, err
		}
		rs, err := operandReg(instr, 1, ins.Class)
		if err != nil {
			return slot, info, err
		}
		slot.setReg(4, rd.Index)
		slot.setReg(7, rs.Index)
		info.Writes = append(info.Writes, rd.Name)
		info.Reads = append(info.Reads, rs.Name)

	case isa.FormatR3:
		if err := requireOperands(instr, 3); err != nil {
			return slot, info, err
		}
		rd, err := operandReg(instr, 0, ins.Class)
		if err != nil {
			return slot, info, err
		}
		rs1, err := operandReg(instr, 1, ins.Class)
		if err != nil {
			return slot, info, err
		}
		slot.setReg(4, rd.Index)
		slot.setReg(7, rs1.Index)
		info.Writes = append(info.Writes, rd.Name)
		info.Reads = append(info.Reads, rs1.Name)
		// Third operand: register (mode 0), or for GPR-class ops a 7-trit
		// balanced immediate (mode +1).
		if reg, ok := instr.Operands[2].(*RegisterNode); ok {
			rs2, err := lookupReg(reg.Name, ins.Class, instr, 3)
			if err != nil {
				return slot, info, err
			}
			slot[10] = 0
			slot.setReg(15, rs2.Index)
			info.Reads = append(info.Reads, rs2.Name)
		} else if ins.Class == isa.ClassGPR {
			imm, err := parseImmediateOperand(instr.Operands[2])
			if err != nil {
				return slot, info, fmt.Errorf("%s operand 3 must be a register or immediate at line %d: %v", instr.Mnemonic, instr.Line, err)
			}
			slot[10] = 1
			if err := slot.setField(11, 7, imm); err != nil {
				return slot, info, fmt.Errorf("%s immediate at line %d: %v", instr.Mnemonic, instr.Line, err)
			}
		} else {
			return slot, info, fmt.Errorf("%s operand 3 is not a register at line %d", instr.Mnemonic, instr.Line)
		}

	case isa.FormatMem:
		if err := requireOperands(instr, 2); err != nil {
			return slot, info, err
		}
		r, err := operandReg(instr, 0, ins.Class)
		if err != nil {
			return slot, info, err
		}
		slot.setReg(4, r.Index)
		memRegs, err := cg.encodeMemField(&slot, instr.Operands[1], instr)
		if err != nil {
			return slot, info, err
		}
		// Loads write the register; stores read it. Address registers are reads.
		switch mnem {
		case "ST", "VST", "FST":
			info.Reads = append(info.Reads, r.Name)
		default: // LD, VLD, FLD, LEA
			info.Writes = append(info.Writes, r.Name)
		}
		info.Reads = append(info.Reads, memRegs...)

	case isa.FormatBranch:
		if err := requireOperands(instr, 3); err != nil {
			return slot, info, err
		}
		rs1, err := operandReg(instr, 0, isa.ClassGPR)
		if err != nil {
			return slot, info, err
		}
		slot.setReg(4, rs1.Index)
		info.Reads = append(info.Reads, rs1.Name)
		if reg, ok := instr.Operands[1].(*RegisterNode); ok {
			rs2, err := lookupReg(reg.Name, isa.ClassGPR, instr, 2)
			if err != nil {
				return slot, info, err
			}
			slot.setReg(7, rs2.Index)
			info.Reads = append(info.Reads, rs2.Name)
		} else {
			imm, err := parseImmediateOperand(instr.Operands[1])
			if err != nil || imm != 0 {
				return slot, info, fmt.Errorf("%s operand 2 must be a register or the immediate 0 at line %d", instr.Mnemonic, instr.Line)
			}
			slot.setReg(7, isa.ZeroRegIndex) // TZ: compare against zero
		}
		target := int64(cg.resolveOperandAddr(instr.Operands[2]))
		rel := target - int64(wordAddr)
		if err := slot.setField(10, 8, rel); err != nil {
			return slot, info, fmt.Errorf("%s branch offset at line %d: target is %d words away, range is +-%d", instr.Mnemonic, instr.Line, rel, isa.MaxRel8)
		}

	case isa.FormatJump:
		if err := requireOperands(instr, 1); err != nil {
			return slot, info, err
		}
		target := int64(cg.resolveOperandAddr(instr.Operands[0]))
		if err := slot.setField(4, 14, target); err != nil {
			return slot, info, fmt.Errorf("%s target at line %d: %v", instr.Mnemonic, instr.Line, err)
		}
	}

	return slot, info, nil
}

// encodeMemField encodes the address operand of a memory instruction
// (mode trit at position 7, 10-trit payload at 8..17) and returns the
// register names it reads.
func (cg *CodeGenerator) encodeMemField(slot *Slot, op OperandNode, instr *InstructionNode) ([]string, error) {
	baseReg := func(name, role string) (*isa.Register, error) {
		r, ok := isa.RegByName[strings.ToUpper(name)]
		if !ok || r.Class != isa.ClassGPR {
			return nil, fmt.Errorf("%s memory operand %s: unknown register %s at line %d", instr.Mnemonic, role, name, instr.Line)
		}
		return r, nil
	}

	switch v := op.(type) {
	case *ImmediateNode:
		imm, err := parseImmediateOperand(v)
		if err != nil {
			return nil, fmt.Errorf("%s immediate: %v at line %d", instr.Mnemonic, err, instr.Line)
		}
		slot[7] = 0
		if err := slot.setField(8, 10, imm); err != nil {
			return nil, fmt.Errorf("%s immediate at line %d: %v", instr.Mnemonic, instr.Line, err)
		}
		return nil, nil
	case *IdentifierNode, *BinaryOpNode:
		addr := int64(cg.resolveOperandAddr(op))
		slot[7] = 0
		if err := slot.setField(8, 10, addr); err != nil {
			return nil, fmt.Errorf("%s direct address at line %d: %v", instr.Mnemonic, instr.Line, err)
		}
		return nil, nil
	case *MemoryOperandNode:
		base, err := baseReg(v.Base, "base")
		if err != nil {
			return nil, err
		}
		reads := []string{base.Name}
		switch {
		case v.Index != "":
			index, err := baseReg(v.Index, "index")
			if err != nil {
				return nil, err
			}
			reads = append(reads, index.Name)
			slot[7] = 1 // base+index
			slot.setReg(8, base.Index)
			slot.setReg(11, index.Index)
		case v.Offset != "":
			off, err := parseImmediateString(v.Offset)
			if err != nil {
				return nil, fmt.Errorf("%s memory offset '%s': %v at line %d", instr.Mnemonic, v.Offset, err, instr.Line)
			}
			slot[7] = -1 // base+offset
			slot.setReg(8, base.Index)
			if err := slot.setField(11, 7, off); err != nil {
				return nil, fmt.Errorf("%s memory offset at line %d: %v", instr.Mnemonic, instr.Line, err)
			}
		default:
			slot[7] = 1 // [base] == base+index with index TZ (zero)
			slot.setReg(8, base.Index)
			slot.setReg(11, isa.ZeroRegIndex)
		}
		return reads, nil
	}
	return nil, fmt.Errorf("%s address operand must be immediate, label, or memory operand at line %d", instr.Mnemonic, instr.Line)
}
