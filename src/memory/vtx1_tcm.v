`timescale 1ns / 1ps
// VTX1 Tightly-Coupled Instruction Memory (TCM)
//
// A single-cycle synchronous-read instruction RAM on the core's dedicated
// fetch port, covering one aligned power-of-two window of the word address
// space (default 0x1000..0x1FFF: the program region). `hit` answers
// combinationally whether the presented address is TCM-resident; when it
// is, THREE consecutive instruction words (addr, addr+1, addr+2) are in
// rdata0..rdata2 one clock later, feeding the core's triple-issue group.
//
// The array is interleaved across four banks on addr[1:0], so any three
// consecutive addresses land in three distinct banks and read in the same
// cycle regardless of alignment.
//
// The write port snoops the CPU's store path (write-through): a CPU
// store into the window updates the TCM alongside main memory, which
// keeps the two coherent for CPU-originated writes and gives the boot
// flow a natural load path (copy the program into the window with
// ordinary stores). Writes from other bus masters (DMA, debug) are NOT
// snooped. In simulation the testbench may also preload the banks
// hierarchically from the assembled image (word w -> bank w[1:0],
// index w[ADDR_BITS-1:2]).

`ifndef VTX1_TCM_V
`define VTX1_TCM_V

module vtx1_tcm #(
    parameter ADDR_BITS = 12,                // 2^12 = 4096 words
    parameter [31:0] BASE = 32'h0000_1000    // must be size-aligned
)(
    input  wire        clk,
    input  wire        en,      // fetch strobe: addr is valid this cycle
    input  wire [31:0] addr,    // binary word address (group base)
    output wire        hit,     // combinational: addr..addr+2 all in window
    output reg  [35:0] rdata0,  // word at addr,   one cycle after en&hit
    output reg  [35:0] rdata1,  // word at addr+1
    output reg  [35:0] rdata2,  // word at addr+2

    // Write-through snoop on the CPU store path
    input  wire        wen,
    input  wire [31:0] waddr,
    input  wire [35:0] wdata
);

    localparam BANK_BITS = ADDR_BITS - 2;

    // (No helper functions that read module signals here: Icarus computes
    // continuous-assign and always-@* sensitivity from function ARGUMENTS
    // only, so hidden signal reads inside a function never re-evaluate.)

    wire [31:0] a0 = addr;
    wire [31:0] a1 = addr + 32'd1;
    wire [31:0] a2 = addr + 32'd2;

    // The whole 3-word group must be resident (the last two words of the
    // window therefore fetch over the bus, which mirrors the program).
    assign hit = (a0[31:ADDR_BITS] == BASE[31:ADDR_BITS])
              && (a2[31:ADDR_BITS] == BASE[31:ADDR_BITS]);

    wire whit = (waddr[31:ADDR_BITS] == BASE[31:ADDR_BITS]);

    reg [35:0] bank0 [0:(1 << BANK_BITS) - 1];
    reg [35:0] bank1 [0:(1 << BANK_BITS) - 1];
    reg [35:0] bank2 [0:(1 << BANK_BITS) - 1];
    reg [35:0] bank3 [0:(1 << BANK_BITS) - 1];

    // Each bank serves the group word whose address lands in it (three
    // consecutive addresses always occupy three distinct banks; the
    // fourth bank reads a don't-care).
    wire [31:0] ab0 = (a0[1:0] == 2'd0) ? a0 : (a1[1:0] == 2'd0) ? a1 : a2;
    wire [31:0] ab1 = (a0[1:0] == 2'd1) ? a0 : (a1[1:0] == 2'd1) ? a1 : a2;
    wire [31:0] ab2 = (a0[1:0] == 2'd2) ? a0 : (a1[1:0] == 2'd2) ? a1 : a2;
    wire [31:0] ab3 = (a0[1:0] == 2'd3) ? a0 : (a1[1:0] == 2'd3) ? a1 : a2;

    reg [35:0] q0, q1, q2, q3;  // per-bank read data
    reg [1:0]  rsel;            // latched low bits of the group base

    always @(posedge clk) begin
        if (en && hit) begin
            q0   <= bank0[ab0[ADDR_BITS-1:2]];
            q1   <= bank1[ab1[ADDR_BITS-1:2]];
            q2   <= bank2[ab2[ADDR_BITS-1:2]];
            q3   <= bank3[ab3[ADDR_BITS-1:2]];
            rsel <= a0[1:0];
        end
        if (wen && whit) begin
            case (waddr[1:0])
                2'd0: bank0[waddr[ADDR_BITS-1:2]] <= wdata;
                2'd1: bank1[waddr[ADDR_BITS-1:2]] <= wdata;
                2'd2: bank2[waddr[ADDR_BITS-1:2]] <= wdata;
                2'd3: bank3[waddr[ADDR_BITS-1:2]] <= wdata;
            endcase
        end
    end

    // Rotate the bank outputs back into group order.
    always @* begin
        case (rsel)
            2'd0: begin rdata0 = q0; rdata1 = q1; rdata2 = q2; end
            2'd1: begin rdata0 = q1; rdata1 = q2; rdata2 = q3; end
            2'd2: begin rdata0 = q2; rdata1 = q3; rdata2 = q0; end
            2'd3: begin rdata0 = q3; rdata1 = q0; rdata2 = q1; end
        endcase
    end

endmodule

`endif // VTX1_TCM_V
