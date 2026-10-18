`timescale 1ns / 1ps
// VTX1 Iterative Divider: native balanced-ternary DIV/MOD.
// Part of the VTX1 Ternary System-on-Chip
//
// A multi-cycle FSM unit (the no-microcode answer to division): one
// quotient trit per cycle over 18 iterations, built from the same
// vtx1_trit_fa cells as the ALU. ~22 cycles per operation.
//
// Semantics (ratified to match vtxsim, which uses Go's rules):
//   DIV truncates toward zero; MOD's sign follows the dividend;
//   a / b * b + a % b == a always. Division by zero is checked by the
//   CORE before start (it traps), so b != 0 may be assumed here.
//
// Algorithm: negation is tritwise in balanced ternary, so the division
// runs on magnitudes and the signs are applied for free at the end.
//   setup:  A = |a| (trit shift register), B = |b|, B3 = 3B (one trit
//           shift, free), B2 = 2B = B3 - B (one adder)
//   loop (18x, MSB-first): R1 = 3R + next trit of A (shift, free);
//           quotient digit magnitude from two comparisons of 2|R1|
//           against B and 3B (digits -2..+2 internally: the +-2 case
//           absorbs the +-1 overshoot that pure {-1,0,+1} selection
//           leaves behind, keeping the invariant 2|R| <= B exact);
//           R -= digit*B; Q = 3Q + digit (adders)
//   correction: the loop yields the ROUND-TO-NEAREST quotient (the
//           balanced-ternary natural rounding) with |R| <= B/2, possibly
//           negative; one conditional step (Q-1, R+B) converts to the
//           truncated quotient with 0 <= R < B
//   sign:   q = (sign a != sign b) ? -Q : Q;  r = (a < 0) ? -R : R
//
// Internal datapath is 20 trits: 2|R1| can reach 3B+2, two trits beyond
// the 18-trit operand range. Zero-trit padding preserves values.

`ifndef VTX1_DIV_V
`define VTX1_DIV_V

// 20-trit ripple adder built from the ALU's full-adder cell
module vtx1_trit_addw (
    input  wire [39:0] a,
    input  wire [39:0] b,
    input  wire        negate_b,
    output wire [39:0] sum
);
    function [1:0] tneg(input [1:0] t);
        case (t)
            2'b00:   tneg = 2'b10;
            2'b10:   tneg = 2'b00;
            default: tneg = t;
        endcase
    endfunction

    wire [39:0] bn;
    wire [1:0]  carry [0:20];
    assign carry[0] = 2'b01; // zero trit

    genvar gi;
    generate
        for (gi = 0; gi < 20; gi = gi + 1) begin : chain
            assign bn[2*gi+1 -: 2] = negate_b ? tneg(b[2*gi+1 -: 2])
                                              : b[2*gi+1 -: 2];
            vtx1_trit_fa fa (
                .a(a[2*gi+1 -: 2]), .b(bn[2*gi+1 -: 2]),
                .cin(carry[gi]), .sum(sum[2*gi+1 -: 2]), .cout(carry[gi+1])
            );
        end
    endgenerate
endmodule

module vtx1_div (
    input  wire        clk,
    input  wire        rst_n,

    input  wire        start,     // one-cycle pulse; operands sampled then
    input  wire        is_mod,    // 0: DIV (quotient), 1: MOD (remainder)
    input  wire [35:0] a,         // dividend (18-trit word)
    input  wire [35:0] b,         // divisor, nonzero (core traps on zero)

    output reg         busy,
    output reg         done,      // one-cycle pulse; result valid
    output reg  [35:0] result
);

    localparam [39:0] ZERO_W = 40'h5555555555; // 20 zero trits
    localparam [35:0] ZERO_N = 36'h555555555;  // 18 zero trits

    // ------------------------------------------------------------------
    // Trit helpers (pure functions of their arguments)
    // ------------------------------------------------------------------
    function [35:0] tneg18(input [35:0] w);
        integer i;
        for (i = 0; i < 18; i = i + 1)
            tneg18[2*i+1 -: 2] = (w[2*i+1 -: 2] == 2'b00) ? 2'b10 :
                                 (w[2*i+1 -: 2] == 2'b10) ? 2'b00 :
                                                            w[2*i+1 -: 2];
    endfunction

    function [39:0] tneg20(input [39:0] w);
        integer i;
        for (i = 0; i < 20; i = i + 1)
            tneg20[2*i+1 -: 2] = (w[2*i+1 -: 2] == 2'b00) ? 2'b10 :
                                 (w[2*i+1 -: 2] == 2'b10) ? 2'b00 :
                                                            w[2*i+1 -: 2];
    endfunction

    function [1:0] sign18(input [35:0] w);
        integer i;
        begin
            sign18 = 2'b01;
            for (i = 17; i >= 0; i = i - 1)
                if (sign18 == 2'b01 && w[2*i+1 -: 2] != 2'b01)
                    sign18 = w[2*i+1 -: 2];
        end
    endfunction

    function [1:0] sign20(input [39:0] w);
        integer i;
        begin
            sign20 = 2'b01;
            for (i = 19; i >= 0; i = i - 1)
                if (sign20 == 2'b01 && w[2*i+1 -: 2] != 2'b01)
                    sign20 = w[2*i+1 -: 2];
        end
    endfunction

    // ------------------------------------------------------------------
    // State
    // ------------------------------------------------------------------
    localparam D_IDLE = 3'd0;
    localparam D_PREP = 3'd1;
    localparam D_LOOP = 3'd2;
    localparam D_CORR = 3'd3;
    localparam D_SIGN = 3'd4;

    reg [2:0]  state;
    reg [4:0]  cnt;
    reg        a_neg, b_neg, mod_r;
    reg [35:0] A_sh;              // |a|, consumed MSB-first
    reg [39:0] R, Q, Bm, B2, B3;

    // ------------------------------------------------------------------
    // Loop datapath (combinational on the registered state)
    // ------------------------------------------------------------------
    // R1 = 3R + next dividend trit: one-trit shift with the trit in the
    // low slot - no adder needed.
    wire [39:0] R1 = {R[37:0], A_sh[35:34]};

    wire [1:0]  r1_sgn  = sign20(R1);
    wire        r1_neg  = (r1_sgn == 2'b00);
    wire [39:0] R1a     = r1_neg ? tneg20(R1) : R1;

    wire [39:0] twoR1;
    vtx1_trit_addw u_double (.a(R1a), .b(R1a), .negate_b(1'b0), .sum(twoR1));

    wire [39:0] sA, sB;
    vtx1_trit_addw u_cmp1 (.a(twoR1), .b(Bm), .negate_b(1'b1), .sum(sA));
    vtx1_trit_addw u_cmp3 (.a(twoR1), .b(B3), .negate_b(1'b1), .sum(sB));

    // Digit magnitude: 0 if 2|R1| <= B, 1 if B < 2|R1| <= 3B, else 2.
    wire       ge1 = (sign20(sA) == 2'b10);
    wire       ge2 = (sign20(sB) == 2'b10);
    wire [1:0] qmag = ge2 ? 2'd2 : ge1 ? 2'd1 : 2'd0;

    // R' = R1 - digit*B (digit carries R1's sign, so the magnitude is
    // subtracted when R1 is positive and added when negative)
    wire [39:0] subw = (qmag == 2'd2) ? B2 :
                       (qmag == 2'd1) ? Bm : ZERO_W;
    wire [39:0] Rnext;
    vtx1_trit_addw u_rsub (.a(R1), .b(subw), .negate_b(!r1_neg), .sum(Rnext));

    // Q' = 3Q + digit (digit constant as a trit word; +-2 = (+-1, -+1))
    function [39:0] qconst(input [1:0] mag, input neg);
        begin
            qconst = ZERO_W;
            case (mag)
                2'd1: qconst[1:0] = neg ? 2'b00 : 2'b10;
                2'd2: begin
                    qconst[3:2] = neg ? 2'b00 : 2'b10; // +-1 at weight 3
                    qconst[1:0] = neg ? 2'b10 : 2'b00; // -+1 at weight 1
                end
                default: ;
            endcase
        end
    endfunction

    wire [39:0] Qsh = {Q[37:0], 2'b01};
    wire [39:0] Qnext;
    vtx1_trit_addw u_qacc (.a(Qsh), .b(qconst(qmag, r1_neg)),
                           .negate_b(1'b0), .sum(Qnext));

    // Correction step: negative remainder -> Q-1, R+B (truncation)
    wire        r_neg = (sign20(R) == 2'b00);
    wire [39:0] Rcorr, Qcorr;
    vtx1_trit_addw u_rcorr (.a(R), .b(Bm), .negate_b(1'b0), .sum(Rcorr));
    vtx1_trit_addw u_qcorr (.a(Q), .b({ZERO_W[39:2], 2'b10}),
                            .negate_b(1'b1), .sum(Qcorr));

    // Setup: B2 = 3B - B
    wire [39:0] B2w;
    vtx1_trit_addw u_bsetup (.a(B3), .b(Bm), .negate_b(1'b1), .sum(B2w));

    // ------------------------------------------------------------------
    // FSM
    // ------------------------------------------------------------------
    wire [35:0] a_abs = (sign18(a) == 2'b00) ? tneg18(a) : a;
    wire [35:0] b_abs = (sign18(b) == 2'b00) ? tneg18(b) : b;

    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            state  <= D_IDLE;
            busy   <= 1'b0;
            done   <= 1'b0;
            result <= ZERO_N;
            cnt    <= 5'd0;
            a_neg  <= 1'b0;
            b_neg  <= 1'b0;
            mod_r  <= 1'b0;
            A_sh   <= ZERO_N;
            R      <= ZERO_W;
            Q      <= ZERO_W;
            Bm     <= ZERO_W;
            B2     <= ZERO_W;
            B3     <= ZERO_W;
        end else begin
            done <= 1'b0;
            case (state)
                D_IDLE: if (start) begin
                    busy  <= 1'b1;
                    a_neg <= (sign18(a) == 2'b00);
                    b_neg <= (sign18(b) == 2'b00);
                    mod_r <= is_mod;
                    A_sh  <= a_abs;
                    Bm    <= {4'b0101, b_abs};             // zero-extend
                    B3    <= {2'b01, b_abs, 2'b01};        // 3|b|: trit shift
                    R     <= ZERO_W;
                    Q     <= ZERO_W;
                    cnt   <= 5'd0;
                    state <= D_PREP;
                end

                D_PREP: begin
                    B2    <= B2w;
                    state <= D_LOOP;
                end

                D_LOOP: begin
                    R    <= Rnext;
                    Q    <= Qnext;
                    A_sh <= {A_sh[33:0], 2'b01};
                    cnt  <= cnt + 5'd1;
                    if (cnt == 5'd17)
                        state <= D_CORR;
                end

                D_CORR: begin
                    if (r_neg) begin
                        R <= Rcorr;
                        Q <= Qcorr;
                    end
                    state <= D_SIGN;
                end

                D_SIGN: begin
                    if (mod_r)
                        result <= a_neg ? tneg18(R[35:0]) : R[35:0];
                    else
                        result <= (a_neg ^ b_neg) ? tneg18(Q[35:0]) : Q[35:0];
                    busy  <= 1'b0;
                    done  <= 1'b1;
                    state <= D_IDLE;
                end

                default: state <= D_IDLE;
            endcase
        end
    end

endmodule

`endif // VTX1_DIV_V
