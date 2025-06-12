`timescale 1ns / 1ps
// =============================================================================
// VTX1 System-on-Chip Top Level - CONSOLIDATED VERSION
// =============================================================================
// Complete VTX1 ternary system integration with consolidated interfaces
// Uses standardized interface definitions to eliminate code duplication
// Compatible with Icarus Verilog simulation
// =============================================================================

module vtx1_top_consolidated (
    // External Clock and Reset
    input wire clk_ext,                // 27 MHz external oscillator
    input wire clk_32k,                // 32.768 kHz RTC clock
    input wire rst_ext_n,              // External reset (active low)
    
    // External Memory Interface
    output wire [35:0] mem_addr,
    inout wire [35:0] mem_data,
    output wire mem_we_n,
    output wire mem_oe_n,
    output wire mem_cs_n,
    
    // Flash Memory Interface (SPI)
    output wire flash_cs_n,
    output wire flash_sck,
    output wire flash_mosi,
    input wire flash_miso,
    output wire flash_wp_n,
    output wire flash_hold_n,
    
    // GPIO Interface
    // ===================================================================
    // ENHANCED GPIO INTERFACE (24-PIN)
    // ===================================================================
    inout wire [23:0] gpio_pins,
      // UART Interface
    input wire uart_rx,
    output wire uart_tx,
    output wire uart_rts,
    input wire uart_cts,
    
    // SPI Interface
    output wire spi_sclk,
    output wire spi_mosi,
    input wire spi_miso,
    output wire [7:0] spi_cs_n,
    // SPI Slave Mode
    input wire spi_ss_n,
    input wire spi_sclk_in,
    input wire spi_mosi_in,
    output wire spi_miso_out,
    
    // I2C Interface
    inout wire i2c_sda,
    inout wire i2c_scl,
    
    // Debug Interface (JTAG)
    input wire jtag_tck,
    input wire jtag_tms,
    input wire jtag_tdi,
    output wire jtag_tdo,
    input wire jtag_trst_n,
    
    // Power Management
    input wire pwr_good,
    output wire pwr_enable,
    input wire wake_event,
    output wire sleep_ack,
    
    // Status LEDs
    output wire led_power,
    output wire led_activity,
    output wire led_error,
    
    // Test Points
    output wire tp_clk_cpu,
    output wire tp_clk_mem,
    output wire tp_rst_n,
    output wire [3:0] tp_state
);

// =============================================================================
// VTX1 STANDARDIZED PARAMETERS - CONSOLIDATED
// =============================================================================
localparam VTX1_WORD_WIDTH = 36;        // 36-bit ternary word
localparam VTX1_CACHE_LINE_WIDTH = 288; // 288-bit cache line (8 words)
localparam VTX1_ADDR_WIDTH = 36;        // 36-bit addressing
localparam VTX1_NUM_REGISTERS = 13;
localparam VTX1_PIPELINE_STAGES = 4;
localparam VTX1_VLIW_SLOTS = 3;
localparam VTX1_VLIW_WIDTH = 108;       // 108-bit VLIW instruction (3 x 36-bit slots)

// MMIO Address Map
localparam VTX1_GPIO_BASE_ADDR  = 16'h1000;
localparam VTX1_UART_BASE_ADDR  = 16'h1001;
localparam VTX1_SPI_BASE_ADDR   = 16'h1002;
localparam VTX1_I2C_BASE_ADDR   = 16'h1003;
localparam VTX1_TIMER_BASE_ADDR = 16'h1004;
localparam VTX1_FLASH_BASE_ADDR = 16'h1005;

// =============================================================================
// CONSOLIDATED SIGNAL DECLARATIONS
// =============================================================================

// Internal Clock and Reset Signals (consolidated)
wire clk_cpu, clk_mem, clk_periph, clk_debug, clk_flash;
wire rst_cpu_n, rst_mem_n, rst_periph_n, rst_debug_n, rst_flash_n, rst_system_n;

// CPU Core Interface (direct CPU interface)
wire cpu_enable;
wire [VTX1_WORD_WIDTH-1:0] cpu_imem_addr, cpu_dmem_addr;
wire [VTX1_VLIW_WIDTH-1:0] cpu_imem_data;
wire [VTX1_WORD_WIDTH-1:0] cpu_dmem_data;
wire cpu_imem_req, cpu_imem_ready, cpu_dmem_req, cpu_dmem_ready;
wire cpu_dmem_we, cpu_dmem_oe;
wire [15:0] cpu_interrupt_req, cpu_interrupt_ack;
wire cpu_nmi_req, cpu_debug_enable, cpu_debug_step;
wire [3:0] cpu_debug_cmd, cpu_debug_status;
wire [VTX1_WORD_WIDTH-1:0] cpu_debug_pc;
wire [31:0] cpu_cycle_count, cpu_instruction_count, cpu_stall_count;
wire cpu_pipeline_stall, cpu_pipeline_flush;
wire [1:0] cpu_pipeline_state;
wire cpu_error, cpu_timeout;
wire [3:0] cpu_error_code;
wire [31:0] cpu_error_count;

// CPU Adapter Interface (CPU to Bus Matrix)
wire cpu_bus_req, cpu_bus_wr, cpu_bus_ready, cpu_bus_error, cpu_bus_timeout, cpu_bus_error_clear;
wire [1:0] cpu_bus_size;
wire [VTX1_ADDR_WIDTH-1:0] cpu_bus_addr;
wire [VTX1_WORD_WIDTH-1:0] cpu_bus_wdata, cpu_bus_rdata;
wire [3:0] cpu_bus_error_code;

// DMA Controller Interface (to Bus Matrix)
wire dma_req, dma_wr, dma_ready, dma_error, dma_timeout, dma_error_clear;
wire [1:0] dma_size;
wire [VTX1_ADDR_WIDTH-1:0] dma_addr;
wire [VTX1_WORD_WIDTH-1:0] dma_wdata, dma_rdata;
wire [3:0] dma_error_code;

// DMA Controller Control Signals
wire dma_enable, dma_start, dma_complete, dma_irq;
wire [7:0] dma_channel, dma_status, dma_error_count;
wire [VTX1_ADDR_WIDTH-1:0] dma_src_addr, dma_dest_addr;
wire [31:0] dma_transfer_count, dma_transfer_progress, dma_operations_count;
wire [2:0] dma_transfer_mode;
wire [3:0] dma_state;

// Debug Controller Interface (to Bus Matrix)
wire debug_req, debug_wr, debug_ready, debug_error, debug_timeout, debug_error_clear;
wire [1:0] debug_size;
wire [VTX1_ADDR_WIDTH-1:0] debug_addr;
wire [VTX1_WORD_WIDTH-1:0] debug_wdata, debug_rdata;
wire [3:0] debug_error_code;

// Debug Controller Control Signals
wire debug_enable_ctrl, debug_halt_cpu, debug_step_cpu, debug_reset_cpu;
wire [VTX1_ADDR_WIDTH-1:0] debug_access_addr;
wire [VTX1_WORD_WIDTH-1:0] debug_access_data;
wire debug_access_wr, debug_access_req, debug_access_ready;
wire [3:0] debug_state;
wire [31:0] debug_operations_count;

// Bus Matrix to Memory Controller Interface
wire bus_mem_req, bus_mem_wr, bus_mem_ready, bus_mem_error;
wire [1:0] bus_mem_size;
wire [VTX1_ADDR_WIDTH-1:0] bus_mem_addr;
wire [VTX1_WORD_WIDTH-1:0] bus_mem_wdata, bus_mem_rdata;

// Bus Matrix Control and Status Signals
wire [1:0] bus_arbitration_mode;
wire [7:0] bus_priority_config;
wire [15:0] bus_timeout_config;
wire bus_deadlock_enable, bus_performance_enable;
wire [3:0] bus_matrix_state;
wire [2:0] bus_current_master, bus_current_slave;
wire bus_arbitration_active, bus_deadlock_detected, bus_deadlock_recovery;

// Bus Matrix Performance Monitoring
wire [31:0] bus_total_transactions, bus_cpu_transactions, bus_dma_transactions;
wire [31:0] bus_debug_transactions, bus_timeout_count, bus_error_count, bus_deadlock_count;
wire [15:0] bus_utilization, bus_avg_latency, bus_max_latency;
wire [31:0] bus_cpu_wait_cycles, bus_dma_wait_cycles, bus_debug_wait_cycles;

// Cache Controller Interface (standardized)
wire cache_req, cache_wr, cache_ready, cache_error;
wire icache_req, icache_hit, icache_ready;
wire dcache_req, dcache_wr, dcache_hit, dcache_ready;
wire [1:0] dcache_size;
wire [VTX1_ADDR_WIDTH-1:0] cache_addr, icache_addr, dcache_addr;
wire [VTX1_CACHE_LINE_WIDTH-1:0] cache_wdata, cache_rdata;
wire [VTX1_WORD_WIDTH-1:0] icache_data, dcache_wdata, dcache_rdata;

// Cache Statistics and Control
wire [31:0] icache_hits, icache_misses, dcache_hits, dcache_misses;

// Memory Controller to Physical Memory (standardized)
wire phy_req, phy_wr, phy_ready, phy_error;
wire [VTX1_ADDR_WIDTH-1:0] phy_addr;
wire [VTX1_WORD_WIDTH-1:0] phy_wdata, phy_rdata;

// MMIO Interface (consolidated with router)
wire mmio_req, mmio_wr;
wire mmio_ready_mc, mmio_error_mc;  // From memory controller
wire [VTX1_ADDR_WIDTH-1:0] mmio_addr;
wire [VTX1_WORD_WIDTH-1:0] mmio_wdata;
wire [VTX1_WORD_WIDTH-1:0] mmio_rdata_mc;  // From memory controller
reg mmio_ready, mmio_error;           // Final MMIO outputs
reg [VTX1_WORD_WIDTH-1:0] mmio_rdata; // Final MMIO output

// GPIO Interface (standardized)
wire gpio_req, gpio_wr, gpio_ready, gpio_error, gpio_irq;
wire [VTX1_ADDR_WIDTH-1:0] gpio_addr;
wire [VTX1_WORD_WIDTH-1:0] gpio_wdata, gpio_rdata;

// Enhanced UART Interface
wire uart_req, uart_wr, uart_ready, uart_error, uart_irq;
wire [VTX1_ADDR_WIDTH-1:0] uart_addr;
wire [VTX1_WORD_WIDTH-1:0] uart_wdata, uart_rdata;
wire [7:0] uart_irq_status;
wire uart_dma_tx_req, uart_dma_rx_req, uart_dma_tx_ack, uart_dma_rx_ack;
wire [VTX1_WORD_WIDTH-1:0] uart_dma_tx_data, uart_dma_rx_data;

// Enhanced SPI Interface
wire spi_req, spi_wr, spi_ready, spi_irq;
wire [3:0] spi_error;
wire [VTX1_ADDR_WIDTH-1:0] spi_addr;
wire [VTX1_WORD_WIDTH-1:0] spi_wdata, spi_rdata;
wire [7:0] spi_irq_vector;
wire spi_dma_tx_req, spi_dma_rx_req, spi_dma_tx_ack, spi_dma_rx_ack;
wire [VTX1_WORD_WIDTH-1:0] spi_dma_tx_data, spi_dma_rx_data;
wire [3:0] spi_state;
wire [31:0] spi_transfer_count;
wire [15:0] spi_error_count;
wire [7:0] spi_fifo_status;
wire [31:0] spi_debug_info;

// Enhanced I2C Interface
wire i2c_req, i2c_wr, i2c_ready, i2c_irq;
wire [3:0] i2c_error;
wire [VTX1_ADDR_WIDTH-1:0] i2c_addr;
wire [VTX1_WORD_WIDTH-1:0] i2c_wdata, i2c_rdata;
wire [7:0] i2c_irq_vector;
wire i2c_dma_tx_req, i2c_dma_rx_req, i2c_dma_tx_ack, i2c_dma_rx_ack;
wire [VTX1_WORD_WIDTH-1:0] i2c_dma_tx_data, i2c_dma_rx_data;
wire [3:0] i2c_state;
wire [31:0] i2c_transfer_count;
wire [15:0] i2c_error_count;
wire i2c_bus_busy, i2c_arbitration_lost;
wire [7:0] i2c_fifo_status;
wire [31:0] i2c_debug_info;

// Flash Controller Interface (standardized)
wire flash_req, flash_wr, flash_ready, flash_error;
wire [1:0] flash_op;
wire [VTX1_ADDR_WIDTH-1:0] flash_addr;
wire [VTX1_WORD_WIDTH-1:0] flash_wdata, flash_rdata;

// Interrupt Controller Interface (standardized)
wire irq;
wire [7:0] irq_vector;
wire [1:0] irq_level;
wire [15:0] ext_irq;
wire timer_irq, cache_irq, mem_irq, flash_irq, debug_irq;

// Control and Configuration (using VTX1 parameters)
wire pll_enable, cache_enable, cache_flush, flash_enable, irq_enable;
wire [1:0] cpu_clk_div, mem_clk_div, periph_clk_div, cache_policy, priority_mode;
wire [31:0] irq_mask, irq_pending, irq_active;
wire [63:0] priority_config;
wire low_power_mode, sleep_req;

// Status and Debug (using VTX1 parameters)
wire pll_locked, clk_gated;
wire [2:0] clk_status;
wire [3:0] mc_state, cache_state, ic_state, cm_state, gpio_state;
wire [31:0] access_count, error_count, irq_count;
wire [31:0] clk_freq_cpu, clk_freq_mem;

// External memory data bus control (consolidated)
reg mem_data_oe;
wire [VTX1_WORD_WIDTH-1:0] mem_data_out;
assign mem_data = mem_data_oe ? mem_data_out : {VTX1_WORD_WIDTH{1'bz}};
assign phy_rdata = mem_data;

// Boot Interface
wire boot_mode, boot_valid;
wire [VTX1_ADDR_WIDTH-1:0] boot_addr;
wire [VTX1_WORD_WIDTH-1:0] boot_data;

// =============================================================================
// SIMPLIFIED CLOCK MANAGER (CONSOLIDATED)
// =============================================================================
// Simple clock distribution for consolidation demo
assign clk_cpu = clk_ext;    // Use external clock directly for simplicity
assign clk_mem = clk_ext;
assign clk_periph = clk_ext;
assign clk_debug = clk_ext;
assign clk_flash = clk_ext;

// Simple reset distribution
reg [3:0] reset_sync;
always @(posedge clk_ext or negedge rst_ext_n) begin
    if (!rst_ext_n) begin
        reset_sync <= 4'b0000;
    end else begin
        reset_sync <= {reset_sync[2:0], 1'b1};
    end
end

assign rst_system_n = reset_sync[3];
assign rst_cpu_n = rst_system_n;
assign rst_mem_n = rst_system_n;
assign rst_periph_n = rst_system_n;
assign rst_debug_n = rst_system_n;
assign rst_flash_n = rst_system_n;

// =============================================================================
// CONSOLIDATED MMIO ROUTER  
// =============================================================================
// Address decode signals
reg mmio_gpio_sel, mmio_uart_sel, mmio_spi_sel, mmio_i2c_sel, mmio_timer_sel, mmio_flash_sel;

// Address decoder (using VTX1 parameters)
always @(*) begin
    mmio_gpio_sel  = (mmio_addr[31:16] == VTX1_GPIO_BASE_ADDR);
    mmio_uart_sel  = (mmio_addr[31:16] == VTX1_UART_BASE_ADDR);
    mmio_spi_sel   = (mmio_addr[31:16] == VTX1_SPI_BASE_ADDR);
    mmio_i2c_sel   = (mmio_addr[31:16] == VTX1_I2C_BASE_ADDR);
    mmio_timer_sel = (mmio_addr[31:16] == VTX1_TIMER_BASE_ADDR);
    mmio_flash_sel = (mmio_addr[31:16] == VTX1_FLASH_BASE_ADDR);
end

// Request routing - GPIO
assign gpio_req = mmio_req && mmio_gpio_sel;
assign gpio_wr = mmio_wr;
assign gpio_addr = mmio_addr;
assign gpio_wdata = mmio_wdata;

// Request routing - UART
assign uart_req = mmio_req && mmio_uart_sel;
assign uart_wr = mmio_wr;
assign uart_addr = mmio_addr;
assign uart_wdata = mmio_wdata;

// Request routing - SPI
assign spi_req = mmio_req && mmio_spi_sel;
assign spi_wr = mmio_wr;
assign spi_addr = mmio_addr;
assign spi_wdata = mmio_wdata;

// Request routing - I2C
assign i2c_req = mmio_req && mmio_i2c_sel;
assign i2c_wr = mmio_wr;
assign i2c_addr = mmio_addr;
assign i2c_wdata = mmio_wdata;

// Request routing - Flash
assign flash_req = mmio_req && mmio_flash_sel;
assign flash_wr = mmio_wr;
assign flash_op = mmio_addr[1:0];  // Operation type in lower bits
assign flash_addr = mmio_addr;
assign flash_wdata = mmio_wdata;

// Response multiplexer (consolidated)
always @(*) begin
    case (1'b1)
        mmio_gpio_sel: begin
            mmio_rdata = gpio_rdata;
            mmio_ready = gpio_ready;
            mmio_error = gpio_error;
        end
        mmio_uart_sel: begin
            mmio_rdata = uart_rdata;
            mmio_ready = uart_ready;
            mmio_error = uart_error;
        end
        mmio_spi_sel: begin
            mmio_rdata = spi_rdata;
            mmio_ready = spi_ready;
            mmio_error = |spi_error;  // Convert 4-bit error to single bit
        end
        mmio_i2c_sel: begin
            mmio_rdata = i2c_rdata;
            mmio_ready = i2c_ready;
            mmio_error = |i2c_error;  // Convert 4-bit error to single bit
        end
        mmio_flash_sel: begin
            mmio_rdata = flash_rdata;
            mmio_ready = flash_ready;
            mmio_error = flash_error;
        end
        default: begin
            // Pass through memory controller response for unmapped addresses
            mmio_rdata = mmio_rdata_mc;
            mmio_ready = mmio_ready_mc;
            mmio_error = mmio_error_mc;
        end
    endcase
end

// =============================================================================
// MEMORY CONTROLLER INSTANTIATION
// =============================================================================
memory_controller mem_ctrl (
    .clk(clk_mem),
    .rst_n(rst_mem_n),
    
    // Bus Matrix Interface (instead of direct CPU interface)
    .mem_req(bus_mem_req),
    .mem_wr(bus_mem_wr),
    .mem_size(bus_mem_size),
    .mem_addr(bus_mem_addr),
    .mem_wdata(bus_mem_wdata),
    .mem_rdata(bus_mem_rdata),    .mem_ready(bus_mem_ready),
    .mem_error(bus_mem_error),
    
    // Cache Interface (standardized)
    .cache_req(cache_req),
    .cache_wr(cache_wr),
    .cache_addr(cache_addr),
    .cache_wdata(cache_wdata),
    .cache_rdata(cache_rdata),
    .cache_ready(cache_ready),
    .cache_error(cache_error),
    
    // Physical Memory Interface (standardized)
    .phy_addr(phy_addr),
    .phy_wdata(phy_wdata),
    .phy_wr(phy_wr),
    .phy_req(phy_req),
    .phy_rdata(phy_rdata),
    .phy_ready(phy_ready),
    .phy_error(phy_error),
      // MMIO Interface (connected to consolidated router)
    .mmio_req(mmio_req),
    .mmio_addr(mmio_addr),
    .mmio_wdata(mmio_wdata),
    .mmio_wr(mmio_wr),
    .mmio_rdata(mmio_rdata_mc),
    .mmio_ready(mmio_ready_mc),
    .mmio_error(mmio_error_mc),
    
    // Debug (standardized)    .mc_state(mc_state),
    .access_count(access_count),
    .error_count(error_count)
);

// =============================================================================
// CACHE CONTROLLER INSTANTIATION
// =============================================================================
cache_controller cache_ctrl (
    .clk(clk_cpu),
    .rst_n(rst_cpu_n),
    
    // CPU Interface - Instruction Cache
    .icache_req(icache_req),
    .icache_addr(icache_addr),
    .icache_data(icache_data),
    .icache_hit(icache_hit),
    .icache_ready(icache_ready),
    
    // CPU Interface - Data Cache  
    .dcache_req(dcache_req),
    .dcache_wr(dcache_wr),
    .dcache_size(dcache_size),
    .dcache_addr(dcache_addr),
    .dcache_wdata(dcache_wdata),
    .dcache_rdata(dcache_rdata),
    .dcache_hit(dcache_hit),
    .dcache_ready(dcache_ready),
    
    // Memory Controller Interface - for cache line fills
    .cache_req(cache_req),
    .cache_wr(cache_wr),
    .cache_addr(cache_addr),
    .cache_wdata(cache_wdata),
    .cache_rdata(cache_rdata),
    .cache_ready(cache_ready),
    .cache_error(cache_error),
    
    // Control and Status
    .cache_enable(cache_enable),
    .cache_flush(cache_flush),
    .cache_state(cache_state),
    
    // Debug and Monitoring
    .icache_hits(icache_hits),
    .icache_misses(icache_misses),
    .dcache_hits(dcache_hits),
    .dcache_misses(dcache_misses)
);

// =============================================================================
// PHYSICAL MEMORY INTERFACE (SIMPLIFIED)
// =============================================================================
// External memory control signals
assign mem_addr = phy_addr;
assign mem_data_out = phy_wdata;
assign mem_we_n = ~phy_wr;
assign mem_oe_n = phy_wr;
assign mem_cs_n = ~phy_req;

always @(*) begin
    mem_data_oe = phy_wr && phy_req;
end

// Simple memory ready/error simulation
reg [1:0] mem_delay;
always @(posedge clk_mem or negedge rst_mem_n) begin
    if (!rst_mem_n) begin
        mem_delay <= 2'b00;
    end else if (phy_req) begin
        mem_delay <= mem_delay + 1;
    end else begin
        mem_delay <= 2'b00;
    end
end

assign phy_ready = (mem_delay == 2'b11);
assign phy_error = 1'b0;

// =============================================================================
// SYSTEM CONFIGURATION AND CONTROL (CONSOLIDATED)
// =============================================================================
// Default configuration values (using VTX1 parameters)
// =============================================================================
// SYSTEM CONFIGURATION ASSIGNMENTS
// =============================================================================
// System Enable Signals
assign pll_enable = pwr_good;                          // PLL enabled when power good
assign cache_enable = rst_system_n;                    // Cache enabled after system reset
assign cache_flush = 1'b0;                            // No cache flush requested
assign flash_enable = rst_system_n;                    // Flash enabled after system reset
assign irq_enable = rst_cpu_n;                        // Interrupts enabled after CPU reset

// Clock Configuration
assign cpu_clk_div = 2'b00;                           // CPU clock full speed
assign mem_clk_div = 2'b00;                           // Memory clock full speed
assign periph_clk_div = 2'b01;                        // Peripheral clock /3 divider

// Cache and Bus Configuration
assign cache_policy = 2'b01;                          // Write-back cache policy
assign priority_mode = 2'b00;                         // Fixed priority mode
assign priority_config = 64'h0000000000000000;        // Default priority configuration

// Power Management Configuration
assign low_power_mode = 1'b0;                         // Low power mode disabled
assign sleep_req = 1'b0;                              // No sleep requested
assign boot_mode = 1'b1;                              // Boot from flash

// Interrupt Configuration
assign irq_mask = 32'h00000000;                       // No interrupts masked
assign timer_irq = 1'b0;                              // Timer controller not implemented
assign cache_irq = cache_error;                       // Cache error interrupts
assign mem_irq = bus_mem_error;                       // Memory error interrupts
assign flash_irq = flash_error;                       // Flash error interrupts
assign debug_irq = 1'b0;                              // Debug interrupts not implemented
assign ext_irq = 16'h0000;                            // No external interrupts

// =============================================================================
// STATUS AND DEBUG OUTPUTS (CONSOLIDATED)
// =============================================================================
assign led_power = pwr_good && rst_system_n;
assign led_activity = cpu_bus_req || dma_req || debug_req || gpio_req || uart_req || spi_req || i2c_req || flash_req;
assign led_error = bus_mem_error || cpu_bus_error || dma_error || debug_error || gpio_error || uart_error || |spi_error || |i2c_error || flash_error;

assign pwr_enable = pwr_good;

// UART interface - managed by enhanced UART controller
// uart_tx, uart_rts outputs are from UART controller

// JTAG interface (not implemented)
assign jtag_tdo = 1'b0;

// Test points
assign tp_clk_cpu = clk_cpu;
assign tp_clk_mem = clk_mem;
assign tp_rst_n = rst_system_n;
assign tp_state = {mc_state[1:0], cache_state[1:0]};

// =============================================================================
// FLASH CONTROLLER INTERFACE (TODO: Phase 4 - Placeholder Implementation)
// =============================================================================
assign flash_rdata = {VTX1_WORD_WIDTH{1'b0}};     // Flash read data (always zero)
assign flash_ready = flash_req;                    // Immediate ready response
assign flash_error = 1'b0;                        // No errors reported
assign flash_cs_n = 1'b1;                         // Flash chip select (inactive)
assign flash_sck = 1'b0;                          // Flash serial clock (idle)
assign flash_mosi = 1'b0;                         // Flash data output (idle)
assign flash_wp_n = 1'b1;                         // Flash write protect (disabled)
assign flash_hold_n = 1'b1;                       // Flash hold (disabled)

// =============================================================================
// INTERRUPT CONTROLLER INTERFACE
// =============================================================================
assign irq = 1'b0;                                // No interrupts generated
assign irq_vector = 8'h00;                        // Default interrupt vector
assign irq_level = 2'b00;                         // Interrupt priority level

// =============================================================================
// CPU CORE INSTANTIATION
// =============================================================================
cpu_core cpu_inst (
    .clk(clk_cpu),
    .rst_n(rst_cpu_n),
    .enable(cpu_enable),
    
    // Instruction Memory Interface
    .imem_addr(cpu_imem_addr),
    .imem_data(cpu_imem_data),
    .imem_req(cpu_imem_req),
    .imem_ready(cpu_imem_ready),
    
    // Data Memory Interface  
    .dmem_addr(cpu_dmem_addr),
    .dmem_data(cpu_dmem_data),
    .dmem_we(cpu_dmem_we),
    .dmem_oe(cpu_dmem_oe),
    .dmem_req(cpu_dmem_req),
    .dmem_ready(cpu_dmem_ready),
    
    // Interrupt Interface
    .interrupt_req(cpu_interrupt_req),
    .interrupt_ack(cpu_interrupt_ack),
    .nmi_req(cpu_nmi_req),
    
    // Debug Interface
    .debug_enable(cpu_debug_enable),
    .debug_step(cpu_debug_step),
    .debug_cmd(cpu_debug_cmd),
    .debug_pc(cpu_debug_pc),
    .debug_status(cpu_debug_status),
    
    // Performance Counters
    .cycle_count(cpu_cycle_count),
    .instruction_count(cpu_instruction_count),
    .stall_count(cpu_stall_count),
    
    // Pipeline Status
    .pipeline_stall(cpu_pipeline_stall),
    .pipeline_flush(cpu_pipeline_flush),
    .pipeline_state(cpu_pipeline_state),
    
    // Error handling interface
    .error(cpu_error),
    .error_code(cpu_error_code),
    .timeout(cpu_timeout),
    .error_count(cpu_error_count)
);

// =============================================================================
// CPU ADAPTER INSTANTIATION
// =============================================================================
cpu_adapter cpu_adapter_inst (
    .clk(clk_cpu),
    .rst_n(rst_cpu_n),
    
    // CPU Core Interface
    .cpu_dmem_addr(cpu_dmem_addr),
    .cpu_dmem_data(cpu_dmem_data),
    .cpu_dmem_we(cpu_dmem_we),
    .cpu_dmem_oe(cpu_dmem_oe),
    .cpu_dmem_req(cpu_dmem_req),
    .cpu_dmem_ready(cpu_dmem_ready),
    .cpu_imem_addr(cpu_imem_addr),
    .cpu_imem_req(cpu_imem_req),
    .cpu_imem_ready(cpu_imem_ready),
    
    // Bus Matrix Interface
    .bus_req(cpu_bus_req),
    .bus_wr(cpu_bus_wr),
    .bus_size(cpu_bus_size),
    .bus_addr(cpu_bus_addr),
    .bus_wdata(cpu_bus_wdata),
    .bus_rdata(cpu_bus_rdata),
    .bus_ready(cpu_bus_ready),
    .bus_error(cpu_bus_error),
    .bus_error_code(cpu_bus_error_code),
    .bus_timeout(cpu_bus_timeout),
    .bus_error_clear(cpu_bus_error_clear)
);

// =============================================================================
// DMA CONTROLLER INSTANTIATION
// =============================================================================
dma_controller dma_inst (
    .clk(clk_cpu),
    .rst_n(rst_cpu_n),
    
    // Bus Matrix Interface
    .dma_req(dma_req),
    .dma_wr(dma_wr),
    .dma_size(dma_size),
    .dma_addr(dma_addr),
    .dma_wdata(dma_wdata),
    .dma_rdata(dma_rdata),
    .dma_ready(dma_ready),
    .dma_error(dma_error),
    .dma_error_code(dma_error_code),
    .dma_timeout(dma_timeout),
    .dma_error_clear(dma_error_clear),
    
    // CPU Control Interface
    .dma_enable(dma_enable),
    .dma_start(dma_start),
    .dma_channel(dma_channel),
    .dma_src_addr(dma_src_addr),
    .dma_dest_addr(dma_dest_addr),
    .dma_transfer_count(dma_transfer_count),
    .dma_transfer_mode(dma_transfer_mode),
    
    // Status and Interrupt
    .dma_complete(dma_complete),
    .dma_irq(dma_irq),
    .dma_status(dma_status),
    .dma_transfer_progress(dma_transfer_progress),
    
    // Debug
    .dma_state(dma_state),
    .dma_operations_count(dma_operations_count),
    .dma_error_count(dma_error_count)
);

// =============================================================================
// DEBUG CONTROLLER INSTANTIATION
// =============================================================================
debug_controller debug_inst (
    .clk(clk_debug),
    .rst_n(rst_debug_n),
    
    // Bus Matrix Interface
    .debug_req(debug_req),
    .debug_wr(debug_wr),
    .debug_size(debug_size),
    .debug_addr(debug_addr),
    .debug_wdata(debug_wdata),
    .debug_rdata(debug_rdata),
    .debug_ready(debug_ready),
    .debug_error(debug_error),
    .debug_error_code(debug_error_code),
    .debug_timeout(debug_timeout),
    .debug_error_clear(debug_error_clear),
    
    // JTAG Interface
    .jtag_tck(jtag_tck),
    .jtag_tms(jtag_tms),
    .jtag_tdi(jtag_tdi),
    .jtag_tdo(jtag_tdo),
    .jtag_trst_n(jtag_trst_n),
    
    // Debug Control Interface
    .debug_enable(debug_enable_ctrl),
    .debug_halt_cpu(debug_halt_cpu),
    .debug_step_cpu(debug_step_cpu),
    .debug_reset_cpu(debug_reset_cpu),
    .debug_access_addr(debug_access_addr),
    .debug_access_data(debug_access_data),
    .debug_access_wr(debug_access_wr),
    .debug_access_req(debug_access_req),
    .debug_access_ready(debug_access_ready),
    
    // Debug Status
    .debug_state(debug_state),
    .debug_operations_count(debug_operations_count)
);

// =============================================================================
// BUS MATRIX INSTANTIATION
// =============================================================================
bus_matrix bus_matrix_inst (
    .clk(clk_cpu),
    .rst_n(rst_cpu_n),
    
    // Master 0: CPU Interface
    .cpu_req(cpu_bus_req),
    .cpu_wr(cpu_bus_wr),
    .cpu_size(cpu_bus_size),
    .cpu_addr(cpu_bus_addr),
    .cpu_wdata(cpu_bus_wdata),
    .cpu_rdata(cpu_bus_rdata),
    .cpu_ready(cpu_bus_ready),
    .cpu_error(cpu_bus_error),
    .cpu_error_code(cpu_bus_error_code),
    .cpu_timeout(cpu_bus_timeout),
    .cpu_error_clear(cpu_bus_error_clear),
    
    // Master 1: DMA Interface
    .dma_req(dma_req),
    .dma_wr(dma_wr),
    .dma_size(dma_size),
    .dma_addr(dma_addr),
    .dma_wdata(dma_wdata),
    .dma_rdata(dma_rdata),
    .dma_ready(dma_ready),
    .dma_error(dma_error),
    .dma_error_code(dma_error_code),
    .dma_timeout(dma_timeout),
    .dma_error_clear(dma_error_clear),
    
    // Master 2: Debug Interface
    .debug_req(debug_req),
    .debug_wr(debug_wr),
    .debug_size(debug_size),
    .debug_addr(debug_addr),
    .debug_wdata(debug_wdata),
    .debug_rdata(debug_rdata),
    .debug_ready(debug_ready),
    .debug_error(debug_error),
    .debug_error_code(debug_error_code),
    .debug_timeout(debug_timeout),
    .debug_error_clear(debug_error_clear),
    
    // Slave 0: Memory Controller Interface
    .mem_req(bus_mem_req),
    .mem_wr(bus_mem_wr),
    .mem_size(bus_mem_size),
    .mem_addr(bus_mem_addr),
    .mem_wdata(bus_mem_wdata),
    .mem_rdata(bus_mem_rdata),
    .mem_ready(bus_mem_ready),
    .mem_error(bus_mem_error),
    
    // Slave 1: MMIO Router Interface
    .mmio_req(mmio_req),
    .mmio_wr(mmio_wr),
    .mmio_addr(mmio_addr),
    .mmio_wdata(mmio_wdata),
    .mmio_rdata(mmio_rdata),
    .mmio_ready(mmio_ready),
    .mmio_error(mmio_error),
    
    // Slave 2: Cache Controller Interface (placeholder)
    .cache_req(cache_req),
    .cache_wr(cache_wr),
    .cache_addr(cache_addr),
    .cache_wdata(cache_wdata),
    .cache_rdata(cache_rdata),
    .cache_ready(cache_ready),
    .cache_error(cache_error),
    
    // Configuration and Control
    .arbitration_mode(bus_arbitration_mode),
    .priority_config(bus_priority_config),
    .timeout_config(bus_timeout_config),
    .deadlock_enable(bus_deadlock_enable),
    .performance_enable(bus_performance_enable),
    
    // Status and Debug Outputs
    .matrix_state(bus_matrix_state),
    .current_master(bus_current_master),
    .current_slave(bus_current_slave),
    .arbitration_active(bus_arbitration_active),
    .deadlock_detected(bus_deadlock_detected),
    .deadlock_recovery(bus_deadlock_recovery),
    
    // Performance Monitoring Counters
    .total_transactions(bus_total_transactions),
    .cpu_transactions(bus_cpu_transactions),
    .dma_transactions(bus_dma_transactions),
    .debug_transactions(bus_debug_transactions),
    .timeout_count(bus_timeout_count),
    .error_count(bus_error_count),
    .deadlock_count(bus_deadlock_count),
    
    // Bus Utilization Statistics
    .bus_utilization(bus_utilization),
    .avg_latency(bus_avg_latency),
    .max_latency(bus_max_latency),
    
    // Individual Master Statistics
    .cpu_wait_cycles(bus_cpu_wait_cycles),
    .dma_wait_cycles(bus_dma_wait_cycles),    .debug_wait_cycles(bus_debug_wait_cycles)
);

// =============================================================================
// ENHANCED PERIPHERAL CONTROLLERS INSTANTIATION
// =============================================================================

// Enhanced GPIO Controller with Advanced Features
gpio_controller gpio_inst (
    .clk(clk_periph),
    .rst_n(rst_periph_n),
    
    // CPU Interface
    .gpio_req(gpio_req),
    .gpio_wr(gpio_wr),
    .gpio_addr(gpio_addr),
    .gpio_wdata(gpio_wdata),
    .gpio_rdata(gpio_rdata),
    .gpio_ready(gpio_ready),
    .gpio_error(gpio_error),
    
    // GPIO Interface (24-bit for enhanced controller)
    .gpio_pins(gpio_pins[23:0]),
    
    // Alternate Function Interface
    .alt_func_out(24'h000000),
    .alt_func_in(),
    .alt_func_enable(24'h000000),
    
    // Power Management Interface
    .sleep_mode(1'b0),
    .wake_request(),
    .wake_source(),
    
    // Interrupt Interface
    .gpio_irq(gpio_irq),
    .gpio_irq_vector(),
    
    // Debug and Status
    .gpio_state(),
    .operation_count(),
    .error_count(),
    .debounce_active(),
    .drive_strength_status()
);

// Enhanced UART Controller
uart_controller uart_inst (
    .clk(clk_periph),
    .rst_n(rst_periph_n),
    
    // CPU Interface
    .uart_req(uart_req),
    .uart_wr(uart_wr),
    .uart_addr(uart_addr),
    .uart_wdata(uart_wdata),
    .uart_rdata(uart_rdata),
    .uart_ready(uart_ready),
    .uart_error(uart_error),
    
    // UART Physical Interface
    .uart_rx(uart_rx),
    .uart_tx(uart_tx),
    .uart_rts(uart_rts),
    .uart_cts(uart_cts),
    
    // Interrupt Interface
    .uart_irq(uart_irq),
    .uart_irq_status(uart_irq_status),
      // DMA Interface (adjusted for 8-bit UART data)
    .uart_dma_tx_req(uart_dma_tx_req),
    .uart_dma_rx_req(uart_dma_rx_req),
    .uart_dma_tx_ack(uart_dma_tx_ack),
    .uart_dma_rx_ack(uart_dma_rx_ack),
    .uart_dma_tx_data(uart_dma_tx_data[7:0]),    // Connect only lower 8 bits
    .uart_dma_rx_data(uart_dma_rx_data[7:0])     // Connect only lower 8 bits
);

// Enhanced SPI Controller
spi_controller spi_inst (
    .clk(clk_periph),
    .rst_n(rst_periph_n),
    
    // CPU Interface
    .spi_req(spi_req),
    .spi_wr(spi_wr),
    .spi_addr(spi_addr),
    .spi_wdata(spi_wdata),
    .spi_rdata(spi_rdata),
    .spi_ready(spi_ready),
    .spi_error(spi_error),
    
    // Enhanced SPI Physical Interface
    .spi_sclk(spi_sclk),
    .spi_mosi(spi_mosi),
    .spi_miso(spi_miso),
    .spi_cs_n(spi_cs_n),
    
    // SPI Slave Mode Interface
    .spi_ss_n(spi_ss_n),
    .spi_sclk_in(spi_sclk_in),
    .spi_mosi_in(spi_mosi_in),
    .spi_miso_out(spi_miso_out),
    
    // Enhanced Interrupt System
    .spi_irq(spi_irq),
    .spi_irq_vector(spi_irq_vector),
    
    // DMA Interface
    .dma_tx_req(spi_dma_tx_req),
    .dma_rx_req(spi_dma_rx_req),
    .dma_tx_ack(spi_dma_tx_ack),
    .dma_rx_ack(spi_dma_rx_ack),
    .dma_tx_data(spi_dma_tx_data),
    .dma_rx_data(spi_dma_rx_data),
    
    // Debug and Status
    .spi_state(spi_state),
    .transfer_count(spi_transfer_count),
    .error_count(spi_error_count),
    .fifo_status(spi_fifo_status),
    .debug_info(spi_debug_info)
);

// Enhanced I2C Controller
i2c_controller i2c_inst (
    .clk(clk_periph),
    .rst_n(rst_periph_n),
    
    // CPU Interface
    .i2c_req(i2c_req),
    .i2c_wr(i2c_wr),
    .i2c_addr(i2c_addr),
    .i2c_wdata(i2c_wdata),
    .i2c_rdata(i2c_rdata),
    .i2c_ready(i2c_ready),
    .i2c_error(i2c_error),
    
    // Enhanced I2C Physical Interface
    .i2c_sda(i2c_sda),
    .i2c_scl(i2c_scl),
    
    // Enhanced Interrupt System
    .i2c_irq(i2c_irq),
    .i2c_irq_vector(i2c_irq_vector),
    
    // DMA Interface
    .dma_tx_req(i2c_dma_tx_req),
    .dma_rx_req(i2c_dma_rx_req),
    .dma_tx_ack(i2c_dma_tx_ack),
    .dma_rx_ack(i2c_dma_rx_ack),
    .dma_tx_data(i2c_dma_tx_data),
    .dma_rx_data(i2c_dma_rx_data),
    
    // Debug and Status
    .i2c_state(i2c_state),
    .transfer_count(i2c_transfer_count),
    .error_count(i2c_error_count),
    .bus_busy(i2c_bus_busy),
    .arbitration_lost(i2c_arbitration_lost),
    .fifo_status(i2c_fifo_status),
    .debug_info(i2c_debug_info)
);

// =============================================================================
// DMA PERIPHERAL CONNECTIONS (TODO: Phase 4 - Placeholder Implementation)
// =============================================================================
// DMA Acknowledge signals - connected to peripheral DMA request lines
assign uart_dma_tx_ack = dma_enable && uart_dma_tx_req; // UART TX DMA acknowledge
assign uart_dma_rx_ack = dma_enable && uart_dma_rx_req; // UART RX DMA acknowledge
assign spi_dma_tx_ack = dma_enable && spi_dma_tx_req;   // SPI TX DMA acknowledge
assign spi_dma_rx_ack = dma_enable && spi_dma_rx_req;   // SPI RX DMA acknowledge
assign i2c_dma_tx_ack = dma_enable && i2c_dma_tx_req;   // I2C TX DMA acknowledge
assign i2c_dma_rx_ack = dma_enable && i2c_dma_rx_req;   // I2C RX DMA acknowledge

// DMA Data connections - placeholder zero data (needs proper arbitration in full implementation)
assign uart_dma_tx_data = {VTX1_WORD_WIDTH{1'b0}};     // UART TX DMA data
assign spi_dma_tx_data = {VTX1_WORD_WIDTH{1'b0}};      // SPI TX DMA data
assign i2c_dma_tx_data = {VTX1_WORD_WIDTH{1'b0}};      // I2C TX DMA data

// =============================================================================
// SYSTEM CONFIGURATION AND CONTROL (CONTINUED)
// =============================================================================
// Bus Matrix Configuration
assign bus_arbitration_mode = 2'b01;    // Round-robin arbitration
assign bus_priority_config = 8'h0F;     // Equal priority for all masters
assign bus_timeout_config = 16'h0FFF;   // Reasonable timeout value
assign bus_deadlock_enable = 1'b1;      // Enable deadlock detection
assign bus_performance_enable = 1'b1;   // Enable performance monitoring

// CPU Configuration
assign cpu_enable = rst_cpu_n;
assign cpu_interrupt_req = 16'h0000;
assign cpu_nmi_req = 1'b0;
assign cpu_debug_enable = 1'b0;
assign cpu_debug_step = 1'b0;
assign cpu_debug_cmd = 4'h0;

// =============================================================================
// DMA CONTROLLER CONFIGURATION
// =============================================================================
assign dma_enable = 1'b0;                         // DMA disabled by default
assign dma_start = 1'b0;                          // No transfers initiated
assign dma_channel = 8'h00;                       // Default channel 0
assign dma_src_addr = {VTX1_ADDR_WIDTH{1'b0}};    // Source address (zero)
assign dma_dest_addr = {VTX1_ADDR_WIDTH{1'b0}};   // Destination address (zero)
assign dma_transfer_count = 32'h0;                // Transfer count (zero)
assign dma_transfer_mode = 3'b000;                // Memory-to-memory mode

// =============================================================================
// DEBUG CONTROLLER CONFIGURATION (TODO: Phase 4 - Placeholder Implementation)
// =============================================================================
assign debug_enable_ctrl = 1'b0;                  // Debug controller disabled
assign debug_halt_cpu = 1'b0;                     // Don't halt CPU
assign debug_step_cpu = 1'b0;                     // No single stepping
assign debug_reset_cpu = 1'b0;                    // Don't reset CPU
assign debug_access_addr = {VTX1_ADDR_WIDTH{1'b0}}; // Debug access address (zero)
assign debug_access_data = {VTX1_WORD_WIDTH{1'b0}}; // Debug access data (zero)
assign debug_access_wr = 1'b0;                    // No debug writes
assign debug_access_req = 1'b0;                   // No debug requests

// Instruction memory - connect to memory controller through bus matrix
assign cpu_imem_data = bus_mem_rdata[VTX1_VLIW_WIDTH-1:0];

// Cache controller connections - now handled by cache_controller module
// Note: All cache signals are now connected through the cache_controller instantiation above
// The cache controller manages both instruction and data caches and interfaces with memory controller
assign cache_wr = 1'b0;
// assign cache_addr = {VTX1_ADDR_WIDTH{1'b0}};
// Remove forced constant assignment to cache_addr. The bus_matrix and cache_controller should drive this signal
assign cache_wdata = {VTX1_CACHE_LINE_WIDTH{1'b0}};

// =============================================================================
// DEBUG OUTPUT (USING VTX1 PARAMETERS)
// =============================================================================
// Synthesis directives for Icarus Verilog
// synthesis translate_off
initial begin
    $display("=================================================");
    $display("VTX1 Ternary System-on-Chip CONSOLIDATED Version");
    $display("=================================================");
    $display("Architecture: VTX1 Ternary VLIW Processor");
    $display("Data Width: %d bits", VTX1_WORD_WIDTH);
    $display("Address Width: %d bits", VTX1_ADDR_WIDTH);
    $display("Cache Line: %d bits", VTX1_CACHE_LINE_WIDTH);
    $display("Registers: %d", VTX1_NUM_REGISTERS);
    $display("Pipeline Stages: %d", VTX1_PIPELINE_STAGES);
    $display("VLIW Slots: %d", VTX1_VLIW_SLOTS);
    $display("=================================================");
end

// System monitoring
always @(posedge clk_cpu) begin
    if (rst_cpu_n && cpu_bus_req) begin
        $display("CPU Bus Access: Addr=0x%09X, Data=0x%09X, WR=%b", 
                cpu_bus_addr, cpu_bus_wdata, cpu_bus_wr);
    end
    if (rst_cpu_n && dma_req) begin
        $display("DMA Bus Access: Addr=0x%09X, Data=0x%09X, WR=%b", 
                dma_addr, dma_wdata, dma_wr);
    end
    if (rst_cpu_n && debug_req) begin
        $display("Debug Bus Access: Addr=0x%09X, Data=0x%09X, WR=%b", 
                debug_addr, debug_wdata, debug_wr);
    end
    if (rst_cpu_n && bus_deadlock_detected) begin
        $display("Bus Matrix Deadlock Detected - Recovery Active");
    end
    if (rst_cpu_n && (spi_irq || i2c_irq || uart_irq || gpio_irq)) begin
        $display("Peripheral Interrupt: SPI=%b I2C=%b UART=%b GPIO=%b", 
                spi_irq, i2c_irq, uart_irq, gpio_irq);
    end
end
// synthesis translate_on

endmodule

