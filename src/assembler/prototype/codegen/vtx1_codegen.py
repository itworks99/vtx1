"""
VTX1 Assembler - Code Generator

This module takes an abstract syntax tree (AST) from the parser and
generates binary machine code according to the VTX1 instruction format.
"""

import sys
import os
import struct
from enum import Enum
from typing import List, Dict, Optional, Union, Tuple, Any, Set

# Import from other modules
sys.path.append('../parser')
sys.path.append('../lexer')
from vtx1_parser import ASTNode, NodeType
from vtx1_lexer import TokenType, ternary_to_decimal

# VLIW instruction encoding parameters
INSTRUCTION_SIZE = 32  # bits per instruction
VLIW_SIZE = 96         # bits for entire VLIW word (3 instructions)

# VTX1 instruction field widths
OPCODE_WIDTH = 6       # Opcode field width
REG_WIDTH = 3          # Register field width (3 bits per register)
IMM_WIDTH = 11         # Immediate/offset field width
TYPE_WIDTH = 3         # Operation type field width
PAR_WIDTH = 3          # Parallel execution flags width

# VTX1 encoding for operation types
class OpType(Enum):
    ALU = 0            # ALU Operation (000)
    MEMORY = 1         # Memory Operation (001)
    CONTROL = 2        # Control Operation (010)
    VECTOR = 3         # Vector Operation (011)
    FPU = 4            # FPU Operation (100)
    SYSTEM = 5         # System Operation (101)
    MICROCODE = 6      # Microcode Operation (110)
    RESERVED = 7       # Reserved (111)

# VTX1 encoding for parallel flags
class ParFlags(Enum):
    SERIAL = 0         # Serial Execution (000)
    PARALLEL_ALU = 1   # Parallel with ALU (001)
    PARALLEL_MEM = 2   # Parallel with Memory (010)
    PARALLEL_CTRL = 3  # Parallel with Control (011)
    FULL_PARALLEL = 4  # Full Parallel (100)
    RESERVED_1 = 5     # Reserved (101)
    RESERVED_2 = 6     # Reserved (110)
    RESERVED_3 = 7     # Reserved (111)

# Register encodings
GPR_ENCODING = {
    'T0': 0b000,
    'T1': 0b001,
    'T2': 0b010,
    'T3': 0b011,
    'T4': 0b100,
    'T5': 0b101,
    'T6': 0b110
}

SPECIAL_REG_ENCODING = {
    'TA': 0b111,        # Accumulator (111)
    'TB': 0b000,        # Base Pointer (accessed via type field)
    'TC': 0b001,        # Program Counter (accessed via type field)
    'TS': 0b010,        # Status Register (accessed via type field)
    'TI': 0b011         # Instruction Register (accessed via type field)
}

# Instruction opcodes and their encoding
ALU_OPCODES = {
    'ADD': 0b000001,    # ADD - Addition
    'SUB': 0b000010,    # SUB - Subtraction
    'MUL': 0b000011,    # MUL - Multiplication
    'AND': 0b000100,    # AND - Logical AND
    'OR':  0b000101,    # OR  - Logical OR
    'NOT': 0b000110,    # NOT - Logical NOT
    'XOR': 0b000111,    # XOR - Logical XOR
    'SHL': 0b001000,    # SHL - Shift Left
    'SHR': 0b001001,    # SHR - Shift Right
    'ROL': 0b001010,    # ROL - Rotate Left
    'ROR': 0b001011,    # ROR - Rotate Right
    'CMP': 0b001100,    # CMP - Compare
    'TEST': 0b001101,   # TEST - Test bits
    'INC': 0b001110,    # INC - Increment
    'DEC': 0b001111,    # DEC - Decrement
    'NEG': 0b010000     # NEG - Negate
}

MEM_OPCODES = {
    'LD':  0b010001,    # LD  - Load
    'ST':  0b010010,    # ST  - Store
    'VLD': 0b010011,    # VLD - Vector Load
    'VST': 0b010100,    # VST - Vector Store
    'FLD': 0b010101,    # FLD - Floating-point Load
    'FST': 0b010110,    # FST - Floating-point Store
    'LEA': 0b010111,    # LEA - Load Effective Address
    'PUSH': 0b011000    # PUSH - Push to stack
}

CTRL_OPCODES = {
    'JMP': 0b011001,    # JMP - Jump
    'JAL': 0b011010,    # JAL - Jump and Link
    'JR':  0b011011,    # JR  - Jump Register
    'JALR': 0b011100,   # JALR - Jump and Link Register
    'BEQ': 0b011101,    # BEQ - Branch if Equal
    'BNE': 0b011110,    # BNE - Branch if Not Equal
    'BLT': 0b011111,    # BLT - Branch if Less Than
    'BGE': 0b100000,    # BGE - Branch if Greater or Equal
    'BLTU': 0b100001,   # BLTU - Branch if Less Than Unsigned
    'BGEU': 0b100010,   # BGEU - Branch if Greater or Equal Unsigned
    'CALL': 0b100011,   # CALL - Call subroutine
    'RET': 0b100100     # RET - Return from subroutine
}

VEC_OPCODES = {
    'VADD': 0b100101,   # VADD - Vector Add
    'VSUB': 0b100110,   # VSUB - Vector Subtract
    'VMUL': 0b100111,   # VMUL - Vector Multiply
    'VAND': 0b101000,   # VAND - Vector AND
    'VOR':  0b101001,   # VOR  - Vector OR
    'VNOT': 0b101010,   # VNOT - Vector NOT
    'VSHL': 0b101011,   # VSHL - Vector Shift Left
    'VSHR': 0b101100    # VSHR - Vector Shift Right
}

FPU_OPCODES = {
    'FADD': 0b101101,   # FADD - Floating-point Add
    'FSUB': 0b101110,   # FSUB - Floating-point Subtract
    'FMUL': 0b101111,   # FMUL - Floating-point Multiply
    'FCMP': 0b110000,   # FCMP - Floating-point Compare
    'FMOV': 0b110001,   # FMOV - Floating-point Move
    'FNEG': 0b110010    # FNEG - Floating-point Negate
}

SYS_OPCODES = {
    'NOP': 0b110011,    # NOP - No Operation
    'WFI': 0b110100     # WFI - Wait for Interrupt
}

MICROCODE_OPCODES = {
    # Complex operations (microcode)
    'DIV':  0b110101,   # DIV - Division
    'MOD':  0b110110,   # MOD - Modulo
    'UDIV': 0b110111,   # UDIV - Unsigned Division
    'UMOD': 0b111000,   # UMOD - Unsigned Modulo
    'SQRT': 0b111001,   # SQRT - Square Root
    'ABS':  0b111010,   # ABS - Absolute Value

    # Transcendental functions
    'SIN':  0b111011,   # SIN - Sine
    'COS':  0b111100,   # COS - Cosine
    'TAN':  0b111101,   # TAN - Tangent
    'ASIN': 0b111110,   # ASIN - Arc Sine
    'ACOS': 0b111111,   # ACOS - Arc Cosine
    'ATAN': 0b000000,   # ATAN - Arc Tangent (using 0b000000)
    'EXP':  0b000001,   # EXP - Exponential (using 0b000001)
    'LOG':  0b000010,   # LOG - Logarithm (using 0b000010)

    # Advanced vector operations
    'VDOT':    0b000011,  # VDOT - Vector Dot Product
    'VREDUCE': 0b000100,  # VREDUCE - Vector Reduction
    'VMAX':    0b000101,  # VMAX - Vector Maximum
    'VMIN':    0b000110,  # VMIN - Vector Minimum
    'VSUM':    0b000111,  # VSUM - Vector Sum
    'VPERM':   0b001000,  # VPERM - Vector Permutation

    # Memory management
    'CACHE':   0b001001,  # CACHE - Cache Control
    'FLUSH':   0b001010,  # FLUSH - Cache Flush
    'MEMBAR':  0b001011,  # MEMBAR - Memory Barrier

    # System control
    'SYSCALL': 0b001100,  # SYSCALL - System Call
    'BREAK':   0b001101,  # BREAK - Debug Breakpoint
    'HALT':    0b001110   # HALT - System Halt
}

# Map instruction mnemonics to their types
INSTRUCTION_TYPES = {
    **{op: OpType.ALU for op in ALU_OPCODES},
    **{op: OpType.MEMORY for op in MEM_OPCODES},
    **{op: OpType.CONTROL for op in CTRL_OPCODES},
    **{op: OpType.VECTOR for op in VEC_OPCODES},
    **{op: OpType.FPU for op in FPU_OPCODES},
    **{op: OpType.SYSTEM for op in SYS_OPCODES},
    **{op: OpType.MICROCODE for op in MICROCODE_OPCODES}
}

# Combine all opcodes for lookup
ALL_OPCODES = {
    **ALU_OPCODES,
    **MEM_OPCODES,
    **CTRL_OPCODES,
    **VEC_OPCODES,
    **FPU_OPCODES,
    **SYS_OPCODES,
    **MICROCODE_OPCODES
}

class CodeGenError(Exception):
    """Exception raised for code generation errors"""
    def __init__(self, message: str, node: ASTNode):
        self.message = message
        self.node = node
        super().__init__(f"{message} at line {node.line}, column {node.column}")

class CodeGenerator:
    def __init__(self, debug: bool = False):
        self.binary_code = bytearray()
        self.current_address = 0
        self.labels = {}          # Maps label names to addresses
        self.symbol_references = []  # List of unresolved symbol references
        self.errors = []
        self.warnings = []
        self.debug = debug

    def generate_code(self, ast: ASTNode) -> bytearray:
        """Generate binary code from the AST"""
        # First pass: collect all labels
        self._collect_labels(ast)

        # Second pass: generate code
        self._generate_code_from_ast(ast)

        # Third pass: resolve symbol references
        self._resolve_symbols()

        return self.binary_code

    def _collect_labels(self, node: ASTNode, address: int = 0):
        """First pass: collect all labels and their addresses"""
        if node.node_type == NodeType.PROGRAM:
            # Process children in program node
            current_addr = address
            for child in node.children:
                current_addr = self._collect_labels(child, current_addr)
            return current_addr

        elif node.node_type == NodeType.LABEL:
            # Register label with current address
            self.labels[node.value] = address
            return address  # Label doesn't change address

        elif node.node_type == NodeType.INSTRUCTION:
            # Regular instruction takes 32 bits
            return address + 4

        elif node.node_type == NodeType.VLIW_INSTRUCTION:
            # VLIW instruction takes 96 bits (3x32)
            return address + 12

        elif node.node_type == NodeType.DIRECTIVE:
            # Handle directives that affect address
            directive = node.value.upper()

            if directive == '.ORG':
                # Set origin address
                if node.children and node.children[0].node_type == NodeType.IMMEDIATE:
                    new_address = self._parse_immediate(node.children[0].value)
                    return new_address

            elif directive == '.DB':
                # Data byte - each value takes 1 byte
                return address + len(node.children)

            elif directive == '.DW':
                # Data word - each value takes 4 bytes
                return address + (len(node.children) * 4)

            elif directive == '.DT':
                # Data tryte - each value takes 4 bytes (36-bit)
                return address + (len(node.children) * 4)

            elif directive == '.SPACE':
                # Reserve space
                if node.children and node.children[0].node_type == NodeType.IMMEDIATE:
                    space_size = self._parse_immediate(node.children[0].value)
                    return address + space_size

        # For nodes that don't affect address directly
        return address

    def _generate_code_from_ast(self, node: ASTNode):
        """Second pass: generate binary code from the AST"""
        if node.node_type == NodeType.PROGRAM:
            # Process all children in the program
            for child in node.children:
                self._generate_code_from_ast(child)

        elif node.node_type == NodeType.INSTRUCTION:
            # Generate code for a single instruction
            try:
                instruction_bytes = self._encode_instruction(node)
                self.binary_code.extend(instruction_bytes)
                self.current_address += len(instruction_bytes)
            except CodeGenError as e:
                self.errors.append(e)

        elif node.node_type == NodeType.VLIW_INSTRUCTION:
            # Generate code for a VLIW instruction (up to 3 operations)
            try:
                vliw_bytes = self._encode_vliw_instruction(node)
                self.binary_code.extend(vliw_bytes)
                self.current_address += len(vliw_bytes)
            except CodeGenError as e:
                self.errors.append(e)

        elif node.node_type == NodeType.DIRECTIVE:
            # Handle directives that generate code
            directive = node.value.upper()

            if directive == '.ORG':
                # Set the current address
                if node.children and node.children[0].node_type == NodeType.IMMEDIATE:
                    new_address = self._parse_immediate(node.children[0].value)
                    # If new address is less than current, we can't go back
                    if new_address < self.current_address:
                        self.warnings.append(
                            f"Warning: .ORG {new_address} is less than current address {self.current_address}, ignoring"
                        )
                    else:
                        # Pad with zeros to reach the new address
                        padding = bytearray([0] * (new_address - self.current_address))
                        self.binary_code.extend(padding)
                        self.current_address = new_address

            elif directive == '.DB':
                # Data byte - each value is 1 byte
                for child in node.children:
                    if child.node_type == NodeType.IMMEDIATE:
                        value = self._parse_immediate(child.value)
                        # Ensure the value fits in a byte
                        if value < 0 or value > 255:
                            self.warnings.append(
                                f"Warning: Value {value} at line {child.line} is outside byte range, truncating"
                            )
                        self.binary_code.append(value & 0xFF)  # Truncate to byte
                        self.current_address += 1

            elif directive == '.DW':
                # Data word - each value is 4 bytes
                for child in node.children:
                    if child.node_type == NodeType.IMMEDIATE:
                        value = self._parse_immediate(child.value)
                        # Pack as 32-bit little-endian
                        word_bytes = struct.pack("<I", value & 0xFFFFFFFF)
                        self.binary_code.extend(word_bytes)
                        self.current_address += 4

            elif directive == '.DT':
                # Data tryte - each value is a 36-bit ternary value
                # Store the ternary value as 4 bytes
                for child in node.children:
                    if child.node_type == NodeType.IMMEDIATE:
                        value = self._parse_immediate(child.value)
                        # Pack as 32-bit little-endian (simplified as placeholder)
                        word_bytes = struct.pack("<I", value & 0xFFFFFFFF)
                        self.binary_code.extend(word_bytes)
                        self.current_address += 4

            elif directive == '.SPACE':
                # Reserve space with zeros
                if node.children and node.children[0].node_type == NodeType.IMMEDIATE:
                    space_size = self._parse_immediate(node.children[0].value)
                    self.binary_code.extend(bytearray([0] * space_size))
                    self.current_address += space_size

    def _resolve_symbols(self):
        """Third pass: resolve all symbol references"""
        for ref in self.symbol_references:
            symbol_name = ref['symbol']
            if symbol_name in self.labels:
                target_address = self.labels[symbol_name]

                # Calculate offset or absolute address based on ref type
                if ref['relative']:
                    # For relative addressing (branch instructions)
                    offset = target_address - ref['pc']
                    # Check if offset fits in the immediate field
                    if offset < -1024 or offset > 1023:  # 11-bit signed immediate
                        self.errors.append(
                            f"Error: Branch target '{symbol_name}' at {target_address} is too far from {ref['pc']}"
                        )

                    # Update the immediate field in the instruction
                    self._patch_immediate(ref['address'], offset & 0x7FF, ref['size'])
                else:
                    # For absolute addressing
                    self._patch_immediate(ref['address'], target_address & 0x7FF, ref['size'])
            else:
                self.errors.append(f"Error: Undefined symbol '{symbol_name}' referenced at {ref['address']}")

    def _encode_instruction(self, node: ASTNode) -> bytearray:
        """Encode a single instruction into binary"""
        instruction = node.value.upper()

        # Check if the instruction is valid
        if instruction not in ALL_OPCODES:
            raise CodeGenError(f"Unknown instruction: {instruction}", node)

        # Get opcode and instruction type
        opcode = ALL_OPCODES[instruction]
        instr_type = INSTRUCTION_TYPES[instruction].value

        # Default fields
        reg1 = 0  # Destination register
        reg2 = 0  # Source register 1
        reg3 = 0  # Source register 2
        immediate = 0  # Immediate value or address offset
        par_flags = ParFlags.SERIAL.value  # Serial execution by default

        # Process operands based on instruction type
        operand_count = len(node.children)

        if instr_type == OpType.ALU.value:
            # ALU instructions: first operand is destination
            if operand_count >= 1:
                reg1 = self._encode_register(node.children[0])
                if operand_count >= 2:
                    reg2 = self._encode_register(node.children[1])
                    if operand_count >= 3:
                        # Third operand can be register or immediate
                        if node.children[2].node_type == NodeType.REGISTER:
                            reg3 = self._encode_register(node.children[2])
                        elif node.children[2].node_type == NodeType.IMMEDIATE:
                            immediate = self._parse_immediate(node.children[2].value)
                        elif node.children[2].node_type == NodeType.SYMBOL_REF:
                            # Symbol references handled later
                            immediate = 0
                            self._add_symbol_reference(
                                node.children[2].value,
                                self.current_address + 1,  # +1 for immediate field
                                False,  # Absolute addressing
                                self.current_address,
                                2  # 2 bytes for immediate field
                            )

        elif instr_type == OpType.MEMORY.value:
            # Memory instructions: first operand is usually register
            if operand_count >= 1:
                reg1 = self._encode_register(node.children[0])
                if operand_count >= 2:
                    # Second operand can be memory reference or immediate
                    if node.children[1].node_type == NodeType.MEMORY_REF:
                        # Memory reference [register+offset]
                        base_reg_name = node.children[1].value
                        reg2 = self._encode_register_name(base_reg_name)

                        # Check for offset
                        if len(node.children[1].children) > 0:
                            offset_node = node.children[1].children[0]
                            if offset_node.node_type == NodeType.IMMEDIATE:
                                immediate = self._parse_immediate(offset_node.value)
                            elif offset_node.node_type == NodeType.REGISTER:
                                reg3 = self._encode_register_name(offset_node.value)
                    elif node.children[1].node_type == NodeType.IMMEDIATE:
                        # Direct address
                        immediate = self._parse_immediate(node.children[1].value)
                    elif node.children[1].node_type == NodeType.SYMBOL_REF:
                        # Symbol address
                        immediate = 0
                        self._add_symbol_reference(
                            node.children[1].value,
                            self.current_address + 1,
                            False,
                            self.current_address,
                            2
                        )

        elif instr_type == OpType.CONTROL.value:
            # Control instructions: branches and jumps
            if instruction in ['BEQ', 'BNE', 'BLT', 'BGE', 'BLTU', 'BGEU']:
                # Branch instructions: compare two registers and branch to target
                if operand_count >= 3:
                    reg1 = self._encode_register(node.children[0])

                    # Second operand can be register or immediate
                    if node.children[1].node_type == NodeType.REGISTER:
                        reg2 = self._encode_register(node.children[1])
                    elif node.children[1].node_type == NodeType.IMMEDIATE:
                        immediate = self._parse_immediate(node.children[1].value)

                    # Third operand is branch target (label)
                    if node.children[2].node_type == NodeType.SYMBOL_REF:
                        # Add to symbol references list to resolve later
                        self._add_symbol_reference(
                            node.children[2].value,
                            self.current_address + 1,  # +1 for immediate field offset
                            True,  # Relative addressing for branches
                            self.current_address + 4,  # PC relative to next instruction
                            2  # 2 bytes for branch offset
                        )

            elif instruction in ['JMP', 'JAL']:
                # Jump instructions with immediate address
                if operand_count >= 1:
                    if node.children[0].node_type == NodeType.SYMBOL_REF:
                        # Jump to label
                        self._add_symbol_reference(
                            node.children[0].value,
                            self.current_address + 1,
                            False,  # Absolute address
                            self.current_address,
                            2
                        )
                    elif node.children[0].node_type == NodeType.IMMEDIATE:
                        # Jump to immediate address
                        immediate = self._parse_immediate(node.children[0].value)

            elif instruction in ['JR', 'JALR']:
                # Jump to register
                if operand_count >= 1:
                    reg1 = self._encode_register(node.children[0])

        elif instr_type == OpType.VECTOR.value:
            # Vector instructions
            if operand_count >= 1:
                reg1 = self._encode_register(node.children[0])
                if operand_count >= 2:
                    reg2 = self._encode_register(node.children[1])
                    if operand_count >= 3:
                        if node.children[2].node_type == NodeType.REGISTER:
                            reg3 = self._encode_register(node.children[2])
                        elif node.children[2].node_type == NodeType.IMMEDIATE:
                            immediate = self._parse_immediate(node.children[2].value)

        elif instr_type == OpType.FPU.value:
            # FPU instructions
            if operand_count >= 1:
                reg1 = self._encode_register(node.children[0])
                if operand_count >= 2:
                    reg2 = self._encode_register(node.children[1])
                    if operand_count >= 3:
                        if node.children[2].node_type == NodeType.REGISTER:
                            reg3 = self._encode_register(node.children[2])
                        elif node.children[2].node_type == NodeType.IMMEDIATE:
                            immediate = self._parse_immediate(node.children[2].value)

        elif instr_type == OpType.SYSTEM.value or instr_type == OpType.MICROCODE.value:
            # System and microcode instructions may have specific operand handling
            # For simplicity, use the same approach as ALU operations
            if operand_count >= 1:
                reg1 = self._encode_register(node.children[0])
                if operand_count >= 2:
                    reg2 = self._encode_register(node.children[1])
                    if operand_count >= 3:
                        if node.children[2].node_type == NodeType.REGISTER:
                            reg3 = self._encode_register(node.children[2])
                        elif node.children[2].node_type == NodeType.IMMEDIATE:
                            immediate = self._parse_immediate(node.children[2].value)

        # Ensure immediate value fits in the field
        immediate &= 0x7FF  # 11 bits for immediate

        # Construct the 32-bit instruction
        instruction_word = (
            (opcode & 0x3F) << 26 |      # Opcode (6 bits)
            (reg1 & 0x7) << 23 |         # Reg1 (3 bits)
            (reg2 & 0x7) << 20 |         # Reg2 (3 bits)
            (reg3 & 0x7) << 17 |         # Reg3 (3 bits)
            (immediate & 0x7FF) << 6 |   # Immediate (11 bits)
            (instr_type & 0x7) << 3 |    # Type (3 bits)
            (par_flags & 0x7)            # Par (3 bits)
        )

        # Convert to byte array (little-endian)
        instruction_bytes = struct.pack("<I", instruction_word)

        if self.debug:
            print(f"Encoded {instruction} to {instruction_word:08x} at {self.current_address:08x}")

        return bytearray(instruction_bytes)

    def _encode_vliw_instruction(self, node: ASTNode) -> bytearray:
        """Encode a VLIW instruction (1-3 operations) into binary"""
        vliw_bytes = bytearray()
        op_count = len(node.children)

        if op_count < 1 or op_count > 3:
            raise CodeGenError(f"VLIW instruction must have 1-3 operations, found {op_count}", node)

        # Process each operation
        for i, op_node in enumerate(node.children):
            # Encode individual operations
            if op_node.node_type == NodeType.INSTRUCTION:
                # Set parallel execution flags
                if i > 0:
                    # Set previous operation's par field to indicate parallel execution
                    # This is a simplification - in a real assembler we'd need more sophisticated scheduling
                    prev_op_bytes = vliw_bytes[-4:]
                    # Extract the instruction word
                    prev_word = struct.unpack("<I", prev_op_bytes)[0]
                    # Update the par field (lowest 3 bits) to indicate parallel execution
                    prev_word = (prev_word & ~0x7) | ParFlags.FULL_PARALLEL.value
                    # Pack it back
                    updated_prev_bytes = struct.pack("<I", prev_word)
                    # Replace in the byte array
                    vliw_bytes[-4:] = updated_prev_bytes

                # Encode this operation
                op_bytes = self._encode_instruction(op_node)
                vliw_bytes.extend(op_bytes)
            else:
                raise CodeGenError(f"Expected instruction in VLIW block, got {op_node.node_type}", op_node)

        # Pad with NOPs if fewer than 3 operations
        while len(vliw_bytes) < 12:  # 12 bytes for full VLIW word (3 x 4 bytes)
            # Create a NOP instruction
            nop_node = ASTNode(NodeType.INSTRUCTION, "NOP", node.line, node.column)
            nop_bytes = self._encode_instruction(nop_node)
            vliw_bytes.extend(nop_bytes)

        return vliw_bytes

    def _encode_register(self, node: ASTNode) -> int:
        """Encode a register operand to its binary representation"""
        if node.node_type != NodeType.REGISTER:
            raise CodeGenError(f"Expected register operand, got {node.node_type}", node)

        register_name = node.value.upper()
        return self._encode_register_name(register_name)

    def _encode_register_name(self, register_name: str) -> int:
        """Encode a register name to its binary representation"""
        register_name = register_name.upper()
        if register_name in GPR_ENCODING:
            return GPR_ENCODING[register_name]
        elif register_name in SPECIAL_REG_ENCODING:
            return SPECIAL_REG_ENCODING[register_name]
        else:
            # For vector and floating-point registers, we'd need special handling
            # This is a simplification - in reality, we'd use different encoding
            return 0  # Default as a placeholder

    def _parse_immediate(self, immediate_str: str) -> int:
        """Parse an immediate value from string to integer"""
        if immediate_str.startswith('0x'):
            # Hexadecimal
            return int(immediate_str[2:], 16)
        elif immediate_str.startswith('0b'):
            # Binary
            return int(immediate_str[2:], 2)
        elif immediate_str.startswith('0t'):
            # Ternary
            return ternary_to_decimal(immediate_str)
        else:
            # Decimal
            return int(immediate_str)

    def _add_symbol_reference(self, symbol: str, address: int, relative: bool, pc: int, size: int):
        """Add a symbol reference to be resolved later"""
        self.symbol_references.append({
            'symbol': symbol,
            'address': address,  # Address where the reference occurs
            'relative': relative,  # Whether it's a relative reference (e.g., branch)
            'pc': pc,  # Program counter for relative addressing
            'size': size  # Size in bytes of the field to patch
        })

    def _patch_immediate(self, address: int, value: int, size: int):
        """Patch an immediate field at the given address"""
        # Calculate offset in binary code
        offset = address - 0  # Assuming binary_code starts at address 0

        if offset < 0 or offset + size > len(self.binary_code):
            self.errors.append(f"Error: Cannot patch address {address}, out of range")
            return

        # For a 2-byte immediate field in little-endian
        if size == 2:
            self.binary_code[offset] = value & 0xFF
            self.binary_code[offset+1] = (value >> 8) & 0xFF
        elif size == 4:
            # For a 4-byte word
            word_bytes = struct.pack("<I", value)
            self.binary_code[offset:offset+4] = word_bytes

    def write_binary_file(self, filename: str) -> bool:
        """Write the generated binary code to a file"""
        try:
            with open(filename, 'wb') as f:
                f.write(self.binary_code)
            return True
        except Exception as e:
            self.errors.append(f"Error writing binary file: {e}")
            return False

    def write_hex_file(self, filename: str) -> bool:
        """Write the generated binary code to an Intel HEX file"""
        try:
            with open(filename, 'w') as f:
                # Simple hex dump format, 16 bytes per line
                for i in range(0, len(self.binary_code), 16):
                    chunk = self.binary_code[i:i+16]
                    hex_values = ' '.join(f"{b:02X}" for b in chunk)
                    ascii_values = ''.join(chr(b) if 32 <= b <= 126 else '.' for b in chunk)
                    f.write(f"{i:08X}: {hex_values:<48} {ascii_values}\n")
            return True
        except Exception as e:
            self.errors.append(f"Error writing hex file: {e}")
            return False

# Usage example
if __name__ == "__main__":
    import vtx1_lexer
    from vtx1_parser import Parser

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

    # Compile the code
    lexer = vtx1_lexer.Lexer()
    tokens = lexer.tokenize(test_code)
    parser = Parser(tokens)
    ast = parser.parse()

    # Generate code
    codegen = CodeGenerator(debug=True)
    binary_code = codegen.generate_code(ast)

    # Print errors and warnings
    if codegen.errors:
        print(f"Found {len(codegen.errors)} code generation error(s):")
        for error in codegen.errors:
            print(f"  {error}")

    if codegen.warnings:
        print(f"Found {len(codegen.warnings)} warning(s):")
        for warning in codegen.warnings:
            print(f"  {warning}")

    # Print binary code summary
    print(f"\nGenerated {len(binary_code)} bytes of binary code")
    for i in range(0, len(binary_code), 4):
        chunk = binary_code[i:i+4]
        if len(chunk) == 4:  # Complete 32-bit word
            word = struct.unpack("<I", chunk)[0]
            print(f"  0x{i:04x}: 0x{word:08x}")

    # Write to files
    codegen.write_binary_file("test_output.bin")
    codegen.write_hex_file("test_output.hex")
    print("Output written to test_output.bin and test_output.hex")
