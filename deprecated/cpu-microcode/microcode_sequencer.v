	`timescale 1ns / 1ps
// VTX1 Microcode Sequencer - Complete Implementation
// Part of the VTX1 Ternary System-on-Chip
// Controls execution of complex microcode operations

`ifndef MICROCODE_SEQUENCER_V
`define MICROCODE_SEQUENCER_V

`include "ternary_constants.v"
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module microcode_sequencer (
    input  wire                     clk,
    input  wire                     rst_n,
    
    // CPU Interface
    input  wire                     enable,
    input  wire                     start,
    input  wire [5:0]               opcode,          // 6-bit microcode operation
    input  wire [`VTX1_WORD_WIDTH-1:0] operand_a,
    input  wire [`VTX1_WORD_WIDTH-1:0] operand_b,
    input  wire [`VTX1_WORD_WIDTH-1:0] operand_c,
    
    output reg  [`VTX1_WORD_WIDTH-1:0] result,
    output reg                      valid,
    output reg                      ready,
    output reg                      error,
    output reg  [3:0]               error_code,
    
    // Microcode ROM Interface
    output reg  [9:0]               rom_addr,
    input  wire [31:0]              rom_data,
    output reg                      rom_enable,
    input  wire                     rom_ready,
    input  wire                     rom_error,
    
    // TCU Interface (for executing microoperations)
    output reg                      tcu_enable,
    output reg  [3:0]               tcu_operation,
    output reg  [`VTX1_WORD_WIDTH-1:0] tcu_operand_a,
    output reg  [`VTX1_WORD_WIDTH-1:0] tcu_operand_b,
    output reg  [`VTX1_WORD_WIDTH-1:0] tcu_operand_c,
    input  wire [`VTX1_WORD_WIDTH-1:0] tcu_result,
    input  wire                     tcu_valid,
    input  wire                     tcu_ready,
    input  wire                     tcu_error,
    
    // Register File Interface (for microcode register access)
    output reg  [3:0]               reg_read_addr_a,
    output reg  [3:0]               reg_read_addr_b,
    output reg  [3:0]               reg_read_addr_c,
    input  wire [`VTX1_WORD_WIDTH-1:0] reg_read_data_a,
    input  wire [`VTX1_WORD_WIDTH-1:0] reg_read_data_b,
    input  wire [`VTX1_WORD_WIDTH-1:0] reg_read_data_c,
    
    output reg                      reg_write_enable,
    output reg  [3:0]               reg_write_addr,
    output reg  [`VTX1_WORD_WIDTH-1:0] reg_write_data,
    
    // Debug Interface
    output reg  [3:0]               microcode_state,
    output reg  [31:0]              instruction_count,
    output reg  [31:0]              cycle_count,
    output reg  [31:0]              operation_count
);

    // ========================================================================
    // MICROCODE SEQUENCER STATE MACHINE
    // ========================================================================
    
    localparam STATE_IDLE          = 4'h0;
    localparam STATE_FETCH         = 4'h1;
    localparam STATE_DECODE        = 4'h2;
    localparam STATE_EXECUTE       = 4'h3;
    localparam STATE_WAIT_TCU      = 4'h4;
    localparam STATE_BRANCH        = 4'h5;
    localparam STATE_COMPLETE      = 4'h6;
    localparam STATE_ERROR         = 4'hF;
    
    reg [3:0] current_state, next_state;
    
    // ========================================================================
    // MICROCODE INSTRUCTION DECODE
    // ========================================================================
    
    // Microcode word format: [31:20] Control, [19:8] Data, [7:0] Next Address
    wire [11:0] control_field = rom_data[31:20];
    wire [11:0] data_field    = rom_data[19:8];
    wire [7:0]  next_addr     = rom_data[7:0];
    
    // Control field decode
    wire [2:0] exec_unit    = control_field[11:9];   // Execution unit selection
    wire [3:0] operation    = control_field[8:5];    // Operation code
    wire [4:0] control_bits = control_field[4:0];    // Control signals
    
    // Execution unit types
    localparam UNIT_NOP     = 3'b000;  // No operation
    localparam UNIT_ALU     = 3'b001;  // Arithmetic Logic Unit
    localparam UNIT_FPU     = 3'b010;  // Floating Point Unit
    localparam UNIT_VEC     = 3'b011;  // Vector Unit
    localparam UNIT_MEM     = 3'b100;  // Memory Unit
    localparam UNIT_CTRL    = 3'b101;  // Control Unit
    localparam UNIT_BRANCH  = 3'b110;  // Branch Unit
    localparam UNIT_RETURN  = 3'b111;  // Return/Exit
      // ========================================================================
    // MICROCODE REGISTERS AND CONTROL
    // ========================================================================
    
    reg [9:0]  microcode_pc;
    reg [9:0]  entry_point;
    reg [31:0] microcode_cycles;
    reg [31:0] timeout_counter;
    reg [31:0] adaptive_timeout_limit;     // Adaptive timeout based on operation
    reg        microcode_active;
    reg        branch_taken;
    reg [9:0]  branch_target;
    
    // Enhanced TCU handshaking
    reg        tcu_request_pending;        // Track pending TCU requests
    reg        tcu_operation_started;      // Track when operation actually started
    reg [2:0]  handshake_state;           // TCU handshaking state machine
    reg [31:0] tcu_operation_cycles;      // Track cycles for current TCU operation
    
    // Handshaking states
    localparam HS_IDLE      = 3'b000;
    localparam HS_REQUEST   = 3'b001;
    localparam HS_WAIT_ACK  = 3'b010;
    localparam HS_EXECUTING = 3'b011;
    localparam HS_WAIT_DONE = 3'b100;
    localparam HS_COMPLETE  = 3'b101;
    
    // Temporary registers for microcode execution
    reg [`VTX1_WORD_WIDTH-1:0] temp_reg_a, temp_reg_b, temp_reg_c;
    reg [`VTX1_WORD_WIDTH-1:0] microcode_result;
      // ========================================================================
    // ADAPTIVE TIMEOUT CALCULATION
    // ========================================================================
    
    // Calculate timeout based on operation complexity
    function [31:0] get_adaptive_timeout;
        input [5:0] op;
        begin
            case (op)
                // Complex arithmetic operations need more time
                6'h00, 6'h01: get_adaptive_timeout = 32'd2000;  // DIV, MOD
                6'h02, 6'h03: get_adaptive_timeout = 32'd1800;  // UDIV, UMOD  
                6'h04:        get_adaptive_timeout = 32'd1500;  // SQRT
                // Transcendental functions need significant time
                6'h08, 6'h09, 6'h0A: get_adaptive_timeout = 32'd3000;  // SIN, COS, TAN
                6'h0B, 6'h0C, 6'h0D: get_adaptive_timeout = 32'd3000;  // ASIN, ACOS, ATAN
                6'h0E, 6'h0F:        get_adaptive_timeout = 32'd2500;  // EXP, LOG
                // Vector operations
                6'h10, 6'h11: get_adaptive_timeout = 32'd1200;  // VDOT, VREDUCE
                6'h12, 6'h13, 6'h14, 6'h15: get_adaptive_timeout = 32'd1000;  // VMAX, VMIN, VSUM, VPERM
                // Memory and system operations
                6'h18, 6'h19, 6'h1A: get_adaptive_timeout = 32'd500;   // CACHE, FLUSH, MEMBAR
                6'h20, 6'h21, 6'h22: get_adaptive_timeout = 32'd800;   // SYSCALL, BREAK, HALT
                // Simple operations
                6'h05:        get_adaptive_timeout = 32'd300;   // ABS
                default:      get_adaptive_timeout = 32'd1000;  // Default timeout
            endcase
        end    endfunction
    
    // Get entry point for microcode operation
    function [9:0] get_entry_point;
        input [5:0] op;
        begin
            case (op)
                6'h00: get_entry_point = 10'h000;  // DIV
                6'h01: get_entry_point = 10'h018;  // MOD
                6'h02: get_entry_point = 10'h030;  // UDIV
                6'h03: get_entry_point = 10'h048;  // UMOD
                6'h04: get_entry_point = 10'h060;  // SQRT
                6'h05: get_entry_point = 10'h078;  // ABS
                6'h08: get_entry_point = 10'h100;  // SIN
                6'h09: get_entry_point = 10'h120;  // COS
                6'h0A: get_entry_point = 10'h140;  // TAN
                6'h0B: get_entry_point = 10'h160;  // ASIN
                6'h0C: get_entry_point = 10'h180;  // ACOS
                6'h0D: get_entry_point = 10'h1A0;  // ATAN
                6'h0E: get_entry_point = 10'h1C0;  // EXP
                6'h0F: get_entry_point = 10'h1E0;  // LOG
                6'h10: get_entry_point = 10'h200;  // VDOT
                6'h11: get_entry_point = 10'h210;  // VREDUCE
                6'h12: get_entry_point = 10'h220;  // VMAX
                6'h13: get_entry_point = 10'h230;  // VMIN
                6'h14: get_entry_point = 10'h240;  // VSUM
                6'h15: get_entry_point = 10'h250;  // VPERM
                6'h18: get_entry_point = 10'h300;  // CACHE
                6'h19: get_entry_point = 10'h310;  // FLUSH
                6'h1A: get_entry_point = 10'h320;  // MEMBAR
                6'h20: get_entry_point = 10'h380;  // SYSCALL
                6'h21: get_entry_point = 10'h390;  // BREAK
                6'h22: get_entry_point = 10'h3A0;  // HALT
                default: get_entry_point = 10'h3F0; // Error handler
            endcase
        end
    endfunction
    
    // ========================================================================
    // STATE MACHINE CONTROL
    // ========================================================================
    
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            current_state <= STATE_IDLE;
        end else begin
            current_state <= next_state;
        end
    end
    
    // Next state logic
    always @(*) begin
        next_state = current_state;
        
        case (current_state)
            STATE_IDLE: begin
                if (enable && start) begin
                    next_state = STATE_FETCH;
                end
            end
            
            STATE_FETCH: begin
                if (rom_ready && !rom_error) begin
                    next_state = STATE_DECODE;
                end else if (rom_error) begin
                    next_state = STATE_ERROR;
                end
            end
            
            STATE_DECODE: begin
                if (exec_unit == UNIT_RETURN) begin
                    next_state = STATE_COMPLETE;
                end else if (exec_unit == UNIT_BRANCH) begin
                    next_state = STATE_BRANCH;
                end else begin
                    next_state = STATE_EXECUTE;
                end
            end
              STATE_EXECUTE: begin
                if (exec_unit == UNIT_NOP) begin
                    next_state = STATE_FETCH;  // Immediate completion for NOP
                end else if (handshake_state == HS_COMPLETE) begin
                    next_state = STATE_FETCH;  // Operation completed successfully
                end else if (handshake_state == HS_IDLE && tcu_ready) begin
                    next_state = STATE_WAIT_TCU;  // Start new operation
                end else if (timeout_counter >= adaptive_timeout_limit) begin
                    next_state = STATE_ERROR;  // Timeout occurred
                end
            end
              STATE_WAIT_TCU: begin
                if (handshake_state == HS_COMPLETE && !tcu_error) begin
                    next_state = STATE_FETCH;  // Next microinstruction
                end else if (tcu_error || handshake_state == HS_IDLE) begin
                    next_state = STATE_ERROR;  // Error or handshake failed
                end else if (timeout_counter >= adaptive_timeout_limit) begin
                    next_state = STATE_ERROR;  // Timeout occurred
                end
            end
            
            STATE_BRANCH: begin
                next_state = STATE_FETCH;  // Branch handled, fetch next
            end
            
            STATE_COMPLETE: begin
                if (!enable) begin
                    next_state = STATE_IDLE;
                end
            end
            
            STATE_ERROR: begin
                if (!enable) begin
                    next_state = STATE_IDLE;
                end
            end
            
            default: begin
                next_state = STATE_IDLE;
            end
        endcase
    end
    
    // ========================================================================
    // MICROCODE EXECUTION LOGIC
    // ========================================================================
    
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            microcode_pc <= 10'h0;
            entry_point <= 10'h0;
            microcode_cycles <= 32'h0;
            timeout_counter <= 32'h0;
            adaptive_timeout_limit <= 32'd1000;
            microcode_active <= 1'b0;
            branch_taken <= 1'b0;
            branch_target <= 10'h0;
            
            // Enhanced handshaking reset
            tcu_request_pending <= 1'b0;
            tcu_operation_started <= 1'b0;
            handshake_state <= HS_IDLE;
            tcu_operation_cycles <= 32'h0;
            
            temp_reg_a <= {`VTX1_WORD_WIDTH{1'b0}};
            temp_reg_b <= {`VTX1_WORD_WIDTH{1'b0}};
            temp_reg_c <= {`VTX1_WORD_WIDTH{1'b0}};
            microcode_result <= {`VTX1_WORD_WIDTH{1'b0}};
            
            rom_addr <= 10'h0;
            rom_enable <= 1'b0;
            
            tcu_enable <= 1'b0;
            tcu_operation <= 4'h0;
            tcu_operand_a <= {`VTX1_WORD_WIDTH{1'b0}};
            tcu_operand_b <= {`VTX1_WORD_WIDTH{1'b0}};
            tcu_operand_c <= {`VTX1_WORD_WIDTH{1'b0}};
            
            reg_read_addr_a <= 4'h0;
            reg_read_addr_b <= 4'h0;
            reg_read_addr_c <= 4'h0;
            reg_write_enable <= 1'b0;
            reg_write_addr <= 4'h0;
            reg_write_data <= {`VTX1_WORD_WIDTH{1'b0}};
              result <= {`VTX1_WORD_WIDTH{1'b0}};
            valid <= 1'b0;
            ready <= 1'b1;  // Ready to accept operations after reset            error <= 1'b0;
            error_code <= 4'h0;
            
            current_state <= STATE_IDLE;
            instruction_count <= 32'h0;
            cycle_count <= 32'h0;
            operation_count <= 32'h0;
            
        end else begin            // Update debug counters
            if (enable) begin
                cycle_count <= cycle_count + 1;
                
                // Adaptive timeout management
                if (timeout_counter < adaptive_timeout_limit) begin
                    timeout_counter <= timeout_counter + 1;
                end else if (!error) begin
                    // Only set timeout error if not already in error state
                    error <= 1'b1;
                    error_code <= `ERR_TIMEOUT;
                end
                
                // Track TCU operation cycles
                if (handshake_state != HS_IDLE && handshake_state != HS_COMPLETE) begin
                    tcu_operation_cycles <= tcu_operation_cycles + 1;
                end
            end
            
            // ========================================================================
            // ENHANCED TCU HANDSHAKING STATE MACHINE
            // ========================================================================
            
            case (handshake_state)
                HS_IDLE: begin
                    tcu_enable <= 1'b0;
                    tcu_request_pending <= 1'b0;
                    tcu_operation_started <= 1'b0;
                    tcu_operation_cycles <= 32'h0;
                end
                
                HS_REQUEST: begin
                    if (tcu_ready && !tcu_request_pending) begin
                        tcu_enable <= 1'b1;
                        tcu_request_pending <= 1'b1;
                        handshake_state <= HS_WAIT_ACK;
                    end
                end
                
                HS_WAIT_ACK: begin
                    if (!tcu_ready && tcu_request_pending) begin
                        // TCU acknowledged the request
                        tcu_operation_started <= 1'b1;
                        handshake_state <= HS_EXECUTING;
                    end else if (tcu_operation_cycles > 32'd10) begin
                        // Handshake timeout - TCU didn't respond
                        handshake_state <= HS_IDLE;
                        error <= 1'b1;
                        error_code <= `ERR_TIMEOUT;
                    end
                end
                
                HS_EXECUTING: begin
                    tcu_enable <= 1'b0;  // Deassert enable after acknowledgment
                    if (tcu_valid && !tcu_error) begin
                        handshake_state <= HS_WAIT_DONE;
                    end else if (tcu_error) begin
                        handshake_state <= HS_IDLE;
                        error <= 1'b1;
                        error_code <= `ERR_BUS_ERROR;  // TCU reported error
                    end
                end
                
                HS_WAIT_DONE: begin
                    if (!tcu_valid) begin
                        // TCU has deasserted valid - operation complete
                        handshake_state <= HS_COMPLETE;
                    end
                end
                
                HS_COMPLETE: begin
                    // Operation completed successfully
                    handshake_state <= HS_IDLE;
                end
                
                default: begin
                    handshake_state <= HS_IDLE;
                end
            endcase
            
            case (current_state)                STATE_IDLE: begin
                    ready <= 1'b1;
                    valid <= 1'b0;
                    error <= 1'b0;
                    timeout_counter <= 32'h0;
                    handshake_state <= HS_IDLE;
                    
                    if (enable && start) begin
                        entry_point <= get_entry_point(opcode);
                        microcode_pc <= get_entry_point(opcode);
                        microcode_active <= 1'b1;
                        microcode_cycles <= 32'h0;
                        operation_count <= operation_count + 1;
                        
                        // Set adaptive timeout based on operation complexity
                        adaptive_timeout_limit <= get_adaptive_timeout(opcode);
                        
                        // Load initial operands
                        temp_reg_a <= operand_a;
                        temp_reg_b <= operand_b;
                        temp_reg_c <= operand_c;
                        
                        ready <= 1'b0;
                    end
                end
                
                STATE_FETCH: begin
                    rom_addr <= microcode_pc;
                    rom_enable <= 1'b1;
                    microcode_cycles <= microcode_cycles + 1;
                    instruction_count <= instruction_count + 1;
                end
                
                STATE_DECODE: begin
                    rom_enable <= 1'b0;
                    
                    // Set up register addresses based on data field
                    reg_read_addr_a <= data_field[11:8];
                    reg_read_addr_b <= data_field[7:4];
                    reg_read_addr_c <= data_field[3:0];
                    
                    // Decode operation for TCU
                    tcu_operation <= operation;
                end
                  STATE_EXECUTE: begin
                    case (exec_unit)
                        UNIT_NOP: begin
                            // No operation - just advance PC
                            microcode_pc <= microcode_pc + 1;
                        end
                        
                        UNIT_ALU, UNIT_FPU, UNIT_VEC: begin
                            // Start enhanced TCU handshaking
                            if (handshake_state == HS_IDLE) begin
                                handshake_state <= HS_REQUEST;
                                // Set up operands for TCU
                                tcu_operand_a <= (data_field[11:8] == 4'hA) ? temp_reg_a : reg_read_data_a;
                                tcu_operand_b <= (data_field[7:4] == 4'hB) ? temp_reg_b : reg_read_data_b;
                                tcu_operand_c <= (data_field[3:0] == 4'hC) ? temp_reg_c : reg_read_data_c;
                            end
                        end
                        
                        UNIT_MEM: begin
                            // Memory operations with enhanced handshaking
                            if (handshake_state == HS_IDLE) begin
                                handshake_state <= HS_REQUEST;
                                tcu_operand_a <= reg_read_data_a;
                                tcu_operand_b <= reg_read_data_b;
                                tcu_operand_c <= {`VTX1_WORD_WIDTH{1'b0}};
                            end
                        end
                        
                        UNIT_CTRL: begin
                            // Control operations - immediate completion
                            microcode_pc <= microcode_pc + 1;
                        end
                    endcase
                end
                
                STATE_WAIT_TCU: begin
                    // Enhanced handshaking completion handling
                    if (handshake_state == HS_COMPLETE && !tcu_error) begin
                        // Store TCU result with proper validation
                        if (data_field[11:8] != 4'h0) begin
                            reg_write_enable <= 1'b1;
                            reg_write_addr <= data_field[11:8];
                            reg_write_data <= tcu_result;
                        end
                        
                        // Update temporary result
                        microcode_result <= tcu_result;
                        
                        // Advance PC with proper address handling
                        microcode_pc <= (next_addr == 8'h00) ? (microcode_pc + 1) : {2'b00, next_addr};
                        
                        // Reset handshaking for next operation
                        handshake_state <= HS_IDLE;
                          end else if (tcu_error || handshake_state == HS_IDLE) begin
                        error <= 1'b1;
                        error_code <= tcu_error ? `ERR_BUS_ERROR : `ERR_TIMEOUT;
                    end
                end
                
                STATE_BRANCH: begin
                    // Handle branch operations
                    case (operation[1:0])
                        2'b00: branch_taken <= 1'b1;  // Unconditional branch
                        2'b01: branch_taken <= (temp_reg_a == {`VTX1_WORD_WIDTH{1'b0}});  // Branch if zero
                        2'b10: branch_taken <= (temp_reg_a[`VTX1_WORD_WIDTH-1] == 1'b1);  // Branch if negative
                        2'b11: branch_taken <= (temp_reg_a != {`VTX1_WORD_WIDTH{1'b0}});  // Branch if not zero
                    endcase
                    
                    if (branch_taken) begin
                        microcode_pc <= {2'b00, next_addr};
                    end else begin
                        microcode_pc <= microcode_pc + 1;
                    end
                end
                
                STATE_COMPLETE: begin
                    microcode_active <= 1'b0;
                    valid <= 1'b1;
                    ready <= 1'b1;
                    result <= microcode_result;
                    rom_enable <= 1'b0;
                    reg_write_enable <= 1'b0;
                end
                
                STATE_ERROR: begin
                    microcode_active <= 1'b0;
                    valid <= 1'b0;
                    ready <= 1'b1;
                    error <= 1'b1;
                    rom_enable <= 1'b0;
                    reg_write_enable <= 1'b0;
                end
            endcase
              // Update state output for debugging
            microcode_state <= current_state;
        end
    end
    
    // synthesis translate_off
    always @(posedge clk) begin
        if (current_state == STATE_FETCH && rom_ready) begin
            $display("Microcode: PC=0x%03X, instr=0x%08X, cycles=%d", 
                     microcode_pc, rom_data, microcode_cycles);
        end
        
        if (current_state == STATE_COMPLETE) begin
            $display("Microcode: Operation 0x%02X completed in %d cycles, result=0x%08X",
                     opcode, microcode_cycles, microcode_result);
        end
    end
    // synthesis translate_on

endmodule

`endif // MICROCODE_SEQUENCER_V

