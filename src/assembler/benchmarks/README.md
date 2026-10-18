# VTX1 Tier A Benchmarks: Balanced Ternary vs Binary (RV32IM)

Architecture-level comparison of VTX1 against a binary RISC-V control core,
measured **without any hardware**: the metrics are substrate-independent
(identical on FPGA, silicon, or in simulation), so both sides run in
ISA-level simulators with identical instrumentation:

- **VTX1**: assembly in this directory, assembled by `vtxasm`, executed by
  `vtxsim` (`src/assembler/sim`)
- **RV32IM**: the same kernels in idiomatic RISC-V assembly, executed by
  `rvsim` (`src/assembler/rvsim`), which shares vtxsim's `Stats` counting
  rules exactly

Both sides are hand-written assembly over the same 16-element data sets, and
every result is asserted by `go test` on both sides, so the table cannot
silently rot.

## Results (dynamic, 16-element inputs)

| Kernel | VTX1 instr | RV32IM instr | Δ | VTX1 loads/stores | RV loads/stores | Branches (taken) — both |
|---|---|---|---|---|---|---|
| `dot_product` | **103** | 123 | **−16%** | 33 / 1 | 32 / 1 | 16 (15) |
| `sign_count`  | **110** | 113 | −3% | 17 / 3 | 16 / 3 | 42 (27) |
| `abs_sum`     | **119** | 171 | **−30%** | 33 / 1 | 32 / 1 | 16 (15) |
| `sign_sum`    | **103** | 155 | **−34%** | 33 / 1 | 32 / 1 | 16 (15) |

Static code size (instructions, and stored bits: VTX1 word = 36 bits at
2 bits/trit; RV32I instruction = 32 bits):

| Kernel | VTX1 instr | RV32IM instr | VTX1 bits | RV bits | Δ bits |
|---|---|---|---|---|---|
| `dot_product` | **13** | 18 | **468** | 576 | −19% |
| `sign_count`  | **20** | 23 | 720 | 736 | −2% |
| `abs_sum`     | **14** | 21 | **504** | 672 | −25% |
| `sign_sum`    | **13** | 20 | **468** | 640 | −27% |

On a true ternary storage medium a VTX1 word is 18 trits (≈28.5 bits of
information), which would widen the code-size advantage further; the table
uses the 2-bit-encoded storage cost, the honest number for binary substrates.

## Cycle-level results (real RTL, both sides)

`task bench-rtl` runs the same kernels on both **RTL cores** in Icarus,
with every result asserted by the testbench:

- **VTX1**: the triple-issue trit-native core (`vtx1_core` + banked TCM),
  fetching from the tightly-coupled instruction RAM (1 cycle, 3 words) with
  data loads/stores through a synchronous RAM with a 2-cycle ready handshake
- **PicoRV32** (`src/third_party/picorv32`, YosysHQ, unmodified): the same
  kernels as real RV32IM machine code (encoded by `rvsim/rvhex` from the
  same instruction lists the interpreter executes), on PicoRV32's **fastest**
  native configuration — single-cycle combinational memory, `ENABLE_FAST_MUL`
  (1-cycle multiply), `BARREL_SHIFTER` (1-cycle shifts)

| Kernel | VTX1 cycles (CPI) | PicoRV32 cycles | Δ cycles |
|---|---|---|---|
| `dot_product` | **187** (1.81) | 517 | **−64%** |
| `sign_count`  | **151** (1.37) | 435 | **−65%** |
| `abs_sum`     | **203** (1.70) | 613 | **−67%** |
| `sign_sum`    | **187** (1.81) | 565 | **−67%** |

The ~3× cycle advantage has two multiplicative sources: the ternary ISA
retires fewer instructions (the Tier A table above, −3%..−34%), and the
microarchitecture issues up to 3 independent words per cycle while PicoRV32
is a deliberately small serial core (its documented CPI is ~4, which these
numbers reproduce). Both factors are real properties of the two designs,
but only the first is an *ISA* claim — see caveat 5.

## Reading the results

- **`sign_sum` (−34%) and `abs_sum` (−30%)** are the designed showcases, and
  they deliver: VTX1's `CMP` produces the three-valued sign as a value in one
  instruction where RV32IM needs `slt`+`slt`+`sub`, and `ABS` is one
  instruction versus the branchless `srai`+`xor`+`sub` idiom.
- **`dot_product` (−16%)** is the neutral workload; the win comes from
  `[base+index]` addressing (one index increment serves both arrays) versus
  RISC-V's two pointer bumps per iteration.
- **`sign_count` (−3%)** is deliberately structurally identical on both ISAs
  (two branches per element) — parity here is evidence the comparison is
  fair, not rigged.
- Wide balanced immediates also show up in setup code: loading an address
  like 0x1100 is one `LEA` on VTX1 (±29524 direct range) but `lui`+`addi`
  on RV32I (12-bit immediates).

## Caveats (read before quoting)

1. The instruction-count tables are **architecture** metrics (dynamic
   instruction count, memory traffic, code size), not performance: the
   simulators do not model cycles. The cycle table is measured on RTL.
2. The RV32IM baseline is granted the M extension (VTX1 has MUL in hardware)
   and is hand-written by the same author as the VTX1 kernels — idiomatic,
   but not compiler-optimized output.
3. Small kernels, one data size. The deltas measure ISA expressiveness for
   signed/three-valued arithmetic, which is exactly VTX1's design thesis —
   workloads outside that thesis (e.g. bit manipulation, unsigned hashing)
   would favor binary.
4. Minor asymmetry: VTX1 kernels load the element count from memory
   (1 extra load) while the RV versions use a build-time end pointer.
5. The cycle comparison pits a triple-issue core with an instruction TCM
   against a size-optimized serial core: PicoRV32 is the realistic *small*
   binary control core (~2–3k cells), not a performance peer — VTX1's CPU
   synthesizes to ~100k cells (three ALUs with multiplier trees). The
   memory models are stated above and deliberately favor PicoRV32
   (combinational 1-cycle memory vs VTX1's 2-cycle data handshake). No
   timing closure was run on either core, so cycles are not seconds: a
   per-ISA comparison at equal issue width would need a triple-issue RV32
   core as the control.

## Running

```sh
task bench                                   # instruction counts (simulators)
task bench-rtl                               # cycle counts (RTL, both cores)
cd src/assembler && go test ./rvsim/ -run TestComparisonStats -v   # RV side
go test ./cmd/ -run TestBenchmarkKernels     # VTX1 results asserted
go test ./rvsim/                             # RV32IM results asserted
```

## Kernels

| Kernel | What it exercises | Result (asserted both sides) |
|---|---|---|
| `dot_product` | signed multiply-accumulate, base+index addressing | −107 |
| `sign_count` | three-way classification by branching (neutral) | 6 / 4 / 6 |
| `abs_sum` | sum of absolute differences | 114 |
| `sign_sum` | sum of three-valued signs — the `CMP` showcase | −2 |
