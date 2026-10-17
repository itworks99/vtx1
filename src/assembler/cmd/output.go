package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/kvany/vtx1/assembler/isa"
)

// writeOutput writes the generated word stream to the output file.
//
// The machine code is a stream of 18-trit words, each stored as 5 bytes
// (36 bits + 4 zero padding bits). Word sizes:
//
//	"36"  (default) - the flat word stream, 5 bytes per word
//	"108"           - words regrouped into VLIW-bundle-sized units of
//	                  54 trits stored as 14 bytes (`VLIW_BYTES)
//
// The legacy "8" and "ternary" sizes from the binary-slot era are gone:
// every word is already balanced ternary.
func writeOutput(binary []byte, outputFile, format, wordSize string) error {
	var data []byte
	switch wordSize {
	case "", "36":
		data = binary
	case "108":
		data = packBundles(binary)
	case "8", "ternary":
		return fmt.Errorf("word size %q is obsolete: the assembler is trit-native; use 36 (words) or 108 (VLIW bundles)", wordSize)
	default:
		return fmt.Errorf("unsupported word size: %s", wordSize)
	}

	switch format {
	case "binary":
		return ioutil.WriteFile(outputFile, data, 0644)
	case "hex":
		return ioutil.WriteFile(outputFile, []byte(formatAsHex(data)), 0644)
	case "objdump":
		// objdump always disassembles the flat word stream; bundle packing
		// is a storage concern, not a listing concern.
		return ioutil.WriteFile(outputFile, []byte(formatAsObjDump(binary)), 0644)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// packBundles regroups the 5-byte word stream into 108-bit units: every
// 3 words (54 trits) are repacked contiguously into 14 bytes, matching the
// hardware's `VLIW_WIDTH/`VLIW_BYTES. A trailing partial group is padded
// with zero-value words.
func packBundles(data []byte) []byte {
	words := len(data) / isa.WordBytes
	out := make([]byte, 0, (words+isa.VLIWSlots-1)/isa.VLIWSlots*isa.BundleBytes)
	for w := 0; w < words; w += isa.VLIWSlots {
		trits := make([]int, 0, isa.BundleTrits)
		for s := 0; s < isa.VLIWSlots; s++ {
			if w+s < words {
				off := (w + s) * isa.WordBytes
				trits = append(trits, isa.DecodeTrits(data[off:off+isa.WordBytes], isa.SlotTrits)...)
			} else {
				trits = append(trits, make([]int, isa.SlotTrits)...) // zero-value word
			}
		}
		out = append(out, isa.EncodeTrits(trits)...)
	}
	return out
}

// --- Intel HEX output ---

// formatAsHex formats binary data as Intel HEX (I8HEX): 16-byte data records
// starting at address 0, terminated by an EOF record.
func formatAsHex(data []byte) string {
	var sb strings.Builder
	const recLen = 16
	for off := 0; off < len(data); off += recLen {
		end := off + recLen
		if end > len(data) {
			end = len(data)
		}
		chunk := data[off:end]
		sum := byte(len(chunk)) + byte(off>>8) + byte(off)
		sb.WriteString(fmt.Sprintf(":%02X%04X00", len(chunk), off&0xFFFF))
		for _, b := range chunk {
			sb.WriteString(fmt.Sprintf("%02X", b))
			sum += b
		}
		sb.WriteString(fmt.Sprintf("%02X\n", byte(-sum)))
	}
	sb.WriteString(":00000001FF\n")
	return sb.String()
}

// --- Objdump output ---

func regNameFor(index int) string {
	if r, ok := isa.RegByIndex[index]; ok {
		return r.Name
	}
	return fmt.Sprintf("R?%d", index)
}

// tritsString renders trits as balanced ternary digits (+, 0, -; ? for the
// undefined pair).
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
		default:
			sb.WriteString("?")
		}
	}
	return sb.String()
}

func formatValue(v int64) string {
	if v < 0 {
		return fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("0x%X", v)
}

// formatAsObjDump disassembles the word stream into a listing of word
// address, trit pattern, and decoded operation.
func formatAsObjDump(data []byte) string {
	var sb strings.Builder
	sb.WriteString("VTX1 objdump (word addresses; 18 balanced trits per word)\n\n")
	sb.WriteString("word      trits               disassembly\n")
	sb.WriteString("--------  ------------------  -----------\n")
	words := len(data) / isa.WordBytes
	for w := 0; w < words; w++ {
		off := w * isa.WordBytes
		trits := isa.DecodeTrits(data[off:off+isa.WordBytes], isa.SlotTrits)
		sb.WriteString(fmt.Sprintf("%08X  %s  %s\n", w, tritsString(trits), disassembleSlot(trits)))
	}
	if rest := len(data) % isa.WordBytes; rest != 0 {
		off := len(data) - rest
		var hexb []string
		for _, b := range data[off:] {
			hexb = append(hexb, fmt.Sprintf("%02X", b))
		}
		sb.WriteString(fmt.Sprintf("--------  partial word: .byte %s\n", strings.Join(hexb, ", ")))
	}
	return sb.String()
}

// fieldValue reads trits [pos, pos+n) as a balanced number.
func fieldValue(trits []int, pos, n int) int64 {
	return isa.TritsValue(trits[pos : pos+n])
}

// fieldReg reads a 3-trit register field.
func fieldReg(trits []int, pos int) int {
	return int(fieldValue(trits, pos, isa.RegTrits)) + isa.RegBias
}

func validTrits(trits []int) bool {
	for _, t := range trits {
		if t < -1 || t > 1 {
			return false
		}
	}
	return true
}

// disassembleSlot decodes one 18-trit word into assembly text. Words that do
// not decode to a known operation (e.g. data) are shown as .word, and words
// containing undefined trit pairs are flagged.
func disassembleSlot(trits []int) string {
	if !validTrits(trits) {
		return ".invalid (undefined trit)"
	}
	opcode := int(fieldValue(trits, 0, isa.OpcodeTrits)) + isa.OpcodeBias
	ins, ok := isa.ByValue[opcode]
	if !ok {
		return fmt.Sprintf(".word %d", isa.TritsValue(trits))
	}
	switch ins.Format {
	case isa.FormatNone:
		// Data words can alias operand-less opcodes; show the value too if
		// any operand trit is set.
		if isa.TritsValue(trits[isa.OpcodeTrits:]) != 0 {
			return fmt.Sprintf(".word %d", isa.TritsValue(trits))
		}
		return ins.Name
	case isa.FormatR1:
		return fmt.Sprintf("%s %s", ins.Name, regNameFor(fieldReg(trits, 4)))
	case isa.FormatR2:
		return fmt.Sprintf("%s %s, %s", ins.Name, regNameFor(fieldReg(trits, 4)), regNameFor(fieldReg(trits, 7)))
	case isa.FormatR3:
		third := ""
		if trits[10] == 1 {
			third = fmt.Sprintf("%d", fieldValue(trits, 11, 7))
		} else {
			third = regNameFor(fieldReg(trits, 15))
		}
		return fmt.Sprintf("%s %s, %s, %s", ins.Name,
			regNameFor(fieldReg(trits, 4)), regNameFor(fieldReg(trits, 7)), third)
	case isa.FormatMem:
		return fmt.Sprintf("%s %s, %s", ins.Name, regNameFor(fieldReg(trits, 4)), disassembleMemField(trits))
	case isa.FormatBranch:
		rs2 := "0"
		if idx := fieldReg(trits, 7); idx != isa.ZeroRegIndex {
			rs2 = regNameFor(idx)
		}
		return fmt.Sprintf("%s %s, %s, PC%+d", ins.Name,
			regNameFor(fieldReg(trits, 4)), rs2, fieldValue(trits, 10, 8))
	case isa.FormatJump:
		return fmt.Sprintf("%s %s", ins.Name, formatValue(fieldValue(trits, 4, 14)))
	}
	return fmt.Sprintf(".word %d", isa.TritsValue(trits))
}

func disassembleMemField(trits []int) string {
	switch trits[7] {
	case 0: // direct
		return formatValue(fieldValue(trits, 8, 10))
	case 1: // base+index ([base] when index is TZ)
		base := regNameFor(fieldReg(trits, 8))
		if idx := fieldReg(trits, 11); idx != isa.ZeroRegIndex {
			return fmt.Sprintf("[%s+%s]", base, regNameFor(idx))
		}
		return fmt.Sprintf("[%s]", base)
	default: // -1: base+offset
		return fmt.Sprintf("[%s+%d]", regNameFor(fieldReg(trits, 8)), fieldValue(trits, 11, 7))
	}
}
