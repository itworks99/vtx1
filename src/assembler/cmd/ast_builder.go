package cmd

import (
	"fmt"
	"reflect"

	parser "github.com/kvany/vtx1/assembler/grammar"

	"github.com/antlr4-go/antlr/v4"
)

// ASTBuilder implements the vtx1_grammarVisitor interface and builds the AST from the parse tree.
type ASTBuilder struct {
	parser.Basevtx1_grammarVisitor
}

// BuildAST builds the AST from the ANTLR parse tree root.
func BuildAST(tree antlr.ParseTree) *AST {
	builder := &ASTBuilder{}
	if programCtx, ok := tree.(*parser.ProgramContext); ok {
		program, ok := builder.VisitProgram(programCtx).(*ProgramNode)
		if !ok || program == nil {
			fmt.Println("[ASTBuilder] Error: Could not build AST. The parse tree did not yield a valid ProgramNode.")
			return nil
		}
		return &AST{Program: program}
	}
	fmt.Println("[ASTBuilder] Error: Root node is not a *ProgramContext.")
	return nil
}

func (b *ASTBuilder) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(b)
}

func (b *ASTBuilder) VisitProgram(ctx *parser.ProgramContext) interface{} {
	lines := []*LineNode{}
	for _, lineCtx := range ctx.AllLine() {
		if lc, ok := lineCtx.(*parser.LineContext); ok {
			line := b.VisitLine(lc)
			if line != nil {
				lines = append(lines, line.(*LineNode))
			}
		}
	}
	return &ProgramNode{Lines: lines}
}

func (b *ASTBuilder) VisitLine(ctx *parser.LineContext) interface{} {
	var label *LabelNode
	var statement StatementNode
	var comment string

	fmt.Printf("[ASTBuilder] VisitLine: ctx.Instruction()=%v, ctx.Directive()=%v, ctx.VliwInstruction()=%v, ctx.LabelledDirective()=%v\n", ctx.Instruction(), ctx.Directive(), ctx.VliwInstruction(), ctx.LabelledDirective())

	if ctx.Label() != nil {
		if l := b.Visit(ctx.Label()); l != nil {
			label = l.(*LabelNode)
		}
	}
	if ctx.Instruction() != nil {
		if s := b.Visit(ctx.Instruction()); s != nil {
			statement = s.(*InstructionNode)
		}
	} else if ctx.Directive() != nil {
		if s := b.Visit(ctx.Directive()); s != nil {
			statement = s.(*DirectiveNode)
		}
	} else if ctx.VliwInstruction() != nil {
		if s := b.Visit(ctx.VliwInstruction()); s != nil {
			statement = s.(*VLIWInstructionNode)
		}
	} else if ctx.LabelledDirective() != nil {
		if s := b.Visit(ctx.LabelledDirective()); s != nil {
			statement = s.(*DirectiveNode)
		}
	}
	// Handle optional end-of-line comment
	if n := len(ctx.AllComment()); n > 0 {
		comment = ctx.Comment(n - 1).GetText()
	}
	fmt.Printf("[ASTBuilder] VisitLine: statement type = %v, line text = %q\n", reflect.TypeOf(statement), ctx.GetText())
	return &LineNode{
		Label:     label,
		Statement: statement,
		Comment:   comment,
		Line:      ctx.GetStart().GetLine(),
		Column:    ctx.GetStart().GetColumn(),
	}
}

func (b *ASTBuilder) VisitLabel(ctx *parser.LabelContext) interface{} {
	return &LabelNode{
		Name:   ctx.IDENTIFIER().GetText(),
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
}

func (b *ASTBuilder) VisitInstruction(ctx *parser.InstructionContext) interface{} {
	mnemonic := ctx.Mnemonic().GetText()
	operands := []OperandNode{}
	for _, opCtx := range ctx.AllOperand() {
		operand := b.Visit(opCtx).(OperandNode)
		operands = append(operands, operand)
	}
	return &InstructionNode{
		Mnemonic: mnemonic,
		Operands: operands,
		Line:     ctx.GetStart().GetLine(),
		Column:   ctx.GetStart().GetColumn(),
	}
}

func (b *ASTBuilder) VisitDirective(ctx *parser.DirectiveContext) interface{} {
	name := ctx.GetText()
	params := []OperandNode{}
	if ctx.Immediate() != nil {
		params = append(params, b.Visit(ctx.Immediate()).(OperandNode))
	}
	if ctx.DataList() != nil {
		params = append(params, b.Visit(ctx.DataList()).([]OperandNode)...)
	}
	if ctx.IDENTIFIER() != nil {
		params = append(params, &IdentifierNode{
			Name:   ctx.IDENTIFIER().GetText(),
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine(),
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn(),
		})
	}
	if ctx.STRING() != nil {
		params = append(params, &ImmediateNode{
			Value:  ctx.STRING().GetText(),
			Line:   ctx.STRING().GetSymbol().GetLine(),
			Column: ctx.STRING().GetSymbol().GetColumn(),
		})
	}
	return &DirectiveNode{
		Name:   name,
		Params: params,
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
}

func (b *ASTBuilder) VisitVliwInstruction(ctx *parser.VliwInstructionContext) interface{} {
	instructions := []*InstructionNode{}
	for _, instrCtx := range ctx.AllInstruction() {
		instr := b.Visit(instrCtx).(*InstructionNode)
		instructions = append(instructions, instr)
	}
	return &VLIWInstructionNode{
		Instructions: instructions,
		Line:         ctx.GetStart().GetLine(),
		Column:       ctx.GetStart().GetColumn(),
	}
}

func (b *ASTBuilder) VisitOperand(ctx *parser.OperandContext) interface{} {
	if ctx.Register() != nil {
		return b.Visit(ctx.Register()).(OperandNode)
	}
	if len(ctx.AllImmediate()) > 0 && ctx.PLUS() == nil {
		return b.Visit(ctx.Immediate(0)).(OperandNode)
	}
	if ctx.MemoryOperand() != nil {
		return b.Visit(ctx.MemoryOperand()).(OperandNode)
	}
	if len(ctx.AllIDENTIFIER()) > 0 && ctx.PLUS() == nil {
		return &IdentifierNode{
			Name:   ctx.IDENTIFIER(0).GetText(),
			Line:   ctx.IDENTIFIER(0).GetSymbol().GetLine(),
			Column: ctx.IDENTIFIER(0).GetSymbol().GetColumn(),
		}
	}
	// Handle binary operations
	if ctx.PLUS() != nil {
		var left OperandNode
		var right OperandNode
		if len(ctx.AllIDENTIFIER()) > 0 {
			left = &IdentifierNode{
				Name:   ctx.IDENTIFIER(0).GetText(),
				Line:   ctx.IDENTIFIER(0).GetSymbol().GetLine(),
				Column: ctx.IDENTIFIER(0).GetSymbol().GetColumn(),
			}
		} else if len(ctx.AllImmediate()) > 0 {
			left = b.Visit(ctx.Immediate(0)).(OperandNode)
		}
		if len(ctx.AllIDENTIFIER()) > 1 {
			right = &IdentifierNode{
				Name:   ctx.IDENTIFIER(1).GetText(),
				Line:   ctx.IDENTIFIER(1).GetSymbol().GetLine(),
				Column: ctx.IDENTIFIER(1).GetSymbol().GetColumn(),
			}
		} else if len(ctx.AllImmediate()) > 1 {
			right = b.Visit(ctx.Immediate(1)).(OperandNode)
		} else if len(ctx.AllImmediate()) > 0 && len(ctx.AllIDENTIFIER()) > 0 {
			if left == nil {
				left = b.Visit(ctx.Immediate(0)).(OperandNode)
			}
			if right == nil {
				right = &IdentifierNode{
					Name:   ctx.IDENTIFIER(0).GetText(),
					Line:   ctx.IDENTIFIER(0).GetSymbol().GetLine(),
					Column: ctx.IDENTIFIER(0).GetSymbol().GetColumn(),
				}
			}
		}
		return &BinaryOpNode{
			Left:   left,
			Op:     "+",
			Right:  right,
			Line:   ctx.GetStart().GetLine(),
			Column: ctx.GetStart().GetColumn(),
		}
	}
	return nil
}

func (b *ASTBuilder) VisitRegister(ctx *parser.RegisterContext) interface{} {
	return &RegisterNode{
		Name:   ctx.GetText(),
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
}

func (b *ASTBuilder) VisitImmediate(ctx *parser.ImmediateContext) interface{} {
	return &ImmediateNode{
		Value:  ctx.GetText(),
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
}

func (b *ASTBuilder) VisitMemoryOperand(ctx *parser.MemoryOperandContext) interface{} {
	base := ""
	index := ""
	offset := ""
	if ctx.BaseRegister() != nil {
		base = ctx.BaseRegister().GetText()
	}
	if ctx.IndexRegister() != nil {
		index = ctx.IndexRegister().GetText()
	}
	if ctx.OffsetImmediate() != nil {
		offset = ctx.OffsetImmediate().GetText()
	}
	return &MemoryOperandNode{
		Base:   base,
		Index:  index,
		Offset: offset,
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
}

func (b *ASTBuilder) VisitDataList(ctx *parser.DataListContext) interface{} {
	operands := []OperandNode{}
	for _, itemCtx := range ctx.AllDataItem() {
		operand := b.Visit(itemCtx).(OperandNode)
		operands = append(operands, operand)
	}
	return operands
}

func (b *ASTBuilder) VisitDataItem(ctx *parser.DataItemContext) interface{} {
	if ctx.Immediate() != nil {
		return b.Visit(ctx.Immediate()).(OperandNode)
	}
	if ctx.STRING() != nil {
		return &ImmediateNode{
			Value:  ctx.STRING().GetText(),
			Line:   ctx.STRING().GetSymbol().GetLine(),
			Column: ctx.STRING().GetSymbol().GetColumn(),
		}
	}
	return nil
}
