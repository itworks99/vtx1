`timescale 1ns / 1ps
// Benchmark harness: one Tier A kernel on the VTX1 RTL core.
//
// Loads a vtxasm memory image (+img=<path>), runs the triple-issue core
// with its tightly-coupled instruction memory and a plain synchronous
// data RAM (2-cycle ready handshake) until WFI, then verifies up to three
// result words at 0x1200 (+nres=<n> +exp0/+exp1/+exp2=<signed value>)
// and prints cycles / retired words / CPI for the comparison table.
//
// Memory model (stated in the table): instruction fetch from the TCM
// (1 cycle, 3 words); data loads/stores through the handshake RAM
// (2 cycles per access).

module tb_vtx1_bench;

    reg clk = 0;
    reg rst_n = 0;
    always #5 clk = ~clk;

    // Data memory: synchronous RAM with a one-cycle ready handshake
    reg [35:0] mem [0:16383];

    wire        mem_req, mem_wr;
    wire [31:0] mem_addr;
    wire [35:0] mem_wdata;
    reg  [35:0] mem_rdata;
    reg         mem_ready = 0;

    always @(posedge clk) begin
        mem_ready <= mem_req && !mem_ready;
        if (mem_req && !mem_wr)
            mem_rdata <= mem[mem_addr[13:0]];
        if (mem_req && mem_wr)
            mem[mem_addr[13:0]] <= mem_wdata;
    end

    // Instruction fetch fast path
    wire        imem_en, imem_hit;
    wire [31:0] imem_addr;
    wire [35:0] imem_rdata0, imem_rdata1, imem_rdata2;

    vtx1_tcm #(.BASE(32'h1000)) tcm (
        .clk(clk), .en(imem_en), .addr(imem_addr),
        .hit(imem_hit),
        .rdata0(imem_rdata0), .rdata1(imem_rdata1), .rdata2(imem_rdata2),
        .wen(mem_req && mem_wr), .waddr(mem_addr), .wdata(mem_wdata)
    );

    wire        halted, trap;
    wire [31:0] trap_pc, retired, cycles;

    vtx1_core core (
        .clk(clk), .rst_n(rst_n),
        .mem_req(mem_req), .mem_wr(mem_wr),
        .mem_addr(mem_addr), .mem_wdata(mem_wdata), .mem_rdata(mem_rdata),
        .mem_ready(mem_ready),
        .imem_en(imem_en), .imem_addr(imem_addr),
        .imem_rdata0(imem_rdata0), .imem_rdata1(imem_rdata1),
        .imem_rdata2(imem_rdata2), .imem_hit(imem_hit),
        .halted(halted), .trap(trap), .trap_pc(trap_pc),
        .retired(retired), .cycles(cycles), .debug_pc()
    );

    // Decode a canonical trit word to its balanced value
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

    reg [1023:0] img;
    integer nres, timeout, errors, i, dummy;
    integer exp0, exp1, exp2;
    integer exp [0:2];
    reg signed [35:0] got;

    initial begin
        errors = 0;
        if (!$value$plusargs("img=%s", img)) begin
            $display("FAIL: missing +img=<memh file>");
            $finish;
        end
        if (!$value$plusargs("nres=%d", nres)) nres = 1;
        exp0 = 0; exp1 = 0; exp2 = 0;
        dummy = $value$plusargs("exp0=%d", exp0);
        dummy = $value$plusargs("exp1=%d", exp1);
        dummy = $value$plusargs("exp2=%d", exp2);
        exp[0] = exp0; exp[1] = exp1; exp[2] = exp2;

        $readmemh(img, mem);
        // Preload the TCM banks (word w -> bank w[1:0], index w[11:2])
        for (i = 0; i < 4096; i = i + 4) begin
            tcm.bank0[i >> 2] = mem['h1000 + i];
            tcm.bank1[i >> 2] = mem['h1000 + i + 1];
            tcm.bank2[i >> 2] = mem['h1000 + i + 2];
            tcm.bank3[i >> 2] = mem['h1000 + i + 3];
        end

        rst_n = 0;
        repeat (4) @(posedge clk);
        rst_n = 1;

        timeout = 0;
        while (!halted && timeout < 100000) begin
            @(posedge clk);
            timeout = timeout + 1;
        end

        if (!halted) begin
            $display("FAIL: core did not halt (pc=0x%0h)", core.debug_pc);
            errors = errors + 1;
        end
        if (trap) begin
            $display("FAIL: core trapped at PC 0x%0h", trap_pc);
            errors = errors + 1;
        end
        for (i = 0; i < nres; i = i + 1) begin
            got = wordval(mem['h1200 + i]);
            if (got !== exp[i]) begin
                $display("FAIL: result[%0d] = %0d, want %0d", i, got, exp[i]);
                errors = errors + 1;
            end
        end

        if (errors == 0)
            $display("VTX1-RTL %0s: cycles=%0d retired=%0d CPI=%0d.%02d PASS",
                     img, cycles, retired,
                     cycles / retired, (cycles * 100 / retired) % 100);
        else
            $display("VTX1-RTL %0s: FAIL (%0d errors)", img, errors);
        $finish;
    end

endmodule
