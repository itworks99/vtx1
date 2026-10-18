`timescale 1ns / 1ps
// Verify the study's fresh ternary units (t_neg, t_sign, t_abs, t_cmp3,
// t_scale) against value semantics on random vectors. t_add and t_mul
// wrap the ALU structures already verified by tb_vtx1_alu.

module tb_funcs;

    reg  [35:0] a, b;
    wire [35:0] y_neg, y_abs, y_scale;
    wire [1:0]  y_sign, y_cmp;

    t_neg   u_neg   (.a(a), .y(y_neg));
    t_sign  u_sign  (.a(a), .y(y_sign));
    t_abs   u_abs   (.a(a), .y(y_abs));
    t_cmp3  u_cmp   (.a(a), .b(b), .y(y_cmp));
    t_scale u_scale (.a(a), .y(y_scale));

    function signed [63:0] wordval(input [35:0] w);
        integer i;
        reg signed [2:0] t;
        begin
            wordval = 64'sd0;
            for (i = 0; i < 18; i = i + 1) begin
                case (w[35-2*i -: 2])
                    2'b00:   t = -3'sd1;
                    2'b01:   t = 3'sd0;
                    2'b10:   t = 3'sd1;
                    default: t = 3'sd0;
                endcase
                wordval = wordval * 3 + t;
            end
        end
    endfunction

    function [35:0] b2t(input signed [35:0] v);
        integer i;
        reg signed [35:0] x, r;
        begin
            x = v;
            b2t = 36'h555555555;
            for (i = 0; i < 18; i = i + 1) begin
                r = x % 3;
                x = x / 3;
                if (r == 2)       begin r = -1; x = x + 1; end
                else if (r == -2) begin r = 1;  x = x - 1; end
                b2t[2*i+1 -: 2] = (r == -1) ? 2'b00 : (r == 1) ? 2'b10 : 2'b01;
            end
        end
    endfunction

    integer errors = 0;
    integer k;
    reg signed [35:0] av, bv;
    reg signed [63:0] want;
    localparam signed [35:0] MAXV = 36'sd193710244;

    task chk(input [127:0] name, input signed [63:0] got,
             input signed [63:0] wantv);
        if (got !== wantv) begin
            $display("FAIL %0s: a=%0d b=%0d got %0d want %0d",
                     name, av, bv, got, wantv);
            errors = errors + 1;
        end
    endtask

    initial begin
        for (k = 0; k < 20000; k = k + 1) begin
            av = ($random % MAXV);
            bv = ($random % MAXV);
            if (k == 0) av = 0;
            if (k == 1) begin av = MAXV; bv = -MAXV; end
            a = b2t(av);
            b = b2t(bv);
            #1;
            chk("neg",  wordval(y_neg), -av);
            chk("abs",  wordval(y_abs), (av < 0) ? -av : av);
            chk("sign", (y_sign == 2'b00) ? -64'sd1 :
                        (y_sign == 2'b01) ? 64'sd0 : 64'sd1,
                        (av < 0) ? -64'sd1 : (av == 0) ? 64'sd0 : 64'sd1);
            chk("cmp3", (y_cmp == 2'b00) ? -64'sd1 :
                        (y_cmp == 2'b01) ? 64'sd0 : 64'sd1,
                        (av < bv) ? -64'sd1 : (av == bv) ? 64'sd0 : 64'sd1);
            // scale: nearest(av/3), exact in balanced ternary
            want = (av - ((a[1:0] == 2'b00) ? -64'sd1 :
                          (a[1:0] == 2'b10) ?  64'sd1 : 64'sd0)) / 3;
            chk("scale", wordval(y_scale), want);
        end
        if (errors == 0)
            $display("=== PASS: 20000 vectors, ternary function units correct ===");
        else
            $display("=== FAIL: %0d errors ===", errors);
        $finish;
    end

endmodule
