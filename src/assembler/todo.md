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
    7. [ ] Implement VLIW instruction packing and validation (3 ops per word, resource conflict detection)
    8. [ ] Implement encoding for balanced ternary and 36/96-bit word formats (stub, all output is currently 8-bit binary)
    9. [ ] Implement support for all assembler directives (.ORG, .DB, .DW, .INCLUDE, etc.) (stub for .INCLUDE, others supported)
    10. [ ] Add/expand support for output formats: binary, hex, objdump (hex/objdump are placeholders)
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
- Next step: Implement macro expansion and include file handling in the assembler logic.

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

## Code Generation Note

- A new `CodeGenerator` has been implemented and integrated into the assembler pipeline.
- The assembler can now emit machine code for `NOP` and `ADD` instructions as a proof of concept.
- Next steps: Expand the code generator to support all instruction types, VLIW packing, and directive handling. 