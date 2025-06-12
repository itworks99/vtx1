	`timescale 1ns / 1ps
// File: vtx1_error_macros.v
// Description: Standardized error handling macros for VTX1 processor
//
// Purpose: Provides common error codes, error handling macros, and utilities
// to reduce code duplication across VTX1 modules and improve maintainability.

`ifndef VTX1_ERROR_MACROS_V
`define VTX1_ERROR_MACROS_V

// Include common state constants
`include "vtx1_state_constants.v"
//==============================================================================
// Standard Error Codes (4-bit encoding)
//==============================================================================
`define VTX1_ERROR_NONE                 4'h0    // No error
`define VTX1_ERROR_TIMEOUT              4'h1    // Operation timeout
`define VTX1_ERROR_INVALID_OP           4'h2    // Invalid operation
`define VTX1_ERROR_INVALID_ADDR         4'h3    // Invalid address
`define VTX1_ERROR_ACCESS_VIOLATION     4'h4    // Access violation
`define VTX1_ERROR_OVERFLOW             4'h5    // Arithmetic overflow
`define VTX1_ERROR_UNDERFLOW            4'h6    // Arithmetic underflow
`define VTX1_ERROR_DIVIDE_BY_ZERO       4'h7    // Division by zero
`define VTX1_ERROR_PARITY               4'h8    // Parity error
`define VTX1_ERROR_BUS                  4'h9    // Bus error
`define VTX1_ERROR_PROTOCOL             4'hA    // Protocol violation
`define VTX1_ERROR_RESOURCE             4'hB    // Resource unavailable
`define VTX1_ERROR_CONFIG               4'hC    // Configuration error
`define VTX1_ERROR_INTERNAL             4'hD    // Internal error
`define VTX1_ERROR_UNKNOWN              4'hE    // Unknown error
`define VTX1_ERROR_CRITICAL             4'hF    // Critical system error

// Extended error codes - aliases for specific use cases
`define VTX1_ERROR_BUS_FAULT            4'h9    // Bus fault (alias for VTX1_ERROR_BUS)
`define VTX1_ERROR_INVALID_STATE        4'h2    // Invalid state (alias for VTX1_ERROR_INVALID_OP)
`define VTX1_ERROR_TCU_FAULT            4'hC    // TCU configuration error
`define VTX1_ERROR_PIPELINE_STALL       4'h5    // Pipeline stall error (alias for VTX1_ERROR_OVERFLOW)
`define VTX1_ERROR_EXCEPTION            4'hD    // Exception occurred (alias for VTX1_ERROR_INTERNAL)
`define VTX1_ERROR_MEMORY_FAULT         4'h9    // Memory fault (alias for VTX1_ERROR_BUS)
`define VTX1_ERROR_COLLISION            4'hA    // Resource collision (alias for VTX1_ERROR_PROTOCOL)
`define VTX1_ERROR_CLOCK_FAIL           4'hB    // Clock failure (alias for VTX1_ERROR_RESOURCE)
`define VTX1_ERROR_PLL_UNLOCK           4'hB    // PLL unlock (alias for VTX1_ERROR_RESOURCE)
`define VTX1_ERROR_INVALID_PARAM        4'h2    // Invalid parameter (alias for VTX1_ERROR_INVALID_OP)
`define VTX1_ERROR_POWER_FAULT          4'hC    // Power fault (alias for VTX1_ERROR_CONFIG)
`define VTX1_ERROR_THERMAL              4'hC    // Thermal fault (alias for VTX1_ERROR_CONFIG)
`define VTX1_ERROR_RECOVERY             4'hF    // Error recovery mode (alias for VTX1_ERROR_CRITICAL)

//==============================================================================
// Error Code Width Definition  
//==============================================================================
`define VTX1_ERROR_CODE_WIDTH           4       // 4-bit error codes

//==============================================================================
// Error Status Register Bit Positions
//==============================================================================
`define VTX1_ERROR_STATUS_ERROR_BIT     0       // Error occurred flag
`define VTX1_ERROR_STATUS_CODE_LOW      1       // Error code bits [3:0]
`define VTX1_ERROR_STATUS_CODE_HIGH     4       
`define VTX1_ERROR_STATUS_PENDING_BIT   5       // Error pending handling
`define VTX1_ERROR_STATUS_CLEARED_BIT   6       // Error cleared flag
`define VTX1_ERROR_STATUS_FATAL_BIT     7       // Fatal error flag

//==============================================================================
// Error Handling Macros
//==============================================================================

// Set error code and transition to error state
`define VTX1_SET_ERROR(error_reg, state_reg, error_code) \
    begin \
        error_reg <= {1'b1, 3'b000, error_code}; \
        state_reg <= `VTX1_STATE_ERROR; \
    end

// Clear error and return to idle state
`define VTX1_CLEAR_ERROR(error_reg, state_reg) \
    begin \
        error_reg <= 8'h00; \
        state_reg <= `VTX1_STATE_IDLE; \
    end

// Check if error is set
`define VTX1_HAS_ERROR(error_reg) \
    (error_reg[`VTX1_ERROR_STATUS_ERROR_BIT])

// Extract error code from error register
`define VTX1_GET_ERROR_CODE(error_reg) \
    (error_reg[`VTX1_ERROR_STATUS_CODE_HIGH:`VTX1_ERROR_STATUS_CODE_LOW])

// Check for specific error code
`define VTX1_IS_ERROR_CODE(error_reg, error_code) \
    (`VTX1_HAS_ERROR(error_reg) && (`VTX1_GET_ERROR_CODE(error_reg) == error_code))

// Set error pending flag
`define VTX1_SET_ERROR_PENDING(error_reg) \
    begin \
        error_reg[`VTX1_ERROR_STATUS_PENDING_BIT] <= 1'b1; \
    end

// Clear error pending flag
`define VTX1_CLEAR_ERROR_PENDING(error_reg) \
    begin \
        error_reg[`VTX1_ERROR_STATUS_PENDING_BIT] <= 1'b0; \
    end

// Set fatal error flag
`define VTX1_SET_FATAL_ERROR(error_reg, state_reg, error_code) \
    begin \
        error_reg <= {1'b1, 1'b1, 2'b00, error_code}; \
        state_reg <= `VTX1_STATE_ERROR; \
    end

//==============================================================================
// Timeout Counter Macros
//==============================================================================

// Initialize timeout counter
`define VTX1_INIT_TIMEOUT(counter, max_count) \
    begin \
        counter <= max_count; \
    end

// Decrement timeout counter and check for timeout
`define VTX1_CHECK_TIMEOUT(counter, error_reg, state_reg, timeout_occurred) \
    begin \
        if (counter == 0) begin \
            timeout_occurred = 1'b1; \
            `VTX1_SET_ERROR(error_reg, state_reg, `VTX1_ERROR_TIMEOUT) \
        end else begin \
            counter <= counter - 1; \
            timeout_occurred = 1'b0; \
        end \
    end

// Reset timeout counter
`define VTX1_RESET_TIMEOUT(counter, max_count) \
    begin \
        counter <= max_count; \
    end

//==============================================================================
// Bus Error Handling Macros
//==============================================================================

// Handle bus error condition
`define VTX1_HANDLE_BUS_ERROR(bus_error, error_reg, state_reg) \
    begin \
        if (bus_error) begin \
            `VTX1_SET_ERROR(error_reg, state_reg, `VTX1_ERROR_BUS) \
        end \
    end

// Handle access violation
`define VTX1_HANDLE_ACCESS_VIOLATION(violation, error_reg, state_reg) \
    begin \
        if (violation) begin \
            `VTX1_SET_ERROR(error_reg, state_reg, `VTX1_ERROR_ACCESS_VIOLATION) \
        end \
    end

//==============================================================================
// Protocol Error Handling Macros
//==============================================================================

// Handle protocol violation
`define VTX1_HANDLE_PROTOCOL_ERROR(protocol_error, error_reg, state_reg) \
    begin \
        if (protocol_error) begin \
            `VTX1_SET_ERROR(error_reg, state_reg, `VTX1_ERROR_PROTOCOL) \
        end \
    end

// Handle invalid operation
`define VTX1_HANDLE_INVALID_OP(invalid_op, error_reg, state_reg) \
    begin \
        if (invalid_op) begin \
            `VTX1_SET_ERROR(error_reg, state_reg, `VTX1_ERROR_INVALID_OP) \
        end \
    end

//==============================================================================
// Arithmetic Error Handling Macros
//==============================================================================

// Handle arithmetic overflow
`define VTX1_HANDLE_OVERFLOW(overflow, error_reg, state_reg) \
    begin \
        if (overflow) begin \
            `VTX1_SET_ERROR(error_reg, state_reg, `VTX1_ERROR_OVERFLOW) \
        end \
    end

// Handle arithmetic underflow
`define VTX1_HANDLE_UNDERFLOW(underflow, error_reg, state_reg) \
    begin \
        if (underflow) begin \
            `VTX1_SET_ERROR(error_reg, state_reg, `VTX1_ERROR_UNDERFLOW) \
        end \
    end

// Handle division by zero
`define VTX1_HANDLE_DIVIDE_BY_ZERO(div_by_zero, error_reg, state_reg) \
    begin \
        if (div_by_zero) begin \
            `VTX1_SET_ERROR(error_reg, state_reg, `VTX1_ERROR_DIVIDE_BY_ZERO) \
        end \
    end

//==============================================================================
// Error Recovery Macros
//==============================================================================

// Standard error recovery sequence
`define VTX1_ERROR_RECOVERY_MACRO(error_reg, state_reg, recovery_state) \
    begin \
        if (`VTX1_HAS_ERROR(error_reg) && !error_reg[`VTX1_ERROR_STATUS_FATAL_BIT]) begin \
            `VTX1_CLEAR_ERROR(error_reg, state_reg) \
            state_reg <= recovery_state; \
        end \
    end

// Conditional error handling based on enable signal
`define VTX1_CONDITIONAL_ERROR(condition, error_enable, error_reg, state_reg, error_code) \
    begin \
        if (condition && error_enable) begin \
            `VTX1_SET_ERROR(error_reg, state_reg, error_code) \
        end \
    end

//==============================================================================
// Debug and Logging Macros (for simulation)
//==============================================================================

`ifdef VTX1_DEBUG_ERRORS
    // Error logging for simulation
    `define VTX1_LOG_ERROR(module_name, error_code, message) \
        begin \
            $display("[ERROR] %s: Code=0x%h, Message=%s, Time=%0t", \
                     module_name, error_code, message, $time); \
        end

    // Warning logging for simulation
    `define VTX1_LOG_WARNING(module_name, message) \
        begin \
            $display("[WARNING] %s: %s, Time=%0t", \
                     module_name, message, $time); \
        end
`else
    `define VTX1_LOG_ERROR(module_name, error_code, message)
    `define VTX1_LOG_WARNING(module_name, message)
`endif

//==============================================================================
// Error Statistics Macros (for debugging and profiling)
//==============================================================================

// Increment error counter
`define VTX1_INC_ERROR_COUNT(error_count) \
    begin \
        if (error_count < 16'hFFFF) begin \
            error_count <= error_count + 1; \
        end \
    end

// Reset error counter
`define VTX1_RESET_ERROR_COUNT(error_count) \
    begin \
        error_count <= 16'h0000; \
    end

`endif // VTX1_ERROR_MACROS_V

