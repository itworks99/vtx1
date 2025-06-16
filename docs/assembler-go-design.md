## VTX1 Assembler in Go

This document outlines the structure and design for reimplementing the VTX1 assembler in Go.

### Project Structure

```
src/
  assembler-go/
    cmd/
      vtxasm/
        main.go           # Entry point for the assembler CLI
    internal/
      lexer/
        lexer.go          # Lexical analyzer
        token.go          # Token definitions
        scanner.go        # Character scanning
      parser/
        parser.go         # Parser implementation
        ast.go            # Abstract Syntax Tree definitions
        error.go          # Error handling
      codegen/
        codegen.go        # Code generator
        binary.go         # Binary output format
      symboltable/
        symboltable.go    # Symbol table for labels & variables
      utils/
        errorhandling.go  # Common error handling utilities
        logger.go         # Logging utilities
    pkg/
      vliw/               # VLIW instruction utilities
        packer.go
        unpacker.go
      ternary/            # Balanced ternary utilities
        ternary.go        # Ternary manipulation functions
      instruction/
        instruction.go    # Instruction definitions
        opcodes.go        # Opcode definitions
    test/
      testdata/           # Test assembly files
        hello_world.asm
        ternary_math.asm
        vliw_example.asm
    go.mod                # Go module definition
    go.sum                # Go module checksum
```

### Components

1. **Lexer**: Converts source code into tokens
   - Recognizes keywords, identifiers, literals, operators
   - Handles balanced ternary literals
   - Preserves source location information for error reporting

2. **Parser**: Builds AST from tokens
   - Handles assembly directives (.ORG, .DB, etc.)
   - Creates instruction nodes with operands
   - Manages label declarations
   - Validates basic syntax

3. **Code Generator**: Converts AST to binary
   - Encodes instructions according to VTX1 architecture
   - Resolves labels and symbol references
   - Handles VLIW instruction packing
   - Produces binary output

4. **Symbol Table**: Manages symbols
   - Tracks label declarations and references
   - Handles scope and context
   - Reports undefined or duplicate symbol errors

### Implementation Advantages in Go

1. **Strong Typing**: Define clear structures for each instruction format
2. **Error Handling**: Go's error model makes it easier to propagate and handle errors
3. **Concurrency**: Could parallelize assembly of multiple files
4. **Testing**: Native testing framework with good coverage tools
5. **Cross-platform**: Single binary deployment across operating systems
6. **Performance**: Faster than Python for large assembly files

### Migration Strategy

1. First implement a basic lexer and validate against example files
2. Build the parser, focusing on correctly handling the syntax that caused issues in Python
3. Implement code generation with extensive testing
4. Add command-line interface compatible with existing workflow
5. Finally, run benchmarks and add optimizations

### CLI Interface (Draft)

```
vtxasm [options] input.asm

Options:
  -o, --output FILE      Output binary file (default: input.bin)
  -l, --listing FILE     Generate assembly listing file
  -v, --verbose          Enable verbose output
  -f, --format FORMAT    Output format: binary, hex, or objdump
  -h, --help             Show help information
```
