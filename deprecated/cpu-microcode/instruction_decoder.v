	`timescale 1ns / 1ps
// VTX1 Instruction Decoder
// Part of the VTX1 Ternary System-on-Chip

`ifndef INSTRUCTION_DECODER_V
`define INSTRUCTION_DECODER_V

// Include paths handled by compiler -I flags (see Taskfile.yml)
`include "ternary_constants.v"
`include "vtx1_interfaces.v"

// ============================================================================
// VTX1 INSTRUCTION DECODER
// ============================================================================
// Decodes VLIW instructions and extracts control signals for the execution
// units. Supports all 78 VTX1 instructions with proper ternary encoding.

module instruction_decoder (
    input  wire                     clk,
    input  wire                     rst_n,
    input  wire                     enable,
      // VLIW instruction input (108 bits = 3 x 36-bit slots)
    input  wire [`VTX1_VLIW_WIDTH-1:0]   instruction,
    
    // Decoded instruction fields for each slot
    output reg  [`OPCODE_WIDTH-1:0] opcode_a,
    output reg  [`OPCODE_WIDTH-1:0] opcode_b,
    output reg  [`OPCODE_WIDTH-1:0] opcode_c,
    
    output reg  [`VTX1_REG_ADDR_WIDTH-1:0] rs1_a, rs2_a, rs3_a, rd_a,
    output reg  [`VTX1_REG_ADDR_WIDTH-1:0] rs1_b, rs2_b, rs3_b, rd_b,
    output reg  [`VTX1_REG_ADDR_WIDTH-1:0] rs1_c, rs2_c, rs3_c, rd_c,
    
    output reg  [11:0]               immediate_a,  // 6 trits = 12 bits
    output reg  [11:0]               immediate_b,
    output reg  [11:0]               immediate_c,
    
    // Control signals for execution units
    output reg                       alu_enable_a,
    output reg                       alu_enable_b,
    output reg                       alu_enable_c,
    
    output reg  [3:0]                alu_operation_a,
    output reg  [3:0]                alu_operation_b,
    output reg  [3:0]                alu_operation_c,
    
    output reg                       memory_enable_a,
    output reg                       memory_enable_b,
    output reg                       memory_enable_c,
    
    output reg                       memory_write_a,
    output reg                       memory_write_b,
    output reg                       memory_write_c,
    
    output reg                       branch_enable,
    output reg  [2:0]                branch_condition,
    output reg  [`WORD_WIDTH-1:0]    branch_target,
      // Instruction validity and error detection
    output reg                       valid_a,
    output reg                       valid_b,
    output reg                       valid_c,
    output reg                       decode_error,
    
    // Enhanced error handling
    output reg                       error,
    output reg  [3:0]                error_code,
    output reg                       timeout,
    input  wire                      error_clear,
    
    // Debug interface
    output reg  [3:0]                state,
    output reg  [31:0]               decode_count,
    output reg  [31:0]               error_count
);

    // ========================================================================
    // INSTRUCTION FORMAT DEFINITIONS
    // ========================================================================
    // VTX1 uses 32-bit instruction slots with the following format:
    // [31:20] - Opcode (6 trits = 12 bits)
    // [19:16] - rs1 (2 trits = 4 bits)
    // [15:12] - rs2 (2 trits = 4 bits)  
    // [11:8]  - rs3/rd (2 trits = 4 bits)
    // [7:0]   - Immediate/Address (4 trits = 8 bits)
    
    // Extract instruction fields
    wire [`OPCODE_WIDTH-1:0] slot_opcode_a = instruction[31:20];
    wire [`OPCODE_WIDTH-1:0] slot_opcode_b = instruction[63:52];
    wire [`OPCODE_WIDTH-1:0] slot_opcode_c = instruction[95:84];
    
    wire [3:0] slot_rs1_a = instruction[19:16];
    wire [3:0] slot_rs2_a = instruction[15:12];
    wire [3:0] slot_rs3_a = instruction[11:8];
    wire [7:0] slot_imm_a = instruction[7:0];
    
    wire [3:0] slot_rs1_b = instruction[51:48];
    wire [3:0] slot_rs2_b = instruction[47:44];
    wire [3:0] slot_rs3_b = instruction[43:40];
    wire [7:0] slot_imm_b = instruction[39:32];
    
    wire [3:0] slot_rs1_c = instruction[83:80];
    wire [3:0] slot_rs2_c = instruction[79:76];
    wire [3:0] slot_rs3_c = instruction[75:72];
    wire [7:0] slot_imm_c = instruction[71:64];
    
    // State machine states
    localparam STATE_IDLE     = 4'b0000;
    localparam STATE_VALIDATE = 4'b0001;
    localparam STATE_DECODE_A = 4'b0010;
    localparam STATE_DECODE_B = 4'b0011;
    localparam STATE_DECODE_C = 4'b0100;
    localparam STATE_COMPLETE = 4'b0101;
    localparam STATE_ERROR    = 4'b1111;
      // Internal signals for state machine
    reg [3:0] current_state, next_state;
    reg [31:0] decoder_timeout_counter;
    reg instruction_format_error;
    reg opcode_validation_error;
    reg register_validation_error;
    reg slot_error_a, slot_error_b, slot_error_c;
    
    // ========================================================================
    // OPCODE DEFINITIONS
    // ========================================================================
    // Major instruction categories (based on opcode[11:8])
    localparam CATEGORY_ALU     = 4'b0001;  // Arithmetic/Logic operations
    localparam CATEGORY_MEMORY  = 4'b0010;  // Memory operations
    localparam CATEGORY_BRANCH  = 4'b0011;  // Branch/Control operations
    localparam CATEGORY_SPECIAL = 4'b0100;  // Special operations
    
    // ALU Operations (opcode[7:0])
    localparam ALU_ADD      = 8'b00010001;  // Ternary addition
    localparam ALU_SUB      = 8'b00010010;  // Ternary subtraction
    localparam ALU_MUL      = 8'b00010011;  // Ternary multiplication
    localparam ALU_DIV      = 8'b00010100;  // Ternary division
    localparam ALU_MOD      = 8'b00010101;  // Ternary modulo
    localparam ALU_NEG      = 8'b00010110;  // Ternary negation
    localparam ALU_ABS      = 8'b00010111;  // Absolute value
    localparam ALU_CMP      = 8'b00011000;  // Compare
    localparam ALU_MAX      = 8'b00011001;  // Maximum
    localparam ALU_MIN      = 8'b00011010;  // Minimum
    localparam ALU_INC      = 8'b00011011;  // Increment
    localparam ALU_DEC      = 8'b00011100;  // Decrement
    
    // Logical Operations
    localparam LOG_AND      = 8'b00100001;  // Ternary AND
    localparam LOG_OR       = 8'b00100010;  // Ternary OR
    localparam LOG_NOT      = 8'b00100011;  // Ternary NOT
    localparam LOG_NAND     = 8'b00100100;  // Ternary NAND
    localparam LOG_NOR      = 8'b00100101;  // Ternary NOR
    localparam LOG_XOR      = 8'b00100110;  // Ternary XOR
    localparam LOG_CONSENSUS = 8'b00100111; // Consensus operation
    localparam LOG_MAJORITY = 8'b00101000; // Majority operation
    
    // Memory Operations
    localparam MEM_LOAD     = 8'b01000001;  // Load from memory
    localparam MEM_STORE    = 8'b01000010;  // Store to memory
    localparam MEM_LOAD_IMM = 8'b01000011;  // Load immediate
    localparam MEM_PUSH     = 8'b01000100;  // Push to stack
    localparam MEM_POP      = 8'b01000101;  // Pop from stack
    
    // Branch Operations
    localparam BR_JUMP      = 8'b01100001;  // Unconditional jump
    localparam BR_JUMP_REG  = 8'b01100010;  // Jump to register
    localparam BR_BRANCH_EQ = 8'b01100011;  // Branch if equal
    localparam BR_BRANCH_NE = 8'b01100100;  // Branch if not equal
    localparam BR_BRANCH_GT = 8'b01100101;  // Branch if greater
    localparam BR_BRANCH_LT = 8'b01100110;  // Branch if less
    localparam BR_BRANCH_GE = 8'b01100111;  // Branch if greater/equal
    localparam BR_BRANCH_LE = 8'b01101000;  // Branch if less/equal
    localparam BR_CALL      = 8'b01101001;  // Function call
    localparam BR_RETURN    = 8'b01101010;  // Function return
      // Special Operations
    localparam SP_NOP       = 8'b10000000;  // No operation
    localparam SP_HALT      = 8'b10000001;  // Halt processor
    localparam SP_INT       = 8'b10000010;  // Software interrupt
    localparam SP_IRET      = 8'b10000011;  // Interrupt return
    
    // ========================================================================
    // VALIDATION TASKS
    // ========================================================================
    
    // Validate instruction format
    task validate_instruction_format;
        input [95:0] instr;
        output reg format_error;
        integer i;
        reg [1:0] current_trit;
        begin
            format_error = 1'b0;
            
            // Check each pair of bits for valid ternary encoding (not 11)
            for (i = 0; i < 48; i = i + 1) begin
                current_trit = instr[i*2 +: 2];
                if (current_trit == 2'b11) begin
                    format_error = 1'b1;
                end
            end
        end
    endtask
    
    // Validate opcode
    task validate_opcode;
        input [11:0] opcode;
        output reg opcode_error;
        reg [3:0] category;
        reg [7:0] operation;
        begin
            opcode_error = 1'b0;
            category = opcode[11:8];
            operation = opcode[7:0];
            
            case (category)
                CATEGORY_ALU: begin
                    case (operation)
                        ALU_ADD, ALU_SUB, ALU_MUL, ALU_DIV, ALU_MOD,
                        ALU_NEG, ALU_ABS, ALU_CMP, ALU_MAX, ALU_MIN,
                        ALU_INC, ALU_DEC, LOG_AND, LOG_OR, LOG_NOT,
                        LOG_NAND, LOG_NOR, LOG_XOR, LOG_CONSENSUS, LOG_MAJORITY:
                            opcode_error = 1'b0;
                        default:
                            opcode_error = 1'b1;
                    endcase
                end
                
                CATEGORY_MEMORY: begin
                    case (operation)
                        MEM_LOAD, MEM_STORE, MEM_LOAD_IMM, MEM_PUSH, MEM_POP:
                            opcode_error = 1'b0;
                        default:
                            opcode_error = 1'b1;
                    endcase
                end
                
                CATEGORY_BRANCH: begin
                    case (operation)
                        BR_JUMP, BR_JUMP_REG, BR_BRANCH_EQ, BR_BRANCH_NE, 
                        BR_BRANCH_GT, BR_BRANCH_LT, BR_BRANCH_GE, BR_BRANCH_LE,
                        BR_CALL, BR_RETURN:
                            opcode_error = 1'b0;
                        default:
                            opcode_error = 1'b1;
                    endcase
                end
                
                CATEGORY_SPECIAL: begin
                    case (operation)
                        SP_NOP, SP_HALT, SP_INT, SP_IRET:
                            opcode_error = 1'b0;
                        default:
                            opcode_error = 1'b1;
                    endcase
                end
                
                default:
                    opcode_error = 1'b1;
            endcase
        end
    endtask
    
    // Validate register addresses
    task validate_register_address;
        input [3:0] reg_addr;
        output reg reg_error;
        begin
            // VTX1 has 13 registers (0-12), encoded as ternary
            // Valid ternary encodings for 0-12: 0000-1022 (excluding 11xx patterns)
            if (reg_addr[3:2] == 2'b11 || reg_addr > 4'hC) begin
                reg_error = 1'b1;
            end else begin
                reg_error = 1'b0;
            end
        end
    endtask
    
    // ========================================================================
    // INSTRUCTION DECODING LOGIC    // ========================================================================
    
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            // Reset all outputs
            current_state <= STATE_IDLE;
            opcode_a <= 12'h0;
            opcode_b <= 12'h0;
            opcode_c <= 12'h0;
            
            rs1_a <= 4'h0; rs2_a <= 4'h0; rs3_a <= 4'h0; rd_a <= 4'h0;
            rs1_b <= 4'h0; rs2_b <= 4'h0; rs3_b <= 4'h0; rd_b <= 4'h0;
            rs1_c <= 4'h0; rs2_c <= 4'h0; rs3_c <= 4'h0; rd_c <= 4'h0;
            
            immediate_a <= 12'h0;
            immediate_b <= 12'h0;
            immediate_c <= 12'h0;
            
            alu_enable_a <= 1'b0;
            alu_enable_b <= 1'b0;
            alu_enable_c <= 1'b0;
            
            alu_operation_a <= 4'h0;
            alu_operation_b <= 4'h0;
            alu_operation_c <= 4'h0;
            
            memory_enable_a <= 1'b0;
            memory_enable_b <= 1'b0;
            memory_enable_c <= 1'b0;
            
            memory_write_a <= 1'b0;
            memory_write_b <= 1'b0;
            memory_write_c <= 1'b0;
            
            branch_enable <= 1'b0;
            branch_condition <= 3'h0;
            branch_target <= `VTX1_WORD_DEFAULT;
            
            valid_a <= 1'b0;
            valid_b <= 1'b0;
            valid_c <= 1'b0;
            decode_error <= 1'b0;
            
            // Enhanced error handling
            error <= 1'b0;
            error_code <= `VTX1_ERROR_NONE;
            timeout <= 1'b0;
            state <= STATE_IDLE;
            decode_count <= 32'h0;            error_count <= 32'h0;
            decoder_timeout_counter <= 32'h0;
            instruction_format_error <= 1'b0;
            opcode_validation_error <= 1'b0;
            register_validation_error <= 1'b0;
            
        end else begin
            // Error clearing
            if (error_clear) begin
                error <= 1'b0;
                error_code <= `VTX1_ERROR_NONE;
                timeout <= 1'b0;
                if (current_state == STATE_ERROR) begin
                    current_state <= STATE_IDLE;
                end
            end
            
            // State machine
            case (current_state)                STATE_IDLE: begin
                    decoder_timeout_counter <= 32'h0;
                    
                    if (enable) begin
                        current_state <= STATE_VALIDATE;
                        decode_count <= decode_count + 1;
                    end
                end
                  STATE_VALIDATE: begin
                    // Timeout detection
                    decoder_timeout_counter <= decoder_timeout_counter + 1;
                    if (decoder_timeout_counter >= `VTX1_TIMEOUT_CYCLES) begin
                        error <= 1'b1;
                        error_code <= `VTX1_ERROR_TIMEOUT;
                        timeout <= 1'b1;
                        error_count <= error_count + 1;
                        current_state <= STATE_ERROR;
                    end else begin
                        // Validate instruction format
                        validate_instruction_format(instruction, instruction_format_error);
                        if (instruction_format_error) begin
                            error <= 1'b1;
                            error_code <= `VTX1_ERROR_INVALID_OP;
                            error_count <= error_count + 1;
                            current_state <= STATE_ERROR;
                        end else begin
                            // Validate opcodes
                            validate_opcode(slot_opcode_a, opcode_validation_error);
                            if (opcode_validation_error) begin
                                error <= 1'b1;
                                error_code <= `VTX1_ERROR_INVALID_OP;
                                error_count <= error_count + 1;
                                current_state <= STATE_ERROR;
                            end else begin
                                validate_opcode(slot_opcode_b, opcode_validation_error);
                                if (opcode_validation_error) begin
                                    error <= 1'b1;
                                    error_code <= `VTX1_ERROR_INVALID_OP;
                                    error_count <= error_count + 1;
                                    current_state <= STATE_ERROR;
                                end else begin
                                    validate_opcode(slot_opcode_c, opcode_validation_error);
                                    if (opcode_validation_error) begin
                                        error <= 1'b1;
                                        error_code <= `VTX1_ERROR_INVALID_OP;
                                        error_count <= error_count + 1;
                                        current_state <= STATE_ERROR;
                                    end else begin
                                        current_state <= STATE_DECODE_A;
                                    end
                                end
                            end
                        end
                    end
                end
                
                STATE_DECODE_A: begin
                    // Validate register addresses for slot A
                    validate_register_address(slot_rs1_a, register_validation_error);
                    if (register_validation_error) begin
                        error <= 1'b1;
                        error_code <= `VTX1_ERROR_INVALID_ADDR;
                        error_count <= error_count + 1;
                        current_state <= STATE_ERROR;
                    end else begin
                        validate_register_address(slot_rs2_a, register_validation_error);
                        if (register_validation_error) begin
                            error <= 1'b1;
                            error_code <= `VTX1_ERROR_INVALID_ADDR;
                            error_count <= error_count + 1;
                            current_state <= STATE_ERROR;
                        end else begin
                            validate_register_address(slot_rs3_a, register_validation_error);
                            if (register_validation_error) begin
                                error <= 1'b1;
                                error_code <= `VTX1_ERROR_INVALID_ADDR;
                                error_count <= error_count + 1;
                                current_state <= STATE_ERROR;                            end else begin
                                // Decode slot A
                                decode_instruction_slot(
                                    slot_opcode_a, slot_rs1_a, slot_rs2_a, slot_rs3_a, slot_imm_a,
                                    alu_enable_a, alu_operation_a, memory_enable_a, memory_write_a,
                                    rd_a, valid_a, slot_error_a
                                );
                                decode_error <= slot_error_a;
                                
                                // Extract fields for slot A
                                opcode_a <= slot_opcode_a;
                                rs1_a <= slot_rs1_a;
                                rs2_a <= slot_rs2_a;
                                rs3_a <= slot_rs3_a;
                                immediate_a <= {4'h0, slot_imm_a};
                                
                                current_state <= STATE_DECODE_B;
                            end
                        end
                    end
                end
                
                STATE_DECODE_B: begin
                    // Similar validation and decoding for slot B
                    validate_register_address(slot_rs1_b, register_validation_error);
                    if (register_validation_error) begin
                        error <= 1'b1;
                        error_code <= `VTX1_ERROR_INVALID_ADDR;
                        error_count <= error_count + 1;
                        current_state <= STATE_ERROR;                    end else begin
                        // Decode slot B
                        decode_instruction_slot(
                            slot_opcode_b, slot_rs1_b, slot_rs2_b, slot_rs3_b, slot_imm_b,
                            alu_enable_b, alu_operation_b, memory_enable_b, memory_write_b,
                            rd_b, valid_b, slot_error_b
                        );
                        decode_error <= decode_error || slot_error_b;
                        
                        // Extract fields for slot B
                        opcode_b <= slot_opcode_b;
                        rs1_b <= slot_rs1_b;
                        rs2_b <= slot_rs2_b;
                        rs3_b <= slot_rs3_b;
                        immediate_b <= {4'h0, slot_imm_b};
                        
                        current_state <= STATE_DECODE_C;
                    end
                end
                
                STATE_DECODE_C: begin
                    // Similar validation and decoding for slot C
                    validate_register_address(slot_rs1_c, register_validation_error);
                    if (register_validation_error) begin
                        error <= 1'b1;
                        error_code <= `VTX1_ERROR_INVALID_ADDR;
                        error_count <= error_count + 1;
                        current_state <= STATE_ERROR;                    end else begin
                        // Decode slot C
                        decode_instruction_slot(
                            slot_opcode_c, slot_rs1_c, slot_rs2_c, slot_rs3_c, slot_imm_c,
                            alu_enable_c, alu_operation_c, memory_enable_c, memory_write_c,
                            rd_c, valid_c, slot_error_c
                        );
                        decode_error <= decode_error || slot_error_c;
                        
                        // Extract fields for slot C
                        opcode_c <= slot_opcode_c;
                        rs1_c <= slot_rs1_c;
                        rs2_c <= slot_rs2_c;
                        rs3_c <= slot_rs3_c;
                        immediate_c <= {4'h0, slot_imm_c};
                        
                        current_state <= STATE_COMPLETE;
                    end
                end
                
                STATE_COMPLETE: begin
                    // Handle branch instructions (only from slot A for simplicity)
                    if (slot_opcode_a[11:8] == CATEGORY_BRANCH) begin
                        branch_enable <= 1'b1;
                        case (slot_opcode_a[7:0])
                            BR_BRANCH_EQ: branch_condition <= 3'b000;
                            BR_BRANCH_NE: branch_condition <= 3'b001;
                            BR_BRANCH_GT: branch_condition <= 3'b010;
                            BR_BRANCH_LT: branch_condition <= 3'b011;
                            BR_BRANCH_GE: branch_condition <= 3'b100;
                            BR_BRANCH_LE: branch_condition <= 3'b101;
                            default: branch_condition <= 3'b111; // Invalid
                        endcase
                        
                        // Branch target from immediate or register
                        if (slot_opcode_a[7:0] == BR_JUMP_REG) begin
                            // Target will be read from register rs1_a
                            branch_target <= `VTX1_WORD_DEFAULT;
                        end else begin
                            // Sign-extend immediate to form branch target
                            branch_target <= {{28{slot_imm_a[7]}}, slot_imm_a};
                        end
                    end
                    
                    current_state <= STATE_IDLE;
                end
                
                STATE_ERROR: begin
                    // Reset all control signals in error state
                    alu_enable_a <= 1'b0;
                    alu_enable_b <= 1'b0;
                    alu_enable_c <= 1'b0;
                    memory_enable_a <= 1'b0;
                    memory_enable_b <= 1'b0;
                    memory_enable_c <= 1'b0;
                    branch_enable <= 1'b0;
                    valid_a <= 1'b0;
                    valid_b <= 1'b0;
                    valid_c <= 1'b0;
                    decode_error <= 1'b1;
                    // Stay in error state until error_clear is asserted
                end
                
                default: begin
                    current_state <= STATE_IDLE;
                end
            endcase
            
            // Update state output for debugging
            state <= current_state;
        end
    end
      // ========================================================================
    // INSTRUCTION SLOT DECODER TASK
    // ========================================================================
    
    task decode_instruction_slot;
        input [`OPCODE_WIDTH-1:0] opcode;
        input [3:0] rs1, rs2, rs3;
        input [7:0] immediate;
        output reg alu_enable;
        output reg [3:0] alu_operation;
        output reg memory_enable;
        output reg memory_write;
        output reg [3:0] rd;
        output reg valid;
        output reg error;
        
        begin
            error = 1'b0;
            alu_enable = 1'b0;
            memory_enable = 1'b0;
            memory_write = 1'b0;
            alu_operation = 4'h0;
            rd = rs3;  // Default destination
            valid = 1'b0;
            
            case (opcode[11:8])
                CATEGORY_ALU: begin
                    alu_enable = 1'b1;
                    valid = 1'b1;
                    
                    case (opcode[7:0])
                        ALU_ADD: alu_operation = 4'b0000;
                        ALU_SUB: alu_operation = 4'b0001;
                        ALU_MUL: alu_operation = 4'b0010;
                        ALU_DIV: alu_operation = 4'b0011;
                        ALU_MOD: alu_operation = 4'b0100;
                        ALU_NEG: alu_operation = 4'b0101;
                        ALU_ABS: alu_operation = 4'b0110;
                        ALU_CMP: alu_operation = 4'b0111;
                        ALU_MAX: alu_operation = 4'b1000;
                        ALU_MIN: alu_operation = 4'b1001;
                        ALU_INC: alu_operation = 4'b1010;
                        ALU_DEC: alu_operation = 4'b1011;
                        
                        LOG_AND: alu_operation = 4'b0000;  // Map to logic unit
                        LOG_OR:  alu_operation = 4'b0001;
                        LOG_NOT: alu_operation = 4'b0010;
                        LOG_NAND: alu_operation = 4'b0011;
                        LOG_NOR:  alu_operation = 4'b0100;
                        LOG_XOR:  alu_operation = 4'b0101;
                        LOG_CONSENSUS: alu_operation = 4'b0110;
                        LOG_MAJORITY:  alu_operation = 4'b0111;
                        
                        default: begin
                            error = 1'b1;
                            valid = 1'b0;
                        end
                    endcase
                end
                
                CATEGORY_MEMORY: begin
                    memory_enable = 1'b1;
                    valid = 1'b1;
                    
                    case (opcode[7:0])
                        MEM_LOAD: begin
                            memory_write = 1'b0;
                            rd = rs3;
                        end
                        MEM_STORE: begin
                            memory_write = 1'b1;
                            rd = 4'h0;  // No destination for store
                        end
                        MEM_LOAD_IMM: begin
                            memory_write = 1'b0;
                            rd = rs1;  // Load immediate to rs1
                        end
                        MEM_PUSH, MEM_POP: begin
                            memory_write = (opcode[7:0] == MEM_PUSH);
                            rd = (opcode[7:0] == MEM_POP) ? rs1 : 4'h0;
                        end
                        default: begin
                            error = 1'b1;
                            valid = 1'b0;
                        end
                    endcase
                end
                
                CATEGORY_BRANCH: begin
                    valid = 1'b1;
                    // Branch handling is done in main decoder
                end
                
                CATEGORY_SPECIAL: begin
                    valid = 1'b1;
                    case (opcode[7:0])
                        SP_NOP: begin
                            // No operation - valid but no enables
                        end
                        SP_HALT: begin
                            // Halt operation - handled by control unit
                        end
                        default: begin
                            error = 1'b1;
                            valid = 1'b0;
                        end
                    endcase
                end
                
                default: begin
                    error = 1'b1;
                    valid = 1'b0;
                end
            endcase
        end
    endtask

endmodule

`endif // INSTRUCTION_DECODER_V

