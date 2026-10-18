`timescale 1ns / 1ps
// MMIO wrapper for the ternary-weight MAC array (vtx1_tmac).
//
// Register map (word offsets inside the peripheral's 64K window; the CPU
// reaches offset k at the negative word address -(0x60000 + k), per the
// ISA rule that negative addresses are memory-mapped I/O):
//
//   0x00..0x0F  write: activation word for lane 0..15
//   0x10        read:  the 18-trit accumulator
//               write: clear the accumulator (data ignored)
//   0x11        write: one weight word = one MAC step - the low 16 trits
//               are the 16 lane weights (lane i at weight 3^i), applied
//               against the loaded activations in a single cycle
//
// A 16x16 matrix-vector product is therefore: 16 activation writes, then
// per output row {clear, one weight-word write, one accumulator read}.

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
    output reg         ready
);

    reg [575:0] acts;
    wire [35:0] acc;

    wire access = req && !ready;
    wire w_act  = access && wr && (offset[15:4] == 12'd0); // 0x00..0x0F
    wire w_clr  = access && wr && (offset == 16'h0010);
    wire w_step = access && wr && (offset == 16'h0011);

    vtx1_tmac u_tmac (
        .clk(clk), .rst_n(rst_n),
        .clear(w_clr),
        .step(w_step),
        .weights(wdata[31:0]),   // the word's low 16 trits
        .acts(acts),
        .acc(acc)
    );

    always @(posedge clk or negedge rst_n) begin
        if (!rst_n) begin
            ready <= 1'b0;
            rdata <= 36'h555555555;
            acts  <= {16{36'h555555555}};
        end else begin
            ready <= req && !ready;
            if (w_act)
                acts[36*offset[3:0] +: 36] <= wdata;
            if (access && !wr)
                rdata <= (offset == 16'h0010) ? acc : 36'h555555555;
        end
    end

endmodule

`endif // VTX1_TMAC_MMIO_V
