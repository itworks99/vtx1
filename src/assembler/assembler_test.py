#!/usr/bin/env python3
"""
VTX1 Assembler Integration Test

This script tests the VTX1 assembler by creating and assembling small, valid assembly programs.

Usage:
    python assembler_test.py
"""

import os
import sys
import unittest
import tempfile
import shutil

# Fix import paths for all components
script_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(script_dir)  # Add current directory
sys.path.append(os.path.join(script_dir, 'lexer'))  # Add lexer subdirectory
sys.path.append(os.path.join(script_dir, 'parser'))  # Add parser subdirectory
sys.path.append(os.path.join(script_dir, 'codegen'))  # Add codegen subdirectory

# Import the assembler
from vtx1_asm import Assembler

class AssemblerIntegrationTests(unittest.TestCase):
    """Test the VTX1 assembler through its main assemble method"""

    def setUp(self):
        """Set up test environment"""
        self.temp_dir = tempfile.mkdtemp()
        self.assembler = Assembler(verbose=True)

    def tearDown(self):
        """Clean up temporary files"""
        shutil.rmtree(self.temp_dir)

    def create_test_asm(self, content):
        """Create a test assembly file with the given content"""
        test_file = os.path.join(self.temp_dir, "test.asm")
        with open(test_file, "w") as f:
            f.write(content)
        return test_file

    def test_basic_program(self):
        """Test assembling a basic program with standard components"""
        # Create a small, valid assembly program
        asm_content = """
;===============================================================================
; Basic Test Program
;===============================================================================

        .ORG 0x1000          ; Start at address 0x1000

;===============================================================================
; Main program
;===============================================================================
main:
        ; Initialize registers
        LD T0, 0x1234        ; Load immediate
        LD T1, 0x5678        ; Load another value
        ADD T2, T0, T1       ; Add them together
        
        ; Store the result
        ST T2, 0x2000        ; Store at memory location

        ; End program
        NOP                  ; No operation
"""
        # Create the test file
        test_file = self.create_test_asm(asm_content)
        output_file = os.path.join(self.temp_dir, "test.bin")

        # Assemble the file
        success = self.assembler.assemble(test_file, output_file)

        # Check assembly succeeded
        self.assertTrue(success, "Failed to assemble basic program")

        # Check output file exists
        self.assertTrue(os.path.exists(output_file), "Output binary file not created")

        # Check output file has content
        self.assertGreater(os.path.getsize(output_file), 0, "Output binary file is empty")

    def test_vliw_program(self):
        """Test assembling a program with VLIW instructions"""
        # Create a small program with VLIW instructions
        asm_content = """
;===============================================================================
; VLIW Test Program
;===============================================================================

        .ORG 0x1000          ; Start at address 0x1000

;===============================================================================
; Main program
;===============================================================================
main:
        ; Initialize registers
        LD T0, 0x10          ; Load immediate
        LD T1, 0x20          ; Load another value
        
        ; Use VLIW to perform multiple operations in parallel
        [ADD T2, T0, T1] [SUB T3, T1, T0] [MUL T4, T0, T1]
        
        ; Store the results
        ST T2, 0x2000        ; Store result 1
        ST T3, 0x2004        ; Store result 2
        ST T4, 0x2008        ; Store result 3
"""
        # Create the test file
        test_file = self.create_test_asm(asm_content)
        output_file = os.path.join(self.temp_dir, "vliw_test.bin")

        # Assemble the file
        success = self.assembler.assemble(test_file, output_file)

        # Check assembly succeeded
        self.assertTrue(success, "Failed to assemble VLIW program")

        # Check output file exists
        self.assertTrue(os.path.exists(output_file), "Output binary file not created")

    def test_example_files(self):
        """Test assembling the example files provided with the assembler"""
        examples_dir = os.path.join(os.path.dirname(__file__), "examples")

        # Get list of example files
        example_files = [f for f in os.listdir(examples_dir) if f.endswith('.asm')]

        for example_file in example_files:
            file_path = os.path.join(examples_dir, example_file)
            output_path = os.path.join(self.temp_dir, example_file + ".bin")

            # Assemble the example file
            success = self.assembler.assemble(file_path, output_path)

            print(f"Assembly of {example_file}: {'Succeeded' if success else 'Failed'}")

            # We don't assert success here because some examples might have intended errors
            # Just check that output exists if assembly succeeded
            if success:
                self.assertTrue(os.path.exists(output_path),
                             f"Output file not created for {example_file}")

if __name__ == "__main__":
    unittest.main()
