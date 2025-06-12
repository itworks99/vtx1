	`timescale 1ns / 1ps
// =============================================================================
// VTX1 Standardized Interface Definitions
// =============================================================================
// Standardized interfaces to eliminate code duplication
// Compatible with Icarus Verilog (no SystemVerilog packages)
// =============================================================================

`ifndef VTX1_INTERFACES_V
`define VTX1_INTERFACES_V

// VTX1 System Parameters
`define VTX1_WORD_WIDTH 36        // 36-bit ternary word
`define VTX1_DWORD_WIDTH 72       // 72-bit double word
`define VTX1_CACHE_LINE_WIDTH 288 // 288-bit cache line (8 words)
`define VTX1_ADDR_WIDTH 36        // 36-bit addressing
`define VTX1_REG_ADDR_WIDTH 4     // 4-bit register addressing (supports 13 registers)
`define VTX1_WORD_DEFAULT 36'b010101010101010101010101010101010101 // All zero trits

// MMIO Address Map
`define VTX1_GPIO_BASE_ADDR  16'h1000
`define VTX1_UART_BASE_ADDR  16'h1001
`define VTX1_SPI_BASE_ADDR   16'h1002
`define VTX1_I2C_BASE_ADDR   16'h1003
`define VTX1_TIMER_BASE_ADDR 16'h1004
`define VTX1_FLASH_BASE_ADDR 16'h1005

// System Configuration
`define VTX1_NUM_REGISTERS 13
`define VTX1_PIPELINE_STAGES 4
`define VTX1_VLIW_SLOTS 3
`define VTX1_VLIW_WIDTH 108       // 108-bit VLIW instruction (3 x 36-bit slots)

// Error Handling Constants
`define VTX1_TIMEOUT_CYCLES 1000     // Maximum cycles to wait for operation
// Note: Error codes are now defined in vtx1_error_macros.v

// Standard Memory Interface Bundle (36-bit ternary)
// Enhanced with comprehensive error handling
`define VTX1_MEMORY_INTERFACE_PORTS \
    input wire req,             /* Request active */ \
    input wire wr,              /* Write enable (1=write, 0=read) */ \
    input wire [1:0] size,      /* Transfer size (00=byte, 01=word, 10=dword) */ \
    input wire [35:0] addr,     /* 36-bit ternary address */ \
    input wire [35:0] wdata,    /* Write data */ \
    output wire [35:0] rdata,   /* Read data */ \
    output wire ready,          /* Operation complete */ \
    output wire error,          /* Error occurred */ \
    output wire [3:0] error_code, /* Specific error code */ \
    output wire timeout,        /* Operation timeout */ \
    input wire error_clear      /* Clear error state */

// Cache Interface Bundle (288-bit cache lines)
// Enhanced with comprehensive error handling
`define VTX1_CACHE_INTERFACE_PORTS \
    input wire req,             /* Cache request */ \
    input wire wr,              /* Write enable */ \
    input wire [35:0] addr,     /* Address */ \
    input wire [287:0] wdata,   /* Cache line write data */ \
    output wire [287:0] rdata,  /* Cache line read data */ \
    output wire hit,            /* Cache hit */ \
    output wire ready,          /* Operation ready */ \
    output wire error,          /* Error occurred */ \
    output wire [3:0] error_code, /* Specific error code */ \
    output wire timeout,        /* Cache operation timeout */ \
    input wire error_clear,     /* Clear error state */ \
    output wire [31:0] hit_count, /* Cache hit counter */ \
    output wire [31:0] miss_count /* Cache miss counter */

// MMIO Interface Bundle
`define VTX1_MMIO_INTERFACE_PORTS \
    input wire req,           /* MMIO request */ \
    input wire wr,            /* Write enable */ \
    input wire [35:0] addr,   /* MMIO address */ \
    input wire [35:0] wdata,  /* Write data */ \
    output wire [35:0] rdata, /* Read data */ \
    output wire ready,        /* Operation ready */ \
    output wire error         /* Error occurred */

// Clock and Reset Bundle
`define VTX1_CLOCK_RESET_PORTS \
    input wire clk,     /* Domain clock */ \
    input wire rst_n    /* Active-low reset */

// Interrupt Interface Bundle
`define VTX1_INTERRUPT_INTERFACE_PORTS \
    output wire irq,                    /* Interrupt request */ \
    output wire [7:0] irq_vector,       /* Interrupt vector */ \
    output wire [1:0] irq_level,        /* Interrupt level */ \
    input wire irq_ack,                 /* Interrupt acknowledge */ \
    input wire [7:0] irq_ack_vector     /* Acknowledged vector */

// Debug Interface Bundle
`define VTX1_DEBUG_INTERFACE_PORTS \
    output wire [3:0] state,        /* Module state */ \
    output wire [31:0] count,       /* Operation count */ \
    output wire [31:0] error_count  /* Error count */

`endif // VTX1_INTERFACES_V

