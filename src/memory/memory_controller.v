	`timescale 1ns / 1ps
// VTX1 Memory Controller - Simplified Implementation
// Part of the VTX1 Ternary System-on-Chip
// Uses standardized VTX1 interfaces

// Include paths handled by compiler -I flags (see Taskfile.yml)
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module memory_controller (
    input  wire                     clk,
    input  wire                     rst_n,
      // CPU Interface (using VTX1 standardized widths)
    input  wire                     mem_req,
    input  wire                     mem_wr,
    input  wire [1:0]               mem_size,
    input  wire [`VTX1_ADDR_WIDTH-1:0] mem_addr,
    input  wire [`VTX1_WORD_WIDTH-1:0] mem_wdata,
    output reg  [`VTX1_WORD_WIDTH-1:0] mem_rdata,
    output reg                      mem_ready,
    output reg                      mem_error,
    
    // Cache Interface (using VTX1 standardized widths)
    input  wire                     cache_req,
    input  wire                     cache_wr,
    input  wire [`VTX1_ADDR_WIDTH-1:0] cache_addr,
    input  wire [`VTX1_CACHE_LINE_WIDTH-1:0] cache_wdata,
    output reg  [`VTX1_CACHE_LINE_WIDTH-1:0] cache_rdata,
    output reg                      cache_ready,
    output reg                      cache_error,
    
    // Physical Memory Interface (using VTX1 standardized widths)
    output reg  [`VTX1_ADDR_WIDTH-1:0] phy_addr,
    output reg  [`VTX1_WORD_WIDTH-1:0] phy_wdata,
    output reg                      phy_wr,
    output reg                      phy_req,
    input  wire [`VTX1_WORD_WIDTH-1:0] phy_rdata,
    input  wire                     phy_ready,
    input  wire                     phy_error,
    
    // MMIO Interface (using VTX1 standardized widths)
    input  wire                     mmio_req,
    input  wire [`VTX1_ADDR_WIDTH-1:0] mmio_addr,
    input  wire [`VTX1_WORD_WIDTH-1:0] mmio_wdata,
    input  wire                     mmio_wr,
    output reg  [`VTX1_WORD_WIDTH-1:0] mmio_rdata,
    output reg                      mmio_ready,
    output reg                      mmio_error,
    
    // Debug
    output reg  [3:0]               mc_state,
    output reg  [31:0]              access_count,
    output reg  [31:0]              error_count
);    // Memory Controller with Comprehensive Error Handling
    localparam DEFAULT_WORD = `VTX1_WORD_DEFAULT;
    
    // Memory Controller Constants
    localparam MAX_MEMORY_ADDR = 32'h10000000;  // 256MB max memory
    localparam MMIO_BASE_ADDR = 32'hF0000000;   // MMIO base address
    localparam TIMEOUT_CYCLES = `VTX1_TIMEOUT_CYCLES; // Use VTX1 standard timeout
      // Internal registers
    reg [3:0] mc_state_reg;
    reg [31:0] access_count_reg;
    reg [31:0] error_count_reg;
    reg [31:0] mc_timeout_counter;
    reg [3:0] last_error_code;
    reg error_recovery_active;
    
    // Address validation
    wire mem_addr_valid;
    wire cache_addr_valid;
    wire mmio_addr_valid;
    wire timeout_detected;
    wire phy_error_detected;
    
    assign mem_addr_valid = (mem_addr < MAX_MEMORY_ADDR) && (mem_addr[1:0] == 2'b00);
    assign cache_addr_valid = (cache_addr < MAX_MEMORY_ADDR) && (cache_addr[3:0] == 4'b0000);
    assign mmio_addr_valid = (mmio_addr >= MMIO_BASE_ADDR) && (mmio_addr < (MMIO_BASE_ADDR + 16'h100));
    assign timeout_detected = (mc_timeout_counter >= TIMEOUT_CYCLES);
    assign phy_error_detected = phy_error;
    
    // State machine with error handling
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin            mc_state_reg <= `VTX1_MEM_STATE_IDLE;
            access_count_reg <= 32'b0;
            error_count_reg <= 32'b0;
            mc_timeout_counter <= 32'b0;
            last_error_code <= `VTX1_ERROR_NONE;
            error_recovery_active <= 1'b0;
        end else begin
            case (mc_state_reg)                `VTX1_MEM_STATE_IDLE: begin
                    mc_state_reg <= `VTX1_MEM_STATE_READY;
                    mc_timeout_counter <= 32'b0;
                    error_recovery_active <= 1'b0;
                end
                `VTX1_MEM_STATE_READY: begin
                    mc_timeout_counter <= 32'b0;
                    if (mem_req) begin
                        if (mem_addr_valid) begin
                            mc_state_reg <= `VTX1_MEM_STATE_CPU_REQ;
                            access_count_reg <= access_count_reg + 1;
                        end else begin
                            mc_state_reg <= `VTX1_MEM_STATE_ERROR;
                            error_count_reg <= error_count_reg + 1;
                            last_error_code <= `VTX1_ERROR_INVALID_ADDR;
                            error_recovery_active <= 1'b1;
                        end
                    end else if (cache_req) begin
                        if (cache_addr_valid) begin
                            mc_state_reg <= `VTX1_MEM_STATE_CACHE_REQ;
                            access_count_reg <= access_count_reg + 1;
                        end else begin
                            mc_state_reg <= `VTX1_MEM_STATE_ERROR;
                            error_count_reg <= error_count_reg + 1;
                            last_error_code <= `VTX1_ERROR_INVALID_ADDR;
                            error_recovery_active <= 1'b1;
                        end
                    end else if (mmio_req) begin
                        if (mmio_addr_valid) begin
                            mc_state_reg <= `VTX1_MEM_STATE_MMIO_REQ;
                            access_count_reg <= access_count_reg + 1;
                        end else begin
                            mc_state_reg <= `VTX1_MEM_STATE_ERROR;
                            error_count_reg <= error_count_reg + 1;
                            last_error_code <= `VTX1_ERROR_INVALID_ADDR;
                            error_recovery_active <= 1'b1;
                        end
                    end
                end                `VTX1_MEM_STATE_CPU_REQ, `VTX1_MEM_STATE_CACHE_REQ: begin
                    if (phy_ready) begin
                        mc_state_reg <= `VTX1_MEM_STATE_READY;
                        mc_timeout_counter <= 32'b0;
                    end else if (phy_error_detected) begin
                        mc_state_reg <= `VTX1_MEM_STATE_ERROR;
                        error_count_reg <= error_count_reg + 1;
                        last_error_code <= `VTX1_ERROR_BUS_FAULT;
                        error_recovery_active <= 1'b1;
                    end else if (timeout_detected) begin
                        mc_state_reg <= `VTX1_MEM_STATE_TIMEOUT;
                        error_count_reg <= error_count_reg + 1;
                        last_error_code <= `VTX1_ERROR_TIMEOUT;
                        error_recovery_active <= 1'b1;
                    end else begin
                        mc_timeout_counter <= mc_timeout_counter + 1;
                    end
                end                `VTX1_MEM_STATE_MMIO_REQ: begin
                    mc_state_reg <= `VTX1_MEM_STATE_READY;
                    mc_timeout_counter <= 32'b0;
                end
                `VTX1_MEM_STATE_ERROR, `VTX1_MEM_STATE_TIMEOUT: begin
                    if (error_recovery_active) begin
                        if (mc_timeout_counter > 32'd10) begin
                            mc_state_reg <= `VTX1_MEM_STATE_READY;
                            error_recovery_active <= 1'b0;
                            mc_timeout_counter <= 32'b0;
                        end else begin
                            mc_timeout_counter <= mc_timeout_counter + 1;
                        end
                    end
                end
                default: begin
                    mc_state_reg <= `VTX1_MEM_STATE_ERROR;
                    error_count_reg <= error_count_reg + 1;
                    last_error_code <= `VTX1_ERROR_INVALID_STATE;
                    error_recovery_active <= 1'b1;
                end
            endcase
        end
    end
    
    // Combinational logic for memory operations
    always @(*) begin
        // Default values
        phy_req = 1'b0;
        phy_wr = 1'b0;
        phy_addr = DEFAULT_WORD;
        phy_wdata = DEFAULT_WORD;
        
        mem_rdata = DEFAULT_WORD;
        mem_ready = 1'b0;
        mem_error = 1'b0;
        
        cache_rdata = {8{DEFAULT_WORD}};
        cache_ready = 1'b0;
        cache_error = 1'b0;
        
        mmio_rdata = DEFAULT_WORD;
        mmio_ready = 1'b0;
        mmio_error = 1'b0;
        
        case (mc_state_reg)
            `VTX1_MEM_STATE_READY: begin
                // Ready to accept new requests
                if (mmio_req && mmio_addr_valid) begin
                    mmio_ready = 1'b1;
                    mmio_rdata = DEFAULT_WORD;
                end
            end
            `VTX1_MEM_STATE_CPU_REQ: begin
                if (mem_req && mem_addr_valid) begin
                    phy_req = 1'b1;
                    phy_wr = mem_wr;
                    phy_addr = mem_addr;
                    phy_wdata = mem_wdata;
                    mem_rdata = phy_rdata;
                    mem_ready = phy_ready && !phy_error_detected;
                    mem_error = phy_error_detected || timeout_detected;
                end else begin
                    mem_error = 1'b1;
                end
            end
            `VTX1_MEM_STATE_CACHE_REQ: begin
                if (cache_req && cache_addr_valid) begin
                    phy_req = 1'b1;
                    phy_wr = cache_wr;
                    phy_addr = cache_addr;
                    phy_wdata = cache_wdata[`VTX1_WORD_WIDTH-1:0];
                    cache_rdata = {8{phy_rdata}};
                    cache_ready = phy_ready && !phy_error_detected;
                    cache_error = phy_error_detected || timeout_detected;
                end else begin
                    cache_error = 1'b1;
                end
            end
            `VTX1_MEM_STATE_MMIO_REQ: begin
                if (mmio_req && mmio_addr_valid) begin
                    mmio_ready = 1'b1;
                    mmio_rdata = DEFAULT_WORD;
                    mmio_error = 1'b0;
                end else begin
                    mmio_error = 1'b1;
                end
            end
            `VTX1_MEM_STATE_ERROR, `VTX1_MEM_STATE_TIMEOUT: begin
                if (mem_req) mem_error = 1'b1;
                if (cache_req) cache_error = 1'b1;
                if (mmio_req) mmio_error = 1'b1;
            end
            default: begin
                if (mem_req) mem_error = 1'b1;
                if (cache_req) cache_error = 1'b1;
                if (mmio_req) mmio_error = 1'b1;
            end
        endcase
    end
      // Output assignments
    always @(*) begin
        mc_state = mc_state_reg;
        access_count = access_count_reg;
        error_count = error_count_reg;
    end

endmodule

