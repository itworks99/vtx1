package rvsim

// The Tier A benchmark kernels, hand-written in idiomatic RV32IM, mirroring
// src/assembler/benchmarks/*.asm: same algorithms, same 16-element data.
// Both sides are hand-written assembly, so neither gets a compiler handicap.
//
// Data layout (byte addresses; RV32 is byte-addressed with 32-bit words):
//   0x1100 vec_a[16] / data[16]
//   0x1140 vec_b[16]
//   0x1200 results

const (
	addrA      = 0x1100
	addrB      = 0x1140
	addrResult = 0x1200
)

// Register conventions (RISC-V ABI names in comments)
const (
	x0 = 0
	a0 = 10 // pointer A
	a1 = 11 // pointer B
	a2 = 12 // end pointer
	a3 = 13 // result pointer
	s0 = 8  // accumulator / neg count
	s1 = 9  // zero count
	s2 = 18 // pos count
	t0 = 5
	t1 = 6
	t2 = 7
	t3 = 28
)

// VecA and VecB match the VTX1 benchmark data exactly.
var VecA = []int32{3, -5, 7, 2, -8, 4, 1, -2, 9, -3, 6, -7, 5, 0, 2, -4}
var VecB = []int32{2, 3, -1, 5, 2, -6, 8, 4, -2, 7, 1, -3, -5, 9, 0, 6}
var SignData = []int32{5, -3, 0, 7, -1, 0, 2, -9, 4, 0, -6, 1, -2, 8, 0, -5}

// LoadVectors places the two input vectors in memory.
func LoadVectors(m *Machine) {
	for i, v := range VecA {
		m.WriteWord(uint32(addrA+4*i), v)
	}
	for i, v := range VecB {
		m.WriteWord(uint32(addrB+4*i), v)
	}
}

// DotProduct: result = sum(a[i]*b[i]). Expected: -107.
func DotProduct() ([]Instr, error) {
	b := NewBuilder()
	b.Li(a0, addrA)
	b.Li(a1, addrB)
	b.Li(a2, addrA+4*16) // end of A
	b.Li(s0, 0)          // acc
	b.Label("loop")
	b.Lw(t0, a0, 0)
	b.Lw(t1, a1, 0)
	b.Mul(t0, t0, t1)
	b.Add(s0, s0, t0)
	b.Addi(a0, a0, 4)
	b.Addi(a1, a1, 4)
	b.Blt(a0, a2, "loop")
	b.Li(a3, addrResult)
	b.Sw(s0, a3, 0)
	b.Ebreak()
	return b.Build()
}

// SignCount: classify elements by sign with two branches.
// Expected: 6 negative, 4 zero, 6 positive at result+0/4/8.
func SignCount() ([]Instr, error) {
	b := NewBuilder()
	b.Li(a0, addrA)
	b.Li(a2, addrA+4*16)
	b.Li(s0, 0) // neg
	b.Li(s1, 0) // zero
	b.Li(s2, 0) // pos
	b.Label("loop")
	b.Lw(t0, a0, 0)
	b.Blt(t0, x0, "neg")
	b.Blt(x0, t0, "pos")
	b.Addi(s1, s1, 1)
	b.J("next")
	b.Label("neg")
	b.Addi(s0, s0, 1)
	b.J("next")
	b.Label("pos")
	b.Addi(s2, s2, 1)
	b.Label("next")
	b.Addi(a0, a0, 4)
	b.Blt(a0, a2, "loop")
	b.Li(a3, addrResult)
	b.Sw(s0, a3, 0)
	b.Sw(s1, a3, 4)
	b.Sw(s2, a3, 8)
	b.Ebreak()
	return b.Build()
}

// LoadSignData places the sign_count input in memory (at addrA).
func LoadSignData(m *Machine) {
	for i, v := range SignData {
		m.WriteWord(uint32(addrA+4*i), v)
	}
}

// AbsSum: result = sum(|a[i]-b[i]|), branchless abs (srai/xor/sub).
// Expected: 114.
func AbsSum() ([]Instr, error) {
	b := NewBuilder()
	b.Li(a0, addrA)
	b.Li(a1, addrB)
	b.Li(a2, addrA+4*16)
	b.Li(s0, 0)
	b.Label("loop")
	b.Lw(t0, a0, 0)
	b.Lw(t1, a1, 0)
	b.Sub(t0, t0, t1)
	b.Srai(t2, t0, 31) // branchless abs: t2 = sign mask
	b.Xor(t0, t0, t2)
	b.Sub(t0, t0, t2)
	b.Add(s0, s0, t0)
	b.Addi(a0, a0, 4)
	b.Addi(a1, a1, 4)
	b.Blt(a0, a2, "loop")
	b.Li(a3, addrResult)
	b.Sw(s0, a3, 0)
	b.Ebreak()
	return b.Build()
}

// SignSum: result = sum(sign(a[i]-b[i])). The three-valued compare costs
// slt+slt+sub on a binary ISA versus one CMP on VTX1. Expected: -2.
func SignSum() ([]Instr, error) {
	b := NewBuilder()
	b.Li(a0, addrA)
	b.Li(a1, addrB)
	b.Li(a2, addrA+4*16)
	b.Li(s0, 0)
	b.Label("loop")
	b.Lw(t0, a0, 0)
	b.Lw(t1, a1, 0)
	b.Slt(t2, t0, t1)  // a < b
	b.Slt(t3, t1, t0)  // b < a
	b.Sub(t2, t3, t2)  // sign(a-b)
	b.Add(s0, s0, t2)
	b.Addi(a0, a0, 4)
	b.Addi(a1, a1, 4)
	b.Blt(a0, a2, "loop")
	b.Li(a3, addrResult)
	b.Sw(s0, a3, 0)
	b.Ebreak()
	return b.Build()
}
