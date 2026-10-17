# VTX1 Tier A Benchmarks

Architecture-level benchmark kernels for evaluating balanced ternary against
a binary control architecture **without any hardware**: the metrics below are
substrate-independent (they would be identical on FPGA, on silicon, or in
simulation), so they are measured with the ISA simulator.

## Running

From the repository root:

```sh
task bench
```

or manually:

```sh
go build -C src/assembler -o ../../build/vtxasm .
go build -C src/assembler -o ../../build/vtxsim ./vtxsim
./build/vtxasm -f memh -o build/dot_product.memh src/assembler/benchmarks/dot_product.asm
./build/vtxsim build/dot_product.memh -dump 0x1200:1
```

Every kernel's result is also asserted by `go test ./cmd/ -run TestBenchmarkKernels`
in `src/assembler`, so the numbers cannot silently rot.

## Kernels

| Kernel | What it exercises | Result (verified) |
|---|---|---|
| `dot_product` | signed multiply-accumulate loop, base+index addressing | -107 |
| `sign_count` | three-way classification: one `CMP` yields the full sign (-1/0/+1) | 6 neg / 4 zero / 6 pos |
| `abs_sum` | sum of absolute differences, microcoded `ABS` | 114 |

## Current VTX1 numbers (dynamic, 16-element inputs)

| Kernel | Instructions | ALU | Memory | Control | Loads | Stores | Branches (taken) |
|---|---|---|---|---|---|---|---|
| `dot_product` | 103 | 48 | 38 | 16 | 33 | 1 | 16 (15) |
| `sign_count`  | 126 | 48 | 25 | 52 | 17 | 3 | 42 (27) |
| `abs_sum`     | 119 | 64 | 38 | 16 | 33 | 1 | 16 (15) |

Code size: count the emitted words (`vtxasm -f objdump`); one word = 18 trits
= 18·log2(3) ≈ 28.5 bits of information for normalized comparison.

## Methodology for the binary comparison (next step)

The control core is RV32I (e.g. PicoRV32 for cycle-level runs). Fairness
rules, following the equal-information convention from the ternary-vs-binary
literature:

1. Same algorithm, same inputs, idiomatic assembly on both sides.
2. Compare dynamic instruction counts, memory traffic, and code size in
   information units (trits x log2(3) vs bits), not raw digit counts.
3. Cycle counts come later from the RTL (Icarus) on both sides; the ISA
   simulator deliberately does not model timing.

The RISC-V side needs a cross toolchain or emulator (not yet in this
repository); the C reference implementations of these kernels are a
straightforward port of the assembly above.

Points of interest already visible: `sign_count` is where balanced ternary
should shine (one CMP produces a three-valued result a binary ISA needs two
compares or arithmetic tricks for), while `dot_product` is the neutral
workload.
