// VTX1 SoC Testbench
//
// Phase 1 checks the consolidated interfaces (clocks, reset, MMIO routing,
// LEDs, peripheral status).
// Phase 2 is the system milestone: hello_world (real assembler output,
// loaded into the external SRAM model) executes on the trit-native CPU
// through the bus matrix and memory controller, and the testbench verifies
// the program's effects in SRAM.

`timescale 1ns/1ps

module tb_vtx1_soc;    // Test signals
    reg clk_ext;
    reg clk_32k;
    reg rst_ext_n;
    reg pwr_good;
    reg uart_rx;
    reg uart_cts;
    reg jtag_tck, jtag_tms, jtag_tdi, jtag_trst_n;
    reg wake_event;
    reg flash_miso;
    // SPI signals
    reg spi_miso;
    reg spi_ss_n;
    reg spi_sclk_in;
    reg spi_mosi_in;
    // I2C signals (pulled up)
    wire i2c_sda;
    wire i2c_scl;
    reg i2c_sda_drive, i2c_scl_drive;
    reg i2c_sda_val, i2c_scl_val;
    
    wire [35:0] mem_addr;
    wire [35:0] mem_data;
    wire mem_we_n, mem_oe_n, mem_cs_n;
    wire flash_cs_n, flash_sck, flash_mosi, flash_wp_n, flash_hold_n;
    wire [23:0] gpio_pins;
    wire uart_tx, uart_rts;
    wire jtag_tdo;
    wire pwr_enable, sleep_ack;
    wire led_power, led_activity, led_error;
    wire tp_clk_cpu, tp_clk_mem, tp_rst_n;
    wire [3:0] tp_state;
    // SPI outputs
    wire spi_sclk, spi_mosi, spi_miso_out;
    wire [7:0] spi_cs_n;
    
    // I2C bidirectional with pullups
    assign i2c_sda = i2c_sda_drive ? i2c_sda_val : 1'bz;
    assign i2c_scl = i2c_scl_drive ? i2c_scl_val : 1'bz;
    pullup(i2c_sda);
    pullup(i2c_scl);    // Instantiate the consolidated VTX1 system
    vtx1_soc uut (
        .clk_ext(clk_ext),
        .clk_32k(clk_32k),
        .rst_ext_n(rst_ext_n),
        .mem_addr(mem_addr),
        .mem_data(mem_data),
        .mem_we_n(mem_we_n),
        .mem_oe_n(mem_oe_n),
        .mem_cs_n(mem_cs_n),
        .flash_cs_n(flash_cs_n),
        .flash_sck(flash_sck),
        .flash_mosi(flash_mosi),
        .flash_miso(flash_miso),
        .flash_wp_n(flash_wp_n),
        .flash_hold_n(flash_hold_n),
        .gpio_pins(gpio_pins),
        .uart_rx(uart_rx),
        .uart_tx(uart_tx),
        .uart_rts(uart_rts),
        .uart_cts(uart_cts),
        // Enhanced SPI Interface
        .spi_sclk(spi_sclk),
        .spi_mosi(spi_mosi),
        .spi_miso(spi_miso),
        .spi_cs_n(spi_cs_n),
        .spi_ss_n(spi_ss_n),
        .spi_sclk_in(spi_sclk_in),
        .spi_mosi_in(spi_mosi_in),
        .spi_miso_out(spi_miso_out),
        // Enhanced I2C Interface
        .i2c_sda(i2c_sda),
        .i2c_scl(i2c_scl),
        .jtag_tck(jtag_tck),
        .jtag_tms(jtag_tms),
        .jtag_tdi(jtag_tdi),
        .jtag_tdo(jtag_tdo),
        .jtag_trst_n(jtag_trst_n),
        .pwr_good(pwr_good),
        .pwr_enable(pwr_enable),
        .wake_event(wake_event),
        .sleep_ack(sleep_ack),
        .led_power(led_power),
        .led_activity(led_activity),
        .led_error(led_error),
        .tp_clk_cpu(tp_clk_cpu),
        .tp_clk_mem(tp_clk_mem),
        .tp_rst_n(tp_rst_n),
        .tp_state(tp_state)
    );

    // ------------------------------------------------------------------
    // External SRAM model (16K x 36), word-addressed, async read through
    // the tri-state data bus. Loaded with the assembled hello_world image.
    // ------------------------------------------------------------------
    reg [35:0] sram [0:16383];
    initial $readmemh("build/hello_world.memh", sram);

    // Preload the SoC's tightly-coupled instruction RAM with the program
    // region (word 0x1000 -> TCM index 0), standing in for the boot flow.
    integer tci;
    initial begin
        #1; // after the SRAM image loads
        for (tci = 0; tci < 4096; tci = tci + 1)
            uut.tcm_inst.mem[tci] = sram['h1000 + tci];
    end

    assign mem_data = (!mem_cs_n && !mem_oe_n && mem_we_n)
                      ? sram[mem_addr[13:0]] : 36'hzzzzzzzzz;

    always @(posedge clk_ext)
        if (!mem_cs_n && !mem_we_n)
            sram[mem_addr[13:0]] <= mem_data;

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

    integer errors = 0;
    integer wi;
    reg [8*16-1:0] hello = "Hello, VTX1!";
    reg [8*16-1:0] got;
    reg signed [35:0] wv;

    // Initialization
    initial begin
        // System initialization
        clk_ext = 0;
        clk_32k = 0;
        rst_ext_n = 0;
        pwr_good = 1;
        
        // UART initialization
        uart_rx = 1;
        uart_cts = 0;
        
        // SPI initialization
        spi_miso = 0;
        spi_ss_n = 1;
        spi_sclk_in = 0;
        spi_mosi_in = 0;
        
        // I2C initialization (released - pullups will pull high)
        i2c_sda_drive = 0;
        i2c_scl_drive = 0;
        i2c_sda_val = 1;
        i2c_scl_val = 1;
        
        // JTAG initialization
        jtag_tck = 0;
        jtag_tms = 0;
        jtag_tdi = 0;
        jtag_trst_n = 1;
        
        // Flash initialization
        flash_miso = 0;
        wake_event = 0;
        
        // Release reset after some time
        #100 rst_ext_n = 1;
    end

    // Clock generation
    initial begin
        forever #18.52 clk_ext = ~clk_ext; // 27 MHz
    end

    initial begin
        clk_32k = 0;
        forever #15259 clk_32k = ~clk_32k; // 32.768 kHz
    end    // Test sequence
    initial begin
        $display("=================================================");
        $display("VTX1 Consolidated Interface Test - Enhanced Peripherals");
        $display("=================================================");
        
        // Add simulation timeout to prevent hanging
        fork
            begin: timeout_block
                #2000000; // 2 ms: covers program execution through the bus
                $display("FAIL: TIMEOUT - simulation exceeded maximum time limit");
                $display("CPU state: halted=%b trap=%b pc=0x%0h retired=%0d",
                         uut.cpu_inst.halted, uut.cpu_inst.trap,
                         uut.cpu_inst.debug_pc, uut.cpu_inst.retired);
                $finish;
            end
            begin: main_test_block
                // Wait for system stabilization
                wait(rst_ext_n);
                #1000;
                
                // Test enhanced peripheral operation
                $display("Testing Enhanced Peripheral Controllers:");
                $display("- GPIO Controller: Ready");
                $display("- UART Controller: Ready");  
                $display("- SPI Controller: Ready");
                $display("- I2C Controller: Ready");
                
                // Test SPI signals
                $display("SPI Interface Status:");
                $display("  SCLK: %b, CS_N[0]: %b", spi_sclk, spi_cs_n[0]);
                
                // Test I2C signals  
                $display("I2C Interface Status:");
                $display("  SDA: %b, SCL: %b", i2c_sda, i2c_scl);
                
                // Test LED status
                $display("System Status LEDs:");
                $display("  Power: %b, Activity: %b, Error: %b", 
                        led_power, led_activity, led_error);
                
                // Test memory interface
                $display("Memory Interface Status:");
                $display("  CS_N: %b, WE_N: %b, OE_N: %b", mem_cs_n, mem_we_n, mem_oe_n);
                
                // Wait and observe
                #10000;
                #1000;
                
                // Check basic functionality
                $display("Power LED: %b", led_power);
                $display("Activity LED: %b", led_activity);
                $display("Error LED: %b", led_error);
                $display("Test Points State: %h", tp_state);                $display("Reset Status: %b", tp_rst_n);
                  // Check interrupt status
                $display("Interrupt Status:");
                $display("  SPI IRQ: %b", uut.spi_irq);
                $display("  I2C IRQ: %b", uut.i2c_irq);
                
                // Detailed error source analysis
                $display("Error Source Analysis:");
                $display("  Bus Memory Error: %b", uut.bus_mem_error);
                $display("  CPU Bus Error: %b", uut.cpu_bus_error);
                $display("  DMA Error: %b", uut.dma_error);
                $display("  Debug Error: %b", uut.debug_error);
                $display("  GPIO Error: %b", uut.gpio_error);
                $display("  UART Error: %b", uut.uart_error);
                $display("  SPI Error: %b (raw: %h)", |uut.spi_error, uut.spi_error);                $display("  I2C Error: %b (raw: %h)", |uut.i2c_error, uut.i2c_error);
                $display("  Flash Error: %b", uut.flash_error);
                  // I2C Controller State Debug
                $display("I2C Debug Information:");
                $display("  I2C State: %h", uut.i2c_inst.i2c_state);
                $display("  I2C Ready: %b", uut.i2c_inst.i2c_ready);
                $display("  I2C Bus Busy: %b", uut.i2c_inst.bus_busy);
                $display("  I2C Transfer Active: %b", uut.i2c_inst.transfer_active);
                $display("  I2C Error Count: %d", uut.i2c_inst.error_count);
                
                // Interface consolidation verification
                if (led_power == 1'b1 && led_error == 1'b0) begin
                    $display("PASS: Consolidated interfaces working correctly");
                end else begin
                    $display("FAIL: Interface consolidation issues detected");
                    errors = errors + 1;
                end

                // ----------------------------------------------------------
                // Phase 2: program execution through the whole SoC
                // ----------------------------------------------------------
                $display("=================================================");
                $display("Phase 2: hello_world on the SoC (CPU -> bus matrix");
                $display("         -> memory controller -> external SRAM)");
                $display("=================================================");

                wait (uut.cpu_inst.halted);
                #200;

                if (uut.cpu_inst.trap) begin
                    $display("FAIL: CPU trapped at PC 0x%0h", uut.cpu_inst.trap_pc);
                    errors = errors + 1;
                end

                got = "";
                for (wi = 0; wi < 12; wi = wi + 1) begin
                    wv = wordval(sram['h2000 + wi]);
                    got = {got[8*15-1:0], wv[7:0]};
                end
                $display("SoC says: \"%0s\"", got);
                for (wi = 0; wi < 12; wi = wi + 1) begin
                    wv = wordval(sram['h2000 + wi]);
                    if (wv !== hello[8*(12-wi)-1 -: 8]) begin
                        $display("FAIL: output[%0d] = %0d, want '%c'",
                                 wi, wv, hello[8*(12-wi)-1 -: 8]);
                        errors = errors + 1;
                    end
                end
                if (wordval(sram['h2100]) !== 12) begin
                    $display("FAIL: character count = %0d, want 12", wordval(sram['h2100]));
                    errors = errors + 1;
                end
                if (wordval(sram['h2104]) !== 1) begin
                    $display("FAIL: completion flag = %0d, want 1", wordval(sram['h2104]));
                    errors = errors + 1;
                end
                if (uut.cpu_inst.retired !== 32'd93) begin
                    $display("FAIL: retired = %0d, want 93 (vtxsim cross-check)",
                             uut.cpu_inst.retired);
                    errors = errors + 1;
                end
                $display("Cycles: %0d  Retired: %0d",
                         uut.cpu_inst.cycles, uut.cpu_inst.retired);

                $display("=================================================");
                if (errors == 0)
                    $display("PASS: VTX1 SoC executed hello_world end to end");
                else
                    $display("FAIL: %0d error(s)", errors);
                $display("=================================================");

                disable timeout_block;
                $finish;
            end
        join
    end
endmodule
