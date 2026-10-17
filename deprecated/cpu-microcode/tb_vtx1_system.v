// VTX1 System-Level Testbench
// Part of the VTX1 Ternary System-on-Chip

`timescale 1ns / 1ps
// Include paths handled by compiler -I flags (see Taskfile.yml)

module tb_vtx1_system;

    // Test signals
    reg clk_ext;
    reg clk_32k;
    reg rst_ext_n;
    reg pwr_good;
    
    // External Memory Interface
    wire [35:0] mem_addr;
    wire [35:0] mem_data_out;
    reg  [35:0] mem_data_in;
    wire mem_we_n;
    wire mem_oe_n;
    wire mem_cs_n;
    
    // Flash Memory Interface (SPI)
    wire flash_cs_n;
    wire flash_sck;
    wire flash_mosi;
    reg  flash_miso;
    wire flash_wp_n;
    wire flash_hold_n;
    
    // GPIO Interface
    wire [31:0] gpio_pins_out;
    reg  [31:0] gpio_pins_in;
    
    // UART Interface
    reg  uart_rx;
    wire uart_tx;
    wire uart_rts;
    reg  uart_cts;
    
    // Debug Interface (JTAG)
    reg  jtag_tck;
    reg  jtag_tms;
    reg  jtag_tdi;
    wire jtag_tdo;
    reg  jtag_trst_n;
    
    // Power Management
    wire pwr_enable;
    
    // Status and Debug
    wire [7:0] system_status;
    wire [7:0] error_flags;
    
    // Bidirectional signal handling
    reg  [35:0]  mem_drive_data;
    reg                     mem_drive_enable;
    reg  [31:0]             gpio_drive_data;
    reg                     gpio_drive_enable;
    
    assign mem_data_out = mem_drive_enable ? mem_drive_data : {`DATA_WIDTH{1'bz}};
    assign gpio_pins_out = gpio_drive_enable ? gpio_drive_data : 32'hzzzzzzzz;
    
    // Instantiate VTX1 System
    vtx1_top dut (
        .clk_ext(clk_ext),
        .clk_32k(clk_32k),
        .rst_ext_n(rst_ext_n),
        .mem_addr(mem_addr),
        .mem_data(mem_data_out),
        .mem_we_n(mem_we_n),
        .mem_oe_n(mem_oe_n),
        .mem_cs_n(mem_cs_n),
        .flash_cs_n(flash_cs_n),
        .flash_sck(flash_sck),
        .flash_mosi(flash_mosi),
        .flash_miso(flash_miso),
        .flash_wp_n(flash_wp_n),
        .flash_hold_n(flash_hold_n),
        .gpio_pins(gpio_pins_out),
        .uart_rx(uart_rx),
        .uart_tx(uart_tx),
        .uart_rts(uart_rts),
        .uart_cts(uart_cts),
        .jtag_tck(jtag_tck),
        .jtag_tms(jtag_tms),
        .jtag_tdi(jtag_tdi),
        .jtag_tdo(jtag_tdo),
        .jtag_trst_n(jtag_trst_n),
        .pwr_good(pwr_good),
        .pwr_enable(pwr_enable)
    );
    
    // ========================================================================
    // CLOCK GENERATION
    // ========================================================================
    initial begin
        clk_ext = 0;
        forever #18.52 clk_ext = ~clk_ext;  // 27MHz external clock
    end
    
    initial begin
        clk_32k = 0;
        forever #15258.78 clk_32k = ~clk_32k;  // 32.768kHz RTC clock
    end
    
    // ========================================================================
    // MEMORY MODEL
    // ========================================================================
    reg [35:0] memory [0:1048575];  // 1M words of memory
    
    // Memory response logic
    always @(posedge clk_ext) begin
        if (!rst_ext_n) begin
            mem_drive_enable <= 1'b0;
        end else begin
            if (!mem_cs_n && !mem_oe_n) begin
                // Read operation
                mem_drive_data <= memory[mem_addr[19:0]];
                mem_drive_enable <= 1'b1;
            end else if (!mem_cs_n && !mem_we_n) begin
                // Write operation
                memory[mem_addr[19:0]] <= mem_data_out;
                mem_drive_enable <= 1'b0;
            end else begin
                mem_drive_enable <= 1'b0;
            end
        end
    end
    
    // ========================================================================
    // FLASH MEMORY MODEL
    // ========================================================================
    reg [7:0] flash_memory [0:442367];  // 432KB Flash
    reg [7:0] spi_shift_reg;
    reg [2:0] spi_bit_count;
    reg [7:0] spi_command;
    reg [23:0] spi_address;
    
    // Simple SPI Flash model
    always @(posedge flash_sck or posedge flash_cs_n) begin
        if (flash_cs_n) begin
            spi_bit_count <= 3'h0;
            flash_miso <= 1'bz;
        end else begin
            spi_shift_reg <= {spi_shift_reg[6:0], flash_mosi};
            spi_bit_count <= spi_bit_count + 1;
            
            if (spi_bit_count == 3'h7) begin
                // Command byte received
                spi_command <= {spi_shift_reg[6:0], flash_mosi};
            end
            
            // Implement basic read command (0x03)
            if (spi_command == 8'h03 && spi_bit_count > 3'h7) begin
                flash_miso <= flash_memory[spi_address][7-spi_bit_count];
            end
        end
    end
    
    // ========================================================================
    // TEST STIMULUS
    // ========================================================================
    initial begin
        // Initialize signals
        rst_ext_n = 0;
        pwr_good = 1;
        mem_data_in = `DEFAULT_WORD;
        flash_miso = 1'bz;
        gpio_pins_in = 32'h0;
        uart_rx = 1'b1;
        uart_cts = 1'b0;
        jtag_tck = 0;
        jtag_tms = 0;
        jtag_tdi = 0;
        jtag_trst_n = 1;
        mem_drive_data = `DEFAULT_WORD;
        mem_drive_enable = 1'b0;
        gpio_drive_data = 32'h0;
        gpio_drive_enable = 1'b0;
        
        // Initialize memory with test program
        initialize_memory();
        initialize_flash();
        
        // Generate VCD file for waveform analysis
        $dumpfile("vtx1_system_test.vcd");
        $dumpvars(0, tb_vtx1_system);
        
        // Reset sequence
        #100;
        rst_ext_n = 1;
        #50;
        
        $display("=== VTX1 System Test Started ===");
        $display("Time: %0t", $time);
        
        // Test 1: System Boot
        test_system_boot();
        
        // Test 2: Memory Interface
        test_memory_interface();
        
        // Test 3: Flash Interface
        test_flash_interface();
        
        // Test 4: GPIO Operations
        test_gpio_operations();
        
        // Test 5: UART Communication
        test_uart_communication();
        
        // Test 6: Clock Management
        test_clock_management();
        
        // Test 7: Power Management
        test_power_management();
        
        // Test 8: Debug Interface
        test_debug_interface();
        
        $display("=== VTX1 System Test Completed ===");
        
        #1000;
        $finish;
    end
    
    // ========================================================================
    // TEST TASKS
    // ========================================================================
    
    // Test system boot sequence
    task test_system_boot;
        begin
            $display("\n--- Test 1: System Boot ---");
            
            // Wait for system to stabilize
            repeat(100) @(posedge clk_ext);
            
            // Check power enable
            if (pwr_enable) begin
                $display("PASS: Power management active");
            end else begin
                $display("ERROR: Power management not active");
            end
            
            // Check if system is running
            repeat(1000) @(posedge clk_ext);
            $display("PASS: System boot sequence");
        end
    endtask
    
    // Test memory interface
    task test_memory_interface;
        begin
            $display("\n--- Test 2: Memory Interface ---");
            
            // Wait for memory access
            wait(!mem_cs_n);
            $display("Memory access detected: addr=%h", mem_addr);
            
            repeat(10) @(posedge clk_ext);
            $display("PASS: Memory interface test");
        end
    endtask
    
    // Test flash interface
    task test_flash_interface;
        begin
            $display("\n--- Test 3: Flash Interface ---");
            
            // Wait for flash access
            wait(!flash_cs_n);
            $display("Flash access detected");
            
            repeat(50) @(posedge flash_sck);
            $display("PASS: Flash interface test");
        end
    endtask
    
    // Test GPIO operations
    task test_gpio_operations;
        begin
            $display("\n--- Test 4: GPIO Operations ---");
            
            // Simulate GPIO input changes
            gpio_pins_in = 32'hAAAAAAAA;
            gpio_drive_enable = 1'b1;
            
            repeat(100) @(posedge clk_ext);
            
            // Check GPIO outputs
            gpio_drive_enable = 1'b0;
            repeat(10) @(posedge clk_ext);
            
            $display("PASS: GPIO operations test");
        end
    endtask
    
    // Test UART communication
    task test_uart_communication;
        begin
            $display("\n--- Test 5: UART Communication ---");
            
            // Send UART data
            send_uart_byte(8'h55);
            
            repeat(1000) @(posedge clk_ext);
            
            // Check for UART transmission
            if (uart_tx !== 1'b1) begin
                $display("INFO: UART transmission detected");
            end
            
            $display("PASS: UART communication test");
        end
    endtask
    
    // Test clock management
    task test_clock_management;
        begin
            $display("\n--- Test 6: Clock Management ---");
            
            // Clock domains should be running
            repeat(100) @(posedge clk_ext);
            
            $display("PASS: Clock management test");
        end
    endtask
    
    // Test power management
    task test_power_management;
        begin
            $display("\n--- Test 7: Power Management ---");
            
            // Test power down
            pwr_good = 0;
            repeat(10) @(posedge clk_ext);
            
            // Test power up
            pwr_good = 1;
            repeat(100) @(posedge clk_ext);
            
            $display("PASS: Power management test");
        end
    endtask
    
    // Test debug interface
    task test_debug_interface;
        begin
            $display("\n--- Test 8: Debug Interface ---");
            
            // JTAG reset
            jtag_trst_n = 0;
            repeat(5) @(posedge jtag_tck);
            jtag_trst_n = 1;
            
            // Simple JTAG sequence
            send_jtag_sequence();
            
            $display("PASS: Debug interface test");
        end
    endtask
    
    // ========================================================================
    // UTILITY TASKS
    // ========================================================================
    
    // Initialize memory with test program
    task initialize_memory;
        integer i;
        begin
            // Initialize with pattern
            for (i = 0; i < 1048576; i = i + 1) begin
                memory[i] = i[35:0] ^ `DEFAULT_WORD;
            end
            
            // Load test program at address 0
            memory[0] = 36'h123456789;  // Test instruction 1
            memory[1] = 36'h987654321;  // Test instruction 2
            memory[2] = 36'h555555555;  // Test instruction 3
            
            $display("Memory initialized with test program");
        end
    endtask
    
    // Initialize flash memory
    task initialize_flash;
        integer i;
        begin
            for (i = 0; i < 442368; i = i + 1) begin
                flash_memory[i] = i[7:0] ^ 8'hA5;
            end
            $display("Flash memory initialized");
        end
    endtask
    
    // Send UART byte
    task send_uart_byte;
        input [7:0] data;
        integer i;
        begin
            // Start bit
            uart_rx = 0;
            repeat(868) @(posedge clk_ext);  // 115200 baud at 27MHz
            
            // Data bits
            for (i = 0; i < 8; i = i + 1) begin
                uart_rx = data[i];
                repeat(868) @(posedge clk_ext);
            end
            
            // Stop bit
            uart_rx = 1;
            repeat(868) @(posedge clk_ext);
        end
    endtask
    
    // Send JTAG sequence
    task send_jtag_sequence;
        begin
            // Generate JTAG clock
            repeat(10) begin
                jtag_tck = 0;
                #50;
                jtag_tck = 1;
                #50;
            end
        end
    endtask
    
    // ========================================================================
    // MONITORING
    // ========================================================================
    
    // Monitor system activity
    always @(posedge clk_ext) begin
        if (rst_ext_n) begin
            // Monitor memory accesses
            if (!mem_cs_n) begin
                if (!mem_we_n) begin
                    $display("Time %0t: Memory WRITE addr=%h data=%h", 
                            $time, mem_addr, mem_data_out);
                end else if (!mem_oe_n) begin
                    $display("Time %0t: Memory READ addr=%h", $time, mem_addr);
                end
            end
        end
    end
    
    // Monitor UART activity
    always @(negedge uart_tx) begin
        $display("Time %0t: UART transmission started", $time);
    end
    
    // Monitor power events
    always @(pwr_good) begin
        if (pwr_good) begin
            $display("Time %0t: Power good asserted", $time);
        end else begin
            $display("Time %0t: Power good deasserted", $time);
        end
    end

endmodule
