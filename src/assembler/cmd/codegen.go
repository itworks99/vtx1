package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kvany/vtx1/assembler/isa"
)

// CodeGenerator traverses the AST and emits machine code as a stream of
// 18-trit words, each stored as 5 bytes (see package isa). All addresses are
// word addresses: a scalar instruction occupies 1 word, a VLIW bundle 3
// words, and each data item one word.
type CodeGenerator struct {
	Output      []byte
	SymbolTable *SymbolTable
	CurrentAddr uint32 // in words
	Labels      map[string]uint32
	Equs        map[string]uint32
}

// NewCodeGenerator creates a new code generator with the given symbol table.
func NewCodeGenerator(symbolTable *SymbolTable) *CodeGenerator {
	return &CodeGenerator{
		Output:      make([]byte, 0),
		SymbolTable: symbolTable,
		CurrentAddr: 0,
		Labels:      make(map[string]uint32),
		Equs:        make(map[string]uint32),
	}
}

// statementWords returns how many machine words a statement occupies.
func statementWords(stmt StatementNode) uint32 {
	switch s := stmt.(type) {
	case *InstructionNode:
		return 1
	case *VLIWInstructionNode:
		return isa.VLIWSlots
	case *DirectiveNode:
		name := strings.ToUpper(s.Name)
		switch name {
		case ".DB", ".DW", ".DT":
			return uint32(dataItemCount(s.Params))
		case ".SPACE":
			if len(s.Params) > 0 {
				if imm, err := parseImmediateOperand(s.Params[0]); err == nil {
					return uint32(imm)
				}
			}
		}
	}
	return 0
}

// Pass 1: Collect labels and .EQUs, handle .ORG for address tracking
func (cg *CodeGenerator) collectSymbols(ast *AST) error {
	addr := uint32(0)
	for _, line := range ast.Program.Lines {
		if line.Label != nil {
			cg.Labels[line.Label.Name] = addr
		}
		if line.Statement == nil {
			continue
		}
		if dir, ok := line.Statement.(*DirectiveNode); ok {
			name := strings.ToUpper(dir.Name)
			switch name {
			case ".ORG":
				if len(dir.Params) > 0 {
					imm, _ := parseImmediateOperand(dir.Params[0])
					addr = uint32(imm)
				}
				continue
			case ".EQU":
				if len(dir.Params) == 2 {
					if id, ok := dir.Params[0].(*IdentifierNode); ok {
						imm, _ := parseImmediateOperand(dir.Params[1])
						cg.Equs[id.Name] = uint32(imm)
					}
				}
				continue
			}
		}
		addr += statementWords(line.Statement)
	}
	return nil
}

// dataItemCount returns the number of emitted words in a data directive,
// expanding quoted strings to their character count (one word per character).
func dataItemCount(params []OperandNode) int {
	count := 0
	for _, op := range params {
		if imm, ok := op.(*ImmediateNode); ok && isQuotedString(imm.Value) {
			count += len(unquoteString(imm.Value))
		} else {
			count++
		}
	}
	return count
}

// Pass 2: Emit code/data, resolving symbols
func (cg *CodeGenerator) Generate(ast *AST) error {
	cg.Output = make([]byte, 0)
	cg.CurrentAddr = 0
	cg.Labels = make(map[string]uint32)
	cg.Equs = make(map[string]uint32)
	if err := cg.collectSymbols(ast); err != nil {
		return err
	}
	cg.CurrentAddr = 0
	for _, line := range ast.Program.Lines {
		if line.Statement == nil {
			continue
		}
		switch stmt := line.Statement.(type) {
		case *InstructionNode:
			if err := cg.emitInstruction(stmt); err != nil {
				return err
			}
		case *VLIWInstructionNode:
			if err := cg.emitVLIWInstruction(stmt); err != nil {
				return err
			}
		case *DirectiveNode:
			if err := cg.emitDirective(stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

// Add a global set of deprecated mnemonics and directives
var deprecatedMnemonics = map[string]bool{
	"OLDOP": true, // Example: replace with real deprecated mnemonics
}
var deprecatedDirectives = map[string]bool{
	".OLD": true, // Example: replace with real deprecated directives
}

// Add a pointer to ErrorManager.Warnings for warning reporting
var codegenWarnings *[]error

func AttachCodegenWarnings(w *[]error) {
	codegenWarnings = w
}

func codegenWarn(format string, args ...interface{}) {
	if codegenWarnings != nil {
		*codegenWarnings = append(*codegenWarnings, fmt.Errorf(format, args...))
	}
}

func (cg *CodeGenerator) emitInstruction(instr *InstructionNode) error {
	if deprecatedMnemonics[strings.ToUpper(instr.Mnemonic)] {
		codegenWarn("warning: instruction '%s' at line %d is deprecated", instr.Mnemonic, instr.Line)
	}
	slot, _, err := cg.encodeOperation(instr, cg.CurrentAddr)
	if err != nil {
		return err
	}
	cg.Output = append(cg.Output, slot.Bytes()...)
	cg.CurrentAddr++
	return nil
}

// emitVLIWInstruction packs up to 3 operations into one VLIW bundle
// (3 x 18-trit words, NOP-padded) and enforces the packing rules from
// docs/vliw_encoding.adoc:
//  1. at most 3 operations per bundle
//  2. at most 1 memory operation per bundle
//  3. at most 1 control operation per bundle
//  4. no two operations may write the same register
//  5. read-after-write between slots is reported as a hazard warning
func (cg *CodeGenerator) emitVLIWInstruction(vliw *VLIWInstructionNode) error {
	if len(vliw.Instructions) > isa.VLIWSlots {
		return fmt.Errorf("VLIW bundle at line %d contains %d operations, maximum is %d",
			vliw.Line, len(vliw.Instructions), isa.VLIWSlots)
	}

	baseAddr := cg.CurrentAddr
	infos := make([]opInfo, isa.VLIWSlots)
	catCount := make(map[isa.Category]int)
	var bundle []byte

	for i := 0; i < isa.VLIWSlots; i++ {
		var instr *InstructionNode
		if i < len(vliw.Instructions) {
			instr = vliw.Instructions[i]
		} else {
			instr = &InstructionNode{Mnemonic: "NOP", Operands: nil, Line: vliw.Line}
		}
		slot, info, err := cg.encodeOperation(instr, baseAddr)
		if err != nil {
			return fmt.Errorf("VLIW error at line %d: %v", instr.Line, err)
		}
		infos[i] = info
		if info.Mnemonic != "NOP" {
			catCount[info.Category]++
		}
		bundle = append(bundle, slot.Bytes()...)
	}

	if catCount[isa.CatMEM] > 1 {
		return fmt.Errorf("VLIW packing violation at line %d: %d memory operations in one bundle, maximum is 1",
			vliw.Line, catCount[isa.CatMEM])
	}
	if catCount[isa.CatCTRL] > 1 {
		return fmt.Errorf("VLIW packing violation at line %d: %d control operations in one bundle, maximum is 1",
			vliw.Line, catCount[isa.CatCTRL])
	}

	// Rule 4: no two slots may write the same register.
	writtenBy := make(map[string]string)
	for _, info := range infos {
		for _, reg := range info.Writes {
			if other, dup := writtenBy[reg]; dup {
				return fmt.Errorf("VLIW resource conflict at line %d: register %s written by both %s and %s in the same bundle",
					vliw.Line, reg, other, info.Mnemonic)
			}
			writtenBy[reg] = info.Mnemonic
		}
	}

	// Rule 5: a slot reading a register written by another slot in the same
	// bundle sees the OLD value (parallel semantics). Warn, since this is a
	// frequent source of surprises when porting sequential code.
	for i, info := range infos {
		for _, reg := range info.Reads {
			if writer, ok := writtenBy[reg]; ok && !slotWritesReg(infos[i], reg) {
				codegenWarn("warning: VLIW hazard at line %d: %s reads %s while %s writes it in the same bundle (reads see the pre-bundle value)",
					vliw.Line, info.Mnemonic, reg, writer)
			}
		}
	}

	cg.Output = append(cg.Output, bundle...)
	cg.CurrentAddr += isa.VLIWSlots
	return nil
}

// slotWritesReg reports whether the given op itself writes the register
// (e.g. INC reads and writes the same register, which is not a cross-slot hazard).
func slotWritesReg(info opInfo, reg string) bool {
	for _, w := range info.Writes {
		if w == reg {
			return true
		}
	}
	return false
}

// --- Directive/Data Emission ---
func (cg *CodeGenerator) emitDirective(dir *DirectiveNode) error {
	if deprecatedDirectives[strings.ToUpper(dir.Name)] {
		codegenWarn("warning: directive '%s' at line %d is deprecated", dir.Name, dir.Line)
	}
	name := strings.ToUpper(dir.Name)
	switch name {
	case ".ORG":
		if len(dir.Params) > 0 {
			imm, _ := parseImmediateOperand(dir.Params[0])
			cg.CurrentAddr = uint32(imm)
		}
	case ".SPACE":
		if len(dir.Params) > 0 {
			imm, _ := parseImmediateOperand(dir.Params[0])
			for i := int64(0); i < imm; i++ {
				if err := cg.emitDataWord(0, dir.Line); err != nil {
					return err
				}
			}
		}
	case ".DB", ".DW", ".DT":
		// All data directives emit one 18-trit word per item; strings emit
		// one word per character.
		for _, op := range dir.Params {
			if v, ok := op.(*ImmediateNode); ok && isQuotedString(v.Value) {
				for _, c := range unquoteString(v.Value) {
					if err := cg.emitDataWord(int64(c), dir.Line); err != nil {
						return err
					}
				}
				continue
			}
			imm, err := parseImmediateOperand(op)
			if err != nil {
				return fmt.Errorf("%s at line %d: %v", name, dir.Line, err)
			}
			if err := cg.emitDataWord(imm, dir.Line); err != nil {
				return err
			}
		}
	case ".EQU":
		// Already handled in pass 1
		return nil
	case ".INCLUDE", ".MACRO", ".ENDM":
		// Handled by the preprocessor; nothing reaches code generation.
		return nil
	// TODO: Add support for other assembler directives as needed
	default:
		// Ignore other directives for now
	}
	return nil
}

// emitDataWord emits one 18-trit data word.
func (cg *CodeGenerator) emitDataWord(v int64, line int) error {
	trits, err := isa.BalancedTrits(v, isa.SlotTrits)
	if err != nil {
		return fmt.Errorf("data value at line %d: %v", line, err)
	}
	cg.Output = append(cg.Output, isa.EncodeTrits(trits)...)
	cg.CurrentAddr++
	return nil
}

// Helper: check if a string is quoted (e.g., '"Hello"')
func isQuotedString(s string) bool {
	return len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"'
}

// Helper: remove quotes and unescape (basic)
func unquoteString(s string) string {
	if isQuotedString(s) {
		return s[1 : len(s)-1]
	}
	return s
}

// --- Helpers ---

// parseImmediateString parses a literal in any supported base:
// decimal, 0x hex, 0b binary, or 0t balanced ternary (digits +, -, 0).
func parseImmediateString(val string) (int64, error) {
	switch {
	case strings.HasPrefix(val, "0x"):
		return strconv.ParseInt(val[2:], 16, 64)
	case strings.HasPrefix(val, "0b"):
		return strconv.ParseInt(val[2:], 2, 64)
	case strings.HasPrefix(val, "0t"):
		return parseBalancedTernary(val[2:])
	default:
		return strconv.ParseInt(val, 10, 64)
	}
}

// parseBalancedTernary converts balanced ternary digits (+, -, 0) to an
// integer; the leftmost digit is the most significant trit:
// "+0-" = (+1)*9 + 0*3 + (-1)*1 = 8.
func parseBalancedTernary(digits string) (int64, error) {
	if digits == "" {
		return 0, fmt.Errorf("empty balanced ternary literal")
	}
	var v int64
	for _, c := range digits {
		v *= 3
		switch c {
		case '+':
			v++
		case '-':
			v--
		case '0':
			// nothing
		default:
			return 0, fmt.Errorf("invalid balanced ternary digit %q", c)
		}
	}
	return v, nil
}

func parseImmediateOperand(op OperandNode) (int64, error) {
	switch v := op.(type) {
	case *ImmediateNode:
		return parseImmediateString(v.Value)
	case *IdentifierNode:
		// Symbol/label reference
		return 0, nil // Will be resolved in pass 2
	}
	return 0, fmt.Errorf("unsupported immediate operand")
}

func (cg *CodeGenerator) resolveSymbol(name string) uint32 {
	if v, ok := cg.Labels[name]; ok {
		return v
	}
	if v, ok := cg.Equs[name]; ok {
		return v
	}
	return 0
}

func (cg *CodeGenerator) resolveOperandAddr(op OperandNode) uint32 {
	switch v := op.(type) {
	case *ImmediateNode:
		imm, _ := parseImmediateOperand(v)
		return uint32(imm)
	case *IdentifierNode:
		return cg.resolveSymbol(v.Name)
	case *BinaryOpNode:
		// label+offset expressions, e.g. "vec_a+1"
		if v.Op == "+" {
			return cg.resolveOperandAddr(v.Left) + cg.resolveOperandAddr(v.Right)
		}
	}
	return 0
}
