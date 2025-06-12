	`timescale 1ns / 1ps
// =============================================================================
// VTX1 DMA Controller
// =============================================================================
// Simplified DMA controller for multi-master bus matrix integration
// Provides basic memory-to-memory and memory-to-peripheral transfers
// =============================================================================

`ifndef DMA_CONTROLLER_V
`define DMA_CONTROLLER_V

// Include VTX1 interface definitions
`include "vtx1_interfaces.v"
// Include VTX1 common infrastructure
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module dma_controller (
    input wire clk,
    input wire rst_n,
    
    // =======================================================================
    // BUS MATRIX INTERFACE - VTX1 STANDARDIZED
    // =======================================================================
    output reg                          dma_req,
    output reg                          dma_wr,
    output reg  [1:0]                   dma_size,
    output reg  [`VTX1_ADDR_WIDTH-1:0] dma_addr,
    output reg  [`VTX1_WORD_WIDTH-1:0] dma_wdata,
    input  wire [`VTX1_WORD_WIDTH-1:0] dma_rdata,
    input  wire                         dma_ready,
    input  wire                         dma_error,
    input  wire [3:0]                   dma_error_code,
    input  wire                         dma_timeout,
    output reg                          dma_error_clear,
    
    // =======================================================================
    // CPU CONTROL INTERFACE
    // =======================================================================
    input  wire                         dma_enable,
    input  wire                         dma_start,
    input  wire [7:0]                   dma_channel,
    input  wire [`VTX1_ADDR_WIDTH-1:0] dma_src_addr,
    input  wire [`VTX1_ADDR_WIDTH-1:0] dma_dest_addr,
    input  wire [31:0]                  dma_transfer_count,
    input  wire [2:0]                   dma_transfer_mode,
    
    // =======================================================================
    // STATUS AND INTERRUPT
    // =======================================================================
    output reg                          dma_complete,
    output reg                          dma_irq,
    output reg  [7:0]                   dma_status,
    output reg  [31:0]                  dma_transfer_progress,
    
    // =======================================================================
    // DEBUG
    // =======================================================================
    output reg  [3:0]                   dma_state,
    output reg  [31:0]                  dma_operations_count,
    output reg  [7:0]                   dma_error_count
);

// =============================================================================
// DMA CONTROLLER PARAMETERS
// =============================================================================
// Remove custom DMA_* state localparams - using VTX1 constants

// Transfer Modes
localparam MODE_MEM_TO_MEM  = 3'b000;
localparam MODE_MEM_TO_PERI = 3'b001;
localparam MODE_PERI_TO_MEM = 3'b010;
localparam MODE_BURST       = 3'b011;

// =============================================================================
// INTERNAL REGISTERS
// =============================================================================
reg [3:0] state_reg, state_next;
reg [`VTX1_ADDR_WIDTH-1:0] current_src_addr;
reg [`VTX1_ADDR_WIDTH-1:0] current_dest_addr;
reg [31:0] current_count;
reg [`VTX1_WORD_WIDTH-1:0] transfer_data;
reg [2:0] current_mode;
reg [7:0] current_channel;
reg transfer_active;
reg read_phase;

// Performance and Error Tracking
reg [31:0] operations_count_reg;
reg [7:0] error_count_reg;
reg [31:0] cycles_since_start;

// =============================================================================
// STATE MACHINE
// =============================================================================

// Next State Logic
always @(*) begin
    state_next = state_reg;
    
    case (state_reg)
        `VTX1_DMA_STATE_IDLE: begin
            if (dma_enable && dma_start && current_count > 0) begin
                state_next = `VTX1_DMA_STATE_SETUP;
            end
        end
        `VTX1_DMA_STATE_SETUP: begin
            state_next = `VTX1_DMA_STATE_READ;
        end
        `VTX1_DMA_STATE_READ: begin
            if (dma_ready && !dma_error) begin
                state_next = `VTX1_DMA_STATE_WRITE;
            end else if (dma_error || dma_timeout) begin
                state_next = `VTX1_DMA_STATE_ERROR;
            end
        end
        `VTX1_DMA_STATE_WRITE: begin
            if (dma_ready && !dma_error) begin
                if (current_count <= 1) begin
                    state_next = `VTX1_DMA_STATE_COMPLETE;
                end else begin
                    state_next = `VTX1_DMA_STATE_READ;
                end
            end else if (dma_error || dma_timeout) begin
                state_next = `VTX1_DMA_STATE_ERROR;
            end
        end
        `VTX1_DMA_STATE_COMPLETE: begin
            state_next = `VTX1_DMA_STATE_IDLE;
        end
        `VTX1_DMA_STATE_ERROR: begin
            if (dma_error_clear) begin
                state_next = `VTX1_DMA_STATE_IDLE;
            end
        end
        default: state_next = `VTX1_DMA_STATE_IDLE;
    endcase
end

// State Register Update
always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        state_reg <= `VTX1_DMA_STATE_IDLE;
        current_src_addr <= {`VTX1_ADDR_WIDTH{1'b0}};
        current_dest_addr <= {`VTX1_ADDR_WIDTH{1'b0}};
        current_count <= 32'b0;
        transfer_data <= {`VTX1_WORD_WIDTH{1'b0}};
        current_mode <= MODE_MEM_TO_MEM;
        current_channel <= 8'h00;
        transfer_active <= 1'b0;
        read_phase <= 1'b1;
        operations_count_reg <= 32'b0;
        error_count_reg <= 8'h00;
        cycles_since_start <= 32'b0;
    end else begin
        state_reg <= state_next;
        
        case (state_reg)
            `VTX1_DMA_STATE_IDLE: begin
                if (dma_enable && dma_start) begin
                    current_src_addr <= dma_src_addr;
                    current_dest_addr <= dma_dest_addr;
                    current_count <= dma_transfer_count;
                    current_mode <= dma_transfer_mode;
                    current_channel <= dma_channel;
                    transfer_active <= 1'b1;
                    read_phase <= 1'b1;
                    cycles_since_start <= 32'b0;
                end else begin
                    transfer_active <= 1'b0;
                end
            end
            `VTX1_DMA_STATE_READ: begin
                if (dma_ready && !dma_error) begin
                    transfer_data <= dma_rdata;
                    read_phase <= 1'b0;
                end
            end
            `VTX1_DMA_STATE_WRITE: begin
                if (dma_ready && !dma_error) begin
                    current_src_addr <= current_src_addr + 1;
                    current_dest_addr <= current_dest_addr + 1;
                    current_count <= current_count - 1;
                    operations_count_reg <= operations_count_reg + 1;
                    read_phase <= 1'b1;
                end
            end
            `VTX1_DMA_STATE_ERROR: begin
                error_count_reg <= error_count_reg + 1;
                transfer_active <= 1'b0;
            end
            `VTX1_DMA_STATE_COMPLETE: begin
                transfer_active <= 1'b0;
            end
        endcase
        
        if (transfer_active) begin
            cycles_since_start <= cycles_since_start + 1;
        end
    end
end

// =============================================================================
// OUTPUT CONTROL LOGIC
// =============================================================================

// Bus Interface Control
always @(*) begin
    // Default values
    dma_req = 1'b0;
    dma_wr = 1'b0;
    dma_size = 2'b10; // Word size
    dma_addr = {`VTX1_ADDR_WIDTH{1'b0}};
    dma_wdata = {`VTX1_WORD_WIDTH{1'b0}};
    dma_error_clear = 1'b0;
    
    case (state_reg)
        `VTX1_DMA_STATE_READ: begin
            dma_req = 1'b1;
            dma_wr = 1'b0; // Read operation
            dma_addr = current_src_addr;
        end
        `VTX1_DMA_STATE_WRITE: begin
            dma_req = 1'b1;
            dma_wr = 1'b1; // Write operation
            dma_addr = current_dest_addr;
            dma_wdata = transfer_data;
        end
        `VTX1_DMA_STATE_ERROR: begin
            dma_error_clear = 1'b1;
        end
    endcase
end

// Status and Control Outputs
always @(*) begin
    // Status outputs
    dma_complete = (state_reg == `VTX1_DMA_STATE_COMPLETE);
    dma_irq = (state_reg == `VTX1_DMA_STATE_COMPLETE) || (state_reg == `VTX1_DMA_STATE_ERROR);
    
    case (state_reg)
        `VTX1_DMA_STATE_IDLE:     dma_status = 8'h00; // Idle
        `VTX1_DMA_STATE_SETUP:    dma_status = 8'h01; // Setting up
        `VTX1_DMA_STATE_READ:     dma_status = 8'h02; // Reading source
        `VTX1_DMA_STATE_WRITE:    dma_status = 8'h03; // Writing destination
        `VTX1_DMA_STATE_COMPLETE: dma_status = 8'h04; // Transfer complete
        `VTX1_DMA_STATE_ERROR:    dma_status = 8'hFF; // Error
        default:                  dma_status = 8'h00;
    endcase
    
    // Transfer progress
    if (dma_transfer_count > 0) begin
        dma_transfer_progress = dma_transfer_count - current_count;
    end else begin
        dma_transfer_progress = 32'b0;
    end
    
    // Debug outputs
    dma_state = state_reg;
    dma_operations_count = operations_count_reg;
    dma_error_count = error_count_reg;
end

endmodule

`endif // DMA_CONTROLLER_V

