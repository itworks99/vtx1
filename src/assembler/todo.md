# VTX1 Assembler Implementation Checklist

1. [x] **Architecture & Design**
    1. [x] Review and finalize overall assembler architecture (multi-pass, modular design)
    2. [x] Document processing flow and module responsibilities
    3. [x] Define clear interfaces between lexer, parser, codegen, and error handling

2. [x] **Grammar & Parsing**
    1. [x] Review and refine ANTLR grammar (`grammar/vtx1_grammar.g4`) for completeness and correctness
    2. [x] Add/expand grammar rules for all VTX1 instructions, directives, and edge cases
    3. [x] Implement or improve parse tree to AST conversion (if needed for codegen)
    4. [x] Add/expand support for macros and include files in the grammar
    5. [x] Add/expand support for error recovery in the parser

3. [x] **Lexical Analysis**
    1. [x] Ensure all VTX1 tokens (mnemonics, registers, literals, directives) are covered in the lexer rules
    2. [x] Add/expand support for comments, whitespace, and line endings
    3. [x] Add/expand support for balanced ternary literals and custom number formats

4. [x] **Symbol Table & Label Management**
    1. [x] Implement robust symbol table for label definition, lookup, and forward references
    2. [x] Add/expand support for global/local symbols and scoping
    3. [x] Add/expand support for symbol redefinition and error reporting

5. [ ] **Code Generation**
    1. [x] Implement code generation for ALU, MEM, CTRL, SYS instructions
    2. [x] Add opcode map for all VTX1 instruction types (ALU, MEM, CTRL, VEC, FP, SYS, COMPLEX, COMPLEX_VEC, COMPLEX_MEM, COMPLEX_SYS)
    3. [x] Add stubs for VEC, FP, COMPLEX, COMPLEX_VEC, COMPLEX_MEM, COMPLEX_SYS in codegen (emitInstruction)
    4. [x] Implement encoding logic for VEC instructions (VA, VT, VB supported)
    5. [x] Implement encoding logic for FP instructions (FA, FT, FB supported)
    6. [x] Implement encoding logic for COMPLEX instructions (unary and binary forms supported)
    7. [x] Implement VLIW instruction packing and validation (3 ops per bundle, max 1 memory / 1 control op, write-write conflicts are errors, read-after-write hazards are warnings)
    8. [x] Implement trit-native encoding: every instruction is one 18-trit word (--wordsize=36 for the word stream, 108 for VLIW bundles; see ISA Note below)
    9. [ ] Implement support for all assembler directives (.ORG, .DB, .DW, .DT, .SPACE, .EQU, .INCLUDE implemented; .ALIGN and .SECTION are parsed but not yet applied in codegen)
    10. [x] Add/expand support for output formats: binary, hex (Intel HEX with checksums), objdump (real trit-level disassembler)
    11. [ ] Add/expand support for symbol table and debug info output (stub only)
    12. [x] Implement encoding logic for COMPLEX_VEC instructions (VA, VT, VB supported)
    13. [x] Implement encoding logic for COMPLEX_MEM, COMPLEX_SYS instructions (system instructions, no operands)

6. [ ] **Error Handling & Reporting**
    1. [ ] Implement comprehensive error reporting with line/column/source context
    2. [ ] Add/expand support for multiple error types: syntax, semantic, symbol, range, value
    3. [ ] Add/expand support for error recovery and continued parsing after errors
    4. [ ] Add/expand support for warnings and suggestions

7. [ ] **Command-Line Interface (CLI)**
    1. [ ] Refine CLI argument parsing and help output
    2. [ ] Add/expand support for all documented CLI options
    3. [ ] Add/expand support for batch processing and scripting
    4. [ ] Add/expand support for verbose and quiet modes

8. [ ] **Testing & Validation**
    1. [ ] Expand unit tests for lexer, parser, codegen, and error handling
    2. [ ] Add/expand integration tests for end-to-end assembly
    3. [ ] Add/expand regression tests for known bugs and edge cases
    4. [ ] Add/expand test coverage for all instruction types and directives
    5. [ ] Add/expand tests for error and warning reporting
    6. [ ] Add/expand tests for output formats and listing generation

9. [ ] **Documentation**
    1. [ ] Update and expand README with usage, features, and examples
    2. [ ] Update and expand AsciiDoc documentation (architecture, syntax, encoding, etc.)
    3. [ ] Add/expand code comments and docstrings
    4. [ ] Add/expand example assembly programs in `examples/`

10. [ ] **Prototype Migration & Cleanup**
    1. [ ] Review Python prototype for any missing features or tests
    2. [ ] Migrate any useful logic or tests from prototype to Go implementation
    3. [ ] Archive obsolete prototype files after migration

11. [ ] **Build & Release**
    1. [ ] Ensure cross-platform builds (Windows, macOS, Linux) work as expected
    2. [ ] Add/expand release packaging and versioning
    3. [ ] Add/expand CI/CD integration for automated builds and tests

12. [ ] **Miscellaneous**
    1. [ ] Review and clean up codebase for unused files, dead code, and TODOs
    2. [ ] Solicit and incorporate user/developer feedback
    3. [ ] Plan for future features (macro system, advanced optimizations, etc.)

---

## Processing Flow and Module Responsibilities

The VTX1 assembler processes input in the following stages:

1. **Lexical Analysis (Lexer/ANTLR)**: Converts source code into tokens using the ANTLR-generated lexer based on `vtx1_grammar.g4`.
2. **Parsing (Parser/ANTLR)**: Parses the token stream into a parse tree using the ANTLR-generated parser. Handles instruction recognition, operand validation, VLIW grouping, and directive processing.
3. **Symbol Table Management**: Tracks label definitions, references, and symbol resolution, including forward references and address calculation.
4. **Code Generation**: Converts the parse tree (or AST) into binary machine code, handling VLIW packing, ternary encoding, and directive processing.
5. **Error Handling**: Collects and reports errors and warnings at each stage, providing line/column/source context.
6. **Output Formatting**: Generates binary, hex, or objdump output, as well as optional listing and debug information.

**Module Responsibilities:**
- `cmd/`: Orchestrates the assembly process, CLI, and high-level flow.
- `grammar/`: Contains ANTLR grammar and generated lexer/parser code.
- `internal/` (future): Will contain symbol table, codegen, and supporting logic.
- `test/`, `examples/`: Contain test cases and sample programs.
- `docs/`: Contains architecture, syntax, and implementation documentation.

---

## Interface Summary

- **Lexer/Parser (ANTLR):**
  - Input: Source code string
  - Output: ANTLR parse tree (`antlr.ParseTree`)
  - Interface: ANTLR-generated Go interfaces/classes

- **Parser → Codegen:**
  - Input: ANTLR parse tree (future: custom AST)
  - Output: Calls to codegen functions to emit binary
  - Interface: Planned Go interface for AST traversal and code emission

- **Symbol Table:**
  - Input: Label definitions, references, and directives from parser
  - Output: Resolved addresses and symbol errors
  - Interface: Go struct with lookup, insert, and resolve methods

- **Error Handling:**
  - Input: Errors from lexer, parser, codegen, and symbol table
  - Output: Aggregated error/warning list with context
  - Interface: Go struct with methods for add, summarize, and report

---

## Grammar Review Note

- The ANTLR grammar in `grammar/vtx1_grammar.g4` covers all major instruction types, operands, directives, and literals for the VTX1 architecture.
- The grammar is well-structured and matches the EBNF specification and implementation requirements.
- Further refinements may be needed as new features or edge cases are discovered during development and testing.

## Grammar Coverage Note

- The grammar covers all documented VTX1 instructions, directives, register types, and literal formats.
- VLIW, ternary, and all directive forms are present.
- Edge cases to monitor: macro syntax (future), complex nested expressions, and error recovery for malformed lines.

## AST Integration Note

- The assembler now constructs a custom AST from the ANTLR parse tree during parsing.
- The AST is available in the CompilationContext for use in code generation and analysis.
- Next step: Refactor codegen to operate on the AST instead of the parse tree.

## Macro & Include Grammar Note

- The grammar now supports macro definition (`.MACRO ... .ENDM`) and macro invocation, as well as include files via `.INCLUDE`.
- Both are implemented in a preprocessing stage (`cmd/preprocess.go`) that runs before lexing: `.INCLUDE` paths resolve relative to the including file with cycle detection; macros support positional parameters, nested invocation (depth-limited), and labelled invocation lines. See `examples/macro_demo.asm`.
- Known limitation: diagnostics after preprocessing refer to line numbers in the expanded source, not the original file.

---

## Parser Error Recovery Note

- A custom ANTLR error listener (`CustomErrorListener`) has been implemented and integrated.
- The parser now captures detailed syntax errors with line/column information and adds them to the central `ErrorManager`.
- This provides more robust diagnostics and lays the foundation for advanced error recovery strategies.

---

## Lexer Token Coverage Note

- A review confirms that the lexer rules in `vtx1_grammar.g4` comprehensively cover all specified VTX1 tokens.
- This includes all mnemonic categories, register types, literal formats (including ternary), and directives.
- The lexer is considered complete for the current language specification.

---

## Whitespace & Comment Handling Note

- The lexer correctly handles whitespace, single-line comments (starting with `;`), and standard line endings (`\n`, `\r\n`).
- The `WHITESPACE` rule correctly uses a hidden channel, which is the standard ANTLR approach.
- This functionality is considered complete.

---

## Ternary Literal Support Note

- The lexer rule `TERNARY: '0t' [+\-0]+;` correctly identifies balanced ternary literals.
- This definition is sufficient for the lexical analysis stage. The semantic validation and conversion to a numeric value will be handled during code generation.
- No further expansion of the lexer rule is needed.

---

## Symbol Table Implementation Note

- A `SymbolTable` has been implemented to manage labels, constants, and forward references.
- The assembler now performs a dedicated symbol pass after parsing to walk the AST, calculate addresses, and define all labels.
- This two-pass system ensures that all symbols are resolved before code generation.

---

## Symbol Redefinition Error Note

- The `Symbol` struct and `SymbolTable` now store the file, line, and column of each symbol's definition.
- The `ASTBuilder` populates this location data during the parsing stage.
- If a symbol is redefined, the assembler will now report a detailed error message, including the locations of both the original and the duplicate definition, significantly improving diagnostics.

---

## ISA Note (trit-native encoding)

- The ISA is now trit-native and single-sourced in `src/assembler/isa/isa.go`; the Verilog constants `src/common/vtx1_isa_generated.vh` are generated from it (`cd src/assembler && go run ./isa/gen`). Full layout documentation: `docs/vliw_encoding.adoc`.
- Every instruction is one 18-trit machine word (36 bits at 2 bits/trit, stored as 5 bytes). A VLIW bundle is 3 words = 54 trits = 108 bits (`VLIW_WIDTH / `VLIW_BYTES match exactly).
- The 4-trit opcode space holds 81 = 3^4 values: 69 assigned, 12 reserved (`isa.ReservedValues`). The all-zero-trit word is NOP.
- MICROCODE IS RETIRED: SIN/COS/TAN/ASIN/ACOS/ATAN/EXP/LOG are software library routines (transcendentals reduce to shift/add/multiply loops and gain nothing from an interpreter); UDIV/UMOD/BLTU/BGEU are deleted (unsigned is meaningless in balanced ternary); ABS is a plain single-cycle ALU op; DIV/MOD/SQRT and the vector reductions are dedicated multi-cycle FSM units that stall the pipeline.
- Memory is word-addressed; immediates, addresses and branch offsets are balanced ternary and therefore naturally signed (branches are PC-relative ±3280 words, direct addresses ±29524, jumps ±2391484; negative addresses reserved for MMIO).
- TZ (index 18) is the architectural zero register; it encodes compare-against-zero branches and plain `[base]` addressing.
- All example programs assemble end to end as valid trit streams (`cmd/integration_test.go` covers every example x word size x format).
- `LD` always reads memory; load-immediate is spelled `LEA` (the effective address of a direct operand is the immediate). The examples use this convention.
- The assembler also emits Verilog memory images (`-f memh`: 36-bit words with @address markers preserving the .ORG layout) for $readmemh and the simulator.

---

## Decoder & Simulator Note

- `src/cpu/vtx1_decoder.v` is the trit-native hardware decoder (combinational slot decoder + 3-slot bundle wrapper), built on the generated headers (`vtx1_isa_generated.vh` defines, `vtx1_isa_decode.vh` per-opcode format/category table). `task test-decoder` assembles hello_world and verifies the RTL decode of every word against the ISA - the toolchain <-> RTL contract test.
- `vtxsim` (src/assembler/vtxsim, engine in src/assembler/sim) is the ISA-level simulator: executes memh images with balanced-ternary semantics and reports dynamic instruction counts, memory traffic and branch behavior. `task sim-hello` runs hello_world ("Hello, VTX1!" lands at 0x2000); `task bench` runs the Tier A kernels in src/assembler/benchmarks (results asserted by `go test -run TestBenchmarkKernels`).
- Provisional ABI decisions to ratify: TA is the link register (JAL/JALR/CALL write it, RET jumps to it); PUSH/POP trap until a stack ABI is defined. (BLTU/BGEU were removed from the ISA.)
- MILESTONE REACHED: `src/cpu/vtx1_core.v` is the trit-native CPU core (v1: multi-cycle word-serial, hardwired control, no microcode). `task test-core` runs hello_world end to end on the RTL: assembler binary -> $readmemh -> execution -> "Hello, VTX1!" at 0x2000, with the retired-word count cross-checked against vtxsim (93 = 93; CPI 2.13).
- vtx1_core v1 notes: fetch model (a) - one word per cycle, bundles execute as sequential slots (parallel bundle issue is a later performance mode); registers/memory hold canonical trit words; the ALU converts trit->binary->trit inline (functionally identical; swapping in native trit datapath logic is a later optimization); DIV/MOD are single-cycle operators in simulation (synthesis needs the iterative divider FSM); SQRT, PUSH/POP and vector/FP ops trap.
- SOC INTEGRATION COMPLETE: the top level is now `src/vtx1_soc.v` (module vtx1_soc, renamed from vtx1_top_consolidated) with vtx1_core as bus master 0 - no adapter needed, the core's memory port carries the bus handshake natively. `task test-soc` executes hello_world through the full chip (CPU -> bus matrix -> memory controller -> external SRAM pins) and verifies the output: 93 retired words (= vtxsim), 911 cycles. The memory controller's byte-era alignment check (addr[1:0]==0) was removed: VTX1 is word-addressed.
- The old binary-format CPU stack (cpu_core, instruction_decoder, microcode_rom/sequencer, forwarding/hazard units, register_file, tcu interface) and its dead testbenches moved to `deprecated/cpu-microcode/`; purge with `git rm -r deprecated` when ready.
- Performance headroom (later): SoC CPI is ~9.8 (bus arbitration + 4-cycle external memory emulation per access); a tightly-coupled instruction RAM or fetch buffering, and parallel VLIW bundle issue, are the obvious wins.

---

## Synthesis Findings (2026-10-18)

- `vtx1_decoder` synthesizes cleanly with Yosys: 5,562 cells, purely
  combinational, CHECK-clean, 17 s (`task synth`).
- `vtx1_core` does NOT practically synthesize: the v1 conversion ALU's
  b2t function (72-bit signed div/mod, 18 iterations, ~20 call sites)
  unrolls into hundreds of wide dividers and ABC never converges. The
  "swap in a native trit datapath later" note is therefore a synthesis
  REQUIREMENT, not an optimization: the next hardware task is a native
  balanced-ternary ALU (per-trit adder cell -> 18-trit ripple/carry-select
  adder; SUB = add(neg); CMP from the subtract sign; value-to-trit
  conversion then only needed at the LEA/immediate boundary, which is
  constant-folded by the assembler anyway... immediates arrive as trits
  in the instruction word, so the datapath can be conversion-free).
- Fixed for synthesis along the way: FIFO memories in spi/i2c/uart were
  written inside async-reset always blocks (Yosys cannot infer RAMs
  there); those blocks now use synchronous reset.
- RESOLVED (same day): the native balanced-ternary ALU landed - see below.
- FULL SOC SYNTHESIS WORKS: 'task synth-soc' completes in ~195 s - 79,122
  cells, ~12.3k FFs, netlist in build/vtx1_soc_synth.v. Getting there also
  flushed out four double-driven nets in the top (cache_wr/cache_wdata
  tied to constants while the bus matrix drives them; debug_access_data
  and uart_dma_tx_data tied while their controllers drive them; jtag_tdo
  tied to ground while the debug controller shifts data onto it) - real
  latent bugs that simulation x-resolution had been masking.
  Remaining from this pass: cache_controller clears its data arrays in
  reset, which forces flip-flop mapping (~4.6k FFs) instead of RAM
  inference - worth fixing eventually.
  RESOLVED (same day): the data/tag arrays moved to a synchronous
  write-only block with no reset (the reset-cleared valid bits gate every
  read, so clearing the arrays was pure waste); the write-enable terms
  are shared wires between the FSM and the array block so the two cannot
  diverge. Yosys now sees plain synchronous RAMs.

---

## Native Balanced-Ternary ALU (2026-10-18)

- `src/cpu/vtx1_alu.v`: vtx1_trit_fa (balanced full-adder cell),
  vtx1_trit_adder (18-trit ripple; negate_b gives subtraction free;
  dropped carry = balanced wrap mod 3^18, matching vtxsim's WrapWord),
  vtx1_trit_mul (18 partial products {-B,0,+B} reduced by an adder tree,
  depth 5), plus tritwise min/max/product/negate, barrel shifter and the
  sign unit. Verified against a binary reference model with 34,612
  directed+random vectors (`task test-alu`).
- vtx1_core v2: conversion-free datapath. Registers, PC, immediates, EAs
  and results are trit words end to end; the decoder emits immediates as
  zero-trit-padded words (balanced ternary needs no sign extension); the
  only trit->binary conversion is t2b toward the memory address bus
  (x3 add chains). The b2t divider chains are gone; the reset vector uses
  a Verilog constant function (elaboration-time only).
- SYNTHESIS UNBLOCKED: the CPU (decoder+ALU+core) synthesizes in 135 s:
  ~37.7k cells, 853 FFs (`task synth`). Previously ABC did not converge
  after 7 CPU-hours.
- Design findings worth keeping:
  - Overflow-correct compare: the adder's dropped carry trit IS the true
    sign when subtraction overflows 18 trits. CMP and branch compares use
    (cout != 0) ? cout : sign(diff) - the balanced-ternary equivalent of
    the binary N^V trick, found by the ALU's randomized testing.
  - NOT and NEG are the same operation in balanced ternary (tritwise
    negation = arithmetic negation). Consider merging them in a future
    ISA revision to free an opcode.
  - Shift/rotate amount semantics RATIFIED: the value of the low 4 trits
    of rs2 (+-40), reduced mod 18. vtxsim and the ALU implement the same
    rule.
- Still trapping in hardware (as designed): SQRT, PUSH/POP, vector/FP
  ops. (DIV/MOD RESOLVED same day: the iterative divider landed - see
  the Iterative Divider section.)

---

## Tightly-Coupled Instruction RAM (2026-10-18)

- `src/memory/vtx1_tcm.v`: a single-cycle synchronous-read instruction RAM
  on a dedicated core fetch port, covering the program region
  (0x1000..0x1FFF, 4096 words). Its write port snoops the CPU's store
  transactions (write-through), so CPU stores into the window stay
  coherent and the boot flow loads code with ordinary stores; DMA/debug
  writes are not snooped. Without the write port Yosys (correctly)
  deleted the never-written array and the netlist TCM was vacuous -
  caught by inspecting the synthesized module, worth remembering: a
  preload-only memory does not survive synthesis.
- Yosys' generic `synth` flow maps every inferred RAM to flip-flops
  (there is no RAM target library), so the full 4096x36 TCM becomes
  ~147k DFFs and swamps the SoC statistics while validating no extra
  logic. `task synth-soc` therefore shrinks the array to 256 words with
  `chparam` for the smoke test; on a real flow (OpenLane/Sky130) the
  array is a RAM macro. The same mechanism is why the cache data arrays
  show up as FFs.
- The core resolves the next PC combinationally during EXEC (branches,
  jumps, RET/JR/JALR included) and issues the TCM read for it in the same
  cycle, consuming the word directly from the TCM's read register next
  cycle - so straight-line TCM code runs EXEC->EXEC at CPI 1 with no fetch
  state at all. A TCM miss (or `imem_hit` tied low) falls back to the
  original bus fetch unchanged; loads/stores still go over the bus.
- Measured on hello_world (93 retired, vtxsim cross-check still exact):
  - core + 1-cycle RAM: 333 -> 175 cycles (CPI 3.58 -> 1.88)
  - full SoC through the bus matrix: 911 -> 310 cycles (CPI 9.8 -> 3.33)
  - remaining cycles are LD/ST bus transactions (~7 cycles each through
    the bus matrix + memory controller) and the 1-cycle refetch after them.
- `cycles` now stops counting in S_HALT (it used to keep incrementing
  after WFI, slightly inflating every earlier CPI number).
- Next performance steps: parallel VLIW bundle issue (fetch 3-slot
  bundles from the TCM), a data fast path or write buffer for LD/ST, then
  the cycle-level comparison against PicoRV32 under the same memory model.

---

## Triple Issue (2026-10-18)

- The core is now v3: three `vtx1_slot` units (decode + single-cycle
  execute of one operation word, combinational against the register file)
  plus dynamic issue logic. Up to three CONSECUTIVE words commit per
  cycle; the hardware enforces the VLIW packing rules itself (max one
  memory op, control/halt ends the group, no intra-group RAW/WAW), so
  co-issued words are provably independent and parallel commit is
  bit-identical to sequential execution. The vtxsim retired-count
  cross-check (93) is preserved exactly, and the assembler's VLIW packing
  becomes the scheduler: a packed bundle is guaranteed to issue in one
  cycle, but plain scalar code multi-issues too whenever it happens to be
  independent.
- The TCM is banked 4-way on addr[1:0] so any three consecutive words
  (arbitrary alignment, branch targets included) read in one cycle;
  rdata0..2 feed the group. The fast path still resolves the next PC
  combinationally during EXEC, so TCM code needs no fetch state: 3 words
  per cycle peak. A bus fetch delivers one word (single issue) - with
  imem_hit low the core is cycle-identical to the single-issue v2.
- hello_world: core+1-cycle RAM 175 -> 136 cycles (CPI 1.88 -> 1.46);
  full SoC 310 -> 271 cycles (CPI 3.33 -> 2.91; was 9.8 before the TCM).
  Remaining cycles are dominated by LD/ST bus transactions.
- ICARUS SENSITIVITY TRAP (cost one debug session, worth remembering):
  a Verilog function that reads module signals NOT passed as arguments
  (e.g. `regword(idx)` reading the register file, or a helper reading
  `addr`) never re-evaluates in a continuous assign or always-@* in
  Icarus - the sensitivity list is computed from the arguments only. The
  v2 core got away with it because its reads sat inside the clocked
  always block. Rule: functions must be pure over their arguments; reads
  of module state go in explicit wires/muxes.

---

## Cycle-Level RTL Comparison vs PicoRV32 (2026-10-18)

- `task bench-rtl`: the four Tier A kernels run on both RTL cores in
  Icarus with results asserted by the testbenches. VTX1 = triple-issue
  core + banked TCM + 2-cycle-handshake data RAM; PicoRV32 (vendored
  unmodified under src/third_party/picorv32, ISC) = its FASTEST config:
  combinational 1-cycle memory, ENABLE_FAST_MUL, BARREL_SHIFTER.
- Results: dot_product 187 vs 517 (-64%), sign_count 151 vs 435 (-65%),
  abs_sum 203 vs 613 (-67%), sign_sum 187 vs 565 (-67%). VTX1 CPI
  1.37-1.81; PicoRV32 reproduces its documented ~4 CPI.
- The RV32 machine code is produced by `rvsim/encode.go` + `rvsim/rvhex`
  from the SAME structural instruction lists the rvsim interpreter
  executes - no external toolchain, and the asserted results validate
  the encoder on real RTL.
- Honest framing (also in benchmarks/README caveat 5): this is a
  triple-issue ~100k-cell core vs a size-optimized ~2-3k-cell serial
  core, memory models stated and favoring PicoRV32, no timing closure
  on either. The ISA claim stays the instruction-count table; the cycle
  table shows the ternary core's microarchitecture is real and fast,
  not that ternary beats binary per se. An issue-width-matched RV32
  control would be the next rigor step.

---

## Issue-Width-Matched RV32 Control (2026-10-18)

- `src/control/rv32_triple/`: rv32_core/rv32_slot/rv32_tcm are vtx1_core
  v3 with the trit datapath swapped for two's-complement RV32IM - same
  three-slot issue logic, same FSM, same banked TCM, same 2-cycle data
  RAM, same author. Cycle differences against it isolate the ISA.
- Results (cycles): dot_product 187 vs 187 (parity), sign_count 151 vs
  151 (parity), abs_sum 203 vs 235 (-14%), sign_sum 187 vs 203 (-8%).
- The lesson, now written into benchmarks/README: issue width buys back
  instruction COUNT (RV32's extra ALU ops hide in empty slots - its CPI
  is lower at equal cycles, and the memory-bound kernels serialize on
  identical load streams), but it cannot buy back DEPENDENT CHAINS. ABS
  and CMP each replace a 3-op dependent sequence with one op on the
  critical path - that is where balanced ternary's cycle win survives a
  fair fight, exactly as the thesis predicts.
- Verilog signedness trap found here: in a mixed signed/unsigned
  conditional expression, >>> evaluates UNSIGNED (the unsigned branch
  poisons the whole expression). SRA now computes in its own signed wire.

---

## Iterative Divider (2026-10-18)

- `src/cpu/vtx1_div.v`: DIV/MOD in hardware, ~22 cycles, built from the
  ALU's vtx1_trit_fa cells on a 20-trit internal datapath (2|R1| reaches
  3B+2, two trits past the operand range; zero-trit padding is free).
- Semantics RATIFIED to match vtxsim (Go rules): truncation toward zero,
  MOD sign follows the dividend, division by zero traps (checked in the
  core before start, so trap_pc points at the DIV).
- Algorithm: magnitudes only (negation is tritwise, so signs are applied
  for free at the end); MSB-first digit recurrence with quotient digits
  in {-2..2} - pure {-1,0,+1} selection leaves a +-1 overshoot at the
  2|R|=B boundary that compounds (3c+1 per iteration), while allowing
  +-2 (the digit 2 = (+1,-1) trit pair, exactly representable) keeps the
  invariant 2|R| <= B exact. The loop naturally yields the ROUND-TO-
  NEAREST quotient (balanced ternary's native rounding) with |R| <= B/2;
  one correction step (Q-1, R+B when R<0) converts to truncation.
- Core integration: DIV/MOD issue alone (group-ending, blocked from
  slots 1/2), run in a new S_DIV state, write back like a load. Verified
  by 6048 random+directed vectors against Verilog's native signed '/'
  and '%' (tb_vtx1_div) and by div_test.asm executing identically in
  vtxsim and on the RTL core (task test-div, in test-all).
- A free byproduct worth remembering: stopping before the correction
  step gives round-to-nearest division - a candidate RDIV opcode for DSP
  work (no binary ISA has a one-instruction nearest division).

---

## Stack ABI (2026-10-18, RATIFIED)

- TB (index 7, "base pointer") is the stack pointer: full-descending,
  word-addressed. PUSH r: TB <- TB-1, mem[TB] <- r (the OLD r, relevant
  for PUSH TB). POP r: r <- mem[TB], TB <- TB+1. POP TB: the loaded
  value supersedes the increment. Programs initialize TB themselves
  (LEA TB, <top>).
- Implemented in vtxsim and the slot/core: PUSH/POP present themselves
  as a store/load plus a same-cycle TB write (the existing commit paths
  handle both); POP has two destinations, so stack ops issue alone
  (group-ending, not issuable from slots 1/2) - same treatment as DIV.
- stack_test.asm (task test-stack, in test-all) runs LIFO order, nested
  CALL/RET with PUSH TA / POP TA, stack balance, and the POP TB edge
  case - identical results in vtxsim and on the RTL core.

---

## Memory-Path Performance Pass (2026-10-18)

- Refetch elimination (both matched cores): a memory op or DIV leaving
  EXEC no longer clears ir_from_imem - the next instruction group was
  fetched during that same EXEC cycle and the TCM output registers hold
  it while imem_en is low, so S_LOAD/S_STORE/S_DIV resume straight into
  EXEC. One cycle saved per memory op / divide; nothing consumes slot
  outputs during the memory states (store data is latched), so it is
  free. Core hello_world 136 -> 109 cycles (CPI 1.17).
- `src/memory/vtx1_scratchpad.v`: the data-side twin of the instruction
  TCM, between the core's memory port and the SoC bus master port.
  In-window READS are served from a local sync RAM (1-cycle handshake);
  WRITES always pass through to the bus (write-through: external memory
  stays the source of truth for DMA/debug/testbench) and are snooped
  into the local copy. Window = low 16K words; MMIO (0xF0000000+)
  untouched; chparam-shrunk for the synth smoke test like the TCM.
- SoC hello_world: 911 (morning) -> 271 -> 244 (refetch) -> 179 cycles
  (CPI 1.92). The remaining gap to the bare core (109) is store traffic
  crossing the bus - by design under write-through.
- bench-rtl after the pass (both matched cores benefit equally):
  dot 153 vs 154, sign_count 131 vs 132, abs_sum 169 vs 202 (-16%),
  sign_sum 153 vs 170 (-10%). Tightening the machines GREW the ISA
  deltas: less fixed overhead makes the dependent-chain difference a
  larger share of what remains.
- Measured and parked: two memory ops per group (banked data memory)
  would halve dot_product's loop but changes a ratified ISA packing rule
  (max 1 mem op per bundle) - an architecture decision for later.
  Software-side bundle packing/unrolling of the kernels also unmeasured.

---

## Ternary-Weight MAC Accelerator (2026-10-18)

- `src/accel/`: the idea from "Ideas Worth Keeping" is now hardware.
  vtx1_tmac = 16 multiplier-free MAC lanes (tritwise select/negate into a
  tree of the ALU's adders; wraps mod 3^18 exactly like CPU ADDs, so
  software and accelerator agree bit for bit). bin_tmac = the matched
  binary control (same function, 30-bit two's-complement). bin_mac8 =
  the INT8 multiplier MAC industry reference. Verified equivalent on 61
  random accumulation runs (tb_tmac).
- Cell counts (task synth-mac, flattened): bin_tmac 4,469 / vtx1_tmac
  7,339 / bin_mac8 9,473. Two honest headlines: (1) ternary WEIGHTS
  halve the MAC array vs the multiplier unit while processing 3.75x
  wider activations - on any substrate; (2) trit logic at 2 bits/trit on
  a binary substrate costs 1.64x its binary equivalent - the concrete
  bar Tier B's real ternary cells must beat.
- NEGATIVE-ADDRESS MMIO RATIFIED: the ISA's "negative addresses are
  MMIO" provision is now implemented - a negative CPU word address maps
  into the bus matrix MMIO window, magnitude m selecting peripheral
  m[19:16] (GPIO 0 ... FLASH 5, TMAC 6) and register m[15:0], so
  peripheral P register k is at address -(P*0x10000 + k). This was dead
  letter before: the CPU's zero-extended address could never reach the
  0x4_xxxx_xxxx MMIO window at all.
- SoC demo (task test-tmac, tb_vtx1_soc_bench markers): 16x16 matvec,
  software 2,489 cycles vs accelerator 605 (4.1x). Only 16 of the 605
  are MAC work (256 MACs); MMIO transport at ~9 cycles/access dominates.
  DMA feeding / a tightly-coupled accelerator port is the next lever.
- tb_vtx1_soc_bench is a new GENERIC full-SoC program harness (+img,
  +exp0..5 at 0x1200, phase timestamps via marker stores to 0x121x).

---

## Ideas Worth Keeping (ternary leverage)

- Ternary CORDIC: CORDIC with the native digit set {-1, 0, +1} can skip
  zero-rotations ("redundant CORDIC"), converging in fewer effective
  iterations. If transcendental throughput ever matters, this is the
  hardware answer: a dedicated CORDIC FSM unit, not microcode.
- Ternary-weight MAC / matrix engine: neural networks quantized to weights
  in {-1, 0, +1} (BitNet-style) map perfectly onto balanced ternary - a
  weight multiply is add/skip/subtract, so a matrix-vector unit needs NO
  multipliers. This is the specialized accelerator that would make VTX1's
  case, and it is what VDOT should grow into.
- The architecture-level ternary wins to measure in Tier A: Booth-free
  multiplication, free negation, truncation-rounds-to-nearest (lower error
  accumulation in MAC chains), symmetric range (no abs(INT_MIN) overflow),
  and native three-way compare (see sign_sum).

---

## Tier A Result (2026-10-18)

- The binary control core exists: `src/assembler/rvsim` is a minimal RV32IM
  interpreter sharing vtxsim's Stats counting rules, with the four kernels
  hand-written in idiomatic RISC-V and results asserted by go test.
- Measured (dynamic instructions, 16-element kernels): sign_sum -34%,
  abs_sum -30%, dot_product -16%, sign_count -3% (deliberately neutral).
  Static code: -19..-27% in stored bits except sign_count (~par).
  Full table and caveats: benchmarks/README.md.
- Verdict: the ternary architecture wins exactly where its thesis says it
  should (three-valued compare as a value, ABS, signed MAC with base+index
  addressing) and holds parity on neutral control flow. This justifies
  proceeding toward Tier B (SPICE-characterized ternary cells) and the
  ternary-weight MAC accelerator.

---

## Visitor Dispatch Note (bug fix)

- The ANTLR parser was generated without the `-visitor` option, so parse tree contexts have no `Accept` methods; `tree.Accept(builder)` silently returned nil for every node and the assembler emitted 0 bytes for any input. `ASTBuilder.Visit` now dispatches on the concrete context type directly.
- Further AST builder fixes: directive names no longer swallow their parameters (`.ORG 0x1000` used to produce the name `.ORG0x1000`, breaking .ORG/.DB/.DW/.EQU handling); `.EQU` parameters keep source order (identifier, value); labelled directives (`name: .DW 10`) no longer lose their label; binary operands (`vec_a+4`) keep their left/right order.