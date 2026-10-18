`timescale 1ns / 1ps
// Benchmark harness: one Tier A kernel on the PicoRV32 RTL (the binary
// control core of the cycle-level comparison).
//
// Loads an rvhex memory image (+img=<path>, one 32-bit word per line,
// code at the reset vector 0x0, data vectors at 0x1100/0x1140), runs
// PicoRV32 until the kernel's EBREAK raises `trap`, then verifies up to
// three result words at byte address 0x1200 (+nres/+exp0..2) and prints
// the cycle count for the comparison table.
//
// Memory model (stated in the table, deliberately generous to the
// control core): single-cycle memory - mem_ready is tied high with
// combinational read data, PicoRV32's fastest native configuration.
// The core itself is configured generously too: ENABLE_FAST_MUL (one
// cycle multiply) and BARREL_SHIFTER (one cycle shifts).

module tb_picorv32_bench;

    reg clk = 0;
    reg resetn = 0;
    always #5 clk = ~clk;

    wire        mem_valid, mem_instr;
    wire [31:0] mem_addr, mem_wdata;
    wire [3:0]  mem_wstrb;
    wire        trap;

    // Single-cycle memory: always ready, combinational read
    reg [31:0] memory [0:4095]; // 16 KB, byte addresses 0x0000..0x3FFF
    wire [31:0] mem_rdata = memory[mem_addr[13:2]];
    wire        mem_ready = 1'b1;

    always @(posedge clk) begin
        if (mem_valid && mem_ready) begin
            if (mem_wstrb[0]) memory[mem_addr[13:2]][7:0]   <= mem_wdata[7:0];
            if (mem_wstrb[1]) memory[mem_addr[13:2]][15:8]  <= mem_wdata[15:8];
            if (mem_wstrb[2]) memory[mem_addr[13:2]][23:16] <= mem_wdata[23:16];
            if (mem_wstrb[3]) memory[mem_addr[13:2]][31:24] <= mem_wdata[31:24];
        end
    end

    picorv32 #(
        .ENABLE_MUL(0),
        .ENABLE_FAST_MUL(1),   // single-cycle multiplier (generous)
        .BARREL_SHIFTER(1),    // single-cycle shifts (generous)
        .ENABLE_COUNTERS(0),
        .PROGADDR_RESET(32'h0)
    ) cpu (
        .clk(clk), .resetn(resetn),
        .trap(trap),
        .mem_valid(mem_valid),
        .mem_instr(mem_instr),
        .mem_ready(mem_ready),
        .mem_addr(mem_addr),
        .mem_wdata(mem_wdata),
        .mem_wstrb(mem_wstrb),
        .mem_rdata(mem_rdata),
        // PCPI / IRQ unused
        .pcpi_wr(1'b0), .pcpi_rd(32'b0), .pcpi_wait(1'b0), .pcpi_ready(1'b0),
        .irq(32'b0)
    );

    reg [1023:0] img;
    integer nres, timeout, errors, i, dummy;
    integer exp0, exp1, exp2;
    integer exp [0:2];
    integer cycles;
    reg signed [31:0] got;

    initial begin
        errors = 0;
        if (!$value$plusargs("img=%s", img)) begin
            $display("FAIL: missing +img=<hex file>");
            $finish;
        end
        if (!$value$plusargs("nres=%d", nres)) nres = 1;
        exp0 = 0; exp1 = 0; exp2 = 0;
        dummy = $value$plusargs("exp0=%d", exp0);
        dummy = $value$plusargs("exp1=%d", exp1);
        dummy = $value$plusargs("exp2=%d", exp2);
        exp[0] = exp0; exp[1] = exp1; exp[2] = exp2;

        $readmemh(img, memory);

        resetn = 0;
        repeat (4) @(posedge clk);
        resetn = 1;

        cycles = 0;
        timeout = 0;
        while (!trap && timeout < 1000000) begin
            @(posedge clk);
            cycles = cycles + 1;
            timeout = timeout + 1;
        end

        if (!trap) begin
            $display("FAIL: PicoRV32 did not reach ebreak");
            errors = errors + 1;
        end
        for (i = 0; i < nres; i = i + 1) begin
            got = memory[('h1200 >> 2) + i];
            if (got !== exp[i]) begin
                $display("FAIL: result[%0d] = %0d, want %0d", i, got, exp[i]);
                errors = errors + 1;
            end
        end

        if (errors == 0)
            $display("PicoRV32 %0s: cycles=%0d PASS", img, cycles);
        else
            $display("PicoRV32 %0s: FAIL (%0d errors)", img, errors);
        $finish;
    end

endmodule
