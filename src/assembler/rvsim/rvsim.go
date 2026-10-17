// Package rvsim is a minimal RV32IM interpreter: the binary control core for
// the Tier A architecture comparison. It mirrors vtxsim's instrumentation
// exactly (same Stats structure, same counting rules), so the two columns of
// the benchmark table are measured with the same stick.
//
// The interpreter is structural: programs are built from typed instructions
// rather than encoded/decoded 32-bit words. That changes nothing for the
// dynamic metrics, and static code size is exact anyway because every RV32I
// instruction occupies exactly 32 bits. Pseudo-instructions are expanded the
// way a real assembler expands them (LI of a constant outside the 12-bit
// immediate range becomes LUI+ADDI and counts as 2 instructions).
package rvsim

import (
	"fmt"

	"github.com/kvany/vtx1/assembler/sim"
)

// Op is an RV32IM operation.
type Op string

const (
	ADDI  Op = "addi"
	LUI   Op = "lui"
	ADD   Op = "add"
	SUB   Op = "sub"
	MUL   Op = "mul" // M extension (VTX1 has MUL, so the control core gets it too)
	SLT   Op = "slt"
	XOR   Op = "xor"
	SRAI  Op = "srai"
	LW    Op = "lw"
	SW    Op = "sw"
	BEQ   Op = "beq"
	BNE   Op = "bne"
	BLT   Op = "blt"
	BGE   Op = "bge"
	JAL   Op = "jal" // used as unconditional jump (rd = x0)
	HALT  Op = "ebreak"
)

// Instr is one instruction. Branch/jump targets are instruction indices
// (resolved from labels by the Builder).
type Instr struct {
	Op            Op
	Rd, Rs1, Rs2  int
	Imm           int32
	Target        int
}

func category(op Op) string {
	switch op {
	case LW, SW:
		return "memory"
	case BEQ, BNE, BLT, BGE, JAL:
		return "control"
	default:
		return "ALU"
	}
}

// Machine is one RV32IM hart with word-addressed data memory.
type Machine struct {
	Reg    [32]int32
	PC     int
	Prog   []Instr
	Mem    map[uint32]int32 // byte addresses, word-aligned accesses only
	Halted bool
	Stats  sim.Stats
}

func New(prog []Instr) *Machine {
	return &Machine{
		Prog: prog,
		Mem:  make(map[uint32]int32),
		Stats: sim.Stats{
			ByCategory: make(map[string]uint64),
			ByMnemonic: make(map[string]uint64),
		},
	}
}

// WriteWord / ReadWord access data memory (byte addresses, word aligned).
func (m *Machine) WriteWord(addr uint32, v int32) { m.Mem[addr] = v }
func (m *Machine) ReadWord(addr uint32) int32     { return m.Mem[addr] }

func (m *Machine) get(r int) int32 {
	if r == 0 {
		return 0
	}
	return m.Reg[r]
}

func (m *Machine) set(r int, v int32) {
	if r != 0 {
		m.Reg[r] = v
	}
}

// Step executes one instruction.
func (m *Machine) Step() error {
	if m.PC < 0 || m.PC >= len(m.Prog) {
		return fmt.Errorf("PC %d outside program", m.PC)
	}
	ins := m.Prog[m.PC]
	m.Stats.Instructions++
	m.Stats.ByCategory[category(ins.Op)]++
	m.Stats.ByMnemonic[string(ins.Op)]++
	next := m.PC + 1

	switch ins.Op {
	case ADDI:
		m.set(ins.Rd, m.get(ins.Rs1)+ins.Imm)
	case LUI:
		m.set(ins.Rd, ins.Imm<<12)
	case ADD:
		m.set(ins.Rd, m.get(ins.Rs1)+m.get(ins.Rs2))
	case SUB:
		m.set(ins.Rd, m.get(ins.Rs1)-m.get(ins.Rs2))
	case MUL:
		m.set(ins.Rd, m.get(ins.Rs1)*m.get(ins.Rs2))
	case SLT:
		if m.get(ins.Rs1) < m.get(ins.Rs2) {
			m.set(ins.Rd, 1)
		} else {
			m.set(ins.Rd, 0)
		}
	case XOR:
		m.set(ins.Rd, m.get(ins.Rs1)^m.get(ins.Rs2))
	case SRAI:
		m.set(ins.Rd, m.get(ins.Rs1)>>uint(ins.Imm))
	case LW:
		m.Stats.Loads++
		addr := uint32(m.get(ins.Rs1) + ins.Imm)
		if addr%4 != 0 {
			return fmt.Errorf("PC %d: misaligned load at 0x%X", m.PC, addr)
		}
		m.set(ins.Rd, m.ReadWord(addr))
	case SW:
		m.Stats.Stores++
		addr := uint32(m.get(ins.Rs1) + ins.Imm)
		if addr%4 != 0 {
			return fmt.Errorf("PC %d: misaligned store at 0x%X", m.PC, addr)
		}
		m.WriteWord(addr, m.get(ins.Rs2))
	case BEQ, BNE, BLT, BGE:
		m.Stats.Branches++
		var take bool
		switch ins.Op {
		case BEQ:
			take = m.get(ins.Rs1) == m.get(ins.Rs2)
		case BNE:
			take = m.get(ins.Rs1) != m.get(ins.Rs2)
		case BLT:
			take = m.get(ins.Rs1) < m.get(ins.Rs2)
		case BGE:
			take = m.get(ins.Rs1) >= m.get(ins.Rs2)
		}
		if take {
			m.Stats.BranchesTaken++
			next = ins.Target
		}
	case JAL:
		m.set(ins.Rd, int32(m.PC+1))
		next = ins.Target
	case HALT:
		m.Halted = true
	default:
		return fmt.Errorf("PC %d: unimplemented op %s", m.PC, ins.Op)
	}

	m.PC = next
	return nil
}

// Run executes until EBREAK or maxSteps instructions.
func (m *Machine) Run(maxSteps uint64) error {
	for !m.Halted {
		if m.Stats.Instructions >= maxSteps {
			return fmt.Errorf("stopped after %d instructions without reaching ebreak", maxSteps)
		}
		if err := m.Step(); err != nil {
			return err
		}
	}
	return nil
}

// --- Program builder with labels and honest pseudo-instruction expansion ---

type Builder struct {
	prog    []Instr
	labels  map[string]int
	fixups  []fixup
}

type fixup struct {
	index int
	label string
}

func NewBuilder() *Builder {
	return &Builder{labels: make(map[string]int)}
}

func (b *Builder) emit(i Instr) { b.prog = append(b.prog, i) }

func (b *Builder) Label(name string) { b.labels[name] = len(b.prog) }

func (b *Builder) branch(op Op, rs1, rs2 int, label string) {
	b.fixups = append(b.fixups, fixup{len(b.prog), label})
	b.emit(Instr{Op: op, Rs1: rs1, Rs2: rs2})
}

// Li loads a constant, expanding to LUI+ADDI when it does not fit the 12-bit
// signed immediate, exactly as a real assembler would (2 instructions).
func (b *Builder) Li(rd int, v int32) {
	if v >= -2048 && v <= 2047 {
		b.emit(Instr{Op: ADDI, Rd: rd, Rs1: 0, Imm: v})
		return
	}
	upper := (v + 0x800) >> 12
	lower := v - (upper << 12)
	b.emit(Instr{Op: LUI, Rd: rd, Imm: upper})
	if lower != 0 {
		b.emit(Instr{Op: ADDI, Rd: rd, Rs1: rd, Imm: lower})
	}
}

func (b *Builder) Addi(rd, rs1 int, imm int32) { b.emit(Instr{Op: ADDI, Rd: rd, Rs1: rs1, Imm: imm}) }
func (b *Builder) Add(rd, rs1, rs2 int)        { b.emit(Instr{Op: ADD, Rd: rd, Rs1: rs1, Rs2: rs2}) }
func (b *Builder) Sub(rd, rs1, rs2 int)        { b.emit(Instr{Op: SUB, Rd: rd, Rs1: rs1, Rs2: rs2}) }
func (b *Builder) Mul(rd, rs1, rs2 int)        { b.emit(Instr{Op: MUL, Rd: rd, Rs1: rs1, Rs2: rs2}) }
func (b *Builder) Slt(rd, rs1, rs2 int)        { b.emit(Instr{Op: SLT, Rd: rd, Rs1: rs1, Rs2: rs2}) }
func (b *Builder) Xor(rd, rs1, rs2 int)        { b.emit(Instr{Op: XOR, Rd: rd, Rs1: rs1, Rs2: rs2}) }
func (b *Builder) Srai(rd, rs1 int, sh int32)  { b.emit(Instr{Op: SRAI, Rd: rd, Rs1: rs1, Imm: sh}) }
func (b *Builder) Lw(rd, rs1 int, off int32)   { b.emit(Instr{Op: LW, Rd: rd, Rs1: rs1, Imm: off}) }
func (b *Builder) Sw(rs2, rs1 int, off int32)  { b.emit(Instr{Op: SW, Rs1: rs1, Rs2: rs2, Imm: off}) }
func (b *Builder) Beq(rs1, rs2 int, l string)  { b.branch(BEQ, rs1, rs2, l) }
func (b *Builder) Bne(rs1, rs2 int, l string)  { b.branch(BNE, rs1, rs2, l) }
func (b *Builder) Blt(rs1, rs2 int, l string)  { b.branch(BLT, rs1, rs2, l) }
func (b *Builder) Bge(rs1, rs2 int, l string)  { b.branch(BGE, rs1, rs2, l) }
func (b *Builder) J(l string) {
	b.fixups = append(b.fixups, fixup{len(b.prog), l})
	b.emit(Instr{Op: JAL, Rd: 0})
}
func (b *Builder) Ebreak() { b.emit(Instr{Op: HALT}) }

// Build resolves labels and returns the program.
func (b *Builder) Build() ([]Instr, error) {
	for _, f := range b.fixups {
		target, ok := b.labels[f.label]
		if !ok {
			return nil, fmt.Errorf("undefined label %q", f.label)
		}
		b.prog[f.index].Target = target
	}
	return b.prog, nil
}
