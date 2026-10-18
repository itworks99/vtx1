`timescale 1ns / 1ps
// VTX1 Tightly-Coupled Instruction Memory (TCM)
//
// A single-cycle synchronous-read RAM on the core's dedicated fetch port,
// covering one aligned power-of-two window of the word address space
// (default 0x1000..0x1FFF: the program region). `hit` answers
// combinationally whether the presented address is TCM-resident; when it
// is, the instruction word is in `rdata` one clock later, so the core
// fetches at full speed instead of arbitrating through the bus matrix.
//
// The write port snoops the CPU's store path (write-through): a CPU
// store into the window updates the TCM alongside main memory, which
// keeps the two coherent for CPU-originated writes and gives the boot
// flow a natural load path (copy the program into the window with
// ordinary stores). Writes from other bus masters (DMA, debug) are NOT
// snooped. In simulation the testbench may also preload `mem`
// hierarchically from the assembled image.

`ifndef VTX1_TCM_V
`define VTX1_TCM_V

module vtx1_tcm #(
    parameter ADDR_BITS = 12,                // 2^12 = 4096 words
    parameter [31:0] BASE = 32'h0000_1000    // must be size-aligned
)(
    input  wire        clk,
    input  wire        en,      // fetch strobe: addr is valid this cycle
    input  wire [31:0] addr,    // binary word address
    output wire        hit,     // combinational: addr is inside the window
    output reg  [35:0] rdata,   // instruction word, one cycle after en&hit

    // Write-through snoop on the CPU store path
    input  wire        wen,
    input  wire [31:0] waddr,
    input  wire [35:0] wdata
);

    assign hit = (addr[31:ADDR_BITS] == BASE[31:ADDR_BITS]);
    wire whit  = (waddr[31:ADDR_BITS] == BASE[31:ADDR_BITS]);

    reg [35:0] mem [0:(1 << ADDR_BITS) - 1];

    always @(posedge clk) begin
        if (en && hit)
            rdata <= mem[addr[ADDR_BITS-1:0]];
        if (wen && whit)
            mem[waddr[ADDR_BITS-1:0]] <= wdata;
    end

endmodule

`endif // VTX1_TCM_V
