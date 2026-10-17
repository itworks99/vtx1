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
- Full-SoC synthesis is parked behind `task synth-soc` (marked BLOCKED)
  until the native-trit ALU lands. Also note: cache_controller clears its
  data arrays in reset, which forces flip-flop mapping (~4.6k FFs) instead
  of RAM inference - worth fixing in the same pass.

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