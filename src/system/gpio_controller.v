	`timescale 1ns / 1ps
// =============================================================================
// VTX1 Enhanced GPIO Controller
// =============================================================================
// Enhanced 24-pin GPIO controller with advanced features:
// - Individual pin interrupt generation with edge/level triggering
// - Configurable drive strength (2mA, 4mA, 8mA, 12mA)
// - Pull-up/pull-down configuration per pin
// - Pin wake-up capability for low-power modes
// - Pin multiplexing for alternate functions
// - Debounce filtering for input pins
// - High-speed digital I/O with configurable slew rate
// =============================================================================

`ifndef GPIO_CONTROLLER_V
`define GPIO_CONTROLLER_V

// Include VTX1 interface definitions
`include "vtx1_interfaces.v"

// Include VTX1 common infrastructure
`include "vtx1_interfaces.v"
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module gpio_controller (
    input  wire                     clk,
    input  wire                     rst_n,
    
    // =======================================================================
    // CPU INTERFACE - VTX1 STANDARDIZED
    // =======================================================================
    input  wire                     gpio_req,
    input  wire                     gpio_wr,
    input  wire [`VTX1_ADDR_WIDTH-1:0] gpio_addr,
    input  wire [`VTX1_WORD_WIDTH-1:0] gpio_wdata,
    output reg  [`VTX1_WORD_WIDTH-1:0] gpio_rdata,
    output reg                      gpio_ready,
    output reg                      gpio_error,
    
    // =======================================================================
    // ENHANCED 24-PIN GPIO INTERFACE
    // =======================================================================
    inout  wire [23:0]              gpio_pins,
    
    // =======================================================================
    // ALTERNATE FUNCTION INTERFACE
    // =======================================================================
    input  wire [23:0]              alt_func_out,      // Alternate function outputs
    output wire [23:0]              alt_func_in,       // Alternate function inputs
    input  wire [23:0]              alt_func_enable,   // Alternate function enable
    
    // =======================================================================
    // POWER MANAGEMENT INTERFACE
    // =======================================================================
    input  wire                     sleep_mode,        // System sleep mode
    output reg                      wake_request,      // Wake-up request from GPIO
    output reg  [23:0]              wake_source,       // Which pins caused wake-up
    
    // =======================================================================
    // INTERRUPT INTERFACE
    // =======================================================================
    output reg                      gpio_irq,
    output reg  [23:0]              gpio_irq_vector,   // Per-pin interrupt status
    
    // =======================================================================
    // DEBUG AND STATUS
    // =======================================================================
    output reg  [3:0]               gpio_state,
    output reg  [31:0]              operation_count,
    output reg  [7:0]               error_count,
    output reg  [23:0]              debounce_active,
    output reg  [23:0]              drive_strength_status
);

// =============================================================================
// ENHANCED GPIO CONTROLLER PARAMETERS
// =============================================================================

// State machine states - Use VTX1 standardized constants
localparam [3:0] GPIO_IDLE  = `VTX1_STATE_IDLE,
                 GPIO_READ  = `VTX1_GPIO_STATE_READ,
                 GPIO_WRITE = `VTX1_GPIO_STATE_WRITE,
                 GPIO_DEBOUNCE = `VTX1_STATE_WAIT,
                 GPIO_INTERRUPT = `VTX1_STATE_ACTIVE,
                 GPIO_ERROR = `VTX1_STATE_ERROR;

// Drive strength configurations
localparam [1:0] DRIVE_2MA  = 2'b00,
                 DRIVE_4MA  = 2'b01,
                 DRIVE_8MA  = 2'b10,
                 DRIVE_12MA = 2'b11;

// Pull-up/pull-down configurations
localparam [1:0] PULL_NONE = 2'b00,
                 PULL_DOWN = 2'b01,
                 PULL_UP   = 2'b10,
                 PULL_KEEP = 2'b11;

// Interrupt trigger types
localparam [1:0] IRQ_LEVEL_LOW  = 2'b00,
                 IRQ_LEVEL_HIGH = 2'b01,
                 IRQ_EDGE_FALL  = 2'b10,
                 IRQ_EDGE_RISE  = 2'b11;

// Enhanced register map for 24-pin GPIO with advanced features
localparam GPIO_CONTROL         = 16'h0000;  // Global control register
localparam GPIO_DIRECTION       = 16'h0004;  // Pin direction (0=input, 1=output)
localparam GPIO_OUTPUT_DATA     = 16'h0008;  // Output data register
localparam GPIO_INPUT_DATA      = 16'h000C;  // Input data register (read-only)
localparam GPIO_ALT_FUNC        = 16'h0010;  // Alternate function enable
localparam GPIO_DRIVE_STR_LOW   = 16'h0014;  // Drive strength [11:0] (2 bits per pin)
localparam GPIO_DRIVE_STR_HIGH  = 16'h0018;  // Drive strength [23:12] (2 bits per pin)
localparam GPIO_PULL_CFG_LOW    = 16'h001C;  // Pull-up/down config [11:0]
localparam GPIO_PULL_CFG_HIGH   = 16'h0020;  // Pull-up/down config [23:12]
localparam GPIO_IRQ_ENABLE      = 16'h0024;  // Interrupt enable per pin
localparam GPIO_IRQ_TYPE_LOW    = 16'h0028;  // IRQ trigger type [11:0]
localparam GPIO_IRQ_TYPE_HIGH   = 16'h002C;  // IRQ trigger type [23:12]
localparam GPIO_IRQ_STATUS      = 16'h0030;  // IRQ status (write 1 to clear)
localparam GPIO_DEBOUNCE_ENABLE = 16'h0034;  // Debounce enable per pin
localparam GPIO_DEBOUNCE_TIME   = 16'h0038;  // Debounce time configuration
localparam GPIO_WAKE_ENABLE     = 16'h003C;  // Wake-up enable per pin
localparam GPIO_SLEW_RATE       = 16'h0040;  // Slew rate control per pin
localparam GPIO_STATUS          = 16'h0044;  // Status register

// =============================================================================
// INTERNAL SIGNALS AND REGISTERS
// =============================================================================

// State machine
reg [3:0] state_reg, next_state;

// Configuration registers
reg [31:0] control_reg;
reg [23:0] direction_reg;
reg [23:0] output_data_reg;
reg [23:0] alt_func_reg;
reg [47:0] drive_strength_reg;     // 2 bits per pin * 24 pins
reg [47:0] pull_config_reg;        // 2 bits per pin * 24 pins
reg [23:0] irq_enable_reg;
reg [47:0] irq_type_reg;           // 2 bits per pin * 24 pins
reg [23:0] irq_status_reg;
reg [23:0] debounce_enable_reg;
reg [15:0] debounce_time_reg;
reg [23:0] wake_enable_reg;
reg [23:0] slew_rate_reg;

// Input processing
reg [23:0] input_sync1, input_sync2, input_sync3;
reg [23:0] input_filtered;
reg [23:0] prev_input;

// Debounce counters (simplified implementation)
reg [7:0] debounce_counter [0:23];
reg [23:0] debounce_stable;

// Interrupt processing
reg [23:0] irq_pending;
reg [23:0] edge_detected;

// Counters and status
reg [31:0] operation_count_reg;
reg [7:0] error_count_reg;

// Pin control signals
wire [23:0] pin_output_enable;
wire [23:0] pin_output_data;
wire [23:0] pin_input_data;

// =============================================================================
// INPUT SYNCHRONIZATION AND FILTERING
// =============================================================================

always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        input_sync1 <= 24'h000000;
        input_sync2 <= 24'h000000;
        input_sync3 <= 24'h000000;
        prev_input <= 24'h000000;
    end else begin
        // Three-stage synchronizer for metastability protection
        input_sync1 <= gpio_pins;
        input_sync2 <= input_sync1;
        input_sync3 <= input_sync2;
        prev_input <= input_filtered;
    end
end

// Debounce processing
integer j;
always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        for (j = 0; j < 24; j = j + 1) begin
            debounce_counter[j] <= 8'h00;
        end
        debounce_stable <= 24'h000000;
        input_filtered <= 24'h000000;
    end else begin
        for (j = 0; j < 24; j = j + 1) begin
            if (debounce_enable_reg[j]) begin
                if (input_sync3[j] == input_sync2[j]) begin
                    if (debounce_counter[j] < debounce_time_reg[7:0]) begin
                        debounce_counter[j] <= debounce_counter[j] + 1;
                    end else begin
                        debounce_stable[j] <= 1'b1;
                        input_filtered[j] <= input_sync3[j];
                    end
                end else begin
                    debounce_counter[j] <= 8'h00;
                    debounce_stable[j] <= 1'b0;
                end
            end else begin
                // No debouncing, pass through immediately
                input_filtered[j] <= input_sync3[j];
                debounce_stable[j] <= 1'b1;
            end
        end
    end
end

// =============================================================================
// INTERRUPT PROCESSING
// =============================================================================

// =============================================================================
// EDGE DETECTION - DRY COMPLIANT IMPLEMENTATION WITH GENERATE BLOCKS
// =============================================================================
// This replaces 600+ lines of repetitive code with a clean, maintainable solution

// Generate edge detection for all pins using parameterized logic
genvar i;
generate
    for (i = 0; i < 24; i = i + 1) begin : edge_detect_gen
        always @(posedge clk or negedge rst_n) begin
            if (!rst_n) begin
                edge_detected[i] <= 1'b0;
            end else begin
                case (irq_type_reg[(i*2)+1:(i*2)])
                    IRQ_EDGE_RISE: edge_detected[i] <= !prev_input[i] && input_filtered[i];
                    IRQ_EDGE_FALL: edge_detected[i] <= prev_input[i] && !input_filtered[i];
                    default:       edge_detected[i] <= 1'b0;  // Level triggers don't generate edges
                endcase
            end
        end
    end
endgenerate

// Generate interrupt pending logic for all pins
generate
    for (i = 0; i < 24; i = i + 1) begin : irq_pending_gen
        always @(posedge clk or negedge rst_n) begin
            if (!rst_n) begin
                irq_pending[i] <= 1'b0;
            end else if (irq_enable_reg[i]) begin
                case (irq_type_reg[(i*2)+1:(i*2)])
                    IRQ_LEVEL_LOW:  irq_pending[i] <= !input_filtered[i];
                    IRQ_LEVEL_HIGH: irq_pending[i] <= input_filtered[i];
                    IRQ_EDGE_FALL,
                    IRQ_EDGE_RISE:  irq_pending[i] <= irq_pending[i] || edge_detected[i];
                    default:        irq_pending[i] <= 1'b0;
                endcase
            end else begin
                irq_pending[i] <= 1'b0;
            end
        end
    end
endgenerate

// Interrupt clearing - separate from main interrupt logic
always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        // Reset handled in generate blocks above
    end else begin        // Clear interrupts when status register is written
        if (state_reg == GPIO_WRITE && gpio_addr[15:0] == GPIO_IRQ_STATUS) begin
            // Clear interrupts by writing 1 to clear
            for (j = 0; j < 24; j = j + 1) begin
                if (gpio_wdata[j]) begin
                    irq_pending[j] <= 1'b0;
                end
            end
        end
    end
end

// =============================================================================
// WAKE-UP PROCESSING
// =============================================================================

always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        wake_request <= 1'b0;
        wake_source <= 24'h000000;
    end else begin
        if (sleep_mode) begin
            // Check for wake-up events
            wake_source <= irq_pending & wake_enable_reg;
            wake_request <= |(irq_pending & wake_enable_reg);
        end else begin
            wake_request <= 1'b0;
            wake_source <= 24'h000000;
        end
    end
end

// =============================================================================
// STATE MACHINE
// =============================================================================

always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        state_reg <= GPIO_IDLE;
    end else begin
        state_reg <= next_state;
    end
end

always @(*) begin
    next_state = state_reg;
    case (state_reg)
        GPIO_IDLE: begin
            if (gpio_req) begin
                if (gpio_wr) begin
                    next_state = GPIO_WRITE;
                end else begin
                    next_state = GPIO_READ;
                end
            end else if (|(irq_pending & irq_enable_reg)) begin
                next_state = GPIO_INTERRUPT;
            end
        end
        
        GPIO_READ: begin
            next_state = GPIO_IDLE;
        end
        
        GPIO_WRITE: begin
            next_state = GPIO_IDLE;
        end
        
        GPIO_INTERRUPT: begin
            next_state = GPIO_IDLE;
        end
        
        default: begin
            next_state = GPIO_IDLE;
        end
    endcase
end

// =============================================================================
// REGISTER INTERFACE
// =============================================================================

always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        control_reg <= 32'h00000000;
        direction_reg <= 24'h000000;
        output_data_reg <= 24'h000000;
        alt_func_reg <= 24'h000000;
        drive_strength_reg <= 48'h000000000000;  // Default to 2mA
        pull_config_reg <= 48'h000000000000;     // Default to no pull
        irq_enable_reg <= 24'h000000;
        irq_type_reg <= 48'h000000000000;        // Default to level low
        debounce_enable_reg <= 24'h000000;
        debounce_time_reg <= 16'h0010;           // Default debounce time
        wake_enable_reg <= 24'h000000;
        slew_rate_reg <= 24'h000000;             // Default to normal slew rate
        operation_count_reg <= 32'h00000000;
        error_count_reg <= 8'h00;
    end else begin        if (state_reg == GPIO_WRITE) begin
            operation_count_reg <= operation_count_reg + 1;
            case (gpio_addr[15:0])
                // Consolidated register writes - DRY improvement
                GPIO_CONTROL:         control_reg <= gpio_wdata[31:0];
                GPIO_DIRECTION:       direction_reg <= gpio_wdata[23:0];
                GPIO_OUTPUT_DATA:     output_data_reg <= gpio_wdata[23:0];
                GPIO_ALT_FUNC:        alt_func_reg <= gpio_wdata[23:0];
                GPIO_DRIVE_STR_LOW:   drive_strength_reg[23:0] <= gpio_wdata[23:0];
                GPIO_DRIVE_STR_HIGH:  drive_strength_reg[47:24] <= gpio_wdata[23:0];
                GPIO_PULL_CFG_LOW:    pull_config_reg[23:0] <= gpio_wdata[23:0];
                GPIO_PULL_CFG_HIGH:   pull_config_reg[47:24] <= gpio_wdata[23:0];
                GPIO_IRQ_ENABLE:      irq_enable_reg <= gpio_wdata[23:0];
                GPIO_IRQ_TYPE_LOW:    irq_type_reg[23:0] <= gpio_wdata[23:0];
                GPIO_IRQ_TYPE_HIGH:   irq_type_reg[47:24] <= gpio_wdata[23:0];
                GPIO_IRQ_STATUS:      irq_status_reg <= gpio_wdata[23:0]; // Write handled in interrupt logic
                GPIO_DEBOUNCE_ENABLE: debounce_enable_reg <= gpio_wdata[23:0];
                GPIO_DEBOUNCE_TIME:   debounce_time_reg <= gpio_wdata[15:0];
                GPIO_WAKE_ENABLE:     wake_enable_reg <= gpio_wdata[23:0];
                GPIO_SLEW_RATE:       slew_rate_reg <= gpio_wdata[23:0];
                default: begin
                    error_count_reg <= error_count_reg + 1;
                end
            endcase
        end
    end
end

// =============================================================================
// OUTPUT CONTROL LOGIC
// =============================================================================

// CPU interface outputs
always @(*) begin
    // Default values
    gpio_ready = 1'b0;
    gpio_error = 1'b0;
    gpio_rdata = {`VTX1_WORD_WIDTH{1'b0}};
    
    case (state_reg)
        GPIO_IDLE: begin
            gpio_ready = !gpio_req;
        end
        
        GPIO_READ: begin
            gpio_ready = 1'b1;
            gpio_error = 1'b0;
              case (gpio_addr[15:0])
                // Consolidated register reads - DRY improvement
                GPIO_CONTROL:         gpio_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, control_reg};
                GPIO_DIRECTION:       gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, direction_reg};
                GPIO_OUTPUT_DATA:     gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, output_data_reg};
                GPIO_INPUT_DATA:      gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, input_filtered};
                GPIO_ALT_FUNC:        gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, alt_func_reg};
                GPIO_DRIVE_STR_LOW:   gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, drive_strength_reg[23:0]};
                GPIO_DRIVE_STR_HIGH:  gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, drive_strength_reg[47:24]};
                GPIO_PULL_CFG_LOW:    gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, pull_config_reg[23:0]};
                GPIO_PULL_CFG_HIGH:   gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, pull_config_reg[47:24]};
                GPIO_IRQ_ENABLE:      gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, irq_enable_reg};
                GPIO_IRQ_TYPE_LOW:    gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, irq_type_reg[23:0]};
                GPIO_IRQ_TYPE_HIGH:   gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, irq_type_reg[47:24]};
                GPIO_IRQ_STATUS:      gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, irq_pending};
                GPIO_DEBOUNCE_ENABLE: gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, debounce_enable_reg};
                GPIO_DEBOUNCE_TIME:   gpio_rdata = {{(`VTX1_WORD_WIDTH-16){1'b0}}, debounce_time_reg};
                GPIO_WAKE_ENABLE:     gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, wake_enable_reg};
                GPIO_SLEW_RATE:       gpio_rdata = {{(`VTX1_WORD_WIDTH-24){1'b0}}, slew_rate_reg};
                GPIO_STATUS:          gpio_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, 
                                                   {8'h00, debounce_stable, operation_count_reg[7:0]}};
                
                default: begin
                    gpio_error = 1'b1;
                    gpio_rdata = {`VTX1_WORD_WIDTH{1'b0}};
                end
            endcase
        end
        
        GPIO_WRITE: begin
            gpio_ready = 1'b1;
            gpio_error = 1'b0;
        end
        
        GPIO_INTERRUPT: begin
            gpio_ready = 1'b1;
            gpio_error = 1'b0;
        end
        
        default: begin
            gpio_ready = 1'b1;
            gpio_error = 1'b1;
        end
    endcase
end

// =============================================================================
// PIN MULTIPLEXING AND DRIVE CONTROL
// =============================================================================

// Pin output enable based on direction and alternate function
assign pin_output_enable = direction_reg & ~alt_func_reg;

// Pin output data multiplexing
assign pin_output_data = alt_func_reg ? alt_func_out : output_data_reg;

// Enhanced GPIO pin control with tri-state buffers
// Note: In real hardware, this would be implemented with configurable I/O cells
// For simulation, we use simplified tri-state logic

// Pin tri-state control (unrolled for compatibility)
assign gpio_pins[0] = pin_output_enable[0] ? pin_output_data[0] : 1'bz;
assign gpio_pins[1] = pin_output_enable[1] ? pin_output_data[1] : 1'bz;
assign gpio_pins[2] = pin_output_enable[2] ? pin_output_data[2] : 1'bz;
assign gpio_pins[3] = pin_output_enable[3] ? pin_output_data[3] : 1'bz;
assign gpio_pins[4] = pin_output_enable[4] ? pin_output_data[4] : 1'bz;
assign gpio_pins[5] = pin_output_enable[5] ? pin_output_data[5] : 1'bz;
assign gpio_pins[6] = pin_output_enable[6] ? pin_output_data[6] : 1'bz;
assign gpio_pins[7] = pin_output_enable[7] ? pin_output_data[7] : 1'bz;
assign gpio_pins[8] = pin_output_enable[8] ? pin_output_data[8] : 1'bz;
assign gpio_pins[9] = pin_output_enable[9] ? pin_output_data[9] : 1'bz;
assign gpio_pins[10] = pin_output_enable[10] ? pin_output_data[10] : 1'bz;
assign gpio_pins[11] = pin_output_enable[11] ? pin_output_data[11] : 1'bz;
assign gpio_pins[12] = pin_output_enable[12] ? pin_output_data[12] : 1'bz;
assign gpio_pins[13] = pin_output_enable[13] ? pin_output_data[13] : 1'bz;
assign gpio_pins[14] = pin_output_enable[14] ? pin_output_data[14] : 1'bz;
assign gpio_pins[15] = pin_output_enable[15] ? pin_output_data[15] : 1'bz;
assign gpio_pins[16] = pin_output_enable[16] ? pin_output_data[16] : 1'bz;
assign gpio_pins[17] = pin_output_enable[17] ? pin_output_data[17] : 1'bz;
assign gpio_pins[18] = pin_output_enable[18] ? pin_output_data[18] : 1'bz;
assign gpio_pins[19] = pin_output_enable[19] ? pin_output_data[19] : 1'bz;
assign gpio_pins[20] = pin_output_enable[20] ? pin_output_data[20] : 1'bz;
assign gpio_pins[21] = pin_output_enable[21] ? pin_output_data[21] : 1'bz;
assign gpio_pins[22] = pin_output_enable[22] ? pin_output_data[22] : 1'bz;
assign gpio_pins[23] = pin_output_enable[23] ? pin_output_data[23] : 1'bz;

// Alternate function input connection
assign alt_func_in = input_filtered;

// Pin input data
assign pin_input_data = input_filtered;

// =============================================================================
// STATUS AND INTERRUPT OUTPUTS
// =============================================================================

always @(*) begin
    // Interrupt outputs
    gpio_irq = |(irq_pending & irq_enable_reg);
    gpio_irq_vector = irq_pending;
    
    // Status outputs
    gpio_state = state_reg;
    operation_count = operation_count_reg;
    error_count = error_count_reg;
    debounce_active = debounce_enable_reg & ~debounce_stable;
    
    // Drive strength status (simplified - show which pins have non-default drive)
    // Unrolled for compatibility
    drive_strength_status[0] = (drive_strength_reg[1:0] != DRIVE_2MA);
    drive_strength_status[1] = (drive_strength_reg[3:2] != DRIVE_2MA);
    drive_strength_status[2] = (drive_strength_reg[5:4] != DRIVE_2MA);
    drive_strength_status[3] = (drive_strength_reg[7:6] != DRIVE_2MA);
    drive_strength_status[4] = (drive_strength_reg[9:8] != DRIVE_2MA);
    drive_strength_status[5] = (drive_strength_reg[11:10] != DRIVE_2MA);
    drive_strength_status[6] = (drive_strength_reg[13:12] != DRIVE_2MA);
    drive_strength_status[7] = (drive_strength_reg[15:14] != DRIVE_2MA);
    drive_strength_status[8] = (drive_strength_reg[17:16] != DRIVE_2MA);
    drive_strength_status[9] = (drive_strength_reg[19:18] != DRIVE_2MA);
    drive_strength_status[10] = (drive_strength_reg[21:20] != DRIVE_2MA);
    drive_strength_status[11] = (drive_strength_reg[23:22] != DRIVE_2MA);
    drive_strength_status[12] = (drive_strength_reg[25:24] != DRIVE_2MA);
    drive_strength_status[13] = (drive_strength_reg[27:26] != DRIVE_2MA);
    drive_strength_status[14] = (drive_strength_reg[29:28] != DRIVE_2MA);
    drive_strength_status[15] = (drive_strength_reg[31:30] != DRIVE_2MA);
    drive_strength_status[16] = (drive_strength_reg[33:32] != DRIVE_2MA);
    drive_strength_status[17] = (drive_strength_reg[35:34] != DRIVE_2MA);
    drive_strength_status[18] = (drive_strength_reg[37:36] != DRIVE_2MA);
    drive_strength_status[19] = (drive_strength_reg[39:38] != DRIVE_2MA);
    drive_strength_status[20] = (drive_strength_reg[41:40] != DRIVE_2MA);
    drive_strength_status[21] = (drive_strength_reg[43:42] != DRIVE_2MA);
    drive_strength_status[22] = (drive_strength_reg[45:44] != DRIVE_2MA);
    drive_strength_status[23] = (drive_strength_reg[47:46] != DRIVE_2MA);
end

endmodule

`endif // GPIO_CONTROLLER_V
