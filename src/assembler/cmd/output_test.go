package cmd

import (
	"strings"
	"testing"

	"github.com/kvany/vtx1/assembler/isa"
)

func TestPackBundles(t *testing.T) {
	cg := newTestCG()
	for _, v := range []int64{1, 2, 3, 4} {
		if err := cg.emitDataWord(v, 1); err != nil {
			t.Fatal(err)
		}
	}
	out := packBundles(cg.Output)
	// 4 words -> 2 bundles of 14 bytes (second padded with zero-value words).
	if len(out) != 2*isa.BundleBytes {
		t.Fatalf("length = %d, want %d", len(out), 2*isa.BundleBytes)
	}
	// Round-trip the first bundle's trits and verify the word values.
	trits := isa.DecodeTrits(out[:isa.BundleBytes], isa.BundleTrits)
	for i, want := range []int64{1, 2, 3} {
		got := isa.TritsValue(trits[i*isa.SlotTrits : (i+1)*isa.SlotTrits])
		if got != want {
			t.Errorf("bundle word %d = %d, want %d", i, got, want)
		}
	}
	// Padding word in the second bundle is zero-value.
	trits = isa.DecodeTrits(out[isa.BundleBytes:], isa.BundleTrits)
	if got := isa.TritsValue(trits[isa.SlotTrits : 2*isa.SlotTrits]); got != 0 {
		t.Errorf("padding word = %d, want 0", got)
	}
}

func TestObsoleteWordSizesRejected(t *testing.T) {
	for _, ws := range []string{"8", "ternary", "64"} {
		if err := writeOutput([]byte{}, "/dev/null", "binary", ws); err == nil {
			t.Errorf("word size %q should be rejected", ws)
		}
	}
}

func TestFormatAsHex(t *testing.T) {
	// Known record: 4 bytes at address 0.
	out := formatAsHex([]byte{0x01, 0x02, 0x03, 0x04})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected data + EOF records, got %d lines: %q", len(lines), out)
	}
	// checksum: -(04+00+00+00+01+02+03+04) = -0x0E = 0xF2
	if lines[0] != ":0400000001020304F2" {
		t.Errorf("data record = %q, want :0400000001020304F2", lines[0])
	}
	if lines[1] != ":00000001FF" {
		t.Errorf("EOF record = %q", lines[1])
	}

	// Records are split every 16 bytes and addresses advance.
	long := formatAsHex(make([]byte, 20))
	if !strings.Contains(long, ":04001000") {
		t.Errorf("second record should start at address 0x0010:\n%s", long)
	}
}

func TestFormatAsObjDump(t *testing.T) {
	cg := newTestCG()
	if err := cg.emitInstruction(ins("ADD", reg("T0"), reg("T1"), reg("T2"))); err != nil {
		t.Fatal(err)
	}
	if err := cg.emitInstruction(ins("WFI")); err != nil {
		t.Fatal(err)
	}
	cg.Output = append(cg.Output, 0xAB) // trailing partial word
	out := formatAsObjDump(cg.Output)
	if !strings.Contains(out, "ADD T0, T1, T2") {
		t.Errorf("objdump missing ADD:\n%s", out)
	}
	if !strings.Contains(out, "WFI") {
		t.Errorf("objdump missing WFI:\n%s", out)
	}
	if !strings.Contains(out, "00000001") {
		t.Errorf("objdump missing word address column:\n%s", out)
	}
	if !strings.Contains(out, "partial word") {
		t.Errorf("objdump missing trailing partial word:\n%s", out)
	}
}
