package cmd

import (
	"strings"
	"testing"

	"github.com/kvany/vtx1/assembler/isa"
)

func newTestCG() *CodeGenerator {
	cg := NewCodeGenerator(NewSymbolTable())
	cg.Labels["loop"] = 0x1010 // far label for absolute jumps
	cg.Labels["near"] = 0x20   // close label for PC-relative branches
	cg.Labels["data"] = 0x1100
	return cg
}

func reg(name string) *RegisterNode     { return &RegisterNode{Name: name} }
func imm(val string) *ImmediateNode     { return &ImmediateNode{Value: val} }
func ident(name string) *IdentifierNode { return &IdentifierNode{Name: name} }

func ins(mnemonic string, ops ...OperandNode) *InstructionNode {
	return &InstructionNode{Mnemonic: mnemonic, Operands: ops, Line: 1}
}

func mustEncode(t *testing.T, cg *CodeGenerator, instr *InstructionNode) (Slot, opInfo) {
	t.Helper()
	slot, info, err := cg.encodeOperation(instr, cg.CurrentAddr)
	if err != nil {
		t.Fatalf("encodeOperation(%s): %v", instr.Mnemonic, err)
	}
	return slot, info
}

// disasm encodes and immediately disassembles, round-tripping the slot
// through its 5-byte stored form.
func disasm(t *testing.T, cg *CodeGenerator, instr *InstructionNode) string {
	t.Helper()
	slot, _ := mustEncode(t, cg, instr)
	return disassembleSlot(isa.DecodeTrits(slot.Bytes(), isa.SlotTrits))
}

func TestEncodeR3(t *testing.T) {
	cg := newTestCG()
	slot, info := mustEncode(t, cg, ins("ADD", reg("T0"), reg("T1"), reg("T2")))
	if got := int(fieldValue(slot[:], 0, 4)) + isa.OpcodeBias; got != isa.ByName["ADD"].Value {
		t.Errorf("opcode value = %d, want %d", got, isa.ByName["ADD"].Value)
	}
	if fieldReg(slot[:], 4) != 0 || fieldReg(slot[:], 7) != 1 || fieldReg(slot[:], 15) != 2 {
		t.Errorf("register fields wrong: %v", slot)
	}
	if len(info.Writes) != 1 || info.Writes[0] != "T0" {
		t.Errorf("ADD writes = %v, want [T0]", info.Writes)
	}
	if len(info.Reads) != 2 {
		t.Errorf("ADD reads = %v, want 2 registers", info.Reads)
	}
}

func TestEncodeR3Immediate(t *testing.T) {
	cg := newTestCG()
	if got := disasm(t, cg, ins("ADD", reg("T0"), reg("T0"), imm("1"))); got != "ADD T0, T0, 1" {
		t.Errorf("disasm = %q", got)
	}
	// Balanced immediates are signed and symmetric.
	if got := disasm(t, cg, ins("ADD", reg("T0"), reg("T0"), imm("-1093"))); got != "ADD T0, T0, -1093" {
		t.Errorf("disasm = %q", got)
	}
	// Out of the 7-trit range must fail.
	if _, _, err := cg.encodeOperation(ins("ADD", reg("T0"), reg("T0"), imm("1094")), 0); err == nil {
		t.Error("ADD with immediate 1094 should fail (range is +-1093)")
	}
}

func TestEncodeSpecialRegisters(t *testing.T) {
	cg := newTestCG()
	if got := disasm(t, cg, ins("ADD", reg("TA"), reg("TB"), reg("T6"))); got != "ADD TA, TB, T6" {
		t.Errorf("disasm = %q", got)
	}
}

func TestEncodeMemDirect(t *testing.T) {
	cg := newTestCG()
	slot, info := mustEncode(t, cg, ins("LD", reg("T1"), imm("0x2000")))
	if got := fieldValue(slot[:], 8, 10); got != 0x2000 {
		t.Errorf("LD direct field = %d, want 8192", got)
	}
	if len(info.Writes) != 1 || info.Writes[0] != "T1" {
		t.Errorf("LD writes = %v, want [T1]", info.Writes)
	}
	// Balanced addressing means negative immediates need no special casing.
	if got := disasm(t, cg, ins("LD", reg("T0"), imm("0t-"))); got != "LD T0, -1" {
		t.Errorf("disasm = %q", got)
	}
	// Out of the 10-trit range must fail.
	if _, _, err := cg.encodeOperation(ins("LD", reg("T0"), imm("29525")), 0); err == nil {
		t.Error("LD with direct value 29525 should fail (range is +-29524)")
	}
}

func TestEncodeMemModes(t *testing.T) {
	cg := newTestCG()
	cases := []struct {
		op   OperandNode
		want string
	}{
		{&MemoryOperandNode{Base: "T0"}, "LD T3, [T0]"},
		{&MemoryOperandNode{Base: "TB", Index: "T2"}, "LD T3, [TB+T2]"},
		{&MemoryOperandNode{Base: "TB", Offset: "1"}, "LD T3, [TB+1]"},
		{&MemoryOperandNode{Base: "TB", Offset: "0t-"}, "LD T3, [TB+-1]"},
	}
	for _, c := range cases {
		if got := disasm(t, cg, ins("LD", reg("T3"), c.op)); got != c.want {
			t.Errorf("disasm = %q, want %q", got, c.want)
		}
	}
	// Address registers are reads.
	_, info := mustEncode(t, cg, ins("LD", reg("T5"), &MemoryOperandNode{Base: "TB", Index: "T2"}))
	if len(info.Reads) != 2 {
		t.Errorf("LD [TB+T2] reads = %v, want [TB T2]", info.Reads)
	}
}

func TestEncodeStoreReadsRegister(t *testing.T) {
	cg := newTestCG()
	_, info := mustEncode(t, cg, ins("ST", reg("T3"), ident("data")))
	if len(info.Writes) != 0 {
		t.Errorf("ST writes = %v, want none", info.Writes)
	}
	if len(info.Reads) != 1 || info.Reads[0] != "T3" {
		t.Errorf("ST reads = %v, want [T3]", info.Reads)
	}
}

func TestEncodeBranchRelative(t *testing.T) {
	cg := newTestCG()
	cg.CurrentAddr = 0x10 // branch at word 0x10, target "near" = 0x20
	if got := disasm(t, cg, ins("BLT", reg("T2"), reg("T0"), ident("near"))); got != "BLT T2, T0, PC+16" {
		t.Errorf("disasm = %q", got)
	}
	// Backward branches are negative offsets, naturally signed in balanced ternary.
	cg.CurrentAddr = 0x30
	if got := disasm(t, cg, ins("BNE", reg("T1"), reg("T2"), ident("near"))); got != "BNE T1, T2, PC-16" {
		t.Errorf("disasm = %q", got)
	}
	// A target beyond +-3280 words must fail.
	cg.CurrentAddr = 0
	if _, _, err := cg.encodeOperation(ins("BEQ", reg("T0"), reg("T1"), ident("loop")), 0); err == nil {
		t.Error("branch to a target 4112 words away should fail (range is +-3280)")
	}
}

func TestEncodeBranchAgainstZero(t *testing.T) {
	cg := newTestCG()
	cg.CurrentAddr = 0x18
	slot, _ := mustEncode(t, cg, ins("BEQ", reg("T3"), imm("0"), ident("near")))
	if got := fieldReg(slot[:], 7); got != isa.ZeroRegIndex {
		t.Errorf("rs2 index = %d, want TZ (%d)", got, isa.ZeroRegIndex)
	}
	if got := disassembleSlot(slot[:]); got != "BEQ T3, 0, PC+8" {
		t.Errorf("disasm = %q", got)
	}
	// Non-zero immediates are not encodable.
	if _, _, err := cg.encodeOperation(ins("BEQ", reg("T3"), imm("5"), ident("near")), 0x18); err == nil {
		t.Error("BEQ with non-zero immediate should fail")
	}
}

func TestEncodeJumpAndVectorAndFP(t *testing.T) {
	cg := newTestCG()
	if got := disasm(t, cg, ins("JMP", ident("loop"))); got != "JMP 0x1010" {
		t.Errorf("disasm = %q", got)
	}
	if got := disasm(t, cg, ins("VADD", reg("VA"), reg("VT"), reg("VB"))); got != "VADD VA, VT, VB" {
		t.Errorf("disasm = %q", got)
	}
	// GPR in a vector slot must fail.
	if _, _, err := cg.encodeOperation(ins("VADD", reg("T0"), reg("VT"), reg("VB")), 0); err == nil {
		t.Error("VADD with GPR operand should fail")
	}
	if got := disasm(t, cg, ins("FMOV", reg("FA"), reg("FB"))); got != "FMOV FA, FB" {
		t.Errorf("disasm = %q", got)
	}
}

func TestEncodeOperandCountErrors(t *testing.T) {
	cg := newTestCG()
	cases := []*InstructionNode{
		ins("ADD", reg("T0"), reg("T1")), // needs 3
		ins("NOP", reg("T0")),            // needs 0
		ins("NOT", reg("T0")),            // needs 2
		ins("INC"),                       // needs 1
		ins("LD", reg("T0")),             // needs 2
		ins("JMP"),                       // needs 1
	}
	for _, c := range cases {
		if _, _, err := cg.encodeOperation(c, 0); err == nil {
			t.Errorf("%s with %d operands should fail", c.Mnemonic, len(c.Operands))
		}
	}
}

func TestEncodeUnknownRegister(t *testing.T) {
	cg := newTestCG()
	if _, _, err := cg.encodeOperation(ins("ADD", reg("T7"), reg("T0"), reg("T1")), 0); err == nil {
		t.Error("register T7 does not exist and should be rejected")
	}
}

func TestParseBalancedTernary(t *testing.T) {
	cases := map[string]int64{
		"0t+":   1,
		"0t-":   -1,
		"0t0":   0,
		"0t+0-": 8,  // 9 + 0 - 1
		"0t+-0": 6,  // 9 - 3 + 0
		"0t++-": 11, // 9 + 3 - 1
		"0t--+": -11,
	}
	for lit, want := range cases {
		got, err := parseImmediateString(lit)
		if err != nil {
			t.Errorf("parseImmediateString(%q): %v", lit, err)
			continue
		}
		if got != want {
			t.Errorf("parseImmediateString(%q) = %d, want %d", lit, got, want)
		}
	}
	if _, err := parseImmediateString("0t+2-"); err == nil {
		t.Error("invalid ternary digit should fail")
	}
}

func TestDisassembleRoundTrip(t *testing.T) {
	cg := newTestCG()
	cg.CurrentAddr = 0x1c
	cases := []struct {
		instr *InstructionNode
		want  string
	}{
		{ins("ADD", reg("T0"), reg("T1"), reg("T2")), "ADD T0, T1, T2"},
		{ins("ADD", reg("T0"), reg("T0"), imm("1")), "ADD T0, T0, 1"},
		{ins("LD", reg("T1"), imm("0x2000")), "LD T1, 0x2000"},
		{ins("LD", reg("T3"), &MemoryOperandNode{Base: "TB", Index: "T2"}), "LD T3, [TB+T2]"},
		{ins("LD", reg("T3"), &MemoryOperandNode{Base: "TB", Offset: "1"}), "LD T3, [TB+1]"},
		{ins("BEQ", reg("T3"), imm("0"), ident("near")), "BEQ T3, 0, PC+4"},
		{ins("JMP", ident("loop")), "JMP 0x1010"},
		{ins("VADD", reg("VA"), reg("VT"), reg("VB")), "VADD VA, VT, VB"},
		{ins("WFI"), "WFI"},
		{ins("INC", reg("T4")), "INC T4"},
		{ins("NOP"), "NOP"},
	}
	for _, c := range cases {
		if got := disasm(t, cg, c.instr); got != c.want {
			t.Errorf("disassembly = %q, want %q", got, c.want)
		}
	}
}

func TestDisassembleDataAndInvalid(t *testing.T) {
	// A data word (value 72 = 'H') must not disassemble as an instruction.
	trits, _ := isa.BalancedTrits(72, isa.SlotTrits)
	if got := disassembleSlot(trits); got != ".word 72" {
		t.Errorf("data word disassembly = %q, want .word 72", got)
	}
	// Undefined trit pairs (0b11) are flagged.
	bad := isa.DecodeTrits([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xF0}, isa.SlotTrits)
	if got := disassembleSlot(bad); !strings.Contains(got, "invalid") {
		t.Errorf("undefined trit disassembly = %q, want .invalid", got)
	}
}
