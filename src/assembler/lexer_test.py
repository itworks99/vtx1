#!/usr/bin/env python3
"""
VTX1 Lexer Test

This script tests the lexer component of the VTX1 assembler by:
1. Processing example assembly files
2. Analyzing token frequency and structure
3. Validating lexical patterns in the code

This helps verify the assembly code structure even if the parser has issues.
"""

import os
import sys
import unittest
import tempfile
from collections import Counter

# Fix import paths
script_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(script_dir)  # Add current directory
sys.path.append(os.path.join(script_dir, 'lexer'))  # Add lexer subdirectory
sys.path.append(os.path.join(script_dir, 'parser'))  # Add parser subdirectory
sys.path.append(os.path.join(script_dir, 'codegen'))  # Add codegen subdirectory

# Import the lexer
from lexer.vtx1_lexer import Lexer, TokenType, Token

class LexerTest(unittest.TestCase):
    """Tests for the VTX1 assembler lexer"""

    def setUp(self):
        """Set up the test environment"""
        self.lexer = Lexer()
        self.examples_dir = os.path.join(script_dir, "examples")

    def test_example_files(self):
        """Test that the lexer can tokenize all example files"""
        # Get list of example files
        example_files = [f for f in os.listdir(self.examples_dir) if f.endswith('.asm')]

        for example_file in example_files:
            file_path = os.path.join(self.examples_dir, example_file)

            try:
                print(f"\nProcessing {example_file}...")

                # Read the file content
                with open(file_path, 'r') as f:
                    content = f.read()

                # Tokenize the file
                tokens = self.lexer.tokenize(content)

                # Count token types for analysis
                type_counts = Counter(token.type for token in tokens)

                # Calculate some statistics
                num_instructions = sum(1 for token in tokens if self._is_instruction_token(token))
                num_registers = sum(1 for token in tokens if self._is_register_token(token))
                num_directives = sum(1 for token in tokens if token.type == TokenType.DIRECTIVE)
                num_labels = self._count_labels(tokens)

                # Print statistics
                print(f"  Total tokens: {len(tokens)}")
                print(f"  Instructions: {num_instructions}")
                print(f"  Registers: {num_registers}")
                print(f"  Directives: {num_directives}")
                print(f"  Labels: {num_labels}")
                print(f"  Most common token types:")

                # Print top 5 most common token types
                for token_type, count in type_counts.most_common(5):
                    print(f"    - {token_type.name}: {count}")

                # Simple validation - check for basic structure
                self.assertGreater(len(tokens), 0, f"No tokens in {example_file}")

                # Advanced validation - check for expected tokens based on file name
                self._validate_file_specific_tokens(example_file, tokens)

                print(f"✓ Successfully tokenized {example_file}")

            except Exception as e:
                self.fail(f"Failed to process {example_file}: {str(e)}")

    def test_token_sequence(self):
        """Test the token sequence patterns in the example files"""
        # Get list of example files
        example_files = [f for f in os.listdir(self.examples_dir) if f.endswith('.asm')]

        for example_file in example_files:
            file_path = os.path.join(self.examples_dir, example_file)

            try:
                # Read the file content
                with open(file_path, 'r') as f:
                    content = f.read()

                # Tokenize the file
                tokens = self.lexer.tokenize(content)

                # Analyze token sequences
                print(f"\nAnalyzing token sequences in {example_file}...")

                # Find instruction patterns
                instruction_patterns = self._find_instruction_patterns(tokens)

                # Print common instruction patterns (up to 3)
                for pattern_desc, count in instruction_patterns.most_common(3):
                    print(f"  - {pattern_desc}: {count} occurrences")

                print(f"✓ Successfully analyzed token sequences in {example_file}")

            except Exception as e:
                self.fail(f"Failed to analyze token sequences in {example_file}: {str(e)}")

    def test_specific_constructs(self):
        """Test specific language constructs in the assembly files"""
        constructs = {
            "vliw_instructions": 0,
            "memory_references": 0,
            "labels": 0,
            "directives": 0,
        }

        # Get list of example files
        example_files = [f for f in os.listdir(self.examples_dir) if f.endswith('.asm')]

        for example_file in example_files:
            file_path = os.path.join(self.examples_dir, example_file)

            try:
                # Read the file content
                with open(file_path, 'r') as f:
                    content = f.read()

                # Tokenize the file
                tokens = self.lexer.tokenize(content)

                # Count specific constructs
                vliw_count = self._count_vliw_instructions(tokens)
                memory_ref_count = self._count_memory_references(tokens)
                label_count = self._count_labels(tokens)
                directive_count = sum(1 for token in tokens if token.type == TokenType.DIRECTIVE)

                constructs["vliw_instructions"] += vliw_count
                constructs["memory_references"] += memory_ref_count
                constructs["labels"] += label_count
                constructs["directives"] += directive_count

            except Exception as e:
                self.fail(f"Failed to analyze constructs in {example_file}: {str(e)}")

        # Print summary of constructs across all files
        print("\nLanguage constructs found across all example files:")
        for construct, count in constructs.items():
            print(f"  - {construct}: {count}")

        # Validate that we found at least some constructs
        self.assertGreater(sum(constructs.values()), 0, "No language constructs found")

    def _is_instruction_token(self, token):
        """Check if a token is an instruction mnemonic"""
        instruction_types = [
            TokenType.ALU_OP, TokenType.MEM_OP, TokenType.CTRL_OP,
            TokenType.VEC_OP, TokenType.FP_OP, TokenType.SYS_OP,
            TokenType.COMPLEX_OP, TokenType.COMPLEX_VEC,
            TokenType.COMPLEX_MEM, TokenType.COMPLEX_SYS
        ]
        return token.type in instruction_types

    def _is_register_token(self, token):
        """Check if a token is a register"""
        register_types = [
            TokenType.GPR, TokenType.SPECIAL_REG,
            TokenType.VECTOR_REG, TokenType.FP_REG
        ]
        return token.type in register_types

    def _count_labels(self, tokens):
        """Count label declarations in token stream"""
        count = 0

        # Look for IDENTIFIER followed by COLON pattern
        for i in range(len(tokens) - 1):
            if tokens[i].type == TokenType.IDENTIFIER and tokens[i+1].type == TokenType.COLON:
                count += 1

        return count

    def _count_vliw_instructions(self, tokens):
        """Count VLIW instruction blocks in token stream"""
        count = 0
        bracket_depth = 0

        # Count LBRACKET tokens at depth 0 (start of VLIW block)
        for token in tokens:
            if token.type == TokenType.LBRACKET:
                if bracket_depth == 0:
                    count += 1
                bracket_depth += 1
            elif token.type == TokenType.RBRACKET:
                bracket_depth -= 1

        return count

    def _count_memory_references(self, tokens):
        """Count memory reference expressions in token stream"""
        count = 0

        # Look for LSQUARE followed by register pattern
        for i in range(len(tokens) - 1):
            if tokens[i].type == TokenType.LSQUARE and self._is_register_token(tokens[i+1]):
                count += 1

        return count

    def _find_instruction_patterns(self, tokens):
        """Find common instruction patterns in token stream"""
        patterns = Counter()

        for i in range(len(tokens) - 2):
            if self._is_instruction_token(tokens[i]):
                # Look for patterns like "INSTRUCTION REG" or "INSTRUCTION REG, REG"
                if i + 1 < len(tokens) and self._is_register_token(tokens[i+1]):
                    if i + 3 < len(tokens) and tokens[i+2].type == TokenType.COMMA and self._is_register_token(tokens[i+3]):
                        patterns[f"{tokens[i].value} REG, REG"] += 1
                    else:
                        patterns[f"{tokens[i].value} REG"] += 1

        return patterns

    def _validate_file_specific_tokens(self, filename, tokens):
        """Validate tokens specific to certain file types"""
        if "hello_world" in filename.lower():
            # Hello world should have string literals
            string_literals = [t for t in tokens if t.type == TokenType.STRING]
            self.assertGreater(len(string_literals), 0, "No string literals found in hello world example")

        elif "vliw" in filename.lower():
            # VLIW examples should have bracket tokens
            brackets = [t for t in tokens if t.type in (TokenType.LBRACKET, TokenType.RBRACKET)]
            self.assertGreater(len(brackets), 0, "No brackets found in VLIW example")

        elif "ternary" in filename.lower():
            # Ternary examples should have ternary literals
            ternary_literals = [t for t in tokens if t.type == TokenType.TERNARY]
            self.assertGreater(len(ternary_literals), 0, "No ternary literals found in ternary example")

if __name__ == "__main__":
    unittest.main()
