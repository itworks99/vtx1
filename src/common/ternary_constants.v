	`timescale 1ns / 1ps
// VTX1 Ternary Constants and Encoding Definitions
// Part of the VTX1 Ternary System-on-Chip

`ifndef TERNARY_CONSTANTS_V
`define TERNARY_CONSTANTS_V

// ============================================================================
// TERNARY ENCODING CONSTANTS
// ============================================================================
// Authoritative ternary encoding (2 bits per trit):
// TRIT_NEG  = 2'b00 = -1 (0.0V ± 0.5V noise margin)
// TRIT_ZERO = 2'b01 =  0 (2.5V ± 0.5V noise margin) 
// TRIT_POS  = 2'b10 = +1 (5.0V ± 0.5V noise margin)
// TRIT_UNDEF= 2'b11 = undefined/invalid state

// Basic trit values
`define TRIT_NEG    2'b00    // -1 value
`define TRIT_ZERO   2'b01    //  0 value
`define TRIT_POS    2'b10    // +1 value
`define TRIT_UNDEF  2'b11    // Undefined/invalid

// Voltage level constants (informational, for documentation)
`define VOLTAGE_NEG   0.0    // -1 corresponds to 0.0V
`define VOLTAGE_ZERO  2.5    //  0 corresponds to 2.5V
`define VOLTAGE_POS   5.0    // +1 corresponds to 5.0V
`define NOISE_MARGIN  0.5    // ±0.5V noise margin

// ============================================================================
// TERNARY WORD SIZES
// ============================================================================
// Standard ternary word sizes used throughout VTX1

`define TRIT_WIDTH      2    // Bits per trit (2-bit encoding)
`define TRYTE_TRITS     9    // Trits per tryte (9 trits)
`define TRYTE_WIDTH    18    // Bits per tryte (9 * 2 = 18 bits)
`define WORD_TRITS     18    // Trits per word (18 trits)
`define WORD_WIDTH     36    // Bits per word (18 * 2 = 36 bits)
`define DWORD_TRITS    36    // Trits per double word (36 trits)
`define DWORD_WIDTH    72    // Bits per double word (36 * 2 = 72 bits)

// ============================================================================
// VTX1 ARCHITECTURE CONSTANTS
// ============================================================================

// Register file configuration
`define NUM_REGISTERS  13    // Total number of registers
`define REG_ADDR_WIDTH  4    // Address width for 13 registers (need 4 bits)

// VTX1-specific register constants
`define VTX1_REG_ADDR_WIDTH  4    // VTX1 register address width (same as REG_ADDR_WIDTH)

// Pipeline stages
`define PIPELINE_STAGES 4    // 4-stage pipeline: Fetch, Decode, Execute, Writeback

// VLIW instruction format
`define VLIW_WIDTH     108   // 108-bit VLIW instruction
`define VLIW_SLOTS      3    // 3 parallel execution slots
`define SLOT_WIDTH     36    // 36 bits per slot
`define VLIW_BYTES     14    // VLIW instruction size in bytes (108/8, rounded up)

// Data width constants
`define DATA_WIDTH     36    // Standard data bus width
`define DATA_BYTES      4    // Data width in bytes (36/8, rounded up)

// VTX1-specific width constants (used throughout the system)
`define VTX1_WORD_WIDTH   36    // VTX1 word width (same as DATA_WIDTH)
`define VTX1_ADDR_WIDTH   36    // VTX1 address width
`define VTX1_VLIW_WIDTH   108   // VTX1 VLIW instruction width (3x36-bit slots)

// Memory configuration
`define FLASH_SIZE_KB  432   // 432KB Flash ROM
`define RAM_SIZE_KB    144   // 144KB RAM
`define CACHE_SIZE_KB    8   // 8KB L1 instruction cache

// Address space
`define ADDR_WIDTH     36    // 36-bit addressing (18 trits)
`define ADDR_TRITS     18    // Address width in trits

// Clock domain frequencies (for reference)
`define CORE_FREQ_MHZ  100   // Core clock: 100MHz
`define MEM_FREQ_MHZ    50   // Memory clock: 50MHz
`define SYS_FREQ_MHZ    25   // System clock: 25MHz
`define DBG_FREQ_MHZ    10   // Debug clock: 10MHz

// ============================================================================
// TIMEOUT AND ERROR HANDLING CONSTANTS
// ============================================================================

`define VTX1_TIMEOUT_CYCLES    1000    // Default timeout in clock cycles
`define VTX1_RETRY_COUNT       3       // Default retry count for operations
`define VTX1_ERROR_LOG_DEPTH   16      // Error log depth (entries)

// Enhanced timeout constants for different operation types
`define VTX1_TIMEOUT_SIMPLE    500     // Simple operations (ADD, SUB, etc.)
`define VTX1_TIMEOUT_COMPLEX   2000    // Complex arithmetic (DIV, MOD, SQRT)
`define VTX1_TIMEOUT_TRANSCENDENTAL 3000 // Transcendental functions (SIN, COS, etc.)
`define VTX1_TIMEOUT_VECTOR    1200    // Vector operations
`define VTX1_TIMEOUT_MEMORY    800     // Memory operations
`define VTX1_TIMEOUT_SYSTEM    1500    // System operations

// Handshaking timeout constants
`define VTX1_HANDSHAKE_TIMEOUT 50      // TCU handshake timeout cycles
`define VTX1_RETRY_DELAY       10      // Delay between retries

// Error codes (4-bit values)
`define ERR_NONE               4'b0000  // No error
`define ERR_TIMEOUT            4'b0001  // Operation timeout
`define ERR_INVALID_TRIT       4'b0010  // Invalid trit encoding detected
`define ERR_OVERFLOW           4'b0011  // Arithmetic overflow
`define ERR_UNDERFLOW          4'b0100  // Arithmetic underflow
`define ERR_DIVISION_BY_ZERO   4'b0101  // Division by zero
`define ERR_INVALID_OPCODE     4'b0110  // Invalid instruction opcode
`define ERR_MEMORY_FAULT       4'b0111  // Memory access fault
`define ERR_BUS_ERROR          4'b1000  // Bus transaction error
`define ERR_CACHE_MISS         4'b1001  // Cache miss (for performance)
`define ERR_HANDSHAKE_FAIL     4'b1010  // TCU handshaking failure
`define ERR_TERNARY_ENCODING   4'b1011  // Ternary encoding violation
`define ERR_INTERRUPT_PENDING  4'b1010  // Interrupt pending
`define ERR_RESERVED_1         4'b1011  // Reserved for future use
`define ERR_RESERVED_2         4'b1100  // Reserved for future use
`define ERR_RESERVED_3         4'b1101  // Reserved for future use
`define ERR_RESERVED_4         4'b1110  // Reserved for future use
`define ERR_CRITICAL           4'b1111  // Critical system error

// ============================================================================
// INSTRUCTION SET CONSTANTS
// ============================================================================

`define TOTAL_INSTRUCTIONS 78   // Total instruction count
`define NATIVE_INSTRUCTIONS 52  // Native instructions
`define MICRO_INSTRUCTIONS  26  // Microcode instructions

// Instruction format fields
`define OPCODE_WIDTH    12   // Opcode field width (6 trits * 2 bits)
`define OPCODE_TRITS     6   // Opcode field width in trits

// ============================================================================
// SYSTEM CONFIGURATION CONSTANTS
// ============================================================================

// Interrupt system
`define NUM_INTERRUPT_SOURCES  32     // Total interrupt sources
`define INTERRUPT_ADDR_WIDTH   5      // Address width for 32 sources
`define INTERRUPT_PRIORITY_LEVELS 8   // Number of priority levels

// GPIO configuration
`define GPIO_PIN_COUNT        24      // Total GPIO pins
`define GPIO_PORT_WIDTH       8       // 8 pins per port (3 ports)
`define GPIO_NUM_PORTS        3       // Number of GPIO ports

// Communication interfaces
`define UART_CHANNELS         2       // Number of UART channels
`define SPI_CHANNELS          2       // Number of SPI channels  
`define I2C_CHANNELS          1       // Number of I2C channels

// Power management
`define POWER_DOMAINS         4       // Number of power domains
`define SLEEP_MODE_LEVELS     4       // Number of sleep levels

// ============================================================================
// UTILITY MACROS
// ============================================================================

// Convert trit count to bit count
`define TRITS_TO_BITS(trits) ((trits) * `TRIT_WIDTH)

// Convert bit count to trit count (must be even)
`define BITS_TO_TRITS(bits)  ((bits) / `TRIT_WIDTH)

// Check if a 2-bit value represents a valid trit
`define IS_VALID_TRIT(trit) ((trit) != `TRIT_UNDEF)

// Extract individual trit from a multi-trit word
`define GET_TRIT(word, trit_index) (word)[((trit_index) * `TRIT_WIDTH + 1) : ((trit_index) * `TRIT_WIDTH)]

// Set individual trit in a multi-trit word
`define SET_TRIT(word, trit_index, value) \
    word[((trit_index) * `TRIT_WIDTH + 1) : ((trit_index) * `TRIT_WIDTH)] = (value)

// Check if entire word contains only valid trits
`define WORD_IS_VALID(word) (\
    `IS_VALID_TRIT(word[1:0]) && `IS_VALID_TRIT(word[3:2]) && \
    `IS_VALID_TRIT(word[5:4]) && `IS_VALID_TRIT(word[7:6]) && \
    `IS_VALID_TRIT(word[9:8]) && `IS_VALID_TRIT(word[11:10]) && \
    `IS_VALID_TRIT(word[13:12]) && `IS_VALID_TRIT(word[15:14]) && \
    `IS_VALID_TRIT(word[17:16]) && `IS_VALID_TRIT(word[19:18]) && \
    `IS_VALID_TRIT(word[21:20]) && `IS_VALID_TRIT(word[23:22]) && \
    `IS_VALID_TRIT(word[25:24]) && `IS_VALID_TRIT(word[27:26]) && \
    `IS_VALID_TRIT(word[29:28]) && `IS_VALID_TRIT(word[31:30]) && \
    `IS_VALID_TRIT(word[33:32]) && `IS_VALID_TRIT(word[35:34]))

// Generate all-zero word/tryte
`define ZERO_WORD   {18{`TRIT_ZERO}}    // 36-bit all-zero word
`define ZERO_TRYTE  {9{`TRIT_ZERO}}     // 18-bit all-zero tryte

// Generate all-positive/all-negative words
`define MAX_POS_WORD   {18{`TRIT_POS}}  // Maximum positive value
`define MAX_NEG_WORD   {18{`TRIT_NEG}}  // Maximum negative value

// Common bit patterns for testing and debugging
`define TEST_PATTERN_01  36'h155555555  // Alternating 01 pattern (all zeros)
`define TEST_PATTERN_10  36'h2aaaaaaaa  // Alternating 10 pattern (all ones)
`define TEST_PATTERN_MIX 36'h1a5a5a5a5  // Mixed pattern for testing

// ============================================================================
// RESET AND INITIALIZATION VALUES
// ============================================================================

// Default values for registers and memory
`define DEFAULT_TRIT    `TRIT_ZERO      // Default trit value (0)
`define DEFAULT_TRYTE   18'b010101010101010101  // All zero tryte (9 trits of 01)
`define DEFAULT_WORD    36'b010101010101010101010101010101010101  // All zero word (18 trits of 01)

// Reset values
`define RESET_PC        `DEFAULT_WORD   // Program counter reset value
`define RESET_SP        `DEFAULT_WORD   // Stack pointer reset value

`endif // TERNARY_CONSTANTS_V

