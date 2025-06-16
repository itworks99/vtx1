"""
VTX1 Assembler - Lexical Analyzer (Lexer)

This module tokenizes VTX1 assembly source code according to the formal grammar.
It transforms the raw text into a stream of tokens that can be processed by the parser.
"""

import re
import os
import sys
from enum import Enum, auto
from typing import List, Tuple, Optional, Dict

# Define token types as an enum for easier reference
class TokenType(Enum):
    # Instruction categories
    ALU_OP = auto()
    MEM_OP = auto()
    CTRL_OP = auto()
    VEC_OP = auto()
    FP_OP = auto()
    SYS_OP = auto()
    COMPLEX_OP = auto()
    COMPLEX_VEC = auto()
    COMPLEX_MEM = auto()
    COMPLEX_SYS = auto()

    # Registers
    GPR = auto()
    SPECIAL_REG = auto()
    VECTOR_REG = auto()
    FP_REG = auto()

    # VLIW delimiters
    LBRACKET = auto()
    RBRACKET = auto()

    # Operand syntax
    COMMA = auto()
    COLON = auto()
    PLUS = auto()
    LSQUARE = auto()
    RSQUARE = auto()

    # Directives
    DIRECTIVE = auto()

    # Literals
    TERNARY = auto()
    BINARY = auto()
    HEXADECIMAL = auto()
    DECIMAL = auto()
    STRING = auto()

    # Identifiers
    IDENTIFIER = auto()

    # Comments and whitespace
    COMMENT = auto()
    NEWLINE = auto()

    # Special tokens
    EOF = auto()
    ERROR = auto()

# Define a Token class to represent individual tokens
class Token:
    def __init__(self, type: TokenType, value: str, line: int, column: int):
        self.type = type
        self.value = value
        self.line = line
        self.column = column

    def __repr__(self):
        return f"Token({self.type}, '{self.value}', line={self.line}, col={self.column})"

class Lexer:
    def __init__(self):
        # Define token patterns
        self.token_specs = [
            # Whitespace and comments
            (r'\s+', None),  # Whitespace is ignored
            (r';.*$', TokenType.COMMENT),

            # Instruction categories (case-insensitive)
            (r'ADD|SUB|MUL|AND|OR|NOT|XOR|SHL|SHR|ROL|ROR|CMP|TEST|INC|DEC|NEG', TokenType.ALU_OP),
            (r'LD|ST|VLD|VST|FLD|FST|LEA|PUSH', TokenType.MEM_OP),
            (r'JMP|JAL|JR|JALR|BEQ|BNE|BLT|BGE|BLTU|BGEU|CALL|RET', TokenType.CTRL_OP),
            (r'VADD|VSUB|VMUL|VAND|VOR|VNOT|VSHL|VSHR', TokenType.VEC_OP),
            (r'FADD|FSUB|FMUL|FCMP|FMOV|FNEG', TokenType.FP_OP),
            (r'NOP|WFI', TokenType.SYS_OP),

            # Complex operations
            (r'DIV|MOD|UDIV|UMOD|SQRT|ABS|SIN|COS|TAN|ASIN|ACOS|ATAN|EXP|LOG', TokenType.COMPLEX_OP),
            (r'VDOT|VREDUCE|VMAX|VMIN|VSUM|VPERM', TokenType.COMPLEX_VEC),
            (r'CACHE|FLUSH|MEMBAR', TokenType.COMPLEX_MEM),
            (r'SYSCALL|BREAK|HALT', TokenType.COMPLEX_SYS),

            # Registers
            (r'T[0-6]', TokenType.GPR),
            (r'TA|TB|TC|TS|TI', TokenType.SPECIAL_REG),
            (r'VA|VT|VB', TokenType.VECTOR_REG),
            (r'FA|FT|FB', TokenType.FP_REG),

            # VLIW and operand syntax
            (r'\[', TokenType.LBRACKET),
            (r'\]', TokenType.RBRACKET),
            (r',', TokenType.COMMA),
            (r':', TokenType.COLON),
            (r'\+', TokenType.PLUS),

            # Directives
            (r'\.(ORG|DB|DW|DT|EQU|INCLUDE|SECTION|ALIGN|SPACE)', TokenType.DIRECTIVE),

            # Literals
            (r'0t[-0+]+', TokenType.TERNARY),
            (r'0b[01]+', TokenType.BINARY),
            (r'0x[0-9A-Fa-f]+', TokenType.HEXADECIMAL),
            (r'[0-9]+', TokenType.DECIMAL),
            (r'"[^"]*"', TokenType.STRING),

            # Identifiers must come after keywords
            (r'[A-Za-z_][A-Za-z0-9_]*', TokenType.IDENTIFIER),

            # Newline (treat separately for line counting)
            (r'\n', TokenType.NEWLINE),

            # Error token (must be last)
            (r'.', TokenType.ERROR)
        ]

        # Compile all patterns into a single regex for efficiency
        self.regex = '|'.join(f'(?P<g{i}>{pattern})' for i, (pattern, _) in enumerate(self.token_specs))
        self.regex = re.compile(self.regex, re.IGNORECASE)

        # Maps group name to token type
        self.group_to_type = {f'g{i}': token_type for i, (_, token_type) in enumerate(self.token_specs)}

        # Current state
        self.source_code = ""
        self.tokens = []
        self.current_pos = 0
        self.line_num = 1
        self.column_num = 1

    def tokenize(self, source_code: str) -> List[Token]:
        """Convert source code text into a list of tokens"""
        self.source_code = source_code
        self.tokens = []
        self.current_pos = 0
        self.line_num = 1
        self.column_num = 1

        # Main tokenization loop
        while self.current_pos < len(source_code):
            match = self.regex.match(source_code, self.current_pos)
            if not match:
                # If no match, report error and skip character
                self._add_token(TokenType.ERROR, source_code[self.current_pos])
                self.current_pos += 1
                self.column_num += 1
                continue

            # Identify which pattern matched
            for group_name, group_value in match.groupdict().items():
                if group_value:
                    token_type = self.group_to_type[group_name]
                    if token_type:  # Skip None token types (like whitespace)
                        if token_type == TokenType.NEWLINE:
                            self.line_num += 1
                            self.column_num = 1
                        else:
                            self._add_token(token_type, group_value)
                    break

            # Move to next position after the match
            match_len = match.end() - match.start()
            self.current_pos += match_len
            if token_type != TokenType.NEWLINE:
                self.column_num += match_len

        # Add EOF token at the end
        self._add_token(TokenType.EOF, "")
        return self.tokens

    def _add_token(self, token_type: TokenType, value: str):
        """Add a new token to the token list"""
        token = Token(token_type, value, self.line_num, self.column_num)
        self.tokens.append(token)

    def tokenize_file(self, filename: str) -> List[Token]:
        """Read a file and tokenize its contents"""
        try:
            with open(filename, 'r') as f:
                source_code = f.read()
            return self.tokenize(source_code)
        except FileNotFoundError:
            print(f"Error: File '{filename}' not found.")
            return []
        except Exception as e:
            print(f"Error reading file '{filename}': {e}")
            return []

# Helper function to convert ternary literals to decimal
def ternary_to_decimal(ternary_str: str) -> int:
    """Convert a balanced ternary string to decimal"""
    # Remove '0t' prefix
    if ternary_str.startswith('0t'):
        ternary_str = ternary_str[2:]

    # Convert each trit and sum
    value = 0
    for i, trit in enumerate(reversed(ternary_str)):
        if trit == '+':
            trit_val = 1
        elif trit == '0':
            trit_val = 0
        elif trit == '-':
            trit_val = -1
        else:
            raise ValueError(f"Invalid trit: {trit}")

        value += trit_val * (3 ** i)

    return value

# Helper functions for token stream processing
def is_instruction(token_type: TokenType) -> bool:
    """Check if the token type represents an instruction mnemonic"""
    return token_type in {
        TokenType.ALU_OP, TokenType.MEM_OP, TokenType.CTRL_OP,
        TokenType.VEC_OP, TokenType.FP_OP, TokenType.SYS_OP,
        TokenType.COMPLEX_OP, TokenType.COMPLEX_VEC,
        TokenType.COMPLEX_MEM, TokenType.COMPLEX_SYS
    }

def is_register(token_type: TokenType) -> bool:
    """Check if the token type represents a register"""
    return token_type in {
        TokenType.GPR, TokenType.SPECIAL_REG,
        TokenType.VECTOR_REG, TokenType.FP_REG
    }

def is_literal(token_type: TokenType) -> bool:
    """Check if the token type represents a literal value"""
    return token_type in {
        TokenType.TERNARY, TokenType.BINARY,
        TokenType.HEXADECIMAL, TokenType.DECIMAL
    }

# Usage example
if __name__ == "__main__":
    # Test the lexer with a simple example
    test_code = """
        ; Simple VTX1 assembly program example
        LD T0, 0x1000        ; Load value from address 0x1000 into T0
        LD T1, 0x1004        ; Load value from address 0x1004 into T1
    loop:   
        ADD T2, T0, T1        ; T2 = T0 + T1
        [ADD T0, T1, 0t+] [SUB T1, T2, T0] [NOP]  ; VLIW instruction with 3 operations
        BNE T0, 0, loop       ; Branch to loop if T0 != 0
        ST T2, [TB+8]         ; Store result to address TB+8
        WFI                   ; Wait for interrupt
    """

    lexer = Lexer()
    tokens = lexer.tokenize(test_code)

    # Print all tokens except whitespace
    for token in tokens:
        if token.type != TokenType.NEWLINE:
            print(token)
