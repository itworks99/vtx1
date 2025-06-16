#!/usr/bin/env python3
"""
VTX1 Assembler Component Tests

This script tests the individual components of the VTX1 assembler (lexer, parser, codegen)
with inputs that match the expected assembly syntax.

Usage:
    python component_tests.py
"""

import os
import sys
import unittest
import tempfile

# Fix import paths
script_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(script_dir)  # Add current directory
sys.path.append(os.path.join(script_dir, 'lexer'))  # Add lexer subdirectory
sys.path.append(os.path.join(script_dir, 'parser'))  # Add parser subdirectory
sys.path.append(os.path.join(script_dir, 'codegen'))  # Add codegen subdirectory

# Import the assembler components
from vtx1_asm import Assembler
from lexer.vtx1_lexer import Lexer, Token, TokenType
from parser.vtx1_parser import Parser, ASTNode, NodeType
from codegen.vtx1_codegen import CodeGenerator

class ComponentTests(unittest.TestCase):
    """Tests individual components of the VTX1 assembler"""

    def setUp(self):
        """Set up test environment"""
        self.assembler = Assembler(verbose=True)

    def test_lexer_directive(self):
        """Test the lexer with a directive"""
        lexer = Lexer()
        code = ".ORG 0x1000"
        tokens = lexer.tokenize(code)

        # Should have 3 tokens: DIRECTIVE, HEXADECIMAL, EOF
        self.assertEqual(len(tokens), 3, "Expected 3 tokens from directive")
        self.assertEqual(tokens[0].type, TokenType.DIRECTIVE, "First token should be a DIRECTIVE")
        self.assertEqual(tokens[0].value, ".ORG", "Directive value should be .ORG")
        self.assertEqual(tokens[1].type, TokenType.HEXADECIMAL, "Second token should be HEXADECIMAL")
        self.assertEqual(tokens[1].value, "0x1000", "Hexadecimal value should be 0x1000")

    def test_lexer_instruction(self):
        """Test the lexer with a basic instruction"""
        lexer = Lexer()
        code = "LD T0, 0x1234"
        tokens = lexer.tokenize(code)

        # Should have: MEM_OP, GPR, COMMA, HEXADECIMAL, EOF
        self.assertEqual(len(tokens), 5, "Expected 5 tokens from instruction")
        self.assertEqual(tokens[0].type, TokenType.MEM_OP, "First token should be MEM_OP")
        self.assertEqual(tokens[0].value, "LD", "Instruction mnemonic should be LD")
        self.assertEqual(tokens[1].type, TokenType.GPR, "Second token should be GPR")
        self.assertEqual(tokens[2].type, TokenType.COMMA, "Third token should be COMMA")

    def test_lexer_comment(self):
        """Test the lexer with a comment"""
        lexer = Lexer()
        code = "; This is a comment"
        tokens = lexer.tokenize(code)

        self.assertEqual(len(tokens), 2, "Expected 2 tokens from comment")
        self.assertEqual(tokens[0].type, TokenType.COMMENT, "First token should be COMMENT")

    def test_lexer_vliw(self):
        """Test the lexer with a VLIW instruction"""
        lexer = Lexer()
        code = "[ADD T0, T1, T2] [SUB T3, T4, 0x1]"
        tokens = lexer.tokenize(code)

        # Should have many tokens including LBRACKET and RBRACKET
        self.assertGreater(len(tokens), 10, "Expected multiple tokens from VLIW")
        self.assertEqual(tokens[0].type, TokenType.LBRACKET, "First token should be LBRACKET")
        self.assertEqual(tokens[1].type, TokenType.ALU_OP, "Second token should be ALU_OP")

    def test_parser_directive(self):
        """Test the parser with a directive"""
        lexer = Lexer()
        # Add indentation to the directive
        code = "        .ORG 0x1000"
        tokens = lexer.tokenize(code)

        parser = Parser(tokens)
        ast = parser.parse()

        # Should have a directive node
        self.assertIsNotNone(ast, "AST should not be None")
        self.assertEqual(len(ast.children), 1, "AST should have one child")
        self.assertEqual(ast.children[0].node_type, NodeType.DIRECTIVE, "Child should be a DIRECTIVE node")

    def test_parser_instruction(self):
        """Test the parser with an instruction"""
        lexer = Lexer()
        # Add indentation to the instruction
        code = "        LD T0, 0x1234"
        tokens = lexer.tokenize(code)

        parser = Parser(tokens)
        ast = parser.parse()

        # Should have an instruction node
        self.assertIsNotNone(ast, "AST should not be None")
        self.assertEqual(len(ast.children), 1, "AST should have one child")
        self.assertEqual(ast.children[0].node_type, NodeType.INSTRUCTION, "Child should be an INSTRUCTION node")

        # The instruction should have 2 operands
        instruction = ast.children[0]
        self.assertEqual(len(instruction.children), 2, "Instruction should have 2 operands")

    def test_parser_label(self):
        """Test the parser with a label"""
        lexer = Lexer()
        # No indentation for label, but indent instruction
        code = "start:\n        LD T0, 0x1234"
        tokens = lexer.tokenize(code)

        parser = Parser(tokens)
        ast = parser.parse()

        # Should have a program node with label and instruction
        self.assertIsNotNone(ast, "AST should not be None")
        self.assertEqual(len(ast.children), 1, "AST should have one child")

        # The program node should have label and instruction
        program_node = ast.children[0]
        self.assertEqual(len(program_node.children), 2, "Program node should have 2 children")
        self.assertEqual(program_node.children[0].node_type, NodeType.LABEL, "First child should be a LABEL node")
        self.assertEqual(program_node.children[1].node_type, NodeType.INSTRUCTION, "Second child should be an INSTRUCTION node")

    def test_parser_vliw(self):
        """Test the parser with a VLIW instruction"""
        lexer = Lexer()
        # Add indentation to the VLIW instruction
        code = "        [ADD T0, T1, T2] [SUB T3, T4, 0x1]"
        tokens = lexer.tokenize(code)

        parser = Parser(tokens)
        ast = parser.parse()

        # Should have a VLIW instruction node
        self.assertIsNotNone(ast, "AST should not be None")
        self.assertEqual(len(ast.children), 1, "AST should have one child")
        self.assertEqual(ast.children[0].node_type, NodeType.VLIW_INSTRUCTION, "Child should be a VLIW_INSTRUCTION node")

        # The VLIW node should have 2 instruction children
        vliw_node = ast.children[0]
        self.assertEqual(len(vliw_node.children), 2, "VLIW node should have 2 instruction children")

if __name__ == "__main__":
    unittest.main()
