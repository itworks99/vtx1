package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kvany/vtx1/assembler/sim"
)

// runProgram assembles a source file to a memh image and executes it in the
// ISA simulator until WFI/HALT.
func runProgram(t *testing.T, source string) *sim.Machine {
	t.Helper()
	image := filepath.Join(t.TempDir(), "prog.memh")

	old := os.Stdout
	devnull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	os.Stdout = devnull
	err := RunAssembler(source, image, "", "memh", false, "", "36")
	os.Stdout = old
	devnull.Close()
	if err != nil {
		t.Fatalf("assembling %s: %v", source, err)
	}

	f, err := os.Open(image)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m := sim.New()
	entry, err := m.LoadMemh(f)
	if err != nil {
		t.Fatal(err)
	}
	m.PC = entry
	if err := m.Run(1_000_000); err != nil {
		t.Fatalf("executing %s: %v", source, err)
	}
	return m
}

// TestHelloWorldExecutes runs the flagship example end to end: assembled to
// balanced ternary and executed, it must write "Hello, VTX1!" to 0x2000,
// the character count to 0x2100, and the completion flag to 0x2104.
func TestHelloWorldExecutes(t *testing.T) {
	m := runProgram(t, "../examples/hello_world.asm")

	want := "Hello, VTX1!"
	for i, c := range want {
		if got := m.ReadWord(uint32(0x2000 + i)); got != int64(c) {
			t.Errorf("output[%d] = %d, want %d (%q)", i, got, c, c)
		}
	}
	if got := m.ReadWord(0x2100); got != int64(len(want)) {
		t.Errorf("character count = %d, want %d", got, len(want))
	}
	if got := m.ReadWord(0x2104); got != 1 {
		t.Errorf("completion flag = %d, want 1", got)
	}
}

// TestBenchmarkKernels executes each Tier A benchmark and checks its
// documented result.
func TestBenchmarkKernels(t *testing.T) {
	t.Run("dot_product", func(t *testing.T) {
		m := runProgram(t, "../benchmarks/dot_product.asm")
		if got := m.ReadWord(0x1200); got != -107 {
			t.Errorf("dot product = %d, want -107", got)
		}
	})
	t.Run("sign_count", func(t *testing.T) {
		m := runProgram(t, "../benchmarks/sign_count.asm")
		if neg := m.ReadWord(0x1200); neg != 6 {
			t.Errorf("negative count = %d, want 6", neg)
		}
		if zero := m.ReadWord(0x1201); zero != 4 {
			t.Errorf("zero count = %d, want 4", zero)
		}
		if pos := m.ReadWord(0x1202); pos != 6 {
			t.Errorf("positive count = %d, want 6", pos)
		}
	})
	t.Run("abs_sum", func(t *testing.T) {
		m := runProgram(t, "../benchmarks/abs_sum.asm")
		if got := m.ReadWord(0x1200); got != 114 {
			t.Errorf("abs sum = %d, want 114", got)
		}
	})
	t.Run("sign_sum", func(t *testing.T) {
		m := runProgram(t, "../benchmarks/sign_sum.asm")
		if got := m.ReadWord(0x1200); got != -2 {
			t.Errorf("sign sum = %d, want -2", got)
		}
	})
}

// TestMacroDemoExecutes runs the preprocessor demo end to end.
func TestMacroDemoExecutes(t *testing.T) {
	m := runProgram(t, "../examples/macro_demo.asm")
	if got := m.ReadWord(0x2000); got != 42 {
		t.Errorf("OUTPUT_PORT = %d, want 42", got)
	}
	if got := m.ReadWord(0x2002); got != 42 {
		t.Errorf("OUTPUT_PORT+2 = %d, want 42 (copied via MOV macro)", got)
	}
	if got := m.ReadWord(0x2104); got != 1 {
		t.Errorf("STATUS_PORT = %d, want 1", got)
	}
}
