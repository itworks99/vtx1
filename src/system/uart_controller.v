	`timescale 1ns / 1ps
// =============================================================================
// VTX1 Enhanced UART Controller
// =============================================================================
// Enhanced UART controller with advanced features:
// - FIFO buffers for transmit and receive
// - Hardware flow control (RTS/CTS)
// - Break detection and generation
// - DMA support for efficient data transfer
// - Configurable baud rate and frame format
// - Comprehensive error detection
// =============================================================================

`ifndef UART_CONTROLLER_V
`define UART_CONTROLLER_V

// Include VTX1 interface definitions
`include "vtx1_interfaces.v"
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module uart_controller (
    input  wire                     clk,
    input  wire                     rst_n,
    
    // =======================================================================
    // CPU INTERFACE - VTX1 STANDARDIZED
    // =======================================================================
    input  wire                     uart_req,
    input  wire                     uart_wr,
    input  wire [`VTX1_ADDR_WIDTH-1:0] uart_addr,
    input  wire [`VTX1_WORD_WIDTH-1:0] uart_wdata,
    output reg  [`VTX1_WORD_WIDTH-1:0] uart_rdata,
    output reg                      uart_ready,
    output reg                      uart_error,
    
    // =======================================================================
    // UART PHYSICAL INTERFACE
    // =======================================================================
    input  wire                     uart_rx,
    output reg                      uart_tx,
    output reg                      uart_rts,
    input  wire                     uart_cts,
    
    // =======================================================================
    // INTERRUPT INTERFACE
    // =======================================================================
    output reg                      uart_irq,
    output reg  [7:0]               uart_irq_status,
    
    // =======================================================================
    // DMA INTERFACE
    // =======================================================================
    output reg                      uart_dma_tx_req,
    output reg                      uart_dma_rx_req,
    input  wire                     uart_dma_tx_ack,
    input  wire                     uart_dma_rx_ack,
    output reg  [7:0]               uart_dma_tx_data,
    input  wire [7:0]               uart_dma_rx_data,
    
    // =======================================================================
    // CONFIGURATION
    // =======================================================================
    input  wire                     uart_enable,
    input  wire [15:0]              uart_clk_div,
    
    // =======================================================================
    // DEBUG AND STATUS
    // =======================================================================
    output reg  [3:0]               uart_state,
    output reg  [31:0]              tx_count,
    output reg  [31:0]              rx_count,
    output reg  [7:0]               error_count,
    output reg  [15:0]              fifo_status
);

// =============================================================================
// UART CONTROLLER PARAMETERS
// =============================================================================

// Register offsets (16-bit addressing within UART space)
localparam UART_DATA        = 16'h0000;  // Data register (R/W)
localparam UART_STATUS      = 16'h0004;  // Status register (RO)
localparam UART_CONTROL     = 16'h0008;  // Control register (RW)
localparam UART_BAUD_DIV    = 16'h000C;  // Baud rate divisor (RW)
localparam UART_FRAME_CFG   = 16'h0010;  // Frame configuration (RW)
localparam UART_FIFO_CTRL   = 16'h0014;  // FIFO control (RW)
localparam UART_INT_CTRL    = 16'h0018;  // Interrupt control (RW)
localparam UART_INT_STATUS  = 16'h001C;  // Interrupt status (RW1C)
localparam UART_DMA_CTRL    = 16'h0020;  // DMA control (RW)
localparam UART_LINE_STATUS = 16'h0024;  // Line status (RO)
localparam UART_MODEM_CTRL  = 16'h0028;  // Modem control (RW)
localparam UART_MODEM_STATUS= 16'h002C;  // Modem status (RO)
localparam UART_BREAK_CTRL  = 16'h0030;  // Break control (RW)
localparam UART_TX_FIFO_LVL = 16'h0034;  // TX FIFO level (RO)
localparam UART_RX_FIFO_LVL = 16'h0038;  // RX FIFO level (RO)
localparam UART_ERROR_COUNT = 16'h003C;  // Error count (RO)

// Control register bits
localparam CTRL_ENABLE      = 0;         // UART enable
localparam CTRL_TX_ENABLE   = 1;         // Transmitter enable
localparam CTRL_RX_ENABLE   = 2;         // Receiver enable
localparam CTRL_LOOPBACK    = 3;         // Loopback mode
localparam CTRL_FLOW_CTRL   = 4;         // Hardware flow control
localparam CTRL_BREAK_EN    = 5;         // Break generation

// Status register bits
localparam STAT_TX_EMPTY    = 0;         // TX FIFO empty
localparam STAT_TX_FULL     = 1;         // TX FIFO full
localparam STAT_RX_EMPTY    = 2;         // RX FIFO empty
localparam STAT_RX_FULL     = 3;         // RX FIFO full
localparam STAT_TX_ACTIVE   = 4;         // Transmitting
localparam STAT_RX_ACTIVE   = 5;         // Receiving
localparam STAT_PARITY_ERR  = 6;         // Parity error
localparam STAT_FRAME_ERR   = 7;         // Frame error
localparam STAT_OVERRUN_ERR = 8;         // Overrun error
localparam STAT_BREAK_DET   = 9;         // Break detected

// Interrupt bits
localparam INT_TX_EMPTY     = 0;         // TX FIFO empty interrupt
localparam INT_RX_READY     = 1;         // RX data ready interrupt
localparam INT_ERROR        = 2;         // Error interrupt
localparam INT_BREAK        = 3;         // Break interrupt
localparam INT_MODEM        = 4;         // Modem status change

// Frame configuration
localparam FRAME_BITS_5     = 2'b00;     // 5 data bits
localparam FRAME_BITS_6     = 2'b01;     // 6 data bits
localparam FRAME_BITS_7     = 2'b10;     // 7 data bits
localparam FRAME_BITS_8     = 2'b11;     // 8 data bits

localparam PARITY_NONE      = 2'b00;     // No parity
localparam PARITY_ODD       = 2'b01;     // Odd parity
localparam PARITY_EVEN      = 2'b10;     // Even parity
localparam PARITY_MARK      = 2'b11;     // Mark parity

// Controller states - Use VTX1 standardized constants
localparam UART_IDLE        = `VTX1_STATE_IDLE;
localparam UART_READ        = `VTX1_STATE_READ;
localparam UART_WRITE       = `VTX1_STATE_WRITE;
localparam UART_TX_START    = `VTX1_UART_STATE_TX_START;
localparam UART_TX_DATA     = `VTX1_UART_STATE_TX_DATA;
localparam UART_TX_STOP     = `VTX1_UART_STATE_TX_STOP;
localparam UART_RX_START    = `VTX1_UART_STATE_RX_START;
localparam UART_RX_DATA     = `VTX1_UART_STATE_RX_DATA;
localparam UART_RX_STOP     = `VTX1_UART_STATE_RX_STOP;
localparam UART_BREAK       = `VTX1_STATE_SPECIAL;
localparam UART_ERROR       = `VTX1_STATE_ERROR;

// FIFO parameters
localparam FIFO_DEPTH       = 64;        // 64-byte FIFOs
localparam FIFO_ADDR_WIDTH  = 6;         // log2(64)

// =============================================================================
// INTERNAL REGISTERS
// =============================================================================

// State machine
reg [3:0] uart_state_reg, uart_state_next;

// VTX1 standardized error tracking
reg [`VTX1_ERROR_CODE_WIDTH-1:0] uart_error_code;
reg uart_timeout_detected;
reg [`VTX1_ERROR_CODE_WIDTH-1:0] vtx1_error_reg;

// Configuration registers
reg [31:0] control_reg;
reg [31:0] baud_div_reg;
reg [31:0] frame_cfg_reg;
reg [31:0] fifo_ctrl_reg;
reg [31:0] int_ctrl_reg;
reg [31:0] int_status_reg;
reg [31:0] dma_ctrl_reg;
reg [31:0] modem_ctrl_reg;
reg [31:0] break_ctrl_reg;

// Status registers
reg [31:0] status_reg;
reg [31:0] line_status_reg;
reg [31:0] modem_status_reg;

// FIFO implementation
reg [7:0] tx_fifo[0:FIFO_DEPTH-1];
reg [7:0] rx_fifo[0:FIFO_DEPTH-1];
reg [FIFO_ADDR_WIDTH-1:0] tx_wr_ptr, tx_rd_ptr;
reg [FIFO_ADDR_WIDTH-1:0] rx_wr_ptr, rx_rd_ptr_reg;
wire [FIFO_ADDR_WIDTH-1:0] rx_rd_ptr = rx_rd_ptr_reg; // Wire alias for compatibility
reg [FIFO_ADDR_WIDTH:0] tx_fifo_count, rx_count_reg;

// Transmitter state
reg [3:0] tx_state;
reg [3:0] tx_bit_count;
reg [15:0] tx_baud_count;
reg [7:0] tx_shift_reg;
reg tx_parity_bit;
reg tx_active;

// Receiver state
reg [3:0] rx_state;
reg [3:0] rx_bit_count;
reg [15:0] rx_baud_count;
reg [7:0] rx_shift_reg;
reg rx_parity_bit;
reg rx_active;
reg expected_parity_bit;  // For parity calculation

// Input synchronizers
reg uart_rx_sync1, uart_rx_sync2, uart_rx_sync3;
reg uart_cts_sync1, uart_cts_sync2;

// Error detection
reg [7:0] error_count_reg;
reg [31:0] tx_count_reg, rx_count_total;

// Break detection
reg [15:0] break_counter;
reg break_detected;

// Timing counters
reg [4:0] uart_timeout_counter;

// =============================================================================
// STATE MACHINE
// =============================================================================

// State register
always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        uart_state_reg <= UART_IDLE;
    end else begin
        uart_state_reg <= uart_state_next;
    end
end

// Next state logic
always @(*) begin
    uart_state_next = uart_state_reg;
    
    case (uart_state_reg)
        UART_IDLE: begin
            if (uart_req) begin
                if (uart_wr) begin
                    uart_state_next = UART_WRITE;
                end else begin
                    uart_state_next = UART_READ;
                end
            end
        end        
        UART_READ: begin
            uart_state_next = UART_IDLE;
        end
        
        UART_WRITE: begin
            uart_state_next = UART_IDLE;
        end
        
        UART_ERROR: begin
            if (uart_timeout_counter == 0) begin
                uart_state_next = UART_IDLE;
            end
        end
        
        default: begin
            uart_state_next = UART_IDLE;
        end
    endcase
end

// =============================================================================
// REGISTER ACCESS LOGIC
// =============================================================================

always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        // Reset all configuration registers
        control_reg <= 32'h00000000;
        baud_div_reg <= 32'h00000364;     // Default 9600 baud (assuming 50MHz clock)
        frame_cfg_reg <= 32'h00000003;   // 8 data bits, no parity, 1 stop bit
        fifo_ctrl_reg <= 32'h00000000;
        int_ctrl_reg <= 32'h00000000;
        int_status_reg <= 32'h00000000;
        dma_ctrl_reg <= 32'h00000000;
        modem_ctrl_reg <= 32'h00000000;
        break_ctrl_reg <= 32'h00000000;
          error_count_reg <= 8'h00;
        uart_timeout_counter <= 5'h00;
          // VTX1 standardized error handling
        uart_error_code <= `VTX1_ERROR_NONE;
        uart_timeout_detected <= 1'b0;
        `VTX1_CLEAR_ERROR(vtx1_error_reg, uart_state_reg)
    end else begin
        case (uart_state_reg)            UART_WRITE: begin
                if (uart_req && uart_wr) begin
                    case (uart_addr[15:0])
                        // FIFO operations handled separately
                        UART_DATA: begin
                            // Write to TX FIFO if not full
                            if (tx_count < FIFO_DEPTH) begin
                                tx_fifo[tx_wr_ptr] <= uart_wdata[7:0];
                                tx_wr_ptr <= tx_wr_ptr + 1;
                                if (tx_wr_ptr == FIFO_DEPTH-1) tx_wr_ptr <= 0;
                            end
                        end
                        
                        // Consolidated register writes - DRY improvement
                        UART_CONTROL:     control_reg <= uart_wdata[31:0];
                        UART_BAUD_DIV:    baud_div_reg <= uart_wdata[31:0];
                        UART_FRAME_CFG:   frame_cfg_reg <= uart_wdata[31:0];
                        UART_FIFO_CTRL:   fifo_ctrl_reg <= uart_wdata[31:0];
                        UART_INT_CTRL:    int_ctrl_reg <= uart_wdata[31:0];
                        UART_DMA_CTRL:    dma_ctrl_reg <= uart_wdata[31:0];
                        UART_MODEM_CTRL:  modem_ctrl_reg <= uart_wdata[31:0];
                        UART_BREAK_CTRL:  break_ctrl_reg <= uart_wdata[31:0];
                        UART_INT_STATUS:  int_status_reg <= int_status_reg & ~uart_wdata[31:0]; // Write 1 to clear
                        
                        default: begin
                            error_count_reg <= error_count_reg + 1;
                        end
                    endcase
                end
            end
            
            UART_ERROR: begin                if (uart_timeout_counter > 0) begin
                    uart_timeout_counter <= uart_timeout_counter - 1;
                end
            end
        endcase
    end
end

// =============================================================================
// INPUT SYNCHRONIZATION
// =============================================================================

always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        uart_rx_sync1 <= 1'b1;
        uart_rx_sync2 <= 1'b1;
        uart_rx_sync3 <= 1'b1;
        uart_cts_sync1 <= 1'b0;
        uart_cts_sync2 <= 1'b0;
    end else begin
        // Three-stage synchronizer for RX (for edge detection)
        uart_rx_sync1 <= uart_rx;
        uart_rx_sync2 <= uart_rx_sync1;
        uart_rx_sync3 <= uart_rx_sync2;
        
        // Two-stage synchronizer for CTS
        uart_cts_sync1 <= uart_cts;
        uart_cts_sync2 <= uart_cts_sync1;
    end
end

// =============================================================================
// UART TRANSMITTER
// =============================================================================

always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        tx_state <= 4'h0;
        tx_bit_count <= 4'h0;
        tx_baud_count <= 16'h0000;
        tx_shift_reg <= 8'h00;
        tx_parity_bit <= 1'b0;
        uart_tx <= 1'b1;
        tx_active <= 1'b0;
        tx_rd_ptr <= 0;
        tx_count_reg <= 32'h00000000;
    end else begin
        if (control_reg[CTRL_ENABLE] && control_reg[CTRL_TX_ENABLE]) begin
            case (tx_state)
                4'h0: begin // IDLE
                    uart_tx <= 1'b1;
                    tx_active <= 1'b0;
                    
                    // Start transmission if FIFO not empty and CTS allows
                    if (tx_count > 0 && (!control_reg[CTRL_FLOW_CTRL] || uart_cts_sync2)) begin
                        tx_shift_reg <= tx_fifo[tx_rd_ptr];
                        tx_rd_ptr <= tx_rd_ptr + 1;
                        if (tx_rd_ptr == FIFO_DEPTH-1) tx_rd_ptr <= 0;
                        tx_state <= 4'h1; // START BIT
                        tx_baud_count <= baud_div_reg[15:0];
                        tx_active <= 1'b1;
                        
                        // Calculate parity
                        case (frame_cfg_reg[3:2]) // Parity type
                            PARITY_NONE: tx_parity_bit <= 1'b0;
                            PARITY_ODD:  tx_parity_bit <= ~(^tx_fifo[tx_rd_ptr]);
                            PARITY_EVEN: tx_parity_bit <= ^tx_fifo[tx_rd_ptr];
                            PARITY_MARK: tx_parity_bit <= 1'b1;
                        endcase
                    end
                end
                
                4'h1: begin // START BIT
                    uart_tx <= 1'b0;
                    if (tx_baud_count > 0) begin
                        tx_baud_count <= tx_baud_count - 1;
                    end else begin
                        tx_state <= 4'h2; // DATA BITS
                        tx_bit_count <= 4'h0;
                        tx_baud_count <= baud_div_reg[15:0];
                    end
                end
                
                4'h2: begin // DATA BITS
                    uart_tx <= tx_shift_reg[0];
                    if (tx_baud_count > 0) begin
                        tx_baud_count <= tx_baud_count - 1;
                    end else begin
                        tx_shift_reg <= {1'b0, tx_shift_reg[7:1]};
                        tx_bit_count <= tx_bit_count + 1;
                        tx_baud_count <= baud_div_reg[15:0];
                        
                        // Check if all data bits sent
                        if (tx_bit_count == (frame_cfg_reg[1:0] + 4)) begin
                            if (frame_cfg_reg[3:2] != PARITY_NONE) begin
                                tx_state <= 4'h3; // PARITY BIT
                            end else begin
                                tx_state <= 4'h4; // STOP BIT
                            end
                        end
                    end
                end
                
                4'h3: begin // PARITY BIT
                    uart_tx <= tx_parity_bit;
                    if (tx_baud_count > 0) begin
                        tx_baud_count <= tx_baud_count - 1;
                    end else begin
                        tx_state <= 4'h4; // STOP BIT
                        tx_baud_count <= baud_div_reg[15:0];
                    end
                end
                
                4'h4: begin // STOP BIT
                    uart_tx <= 1'b1;
                    if (tx_baud_count > 0) begin
                        tx_baud_count <= tx_baud_count - 1;
                    end else begin
                        tx_state <= 4'h0; // IDLE
                        tx_count_reg <= tx_count_reg + 1;
                        
                        // Generate interrupt if enabled
                        if (int_ctrl_reg[INT_TX_EMPTY] && tx_count == 1) begin
                            int_status_reg[INT_TX_EMPTY] <= 1'b1;
                        end
                    end
                end
            endcase
        end else begin
            tx_state <= 4'h0;
            uart_tx <= 1'b1;
            tx_active <= 1'b0;
        end
    end
end

// =============================================================================
// UART RECEIVER
// =============================================================================

always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        rx_state <= 4'h0;
        rx_bit_count <= 4'h0;
        rx_baud_count <= 16'h0000;
        rx_shift_reg <= 8'h00;
        rx_parity_bit <= 1'b0;
        rx_active <= 1'b0;
        rx_wr_ptr <= 0;
        rx_count_total <= 32'h00000000;
        break_counter <= 16'h0000;
        break_detected <= 1'b0;
    end else begin
        if (control_reg[CTRL_ENABLE] && control_reg[CTRL_RX_ENABLE]) begin
            case (rx_state)
                4'h0: begin // IDLE
                    rx_active <= 1'b0;
                    
                    // Detect start bit (falling edge)
                    if (uart_rx_sync3 && !uart_rx_sync2) begin
                        rx_state <= 4'h1; // START BIT
                        rx_baud_count <= baud_div_reg[15:1]; // Half bit time
                        rx_active <= 1'b1;
                        break_counter <= 16'h0000;
                    end
                    
                    // Break detection (RX low for extended period)
                    if (!uart_rx_sync2) begin
                        if (break_counter < baud_div_reg[15:0] * 10) begin
                            break_counter <= break_counter + 1;
                        end else begin
                            break_detected <= 1'b1;
                            if (int_ctrl_reg[INT_BREAK]) begin
                                int_status_reg[INT_BREAK] <= 1'b1;
                            end
                        end
                    end else begin
                        break_counter <= 16'h0000;
                        break_detected <= 1'b0;
                    end
                end
                
                4'h1: begin // START BIT
                    if (rx_baud_count > 0) begin
                        rx_baud_count <= rx_baud_count - 1;
                    end else begin
                        // Sample in middle of start bit
                        if (!uart_rx_sync2) begin
                            rx_state <= 4'h2; // DATA BITS
                            rx_bit_count <= 4'h0;
                            rx_baud_count <= baud_div_reg[15:0];
                            rx_shift_reg <= 8'h00;
                        end else begin
                            // False start, return to idle
                            rx_state <= 4'h0;
                        end
                    end
                end
                
                4'h2: begin // DATA BITS
                    if (rx_baud_count > 0) begin
                        rx_baud_count <= rx_baud_count - 1;
                    end else begin
                        // Sample data bit
                        rx_shift_reg <= {uart_rx_sync2, rx_shift_reg[7:1]};
                        rx_bit_count <= rx_bit_count + 1;
                        rx_baud_count <= baud_div_reg[15:0];
                        
                        // Check if all data bits received
                        if (rx_bit_count == (frame_cfg_reg[1:0] + 4)) begin
                            if (frame_cfg_reg[3:2] != PARITY_NONE) begin
                                rx_state <= 4'h3; // PARITY BIT
                            end else begin
                                rx_state <= 4'h4; // STOP BIT
                            end
                        end
                    end
                end
                
                4'h3: begin // PARITY BIT
                    if (rx_baud_count > 0) begin
                        rx_baud_count <= rx_baud_count - 1;
                    end else begin
                        rx_parity_bit <= uart_rx_sync2;
                        rx_state <= 4'h4; // STOP BIT
                        rx_baud_count <= baud_div_reg[15:0];
                    end
                end
                
                4'h4: begin // STOP BIT
                    if (rx_baud_count > 0) begin
                        rx_baud_count <= rx_baud_count - 1;
                    end else begin
                        rx_state <= 4'h0; // IDLE
                        
                        // Check stop bit and store data if valid
                        if (uart_rx_sync2) begin
                            // Valid frame, store in FIFO if not full
                            if (rx_count_reg < FIFO_DEPTH) begin
                                rx_fifo[rx_wr_ptr] <= rx_shift_reg;
                                rx_wr_ptr <= rx_wr_ptr + 1;
                                if (rx_wr_ptr == FIFO_DEPTH-1) rx_wr_ptr <= 0;
                                rx_count_total <= rx_count_total + 1;
                                
                                // Generate interrupt if enabled
                                if (int_ctrl_reg[INT_RX_READY]) begin
                                    int_status_reg[INT_RX_READY] <= 1'b1;
                                end
                            end else begin
                                // Overrun error
                                status_reg[STAT_OVERRUN_ERR] <= 1'b1;
                                if (int_ctrl_reg[INT_ERROR]) begin
                                    int_status_reg[INT_ERROR] <= 1'b1;
                                end
                            end
                        end else begin
                            // Frame error
                            status_reg[STAT_FRAME_ERR] <= 1'b1;
                            if (int_ctrl_reg[INT_ERROR]) begin
                                int_status_reg[INT_ERROR] <= 1'b1;
                            end
                        end
                          // Check parity if enabled
                        if (frame_cfg_reg[3:2] != PARITY_NONE) begin
                            // Use module-level expected_parity_bit register
                            case (frame_cfg_reg[3:2])
                                PARITY_ODD:  expected_parity_bit = ~(^rx_shift_reg);
                                PARITY_EVEN: expected_parity_bit = ^rx_shift_reg;
                                PARITY_MARK: expected_parity_bit = 1'b1;
                                default:     expected_parity_bit = 1'b0;
                            endcase
                            
                            if (rx_parity_bit != expected_parity_bit) begin
                                status_reg[STAT_PARITY_ERR] <= 1'b1;
                                if (int_ctrl_reg[INT_ERROR]) begin
                                    int_status_reg[INT_ERROR] <= 1'b1;
                                end
                            end
                        end
                    end
                end
            endcase
        end else begin
            rx_state <= 4'h0;
            rx_active <= 1'b0;
        end
    end
end

// =============================================================================
// FIFO MANAGEMENT
// =============================================================================

// FIFO counters
always @(*) begin
    // TX FIFO count
    if (tx_wr_ptr >= tx_rd_ptr) begin
        tx_count = tx_wr_ptr - tx_rd_ptr;
    end else begin
        tx_count = FIFO_DEPTH - tx_rd_ptr + tx_wr_ptr;
    end
    
    // RX FIFO count
    if (rx_wr_ptr >= rx_rd_ptr) begin
        rx_count_reg = rx_wr_ptr - rx_rd_ptr;
    end else begin
        rx_count_reg = FIFO_DEPTH - rx_rd_ptr + rx_wr_ptr;
    end
end

// =============================================================================
// OUTPUT CONTROL LOGIC
// =============================================================================

// CPU interface outputs
always @(*) begin
    // Default values
    uart_ready = 1'b0;
    uart_error = 1'b0;
    uart_rdata = {`VTX1_WORD_WIDTH{1'b0}};
    
    case (uart_state_reg)
        UART_IDLE: begin
            uart_ready = !uart_req;
        end
          UART_READ: begin
            uart_ready = 1'b1;
            uart_error = 1'b0;
            
            case (uart_addr[15:0])
                // FIFO operations handled separately  
                UART_DATA: begin
                    // Read from RX FIFO if not empty
                    if (rx_count_reg > 0) begin
                        uart_rdata = {{(`VTX1_WORD_WIDTH-8){1'b0}}, rx_fifo[rx_rd_ptr]};
                        // Note: FIFO read pointer updated in separate always block
                    end else begin
                        uart_rdata = {`VTX1_WORD_WIDTH{1'b0}};
                    end
                end
                
                // Consolidated register reads - DRY improvement
                UART_STATUS:      uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, status_reg};
                UART_CONTROL:     uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, control_reg};
                UART_BAUD_DIV:    uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, baud_div_reg};
                UART_FRAME_CFG:   uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, frame_cfg_reg};
                UART_FIFO_CTRL:   uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, fifo_ctrl_reg};
                UART_INT_CTRL:    uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, int_ctrl_reg};
                UART_INT_STATUS:  uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, int_status_reg};
                UART_DMA_CTRL:    uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, dma_ctrl_reg};
                UART_LINE_STATUS: uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, line_status_reg};
                UART_MODEM_CTRL:  uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, modem_ctrl_reg};
                UART_MODEM_STATUS:uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, modem_status_reg};
                UART_BREAK_CTRL:  uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, break_ctrl_reg};
                UART_TX_FIFO_LVL: uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, 24'h000000, tx_count[7:0]};
                UART_RX_FIFO_LVL: uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, 24'h000000, rx_count_reg[7:0]};
                UART_ERROR_COUNT: uart_rdata = {{(`VTX1_WORD_WIDTH-32){1'b0}}, 24'h000000, error_count_reg};
                
                default: begin
                    uart_error = 1'b1;
                    uart_rdata = {`VTX1_WORD_WIDTH{1'b0}};
                end
            endcase
        end
        
        UART_WRITE: begin
            uart_ready = 1'b1;
            uart_error = 1'b0;
        end
        
        UART_ERROR: begin
            uart_ready = 1'b1;
            uart_error = 1'b1;
        end
        
        default: begin
            uart_ready = 1'b1;
            uart_error = 1'b0;
        end
    endcase
end

// RX FIFO read pointer management
always @(posedge clk or negedge rst_n) begin
    if (!rst_n) begin
        rx_rd_ptr_reg <= 0;
    end else begin
        if (uart_state_reg == UART_READ && uart_addr[15:0] == UART_DATA && rx_count_reg > 0) begin
            rx_rd_ptr_reg <= rx_rd_ptr_reg + 1;
            if (rx_rd_ptr_reg == FIFO_DEPTH-1) rx_rd_ptr_reg <= 0;
        end
    end
end

// Status register updates
always @(*) begin
    status_reg[STAT_TX_EMPTY] = (tx_count == 0);
    status_reg[STAT_TX_FULL] = (tx_count == FIFO_DEPTH);
    status_reg[STAT_RX_EMPTY] = (rx_count_reg == 0);
    status_reg[STAT_RX_FULL] = (rx_count_reg == FIFO_DEPTH);
    status_reg[STAT_TX_ACTIVE] = tx_active;
    status_reg[STAT_RX_ACTIVE] = rx_active;
    status_reg[STAT_BREAK_DET] = break_detected;
end

// Hardware flow control
always @(*) begin
    if (control_reg[CTRL_FLOW_CTRL]) begin
        uart_rts = (rx_count_reg < (FIFO_DEPTH - 8)); // Assert RTS when FIFO has space
    end else begin
        uart_rts = 1'b0;
    end
end

// Interrupt output
always @(*) begin
    uart_irq = |int_status_reg[7:0];
    uart_irq_status = int_status_reg[7:0];
end

// DMA interface
always @(*) begin
    uart_dma_tx_req = dma_ctrl_reg[0] && (tx_count < (FIFO_DEPTH / 2));
    uart_dma_rx_req = dma_ctrl_reg[1] && (rx_count_reg > (FIFO_DEPTH / 2));
    uart_dma_tx_data = tx_fifo[tx_rd_ptr];
end

// Status and debug outputs
always @(*) begin
    uart_state = uart_state_reg;
    tx_count = tx_count_reg;
    rx_count = rx_count_total;
    error_count = error_count_reg;
    fifo_status = {tx_count[7:0], rx_count_reg[7:0]};
end

endmodule

`endif // UART_CONTROLLER_V

