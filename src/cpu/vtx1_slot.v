`timescale 1ns / 1ps
// VTX1 Issue Slot: decode + single-cycle execute of one operation word
// Part of the VTX1 Ternary System-on-Chip
//
// One slot = one 18-trit operation word evaluated combinationally against
// the architectural register file: its ALU result, effective address,
// branch decision and control-flow target, plus the classification and
// hazard metadata (destination and source register indices) the core's
// issue logic needs to decide how many consecutive words can commit in
// the same cycle.
//
// The slot never writes state itself: the core commits wr_*/ld_*/st_*
// only for the slots it actually issues, so an un-issued slot is free.
//
// `pc_self` is this slot's own PC (the group base + slot position), which
// branch targets and link values are relative to; `pc_next` = pc_self + 1.

`ifndef VTX1_SLOT_V
`define VTX1_SLOT_V

`include "vtx1_isa_generated.vh"

module vtx1_slot (
    input  wire [35:0]  word,       // the operation word
    input  wire [647:0] regs_flat,  // registers 0..17, 36 bits each
    input  wire [35:0]  pc_self,    // this operation's PC (trit word)
    input  wire [35:0]  pc_next,    // pc_self + 1 (link value, trit word)

    // Classification for the issue logic
    output wire         traps,      // would trap if executed
    output wire         is_ld,
    output wire         is_st,
    output wire         is_div,     // DIV/MOD: multi-cycle, issues alone
    output wire         is_mod,     // distinguishes MOD within is_div
    output wire         is_stack,   // PUSH/POP: issues alone (two dests)
    output wire         is_ctrl,    // branch/jump/ret/jr/jalr: ends the group
    output wire         is_halt,    // WFI/HALT/SYSCALL/BREAK: ends the group
    output wire         redirect,   // control flow leaves pc_self + 1
    output wire [35:0]  ctrl_target,// next PC when redirect (trit word)

    // Single-cycle register write (ALU, LEA, INC/DEC, link)
    output reg          wr_en,
    output reg  [4:0]   wr_idx,
    output reg  [35:0]  wr_data,

    // Memory operation payload (LD/ST)
    output wire [4:0]   ld_rd,
    output wire [31:0]  ea_bin,
    output wire [35:0]  st_data,

    // Operands for the core's shared multi-cycle units (DIV/MOD)
    output wire [35:0]  op_a,
    output wire [35:0]  op_b,

    // Hazard metadata: architected destination (incl. the LD result) and
    // up to three source registers
    output wire         dst_v,
    output wire [4:0]   dst,
    output reg          src1_v, src2_v, src3_v,
    output reg  [4:0]   src1, src2, src3
);

    localparam [35:0] ZERO_WORD = 36'h555555555; // 18 zero trits
    localparam [4:0]  REG_TA    = 5'd8;          // link register

    // NOTE: no helper function reads regs_flat here. Icarus computes
    // continuous-assign sensitivity from function ARGUMENTS only, so a
    // regword(idx) function would never re-evaluate when the register
    // file changes. Register reads are explicit wires naming regs_flat.

    // ------------------------------------------------------------------
    // Trit helpers
    // ------------------------------------------------------------------
    function signed [2:0] tv(input [1:0] pair);
        case (pair)
            2'b00:   tv = -3'sd1;
            2'b10:   tv = 3'sd1;
            default: tv = 3'sd0;
        endcase
    endfunction

    // Trit word -> signed binary (x3 adds only); for the memory address.
    function signed [35:0] t2b(input [35:0] w);
        integer i;
        begin
            t2b = 36'sd0;
            for (i = 17; i >= 0; i = i - 1)
                t2b = t2b * 3 + tv(w[2*i+1 -: 2]);
        end
    endfunction

    // Sign of a word: its most significant nonzero trit.
    function [1:0] word_sign(input [35:0] w);
        integer i;
        begin
            word_sign = 2'b01;
            for (i = 17; i >= 0; i = i - 1)
                if (word_sign == 2'b01 && w[2*i+1 -: 2] != 2'b01)
                    word_sign = w[2*i+1 -: 2];
        end
    endfunction

    // ------------------------------------------------------------------
    // Decode
    // ------------------------------------------------------------------
    wire        dec_valid;
    wire [6:0]  opcode;
    wire [2:0]  format, category;
    wire        multicycle;
    wire [4:0]  reg1, reg2, reg3;
    wire        uses_imm;
    wire [1:0]  mem_mode;
    wire [35:0] imm_trits;

    vtx1_decoder dec (
        .word(word), .valid(dec_valid), .opcode(opcode),
        .format(format), .category(category), .multicycle(multicycle),
        .reg1(reg1), .reg2(reg2), .reg3(reg3),
        .uses_imm(uses_imm), .imm(), .mem_mode(mem_mode),
        .imm_trits(imm_trits), .word_value()
    );

    // ------------------------------------------------------------------
    // Execute (combinational)
    // ------------------------------------------------------------------

    // Register reads (TZ = 18 and above read as zero)
    wire [35:0] r1val = (reg1 >= 5'd18) ? ZERO_WORD : regs_flat[36*reg1 +: 36];
    wire [35:0] r2val = (reg2 >= 5'd18) ? ZERO_WORD : regs_flat[36*reg2 +: 36];
    wire [35:0] r3val = (reg3 >= 5'd18) ? ZERO_WORD : regs_flat[36*reg3 +: 36];
    wire [35:0] taval = regs_flat[36*REG_TA +: 36];

    // ALU
    wire [35:0] alu_a =
        (format == `VTX1_FMT_R1) ? r1val : r2val;
    wire [35:0] alu_b =
        (format == `VTX1_FMT_R1) ? 36'h555555556 /* +1 */ :
        uses_imm                 ? imm_trits : r3val;

    wire [35:0] alu_result;
    wire        alu_unsup;
    vtx1_alu u_alu (.a(alu_a), .b(alu_b), .opcode(opcode),
                    .result(alu_result), .unsupported(alu_unsup));

    assign op_a = alu_a;
    assign op_b = alu_b;

    // Branch compare: sign of (rs1 - rs2); a nonzero carry trit is the
    // true sign (overflow-correct compare).
    wire [35:0] br_diff;
    wire [1:0]  br_cout;
    vtx1_trit_adder u_brsub (.a(r1val), .b(r2val),
                             .negate_b(1'b1), .sum(br_diff), .cout(br_cout));
    wire [1:0] br_sign = (br_cout != 2'b01) ? br_cout : word_sign(br_diff);

    reg branch_take;
    always @* begin
        branch_take = 1'b0;
        case (opcode)
            `VTX1_OPV_BEQ: branch_take = (r1val == r2val);
            `VTX1_OPV_BNE: branch_take = (r1val != r2val);
            `VTX1_OPV_BLT: branch_take = (br_sign == 2'b00);
            `VTX1_OPV_BGE: branch_take = (br_sign != 2'b00);
            `VTX1_OPV_BGT: branch_take = (br_sign == 2'b10);
            `VTX1_OPV_BLE: branch_take = (br_sign != 2'b10);
            default: ;
        endcase
    end

    // Branch target: pc_self + rel8
    wire [35:0] branch_target;
    vtx1_trit_adder u_btgt (.a(pc_self), .b(imm_trits), .negate_b(1'b0),
                            .sum(branch_target), .cout());

    // Stack ops (format R1, category MEM); declared here because the EA
    // and stack-pointer datapaths below use them.
    wire op_push = (format == `VTX1_FMT_R1) && (opcode == `VTX1_OPV_PUSH);
    wire op_pop  = (format == `VTX1_FMT_R1) && (opcode == `VTX1_OPV_POP);

    // Stack ABI (ratified): TB (index 7) is the stack pointer, full-
    // descending, word-addressed. PUSH: TB <- TB-1, mem[TB] <- r (old r).
    // POP: r <- mem[TB], TB <- TB+1 (for POP TB the loaded value wins).
    // PUSH/POP present themselves as a store/load plus a same-cycle TB
    // write; POP has two destinations, so stack ops issue alone.
    localparam [4:0] REG_TB = 5'd7;
    wire [35:0] tbval = regs_flat[36*REG_TB +: 36];
    wire [35:0] sp_adj;
    vtx1_trit_adder u_sp (.a(tbval), .b(36'h555555556 /* +1 */),
                          .negate_b(op_push), .sum(sp_adj), .cout());

    assign is_stack = op_push || op_pop;

    // Effective address (trit add; binary only at the port)
    wire [35:0] ea_sum;
    vtx1_trit_adder u_ea (
        .a(r2val),
        .b((mem_mode == 2'd1) ? r3val : imm_trits),
        .negate_b(1'b0), .sum(ea_sum), .cout());
    wire [35:0] ea_trits = (mem_mode == 2'd0) ? imm_trits : ea_sum;
    wire [35:0] ea_src = op_push ? sp_adj : op_pop ? tbval : ea_trits;
    assign ea_bin = t2b(ea_src);

    assign ld_rd   = reg1;
    assign st_data = r1val;

    // ------------------------------------------------------------------
    // Classification
    // ------------------------------------------------------------------
    wire f_none   = (format == `VTX1_FMT_NONE);
    wire f_r1     = (format == `VTX1_FMT_R1);
    wire f_r23    = (format == `VTX1_FMT_R2) || (format == `VTX1_FMT_R3);
    wire f_mem    = (format == `VTX1_FMT_MEM);
    wire f_branch = (format == `VTX1_FMT_BRANCH);
    wire f_jump   = (format == `VTX1_FMT_JUMP);

    wire op_ret  = f_none && (opcode == `VTX1_OPV_RET);
    wire op_jr   = f_r1   && (opcode == `VTX1_OPV_JR);
    wire op_jalr = f_r1   && (opcode == `VTX1_OPV_JALR);
    wire op_lea  = f_mem  && (opcode == `VTX1_OPV_LEA);
    wire op_link = f_jump && (opcode == `VTX1_OPV_JAL ||
                              opcode == `VTX1_OPV_CALL);

    assign is_ld   = (f_mem && (opcode == `VTX1_OPV_LD)) || op_pop;
    assign is_st   = (f_mem && (opcode == `VTX1_OPV_ST)) || op_push;
    assign is_div  = f_r23 && (opcode == `VTX1_OPV_DIV ||
                               opcode == `VTX1_OPV_MOD);
    assign is_mod  = f_r23 && (opcode == `VTX1_OPV_MOD);
    assign is_halt = f_none && (opcode == `VTX1_OPV_WFI  ||
                                opcode == `VTX1_OPV_HALT ||
                                opcode == `VTX1_OPV_SYSCALL ||
                                opcode == `VTX1_OPV_BREAK);
    assign is_ctrl = f_branch || f_jump || op_jr || op_jalr || op_ret;

    // Trap conditions: mirror of the core's per-format legality rules.
    assign traps = !dec_valid
        || (f_none && !(opcode == `VTX1_OPV_NOP || is_halt || op_ret ||
                        opcode == `VTX1_OPV_CACHE ||
                        opcode == `VTX1_OPV_FLUSH ||
                        opcode == `VTX1_OPV_MEMBAR))
        || (f_r1 && !(opcode == `VTX1_OPV_INC || opcode == `VTX1_OPV_DEC ||
                      op_jr || op_jalr || op_push || op_pop))
        || (f_r23 && alu_unsup && !is_div)          // SQRT, vec/FP (DIV/MOD
                                                    // go to the divider FSM)
        || (f_mem && !(is_ld || is_st || op_lea))   // VLD/VST/FLD/FST
        || !(f_none || f_r1 || f_r23 || f_mem || f_branch || f_jump);

    // Control-flow resolution
    assign redirect = (f_branch && branch_take) || f_jump || op_jr || op_jalr || op_ret;
    assign ctrl_target =
        f_branch ? branch_target :
        f_jump   ? imm_trits :          // absolute target as trits
        op_ret   ? taval :
                   r1val;               // JR / JALR

    // ------------------------------------------------------------------
    // Single-cycle register write
    // ------------------------------------------------------------------
    always @* begin
        wr_en   = 1'b0;
        wr_idx  = reg1;
        wr_data = alu_result;
        if (!traps) begin
            if (op_push || op_pop) begin // stack pointer update
                wr_en   = 1'b1;
                wr_idx  = REG_TB;
                wr_data = sp_adj;
            end else if (f_r23 && !is_div) begin // DIV/MOD write back from the FSM
                wr_en = 1'b1;
            end else if (f_r1 && (opcode == `VTX1_OPV_INC ||
                                  opcode == `VTX1_OPV_DEC)) begin
                wr_en = 1'b1;
            end else if (op_lea) begin
                wr_en   = 1'b1;
                wr_data = ea_trits;
            end else if (op_jalr || op_link) begin
                wr_en   = 1'b1;
                wr_idx  = REG_TA;
                wr_data = pc_next;
            end
        end
    end

    // ------------------------------------------------------------------
    // Hazard metadata
    // ------------------------------------------------------------------
    assign dst_v = wr_en || is_ld;
    assign dst   = is_ld ? reg1 : wr_idx;

    always @* begin
        src1_v = 1'b0; src1 = 5'd0;
        src2_v = 1'b0; src2 = 5'd0;
        src3_v = 1'b0; src3 = 5'd0;
        case (format)
            `VTX1_FMT_R3: begin
                src1_v = 1'b1; src1 = reg2;
                if (!uses_imm) begin src2_v = 1'b1; src2 = reg3; end
            end
            `VTX1_FMT_R2: begin
                src1_v = 1'b1; src1 = reg2;
            end
            `VTX1_FMT_R1: begin
                if (op_pop) begin
                    src1_v = 1'b1; src1 = REG_TB;
                end else begin
                    src1_v = 1'b1; src1 = reg1; // INC/DEC, JR/JALR, PUSH data
                    if (op_push) begin src2_v = 1'b1; src2 = REG_TB; end
                end
            end
            `VTX1_FMT_NONE: begin
                if (op_ret) begin src1_v = 1'b1; src1 = REG_TA; end
            end
            `VTX1_FMT_MEM: begin
                if (mem_mode != 2'd0) begin
                    src1_v = 1'b1; src1 = reg2;           // base
                    if (mem_mode == 2'd1) begin
                        src2_v = 1'b1; src2 = reg3;        // index
                    end
                end
                if (is_st) begin src3_v = 1'b1; src3 = reg1; end // data
            end
            `VTX1_FMT_BRANCH: begin
                src1_v = 1'b1; src1 = reg1;
                src2_v = 1'b1; src2 = reg2;
            end
            default: ; // JUMP: no register sources
        endcase
    end

endmodule

`endif // VTX1_SLOT_V
