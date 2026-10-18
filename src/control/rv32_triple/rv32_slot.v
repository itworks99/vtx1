`timescale 1ns / 1ps
// RV32 Issue Slot: decode + single-cycle execute of one RV32IM instruction.
//
// This is the binary half of the issue-width-matched control experiment:
// a mirror of vtx1_slot with the trit datapath replaced by two's-complement
// RV32. The core (rv32_core) uses the identical issue rules and FSM as
// vtx1_core, so a cycle comparison between the two isolates the ISA.
//
// Implemented: the RV32I base integer set minus byte/halfword memory
// access (LW/SW only; LB/LH/SB/SH trap), plus MUL from M (the VTX1 side
// has hardware MUL). FENCE is a NOP; ECALL/EBREAK halt; CSR, the rest of
// M, and anything unknown traps. This covers everything the Tier A
// kernels and their foreseeable extensions use, and trapping (rather than
// silently misexecuting) keeps the control honest.
//
// NOTE (Icarus): no function here reads module signals - sensitivity is
// computed from function arguments only, so such reads go stale.

`ifndef RV32_SLOT_V
`define RV32_SLOT_V

module rv32_slot (
    input  wire [31:0]  word,       // the instruction
    input  wire [1023:0] regs_flat, // registers x0..x31, 32 bits each
    input  wire [31:0]  pc_self,    // this instruction's PC (byte address)
    input  wire [31:0]  pc_next,    // pc_self + 4 (link value)

    // Classification for the issue logic
    output wire         traps,      // would trap if executed
    output wire         is_ld,
    output wire         is_st,
    output wire         is_ctrl,    // branch/jal/jalr: ends the group
    output wire         is_halt,    // ecall/ebreak: ends the group
    output wire         redirect,   // control flow leaves pc_self + 4
    output wire [31:0]  ctrl_target,

    // Single-cycle register write (ALU, LUI/AUIPC, link)
    output reg          wr_en,
    output reg  [4:0]   wr_idx,
    output reg  [31:0]  wr_data,

    // Memory operation payload (LW/SW)
    output wire [4:0]   ld_rd,
    output wire [31:0]  ea,         // byte address
    output wire [31:0]  st_data,

    // Hazard metadata
    output wire         dst_v,
    output wire [4:0]   dst,
    output reg          src1_v, src2_v,
    output reg  [4:0]   src1, src2
);

    // ------------------------------------------------------------------
    // Decode
    // ------------------------------------------------------------------
    wire [6:0] opc = word[6:0];
    wire [4:0] rd  = word[11:7];
    wire [4:0] rs1 = word[19:15];
    wire [4:0] rs2 = word[24:20];
    wire [2:0] f3  = word[14:12];
    wire [6:0] f7  = word[31:25];

    wire [31:0] imm_i = {{20{word[31]}}, word[31:20]};
    wire [31:0] imm_s = {{20{word[31]}}, word[31:25], word[11:7]};
    wire [31:0] imm_b = {{19{word[31]}}, word[31], word[7], word[30:25],
                         word[11:8], 1'b0};
    wire [31:0] imm_u = {word[31:12], 12'b0};
    wire [31:0] imm_j = {{11{word[31]}}, word[31], word[19:12], word[20],
                         word[30:21], 1'b0};

    localparam OP_LUI    = 7'h37;
    localparam OP_AUIPC  = 7'h17;
    localparam OP_JAL    = 7'h6F;
    localparam OP_JALR   = 7'h67;
    localparam OP_BRANCH = 7'h63;
    localparam OP_LOAD   = 7'h03;
    localparam OP_STORE  = 7'h23;
    localparam OP_IMM    = 7'h13;
    localparam OP_OP     = 7'h33;
    localparam OP_FENCE  = 7'h0F;
    localparam OP_SYSTEM = 7'h73;

    // Register reads (x0 reads zero)
    wire [31:0] rs1val = (rs1 == 5'd0) ? 32'd0 : regs_flat[32*rs1 +: 32];
    wire [31:0] rs2val = (rs2 == 5'd0) ? 32'd0 : regs_flat[32*rs2 +: 32];

    // ------------------------------------------------------------------
    // Execute (combinational)
    // ------------------------------------------------------------------
    wire signed [31:0] s1 = rs1val;
    wire signed [31:0] s2 = rs2val;
    wire signed [31:0] simm = imm_i;

    wire is_op  = (opc == OP_OP);
    wire is_imm = (opc == OP_IMM);
    wire [31:0] opb = is_imm ? imm_i : rs2val;
    wire [4:0]  sh  = is_imm ? word[24:20] : rs2val[4:0];

    // Arithmetic shift in its own signed context: inside a mixed
    // signed/unsigned expression, Verilog would evaluate >>> unsigned.
    wire signed [31:0] sra = s1 >>> sh;

    reg  [31:0] alu_result;
    reg         alu_bad;    // unimplemented funct combination
    always @* begin
        alu_result = 32'd0;
        alu_bad    = 1'b0;
        if (is_op && f7 == 7'h01) begin
            if (f3 == 3'd0) alu_result = rs1val * rs2val;  // MUL
            else alu_bad = 1'b1;                           // rest of M
        end else begin
            case (f3)
                3'd0: alu_result = (is_op && f7 == 7'h20) ? rs1val - rs2val
                                                          : rs1val + opb;
                3'd1: alu_result = rs1val << sh;                     // SLL(I)
                3'd2: alu_result = (s1 < (is_imm ? simm : s2)) ? 32'd1 : 32'd0;
                3'd3: alu_result = (rs1val < opb) ? 32'd1 : 32'd0;   // SLTU
                3'd4: alu_result = rs1val ^ opb;
                3'd5: alu_result = (f7 == 7'h20) ? sra               // SRA(I)
                                                 : (rs1val >> sh);   // SRL(I)
                3'd6: alu_result = rs1val | opb;
                3'd7: alu_result = rs1val & opb;
            endcase
            if ((f3 == 3'd1 || f3 == 3'd5) && f7[4:0] != 5'd0)
                alu_bad = 1'b1; // malformed shift funct7
            if (is_op && f3 != 3'd0 && f3 != 3'd5 && f7 != 7'd0)
                alu_bad = 1'b1;
        end
    end

    // Branch decision
    reg branch_take;
    always @* begin
        case (f3)
            3'd0: branch_take = (rs1val == rs2val);          // BEQ
            3'd1: branch_take = (rs1val != rs2val);          // BNE
            3'd4: branch_take = (s1 < s2);                   // BLT
            3'd5: branch_take = (s1 >= s2);                  // BGE
            3'd6: branch_take = (rs1val < rs2val);           // BLTU
            3'd7: branch_take = (rs1val >= rs2val);          // BGEU
            default: branch_take = 1'b0;
        endcase
    end
    wire branch_bad = (opc == OP_BRANCH) && (f3 == 3'd2 || f3 == 3'd3);

    // ------------------------------------------------------------------
    // Classification
    // ------------------------------------------------------------------
    assign is_ld   = (opc == OP_LOAD);
    assign is_st   = (opc == OP_STORE);
    wire   f_br    = (opc == OP_BRANCH);
    wire   f_jal   = (opc == OP_JAL);
    wire   f_jalr  = (opc == OP_JALR);
    assign is_ctrl = f_br || f_jal || f_jalr;
    assign is_halt = (opc == OP_SYSTEM) && (f3 == 3'd0); // ECALL/EBREAK

    assign traps =
           (is_ld && f3 != 3'd2)                          // LW only
        || (is_st && f3 != 3'd2)                          // SW only
        || ((is_op || is_imm) && alu_bad)
        || branch_bad
        || ((opc == OP_SYSTEM) && (f3 != 3'd0))           // CSR ops
        || !(opc == OP_LUI || opc == OP_AUIPC || f_jal || f_jalr || f_br
             || is_ld || is_st || is_imm || is_op
             || opc == OP_FENCE || opc == OP_SYSTEM);

    assign redirect = (f_br && branch_take) || f_jal || f_jalr;
    assign ctrl_target =
        f_br  ? pc_self + imm_b :
        f_jal ? pc_self + imm_j :
                (rs1val + imm_i) & 32'hFFFF_FFFE; // JALR

    // ------------------------------------------------------------------
    // Single-cycle register write
    // ------------------------------------------------------------------
    always @* begin
        wr_en   = 1'b0;
        wr_idx  = rd;
        wr_data = alu_result;
        if (!traps) begin
            if (is_op || is_imm) begin
                wr_en = 1'b1;
            end else if (opc == OP_LUI) begin
                wr_en   = 1'b1;
                wr_data = imm_u;
            end else if (opc == OP_AUIPC) begin
                wr_en   = 1'b1;
                wr_data = pc_self + imm_u;
            end else if (f_jal || f_jalr) begin
                wr_en   = 1'b1;
                wr_data = pc_next;
            end
        end
    end

    // ------------------------------------------------------------------
    // Memory payload and hazard metadata
    // ------------------------------------------------------------------
    assign ld_rd   = rd;
    assign ea      = rs1val + (is_st ? imm_s : imm_i);
    assign st_data = rs2val;

    assign dst_v = wr_en || is_ld;
    assign dst   = is_ld ? rd : wr_idx;

    always @* begin
        src1_v = 1'b0; src1 = 5'd0;
        src2_v = 1'b0; src2 = 5'd0;
        if (is_op || f_br) begin
            src1_v = 1'b1; src1 = rs1;
            src2_v = 1'b1; src2 = rs2;
        end else if (is_imm || is_ld || f_jalr) begin
            src1_v = 1'b1; src1 = rs1;
        end else if (is_st) begin
            src1_v = 1'b1; src1 = rs1;  // base
            src2_v = 1'b1; src2 = rs2;  // data
        end
        // LUI/AUIPC/JAL/FENCE/SYSTEM: no register sources
    end

endmodule

`endif // RV32_SLOT_V
