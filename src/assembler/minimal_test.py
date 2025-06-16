#!/usr/bin/env python3
"""
VTX1 Assembler Minimal Test

This script performs a minimal test of the VTX1 assembler to verify it can be
initialized and run without errors.
"""

import os
import sys
import unittest
import tempfile

# Fix import paths for all components
script_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(script_dir)  # Add current directory
sys.path.append(os.path.join(script_dir, 'lexer'))  # Add lexer subdirectory
sys.path.append(os.path.join(script_dir, 'parser'))  # Add parser subdirectory
sys.path.append(os.path.join(script_dir, 'codegen'))  # Add codegen subdirectory

# Import the assembler
from vtx1_asm import Assembler
from lexer.vtx1_lexer import Lexer, TokenType, Token

class MinimalAssemblerTest(unittest.TestCase):
    """Minimal tests for the VTX1 assembler"""

    def test_assembler_initialization(self):
        """Test that we can create an assembler instance"""
        try:
            assembler = Assembler(verbose=True)
            self.assertIsNotNone(assembler, "Assembler instance should not be None")
            print("✓ Successfully created Assembler instance")
        except Exception as e:
            self.fail(f"Failed to create Assembler instance: {str(e)}")

    def test_lexer_tokenization(self):
        """Test that the lexer can tokenize a simple string"""
        try:
            lexer = Lexer()
            tokens = lexer.tokenize(".ORG 0x1000")

            self.assertGreater(len(tokens), 0, "Should generate tokens")

            # Print out the tokens to see what we have
            print("\nTokens generated:")
            for i, token in enumerate(tokens):
                print(f"  {i}: {token.type.name} = '{token.value}'")

            print(f"\n✓ Successfully tokenized input, generated {len(tokens)} tokens")
        except Exception as e:
            self.fail(f"Failed to tokenize: {str(e)}")

if __name__ == "__main__":
    unittest.main()
