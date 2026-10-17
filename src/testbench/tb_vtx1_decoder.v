`timescale 1ns / 1ps
// Testbench: trit-native decoder vs. real assembler output
//
// Loads build/hello_world.memh -- produced by the VTX1 assembler from
// examples/hello_world.asm (`task test-decoder` regenerates it) -- and
// verifies that vtx1_decoder decodes every instruction word to the expected
// operation, registers and immediates, and that every emitted word is a
// valid trit stream. This is the toolchain <-> RTL contract test: both sides
// are generated from src/assembler/isa/isa.go.

`include "vtx1_isa_generated.vh"

module tb_vtx1_decoder;

    reg [35:0] mem [0:16383];

    reg  [35:0] word;
    wire        valid;
    wire [6:0]  opcode;
    wire [2:0]  format;
    wire [2:0]  category;
    wire        multicycle;
    wire [4:0]  reg1, reg2, reg3;
    wire        uses_imm;
    wire signed [31:0] imm;
    wire [1:0]  mem_mode;
    wire signed [35:0] word_value;

    vtx1_decoder dut (
        .word(word), .valid(valid), .opcode(opcode), .format(format),
        .category(category), .multicycle(multicycle),
        .reg1(reg1), .reg2(reg2), .reg3(reg3),
        .uses_imm(uses_imm), .imm(imm), .mem_mode(mem_mode),
        .word_value(word_value)
    );

    integer errors;
    integer i;

    // hello_world.asm code section: word addresses 0x1000..0x100D
    localparam CODE_BASE  = 'h1000;
    localparam CODE_WORDS = 14;
    localparam DATA_BASE  = 'h100E; // "Hello, VTX1!\0" -- one word per char
    localparam DATA_WORDS = 13;

    // Expected opcode value per code word (from the ISA definition).
    reg [6:0] expected_opcode [0:CODE_WORDS-1];

    // Expected data values: "Hello, VTX1!" + NUL
    reg signed [35:0] expected_data [0:DATA_WORDS-1];

    task check(input [255:0] name, input [31:0] got, input [31:0] want);
        begin
            if (got !== want) begin
                $display("  FAIL: %0s = %0d, want %0d", name, $signed(got), $signed(want));
                errors = errors + 1;
            end
        end
    endtask

    initial begin
        errors = 0;
        $readmemh("build/hello_world.memh", mem);

        expected_opcode[0]  = 7'd23; // LEA T0, string
        expected_opcode[1]  = 7'd23; // LEA T1, 0x2000
        expected_opcode[2]  = 7'd23; // LEA T2, 0
        expected_opcode[3]  = 7'd17; // LD  T3, [T0]
        expected_opcode[4]  = 7'd30; // BEQ T3, 0, done
        expected_opcode[5]  = 7'd18; // ST  T3, [T1]
        expected_opcode[6]  = 7'd1;  // ADD T0, T0, 1   (VLIW slot 0)
        expected_opcode[7]  = 7'd1;  // ADD T1, T1, 1   (VLIW slot 1)
        expected_opcode[8]  = 7'd1;  // ADD T2, T2, 1   (VLIW slot 2)
        expected_opcode[9]  = 7'd26; // JMP loop
        expected_opcode[10] = 7'd18; // ST  T2, 0x2100
        expected_opcode[11] = 7'd23; // LEA T0, 1
        expected_opcode[12] = 7'd18; // ST  T0, 0x2104
        expected_opcode[13] = 7'd54; // WFI

        expected_data[0]  = "H"; expected_data[1]  = "e"; expected_data[2]  = "l";
        expected_data[3]  = "l"; expected_data[4]  = "o"; expected_data[5]  = ",";
        expected_data[6]  = " "; expected_data[7]  = "V"; expected_data[8]  = "T";
        expected_data[9]  = "X"; expected_data[10] = "1"; expected_data[11] = "!";
        expected_data[12] = 0;

        $display("=== VTX1 trit-native decoder vs assembler output ===");

        // Every code word must be a valid trit stream with the right opcode.
        for (i = 0; i < CODE_WORDS; i = i + 1) begin
            word = mem[CODE_BASE + i];
            #1;
            if (!valid) begin
                $display("  FAIL: code word %0d is not a valid trit stream (%09h)", i, word);
                errors = errors + 1;
            end
            if (opcode !== expected_opcode[i]) begin
                $display("  FAIL: code word %0d opcode = %0d, want %0d", i, opcode, expected_opcode[i]);
                errors = errors + 1;
            end
        end
        $display("Code words decoded: %0d", CODE_WORDS);

        // Spot-check full field decode of representative words.

        // word 0: LEA T0, string (direct address 0x100E)
        word = mem[CODE_BASE + 0]; #1;
        check("LEA fmt",      format,   `VTX1_FMT_MEM);
        check("LEA cat",      category, `VTX1_CAT_MEM);
        check("LEA reg1(T0)", reg1,     5'd0);
        check("LEA mode",     mem_mode, 2'd0);
        check("LEA imm",      imm,      32'h100E);

        // word 3: LD T3, [T0] (base+index with index TZ)
        word = mem[CODE_BASE + 3]; #1;
        check("LD[T0] mode",      mem_mode, 2'd1);
        check("LD[T0] reg1(T3)",  reg1,     5'd3);
        check("LD[T0] base(T0)",  reg2,     5'd0);
        check("LD[T0] index(TZ)", reg3,     5'd18);

        // word 4: BEQ T3, 0, done (PC-relative +6, rs2 = TZ)
        word = mem[CODE_BASE + 4]; #1;
        check("BEQ fmt",      format, `VTX1_FMT_BRANCH);
        check("BEQ cat",      category, `VTX1_CAT_CTRL);
        check("BEQ rs1(T3)",  reg1,   5'd3);
        check("BEQ rs2(TZ)",  reg2,   5'd18);
        check("BEQ offset",   imm,    32'sd6);

        // word 6: ADD T0, T0, 1 (immediate mode)
        word = mem[CODE_BASE + 6]; #1;
        check("ADD fmt",      format,   `VTX1_FMT_R3);
        check("ADD uses_imm", uses_imm, 1'b1);
        check("ADD imm",      imm,      32'sd1);
        check("ADD rd(T0)",   reg1,     5'd0);
        check("ADD rs1(T0)",  reg2,     5'd0);

        // word 9: JMP loop (absolute 0x1003)
        word = mem[CODE_BASE + 9]; #1;
        check("JMP fmt", format, `VTX1_FMT_JUMP);
        check("JMP tgt", imm,    32'h1003);

        // word 13: WFI
        word = mem[CODE_BASE + 13]; #1;
        check("WFI fmt", format,   `VTX1_FMT_NONE);
        check("WFI cat", category, `VTX1_CAT_SYS);

        // Data section: valid trit words carrying "Hello, VTX1!\0"
        for (i = 0; i < DATA_WORDS; i = i + 1) begin
            word = mem[DATA_BASE + i];
            #1;
            if (!valid && format != `VTX1_FMT_BAD) begin
                $display("  FAIL: data word %0d has undefined trits", i);
                errors = errors + 1;
            end
            if (word_value !== expected_data[i]) begin
                $display("  FAIL: data word %0d value = %0d, want %0d ('%s')",
                         i, word_value, expected_data[i], expected_data[i]);
                errors = errors + 1;
            end
        end
        $display("Data words checked: %0d (\"Hello, VTX1!\" + NUL)", DATA_WORDS);

        // An undefined trit pair must be rejected.
        word = 36'hFFFFFFFFF; #1;
        if (valid) begin
            $display("  FAIL: word of undefined trits reported valid");
            errors = errors + 1;
        end

        if (errors == 0)
            $display("=== PASS: decoder agrees with assembler on all %0d words ===",
                     CODE_WORDS + DATA_WORDS);
        else
            $display("=== FAIL: %0d error(s) ===", errors);
        $finish;
    end

endmodule
