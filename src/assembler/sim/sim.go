// Package sim is the VTX1 ISA-level simulator: it executes the assembler's
// memh images word by word with balanced-ternary semantics and collects the
// dynamic metrics used for architecture evaluation (instruction counts by
// category, memory traffic, branch behavior).
//
// It models architecture, not microarchitecture: one instruction word per
// step, no pipeline, no cycle timing. VLIW bundles execute as their slots in
// sequence, which matches parallel semantics for any bundle the assembler
// accepts (the assembler rejects intra-bundle write conflicts and warns on
// read-after-write hazards).
package sim

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/kvany/vtx1/assembler/isa"
)

// wordMod is 3^18, the number of values an 18-trit word can hold.
const wordMod = int64(387420489)

// WrapWord reduces v to the balanced 18-trit range ±(3^18-1)/2, the way the
// hardware register file and memory would.
func WrapWord(v int64) int64 {
	r := ((v % wordMod) + wordMod) % wordMod
	if r > (wordMod-1)/2 {
		r -= wordMod
	}
	return r
}

// Stats holds the dynamic execution metrics.
type Stats struct {
	Instructions  uint64
	ByCategory    map[string]uint64
	ByMnemonic    map[string]uint64
	Loads, Stores uint64
	Branches      uint64
	BranchesTaken uint64
}

// Machine is one VTX1 processor with word-addressed memory.
type Machine struct {
	Mem    map[uint32]uint64 // raw 36-bit words (2-bit trit encoding)
	Reg    [19]int64         // register file; index 18 is TZ (always 0)
	PC     uint32
	Halted bool
	Stats  Stats
	Trace  io.Writer // optional: one line per executed instruction
}

func New() *Machine {
	return &Machine{
		Mem: make(map[uint32]uint64),
		Stats: Stats{
			ByCategory: make(map[string]uint64),
			ByMnemonic: make(map[string]uint64),
		},
	}
}

// --- memh image loading ---

// LoadMemh reads a $readmemh image produced by the assembler (-f memh):
// 9-hex-digit 36-bit words with @ADDR markers. It returns the first loaded
// address, the default entry point.
func (m *Machine) LoadMemh(r io.Reader) (uint32, error) {
	scanner := bufio.NewScanner(r)
	addr := uint32(0)
	first := uint32(0)
	haveFirst := false
	for scanner.Scan() {
		line := scanner.Text()
		if i := strings.Index(line, "//"); i >= 0 {
			line = line[:i]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "@") {
			v, err := strconv.ParseUint(line[1:], 16, 32)
			if err != nil {
				return 0, fmt.Errorf("bad address marker %q: %v", line, err)
			}
			addr = uint32(v)
			continue
		}
		v, err := strconv.ParseUint(line, 16, 64)
		if err != nil {
			return 0, fmt.Errorf("bad word %q: %v", line, err)
		}
		m.Mem[addr] = v & 0xFFFFFFFFF
		if !haveFirst {
			first, haveFirst = addr, true
		}
		addr++
	}
	return first, scanner.Err()
}

// --- raw word <-> trits/values ---

func rawToTrits(raw uint64) []int {
	b := []byte{
		byte(raw >> 28), byte(raw >> 20), byte(raw >> 12), byte(raw >> 4),
		byte(raw&0xF) << 4,
	}
	return isa.DecodeTrits(b, isa.SlotTrits)
}

func valueToRaw(v int64) uint64 {
	trits, _ := isa.BalancedTrits(WrapWord(v), isa.SlotTrits)
	b := isa.EncodeTrits(trits)
	return (uint64(b[0])<<32 | uint64(b[1])<<24 | uint64(b[2])<<16 |
		uint64(b[3])<<8 | uint64(b[4])) >> 4
}

func fieldVal(trits []int, pos, n int) int64 { return isa.TritsValue(trits[pos : pos+n]) }
func fieldReg(trits []int, pos int) int      { return int(fieldVal(trits, pos, isa.RegTrits)) + isa.RegBias }

// --- register access (TZ reads 0, writes discarded) ---

func (m *Machine) getReg(idx int) (int64, error) {
	if idx < 0 || idx >= len(m.Reg) {
		return 0, fmt.Errorf("register index %d out of range", idx)
	}
	if idx == isa.ZeroRegIndex {
		return 0, nil
	}
	return m.Reg[idx], nil
}

func (m *Machine) setReg(idx int, v int64) error {
	if idx < 0 || idx >= len(m.Reg) {
		return fmt.Errorf("register index %d out of range", idx)
	}
	if idx != isa.ZeroRegIndex {
		m.Reg[idx] = WrapWord(v)
	}
	return nil
}

// --- memory access ---

// ReadWord returns the value of the word at addr (0 when never written).
func (m *Machine) ReadWord(addr uint32) int64 {
	raw, ok := m.Mem[addr]
	if !ok {
		return 0
	}
	return isa.TritsValue(rawToTrits(raw))
}

func (m *Machine) writeWord(addr uint32, v int64) {
	m.Mem[addr] = valueToRaw(v)
}

// --- tritwise logic (balanced ternary conventions) ---

func tritwise(a, b int64, op func(x, y int) int) int64 {
	ta, _ := isa.BalancedTrits(WrapWord(a), isa.SlotTrits)
	tb, _ := isa.BalancedTrits(WrapWord(b), isa.SlotTrits)
	out := make([]int, isa.SlotTrits)
	for i := range out {
		out[i] = op(ta[i], tb[i])
	}
	return isa.TritsValue(out)
}

func tritMin(x, y int) int { if x < y { return x }; return y }
func tritMax(x, y int) int { if x > y { return x }; return y }
func tritMul(x, y int) int { return x * y }

func shiftTrits(v int64, n int64, left bool, rotate bool) int64 {
	trits, _ := isa.BalancedTrits(WrapWord(v), isa.SlotTrits)
	k := int(((n % 18) + 18) % 18)
	out := make([]int, isa.SlotTrits)
	for i := range out {
		var src int
		if left {
			src = i + k
		} else {
			src = i - k
		}
		switch {
		case src >= 0 && src < isa.SlotTrits:
			out[i] = trits[src]
		case rotate:
			out[i] = trits[((src%18)+18)%18]
		default:
			out[i] = 0
		}
	}
	return isa.TritsValue(out)
}

func sign(v int64) int64 {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

func isqrt(v int64) int64 {
	if v < 0 {
		return 0
	}
	r := int64(0)
	for (r+1)*(r+1) <= v {
		r++
	}
	return r
}

// --- execution ---

// Step executes one instruction word. It returns an error on traps
// (undefined trits, unsupported operations, division by zero).
func (m *Machine) Step() error {
	raw, ok := m.Mem[m.PC]
	if !ok {
		return fmt.Errorf("PC 0x%X: fetch from uninitialized memory", m.PC)
	}
	trits := rawToTrits(raw)
	for _, tr := range trits {
		if tr == 2 {
			return fmt.Errorf("PC 0x%X: undefined trit in instruction word", m.PC)
		}
	}
	opcode := int(fieldVal(trits, 0, isa.OpcodeTrits)) + isa.OpcodeBias
	ins, ok := isa.ByValue[opcode]
	if !ok {
		return fmt.Errorf("PC 0x%X: no such opcode %d", m.PC, opcode)
	}
	if m.Trace != nil {
		fmt.Fprintf(m.Trace, "0x%04X  %s\n", m.PC, ins.Name)
	}
	m.Stats.Instructions++
	m.Stats.ByCategory[ins.Category.String()]++
	m.Stats.ByMnemonic[ins.Name]++

	nextPC := m.PC + 1

	reg1 := fieldReg(trits, 4)

	switch ins.Format {
	case isa.FormatNone:
		switch ins.Name {
		case "NOP":
		case "WFI", "HALT":
			m.Halted = true
		case "RET":
			ta, _ := m.getReg(8) // TA is the link register (provisional ABI)
			nextPC = uint32(ta)
		case "CACHE", "FLUSH", "MEMBAR":
			// no-ops at architecture level
		case "SYSCALL", "BREAK":
			m.Halted = true
		}

	case isa.FormatR1:
		v, err := m.getReg(reg1)
		if err != nil {
			return err
		}
		switch ins.Name {
		case "INC":
			m.setReg(reg1, v+1)
		case "DEC":
			m.setReg(reg1, v-1)
		case "JR":
			nextPC = uint32(v)
		case "JALR":
			m.setReg(8, int64(m.PC+1))
			nextPC = uint32(v)
		case "PUSH", "POP":
			return fmt.Errorf("PC 0x%X: %s needs a stack ABI (not defined yet)", m.PC, ins.Name)
		}

	case isa.FormatR2:
		src, err := m.getReg(fieldReg(trits, 7))
		if err != nil {
			return err
		}
		var v int64
		switch ins.Name {
		case "NOT":
			v = tritwise(src, 0, func(x, _ int) int { return -x })
		case "NEG":
			v = -src
		case "ABS":
			if src < 0 {
				v = -src
			} else {
				v = src
			}
		case "SQRT":
			if src < 0 {
				return fmt.Errorf("PC 0x%X: SQRT of negative value %d", m.PC, src)
			}
			v = isqrt(src)
		case "FMOV", "FNEG", "VNOT":
			return fmt.Errorf("PC 0x%X: %s: vector/FP units not modeled", m.PC, ins.Name)
		default: // SIN, COS, TAN, ASIN, ACOS, ATAN, EXP, LOG
			return fmt.Errorf("PC 0x%X: %s: transcendental microcode not modeled", m.PC, ins.Name)
		}
		m.setReg(reg1, v)

	case isa.FormatR3:
		if ins.Class != isa.ClassGPR {
			return fmt.Errorf("PC 0x%X: %s: vector/FP units not modeled", m.PC, ins.Name)
		}
		a, err := m.getReg(fieldReg(trits, 7))
		if err != nil {
			return err
		}
		var b int64
		if trits[10] == 1 {
			b = fieldVal(trits, 11, 7)
		} else {
			b, err = m.getReg(fieldReg(trits, 15))
			if err != nil {
				return err
			}
		}
		var v int64
		switch ins.Name {
		case "ADD":
			v = a + b
		case "SUB":
			v = a - b
		case "MUL":
			v = a * b
		case "AND":
			v = tritwise(a, b, tritMin)
		case "OR":
			v = tritwise(a, b, tritMax)
		case "XOR":
			v = tritwise(a, b, tritMul)
		case "SHL":
			v = shiftTrits(a, b, true, false)
		case "SHR":
			v = shiftTrits(a, b, false, false)
		case "ROL":
			v = shiftTrits(a, b, true, true)
		case "ROR":
			v = shiftTrits(a, b, false, true)
		case "CMP":
			v = sign(a - b)
		case "TEST":
			v = sign(tritwise(a, b, tritMin))
		case "DIV", "UDIV":
			if b == 0 {
				return fmt.Errorf("PC 0x%X: division by zero", m.PC)
			}
			v = a / b
		case "MOD", "UMOD":
			if b == 0 {
				return fmt.Errorf("PC 0x%X: division by zero", m.PC)
			}
			v = a % b
		default:
			return fmt.Errorf("PC 0x%X: %s not modeled", m.PC, ins.Name)
		}
		m.setReg(reg1, v)

	case isa.FormatMem:
		if ins.Class != isa.ClassGPR {
			return fmt.Errorf("PC 0x%X: %s: vector/FP units not modeled", m.PC, ins.Name)
		}
		var ea int64
		switch trits[7] {
		case 0: // direct
			ea = fieldVal(trits, 8, 10)
		case 1: // base + index
			base, _ := m.getReg(fieldReg(trits, 8))
			index, _ := m.getReg(fieldReg(trits, 11))
			ea = base + index
		default: // base + offset
			base, _ := m.getReg(fieldReg(trits, 8))
			ea = base + fieldVal(trits, 11, 7)
		}
		switch ins.Name {
		case "LD":
			// LD always reads memory; load-immediate is spelled LEA (the
			// effective address of a direct operand IS the immediate).
			m.Stats.Loads++
			m.setReg(reg1, m.ReadWord(uint32(ea)))
		case "ST":
			m.Stats.Stores++
			v, _ := m.getReg(reg1)
			m.writeWord(uint32(ea), v)
		case "LEA":
			m.setReg(reg1, ea)
		default:
			return fmt.Errorf("PC 0x%X: %s not modeled", m.PC, ins.Name)
		}

	case isa.FormatBranch:
		m.Stats.Branches++
		a, _ := m.getReg(reg1)
		b, _ := m.getReg(fieldReg(trits, 7))
		var take bool
		switch ins.Name {
		case "BEQ":
			take = a == b
		case "BNE":
			take = a != b
		case "BLT", "BLTU":
			take = a < b
		case "BGE", "BGEU":
			take = a >= b
		case "BGT":
			take = a > b
		case "BLE":
			take = a <= b
		}
		if take {
			m.Stats.BranchesTaken++
			nextPC = uint32(int64(m.PC) + fieldVal(trits, 10, 8))
		}

	case isa.FormatJump:
		target := uint32(fieldVal(trits, 4, 14))
		switch ins.Name {
		case "JAL", "CALL":
			m.setReg(8, int64(m.PC+1)) // link in TA (provisional ABI)
		}
		nextPC = target
	}

	m.PC = nextPC
	return nil
}

// Run executes until WFI/HALT or maxSteps instructions.
func (m *Machine) Run(maxSteps uint64) error {
	for !m.Halted {
		if m.Stats.Instructions >= maxSteps {
			return fmt.Errorf("stopped after %d instructions without reaching WFI/HALT", maxSteps)
		}
		if err := m.Step(); err != nil {
			return err
		}
	}
	return nil
}
