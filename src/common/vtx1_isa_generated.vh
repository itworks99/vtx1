// VTX1 trit-native ISA constants
// GENERATED FILE - DO NOT EDIT.
// Source of truth: src/assembler/isa/isa.go
// Regenerate with:  cd src/assembler && go run ./isa/gen
//
// Instruction word: 18 trits = 36 bits (2 bits/trit: 00=-1, 01=0, 10=+1).
// Trit t (0 = most significant) occupies bits [35-2t : 34-2t].
//
// Field layout (trit positions, MSB first):
//   all:     opcode  trits 0..3   (bits [35:28])
//   R3:      rd 4..6 | rs1 7..9 | mode 10 | payload 11..17
//            (mode 0: rs2 = trits 15..17; mode +1: 7-trit immediate)
//   R2:      rd 4..6 | rs 7..9
//   R1:      r  4..6
//   Mem:     reg 4..6 | mode 7 | payload 8..17
//            (mode 0: 10-trit address; mode +1: base 8..10, index 11..13;
//             mode -1: base 8..10, 7-trit offset 11..17)
//   Branch:  rs1 4..6 | rs2 7..9 | 8-trit PC-relative word offset 10..17
//            (rs2 = TZ means compare against zero)
//   Jump:    14-trit absolute word address, trits 4..17
//
// A VLIW bundle is 3 instruction words = 54 trits = 108 bits (`VLIW_WIDTH).
// The all-zero-trit word is NOP.

`ifndef VTX1_ISA_GENERATED_VH
`define VTX1_ISA_GENERATED_VH

// Geometry
`define VTX1_SLOT_TRITS   18
`define VTX1_SLOT_BITS    36
`define VTX1_OPCODE_BITS  8
`define VTX1_REG_BITS     6

// Opcodes (4 trits as 8 bits; comment shows balanced-trit pattern and value)
`define VTX1_OPC_VADD     8'b00_00_00_00  // ---- (value 0, vector)
`define VTX1_OPC_ADD      8'b00_00_00_01  // ---0 (value 1, ALU)
`define VTX1_OPC_SUB      8'b00_00_00_10  // ---+ (value 2, ALU)
`define VTX1_OPC_MUL      8'b00_00_01_00  // --0- (value 3, ALU)
`define VTX1_OPC_AND      8'b00_00_01_01  // --00 (value 4, ALU)
`define VTX1_OPC_OR       8'b00_00_01_10  // --0+ (value 5, ALU)
`define VTX1_OPC_NOT      8'b00_00_10_00  // --+- (value 6, ALU)
`define VTX1_OPC_XOR      8'b00_00_10_01  // --+0 (value 7, ALU)
`define VTX1_OPC_SHL      8'b00_00_10_10  // --++ (value 8, ALU)
`define VTX1_OPC_SHR      8'b00_01_00_00  // -0-- (value 9, ALU)
`define VTX1_OPC_ROL      8'b00_01_00_01  // -0-0 (value 10, ALU)
`define VTX1_OPC_ROR      8'b00_01_00_10  // -0-+ (value 11, ALU)
`define VTX1_OPC_CMP      8'b00_01_01_00  // -00- (value 12, ALU)
`define VTX1_OPC_TEST     8'b00_01_01_01  // -000 (value 13, ALU)
`define VTX1_OPC_INC      8'b00_01_01_10  // -00+ (value 14, ALU)
`define VTX1_OPC_DEC      8'b00_01_10_00  // -0+- (value 15, ALU)
`define VTX1_OPC_NEG      8'b00_01_10_01  // -0+0 (value 16, ALU)
`define VTX1_OPC_LD       8'b00_01_10_10  // -0++ (value 17, memory)
`define VTX1_OPC_ST       8'b00_10_00_00  // -+-- (value 18, memory)
`define VTX1_OPC_VLD      8'b00_10_00_01  // -+-0 (value 19, memory)
`define VTX1_OPC_VST      8'b00_10_00_10  // -+-+ (value 20, memory)
`define VTX1_OPC_FLD      8'b00_10_01_00  // -+0- (value 21, memory)
`define VTX1_OPC_FST      8'b00_10_01_01  // -+00 (value 22, memory)
`define VTX1_OPC_LEA      8'b00_10_01_10  // -+0+ (value 23, memory)
`define VTX1_OPC_PUSH     8'b00_10_10_00  // -++- (value 24, memory)
`define VTX1_OPC_POP      8'b00_10_10_01  // -++0 (value 25, memory)
`define VTX1_OPC_JMP      8'b00_10_10_10  // -+++ (value 26, control)
`define VTX1_OPC_JAL      8'b01_00_00_00  // 0--- (value 27, control)
`define VTX1_OPC_JR       8'b01_00_00_01  // 0--0 (value 28, control)
`define VTX1_OPC_JALR     8'b01_00_00_10  // 0--+ (value 29, control)
`define VTX1_OPC_BEQ      8'b01_00_01_00  // 0-0- (value 30, control)
`define VTX1_OPC_BNE      8'b01_00_01_01  // 0-00 (value 31, control)
`define VTX1_OPC_BLT      8'b01_00_01_10  // 0-0+ (value 32, control)
`define VTX1_OPC_BGE      8'b01_00_10_00  // 0-+- (value 33, control)
`define VTX1_OPC_BLTU     8'b01_00_10_01  // 0-+0 (value 34, control)
`define VTX1_OPC_BGEU     8'b01_00_10_10  // 0-++ (value 35, control)
`define VTX1_OPC_BGT      8'b01_01_00_00  // 00-- (value 36, control)
`define VTX1_OPC_BLE      8'b01_01_00_01  // 00-0 (value 37, control)
`define VTX1_OPC_CALL     8'b01_01_00_10  // 00-+ (value 38, control)
`define VTX1_OPC_RET      8'b01_01_01_00  // 000- (value 39, control)
`define VTX1_OPC_NOP      8'b01_01_01_01  // 0000 (value 40, ALU)
`define VTX1_OPC_VSUB     8'b01_01_01_10  // 000+ (value 41, vector)
`define VTX1_OPC_VMUL     8'b01_01_10_00  // 00+- (value 42, vector)
`define VTX1_OPC_VAND     8'b01_01_10_01  // 00+0 (value 43, vector)
`define VTX1_OPC_VOR      8'b01_01_10_10  // 00++ (value 44, vector)
`define VTX1_OPC_VNOT     8'b01_10_00_00  // 0+-- (value 45, vector)
`define VTX1_OPC_VSHL     8'b01_10_00_01  // 0+-0 (value 46, vector)
`define VTX1_OPC_VSHR     8'b01_10_00_10  // 0+-+ (value 47, vector)
`define VTX1_OPC_FADD     8'b01_10_01_00  // 0+0- (value 48, floating point)
`define VTX1_OPC_FSUB     8'b01_10_01_01  // 0+00 (value 49, floating point)
`define VTX1_OPC_FMUL     8'b01_10_01_10  // 0+0+ (value 50, floating point)
`define VTX1_OPC_FCMP     8'b01_10_10_00  // 0++- (value 51, floating point)
`define VTX1_OPC_FMOV     8'b01_10_10_01  // 0++0 (value 52, floating point)
`define VTX1_OPC_FNEG     8'b01_10_10_10  // 0+++ (value 53, floating point)
`define VTX1_OPC_WFI      8'b10_00_00_00  // +--- (value 54, system)
`define VTX1_OPC_DIV      8'b10_00_00_01  // +--0 (value 55, ALU, microcoded)
`define VTX1_OPC_MOD      8'b10_00_00_10  // +--+ (value 56, ALU, microcoded)
`define VTX1_OPC_UDIV     8'b10_00_01_00  // +-0- (value 57, ALU, microcoded)
`define VTX1_OPC_UMOD     8'b10_00_01_01  // +-00 (value 58, ALU, microcoded)
`define VTX1_OPC_SQRT     8'b10_00_01_10  // +-0+ (value 59, ALU, microcoded)
`define VTX1_OPC_ABS      8'b10_00_10_00  // +-+- (value 60, ALU, microcoded)
`define VTX1_OPC_SIN      8'b10_00_10_01  // +-+0 (value 61, ALU, microcoded)
`define VTX1_OPC_COS      8'b10_00_10_10  // +-++ (value 62, ALU, microcoded)
`define VTX1_OPC_TAN      8'b10_01_00_00  // +0-- (value 63, ALU, microcoded)
`define VTX1_OPC_ASIN     8'b10_01_00_01  // +0-0 (value 64, ALU, microcoded)
`define VTX1_OPC_ACOS     8'b10_01_00_10  // +0-+ (value 65, ALU, microcoded)
`define VTX1_OPC_ATAN     8'b10_01_01_00  // +00- (value 66, ALU, microcoded)
`define VTX1_OPC_EXP      8'b10_01_01_01  // +000 (value 67, ALU, microcoded)
`define VTX1_OPC_LOG      8'b10_01_01_10  // +00+ (value 68, ALU, microcoded)
`define VTX1_OPC_VDOT     8'b10_01_10_00  // +0+- (value 69, vector, microcoded)
`define VTX1_OPC_VREDUCE  8'b10_01_10_01  // +0+0 (value 70, vector, microcoded)
`define VTX1_OPC_VMAX     8'b10_01_10_10  // +0++ (value 71, vector, microcoded)
`define VTX1_OPC_VMIN     8'b10_10_00_00  // ++-- (value 72, vector, microcoded)
`define VTX1_OPC_VSUM     8'b10_10_00_01  // ++-0 (value 73, vector, microcoded)
`define VTX1_OPC_VPERM    8'b10_10_00_10  // ++-+ (value 74, vector, microcoded)
`define VTX1_OPC_CACHE    8'b10_10_01_00  // ++0- (value 75, memory, microcoded)
`define VTX1_OPC_FLUSH    8'b10_10_01_01  // ++00 (value 76, memory, microcoded)
`define VTX1_OPC_MEMBAR   8'b10_10_01_10  // ++0+ (value 77, memory, microcoded)
`define VTX1_OPC_SYSCALL  8'b10_10_10_00  // +++- (value 78, system, microcoded)
`define VTX1_OPC_BREAK    8'b10_10_10_01  // +++0 (value 79, system, microcoded)
`define VTX1_OPC_HALT     8'b10_10_10_10  // ++++ (value 80, system, microcoded)

// Registers (3 trits as 6 bits)
`define VTX1_REG_T0  6'b00_00_00  // --- (index 0)
`define VTX1_REG_T1  6'b00_00_01  // --0 (index 1)
`define VTX1_REG_T2  6'b00_00_10  // --+ (index 2)
`define VTX1_REG_T3  6'b00_01_00  // -0- (index 3)
`define VTX1_REG_T4  6'b00_01_01  // -00 (index 4)
`define VTX1_REG_T5  6'b00_01_10  // -0+ (index 5)
`define VTX1_REG_T6  6'b00_10_00  // -+- (index 6)
`define VTX1_REG_TB  6'b00_10_01  // -+0 (index 7)
`define VTX1_REG_TA  6'b00_10_10  // -++ (index 8)
`define VTX1_REG_TC  6'b01_00_00  // 0-- (index 9)
`define VTX1_REG_TS  6'b01_00_01  // 0-0 (index 10)
`define VTX1_REG_TI  6'b01_00_10  // 0-+ (index 11)
`define VTX1_REG_VA  6'b01_01_00  // 00- (index 12)
`define VTX1_REG_VT  6'b01_01_01  // 000 (index 13)
`define VTX1_REG_VB  6'b01_01_10  // 00+ (index 14)
`define VTX1_REG_FA  6'b01_10_00  // 0+- (index 15)
`define VTX1_REG_FT  6'b01_10_01  // 0+0 (index 16)
`define VTX1_REG_FB  6'b01_10_10  // 0++ (index 17)
`define VTX1_REG_TZ  6'b10_00_00  // +-- (index 18)

`endif // VTX1_ISA_GENERATED_VH
