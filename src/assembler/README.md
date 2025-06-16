# VTX1 Assembler

The VTX1 Assembler is a modern assembly language tool for the VTX1 architecture, implemented in Go. It combines lexical analysis, parsing, and code generation into a single efficient tool.

## Architecture

The VTX1 assembly language is designed for the VTX1 processor architecture, featuring:

- Balanced ternary arithmetic support
- VLIW (Very Long Instruction Word) capabilities
- Register-based operations with a rich instruction set
- Memory-mapped I/O
- Support for both absolute and relative addressing modes

## Features

- Single-pass assembly for efficient compilation
- Comprehensive error reporting with source location information
- Support for macros and include files
- Listing file generation with binary output
- Multiple output formats (binary, hex, objdump)
- Cross-platform support (Windows, macOS, Linux)

## Project Structure

- `/cmd` - Command line interface and entry points
- `/internal` - Core assembler components (lexer, parser, codegen)
- `/pkg` - Reusable packages for VTX1-specific operations
- `/docs` - Documentation in AsciiDoc format
- `/test` - Test files and test data

## Building

```bash
# Build the assembler
go build -o vtxasm ./cmd/vtxasm

# Run tests
go test ./...

# Install locally
go install ./cmd/vtxasm
```

## Usage

```bash
vtxasm [options] input.asm

Options:
  -o, --output FILE      Output binary file (default: input.bin)
  -l, --listing FILE     Generate assembly listing file
  -v, --verbose          Enable verbose output
  -f, --format FORMAT    Output format: binary, hex, or objdump
  -h, --help             Show help information
```

## Assembly Language Syntax

The VTX1 assembly language uses a clean, modern syntax:

```assembly
; Comments start with a semicolon

; Directives
.ORG 0x1000          ; Set origin address
.DB 0x42, 0x43       ; Define bytes
.DW 0xABCD           ; Define word
.INCLUDE "file.inc"  ; Include another file

; Labels
main:
        ; Instructions must be indented
        LD T0, 0x1234        ; Load immediate value
        LD T1, [T0]          ; Load from memory
        ADD T2, T0, T1       ; Add registers
        
        ; VLIW instructions use brackets
        [ADD T3, T0, T1] [SUB T4, T0, T1] [MUL T5, T0, T1]
        
        ; Balanced ternary literals
        LD T0, %+-0          ; Load balanced ternary value
```

## License

This project is licensed under the terms of the license in the project root directory.

## Documentation

For detailed documentation on the VTX1 architecture and assembly language, see the AsciiDoc files in the `docs` directory.
