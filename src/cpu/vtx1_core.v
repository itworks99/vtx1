`timescale 1ns / 1ps
// VTX1 Trit-Native CPU Core (v2: native balanced-ternary datapath)
// Part of the VTX1 Ternary System-on-Chip
//
// A clean core built around vtx1_decoder and vtx1_alu: hardwired control,
// no microcode, and NO binary<->ternary conversion in the datapath:
//
//   - registers, the PC, immediates, effective addresses and all ALU
//     results are trit words end to end;
//   - immediates come out of the decoder as trit fields placed in
//     zero-trit-padded words (padding with zero trits preserves any
//     balanced value - no sign extension exists or is needed);
//   - the only trit->binary conversion is t2b() toward the memory address
//     bus, which is a chain of x3 adds (synthesizes cleanly); there is no
//     binary->trit conversion anywhere (the b2t divider chains of the v1
//     core are gone - they were a synthesis blocker).
//
// Fetch model (a): one 18-trit instruction word per fetch, PC steps by one
// word; a VLIW bundle executes as its three slots in sequence.
//
// Instruction fetch has a fast path through a tightly-coupled instruction
// memory (vtx1_tcm): during EXEC the next PC is known combinationally, so
// the core issues the TCM read for it in the same cycle and executes the
// next instruction directly from imem_rdata - straight-line TCM code runs
// at CPI 1 with no fetch state at all. Fetches outside the TCM window
// (imem_hit low) fall back to the bus port exactly as before; tying
// imem_hit to 0 gives the original bus-only core.
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
    // combinationally and delivers the word in imem_rdata one clock
    // later. Tie imem_hit low to fetch everything over the bus port.
    output reg         imem_en,
    output reg  [31:0] imem_addr,
    input  wire [35:0] imem_rdata,
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

    function [35:0] regword(input [4:0] idx);
        regword = (idx >= 5'd18) ? ZERO_WORD : regs[idx];
    endfunction

    task writereg(input [4:0] idx, input [35:0] w);
        if (idx < 5'd18)
            regs[idx] <= w;
    endtask

    localparam REG_TA = 5'd8; // link register (provisional ABI)

    // ------------------------------------------------------------------
    // Decoder (combinational on the current instruction word)
    // ------------------------------------------------------------------
    // On the fast path the instruction is consumed directly from the TCM
    // read register (ir_from_imem = 1); otherwise it is the latched ir.
    // EXEC always latches the executed word into ir so it stays stable
    // through the LOAD/STORE states (mem_wdata reads a decoder field).
    reg [35:0] ir;
    reg        ir_from_imem;
    wire [35:0] iword = ir_from_imem ? imem_rdata : ir;

    wire        dec_valid;
    wire [6:0]  opcode;
    wire [2:0]  format, category;
    wire        multicycle;
    wire [4:0]  reg1, reg2, reg3;
    wire        uses_imm;
    wire signed [31:0] imm;
    wire [1:0]  mem_mode;
    wire [35:0] imm_trits;

    vtx1_decoder dec (
        .word(iword), .valid(dec_valid), .opcode(opcode),
        .format(format), .category(category), .multicycle(multicycle),
        .reg1(reg1), .reg2(reg2), .reg3(reg3),
        .uses_imm(uses_imm), .imm(imm), .mem_mode(mem_mode),
        .imm_trits(imm_trits), .word_value()
    );

    // ------------------------------------------------------------------
    // Native-trit datapath
    // ------------------------------------------------------------------

    // Program counter as a trit word; binary only at the address bus.
    reg  [35:0] pc;
    wire [35:0] pc_plus1;
    vtx1_trit_adder u_pcinc (.a(pc), .b(ONE_WORD), .negate_b(1'b0), .sum(pc_plus1), .cout());

    wire [35:0] branch_target;
    vtx1_trit_adder u_btgt (.a(pc), .b(imm_trits), .negate_b(1'b0), .sum(branch_target), .cout());

    wire [31:0] pc_bin = t2b(pc);
    assign debug_pc = pc_bin;

    // ALU operand selection
    wire [35:0] alu_a =
        (format == `VTX1_FMT_R1) ? regword(reg1) : regword(reg2);
    wire [35:0] alu_b =
        (format == `VTX1_FMT_R1) ? ONE_WORD :
        uses_imm                 ? imm_trits : regword(reg3);

    wire [35:0] alu_result;
    wire        alu_unsup;
    vtx1_alu u_alu (.a(alu_a), .b(alu_b), .opcode(opcode),
                    .result(alu_result), .unsupported(alu_unsup));

    // Branch compare: sign of (rs1 - rs2)
    wire [35:0] br_diff;
    wire [1:0]  br_cout;
    vtx1_trit_adder u_brsub (.a(regword(reg1)), .b(regword(reg2)),
                             .negate_b(1'b1), .sum(br_diff), .cout(br_cout));
    // Overflow-correct compare: a nonzero carry trit is the true sign.
    wire [1:0] br_sign = (br_cout != 2'b01) ? br_cout : word_sign(br_diff);

    reg branch_take;
    always @* begin
        branch_take = 1'b0;
        case (opcode)
            `VTX1_OPV_BEQ: branch_take = (regword(reg1) == regword(reg2));
            `VTX1_OPV_BNE: branch_take = (regword(reg1) != regword(reg2));
            `VTX1_OPV_BLT: branch_take = (br_sign == 2'b00);
            `VTX1_OPV_BGE: branch_take = (br_sign != 2'b00);
            `VTX1_OPV_BGT: branch_take = (br_sign == 2'b10);
            `VTX1_OPV_BLE: branch_take = (br_sign != 2'b10);
            default: ;
        endcase
    end

    // Next PC, resolved combinationally during EXEC (control flow
    // included) so the TCM fetch of the next instruction can be issued
    // in the same cycle.
    reg [35:0] nxt_pc;
    always @* begin
        nxt_pc = pc_plus1;
        case (format)
            `VTX1_FMT_NONE:   if (opcode == `VTX1_OPV_RET)
                                  nxt_pc = regword(REG_TA);
            `VTX1_FMT_R1:     if (opcode == `VTX1_OPV_JR ||
                                  opcode == `VTX1_OPV_JALR)
                                  nxt_pc = regword(reg1);
            `VTX1_FMT_BRANCH: if (branch_take) nxt_pc = branch_target;
            `VTX1_FMT_JUMP:   nxt_pc = imm_trits; // absolute target as trits
            default: ;
        endcase
    end
    wire [31:0] nxt_pc_bin = t2b(nxt_pc);

    // Effective address (trit add; binary only at the port)
    wire [35:0] ea_sum;
    vtx1_trit_adder u_ea (
        .a(regword(reg2)),
        .b((mem_mode == 2'd1) ? regword(reg3) : imm_trits),
        .negate_b(1'b0), .sum(ea_sum), .cout());
    wire [35:0] ea_trits = (mem_mode == 2'd0) ? imm_trits : ea_sum;
    wire [31:0] ea_bin   = t2b(ea_trits);

    // ------------------------------------------------------------------
    // Control state machine
    // ------------------------------------------------------------------
    localparam S_FETCH = 3'd0;
    localparam S_EXEC  = 3'd1;
    localparam S_LOAD  = 3'd2;
    localparam S_STORE = 3'd3;
    localparam S_HALT  = 3'd4;

    reg [2:0]  state;
    reg [4:0]  ld_rd;  // destination latched for the load writeback cycle
    reg [31:0] mem_ea; // effective address latched for LOAD/STORE

    // Instruction-TCM port control (combinational): fetch the next
    // instruction during EXEC (fast path - next PC known this cycle) or
    // the current PC during FETCH (after reset, loads, stores and bus
    // fetches). Reads the TCM ignores are harmless, so EXEC may issue
    // speculatively even when it ends up in LOAD/STORE/HALT instead.
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
        mem_wdata = regword(reg1);
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

    task do_trap;
        begin
            trap         <= 1'b1;
            trap_pc      <= pc_bin;
            state        <= S_HALT;
            ir_from_imem <= 1'b0;
        end
    endtask

    integer ri;
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            state   <= S_FETCH;
            pc      <= RESET_PC;
            ir      <= ZERO_WORD;
            ir_from_imem <= 1'b0;
            halted  <= 1'b0;
            trap    <= 1'b0;
            trap_pc <= 32'h0;
            retired <= 32'h0;
            cycles  <= 32'h0;
            for (ri = 0; ri < 18; ri = ri + 1)
                regs[ri] <= ZERO_WORD;
        end else begin
            if (state != S_HALT)
                cycles <= cycles + 1;
            case (state)
                S_FETCH:
                    if (imem_hit) begin
                        // TCM fetch: the word read this cycle is in
                        // imem_rdata next cycle; execute it from there.
                        state        <= S_EXEC;
                        ir_from_imem <= 1'b1;
                    end else if (mem_ready) begin
                        ir           <= mem_rdata;
                        ir_from_imem <= 1'b0;
                        state        <= S_EXEC;
                    end

                S_EXEC: begin
                    retired <= retired + 1;
                    pc      <= nxt_pc;   // control flow resolved combinationally
                    ir      <= iword;    // keep the word stable for LOAD/STORE
                    // Default next state: when the next PC hits the TCM, the
                    // fetch issued this cycle lands in imem_rdata and we stay
                    // in EXEC (CPI 1); otherwise fetch over the bus.
                    if (imem_hit) begin
                        state        <= S_EXEC;
                        ir_from_imem <= 1'b1;
                    end else begin
                        state        <= S_FETCH;
                        ir_from_imem <= 1'b0;
                    end

                    if (!dec_valid) do_trap;
                    else case (format)

                        `VTX1_FMT_NONE: case (opcode)
                            `VTX1_OPV_NOP: ;
                            `VTX1_OPV_WFI, `VTX1_OPV_HALT,
                            `VTX1_OPV_SYSCALL, `VTX1_OPV_BREAK:
                                begin
                                    halted       <= 1'b1;
                                    state        <= S_HALT;
                                    ir_from_imem <= 1'b0;
                                end
                            `VTX1_OPV_RET: ; // nxt_pc = TA
                            `VTX1_OPV_CACHE, `VTX1_OPV_FLUSH, `VTX1_OPV_MEMBAR: ;
                            default: do_trap;
                        endcase

                        `VTX1_FMT_R1: case (opcode)
                            `VTX1_OPV_INC, `VTX1_OPV_DEC:
                                writereg(reg1, alu_result);
                            `VTX1_OPV_JR: ; // nxt_pc = reg1
                            `VTX1_OPV_JALR:
                                writereg(REG_TA, pc_plus1); // nxt_pc = reg1
                            default: do_trap; // PUSH/POP: no stack ABI yet
                        endcase

                        `VTX1_FMT_R2, `VTX1_FMT_R3: begin
                            if (alu_unsup) do_trap; // DIV/MOD (FSM pending), SQRT, vec/FP
                            else writereg(reg1, alu_result);
                        end

                        `VTX1_FMT_MEM: case (opcode)
                            `VTX1_OPV_LD: begin
                                ld_rd        <= reg1;
                                mem_ea       <= ea_bin;
                                state        <= S_LOAD;
                                ir_from_imem <= 1'b0;
                            end
                            `VTX1_OPV_ST: begin
                                mem_ea       <= ea_bin;
                                state        <= S_STORE;
                                ir_from_imem <= 1'b0;
                            end
                            `VTX1_OPV_LEA: writereg(reg1, ea_trits);
                            default: do_trap; // VLD/VST/FLD/FST
                        endcase

                        `VTX1_FMT_BRANCH: ; // nxt_pc = taken ? target : pc+1

                        `VTX1_FMT_JUMP:
                            if (opcode == `VTX1_OPV_JAL || opcode == `VTX1_OPV_CALL)
                                writereg(REG_TA, pc_plus1); // nxt_pc = target

                        default: do_trap;
                    endcase
                end

                S_LOAD: if (mem_ready) begin
                    writereg(ld_rd, mem_rdata);
                    state <= S_FETCH;
                end

                S_STORE: if (mem_ready)
                    state <= S_FETCH;

                S_HALT: ; // stay halted

            endcase
        end
    end

endmodule

`endif // VTX1_CORE_V
