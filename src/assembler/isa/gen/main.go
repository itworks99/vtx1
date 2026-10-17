// Command gen generates the Verilog ISA constants header from the Go ISA
// definition (package isa), so the RTL and the assembler share one source
// of truth. Run from src/assembler:
//
//	go run ./isa/gen -out ../common/vtx1_isa_generated.vh
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/kvany/vtx1/assembler/isa"
)

// bits renders a trit pattern as a Verilog binary literal using the hardware
// trit encoding (00 = -1, 01 = 0, 10 = +1), most significant trit first.
func bits(trits []int) string {
	var sb strings.Builder
	for i, t := range trits {
		if i > 0 {
			sb.WriteString("_")
		}
		switch t {
		case -1:
			sb.WriteString("00")
		case 0:
			sb.WriteString("01")
		case 1:
			sb.WriteString("10")
		default:
			sb.WriteString("11")
		}
	}
	return fmt.Sprintf("%d'b%s", 2*len(trits), sb.String())
}

func tritsString(trits []int) string {
	var sb strings.Builder
	for _, t := range trits {
		switch t {
		case -1:
			sb.WriteString("-")
		case 0:
			sb.WriteString("0")
		case 1:
			sb.WriteString("+")
		}
	}
	return sb.String()
}

func main() {
	out := flag.String("out", "../common/vtx1_isa_generated.vh", "output header path")
	flag.Parse()

	var sb strings.Builder
	// "~" stands in for the Verilog backtick, which would terminate a Go raw string.
	header := `// VTX1 trit-native ISA constants
// GENERATED FILE - DO NOT EDIT.
// Source of truth: src/assembler/isa/isa.go
// Regenerate with:  cd src/assembler && go run ./isa/gen
//
// Instruction word: 18 trits = 36 bits (2 bits/trit: 00=-1, 01=0, 10=+1).
// Trit t (0 = most significant) occupies bits [35-2t : 34-2t].
//
// Field layout (trit positions, MSB first):
//   all:     opcode  trits 0..3   (bits [35:28])
//   R3:      rd 4..6 | rs1 7..9 | mode 10 | payload 11..17
//            (mode 0: rs2 = trits 15..17; mode +1: 7-trit immediate)
//   R2:      rd 4..6 | rs 7..9
//   R1:      r  4..6
//   Mem:     reg 4..6 | mode 7 | payload 8..17
//            (mode 0: 10-trit address; mode +1: base 8..10, index 11..13;
//             mode -1: base 8..10, 7-trit offset 11..17)
//   Branch:  rs1 4..6 | rs2 7..9 | 8-trit PC-relative word offset 10..17
//            (rs2 = TZ means compare against zero)
//   Jump:    14-trit absolute word address, trits 4..17
//
// A VLIW bundle is 3 instruction words = 54 trits = 108 bits (~VLIW_WIDTH).
// The all-zero-trit word is NOP.

~ifndef VTX1_ISA_GENERATED_VH
~define VTX1_ISA_GENERATED_VH

// Geometry
`
	sb.WriteString(strings.ReplaceAll(header, "~", "`"))
	fmt.Fprintf(&sb, "`define VTX1_SLOT_TRITS   %d\n", isa.SlotTrits)
	fmt.Fprintf(&sb, "`define VTX1_SLOT_BITS    %d\n", isa.SlotTrits*2)
	fmt.Fprintf(&sb, "`define VTX1_OPCODE_BITS  %d\n", isa.OpcodeTrits*2)
	fmt.Fprintf(&sb, "`define VTX1_REG_BITS     %d\n\n", isa.RegTrits*2)

	sb.WriteString("// Opcodes (4 trits as 8 bits; comment shows balanced-trit pattern and value)\n")
	instrs := make([]isa.Instruction, len(isa.Instructions))
	copy(instrs, isa.Instructions)
	sort.Slice(instrs, func(i, j int) bool { return instrs[i].Value < instrs[j].Value })
	for _, ins := range instrs {
		p := isa.OpcodePattern(ins.Value)
		micro := ""
		if ins.MultiCycle {
			micro = ", multi-cycle"
		}
		fmt.Fprintf(&sb, "`define VTX1_OPC_%-8s %s  // %s (value %d, %s%s)\n",
			ins.Name, bits(p), tritsString(p), ins.Value, ins.Category, micro)
	}

	sb.WriteString("\n// Registers (3 trits as 6 bits)\n")
	for _, r := range isa.Registers {
		p := isa.RegPattern(r.Index)
		fmt.Fprintf(&sb, "`define VTX1_REG_%-3s %s  // %s (index %d)\n",
			r.Name, bits(p), tritsString(p), r.Index)
	}

	// Decode helper defines live in the guarded header; the function itself
	// goes to a separate unguarded file (vtx1_isa_decode.vh) because Verilog
	// functions must be included inside each module body that uses them.
	sb.WriteString(`
// Operand formats
` + "`define" + ` VTX1_FMT_NONE   3'd0
` + "`define" + ` VTX1_FMT_R3     3'd1
` + "`define" + ` VTX1_FMT_R2     3'd2
` + "`define" + ` VTX1_FMT_R1     3'd3
` + "`define" + ` VTX1_FMT_MEM    3'd4
` + "`define" + ` VTX1_FMT_BRANCH 3'd5
` + "`define" + ` VTX1_FMT_JUMP   3'd6
` + "`define" + ` VTX1_FMT_BAD    3'd7

// Operation categories
` + "`define" + ` VTX1_CAT_ALU  3'd0
` + "`define" + ` VTX1_CAT_MEM  3'd1
` + "`define" + ` VTX1_CAT_CTRL 3'd2
` + "`define" + ` VTX1_CAT_VEC  3'd3
` + "`define" + ` VTX1_CAT_FP   3'd4
` + "`define" + ` VTX1_CAT_SYS  3'd5

`)
	sb.WriteString("\n`endif // VTX1_ISA_GENERATED_VH\n")

	// Second file: the decode function, unguarded, for inclusion inside a
	// module body.
	var fb strings.Builder
	fb.WriteString(`// VTX1 opcode decode helper
// GENERATED FILE - DO NOT EDIT. Source of truth: src/assembler/isa/isa.go
// Include INSIDE a module body (no include guard: each module needs its own
// copy of the function). Requires vtx1_isa_generated.vh for the defines.
//
// vtx1_opcode_info(value) = {multicycle, category[2:0], format[2:0]}
// where value = 40 + 27*t0 + 9*t1 + 3*t2 + t3 over the signed opcode trits.
`)
	fb.WriteString("function [6:0] vtx1_opcode_info(input [6:0] value);\n")
	fb.WriteString("    case (value)\n")
	fmtName := map[isa.Format]string{
		isa.FormatNone: "`VTX1_FMT_NONE", isa.FormatR3: "`VTX1_FMT_R3",
		isa.FormatR2: "`VTX1_FMT_R2", isa.FormatR1: "`VTX1_FMT_R1",
		isa.FormatMem: "`VTX1_FMT_MEM", isa.FormatBranch: "`VTX1_FMT_BRANCH",
		isa.FormatJump: "`VTX1_FMT_JUMP",
	}
	catName := map[isa.Category]string{
		isa.CatALU: "`VTX1_CAT_ALU", isa.CatMEM: "`VTX1_CAT_MEM",
		isa.CatCTRL: "`VTX1_CAT_CTRL", isa.CatVEC: "`VTX1_CAT_VEC",
		isa.CatFP: "`VTX1_CAT_FP", isa.CatSYS: "`VTX1_CAT_SYS",
	}
	for _, ins := range instrs {
		micro := "1'b0"
		if ins.MultiCycle {
			micro = "1'b1"
		}
		fmt.Fprintf(&fb, "        7'd%-2d: vtx1_opcode_info = {%s, %s, %s}; // %s\n",
			ins.Value, micro, catName[ins.Category], fmtName[ins.Format], ins.Name)
	}
	fb.WriteString("        default: vtx1_opcode_info = {1'b0, 3'd7, `VTX1_FMT_BAD};\n")
	fb.WriteString("    endcase\nendfunction\n")

	if err := ioutil.WriteFile(*out, []byte(sb.String()), 0644); err != nil {
		log.Fatal(err)
	}
	decodeOut := strings.TrimSuffix(*out, ".vh") + "_decode.vh"
	decodeOut = strings.Replace(decodeOut, "generated_decode", "decode", 1)
	if err := ioutil.WriteFile(decodeOut, []byte(fb.String()), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s and %s (%d opcodes, %d registers)\n", *out, decodeOut, len(instrs), len(isa.Registers))
}
