package isa

import "testing"

func TestOpcodeSpaceIsExactlyFull(t *testing.T) {
	if len(Instructions) != 81 {
		t.Fatalf("instruction count = %d, want 81 (= 3^4, the full 4-trit opcode space)", len(Instructions))
	}
	seen := make(map[int]string)
	for _, ins := range Instructions {
		if ins.Value < 0 || ins.Value > 80 {
			t.Errorf("%s: opcode value %d out of range 0..80", ins.Name, ins.Value)
		}
		if other, dup := seen[ins.Value]; dup {
			t.Errorf("opcode value %d assigned to both %s and %s", ins.Value, other, ins.Name)
		}
		seen[ins.Value] = ins.Name
	}
}

func TestNopIsAllZeroTrits(t *testing.T) {
	nop := ByName["NOP"]
	for i, tr := range OpcodePattern(nop.Value) {
		if tr != 0 {
			t.Fatalf("NOP opcode trit %d = %d, want 0 (the all-zero word must be NOP)", i, tr)
		}
	}
}

func TestBalancedTritsRoundTrip(t *testing.T) {
	cases := []int64{0, 1, -1, 2, -2, 40, -40, 1093, -1093, 29524, -29524, MaxWord18, -MaxWord18}
	for _, v := range cases {
		trits, err := BalancedTrits(v, SlotTrits)
		if err != nil {
			t.Errorf("BalancedTrits(%d): %v", v, err)
			continue
		}
		if got := TritsValue(trits); got != v {
			t.Errorf("round trip of %d = %d", v, got)
		}
	}
	if _, err := BalancedTrits(MaxWord18+1, SlotTrits); err == nil {
		t.Error("out-of-range value should fail")
	}
	if _, err := BalancedTrits(2, 1); err == nil {
		t.Error("2 does not fit in one balanced trit")
	}
}

func TestEncodeDecodeTrits(t *testing.T) {
	trits, _ := BalancedTrits(-123456, SlotTrits)
	packed := EncodeTrits(trits)
	if len(packed) != WordBytes {
		t.Fatalf("packed word = %d bytes, want %d", len(packed), WordBytes)
	}
	back := DecodeTrits(packed, SlotTrits)
	for i := range trits {
		if back[i] != trits[i] {
			t.Fatalf("byte round trip mismatch: %v != %v", back, trits)
		}
	}
}

func TestPatternsContainNoUndefinedPairs(t *testing.T) {
	for _, ins := range Instructions {
		packed := EncodeTrits(OpcodePattern(ins.Value))
		for _, tr := range DecodeTrits(packed, OpcodeTrits) {
			if tr == 2 {
				t.Errorf("%s opcode pattern contains an undefined trit pair", ins.Name)
			}
		}
	}
	for _, r := range Registers {
		packed := EncodeTrits(RegPattern(r.Index))
		for _, tr := range DecodeTrits(packed, RegTrits) {
			if tr == 2 {
				t.Errorf("register %s pattern contains an undefined trit pair", r.Name)
			}
		}
	}
}

func TestRegisterIndicesUnique(t *testing.T) {
	seen := make(map[int]string)
	for _, r := range Registers {
		if other, dup := seen[r.Index]; dup {
			t.Errorf("register index %d assigned to both %s and %s", r.Index, other, r.Name)
		}
		seen[r.Index] = r.Name
	}
}

func TestBundleGeometry(t *testing.T) {
	if BundleTrits != 54 {
		t.Errorf("bundle = %d trits, want 54", BundleTrits)
	}
	if got := len(EncodeTrits(make([]int, BundleTrits))); got != BundleBytes {
		t.Errorf("packed bundle = %d bytes, want %d (`VLIW_BYTES)", got, BundleBytes)
	}
}
