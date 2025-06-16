# VTX1 Assembler Prototype

## DEPRECATED

**This Python implementation of the VTX1 assembler is DEPRECATED and for reference only.**

This directory contains the prototype version of the VTX1 assembler, implemented in Python. It was created as a proof-of-concept and initial exploration of the VTX1 assembly language design and semantics.

## Purpose

The code in this directory is maintained for:
- Historical reference
- Understanding the original design decisions
- Testing and validation of the assembly language specification

## Limitations

The Python prototype has several known limitations:
- Parser stability issues with certain syntax constructs
- Performance constraints with larger assembly files
- Limited error reporting and recovery
- Challenges with lexical analysis of certain constructs

## Future Direction

A new implementation of the VTX1 assembler is being developed in Go, which addresses these limitations and provides:
- Improved performance
- Better error handling and reporting
- More robust parsing of complex syntax
- Cross-platform binary distribution

See `/docs/assembler-go-design.md` for details on the new Go-based implementation.

## Structure

- `lexer/` - Lexical analysis components
- `parser/` - Parser implementation
- `codegen/` - Code generation and binary output
- Other supporting modules and examples

## Usage

While deprecated, the prototype can still be executed for reference:

```bash
python vtx1_asm.py input.asm [-o output.bin] [-l listing.txt] [-v]
```

Please refer to the Go implementation for production use.
