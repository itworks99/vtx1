package rvsim

// Real RV32IM machine-code encoding for the structural instruction list, so
// the same kernel definitions that drive the interpreter also produce the
// binary image executed by the PicoRV32 RTL in the cycle-level comparison.
// Branch/jump targets are instruction indices; code is placed contiguously,
// so the byte offset is 4*(target - pc).

import "fmt"

func enc(op Op, i Instr, pc int) (uint32, error) {
	r := func(x int) uint32 { return uint32(x) & 31 }
	rd, rs1, rs2 := r(i.Rd), r(i.Rs1), r(i.Rs2)
	imm := uint32(i.Imm)

	iType := func(opc, f3 uint32) uint32 {
		return (imm&0xFFF)<<20 | rs1<<15 | f3<<12 | rd<<7 | opc
	}
	rType := func(f7, f3 uint32) uint32 {
		return f7<<25 | rs2<<20 | rs1<<15 | f3<<12 | rd<<7 | 0x33
	}
	bType := func(f3 uint32) uint32 {
		off := uint32(4 * (i.Target - pc))
		return (off>>12&1)<<31 | (off>>5&0x3F)<<25 | rs2<<20 | rs1<<15 |
			f3<<12 | (off>>1&0xF)<<8 | (off>>11&1)<<7 | 0x63
	}

	switch op {
	case LUI:
		return (imm&0xFFFFF)<<12 | rd<<7 | 0x37, nil
	case ADDI:
		return iType(0x13, 0), nil
	case SRAI:
		return 0x20<<25 | (imm&31)<<20 | rs1<<15 | 5<<12 | rd<<7 | 0x13, nil
	case ADD:
		return rType(0x00, 0), nil
	case SUB:
		return rType(0x20, 0), nil
	case SLT:
		return rType(0x00, 2), nil
	case XOR:
		return rType(0x00, 4), nil
	case MUL:
		return rType(0x01, 0), nil
	case LW:
		return iType(0x03, 2), nil
	case SW:
		return (imm>>5&0x7F)<<25 | rs2<<20 | rs1<<15 | 2<<12 |
			(imm&0x1F)<<7 | 0x23, nil
	case BEQ:
		return bType(0), nil
	case BNE:
		return bType(1), nil
	case BLT:
		return bType(4), nil
	case BGE:
		return bType(5), nil
	case JAL:
		off := uint32(4 * (i.Target - pc))
		return (off>>20&1)<<31 | (off>>1&0x3FF)<<21 | (off>>11&1)<<20 |
			(off>>12&0xFF)<<12 | rd<<7 | 0x6F, nil
	case HALT:
		return 0x00100073, nil // EBREAK
	}
	return 0, fmt.Errorf("encode: unimplemented op %s", op)
}

// Encode produces the RV32IM machine words for a built program.
func Encode(prog []Instr) ([]uint32, error) {
	words := make([]uint32, len(prog))
	for pc, ins := range prog {
		w, err := enc(ins.Op, ins, pc)
		if err != nil {
			return nil, fmt.Errorf("pc %d: %w", pc, err)
		}
		words[pc] = w
	}
	return words, nil
}
