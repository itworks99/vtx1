`timescale 1ns / 1ps
// RV32 Tightly-Coupled Instruction Memory: the binary mirror of vtx1_tcm.
//
// Same structure as the VTX1 TCM so the memory models of the two cores in
// the issue-width-matched comparison are identical: four banks interleaved
// on the instruction index, delivering THREE consecutive instructions
// (addr, addr+4, addr+8; byte addresses) one clock after the fetch strobe,
// with a write-through snoop on the store path.

`ifndef RV32_TCM_V
`define RV32_TCM_V

module rv32_tcm #(
    parameter ADDR_BITS = 12,              // 2^12 = 4096 instructions
    parameter [31:0] BASE = 32'h0000_0000  // byte address, size-aligned
)(
    input  wire        clk,
    input  wire        en,
    input  wire [31:0] addr,    // byte address of the group base
    output wire        hit,     // addr..addr+8 all inside the window
    output reg  [31:0] rdata0,
    output reg  [31:0] rdata1,
    output reg  [31:0] rdata2,

    // Write-through snoop on the store path (byte address, word stores)
    input  wire        wen,
    input  wire [31:0] waddr,
    input  wire [31:0] wdata
);

    localparam BANK_BITS = ADDR_BITS - 2;

    // Instruction indices (word addresses)
    wire [31:0] a0 = addr >> 2;
    wire [31:0] a1 = (addr >> 2) + 32'd1;
    wire [31:0] a2 = (addr >> 2) + 32'd2;
    wire [31:0] wa = waddr >> 2;

    assign hit = (a0[29:ADDR_BITS] == BASE[31:ADDR_BITS+2])
              && (a2[29:ADDR_BITS] == BASE[31:ADDR_BITS+2]);

    wire whit = (wa[29:ADDR_BITS] == BASE[31:ADDR_BITS+2]);

    reg [31:0] bank0 [0:(1 << BANK_BITS) - 1];
    reg [31:0] bank1 [0:(1 << BANK_BITS) - 1];
    reg [31:0] bank2 [0:(1 << BANK_BITS) - 1];
    reg [31:0] bank3 [0:(1 << BANK_BITS) - 1];

    wire [31:0] ab0 = (a0[1:0] == 2'd0) ? a0 : (a1[1:0] == 2'd0) ? a1 : a2;
    wire [31:0] ab1 = (a0[1:0] == 2'd1) ? a0 : (a1[1:0] == 2'd1) ? a1 : a2;
    wire [31:0] ab2 = (a0[1:0] == 2'd2) ? a0 : (a1[1:0] == 2'd2) ? a1 : a2;
    wire [31:0] ab3 = (a0[1:0] == 2'd3) ? a0 : (a1[1:0] == 2'd3) ? a1 : a2;

    reg [31:0] q0, q1, q2, q3;
    reg [1:0]  rsel;

    always @(posedge clk) begin
        if (en && hit) begin
            q0   <= bank0[ab0[ADDR_BITS-1:2]];
            q1   <= bank1[ab1[ADDR_BITS-1:2]];
            q2   <= bank2[ab2[ADDR_BITS-1:2]];
            q3   <= bank3[ab3[ADDR_BITS-1:2]];
            rsel <= a0[1:0];
        end
        if (wen && whit) begin
            case (wa[1:0])
                2'd0: bank0[wa[ADDR_BITS-1:2]] <= wdata;
                2'd1: bank1[wa[ADDR_BITS-1:2]] <= wdata;
                2'd2: bank2[wa[ADDR_BITS-1:2]] <= wdata;
                2'd3: bank3[wa[ADDR_BITS-1:2]] <= wdata;
            endcase
        end
    end

    always @* begin
        case (rsel)
            2'd0: begin rdata0 = q0; rdata1 = q1; rdata2 = q2; end
            2'd1: begin rdata0 = q1; rdata1 = q2; rdata2 = q3; end
            2'd2: begin rdata0 = q2; rdata1 = q3; rdata2 = q0; end
            2'd3: begin rdata0 = q3; rdata1 = q0; rdata2 = q1; end
        endcase
    end

endmodule

`endif // RV32_TCM_V
