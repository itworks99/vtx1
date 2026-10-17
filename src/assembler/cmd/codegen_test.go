package cmd

import (
	"strings"
	"testing"

	"github.com/kvany/vtx1/assembler/isa"
)

func vliw(instrs ...*InstructionNode) *VLIWInstructionNode {
	return &VLIWInstructionNode{Instructions: instrs, Line: 1}
}

func withWarnings(t *testing.T) *[]error {
	t.Helper()
	warnings := []error{}
	AttachCodegenWarnings(&warnings)
	t.Cleanup(func() { AttachCodegenWarnings(nil) })
	return &warnings
}

// wordValue decodes the n-th emitted word as a balanced ternary value.
func wordValue(t *testing.T, cg *CodeGenerator, n int) int64 {
	t.Helper()
	off := n * isa.WordBytes
	if off+isa.WordBytes > len(cg.Output) {
		t.Fatalf("output has no word %d (len %d)", n, len(cg.Output))
	}
	return isa.TritsValue(isa.DecodeTrits(cg.Output[off:off+isa.WordBytes], isa.SlotTrits))
}

// wordDisasm disassembles the n-th emitted word.
func wordDisasm(t *testing.T, cg *CodeGenerator, n int) string {
	t.Helper()
	off := n * isa.WordBytes
	return disassembleSlot(isa.DecodeTrits(cg.Output[off:off+isa.WordBytes], isa.SlotTrits))
}

func TestEmitInstruction(t *testing.T) {
	cg := newTestCG()
	if err := cg.emitInstruction(ins("ADD", reg("T0"), reg("T1"), reg("T2"))); err != nil {
		t.Fatal(err)
	}
	if len(cg.Output) != isa.WordBytes {
		t.Errorf("output length = %d, want %d (one word)", len(cg.Output), isa.WordBytes)
	}
	if cg.CurrentAddr != 1 {
		t.Errorf("current address = %d words, want 1", cg.CurrentAddr)
	}
}

func TestVLIWEmitsThreeSlots(t *testing.T) {
	withWarnings(t)
	cg := newTestCG()
	// A single op must be padded with NOPs to a full 3-word bundle.
	if err := cg.emitVLIWInstruction(vliw(ins("ADD", reg("T0"), reg("T1"), reg("T2")))); err != nil {
		t.Fatal(err)
	}
	if len(cg.Output) != 3*isa.WordBytes {
		t.Fatalf("VLIW output length = %d, want %d", len(cg.Output), 3*isa.WordBytes)
	}
	if cg.CurrentAddr != 3 {
		t.Errorf("current address = %d words, want 3", cg.CurrentAddr)
	}
	for slot := 1; slot < 3; slot++ {
		if got := wordDisasm(t, cg, slot); got != "NOP" {
			t.Errorf("padding slot %d = %q, want NOP", slot, got)
		}
	}
}

func TestVLIWMaxOneMemoryOp(t *testing.T) {
	withWarnings(t)
	cg := newTestCG()
	err := cg.emitVLIWInstruction(vliw(
		ins("ST", reg("T1"), ident("data")),
		ins("ST", reg("T2"), ident("data")),
	))
	if err == nil || !strings.Contains(err.Error(), "memory operations") {
		t.Errorf("two stores in one bundle should fail with a memory-op violation, got: %v", err)
	}
}

func TestVLIWMaxOneControlOp(t *testing.T) {
	withWarnings(t)
	cg := newTestCG()
	cg.CurrentAddr = 0x1e
	err := cg.emitVLIWInstruction(vliw(
		ins("JMP", ident("loop")),
		ins("BEQ", reg("T0"), imm("0"), ident("near")),
	))
	if err == nil || !strings.Contains(err.Error(), "control operations") {
		t.Errorf("two control ops in one bundle should fail, got: %v", err)
	}
}

func TestVLIWDestinationConflict(t *testing.T) {
	withWarnings(t)
	cg := newTestCG()
	err := cg.emitVLIWInstruction(vliw(
		ins("ADD", reg("T0"), reg("T1"), reg("T2")),
		ins("SUB", reg("T0"), reg("T3"), reg("T4")),
	))
	if err == nil || !strings.Contains(err.Error(), "resource conflict") {
		t.Errorf("two writes to T0 in one bundle should fail, got: %v", err)
	}
}

func TestVLIWNoConflictAcrossRegisterClasses(t *testing.T) {
	withWarnings(t)
	cg := newTestCG()
	// T0 and VA are different architectural registers.
	err := cg.emitVLIWInstruction(vliw(
		ins("ADD", reg("T0"), reg("T1"), reg("T2")),
		ins("VADD", reg("VA"), reg("VT"), reg("VB")),
	))
	if err != nil {
		t.Errorf("ALU and vector ops writing different register files should pack, got: %v", err)
	}
}

func TestVLIWHazardWarning(t *testing.T) {
	warnings := withWarnings(t)
	cg := newTestCG()
	// SUB reads T0, which ADD writes in the same bundle: legal, but a hazard warning.
	if err := cg.emitVLIWInstruction(vliw(
		ins("ADD", reg("T0"), reg("T1"), reg("T2")),
		ins("SUB", reg("T3"), reg("T0"), reg("T4")),
	)); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, w := range *warnings {
		if strings.Contains(w.Error(), "hazard") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a read-after-write hazard warning, warnings: %v", *warnings)
	}
}

func TestVLIWSelfReadIsNotAHazard(t *testing.T) {
	warnings := withWarnings(t)
	cg := newTestCG()
	// Each op reads only registers it writes itself: no cross-slot hazard.
	if err := cg.emitVLIWInstruction(vliw(
		ins("ADD", reg("T0"), reg("T0"), imm("1")),
		ins("ADD", reg("T1"), reg("T1"), imm("1")),
		ins("ADD", reg("T2"), reg("T2"), imm("1")),
	)); err != nil {
		t.Fatal(err)
	}
	for _, w := range *warnings {
		if strings.Contains(w.Error(), "hazard") {
			t.Errorf("unexpected hazard warning: %v", w)
		}
	}
}

func TestVLIWTooManyOps(t *testing.T) {
	withWarnings(t)
	cg := newTestCG()
	err := cg.emitVLIWInstruction(vliw(
		ins("NOP"), ins("NOP"), ins("NOP"), ins("NOP"),
	))
	if err == nil {
		t.Error("four ops in one bundle should fail")
	}
}

func TestEmitDirective(t *testing.T) {
	// .DB with a string and a terminator: one word per character.
	cg := newTestCG()
	err := cg.emitDirective(&DirectiveNode{
		Name:   ".DB",
		Params: []OperandNode{imm(`"Hi"`), imm("0")},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cg.Output) != 3*isa.WordBytes {
		t.Fatalf(".DB output = %d bytes, want 3 words", len(cg.Output))
	}
	for i, want := range []int64{'H', 'i', 0} {
		if got := wordValue(t, cg, i); got != want {
			t.Errorf(".DB word %d = %d, want %d", i, got, want)
		}
	}

	// .DW with values, including a negative ternary literal.
	cg = newTestCG()
	if err := cg.emitDirective(&DirectiveNode{
		Name:   ".DW",
		Params: []OperandNode{imm("5"), imm("0t--+")},
	}); err != nil {
		t.Fatal(err)
	}
	for i, want := range []int64{5, -11} {
		if got := wordValue(t, cg, i); got != want {
			t.Errorf(".DW word %d = %d, want %d", i, got, want)
		}
	}

	// Values outside the 18-trit word range must fail.
	cg = newTestCG()
	if err := cg.emitDirective(&DirectiveNode{
		Name:   ".DW",
		Params: []OperandNode{imm("193710245")},
	}); err == nil {
		t.Error(".DW with a value beyond +-(3^18-1)/2 should fail")
	}
}

func TestLabelResolution(t *testing.T) {
	ast := &AST{Program: &ProgramNode{Lines: []*LineNode{
		{Statement: &DirectiveNode{Name: ".ORG", Params: []OperandNode{imm("0x1000")}}},
		{Label: &LabelNode{Name: "start"}, Statement: ins("NOP")},
		{Statement: ins("JMP", ident("start"))},
	}}}
	cg := NewCodeGenerator(NewSymbolTable())
	if err := cg.Generate(ast); err != nil {
		t.Fatal(err)
	}
	if cg.Labels["start"] != 0x1000 {
		t.Errorf("label start = 0x%X, want 0x1000", cg.Labels["start"])
	}
	if got := wordDisasm(t, cg, 1); got != "JMP 0x1000" {
		t.Errorf("JMP disassembly = %q, want JMP 0x1000", got)
	}
}

func TestEquResolution(t *testing.T) {
	ast := &AST{Program: &ProgramNode{Lines: []*LineNode{
		{Statement: &DirectiveNode{Name: ".EQU", Params: []OperandNode{ident("IOBASE"), imm("0x2000")}}},
		{Statement: ins("LD", reg("T0"), ident("IOBASE"))},
	}}}
	cg := NewCodeGenerator(NewSymbolTable())
	if err := cg.Generate(ast); err != nil {
		t.Fatal(err)
	}
	if got := wordDisasm(t, cg, 0); got != "LD T0, 0x2000" {
		t.Errorf("LD disassembly = %q, want LD T0, 0x2000", got)
	}
}

func TestLabelPlusOffset(t *testing.T) {
	ast := &AST{Program: &ProgramNode{Lines: []*LineNode{
		{Statement: &DirectiveNode{Name: ".ORG", Params: []OperandNode{imm("0x1100")}}},
		{Label: &LabelNode{Name: "vec_a"}, Statement: &DirectiveNode{Name: ".DW", Params: []OperandNode{imm("1"), imm("2")}}},
	}}}
	cg := NewCodeGenerator(NewSymbolTable())
	if err := cg.Generate(ast); err != nil {
		t.Fatal(err)
	}
	// Word addressing: vec_a+1 is the second element.
	addr := cg.resolveOperandAddr(&BinaryOpNode{Left: ident("vec_a"), Op: "+", Right: imm("1")})
	if addr != 0x1101 {
		t.Errorf("vec_a+1 = 0x%X, want 0x1101", addr)
	}
}
