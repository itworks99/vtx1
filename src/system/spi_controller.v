// VTX1 Enhanced SPI Controller
// Part of the VTX1 Ternary System-on-Chip
// Enhanced with master/slave modes, multi-slave support, DMA integration

`timescale 1ns / 1ps

// Include VTX1 common infrastructure
`include "vtx1_interfaces.v"
`include "vtx1_state_constants.v"
`include "vtx1_error_macros.v"

module spi_controller (
    input  wire                     clk,
    input  wire                     rst_n,

    // CPU Interface (using VTX1 standardized widths)
    input  wire                     spi_req,
    input  wire                     spi_wr,
    input  wire [`VTX1_ADDR_WIDTH-1:0] spi_addr,
    input  wire [`VTX1_WORD_WIDTH-1:0] spi_wdata,
    output reg  [`VTX1_WORD_WIDTH-1:0] spi_rdata,
    output reg                      spi_ready,
    output reg  [3:0]               spi_error,

    // Enhanced SPI Interface with Multi-Slave Support
    output reg                      spi_sclk,
    output reg                      spi_mosi,
    input  wire                     spi_miso,
    output reg  [7:0]               spi_cs_n,    // 8 chip selects

    // Slave Mode Interface
    input  wire                     spi_ss_n,    // Slave select input
    input  wire                     spi_sclk_in, // External clock input
    input  wire                     spi_mosi_in, // Master out, slave in
    output reg                      spi_miso_out, // Master in, slave out

    // Enhanced Interrupt System
    output reg                      spi_irq,
    output reg  [7:0]               spi_irq_vector,

    // DMA Interface
    output reg                      dma_tx_req,
    output reg                      dma_rx_req,
    input  wire                     dma_tx_ack,
    input  wire                     dma_rx_ack,
    input  wire [`VTX1_WORD_WIDTH-1:0] dma_tx_data,
    output reg  [`VTX1_WORD_WIDTH-1:0] dma_rx_data,

    // Debug and Status
    output reg  [3:0]               spi_state,
    output reg  [31:0]              transfer_count,
    output reg  [15:0]              error_count,
    output reg  [7:0]               fifo_status,
    output reg  [31:0]              debug_info
  );    // Enhanced SPI Controller State Machine - Use VTX1 standardized constants
  localparam IDLE         = `VTX1_STATE_IDLE;
  localparam SETUP        = `VTX1_SPI_STATE_SETUP;
  localparam START        = `VTX1_STATE_ACTIVE;
  localparam TRANSFER     = `VTX1_SPI_STATE_TRANSFER;
  localparam WAIT_CS      = `VTX1_SPI_STATE_HOLD;
  localparam FINISH       = `VTX1_SPI_STATE_COMPLETE;
  localparam SLAVE_IDLE   = `VTX1_STATE_IDLE;
  localparam SLAVE_ACTIVE = `VTX1_STATE_ACTIVE;
  localparam DMA_WAIT     = `VTX1_STATE_WAIT;
  localparam ERROR        = `VTX1_STATE_ERROR;

  // Register Map (6-bit addressing for 64 registers)
  localparam REG_CTRL         = 6'h00;  // Control register
  localparam REG_STATUS       = 6'h01;  // Status register
  localparam REG_DATA         = 6'h02;  // Data register
  localparam REG_CLK_DIV      = 6'h03;  // Clock divider
  localparam REG_CS_CTRL      = 6'h04;  // Chip select control
  localparam REG_TIMING       = 6'h05;  // Timing configuration
  localparam REG_IRQ_CTRL     = 6'h06;  // Interrupt control
  localparam REG_IRQ_STATUS   = 6'h07;  // Interrupt status
  localparam REG_FIFO_CTRL    = 6'h08;  // FIFO control
  localparam REG_FIFO_STATUS  = 6'h09;  // FIFO status
  localparam REG_DMA_CTRL     = 6'h0A;  // DMA control
  localparam REG_ERROR_STATUS = 6'h0B;  // Error status
  localparam REG_DEBUG        = 6'h0C;  // Debug register    // Configuration Registers - DRY improved with unified handling
  reg [35:0] ctrl_reg;          // Control register
  reg [35:0] clk_div_reg;       // Clock divider register
  reg [35:0] cs_ctrl_reg;       // Chip select control
  reg [35:0] timing_reg;        // Timing configuration
  reg [35:0] irq_ctrl_reg;      // Interrupt control
  reg [35:0] fifo_ctrl_reg;     // FIFO control
  reg [35:0] dma_ctrl_reg;      // DMA control

  // Status registers (read-only, updated by logic)
  reg [35:0] status_reg;        // Status register
  reg [35:0] irq_status_reg;    // Interrupt status
  reg [35:0] error_status_reg;  // Error status    // Unified register write logic - DRY improvement
  always @(posedge clk or negedge rst_n)
  begin
    if (!rst_n)
    begin
      ctrl_reg <= 36'h015555555;
      clk_div_reg <= 36'h000000010;  // Default divider = 16
      cs_ctrl_reg <= 36'h015555555;
      timing_reg <= 36'h000040404;   // Default timing
      irq_ctrl_reg <= 36'h015555555;
      fifo_ctrl_reg <= 36'h015555555;
      dma_ctrl_reg <= 36'h015555555;
      spi_ready <= 1'b0;
    end
    else
    begin
      spi_ready <= spi_req; // Ready next cycle

      // Consolidated register write handling
      if (spi_req && spi_wr)
      begin
        case (spi_addr[5:0])
          REG_CTRL:
            ctrl_reg <= spi_wdata;
          REG_CLK_DIV:
            clk_div_reg <= spi_wdata;
          REG_CS_CTRL:
            cs_ctrl_reg <= spi_wdata;
          REG_TIMING:
            timing_reg <= spi_wdata;
          REG_IRQ_CTRL:
            irq_ctrl_reg <= spi_wdata;
          REG_FIFO_CTRL:
            fifo_ctrl_reg <= spi_wdata;
          REG_DMA_CTRL:
            dma_ctrl_reg <= spi_wdata;
          // DATA register handled separately for FIFO operations
        endcase
      end
    end
  end

  // Unified register read logic - DRY improvement with status updates
  always @(*)
  begin
    case (spi_addr[5:0])
      REG_CTRL:
        spi_rdata = ctrl_reg;
      REG_STATUS:
        spi_rdata = status_reg;
      REG_CLK_DIV:
        spi_rdata = clk_div_reg;
      REG_CS_CTRL:
        spi_rdata = cs_ctrl_reg;
      REG_TIMING:
        spi_rdata = timing_reg;
      REG_IRQ_CTRL:
        spi_rdata = irq_ctrl_reg;
      REG_IRQ_STATUS:
        spi_rdata = irq_status_reg;
      REG_FIFO_CTRL:
        spi_rdata = fifo_ctrl_reg;
      REG_FIFO_STATUS:
        spi_rdata = {16'h0, rx_fifo_count, tx_fifo_count, 2'b00, rx_fifo_full, rx_fifo_empty, tx_fifo_full, tx_fifo_empty};
      REG_DMA_CTRL:
        spi_rdata = dma_ctrl_reg;
      REG_ERROR_STATUS:
        spi_rdata = error_status_reg;
      REG_DEBUG:
        spi_rdata = debug_info;
      REG_DATA:
        spi_rdata = rx_fifo_empty ? 36'h015555555 : rx_fifo[rx_fifo_rd_ptr];
      default:
        spi_rdata = 36'h015555555;
    endcase
  end

  // Control register bit definitions
  wire master_mode     = ctrl_reg[0];
  wire spi_enable      = ctrl_reg[1];
  wire cpol            = ctrl_reg[2];
  wire cpha            = ctrl_reg[3];
  wire lsb_first       = ctrl_reg[4];
  wire [2:0] data_width = ctrl_reg[7:5];   // 5-8 bits
  wire loopback_mode   = ctrl_reg[8];
  wire auto_cs         = ctrl_reg[9];
  wire dma_enable      = ctrl_reg[10];
  wire [2:0] cs_select = ctrl_reg[13:11];  // Active chip select
  wire continuous_mode = ctrl_reg[14];
  wire duplex_mode     = ctrl_reg[15];

  // Timing register bit definitions
  wire [7:0] cs_setup_time   = timing_reg[7:0];
  wire [7:0] cs_hold_time    = timing_reg[15:8];
  wire [7:0] inter_byte_delay = timing_reg[23:16];

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
  reg [15:0] spi_timeout_counter;
  reg [7:0] cs_timer;

  // Clock generation
  reg spi_clk_enable;
  reg [15:0] clk_div_counter;
  wire [15:0] clk_divisor = clk_div_reg[15:0];

  // Input synchronizers for slave mode
  reg [2:0] spi_ss_sync;
  reg [2:0] spi_sclk_sync;
  reg [2:0] spi_mosi_sync;

  always @(posedge clk or negedge rst_n)
  begin
    if (!rst_n)
    begin
      spi_ss_sync <= 3'b111;
      spi_sclk_sync <= 3'b000;
      spi_mosi_sync <= 3'b000;
    end
    else
    begin
      spi_ss_sync <= {spi_ss_sync[1:0], spi_ss_n};
      spi_sclk_sync <= {spi_sclk_sync[1:0], spi_sclk_in};
      spi_mosi_sync <= {spi_mosi_sync[1:0], spi_mosi_in};
    end
  end
  wire spi_ss_active = !spi_ss_sync[2];
  wire spi_sclk_edge = spi_sclk_sync[2] ^ spi_sclk_sync[1];
  wire spi_sclk_rise = spi_sclk_edge & spi_sclk_sync[2];
  wire spi_sclk_fall = spi_sclk_edge & !spi_sclk_sync[2];

  // Controller status signals
  reg spi_busy;
  reg transfer_complete;
  reg transfer_error;

  // Default ternary word (all zeros)
  localparam DEFAULT_WORD = 36'h015555555;

  // Special handling for DATA register FIFO operations
  always @(posedge clk or negedge rst_n)
  begin
    if (!rst_n)
    begin
      // FIFO pointers and counters reset
      tx_fifo_wr_ptr <= 6'h0;
      tx_fifo_rd_ptr <= 6'h0;
      rx_fifo_wr_ptr <= 6'h0;
      rx_fifo_rd_ptr <= 6'h0;
      tx_fifo_count <= 7'h0;
      rx_fifo_count <= 7'h0;
    end
    else
    begin
      // Handle DATA register writes (TX FIFO)
      if (spi_req && spi_wr && (spi_addr[5:0] == REG_DATA))
      begin
        if (!tx_fifo_full)
        begin
          tx_fifo[tx_fifo_wr_ptr] <= spi_wdata;
          tx_fifo_wr_ptr <= tx_fifo_wr_ptr + 1;
          tx_fifo_count <= tx_fifo_count + 1;
        end
      end
      // Handle DATA register reads (RX FIFO)
      if (spi_req && !spi_wr && (spi_addr[5:0] == REG_DATA))
      begin
        if (!rx_fifo_empty)
        begin
          rx_fifo_rd_ptr <= rx_fifo_rd_ptr + 1;
          rx_fifo_count <= rx_fifo_count - 1;
        end
      end

      // FIFO operations from SPI transfer logic would go here
      // This is simplified for register bank integration
    end
  end
  // Main state machine logic
  always @(posedge clk or negedge rst_n)
  begin
    if (!rst_n)
    begin
      state <= IDLE;
      spi_sclk <= 1'b0;
      spi_mosi <= 1'b0;
      spi_cs_n <= 8'hFF;
      spi_miso_out <= 1'b0;
      spi_ready <= 1'b0;
      spi_error <= `VTX1_ERROR_NONE;
      spi_irq <= 1'b0;
      spi_irq_vector <= 8'h00;
      bit_counter <= 8'h0;
      clk_counter <= 16'h0;
      shift_reg_tx <= DEFAULT_WORD;
      shift_reg_rx <= DEFAULT_WORD;
      byte_counter <= 8'h0;
      transfer_active <= 1'b0;
      transfer_count <= 32'h0;
      error_count <= 16'h0;
      spi_timeout_counter <= 16'h0;
      cs_timer <= 8'h0;
      spi_clk_enable <= 1'b0;
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

    end
    else
    begin
      state <= next_state;
      // Timeout management
      if (transfer_active)
      begin
        spi_timeout_counter <= spi_timeout_counter + 1;
        if (spi_timeout_counter >= `VTX1_TIMEOUT_CYCLES)
        begin
          spi_error <= `VTX1_ERROR_TIMEOUT;
          error_count <= error_count + 1;
        end
      end
      else
      begin
        spi_timeout_counter <= 16'h0;
      end

      // Clock generation for master mode
      if (master_mode && spi_clk_enable)
      begin
        if (clk_div_counter >= clk_divisor)
        begin
          clk_div_counter <= 16'h0;
          if (state == TRANSFER)
          begin
            spi_sclk <= ~spi_sclk;
          end
        end
        else
        begin
          clk_div_counter <= clk_div_counter + 1;
        end
      end

      case (state)
        IDLE:
        begin
          spi_ready <= 1'b1;
          spi_error <= `VTX1_ERROR_NONE;
          transfer_active <= 1'b0;
          spi_clk_enable <= 1'b0;
          spi_sclk <= cpol;

          if (master_mode)
          begin
            spi_cs_n <= 8'hFF; // Deassert all chip selects

            // Check for data to transmit
            if (spi_enable && !tx_fifo_empty)
            begin
              // Load data from FIFO
              shift_reg_tx <= tx_fifo[tx_fifo_rd_ptr];
              tx_fifo_rd_ptr <= tx_fifo_rd_ptr + 1;
              tx_fifo_count <= tx_fifo_count - 1;
              transfer_count <= transfer_count + 1;
            end
          end
          else
          begin
            // Slave mode initialization
            if (spi_ss_active)
            begin
              transfer_active <= 1'b1;
              bit_counter <= 8'h0;
              shift_reg_rx <= DEFAULT_WORD;
              shift_reg_tx <= tx_fifo_empty ? DEFAULT_WORD : tx_fifo[tx_fifo_rd_ptr];
            end
          end
        end

        SETUP:
        begin
          spi_ready <= 1'b0;
          transfer_active <= 1'b1;

          // Setup chip select
          if (auto_cs)
          begin
            spi_cs_n[cs_select] <= 1'b0;
          end

          // CS setup time
          if (cs_timer >= cs_setup_time)
          begin
            cs_timer <= 8'h0;
            bit_counter <= 8'h0;
            clk_counter <= 16'h0;
            spi_clk_enable <= 1'b1;
          end
          else
          begin
            cs_timer <= cs_timer + 1;
          end
        end

        START:
        begin
          // Initialize transfer
          spi_mosi <= lsb_first ? shift_reg_tx[0] : shift_reg_tx[35];
          spi_sclk <= cpol;
        end

        TRANSFER:
        begin
          // SPI transfer with configurable CPOL/CPHA
          if (clk_div_counter == 0)
          begin
            // Clock edge handling based on CPHA
            if ((cpha && spi_sclk != cpol) || (!cpha && spi_sclk == cpol))
            begin
              // Sample edge
              if (lsb_first)
              begin
                shift_reg_rx <= {spi_miso, shift_reg_rx[35:1]};
              end
              else
              begin
                shift_reg_rx <= {shift_reg_rx[34:0], spi_miso};
              end
            end
            else
            begin
              // Setup edge
              bit_counter <= bit_counter + 1;
              if (lsb_first)
              begin
                shift_reg_tx <= {1'b0, shift_reg_tx[35:1]};
                spi_mosi <= shift_reg_tx[1];
              end
              else
              begin
                shift_reg_tx <= {shift_reg_tx[34:0], 1'b0};
                spi_mosi <= shift_reg_tx[34];
              end
            end
          end
        end

        WAIT_CS:
        begin
          // CS hold time
          if (cs_timer >= cs_hold_time)
          begin
            cs_timer <= 8'h0;
          end
          else
          begin
            cs_timer <= cs_timer + 1;
          end
        end        FINISH:
        begin
          if (auto_cs)
          begin
            spi_cs_n <= 8'hFF; // Deassert chip select
          end
          spi_sclk <= cpol;
          spi_ready <= 1'b1;
          transfer_active <= 1'b0;
          spi_clk_enable <= 1'b0;

          // Store received data in RX FIFO
          if (!rx_fifo_full)
          begin
            rx_fifo[rx_fifo_wr_ptr] <= shift_reg_rx;
            rx_fifo_wr_ptr <= rx_fifo_wr_ptr + 1;
            rx_fifo_count <= rx_fifo_count + 1;
          end

          // Generate interrupt only if interrupts are enabled and transfer was actually active
          if (irq_ctrl_reg[0] && transfer_active) begin
            spi_irq <= 1'b1;
            spi_irq_vector <= 8'h01; // Transfer complete
            irq_status_reg[0] <= 1'b1;
          end
        end

        SLAVE_IDLE:
        begin
          spi_ready <= 1'b1;
          transfer_active <= 1'b0;

          if (spi_ss_active)
          begin
            transfer_active <= 1'b1;
            bit_counter <= 8'h0;
            shift_reg_rx <= DEFAULT_WORD;
            if (!tx_fifo_empty)
            begin
              shift_reg_tx <= tx_fifo[tx_fifo_rd_ptr];
              tx_fifo_rd_ptr <= tx_fifo_rd_ptr + 1;
              tx_fifo_count <= tx_fifo_count - 1;
            end
            else
            begin
              shift_reg_tx <= DEFAULT_WORD;
            end
          end
        end

        SLAVE_ACTIVE:
        begin
          spi_ready <= 1'b0;

          // Handle slave mode transfers
          if ((cpha && spi_sclk_rise) || (!cpha && spi_sclk_fall))
          begin
            // Sample edge
            if (lsb_first)
            begin
              shift_reg_rx <= {spi_mosi_sync[2], shift_reg_rx[35:1]};
            end
            else
            begin
              shift_reg_rx <= {shift_reg_rx[34:0], spi_mosi_sync[2]};
            end
            bit_counter <= bit_counter + 1;
          end
          else if ((cpha && spi_sclk_fall) || (!cpha && spi_sclk_rise))
          begin
            // Setup edge
            if (lsb_first)
            begin
              spi_miso_out <= shift_reg_tx[0];
              shift_reg_tx <= {1'b0, shift_reg_tx[35:1]};
            end
            else
            begin
              spi_miso_out <= shift_reg_tx[35];
              shift_reg_tx <= {shift_reg_tx[34:0], 1'b0};
            end
          end          if (!spi_ss_active)
          begin
            // Transfer complete
            if (!rx_fifo_full)
            begin
              rx_fifo[rx_fifo_wr_ptr] <= shift_reg_rx;
              rx_fifo_wr_ptr <= rx_fifo_wr_ptr + 1;
              rx_fifo_count <= rx_fifo_count + 1;
            end
            transfer_count <= transfer_count + 1;
            // Generate interrupt only if interrupts are enabled
            if (irq_ctrl_reg[0]) begin
              spi_irq <= 1'b1;
              spi_irq_vector <= 8'h01;
            end
          end
        end

        DMA_WAIT:
        begin
          if (dma_enable)
          begin
            if (tx_fifo_count < 32)
            begin // TX FIFO half empty
              dma_tx_req <= 1'b1;
              if (dma_tx_ack)
              begin
                if (!tx_fifo_full)
                begin
                  tx_fifo[tx_fifo_wr_ptr] <= dma_tx_data;
                  tx_fifo_wr_ptr <= tx_fifo_wr_ptr + 1;
                  tx_fifo_count <= tx_fifo_count + 1;
                end
                dma_tx_req <= 1'b0;
              end
            end

            if (rx_fifo_count >= 32)
            begin // RX FIFO half full
              dma_rx_req <= 1'b1;
              dma_rx_data <= rx_fifo[rx_fifo_rd_ptr];
              if (dma_rx_ack)
              begin
                rx_fifo_rd_ptr <= rx_fifo_rd_ptr + 1;
                rx_fifo_count <= rx_fifo_count - 1;
                dma_rx_req <= 1'b0;
              end
            end
          end
        end

        ERROR:
        begin
          spi_error <= `VTX1_ERROR_BUS_FAULT;
          error_count <= error_count + 1;
          spi_cs_n <= 8'hFF;
          transfer_active <= 1'b0;
          spi_irq <= 1'b1;
          spi_irq_vector <= 8'h80; // Error interrupt
        end

        default:
        begin
          state <= IDLE;
        end
      endcase
    end
  end
  // Next state logic
  always @(*)
  begin
    case (state)
      IDLE:
      begin
        if (master_mode)
        begin
          if (spi_enable && !tx_fifo_empty)
          begin
            next_state = SETUP;
          end
          else if (dma_enable && (tx_fifo_count < 32 || rx_fifo_count >= 32))
          begin
            next_state = DMA_WAIT;
          end
          else
          begin
            next_state = IDLE;
          end
        end
        else
        begin
          if (spi_ss_active)
          begin
            next_state = SLAVE_ACTIVE;
          end
          else
          begin
            next_state = SLAVE_IDLE;
          end
        end
      end

      SETUP:
      begin
        if (cs_timer >= cs_setup_time)
        begin
          next_state = START;
        end
        else
        begin
          next_state = SETUP;
        end
      end

      START:
        next_state = TRANSFER;
      TRANSFER:
      begin
        if (bit_counter >= 36)
        begin
          next_state = WAIT_CS;
        end
        else if (spi_timeout_counter >= `VTX1_TIMEOUT_CYCLES)
        begin
          next_state = ERROR;
        end
        else
        begin
          next_state = TRANSFER;
        end
      end

      WAIT_CS:
      begin
        if (cs_timer >= cs_hold_time)
        begin
          next_state = FINISH;
        end
        else
        begin
          next_state = WAIT_CS;
        end
      end

      FINISH:
      begin
        if (continuous_mode && !tx_fifo_empty)
        begin
          next_state = SETUP;
        end
        else
        begin
          next_state = IDLE;
        end
      end

      SLAVE_IDLE:
      begin
        if (spi_ss_active)
        begin
          next_state = SLAVE_ACTIVE;
        end
        else
        begin
          next_state = SLAVE_IDLE;
        end
      end

      SLAVE_ACTIVE:
      begin
        if (!spi_ss_active)
        begin
          next_state = SLAVE_IDLE;
        end
        else if (bit_counter >= 36)
        begin
          next_state = SLAVE_IDLE;
        end
        else
        begin
          next_state = SLAVE_ACTIVE;
        end
      end

      DMA_WAIT:
      begin
        if (!dma_enable)
        begin
          next_state = IDLE;
        end
        else if (master_mode && !tx_fifo_empty)
        begin
          next_state = SETUP;
        end
        else
        begin
          next_state = DMA_WAIT;
        end
      end

      ERROR:
        next_state = IDLE;

      default:
        next_state = IDLE;
    endcase
  end

  // Status register updates
  always @(posedge clk or negedge rst_n)
  begin
    if (!rst_n)
    begin
      status_reg <= DEFAULT_WORD;
      fifo_status <= 8'h00;
      debug_info <= 32'h00000000;
    end
    else
    begin
      // Update status register
      status_reg[0] <= transfer_active;
      status_reg[1] <= tx_fifo_empty;
      status_reg[2] <= tx_fifo_full;
      status_reg[3] <= rx_fifo_empty;
      status_reg[4] <= rx_fifo_full;
      status_reg[5] <= (spi_error != `VTX1_ERROR_NONE);
      status_reg[6] <= spi_irq;
      status_reg[7] <= master_mode;
      status_reg[11:8] <= state;
      status_reg[19:12] <= bit_counter;
      status_reg[35:20] <= spi_timeout_counter;

      // Update FIFO status
      fifo_status <= {rx_fifo_full, rx_fifo_empty, tx_fifo_full, tx_fifo_empty,
                      rx_fifo_count[3:0]};

      // Update debug info
      debug_info <= {transfer_count[15:0], error_count};

      // Clear interrupt on read of status register
      if (spi_req && !spi_wr && spi_addr[5:0] == REG_IRQ_STATUS)
      begin
        spi_irq <= 1'b0;
        irq_status_reg <= DEFAULT_WORD;
      end
    end
  end

  // Debug state output
  always @(*)
  begin
    spi_state = state;
  end

endmodule
