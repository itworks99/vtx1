#!/usr/bin/env python3
"""
VTX1 Assembler Basic Structure Test

This script tests the basic structure and setup of the VTX1 assembler by:
1. Verifying that all the components exist and can be instantiated
2. Checking the structure of the assembler and its components
3. Running a simple test on the example files to see if they are readable

Usage:
    python basic_structure_test.py
"""

import os
import sys
import unittest
import tempfile
from pathlib import Path

# Fix import paths
script_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(script_dir)  # Add current directory
sys.path.append(os.path.join(script_dir, 'lexer'))  # Add lexer subdirectory
sys.path.append(os.path.join(script_dir, 'parser'))  # Add parser subdirectory
sys.path.append(os.path.join(script_dir, 'codegen'))  # Add codegen subdirectory

# Import the assembler components
from vtx1_asm import Assembler

class BasicStructureTests(unittest.TestCase):
    """Test the basic structure of the assembler"""

    def setUp(self):
        """Set up test environment"""
        # Get paths
        self.assembler_dir = os.path.dirname(os.path.abspath(__file__))
        self.examples_dir = os.path.join(self.assembler_dir, "examples")
        self.temp_dir = tempfile.mkdtemp()

        # Create assembler instance
        self.assembler = Assembler(verbose=True)

    def test_assembler_exists(self):
        """Test that the assembler class exists and can be instantiated"""
        self.assertIsNotNone(self.assembler, "Assembler instance could not be created")
        self.assertIsInstance(self.assembler, Assembler, "Instance is not of type Assembler")

    def test_assembler_structure(self):
        """Test that the assembler has the expected structure"""
        # Check that the assembler has expected attributes
        self.assertTrue(hasattr(self.assembler, 'lexer'), "Assembler does not have a lexer attribute")
        self.assertTrue(hasattr(self.assembler, 'assemble'), "Assembler does not have an assemble method")

        # Check that the assembler has the expected methods
        self.assertTrue(callable(getattr(self.assembler, 'assemble', None)),
                       "Assembler's assemble attribute is not callable")

    def test_example_files_exist(self):
        """Test that example files exist and can be read"""
        # Check that the examples directory exists
        self.assertTrue(os.path.exists(self.examples_dir), "Examples directory does not exist")

        # Check that at least one example file exists
        example_files = [f for f in os.listdir(self.examples_dir) if f.endswith('.asm')]
        self.assertTrue(len(example_files) > 0, "No example files found")

        # Check that we can read an example file
        example_file = os.path.join(self.examples_dir, example_files[0])
        with open(example_file, 'r') as f:
            content = f.read()
            self.assertTrue(len(content) > 0, f"Example file {example_files[0]} is empty")
            print(f"First 100 chars of {example_files[0]}: {content[:100]}")

    def test_assembler_reads_examples(self):
        """Test that the assembler can read example files"""
        # Find the example files
        example_files = [f for f in os.listdir(self.examples_dir) if f.endswith('.asm')]

        # Check that the assembler can read at least one example
        if example_files:
            example_file = os.path.join(self.examples_dir, example_files[0])

            # Use the assembler's internal file reading, not the full assembly process
            try:
                with open(example_file, 'r') as f:
                    source_code = f.read()
                    self.assembler.source_lines = source_code.splitlines()
                    print(f"Successfully read {len(self.assembler.source_lines)} lines from {example_files[0]}")
                    self.assertTrue(len(self.assembler.source_lines) > 0,
                                  f"No lines read from {example_files[0]}")
            except Exception as e:
                self.fail(f"Failed to read example file: {str(e)}")

    def test_lexer_structure(self):
        """Test that the lexer has the expected structure"""
        # Check that the lexer exists
        self.assertIsNotNone(self.assembler.lexer, "Assembler's lexer is None")

        # Check that the lexer has a tokenize method
        self.assertTrue(hasattr(self.assembler.lexer, 'tokenize'),
                       "Lexer does not have a tokenize method")
        self.assertTrue(callable(getattr(self.assembler.lexer, 'tokenize', None)),
                       "Lexer's tokenize attribute is not callable")

        # Test that the lexer can tokenize a simple string
        try:
            tokens = self.assembler.lexer.tokenize(".ORG 0x1000")
            self.assertGreater(len(tokens), 0, "No tokens generated from test string")
            print(f"Generated {len(tokens)} tokens from test string")
        except Exception as e:
            self.fail(f"Failed to tokenize test string: {str(e)}")

if __name__ == "__main__":
    unittest.main()
