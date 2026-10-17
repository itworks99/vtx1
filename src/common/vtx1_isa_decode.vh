// VTX1 opcode decode helper
// GENERATED FILE - DO NOT EDIT. Source of truth: src/assembler/isa/isa.go
// Include INSIDE a module body (no include guard: each module needs its own
// copy of the function). Requires vtx1_isa_generated.vh for the defines.
//
// vtx1_opcode_info(value) = {microcoded, category[2:0], format[2:0]}
// where value = 40 + 27*t0 + 9*t1 + 3*t2 + t3 over the signed opcode trits.
function [6:0] vtx1_opcode_info(input [6:0] value);
    case (value)
        7'd0 : vtx1_opcode_info = {1'b0, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VADD
        7'd1 : vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // ADD
        7'd2 : vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // SUB
        7'd3 : vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // MUL
        7'd4 : vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // AND
        7'd5 : vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // OR
        7'd6 : vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // NOT
        7'd7 : vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // XOR
        7'd8 : vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // SHL
        7'd9 : vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // SHR
        7'd10: vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // ROL
        7'd11: vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // ROR
        7'd12: vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // CMP
        7'd13: vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // TEST
        7'd14: vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R1}; // INC
        7'd15: vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R1}; // DEC
        7'd16: vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // NEG
        7'd17: vtx1_opcode_info = {1'b0, `VTX1_CAT_MEM, `VTX1_FMT_MEM}; // LD
        7'd18: vtx1_opcode_info = {1'b0, `VTX1_CAT_MEM, `VTX1_FMT_MEM}; // ST
        7'd19: vtx1_opcode_info = {1'b0, `VTX1_CAT_MEM, `VTX1_FMT_MEM}; // VLD
        7'd20: vtx1_opcode_info = {1'b0, `VTX1_CAT_MEM, `VTX1_FMT_MEM}; // VST
        7'd21: vtx1_opcode_info = {1'b0, `VTX1_CAT_MEM, `VTX1_FMT_MEM}; // FLD
        7'd22: vtx1_opcode_info = {1'b0, `VTX1_CAT_MEM, `VTX1_FMT_MEM}; // FST
        7'd23: vtx1_opcode_info = {1'b0, `VTX1_CAT_MEM, `VTX1_FMT_MEM}; // LEA
        7'd24: vtx1_opcode_info = {1'b0, `VTX1_CAT_MEM, `VTX1_FMT_R1}; // PUSH
        7'd25: vtx1_opcode_info = {1'b0, `VTX1_CAT_MEM, `VTX1_FMT_R1}; // POP
        7'd26: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_JUMP}; // JMP
        7'd27: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_JUMP}; // JAL
        7'd28: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_R1}; // JR
        7'd29: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_R1}; // JALR
        7'd30: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_BRANCH}; // BEQ
        7'd31: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_BRANCH}; // BNE
        7'd32: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_BRANCH}; // BLT
        7'd33: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_BRANCH}; // BGE
        7'd34: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_BRANCH}; // BLTU
        7'd35: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_BRANCH}; // BGEU
        7'd36: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_BRANCH}; // BGT
        7'd37: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_BRANCH}; // BLE
        7'd38: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_JUMP}; // CALL
        7'd39: vtx1_opcode_info = {1'b0, `VTX1_CAT_CTRL, `VTX1_FMT_NONE}; // RET
        7'd40: vtx1_opcode_info = {1'b0, `VTX1_CAT_ALU, `VTX1_FMT_NONE}; // NOP
        7'd41: vtx1_opcode_info = {1'b0, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VSUB
        7'd42: vtx1_opcode_info = {1'b0, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VMUL
        7'd43: vtx1_opcode_info = {1'b0, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VAND
        7'd44: vtx1_opcode_info = {1'b0, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VOR
        7'd45: vtx1_opcode_info = {1'b0, `VTX1_CAT_VEC, `VTX1_FMT_R2}; // VNOT
        7'd46: vtx1_opcode_info = {1'b0, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VSHL
        7'd47: vtx1_opcode_info = {1'b0, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VSHR
        7'd48: vtx1_opcode_info = {1'b0, `VTX1_CAT_FP, `VTX1_FMT_R3}; // FADD
        7'd49: vtx1_opcode_info = {1'b0, `VTX1_CAT_FP, `VTX1_FMT_R3}; // FSUB
        7'd50: vtx1_opcode_info = {1'b0, `VTX1_CAT_FP, `VTX1_FMT_R3}; // FMUL
        7'd51: vtx1_opcode_info = {1'b0, `VTX1_CAT_FP, `VTX1_FMT_R3}; // FCMP
        7'd52: vtx1_opcode_info = {1'b0, `VTX1_CAT_FP, `VTX1_FMT_R2}; // FMOV
        7'd53: vtx1_opcode_info = {1'b0, `VTX1_CAT_FP, `VTX1_FMT_R2}; // FNEG
        7'd54: vtx1_opcode_info = {1'b0, `VTX1_CAT_SYS, `VTX1_FMT_NONE}; // WFI
        7'd55: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // DIV
        7'd56: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // MOD
        7'd57: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // UDIV
        7'd58: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R3}; // UMOD
        7'd59: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // SQRT
        7'd60: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // ABS
        7'd61: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // SIN
        7'd62: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // COS
        7'd63: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // TAN
        7'd64: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // ASIN
        7'd65: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // ACOS
        7'd66: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // ATAN
        7'd67: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // EXP
        7'd68: vtx1_opcode_info = {1'b1, `VTX1_CAT_ALU, `VTX1_FMT_R2}; // LOG
        7'd69: vtx1_opcode_info = {1'b1, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VDOT
        7'd70: vtx1_opcode_info = {1'b1, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VREDUCE
        7'd71: vtx1_opcode_info = {1'b1, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VMAX
        7'd72: vtx1_opcode_info = {1'b1, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VMIN
        7'd73: vtx1_opcode_info = {1'b1, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VSUM
        7'd74: vtx1_opcode_info = {1'b1, `VTX1_CAT_VEC, `VTX1_FMT_R3}; // VPERM
        7'd75: vtx1_opcode_info = {1'b1, `VTX1_CAT_MEM, `VTX1_FMT_NONE}; // CACHE
        7'd76: vtx1_opcode_info = {1'b1, `VTX1_CAT_MEM, `VTX1_FMT_NONE}; // FLUSH
        7'd77: vtx1_opcode_info = {1'b1, `VTX1_CAT_MEM, `VTX1_FMT_NONE}; // MEMBAR
        7'd78: vtx1_opcode_info = {1'b1, `VTX1_CAT_SYS, `VTX1_FMT_NONE}; // SYSCALL
        7'd79: vtx1_opcode_info = {1'b1, `VTX1_CAT_SYS, `VTX1_FMT_NONE}; // BREAK
        7'd80: vtx1_opcode_info = {1'b1, `VTX1_CAT_SYS, `VTX1_FMT_NONE}; // HALT
        default: vtx1_opcode_info = {1'b0, 3'd7, `VTX1_FMT_BAD};
    endcase
endfunction
