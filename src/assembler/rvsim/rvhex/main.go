// rvhex emits the Tier A kernels as RV32IM memory images for the PicoRV32
// RTL testbench: one 32-bit word per line ($readmemh format), code at byte
// address 0 (PicoRV32's reset vector), input vectors at their kernel
// addresses (0x1100/0x1140), the result area zeroed.
//
// Usage: rvhex <outdir>
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kvany/vtx1/assembler/rvsim"
)

const imageWords = 4096 // matches the testbench RAM (16 KB)

func emit(outdir, name string, prog []rvsim.Instr, data map[uint32]int32) error {
	code, err := rvsim.Encode(prog)
	if err != nil {
		return err
	}
	img := make([]uint32, imageWords)
	for i, w := range code {
		img[i] = w
	}
	for addr, v := range data {
		img[addr/4] = uint32(v)
	}
	var sb []byte
	for _, w := range img {
		sb = append(sb, fmt.Sprintf("%08x\n", w)...)
	}
	path := filepath.Join(outdir, "rv_"+name+".hex")
	if err := os.WriteFile(path, sb, 0o644); err != nil {
		return err
	}
	fmt.Printf("%s: %d instructions\n", path, len(code))
	return nil
}

func vectorData() map[uint32]int32 {
	d := make(map[uint32]int32)
	for i, v := range rvsim.VecA {
		d[uint32(0x1100+4*i)] = v
	}
	for i, v := range rvsim.VecB {
		d[uint32(0x1140+4*i)] = v
	}
	return d
}

func signData() map[uint32]int32 {
	d := make(map[uint32]int32)
	for i, v := range rvsim.SignData {
		d[uint32(0x1100+4*i)] = v
	}
	return d
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: rvhex <outdir>")
		os.Exit(1)
	}
	outdir := os.Args[1]

	kernels := []struct {
		name  string
		build func() ([]rvsim.Instr, error)
		data  func() map[uint32]int32
	}{
		{"dot_product", rvsim.DotProduct, vectorData},
		{"sign_count", rvsim.SignCount, signData},
		{"abs_sum", rvsim.AbsSum, vectorData},
		{"sign_sum", rvsim.SignSum, vectorData},
	}
	for _, k := range kernels {
		prog, err := k.build()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", k.name, err)
			os.Exit(1)
		}
		if err := emit(outdir, k.name, prog, k.data()); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", k.name, err)
			os.Exit(1)
		}
	}
}
