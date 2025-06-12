`timescale 1ns / 1ps
// VTX1 Ternary Arithmetic Unit
// Fixed version compatible with Icarus Verilog
// Part of the VTX1 Ternary System-on-Chip

// Include paths handled by compiler -I flags (see Taskfile.yml)
`include "vtx1_interfaces.v"
`include "vtx1_error_macros.v"

(* noopt *)
(* keep_hierarchy = "yes" *)
module ternary_arithmetic_unit (
    input wire clk,
    input wire rst_n,
    
    // Input operands
    input wire [35:0] operand_a,
    input wire [35:0] operand_b,
    
    // Operation control
    input wire [3:0] operation,
    input wire operation_enable,
    
    // Output
    output reg [35:0] result,
    output reg operation_complete,
    output reg arithmetic_overflow,
    output reg arithmetic_underflow,
    
    // Enhanced error handling
    output reg error,
    output reg [3:0] error_code,
    output reg timeout,
    input wire error_clear,
    
    // Debug interface
    output reg [3:0] state,
    output reg [31:0] operation_count,
    output reg [31:0] error_count
);

    // Operation codes
    localparam OP_ADD    = 4'b0000;
    localparam OP_SUB    = 4'b0001;
    localparam OP_MUL    = 4'b0010;
    localparam OP_DIV    = 4'b0011;
    localparam OP_MOD    = 4'b0100;
    localparam OP_NEG    = 4'b0101;
    localparam OP_ABS    = 4'b0110;
    localparam OP_CMP    = 4'b0111;
    localparam OP_MAX    = 4'b1000;
    localparam OP_MIN    = 4'b1001;    localparam OP_INC    = 4'b1010;
    localparam OP_DEC    = 4'b1011;
    
    // State machine states
    localparam STATE_IDLE     = 4'b0000;
    localparam STATE_DECODE   = 4'b0001;
    localparam STATE_CONVERT  = 4'b0010;
    localparam STATE_EXECUTE  = 4'b0011;
    localparam STATE_RESULT   = 4'b0100;
    localparam STATE_ERROR    = 4'b1111;
    
    // Internal signals
    reg [71:0] extended_result;  // Double word = 36 trits = 72 bits
    reg signed [31:0] decimal_a, decimal_b, decimal_result;
    reg operation_valid;
    reg [31:0] timeout_counter;
    reg [3:0] current_state, next_state;
    reg input_validation_error;
    reg range_check_error;
    
    // Enhanced error handling and performance monitoring
    reg [31:0] operation_start_time;
    reg [31:0] max_operation_time;
    reg [31:0] min_operation_time;
    reg [31:0] avg_operation_time;
    reg [31:0] performance_counter;
    
    // Additional debugging signals
    reg [3:0] last_operation;
    reg [35:0] last_operand_a;
    reg [35:0] last_operand_b;
    reg [35:0] last_result;
    reg last_overflow;
    reg last_underflow;
    
    // Tasks for conversion (instead of functions for Icarus Verilog compatibility)
    task ternary_to_decimal_task;
        input [35:0] ternary_value;
        output signed [31:0] decimal_out;
        reg [1:0] current_trit;
        reg signed [31:0] accumulator;
        begin
            accumulator = 32'sd0;
            
            // Unroll the loop for each trit position (powers of 3)
            // Trit 0: 3^0 = 1
            current_trit = ternary_value[1:0];
            case (current_trit)
                2'b00:  accumulator = accumulator - 1;
                2'b10:  accumulator = accumulator + 1;
                default: ; // Zero or invalid - add nothing
            endcase
            
            // Trit 1: 3^1 = 3
            current_trit = ternary_value[3:2];
            case (current_trit)
                2'b00:  accumulator = accumulator - 3;
                2'b10:  accumulator = accumulator + 3;
                default: ;
            endcase
            
            // Trit 2: 3^2 = 9
            current_trit = ternary_value[5:4];
            case (current_trit)
                2'b00:  accumulator = accumulator - 9;
                2'b10:  accumulator = accumulator + 9;
                default: ;
            endcase
            
            // Trit 3: 3^3 = 27
            current_trit = ternary_value[7:6];
            case (current_trit)
                2'b00:  accumulator = accumulator - 27;
                2'b10:  accumulator = accumulator + 27;
                default: ;
            endcase
            
            // Trit 4: 3^4 = 81
            current_trit = ternary_value[9:8];
            case (current_trit)
                2'b00:  accumulator = accumulator - 81;
                2'b10:  accumulator = accumulator + 81;
                default: ;
            endcase
            
            // Trit 5: 3^5 = 243
            current_trit = ternary_value[11:10];
            case (current_trit)
                2'b00:  accumulator = accumulator - 243;
                2'b10:  accumulator = accumulator + 243;
                default: ;
            endcase
            
            // Trit 6: 3^6 = 729
            current_trit = ternary_value[13:12];
            case (current_trit)
                2'b00:  accumulator = accumulator - 729;
                2'b10:  accumulator = accumulator + 729;
                default: ;
            endcase
            
            // Trit 7: 3^7 = 2187
            current_trit = ternary_value[15:14];
            case (current_trit)
                2'b00:  accumulator = accumulator - 2187;
                2'b10:  accumulator = accumulator + 2187;
                default: ;
            endcase
            
            // Trit 8: 3^8 = 6561
            current_trit = ternary_value[17:16];
            case (current_trit)
                2'b00:  accumulator = accumulator - 6561;
                2'b10:  accumulator = accumulator + 6561;
                default: ;
            endcase
            
            // Trit 9: 3^9 = 19683
            current_trit = ternary_value[19:18];
            case (current_trit)
                2'b00:  accumulator = accumulator - 19683;
                2'b10:  accumulator = accumulator + 19683;
                default: ;
            endcase
            
            // Trit 10: 3^10 = 59049
            current_trit = ternary_value[21:20];
            case (current_trit)
                2'b00:  accumulator = accumulator - 59049;
                2'b10:  accumulator = accumulator + 59049;
                default: ;
            endcase
            
            // Trit 11: 3^11 = 177147
            current_trit = ternary_value[23:22];
            case (current_trit)
                2'b00:  accumulator = accumulator - 177147;
                2'b10:  accumulator = accumulator + 177147;
                default: ;
            endcase
            
            // Trit 12: 3^12 = 531441
            current_trit = ternary_value[25:24];
            case (current_trit)
                2'b00:  accumulator = accumulator - 531441;
                2'b10:  accumulator = accumulator + 531441;
                default: ;
            endcase
            
            // Trit 13: 3^13 = 1594323
            current_trit = ternary_value[27:26];
            case (current_trit)
                2'b00:  accumulator = accumulator - 1594323;
                2'b10:  accumulator = accumulator + 1594323;
                default: ;
            endcase
            
            // Trit 14: 3^14 = 4782969
            current_trit = ternary_value[29:28];
            case (current_trit)
                2'b00:  accumulator = accumulator - 4782969;
                2'b10:  accumulator = accumulator + 4782969;
                default: ;
            endcase
            
            // Trit 15: 3^15 = 14348907
            current_trit = ternary_value[31:30];
            case (current_trit)
                2'b00:  accumulator = accumulator - 14348907;
                2'b10:  accumulator = accumulator + 14348907;
                default: ;
            endcase
            
            // Trit 16: 3^16 = 43046721
            current_trit = ternary_value[33:32];
            case (current_trit)
                2'b00:  accumulator = accumulator - 43046721;
                2'b10:  accumulator = accumulator + 43046721;
                default: ;
            endcase
            
            // Trit 17: 3^17 = 129140163
            current_trit = ternary_value[35:34];
            case (current_trit)
                2'b00:  accumulator = accumulator - 129140163;
                2'b10:  accumulator = accumulator + 129140163;
                default: ;
            endcase
            
            decimal_out = accumulator;
        end
    endtask    task decimal_to_ternary_task;
        input signed [31:0] decimal_value;
        output [35:0] ternary_out;
        reg signed [31:0] working_value;
        reg [35:0] result_bits;
        reg [1:0] trit_value;
        begin
            working_value = decimal_value;
            result_bits = 36'h0;
            
            // Convert using balanced ternary representation
            // Work from highest to lowest power of 3
            
            // Trit 17: 3^17 = 129140163
            if (working_value >= 129140163) begin
                trit_value = 2'b10;
                working_value = working_value - 129140163;
            end else if (working_value <= -129140163) begin
                trit_value = 2'b00;
                working_value = working_value + 129140163;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[35:34] = trit_value;
            
            // Trit 16: 3^16 = 43046721
            if (working_value >= 43046721) begin
                trit_value = 2'b10;
                working_value = working_value - 43046721;
            end else if (working_value <= -43046721) begin
                trit_value = 2'b00;
                working_value = working_value + 43046721;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[33:32] = trit_value;
            
            // Trit 15: 3^15 = 14348907
            if (working_value >= 14348907) begin
                trit_value = 2'b10;
                working_value = working_value - 14348907;
            end else if (working_value <= -14348907) begin
                trit_value = 2'b00;
                working_value = working_value + 14348907;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[31:30] = trit_value;
            
            // Trit 14: 3^14 = 4782969
            if (working_value >= 4782969) begin
                trit_value = 2'b10;
                working_value = working_value - 4782969;
            end else if (working_value <= -4782969) begin
                trit_value = 2'b00;
                working_value = working_value + 4782969;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[29:28] = trit_value;
            
            // Trit 13: 3^13 = 1594323
            if (working_value >= 1594323) begin
                trit_value = 2'b10;
                working_value = working_value - 1594323;
            end else if (working_value <= -1594323) begin
                trit_value = 2'b00;
                working_value = working_value + 1594323;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[27:26] = trit_value;
            
            // Trit 12: 3^12 = 531441
            if (working_value >= 531441) begin
                trit_value = 2'b10;
                working_value = working_value - 531441;
            end else if (working_value <= -531441) begin
                trit_value = 2'b00;
                working_value = working_value + 531441;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[25:24] = trit_value;
            
            // Trit 11: 3^11 = 177147
            if (working_value >= 177147) begin
                trit_value = 2'b10;
                working_value = working_value - 177147;
            end else if (working_value <= -177147) begin
                trit_value = 2'b00;
                working_value = working_value + 177147;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[23:22] = trit_value;
            
            // Trit 10: 3^10 = 59049
            if (working_value >= 59049) begin
                trit_value = 2'b10;
                working_value = working_value - 59049;
            end else if (working_value <= -59049) begin
                trit_value = 2'b00;
                working_value = working_value + 59049;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[21:20] = trit_value;
            
            // Trit 9: 3^9 = 19683
            if (working_value >= 19683) begin
                trit_value = 2'b10;
                working_value = working_value - 19683;
            end else if (working_value <= -19683) begin
                trit_value = 2'b00;
                working_value = working_value + 19683;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[19:18] = trit_value;
            
            // Trit 8: 3^8 = 6561
            if (working_value >= 6561) begin
                trit_value = 2'b10;
                working_value = working_value - 6561;
            end else if (working_value <= -6561) begin
                trit_value = 2'b00;
                working_value = working_value + 6561;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[17:16] = trit_value;
            
            // Trit 7: 3^7 = 2187
            if (working_value >= 2187) begin
                trit_value = 2'b10;
                working_value = working_value - 2187;
            end else if (working_value <= -2187) begin
                trit_value = 2'b00;
                working_value = working_value + 2187;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[15:14] = trit_value;
            
            // Trit 6: 3^6 = 729
            if (working_value >= 729) begin
                trit_value = 2'b10;
                working_value = working_value - 729;
            end else if (working_value <= -729) begin
                trit_value = 2'b00;
                working_value = working_value + 729;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[13:12] = trit_value;
            
            // Trit 5: 3^5 = 243
            if (working_value >= 243) begin
                trit_value = 2'b10;
                working_value = working_value - 243;
            end else if (working_value <= -243) begin
                trit_value = 2'b00;
                working_value = working_value + 243;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[11:10] = trit_value;
            
            // Trit 4: 3^4 = 81
            if (working_value >= 81) begin
                trit_value = 2'b10;
                working_value = working_value - 81;
            end else if (working_value <= -81) begin
                trit_value = 2'b00;
                working_value = working_value + 81;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[9:8] = trit_value;
            
            // Trit 3: 3^3 = 27
            if (working_value >= 27) begin
                trit_value = 2'b10;
                working_value = working_value - 27;
            end else if (working_value <= -27) begin
                trit_value = 2'b00;
                working_value = working_value + 27;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[7:6] = trit_value;
            
            // Trit 2: 3^2 = 9
            if (working_value >= 9) begin
                trit_value = 2'b10;
                working_value = working_value - 9;
            end else if (working_value <= -9) begin
                trit_value = 2'b00;
                working_value = working_value + 9;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[5:4] = trit_value;
            
            // Trit 1: 3^1 = 3
            if (working_value >= 3) begin
                trit_value = 2'b10;
                working_value = working_value - 3;
            end else if (working_value <= -3) begin
                trit_value = 2'b00;
                working_value = working_value + 3;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[3:2] = trit_value;
            
            // Trit 0: 3^0 = 1
            if (working_value >= 1) begin
                trit_value = 2'b10;
                working_value = working_value - 1;
            end else if (working_value <= -1) begin
                trit_value = 2'b00;
                working_value = working_value + 1;
            end else begin
                trit_value = 2'b01;
            end
            result_bits[1:0] = trit_value;
              ternary_out = result_bits;
        end
    endtask
    
    // Input validation task
    task validate_ternary_input;
        input [35:0] ternary_value;
        output reg validation_error;
        integer i;
        reg [1:0] current_trit;
        begin
            validation_error = 1'b0;
            
            // Check each trit for valid encoding (00=-1, 01=0, 10=+1, 11=invalid)
            for (i = 0; i < 18; i = i + 1) begin
                current_trit = ternary_value[i*2 +: 2];
                if (current_trit == 2'b11) begin
                    validation_error = 1'b1;
                end
            end
        end
    endtask
    
    // Range checking task
    task check_result_range;
        input signed [31:0] decimal_value;
        output reg range_error;
        output reg overflow_flag;
        output reg underflow_flag;
        begin
            range_error = 1'b0;
            overflow_flag = 1'b0;
            underflow_flag = 1'b0;
            
            if (decimal_value > 32'sd129140163) begin  // Max ternary value
                range_error = 1'b1;
                overflow_flag = 1'b1;
            end else if (decimal_value < -32'sd129140163) begin  // Min ternary value
                range_error = 1'b1;
                underflow_flag = 1'b1;
            end
        end
    endtask
      // State machine and error handling
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            current_state <= STATE_IDLE;
            result <= `VTX1_WORD_DEFAULT;
            operation_complete <= 1'b0;
            arithmetic_overflow <= 1'b0;
            arithmetic_underflow <= 1'b0;
            error <= 1'b0;
            error_code <= `VTX1_ERROR_NONE;
            timeout <= 1'b0;
            state <= STATE_IDLE;
            operation_count <= 32'h0;
            error_count <= 32'h0;
            decimal_a <= 32'sd0;
            decimal_b <= 32'sd0;
            decimal_result <= 32'sd0;
            operation_valid <= 1'b0;
            timeout_counter <= 32'h0;
            input_validation_error <= 1'b0;
            range_check_error <= 1'b0;
            max_operation_time <= 32'h0;
            min_operation_time <= 32'hFFFFFFFF;
            avg_operation_time <= 32'h0;
            performance_counter <= 32'h0;
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
            case (current_state)
                STATE_IDLE: begin
                    operation_complete <= 1'b0;
                    timeout_counter <= 32'h0;
                    
                    if (operation_enable) begin
                        current_state <= STATE_DECODE;
                        operation_count <= operation_count + 1;
                    end
                end
                
                STATE_DECODE: begin
                    // Validate operation code
                    if (operation > OP_DEC) begin
                        error <= 1'b1;
                        error_code <= `VTX1_ERROR_INVALID_OP;
                        error_count <= error_count + 1;
                        current_state <= STATE_ERROR;
                    end else begin
                        current_state <= STATE_CONVERT;
                    end
                end
                
                STATE_CONVERT: begin
                    // Validate input operands
                    validate_ternary_input(operand_a, input_validation_error);
                    if (input_validation_error) begin
                        error <= 1'b1;
                        error_code <= `VTX1_ERROR_INVALID_ADDR; // Reusing for invalid data
                        error_count <= error_count + 1;
                        current_state <= STATE_ERROR;
                    end else begin
                        validate_ternary_input(operand_b, input_validation_error);
                        if (input_validation_error && (operation != OP_NEG && operation != OP_ABS && operation != OP_INC && operation != OP_DEC)) begin
                            error <= 1'b1;
                            error_code <= `VTX1_ERROR_INVALID_ADDR; // Reusing for invalid data
                            error_count <= error_count + 1;
                            current_state <= STATE_ERROR;
                        end else begin
                            // Convert operands
                            ternary_to_decimal_task(operand_a, decimal_a);
                            if (operation != OP_NEG && operation != OP_ABS && operation != OP_INC && operation != OP_DEC) begin
                                ternary_to_decimal_task(operand_b, decimal_b);
                            end
                            current_state <= STATE_EXECUTE;
                        end
                    end
                end
                
                STATE_EXECUTE: begin
                    // Timeout detection
                    timeout_counter <= timeout_counter + 1;
                    if (timeout_counter >= `VTX1_TIMEOUT_CYCLES) begin
                        error <= 1'b1;
                        error_code <= `VTX1_ERROR_TIMEOUT;
                        timeout <= 1'b1;
                        error_count <= error_count + 1;
                        current_state <= STATE_ERROR;
                    end else begin
                        // Perform arithmetic operation with comprehensive error checking
                        case (operation)
                            OP_ADD: begin
                                decimal_result = decimal_a + decimal_b;
                                operation_valid = 1'b1;
                            end
                            
                            OP_SUB: begin
                                decimal_result = decimal_a - decimal_b;
                                operation_valid = 1'b1;
                            end
                            
                            OP_MUL: begin
                                decimal_result = decimal_a * decimal_b;
                                operation_valid = 1'b1;
                            end
                            
                            OP_DIV: begin
                                if (decimal_b != 32'sd0) begin
                                    decimal_result = decimal_a / decimal_b;
                                    operation_valid = 1'b1;
                                end else begin
                                    error <= 1'b1;
                                    error_code <= `VTX1_ERROR_INVALID_OP;
                                    error_count <= error_count + 1;
                                    current_state <= STATE_ERROR;
                                    operation_valid = 1'b0;
                                end
                            end
                            
                            OP_MOD: begin
                                if (decimal_b != 32'sd0) begin
                                    decimal_result = decimal_a % decimal_b;
                                    operation_valid = 1'b1;
                                end else begin
                                    error <= 1'b1;
                                    error_code <= `VTX1_ERROR_INVALID_OP;
                                    error_count <= error_count + 1;
                                    current_state <= STATE_ERROR;
                                    operation_valid = 1'b0;
                                end
                            end
                            
                            OP_NEG: begin
                                decimal_result = -decimal_a;
                                operation_valid = 1'b1;
                            end
                            
                            OP_ABS: begin
                                decimal_result = (decimal_a < 32'sd0) ? -decimal_a : decimal_a;
                                operation_valid = 1'b1;
                            end
                            
                            OP_CMP: begin
                                if (decimal_a > decimal_b)
                                    decimal_result = 32'sd1;
                                else if (decimal_a < decimal_b)
                                    decimal_result = -32'sd1;
                                else
                                    decimal_result = 32'sd0;
                                operation_valid = 1'b1;
                            end
                            
                            OP_MAX: begin
                                decimal_result = (decimal_a > decimal_b) ? decimal_a : decimal_b;
                                operation_valid = 1'b1;
                            end
                            
                            OP_MIN: begin
                                decimal_result = (decimal_a < decimal_b) ? decimal_a : decimal_b;
                                operation_valid = 1'b1;
                            end
                            
                            OP_INC: begin
                                decimal_result = decimal_a + 32'sd1;
                                operation_valid = 1'b1;
                            end
                            
                            OP_DEC: begin
                                decimal_result = decimal_a - 32'sd1;
                                operation_valid = 1'b1;
                            end
                            
                            default: begin
                                error <= 1'b1;
                                error_code <= `VTX1_ERROR_INVALID_OP;
                                error_count <= error_count + 1;
                                current_state <= STATE_ERROR;
                                operation_valid = 1'b0;
                            end
                        endcase
                        
                        if (operation_valid) begin
                            current_state <= STATE_RESULT;
                        end
                    end
                end
                
                STATE_RESULT: begin
                    // Check for overflow/underflow
                    check_result_range(decimal_result, range_check_error, arithmetic_overflow, arithmetic_underflow);
                    
                    if (range_check_error) begin
                        if (arithmetic_overflow) begin
                            error <= 1'b1;
                            error_code <= `VTX1_ERROR_OVERFLOW;
                            error_count <= error_count + 1;
                            result <= 36'h222222222; // Max positive ternary value
                        end else if (arithmetic_underflow) begin
                            error <= 1'b1;
                            error_code <= `VTX1_ERROR_UNDERFLOW;
                            error_count <= error_count + 1;
                            result <= 36'h000000000; // Max negative ternary value
                        end
                    end else begin
                        // Convert result back to ternary
                        decimal_to_ternary_task(decimal_result, result);
                    end
                    
                    operation_complete <= 1'b1;
                    current_state <= STATE_IDLE;
                end
                
                STATE_ERROR: begin
                    operation_complete <= 1'b0;
                    result <= `VTX1_WORD_DEFAULT;
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

endmodule

// ============================================================================
// FAST TERNARY ARITHMETIC FUNCTIONS
// ============================================================================
// Optimized combinational functions for simple ternary arithmetic operations

(* keep_hierarchy = "yes" *)
(* noopt *)
module fast_ternary_adder (
    input wire [35:0] a,
    input wire [35:0] b,
    output reg [35:0] sum,
    output reg carry_out,
    output reg valid
);// Internal signals for trit-wise addition
    reg [1:0] trit_a, trit_b, trit_sum;
    reg [1:0] carry [0:18]; // Extended to include carry[18]
    integer i;
    
    always @(*) begin
        valid = 1'b1;
        carry_out = 1'b0;
        carry[0] = 2'b01; // Initial carry is zero
        
        // Perform trit-wise addition
        for (i = 0; i < 18; i = i + 1) begin
            trit_a = a[i*2 +: 2];
            trit_b = b[i*2 +: 2];
            
            // Check for invalid trits
            if (trit_a == 2'b11 || trit_b == 2'b11) begin
                valid = 1'b0;
                sum[i*2 +: 2] = 2'b01; // Default to zero
                carry[i+1] = 2'b01;
            end else begin
                // Ternary addition with carry
                case ({trit_a, trit_b, carry[i]})
                    // {a, b, carry_in} -> {sum, carry_out}
                    6'b000001: begin sum[i*2 +: 2] = 2'b00; carry[i+1] = 2'b01; end // -1 + -1 + 0 = -2 -> -1, carry -1
                    6'b000101: begin sum[i*2 +: 2] = 2'b01; carry[i+1] = 2'b00; end // -1 + 0 + 0 = -1
                    6'b001001: begin sum[i*2 +: 2] = 2'b01; carry[i+1] = 2'b01; end // -1 + +1 + 0 = 0
                    6'b010001: begin sum[i*2 +: 2] = 2'b00; carry[i+1] = 2'b01; end // 0 + -1 + 0 = -1
                    6'b010101: begin sum[i*2 +: 2] = 2'b01; carry[i+1] = 2'b01; end // 0 + 0 + 0 = 0
                    6'b011001: begin sum[i*2 +: 2] = 2'b10; carry[i+1] = 2'b01; end // 0 + +1 + 0 = +1
                    6'b100001: begin sum[i*2 +: 2] = 2'b01; carry[i+1] = 2'b01; end // +1 + -1 + 0 = 0
                    6'b100101: begin sum[i*2 +: 2] = 2'b10; carry[i+1] = 2'b01; end // +1 + 0 + 0 = +1
                    6'b101001: begin sum[i*2 +: 2] = 2'b01; carry[i+1] = 2'b10; end // +1 + +1 + 0 = +2 -> 0, carry +1
                    default: begin sum[i*2 +: 2] = 2'b01; carry[i+1] = 2'b01; end
                endcase
            end
        end
        
        // Final carry indicates overflow
        if (carry[18] != 2'b01) begin
            carry_out = 1'b1;
        end
    end

endmodule

// ============================================================================
// TERNARY MAGNITUDE COMPARATOR
// ============================================================================

(* keep_hierarchy = "yes" *)
(* noopt *)
module ternary_comparator (
    input wire [35:0] a,
    input wire [35:0] b,
    output reg a_greater,
    output reg a_equal,
    output reg a_less,
    output reg valid
);

    reg [1:0] trit_a, trit_b;
    reg comparison_found;
    integer i;
    
    always @(*) begin
        valid = 1'b1;
        a_greater = 1'b0;
        a_equal = 1'b1;
        a_less = 1'b0;
        comparison_found = 1'b0;
        // Compare from most significant trit to least significant
        for (i = 17; i >= 0; i = i - 1) begin
            if (!comparison_found) begin
                trit_a = a[i*2 +: 2];
                trit_b = b[i*2 +: 2];
                // Check for invalid trits
                if (trit_a == 2'b11 || trit_b == 2'b11) begin
                    valid = 1'b0;
                    a_equal = 1'b0;
                    comparison_found = 1'b1;
                end else begin
                    // Compare trits: NEG=00 < ZERO=01 < POS=10
                    if (trit_a > trit_b) begin
                        a_greater = 1'b1;
                        a_equal = 1'b0;
                        comparison_found = 1'b1;
                    end else if (trit_a < trit_b) begin
                        a_less = 1'b1;
                        a_equal = 1'b0;
                        comparison_found = 1'b1;
                    end
                    // If equal, continue to next trit
                end
            end
        end
    end

endmodule

// ============================================================================
// TERNARY MULTIPLIER (SIMPLIFIED)
// ============================================================================

(* keep_hierarchy = "yes" *)
(* noopt *)
module ternary_multiplier (
    input wire clk,
    input wire rst_n,
    input wire [17:0] multiplicand,  // 9-trit multiplicand
    input wire [17:0] multiplier,    // 9-trit multiplier  
    input wire start,
    output reg [35:0] product,       // 18-trit product
    output reg complete,
    output reg valid
);

    // Internal state machine
    localparam IDLE = 2'b00;
    localparam MULTIPLY = 2'b01;
    localparam DONE = 2'b10;
    
    reg [1:0] state;
    reg [4:0] bit_counter;
    reg [35:0] accumulator;
    reg [17:0] shift_multiplicand;
    reg [1:0] current_multiplier_trit;
    
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            state <= IDLE;
            product <= 36'h0;
            complete <= 1'b0;
            valid <= 1'b1;
            bit_counter <= 5'h0;
            accumulator <= 36'h0;
            shift_multiplicand <= 18'h0;
        end else begin
            case (state)
                IDLE: begin
                    complete <= 1'b0;
                    if (start) begin
                        accumulator <= 36'h015555555; // All zeros
                        shift_multiplicand <= multiplicand;
                        bit_counter <= 5'h0;
                        valid <= 1'b1;
                        state <= MULTIPLY;
                    end
                end
                
                MULTIPLY: begin
                    if (bit_counter < 9) begin
                        current_multiplier_trit <= multiplier[bit_counter*2 +: 2];
                        
                        // Add or subtract based on multiplier trit
                        case (multiplier[bit_counter*2 +: 2])
                            2'b10: begin // POS: add
                                // accumulator <= accumulator + shift_multiplicand;
                                // Simplified for this example
                            end
                            2'b00: begin // NEG: subtract
                                // accumulator <= accumulator - shift_multiplicand;
                                // Simplified for this example
                            end
                            2'b01: begin // ZERO: do nothing
                                // No operation needed
                            end
                            2'b11: begin // Invalid
                                valid <= 1'b0;
                            end
                        endcase
                        
                        // Shift multiplicand left by one trit (multiply by 3)
                        shift_multiplicand <= {shift_multiplicand[15:0], 2'b01};
                        bit_counter <= bit_counter + 1;
                    end else begin
                        product <= accumulator;
                        state <= DONE;
                    end
                end
                
                DONE: begin
                    complete <= 1'b1;
                    if (!start) begin
                        state <= IDLE;
                    end
                end
            endcase
        end
    end

endmodule

