`timescale 1ns / 1ps
// VTX1 Native Balanced-Ternary ALU
// Part of the VTX1 Ternary System-on-Chip
//
// All arithmetic is performed directly on trits (2-bit encoding: 00 = -1,
// 01 = 0, 10 = +1) with NO binary conversion anywhere in the datapath:
//
//   - vtx1_trit_fa:    the balanced ternary full adder cell
//                      (a + b + cin -> sum trit + carry trit)
//   - vtx1_trit_adder: 18-trit ripple adder with optional tritwise negation
//                      of B (subtraction costs nothing extra). The final
//                      carry is dropped, which IS balanced wrap-around
//                      modulo 3^18 - the same semantics as the ISA
//                      simulator's WrapWord.
//   - vtx1_trit_mul:   combinational multiplier: 18 partial products
//                      (each is -B, 0, or +B, shifted) reduced by a tree
//                      of trit adders, truncated to 18 trits (mod wrap).
//   - vtx1_alu:        operation mux: ADD/SUB/MUL, tritwise AND(min)/
//                      OR(max)/XOR(product)/NOT(negate), trit shifts and
//                      rotates, CMP/TEST (three-valued sign as a value),
//                      ABS, NEG, INC, DEC.
//
// Shift/rotate amount semantics (ISA note): the amount is the value of the
// LOW 4 TRITS of operand B (range +-40), reduced modulo 18; negative
// amounts therefore shift the other way around the ring. The simulator
// implements the identical rule.
//
// DIV/MOD are not combinational operations; they belong to the iterative
// divider FSM (pending) and are not handled here.

`ifndef VTX1_ALU_V
`define VTX1_ALU_V

`include "vtx1_isa_generated.vh"

// ---------------------------------------------------------------------------
// Balanced ternary full adder cell: s = a + b + cin, |s| <= 3
//   -3 -> carry -1, sum  0        +3 -> carry +1, sum  0
//   -2 -> carry -1, sum +1        +2 -> carry +1, sum -1
//   -1 -> carry  0, sum -1        +1 -> carry  0, sum +1
//    0 -> carry  0, sum  0
// ---------------------------------------------------------------------------
module vtx1_trit_fa (
    input  wire [1:0] a,
    input  wire [1:0] b,
    input  wire [1:0] cin,
    output reg  [1:0] sum,
    output reg  [1:0] cout
);
    // Trit value: 00 = -1, 01 = 0, 10 = +1 (11 treated as 0)
    function signed [2:0] tv(input [1:0] t);
        case (t)
            2'b00:   tv = -3'sd1;
            2'b10:   tv = 3'sd1;
            default: tv = 3'sd0;
        endcase
    endfunction

    reg signed [2:0] s;
    always @* begin
        s = tv(a) + tv(b) + tv(cin);
        case (s)
            -3'sd3:  begin cout = 2'b00; sum = 2'b01; end
            -3'sd2:  begin cout = 2'b00; sum = 2'b10; end
            -3'sd1:  begin cout = 2'b01; sum = 2'b00; end
            3'sd1:   begin cout = 2'b01; sum = 2'b10; end
            3'sd2:   begin cout = 2'b10; sum = 2'b00; end
            3'sd3:   begin cout = 2'b10; sum = 2'b01; end
            default: begin cout = 2'b01; sum = 2'b01; end
        endcase
    end
endmodule

// ---------------------------------------------------------------------------
// 18-trit ripple adder. negate_b turns it into a subtractor for free
// (balanced ternary negation is tritwise). Carry-out is dropped: balanced
// wrap modulo 3^18.
// ---------------------------------------------------------------------------
module vtx1_trit_adder (
    input  wire [35:0] a,
    input  wire [35:0] b,
    input  wire        negate_b,
    output wire [35:0] sum,
    output wire [1:0]  cout   // final carry trit: nonzero means the true
                              // result overflowed 18 trits, and the carry's
                              // sign IS the true sign of the result
);
    // Tritwise negation: swap the -1 and +1 encodings (00 <-> 10)
    function [1:0] tneg(input [1:0] t);
        case (t)
            2'b00:   tneg = 2'b10;
            2'b10:   tneg = 2'b00;
            default: tneg = t;
        endcase
    endfunction

    wire [35:0] bn;
    genvar gi;
    generate
        for (gi = 0; gi < 18; gi = gi + 1) begin : negb
            assign bn[2*gi+1:2*gi] = negate_b ? tneg(b[2*gi+1:2*gi]) : b[2*gi+1:2*gi];
        end
    endgenerate

    // Ripple from the least significant trit (bit pair [1:0] is trit 17).
    wire [1:0] carry [0:18];
    assign carry[0] = 2'b01; // zero
    generate
        for (gi = 0; gi < 18; gi = gi + 1) begin : stage
            vtx1_trit_fa fa (
                .a(a[2*gi+1:2*gi]), .b(bn[2*gi+1:2*gi]), .cin(carry[gi]),
                .sum(sum[2*gi+1:2*gi]), .cout(carry[gi+1])
            );
        end
    endgenerate
    assign cout = carry[18];
endmodule

// ---------------------------------------------------------------------------
// Combinational multiplier: result = (a * b) wrapped to 18 trits.
// Partial product i is b selected by trit i of a (-B, 0, +B), shifted left
// by i trit positions and truncated; partials are reduced by a binary tree
// of trit adders (depth 5).
// ---------------------------------------------------------------------------
module vtx1_trit_mul (
    input  wire [35:0] a,
    input  wire [35:0] b,
    output wire [35:0] product
);
    localparam [35:0] ZERO_WORD = 36'h555555555;

    function [1:0] tneg(input [1:0] t);
        case (t)
            2'b00:   tneg = 2'b10;
            2'b10:   tneg = 2'b00;
            default: tneg = t;
        endcase
    endfunction

    // Partial products (trit i of a is at bits [2i+1:2i], i counted from LSB)
    wire [35:0] pp [0:17];
    genvar gi, gj;
    generate
        for (gi = 0; gi < 18; gi = gi + 1) begin : part
            wire [1:0] sel = a[2*gi+1:2*gi];
            // b shifted left by gi trits, truncated to 18 trits, with
            // zero-trit fill; then negated/zeroed by the selector trit.
            for (gj = 0; gj < 18; gj = gj + 1) begin : bits
                wire [1:0] src = (gj >= gi) ? b[2*(gj-gi)+1:2*(gj-gi)] : 2'b01;
                assign pp[gi][2*gj+1:2*gj] =
                    (sel == 2'b10) ? src :
                    (sel == 2'b00) ? tneg(src) : 2'b01;
            end
        end
    endgenerate

    // Reduction tree: 18 -> 9 -> 5 -> 3 -> 2 -> 1
    wire [35:0] l1 [0:8];
    wire [35:0] l2 [0:4];
    wire [35:0] l3 [0:2];
    wire [35:0] l4 [0:1];
    generate
        for (gi = 0; gi < 9; gi = gi + 1) begin : red1
            vtx1_trit_adder add (.a(pp[2*gi]), .b(pp[2*gi+1]), .negate_b(1'b0), .sum(l1[gi]), .cout());
        end
        for (gi = 0; gi < 4; gi = gi + 1) begin : red2
            vtx1_trit_adder add (.a(l1[2*gi]), .b(l1[2*gi+1]), .negate_b(1'b0), .sum(l2[gi]), .cout());
        end
    endgenerate
    assign l2[4] = l1[8];
    vtx1_trit_adder r3a (.a(l2[0]), .b(l2[1]), .negate_b(1'b0), .sum(l3[0]), .cout());
    vtx1_trit_adder r3b (.a(l2[2]), .b(l2[3]), .negate_b(1'b0), .sum(l3[1]), .cout());
    assign l3[2] = l2[4];
    vtx1_trit_adder r4a (.a(l3[0]), .b(l3[1]), .negate_b(1'b0), .sum(l4[0]), .cout());
    assign l4[1] = l3[2];
    vtx1_trit_adder r5  (.a(l4[0]), .b(l4[1]), .negate_b(1'b0), .sum(product), .cout());
endmodule

// ---------------------------------------------------------------------------
// The ALU
// ---------------------------------------------------------------------------
module vtx1_alu (
    input  wire [35:0] a,       // first operand (rs1)
    input  wire [35:0] b,       // second operand (rs2 / immediate as trits)
    input  wire [6:0]  opcode,  // ISA opcode value (`VTX1_OPV_*)
    output reg  [35:0] result,
    output wire        unsupported // opcode is not a combinational ALU op
);
    localparam [35:0] ZERO_WORD = 36'h555555555;
    localparam [35:0] ONE_WORD  = 36'h555555556; // trits 0...0,+1

    function signed [2:0] tv(input [1:0] t);
        case (t)
            2'b00:   tv = -3'sd1;
            2'b10:   tv = 3'sd1;
            default: tv = 3'sd0;
        endcase
    endfunction

    function [1:0] tenc(input signed [2:0] v);
        case (v)
            -3'sd1:  tenc = 2'b00;
            3'sd1:   tenc = 2'b10;
            default: tenc = 2'b01;
        endcase
    endfunction

    // --- Sign of a word: the most significant nonzero trit ---
    function [1:0] word_sign(input [35:0] w);
        integer i;
        begin
            word_sign = 2'b01; // zero
            for (i = 17; i >= 0; i = i - 1) // from MST (bits [35:34]) down
                if (word_sign == 2'b01 && w[2*i+1 -: 2] != 2'b01)
                    word_sign = w[2*i+1 -: 2];
        end
    endfunction

    // --- Tritwise logic ---
    function [35:0] tw_min(input [35:0] x, input [35:0] y);
        integer i;
        for (i = 0; i < 18; i = i + 1)
            tw_min[2*i+1 -: 2] = (tv(x[2*i+1 -: 2]) < tv(y[2*i+1 -: 2]))
                                 ? x[2*i+1 -: 2] : y[2*i+1 -: 2];
    endfunction

    function [35:0] tw_max(input [35:0] x, input [35:0] y);
        integer i;
        for (i = 0; i < 18; i = i + 1)
            tw_max[2*i+1 -: 2] = (tv(x[2*i+1 -: 2]) > tv(y[2*i+1 -: 2]))
                                 ? x[2*i+1 -: 2] : y[2*i+1 -: 2];
    endfunction

    function [35:0] tw_mul(input [35:0] x, input [35:0] y);
        integer i;
        for (i = 0; i < 18; i = i + 1)
            tw_mul[2*i+1 -: 2] = tenc(tv(x[2*i+1 -: 2]) * tv(y[2*i+1 -: 2]));
    endfunction

    function [35:0] tw_neg(input [35:0] x);
        integer i;
        for (i = 0; i < 18; i = i + 1)
            tw_neg[2*i+1 -: 2] = tenc(-tv(x[2*i+1 -: 2]));
    endfunction

    // --- Shift amount: value of B's low 4 trits (range -40..+40), mod 18.
    // Offset by +54 (= 0 mod 18) into 14..94, then reduce by subtraction.
    wire signed [7:0] amt_raw = 27*tv(b[7:6]) + 9*tv(b[5:4]) + 3*tv(b[3:2]) + tv(b[1:0]);
    wire [6:0] amt_pos = amt_raw + 8'sd54;
    wire [6:0] amt_s1  = (amt_pos >= 7'd72) ? amt_pos - 7'd72 : amt_pos;
    wire [6:0] amt_s2  = (amt_s1  >= 7'd54) ? amt_s1  - 7'd54 : amt_s1;
    wire [6:0] amt_s3  = (amt_s2  >= 7'd36) ? amt_s2  - 7'd36 : amt_s2;
    wire [4:0] amt     = (amt_s3  >= 7'd18) ? amt_s3[4:0] - 5'd18 : amt_s3[4:0];

    // --- Barrel shifter over trit positions (zero-trit fill / rotate) ---
    function [35:0] shift(input [35:0] x, input [4:0] k,
                          input left, input rotate);
        integer i, src;
        begin
            for (i = 0; i < 18; i = i + 1) begin
                // i counts trits from LSB here; "shift left" moves toward
                // more significant trits, matching the simulator.
                src = left ? i - k : i + k;
                if (src >= 0 && src < 18)
                    shift[2*i+1 -: 2] = x[2*src+1 -: 2];
                else if (rotate)
                    shift[2*i+1 -: 2] = x[2*(((src % 18) + 18) % 18)+1 -: 2];
                else
                    shift[2*i+1 -: 2] = 2'b01;
            end
        end
    endfunction

    // --- Shared units ---
    wire [35:0] sum, diff, prod;
    wire [1:0]  diff_cout;
    vtx1_trit_adder u_add (.a(a), .b(b), .negate_b(1'b0), .sum(sum), .cout());
    vtx1_trit_adder u_sub (.a(a), .b(b), .negate_b(1'b1), .sum(diff), .cout(diff_cout));
    vtx1_trit_mul   u_mul (.a(a), .b(b), .product(prod));

    // True sign of (a-b): when the subtract overflows 18 trits, the carry
    // trit carries the sign; otherwise the sign of the wrapped difference.
    wire [1:0] diff_sign = (diff_cout != 2'b01) ? diff_cout : word_sign(diff);
    wire [1:0] a_sign    = word_sign(a);

    reg unsup;
    assign unsupported = unsup;

    always @* begin
        unsup  = 1'b0;
        result = ZERO_WORD;
        case (opcode)
            `VTX1_OPV_ADD, `VTX1_OPV_INC: result = sum;  // INC: core passes b = +1
            `VTX1_OPV_SUB, `VTX1_OPV_DEC: result = diff; // DEC: core passes b = +1
            `VTX1_OPV_MUL:  result = prod;
            `VTX1_OPV_AND:  result = tw_min(a, b);
            `VTX1_OPV_OR:   result = tw_max(a, b);
            `VTX1_OPV_XOR:  result = tw_mul(a, b);
            `VTX1_OPV_NOT:  result = tw_neg(a);
            `VTX1_OPV_NEG:  result = tw_neg(a);
            `VTX1_OPV_ABS:  result = (a_sign == 2'b00) ? tw_neg(a) : a;
            `VTX1_OPV_CMP:  result = {ZERO_WORD[35:2], diff_sign};        // sign(a-b)
            `VTX1_OPV_TEST: result = {ZERO_WORD[35:2], word_sign(tw_min(a, b))};
            `VTX1_OPV_SHL:  result = shift(a, amt, 1'b1, 1'b0);
            `VTX1_OPV_SHR:  result = shift(a, amt, 1'b0, 1'b0);
            `VTX1_OPV_ROL:  result = shift(a, amt, 1'b1, 1'b1);
            `VTX1_OPV_ROR:  result = shift(a, amt, 1'b0, 1'b1);
            default: unsup = 1'b1;
        endcase
    end

endmodule

`endif // VTX1_ALU_V
