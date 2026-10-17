	`timescale 1ns / 1ps
// ============================================================================
// VTX1 Forwarding Unit
// ============================================================================
// Data forwarding logic for VTX1 4-stage VLIW pipeline
// Implements bypass paths to resolve RAW hazards without stalling
// ============================================================================

`ifndef FORWARDING_UNIT_V
`define FORWARDING_UNIT_V

`include "vtx1_interfaces.v"

module forwarding_unit (
    input wire clk,
    input wire rst_n,
    
    // Pipeline stage control
    input wire decode_valid,
    input wire execute_valid,
    input wire writeback_valid,
    
    // Decode stage register addresses
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] decode_rs1,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] decode_rs2,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] decode_rs3,
    
    // Execute stage data and control
    input wire [`VTX1_WORD_WIDTH-1:0] execute_result_a,
    input wire [`VTX1_WORD_WIDTH-1:0] execute_result_b,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] execute_rd1,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] execute_rd2,
    input wire execute_wr_en1,
    input wire execute_wr_en2,
    input wire execute_mem_read,
    
    // Writeback stage data and control
    input wire [`VTX1_WORD_WIDTH-1:0] writeback_result_a,
    input wire [`VTX1_WORD_WIDTH-1:0] writeback_result_b,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] writeback_rd1,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] writeback_rd2,
    input wire writeback_wr_en1,
    input wire writeback_wr_en2,
    
    // Memory stage data (for load forwarding)
    input wire [`VTX1_WORD_WIDTH-1:0] memory_result,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] memory_rd,
    input wire memory_wr_en,
    input wire memory_valid,
    
    // Register file data (default path)
    input wire [`VTX1_WORD_WIDTH-1:0] rf_data_a,
    input wire [`VTX1_WORD_WIDTH-1:0] rf_data_b,
    input wire [`VTX1_WORD_WIDTH-1:0] rf_data_c,
    
    // Forwarded data outputs
    output reg [`VTX1_WORD_WIDTH-1:0] forward_data_a,
    output reg [`VTX1_WORD_WIDTH-1:0] forward_data_b,
    output reg [`VTX1_WORD_WIDTH-1:0] forward_data_c,
    
    // Forwarding control signals
    output reg [2:0] forward_sel_a,    // 000=RF, 001=EXE_A, 010=EXE_B, 011=WB_A, 100=WB_B, 101=MEM
    output reg [2:0] forward_sel_b,
    output reg [2:0] forward_sel_c,
    
    // Forwarding status
    output reg forward_valid_a,
    output reg forward_valid_b,
    output reg forward_valid_c,
    output reg load_forward_stall,     // Stall needed for load forwarding
    
    // Performance monitoring
    output reg [31:0] forwards_from_execute,
    output reg [31:0] forwards_from_writeback,
    output reg [31:0] forwards_from_memory,
    output reg [31:0] load_forward_stalls
);

// ============================================================================
// CONSTANTS AND PARAMETERS
// ============================================================================

// Forwarding source selection
localparam FORWARD_RF      = 3'b000;  // Register file (no forwarding)
localparam FORWARD_EXE_A   = 3'b001;  // Execute stage result A
localparam FORWARD_EXE_B   = 3'b010;  // Execute stage result B
localparam FORWARD_WB_A    = 3'b011;  // Writeback stage result A
localparam FORWARD_WB_B    = 3'b100;  // Writeback stage result B
localparam FORWARD_MEM     = 3'b101;  // Memory stage result
localparam FORWARD_STALL   = 3'b110;  // Forwarding not possible, stall needed
localparam FORWARD_INVALID = 3'b111;  // Invalid forwarding condition

// ============================================================================
// INTERNAL SIGNALS
// ============================================================================

reg exe_to_rs1, exe_to_rs2, exe_to_rs3;
reg wb_to_rs1, wb_to_rs2, wb_to_rs3;
reg mem_to_rs1, mem_to_rs2, mem_to_rs3;
reg load_stall_rs1, load_stall_rs2, load_stall_rs3;

// ============================================================================
// RS1 FORWARDING LOGIC
// ============================================================================

always @(*) begin
    forward_sel_a = FORWARD_RF;
    forward_data_a = rf_data_a;
    forward_valid_a = 1'b1;
    exe_to_rs1 = 1'b0;
    wb_to_rs1 = 1'b0;
    mem_to_rs1 = 1'b0;
    load_stall_rs1 = 1'b0;
    
    if (decode_valid && (decode_rs1 != 5'h00)) begin
        // Priority 1: Memory stage forwarding (for loads)
        if (memory_valid && memory_wr_en && (decode_rs1 == memory_rd)) begin
            forward_sel_a = FORWARD_MEM;
            forward_data_a = memory_result;
            mem_to_rs1 = 1'b1;
        end
        
        // Priority 2: Execute stage forwarding
        else if (execute_valid && execute_wr_en1 && (decode_rs1 == execute_rd1)) begin
            if (execute_mem_read) begin
                // Load instruction in execute - need to stall
                forward_sel_a = FORWARD_STALL;
                load_stall_rs1 = 1'b1;
                forward_valid_a = 1'b0;
            end else begin
                forward_sel_a = FORWARD_EXE_A;
                forward_data_a = execute_result_a;
                exe_to_rs1 = 1'b1;
            end
        end
        else if (execute_valid && execute_wr_en2 && (decode_rs1 == execute_rd2)) begin
            if (execute_mem_read) begin
                // Load instruction in execute - need to stall
                forward_sel_a = FORWARD_STALL;
                load_stall_rs1 = 1'b1;
                forward_valid_a = 1'b0;
            end else begin
                forward_sel_a = FORWARD_EXE_B;
                forward_data_a = execute_result_b;
                exe_to_rs1 = 1'b1;
            end
        end
        
        // Priority 3: Writeback stage forwarding
        else if (writeback_valid && writeback_wr_en1 && (decode_rs1 == writeback_rd1)) begin
            forward_sel_a = FORWARD_WB_A;
            forward_data_a = writeback_result_a;
            wb_to_rs1 = 1'b1;
        end
        else if (writeback_valid && writeback_wr_en2 && (decode_rs1 == writeback_rd2)) begin
            forward_sel_a = FORWARD_WB_B;
            forward_data_a = writeback_result_b;
            wb_to_rs1 = 1'b1;
        end
    end
end

// ============================================================================
// RS2 FORWARDING LOGIC
// ============================================================================

always @(*) begin
    forward_sel_b = FORWARD_RF;
    forward_data_b = rf_data_b;
    forward_valid_b = 1'b1;
    exe_to_rs2 = 1'b0;
    wb_to_rs2 = 1'b0;
    mem_to_rs2 = 1'b0;
    load_stall_rs2 = 1'b0;
    
    if (decode_valid && (decode_rs2 != 5'h00)) begin
        // Priority 1: Memory stage forwarding (for loads)
        if (memory_valid && memory_wr_en && (decode_rs2 == memory_rd)) begin
            forward_sel_b = FORWARD_MEM;
            forward_data_b = memory_result;
            mem_to_rs2 = 1'b1;
        end
        
        // Priority 2: Execute stage forwarding
        else if (execute_valid && execute_wr_en1 && (decode_rs2 == execute_rd1)) begin
            if (execute_mem_read) begin
                // Load instruction in execute - need to stall
                forward_sel_b = FORWARD_STALL;
                load_stall_rs2 = 1'b1;
                forward_valid_b = 1'b0;
            end else begin
                forward_sel_b = FORWARD_EXE_A;
                forward_data_b = execute_result_a;
                exe_to_rs2 = 1'b1;
            end
        end
        else if (execute_valid && execute_wr_en2 && (decode_rs2 == execute_rd2)) begin
            if (execute_mem_read) begin
                // Load instruction in execute - need to stall
                forward_sel_b = FORWARD_STALL;
                load_stall_rs2 = 1'b1;
                forward_valid_b = 1'b0;
            end else begin
                forward_sel_b = FORWARD_EXE_B;
                forward_data_b = execute_result_b;
                exe_to_rs2 = 1'b1;
            end
        end
        
        // Priority 3: Writeback stage forwarding
        else if (writeback_valid && writeback_wr_en1 && (decode_rs2 == writeback_rd1)) begin
            forward_sel_b = FORWARD_WB_A;
            forward_data_b = writeback_result_a;
            wb_to_rs2 = 1'b1;
        end
        else if (writeback_valid && writeback_wr_en2 && (decode_rs2 == writeback_rd2)) begin
            forward_sel_b = FORWARD_WB_B;
            forward_data_b = writeback_result_b;
            wb_to_rs2 = 1'b1;
        end
    end
end

// ============================================================================
// RS3 FORWARDING LOGIC
// ============================================================================

always @(*) begin
    forward_sel_c = FORWARD_RF;
    forward_data_c = rf_data_c;
    forward_valid_c = 1'b1;
    exe_to_rs3 = 1'b0;
    wb_to_rs3 = 1'b0;
    mem_to_rs3 = 1'b0;
    load_stall_rs3 = 1'b0;
    
    if (decode_valid && (decode_rs3 != 5'h00)) begin
        // Priority 1: Memory stage forwarding (for loads)
        if (memory_valid && memory_wr_en && (decode_rs3 == memory_rd)) begin
            forward_sel_c = FORWARD_MEM;
            forward_data_c = memory_result;
            mem_to_rs3 = 1'b1;
        end
        
        // Priority 2: Execute stage forwarding
        else if (execute_valid && execute_wr_en1 && (decode_rs3 == execute_rd1)) begin
            if (execute_mem_read) begin
                // Load instruction in execute - need to stall
                forward_sel_c = FORWARD_STALL;
                load_stall_rs3 = 1'b1;
                forward_valid_c = 1'b0;
            end else begin
                forward_sel_c = FORWARD_EXE_A;
                forward_data_c = execute_result_a;
                exe_to_rs3 = 1'b1;
            end
        end
        else if (execute_valid && execute_wr_en2 && (decode_rs3 == execute_rd2)) begin
            if (execute_mem_read) begin
                // Load instruction in execute - need to stall
                forward_sel_c = FORWARD_STALL;
                load_stall_rs3 = 1'b1;
                forward_valid_c = 1'b0;
            end else begin
                forward_sel_c = FORWARD_EXE_B;
                forward_data_c = execute_result_b;
                exe_to_rs3 = 1'b1;
            end
        end
        
        // Priority 3: Writeback stage forwarding
        else if (writeback_valid && writeback_wr_en1 && (decode_rs3 == writeback_rd1)) begin
            forward_sel_c = FORWARD_WB_A;
            forward_data_c = writeback_result_a;
            wb_to_rs3 = 1'b1;
        end
        else if (writeback_valid && writeback_wr_en2 && (decode_rs3 == writeback_rd2)) begin
            forward_sel_c = FORWARD_WB_B;
            forward_data_c = writeback_result_b;
            wb_to_rs3 = 1'b1;
        end
    end
end

// ============================================================================
// LOAD FORWARDING STALL DETECTION
// ============================================================================

always @(*) begin
    load_forward_stall = load_stall_rs1 || load_stall_rs2 || load_stall_rs3;
end

// ============================================================================
// PERFORMANCE MONITORING COUNTERS
// ============================================================================

always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        forwards_from_execute <= 32'h0;
        forwards_from_writeback <= 32'h0;
        forwards_from_memory <= 32'h0;
        load_forward_stalls <= 32'h0;
    end else begin
        // Count execute stage forwards
        if (exe_to_rs1 || exe_to_rs2 || exe_to_rs3) begin
            forwards_from_execute <= forwards_from_execute + 1;
        end
        
        // Count writeback stage forwards
        if (wb_to_rs1 || wb_to_rs2 || wb_to_rs3) begin
            forwards_from_writeback <= forwards_from_writeback + 1;
        end
        
        // Count memory stage forwards
        if (mem_to_rs1 || mem_to_rs2 || mem_to_rs3) begin
            forwards_from_memory <= forwards_from_memory + 1;
        end
        
        // Count load forwarding stalls
        if (load_forward_stall) begin
            load_forward_stalls <= load_forward_stalls + 1;
        end
    end
end

// ============================================================================
// FORWARDING VALIDATION AND ERROR CHECKING
// ============================================================================

// Synthesis directive for simulation debug
// synthesis translate_off
always @(posedge clk) begin
    if (rst_n && decode_valid) begin
        // Check for invalid forwarding conditions
        if ((forward_sel_a == FORWARD_STALL) && forward_valid_a) begin
            $display("WARNING: Forwarding unit - Invalid stall condition for RS1");
        end
        if ((forward_sel_b == FORWARD_STALL) && forward_valid_b) begin
            $display("WARNING: Forwarding unit - Invalid stall condition for RS2");
        end
        if ((forward_sel_c == FORWARD_STALL) && forward_valid_c) begin
            $display("WARNING: Forwarding unit - Invalid stall condition for RS3");
        end
        
        // Debug forwarding activity
        if (exe_to_rs1 || exe_to_rs2 || exe_to_rs3) begin
            $display("FORWARD: Execute stage forward - RS1:%b RS2:%b RS3:%b", 
                    exe_to_rs1, exe_to_rs2, exe_to_rs3);
        end
        if (wb_to_rs1 || wb_to_rs2 || wb_to_rs3) begin
            $display("FORWARD: Writeback stage forward - RS1:%b RS2:%b RS3:%b", 
                    wb_to_rs1, wb_to_rs2, wb_to_rs3);
        end
    end
end
// synthesis translate_on

endmodule

`endif // FORWARDING_UNIT_V

