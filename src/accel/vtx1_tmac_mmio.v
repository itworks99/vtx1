`timescale 1ns / 1ps
// MMIO wrapper + DMA engine for the ternary-weight MAC array (vtx1_tmac).
//
// Register map (word offsets inside the peripheral's 64K window; the CPU
// reaches offset k at the negative word address -(0x60000 + k), per the
// ISA rule that negative addresses are memory-mapped I/O):
//
//   0x00..0x0F  write: activation word for lane 0..15 (direct mode)
//   0x10        read:  the 18-trit accumulator
//               write: clear the accumulator (data ignored)
//   0x11        write: one weight word = one MAC step (direct mode)
//   0x12        write: ACT_BASE - scratchpad word address of the 16
//                      activation words (VALUE of the written trit word)
//   0x13        write: W_BASE - base of the packed weight rows (value)
//   0x14        write: R_BASE - base for the result words (value)
//   0x15        write: row count (value), and START the DMA run
//               read:  status (+1 while busy, 0 when idle/done)
//
// The DMA engine streams from the data scratchpad: 16 activation reads
// at one word per cycle, then 2 cycles per output row (one packed-weight
// read, one result write of the row's 16-lane dot product). A 16xN
// matrix-vector product is a 4-register descriptor and a status poll.
// Results are written INTO THE SCRATCHPAD (see its ownership caveat).

`ifndef VTX1_TMAC_MMIO_V
`define VTX1_TMAC_MMIO_V

module vtx1_tmac_mmio (
    input  wire        clk,
    input  wire        rst_n,

    // MMIO slave port (bus matrix handshake)
    input  wire        req,
    input  wire        wr,
    input  wire [15:0] offset,   // word offset inside the window
    input  wire [35:0] wdata,
    output reg  [35:0] rdata,
    output reg         ready,

    // Scratchpad DMA port
    output wire        sp_rd_en,
    output wire [31:0] sp_rd_addr,
    input  wire [35:0] sp_rd_data,
    output wire        sp_wr_en,
    output wire [31:0] sp_wr_addr,
    output wire [35:0] sp_wr_data,
    input  wire        sp_wr_stall
);

    localparam [35:0] ZERO_WORD = 36'h555555555;
    localparam [35:0] ONE_WORD  = 36'h555555556;

    // Trit word -> signed binary value (x3 adds only; for the descriptor
    // registers, which the CPU writes as trit words)
    function signed [35:0] t2b(input [35:0] w);
        integer i;
        reg signed [2:0] t;
        begin
            t2b = 36'sd0;
            for (i = 17; i >= 0; i = i - 1) begin
                case (w[2*i+1 -: 2])
                    2'b00:   t = -3'sd1;
                    2'b10:   t = 3'sd1;
                    default: t = 3'sd0;
                endcase
                t2b = t2b * 3 + t;
            end
        end
    endfunction

    // ------------------------------------------------------------------
    // DMA engine state
    // ------------------------------------------------------------------
    localparam D_IDLE = 2'd0;
    localparam D_ACT  = 2'd1;   // stream the 16 activation words
    localparam D_W    = 2'd2;   // present a weight-row address
    localparam D_RES  = 2'd3;   // weight word in sp_rd_data: write result

    reg [1:0]  dst;
    reg [31:0] act_base, w_base, r_base;
    reg [5:0]  rows_left;
    reg [31:0] row_addr, res_addr;
    reg [4:0]  ai;              // activation stream index (0..16)
    wire busy = (dst != D_IDLE);

    // ------------------------------------------------------------------
    // The MAC array (shared by direct MMIO mode and the DMA engine).
    // During D_RES the scratchpad read register holds the packed weight
    // row (no new read is issued in that state), so it feeds the tree
    // directly; sum16 is the row's complete 16-lane dot product.
    // ------------------------------------------------------------------
    reg [575:0] acts;
    wire [35:0] acc, sum16;

    // The written word's VALUE (for the descriptor registers)
    wire signed [35:0] wval = t2b(wdata);

    wire access = req && !ready;
    wire w_act  = access && wr && (offset[15:4] == 12'd0); // 0x00..0x0F
    wire w_clr  = access && wr && (offset == 16'h0010);
    wire w_step = access && wr && (offset == 16'h0011) && !busy;

    vtx1_tmac u_tmac (
        .clk(clk), .rst_n(rst_n),
        .clear(w_clr),
        .step(w_step),
        .weights((dst == D_RES) ? sp_rd_data[31:0] : wdata[31:0]),
        .acts(acts),
        .acc(acc), .sum16(sum16)
    );

    // ------------------------------------------------------------------
    // Scratchpad streaming
    // ------------------------------------------------------------------
    assign sp_rd_en   = (dst == D_ACT && ai != 5'd16) || (dst == D_W);
    assign sp_rd_addr = (dst == D_ACT) ? act_base + {27'd0, ai}
                                       : row_addr;
    assign sp_wr_en   = (dst == D_RES);
    assign sp_wr_addr = res_addr;
    assign sp_wr_data = sum16;

    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            ready     <= 1'b0;
            rdata     <= ZERO_WORD;
            acts      <= {16{36'h555555555}};
            dst       <= D_IDLE;
            act_base  <= 32'd0;
            w_base    <= 32'd0;
            r_base    <= 32'd0;
            rows_left <= 6'd0;
            row_addr  <= 32'd0;
            res_addr  <= 32'd0;
            ai        <= 5'd0;
        end else begin
            ready <= req && !ready;

            // MMIO register writes
            if (w_act)
                acts[36*offset[3:0] +: 36] <= wdata;
            if (access && wr && (offset == 16'h0012))
                act_base <= wval[31:0];
            if (access && wr && (offset == 16'h0013))
                w_base <= wval[31:0];
            if (access && wr && (offset == 16'h0014))
                r_base <= wval[31:0];
            if (access && wr && (offset == 16'h0015) && !busy) begin
                rows_left <= wval[5:0];
                dst       <= D_ACT;
                ai        <= 5'd0;
            end

            // MMIO reads
            if (access && !wr)
                rdata <= (offset == 16'h0010) ? acc :
                         (offset == 16'h0015) ? (busy ? ONE_WORD : ZERO_WORD) :
                                                ZERO_WORD;

            // DMA engine
            case (dst)
                D_ACT: begin
                    // the address for word ai goes out this cycle; the
                    // data for word ai-1 arrives this cycle
                    if (ai != 5'd0)
                        acts[36*(ai - 5'd1) +: 36] <= sp_rd_data;
                    if (ai == 5'd16) begin
                        row_addr <= w_base;
                        res_addr <= r_base;
                        dst      <= (rows_left == 6'd0) ? D_IDLE : D_W;
                    end else
                        ai <= ai + 5'd1;
                end

                D_W: dst <= D_RES;

                D_RES: if (!sp_wr_stall) begin
                    rows_left <= rows_left - 6'd1;
                    row_addr  <= row_addr + 32'd1;
                    res_addr  <= res_addr + 32'd1;
                    dst       <= (rows_left == 6'd1) ? D_IDLE : D_W;
                end

                default: ;
            endcase
        end
    end

endmodule

`endif // VTX1_TMAC_MMIO_V
