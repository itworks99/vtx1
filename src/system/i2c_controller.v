// VTX1 Enhanced I2C Controller
// Part of the VTX1 Ternary System-on-Chip
// Enhanced with multi-master capability, multiple speeds, clock stretching

`timescale 1ns / 1ps

// Include VTX1 common infrastructure
`include "vtx1_interfaces.v"
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module i2c_controller (
    input  wire                     clk,
    input  wire                     rst_n,
    
    // CPU Interface (using VTX1 standardized widths)
    input  wire                     i2c_req,
    input  wire                     i2c_wr,
    input  wire [`VTX1_ADDR_WIDTH-1:0] i2c_addr,
    input  wire [`VTX1_WORD_WIDTH-1:0] i2c_wdata,
    output reg  [`VTX1_WORD_WIDTH-1:0] i2c_rdata,
    output reg                      i2c_ready,
    output reg  [3:0]               i2c_error,
    
    // Enhanced I2C Interface
    inout  wire                     i2c_sda,
    inout  wire                     i2c_scl,
    
    // Enhanced Interrupt System
    output reg                      i2c_irq,
    output reg  [7:0]               i2c_irq_vector,
    
    // DMA Interface
    output reg                      dma_tx_req,
    output reg                      dma_rx_req,
    input  wire                     dma_tx_ack,
    input  wire                     dma_rx_ack,
    input  wire [`VTX1_WORD_WIDTH-1:0] dma_tx_data,
    output reg  [`VTX1_WORD_WIDTH-1:0] dma_rx_data,
    
    // Debug and Status
    output reg  [3:0]               i2c_state,
    output reg  [31:0]              transfer_count,
    output reg  [15:0]              error_count,
    output reg                      bus_busy,
    output reg                      arbitration_lost,
    output reg  [7:0]               fifo_status,    output reg  [31:0]              debug_info
);    // Enhanced I2C Controller State Machine - Use VTX1 standardized constants
    localparam IDLE              = `VTX1_STATE_IDLE;
    localparam START             = `VTX1_I2C_STATE_START;
    localparam SEND_ADDR         = `VTX1_I2C_STATE_ADDR;
    localparam WAIT_ACK          = `VTX1_I2C_STATE_ACK;
    localparam SEND_DATA         = `VTX1_I2C_STATE_SEND_DATA;
    localparam RECV_DATA         = `VTX1_I2C_STATE_RECV_DATA;
    localparam SEND_ACK          = `VTX1_I2C_STATE_ACK;
    localparam STOP              = `VTX1_I2C_STATE_STOP;
    localparam RESTART           = `VTX1_STATE_SETUP;
    localparam SLAVE_LISTEN      = `VTX1_I2C_STATE_SLAVE_IDLE;
    localparam SLAVE_ADDR_MATCH  = `VTX1_I2C_STATE_SLAVE_ACTIVE;
    localparam ARBITRATION_LOST  = `VTX1_STATE_RECOVERY;
    localparam CLOCK_STRETCH     = `VTX1_STATE_WAIT;
    localparam DMA_WAIT          = `VTX1_STATE_WAIT;
    localparam ERROR             = `VTX1_STATE_ERROR;
    
    // Register Map (6-bit addressing for 64 registers)
    localparam REG_CTRL          = 6'h00;  // Control register
    localparam REG_STATUS        = 6'h01;  // Status register  
    localparam REG_DATA          = 6'h02;  // Data register
    localparam REG_SLAVE_ADDR    = 6'h03;  // Slave address
    localparam REG_CLK_DIV       = 6'h04;  // Clock divider
    localparam REG_TIMING        = 6'h05;  // Timing configuration
    localparam REG_IRQ_CTRL      = 6'h06;  // Interrupt control
    localparam REG_IRQ_STATUS    = 6'h07;  // Interrupt status
    localparam REG_FIFO_CTRL     = 6'h08;  // FIFO control
    localparam REG_FIFO_STATUS   = 6'h09;  // FIFO status
    localparam REG_DMA_CTRL      = 6'h0A;  // DMA control
    localparam REG_ERROR_STATUS  = 6'h0B;  // Error status
    localparam REG_DEBUG         = 6'h0C;  // Debug register
    localparam REG_ADDR_FILTER   = 6'h0D;  // Address filter
    localparam REG_SPEED_CTRL    = 6'h0E;  // Speed control
    
    // Configuration Registers
    reg [35:0] ctrl_reg;          // Control register
    reg [35:0] status_reg;        // Status register
    reg [35:0] slave_addr_reg;    // Slave address register
    reg [35:0] clk_div_reg;       // Clock divider register
    reg [35:0] timing_reg;        // Timing configuration
    reg [35:0] irq_ctrl_reg;      // Interrupt control
    reg [35:0] irq_status_reg;    // Interrupt status
    reg [35:0] fifo_ctrl_reg;     // FIFO control
    reg [35:0] dma_ctrl_reg;      // DMA control
    reg [35:0] error_status_reg;  // Error status
    reg [35:0] addr_filter_reg;   // Address filter
    reg [35:0] speed_ctrl_reg;    // Speed control
    
    // Control register bit definitions
    wire master_mode         = ctrl_reg[0];
    wire i2c_enable          = ctrl_reg[1];
    wire addr_10bit          = ctrl_reg[2];
    wire general_call_en     = ctrl_reg[3];
    wire clock_stretch_en    = ctrl_reg[4];
    wire smbus_mode          = ctrl_reg[5];
    wire dma_enable          = ctrl_reg[6];
    wire auto_nack           = ctrl_reg[7];
    wire repeated_start_en   = ctrl_reg[8];
    wire [2:0] speed_mode    = ctrl_reg[11:9];  // Standard/Fast/Fast+
    wire multi_master_en     = ctrl_reg[12];
    wire addr_filter_en      = ctrl_reg[13];
    wire stretch_timeout_en  = ctrl_reg[14];
    wire bus_timeout_en      = ctrl_reg[15];
    
    // Speed mode definitions
    localparam SPEED_STANDARD = 3'b000;  // 100 kHz
    localparam SPEED_FAST     = 3'b001;  // 400 kHz  
    localparam SPEED_FAST_PLUS = 3'b010; // 1 MHz
    localparam SPEED_HIGH     = 3'b011;  // 3.4 MHz
    
    // Timing register bit definitions
    wire [7:0] setup_time      = timing_reg[7:0];
    wire [7:0] hold_time       = timing_reg[15:8];
    wire [15:0] stretch_timeout = timing_reg[31:16];
    
    // FIFO Implementation (64-entry deep)
    reg [35:0] tx_fifo [0:63];
    reg [35:0] rx_fifo [0:63];
    reg [5:0] tx_fifo_wr_ptr, tx_fifo_rd_ptr;
    reg [5:0] rx_fifo_wr_ptr, rx_fifo_rd_ptr;
    reg [6:0] tx_fifo_count, rx_fifo_count;
    
    wire tx_fifo_empty = (tx_fifo_count == 0);
    wire tx_fifo_full  = (tx_fifo_count == 64);
    wire rx_fifo_empty = (rx_fifo_count == 0);
    wire rx_fifo_full  = (rx_fifo_count == 64);
    
    // Internal registers
    reg [3:0] state, next_state;
    reg [7:0] bit_counter;
    reg [15:0] clk_counter;
    reg [35:0] shift_reg_tx, shift_reg_rx;
    reg [7:0] byte_counter;
    reg transfer_active;
    reg [15:0] timeout_counter;
    reg [15:0] stretch_counter;
    
    // I2C bus control
    reg sda_out, scl_out;
    reg sda_oe, scl_oe;
    
    assign i2c_sda = sda_oe ? sda_out : 1'bz;
    assign i2c_scl = scl_oe ? scl_out : 1'bz;
    
    // Bus state detection with synchronizers
    reg [2:0] sda_sync, scl_sync;
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            sda_sync <= 3'b111;
            scl_sync <= 3'b111;
        end else begin
            sda_sync <= {sda_sync[1:0], i2c_sda};
            scl_sync <= {scl_sync[1:0], i2c_scl};
        end
    end
    
    wire sda_in = sda_sync[2];
    wire scl_in = scl_sync[2];
    wire sda_fall = sda_sync[2:1] == 2'b01;
    wire sda_rise = sda_sync[2:1] == 2'b10;
    wire scl_fall = scl_sync[2:1] == 2'b01;
    wire scl_rise = scl_sync[2:1] == 2'b10;
    
    // Start/Stop condition detection
    wire start_condition = sda_fall && scl_in;
    wire stop_condition = sda_rise && scl_in;
    
    // Clock generation
    reg [15:0] clk_div_counter;
    wire [15:0] clk_divisor = clk_div_reg[15:0];
    reg scl_enable;
    
    // Address matching logic
    wire [9:0] current_addr = addr_10bit ? shift_reg_rx[9:0] : {3'b000, shift_reg_rx[6:0]};
    wire [9:0] slave_address = addr_10bit ? slave_addr_reg[9:0] : {3'b000, slave_addr_reg[6:0]};
    wire addr_match = (current_addr == slave_address) || 
                      (general_call_en && current_addr == 0) ||
                      (addr_filter_en && addr_filter_reg[current_addr[6:0]]);
    
    // Default ternary word (all zeros)
    localparam DEFAULT_WORD = 36'h000000000; // Initialize to all zeros for proper startup
      // Unified register write logic - DRY improvement
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            ctrl_reg <= DEFAULT_WORD;
            slave_addr_reg <= DEFAULT_WORD;
            clk_div_reg <= 36'h000000100;  // Default divider = 256 (100kHz @ 25.6MHz)
            timing_reg <= 36'h001000404;   // Default timing
            irq_ctrl_reg <= 36'h015555554;  // Interrupts disabled by default (bit 0 = 0)
            fifo_ctrl_reg <= DEFAULT_WORD;
            dma_ctrl_reg <= DEFAULT_WORD;
            addr_filter_reg <= DEFAULT_WORD;
            speed_ctrl_reg <= DEFAULT_WORD;
            i2c_ready <= 1'b0;
        end else begin
            i2c_ready <= i2c_req; // Ready next cycle

            // Consolidated register write handling
            if (i2c_req && i2c_wr) begin
                case (i2c_addr[5:0])
                    REG_CTRL:         ctrl_reg <= i2c_wdata;
                    REG_SLAVE_ADDR:   slave_addr_reg <= i2c_wdata;
                    REG_CLK_DIV:      clk_div_reg <= i2c_wdata;
                    REG_TIMING:       timing_reg <= i2c_wdata;
                    REG_IRQ_CTRL:     irq_ctrl_reg <= i2c_wdata;
                    REG_FIFO_CTRL:    fifo_ctrl_reg <= i2c_wdata;
                    REG_DMA_CTRL:     dma_ctrl_reg <= i2c_wdata;
                    REG_ADDR_FILTER:  addr_filter_reg <= i2c_wdata;
                    REG_SPEED_CTRL:   speed_ctrl_reg <= i2c_wdata;
                    // DATA register handled separately for FIFO operations
                endcase
            end
        end
    end

    // Unified register read logic - DRY improvement
    always @(*) begin
        case (i2c_addr[5:0])
            REG_CTRL:         i2c_rdata = ctrl_reg;
            REG_STATUS:       i2c_rdata = status_reg;
            REG_SLAVE_ADDR:   i2c_rdata = slave_addr_reg;
            REG_CLK_DIV:      i2c_rdata = clk_div_reg;
            REG_TIMING:       i2c_rdata = timing_reg;
            REG_IRQ_CTRL:     i2c_rdata = irq_ctrl_reg;
            REG_IRQ_STATUS:   i2c_rdata = irq_status_reg;
            REG_FIFO_CTRL:    i2c_rdata = fifo_ctrl_reg;
            REG_FIFO_STATUS:  i2c_rdata = {16'h0, rx_fifo_count, tx_fifo_count, 2'b00, rx_fifo_full, rx_fifo_empty, tx_fifo_full, tx_fifo_empty};
            REG_DMA_CTRL:     i2c_rdata = dma_ctrl_reg;
            REG_ERROR_STATUS: i2c_rdata = error_status_reg;
            REG_DEBUG:        i2c_rdata = debug_info;
            REG_ADDR_FILTER:  i2c_rdata = addr_filter_reg;
            REG_SPEED_CTRL:   i2c_rdata = speed_ctrl_reg;
            REG_DATA:         i2c_rdata = rx_fifo_empty ? DEFAULT_WORD : rx_fifo[rx_fifo_rd_ptr];
            default:          i2c_rdata = DEFAULT_WORD;
        endcase
    end

    // Dedicated FIFO operations - separated for clarity
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            tx_fifo_wr_ptr <= 6'h0;
            rx_fifo_rd_ptr <= 6'h0;
            tx_fifo_count <= 7'h0;
            rx_fifo_count <= 7'h0;
        end else begin
            // TX FIFO write operation
            if (i2c_req && i2c_wr && (i2c_addr[5:0] == REG_DATA)) begin
                if (!tx_fifo_full) begin
                    tx_fifo[tx_fifo_wr_ptr] <= i2c_wdata;
                    tx_fifo_wr_ptr <= tx_fifo_wr_ptr + 1;
                    tx_fifo_count <= tx_fifo_count + 1;
                end
            end

            // RX FIFO read operation
            if (i2c_req && !i2c_wr && (i2c_addr[5:0] == REG_DATA)) begin
                if (!rx_fifo_empty) begin
                    rx_fifo_rd_ptr <= rx_fifo_rd_ptr + 1;
                    rx_fifo_count <= rx_fifo_count - 1;
                end
            end

            // FIFO operations from I2C transfer logic would go here
            // This is simplified for register bank integration
        end
    end
      // Main state machine logic
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            state <= IDLE;
            i2c_ready <= 1'b0;
            i2c_error <= `VTX1_ERROR_NONE;
            i2c_irq <= 1'b0;
            i2c_irq_vector <= 8'h00;
            bit_counter <= 8'h0;
            clk_counter <= 16'h0;
            shift_reg_tx <= DEFAULT_WORD;
            shift_reg_rx <= DEFAULT_WORD;
            byte_counter <= 8'h0;
            transfer_active <= 1'b0;
            transfer_count <= 32'h0;
            error_count <= 16'h0;
            bus_busy <= 1'b0;
            arbitration_lost <= 1'b0;
            timeout_counter <= 16'h0;
            stretch_counter <= 16'h0;
            sda_out <= 1'b1;
            scl_out <= 1'b1;
            sda_oe <= 1'b0;
            scl_oe <= 1'b0;
            scl_enable <= 1'b0;
            clk_div_counter <= 16'h0;
            
            // Initialize FIFOs
            tx_fifo_wr_ptr <= 6'h0;
            tx_fifo_rd_ptr <= 6'h0;
            rx_fifo_wr_ptr <= 6'h0;
            rx_fifo_rd_ptr <= 6'h0;
            tx_fifo_count <= 7'h0;
            rx_fifo_count <= 7'h0;
            
            // Initialize status registers
            status_reg <= DEFAULT_WORD;
            irq_status_reg <= DEFAULT_WORD;
            error_status_reg <= DEFAULT_WORD;
            
            // DMA signals
            dma_tx_req <= 1'b0;
            dma_rx_req <= 1'b0;
            dma_rx_data <= DEFAULT_WORD;
            
        end else begin
            state <= next_state;            // Timeout management - only during active transfers
            if (transfer_active && i2c_enable && (state != ERROR) && (state != IDLE) && (state != DMA_WAIT)) begin
                timeout_counter <= timeout_counter + 1;
            end else begin
                timeout_counter <= 16'h0;
            end
            
            // Clock stretching timeout - only during clock stretch state  
            if (clock_stretch_en && state == CLOCK_STRETCH) begin
                stretch_counter <= stretch_counter + 1;
            end else begin
                stretch_counter <= 16'h0;
            end
            
            // Clock generation for master mode
            if (master_mode && scl_enable) begin
                if (clk_div_counter >= clk_divisor) begin
                    clk_div_counter <= 16'h0;
                    if (state == SEND_ADDR || state == SEND_DATA || state == RECV_DATA) begin
                        scl_out <= ~scl_out;
                    end
                end else begin
                    clk_div_counter <= clk_div_counter + 1;
                end
            end
              // Bus busy detection - more conservative to avoid false positives during init
            bus_busy <= (state != IDLE) || start_condition || (transfer_active && (!sda_in && scl_in));
            
            case (state)                IDLE: begin
                    i2c_ready <= 1'b1;
                    // Clear errors when in IDLE state (always clear in IDLE)
                    i2c_error <= `VTX1_ERROR_NONE;
                    transfer_active <= 1'b0;
                    scl_enable <= 1'b0;
                    sda_oe <= 1'b0;
                    scl_oe <= 1'b0;
                    arbitration_lost <= 1'b0;
                    timeout_counter <= 16'h0;
                    stretch_counter <= 16'h0;
                    
                    if (master_mode && i2c_enable) begin
                        // Master mode operations - only start if explicitly commanded
                        if (!tx_fifo_empty && i2c_req) begin // Add i2c_req condition
                            // Load data from FIFO
                            shift_reg_tx <= tx_fifo[tx_fifo_rd_ptr];
                            tx_fifo_rd_ptr <= tx_fifo_rd_ptr + 1;
                            tx_fifo_count <= tx_fifo_count - 1;
                            transfer_count <= transfer_count + 1;
                        end
                    end else begin
                        // Slave mode - listen for start condition
                        if (start_condition) begin
                            transfer_active <= 1'b1;
                            bit_counter <= 8'h0;
                            shift_reg_rx <= DEFAULT_WORD;
                        end
                    end
                end
                
                START: begin
                    i2c_ready <= 1'b0;
                    transfer_active <= 1'b1;
                    
                    // Generate START condition
                    if (clk_counter == 0) begin
                        sda_oe <= 1'b1;
                        scl_oe <= 1'b1;
                        sda_out <= 1'b1;
                        scl_out <= 1'b1;
                    end else if (clk_counter == setup_time) begin
                        sda_out <= 1'b0; // SDA goes low while SCL is high
                    end else if (clk_counter >= (setup_time + hold_time)) begin
                        scl_out <= 1'b0; // SCL goes low after SDA
                        clk_counter <= 16'h0;
                        bit_counter <= 8'h0;
                        scl_enable <= 1'b1;
                    end
                    clk_counter <= clk_counter + 1;
                end
                
                SEND_ADDR: begin
                    // Send slave address + R/W bit
                    if (clk_div_counter == 0) begin
                        if (scl_out) begin
                            // Setup data on SCL high
                            if (addr_10bit && bit_counter < 10) begin
                                sda_out <= shift_reg_tx[9 - bit_counter];
                            end else begin
                                sda_out <= shift_reg_tx[7 - bit_counter];
                            end
                        end else begin
                            // SCL low phase
                            bit_counter <= bit_counter + 1;
                        end
                    end
                    
                    // Arbitration check
                    if (multi_master_en && sda_oe && sda_out && !sda_in) begin
                        arbitration_lost <= 1'b1;
                    end
                end
                
                WAIT_ACK: begin
                    // Wait for slave acknowledgment
                    sda_oe <= 1'b0; // Release SDA for slave ACK
                    
                    if (clk_div_counter == 0 && scl_out) begin
                        if (!sda_in) begin // ACK received (SDA low)
                            // ACK received, continue
                            bit_counter <= 8'h0;
                        end else begin // NACK received
                            if (auto_nack) begin
                                // Generate STOP condition
                            end else begin
                                i2c_error <= `VTX1_ERROR_BUS_FAULT;
                                error_count <= error_count + 1;
                            end
                        end
                    end
                end
                
                SEND_DATA: begin
                    // Send data bytes
                    if (clk_div_counter == 0) begin
                        if (scl_out) begin
                            // Setup data on SCL high
                            sda_out <= shift_reg_tx[7 - bit_counter];
                            sda_oe <= 1'b1;
                        end else begin
                            // SCL low phase
                            bit_counter <= bit_counter + 1;
                            if (bit_counter >= 7) begin
                                bit_counter <= 8'h0;
                                byte_counter <= byte_counter + 1;
                            end
                        end
                    end
                    
                    // Arbitration check
                    if (multi_master_en && sda_oe && sda_out && !sda_in) begin
                        arbitration_lost <= 1'b1;
                    end
                end
                
                RECV_DATA: begin
                    // Receive data bytes
                    sda_oe <= 1'b0; // Release SDA for slave data
                    
                    if (clk_div_counter == 0 && scl_out) begin
                        // Sample data on SCL high
                        shift_reg_rx <= {shift_reg_rx[34:0], sda_in};
                        bit_counter <= bit_counter + 1;
                        
                        if (bit_counter >= 7) begin
                            bit_counter <= 8'h0;
                            byte_counter <= byte_counter + 1;
                            // Store received byte in RX FIFO
                            if (!rx_fifo_full) begin
                                rx_fifo[rx_fifo_wr_ptr] <= {28'h0, shift_reg_rx[7:0]};
                                rx_fifo_wr_ptr <= rx_fifo_wr_ptr + 1;
                                rx_fifo_count <= rx_fifo_count + 1;
                            end
                        end
                    end
                end
                
                SEND_ACK: begin
                    // Send ACK/NACK
                    sda_oe <= 1'b1;
                    sda_out <= auto_nack ? 1'b1 : 1'b0; // NACK on last byte if auto_nack
                end
                  STOP: begin
                    // Generate STOP condition
                    if (clk_counter == 0) begin
                        sda_out <= 1'b0;
                        scl_out <= 1'b0;
                        sda_oe <= 1'b1;
                        scl_oe <= 1'b1;
                    end else if (clk_counter == setup_time) begin
                        scl_out <= 1'b1; // SCL goes high first
                    end else if (clk_counter >= (setup_time + hold_time)) begin
                        sda_out <= 1'b1; // SDA goes high while SCL is high
                        i2c_ready <= 1'b1;
                        transfer_active <= 1'b0;
                        bus_busy <= 1'b0;
                        scl_enable <= 1'b0;
                        sda_oe <= 1'b0;
                        scl_oe <= 1'b0;
                        // Generate interrupt only if interrupts are enabled and transfer was actually active
                        if (irq_ctrl_reg[0] && transfer_active) begin
                            i2c_irq <= 1'b1; // Transfer complete interrupt
                            i2c_irq_vector <= 8'h01;
                        end
                        clk_counter <= 16'h0;
                    end
                    clk_counter <= clk_counter + 1;
                end
                
                RESTART: begin
                    // Generate repeated START condition
                    if (clk_counter == 0) begin
                        sda_out <= 1'b1;
                        scl_out <= 1'b1;
                    end else if (clk_counter == setup_time) begin
                        sda_out <= 1'b0; // SDA goes low while SCL is high
                    end else if (clk_counter >= (setup_time + hold_time)) begin
                        scl_out <= 1'b0; // SCL goes low after SDA
                        bit_counter <= 8'h0;
                        clk_counter <= 16'h0;
                    end
                    clk_counter <= clk_counter + 1;
                end
                
                SLAVE_LISTEN: begin
                    // Slave mode - listen for address
                    sda_oe <= 1'b0;
                    scl_oe <= 1'b0;
                    
                    if (start_condition) begin
                        bit_counter <= 8'h0;
                        shift_reg_rx <= DEFAULT_WORD;
                        transfer_active <= 1'b1;
                    end else if (scl_rise && transfer_active) begin
                        // Sample address bits
                        shift_reg_rx <= {shift_reg_rx[34:0], sda_in};
                        bit_counter <= bit_counter + 1;
                    end
                end
                
                SLAVE_ADDR_MATCH: begin
                    // Check if address matches
                    if (addr_match) begin
                        // Send ACK
                        sda_oe <= 1'b1;
                        sda_out <= 1'b0;
                        
                        // Prepare for data transfer
                        bit_counter <= 8'h0;
                        if (!tx_fifo_empty) begin
                            shift_reg_tx <= tx_fifo[tx_fifo_rd_ptr];
                            tx_fifo_rd_ptr <= tx_fifo_rd_ptr + 1;
                            tx_fifo_count <= tx_fifo_count - 1;
                        end
                    end else begin
                        // Address doesn't match, ignore
                        sda_oe <= 1'b0;
                    end
                end
                  ARBITRATION_LOST: begin
                    // Lost arbitration, switch to slave mode
                    arbitration_lost <= 1'b1;
                    sda_oe <= 1'b0;
                    scl_oe <= 1'b0;
                    transfer_active <= 1'b0;
                    i2c_error <= `VTX1_ERROR_COLLISION;
                    error_count <= error_count + 1;
                    // Generate interrupt only if interrupts are enabled
                    if (irq_ctrl_reg[0]) begin
                        i2c_irq <= 1'b1;
                        i2c_irq_vector <= 8'h40; // Arbitration lost interrupt
                    end
                end
                
                CLOCK_STRETCH: begin
                    // Wait for slave to release SCL
                    if (scl_in) begin
                        // Clock stretching ended
                        scl_enable <= 1'b1;
                    end
                end
                
                DMA_WAIT: begin
                    if (dma_enable) begin
                        if (tx_fifo_count < 32) begin // TX FIFO half empty
                            dma_tx_req <= 1'b1;
                            if (dma_tx_ack) begin
                                if (!tx_fifo_full) begin
                                    tx_fifo[tx_fifo_wr_ptr] <= dma_tx_data;
                                    tx_fifo_wr_ptr <= tx_fifo_wr_ptr + 1;
                                    tx_fifo_count <= tx_fifo_count + 1;
                                end
                                dma_tx_req <= 1'b0;
                            end
                        end
                        
                        if (rx_fifo_count >= 32) begin // RX FIFO half full
                            dma_rx_req <= 1'b1;
                            dma_rx_data <= rx_fifo[rx_fifo_rd_ptr];
                            if (dma_rx_ack) begin
                                rx_fifo_rd_ptr <= rx_fifo_rd_ptr + 1;
                                rx_fifo_count <= rx_fifo_count - 1;
                                dma_rx_req <= 1'b0;
                            end
                        end
                    end
                end                ERROR: begin
                    // Reset all I2C signals and state for recovery
                    bus_busy <= 1'b0;
                    sda_oe <= 1'b0;
                    scl_oe <= 1'b0;
                    transfer_active <= 1'b0;
                    scl_enable <= 1'b0;
                    bit_counter <= 8'h0;
                    byte_counter <= 8'h0;
                    clk_counter <= 16'h0;
                    timeout_counter <= 16'h0;
                    stretch_counter <= 16'h0;
                    arbitration_lost <= 1'b0;
                    
                    // Only set error code and increment count once when first entering ERROR state
                    // Use a simple state-entry detection mechanism
                    if (state != ERROR) begin
                        i2c_error <= `VTX1_ERROR_BUS_FAULT;
                        error_count <= error_count + 1;
                        
                        // Generate interrupt only if interrupts are enabled
                        if (irq_ctrl_reg[0]) begin
                            i2c_irq <= 1'b1;
                            i2c_irq_vector <= 8'h80; // Error interrupt
                        end
                    end
                end
                
                default: begin
                    state <= IDLE;
                end
            endcase
        end
    end
                        // Next state logic
    always @(*) begin
        case (state)
            IDLE: begin
                if (master_mode) begin
                    if (i2c_enable && !tx_fifo_empty) begin
                        next_state = START;
                    end else if (dma_enable && (tx_fifo_count < 32 || rx_fifo_count >= 32)) begin
                        next_state = DMA_WAIT;
                    end else begin
                        next_state = IDLE;
                    end
                end else begin
                    next_state = SLAVE_LISTEN;
                end
            end
            
            START: begin
                if (clk_counter >= (setup_time + hold_time)) begin
                    next_state = SEND_ADDR;
                end else begin
                    next_state = START;
                end
            end
              SEND_ADDR: begin
                if (arbitration_lost) begin
                    next_state = ARBITRATION_LOST;
                end else if (timeout_counter >= `VTX1_TIMEOUT_CYCLES) begin
                    next_state = ERROR;
                end else if (bit_counter >= (addr_10bit ? 10 : 8)) begin
                    next_state = WAIT_ACK;
                end else begin
                    next_state = SEND_ADDR;
                end
            end
              WAIT_ACK: begin
                if (timeout_counter >= `VTX1_TIMEOUT_CYCLES) begin
                    next_state = ERROR;
                end else if (clk_div_counter == 0 && scl_out) begin
                    if (!sda_in) begin // ACK received
                        if (shift_reg_tx[0]) begin // Read operation
                            next_state = RECV_DATA;
                        end else begin // Write operation
                            next_state = SEND_DATA;
                        end
                    end else begin // NACK received
                        if (auto_nack) begin
                            next_state = STOP;
                        end else begin
                            next_state = ERROR;
                        end
                    end
                end else begin
                    next_state = WAIT_ACK;
                end
            end
              SEND_DATA: begin
                if (arbitration_lost) begin
                    next_state = ARBITRATION_LOST;
                end else if (timeout_counter >= `VTX1_TIMEOUT_CYCLES) begin
                    next_state = ERROR;
                end else if (byte_counter >= 4) begin // Sent all bytes
                    next_state = WAIT_ACK;
                end else if (bit_counter >= 7) begin
                    next_state = WAIT_ACK;
                end else begin
                    next_state = SEND_DATA;
                end
            end
              RECV_DATA: begin
                if (timeout_counter >= `VTX1_TIMEOUT_CYCLES) begin
                    next_state = ERROR;
                end else if (byte_counter >= 4) begin // Received all bytes
                    next_state = SEND_ACK;
                end else if (bit_counter >= 7) begin
                    next_state = SEND_ACK;
                end else begin
                    next_state = RECV_DATA;
                end
            end
            
            SEND_ACK: begin
                if (repeated_start_en && !tx_fifo_empty) begin
                    next_state = RESTART;
                end else begin
                    next_state = STOP;
                end
            end
            
            STOP: begin
                if (clk_counter >= (setup_time + hold_time)) begin
                    next_state = IDLE;
                end else begin
                    next_state = STOP;
                end
            end
            
            RESTART: begin
                if (clk_counter >= (setup_time + hold_time)) begin
                    next_state = SEND_ADDR;
                end else begin
                    next_state = RESTART;
                end
            end
            
            SLAVE_LISTEN: begin
                if (start_condition) begin
                    next_state = SLAVE_ADDR_MATCH;
                end else if (stop_condition) begin
                    next_state = IDLE;
                end else begin
                    next_state = SLAVE_LISTEN;
                end
            end
            
            SLAVE_ADDR_MATCH: begin
                if (bit_counter >= (addr_10bit ? 10 : 8)) begin
                    if (addr_match) begin
                        if (shift_reg_rx[0]) begin // Read operation
                            next_state = SEND_DATA;
                        end else begin // Write operation
                            next_state = RECV_DATA;
                        end
                    end else begin
                        next_state = SLAVE_LISTEN;
                    end
                end else begin
                    next_state = SLAVE_ADDR_MATCH;
                end
            end
            
            ARBITRATION_LOST: next_state = SLAVE_LISTEN;
              CLOCK_STRETCH: begin
                if (stretch_timeout_en && stretch_counter >= stretch_timeout) begin
                    next_state = ERROR;
                end else if (scl_in) begin
                    next_state = SEND_DATA; // Resume from where we left off
                end else begin
                    next_state = CLOCK_STRETCH;
                end
            end
            
            DMA_WAIT: begin
                if (!dma_enable) begin
                    next_state = IDLE;
                end else if (master_mode && !tx_fifo_empty) begin
                    next_state = START;
                end else begin
                    next_state = DMA_WAIT;
                end
            end
            
            ERROR: next_state = IDLE;
            
            default: next_state = IDLE;
        endcase
    end
    
    // Status register updates
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            status_reg <= DEFAULT_WORD;
            fifo_status <= 8'h00;
            debug_info <= 32'h00000000;
        end else begin
            // Update status register
            status_reg[0] <= transfer_active;
            status_reg[1] <= tx_fifo_empty;
            status_reg[2] <= tx_fifo_full;
            status_reg[3] <= rx_fifo_empty;
            status_reg[4] <= rx_fifo_full;
            status_reg[5] <= (i2c_error != `VTX1_ERROR_NONE);
            status_reg[6] <= i2c_irq;
            status_reg[7] <= master_mode;
            status_reg[8] <= bus_busy;
            status_reg[9] <= arbitration_lost;
            status_reg[10] <= start_condition;
            status_reg[11] <= stop_condition;
            status_reg[15:12] <= state;
            status_reg[23:16] <= bit_counter;
            status_reg[35:24] <= clk_counter[11:0];
            
            // Update FIFO status
            fifo_status <= {rx_fifo_full, rx_fifo_empty, tx_fifo_full, tx_fifo_empty, 
                           rx_fifo_count[3:0]};
            
            // Update debug info
            debug_info <= {transfer_count[15:0], error_count};
            
            // Clear interrupt on read of status register
            if (i2c_req && !i2c_wr && i2c_addr[5:0] == REG_IRQ_STATUS) begin
                i2c_irq <= 1'b0;
                irq_status_reg <= DEFAULT_WORD;
            end
        end
    end
    
    // Debug state output
    always @(*) begin
        i2c_state = state;
    end

endmodule
