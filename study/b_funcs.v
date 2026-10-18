`timescale 1ns / 1ps
// Function-level study: the binary counterparts, at the minimal width
// that covers the 18-trit value range (+-193,710,244 < 2^28, so 29-bit
// two's complement). Written as plain behavioral Verilog so synthesis
// gives binary its best idiomatic implementation.

`ifndef B_FUNCS_V
`define B_FUNCS_V

// Negation: invert + increment - a full carry chain.
module b_neg (
    input  wire signed [28:0] a,
    output wire signed [28:0] y
);
    assign y = -a;
endmodule

// Sign: {-1, 0, +1}, encoded in 2 bits like the trit output (00/01/10).
module b_sign (
    input  wire signed [28:0] a,
    output wire [1:0] y
);
    assign y = (a == 29'sd0) ? 2'b01 : a[28] ? 2'b00 : 2'b10;
endmodule

// Absolute value: conditional negation - the carry chain again.
module b_abs (
    input  wire signed [28:0] a,
    output wire signed [28:0] y
);
    assign y = a[28] ? -a : a;
endmodule

// Three-way compare: sign(a - b) as a 2-bit code (00 lt / 01 eq / 10 gt).
module b_cmp3 (
    input  wire signed [28:0] a,
    input  wire signed [28:0] b,
    output wire [1:0] y
);
    assign y = (a == b) ? 2'b01 : (a < b) ? 2'b00 : 2'b10;
endmodule

// Addition mod 2^29.
module b_add (
    input  wire signed [28:0] a,
    input  wire signed [28:0] b,
    output wire signed [28:0] y
);
    assign y = a + b;
endmodule

// Multiplication mod 2^29.
module b_mul (
    input  wire signed [28:0] a,
    input  wire signed [28:0] b,
    output wire signed [28:0] y
);
    assign y = a * b;
endmodule

// Scale down by the radix with round-to-nearest (ties away from zero):
// shift plus a conditional increment - a carry chain where ternary
// needs zero gates.
module b_scale (
    input  wire signed [28:0] a,
    output wire signed [28:0] y
);
    assign y = (a >>> 1) + {28'd0, (a[0] & ~a[28])};
endmodule

`endif // B_FUNCS_V
