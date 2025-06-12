`timescale 1ns / 1ps
// VTX1 CPU Core - 4-Stage Pipeline Implementation
// Part of the VTX1 Ternary System-on-Chip

`ifndef CPU_CORE_V
`define CPU_CORE_V

// Include VTX1 interface definitions
// Note: Include paths are handled by compiler flags in Taskfile.yml
`include "ternary_constants.v"
`include "vtx1_interfaces.v"
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

// ============================================================================
// VTX1 CPU CORE - 4-STAGE PIPELINE
// ============================================================================
// Implements the complete VTX1 CPU with:
// - 4-stage pipeline (Fetch, Decode, Execute, Writeback)
// - VLIW instruction support (3 operations per cycle)
// - Integrated TCU (Ternary Computing Unit)
// - Hazard detection and forwarding
// - Microcode offloading for complex operations

module cpu_core (
    input  wire                     clk,
    input  wire                     rst_n,
    input  wire                     enable,
      // Instruction Memory Interface
    output reg  [`VTX1_WORD_WIDTH-1:0]  imem_addr,
    input  wire [`VTX1_VLIW_WIDTH-1:0]  imem_data,
    output reg                           imem_req,
    input  wire                          imem_ready,
    
    // Data Memory Interface
    output reg  [`VTX1_WORD_WIDTH-1:0]  dmem_addr,
    inout  wire [`VTX1_WORD_WIDTH-1:0]  dmem_data,
    output reg                           dmem_we,
    output reg                           dmem_oe,
    output reg                           dmem_req,
    input  wire                          dmem_ready,
    
    // Interrupt Interface
    input  wire [15:0]                   interrupt_req,
    output reg  [15:0]                   interrupt_ack,
    input  wire                          nmi_req,
    
    // Debug Interface
    input  wire                          debug_enable,
    input  wire                          debug_step,
    input  wire [3:0]                    debug_cmd,
    output reg  [`VTX1_WORD_WIDTH-1:0]  debug_pc,
    output reg  [3:0]                    debug_status,
    
    // Performance Counters
    output reg  [31:0]                   cycle_count,
    output reg  [31:0]                   instruction_count,
    output reg  [31:0]                   stall_count,
    
    // Pipeline Status
    output reg                           pipeline_stall,
    output reg                           pipeline_flush,
    output reg  [1:0]                    pipeline_state,
    
    // Error handling interface
    output reg                           error,
    output reg  [3:0]                    error_code,
    output reg                           timeout,
    output reg  [31:0]                   error_count
);

    // ========================================================================
    // PIPELINE STAGE PARAMETERS
    // ========================================================================
    localparam STAGE_FETCH     = 2'b00;
    localparam STAGE_DECODE    = 2'b01;
    localparam STAGE_EXECUTE   = 2'b10;
    localparam STAGE_WRITEBACK = 2'b11;
    
    // Debug command codes
    localparam DEBUG_RUN       = 4'b0000;
    localparam DEBUG_HALT      = 4'b0001;
    localparam DEBUG_STEP      = 4'b0010;
    localparam DEBUG_RESET     = 4'b0011;
      // ========================================================================
    // PIPELINE REGISTERS
    // ========================================================================
    
    // Fetch Stage
    reg [`VTX1_WORD_WIDTH-1:0]   pc_fetch;
    reg [`VTX1_VLIW_WIDTH-1:0]   instruction_fetch;
    reg                          fetch_valid;
    
    // Decode Stage  
    reg [`VTX1_WORD_WIDTH-1:0]   pc_decode;
    reg [`VTX1_VLIW_WIDTH-1:0]   instruction_decode;
    reg                          decode_valid;
    
    // Execute Stage
    reg [`VTX1_WORD_WIDTH-1:0]   pc_execute;
    reg [`VTX1_VLIW_WIDTH-1:0]   instruction_execute;
    reg                          execute_valid;
    
    // Writeback Stage
    reg [`VTX1_WORD_WIDTH-1:0]   pc_writeback;
    reg [`VTX1_VLIW_WIDTH-1:0]   instruction_writeback;
    reg                          writeback_valid;
      // ========================================================================
    // ERROR HANDLING AND TIMEOUTS
    // ========================================================================
    
    reg [31:0] cpu_timeout_counter;
    reg [31:0] cpu_fetch_timeout_counter;
    reg [31:0] cpu_memory_timeout_counter;    reg        fetch_error, decode_error, execute_error, memory_error;
    reg        pipeline_error;
    // CONTROL SIGNALS
    // ========================================================================
    reg                     hazard_detected;
    reg                     branch_taken;
    reg [`VTX1_WORD_WIDTH-1:0]  branch_target;
    wire                    flush_pipeline;
    wire                    stall_pipeline;
      // Register file interface
    wire [`VTX1_REG_ADDR_WIDTH-1:0] rf_read_addr_a, rf_read_addr_b, rf_read_addr_c;
    wire [`VTX1_WORD_WIDTH-1:0]     rf_read_data_a, rf_read_data_b, rf_read_data_c;
    wire                            rf_write_enable_a, rf_write_enable_b;
    wire [`VTX1_REG_ADDR_WIDTH-1:0] rf_write_addr_a, rf_write_addr_b;
    wire [`VTX1_WORD_WIDTH-1:0]     rf_write_data_a, rf_write_data_b;      // TCU interface (legacy execution unit signals - simplified for basic pipeline)
    wire                    tcu_enable;
    wire [`VTX1_WORD_WIDTH-1:0]  tcu_operand_a, tcu_operand_b, tcu_operand_c;
    wire [3:0]              tcu_operation;
    wire [`VTX1_WORD_WIDTH-1:0]  tcu_result;
    wire                    tcu_valid;
    wire                    tcu_ready;
    wire                    tcu_error;
    
    // Enhanced TCU Interface
    wire                    microcode_enable;
    wire [3:0]              microcode_operation;
    wire [`VTX1_WORD_WIDTH-1:0] microcode_operand_a, microcode_operand_b, microcode_operand_c;
    wire [`VTX1_WORD_WIDTH-1:0] microcode_result;
    wire                    microcode_valid;
    wire                    microcode_ready;
    wire                    microcode_error;
    wire [3:0]              enhanced_interface_state;
    wire [31:0]             enhanced_operation_cycles;
    wire [31:0]             enhanced_total_operations;
    wire [31:0]             enhanced_error_count;
    
    // TCU system control interface signals
    wire                    mem_req, mem_wr;
    wire [1:0]              mem_size;
    wire [`VTX1_ADDR_WIDTH-1:0] mem_addr;
    wire [`VTX1_WORD_WIDTH-1:0] mem_wdata, mem_rdata;
    wire                    mem_ready, mem_error;
    
    wire                    icache_req, icache_hit, icache_ready;
    wire [`VTX1_ADDR_WIDTH-1:0] icache_addr;
    wire [`VTX1_WORD_WIDTH-1:0] icache_data;
    
    wire                    dcache_req, dcache_wr, dcache_hit, dcache_ready;
    wire [1:0]              dcache_size;
    wire [`VTX1_ADDR_WIDTH-1:0] dcache_addr;
    wire [`VTX1_WORD_WIDTH-1:0] dcache_wdata, dcache_rdata;
      // Interrupt signals removed - not implemented
    // Interrupt vectors removed - not implemented
    // Interrupt level removed - not implemented
    
    // Simplified control signals for arithmetic unit
    wire                    cpu_enable, debug_mode, single_step;
    wire [3:0]              tcu_error_code;// Decoder outputs
    wire [2:0]              decoded_op_type;
    wire [`VTX1_REG_ADDR_WIDTH-1:0] decoded_src_reg_a, decoded_src_reg_b, decoded_src_reg_c;
    wire [`VTX1_REG_ADDR_WIDTH-1:0] decoded_dst_reg_a, decoded_dst_reg_b;
    wire [`VTX1_WORD_WIDTH-1:0]     decoded_immediate;
    wire                            decoded_valid;
      // Additional decoder outputs
    wire [`OPCODE_WIDTH-1:0] decoded_opcode_a, decoded_opcode_b, decoded_opcode_c;
    wire [11:0]             decoded_immediate_a, decoded_immediate_b, decoded_immediate_c;
    wire                    decoded_alu_enable_a, decoded_alu_enable_b, decoded_alu_enable_c;
    wire [3:0]              decoded_alu_operation_a, decoded_alu_operation_b, decoded_alu_operation_c;
    wire                    decoded_memory_enable_a, decoded_memory_enable_b, decoded_memory_enable_c;
    wire                    decoded_memory_write_a, decoded_memory_write_b, decoded_memory_write_c;
    wire                    decoded_branch_enable;
    wire [2:0]              decoded_branch_condition;
    wire [`VTX1_WORD_WIDTH-1:0]  decoded_branch_target;// Error handling signals
    wire rf_error, rf_timeout, tcu_timeout;
    wire decoder_error, decoder_timeout;
    
    // Pipeline control signals
    wire pipeline_stall_hazard, pipeline_flush_hazard;
    wire [1:0] stall_cycles_needed;
    wire load_forward_stall;
    
    // Forwarding control signals
    wire [2:0] forward_sel_a, forward_sel_b, forward_sel_c;
    wire [`VTX1_WORD_WIDTH-1:0] forward_data_a, forward_data_b, forward_data_c;
    wire forward_valid_a, forward_valid_b, forward_valid_c;
    
    // Hazard detection signals
    wire raw_hazard, waw_hazard, war_hazard, structural_hazard;
    wire control_hazard, memory_hazard, load_use_hazard;
    wire [2:0] hazard_stage;
    
    // ========================================================================
    // COMPONENT INSTANTIATIONS
    // ========================================================================
      // Register File
    register_file rf_inst (
        .clk(clk),
        .rst_n(rst_n),
        .read_addr_a(rf_read_addr_a),
        .read_addr_b(rf_read_addr_b),
        .read_addr_c(rf_read_addr_c),
        .read_data_a(rf_read_data_a),
        .read_data_b(rf_read_data_b),
        .read_data_c(rf_read_data_c),
        .write_enable_a(rf_write_enable_a),
        .write_enable_b(rf_write_enable_b),
        .write_addr_a(rf_write_addr_a),
        .write_addr_b(rf_write_addr_b),
        .write_data_a(rf_write_data_a),
        .write_data_b(rf_write_data_b),
        .debug_enable(debug_enable),
        .debug_addr(debug_cmd),
        .debug_data(),
        .error(rf_error),
        .error_code(),
        .timeout(rf_timeout),
        .write_conflict(),
        .invalid_address(),
        .current_state(),
        .operation_count(),
        .read_count(),
        .write_count(),        .error_count()
    );
    
    // Hazard Detection Unit
    hazard_detection hazard_unit (
        .clk(clk),
        .rst_n(rst_n),
        .fetch_valid(fetch_valid),
        .decode_valid(decode_valid),
        .execute_valid(execute_valid),
        .writeback_valid(writeback_valid),
        .decode_rs1(decoded_src_reg_a),
        .decode_rs2(decoded_src_reg_b),        .decode_rs3(decoded_src_reg_c),
        .decode_rd1(decoded_dst_reg_a),
        .decode_rd2(decoded_dst_reg_b),
        .decode_wr_en1(decoded_valid),
        .decode_wr_en2(1'b0), // Single write port in basic implementation
        .decode_mem_read(decoded_op_type == 3'b100), // Memory read operation
        .decode_mem_write(decoded_op_type == 3'b101), // Memory write operation
        .decode_branch(decoded_op_type == 3'b110), // Branch operation
        .decode_op_a(decoded_op_type),
        .decode_op_b(3'b000), // VLIW operation B (simplified)
        .decode_op_c(3'b000), // VLIW operation C (simplified)
        .execute_rd1(decoded_dst_reg_a), // Execute stage destination (delayed)
        .execute_rd2(decoded_dst_reg_b),
        .execute_wr_en1(execute_valid),
        .execute_wr_en2(1'b0),
        .execute_mem_read(decoded_op_type == 3'b100),
        .execute_mem_write(decoded_op_type == 3'b101),
        .execute_branch(decoded_op_type == 3'b110),
        .execute_op_a(decoded_op_type),
        .execute_op_b(3'b000),        .execute_op_c(3'b000),
        .writeback_rd1(decoded_dst_reg_a), // Writeback stage destination (delayed)
        .writeback_rd2(decoded_dst_reg_b),
        .writeback_wr_en1(writeback_valid),
        .writeback_wr_en2(1'b0),
        .mem_ready(dmem_ready && imem_ready),
        .mem_error(memory_error),
        .cache_miss(1'b0), // Cache miss detection (simplified)
        .branch_mispredict(branch_taken), // Branch misprediction (simplified)
        .branch_taken(branch_taken),
        .raw_hazard(raw_hazard),
        .waw_hazard(waw_hazard),
        .war_hazard(war_hazard),
        .structural_hazard(structural_hazard),
        .control_hazard(control_hazard),
        .memory_hazard(memory_hazard),
        .load_use_hazard(load_use_hazard),
        .raw_source(),
        .waw_source(),
        .hazard_stage(hazard_stage),
        .pipeline_stall(pipeline_stall_hazard),
        .pipeline_flush(pipeline_flush_hazard),
        .stall_cycles(stall_cycles_needed),
        .raw_hazards_detected(),
        .waw_hazards_detected(),
        .structural_hazards_detected(),
        .memory_hazards_detected(),
        .control_hazards_detected()
    );
    
    // Forwarding Unit
    forwarding_unit forward_unit (
        .clk(clk),
        .rst_n(rst_n),
        .decode_valid(decode_valid),
        .execute_valid(execute_valid),
        .writeback_valid(writeback_valid),
        .decode_rs1(decoded_src_reg_a),
        .decode_rs2(decoded_src_reg_b),
        .decode_rs3(decoded_src_reg_c),
        .execute_result_a(tcu_result),        .execute_result_b({`VTX1_WORD_WIDTH{1'b0}}), // Second result (simplified)
        .execute_rd1(decoded_dst_reg_a), // Execute stage destination
        .execute_rd2(decoded_dst_reg_b),
        .execute_wr_en1(execute_valid),
        .execute_wr_en2(1'b0),
        .execute_mem_read(decoded_op_type == 3'b100),        .writeback_result_a(tcu_result), // Writeback result (delayed)
        .writeback_result_b({`VTX1_WORD_WIDTH{1'b0}}),
        .writeback_rd1(decoded_dst_reg_a), // Writeback destination (delayed)
        .writeback_rd2(decoded_dst_reg_b),
        .writeback_wr_en1(writeback_valid),
        .writeback_wr_en2(1'b0),        .memory_result({`VTX1_WORD_WIDTH{1'b0}}), // Memory result (simplified)
        .memory_rd(4'h0),
        .memory_wr_en(1'b0),
        .memory_valid(1'b0),
        .rf_data_a(rf_read_data_a),
        .rf_data_b(rf_read_data_b),
        .rf_data_c(rf_read_data_c),
        .forward_data_a(forward_data_a),
        .forward_data_b(forward_data_b),
        .forward_data_c(forward_data_c),
        .forward_sel_a(forward_sel_a),
        .forward_sel_b(forward_sel_b),
        .forward_sel_c(forward_sel_c),
        .forward_valid_a(forward_valid_a),
        .forward_valid_b(forward_valid_b),
        .forward_valid_c(forward_valid_c),
        .load_forward_stall(load_forward_stall),
        .forwards_from_execute(),
        .forwards_from_writeback(),
        .forwards_from_memory(),
        .load_forward_stalls()
    );    // Instruction Decoder
    instruction_decoder decoder_inst (
        .clk(clk),
        .rst_n(rst_n),
        .enable(decode_valid),
        .instruction(instruction_decode),
        .opcode_a(decoded_opcode_a),
        .opcode_b(decoded_opcode_b), 
        .opcode_c(decoded_opcode_c),
        .rs1_a(decoded_src_reg_a),
        .rs2_a(decoded_src_reg_b),
        .rs3_a(decoded_src_reg_c),
        .rd_a(decoded_dst_reg_a),
        .rd_b(decoded_dst_reg_b),
        .rs1_b(),  // Slot B not used in current implementation
        .rs2_b(),
        .rs3_b(),
        .rs1_c(),  // Slot C not used in current implementation
        .rs2_c(), 
        .rs3_c(),
        .rd_c(),
        .immediate_a(decoded_immediate_a),
        .immediate_b(decoded_immediate_b),
        .immediate_c(decoded_immediate_c),
        .alu_enable_a(decoded_alu_enable_a),
        .alu_enable_b(decoded_alu_enable_b),
        .alu_enable_c(decoded_alu_enable_c),
        .alu_operation_a(decoded_alu_operation_a),
        .alu_operation_b(decoded_alu_operation_b),
        .alu_operation_c(decoded_alu_operation_c),
        .memory_enable_a(decoded_memory_enable_a),
        .memory_enable_b(decoded_memory_enable_b),
        .memory_enable_c(decoded_memory_enable_c),
        .memory_write_a(decoded_memory_write_a),
        .memory_write_b(decoded_memory_write_b),
        .memory_write_c(decoded_memory_write_c),
        .branch_enable(decoded_branch_enable),
        .branch_condition(decoded_branch_condition),
        .branch_target(decoded_branch_target),
        .valid_a(decoded_valid),
        .valid_b(),  // Slot B not used
        .valid_c(),  // Slot C not used
        .decode_error(),
        .error(decoder_error),
        .error_code(),
        .timeout(decoder_timeout),
        .error_clear(1'b0),
        .state(),        .decode_count(),
        .error_count()
    );
      // ========================================================================
    // OPERATION TYPE CLASSIFICATION
    // ========================================================================
    // Map instruction categories to operation types for proper TCU control
    
    // Instruction category constants (from instruction_decoder.v)
    localparam CATEGORY_ALU     = 4'b0001;  // Arithmetic/Logic operations
    localparam CATEGORY_MEMORY  = 4'b0010;  // Memory operations
    localparam CATEGORY_BRANCH  = 4'b0011;  // Branch/Control operations
    localparam CATEGORY_SPECIAL = 4'b0100;  // Special operations (NOP, HALT)
    
    // Operation type encoding for CPU control
    localparam OP_TYPE_SPECIAL  = 3'b000;  // Special/NOP - no TCU
    localparam OP_TYPE_ALU      = 3'b001;  // Ternary arithmetic - enable TCU
    localparam OP_TYPE_MEMORY_R = 3'b100;  // Memory read
    localparam OP_TYPE_MEMORY_W = 3'b101;  // Memory write
    localparam OP_TYPE_BRANCH   = 3'b110;  // Branch - use microcode    // Extract instruction category from opcode_a[11:8] 
    wire [3:0] instruction_category = decoded_opcode_a[11:8];
    
    // Registered operation type for stable classification
    reg [2:0] decoded_op_type_reg;
    
    // Classify operation type based on instruction category and specific operation
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            decoded_op_type_reg <= OP_TYPE_SPECIAL; // Default to no TCU activation
        end else if (enable && decode_valid) begin
            case (instruction_category)
                CATEGORY_SPECIAL: decoded_op_type_reg <= OP_TYPE_SPECIAL;
                CATEGORY_ALU:     decoded_op_type_reg <= OP_TYPE_ALU;
                CATEGORY_MEMORY:  decoded_op_type_reg <= decoded_memory_write_a ? OP_TYPE_MEMORY_W : OP_TYPE_MEMORY_R;
                CATEGORY_BRANCH:  decoded_op_type_reg <= OP_TYPE_BRANCH;
                default:          decoded_op_type_reg <= OP_TYPE_SPECIAL; // Safe default
            endcase
        end
    end
    
    // Use registered version for TCU control
    assign decoded_op_type = decoded_op_type_reg;
    
    // Enhanced immediate value assignment (use slot A immediate for now)
    assign decoded_immediate = {{(`VTX1_WORD_WIDTH-12){decoded_immediate_a[11]}}, decoded_immediate_a};
      // ========================================================================
    // TERNARY ARITHMETIC EXECUTION UNIT
    // ========================================================================
    ternary_arithmetic_unit ternary_exec_unit (
        .clk(clk),
        .rst_n(rst_n),
        
        // Input operands (from register file)
        .operand_a(tcu_operand_a),
        .operand_b(tcu_operand_b),
        
        // Operation control
        .operation(tcu_operation[3:0]),  // Map operation field
        .operation_enable(tcu_enable && (decoded_op_type == 3'b001)), // Enable only for ALU operations
        
        // Output
        .result(tcu_result),
        .operation_complete(tcu_ready),
        .arithmetic_overflow(),  // TODO: Connect to status flags
        .arithmetic_underflow(), // TODO: Connect to status flags
        
        // Enhanced error handling
        .error(tcu_error),
        .error_code(tcu_error_code[3:0]),
        .timeout(tcu_timeout),
        .error_clear(1'b0),
        
        // Debug interface
        .state(),  // Not connected for now
        .operation_count(),
        .error_count()
    );
    
    // Enhanced TCU Interface
    tcu_enhanced_interface tcu_enhanced_inst (
        .clk(clk),
        .rst_n(rst_n),
        
        // Microcode Sequencer Interface (Enhanced)
        .microcode_enable(microcode_enable),
        .microcode_operation(microcode_operation),
        .microcode_operand_a(microcode_operand_a),
        .microcode_operand_b(microcode_operand_b),
        .microcode_operand_c(microcode_operand_c),
        .microcode_result(microcode_result),
        .microcode_valid(microcode_valid),
        .microcode_ready(microcode_ready),
        .microcode_error(microcode_error),
        
        // Legacy TCU Interface (Existing CPU Core)
        .tcu_enable(tcu_enable),
        .tcu_operation(tcu_operation),
        .tcu_operand_a(tcu_operand_a),
        .tcu_operand_b(tcu_operand_b),
        .tcu_operand_c(tcu_operand_c),
        .tcu_result(tcu_result),
        .tcu_valid(tcu_valid),
        .tcu_ready(tcu_ready),
        .tcu_error(tcu_error),
        
        // Enhanced Control and Status
        .interface_state(enhanced_interface_state),
        .operation_cycles(enhanced_operation_cycles),
        .total_operations(enhanced_total_operations),
        .error_count(enhanced_error_count)
    );

    // Microcode Sequencer
    microcode_sequencer microcode_inst (
        .clk(clk),
        .rst_n(rst_n),
        
        // CPU Interface
        .enable(execute_valid && (decoded_op_type == 3'b110)), // Enable for microcode operations
        .start(execute_valid && (decoded_op_type == 3'b110)),
        .opcode({2'b00, decoded_op_type, 1'b0}), // Convert to 6-bit microcode opcode
        .operand_a(tcu_operand_a),
        .operand_b(tcu_operand_b),
        .operand_c(tcu_operand_c),
        
        .result(microcode_result),
        .valid(microcode_valid),
        .ready(microcode_ready),
        .error(microcode_error),
        .error_code(),
        
        // Microcode ROM Interface
        .rom_addr(),
        .rom_data(32'h0), // Simplified for now
        .rom_enable(),
        .rom_ready(1'b1),
        .rom_error(1'b0),
        
        // Enhanced TCU Interface
        .tcu_enable(microcode_enable),
        .tcu_operation(microcode_operation),
        .tcu_operand_a(microcode_operand_a),
        .tcu_operand_b(microcode_operand_b),
        .tcu_operand_c(microcode_operand_c),
        .tcu_result(microcode_result),
        .tcu_valid(microcode_valid),
        .tcu_ready(microcode_ready),
        .tcu_error(microcode_error),
        
        // Register File Interface (simplified for basic integration)
        .reg_read_addr_a(),
        .reg_read_addr_b(),
        .reg_read_addr_c(),
        .reg_read_data_a(rf_read_data_a),
        .reg_read_data_b(rf_read_data_b),
        .reg_read_data_c(rf_read_data_c),
        .reg_write_enable(),
        .reg_write_addr(),
        .reg_write_data(),
          // Debug Interface
        .microcode_state(),
        .instruction_count(),
        .cycle_count(),
        .operation_count()
    );

    // ========================================================================
    // TCU SYSTEM CONTROL INTERFACE DEFAULTS
    // ========================================================================
    // Simplified assignments for current pipeline implementation
    // These would be enhanced in a full system implementation
    assign mem_rdata = dmem_data;  // Connect data memory read data
    assign mem_ready = dmem_ready;
    assign mem_error = 1'b0;       // Simplified error handling
    
    assign icache_data = imem_data[`VTX1_WORD_WIDTH-1:0]; // Use lower portion of VLIW
    assign icache_hit = 1'b1;      // Always hit for simplified implementation
    assign icache_ready = imem_ready;
    
    assign dcache_rdata = dmem_data;
    assign dcache_hit = 1'b1;      // Always hit for simplified implementation
    assign dcache_ready = dmem_ready;
    
    assign irq = |interrupt_req;   // Any interrupt request
    assign irq_vector = 8'h0;      // Simplified interrupt vector
    assign irq_level = 2'h0;       // Simplified interrupt level
    
    assign cpu_enable = enable;
    assign debug_mode = debug_enable;
    assign single_step = debug_step;
      // Enhanced execution unit with microcode integration
    assign tcu_result = (decoded_op_type == 3'b110) ? microcode_result : tcu_operand_a;  // Use microcode result for complex ops
    assign tcu_valid = (decoded_op_type == 3'b110) ? microcode_valid : tcu_enable;
    assign tcu_ready = (decoded_op_type == 3'b110) ? microcode_ready : 1'b1;             // Use enhanced ready signal
      // ========================================================================
    // PIPELINE CONTROL
    // ========================================================================
    assign stall_pipeline = pipeline_stall_hazard || load_forward_stall || 
                           !imem_ready || !dmem_ready || !tcu_ready;
    assign flush_pipeline = pipeline_flush_hazard || control_hazard || 
                           (debug_enable && debug_cmd == DEBUG_RESET);
    
    // Enhanced hazard-aware pipeline control
    always @(*) begin
        hazard_detected = raw_hazard || waw_hazard || structural_hazard || 
                         memory_hazard || load_use_hazard;
        branch_taken = control_hazard || (debug_enable && debug_cmd == DEBUG_STEP);
        branch_target = pc_fetch + 4; // Simplified branch target calculation
    end
    
    // ========================================================================
    // ERROR HANDLING LOGIC
    // ========================================================================
    
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            error <= 1'b0;            error_code <= `VTX1_ERROR_NONE;
            timeout <= 1'b0;
            error_count <= 32'h0;
            cpu_timeout_counter <= 32'h0;
            cpu_fetch_timeout_counter <= 32'h0;
            cpu_memory_timeout_counter <= 32'h0;
        end else begin
            // Aggregate error detection
            pipeline_error = rf_error || tcu_error || decoder_error;
              // Timeout management
            if (imem_req && !imem_ready) begin
                cpu_fetch_timeout_counter <= cpu_fetch_timeout_counter + 1;
                if (cpu_fetch_timeout_counter >= `VTX1_TIMEOUT_CYCLES) begin
                    fetch_error <= 1'b1;
                    error <= 1'b1;
                    error_code <= `VTX1_ERROR_TIMEOUT;
                    timeout <= 1'b1;
                    error_count <= error_count + 1;
                end            end else begin
                cpu_fetch_timeout_counter <= 32'h0;
                fetch_error <= 1'b0;
            end
            
            if (dmem_req && !dmem_ready) begin
                cpu_memory_timeout_counter <= cpu_memory_timeout_counter + 1;
                if (cpu_memory_timeout_counter >= `VTX1_TIMEOUT_CYCLES) begin
                    memory_error <= 1'b1;
                    error <= 1'b1;
                    error_code <= `VTX1_ERROR_TIMEOUT;
                    timeout <= 1'b1;
                    error_count <= error_count + 1;
                end            end else begin
                cpu_memory_timeout_counter <= 32'h0;
                memory_error <= 1'b0;
            end
              // Component error aggregation
            if (pipeline_error) begin
                error <= 1'b1;
                error_count <= error_count + 1;
                if (rf_error) begin
                    error_code <= `VTX1_ERROR_INVALID_ADDR;                end else if (tcu_error) begin
                    error_code <= `VTX1_ERROR_INVALID_OP;   // Arithmetic unit errors
                end else if (decoder_error) begin
                    error_code <= `VTX1_ERROR_INVALID_OP;   // Keep original code for decoder errors
                end
            end else if (!fetch_error && !memory_error) begin
                error <= 1'b0;
                error_code <= `VTX1_ERROR_NONE;
                timeout <= 1'b0;
            end
        end
    end
    
    // ========================================================================    // FETCH STAGE
    // ========================================================================
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            pc_fetch <= {`VTX1_WORD_WIDTH{1'b0}};
            instruction_fetch <= {`VTX1_VLIW_WIDTH{1'b0}};
            fetch_valid <= 1'b0;
            imem_req <= 1'b0;
            imem_addr <= {`VTX1_WORD_WIDTH{1'b0}};
        end else if (enable && !stall_pipeline) begin
            if (flush_pipeline) begin
                // Pipeline flush - invalidate fetch stage
                fetch_valid <= 1'b0;
                imem_req <= 1'b0;
            end else if (branch_taken) begin
                // Branch taken - update PC to branch target
                pc_fetch <= branch_target;
                imem_addr <= branch_target;
                imem_req <= 1'b1;
                fetch_valid <= 1'b0;  // Invalidate current fetch
            end else if (imem_ready) begin
                // Normal instruction fetch
                instruction_fetch <= imem_data;
                fetch_valid <= 1'b1;
                imem_req <= 1'b1;
                
                // Increment PC for next instruction (assuming 4-byte VLIW)
                pc_fetch <= pc_fetch + 4;
                imem_addr <= pc_fetch + 4;
            end else begin
                // Memory not ready - maintain request
                imem_req <= 1'b1;
                fetch_valid <= 1'b0;
            end
        end
    end
    
    // ========================================================================    // DECODE STAGE
    // ========================================================================
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            pc_decode <= {`VTX1_WORD_WIDTH{1'b0}};
            instruction_decode <= {`VTX1_VLIW_WIDTH{1'b0}};
            decode_valid <= 1'b0;
        end else if (enable && !stall_pipeline) begin
            if (flush_pipeline) begin
                decode_valid <= 1'b0;
            end else begin
                pc_decode <= pc_fetch;
                instruction_decode <= instruction_fetch;
                decode_valid <= fetch_valid;
            end
        end
    end
    
    // Register file read address assignment
    assign rf_read_addr_a = decoded_src_reg_a;
    assign rf_read_addr_b = decoded_src_reg_b;
    assign rf_read_addr_c = decoded_src_reg_c;
    
    // ========================================================================    // EXECUTE STAGE
    // ========================================================================
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            pc_execute <= {`VTX1_WORD_WIDTH{1'b0}};
            instruction_execute <= {`VTX1_VLIW_WIDTH{1'b0}};
            execute_valid <= 1'b0;
        end else if (enable && !stall_pipeline) begin
            if (flush_pipeline) begin
                execute_valid <= 1'b0;
            end else begin
                pc_execute <= pc_decode;
                instruction_execute <= instruction_decode;
                execute_valid <= decode_valid;
            end
        end
    end    // TCU operand assignment with forwarding
    assign tcu_operand_a = forward_valid_a ? forward_data_a : rf_read_data_a;
    assign tcu_operand_b = forward_valid_b ? forward_data_b : rf_read_data_b;
    assign tcu_operand_c = forward_valid_c ? forward_data_c : rf_read_data_c;
    // TEMPORARY: Force TCU disable for debugging
    // assign tcu_enable = 1'b0; // execute_valid && (decoded_op_type != 3'b000);
    // Restored intended logic for tcu_enable:
    // For legacy TCU, enable on execute_valid and ALU op; for enhanced, driven by microcode logic
    // If both legacy and enhanced are present, this assignment should be removed entirely and tcu_enable should be driven by the appropriate submodule
    assign tcu_enable = execute_valid && (decoded_op_type == 3'b001);
    
    // ========================================================================    // WRITEBACK STAGE
    // ========================================================================
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            pc_writeback <= {`VTX1_WORD_WIDTH{1'b0}};
            instruction_writeback <= {`VTX1_VLIW_WIDTH{1'b0}};
            writeback_valid <= 1'b0;
        end else if (enable && !stall_pipeline) begin
            if (flush_pipeline) begin
                writeback_valid <= 1'b0;
            end else begin
                pc_writeback <= pc_execute;
                instruction_writeback <= instruction_execute;
                writeback_valid <= execute_valid;
            end
        end
    end
    
    // Register file write assignment
    assign rf_write_enable_a = writeback_valid && tcu_valid;    assign rf_write_enable_b = 1'b0;  // Second write port not used in this implementation
    assign rf_write_addr_a = decoded_dst_reg_a;
    assign rf_write_addr_b = decoded_dst_reg_b;    assign rf_write_data_a = tcu_result;
    assign rf_write_data_b = {`VTX1_WORD_WIDTH{1'b0}};
      // ========================================================================
    // PIPELINE CONTROL
    // ========================================================================
    assign stall_pipeline = pipeline_stall_hazard || load_forward_stall || 
                           !imem_ready || !dmem_ready || !tcu_ready;
    assign flush_pipeline = pipeline_flush_hazard || control_hazard || 
                           (debug_enable && debug_cmd == DEBUG_RESET);
    
    // Enhanced hazard-aware pipeline control
    always @(*) begin
        hazard_detected = raw_hazard || waw_hazard || structural_hazard || 
                         memory_hazard || load_use_hazard;
        branch_taken = control_hazard || (debug_enable && debug_cmd == DEBUG_STEP);
        branch_target = pc_fetch + 4; // Simplified branch target calculation
    end

    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            pipeline_stall <= 1'b0;
            pipeline_flush <= 1'b0;
            pipeline_state <= STAGE_FETCH;
        end else begin
            pipeline_stall <= stall_pipeline;
            pipeline_flush <= flush_pipeline;
            
            // Update pipeline state for debugging
            if (writeback_valid)
                pipeline_state <= STAGE_WRITEBACK;
            else if (execute_valid)
                pipeline_state <= STAGE_EXECUTE;
            else if (decode_valid)
                pipeline_state <= STAGE_DECODE;
            else
                pipeline_state <= STAGE_FETCH;
        end
    end
    
    // ========================================================================
    // PERFORMANCE COUNTERS
    // ========================================================================
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            cycle_count <= 32'h0;
            instruction_count <= 32'h0;
            stall_count <= 32'h0;
        end else if (enable) begin
            cycle_count <= cycle_count + 1;
            
            if (writeback_valid && !flush_pipeline)
                instruction_count <= instruction_count + 1;
                
            if (stall_pipeline)
                stall_count <= stall_count + 1;
        end
    end
      // ========================================================================
    // DEBUG INTERFACE
    // ========================================================================
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            debug_pc <= {`VTX1_WORD_WIDTH{1'b0}};
            debug_status <= 4'h0;
        end else begin
            debug_pc <= pc_fetch;
            debug_status <= {pipeline_state, stall_pipeline, flush_pipeline};
        end
    end
    
    // ========================================================================
    // INTERRUPT HANDLING
    // ========================================================================
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            interrupt_ack <= 16'h0;
        end else begin
            // Simple interrupt acknowledgment - can be enhanced
            interrupt_ack <= interrupt_req;
        end
    end

endmodule

`endif // CPU_CORE_V

