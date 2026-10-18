`timescale 1ns/1ps
// Generic SoC benchmark harness: runs any assembled program (+img=<memh>)
// on the FULL vtx1_soc (bus matrix, memory controller, MMIO peripherals,
// external SRAM model), waits for WFI, and verifies up to six result
// words at 0x1200 (+nres/+exp0..5).
//
// Programs can timestamp phases by storing to the marker addresses
// 0x1210..0x121F: the harness snoops those stores on the external memory
// pins and prints the CPU cycle counter at each, so multi-phase programs
// (like the software-vs-accelerator matvec) report per-phase cycles.

module tb_vtx1_soc_bench;
    reg clk_ext = 0;
    reg clk_32k = 0;
    reg rst_ext_n = 0;
    reg pwr_good = 1;
    reg uart_rx = 1, uart_cts = 0;
    reg jtag_tck = 0, jtag_tms = 0, jtag_tdi = 0, jtag_trst_n = 1;
    reg wake_event = 0, flash_miso = 0;
    reg spi_miso = 0, spi_ss_n = 1, spi_sclk_in = 0, spi_mosi_in = 0;

    wire i2c_sda, i2c_scl;
    pullup(i2c_sda);
    pullup(i2c_scl);

    wire [35:0] mem_addr;
    wire [35:0] mem_data;
    wire mem_we_n, mem_oe_n, mem_cs_n;

    vtx1_soc uut (
        .clk_ext(clk_ext), .clk_32k(clk_32k), .rst_ext_n(rst_ext_n),
        .mem_addr(mem_addr), .mem_data(mem_data),
        .mem_we_n(mem_we_n), .mem_oe_n(mem_oe_n), .mem_cs_n(mem_cs_n),
        .flash_miso(flash_miso),
        .uart_rx(uart_rx), .uart_cts(uart_cts),
        .spi_miso(spi_miso), .spi_ss_n(spi_ss_n),
        .spi_sclk_in(spi_sclk_in), .spi_mosi_in(spi_mosi_in),
        .i2c_sda(i2c_sda), .i2c_scl(i2c_scl),
        .jtag_tck(jtag_tck), .jtag_tms(jtag_tms), .jtag_tdi(jtag_tdi),
        .jtag_trst_n(jtag_trst_n),
        .pwr_good(pwr_good), .wake_event(wake_event)
    );

    // External SRAM model (16K x 36), loaded with the program image
    reg [35:0] sram [0:16383];

    assign mem_data = (!mem_cs_n && !mem_oe_n && mem_we_n)
                      ? sram[mem_addr[13:0]] : 36'hzzzzzzzzz;

    always @(posedge clk_ext)
        if (!mem_cs_n && !mem_we_n)
            sram[mem_addr[13:0]] <= mem_data;

    initial forever #18.52 clk_ext = ~clk_ext; // 27 MHz
    initial forever #15259 clk_32k = ~clk_32k;

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

    // Marker snooping: stores to 0x1210..0x121F timestamp phases
    integer marker_cycles [0:15];
    integer last_marker;
    initial last_marker = -1;
    always @(posedge clk_ext) begin
        if (!mem_cs_n && !mem_we_n && mem_addr[13:4] == 10'h121) begin
            marker_cycles[mem_addr[3:0]] = uut.cpu_inst.cycles;
            if ($signed(mem_addr[3:0]) > last_marker)
                last_marker = mem_addr[3:0];
        end
    end

    reg [1023:0] img;
    integer nres, errors, i, timeout, dummy, prev;
    integer exp0, exp1, exp2, exp3, exp4, exp5;
    integer exp [0:5];
    reg signed [35:0] got;

    initial begin
        errors = 0;
        if (!$value$plusargs("img=%s", img)) begin
            $display("FAIL: missing +img=<memh file>");
            $finish;
        end
        if (!$value$plusargs("nres=%d", nres)) nres = 1;
        exp0 = 0; exp1 = 0; exp2 = 0; exp3 = 0; exp4 = 0; exp5 = 0;
        dummy = $value$plusargs("exp0=%d", exp0);
        dummy = $value$plusargs("exp1=%d", exp1);
        dummy = $value$plusargs("exp2=%d", exp2);
        dummy = $value$plusargs("exp3=%d", exp3);
        dummy = $value$plusargs("exp4=%d", exp4);
        dummy = $value$plusargs("exp5=%d", exp5);
        exp[0] = exp0; exp[1] = exp1; exp[2] = exp2;
        exp[3] = exp3; exp[4] = exp4; exp[5] = exp5;

        $readmemh(img, sram);
        #1;
        // Preload the instruction TCM banks and the data scratchpad
        for (i = 0; i < 4096; i = i + 4) begin
            uut.tcm_inst.bank0[i >> 2] = sram['h1000 + i];
            uut.tcm_inst.bank1[i >> 2] = sram['h1000 + i + 1];
            uut.tcm_inst.bank2[i >> 2] = sram['h1000 + i + 2];
            uut.tcm_inst.bank3[i >> 2] = sram['h1000 + i + 3];
        end
        for (i = 0; i < 16384; i = i + 1)
            uut.spad_inst.mem[i] = sram[i];

        #100 rst_ext_n = 1;

        timeout = 0;
        while (!uut.cpu_inst.halted && timeout < 400000) begin
            @(posedge clk_ext);
            timeout = timeout + 1;
        end
        #200;

        if (!uut.cpu_inst.halted) begin
            $display("FAIL: CPU did not halt (pc=0x%0h)", uut.cpu_inst.debug_pc);
            errors = errors + 1;
        end
        if (uut.cpu_inst.trap) begin
            $display("FAIL: CPU trapped at PC 0x%0h", uut.cpu_inst.trap_pc);
            errors = errors + 1;
        end
        for (i = 0; i < nres; i = i + 1) begin
            got = wordval(sram['h1200 + i]);
            if (got !== exp[i]) begin
                $display("FAIL: result[%0d] = %0d, want %0d", i, got, exp[i]);
                errors = errors + 1;
            end
        end

        // Phase report from the markers
        prev = 0;
        for (i = 0; i <= last_marker; i = i + 1) begin
            $display("phase %0d: %0d cycles (marker 0x121%0x at cycle %0d)",
                     i, marker_cycles[i] - prev, i, marker_cycles[i]);
            prev = marker_cycles[i];
        end

        if (errors == 0)
            $display("SoC-bench %0s: cycles=%0d retired=%0d PASS",
                     img, uut.cpu_inst.cycles, uut.cpu_inst.retired);
        else
            $display("SoC-bench %0s: FAIL (%0d errors)", img, errors);
        $finish;
    end

endmodule
