`timescale 1ns / 1ps
// VTX1 Trit-Native CPU Core (v1: multi-cycle, word-serial)
// Part of the VTX1 Ternary System-on-Chip
//
// A clean core built around vtx1_decoder: hardwired control, no microcode.
// Fetch model (a): one 18-trit instruction word per fetch, PC steps by one
// word; a VLIW bundle executes as its three slots in sequence. (Parallel
// bundle issue is a later performance mode; the assembler's packing rules
// already guarantee bundle slots are independent.)
//
// Registers and memory hold trit-encoded words (36 bits, 2 bits/trit,
// canonical encoding). The v1 ALU converts trits -> signed binary ->
// trits inline; this is functionally identical to native ternary
// arithmetic (the ternary_arithmetic_unit in src/common does the same
// conversion internally) and is the part to swap for native trit logic
// when the datapath is optimized.
//
// Multi-cycle ISA ops: DIV/MOD execute here in one EXEC cycle using
// operators (fine for simulation; synthesis wants the iterative divider
// FSM noted in todo.md). SQRT, PUSH/POP (no stack ABI yet) and all
// vector/FP operations trap.
//
// Memory interface: one word-addressed port, synchronous 1-cycle RAM
// (read data valid the cycle after mem_req).

`ifndef VTX1_CORE_V
`define VTX1_CORE_V

`include "vtx1_isa_generated.vh"

module vtx1_core (
    input  wire        clk,
    input  wire        rst_n,

    // Unified word-addressed memory port (synchronous, always ready)
    output reg         mem_req,
    output reg         mem_wr,
    output reg  [31:0] mem_addr,
    output reg  [35:0] mem_wdata,
    input  wire [35:0] mem_rdata,

    // Execution state
    output reg         halted,
    output reg         trap,        // invalid word or unimplemented operation
    output reg  [31:0] trap_pc,

    // Metrics (cross-checked against the ISA simulator)
    output reg  [31:0] retired,     // instruction words executed
    output reg  [31:0] cycles
);

    // Canonical all-zero-trit word (18 x encoding 01)
    localparam [35:0] ZERO_WORD = 36'h555555555;

    // ------------------------------------------------------------------
    // Balanced ternary <-> binary value helpers
    // ------------------------------------------------------------------
    function signed [2:0] tv(input [1:0] pair);
        case (pair)
            2'b00:   tv = -3'sd1;
            2'b01:   tv = 3'sd0;
            2'b10:   tv = 3'sd1;
            default: tv = 3'sd0;
        endcase
    endfunction

    // trit word -> signed value
    function signed [35:0] t2b(input [35:0] w);
        integer i;
        begin
            t2b = 36'sd0;
            for (i = 0; i < 18; i = i + 1)
                t2b = t2b * 3 + tv(w[35-2*i -: 2]);
        end
    endfunction

    // signed value -> canonical trit word (wraps modulo 3^18, balanced)
    function [35:0] b2t(input signed [71:0] v);
        integer i;
        reg signed [71:0] x;
        reg signed [3:0] r;
        begin
            x = v;
            b2t = 36'h0;
            for (i = 17; i >= 0; i = i - 1) begin
                r = x % 3;
                x = x / 3;
                if (r == 2) begin r = -1; x = x + 1; end
                if (r == -2) begin r = 1; x = x - 1; end
                case (r)
                    -1:      b2t[35-2*i -: 2] = 2'b00;
                    0:       b2t[35-2*i -: 2] = 2'b01;
                    default: b2t[35-2*i -: 2] = 2'b10;
                endcase
            end
        end
    endfunction

    // Tritwise logic (balanced ternary conventions: AND=min, OR=max, XOR=product)
    function [1:0] trit_enc(input signed [2:0] t);
        case (t)
            -3'sd1:  trit_enc = 2'b00;
            3'sd0:   trit_enc = 2'b01;
            default: trit_enc = 2'b10;
        endcase
    endfunction

    function [35:0] tritwise_min(input [35:0] a, input [35:0] b);
        integer i;
        for (i = 0; i < 18; i = i + 1)
            tritwise_min[35-2*i -: 2] = trit_enc(
                (tv(a[35-2*i -: 2]) < tv(b[35-2*i -: 2])) ? tv(a[35-2*i -: 2]) : tv(b[35-2*i -: 2]));
    endfunction

    function [35:0] tritwise_max(input [35:0] a, input [35:0] b);
        integer i;
        for (i = 0; i < 18; i = i + 1)
            tritwise_max[35-2*i -: 2] = trit_enc(
                (tv(a[35-2*i -: 2]) > tv(b[35-2*i -: 2])) ? tv(a[35-2*i -: 2]) : tv(b[35-2*i -: 2]));
    endfunction

    function [35:0] tritwise_mul(input [35:0] a, input [35:0] b);
        integer i;
        for (i = 0; i < 18; i = i + 1)
            tritwise_mul[35-2*i -: 2] = trit_enc(tv(a[35-2*i -: 2]) * tv(b[35-2*i -: 2]));
    endfunction

    function [35:0] tritwise_neg(input [35:0] a);
        integer i;
        for (i = 0; i < 18; i = i + 1)
            tritwise_neg[35-2*i -: 2] = trit_enc(-tv(a[35-2*i -: 2]));
    endfunction

    // Trit shift/rotate: k trit positions (zero-trit fill for shifts)
    function [35:0] trit_shift(input [35:0] a, input [4:0] k,
                               input left, input rotate);
        integer i, src;
        begin
            for (i = 0; i < 18; i = i + 1) begin
                src = left ? i + k : i - k;
                if (src >= 0 && src < 18)
                    trit_shift[35-2*i -: 2] = a[35-2*src -: 2];
                else if (rotate)
                    trit_shift[35-2*i -: 2] = a[35-2*(((src % 18) + 18) % 18) -: 2];
                else
                    trit_shift[35-2*i -: 2] = 2'b01; // zero trit
            end
        end
    endfunction

    // ------------------------------------------------------------------
    // Register file: unified index space 0..18 (TZ = 18 reads zero,
    // writes discarded). Registers hold canonical trit-encoded words.
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
    // Decoder (combinational on the fetched word)
    // ------------------------------------------------------------------
    wire        dec_valid;
    wire [6:0]  opcode;
    wire [2:0]  format, category;
    wire        multicycle;
    wire [4:0]  reg1, reg2, reg3;
    wire        uses_imm;
    wire signed [31:0] imm;
    wire [1:0]  mem_mode;
    wire signed [35:0] word_value;

    vtx1_decoder dec (
        .word(mem_rdata), .valid(dec_valid), .opcode(opcode),
        .format(format), .category(category), .multicycle(multicycle),
        .reg1(reg1), .reg2(reg2), .reg3(reg3),
        .uses_imm(uses_imm), .imm(imm), .mem_mode(mem_mode),
        .word_value(word_value)
    );

    // Operand values
    wire signed [35:0] val1 = t2b(regword(reg1));
    wire signed [35:0] val2 = t2b(regword(reg2));
    wire signed [35:0] val3 = t2b(regword(reg3));
    wire signed [35:0] opb  = (format == `VTX1_FMT_R3 && uses_imm) ? imm : val3;

    // Memory effective address
    wire signed [35:0] ea =
        (mem_mode == 2'd0) ? imm :
        (mem_mode == 2'd1) ? val2 + val3 :
                             val2 + imm;

    // ------------------------------------------------------------------
    // Control state machine
    // ------------------------------------------------------------------
    localparam S_FETCH = 2'd0;
    localparam S_EXEC  = 2'd1;
    localparam S_LOAD  = 2'd2;
    localparam S_HALT  = 2'd3;

    reg [1:0]  state;
    reg [31:0] pc;
    reg [4:0]  ld_rd; // destination latched for the load writeback cycle

    // Memory port control (combinational)
    always @* begin
        mem_req   = 1'b0;
        mem_wr    = 1'b0;
        mem_addr  = pc;
        mem_wdata = regword(reg1);
        case (state)
            S_FETCH: begin
                mem_req  = 1'b1;
                mem_addr = pc;
            end
            S_EXEC: begin
                if (dec_valid && format == `VTX1_FMT_MEM) begin
                    if (opcode == 7'd17) begin // LD
                        mem_req  = 1'b1;
                        mem_addr = ea[31:0];
                    end else if (opcode == 7'd18) begin // ST
                        mem_req  = 1'b1;
                        mem_wr   = 1'b1;
                        mem_addr = ea[31:0];
                    end
                end
            end
            default: ;
        endcase
    end

    integer ri;
    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            state   <= S_FETCH;
            pc      <= 32'h1000;
            halted  <= 1'b0;
            trap    <= 1'b0;
            trap_pc <= 32'h0;
            retired <= 32'h0;
            cycles  <= 32'h0;
            for (ri = 0; ri < 18; ri = ri + 1)
                regs[ri] <= ZERO_WORD;
        end else begin
            cycles <= cycles + 1;
            case (state)
                S_FETCH: state <= S_EXEC;

                S_EXEC: begin
                    retired <= retired + 1;
                    pc      <= pc + 1; // default; control flow overrides
                    state   <= S_FETCH;

                    if (!dec_valid) begin
                        trap <= 1'b1; trap_pc <= pc; state <= S_HALT;
                    end else case (format)

                        `VTX1_FMT_NONE: case (opcode)
                            7'd40: ;                                   // NOP
                            7'd54, 7'd80, 7'd78, 7'd79:                // WFI/HALT/SYSCALL/BREAK
                                begin halted <= 1'b1; state <= S_HALT; end
                            7'd39: pc <= t2b(regword(REG_TA));         // RET
                            7'd75, 7'd76, 7'd77: ;                     // CACHE/FLUSH/MEMBAR
                            default:
                                begin trap <= 1'b1; trap_pc <= pc; state <= S_HALT; end
                        endcase

                        `VTX1_FMT_R1: case (opcode)
                            7'd14: writereg(reg1, b2t(val1 + 1));      // INC
                            7'd15: writereg(reg1, b2t(val1 - 1));      // DEC
                            7'd28: pc <= val1[31:0];                   // JR
                            7'd29: begin                               // JALR
                                writereg(REG_TA, b2t({4'h0, pc} + 1));
                                pc <= val1[31:0];
                            end
                            default: // PUSH/POP: no stack ABI yet
                                begin trap <= 1'b1; trap_pc <= pc; state <= S_HALT; end
                        endcase

                        `VTX1_FMT_R2: case (opcode)
                            7'd6:  writereg(reg1, tritwise_neg(regword(reg2)));  // NOT
                            7'd16: writereg(reg1, b2t(-val2));                   // NEG
                            7'd60: writereg(reg1, b2t(val2 < 0 ? -val2 : val2)); // ABS
                            default: // SQRT (FSM pending), vector/FP
                                begin trap <= 1'b1; trap_pc <= pc; state <= S_HALT; end
                        endcase

                        `VTX1_FMT_R3: case (opcode)
                            7'd1:  writereg(reg1, b2t(val2 + opb));             // ADD
                            7'd2:  writereg(reg1, b2t(val2 - opb));             // SUB
                            7'd3:  writereg(reg1, b2t(val2 * opb));             // MUL
                            7'd12: writereg(reg1, b2t(                          // CMP: sign(a-b)
                                       (val2 > opb) ? 72'sd1 :
                                       (val2 < opb) ? -72'sd1 : 72'sd0));
                            7'd13: writereg(reg1, b2t(                          // TEST: sign(a AND b)
                                       t2b(tritwise_min(regword(reg2),
                                           uses_imm ? b2t(imm) : regword(reg3)))));
                            7'd4:  writereg(reg1, tritwise_min(regword(reg2),   // AND
                                       uses_imm ? b2t(imm) : regword(reg3)));
                            7'd5:  writereg(reg1, tritwise_max(regword(reg2),   // OR
                                       uses_imm ? b2t(imm) : regword(reg3)));
                            7'd7:  writereg(reg1, tritwise_mul(regword(reg2),   // XOR
                                       uses_imm ? b2t(imm) : regword(reg3)));
                            7'd8:  writereg(reg1, trit_shift(regword(reg2), opb[4:0], 1'b1, 1'b0)); // SHL
                            7'd9:  writereg(reg1, trit_shift(regword(reg2), opb[4:0], 1'b0, 1'b0)); // SHR
                            7'd10: writereg(reg1, trit_shift(regword(reg2), opb[4:0], 1'b1, 1'b1)); // ROL
                            7'd11: writereg(reg1, trit_shift(regword(reg2), opb[4:0], 1'b0, 1'b1)); // ROR
                            // DIV/MOD: single-cycle here for simulation;
                            // synthesis needs the iterative divider FSM.
                            7'd55: if (opb == 0) begin trap <= 1'b1; trap_pc <= pc; state <= S_HALT; end
                                   else writereg(reg1, b2t(val2 / opb));        // DIV
                            7'd56: if (opb == 0) begin trap <= 1'b1; trap_pc <= pc; state <= S_HALT; end
                                   else writereg(reg1, b2t(val2 % opb));        // MOD
                            default: // vector ops
                                begin trap <= 1'b1; trap_pc <= pc; state <= S_HALT; end
                        endcase

                        `VTX1_FMT_MEM: case (opcode)
                            7'd17: begin                               // LD
                                ld_rd <= reg1;
                                state <= S_LOAD;
                            end
                            7'd18: ;                                   // ST (port handles it)
                            7'd23: writereg(reg1, b2t(ea));            // LEA
                            default: // VLD/VST/FLD/FST
                                begin trap <= 1'b1; trap_pc <= pc; state <= S_HALT; end
                        endcase

                        `VTX1_FMT_BRANCH: begin
                            case (opcode)
                                7'd30: if (regword(reg1) == regword(reg2)) pc <= pc + imm; // BEQ
                                7'd31: if (regword(reg1) != regword(reg2)) pc <= pc + imm; // BNE
                                7'd32: if (val1 <  val2) pc <= pc + imm;                   // BLT
                                7'd33: if (val1 >= val2) pc <= pc + imm;                   // BGE
                                7'd36: if (val1 >  val2) pc <= pc + imm;                   // BGT
                                7'd37: if (val1 <= val2) pc <= pc + imm;                   // BLE
                                default:
                                    begin trap <= 1'b1; trap_pc <= pc; state <= S_HALT; end
                            endcase
                        end

                        `VTX1_FMT_JUMP: begin
                            if (opcode == 7'd27 || opcode == 7'd38)    // JAL/CALL
                                writereg(REG_TA, b2t({4'h0, pc} + 1));
                            pc <= imm[31:0];                           // JMP/JAL/CALL
                        end

                        default:
                            begin trap <= 1'b1; trap_pc <= pc; state <= S_HALT; end
                    endcase
                end

                S_LOAD: begin
                    writereg(ld_rd, mem_rdata);
                    state <= S_FETCH;
                end

                S_HALT: ; // stay halted

            endcase
        end
    end

endmodule

`endif // VTX1_CORE_V
