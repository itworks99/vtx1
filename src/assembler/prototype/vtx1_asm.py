#!/usr/bin/env python3
"""
VTX1 Assembler - Main Program

This module integrates the lexer, parser, and code generator to assemble
VTX1 assembly source code into binary machine code.

Usage:
    python vtx1_asm.py input.asm [-o output.bin] [-l listing.txt] [-v]

Options:
    -o, --output FILE     Specify output binary file (default: input.bin)
    -l, --listing FILE    Generate assembly listing file
    -v, --verbose         Enable verbose output
    -h, --help            Show this help message and exit
"""

import sys
import os
import argparse
import time
from typing import List, Dict, Optional, Tuple

# Add module paths for imports
sys.path.append(os.path.join(os.path.dirname(__file__), '../lexer'))
sys.path.append(os.path.join(os.path.dirname(__file__), '../parser'))
sys.path.append(os.path.join(os.path.dirname(__file__), '../codegen'))

# Import components
from vtx1_lexer import Lexer, Token, TokenType
from vtx1_parser import Parser, ASTNode, NodeType, print_ast
from vtx1_codegen import CodeGenerator

class Assembler:
    def __init__(self, verbose=False):
        self.lexer = Lexer()
        self.verbose = verbose
        self.source_lines = []
        self.output_binary = bytearray()
        self.listing_lines = []

    def assemble(self, input_file: str, output_file: str, listing_file: Optional[str] = None) -> bool:
        """Assemble source file to binary"""
        start_time = time.time()

        # Read source file
        try:
            with open(input_file, 'r') as f:
                source_code = f.read()
                self.source_lines = source_code.splitlines()
        except Exception as e:
            print(f"Error reading input file: {e}")
            return False

        if self.verbose:
            print(f"Assembling {input_file}...")
            print(f"Read {len(self.source_lines)} source lines")

        # Run the lexer
        tokens = self.lexer.tokenize(source_code)

        if self.verbose:
            print(f"Generated {len(tokens)} tokens")

        # Run the parser
        parser = Parser(tokens)
        ast = parser.parse()

        if parser.errors:
            print(f"Found {len(parser.errors)} parsing error(s):")
            for error in parser.errors:
                print(f"  {error}")
            return False

        if self.verbose:
            print("Parsing completed successfully")
            if self.verbose > 1:
                print("\nAST structure:")
                print_ast(ast)

        # Run the code generator
        codegen = CodeGenerator(debug=self.verbose)
        binary_code = codegen.generate_code(ast)

        if codegen.errors:
            print(f"Found {len(codegen.errors)} code generation error(s):")
            for error in codegen.errors:
                print(f"  {error}")
            return False

        if codegen.warnings:
            print(f"Found {len(codegen.warnings)} warning(s):")
            for warning in codegen.warnings:
                print(f"  {warning}")

        # Save binary output
        self.output_binary = binary_code
        try:
            with open(output_file, 'wb') as f:
                f.write(binary_code)
        except Exception as e:
            print(f"Error writing output file: {e}")
            return False

        if self.verbose:
            print(f"Generated {len(binary_code)} bytes of machine code")
            print(f"Binary output saved to {output_file}")

        # Generate listing file if requested
        if listing_file:
            self._generate_listing(listing_file, codegen.labels)

        # Print statistics
        end_time = time.time()
        elapsed = end_time - start_time
        print(f"Assembly completed in {elapsed:.3f} seconds")
        print(f"Code size: {len(binary_code)} bytes")
        print(f"Labels defined: {len(codegen.labels)}")

        return True

    def _generate_listing(self, listing_file: str, labels: Dict[str, int]) -> bool:
        """Generate an assembly listing file with address, binary, and source"""
        try:
            with open(listing_file, 'w') as f:
                f.write("VTX1 Assembly Listing\n")
                f.write("=====================\n\n")

                # Write symbol table
                f.write("Symbol Table:\n")
                f.write("------------\n")
                for symbol, address in sorted(labels.items()):
                    f.write(f"{symbol:20s} 0x{address:08X}\n")
                f.write("\n")

                # Write code listing with addresses
                f.write("Code Listing:\n")
                f.write("------------\n")
                f.write("Address    Machine Code                   Source\n")
                f.write("--------  -----------------------------  --------------------------\n")

                # For a proper listing, we'd need to track which source lines
                # correspond to which binary bytes, which is complex and beyond
                # the scope of this implementation. For now, we'll show a simplified listing.

                # Write the binary in chunks with corresponding line numbers
                address = 0
                for i, line in enumerate(self.source_lines):
                    line = line.strip()
                    if not line or line.startswith(';'):
                        # Just show source line for comments/blank lines
                        f.write(f"          {'':<30}  {line}\n")
                    else:
                        # For code lines, show 4-byte chunks (simplification)
                        chunk_size = 4
                        if address < len(self.output_binary):
                            chunk = self.output_binary[address:address+chunk_size]
                            hex_str = ' '.join(f"{b:02X}" for b in chunk)
                            f.write(f"0x{address:06X}  {hex_str:<30}  {line}\n")
                            address += chunk_size
                        else:
                            f.write(f"          {'':<30}  {line}\n")

            if self.verbose:
                print(f"Listing file saved to {listing_file}")

            return True
        except Exception as e:
            print(f"Error writing listing file: {e}")
            return False


def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='VTX1 Assembler')
    parser.add_argument('input', help='Input assembly source file')
    parser.add_argument('-o', '--output', help='Output binary file')
    parser.add_argument('-l', '--listing', help='Generate assembly listing file')
    parser.add_argument('-v', '--verbose', action='count', default=0, help='Enable verbose output')
    args = parser.parse_args()

    # Determine output filename if not specified
    if not args.output:
        base_name = os.path.splitext(args.input)[0]
        args.output = f"{base_name}.bin"

    # Create assembler instance
    assembler = Assembler(verbose=args.verbose)

    # Run the assembler
    result = assembler.assemble(args.input, args.output, args.listing)

    # Return exit code
    sys.exit(0 if result else 1)

if __name__ == "__main__":
    main()
