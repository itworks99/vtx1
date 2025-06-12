	`timescale 1ns / 1ps
// =============================================================================
// VTX1 CPU Adapter Interface
// =============================================================================
// Bridges CPU core direct memory interface to bus matrix master interface
// Converts bidirectional data bus to separate read/write data paths
// Provides standardized VTX1 bus matrix interface
// =============================================================================

`ifndef CPU_ADAPTER_V
`define CPU_ADAPTER_V

// Include VTX1 interface definitions
`include "vtx1_interfaces.v"
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module cpu_adapter (
    input wire clk,
    input wire rst_n,
    
    // =======================================================================
    // CPU CORE INTERFACE (Direct Memory Interface)
    // =======================================================================
    input  wire [`VTX1_WORD_WIDTH-1:0]  cpu_dmem_addr,
    inout  wire [`VTX1_WORD_WIDTH-1:0]  cpu_dmem_data,
    input  wire                          cpu_dmem_we,
    input  wire                          cpu_dmem_oe,
    input  wire                          cpu_dmem_req,
    output reg                           cpu_dmem_ready,
    
    // Additional CPU interface signals
    input  wire [`VTX1_WORD_WIDTH-1:0]  cpu_imem_addr,
    input  wire                          cpu_imem_req,
    output reg                           cpu_imem_ready,
    
    // =======================================================================
    // BUS MATRIX INTERFACE (Standardized Master Interface)
    // =======================================================================
    output reg                           bus_req,
    output reg                           bus_wr,
    output reg  [1:0]                    bus_size,
    output reg  [`VTX1_ADDR_WIDTH-1:0]  bus_addr,
    output reg  [`VTX1_WORD_WIDTH-1:0]  bus_wdata,
    input  wire [`VTX1_WORD_WIDTH-1:0]  bus_rdata,
    input  wire                          bus_ready,
    input  wire                          bus_error,
    input  wire [3:0]                    bus_error_code,
    input  wire                          bus_timeout,
    output reg                           bus_error_clear,
    
    // =======================================================================
    // STATUS AND DEBUG
    // =======================================================================
    output reg  [2:0]                    adapter_state,
    output reg  [31:0]                   data_transactions,
    output reg  [31:0]                   instruction_transactions,
    output reg  [31:0]                   error_count,
    output reg                           adapter_error,
    output reg  [3:0]                    adapter_error_code
);

// =============================================================================
// ADAPTER STATE MACHINE - Use VTX1 standardized states
// =============================================================================
localparam STATE_IDLE           = `VTX1_STATE_IDLE;
localparam STATE_DATA_REQUEST   = `VTX1_STATE_ACTIVE;
localparam STATE_DATA_WAIT      = `VTX1_STATE_WAIT;
localparam STATE_INSTR_REQUEST  = `VTX1_STATE_EXECUTE;  // Map instruction handling to execute
localparam STATE_INSTR_WAIT     = `VTX1_STATE_VALIDATE; // Map instruction wait to validate
localparam STATE_ERROR          = `VTX1_STATE_ERROR;

// =============================================================================
// INTERNAL REGISTERS
// =============================================================================
reg [2:0] state_reg, state_next;
reg [`VTX1_WORD_WIDTH-1:0] data_write_reg;
reg [`VTX1_WORD_WIDTH-1:0] data_read_reg;
reg data_request_pending;
reg instruction_request_pending;
reg [31:0] data_transactions_reg;
reg [31:0] instruction_transactions_reg;
reg [31:0] error_count_reg;
reg [3:0] adapter_timeout_counter;

// VTX1 Error Handling Variables
reg [3:0] vtx1_error_reg;
reg [31:0] vtx1_error_info;
reg vtx1_error_valid;

// Bidirectional data bus control
reg dmem_data_oe;
wire [`VTX1_WORD_WIDTH-1:0] dmem_data_in;
reg [`VTX1_WORD_WIDTH-1:0] dmem_data_out;

assign cpu_dmem_data = dmem_data_oe ? dmem_data_out : {`VTX1_WORD_WIDTH{1'bz}};
assign dmem_data_in = cpu_dmem_data;

// =============================================================================
// STATE MACHINE LOGIC
// =============================================================================

// Next state logic
always @(*) begin
    state_next = state_reg;
    
    case (state_reg)
        STATE_IDLE: begin
            if (cpu_dmem_req) begin
                state_next = STATE_DATA_REQUEST;
            end else if (cpu_imem_req) begin
                state_next = STATE_INSTR_REQUEST;
            end
        end
        
        STATE_DATA_REQUEST: begin
            if (bus_ready) begin
                state_next = STATE_IDLE;
            end else if (bus_error || bus_timeout) begin
                state_next = STATE_ERROR;
            end else begin
                state_next = STATE_DATA_WAIT;
            end
        end
          STATE_DATA_WAIT: begin
            if (bus_ready) begin
                state_next = STATE_IDLE;
            end else if (bus_error || bus_timeout || adapter_timeout_counter >= 4'hF) begin
                state_next = STATE_ERROR;
            end
        end
        
        STATE_INSTR_REQUEST: begin
            if (bus_ready) begin
                state_next = STATE_IDLE;
            end else if (bus_error || bus_timeout) begin
                state_next = STATE_ERROR;
            end else begin
                state_next = STATE_INSTR_WAIT;
            end
        end
          STATE_INSTR_WAIT: begin
            if (bus_ready) begin
                state_next = STATE_IDLE;
            end else if (bus_error || bus_timeout || adapter_timeout_counter >= 4'hF) begin
                state_next = STATE_ERROR;
            end
        end
        
        STATE_ERROR: begin
            // Stay in error state until next request
            if (!cpu_dmem_req && !cpu_imem_req) begin
                state_next = STATE_IDLE;
            end
        end
        
        default: begin
            state_next = STATE_IDLE;
        end
    endcase
end

// State register update
always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        state_reg <= STATE_IDLE;
        data_write_reg <= {`VTX1_WORD_WIDTH{1'b0}};
        data_read_reg <= {`VTX1_WORD_WIDTH{1'b0}};
        data_request_pending <= 1'b0;
        instruction_request_pending <= 1'b0;
        data_transactions_reg <= 32'b0;
        instruction_transactions_reg <= 32'b0;        error_count_reg <= 32'b0;
        adapter_timeout_counter <= 4'b0;
    end else begin
        state_reg <= state_next;
        
        case (state_reg)            STATE_IDLE: begin
                adapter_timeout_counter <= 4'b0;
                if (cpu_dmem_req) begin
                    data_request_pending <= 1'b1;
                    if (cpu_dmem_we) begin
                        data_write_reg <= dmem_data_in;
                    end
                end else if (cpu_imem_req) begin
                    instruction_request_pending <= 1'b1;
                end
            end            STATE_DATA_REQUEST, STATE_DATA_WAIT: begin
                adapter_timeout_counter <= adapter_timeout_counter + 1;
                if (bus_ready && !bus_error) begin
                    data_request_pending <= 1'b0;
                    data_transactions_reg <= data_transactions_reg + 1;
                                    if (!cpu_dmem_we) begin
                        data_read_reg <= bus_rdata;
                    end
                    `VTX1_CLEAR_ERROR(vtx1_error_reg, state_reg)
                    vtx1_error_info <= 32'h0;
                    vtx1_error_valid <= 1'b0;
                end else if (bus_error || bus_timeout) begin
                    `VTX1_SET_ERROR(vtx1_error_reg, state_reg, `VTX1_ERROR_BUS_FAULT)
                    vtx1_error_info <= 32'h0000_0001;
                    vtx1_error_valid <= 1'b1;
                    error_count_reg <= error_count_reg + 1;
                    data_request_pending <= 1'b0;
                end            end
              STATE_INSTR_REQUEST, STATE_INSTR_WAIT: begin
                adapter_timeout_counter <= adapter_timeout_counter + 1;
                if (bus_ready && !bus_error) begin
                    instruction_request_pending <= 1'b0;
                    instruction_transactions_reg <= instruction_transactions_reg + 1;
                    `VTX1_CLEAR_ERROR(vtx1_error_reg, state_reg)
                    vtx1_error_info <= 32'h0;
                    vtx1_error_valid <= 1'b0;
                end else if (bus_error || bus_timeout) begin
                    `VTX1_SET_ERROR(vtx1_error_reg, state_reg, `VTX1_ERROR_BUS_FAULT)
                    vtx1_error_info <= 32'h0000_0002;
                    vtx1_error_valid <= 1'b1;
                    error_count_reg <= error_count_reg + 1;
                    instruction_request_pending <= 1'b0;
                end
            end
              STATE_ERROR: begin
                if (!cpu_dmem_req && !cpu_imem_req) begin
                    data_request_pending <= 1'b0;
                    instruction_request_pending <= 1'b0;
                    adapter_timeout_counter <= 4'b0;
                    `VTX1_CLEAR_ERROR(vtx1_error_reg, state_reg)
                    vtx1_error_info <= 32'h0;
                    vtx1_error_valid <= 1'b0;
                end
            end
        endcase
    end
end

// =============================================================================
// OUTPUT CONTROL LOGIC
// =============================================================================

// Bus interface control
always @(*) begin
    // Default values
    bus_req = 1'b0;
    bus_wr = 1'b0;
    bus_size = 2'b10; // Word size by default
    bus_addr = {`VTX1_ADDR_WIDTH{1'b0}};
    bus_wdata = {`VTX1_WORD_WIDTH{1'b0}};
    bus_error_clear = 1'b0;
    
    case (state_reg)
        STATE_DATA_REQUEST, STATE_DATA_WAIT: begin
            bus_req = data_request_pending;
            bus_wr = cpu_dmem_we;
            bus_addr = cpu_dmem_addr;
            bus_wdata = data_write_reg;
        end
        
        STATE_INSTR_REQUEST, STATE_INSTR_WAIT: begin
            bus_req = instruction_request_pending;
            bus_wr = 1'b0; // Instructions are always read
            bus_addr = cpu_imem_addr;
        end
        
        STATE_ERROR: begin
            bus_error_clear = 1'b1;
        end
    endcase
end

// CPU interface control
always @(*) begin
    // Default values
    cpu_dmem_ready = 1'b0;
    cpu_imem_ready = 1'b0;
    dmem_data_oe = 1'b0;
    dmem_data_out = {`VTX1_WORD_WIDTH{1'b0}};
    
    case (state_reg)
        STATE_IDLE: begin
            // Ready when bus transaction completes successfully
            if (!data_request_pending && !instruction_request_pending) begin
                cpu_dmem_ready = !cpu_dmem_req; // Ready when no request pending
                cpu_imem_ready = !cpu_imem_req;
            end
        end
        
        STATE_DATA_REQUEST, STATE_DATA_WAIT: begin
            if (bus_ready && !bus_error) begin
                cpu_dmem_ready = 1'b1;
                if (!cpu_dmem_we) begin
                    dmem_data_oe = 1'b1;
                    dmem_data_out = bus_rdata;
                end
            end
        end
        
        STATE_INSTR_REQUEST, STATE_INSTR_WAIT: begin
            if (bus_ready && !bus_error) begin
                cpu_imem_ready = 1'b1;
            end
        end
        
        STATE_ERROR: begin
            // Signal error condition to CPU
            cpu_dmem_ready = 1'b0;
            cpu_imem_ready = 1'b0;
        end
    endcase
end

// Status outputs - Enhanced with VTX1 error tracking
always @(*) begin
    adapter_state = state_reg;
    data_transactions = data_transactions_reg;
    instruction_transactions = instruction_transactions_reg;
    error_count = error_count_reg;
    adapter_error = vtx1_error_valid || (state_reg == STATE_ERROR) || bus_error;
    adapter_error_code = vtx1_error_valid ? vtx1_error_reg : (bus_error ? bus_error_code : 4'h0);
end

endmodule

`endif // CPU_ADAPTER_V

