`timescale 1ns / 1ps
// Testbench: native balanced-ternary ALU vs a binary reference model
//
// Drives the ALU with directed edge cases and random trit words, and checks
// every operation against a reference computed in plain binary arithmetic
// (division is fine HERE - the whole point is that it exists only in the
// testbench, never in the datapath).

`include "vtx1_isa_generated.vh"

module tb_vtx1_alu;

    localparam [35:0] ZERO_WORD = 36'h555555555;
    localparam signed [63:0] WORD_MOD = 64'sd387420489; // 3^18
    localparam signed [63:0] WORD_MAX = (WORD_MOD - 1) / 2;

    reg  [35:0] a, b;
    reg  [6:0]  opcode;
    wire [35:0] result;
    wire        unsupported;

    vtx1_alu dut (.a(a), .b(b), .opcode(opcode),
                  .result(result), .unsupported(unsupported));

    integer errors;
    integer vectors;

    // --- Reference helpers (binary; simulation only) ---
    function signed [2:0] tv(input [1:0] t);
        case (t)
            2'b00:   tv = -3'sd1;
            2'b10:   tv = 3'sd1;
            default: tv = 3'sd0;
        endcase
    endfunction

    function signed [63:0] t2b(input [35:0] w);
        integer i;
        begin
            t2b = 0;
            for (i = 17; i >= 0; i = i - 1)
                t2b = t2b * 3 + tv(w[2*i+1 -: 2]);
        end
    endfunction

    // Balanced wrap modulo 3^18 (the ISA simulator's WrapWord)
    function signed [63:0] wrap(input signed [63:0] v);
        begin
            wrap = ((v % WORD_MOD) + WORD_MOD) % WORD_MOD;
            if (wrap > WORD_MAX) wrap = wrap - WORD_MOD;
        end
    endfunction

    function signed [63:0] sgn(input signed [63:0] v);
        sgn = (v > 0) ? 1 : (v < 0) ? -1 : 0;
    endfunction

    // Tritwise reference ops as values
    function signed [63:0] tw_ref(input [35:0] x, input [35:0] y,
                                  input [1:0] which); // 0 min, 1 max, 2 mul
        integer i;
        reg signed [2:0] xi, yi, r;
        begin
            tw_ref = 0;
            for (i = 17; i >= 0; i = i - 1) begin
                xi = tv(x[2*i+1 -: 2]);
                yi = tv(y[2*i+1 -: 2]);
                case (which)
                    2'd0: r = (xi < yi) ? xi : yi;
                    2'd1: r = (xi > yi) ? xi : yi;
                    default: r = xi * yi;
                endcase
                tw_ref = tw_ref * 3 + r;
            end
        end
    endfunction

    // Shift reference: amount = low 4 trits of b, mod 18
    function signed [63:0] shift_ref(input [35:0] x, input [35:0] y,
                                     input left, input rotate);
        integer i, k, src;
        reg signed [2:0] trits [0:17]; // index 0 = most significant
        reg signed [63:0] v;
        reg signed [31:0] m;
        begin
            for (i = 0; i < 18; i = i + 1)
                trits[i] = tv(x[2*(17-i)+1 -: 2]);
            m = 27*tv(y[7:6]) + 9*tv(y[5:4]) + 3*tv(y[3:2]) + tv(y[1:0]);
            k = ((m % 18) + 18) % 18;
            v = 0;
            for (i = 0; i < 18; i = i + 1) begin
                src = left ? i + k : i - k;
                if (src >= 0 && src < 18)
                    v = v * 3 + trits[src];
                else if (rotate)
                    v = v * 3 + trits[((src % 18) + 18) % 18];
                else
                    v = v * 3;
            end
            shift_ref = v;
        end
    endfunction

    // Random canonical trit word
    function [35:0] rand_word;
        integer i, r;
        begin
            for (i = 0; i < 18; i = i + 1) begin
                r = $urandom % 3;
                rand_word[2*i+1 -: 2] = (r == 0) ? 2'b00 : (r == 1) ? 2'b01 : 2'b10;
            end
        end
    endfunction

    task check(input [6:0] op, input [127:0] name, input signed [63:0] want_v);
        begin
            opcode = op; #1;
            vectors = vectors + 1;
            if (unsupported) begin
                $display("  FAIL: %0s reported unsupported", name);
                errors = errors + 1;
            end else if (t2b(result) !== want_v) begin
                $display("  FAIL: %0s(a=%0d, b=%0d) = %0d, want %0d",
                         name, t2b(a), t2b(b), t2b(result), want_v);
                errors = errors + 1;
            end
        end
    endtask

    task check_all_ops;
        begin
            check(`VTX1_OPV_ADD, "ADD", wrap(t2b(a) + t2b(b)));
            check(`VTX1_OPV_SUB, "SUB", wrap(t2b(a) - t2b(b)));
            check(`VTX1_OPV_MUL, "MUL", wrap(t2b(a) * t2b(b)));
            check(`VTX1_OPV_CMP, "CMP", sgn(t2b(a) - t2b(b)));
            check(`VTX1_OPV_TEST, "TEST", sgn(tw_ref(a, b, 2'd0)));
            check(`VTX1_OPV_AND, "AND", tw_ref(a, b, 2'd0));
            check(`VTX1_OPV_OR,  "OR",  tw_ref(a, b, 2'd1));
            check(`VTX1_OPV_XOR, "XOR", tw_ref(a, b, 2'd2));
            check(`VTX1_OPV_NOT, "NOT", -t2b(a));
            check(`VTX1_OPV_NEG, "NEG", -t2b(a));
            check(`VTX1_OPV_ABS, "ABS", (t2b(a) < 0) ? -t2b(a) : t2b(a));
            check(`VTX1_OPV_SHL, "SHL", shift_ref(a, b, 1'b1, 1'b0));
            check(`VTX1_OPV_SHR, "SHR", shift_ref(a, b, 1'b0, 1'b0));
            check(`VTX1_OPV_ROL, "ROL", shift_ref(a, b, 1'b1, 1'b1));
            check(`VTX1_OPV_ROR, "ROR", shift_ref(a, b, 1'b0, 1'b1));
            check(`VTX1_OPV_INC, "INC", wrap(t2b(a) + t2b(b))); // core passes b=+1
            check(`VTX1_OPV_DEC, "DEC", wrap(t2b(a) - t2b(b)));
        end
    endtask

    // Directed values: zero, +-1, extremes, and a few structured words
    reg [35:0] directed [0:5];
    integer i, j;

    initial begin
        errors = 0;
        vectors = 0;

        directed[0] = ZERO_WORD;                 // 0
        directed[1] = 36'h555555556;             // +1
        directed[2] = 36'h555555554;             // -1
        directed[3] = {18{2'b10}};               // WORD_MAX (all +1)
        directed[4] = {18{2'b00}};               // -WORD_MAX (all -1)
        directed[5] = 36'h565154565;             // arbitrary mixed word

        $display("=== native trit ALU vs binary reference ===");

        // Directed pairs
        for (i = 0; i < 6; i = i + 1)
            for (j = 0; j < 6; j = j + 1) begin
                a = directed[i]; b = directed[j];
                check_all_ops;
            end

        // Random vectors
        for (i = 0; i < 2000; i = i + 1) begin
            a = rand_word(); b = rand_word();
            check_all_ops;
        end

        // DIV must report unsupported (it belongs to the divider FSM)
        a = directed[5]; b = directed[1];
        opcode = `VTX1_OPV_DIV; #1;
        if (!unsupported) begin
            $display("  FAIL: DIV should be unsupported in the combinational ALU");
            errors = errors + 1;
        end

        if (errors == 0)
            $display("=== PASS: %0d vectors, ALU matches the reference ===", vectors);
        else
            $display("=== FAIL: %0d error(s) in %0d vectors ===", errors, vectors);
        $finish;
    end

endmodule
