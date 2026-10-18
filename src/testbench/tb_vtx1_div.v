`timescale 1ns / 1ps
// Testbench: the balanced-ternary iterative divider vs a binary reference.
//
// Directed + random operand pairs; for each, DIV and MOD results are
// compared against Verilog's native signed '/' and '%', which share the
// semantics ratified for VTX1 (and implemented by vtxsim via Go):
// truncation toward zero, remainder sign follows the dividend.
// The reference conversion uses division - testbench only, never RTL.

module tb_vtx1_div;

    reg clk = 0;
    reg rst_n = 0;
    always #5 clk = ~clk;

    reg         start = 0, is_mod = 0;
    reg  [35:0] a, b;
    wire        busy, done;
    wire [35:0] result;

    vtx1_div dut (
        .clk(clk), .rst_n(rst_n),
        .start(start), .is_mod(is_mod), .a(a), .b(b),
        .busy(busy), .done(done), .result(result)
    );

    localparam signed [35:0] MAXV = 36'sd193710244; // (3^18-1)/2

    // Trit word -> signed value
    function signed [35:0] wordval(input [35:0] w);
        integer i;
        reg signed [2:0] t;
        begin
            wordval = 36'sd0;
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

    // Signed value -> trit word (reference encoding; TB only)
    function [35:0] b2t(input signed [35:0] v);
        integer i;
        reg signed [35:0] x;
        reg signed [35:0] r;
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
    integer tests  = 0;
    integer k;
    reg signed [35:0] av, bv, want, got;
    integer guard;

    task run_one(input signed [35:0] ain, input signed [35:0] bin,
                 input mod_in);
        begin
            a      = b2t(ain);
            b      = b2t(bin);
            is_mod = mod_in;
            @(negedge clk);
            start = 1;
            @(negedge clk);
            start = 0;
            guard = 0;
            while (!done && guard < 100) begin
                @(negedge clk);
                guard = guard + 1;
            end
            want = mod_in ? (ain % bin) : (ain / bin);
            got  = wordval(result);
            tests = tests + 1;
            if (!done) begin
                $display("FAIL: no done for a=%0d b=%0d", ain, bin);
                errors = errors + 1;
            end else if (got !== want) begin
                $display("FAIL: %0d %s %0d = %0d, want %0d",
                         ain, mod_in ? "MOD" : "DIV", bin, got, want);
                errors = errors + 1;
            end
        end
    endtask

    task run_pair(input signed [35:0] ain, input signed [35:0] bin);
        begin
            run_one(ain, bin, 1'b0);
            run_one(ain, bin, 1'b1);
        end
    endtask

    reg signed [35:0] ra, rb;
    integer s1, s2;

    initial begin
        rst_n = 0;
        repeat (4) @(posedge clk);
        rst_n = 1;
        @(negedge clk);

        // Directed: boundaries, exact divisions, the sign grid
        run_pair(36'sd0, 36'sd1);
        run_pair(36'sd0, -36'sd1);
        run_pair(36'sd1, 36'sd1);
        run_pair(-36'sd1, 36'sd1);
        run_pair(36'sd7, 36'sd2);
        run_pair(-36'sd7, 36'sd2);
        run_pair(36'sd7, -36'sd2);
        run_pair(-36'sd7, -36'sd2);
        run_pair(36'sd6, 36'sd2);
        run_pair(36'sd6, 36'sd3);
        run_pair(36'sd5, 36'sd3);
        run_pair(36'sd1, 36'sd2);
        run_pair(36'sd2, 36'sd2);
        run_pair(MAXV, 36'sd1);
        run_pair(MAXV, -36'sd1);
        run_pair(-MAXV, 36'sd1);
        run_pair(MAXV, MAXV);
        run_pair(MAXV, 36'sd2);
        run_pair(MAXV, 36'sd3);
        run_pair(-MAXV, -MAXV);
        run_pair(36'sd1, MAXV);
        run_pair(-36'sd1, MAXV);
        run_pair(36'sd100, 36'sd10);
        run_pair(36'sd1000000, 36'sd7);

        // Random sweep (b != 0), full range and small-divisor heavy
        for (k = 0; k < 3000; k = k + 1) begin
            s1 = $random;
            s2 = $random;
            ra = ($random % MAXV);
            if (k % 3 == 0)
                rb = ($random % 36'sd100);     // small divisors
            else
                rb = ($random % MAXV);
            if (rb == 0) rb = 36'sd1;
            run_pair(ra, rb);
        end

        if (errors == 0)
            $display("=== PASS: %0d division vectors match the reference ===",
                     tests);
        else
            $display("=== FAIL: %0d/%0d vectors wrong ===", errors, tests);
        $finish;
    end

endmodule
