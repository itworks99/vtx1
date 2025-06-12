// VTX1 Simple Microcode Test
// Simplified test to debug microcode system issues

`timescale 1ns / 1ps

`include "ternary_constants.v"

module tb_microcode_simple;

    // ========================================================================
    // TESTBENCH SIGNALS
    // ========================================================================
    
    reg clk;
    reg rst_n;
    
    // ROM Interface Test
    reg  [9:0] rom_addr;
    wire [31:0] rom_data;
    reg  rom_enable;
    wire rom_ready, rom_error;
    
    // ========================================================================
    // DEVICE UNDER TEST - ROM ONLY
    // ========================================================================
    
    microcode_rom rom_inst (
        .clk(clk),
        .rst_n(rst_n),
        .addr(rom_addr),
        .data(rom_data),
        .enable(rom_enable),
        .ready(rom_ready),
        .error(rom_error),
        .last_addr(),
        .access_count()
    );
    
    // ========================================================================
    // TEST PROCEDURES
    // ========================================================================
    
    // Clock generation
    initial begin
        clk = 0;
        forever #5 clk = ~clk;
    end
    
    // Test ROM access
    task test_rom_read;
        input [9:0] address;
        begin
            $display("Testing ROM read at address 0x%03X", address);
            
            @(posedge clk);
            rom_addr = address;
            rom_enable = 1'b1;
            
            @(posedge clk);
            wait(rom_ready);
            
            if (rom_error) begin
                $display("  ERROR: ROM read failed");
            end else begin
                $display("  SUCCESS: ROM[0x%03X] = 0x%08X", address, rom_data);
            end
            
            rom_enable = 1'b0;
            @(posedge clk);
        end
    endtask
    
    // ========================================================================
    // MAIN TEST SEQUENCE
    // ========================================================================
    
    initial begin
        $display("=== VTX1 Simple Microcode ROM Test ===");
        
        // Initialize
        rst_n = 0;
        rom_enable = 0;
        rom_addr = 10'h0;
        
        // Reset
        #20 rst_n = 1;
        #10;
        
        $display("1. Testing ROM initialization");
        wait(rom_ready);
        $display("   ROM ready signal detected");
        
        $display("2. Testing ROM access patterns");
        
        // Test key entry points
        test_rom_read(10'h000);  // DIV entry
        test_rom_read(10'h018);  // MOD entry
        test_rom_read(10'h100);  // SIN entry
        test_rom_read(10'h3F0);  // Error handler
        
        $display("3. Testing ROM boundary conditions");
        test_rom_read(10'h3FF);  // Last valid address
        test_rom_read(10'h3FF);  // Should work
        
        $display("=== Simple Microcode ROM Test Complete ===");
        $finish;
    end
    
    // Timeout protection
    initial begin
        #10000;
        $display("ERROR: Test timeout!");
        $finish;
    end

endmodule
