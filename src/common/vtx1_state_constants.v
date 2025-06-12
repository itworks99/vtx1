	`timescale 1ns / 1ps
// =============================================================================
// VTX1 Common State Machine Constants
// =============================================================================
// Standardized state constants to eliminate duplication across modules
// Compatible with Icarus Verilog - no SystemVerilog packages required
// =============================================================================

`ifndef VTX1_STATE_CONSTANTS_V
`define VTX1_STATE_CONSTANTS_V

// Include base VTX1 interfaces for error codes
`include "vtx1_interfaces.v"
// =============================================================================
// STANDARD STATE MACHINE STATES
// =============================================================================
// These states should be used consistently across all VTX1 modules

// Basic Control States (4-bit encoding)
`define VTX1_STATE_IDLE      4'h0  // Module idle, ready for operation
`define VTX1_STATE_ACTIVE    4'h1  // Module active, processing
`define VTX1_STATE_READ      4'h2  // Read operation state
`define VTX1_STATE_WRITE     4'h3  // Write operation state
`define VTX1_STATE_WAIT      4'h4  // Module waiting for external signal
`define VTX1_STATE_COMPLETE  4'h5  // Operation completed successfully
`define VTX1_STATE_SPECIAL   4'h6  // Special operation mode
`define VTX1_STATE_ERROR     4'hF  // Error state - always 0xF for easy identification

// Extended States for Complex Operations (4-bit encoding)
`define VTX1_STATE_RESET     4'h0  // Reset state (alias for IDLE)
`define VTX1_STATE_INIT      4'h1  // Initialization state
`define VTX1_STATE_SETUP     4'h2  // Setup/configuration state
`define VTX1_STATE_EXECUTE   4'h3  // Execution state
`define VTX1_STATE_VALIDATE  4'h4  // Validation state
`define VTX1_STATE_WRITEBACK 4'h5  // Writeback state
`define VTX1_STATE_RECOVERY  4'h6  // Error recovery state
`define VTX1_STATE_DEBUG     4'h7  // Debug mode state

// =============================================================================
// CPU-SPECIFIC STATES
// =============================================================================
// Pipeline stage states
`define VTX1_CPU_STATE_FETCH    4'h1  // Instruction fetch
`define VTX1_CPU_STATE_DECODE   4'h2  // Instruction decode
`define VTX1_CPU_STATE_EXECUTE  4'h3  // Instruction execute
`define VTX1_CPU_STATE_BRANCH   4'h5  // Branch handling
`define VTX1_CPU_STATE_STALL    4'h6  // Pipeline stall

// Decoder states for VLIW processing
`define VTX1_DECODE_STATE_VALIDATE  4'h1  // Validate instruction format
`define VTX1_DECODE_STATE_A         4'h2  // Decode slot A
`define VTX1_DECODE_STATE_B         4'h3  // Decode slot B  
`define VTX1_DECODE_STATE_C         4'h4  // Decode slot C

// =============================================================================
// MEMORY SUBSYSTEM STATES
// =============================================================================
// Memory controller states
`define VTX1_MC_STATE_IDLE       4'h0  // Memory controller idle
`define VTX1_MC_STATE_CPU_REQ    4'h2  // CPU request processing
`define VTX1_MC_STATE_CACHE_REQ  4'h3  // Cache request processing
`define VTX1_MC_STATE_DMA_REQ    4'h4  // DMA request processing
`define VTX1_MC_STATE_REFRESH    4'h5  // Memory refresh
`define VTX1_MC_STATE_TIMEOUT    4'hE  // Timeout condition

// Memory controller extended states
`define VTX1_MEM_STATE_IDLE       4'h0  // Memory controller idle
`define VTX1_MEM_STATE_READY      4'h1  // Memory ready for requests
`define VTX1_MEM_STATE_CPU_REQ    4'h2  // CPU request processing
`define VTX1_MEM_STATE_CACHE_REQ  4'h3  // Cache request processing
`define VTX1_MEM_STATE_DMA_REQ    4'h4  // DMA request processing
`define VTX1_MEM_STATE_MMIO_REQ   4'h5  // MMIO request processing
`define VTX1_MEM_STATE_REFRESH    4'h6  // Memory refresh
`define VTX1_MEM_STATE_ERROR      4'hF  // Memory error state
`define VTX1_MEM_STATE_TIMEOUT    4'hE  // Memory timeout state

// Cache controller states
`define VTX1_CACHE_STATE_IDLE      4'h0  // Cache idle
`define VTX1_CACHE_STATE_READY     4'h1  // Cache ready for requests
`define VTX1_CACHE_STATE_ACCESS    4'h2  // Cache access in progress
`define VTX1_CACHE_STATE_LOOKUP    4'h3  // Tag lookup
`define VTX1_CACHE_STATE_HIT       4'h4  // Cache hit processing
`define VTX1_CACHE_STATE_MISS      4'h5  // Cache miss processing
`define VTX1_CACHE_STATE_FILL      4'h6  // Cache line fill
`define VTX1_CACHE_STATE_WRITEBACK 4'h7  // Cache writeback
`define VTX1_CACHE_STATE_FLUSH     4'h8  // Cache flush
`define VTX1_CACHE_STATE_ERROR     4'hF  // Cache error state
`define VTX1_CACHE_STATE_TIMEOUT   4'hE  // Cache timeout

// =============================================================================
// PERIPHERAL CONTROLLER STATES
// =============================================================================
// UART controller states
`define VTX1_UART_STATE_IDLE     4'h0  // UART idle
`define VTX1_UART_STATE_TX_START 4'h1  // Start transmission
`define VTX1_UART_STATE_TX_DATA  4'h2  // Transmitting data
`define VTX1_UART_STATE_TX_STOP  4'h3  // Stop bit transmission
`define VTX1_UART_STATE_RX_START 4'h4  // Receive start bit
`define VTX1_UART_STATE_RX_DATA  4'h5  // Receiving data
`define VTX1_UART_STATE_RX_STOP  4'h6  // Receive stop bit

// I2C controller states
`define VTX1_I2C_STATE_IDLE       4'h0  // I2C idle
`define VTX1_I2C_STATE_START      4'h1  // Start condition
`define VTX1_I2C_STATE_ADDR       4'h2  // Address transmission
`define VTX1_I2C_STATE_DATA       4'h3  // Data transmission/reception (deprecated - use specific ones below)
`define VTX1_I2C_STATE_ACK        4'h4  // Acknowledge processing
`define VTX1_I2C_STATE_STOP       4'h5  // Stop condition
`define VTX1_I2C_STATE_SLAVE_IDLE 4'h6  // Slave mode idle
`define VTX1_I2C_STATE_SLAVE_ACTIVE 4'h7 // Slave mode active
`define VTX1_I2C_STATE_SEND_DATA  4'h8  // Data transmission (send)
`define VTX1_I2C_STATE_RECV_DATA  4'h9  // Data reception (receive)

// SPI controller states
`define VTX1_SPI_STATE_IDLE       4'h0  // SPI idle
`define VTX1_SPI_STATE_SETUP      4'h1  // Setup phase
`define VTX1_SPI_STATE_TRANSFER   4'h2  // Data transfer
`define VTX1_SPI_STATE_HOLD       4'h3  // Hold phase
`define VTX1_SPI_STATE_COMPLETE   4'h4  // Transfer complete

// GPIO controller states
`define VTX1_GPIO_STATE_IDLE      4'h0  // GPIO idle
`define VTX1_GPIO_STATE_READ      4'h1  // Pin read operation
`define VTX1_GPIO_STATE_WRITE     4'h2  // Pin write operation
`define VTX1_GPIO_STATE_CONFIG    4'h3  // Pin configuration

// =============================================================================
// HANDSHAKING STATES
// =============================================================================
// TCU handshake states
`define VTX1_HS_STATE_IDLE        3'h0  // Handshake idle
`define VTX1_HS_STATE_REQUEST     3'h1  // Request sent
`define VTX1_HS_STATE_WAIT_ACK    3'h2  // Waiting for acknowledge
`define VTX1_HS_STATE_EXECUTING   3'h3  // Operation executing
`define VTX1_HS_STATE_WAIT_DONE   3'h4  // Waiting for completion
`define VTX1_HS_STATE_COMPLETE    3'h5  // Handshake complete
`define VTX1_HS_STATE_ERROR       3'h7  // Handshake error

// =============================================================================
// BUS ARBITRATION STATES
// =============================================================================
// Bus matrix states
`define VTX1_BUS_STATE_IDLE        4'h0  // Bus idle
`define VTX1_BUS_STATE_ARBITRATE   4'h1  // Arbitration in progress
`define VTX1_BUS_STATE_GRANT       4'h2  // Grant issued
`define VTX1_BUS_STATE_TRANSFER    4'h3  // Data transfer
`define VTX1_BUS_STATE_WAIT        4'h4  // Wait for response
`define VTX1_BUS_STATE_DEADLOCK    4'hD  // Deadlock detected
`define VTX1_BUS_STATE_ERROR       4'hF  // Bus error

// Extended bus matrix states
`define VTX1_BUS_STATE_CPU_ACTIVE     4'h5  // CPU actively using bus
`define VTX1_BUS_STATE_DMA_ACTIVE     4'h6  // DMA actively using bus
`define VTX1_BUS_STATE_DEBUG_ACTIVE   4'h7  // Debug actively using bus

// =============================================================================
// CLOCK AND POWER MANAGEMENT STATES
// =============================================================================
// Clock manager states
`define VTX1_CM_STATE_RESET       4'h0  // Clock reset
`define VTX1_CM_STATE_OSC_START   4'h1  // Oscillator startup
`define VTX1_CM_STATE_PLL_LOCK    4'h2  // PLL lock sequence
`define VTX1_CM_STATE_STABILIZE   4'h3  // Clock stabilization
`define VTX1_CM_STATE_NORMAL      4'h4  // Normal operation
`define VTX1_CM_STATE_LOW_POWER   4'h5  // Low power mode
`define VTX1_CM_STATE_SLEEP       4'h6  // Sleep mode
`define VTX1_CM_STATE_SHUTDOWN    4'h7  // Shutdown mode
`define VTX1_CM_STATE_ERROR       4'hF  // Clock error

// =============================================================================
// DMA CONTROLLER STATES
// =============================================================================
// DMA controller states
`define VTX1_DMA_STATE_IDLE       4'h0  // DMA idle
`define VTX1_DMA_STATE_SETUP      4'h1  // DMA setup
`define VTX1_DMA_STATE_READ       4'h2  // DMA read operation
`define VTX1_DMA_STATE_WRITE      4'h3  // DMA write operation
`define VTX1_DMA_STATE_COMPLETE   4'h4  // DMA transfer complete
`define VTX1_DMA_STATE_ERROR      4'hF  // DMA error state

// =============================================================================
// DEBUG CONTROLLER STATES
// =============================================================================
// Debug controller states
`define VTX1_DEBUG_STATE_IDLE         4'h0  // Debug idle
`define VTX1_DEBUG_STATE_WAIT_ENABLE  4'h1  // Waiting for debug enable
`define VTX1_DEBUG_STATE_ACCESS       4'h2  // Debug access in progress
`define VTX1_DEBUG_STATE_READ         4'h3  // Debug read operation
`define VTX1_DEBUG_STATE_WRITE        4'h4  // Debug write operation
`define VTX1_DEBUG_STATE_BREAKPOINT   4'h5  // Breakpoint hit
`define VTX1_DEBUG_STATE_WATCHPOINT   4'h6  // Watchpoint hit
`define VTX1_DEBUG_STATE_HALT         4'h7  // CPU halt state
`define VTX1_DEBUG_STATE_ERROR        4'hF  // Debug error state

// =============================================================================
// STATE VALIDATION MACROS
// =============================================================================
// Macro to check if a state is valid
`define VTX1_VALID_STATE(state) \
    ((state) <= 4'hF)

// Macro to check if a state is an error state
`define VTX1_IS_ERROR_STATE(state) \
    ((state) == `VTX1_STATE_ERROR || (state) == 4'hF)

// Macro to check if a state is idle
`define VTX1_IS_IDLE_STATE(state) \
    ((state) == `VTX1_STATE_IDLE || (state) == 4'h0)

// =============================================================================
// STATE TRANSITION HELPERS
// =============================================================================
// Standard state transition patterns as localparam definitions
// These can be included in modules to ensure consistent behavior

// Standard idle-to-active transition check
`define VTX1_IDLE_TO_ACTIVE_CONDITION(enable, rst_n) \
    ((enable) && (rst_n))

// Standard error-to-recovery transition check  
`define VTX1_ERROR_TO_RECOVERY_CONDITION(error_clear, timeout_expired) \
    ((error_clear) || (timeout_expired))

// Standard completion check
`define VTX1_OPERATION_COMPLETE_CONDITION(valid, ready, done) \
    ((valid) && (ready) && (done))

`endif // VTX1_STATE_CONSTANTS_V

