`timescale 1ns / 1ps
// VTX1 Ternary Register File
// Part of the VTX1 Ternary System-on-Chip

`ifndef REGISTER_FILE_V
`define REGISTER_FILE_V

// Include VTX1 interface definitions
// Note: Include paths are handled by compiler -I flags in Taskfile.yml
`include "vtx1_interfaces.v"
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

// ============================================================================
// VTX1 TERNARY REGISTER FILE
// ============================================================================
// Implements the 13-register ternary register file with proper read/write
// ports and register addressing for the VTX1 architecture.

module register_file (
    input  wire                     clk,
    input  wire                     rst_n,
      // Read ports (3 ports for VLIW support)
    input  wire [`VTX1_REG_ADDR_WIDTH-1:0] read_addr_a,
    input  wire [`VTX1_REG_ADDR_WIDTH-1:0] read_addr_b,
    input  wire [`VTX1_REG_ADDR_WIDTH-1:0] read_addr_c,
    output reg  [`VTX1_WORD_WIDTH-1:0]     read_data_a,
    output reg  [`VTX1_WORD_WIDTH-1:0]     read_data_b,
    output reg  [`VTX1_WORD_WIDTH-1:0]     read_data_c,
    
    // Write ports (2 ports for VLIW support)
    input  wire                        write_enable_a,
    input  wire                        write_enable_b,
    input  wire [`VTX1_REG_ADDR_WIDTH-1:0] write_addr_a,
    input  wire [`VTX1_REG_ADDR_WIDTH-1:0] write_addr_b,
    input  wire [`VTX1_WORD_WIDTH-1:0]     write_data_a,
    input  wire [`VTX1_WORD_WIDTH-1:0]     write_data_b,
    
    // Debug interface
    input  wire                        debug_enable,
    input  wire [`VTX1_REG_ADDR_WIDTH-1:0] debug_addr,
    output reg  [`VTX1_WORD_WIDTH-1:0]     debug_data,
    
    // Error handling interface
    output reg                        error,
    output reg  [3:0]                 error_code,
    output reg                        timeout,
    output reg  [31:0]                error_count,
    
    // Status outputs
    output reg                        write_conflict,
    output reg                        invalid_address,
    
    // Debug outputs
    output reg  [3:0]                 current_state,
    output reg  [31:0]                operation_count,
    output reg  [31:0]                read_count,
    output reg  [31:0]                write_count
);    // ========================================================================
    // STATE MACHINE AND ERROR HANDLING - Use VTX1 standardized states
    // ========================================================================
    
    // State definitions using VTX1 constants
    localparam STATE_IDLE     = `VTX1_STATE_IDLE;
    localparam STATE_READ     = `VTX1_STATE_ACTIVE;  // Map read to active
    localparam STATE_WRITE    = `VTX1_STATE_EXECUTE; // Map write to execute
    localparam STATE_DEBUG    = `VTX1_STATE_DEBUG;
    localparam STATE_ERROR    = `VTX1_STATE_ERROR;
    
    // Internal state and error tracking
    reg [3:0] next_state;
    reg [31:0] timeout_counter;
    reg read_error, write_error, debug_error;
    
    // VTX1 Error Handling Variables
    reg [3:0] vtx1_error_reg;
    reg [31:0] vtx1_error_info;
    reg vtx1_error_valid;
    
    // ========================================================================
    // REGISTER DEFINITIONS
    // ========================================================================
    // VTX1 Register mapping:
    // 0-6:  T0-T6    - General purpose ternary registers
    // 7:    TA       - Accumulator register
    // 8:    TB       - Base register  
    // 9:    TC       - Counter register
    // 10:   TS       - Status register
    // 11:   TI       - Index register
    // 12:   VA       - Vector accumulator (reserved for future)
    
    // Note: Additional registers VT, VB, FA, FT, FB are planned for future
    // vector and floating-point extensions but not implemented in this version
    
    localparam REG_T0 = 4'd0;   // General purpose register 0
    localparam REG_T1 = 4'd1;   // General purpose register 1
    localparam REG_T2 = 4'd2;   // General purpose register 2
    localparam REG_T3 = 4'd3;   // General purpose register 3
    localparam REG_T4 = 4'd4;   // General purpose register 4
    localparam REG_T5 = 4'd5;   // General purpose register 5
    localparam REG_T6 = 4'd6;   // General purpose register 6
    localparam REG_TA = 4'd7;   // Accumulator register
    localparam REG_TB = 4'd8;   // Base register
    localparam REG_TC = 4'd9;   // Counter register
    localparam REG_TS = 4'd10;  // Status register
    localparam REG_TI = 4'd11;  // Index register
    localparam REG_VA = 4'd12;  // Vector accumulator (future)
    
    // Physical register storage
    reg [`VTX1_WORD_WIDTH-1:0] registers [0:`VTX1_NUM_REGISTERS-1];
    integer i;
    
    // Address validation tasks
    task validate_read_address;
        input [3:0] addr;
        output reg addr_error;
        begin
            addr_error = (addr >= `VTX1_NUM_REGISTERS);
        end
    endtask
    
    task validate_write_address;
        input [3:0] addr;
        output reg addr_error;
        begin
            addr_error = (addr >= `VTX1_NUM_REGISTERS);
        end
    endtask
    
    // Write conflict detection
    wire addr_conflict = (write_enable_a && write_enable_b && 
                         (write_addr_a == write_addr_b));
    
    // ========================================================================
    // MAIN STATE MACHINE
    // ========================================================================
      always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            // Initialize all registers to default ternary zero
            for (i = 0; i < `VTX1_NUM_REGISTERS; i = i + 1) begin
                registers[i] <= `VTX1_WORD_DEFAULT;
            end
            
            // Initialize state and counters
            current_state <= STATE_IDLE;
            timeout_counter <= 0;
            operation_count <= 0;
            read_count <= 0;
            write_count <= 0;
            error_count <= 0;
              // Initialize VTX1 error handling
            `VTX1_CLEAR_ERROR(vtx1_error_reg, current_state)
            
            // Initialize outputs
            error <= 1'b0;
            error_code <= `VTX1_ERROR_NONE;
            timeout <= 1'b0;
            write_conflict <= 1'b0;
            invalid_address <= 1'b0;
            
            read_data_a <= `VTX1_WORD_DEFAULT;
            read_data_b <= `VTX1_WORD_DEFAULT;
            read_data_c <= `VTX1_WORD_DEFAULT;
            debug_data <= `VTX1_WORD_DEFAULT;
        end else begin            // Timeout detection with VTX1 error handling
            if (timeout_counter >= `VTX1_TIMEOUT_CYCLES) begin
                `VTX1_SET_ERROR(vtx1_error_reg, current_state, `VTX1_ERROR_TIMEOUT)
                vtx1_error_info <= 32'h0000_0001;
                vtx1_error_valid <= 1'b1;
                error <= 1'b1;
                error_code <= `VTX1_ERROR_TIMEOUT;
                timeout <= 1'b1;
                current_state <= STATE_ERROR;
                error_count <= error_count + 1;
            end else begin
                timeout_counter <= timeout_counter + 1;
                
                case (current_state)
                    STATE_IDLE: begin                        timeout_counter <= 0;
                        `VTX1_CLEAR_ERROR(vtx1_error_reg, current_state)
                        error <= 1'b0;
                        timeout <= 1'b0;
                        
                        // Determine next operation
                        if (debug_enable) begin
                            current_state <= STATE_DEBUG;
                        end else if (write_enable_a || write_enable_b) begin
                            current_state <= STATE_WRITE;
                        end else begin
                            current_state <= STATE_READ;
                        end
                        
                        operation_count <= operation_count + 1;                    end
                    
                    STATE_READ: begin                        // Validate read addresses
                        validate_read_address(read_addr_a, read_error);
                        if (read_error) begin
                            `VTX1_SET_ERROR(vtx1_error_reg, current_state, `VTX1_ERROR_INVALID_ADDR)
                            vtx1_error_info <= {16'h0000, read_addr_a, 12'h000};
                            vtx1_error_valid <= 1'b1;
                            error <= 1'b1;
                            error_code <= `VTX1_ERROR_INVALID_ADDR;
                            invalid_address <= 1'b1;
                            current_state <= STATE_ERROR;
                            error_count <= error_count + 1;
                        end else begin
                            // Perform reads
                            read_data_a <= registers[read_addr_a];
                            read_data_b <= registers[read_addr_b];
                            read_data_c <= registers[read_addr_c];
                            read_count <= read_count + 1;
                            current_state <= STATE_IDLE;
                        end
                    end
                    
                    STATE_WRITE: begin                        // Check for write conflicts
                        if (addr_conflict) begin
                            `VTX1_SET_ERROR(vtx1_error_reg, current_state, `VTX1_ERROR_COLLISION)
                            vtx1_error_info <= {16'h0000, write_addr_a, write_addr_b, 8'h00};
                            vtx1_error_valid <= 1'b1;
                            error <= 1'b1;
                            error_code <= `VTX1_ERROR_COLLISION;
                            write_conflict <= 1'b1;
                            current_state <= STATE_ERROR;
                            error_count <= error_count + 1;
                        end else begin
                            // Validate write addresses and perform writes
                            if (write_enable_a) begin                                validate_write_address(write_addr_a, write_error);
                                if (write_error) begin
                                    `VTX1_SET_ERROR(vtx1_error_reg, current_state, `VTX1_ERROR_INVALID_ADDR)
                                    vtx1_error_info <= {16'h0000, write_addr_a, 12'h000};
                                    vtx1_error_valid <= 1'b1;
                                    error <= 1'b1;
                                    error_code <= `VTX1_ERROR_INVALID_ADDR;
                                    invalid_address <= 1'b1;
                                    current_state <= STATE_ERROR;
                                    error_count <= error_count + 1;
                                end else begin
                                    registers[write_addr_a] <= write_data_a;
                                    write_count <= write_count + 1;
                                end
                            end
                            
                            if (write_enable_b && !write_error) begin                                validate_write_address(write_addr_b, write_error);
                                if (write_error) begin
                                    `VTX1_SET_ERROR(vtx1_error_reg, current_state, `VTX1_ERROR_INVALID_ADDR)
                                    vtx1_error_info <= {16'h0000, write_addr_b, 12'h000};
                                    vtx1_error_valid <= 1'b1;
                                    error <= 1'b1;
                                    error_code <= `VTX1_ERROR_INVALID_ADDR;
                                    invalid_address <= 1'b1;
                                    current_state <= STATE_ERROR;
                                    error_count <= error_count + 1;
                                end else begin
                                    registers[write_addr_b] <= write_data_b;
                                    write_count <= write_count + 1;
                                end
                            end
                            
                            if (!write_error) begin
                                current_state <= STATE_IDLE;
                            end
                        end                    end
                    
                    STATE_DEBUG: begin                        validate_read_address(debug_addr, debug_error);
                        if (debug_error) begin
                            `VTX1_SET_ERROR(vtx1_error_reg, current_state, `VTX1_ERROR_INVALID_ADDR)
                            vtx1_error_info <= {16'h0000, debug_addr, 12'h000};
                            vtx1_error_valid <= 1'b1;
                            error <= 1'b1;
                            error_code <= `VTX1_ERROR_INVALID_ADDR;
                            invalid_address <= 1'b1;
                            current_state <= STATE_ERROR;
                            error_count <= error_count + 1;
                        end else begin
                            debug_data <= registers[debug_addr];
                            current_state <= STATE_IDLE;
                        end
                    end
                    
                    STATE_ERROR: begin
                        // Error recovery - stay in error state until reset or explicit clear
                        timeout_counter <= 0;
                        // Error state persists until reset
                    end
                      default: begin
                        `VTX1_SET_ERROR(vtx1_error_reg, current_state, `VTX1_ERROR_INVALID_OP)
                        vtx1_error_info <= 32'h0000_0001;
                        vtx1_error_valid <= 1'b1;
                        error <= 1'b1;
                        error_code <= `VTX1_ERROR_INVALID_OP;
                        current_state <= STATE_ERROR;
                        error_count <= error_count + 1;
                    end
                endcase
            end
        end
    end

endmodule

// ============================================================================
// REGISTER FILE CONTROLLER
// ============================================================================
// Provides higher-level control for register operations

module register_file_controller (
    input  wire                     clk,
    input  wire                     rst_n,
    
    // Pipeline interface
    input  wire                     decode_valid,
    input  wire                     execute_valid,
    input  wire                     writeback_valid,
      // Instruction fields
    input  wire [`VTX1_REG_ADDR_WIDTH-1:0] inst_rs1,    // Source register 1
    input  wire [`VTX1_REG_ADDR_WIDTH-1:0] inst_rs2,    // Source register 2
    input  wire [`VTX1_REG_ADDR_WIDTH-1:0] inst_rs3,    // Source register 3
    input  wire [`VTX1_REG_ADDR_WIDTH-1:0] inst_rd1,    // Destination register 1
    input  wire [`VTX1_REG_ADDR_WIDTH-1:0] inst_rd2,    // Destination register 2
    
    // Register file interface
    output reg  [`VTX1_REG_ADDR_WIDTH-1:0] rf_read_addr_a,
    output reg  [`VTX1_REG_ADDR_WIDTH-1:0] rf_read_addr_b,
    output reg  [`VTX1_REG_ADDR_WIDTH-1:0] rf_read_addr_c,
    input  wire [`VTX1_WORD_WIDTH-1:0]     rf_read_data_a,
    input  wire [`VTX1_WORD_WIDTH-1:0]     rf_read_data_b,
    input  wire [`VTX1_WORD_WIDTH-1:0]     rf_read_data_c,
    
    output reg                        rf_write_enable_a,
    output reg                        rf_write_enable_b,
    output reg  [`VTX1_REG_ADDR_WIDTH-1:0] rf_write_addr_a,
    output reg  [`VTX1_REG_ADDR_WIDTH-1:0] rf_write_addr_b,
    output reg  [`VTX1_WORD_WIDTH-1:0]     rf_write_data_a,
    output reg  [`VTX1_WORD_WIDTH-1:0]     rf_write_data_b,
    
    // Execute stage inputs
    input  wire [`VTX1_WORD_WIDTH-1:0]     execute_result_a,
    input  wire [`VTX1_WORD_WIDTH-1:0]     execute_result_b,
    input  wire                       execute_write_enable_a,
    input  wire                       execute_write_enable_b,
    
    // Hazard detection
    output reg                        read_after_write_hazard,
    output reg                        write_after_write_hazard,
    
    // Status outputs
    input  wire                       rf_write_conflict,
    input  wire                       rf_invalid_address,    output reg                        register_error
);

    // Register address constants (for controller use)
    localparam REG_T0 = 4'd0;   // General purpose register 0

    // ========================================================================
    // READ ADDRESS GENERATION
    // ========================================================================
    
    always @(*) begin
        if (decode_valid) begin
            rf_read_addr_a = inst_rs1;
            rf_read_addr_b = inst_rs2;
            rf_read_addr_c = inst_rs3;
        end else begin
            rf_read_addr_a = REG_T0;  // Default to T0
            rf_read_addr_b = REG_T0;
            rf_read_addr_c = REG_T0;
        end
    end
    
    // ========================================================================
    // WRITE CONTROL
    // ========================================================================
    
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            rf_write_enable_a <= 1'b0;
            rf_write_enable_b <= 1'b0;
            rf_write_addr_a <= REG_T0;
            rf_write_addr_b <= REG_T0;            rf_write_data_a <= `VTX1_WORD_DEFAULT;
            rf_write_data_b <= `VTX1_WORD_DEFAULT;
            register_error <= 1'b0;
        end else begin
            // Update error status
            register_error <= rf_write_conflict || rf_invalid_address;
            
            // Handle writeback stage
            if (writeback_valid) begin
                rf_write_enable_a <= execute_write_enable_a;
                rf_write_enable_b <= execute_write_enable_b;
                rf_write_addr_a <= inst_rd1;
                rf_write_addr_b <= inst_rd2;
                rf_write_data_a <= execute_result_a;
                rf_write_data_b <= execute_result_b;
            end else begin
                rf_write_enable_a <= 1'b0;
                rf_write_enable_b <= 1'b0;
            end
        end
    end
    
    // ========================================================================
    // HAZARD DETECTION
    // ========================================================================
    
    always @(*) begin
        // RAW (Read After Write) hazard detection
        read_after_write_hazard = 1'b0;
        
        if (decode_valid && execute_valid) begin
            // Check if current instruction reads a register that the
            // previous instruction is writing to
            if (execute_write_enable_a && 
                ((inst_rs1 == inst_rd1) || (inst_rs2 == inst_rd1) || (inst_rs3 == inst_rd1))) begin
                read_after_write_hazard = 1'b1;
            end
            
            if (execute_write_enable_b && 
                ((inst_rs1 == inst_rd2) || (inst_rs2 == inst_rd2) || (inst_rs3 == inst_rd2))) begin
                read_after_write_hazard = 1'b1;
            end
        end
        
        // WAW (Write After Write) hazard detection
        write_after_write_hazard = 1'b0;
        
        if (decode_valid && execute_valid) begin
            if (execute_write_enable_a && execute_write_enable_b &&
                (inst_rd1 == inst_rd2)) begin
                write_after_write_hazard = 1'b1;
            end
        end
    end

endmodule

`endif // REGISTER_FILE_V

