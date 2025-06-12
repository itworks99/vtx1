	`timescale 1ns / 1ps
// VTX1 Enhanced TCU Interface
// Provides enhanced handshaking and timing for microcode integration

`ifndef TCU_ENHANCED_INTERFACE_V
`define TCU_ENHANCED_INTERFACE_V

`include "ternary_constants.v"
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module tcu_enhanced_interface (
    input  wire                     clk,
    input  wire                     rst_n,
    
    // Microcode Sequencer Interface (Enhanced)
    input  wire                     microcode_enable,
    input  wire [3:0]               microcode_operation,
    input  wire [`VTX1_WORD_WIDTH-1:0] microcode_operand_a,
    input  wire [`VTX1_WORD_WIDTH-1:0] microcode_operand_b,
    input  wire [`VTX1_WORD_WIDTH-1:0] microcode_operand_c,
    output reg  [`VTX1_WORD_WIDTH-1:0] microcode_result,
    output reg                      microcode_valid,
    output reg                      microcode_ready,
    output reg                      microcode_error,
    
    // Legacy TCU Interface (Existing CPU Core)
    output reg                      tcu_enable,
    output reg [3:0]                tcu_operation,
    output reg [`VTX1_WORD_WIDTH-1:0] tcu_operand_a,
    output reg [`VTX1_WORD_WIDTH-1:0] tcu_operand_b,
    output reg [`VTX1_WORD_WIDTH-1:0] tcu_operand_c,
    input  wire [`VTX1_WORD_WIDTH-1:0] tcu_result,
    input  wire                     tcu_valid,
    input  wire                     tcu_ready,
    input  wire                     tcu_error,
    
    // Enhanced Control and Status
    output reg [3:0]                interface_state,
    output reg [31:0]               operation_cycles,
    output reg [31:0]               total_operations,
    output reg [31:0]               error_count
);    // ========================================================================
    // ENHANCED HANDSHAKING STATE MACHINE - Use VTX1 standardized states
    // ========================================================================
    
    localparam STATE_IDLE           = `VTX1_HS_STATE_IDLE;
    localparam STATE_REQUEST        = `VTX1_HS_STATE_REQUEST;
    localparam STATE_WAIT_ACK       = `VTX1_HS_STATE_WAIT_ACK;
    localparam STATE_EXECUTING      = `VTX1_HS_STATE_EXECUTING;
    localparam STATE_WAIT_RESULT    = `VTX1_HS_STATE_WAIT_DONE;
    localparam STATE_COMPLETE       = `VTX1_HS_STATE_COMPLETE;
    localparam STATE_ERROR          = `VTX1_HS_STATE_ERROR;
      reg [3:0] current_state, next_state;
    reg [31:0] state_timer;
    reg [31:0] handshake_timeout;
    reg        operation_pending;
    
    // VTX1 Error Handling Variables
    reg [3:0] vtx1_error_reg;
    reg [31:0] vtx1_error_info;
    reg vtx1_error_valid;
    
    // ========================================================================
    // ADAPTIVE TIMEOUT CALCULATION
    // ========================================================================
    
    function [31:0] get_operation_timeout;
        input [3:0] op;
        begin
            case (op)
                4'h0, 4'h1, 4'h2: get_operation_timeout = `VTX1_TIMEOUT_SIMPLE;    // ADD, SUB, MUL
                4'h3, 4'h4, 4'h5: get_operation_timeout = `VTX1_TIMEOUT_COMPLEX;   // DIV, MOD, SQRT
                4'h6, 4'h7:       get_operation_timeout = `VTX1_TIMEOUT_TRANSCENDENTAL; // TRIG functions
                4'h8, 4'h9, 4'hA: get_operation_timeout = `VTX1_TIMEOUT_VECTOR;    // Vector ops
                4'hB, 4'hC:       get_operation_timeout = `VTX1_TIMEOUT_MEMORY;     // Memory ops
                4'hD, 4'hE, 4'hF: get_operation_timeout = `VTX1_TIMEOUT_SYSTEM;    // System ops
                default:          get_operation_timeout = `VTX1_TIMEOUT_CYCLES;     // Default
            endcase
        end
    endfunction
    
    // ========================================================================
    // STATE MACHINE NEXT STATE LOGIC
    // ========================================================================
    
    always @(*) begin
        next_state = current_state;
        
        case (current_state)
            STATE_IDLE: begin
                if (microcode_enable && !operation_pending) begin
                    next_state = STATE_REQUEST;
                end
            end
            
            STATE_REQUEST: begin
                if (tcu_ready) begin
                    next_state = STATE_WAIT_ACK;
                end else if (state_timer >= `VTX1_HANDSHAKE_TIMEOUT) begin
                    next_state = STATE_ERROR;
                end
            end
            
            STATE_WAIT_ACK: begin
                if (!tcu_ready && tcu_enable) begin
                    next_state = STATE_EXECUTING;
                end else if (state_timer >= `VTX1_HANDSHAKE_TIMEOUT) begin
                    next_state = STATE_ERROR;
                end
            end
            
            STATE_EXECUTING: begin
                if (tcu_valid && !tcu_error) begin
                    next_state = STATE_WAIT_RESULT;
                end else if (tcu_error) begin
                    next_state = STATE_ERROR;
                end else if (state_timer >= handshake_timeout) begin
                    next_state = STATE_ERROR;
                end
            end
            
            STATE_WAIT_RESULT: begin
                if (!tcu_valid) begin
                    next_state = STATE_COMPLETE;
                end else if (state_timer >= `VTX1_HANDSHAKE_TIMEOUT) begin
                    next_state = STATE_ERROR;
                end
            end
            
            STATE_COMPLETE: begin
                if (!microcode_enable) begin
                    next_state = STATE_IDLE;
                end
            end
            
            STATE_ERROR: begin
                if (!microcode_enable) begin
                    next_state = STATE_IDLE;
                end
            end
            
            default: begin
                next_state = STATE_IDLE;
            end
        endcase
    end
    
    // ========================================================================
    // STATE MACHINE SEQUENTIAL LOGIC
    // ========================================================================
    
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            current_state <= STATE_IDLE;
            state_timer <= 32'h0;
            handshake_timeout <= 32'd1000;
            operation_pending <= 1'b0;
            operation_cycles <= 32'h0;
            total_operations <= 32'h0;
            error_count <= 32'h0;
            
            microcode_result <= {`VTX1_WORD_WIDTH{1'b0}};
            microcode_valid <= 1'b0;
            microcode_ready <= 1'b1;
            microcode_error <= 1'b0;
            
            tcu_enable <= 1'b0;
            tcu_operation <= 4'h0;
            tcu_operand_a <= {`VTX1_WORD_WIDTH{1'b0}};
            tcu_operand_b <= {`VTX1_WORD_WIDTH{1'b0}};
            tcu_operand_c <= {`VTX1_WORD_WIDTH{1'b0}};
            
            interface_state <= STATE_IDLE;
            
        end else begin
            current_state <= next_state;
            state_timer <= state_timer + 1;
            interface_state <= current_state;
            
            case (current_state)
                STATE_IDLE: begin
                    microcode_ready <= 1'b1;
                    microcode_valid <= 1'b0;
                    microcode_error <= 1'b0;
                    operation_pending <= 1'b0;
                    tcu_enable <= 1'b0;
                    state_timer <= 32'h0;
                    operation_cycles <= 32'h0;
                end
                
                STATE_REQUEST: begin
                    if (!operation_pending) begin
                        operation_pending <= 1'b1;
                        microcode_ready <= 1'b0;
                        
                        // Set adaptive timeout
                        handshake_timeout <= get_operation_timeout(microcode_operation);
                        
                        // Forward operation to TCU
                        tcu_operation <= microcode_operation;
                        tcu_operand_a <= microcode_operand_a;
                        tcu_operand_b <= microcode_operand_b;
                        tcu_operand_c <= microcode_operand_c;
                        
                        if (tcu_ready) begin
                            tcu_enable <= 1'b1;
                        end
                    end
                end
                
                STATE_WAIT_ACK: begin
                    if (!tcu_ready) begin
                        // TCU acknowledged request
                        state_timer <= 32'h0;  // Reset timer for execution phase
                    end
                end
                  STATE_EXECUTING: begin
                    tcu_enable <= 1'b0;  // Deassert enable after acknowledgment
                    operation_cycles <= operation_cycles + 1;
                      if (tcu_valid && !tcu_error) begin
                        microcode_result <= tcu_result;
                        state_timer <= 32'h0;  // Reset timer for result phase
                        `VTX1_CLEAR_ERROR(vtx1_error_reg, interface_state)
                    end else if (tcu_error) begin
                        `VTX1_SET_ERROR(vtx1_error_reg, interface_state, `VTX1_ERROR_TCU_FAULT)
                        vtx1_error_info <= 32'h0000_0001;
                        vtx1_error_valid <= 1'b1;
                    end
                end
                
                STATE_WAIT_RESULT: begin
                    if (!tcu_valid) begin
                        // Operation completed successfully
                        microcode_valid <= 1'b1;
                        total_operations <= total_operations + 1;
                    end
                end
                
                STATE_COMPLETE: begin
                    if (!microcode_enable) begin
                        microcode_valid <= 1'b0;
                        operation_pending <= 1'b0;
                    end                end
                
                STATE_ERROR: begin
                    `VTX1_SET_ERROR(vtx1_error_reg, interface_state, `VTX1_ERROR_TIMEOUT)
                    vtx1_error_info <= 32'h0000_0002;
                    vtx1_error_valid <= 1'b1;
                    microcode_error <= 1'b1;
                    microcode_ready <= 1'b1;
                    microcode_valid <= 1'b0;
                    tcu_enable <= 1'b0;
                    operation_pending <= 1'b0;
                    error_count <= error_count + 1;
                      if (!microcode_enable) begin
                        microcode_error <= 1'b0;
                        `VTX1_CLEAR_ERROR(vtx1_error_reg, interface_state)
                    end
                end
            endcase
            
            // Reset state timer on state changes
            if (current_state != next_state) begin
                state_timer <= 32'h0;
            end
        end
    end

endmodule

`endif // TCU_ENHANCED_INTERFACE_V

