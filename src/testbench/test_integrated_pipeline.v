// ============================================================================
// VTX1 Integrated Pipeline Testbench
// ============================================================================
// Basic functionality test for the enhanced CPU pipeline with:
// - Hazard detection unit
// - Forwarding unit  
// - TCU integration
// ============================================================================

`timescale 1ns/1ps

`include "vtx1_interfaces.v"

module test_integrated_pipeline;

    // Clock and reset
    reg clk;
    reg rst_n;
    reg enable;
    
    // Instruction memory interface
    wire [`VTX1_WORD_WIDTH-1:0] imem_addr;
    reg [`VTX1_VLIW_WIDTH-1:0] imem_data;
    wire imem_req;
    reg imem_ready;
    
    // Data memory interface
    wire [`VTX1_WORD_WIDTH-1:0] dmem_addr;
    wire [`VTX1_WORD_WIDTH-1:0] dmem_data;
    wire dmem_we, dmem_oe, dmem_req;
    reg dmem_ready;
    
    // Debug and status
    wire [`VTX1_WORD_WIDTH-1:0] debug_pc;
    wire [3:0] debug_status;
    wire [31:0] cycle_count, instruction_count, stall_count;
    wire pipeline_stall, pipeline_flush;
    wire [1:0] pipeline_state;
    
    // Error handling
    wire error, timeout;
    wire [3:0] error_code;
    wire [31:0] error_count;
    
    // Test control
    integer test_cycle;
    reg [31:0] test_instruction;
    
    // ========================================================================
    // DUT INSTANTIATION
    // ========================================================================
    
    cpu_core dut (
        .clk(clk),
        .rst_n(rst_n),
        .enable(enable),
        
        // Instruction memory
        .imem_addr(imem_addr),
        .imem_data(imem_data),
        .imem_req(imem_req),
        .imem_ready(imem_ready),
        
        // Data memory
        .dmem_addr(dmem_addr),
        .dmem_data(dmem_data),
        .dmem_we(dmem_we),
        .dmem_oe(dmem_oe),
        .dmem_req(dmem_req),
        .dmem_ready(dmem_ready),
        
        // Interrupt interface (unused for basic test)
        .interrupt_req(16'h0),
        .interrupt_ack(),
        .nmi_req(1'b0),
        
        // Debug interface
        .debug_enable(1'b0),
        .debug_step(1'b0),
        .debug_cmd(4'h0),
        .debug_pc(debug_pc),
        .debug_status(debug_status),
        
        // Performance counters
        .cycle_count(cycle_count),
        .instruction_count(instruction_count),
        .stall_count(stall_count),
        
        // Pipeline status
        .pipeline_stall(pipeline_stall),
        .pipeline_flush(pipeline_flush),
        .pipeline_state(pipeline_state),
        
        // Error handling
        .error(error),
        .error_code(error_code),
        .timeout(timeout),
        .error_count(error_count)
    );
    
    // ========================================================================
    // CLOCK GENERATION
    // ========================================================================
    
    initial begin
        clk = 0;
        forever #5 clk = ~clk; // 100MHz clock
    end
    
    // ========================================================================
    // STIMULUS AND TEST SEQUENCE
    // ========================================================================
    
    initial begin
        $dumpfile("test_integrated_pipeline.vcd");
        $dumpvars(0, test_integrated_pipeline);
        
        // Initialize signals
        rst_n = 0;
        enable = 0;
        imem_ready = 1;
        dmem_ready = 1;
        imem_data = {`VTX1_VLIW_WIDTH{1'b0}};
        test_cycle = 0;
        
        $display("========================================");
        $display("VTX1 Integrated Pipeline Test");
        $display("========================================");
        
        // Reset sequence
        #20;
        rst_n = 1;
        #10;
        enable = 1;
        
        // Test 1: Basic pipeline operation
        $display("\nTest 1: Basic Pipeline Operation");
        test_basic_pipeline();
        
        // Test 2: Hazard detection
        $display("\nTest 2: Hazard Detection");
        test_hazard_detection();
        
        // Test 3: Forwarding operation
        $display("\nTest 3: Data Forwarding");
        test_data_forwarding();
        
        // Test 4: Pipeline stall and flush
        $display("\nTest 4: Pipeline Control");
        test_pipeline_control();
        
        // Test completion
        #100;
        $display("\n========================================");
        $display("Test Completed");
        $display("Cycles: %d, Instructions: %d, Stalls: %d", 
                 cycle_count, instruction_count, stall_count);
        $display("Error Count: %d", error_count);
        $display("========================================");
        
        $finish;
    end
    
    // ========================================================================
    // TEST PROCEDURES
    // ========================================================================
    
    task test_basic_pipeline;
        begin
            $display("  Running basic pipeline test...");
            
            // Provide simple instructions
            repeat (10) begin
                @(posedge clk);
                if (imem_req) begin
                    // Simple NOP-like instruction
                    imem_data <= generate_test_instruction(3'b000, 4'h0, 4'h1, 4'h2);
                end
                test_cycle <= test_cycle + 1;
            end
            
            $display("  Basic pipeline test completed");
            $display("    Cycles: %d", test_cycle);
            $display("    Pipeline State: %d", pipeline_state);
            $display("    Stalls: %d", stall_count);
        end
    endtask
    
    task test_hazard_detection;
        begin
            $display("  Running hazard detection test...");
            
            // Create RAW hazard scenario
            @(posedge clk);
            if (imem_req) begin
                // Instruction that writes to register 1
                imem_data <= generate_test_instruction(3'b001, 4'h0, 4'h2, 4'h1);
            end
            
            @(posedge clk);
            if (imem_req) begin
                // Instruction that reads from register 1 (RAW hazard)
                imem_data <= generate_test_instruction(3'b010, 4'h1, 4'h2, 4'h3);
            end
            
            // Monitor for stalls
            repeat (5) begin
                @(posedge clk);
                if (pipeline_stall) begin
                    $display("    Hazard detected, pipeline stalled");
                end
            end
            
            $display("  Hazard detection test completed");
        end
    endtask
    
    task test_data_forwarding;
        begin
            $display("  Running data forwarding test...");
            
            // Test forwarding scenarios
            repeat (8) begin
                @(posedge clk);
                if (imem_req) begin
                    // Alternating read/write pattern
                    if (test_cycle % 2 == 0) begin
                        imem_data <= generate_test_instruction(3'b001, 4'h0, 4'h1, 4'h2);
                    end else begin
                        imem_data <= generate_test_instruction(3'b010, 4'h2, 4'h0, 4'h3);
                    end
                end
                test_cycle <= test_cycle + 1;
            end
            
            $display("  Data forwarding test completed");
        end
    endtask
    
    task test_pipeline_control;
        begin
            $display("  Running pipeline control test...");
            
            // Test memory stall scenarios
            dmem_ready <= 0; // Simulate memory not ready
            
            repeat (3) begin
                @(posedge clk);
                if (pipeline_stall) begin
                    $display("    Memory stall detected");
                end
            end
            
            dmem_ready <= 1; // Resume normal operation
            
            repeat (5) begin
                @(posedge clk);
                test_cycle <= test_cycle + 1;
            end
            
            $display("  Pipeline control test completed");
        end
    endtask
    
    // ========================================================================
    // HELPER FUNCTIONS
    // ========================================================================
    
    function [`VTX1_VLIW_WIDTH-1:0] generate_test_instruction;
        input [2:0] opcode;
        input [3:0] src_reg;
        input [3:0] dst_reg;
        input [3:0] aux_reg;
        begin
            // Simple instruction format for testing
            // This is a simplified format - real VLIW would be more complex
            generate_test_instruction = {
                {(`VTX1_VLIW_WIDTH-16){1'b0}}, // Padding
                opcode,                        // Operation type
                src_reg,                       // Source register
                dst_reg,                       // Destination register  
                aux_reg,                       // Auxiliary register
                1'b1                          // Valid bit
            };
        end
    endfunction
    
    // ========================================================================
    // MONITORING
    // ========================================================================
    
    always @(posedge clk) begin
        if (error) begin
            $display("ERROR at cycle %d: Code %d", test_cycle, error_code);
        end
        
        if (timeout) begin
            $display("TIMEOUT at cycle %d", test_cycle);
        end
    end

endmodule
