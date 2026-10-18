# Archived specification chapters (pre-2026 design)

These chapters describe the VTX1 as designed before the 2026 trit-native
rework: a microcode-based machine with binary instruction encodings, an
FPGA prototyping target, and aspirational gate counts. That design was
superseded wholesale — the microcode layer was retired, the ISA became
trit-native end to end, the core was rebuilt twice (conversion-free
datapath, then triple issue), and the evaluation moved from FPGA plans
to measured simulation tiers.

They are kept for history only. The current specification is
`docs/Architecture.adoc`; the current design-decision log is
`src/assembler/todo.md`. The peripheral chapters (`system/`, `boot`,
`debug`) describe RTL that still largely exists but have not been
re-verified against it.
