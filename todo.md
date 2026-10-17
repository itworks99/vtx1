# TODO

- [ ] Finish compiler documentation

## For the future projects - lessons learned

- Migrate from VLIW to a more traditional RISC architecture
- Use SystemVerilog for hardware description rather than Verilog
- Migrate from microcode to a more traditional instruction set architecture - current architecture adds a design/code layer that can be avoided for the simplification of the architecture
  - APPLIED to VTX1 (2026-10): the microcode layer is retired from the ISA; transcendentals moved to software, DIV/MOD/SQRT and vector reductions became dedicated multi-cycle FSM units, unsigned ops were deleted. The RTL microcode files are deleted with the cpu_core port.
