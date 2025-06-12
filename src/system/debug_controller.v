	`timescale 1ns / 1ps
// =============================================================================
// VTX1 Debug Interface Controller
// =============================================================================
// Debug interface for multi-master bus matrix integration
// Provides basic debug access to system memory and registers
// Simplified JTAG-compatible interface for development
// =============================================================================

`ifndef DEBUG_CONTROLLER_V
`define DEBUG_CONTROLLER_V

// Include VTX1 interface definitions
`include "vtx1_interfaces.v"
// Include VTX1 common infrastructure
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module debug_controller (
    input wire clk,
    input wire rst_n,
    
    // =======================================================================
    // BUS MATRIX INTERFACE - VTX1 STANDARDIZED
    // =======================================================================
    output reg                          debug_req,
    output reg                          debug_wr,
    output reg  [1:0]                   debug_size,
    output reg  [`VTX1_ADDR_WIDTH-1:0] debug_addr,
    output reg  [`VTX1_WORD_WIDTH-1:0] debug_wdata,
    input  wire [`VTX1_WORD_WIDTH-1:0] debug_rdata,
    input  wire                         debug_ready,
    input  wire                         debug_error,
    input  wire [3:0]                   debug_error_code,
    input  wire                         debug_timeout,
    output reg                          debug_error_clear,
    
    // =======================================================================
    // JTAG INTERFACE (Simplified)
    // =======================================================================
    input  wire                         jtag_tck,
    input  wire                         jtag_tms,
    input  wire                         jtag_tdi,
    output reg                          jtag_tdo,
    input  wire                         jtag_trst_n,
    
    // =======================================================================
    // DEBUG CONTROL INTERFACE
    // =======================================================================
    input  wire                         debug_enable,
    input  wire                         debug_halt_cpu,
    input  wire                         debug_step_cpu,
    input  wire                         debug_reset_cpu,
    input  wire [`VTX1_ADDR_WIDTH-1:0] debug_access_addr,
    input  wire [`VTX1_WORD_WIDTH-1:0] debug_access_wdata,
    input  wire                         debug_access_wr,
    input  wire                         debug_access_req,
    
    // =======================================================================
    // CPU DEBUG INTERFACE
    // =======================================================================
    output reg                          cpu_debug_halt,
    output reg                          cpu_debug_step,
    output reg                          cpu_debug_reset,
    input  wire [`VTX1_ADDR_WIDTH-1:0] cpu_debug_pc,
    input  wire [3:0]                   cpu_debug_state,
    input  wire [`VTX1_WORD_WIDTH-1:0] cpu_debug_reg_data,
    output reg  [4:0]                   cpu_debug_reg_addr,
    output reg                          cpu_debug_reg_req,
      // =======================================================================
    // STATUS AND INTERRUPT
    // =======================================================================
    output reg                          debug_halted,
    output reg                          debug_break_hit,
    output reg                          debug_irq,
    output reg  [7:0]                   debug_status,
    output reg  [`VTX1_ADDR_WIDTH-1:0] debug_last_access_addr,
    output reg  [`VTX1_WORD_WIDTH-1:0] debug_last_read_data,
    
    // =======================================================================
    // MISSING PORTS EXPECTED BY TOP LEVEL
    // =======================================================================
    output reg  [`VTX1_WORD_WIDTH-1:0] debug_access_data,
    output reg                          debug_access_ready,
    output reg  [3:0]                   debug_state,
    output reg  [31:0]                  debug_operations_count,
    
    // =======================================================================
    // DEBUG CONFIGURATION
    // =======================================================================    input  wire [15:0]                  breakpoint_enable,
    input  wire [`VTX1_ADDR_WIDTH*16-1:0] breakpoint_addr_packed,
    input  wire [15:0]                  watchpoint_enable,
    input  wire [`VTX1_ADDR_WIDTH*16-1:0] watchpoint_addr_packed,
    
    // =======================================================================
    // DEBUG STATE OUTPUT
    // =======================================================================
    output reg  [3:0]                   debug_controller_state,
    output reg  [31:0]                  debug_access_count,
    output reg  [7:0]                   debug_error_count
);

// =============================================================================
// DEBUG CONTROLLER PARAMETERS
// =============================================================================
// Remove custom DBG_* state localparams - using VTX1 constants

// JTAG TAP States (simplified)
localparam JTAG_IDLE        = 4'h0;
localparam JTAG_SHIFT_DR    = 4'h1;
localparam JTAG_SHIFT_IR    = 4'h2;
localparam JTAG_UPDATE_DR   = 4'h3;

// =============================================================================
// INTERNAL REGISTERS
// =============================================================================
reg [3:0] debug_state_reg, debug_state_next;
reg [3:0] jtag_state_reg, jtag_state_next;

// Debug Access Control
reg [`VTX1_ADDR_WIDTH-1:0] pending_addr;
reg [`VTX1_WORD_WIDTH-1:0] pending_wdata;
reg pending_wr;
reg access_pending;

// JTAG Shift Registers
reg [31:0] jtag_dr; // Data register
reg [7:0] jtag_ir;  // Instruction register
reg [5:0] jtag_bit_count;

// Breakpoint and Watchpoint Detection
reg [15:0] breakpoint_hit;
reg [15:0] watchpoint_hit;
reg [`VTX1_ADDR_WIDTH-1:0] last_cpu_pc;

// Missing variable declarations
reg [15:0] breakpoint_enable;
reg [`VTX1_ADDR_WIDTH-1:0] breakpoint_addr [0:15];
reg [`VTX1_ADDR_WIDTH-1:0] watchpoint_addr [0:15];

// Performance and Error Tracking
reg [31:0] access_count_reg;
reg [7:0] debug_error_count_reg;

// CPU Control State
reg cpu_halted_state;
reg step_pending;

// =============================================================================
// STATE MACHINE - MAIN DEBUG CONTROLLER
// =============================================================================

// Next State Logic
always @(*) begin
    debug_state_next = debug_state_reg;
    
    case (debug_state_reg)
        `VTX1_DEBUG_STATE_IDLE: begin            if (!debug_enable) begin
                debug_state_next = `VTX1_DEBUG_STATE_WAIT_ENABLE;
            end else if (|breakpoint_hit && |breakpoint_enable) begin
                debug_state_next = `VTX1_DEBUG_STATE_BREAKPOINT;
            end else if (|watchpoint_hit && |watchpoint_enable) begin
                debug_state_next = `VTX1_DEBUG_STATE_WATCHPOINT;
            end else if (debug_halt_cpu) begin
                debug_state_next = `VTX1_DEBUG_STATE_HALT;
            end else if (debug_access_req || access_pending) begin
                debug_state_next = `VTX1_DEBUG_STATE_ACCESS;
            end
        end        `VTX1_DEBUG_STATE_WAIT_ENABLE: begin
            if (debug_enable) begin
                debug_state_next = `VTX1_DEBUG_STATE_IDLE;
            end
        end
        `VTX1_DEBUG_STATE_ACCESS: begin
            if (pending_wr) begin
                debug_state_next = `VTX1_DEBUG_STATE_WRITE;
            end else begin
                debug_state_next = `VTX1_DEBUG_STATE_READ;
            end
        end        `VTX1_DEBUG_STATE_READ: begin
            if (debug_ready || debug_error || debug_timeout) begin
                debug_state_next = `VTX1_DEBUG_STATE_IDLE;
            end
        end
        `VTX1_DEBUG_STATE_WRITE: begin
            if (debug_ready || debug_error || debug_timeout) begin
                debug_state_next = `VTX1_DEBUG_STATE_IDLE;
            end
        end
        `VTX1_DEBUG_STATE_BREAKPOINT: begin
            if (!debug_enable || debug_step_cpu) begin
                debug_state_next = `VTX1_DEBUG_STATE_IDLE;
            end
        end
        `VTX1_DEBUG_STATE_WATCHPOINT: begin
            if (!debug_enable || debug_step_cpu) begin
                debug_state_next = `VTX1_DEBUG_STATE_IDLE;
            end
        end
        `VTX1_DEBUG_STATE_HALT: begin
            if (!debug_halt_cpu || debug_step_cpu) begin
                debug_state_next = `VTX1_DEBUG_STATE_IDLE;
            end
        end
        `VTX1_DEBUG_STATE_ERROR: begin
            if (debug_error_clear) begin
                debug_state_next = `VTX1_DEBUG_STATE_IDLE;
            end
        end
        default: debug_state_next = `VTX1_DEBUG_STATE_IDLE;
    endcase
end

// State Register Update
always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        debug_state_reg <= `VTX1_DEBUG_STATE_WAIT_ENABLE;
        pending_addr <= {`VTX1_ADDR_WIDTH{1'b0}};
        pending_wdata <= {`VTX1_WORD_WIDTH{1'b0}};
        pending_wr <= 1'b0;
        access_pending <= 1'b0;
        access_count_reg <= 32'b0;
        debug_error_count_reg <= 8'h00;
        cpu_halted_state <= 1'b0;
        step_pending <= 1'b0;
        last_cpu_pc <= {`VTX1_ADDR_WIDTH{1'b0}};
    end else begin
        debug_state_reg <= debug_state_next;
        
        case (debug_state_reg)
            `VTX1_DEBUG_STATE_IDLE: begin
                if (debug_access_req) begin
                    pending_addr <= debug_access_addr;
                    pending_wdata <= debug_access_wdata;
                    pending_wr <= debug_access_wr;
                    access_pending <= 1'b1;
                end
            end
            `VTX1_DEBUG_STATE_ACCESS: begin
                access_pending <= 1'b0;
            end
            `VTX1_DEBUG_STATE_READ, `VTX1_DEBUG_STATE_WRITE: begin
                if (debug_ready && !debug_error) begin
                    access_count_reg <= access_count_reg + 1;
                end else if (debug_error || debug_timeout) begin
                    debug_error_count_reg <= debug_error_count_reg + 1;
                end
            end
            `VTX1_DEBUG_STATE_HALT: begin
                cpu_halted_state <= 1'b1;
            end            default: begin
                if (debug_state_next == `VTX1_DEBUG_STATE_IDLE) begin
                    cpu_halted_state <= 1'b0;
                end
            end
        endcase
        
        // Update last PC for breakpoint detection
        last_cpu_pc <= cpu_debug_pc;
        
        // Single step control
        if (debug_step_cpu) begin
            step_pending <= 1'b1;
        end else if (cpu_debug_pc != last_cpu_pc) begin
            step_pending <= 1'b0;
        end
    end
end

// =============================================================================
// JTAG TAP CONTROLLER (Simplified)
// =============================================================================

// JTAG State Machine
always @(posedge jtag_tck or negedge jtag_trst_n) begin
    if (!jtag_trst_n) begin
        jtag_state_reg <= JTAG_IDLE;
        jtag_dr <= 32'h00000000;
        jtag_ir <= 8'h00;
        jtag_bit_count <= 6'h00;
    end else begin
        case (jtag_state_reg)
            JTAG_IDLE: begin
                if (!jtag_tms) begin
                    jtag_state_reg <= JTAG_SHIFT_DR;
                    jtag_bit_count <= 6'h00;
                end
            end
            
            JTAG_SHIFT_DR: begin
                if (jtag_bit_count < 32) begin
                    jtag_dr <= {jtag_tdi, jtag_dr[31:1]};
                    jtag_bit_count <= jtag_bit_count + 1;
                end
                if (jtag_tms) begin
                    jtag_state_reg <= JTAG_UPDATE_DR;
                end
            end
            
            JTAG_UPDATE_DR: begin
                jtag_state_reg <= JTAG_IDLE;
                // Process JTAG command (simplified)
            end
            
            default: jtag_state_reg <= JTAG_IDLE;
        endcase
    end
end

// =============================================================================
// BREAKPOINT AND WATCHPOINT DETECTION
// =============================================================================

// Breakpoint Detection
integer bp_i;
always @(*) begin
    breakpoint_hit = 16'h0000;
    for (bp_i = 0; bp_i < 16; bp_i = bp_i + 1) begin
        if (breakpoint_enable[bp_i] && (cpu_debug_pc == breakpoint_addr[bp_i])) begin
            breakpoint_hit[bp_i] = 1'b1;
        end
    end
end

// Watchpoint Detection (on debug access)
integer wp_i;
always @(*) begin
    watchpoint_hit = 16'h0000;
    for (wp_i = 0; wp_i < 16; wp_i = wp_i + 1) begin
        if (watchpoint_enable[wp_i] && (pending_addr == watchpoint_addr[wp_i])) begin
            watchpoint_hit[wp_i] = 1'b1;
        end
    end
end

// =============================================================================
// OUTPUT CONTROL LOGIC
// =============================================================================

// Bus Interface Control
always @(*) begin
    // Default values
    debug_req = 1'b0;
    debug_wr = 1'b0;
    debug_size = 2'b10; // Word size
    debug_addr = {`VTX1_ADDR_WIDTH{1'b0}};
    debug_wdata = {`VTX1_WORD_WIDTH{1'b0}};
    debug_error_clear = 1'b0;
      case (debug_state_reg)
        `VTX1_DEBUG_STATE_ACCESS: begin
            debug_req = 1'b1;
            debug_wr = pending_wr;
            debug_addr = pending_addr;
            debug_wdata = pending_wdata;
        end
        `VTX1_DEBUG_STATE_READ: begin
            debug_req = 1'b1;
            debug_wr = 1'b0;
            debug_addr = pending_addr;
        end
        `VTX1_DEBUG_STATE_WRITE: begin
            debug_req = 1'b1;
            debug_wr = 1'b1;
            debug_addr = pending_addr;
            debug_wdata = pending_wdata;
        end
        `VTX1_DEBUG_STATE_ERROR: begin
            debug_error_clear = 1'b1;
        end
    endcase
end

// CPU Control Outputs
always @(*) begin
    cpu_debug_halt = (debug_state_reg == `VTX1_DEBUG_STATE_HALT) || (debug_state_reg == `VTX1_DEBUG_STATE_BREAKPOINT);
    cpu_debug_step = debug_step_cpu && step_pending;
    cpu_debug_reset = debug_reset_cpu;
    
    // Register access (placeholder)
    cpu_debug_reg_addr = 5'h00;
    cpu_debug_reg_req = 1'b0;
end

// Status and Control Outputs
always @(*) begin
    // Status outputs
    debug_halted = cpu_halted_state || (|breakpoint_hit) || (|watchpoint_hit);
    debug_break_hit = |breakpoint_hit;
    debug_irq = debug_break_hit || (|watchpoint_hit);
      case (debug_state_reg)
        `VTX1_DEBUG_STATE_IDLE:         debug_status = 8'h00; // Idle
        `VTX1_DEBUG_STATE_WAIT_ENABLE:  debug_status = 8'h01; // Waiting for enable
        `VTX1_DEBUG_STATE_ACCESS:       debug_status = 8'h02; // Accessing bus
        `VTX1_DEBUG_STATE_READ:         debug_status = 8'h03; // Reading
        `VTX1_DEBUG_STATE_WRITE:        debug_status = 8'h04; // Writing
        `VTX1_DEBUG_STATE_BREAKPOINT:   debug_status = 8'h05; // Breakpoint hit
        `VTX1_DEBUG_STATE_WATCHPOINT:   debug_status = 8'h06; // Watchpoint hit
        `VTX1_DEBUG_STATE_HALT:         debug_status = 8'h07; // CPU halted
        `VTX1_DEBUG_STATE_ERROR:        debug_status = 8'hFF; // Error
        default:                        debug_status = 8'h00;
    endcase
    
    // Debug state and statistics
    debug_controller_state = debug_state_reg;
    debug_access_count = access_count_reg;
    debug_error_count = debug_error_count_reg;
    
    // Last access tracking
    debug_last_access_addr = pending_addr;
    debug_last_read_data = debug_rdata;
end

// JTAG Output
always @(*) begin
    jtag_tdo = jtag_dr[0]; // Shift out LSB
end

endmodule

`endif // DEBUG_CONTROLLER_V

