`timescale 1ns / 1ps
// RV32 Triple-Issue Core: the issue-width-matched binary control for the
// VTX1 cycle-level comparison.
//
// This is vtx1_core v3 with the trit datapath replaced by two's-complement
// RV32IM (see rv32_slot.v for the implemented subset): the same three-slot
// structure, the SAME dynamic issue rules (at most one memory operation
// per group, a control op or halt ends the group, no intra-group RAW/WAW),
// the same FSM (FETCH/EXEC/LOAD/STORE/HALT), the same banked TCM fast path
// delivering 3 consecutive instructions with the next PC resolved
// combinationally during EXEC, and the same bus fallback. A cycle
// comparison between this core and vtx1_core therefore isolates the ISA:
// microarchitecture, memory model and design effort are held equal.
//
// ECALL/EBREAK halt (the kernels' exit); retired counts instructions,
// cycles stops at halt - identical metric semantics to vtx1_core.

`ifndef RV32_CORE_V
`define RV32_CORE_V

module rv32_core (
    input  wire        clk,
    input  wire        rst_n,

    // Unified memory port with ready handshake (byte address, word access)
    output reg         mem_req,
    output reg         mem_wr,
    output reg  [31:0] mem_addr,
    output reg  [31:0] mem_wdata,
    input  wire [31:0] mem_rdata,
    input  wire        mem_ready,

    // Tightly-coupled instruction memory port (synchronous read, 3 words)
    output reg         imem_en,
    output reg  [31:0] imem_addr,
    input  wire [31:0] imem_rdata0,
    input  wire [31:0] imem_rdata1,
    input  wire [31:0] imem_rdata2,
    input  wire        imem_hit,

    // Execution state
    output reg         halted,
    output reg         trap,
    output reg  [31:0] trap_pc,

    // Metrics (same semantics as vtx1_core)
    output reg  [31:0] retired,
    output reg  [31:0] cycles,

    output wire [31:0] debug_pc
);

    parameter [31:0] RESET_PC = 32'h0000_0000;

    // ------------------------------------------------------------------
    // Register file: x0 reads zero, writes discarded
    // ------------------------------------------------------------------
    reg [31:0] regs [0:31];

    wire [1023:0] regs_flat;
    genvar gi;
    generate
        for (gi = 0; gi < 32; gi = gi + 1) begin : flat
            assign regs_flat[32*gi +: 32] = regs[gi];
        end
    endgenerate

    task writereg(input [4:0] idx, input [31:0] w);
        if (idx != 5'd0)
            regs[idx] <= w;
    endtask

    // ------------------------------------------------------------------
    // Program counter chain
    // ------------------------------------------------------------------
    reg  [31:0] pc;
    wire [31:0] p1 = pc + 32'd4;
    wire [31:0] p2 = pc + 32'd8;
    wire [31:0] p3 = pc + 32'd12;

    assign debug_pc = pc;

    // ------------------------------------------------------------------
    // Issue slots
    // ------------------------------------------------------------------
    reg  [31:0] ir;
    reg         ir_from_imem;

    wire [31:0] w0 = ir_from_imem ? imem_rdata0 : ir;
    wire [31:0] w1 = imem_rdata1;
    wire [31:0] w2 = imem_rdata2;
    wire        avail12 = ir_from_imem;

    wire        s0_traps, s0_is_ld, s0_is_st, s0_is_ctrl, s0_is_halt, s0_redirect;
    wire [31:0] s0_ctrl_target;
    wire        s0_wr_en;
    wire [4:0]  s0_wr_idx;
    wire [31:0] s0_wr_data;
    wire [4:0]  s0_ld_rd;
    wire [31:0] s0_ea, s0_st_data;
    wire        s0_dst_v, s0_src1_v, s0_src2_v;
    wire [4:0]  s0_dst, s0_src1, s0_src2;

    rv32_slot slot0 (
        .word(w0), .regs_flat(regs_flat), .pc_self(pc), .pc_next(p1),
        .traps(s0_traps), .is_ld(s0_is_ld), .is_st(s0_is_st),
        .is_ctrl(s0_is_ctrl), .is_halt(s0_is_halt),
        .redirect(s0_redirect), .ctrl_target(s0_ctrl_target),
        .wr_en(s0_wr_en), .wr_idx(s0_wr_idx), .wr_data(s0_wr_data),
        .ld_rd(s0_ld_rd), .ea(s0_ea), .st_data(s0_st_data),
        .dst_v(s0_dst_v), .dst(s0_dst),
        .src1_v(s0_src1_v), .src2_v(s0_src2_v),
        .src1(s0_src1), .src2(s0_src2)
    );

    wire        s1_traps, s1_is_ld, s1_is_st, s1_is_ctrl, s1_is_halt, s1_redirect;
    wire [31:0] s1_ctrl_target;
    wire        s1_wr_en;
    wire [4:0]  s1_wr_idx;
    wire [31:0] s1_wr_data;
    wire [4:0]  s1_ld_rd;
    wire [31:0] s1_ea, s1_st_data;
    wire        s1_dst_v, s1_src1_v, s1_src2_v;
    wire [4:0]  s1_dst, s1_src1, s1_src2;

    rv32_slot slot1 (
        .word(w1), .regs_flat(regs_flat), .pc_self(p1), .pc_next(p2),
        .traps(s1_traps), .is_ld(s1_is_ld), .is_st(s1_is_st),
        .is_ctrl(s1_is_ctrl), .is_halt(s1_is_halt),
        .redirect(s1_redirect), .ctrl_target(s1_ctrl_target),
        .wr_en(s1_wr_en), .wr_idx(s1_wr_idx), .wr_data(s1_wr_data),
        .ld_rd(s1_ld_rd), .ea(s1_ea), .st_data(s1_st_data),
        .dst_v(s1_dst_v), .dst(s1_dst),
        .src1_v(s1_src1_v), .src2_v(s1_src2_v),
        .src1(s1_src1), .src2(s1_src2)
    );

    wire        s2_traps, s2_is_ld, s2_is_st, s2_is_ctrl, s2_is_halt, s2_redirect;
    wire [31:0] s2_ctrl_target;
    wire        s2_wr_en;
    wire [4:0]  s2_wr_idx;
    wire [31:0] s2_wr_data;
    wire [4:0]  s2_ld_rd;
    wire [31:0] s2_ea, s2_st_data;
    wire        s2_dst_v, s2_src1_v, s2_src2_v;
    wire [4:0]  s2_dst, s2_src1, s2_src2;

    rv32_slot slot2 (
        .word(w2), .regs_flat(regs_flat), .pc_self(p2), .pc_next(p3),
        .traps(s2_traps), .is_ld(s2_is_ld), .is_st(s2_is_st),
        .is_ctrl(s2_is_ctrl), .is_halt(s2_is_halt),
        .redirect(s2_redirect), .ctrl_target(s2_ctrl_target),
        .wr_en(s2_wr_en), .wr_idx(s2_wr_idx), .wr_data(s2_wr_data),
        .ld_rd(s2_ld_rd), .ea(s2_ea), .st_data(s2_st_data),
        .dst_v(s2_dst_v), .dst(s2_dst),
        .src1_v(s2_src1_v), .src2_v(s2_src2_v),
        .src1(s2_src1), .src2(s2_src2)
    );

    // ------------------------------------------------------------------
    // Issue logic: identical rules to vtx1_core (x0 never conflicts)
    // ------------------------------------------------------------------
    function dep(
        input av,  input [4:0] a,
        input s1v, input [4:0] s1,
        input s2v, input [4:0] s2,
        input dv,  input [4:0] d
    );
        dep = av && (a != 5'd0) && (
              (s1v && s1 == a) || (s2v && s2 == a) ||
              (dv  && d  == a));
    endfunction

    wire grp_end0 = s0_traps || s0_is_ctrl || s0_is_halt;
    wire mem0     = s0_is_ld || s0_is_st;
    wire h01 = dep(s0_dst_v, s0_dst, s1_src1_v, s1_src1, s1_src2_v, s1_src2,
                   s1_dst_v, s1_dst);

    wire issue1 = avail12 && !grp_end0 && !s1_traps && !h01
               && !(mem0 && (s1_is_ld || s1_is_st || s1_is_halt));

    wire grp_end1 = s1_is_ctrl || s1_is_halt;
    wire mem01    = mem0 || s1_is_ld || s1_is_st;
    wire h02 = dep(s0_dst_v, s0_dst, s2_src1_v, s2_src1, s2_src2_v, s2_src2,
                   s2_dst_v, s2_dst);
    wire h12 = dep(s1_dst_v, s1_dst, s2_src1_v, s2_src1, s2_src2_v, s2_src2,
                   s2_dst_v, s2_dst);

    wire issue2 = issue1 && !grp_end1 && !s2_traps && !h02 && !h12
               && !(mem01 && (s2_is_ld || s2_is_st || s2_is_halt));

    wire [2:0] n_issue = 3'd1 + {2'd0, issue1} + {2'd0, issue2};

    reg [31:0] nxt_pc;
    always @* begin
        if (!s0_traps && s0_redirect) nxt_pc = s0_ctrl_target;
        else if (issue1 && s1_redirect) nxt_pc = s1_ctrl_target;
        else if (issue2 && s2_redirect) nxt_pc = s2_ctrl_target;
        else if (issue2) nxt_pc = p3;
        else if (issue1) nxt_pc = p2;
        else nxt_pc = p1;
    end

    // ------------------------------------------------------------------
    // Control state machine (identical to vtx1_core)
    // ------------------------------------------------------------------
    localparam S_FETCH = 3'd0;
    localparam S_EXEC  = 3'd1;
    localparam S_LOAD  = 3'd2;
    localparam S_STORE = 3'd3;
    localparam S_HALT  = 3'd4;

    reg [2:0]  state;
    reg [4:0]  ld_rd;
    reg [31:0] mem_ea;
    reg [31:0] st_data_r;

    always @* begin
        imem_en   = 1'b0;
        imem_addr = pc;
        case (state)
            S_FETCH: begin imem_en = 1'b1; imem_addr = pc; end
            S_EXEC:  begin imem_en = 1'b1; imem_addr = nxt_pc; end
            default: ;
        endcase
    end

    always @* begin
        mem_req   = 1'b0;
        mem_wr    = 1'b0;
        mem_addr  = pc;
        mem_wdata = st_data_r;
        case (state)
            S_FETCH: begin
                mem_req  = !imem_hit;
                mem_addr = pc;
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

    task do_trap;
        begin
            trap         <= 1'b1;
            trap_pc      <= pc;
            state        <= S_HALT;
            ir_from_imem <= 1'b0;
        end
    endtask

    task commit_slot(
        input        c_wr_en,
        input [4:0]  c_wr_idx,
        input [31:0] c_wr_data,
        input        c_is_ld,
        input [4:0]  c_ld_rd,
        input        c_is_st,
        input [31:0] c_ea,
        input [31:0] c_st_data,
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
                mem_ea       <= c_ea;
                state        <= S_LOAD;
            end
            if (c_is_st) begin
                mem_ea       <= c_ea;
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
            ir           <= 32'd0;
            ir_from_imem <= 1'b0;
            halted       <= 1'b0;
            trap         <= 1'b0;
            trap_pc      <= 32'h0;
            retired      <= 32'h0;
            cycles       <= 32'h0;
            st_data_r    <= 32'd0;
            for (ri = 0; ri < 32; ri = ri + 1)
                regs[ri] <= 32'd0;
        end else begin
            if (state != S_HALT)
                cycles <= cycles + 1;
            case (state)
                S_FETCH:
                    if (imem_hit) begin
                        state        <= S_EXEC;
                        ir_from_imem <= 1'b1;
                    end else if (mem_ready) begin
                        ir           <= mem_rdata;
                        ir_from_imem <= 1'b0;
                        state        <= S_EXEC;
                    end

                S_EXEC: begin
                    retired <= retired + {29'd0, n_issue};
                    pc      <= nxt_pc;
                    if (imem_hit) begin
                        state        <= S_EXEC;
                        ir_from_imem <= 1'b1;
                    end else begin
                        state        <= S_FETCH;
                        ir_from_imem <= 1'b0;
                    end

                    if (s0_traps) do_trap;
                    else begin
                        commit_slot(s0_wr_en, s0_wr_idx, s0_wr_data,
                                    s0_is_ld, s0_ld_rd,
                                    s0_is_st, s0_ea, s0_st_data,
                                    s0_is_halt);
                        if (issue1)
                            commit_slot(s1_wr_en, s1_wr_idx, s1_wr_data,
                                        s1_is_ld, s1_ld_rd,
                                        s1_is_st, s1_ea, s1_st_data,
                                        s1_is_halt);
                        if (issue2)
                            commit_slot(s2_wr_en, s2_wr_idx, s2_wr_data,
                                        s2_is_ld, s2_ld_rd,
                                        s2_is_st, s2_ea, s2_st_data,
                                        s2_is_halt);
                    end
                end

                S_LOAD: if (mem_ready) begin
                    writereg(ld_rd, mem_rdata);
                    state <= ir_from_imem ? S_EXEC : S_FETCH;
                end

                S_STORE: if (mem_ready)
                    state <= ir_from_imem ? S_EXEC : S_FETCH;

                S_HALT: ;

            endcase
        end
    end

endmodule

`endif // RV32_CORE_V
