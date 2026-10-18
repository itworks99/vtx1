`timescale 1ns / 1ps
// Testbench: the ternary-weight MAC array (vtx1_tmac) against its binary
// control (bin_tmac) and a golden model.
//
// Random weight words and activation vectors over many accumulation
// runs: on non-overflowing data both units must equal the golden signed
// accumulation AND each other. A directed check verifies vtx1_tmac's
// wrap semantics (mod 3^18 via dropped carries - identical to a chain of
// CPU ADD instructions on the same data).

module tb_tmac;

    reg clk = 0;
    reg rst_n = 0;
    always #5 clk = ~clk;

    reg         clear = 0, step = 0;
    reg  [31:0] weights;
    reg  [575:0] acts_t;
    reg  [479:0] acts_b;
    wire [35:0] acc_t;
    wire signed [29:0] acc_b;

    vtx1_tmac dut_t (.clk(clk), .rst_n(rst_n), .clear(clear), .step(step),
                     .weights(weights), .acts(acts_t), .acc(acc_t));
    bin_tmac  dut_b (.clk(clk), .rst_n(rst_n), .clear(clear), .step(step),
                     .weights(weights), .acts(acts_b), .acc(acc_b));

    // Trit word -> signed value
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

    // Balanced wrap to 18 trits (golden model of the dropped-carry
    // semantics): value mod 3^18, centered on zero
    localparam signed [63:0] P18 = 64'sd387420489; // 3^18
    function signed [63:0] wrap18(input signed [63:0] v);
        begin
            wrap18 = v % P18;
            if (wrap18 >  (P18 - 1) / 2) wrap18 = wrap18 - P18;
            if (wrap18 < -(P18 - 1) / 2) wrap18 = wrap18 + P18;
        end
    endfunction

    integer errors = 0;
    integer tests  = 0;
    integer run, k, i;
    reg signed [35:0] a_val [0:15];
    reg signed [1:0]  w_val [0:15];
    reg signed [63:0] golden;
    integer wsel;
    integer nsteps;
    integer amax;

    task do_step;
        begin
            @(negedge clk);
            step = 1;
            @(negedge clk);
            step = 0;
        end
    endtask

    task do_clear;
        begin
            @(negedge clk);
            clear = 1;
            @(negedge clk);
            clear = 0;
        end
    endtask

    initial begin
        rst_n = 0;
        repeat (4) @(posedge clk);
        rst_n = 1;

        // ------------------------------------------------------------
        // Random equivalence runs (no overflow): both units == golden
        // ------------------------------------------------------------
        for (run = 0; run < 60; run = run + 1) begin
            do_clear;
            golden = 0;
            // keep |acc| < 2^29 and < (3^18-1)/2: 8 steps x 16 lanes
            // x 1e6 max < 1.3e8
            nsteps = 8;
            amax   = 1000000;
            for (k = 0; k < nsteps; k = k + 1) begin
                for (i = 0; i < 16; i = i + 1) begin
                    a_val[i] = ($random % amax);
                    wsel = $random % 3;
                    w_val[i] = (wsel == 1) ? 2'sd1 :
                               (wsel == -1 || wsel == 2) ? -2'sd1 : 2'sd0;
                    acts_t[36*i +: 36] = b2t(a_val[i]);
                    acts_b[30*i +: 30] = a_val[i][29:0];
                    weights[2*i +: 2]  = (w_val[i] == 1)  ? 2'b10 :
                                         (w_val[i] == -1) ? 2'b00 : 2'b01;
                    golden = golden + w_val[i] * a_val[i];
                end
                do_step;
            end
            tests = tests + 1;
            if (wordval(acc_t) !== golden) begin
                $display("FAIL run %0d: tmac acc = %0d, golden %0d",
                         run, wordval(acc_t), golden);
                errors = errors + 1;
            end
            if (acc_b !== golden[29:0]) begin
                $display("FAIL run %0d: bin acc = %0d, golden %0d",
                         run, acc_b, golden);
                errors = errors + 1;
            end
        end

        // ------------------------------------------------------------
        // Directed: ternary wrap semantics (mod 3^18, like CPU ADDs)
        // ------------------------------------------------------------
        do_clear;
        golden = 0;
        for (k = 0; k < 40; k = k + 1) begin
            for (i = 0; i < 16; i = i + 1) begin
                a_val[i] = 36'sd193710244; // max word value
                acts_t[36*i +: 36] = b2t(a_val[i]);
                weights[2*i +: 2]  = 2'b10; // all +1
                golden = golden + a_val[i];
            end
            do_step;
        end
        tests = tests + 1;
        if (wordval(acc_t) !== wrap18(golden)) begin
            $display("FAIL wrap: tmac acc = %0d, want %0d (golden %0d)",
                     wordval(acc_t), wrap18(golden), golden);
            errors = errors + 1;
        end

        if (errors == 0)
            $display("=== PASS: %0d MAC runs, ternary == binary == golden ===",
                     tests);
        else
            $display("=== FAIL: %0d error(s) ===", errors);
        $finish;
    end

endmodule
