package rvsim

import "testing"

func run(t *testing.T, prog []Instr, err error, load func(*Machine)) *Machine {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	m := New(prog)
	load(m)
	if err := m.Run(1_000_000); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestDotProduct(t *testing.T) {
	prog, err := DotProduct()
	m := run(t, prog, err, LoadVectors)
	if got := m.ReadWord(addrResult); got != -107 {
		t.Errorf("dot product = %d, want -107", got)
	}
}

func TestSignCount(t *testing.T) {
	prog, err := SignCount()
	m := run(t, prog, err, LoadSignData)
	if neg := m.ReadWord(addrResult); neg != 6 {
		t.Errorf("negative count = %d, want 6", neg)
	}
	if zero := m.ReadWord(addrResult + 4); zero != 4 {
		t.Errorf("zero count = %d, want 4", zero)
	}
	if pos := m.ReadWord(addrResult + 8); pos != 6 {
		t.Errorf("positive count = %d, want 6", pos)
	}
}

func TestAbsSum(t *testing.T) {
	prog, err := AbsSum()
	m := run(t, prog, err, LoadVectors)
	if got := m.ReadWord(addrResult); got != 114 {
		t.Errorf("abs sum = %d, want 114", got)
	}
}

func TestSignSum(t *testing.T) {
	prog, err := SignSum()
	m := run(t, prog, err, LoadVectors)
	if got := m.ReadWord(addrResult); got != -2 {
		t.Errorf("sign sum = %d, want -2", got)
	}
}

// TestComparisonStats prints the RV32IM column of the Tier A table
// (go test -run TestComparisonStats -v).
func TestComparisonStats(t *testing.T) {
	cases := []struct {
		name string
		prog func() ([]Instr, error)
		load func(*Machine)
	}{
		{"dot_product", DotProduct, LoadVectors},
		{"sign_count", SignCount, LoadSignData},
		{"abs_sum", AbsSum, LoadVectors},
		{"sign_sum", SignSum, LoadVectors},
	}
	for _, c := range cases {
		prog, err := c.prog()
		m := run(t, prog, err, c.load)
		s := m.Stats
		t.Logf("%-12s static=%2d dyn=%3d ALU=%3d mem=%2d ctrl=%2d loads=%2d stores=%d branches=%2d(taken %2d)",
			c.name, len(prog), s.Instructions,
			s.ByCategory["ALU"], s.ByCategory["memory"], s.ByCategory["control"],
			s.Loads, s.Stores, s.Branches, s.BranchesTaken)
	}
}
