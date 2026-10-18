`timescale 1ns / 1ps
// Binary Ternary-Weight MAC Array: the matched control for vtx1_tmac.
//
// The same function in two's-complement: 16 lanes of ternary weights
// {-1, 0, +1} against 30-bit signed activations (the smallest power-free
// width that holds the 18-trit value range, +-193710244 < 2^29), a
// select/negate per lane and a binary adder tree. No multipliers here
// either - ternary WEIGHTS eliminate the multiplier on any substrate;
// what this control measures is the cost of the same unit in binary
// logic versus trit logic.
//
// On data that does not overflow either accumulator, this unit and
// vtx1_tmac produce identical numeric results (the testbench enforces
// it). Overflow wraps mod 2^30 here vs mod 3^18 there - stated, and
// excluded from the equivalence checks.

`ifndef BIN_TMAC_V
`define BIN_TMAC_V

module bin_tmac (
    input  wire         clk,
    input  wire         rst_n,

    input  wire         clear,
    input  wire         step,
    input  wire [31:0]  weights,  // same encoding: 10 = +1, 00 = -1, else 0
    input  wire [479:0] acts,     // 16 x 30-bit signed, lane i = [30i+29:30i]
    output reg signed [29:0] acc
);

    function signed [29:0] bsel(input [1:0] w, input signed [29:0] a);
        case (w)
            2'b10:   bsel = a;
            2'b00:   bsel = -a;
            default: bsel = 30'sd0;
        endcase
    endfunction

    wire signed [29:0] sel [0:15];
    genvar gi;
    generate
        for (gi = 0; gi < 16; gi = gi + 1) begin : lanes
            assign sel[gi] = bsel(weights[2*gi+1 -: 2],
                                  acts[30*gi + 29 -: 30]);
        end
    endgenerate

    wire signed [29:0] l1 [0:7];
    wire signed [29:0] l2 [0:3];
    wire signed [29:0] l3 [0:1];
    wire signed [29:0] l4;

    generate
        for (gi = 0; gi < 8; gi = gi + 1) begin : t1
            assign l1[gi] = sel[2*gi] + sel[2*gi+1];
        end
        for (gi = 0; gi < 4; gi = gi + 1) begin : t2
            assign l2[gi] = l1[2*gi] + l1[2*gi+1];
        end
        for (gi = 0; gi < 2; gi = gi + 1) begin : t3
            assign l3[gi] = l2[2*gi] + l2[2*gi+1];
        end
    endgenerate

    assign l4 = l3[0] + l3[1];

    always @(posedge clk or negedge rst_n) begin
        if (!rst_n)
            acc <= 30'sd0;
        else if (clear)
            acc <= 30'sd0;
        else if (step)
            acc <= acc + l4;
    end

endmodule

`endif // BIN_TMAC_V
