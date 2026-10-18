# VTX1 Ternary-Weight MAC Accelerator

The specialized unit the architecture was designed to justify: neural
networks quantized to weights in {−1, 0, +1} (BitNet-style) map onto
balanced ternary with **no multipliers** — a weight "multiply" is
select / skip / negate, and negation in balanced ternary is a wire swap.

## The units

| File | What it is |
|---|---|
| `vtx1_tmac.v` | 16-lane ternary-weight MAC: tritwise select/negate into a tree of the ALU's 18-trit adders; wraps mod 3¹⁸ exactly like a chain of CPU `ADD`s, so software and accelerator results are bit-identical |
| `bin_tmac.v` | the matched control: the same function in two's-complement (30-bit signed lanes = the same value range), select/negate + binary adder tree — also multiplier-free |
| `bin_mac8.v` | the industry reference: 16 lanes of signed INT8×INT8 multiply + adder tree, what a conventional binary NPU builds (different function — labeled reference, not a control) |
| `vtx1_tmac_mmio.v` | MMIO wrapper + DMA engine: direct mode (activations at offsets 0x00–0x0F, accumulator at 0x10, weight-word write to 0x11 = one step) and descriptor mode (ACT_BASE/W_BASE/R_BASE at 0x12–0x14, row count + start at 0x15, status poll on 0x15) — the DMA engine streams from the data scratchpad at one word per cycle and writes each row's dot product back |

One 18-trit machine word carries all 16 lane weights in its low 16 trits
(lane *i* at weight 3^i), so a weight word is simply the balanced number
Σ wᵢ·3^i — the assembler's ordinary `.DW` emits it, and packing is a
MUL/ADD loop in software.

## Cell counts (`task synth-mac`, Yosys generic flow, flattened)

| Unit | Function | Cells |
|---|---|---|
| `bin_tmac` | ternary weights × 30-bit activations, binary logic | **4,469** |
| `vtx1_tmac` | ternary weights × 18-trit activations, trit logic (2 bits/trit) | **7,339** |
| `bin_mac8` | INT8 × INT8 multiplier MAC | **9,473** |

Two findings, both honest:

1. **Ternary weights halve the MAC array.** The multiplier-free binary
   unit (4,469) is less than half the INT8 multiplier unit (9,473) while
   processing 3.75× wider activations. This is the BitNet hardware
   argument, quantified — and it holds on *any* substrate.
2. **Trit logic emulated on a binary substrate costs 1.64×** its binary
   equivalent (7,339 vs 4,469) for the identical function. This number is
   the concrete bar for Tier B: real ternary devices (SPICE-characterized
   cells on SkyWater 130nm) must beat ~1.6× to win on area — plus the
   storage side, where a weight word is 16 trits against 32 bits
   (log₂3 ≈ 1.585 bits of information per trit).

## System demo (`task test-tmac`)

`examples/tmac_matvec.asm` computes a 16×16 ternary-weight matrix-vector
product twice on the full SoC — a software MUL/ADD loop, then the
accelerator through MMIO — and compares the vectors itself (the ISA's
negative-address MMIO space, ratified today: peripheral *P* register *k*
lives at word address −(P·0x10000 + k); TMAC is peripheral 6).

| Phase | Cycles | vs software |
|---|---|---|
| software matvec (MUL/ADD loop) | 2,489 | 1× |
| accelerator, direct MMIO mode | 605 | 4.1× |
| accelerator, **DMA-fed from the scratchpad** | **86** | **29×** |

The DMA number includes everything: the 4-register descriptor setup over
MMIO, the stream itself (16 activation reads + per row one weight read
and one result write, at scratchpad speed), and the status poll. The
MMIO-mode row remains in the table as the measured cost of register-level
transport — 97% of its 605 cycles.

Caveats: results are written into the data scratchpad, which owns them
until a CPU store overwrites that address (external memory is not
updated — fine for CPU-consumed results, stated for everything else);
the scratchpad gains a second read port and a shared write port for the
accelerator, with the CPU write-through snoop taking priority; rows are
currently one weight word (16 lanes) wide — tiling wider layers across
multiple steps per row is the designed extension (the `step` interface
already accumulates).
