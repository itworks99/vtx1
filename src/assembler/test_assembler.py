#!/usr/bin/env python3
"""
VTX1 Assembler Test Suite

This script tests the VTX1 assembler by:
1. Running the assembler on example files
2. Verifying that assembly succeeds without errors
3. Checking specific features of the generated binary

Usage:
    python test_assembler.py
"""

import os
import sys
import unittest
import tempfile
import subprocess
import binascii
from pathlib import Path

# Fix import paths to match vtx1_asm.py
script_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(script_dir)  # Add current directory
sys.path.append(os.path.join(script_dir, 'lexer'))  # Add lexer subdirectory
sys.path.append(os.path.join(script_dir, 'parser'))  # Add parser subdirectory
sys.path.append(os.path.join(script_dir, 'codegen'))  # Add codegen subdirectory

# Import the assembler components
from vtx1_asm import Assembler
from lexer.vtx1_lexer import Lexer, TokenType, Token
from parser.vtx1_parser import Parser, ASTNode, NodeType
from codegen.vtx1_codegen import CodeGenerator

class AssemblerTestCase(unittest.TestCase):
    """Base test case for assembler tests with utility methods"""

    def setUp(self):
        """Set up test environment"""
        # Get paths
        self.assembler_dir = os.path.dirname(os.path.abspath(__file__))
        self.examples_dir = os.path.join(self.assembler_dir, "examples")
        self.temp_dir = tempfile.mkdtemp()

        # Create assembler instance
        self.assembler = Assembler(verbose=True)

    def assemble_file(self, input_path, output_path=None):
        """Assemble a file and return success status"""
        if output_path is None:
            output_path = os.path.join(self.temp_dir,
                            os.path.basename(input_path) + ".bin")

        return self.assembler.assemble(input_path, output_path)

    def get_binary_content(self, file_path):
        """Read binary file and return contents"""
        with open(file_path, 'rb') as f:
            return f.read()


class ExampleFileTests(AssemblerTestCase):
    """Test the assembler on example files"""

    def test_hello_world(self):
        """Test assembling the hello world example"""
        input_path = os.path.join(self.examples_dir, "hello_world.asm")
        output_path = os.path.join(self.temp_dir, "hello_world.bin")

        # Assemble the file
        success = self.assemble_file(input_path, output_path)

        # Check if assembly was successful
        self.assertTrue(success, "Failed to assemble hello_world.asm")

        # Check if output file exists
        self.assertTrue(os.path.exists(output_path), "Output file not created")

        # Check if file is not empty
        binary_data = self.get_binary_content(output_path)
        self.assertGreater(len(binary_data), 0, "Output binary is empty")

    def test_ternary_math(self):
        """Test assembling the ternary math example"""
        input_path = os.path.join(self.examples_dir, "ternary_math.asm")
        output_path = os.path.join(self.temp_dir, "ternary_math.bin")

        # Assemble the file
        success = self.assemble_file(input_path, output_path)

        # Check if assembly was successful
        self.assertTrue(success, "Failed to assemble ternary_math.asm")

        # Check if output file exists
        self.assertTrue(os.path.exists(output_path), "Output file not created")

        # Check if file is not empty
        binary_data = self.get_binary_content(output_path)
        self.assertGreater(len(binary_data), 0, "Output binary is empty")

    def test_vliw_example(self):
        """Test assembling the VLIW example"""
        input_path = os.path.join(self.examples_dir, "vliw_example.asm")
        output_path = os.path.join(self.temp_dir, "vliw_example.bin")

        # Assemble the file
        success = self.assemble_file(input_path, output_path)

        # Check if assembly was successful
        self.assertTrue(success, "Failed to assemble vliw_example.asm")

        # Check if output file exists
        self.assertTrue(os.path.exists(output_path), "Output file not created")

        # Check if file is not empty
        binary_data = self.get_binary_content(output_path)
        self.assertGreater(len(binary_data), 0, "Output binary is empty")


class AssemblerComponentTests(AssemblerTestCase):
    """Test individual components of the assembler"""

    def test_lexer(self):
        """Test the lexer component"""
        lexer = Lexer()
        tokens = lexer.tokenize("ADD R0, R1, R2")

        # Check if tokens are generated
        self.assertGreater(len(tokens), 0, "No tokens generated")

        # Check if first token is an ALU operation
        self.assertEqual(tokens[0].type, TokenType.ALU_OP,
                        "First token should be an ALU_OP")

        # Check the opcode value
        self.assertEqual(tokens[0].value, "ADD",
                        "Opcode should be ADD")

    def test_parser(self):
        """Test the parser component"""
        lexer = Lexer()
        tokens = lexer.tokenize("ADD R0, R1, R2")

        # Create parser with tokens
        parser = Parser(tokens)
        ast = parser.parse()

        # Check if AST was generated
        self.assertIsNotNone(ast, "No AST generated")

        # Check if the AST has nodes
        self.assertTrue(hasattr(ast, 'children') or hasattr(ast, 'nodes'),
                      "AST has no instruction nodes")

    def test_code_generator(self):
        """Test the code generator component"""
        lexer = Lexer()
        tokens = lexer.tokenize("ADD R0, R1, R2")

        # Create parser with tokens
        parser = Parser(tokens)
        ast = parser.parse()

        codegen = CodeGenerator()
        binary = codegen.generate(ast)

        # Check if binary was generated
        self.assertGreater(len(binary), 0, "No binary generated")


class FeatureTests(AssemblerTestCase):
    """Test specific features of the VTX1 assembler"""

    def test_ternary_literals(self):
        """Test balanced ternary literal handling"""
        # Create a temporary assembly file with ternary literals
        test_asm = os.path.join(self.temp_dir, "test_ternary.asm")
        with open(test_asm, "w") as f:
            f.write("""
            ; Test ternary literals
            .ORG 0x1000
            LD R0, %+0-+    ; Load balanced ternary literal
            LD R1, %+-0     ; Another ternary literal
            ADD R2, R0, R1  ; Add them
            """)

        # Assemble the file
        output_path = os.path.join(self.temp_dir, "test_ternary.bin")
        success = self.assemble_file(test_asm, output_path)

        # Check if assembly was successful
        self.assertTrue(success, "Failed to assemble ternary literals")

    def test_vliw_instructions(self):
        """Test VLIW instruction packing"""
        # Create a temporary assembly file with VLIW instructions
        test_asm = os.path.join(self.temp_dir, "test_vliw.asm")
        with open(test_asm, "w") as f:
            f.write("""
            ; Test VLIW instruction packing
            .ORG 0x1000
            [ADD R0, R1, R2] [SUB R3, R4, R5] [MUL R6, R7, R8]
            """)

        # Assemble the file
        output_path = os.path.join(self.temp_dir, "test_vliw.bin")
        success = self.assemble_file(test_asm, output_path)

        # Check if assembly was successful
        self.assertTrue(success, "Failed to assemble VLIW instructions")

    def test_symbols_and_labels(self):
        """Test symbol resolution and label handling"""
        # Create a temporary assembly file with symbols and labels
        test_asm = os.path.join(self.temp_dir, "test_symbols.asm")
        with open(test_asm, "w") as f:
            f.write("""
            ; Test symbols and labels
            .ORG 0x1000
            start:
                LD R0, 0x1234
                JMP end
            middle:
                ADD R0, R0, 1
            end:
                JMP middle
            """)

        # Assemble the file
        output_path = os.path.join(self.temp_dir, "test_symbols.bin")
        success = self.assemble_file(test_asm, output_path)

        # Check if assembly was successful
        self.assertTrue(success, "Failed to assemble with symbols and labels")


class ErrorHandlingTests(AssemblerTestCase):
    """Test error handling in the assembler"""

    def test_syntax_error(self):
        """Test handling of syntax errors"""
        # Create a temporary assembly file with syntax errors
        test_asm = os.path.join(self.temp_dir, "test_syntax_error.asm")
        with open(test_asm, "w") as f:
            f.write("""
            ; Test syntax error handling
            .ORG 0x1000
            ADD R0, , R1  ; Missing operand
            """)

        # Assemble the file
        output_path = os.path.join(self.temp_dir, "test_syntax_error.bin")
        success = self.assemble_file(test_asm, output_path)

        # Check that assembly failed due to syntax error
        self.assertFalse(success, "Assembly should fail with syntax error")

    def test_undefined_symbol(self):
        """Test handling of undefined symbols"""
        # Create a temporary assembly file with undefined symbol
        test_asm = os.path.join(self.temp_dir, "test_undef_symbol.asm")
        with open(test_asm, "w") as f:
            f.write("""
            ; Test undefined symbol handling
            .ORG 0x1000
            JMP nonexistent_label  ; Jump to undefined label
            """)

        # Assemble the file
        output_path = os.path.join(self.temp_dir, "test_undef_symbol.bin")
        success = self.assemble_file(test_asm, output_path)

        # Check that assembly failed due to undefined symbol
        self.assertFalse(success, "Assembly should fail with undefined symbol")


if __name__ == "__main__":
    unittest.main()
