package cmd

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"
	parser "github.com/kvany/vtx1/assembler/grammar"
)

func parseSource(t *testing.T, src string) *AST {
	t.Helper()
	input := antlr.NewInputStream(src)
	lexer := parser.Newvtx1_grammarLexer(input)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.Newvtx1_grammarParser(tokens)
	ast := BuildAST(p.Program())
	if ast == nil {
		t.Fatalf("BuildAST returned nil for source:\n%s", src)
	}
	return ast
}

func firstStatement(t *testing.T, ast *AST) StatementNode {
	t.Helper()
	for _, line := range ast.Program.Lines {
		if line.Statement != nil {
			return line.Statement
		}
	}
	t.Fatal("no statement in AST")
	return nil
}

func TestVisitInstruction(t *testing.T) {
	ast := parseSource(t, "        ADD T0, T1, T2\n")
	instr, ok := firstStatement(t, ast).(*InstructionNode)
	if !ok {
		t.Fatalf("statement is %T, want *InstructionNode", firstStatement(t, ast))
	}
	if instr.Mnemonic != "ADD" || len(instr.Operands) != 3 {
		t.Errorf("instruction = %s with %d operands", instr.Mnemonic, len(instr.Operands))
	}
}

func TestVisitDirectiveEquKeepsParamOrder(t *testing.T) {
	ast := parseSource(t, "        .EQU IOBASE, 0x2000\n")
	dir, ok := firstStatement(t, ast).(*DirectiveNode)
	if !ok || dir.Name != ".EQU" {
		t.Fatalf("statement = %#v, want .EQU directive", firstStatement(t, ast))
	}
	if len(dir.Params) != 2 {
		t.Fatalf(".EQU params = %d, want 2", len(dir.Params))
	}
	if id, ok := dir.Params[0].(*IdentifierNode); !ok || id.Name != "IOBASE" {
		t.Errorf(".EQU param 0 = %#v, want identifier IOBASE", dir.Params[0])
	}
	if imm, ok := dir.Params[1].(*ImmediateNode); !ok || imm.Value != "0x2000" {
		t.Errorf(".EQU param 1 = %#v, want immediate 0x2000", dir.Params[1])
	}
}

func TestVisitLabelledDirectiveKeepsLabel(t *testing.T) {
	ast := parseSource(t, "array_length:  .DW 10\n")
	var found *LineNode
	for _, line := range ast.Program.Lines {
		if line.Statement != nil {
			found = line
			break
		}
	}
	if found == nil || found.Label == nil || found.Label.Name != "array_length" {
		t.Fatalf("labelled directive lost its label: %#v", found)
	}
	if dir, ok := found.Statement.(*DirectiveNode); !ok || dir.Name != ".DW" {
		t.Errorf("statement = %#v, want .DW directive", found.Statement)
	}
}

func TestVisitOperandBinaryExpressionOrder(t *testing.T) {
	ast := parseSource(t, "        LD T1, vec_a+4\n")
	instr := firstStatement(t, ast).(*InstructionNode)
	if len(instr.Operands) != 2 {
		t.Fatalf("LD operands = %d, want 2", len(instr.Operands))
	}
	bin, ok := instr.Operands[1].(*BinaryOpNode)
	if !ok {
		t.Fatalf("operand 2 = %#v, want *BinaryOpNode", instr.Operands[1])
	}
	if id, ok := bin.Left.(*IdentifierNode); !ok || id.Name != "vec_a" {
		t.Errorf("left = %#v, want identifier vec_a", bin.Left)
	}
	if imm, ok := bin.Right.(*ImmediateNode); !ok || imm.Value != "4" {
		t.Errorf("right = %#v, want immediate 4", bin.Right)
	}
}

func TestVisitDataList(t *testing.T) {
	ast := parseSource(t, "        .DW 5, 9, 3\n")
	dir := firstStatement(t, ast).(*DirectiveNode)
	if len(dir.Params) != 3 {
		t.Errorf(".DW params = %d, want 3", len(dir.Params))
	}
}

func TestVisitVLIW(t *testing.T) {
	ast := parseSource(t, "        [ADD T0, T0, 1] [ADD T1, T1, 1] [NOP]\n")
	vl, ok := firstStatement(t, ast).(*VLIWInstructionNode)
	if !ok {
		t.Fatalf("statement = %T, want *VLIWInstructionNode", firstStatement(t, ast))
	}
	if len(vl.Instructions) != 3 {
		t.Errorf("VLIW ops = %d, want 3", len(vl.Instructions))
	}
}
