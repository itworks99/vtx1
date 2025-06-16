"""
VTX1 Assembler - Parser

This module converts a stream of tokens into an abstract syntax tree (AST)
representing the structure of a VTX1 assembly program.
"""

import sys
from typing import List, Dict, Optional, Union, Tuple, Any
from enum import Enum

# Import token definitions from the lexer
sys.path.append('../lexer')
from vtx1_lexer import Token, TokenType, is_instruction, is_register, is_literal

# Define AST node types
class NodeType(Enum):
    PROGRAM = "Program"
    INSTRUCTION = "Instruction"
    VLIW_INSTRUCTION = "VliwInstruction"
    LABEL = "Label"
    OPERAND = "Operand"
    REGISTER = "Register"
    MEMORY_REF = "MemoryReference"
    IMMEDIATE = "Immediate"
    SYMBOL_REF = "SymbolReference"
    DIRECTIVE = "Directive"
    COMMENT = "Comment"

class ASTNode:
    """Base class for all AST nodes"""
    def __init__(self, node_type: NodeType, value: Any = None, line: int = 0, column: int = 0):
        self.node_type = node_type
        self.value = value
        self.line = line
        self.column = column
        self.children = []

    def add_child(self, child):
        self.children.append(child)
        return self

    def __repr__(self):
        return f"{self.node_type.value}({self.value}, line={self.line})"

class ParseError(Exception):
    """Exception raised for parsing errors"""
    def __init__(self, message: str, token: Token):
        self.message = message
        self.token = token
        super().__init__(f"{message} at line {token.line}, column {token.column}, token: {token.value}")

class Parser:
    def __init__(self, tokens: List[Token]):
        self.tokens = tokens
        self.current = 0
        self.symbol_table = {}  # Maps labels to their positions
        self.errors = []
        self.warnings = []

    def parse(self) -> ASTNode:
        """Parse the token stream and construct an AST"""
        program_node = ASTNode(NodeType.PROGRAM, "Program", 0, 0)

        # Process until we reach the end of file
        while not self.is_at_end():
            try:
                # Skip newlines
                while self.match(TokenType.NEWLINE):
                    pass

                if self.is_at_end():
                    break

                # Parse a line of assembly code
                line_node = self.parse_line()
                if line_node:
                    program_node.add_child(line_node)
            except ParseError as e:
                self.errors.append(e)
                self.synchronize()  # Skip to next valid point

        return program_node

    def parse_line(self) -> Optional[ASTNode]:
        """Parse a single line of assembly code"""
        # Check for label
        label_node = None
        if self.check(TokenType.IDENTIFIER) and self.check_next(TokenType.COLON):
            identifier = self.advance()  # Consume identifier
            self.advance()  # Consume colon

            label_node = ASTNode(NodeType.LABEL, identifier.value, identifier.line, identifier.column)

            # Register the label in the symbol table
            self.symbol_table[identifier.value] = {
                "line": identifier.line,
                "column": identifier.column
            }

            # Skip any newlines after the label
            while self.match(TokenType.NEWLINE):
                pass

        # Check for instruction, directive, or comment
        if self.is_at_end() or self.check(TokenType.NEWLINE):
            # Line has only a label or is empty
            return label_node

        if self.check(TokenType.LBRACKET):
            # This is a VLIW instruction
            vliw_node = self.parse_vliw_instruction()
            if label_node:
                vliw_node.line = label_node.line
                vliw_node.column = label_node.column
                return ASTNode(NodeType.PROGRAM, None, label_node.line, label_node.column).add_child(label_node).add_child(vliw_node)
            return vliw_node
        elif is_instruction(self.peek().type):
            # This is a regular instruction
            instruction_node = self.parse_instruction()
            if label_node:
                return ASTNode(NodeType.PROGRAM, None, label_node.line, label_node.column).add_child(label_node).add_child(instruction_node)
            return instruction_node
        elif self.check(TokenType.DIRECTIVE):
            # This is an assembler directive
            directive_node = self.parse_directive()
            if label_node:
                return ASTNode(NodeType.PROGRAM, None, label_node.line, label_node.column).add_child(label_node).add_child(directive_node)
            return directive_node
        elif self.check(TokenType.COMMENT):
            # This is just a comment
            comment = self.advance()
            comment_node = ASTNode(NodeType.COMMENT, comment.value, comment.line, comment.column)
            if label_node:
                return ASTNode(NodeType.PROGRAM, None, label_node.line, label_node.column).add_child(label_node).add_child(comment_node)
            return comment_node
        else:
            # Unexpected token
            token = self.peek()
            raise ParseError(f"Unexpected token: {token.value}", token)

    def parse_vliw_instruction(self) -> ASTNode:
        """Parse a VLIW instruction (group of instructions in brackets)"""
        # Create a new VLIW instruction node
        first_token = self.peek()
        vliw_node = ASTNode(NodeType.VLIW_INSTRUCTION, "VLIW", first_token.line, first_token.column)

        # Up to 3 operations in a VLIW instruction
        for _ in range(3):
            if not self.check(TokenType.LBRACKET):
                break

            # Parse operation in brackets
            self.advance()  # Consume [
            if not is_instruction(self.peek().type):
                raise ParseError("Expected instruction mnemonic after '['", self.peek())

            # Parse the instruction within brackets
            instruction_node = self.parse_instruction()
            vliw_node.add_child(instruction_node)

            if not self.check(TokenType.RBRACKET):
                raise ParseError("Expected ']' after instruction in VLIW", self.peek())

            self.advance()  # Consume ]

        # There should be 1-3 operations in a VLIW instruction
        if len(vliw_node.children) < 1 or len(vliw_node.children) > 3:
            raise ParseError(
                f"VLIW instruction must have 1-3 operations, found {len(vliw_node.children)}",
                first_token
            )

        return vliw_node

    def parse_instruction(self) -> ASTNode:
        """Parse a single instruction"""
        mnemonic = self.advance()  # Consume the instruction mnemonic
        instruction_node = ASTNode(
            NodeType.INSTRUCTION,
            mnemonic.value,
            mnemonic.line,
            mnemonic.column
        )

        # Parse operands (if any)
        if not self.check(TokenType.NEWLINE) and not self.check(TokenType.COMMENT) and not self.check(TokenType.RBRACKET) and not self.is_at_end():
            # First operand
            operand_node = self.parse_operand()
            instruction_node.add_child(operand_node)

            # Additional operands separated by commas
            while self.match(TokenType.COMMA):
                operand_node = self.parse_operand()
                instruction_node.add_child(operand_node)

        return instruction_node

    def parse_operand(self) -> ASTNode:
        """Parse an instruction operand"""
        token = self.peek()

        # Register operand
        if is_register(token.type):
            register = self.advance()
            return ASTNode(NodeType.REGISTER, register.value, register.line, register.column)

        # Immediate value
        elif is_literal(token.type):
            immediate = self.advance()
            return ASTNode(NodeType.IMMEDIATE, immediate.value, immediate.line, immediate.column)

        # Memory reference [reg+offset]
        elif self.match(TokenType.LSQUARE):
            mem_token = token
            base_reg = self.advance()  # Should be a register

            if not is_register(base_reg.type):
                raise ParseError("Expected register as base in memory reference", base_reg)

            mem_node = ASTNode(NodeType.MEMORY_REF, base_reg.value, mem_token.line, mem_token.column)

            # Check for optional offset
            if self.match(TokenType.PLUS):
                if is_register(self.peek().type):
                    # Index register
                    index_reg = self.advance()
                    mem_node.add_child(ASTNode(NodeType.REGISTER, index_reg.value, index_reg.line, index_reg.column))
                elif is_literal(self.peek().type):
                    # Immediate offset
                    offset = self.advance()
                    mem_node.add_child(ASTNode(NodeType.IMMEDIATE, offset.value, offset.line, offset.column))
                else:
                    raise ParseError("Expected register or immediate after '+' in memory reference", self.peek())

            # Expect closing bracket
            if not self.match(TokenType.RSQUARE):
                raise ParseError("Expected ']' after memory reference", self.peek())

            return mem_node

        # Label reference
        elif self.check(TokenType.IDENTIFIER):
            identifier = self.advance()
            return ASTNode(NodeType.SYMBOL_REF, identifier.value, identifier.line, identifier.column)

        else:
            raise ParseError("Expected operand (register, immediate, memory reference, or label)", token)

    def parse_directive(self) -> ASTNode:
        """Parse an assembler directive"""
        directive = self.advance()  # Consume directive token
        directive_node = ASTNode(NodeType.DIRECTIVE, directive.value, directive.line, directive.column)

        # Different directives have different operand requirements
        directive_name = directive.value.upper()

        if directive_name in ['.DB', '.DW', '.DT']:
            # Data directives take a list of values
            while True:
                if is_literal(self.peek().type) or self.check(TokenType.STRING):
                    value = self.advance()
                    directive_node.add_child(ASTNode(NodeType.IMMEDIATE, value.value, value.line, value.column))
                else:
                    raise ParseError(f"Expected value after {directive_name}", self.peek())

                if not self.match(TokenType.COMMA):
                    break

        elif directive_name in ['.ORG', '.ALIGN', '.SPACE']:
            # These directives take a single immediate value
            if is_literal(self.peek().type):
                value = self.advance()
                directive_node.add_child(ASTNode(NodeType.IMMEDIATE, value.value, value.line, value.column))
            else:
                raise ParseError(f"Expected immediate value after {directive_name}", self.peek())

        elif directive_name == '.EQU':
            # .EQU takes an identifier and a value
            if self.check(TokenType.IDENTIFIER):
                identifier = self.advance()
                directive_node.add_child(ASTNode(NodeType.SYMBOL_REF, identifier.value, identifier.line, identifier.column))

                if self.match(TokenType.COMMA):
                    if is_literal(self.peek().type):
                        value = self.advance()
                        directive_node.add_child(ASTNode(NodeType.IMMEDIATE, value.value, value.line, value.column))
                    else:
                        raise ParseError(f"Expected value after comma in {directive_name}", self.peek())
                else:
                    raise ParseError(f"Expected comma after identifier in {directive_name}", self.peek())
            else:
                raise ParseError(f"Expected identifier after {directive_name}", self.peek())

        elif directive_name == '.INCLUDE':
            # .INCLUDE takes a string
            if self.check(TokenType.STRING):
                string = self.advance()
                directive_node.add_child(ASTNode(NodeType.IMMEDIATE, string.value, string.line, string.column))
            else:
                raise ParseError(f"Expected string after {directive_name}", self.peek())

        elif directive_name == '.SECTION':
            # .SECTION takes an identifier
            if self.check(TokenType.IDENTIFIER):
                identifier = self.advance()
                directive_node.add_child(ASTNode(NodeType.SYMBOL_REF, identifier.value, identifier.line, identifier.column))
            else:
                raise ParseError(f"Expected identifier after {directive_name}", self.peek())

        return directive_node

    # Helper methods for the parser

    def is_at_end(self) -> bool:
        """Check if we've reached the end of the token stream"""
        return self.peek().type == TokenType.EOF

    def peek(self) -> Token:
        """Return the current token without advancing"""
        return self.tokens[self.current]

    def previous(self) -> Token:
        """Return the most recently consumed token"""
        return self.tokens[self.current - 1]

    def advance(self) -> Token:
        """Consume and return the current token"""
        if not self.is_at_end():
            self.current += 1
        return self.previous()

    def check(self, type: TokenType) -> bool:
        """Check if the current token is of the given type"""
        if self.is_at_end():
            return False
        return self.peek().type == type

    def check_next(self, type: TokenType) -> bool:
        """Check if the next token is of the given type"""
        if self.current + 1 >= len(self.tokens):
            return False
        return self.tokens[self.current + 1].type == type

    def match(self, type: TokenType) -> bool:
        """Check if the current token matches the given type and advance if it does"""
        if self.check(type):
            self.advance()
            return True
        return False

    def synchronize(self):
        """Recover from a parse error by advancing to the next statement"""
        self.advance()  # Skip the problematic token

        while not self.is_at_end():
            # Skip until we reach a newline or another synchronization point
            if self.previous().type == TokenType.NEWLINE:
                return

            if is_instruction(self.peek().type) or self.check(TokenType.DIRECTIVE):
                return

            self.advance()

# Function to print the AST for debugging
def print_ast(node: ASTNode, indent: int = 0):
    """Print an AST in a readable format for debugging"""
    print("  " * indent + str(node))
    for child in node.children:
        print_ast(child, indent + 1)

# Usage example
if __name__ == "__main__":
    import vtx1_lexer

    # Test code
    test_code = """
        ; Simple VTX1 assembly program example
        .ORG 0x1000
        LD T0, 0x1000        ; Load value from address 0x1000 into T0
        LD T1, 0x1004        ; Load value from address 0x1004 into T1
    loop:   
        ADD T2, T0, T1        ; T2 = T0 + T1
        [ADD T0, T1, 0t+] [SUB T1, T2, T0] [NOP]  ; VLIW instruction with 3 operations
        BNE T0, 0, loop       ; Branch to loop if T0 != 0
        ST T2, [TB+8]         ; Store result to address TB+8
        WFI                   ; Wait for interrupt
    """

    # Tokenize the test code
    lexer = vtx1_lexer.Lexer()
    tokens = lexer.tokenize(test_code)

    # Parse the tokens
    parser = Parser(tokens)
    ast = parser.parse()

    # Print any parsing errors
    if parser.errors:
        print(f"Found {len(parser.errors)} parsing error(s):")
        for error in parser.errors:
            print(f"  {error}")

    # Print the AST
    print("\nAbstract Syntax Tree:")
    print_ast(ast)

    # Print the symbol table
    print("\nSymbol Table:")
    for symbol, info in parser.symbol_table.items():
        print(f"  {symbol}: line {info['line']}, column {info['column']}")
