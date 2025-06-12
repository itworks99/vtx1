	`timescale 1ns / 1ps
// =============================================================================
// VTX1 Bus Matrix Implementation
// =============================================================================
// Comprehensive Bus Matrix System supporting:
// - 3 Masters: CPU, DMA, Debug Interface  
// - Multiple Slaves: Memory Controller, MMIO Router, Cache Controller
// - Round-robin and priority arbitration schemes
// - Deadlock detection and recovery
// - Performance monitoring and timeout handling
// - Full VTX1 interface compliance
// =============================================================================

`ifndef BUS_MATRIX_V
`define BUS_MATRIX_V

// Include VTX1 interface definitions
`include "vtx1_interfaces.v"

// Include VTX1 common infrastructure
`include "vtx1_interfaces.v"
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module bus_matrix (
    input wire clk,
    input wire rst_n,
    
    // =======================================================================
    // MASTER INTERFACES - VTX1 STANDARDIZED
    // =======================================================================
    
    // Master 0: CPU Interface
    input  wire                         cpu_req,
    input  wire                         cpu_wr,
    input  wire [1:0]                   cpu_size,
    input  wire [`VTX1_ADDR_WIDTH-1:0] cpu_addr,
    input  wire [`VTX1_WORD_WIDTH-1:0] cpu_wdata,
    output reg  [`VTX1_WORD_WIDTH-1:0] cpu_rdata,
    output reg                          cpu_ready,
    output reg                          cpu_error,
    output reg  [3:0]                   cpu_error_code,
    output reg                          cpu_timeout,
    input  wire                         cpu_error_clear,
    
    // Master 1: DMA Interface  
    input  wire                         dma_req,
    input  wire                         dma_wr,
    input  wire [1:0]                   dma_size,
    input  wire [`VTX1_ADDR_WIDTH-1:0] dma_addr,
    input  wire [`VTX1_WORD_WIDTH-1:0] dma_wdata,
    output reg  [`VTX1_WORD_WIDTH-1:0] dma_rdata,
    output reg                          dma_ready,
    output reg                          dma_error,
    output reg  [3:0]                   dma_error_code,
    output reg                          dma_timeout,
    input  wire                         dma_error_clear,
    
    // Master 2: Debug Interface
    input  wire                         debug_req,
    input  wire                         debug_wr,
    input  wire [1:0]                   debug_size,
    input  wire [`VTX1_ADDR_WIDTH-1:0] debug_addr,
    input  wire [`VTX1_WORD_WIDTH-1:0] debug_wdata,
    output reg  [`VTX1_WORD_WIDTH-1:0] debug_rdata,
    output reg                          debug_ready,
    output reg                          debug_error,
    output reg  [3:0]                   debug_error_code,
    output reg                          debug_timeout,
    input  wire                         debug_error_clear,
    
    // =======================================================================
    // SLAVE INTERFACES - VTX1 STANDARDIZED  
    // =======================================================================
    
    // Slave 0: Memory Controller Interface
    output reg                          mem_req,
    output reg                          mem_wr,
    output reg  [1:0]                   mem_size,
    output reg  [`VTX1_ADDR_WIDTH-1:0] mem_addr,
    output reg  [`VTX1_WORD_WIDTH-1:0] mem_wdata,
    input  wire [`VTX1_WORD_WIDTH-1:0] mem_rdata,
    input  wire                         mem_ready,
    input  wire                         mem_error,
    
    // Slave 1: MMIO Router Interface
    output reg                          mmio_req,
    output reg                          mmio_wr,
    output reg  [`VTX1_ADDR_WIDTH-1:0] mmio_addr,
    output reg  [`VTX1_WORD_WIDTH-1:0] mmio_wdata,
    input  wire [`VTX1_WORD_WIDTH-1:0] mmio_rdata,
    input  wire                         mmio_ready,
    input  wire                         mmio_error,
    
    // Slave 2: Cache Controller Interface
    output reg                          cache_req,
    output reg                          cache_wr,
    output reg  [`VTX1_ADDR_WIDTH-1:0] cache_addr,
    output reg  [`VTX1_CACHE_LINE_WIDTH-1:0] cache_wdata,
    input  wire [`VTX1_CACHE_LINE_WIDTH-1:0] cache_rdata,
    input  wire                         cache_ready,
    input  wire                         cache_error,
    
    // =======================================================================
    // CONFIGURATION AND CONTROL
    // =======================================================================
    input  wire [1:0]                   arbitration_mode,  // 00=Fixed, 01=Round-Robin, 10=Weighted, 11=Priority
    input  wire [7:0]                   priority_config,   // Priority levels for each master
    input  wire [15:0]                  timeout_config,    // Timeout configuration
    input  wire                         deadlock_enable,   // Enable deadlock detection
    input  wire                         performance_enable, // Enable performance monitoring
    
    // =======================================================================
    // STATUS AND DEBUG OUTPUTS
    // =======================================================================
    output reg  [3:0]                   matrix_state,
    output reg  [2:0]                   current_master,
    output reg  [2:0]                   current_slave,
    output reg                          arbitration_active,
    output reg                          deadlock_detected,
    output reg                          deadlock_recovery,
    
    // Performance Monitoring Counters
    output reg  [31:0]                  total_transactions,
    output reg  [31:0]                  cpu_transactions,
    output reg  [31:0]                  dma_transactions,
    output reg  [31:0]                  debug_transactions,
    output reg  [31:0]                  timeout_count,
    output reg  [31:0]                  error_count,
    output reg  [31:0]                  deadlock_count,
    
    // Bus Utilization Statistics
    output reg  [15:0]                  bus_utilization,    // Percentage * 100
    output reg  [15:0]                  avg_latency,        // Average cycles
    output reg  [15:0]                  max_latency,        // Maximum cycles observed
    
    // Individual Master Statistics
    output reg  [31:0]                  cpu_wait_cycles,
    output reg  [31:0]                  dma_wait_cycles,
    output reg  [31:0]                  debug_wait_cycles
);

// =============================================================================
// PARAMETERS AND CONSTANTS
// =============================================================================

// Bus Matrix States - Using VTX1 Standardized Constants
// Note: Bus matrix uses VTX1_BUS_STATE_* constants for standardization

// Master and Slave Indices
localparam MASTER_CPU           = 3'b000;
localparam MASTER_DMA           = 3'b001;
localparam MASTER_DEBUG         = 3'b010;
localparam MASTER_NONE          = 3'b111;

localparam SLAVE_MEMORY         = 3'b000;
localparam SLAVE_MMIO           = 3'b001;
localparam SLAVE_CACHE          = 3'b010;
localparam SLAVE_NONE           = 3'b111;

// Arbitration Modes
localparam ARB_FIXED_PRIORITY   = 2'b00;
localparam ARB_ROUND_ROBIN      = 2'b01;
localparam ARB_WEIGHTED         = 2'b10;
localparam ARB_PRIORITY_RR      = 2'b11;

// Address Ranges (based on VTX1 memory map)
localparam MMIO_BASE_ADDR       = 36'h400000000;  // 0x4000_0000
localparam MMIO_END_ADDR        = 36'h4FFFFFFFF;  // 0x4FFF_FFFF
localparam CACHE_CONTROL_ADDR   = 36'h500000000;  // 0x5000_0000 (cache control space)

// Timeout and Performance Constants
localparam DEFAULT_TIMEOUT      = 16'd1000;       // Default timeout cycles
localparam DEADLOCK_THRESHOLD   = 16'd500;        // Cycles before deadlock detection
localparam PERFORMANCE_WINDOW   = 32'd10000;      // Performance measurement window

// =============================================================================
// INTERNAL REGISTERS AND WIRES
// =============================================================================

// State Machine
reg [3:0] state_reg, state_next;

// Arbitration Logic
reg [2:0] granted_master;
reg [2:0] next_master;
reg [2:0] last_granted_master;
reg [2:0] round_robin_ptr;

// Request and Grant Signals
wire [2:0] master_requests;
reg  [2:0] master_grants;

// Address Decoding
reg [2:0] selected_slave;
wire [2:0] cpu_slave_select, dma_slave_select, debug_slave_select;

// Timeout Handling
reg [15:0] bus_timeout_counter;
reg [15:0] timeout_threshold;
reg timeout_active;

// Deadlock Detection
reg [15:0] deadlock_counter;
reg [2:0] deadlock_masters;
reg deadlock_state;

// Performance Monitoring
reg [31:0] cycle_counter;
reg [31:0] active_cycles;
reg [31:0] total_requests;
reg [15:0] current_latency;
reg [31:0] latency_accumulator;

// Master Wait Time Tracking
reg [31:0] cpu_wait_start, dma_wait_start, debug_wait_start;
reg cpu_waiting, dma_waiting, debug_waiting;

// Error Tracking
reg [3:0] last_error_code;
reg [2:0] error_master;

// =============================================================================
// COMBINATIONAL LOGIC
// =============================================================================

// Master Request Detection
assign master_requests = {debug_req, dma_req, cpu_req};

// Address Decoding for each Master
assign cpu_slave_select = (cpu_addr >= MMIO_BASE_ADDR && cpu_addr <= MMIO_END_ADDR) ? SLAVE_MMIO :
                          (cpu_addr >= CACHE_CONTROL_ADDR) ? SLAVE_CACHE : SLAVE_MEMORY;

assign dma_slave_select = (dma_addr >= MMIO_BASE_ADDR && dma_addr <= MMIO_END_ADDR) ? SLAVE_MMIO :
                          (dma_addr >= CACHE_CONTROL_ADDR) ? SLAVE_CACHE : SLAVE_MEMORY;

assign debug_slave_select = (debug_addr >= MMIO_BASE_ADDR && debug_addr <= MMIO_END_ADDR) ? SLAVE_MMIO :
                            (debug_addr >= CACHE_CONTROL_ADDR) ? SLAVE_CACHE : SLAVE_MEMORY;

// =============================================================================
// ARBITRATION LOGIC
// =============================================================================

// Priority-based Arbitration (CPU > DMA > Debug)
function [2:0] fixed_priority_arbitration;
    input [2:0] requests;
    begin
        casex (requests)
            3'bxx1: fixed_priority_arbitration = MASTER_CPU;
            3'bx10: fixed_priority_arbitration = MASTER_DMA;
            3'b100: fixed_priority_arbitration = MASTER_DEBUG;
            default: fixed_priority_arbitration = MASTER_NONE;
        endcase
    end
endfunction

// Round-Robin Arbitration
function [2:0] round_robin_arbitration;
    input [2:0] requests;
    input [2:0] last_master;
    begin
        case (last_master)
            MASTER_CPU: begin
                casex (requests)
                    3'bx1x: round_robin_arbitration = MASTER_DMA;
                    3'b1xx: round_robin_arbitration = MASTER_DEBUG;
                    3'bxx1: round_robin_arbitration = MASTER_CPU;
                    default: round_robin_arbitration = MASTER_NONE;
                endcase
            end
            MASTER_DMA: begin
                casex (requests)
                    3'b1xx: round_robin_arbitration = MASTER_DEBUG;
                    3'bxx1: round_robin_arbitration = MASTER_CPU;
                    3'bx1x: round_robin_arbitration = MASTER_DMA;
                    default: round_robin_arbitration = MASTER_NONE;
                endcase
            end
            MASTER_DEBUG: begin
                casex (requests)
                    3'bxx1: round_robin_arbitration = MASTER_CPU;
                    3'bx1x: round_robin_arbitration = MASTER_DMA;
                    3'b1xx: round_robin_arbitration = MASTER_DEBUG;
                    default: round_robin_arbitration = MASTER_NONE;
                endcase
            end
            default: round_robin_arbitration = fixed_priority_arbitration(requests);
        endcase
    end
endfunction

// Weighted Arbitration (CPU:70%, DMA:20%, Debug:10%)
function [2:0] weighted_arbitration;
    input [2:0] requests;
    input [31:0] cycle_count;
    reg [7:0] weight_selection;
    begin
        weight_selection = cycle_count[7:0]; // Use lower 8 bits for selection
        
        casex (requests)
            3'bxx1: begin // CPU request
                if (weight_selection < 8'd179) // 70% = 179/256
                    weighted_arbitration = MASTER_CPU;
                else if (weight_selection < 8'd230 && requests[1]) // 20% = 51/256  
                    weighted_arbitration = MASTER_DMA;
                else if (requests[2]) // 10%
                    weighted_arbitration = MASTER_DEBUG;
                else
                    weighted_arbitration = MASTER_CPU;
            end
            3'bx10: begin // DMA request (no CPU)
                if (weight_selection < 8'd171 && requests[1]) // DMA gets higher weight
                    weighted_arbitration = MASTER_DMA;
                else if (requests[2])
                    weighted_arbitration = MASTER_DEBUG;
                else
                    weighted_arbitration = MASTER_DMA;
            end
            3'b100: weighted_arbitration = MASTER_DEBUG;
            default: weighted_arbitration = MASTER_NONE;
        endcase
    end
endfunction

// Master Arbitration Decision
always @(*) begin
    case (arbitration_mode)
        ARB_FIXED_PRIORITY: next_master = fixed_priority_arbitration(master_requests);
        ARB_ROUND_ROBIN:    next_master = round_robin_arbitration(master_requests, last_granted_master);
        ARB_WEIGHTED:       next_master = weighted_arbitration(master_requests, cycle_counter);
        ARB_PRIORITY_RR:    next_master = (|priority_config) ? 
                                         fixed_priority_arbitration(master_requests) :
                                         round_robin_arbitration(master_requests, last_granted_master);
        default:            next_master = fixed_priority_arbitration(master_requests);
    endcase
end

// =============================================================================
// STATE MACHINE - MAIN BUS MATRIX CONTROL
// =============================================================================

// State Machine Next State Logic
always @(*) begin
    state_next = state_reg;
    case (state_reg)
        `VTX1_BUS_STATE_IDLE: begin
            if (deadlock_detected && deadlock_enable) begin
                state_next = `VTX1_BUS_STATE_DEADLOCK;
            end else if (|master_requests) begin
                state_next = `VTX1_BUS_STATE_ARBITRATE;
            end
        end
        `VTX1_BUS_STATE_ARBITRATE: begin
            case (next_master)
                MASTER_CPU, MASTER_DMA, MASTER_DEBUG: state_next = `VTX1_BUS_STATE_GRANT;
                default: state_next = `VTX1_BUS_STATE_IDLE;
            endcase
        end
        `VTX1_BUS_STATE_GRANT: begin
            // Grant phase: wait for ready/error/timeout from selected master
            if ((granted_master == MASTER_CPU && (cpu_ready || cpu_error || cpu_timeout)) ||
                (granted_master == MASTER_DMA && (dma_ready || dma_error || dma_timeout)) ||
                (granted_master == MASTER_DEBUG && (debug_ready || debug_error || debug_timeout))) begin
                if (|master_requests) begin
                    state_next = `VTX1_BUS_STATE_ARBITRATE;
                end else begin
                    state_next = `VTX1_BUS_STATE_IDLE;
                end
            end else if (bus_timeout_counter >= timeout_threshold) begin
                state_next = `VTX1_BUS_STATE_ERROR;
            end
        end
        `VTX1_BUS_STATE_DEADLOCK: begin
            if (deadlock_counter >= 16'd10) begin
                state_next = `VTX1_BUS_STATE_IDLE;
            end
        end        `VTX1_BUS_STATE_ERROR: begin
            if (bus_timeout_counter >= 16'd50) begin
                state_next = `VTX1_BUS_STATE_IDLE;
            end
        end
        default: state_next = `VTX1_BUS_STATE_IDLE;
    endcase
end

// State Machine Register Update
always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        state_reg <= `VTX1_BUS_STATE_IDLE;
        granted_master <= MASTER_NONE;
        last_granted_master <= MASTER_NONE;        round_robin_ptr <= MASTER_CPU;
        bus_timeout_counter <= 16'b0;
        timeout_threshold <= DEFAULT_TIMEOUT;
        deadlock_counter <= 16'b0;
        deadlock_state <= 1'b0;
    end else begin
        state_reg <= state_next;
        
        // Update granted master
        if (state_reg == `VTX1_BUS_STATE_ARBITRATE) begin
            granted_master <= next_master;
            if (next_master != MASTER_NONE) begin
                last_granted_master <= next_master;
            end
        end
          // Timeout handling
        if (state_reg == `VTX1_BUS_STATE_IDLE || state_reg == `VTX1_BUS_STATE_ARBITRATE) begin
            bus_timeout_counter <= 16'b0;
        end else begin
            bus_timeout_counter <= bus_timeout_counter + 1;
        end
        
        // Update timeout threshold from configuration
        if (timeout_config != 16'b0) begin
            timeout_threshold <= timeout_config;
        end
        
        // Deadlock detection
        if (deadlock_enable) begin
            if (state_reg != `VTX1_BUS_STATE_IDLE && state_next == state_reg && |master_requests) begin
                deadlock_counter <= deadlock_counter + 1;
            end else begin
                deadlock_counter <= 16'b0;
            end
            
            deadlock_state <= (deadlock_counter >= DEADLOCK_THRESHOLD);
        end
    end
end

// =============================================================================
// SLAVE INTERFACE ROUTING
// =============================================================================

// Slave Address Decoding and Request Routing
always @(*) begin
    // Default values
    mem_req = 1'b0;
    mem_wr = 1'b0;
    mem_size = 2'b00;
    mem_addr = {`VTX1_ADDR_WIDTH{1'b0}};
    mem_wdata = {`VTX1_WORD_WIDTH{1'b0}};
    
    mmio_req = 1'b0;
    mmio_wr = 1'b0;
    mmio_addr = {`VTX1_ADDR_WIDTH{1'b0}};
    mmio_wdata = {`VTX1_WORD_WIDTH{1'b0}};
    
    cache_req = 1'b0;
    cache_wr = 1'b0;
    cache_addr = {`VTX1_ADDR_WIDTH{1'b0}};
    cache_wdata = {`VTX1_CACHE_LINE_WIDTH{1'b0}};
    
    selected_slave = SLAVE_NONE;
    if (state_reg == `VTX1_BUS_STATE_GRANT) begin
        case (granted_master)
            MASTER_CPU: begin
                selected_slave = cpu_slave_select;
                case (cpu_slave_select)
                    SLAVE_MEMORY: begin
                        mem_req = cpu_req;
                        mem_wr = cpu_wr;
                        mem_size = cpu_size;
                        mem_addr = cpu_addr;
                        mem_wdata = cpu_wdata;
                    end
                    SLAVE_MMIO: begin
                        mmio_req = cpu_req;
                        mmio_wr = cpu_wr;
                        mmio_addr = cpu_addr;
                        mmio_wdata = cpu_wdata;
                    end
                    SLAVE_CACHE: begin
                        cache_req = cpu_req;
                        cache_wr = cpu_wr;
                        cache_addr = cpu_addr;
                        cache_wdata = {{(`VTX1_CACHE_LINE_WIDTH-`VTX1_WORD_WIDTH){1'b0}}, cpu_wdata};
                    end
                endcase
            end
            MASTER_DMA: begin
                selected_slave = dma_slave_select;
                case (dma_slave_select)
                    SLAVE_MEMORY: begin
                        mem_req = dma_req;
                        mem_wr = dma_wr;
                        mem_size = dma_size;
                        mem_addr = dma_addr;
                        mem_wdata = dma_wdata;
                    end
                    SLAVE_MMIO: begin
                        mmio_req = dma_req;
                        mmio_wr = dma_wr;
                        mmio_addr = dma_addr;
                        mmio_wdata = dma_wdata;
                    end
                    SLAVE_CACHE: begin
                        cache_req = dma_req;
                        cache_wr = dma_wr;
                        cache_addr = dma_addr;
                        cache_wdata = {{(`VTX1_CACHE_LINE_WIDTH-`VTX1_WORD_WIDTH){1'b0}}, dma_wdata};
                    end
                endcase
            end
            MASTER_DEBUG: begin
                selected_slave = debug_slave_select;
                case (debug_slave_select)
                    SLAVE_MEMORY: begin
                        mem_req = debug_req;
                        mem_wr = debug_wr;
                        mem_size = debug_size;
                        mem_addr = debug_addr;
                        mem_wdata = debug_wdata;
                    end
                    SLAVE_MMIO: begin
                        mmio_req = debug_req;
                        mmio_wr = debug_wr;
                        mmio_addr = debug_addr;
                        mmio_wdata = debug_wdata;
                    end
                    SLAVE_CACHE: begin
                        cache_req = debug_req;
                        cache_wr = debug_wr;
                        cache_addr = debug_addr;
                        cache_wdata = {{(`VTX1_CACHE_LINE_WIDTH-`VTX1_WORD_WIDTH){1'b0}}, debug_wdata};
                    end
                endcase
            end
        endcase
    end
end

// =============================================================================
// MASTER RESPONSE ROUTING
// =============================================================================

// Response Routing to Masters
always @(*) begin
    // Default values for all masters
    cpu_rdata = {`VTX1_WORD_WIDTH{1'b0}};
    cpu_ready = 1'b0;
    cpu_error = 1'b0;
    cpu_error_code = `VTX1_ERROR_NONE;
    cpu_timeout = 1'b0;
    
    dma_rdata = {`VTX1_WORD_WIDTH{1'b0}};
    dma_ready = 1'b0;
    dma_error = 1'b0;
    dma_error_code = `VTX1_ERROR_NONE;
    dma_timeout = 1'b0;
    
    debug_rdata = {`VTX1_WORD_WIDTH{1'b0}};
    debug_ready = 1'b0;
    debug_error = 1'b0;
    debug_error_code = `VTX1_ERROR_NONE;
    debug_timeout = 1'b0;
    
    case (state_reg)
        `VTX1_BUS_STATE_GRANT: begin
            case (granted_master)
                MASTER_CPU: begin
                    case (selected_slave)
                        SLAVE_MEMORY: begin
                            cpu_rdata = mem_rdata;
                            cpu_ready = mem_ready;
                            cpu_error = mem_error;
                            cpu_error_code = mem_error ? `VTX1_ERROR_BUS_FAULT : `VTX1_ERROR_NONE;
                        end
                        SLAVE_MMIO: begin
                            cpu_rdata = mmio_rdata;
                            cpu_ready = mmio_ready;
                            cpu_error = mmio_error;
                            cpu_error_code = mmio_error ? `VTX1_ERROR_INVALID_ADDR : `VTX1_ERROR_NONE;
                        end
                        SLAVE_CACHE: begin
                            cpu_rdata = cache_rdata[`VTX1_WORD_WIDTH-1:0];
                            cpu_ready = cache_ready;
                            cpu_error = cache_error;
                            cpu_error_code = cache_error ? `VTX1_ERROR_BUS_FAULT : `VTX1_ERROR_NONE;
                        end                    endcase
                    cpu_timeout = (bus_timeout_counter >= timeout_threshold);
                end
                MASTER_DMA: begin
                    case (selected_slave)
                        SLAVE_MEMORY: begin
                            dma_rdata = mem_rdata;
                            dma_ready = mem_ready;
                            dma_error = mem_error;
                            dma_error_code = mem_error ? `VTX1_ERROR_BUS_FAULT : `VTX1_ERROR_NONE;
                        end
                        SLAVE_MMIO: begin
                            dma_rdata = mmio_rdata;
                            dma_ready = mmio_ready;
                            dma_error = mmio_error;
                            dma_error_code = mmio_error ? `VTX1_ERROR_INVALID_ADDR : `VTX1_ERROR_NONE;
                        end
                        SLAVE_CACHE: begin
                            dma_rdata = cache_rdata[`VTX1_WORD_WIDTH-1:0];
                            dma_ready = cache_ready;
                            dma_error = cache_error;
                            dma_error_code = cache_error ? `VTX1_ERROR_BUS_FAULT : `VTX1_ERROR_NONE;
                        end                    endcase
                    dma_timeout = (bus_timeout_counter >= timeout_threshold);
                end
                MASTER_DEBUG: begin
                    case (selected_slave)
                        SLAVE_MEMORY: begin
                            debug_rdata = mem_rdata;
                            debug_ready = mem_ready;
                            debug_error = mem_error;
                            debug_error_code = mem_error ? `VTX1_ERROR_BUS_FAULT : `VTX1_ERROR_NONE;
                        end
                        SLAVE_MMIO: begin
                            debug_rdata = mmio_rdata;
                            debug_ready = mmio_ready;
                            debug_error = mmio_error;
                            debug_error_code = mmio_error ? `VTX1_ERROR_INVALID_ADDR : `VTX1_ERROR_NONE;
                        end
                        SLAVE_CACHE: begin
                            debug_rdata = cache_rdata[`VTX1_WORD_WIDTH-1:0];
                            debug_ready = cache_ready;
                            debug_error = cache_error;
                            debug_error_code = cache_error ? `VTX1_ERROR_BUS_FAULT : `VTX1_ERROR_NONE;
                        end                    endcase
                    debug_timeout = (bus_timeout_counter >= timeout_threshold);
                end
            endcase
        end
        `VTX1_BUS_STATE_ERROR: begin
            if (cpu_req) begin
                cpu_error = 1'b1;
                cpu_error_code = `VTX1_ERROR_TIMEOUT;
                cpu_timeout = 1'b1;
            end
            if (dma_req) begin
                dma_error = 1'b1;
                dma_error_code = `VTX1_ERROR_TIMEOUT;
                dma_timeout = 1'b1;
            end
            if (debug_req) begin
                debug_error = 1'b1;
                debug_error_code = `VTX1_ERROR_TIMEOUT;
                debug_timeout = 1'b1;
            end
        end
        `VTX1_BUS_STATE_DEADLOCK: begin
            if (cpu_req) begin
                cpu_error = 1'b1;
                cpu_error_code = `VTX1_ERROR_COLLISION;
            end
            if (dma_req) begin
                dma_error = 1'b1;
                dma_error_code = `VTX1_ERROR_COLLISION;
            end
            if (debug_req) begin
                debug_error = 1'b1;
                debug_error_code = `VTX1_ERROR_COLLISION;
            end
        end
    endcase
end

// =============================================================================
// PERFORMANCE MONITORING AND STATISTICS
// =============================================================================

// Performance Counter Updates
always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        cycle_counter <= 32'b0;
        active_cycles <= 32'b0;
        total_transactions <= 32'b0;
        cpu_transactions <= 32'b0;
        dma_transactions <= 32'b0;
        debug_transactions <= 32'b0;
        timeout_count <= 32'b0;
        error_count <= 32'b0;
        deadlock_count <= 32'b0;
        bus_utilization <= 16'b0;
        avg_latency <= 16'b0;
        max_latency <= 16'b0;
        current_latency <= 16'b0;
        latency_accumulator <= 32'b0;
        cpu_wait_cycles <= 32'b0;
        dma_wait_cycles <= 32'b0;
        debug_wait_cycles <= 32'b0;
        cpu_wait_start <= 32'b0;
        dma_wait_start <= 32'b0;
        debug_wait_start <= 32'b0;
        cpu_waiting <= 1'b0;
        dma_waiting <= 1'b0;
        debug_waiting <= 1'b0;
    end else if (performance_enable) begin
        cycle_counter <= cycle_counter + 1;
        if (state_reg != `VTX1_BUS_STATE_IDLE) begin
            active_cycles <= active_cycles + 1;
        end
        if ((state_reg == `VTX1_BUS_STATE_GRANT && granted_master == MASTER_CPU && cpu_ready) ||
            (state_reg == `VTX1_BUS_STATE_GRANT && granted_master == MASTER_DMA && dma_ready) ||
            (state_reg == `VTX1_BUS_STATE_GRANT && granted_master == MASTER_DEBUG && debug_ready)) begin
            total_transactions <= total_transactions + 1;
            case (granted_master)
                MASTER_CPU:   cpu_transactions <= cpu_transactions + 1;
                MASTER_DMA:   dma_transactions <= dma_transactions + 1;
                MASTER_DEBUG: debug_transactions <= debug_transactions + 1;
            endcase
        end
        // Error and timeout counting
        if (cpu_timeout || dma_timeout || debug_timeout) begin
            timeout_count <= timeout_count + 1;
        end
        
        if (cpu_error || dma_error || debug_error) begin
            error_count <= error_count + 1;
        end
        
        if (state_reg == `VTX1_BUS_STATE_DEADLOCK) begin
            deadlock_count <= deadlock_count + 1;
        end
        
        // Wait time tracking
        if (cpu_req && !cpu_waiting) begin
            cpu_wait_start <= cycle_counter;
            cpu_waiting <= 1'b1;
        end else if (cpu_waiting && (cpu_ready || cpu_error || cpu_timeout)) begin
            cpu_wait_cycles <= cpu_wait_cycles + (cycle_counter - cpu_wait_start);
            cpu_waiting <= 1'b0;
        end
        
        if (dma_req && !dma_waiting) begin
            dma_wait_start <= cycle_counter;
            dma_waiting <= 1'b1;
        end else if (dma_waiting && (dma_ready || dma_error || dma_timeout)) begin
            dma_wait_cycles <= dma_wait_cycles + (cycle_counter - dma_wait_start);
            dma_waiting <= 1'b0;
        end
        
        if (debug_req && !debug_waiting) begin
            debug_wait_start <= cycle_counter;
            debug_waiting <= 1'b1;
        end else if (debug_waiting && (debug_ready || debug_error || debug_timeout)) begin
            debug_wait_cycles <= debug_wait_cycles + (cycle_counter - debug_wait_start);
            debug_waiting <= 1'b0;
        end
        
        // Bus utilization calculation (every 1000 cycles)
        if (cycle_counter[9:0] == 10'b0 && cycle_counter != 32'b0) begin
            bus_utilization <= (active_cycles[9:0] * 16'd100) >> 10; // Percentage * 100
        end
        
        // Latency tracking
        if (state_reg == `VTX1_BUS_STATE_ARBITRATE) begin
            current_latency <= 16'd1;
        end else if (state_reg != `VTX1_BUS_STATE_IDLE && state_reg == state_next) begin
            current_latency <= current_latency + 1;
        end else if ((state_reg == `VTX1_BUS_STATE_CPU_ACTIVE || state_reg == `VTX1_BUS_STATE_DMA_ACTIVE || 
                     state_reg == `VTX1_BUS_STATE_DEBUG_ACTIVE) && 
                    (cpu_ready || dma_ready || debug_ready)) begin
            latency_accumulator <= latency_accumulator + current_latency;
            if (current_latency > max_latency) begin
                max_latency <= current_latency;
            end
            // Calculate average latency
            if (total_transactions != 32'b0) begin
                avg_latency <= latency_accumulator[31:16] / total_transactions[15:0];
            end
        end
    end
end

// =============================================================================
// STATUS OUTPUTS
// =============================================================================

// Status Output Assignments
always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        matrix_state <= `VTX1_BUS_STATE_IDLE;
        current_master <= MASTER_NONE;
        current_slave <= SLAVE_NONE;
        arbitration_active <= 1'b0;
        deadlock_detected <= 1'b0;
        deadlock_recovery <= 1'b0;
    end else begin
        matrix_state <= state_reg;
        current_master <= granted_master;
        current_slave <= selected_slave;
        arbitration_active <= (state_reg == `VTX1_BUS_STATE_ARBITRATE);
        deadlock_detected <= deadlock_state;
        deadlock_recovery <= (state_reg == `VTX1_BUS_STATE_DEADLOCK);
    end
end

endmodule

`endif // BUS_MATRIX_V

