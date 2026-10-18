`timescale 1ns / 1ps
// Conventional INT8 Multiplier MAC Array: the industry reference point.
//
// 16 lanes of signed 8x8 multiply feeding a binary adder tree and a
// 32-bit accumulator - the unit a conventional binary NPU builds per
// lane. This is NOT function-equivalent to the ternary-weight units
// (full INT8 weights vs {-1,0,+1} weights; INT8 activations vs the wider
// activation words); it is in the synthesis comparison purely to show
// what the ternary-weight quantization removes: the multipliers.

`ifndef BIN_MAC8_V
`define BIN_MAC8_V

module bin_mac8 (
    input  wire         clk,
    input  wire         rst_n,

    input  wire         clear,
    input  wire         step,
    input  wire [127:0] weights,  // 16 x int8
    input  wire [127:0] acts,     // 16 x int8
    output reg signed [31:0] acc
);

    wire signed [15:0] prod [0:15];
    genvar gi;
    generate
        for (gi = 0; gi < 16; gi = gi + 1) begin : lanes
            assign prod[gi] = $signed(weights[8*gi + 7 -: 8]) *
                              $signed(acts[8*gi + 7 -: 8]);
        end
    endgenerate

    wire signed [16:0] l1 [0:7];
    wire signed [17:0] l2 [0:3];
    wire signed [18:0] l3 [0:1];
    wire signed [19:0] l4;

    generate
        for (gi = 0; gi < 8; gi = gi + 1) begin : t1
            assign l1[gi] = prod[2*gi] + prod[2*gi+1];
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
            acc <= 32'sd0;
        else if (clear)
            acc <= 32'sd0;
        else if (step)
            acc <= acc + l4;
    end

endmodule

`endif // BIN_MAC8_V
