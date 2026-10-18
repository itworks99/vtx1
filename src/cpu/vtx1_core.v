`timescale 1ns / 1ps
// VTX1 Trit-Native CPU Core (v3: triple-issue native balanced-ternary)
// Part of the VTX1 Ternary System-on-Chip
//
// A core built from three vtx1_slot units (decode + single-cycle execute
// of one operation word each): hardwired control, no microcode, and NO
// binary<->ternary conversion in the datapath - registers, the PC,
// immediates, effective addresses and all ALU results are trit words end
// to end; the only trit->binary conversion is t2b() toward the memory
// address bus.
//
// Issue model: up to three CONSECUTIVE words commit in the same cycle.
// The hardware enforces the VLIW packing rules dynamically, so they hold
// whether the assembler packed the words into a bundle or not:
//   - at most one memory operation per group
//   - a control operation (branch/jump/ret/jr/jalr) or halt ends the group
//   - no write-write conflicts and no read-after-write within a group
//     (dependent words simply start the next group)
// Co-issued words are therefore provably independent, so parallel commit
// is bit-identical to sequential execution - the ISA simulator cross-check
// (retired word count) is preserved exactly.
//
// Fetch: the tightly-coupled instruction memory delivers the 3-word group
// at the next PC, which is resolved combinationally during EXEC (branches
// included) - so TCM-resident code needs no fetch state at all and
// straight-line independent code runs at up to 3 words per cycle. A TCM
// miss (or imem_hit tied low) falls back to the bus port, one word per
// fetch, single issue: the original scalar core behavior.
//
// DIV/MOD belong to the pending iterative divider FSM and trap for now;
// SQRT, PUSH/POP (no stack ABI yet) and vector/FP operations trap.
//
// Memory interface: one word-addressed port with a ready handshake, so the
// core connects equally to a simple synchronous RAM or to the SoC bus matrix.

`ifndef VTX1_CORE_V
`define VTX1_CORE_V

`include "vtx1_isa_generated.vh"

module vtx1_core (
    input  wire        clk,
    input  wire        rst_n,

    // Unified word-addressed memory port with ready handshake
    output reg         mem_req,
    output reg         mem_wr,
    output reg  [31:0] mem_addr,
    output reg  [35:0] mem_wdata,
    input  wire [35:0] mem_rdata,
    input  wire        mem_ready,

    // Tightly-coupled instruction memory port (synchronous read):
    // imem_addr is valid while imem_en is high; the TCM answers hit
    // combinationally and delivers the 3-word group (imem_addr,+1,+2)
    // in imem_rdata0..2 one clock later. Tie imem_hit low to fetch
    // everything over the bus port, one word at a time.
    output reg         imem_en,
    output reg  [31:0] imem_addr,
    input  wire [35:0] imem_rdata0,
    input  wire [35:0] imem_rdata1,
    input  wire [35:0] imem_rdata2,
    input  wire        imem_hit,

    // Execution state
    output reg         halted,
    output reg         trap,        // invalid word or unimplemented operation
    output reg  [31:0] trap_pc,

    // Metrics (cross-checked against the ISA simulator)
    output reg  [31:0] retired,     // instruction words executed
    output reg  [31:0] cycles,

    // Debug observation
    output wire [31:0] debug_pc
);

    // Canonical trit words
    localparam [35:0] ZERO_WORD = 36'h555555555; // 18 zero trits
    localparam [35:0] ONE_WORD  = 36'h555555556; // 0...0,+1

    // ------------------------------------------------------------------
    // Trit helpers (no division anywhere: synthesizable)
    // ------------------------------------------------------------------
    function signed [2:0] tv(input [1:0] pair);
        case (pair)
            2'b00:   tv = -3'sd1;
            2'b10:   tv = 3'sd1;
            default: tv = 3'sd0;
        endcase
    endfunction

    // Trit word -> signed binary (x3 adds only). Used solely for the
    // memory address bus and debug.
    function signed [35:0] t2b(input [35:0] w);
        integer i;
        begin
            t2b = 36'sd0;
            for (i = 17; i >= 0; i = i - 1) // from MST (bits [35:34]) down
                t2b = t2b * 3 + tv(w[2*i+1 -: 2]);
        end
    endfunction

    // CONSTANT function: binary -> trit word. Contains division, but it is
    // evaluated at elaboration time only (reset vector), costing no gates.
    function [35:0] const_b2t(input integer v);
        integer i, x, r;
        begin
            x = v;
            const_b2t = ZERO_WORD;
            for (i = 0; i < 18; i = i + 1) begin
                r = x % 3;
                x = x / 3;
                if (r == 2) begin r = -1; x = x + 1; end
                if (r == -2) begin r = 1; x = x - 1; end
                const_b2t[2*i+1 -: 2] = (r == -1) ? 2'b00 : (r == 1) ? 2'b10 : 2'b01;
            end
        end
    endfunction

    localparam [35:0] RESET_PC = const_b2t(32'h1000);

    // ------------------------------------------------------------------
    // Register file: unified index space 0..18 (TZ = 18 reads zero,
    // writes discarded). Registers hold canonical trit words.
    // ------------------------------------------------------------------
    reg [35:0] regs [0:17];

    wire [647:0] regs_flat;
    genvar gi;
    generate
        for (gi = 0; gi < 18; gi = gi + 1) begin : flat
            assign regs_flat[36*gi +: 36] = regs[gi];
        end
    endgenerate

    task writereg(input [4:0] idx, input [35:0] w);
        if (idx < 5'd18)
            regs[idx] <= w;
    endtask

    // ------------------------------------------------------------------
    // Program counter chain: p0 = PC of slot 0, p1..p3 = +1, +2, +3
    // (trit words; binary only at the address buses)
    // ------------------------------------------------------------------
    reg  [35:0] pc;
    wire [35:0] p1, p2, p3;
    vtx1_trit_adder u_pc1 (.a(pc), .b(ONE_WORD), .negate_b(1'b0), .sum(p1), .cout());
    vtx1_trit_adder u_pc2 (.a(p1), .b(ONE_WORD), .negate_b(1'b0), .sum(p2), .cout());
    vtx1_trit_adder u_pc3 (.a(p2), .b(ONE_WORD), .negate_b(1'b0), .sum(p3), .cout());

    wire [31:0] pc_bin = t2b(pc);
    assign debug_pc = pc_bin;

    // ------------------------------------------------------------------
    // Issue slots: the group's three consecutive words. On the fast path
    // the words come from the TCM read registers; a bus fetch delivers a
    // single word into ir (slots 1 and 2 are then unavailable).
    // ------------------------------------------------------------------
    reg  [35:0] ir;           // bus-fetched instruction word
    reg         ir_from_imem; // 1: group comes from the TCM (3 words)

    wire [35:0] w0 = ir_from_imem ? imem_rdata0 : ir;
    wire [35:0] w1 = imem_rdata1;
    wire [35:0] w2 = imem_rdata2;
    wire        avail12 = ir_from_imem;

    wire        s0_traps, s0_is_ld, s0_is_st, s0_is_div, s0_is_mod, s0_is_stack, s0_is_ctrl, s0_is_halt, s0_redirect;
    wire [35:0] s0_ctrl_target;
    wire        s0_wr_en;
    wire [4:0]  s0_wr_idx;
    wire [35:0] s0_wr_data;
    wire [4:0]  s0_ld_rd;
    wire [31:0] s0_ea_bin;
    wire [35:0] s0_st_data;
    wire [35:0] s0_op_a, s0_op_b;
    wire        s0_dst_v, s0_src1_v, s0_src2_v, s0_src3_v;
    wire [4:0]  s0_dst, s0_src1, s0_src2, s0_src3;

    vtx1_slot slot0 (
        .word(w0), .regs_flat(regs_flat), .pc_self(pc), .pc_next(p1),
        .traps(s0_traps), .is_ld(s0_is_ld), .is_st(s0_is_st),
        .is_div(s0_is_div), .is_mod(s0_is_mod), .is_stack(s0_is_stack),
        .is_ctrl(s0_is_ctrl), .is_halt(s0_is_halt),
        .redirect(s0_redirect), .ctrl_target(s0_ctrl_target),
        .wr_en(s0_wr_en), .wr_idx(s0_wr_idx), .wr_data(s0_wr_data),
        .ld_rd(s0_ld_rd), .ea_bin(s0_ea_bin), .st_data(s0_st_data),
        .op_a(s0_op_a), .op_b(s0_op_b),
        .dst_v(s0_dst_v), .dst(s0_dst),
        .src1_v(s0_src1_v), .src2_v(s0_src2_v), .src3_v(s0_src3_v),
        .src1(s0_src1), .src2(s0_src2), .src3(s0_src3)
    );

    wire        s1_traps, s1_is_ld, s1_is_st, s1_is_div, s1_is_stack, s1_is_ctrl, s1_is_halt, s1_redirect;
    wire [35:0] s1_ctrl_target;
    wire        s1_wr_en;
    wire [4:0]  s1_wr_idx;
    wire [35:0] s1_wr_data;
    wire [4:0]  s1_ld_rd;
    wire [31:0] s1_ea_bin;
    wire [35:0] s1_st_data;
    wire [35:0] s1_op_a, s1_op_b;
    wire        s1_dst_v, s1_src1_v, s1_src2_v, s1_src3_v;
    wire [4:0]  s1_dst, s1_src1, s1_src2, s1_src3;

    vtx1_slot slot1 (
        .word(w1), .regs_flat(regs_flat), .pc_self(p1), .pc_next(p2),
        .traps(s1_traps), .is_ld(s1_is_ld), .is_st(s1_is_st),
        .is_div(s1_is_div), .is_mod(), .is_stack(s1_is_stack),
        .is_ctrl(s1_is_ctrl), .is_halt(s1_is_halt),
        .redirect(s1_redirect), .ctrl_target(s1_ctrl_target),
        .wr_en(s1_wr_en), .wr_idx(s1_wr_idx), .wr_data(s1_wr_data),
        .ld_rd(s1_ld_rd), .ea_bin(s1_ea_bin), .st_data(s1_st_data),
        .op_a(s1_op_a), .op_b(s1_op_b),
        .dst_v(s1_dst_v), .dst(s1_dst),
        .src1_v(s1_src1_v), .src2_v(s1_src2_v), .src3_v(s1_src3_v),
        .src1(s1_src1), .src2(s1_src2), .src3(s1_src3)
    );

    wire        s2_traps, s2_is_ld, s2_is_st, s2_is_div, s2_is_stack, s2_is_ctrl, s2_is_halt, s2_redirect;
    wire [35:0] s2_ctrl_target;
    wire        s2_wr_en;
    wire [4:0]  s2_wr_idx;
    wire [35:0] s2_wr_data;
    wire [4:0]  s2_ld_rd;
    wire [31:0] s2_ea_bin;
    wire [35:0] s2_st_data;
    wire [35:0] s2_op_a, s2_op_b;
    wire        s2_dst_v, s2_src1_v, s2_src2_v, s2_src3_v;
    wire [4:0]  s2_dst, s2_src1, s2_src2, s2_src3;

    vtx1_slot slot2 (
        .word(w2), .regs_flat(regs_flat), .pc_self(p2), .pc_next(p3),
        .traps(s2_traps), .is_ld(s2_is_ld), .is_st(s2_is_st),
        .is_div(s2_is_div), .is_mod(), .is_stack(s2_is_stack),
        .is_ctrl(s2_is_ctrl), .is_halt(s2_is_halt),
        .redirect(s2_redirect), .ctrl_target(s2_ctrl_target),
        .wr_en(s2_wr_en), .wr_idx(s2_wr_idx), .wr_data(s2_wr_data),
        .ld_rd(s2_ld_rd), .ea_bin(s2_ea_bin), .st_data(s2_st_data),
        .op_a(s2_op_a), .op_b(s2_op_b),
        .dst_v(s2_dst_v), .dst(s2_dst),
        .src1_v(s2_src1_v), .src2_v(s2_src2_v), .src3_v(s2_src3_v),
        .src1(s2_src1), .src2(s2_src2), .src3(s2_src3)
    );

    // ------------------------------------------------------------------
    // Issue logic: the dynamic packing-rule check
    // ------------------------------------------------------------------

    // Does a later slot depend on an earlier slot's destination?
    // (RAW on any source, or WAW on the destination; TZ never conflicts.)
    function dep(
        input av,  input [4:0] a,                    // earlier slot's dest
        input s1v, input [4:0] s1,                   // later slot's sources
        input s2v, input [4:0] s2,
        input s3v, input [4:0] s3,
        input dv,  input [4:0] d                     // later slot's dest
    );
        dep = av && (a < 5'd18) && (
              (s1v && s1 == a) || (s2v && s2 == a) || (s3v && s3 == a) ||
              (dv  && d  == a));
    endfunction

    wire grp_end0 = s0_traps || s0_is_ctrl || s0_is_halt || s0_is_div || s0_is_stack;
    wire mem0     = s0_is_ld || s0_is_st;
    wire h01 = dep(s0_dst_v, s0_dst, s1_src1_v, s1_src1, s1_src2_v, s1_src2,
                   s1_src3_v, s1_src3, s1_dst_v, s1_dst);

    // A halt may not follow a memory op in the same group (the LD/ST
    // states must run to completion before the core stops).
    wire issue1 = avail12 && !grp_end0 && !s1_traps && !s1_is_div && !s1_is_stack && !h01
               && !(mem0 && (s1_is_ld || s1_is_st || s1_is_halt));

    wire grp_end1 = s1_is_ctrl || s1_is_halt;
    wire mem01    = mem0 || s1_is_ld || s1_is_st;
    wire h02 = dep(s0_dst_v, s0_dst, s2_src1_v, s2_src1, s2_src2_v, s2_src2,
                   s2_src3_v, s2_src3, s2_dst_v, s2_dst);
    wire h12 = dep(s1_dst_v, s1_dst, s2_src1_v, s2_src1, s2_src2_v, s2_src2,
                   s2_src3_v, s2_src3, s2_dst_v, s2_dst);

    wire issue2 = issue1 && !grp_end1 && !s2_traps && !s2_is_div && !s2_is_stack && !h02 && !h12
               && !(mem01 && (s2_is_ld || s2_is_st || s2_is_halt));

    wire [2:0] n_issue = 3'd1 + {2'd0, issue1} + {2'd0, issue2};

    // Next PC: the issued control slot's target, or the fall-through past
    // the last issued word. (At most one control op issues per group.)
    reg [35:0] nxt_pc;
    always @* begin
        if (!s0_traps && s0_redirect) nxt_pc = s0_ctrl_target;
        else if (issue1 && s1_redirect) nxt_pc = s1_ctrl_target;
        else if (issue2 && s2_redirect) nxt_pc = s2_ctrl_target;
        else if (issue2) nxt_pc = p3;
        else if (issue1) nxt_pc = p2;
        else nxt_pc = p1;
    end
    wire [31:0] nxt_pc_bin = t2b(nxt_pc);

    // ------------------------------------------------------------------
    // Control state machine
    // ------------------------------------------------------------------
    localparam S_FETCH = 3'd0;
    localparam S_EXEC  = 3'd1;
    localparam S_LOAD  = 3'd2;
    localparam S_STORE = 3'd3;
    localparam S_HALT  = 3'd4;
    localparam S_DIV   = 3'd5;  // waiting on the iterative divider

    reg [2:0]  state;
    reg [4:0]  ld_rd;      // writeback destination for LOAD and DIV     // destination latched for the load writeback cycle
    reg [31:0] mem_ea;    // effective address latched for LOAD/STORE
    reg [35:0] st_data_r; // store data latched at issue

    // Instruction-TCM port control (combinational): fetch the next group
    // during EXEC (fast path - next PC known this cycle) or the current
    // PC during FETCH (after reset, loads, stores and bus fetches).
    // Reads the TCM ignores are harmless, so EXEC may issue speculatively
    // even when it ends up in LOAD/STORE/HALT instead.
    always @* begin
        imem_en   = 1'b0;
        imem_addr = pc_bin;
        case (state)
            S_FETCH: begin imem_en = 1'b1; imem_addr = pc_bin; end
            S_EXEC:  begin imem_en = 1'b1; imem_addr = nxt_pc_bin; end
            default: ;
        endcase
    end

    // Memory port control (combinational); the request is held until ready.
    // A fetch goes to the bus only when the PC misses the TCM window.
    always @* begin
        mem_req   = 1'b0;
        mem_wr    = 1'b0;
        mem_addr  = pc_bin;
        mem_wdata = st_data_r;
        case (state)
            S_FETCH: begin
                mem_req  = !imem_hit;
                mem_addr = pc_bin;
            end
            S_LOAD: begin
                mem_req  = 1'b1;
                mem_addr = mem_ea;
            end
            S_STORE: begin
                mem_req  = 1'b1;
                mem_wr   = 1'b1;
                mem_addr = mem_ea;
            end
            default: ;
        endcase
    end

    // Iterative divider (DIV/MOD): started from EXEC, result written
    // back in S_DIV. Division by zero traps before start (vtxsim errors).
    wire div_bz    = (s0_op_b == ZERO_WORD);
    wire div_start = (state == S_EXEC) && !s0_traps && s0_is_div && !div_bz;
    wire        div_busy, div_done;
    wire [35:0] div_result;

    vtx1_div u_div (
        .clk(clk), .rst_n(rst_n),
        .start(div_start),
        .is_mod(s0_is_mod),
        .a(s0_op_a), .b(s0_op_b),
        .busy(div_busy), .done(div_done), .result(div_result)
    );

    task do_trap;
        begin
            trap         <= 1'b1;
            trap_pc      <= pc_bin;
            state        <= S_HALT;
            ir_from_imem <= 1'b0;
        end
    endtask

    // Commit actions shared by all three slots
    task commit_slot(
        input        c_wr_en,
        input [4:0]  c_wr_idx,
        input [35:0] c_wr_data,
        input        c_is_ld,
        input [4:0]  c_ld_rd,
        input        c_is_st,
        input [31:0] c_ea_bin,
        input [35:0] c_st_data,
        input        c_is_halt
    );
        begin
            if (c_wr_en) writereg(c_wr_idx, c_wr_data);
            // LD/ST keep the default ir_from_imem: the next group was
            // fetched during this EXEC and the TCM output registers hold
            // it through the memory state (imem_en is low there), so the
            // core resumes in EXEC with no refetch cycle.
            if (c_is_ld) begin
                ld_rd        <= c_ld_rd;
                mem_ea       <= c_ea_bin;
                state        <= S_LOAD;
            end
            if (c_is_st) begin
                mem_ea       <= c_ea_bin;
                st_data_r    <= c_st_data;
                state        <= S_STORE;
            end
            if (c_is_halt) begin
                halted       <= 1'b1;
                state        <= S_HALT;
                ir_from_imem <= 1'b0;
            end
        end
    endtask

    integer ri;
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            state        <= S_FETCH;
            pc           <= RESET_PC;
            ir           <= ZERO_WORD;
            ir_from_imem <= 1'b0;
            halted       <= 1'b0;
            trap         <= 1'b0;
            trap_pc      <= 32'h0;
            retired      <= 32'h0;
            cycles       <= 32'h0;
            st_data_r    <= ZERO_WORD;
            for (ri = 0; ri < 18; ri = ri + 1)
                regs[ri] <= ZERO_WORD;
        end else begin
            if (state != S_HALT)
                cycles <= cycles + 1;
            case (state)
                S_FETCH:
                    if (imem_hit) begin
                        // TCM fetch: the 3-word group read this cycle is in
                        // imem_rdata0..2 next cycle; execute from there.
                        state        <= S_EXEC;
                        ir_from_imem <= 1'b1;
                    end else if (mem_ready) begin
                        ir           <= mem_rdata;
                        ir_from_imem <= 1'b0;
                        state        <= S_EXEC;
                    end

                S_EXEC: begin
                    retired <= retired + {29'd0, n_issue};
                    pc      <= nxt_pc;  // control flow resolved combinationally
                    // Default next state: when the next PC hits the TCM, the
                    // group fetch issued this cycle lands in imem_rdata0..2
                    // and we stay in EXEC; otherwise fetch over the bus.
                    // LD/ST/halt commits below override this.
                    if (imem_hit) begin
                        state        <= S_EXEC;
                        ir_from_imem <= 1'b1;
                    end else begin
                        state        <= S_FETCH;
                        ir_from_imem <= 1'b0;
                    end

                    if (s0_traps) do_trap;
                    else if (s0_is_div) begin
                        if (div_bz) do_trap; // division by zero
                        else begin
                            ld_rd <= s0_ld_rd; // rd (reg1)
                            state <= S_DIV;    // default ir_from_imem kept
                        end
                    end
                    else begin
                        commit_slot(s0_wr_en, s0_wr_idx, s0_wr_data,
                                    s0_is_ld, s0_ld_rd,
                                    s0_is_st, s0_ea_bin, s0_st_data,
                                    s0_is_halt);
                        if (issue1)
                            commit_slot(s1_wr_en, s1_wr_idx, s1_wr_data,
                                        s1_is_ld, s1_ld_rd,
                                        s1_is_st, s1_ea_bin, s1_st_data,
                                        s1_is_halt);
                        if (issue2)
                            commit_slot(s2_wr_en, s2_wr_idx, s2_wr_data,
                                        s2_is_ld, s2_ld_rd,
                                        s2_is_st, s2_ea_bin, s2_st_data,
                                        s2_is_halt);
                    end
                end

                S_LOAD: if (mem_ready) begin
                    writereg(ld_rd, mem_rdata);
                    state <= ir_from_imem ? S_EXEC : S_FETCH;
                end

                S_STORE: if (mem_ready)
                    state <= ir_from_imem ? S_EXEC : S_FETCH;

                S_DIV: if (div_done) begin
                    writereg(ld_rd, div_result);
                    state <= ir_from_imem ? S_EXEC : S_FETCH;
                end

                S_HALT: ; // stay halted

            endcase
        end
    end

endmodule

`endif // VTX1_CORE_V
