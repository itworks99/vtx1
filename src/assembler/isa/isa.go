// Package isa is the single source of truth for the VTX1 trit-native
// instruction set architecture. The assembler encodes from these tables, and
// the generator in ./gen emits the matching Verilog constants
// (src/common/vtx1_isa_generated.vh) so the toolchain and the RTL can never
// drift apart.
//
// Architecture summary
// ====================
// VTX1 is a balanced-ternary machine. One machine word is 18 trits,
// physically stored as 36 bits with 2 bits per trit (hardware encoding:
// 00 = -1, 01 = 0, 10 = +1, 11 = undefined/invalid).
//
// Every instruction occupies exactly one 18-trit word. A VLIW bundle is
// 3 instruction words = 54 trits = 108 bits (matching `VLIW_WIDTH).
// Memory is word-addressed; addresses and immediates are balanced-ternary
// and therefore naturally signed and symmetric around zero.
//
// Instruction word layout (trit 0 = most significant):
//
//	all formats:  trits[0:4]   opcode (4 trits: 81 = 3^4 values, of which
//	                           69 are assigned and 12 reserved)
//
//	R3:      trits[4:7] rd | [7:10] rs1 | [10] mode | [11:18] payload
//	         mode  0: payload[4:7] = rs2 (last 3 trits)
//	         mode +1: payload = 7-trit balanced immediate (+-1093)
//	R2:      trits[4:7] rd | [7:10] rs
//	R1:      trits[4:7] r
//	None:    (all operand trits zero)
//	Mem:     trits[4:7] reg | [7] mode | [8:18] payload
//	         mode  0: payload = 10-trit balanced address/immediate (+-29524)
//	         mode +1: payload[0:3] = base, payload[3:6] = index
//	         mode -1: payload[0:3] = base, payload[3:10] = 7-trit offset
//	Branch:  trits[4:7] rs1 | [7:10] rs2 | [10:18] 8-trit PC-relative
//	         word offset (+-3280). rs2 = TZ encodes "compare against zero".
//	Jump:    trits[4:18] 14-trit absolute word address (+-2391484;
//	         negative addresses are architecturally legal, reserved for
//	         memory-mapped I/O)
//
// The all-zero-trit word decodes as NOP (opcode value 40, pattern 0000).
package isa

import "fmt"

// Geometry of the machine.
const (
	SlotTrits  = 18 // trits per instruction/machine word
	WordBytes  = 5  // bytes per stored word (36 bits + 4 padding bits)
	VLIWSlots  = 3  // instruction words per VLIW bundle
	BundleTrits = SlotTrits * VLIWSlots
	BundleBytes = 14 // 108 bits rounded up to bytes (matches `VLIW_BYTES)

	OpcodeTrits = 4
	RegTrits    = 3

	// Field value biases: a field's trit pattern is BalancedTrits(value-bias).
	OpcodeBias = 40 // opcode values 0..80 map to balanced -40..+40
	RegBias    = 13 // register indices 0..26 map to balanced -13..+13
)

// Balanced field ranges.
const (
	MaxImm7     = 1093      // (3^7-1)/2:  R3 immediate, memory offset
	MaxAddr10   = 29524     // (3^10-1)/2: direct memory address/immediate
	MaxRel8     = 3280      // (3^8-1)/2:  branch word offset
	MaxTarget14 = 2391484   // (3^14-1)/2: jump target word address
	MaxWord18   = 193710244 // (3^18-1)/2: data word value
)

// Format describes an instruction's operand encoding.
type Format int

const (
	FormatNone Format = iota
	FormatR3
	FormatR2
	FormatR1
	FormatMem
	FormatBranch
	FormatJump
)

// RegClass is the register file an operand must come from.
type RegClass int

const (
	ClassGPR RegClass = iota
	ClassVec
	ClassFP
)

// Category drives the VLIW packing rules.
type Category int

const (
	CatALU Category = iota
	CatMEM
	CatCTRL
	CatVEC
	CatFP
	CatSYS
)

func (c Category) String() string {
	switch c {
	case CatALU:
		return "ALU"
	case CatMEM:
		return "memory"
	case CatCTRL:
		return "control"
	case CatVEC:
		return "vector"
	case CatFP:
		return "floating point"
	case CatSYS:
		return "system"
	}
	return "unknown"
}

// Instruction is one row of the ISA.
type Instruction struct {
	Name       string
	Value      int // opcode value 0..80; trit pattern = BalancedTrits(Value-OpcodeBias, 4)
	Format     Format
	Class      RegClass // register class of the register operands
	Category   Category
	MultiCycle bool // executed by a dedicated multi-cycle unit (stalls the pipeline)
}

// Instructions lists the 69 operations of the 4-trit opcode space (81 = 3^4
// values). NOP holds value 40 so its pattern - and with zeroed operands the
// entire word - is all-zero trits. VADD takes NOP's former block slot
// (value 0) in exchange.
//
// There is no microcode: every operation is either single-cycle hardwired
// control or a dedicated multi-cycle unit (divider, vector reduction).
// Deliberately absent, with their opcode values reserved (see
// ReservedValues):
//   - SIN/COS/TAN/ASIN/ACOS/ATAN/EXP/LOG (61-68): transcendentals are
//     software library routines over the (excellent, Booth-free) ternary
//     primitives; hardware acceleration, if ever needed, is a CORDIC FSM
//     unit, not an interpreter.
//   - UDIV/UMOD (57-58) and BLTU/BGEU (34-35): "unsigned" is meaningless
//     in balanced ternary; the number system is inherently signed and
//     symmetric.
var Instructions = []Instruction{
	// ALU
	{"NOP", 40, FormatNone, ClassGPR, CatALU, false},
	{"ADD", 1, FormatR3, ClassGPR, CatALU, false},
	{"SUB", 2, FormatR3, ClassGPR, CatALU, false},
	{"MUL", 3, FormatR3, ClassGPR, CatALU, false},
	{"AND", 4, FormatR3, ClassGPR, CatALU, false},
	{"OR", 5, FormatR3, ClassGPR, CatALU, false},
	{"NOT", 6, FormatR2, ClassGPR, CatALU, false},
	{"XOR", 7, FormatR3, ClassGPR, CatALU, false},
	{"SHL", 8, FormatR3, ClassGPR, CatALU, false},
	{"SHR", 9, FormatR3, ClassGPR, CatALU, false},
	{"ROL", 10, FormatR3, ClassGPR, CatALU, false},
	{"ROR", 11, FormatR3, ClassGPR, CatALU, false},
	{"CMP", 12, FormatR3, ClassGPR, CatALU, false},
	{"TEST", 13, FormatR3, ClassGPR, CatALU, false},
	{"INC", 14, FormatR1, ClassGPR, CatALU, false},
	{"DEC", 15, FormatR1, ClassGPR, CatALU, false},
	{"NEG", 16, FormatR2, ClassGPR, CatALU, false},
	// MEM
	{"LD", 17, FormatMem, ClassGPR, CatMEM, false},
	{"ST", 18, FormatMem, ClassGPR, CatMEM, false},
	{"VLD", 19, FormatMem, ClassVec, CatMEM, false},
	{"VST", 20, FormatMem, ClassVec, CatMEM, false},
	{"FLD", 21, FormatMem, ClassFP, CatMEM, false},
	{"FST", 22, FormatMem, ClassFP, CatMEM, false},
	{"LEA", 23, FormatMem, ClassGPR, CatMEM, false},
	{"PUSH", 24, FormatR1, ClassGPR, CatMEM, false},
	{"POP", 25, FormatR1, ClassGPR, CatMEM, false},
	// CTRL
	{"JMP", 26, FormatJump, ClassGPR, CatCTRL, false},
	{"JAL", 27, FormatJump, ClassGPR, CatCTRL, false},
	{"JR", 28, FormatR1, ClassGPR, CatCTRL, false},
	{"JALR", 29, FormatR1, ClassGPR, CatCTRL, false},
	{"BEQ", 30, FormatBranch, ClassGPR, CatCTRL, false},
	{"BNE", 31, FormatBranch, ClassGPR, CatCTRL, false},
	{"BLT", 32, FormatBranch, ClassGPR, CatCTRL, false},
	{"BGE", 33, FormatBranch, ClassGPR, CatCTRL, false},
	{"BGT", 36, FormatBranch, ClassGPR, CatCTRL, false},
	{"BLE", 37, FormatBranch, ClassGPR, CatCTRL, false},
	{"CALL", 38, FormatJump, ClassGPR, CatCTRL, false},
	{"RET", 39, FormatNone, ClassGPR, CatCTRL, false},
	// VEC (VADD takes value 0; see NOP note above)
	{"VADD", 0, FormatR3, ClassVec, CatVEC, false},
	{"VSUB", 41, FormatR3, ClassVec, CatVEC, false},
	{"VMUL", 42, FormatR3, ClassVec, CatVEC, false},
	{"VAND", 43, FormatR3, ClassVec, CatVEC, false},
	{"VOR", 44, FormatR3, ClassVec, CatVEC, false},
	{"VNOT", 45, FormatR2, ClassVec, CatVEC, false},
	{"VSHL", 46, FormatR3, ClassVec, CatVEC, false},
	{"VSHR", 47, FormatR3, ClassVec, CatVEC, false},
	// FP
	{"FADD", 48, FormatR3, ClassFP, CatFP, false},
	{"FSUB", 49, FormatR3, ClassFP, CatFP, false},
	{"FMUL", 50, FormatR3, ClassFP, CatFP, false},
	{"FCMP", 51, FormatR3, ClassFP, CatFP, false},
	{"FMOV", 52, FormatR2, ClassFP, CatFP, false},
	{"FNEG", 53, FormatR2, ClassFP, CatFP, false},
	// SYS
	{"WFI", 54, FormatNone, ClassGPR, CatSYS, false},
	// Multi-cycle ALU units (iterative divider / root FSMs)
	{"DIV", 55, FormatR3, ClassGPR, CatALU, true},
	{"MOD", 56, FormatR3, ClassGPR, CatALU, true},
	{"SQRT", 59, FormatR2, ClassGPR, CatALU, true},
	{"ABS", 60, FormatR2, ClassGPR, CatALU, false},
	// Multi-cycle vector reductions
	{"VDOT", 69, FormatR3, ClassVec, CatVEC, true},
	{"VREDUCE", 70, FormatR3, ClassVec, CatVEC, true},
	{"VMAX", 71, FormatR3, ClassVec, CatVEC, true},
	{"VMIN", 72, FormatR3, ClassVec, CatVEC, true},
	{"VSUM", 73, FormatR3, ClassVec, CatVEC, true},
	{"VPERM", 74, FormatR3, ClassVec, CatVEC, true},
	// Memory system control (single-cycle)
	{"CACHE", 75, FormatNone, ClassGPR, CatMEM, false},
	{"FLUSH", 76, FormatNone, ClassGPR, CatMEM, false},
	{"MEMBAR", 77, FormatNone, ClassGPR, CatMEM, false},
	// System control (single-cycle)
	{"SYSCALL", 78, FormatNone, ClassGPR, CatSYS, false},
	{"BREAK", 79, FormatNone, ClassGPR, CatSYS, false},
	{"HALT", 80, FormatNone, ClassGPR, CatSYS, false},
}

// ReservedValues are opcode values intentionally left unassigned (see the
// note above Instructions). They decode as invalid.
var ReservedValues = []int{34, 35, 57, 58, 61, 62, 63, 64, 65, 66, 67, 68}

// Register is one architectural register in the unified 3-trit index space.
type Register struct {
	Name  string
	Index int // 0..26; trit pattern = BalancedTrits(Index-RegBias, 3)
	Class RegClass
}

// ZeroRegIndex is TZ, the architectural zero register: reads as 0, writes
// are discarded. In branch rs2 position it encodes "compare against zero".
const ZeroRegIndex = 18

var Registers = []Register{
	{"T0", 0, ClassGPR}, {"T1", 1, ClassGPR}, {"T2", 2, ClassGPR},
	{"T3", 3, ClassGPR}, {"T4", 4, ClassGPR}, {"T5", 5, ClassGPR},
	{"T6", 6, ClassGPR},
	{"TB", 7, ClassGPR},  // base pointer
	{"TA", 8, ClassGPR},  // accumulator
	{"TC", 9, ClassGPR},  // program counter
	{"TS", 10, ClassGPR}, // status
	{"TI", 11, ClassGPR}, // instruction
	{"VA", 12, ClassVec}, {"VT", 13, ClassVec}, {"VB", 14, ClassVec},
	{"FA", 15, ClassFP}, {"FT", 16, ClassFP}, {"FB", 17, ClassFP},
	{"TZ", ZeroRegIndex, ClassGPR}, // zero register
}

// Lookup tables.
var (
	ByName     = map[string]*Instruction{}
	ByValue    = map[int]*Instruction{}
	RegByName  = map[string]*Register{}
	RegByIndex = map[int]*Register{}
)

func init() {
	for i := range Instructions {
		ins := &Instructions[i]
		ByName[ins.Name] = ins
		ByValue[ins.Value] = ins
	}
	for i := range Registers {
		r := &Registers[i]
		RegByName[r.Name] = r
		RegByIndex[r.Index] = r
	}
}

// --- Balanced ternary helpers ---

// BalancedTrits converts v to n balanced trits (-1, 0, +1), most significant
// first. It errors when v is outside +-(3^n-1)/2.
func BalancedTrits(v int64, n int) ([]int, error) {
	limit := (pow3(n) - 1) / 2
	if v < -limit || v > limit {
		return nil, fmt.Errorf("value %d does not fit in %d balanced trits (range +-%d)", v, n, limit)
	}
	trits := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		r := v % 3
		v /= 3
		switch r {
		case 2:
			r = -1
			v++
		case -2:
			r = 1
			v--
		}
		trits[i] = int(r)
	}
	return trits, nil
}

// TritsValue is the inverse of BalancedTrits.
func TritsValue(trits []int) int64 {
	var v int64
	for _, t := range trits {
		v = v*3 + int64(t)
	}
	return v
}

func pow3(n int) int64 {
	v := int64(1)
	for i := 0; i < n; i++ {
		v *= 3
	}
	return v
}

// --- Physical trit encoding (2 bits per trit) ---

// EncodeTrits packs trits into bytes using the hardware encoding
// (00 = -1, 01 = 0, 10 = +1), most significant trit first; the final byte is
// zero-padded. 18 trits yield 5 bytes (one stored machine word).
func EncodeTrits(trits []int) []byte {
	out := make([]byte, 0, (len(trits)+3)/4)
	var cur byte
	for i, t := range trits {
		var enc byte
		switch t {
		case -1:
			enc = 0b00
		case 0:
			enc = 0b01
		case 1:
			enc = 0b10
		default:
			enc = 0b11 // undefined
		}
		cur = cur<<2 | enc
		if i%4 == 3 {
			out = append(out, cur)
			cur = 0
		}
	}
	if len(trits)%4 != 0 {
		cur <<= 2 * (4 - len(trits)%4)
		out = append(out, cur)
	}
	return out
}

// DecodeTrits is the inverse of EncodeTrits: it reads n trits from data.
// An undefined pair (0b11) decodes as the sentinel value 2.
func DecodeTrits(data []byte, n int) []int {
	trits := make([]int, 0, n)
	for _, b := range data {
		for shift := 6; shift >= 0 && len(trits) < n; shift -= 2 {
			switch (b >> shift) & 0b11 {
			case 0b00:
				trits = append(trits, -1)
			case 0b01:
				trits = append(trits, 0)
			case 0b10:
				trits = append(trits, 1)
			default:
				trits = append(trits, 2) // undefined
			}
		}
	}
	return trits
}

// OpcodePattern returns the 4-trit pattern for an opcode value.
func OpcodePattern(value int) []int {
	trits, _ := BalancedTrits(int64(value-OpcodeBias), OpcodeTrits)
	return trits
}

// RegPattern returns the 3-trit pattern for a register index.
func RegPattern(index int) []int {
	trits, _ := BalancedTrits(int64(index-RegBias), RegTrits)
	return trits
}
