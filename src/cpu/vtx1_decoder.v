`timescale 1ns / 1ps
// VTX1 Trit-Native Instruction Decoder
// Part of the VTX1 Ternary System-on-Chip
//
// Decodes one 18-trit instruction word (36 bits, 2 bits/trit) according to
// the trit-native ISA. The ISA single source of truth is
// src/assembler/isa/isa.go; the opcode constants and the per-opcode
// format/category table used here are GENERATED from it
// (vtx1_isa_generated.vh / vtx1_isa_decode.vh), so this decoder cannot
// drift from the assembler.
//
// Field layout (trit t, 0 = most significant, occupies bits [35-2t:34-2t]):
//   all:     opcode  t0..t3
//   R3:      rd t4..t6 | rs1 t7..t9 | mode t10 | (rs2 t15..t17 or imm7 t11..t17)
//   R2:      rd t4..t6 | rs t7..t9
//   R1:      r  t4..t6
//   Mem:     reg t4..t6 | mode t7 | (addr10 t8..t17, or base t8..t10 with
//            index t11..t13 or offset7 t11..t17)
//   Branch:  rs1 t4..t6 | rs2 t7..t9 | rel8 t10..t17 (rs2 = TZ: compare to 0)
//   Jump:    target14 t4..t17
//
// Outputs use binary register indices (0..26; TZ = 18) and sign-extended
// binary immediates, so downstream units keep conventional interfaces while
// the instruction stream itself is pure balanced ternary.

`ifndef VTX1_DECODER_V
`define VTX1_DECODER_V

`include "vtx1_isa_generated.vh"

module vtx1_decoder (
    input  wire [35:0] word,          // one 18-trit instruction word

    output wire        valid,         // no undefined (2'b11) trit pairs
    output wire [6:0]  opcode,        // opcode value 0..80
    output wire [2:0]  format,        // `VTX1_FMT_*
    output wire [2:0]  category,      // `VTX1_CAT_*
    output wire        multicycle,

    // Register indices (0..26 in the unified space; TZ = 18).
    // reg1: rd (R3/R2), r (R1), data reg (Mem), rs1 (Branch)
    // reg2: rs1 (R3), rs (R2), base (Mem), rs2 (Branch; TZ = compare to 0)
    // reg3: rs2 (R3 register mode), index (Mem base+index)
    output wire [4:0]  reg1,
    output wire [4:0]  reg2,
    output wire [4:0]  reg3,

    // Immediate: R3 imm7 / Mem direct addr10 or offset7 / Branch rel8
    // (PC-relative, in words) / Jump target14 (absolute word address).
    output wire        uses_imm,
    output wire signed [31:0] imm,

    // Memory addressing mode (FMT_MEM only)
    output wire [1:0]  mem_mode,      // 0 direct, 1 base+index, 2 base+offset

    // Full word interpreted as an 18-trit balanced number (for data words)
    output wire signed [35:0] word_value
);

    `include "vtx1_isa_decode.vh"

    localparam MEM_DIRECT     = 2'd0;
    localparam MEM_BASE_INDEX = 2'd1;
    localparam MEM_BASE_OFF   = 2'd2;

    // ------------------------------------------------------------------
    // Trit extraction: trit t occupies bits [35-2t : 34-2t]
    // ------------------------------------------------------------------
    function signed [2:0] tv(input [1:0] pair);
        case (pair)
            2'b00:   tv = -3'sd1;
            2'b01:   tv = 3'sd0;
            2'b10:   tv = 3'sd1;
            default: tv = 3'sd0; // undefined: masked by `valid`
        endcase
    endfunction

    wire signed [2:0] t [0:17];
    genvar gi;
    generate
        for (gi = 0; gi < 18; gi = gi + 1) begin : trits
            assign t[gi] = tv(word[35-2*gi -: 2]);
        end
    endgenerate

    // Validity: no 2'b11 pair anywhere in the word.
    wire [17:0] undef_pair;
    generate
        for (gi = 0; gi < 18; gi = gi + 1) begin : check
            assign undef_pair[gi] = (word[35-2*gi -: 2] == 2'b11);
        end
    endgenerate
    assign valid = ~|undef_pair && (format != `VTX1_FMT_BAD);

    // ------------------------------------------------------------------
    // Field values (balanced ternary, signed)
    // ------------------------------------------------------------------
    wire signed [31:0] opc_signed = 27*t[0] + 9*t[1] + 3*t[2] + t[3];
    assign opcode = opc_signed[6:0] + 7'd40; // value = balanced + 40

    wire [6:0] info = vtx1_opcode_info(opcode);
    assign multicycle = info[6];
    assign category   = info[5:3];
    assign format     = info[2:0];

    // Register fields at the three positions the formats use.
    wire signed [31:0] idx4  = 9*t[4]  + 3*t[5]  + t[6]  + 13; // t4..t6
    wire signed [31:0] idx7  = 9*t[7]  + 3*t[8]  + t[9]  + 13; // t7..t9
    wire signed [31:0] idx8  = 9*t[8]  + 3*t[9]  + t[10] + 13; // t8..t10 (mem base)
    wire signed [31:0] idx11 = 9*t[11] + 3*t[12] + t[13] + 13; // t11..t13 (mem index)
    wire signed [31:0] idx15 = 9*t[15] + 3*t[16] + t[17] + 13; // t15..t17 (R3 rs2)

    // Immediate candidates.
    wire signed [31:0] imm7_11 = 729*t[11] + 243*t[12] + 81*t[13] + 27*t[14]
                               + 9*t[15] + 3*t[16] + t[17];            // t11..t17
    wire signed [31:0] addr10  = 19683*t[8] + 6561*t[9] + 2187*t[10]
                               + 729*t[11] + 243*t[12] + 81*t[13]
                               + 27*t[14] + 9*t[15] + 3*t[16] + t[17]; // t8..t17
    wire signed [31:0] rel8    = 2187*t[10] + 729*t[11] + 243*t[12]
                               + 81*t[13] + 27*t[14] + 9*t[15]
                               + 3*t[16] + t[17];                      // t10..t17
    wire signed [31:0] tgt14   = 1594323*t[4] + 531441*t[5] + 177147*t[6]
                               + 59049*t[7] + 19683*t[8] + 6561*t[9]
                               + 2187*t[10] + 729*t[11] + 243*t[12]
                               + 81*t[13] + 27*t[14] + 9*t[15]
                               + 3*t[16] + t[17];                      // t4..t17

    // ------------------------------------------------------------------
    // Format-dependent output muxing
    // ------------------------------------------------------------------
    wire r3_imm_mode = (format == `VTX1_FMT_R3) && (t[10] == 3'sd1);

    assign mem_mode = (t[7] == 3'sd0)  ? MEM_DIRECT :
                      (t[7] == 3'sd1)  ? MEM_BASE_INDEX :
                                         MEM_BASE_OFF;

    assign reg1 = idx4[4:0];

    assign reg2 = (format == `VTX1_FMT_MEM) ? idx8[4:0] : idx7[4:0];

    assign reg3 = (format == `VTX1_FMT_R3)  ? idx15[4:0] :
                  (format == `VTX1_FMT_MEM) ? idx11[4:0] : 5'd0;

    assign uses_imm = r3_imm_mode
                   || (format == `VTX1_FMT_BRANCH)
                   || (format == `VTX1_FMT_JUMP)
                   || ((format == `VTX1_FMT_MEM) && (mem_mode != MEM_BASE_INDEX));

    assign imm = (format == `VTX1_FMT_JUMP)   ? tgt14 :
                 (format == `VTX1_FMT_BRANCH) ? rel8 :
                 (format == `VTX1_FMT_MEM)    ? ((mem_mode == MEM_DIRECT) ? addr10 : imm7_11) :
                 r3_imm_mode                  ? imm7_11 : 32'sd0;

    // Full-word balanced value (for data words). Indexes `word` directly
    // (not the trit array) so the @* sensitivity is just the input vector.
    integer wi;
    reg signed [35:0] wv;
    always @* begin
        wv = 36'sd0;
        for (wi = 0; wi < 18; wi = wi + 1)
            wv = wv * 3 + tv(word[35-2*wi -: 2]);
    end
    assign word_value = wv;

endmodule

// ----------------------------------------------------------------------
// VLIW bundle decoder: 3 slots = 54 trits = 108 bits (`VLIW_WIDTH).
// Slot 0 (first operation) occupies the most significant bits.
// ----------------------------------------------------------------------
module vtx1_bundle_decoder (
    input  wire [107:0] bundle,

    output wire [2:0]   valid,
    output wire [6:0]   opcode_a, opcode_b, opcode_c,
    output wire [2:0]   format_a, format_b, format_c,
    output wire [2:0]   category_a, category_b, category_c,
    output wire [2:0]   multicycle,
    output wire [4:0]   reg1_a, reg2_a, reg3_a,
    output wire [4:0]   reg1_b, reg2_b, reg3_b,
    output wire [4:0]   reg1_c, reg2_c, reg3_c,
    output wire [2:0]   uses_imm,
    output wire signed [31:0] imm_a, imm_b, imm_c,
    output wire [1:0]   mem_mode_a, mem_mode_b, mem_mode_c
);

    vtx1_decoder slot_a (
        .word(bundle[107:72]), .valid(valid[0]), .opcode(opcode_a),
        .format(format_a), .category(category_a), .multicycle(multicycle[0]),
        .reg1(reg1_a), .reg2(reg2_a), .reg3(reg3_a),
        .uses_imm(uses_imm[0]), .imm(imm_a), .mem_mode(mem_mode_a),
        .word_value()
    );
    vtx1_decoder slot_b (
        .word(bundle[71:36]), .valid(valid[1]), .opcode(opcode_b),
        .format(format_b), .category(category_b), .multicycle(multicycle[1]),
        .reg1(reg1_b), .reg2(reg2_b), .reg3(reg3_b),
        .uses_imm(uses_imm[1]), .imm(imm_b), .mem_mode(mem_mode_b),
        .word_value()
    );
    vtx1_decoder slot_c (
        .word(bundle[35:0]), .valid(valid[2]), .opcode(opcode_c),
        .format(format_c), .category(category_c), .multicycle(multicycle[2]),
        .reg1(reg1_c), .reg2(reg2_c), .reg3(reg3_c),
        .uses_imm(uses_imm[2]), .imm(imm_c), .mem_mode(mem_mode_c),
        .word_value()
    );

endmodule

`endif // VTX1_DECODER_V
