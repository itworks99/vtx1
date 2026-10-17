	`timescale 1ns / 1ps
// ============================================================================
// VTX1 Hazard Detection Unit
// ============================================================================
// Comprehensive hazard detection for VTX1 4-stage VLIW pipeline
// Handles RAW, WAW, WAR, structural, and control hazards
// ============================================================================

`ifndef HAZARD_DETECTION_V
`define HAZARD_DETECTION_V

`include "vtx1_interfaces.v"

module hazard_detection (
    input wire clk,
    input wire rst_n,
    
    // Pipeline stage control
    input wire fetch_valid,
    input wire decode_valid,
    input wire execute_valid,
    input wire writeback_valid,
    
    // Decode stage instruction fields
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] decode_rs1,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] decode_rs2,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] decode_rs3,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] decode_rd1,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] decode_rd2,
    input wire decode_wr_en1,
    input wire decode_wr_en2,
    input wire decode_mem_read,
    input wire decode_mem_write,
    input wire decode_branch,
    input wire [2:0] decode_op_a,
    input wire [2:0] decode_op_b,
    input wire [2:0] decode_op_c,
    
    // Execute stage instruction fields
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] execute_rd1,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] execute_rd2,
    input wire execute_wr_en1,
    input wire execute_wr_en2,
    input wire execute_mem_read,
    input wire execute_mem_write,
    input wire execute_branch,
    input wire [2:0] execute_op_a,
    input wire [2:0] execute_op_b,
    input wire [2:0] execute_op_c,
    
    // Writeback stage instruction fields
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] writeback_rd1,
    input wire [`VTX1_REG_ADDR_WIDTH-1:0] writeback_rd2,
    input wire writeback_wr_en1,
    input wire writeback_wr_en2,
    
    // Memory system status
    input wire mem_ready,
    input wire mem_error,
    input wire cache_miss,
    
    // Branch prediction status
    input wire branch_mispredict,
    input wire branch_taken,
    
    // Hazard detection outputs
    output reg raw_hazard,          // Read After Write
    output reg waw_hazard,          // Write After Write
    output reg war_hazard,          // Write After Read
    output reg structural_hazard,   // Resource conflicts
    output reg control_hazard,      // Branch/jump hazards
    output reg memory_hazard,       // Memory dependency hazards
    output reg load_use_hazard,     // Load-use data hazard
    
    // Hazard sources for debugging
    output reg [2:0] raw_source,    // Which register causes RAW
    output reg [2:0] waw_source,    // Which register causes WAW
    output reg [2:0] hazard_stage,  // Which stage has hazard
    
    // Pipeline control outputs
    output reg pipeline_stall,
    output reg pipeline_flush,
    output reg [1:0] stall_cycles,
    
    // Performance monitoring
    output reg [31:0] raw_hazards_detected,
    output reg [31:0] waw_hazards_detected,
    output reg [31:0] structural_hazards_detected,
    output reg [31:0] memory_hazards_detected,
    output reg [31:0] control_hazards_detected
);

// ============================================================================
// CONSTANTS AND PARAMETERS
// ============================================================================

// Hazard source encoding
localparam SRC_RS1 = 3'b001;
localparam SRC_RS2 = 3'b010;
localparam SRC_RS3 = 3'b011;
localparam SRC_RD1 = 3'b100;
localparam SRC_RD2 = 3'b101;

// Stage encoding
localparam STAGE_DECODE    = 3'b001;
localparam STAGE_EXECUTE   = 3'b010;
localparam STAGE_WRITEBACK = 3'b011;

// Operation types for structural hazard detection
localparam OP_ALU     = 3'b000;
localparam OP_MUL     = 3'b001;
localparam OP_DIV     = 3'b010;
localparam OP_FPU     = 3'b011;
localparam OP_MEMORY  = 3'b100;
localparam OP_BRANCH  = 3'b101;
localparam OP_SYSTEM  = 3'b110;

// ============================================================================
// INTERNAL SIGNALS
// ============================================================================

reg raw_decode_execute, raw_decode_writeback;
reg waw_decode_execute, waw_decode_writeback;
reg war_execute_decode;
reg structural_vliw_conflict;
reg load_dependency;

// ============================================================================
// RAW HAZARD DETECTION (Read After Write)
// ============================================================================

always @(*) begin
    raw_decode_execute = 1'b0;
    raw_decode_writeback = 1'b0;
    raw_source = 3'b000;
    
    if (decode_valid && execute_valid) begin
        // Check RS1 against execute destination registers
        if (decode_rs1 != 5'h00 && execute_wr_en1 && (decode_rs1 == execute_rd1)) begin
            raw_decode_execute = 1'b1;
            raw_source = SRC_RS1;
        end else if (decode_rs1 != 5'h00 && execute_wr_en2 && (decode_rs1 == execute_rd2)) begin
            raw_decode_execute = 1'b1;
            raw_source = SRC_RS1;
        end
        
        // Check RS2 against execute destination registers
        if (decode_rs2 != 5'h00 && execute_wr_en1 && (decode_rs2 == execute_rd1)) begin
            raw_decode_execute = 1'b1;
            raw_source = SRC_RS2;
        end else if (decode_rs2 != 5'h00 && execute_wr_en2 && (decode_rs2 == execute_rd2)) begin
            raw_decode_execute = 1'b1;
            raw_source = SRC_RS2;
        end
        
        // Check RS3 against execute destination registers
        if (decode_rs3 != 5'h00 && execute_wr_en1 && (decode_rs3 == execute_rd1)) begin
            raw_decode_execute = 1'b1;
            raw_source = SRC_RS3;
        end else if (decode_rs3 != 5'h00 && execute_wr_en2 && (decode_rs3 == execute_rd2)) begin
            raw_decode_execute = 1'b1;
            raw_source = SRC_RS3;
        end
    end
    
    // Check decode against writeback (can be forwarded)
    if (decode_valid && writeback_valid) begin
        if (decode_rs1 != 5'h00 && writeback_wr_en1 && (decode_rs1 == writeback_rd1)) begin
            raw_decode_writeback = 1'b1;
        end else if (decode_rs1 != 5'h00 && writeback_wr_en2 && (decode_rs1 == writeback_rd2)) begin
            raw_decode_writeback = 1'b1;
        end
        
        if (decode_rs2 != 5'h00 && writeback_wr_en1 && (decode_rs2 == writeback_rd1)) begin
            raw_decode_writeback = 1'b1;
        end else if (decode_rs2 != 5'h00 && writeback_wr_en2 && (decode_rs2 == writeback_rd2)) begin
            raw_decode_writeback = 1'b1;
        end
        
        if (decode_rs3 != 5'h00 && writeback_wr_en1 && (decode_rs3 == writeback_rd1)) begin
            raw_decode_writeback = 1'b1;
        end else if (decode_rs3 != 5'h00 && writeback_wr_en2 && (decode_rs3 == writeback_rd2)) begin
            raw_decode_writeback = 1'b1;
        end
    end
    
    raw_hazard = raw_decode_execute; // Writeback can be forwarded
end

// ============================================================================
// WAW HAZARD DETECTION (Write After Write)
// ============================================================================

always @(*) begin
    waw_decode_execute = 1'b0;
    waw_decode_writeback = 1'b0;
    waw_source = 3'b000;
    
    if (decode_valid && execute_valid) begin
        // Check decode destination against execute destinations
        if (decode_wr_en1 && execute_wr_en1 && 
            (decode_rd1 != 5'h00) && (decode_rd1 == execute_rd1)) begin
            waw_decode_execute = 1'b1;
            waw_source = SRC_RD1;
        end else if (decode_wr_en1 && execute_wr_en2 && 
                    (decode_rd1 != 5'h00) && (decode_rd1 == execute_rd2)) begin
            waw_decode_execute = 1'b1;
            waw_source = SRC_RD1;
        end
        
        if (decode_wr_en2 && execute_wr_en1 && 
            (decode_rd2 != 5'h00) && (decode_rd2 == execute_rd1)) begin
            waw_decode_execute = 1'b1;
            waw_source = SRC_RD2;
        end else if (decode_wr_en2 && execute_wr_en2 && 
                    (decode_rd2 != 5'h00) && (decode_rd2 == execute_rd2)) begin
            waw_decode_execute = 1'b1;
            waw_source = SRC_RD2;
        end
    end
    
    waw_hazard = waw_decode_execute;
end

// ============================================================================
// WAR HAZARD DETECTION (Write After Read) - Less common in our pipeline
// ============================================================================

always @(*) begin
    war_execute_decode = 1'b0;
    
    // In our 4-stage pipeline, WAR hazards are less critical
    // but can occur with out-of-order completion
    war_hazard = war_execute_decode;
end

// ============================================================================
// STRUCTURAL HAZARD DETECTION
// ============================================================================

always @(*) begin
    structural_vliw_conflict = 1'b0;
    
    if (decode_valid) begin
        // Check for VLIW functional unit conflicts
        // Each VLIW instruction can have up to 3 operations
        
        // Check if multiple operations require the same functional unit
        if ((decode_op_a == decode_op_b) && (decode_op_a != OP_ALU)) begin
            structural_vliw_conflict = 1'b1;
        end else if ((decode_op_a == decode_op_c) && (decode_op_a != OP_ALU)) begin
            structural_vliw_conflict = 1'b1;
        end else if ((decode_op_b == decode_op_c) && (decode_op_b != OP_ALU)) begin
            structural_vliw_conflict = 1'b1;
        end
        
        // Check for memory port conflicts
        if ((decode_op_a == OP_MEMORY) && (decode_op_b == OP_MEMORY)) begin
            structural_vliw_conflict = 1'b1;
        end else if ((decode_op_a == OP_MEMORY) && (decode_op_c == OP_MEMORY)) begin
            structural_vliw_conflict = 1'b1;
        end else if ((decode_op_b == OP_MEMORY) && (decode_op_c == OP_MEMORY)) begin
            structural_vliw_conflict = 1'b1;
        end
    end
    
    structural_hazard = structural_vliw_conflict;
end

// ============================================================================
// MEMORY HAZARD DETECTION
// ============================================================================

always @(*) begin
    load_dependency = 1'b0;
    
    // Load-use hazard: previous instruction loads, current instruction uses
    if (decode_valid && execute_valid && execute_mem_read) begin
        if ((decode_rs1 == execute_rd1) || (decode_rs2 == execute_rd1) || 
            (decode_rs3 == execute_rd1)) begin
            load_dependency = 1'b1;
        end
        if ((decode_rs1 == execute_rd2) || (decode_rs2 == execute_rd2) || 
            (decode_rs3 == execute_rd2)) begin
            load_dependency = 1'b1;
        end
    end
    
    load_use_hazard = load_dependency;
    
    // Memory system hazards
    memory_hazard = !mem_ready || mem_error || cache_miss || 
                   (decode_valid && execute_valid && 
                    decode_mem_write && execute_mem_read) ||
                   (decode_valid && execute_valid && 
                    decode_mem_read && execute_mem_write);
end

// ============================================================================
// CONTROL HAZARD DETECTION
// ============================================================================

always @(*) begin
    control_hazard = branch_mispredict || 
                    (decode_valid && decode_branch) ||
                    (execute_valid && execute_branch && branch_taken);
end

// ============================================================================
// PIPELINE CONTROL LOGIC
// ============================================================================

always @(*) begin
    pipeline_stall = raw_hazard || waw_hazard || structural_hazard || 
                    load_use_hazard || memory_hazard;
    
    pipeline_flush = control_hazard || branch_mispredict;
    
    // Determine stall cycles needed
    if (load_use_hazard || memory_hazard) begin
        stall_cycles = 2'b10; // 2 cycles
    end else if (raw_hazard || waw_hazard) begin
        stall_cycles = 2'b01; // 1 cycle
    end else begin
        stall_cycles = 2'b00; // No stall
    end
    
    // Set hazard stage indicator
    if (raw_hazard || waw_hazard) begin
        hazard_stage = STAGE_EXECUTE;
    end else if (memory_hazard) begin
        hazard_stage = STAGE_EXECUTE;
    end else if (control_hazard) begin
        hazard_stage = STAGE_DECODE;
    end else begin
        hazard_stage = 3'b000;
    end
end

// ============================================================================
// PERFORMANCE MONITORING COUNTERS
// ============================================================================

always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        raw_hazards_detected <= 32'h0;
        waw_hazards_detected <= 32'h0;
        structural_hazards_detected <= 32'h0;
        memory_hazards_detected <= 32'h0;
        control_hazards_detected <= 32'h0;
    end else begin
        if (raw_hazard) raw_hazards_detected <= raw_hazards_detected + 1;
        if (waw_hazard) waw_hazards_detected <= waw_hazards_detected + 1;
        if (structural_hazard) structural_hazards_detected <= structural_hazards_detected + 1;
        if (memory_hazard) memory_hazards_detected <= memory_hazards_detected + 1;
        if (control_hazard) control_hazards_detected <= control_hazards_detected + 1;
    end
end

endmodule

`endif // HAZARD_DETECTION_V

