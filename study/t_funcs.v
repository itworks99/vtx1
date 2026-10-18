`timescale 1ns / 1ps
// Function-level study: balanced-ternary math units (18-trit words,
// 2-bit encoding), each paired with an equal-value-range binary unit in
// b_funcs.v. The adder and multiplier are the ALU's own verified
// structures (vtx1_trit_adder / vtx1_trit_mul from src/cpu/vtx1_alu.v).
//
// The point of the study: compare FUNCTIONS over equal ranges, not
// single-digit signal operations - this is where balanced ternary's
// structural properties (carry-free negation, no sign extension,
// round-to-nearest truncation, fewer digits per range) actually live.

`ifndef T_FUNCS_V
`define T_FUNCS_V

// Tritwise negation of a word: NO CARRY CHAIN - 18 independent swaps.
module t_neg (
    input  wire [35:0] a,
    output wire [35:0] y
);
    genvar gi;
    generate
        for (gi = 0; gi < 18; gi = gi + 1) begin : tn
            // 00 <-> 10, 01 stays: y[2i+1] = ~b1 & ~b0, y[2i] = b0
            assign y[2*gi+1] = ~a[2*gi+1] & ~a[2*gi];
            assign y[2*gi]   = a[2*gi];
        end
    endgenerate
endmodule

// Sign of a word: its most significant nonzero trit (one trit out).
module t_sign (
    input  wire [35:0] a,
    output wire [1:0]  y
);
    function [1:0] ws(input [35:0] w);
        integer i;
        begin
            ws = 2'b01;
            for (i = 17; i >= 0; i = i - 1)
                if (ws == 2'b01 && w[2*i+1 -: 2] != 2'b01)
                    ws = w[2*i+1 -: 2];
        end
    endfunction
    assign y = ws(a);
endmodule

// Absolute value: conditional tritwise negation - the only "arithmetic"
// is the sign scan; there is still no carry chain.
module t_abs (
    input  wire [35:0] a,
    output wire [35:0] y
);
    wire [1:0]  s;
    wire [35:0] n;
    t_sign u_s (.a(a), .y(s));
    t_neg  u_n (.a(a), .y(n));
    assign y = (s == 2'b00) ? n : a;
endmodule

// Three-way compare: sign(a - b) as one trit, overflow-correct (a
// nonzero carry-out trit of the subtraction is the true sign).
module t_cmp3 (
    input  wire [35:0] a,
    input  wire [35:0] b,
    output wire [1:0]  y
);
    wire [35:0] d;
    wire [1:0]  co;
    vtx1_trit_adder u_sub (.a(a), .b(b), .negate_b(1'b1), .sum(d), .cout(co));
    wire [1:0] ds;
    t_sign u_s (.a(d), .y(ds));
    assign y = (co != 2'b01) ? co : ds;
endmodule

// Addition mod 3^18 (the ALU's ripple structure).
module t_add (
    input  wire [35:0] a,
    input  wire [35:0] b,
    output wire [35:0] y
);
    vtx1_trit_adder u (.a(a), .b(b), .negate_b(1'b0), .sum(y), .cout());
endmodule

// Multiplication mod 3^18 (the ALU's 18-partial tree; digits are already
// in {-1,0,+1}, so partial products are select/negate - Booth for free).
module t_mul (
    input  wire [35:0] a,
    input  wire [35:0] b,
    output wire [35:0] y
);
    vtx1_trit_mul u (.a(a), .b(b), .product(y));
endmodule

// Scale down by the radix with round-to-nearest: in balanced ternary,
// dropping the least significant trit IS nearest rounding (the remainder
// is -1, 0 or +1 - never a tie). Zero gates, zero depth.
module t_scale (
    input  wire [35:0] a,
    output wire [35:0] y
);
    assign y = {2'b01, a[35:2]}; // shift one trit, zero-trit fill
endmodule

`endif // T_FUNCS_V
