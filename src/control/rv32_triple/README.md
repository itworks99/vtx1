# rv32_triple — the issue-width-matched binary control core

A triple-issue RV32IM core that is `vtx1_core` with the trit datapath
replaced by two's-complement RV32: the same three-slot structure, the
same dynamic issue rules (at most one memory operation per group, a
control op or halt ends the group, no intra-group RAW/WAW), the same
FSM, the same banked 4-way TCM delivering three consecutive instructions
per cycle with the next PC resolved combinationally during EXEC, and the
same bus fallback. Both cores are written by the same author with the
same design effort.

Purpose: in `task bench-rtl`, a cycle comparison between this core and
`vtx1_core` holds microarchitecture, memory model and design effort
equal, so the remaining difference is the ISA — balanced ternary versus
binary. (PicoRV32 remains in the table as the realistic size-optimized
reference point.)

Implemented subset (anything else traps rather than misexecuting):
RV32I base integer ALU/branch/jump ops, LW/SW (no byte/halfword access),
MUL from M, FENCE as NOP, ECALL/EBREAK halt. Programs come from
`src/assembler/rvsim/rvhex`, encoded from the same instruction lists the
`rvsim` interpreter executes and verifies.
