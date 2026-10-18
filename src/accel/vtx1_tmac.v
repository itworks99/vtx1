`timescale 1ns / 1ps
// VTX1 Ternary-Weight MAC Array (the VDOT growth path)
// Part of the VTX1 Ternary System-on-Chip
//
// 16 multiply-accumulate lanes with NO multipliers: weights are trits
// {-1, 0, +1} (BitNet-style ternary quantization), so a weight "multiply"
// is tritwise select/negate - negation is a wire swap in balanced
// ternary, zero is the constant zero word. The 16 selected words reduce
// through a tree of the ALU's 18-trit adders into an 18-trit accumulator.
//
// One 18-trit machine word carries all 16 lane weights in its low 16
// trits (lane i = trit value at weight 3^i, i.e. bits [2i+1:2i]), so a
// weight word is written as the balanced number sum(w_i * 3^i).
//
// Arithmetic wraps mod 3^18 by dropped carries, exactly like the CPU's
// ADD chain - a software MAC loop and this unit produce bit-identical
// accumulators on the same data.
//
// Per step (one cycle): acc += sum_{i=0}^{15} w_i (*) act_i

`ifndef VTX1_TMAC_V
`define VTX1_TMAC_V

module vtx1_tmac (
    input  wire         clk,
    input  wire         rst_n,

    input  wire         clear,    // zero the accumulator (wins over step)
    input  wire         step,     // perform one 16-lane MAC step
    input  wire [31:0]  weights,  // 16 trit weights, lane i = bits [2i+1:2i]
    input  wire [575:0] acts,     // 16 activation words, lane i = [36i+35:36i]
    output reg  [35:0]  acc
);

    localparam [35:0] ZERO_WORD = 36'h555555555; // 18 zero trits

    // Tritwise negation: swap the -1 and +1 encodings
    function [35:0] tneg(input [35:0] w);
        integer i;
        for (i = 0; i < 18; i = i + 1)
            tneg[2*i+1 -: 2] = (w[2*i+1 -: 2] == 2'b00) ? 2'b10 :
                               (w[2*i+1 -: 2] == 2'b10) ? 2'b00 :
                                                          w[2*i+1 -: 2];
    endfunction

    // The "multiplier": one 2-bit decode per lane
    function [35:0] tsel(input [1:0] w, input [35:0] a);
        case (w)
            2'b10:   tsel = a;         // +1
            2'b00:   tsel = tneg(a);   // -1: a wire swap
            default: tsel = ZERO_WORD; //  0
        endcase
    endfunction

    wire [35:0] sel [0:15];
    genvar gi;
    generate
        for (gi = 0; gi < 16; gi = gi + 1) begin : lanes
            assign sel[gi] = tsel(weights[2*gi+1 -: 2],
                                  acts[36*gi + 35 -: 36]);
        end
    endgenerate

    // Reduction tree: 8 + 4 + 2 + 1 adders, then the accumulate adder.
    // Dropped carries = wrap mod 3^18 (matches the CPU ADD semantics).
    wire [35:0] l1 [0:7];
    wire [35:0] l2 [0:3];
    wire [35:0] l3 [0:1];
    wire [35:0] l4;
    wire [35:0] acc_next;

    generate
        for (gi = 0; gi < 8; gi = gi + 1) begin : t1
            vtx1_trit_adder a (.a(sel[2*gi]), .b(sel[2*gi+1]),
                               .negate_b(1'b0), .sum(l1[gi]), .cout());
        end
        for (gi = 0; gi < 4; gi = gi + 1) begin : t2
            vtx1_trit_adder a (.a(l1[2*gi]), .b(l1[2*gi+1]),
                               .negate_b(1'b0), .sum(l2[gi]), .cout());
        end
        for (gi = 0; gi < 2; gi = gi + 1) begin : t3
            vtx1_trit_adder a (.a(l2[2*gi]), .b(l2[2*gi+1]),
                               .negate_b(1'b0), .sum(l3[gi]), .cout());
        end
    endgenerate

    vtx1_trit_adder t4  (.a(l3[0]), .b(l3[1]),
                         .negate_b(1'b0), .sum(l4), .cout());
    vtx1_trit_adder tacc (.a(acc), .b(l4),
                          .negate_b(1'b0), .sum(acc_next), .cout());

    always @(posedge clk or negedge rst_n) begin
        if (!rst_n)
            acc <= ZERO_WORD;
        else if (clear)
            acc <= ZERO_WORD;
        else if (step)
            acc <= acc_next;
    end

endmodule

`endif // VTX1_TMAC_V
