`timescale 1ns / 1ps
// VTX1 Data Scratchpad: the data-side twin of the instruction TCM.
//
// Sits between the core's memory port and the SoC bus master port and
// serves in-window READS from a local synchronous RAM with a one-cycle
// ready handshake, instead of the ~7-cycle trip through the bus matrix
// and memory controller. WRITES always pass through to the bus
// (write-through: external memory remains the source of truth, so DMA,
// debug and the testbench see every store), and in-window writes are
// snooped into the local copy on completion. Out-of-window accesses pass
// through untouched, so MMIO (0xF0000000+) and anything beyond the
// window behave exactly as before.
//
// Coherence limits (stated): writes by other bus masters (DMA, debug)
// are NOT snooped. In simulation the testbench preloads `mem` with the
// same image as the external SRAM; on silicon the boot flow would.

`ifndef VTX1_SCRATCHPAD_V
`define VTX1_SCRATCHPAD_V

module vtx1_scratchpad #(
    parameter ADDR_BITS = 14,               // 2^14 = 16384 words
    parameter [31:0] BASE = 32'h0000_0000   // must be size-aligned
)(
    input  wire        clk,
    input  wire        rst_n,

    // CPU side (the core's memory port)
    input  wire        cpu_req,
    input  wire        cpu_wr,
    input  wire [31:0] cpu_addr,
    input  wire [35:0] cpu_wdata,
    output wire [35:0] cpu_rdata,
    output wire        cpu_ready,

    // Bus side (the SoC bus master port)
    output wire        bus_req,
    output wire        bus_wr,
    output wire [31:0] bus_addr,
    output wire [35:0] bus_wdata,
    input  wire [35:0] bus_rdata,
    input  wire        bus_ready
);

    wire in_win    = (cpu_addr[31:ADDR_BITS] == BASE[31:ADDR_BITS]);
    wire local_rd  = cpu_req && !cpu_wr && in_win;

    // Pass-through for writes and out-of-window reads
    assign bus_req   = cpu_req && !local_rd;
    assign bus_wr    = cpu_wr;
    assign bus_addr  = cpu_addr;
    assign bus_wdata = cpu_wdata;

    // Local memory: synchronous read, write-snoop, no reset (RAM inference)
    reg [35:0] mem [0:(1 << ADDR_BITS) - 1];
    reg [35:0] rdata_r;

    always @(posedge clk) begin
        if (local_rd)
            rdata_r <= mem[cpu_addr[ADDR_BITS-1:0]];
        if (cpu_req && cpu_wr && bus_ready && in_win)
            mem[cpu_addr[ADDR_BITS-1:0]] <= cpu_wdata;
    end

    // One-cycle ready handshake for local reads
    reg ready_r;
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n)
            ready_r <= 1'b0;
        else
            ready_r <= local_rd && !ready_r;
    end

    assign cpu_rdata = local_rd ? rdata_r : bus_rdata;
    assign cpu_ready = local_rd ? ready_r : bus_ready;

endmodule

`endif // VTX1_SCRATCHPAD_V
