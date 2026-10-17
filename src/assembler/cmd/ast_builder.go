package cmd

import (
	"fmt"

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

// Visit dispatches on the concrete parse tree node type. The generated parser
// was produced without ANTLR's -visitor option, so the contexts have no Accept
// methods of their own: tree.Accept(b) would fall through to the runtime's
// default VisitChildren and silently return nil. Dispatching here keeps the
// builder independent of how the parser was generated.
func (b *ASTBuilder) Visit(tree antlr.ParseTree) interface{} {
	switch t := tree.(type) {
	case *parser.ProgramContext:
		return b.VisitProgram(t)
	case *parser.LineContext:
		return b.VisitLine(t)
	case *parser.LabelContext:
		return b.VisitLabel(t)
	case *parser.InstructionContext:
		return b.VisitInstruction(t)
	case *parser.DirectiveContext:
		return b.VisitDirective(t)
	case *parser.LabelledDirectiveContext:
		return b.VisitLabelledDirective(t)
	case *parser.VliwInstructionContext:
		return b.VisitVliwInstruction(t)
	case *parser.OperandContext:
		return b.VisitOperand(t)
	case *parser.RegisterContext:
		return b.VisitRegister(t)
	case *parser.ImmediateContext:
		return b.VisitImmediate(t)
	case *parser.MemoryOperandContext:
		return b.VisitMemoryOperand(t)
	case *parser.DataListContext:
		return b.VisitDataList(t)
	case *parser.DataItemContext:
		return b.VisitDataItem(t)
	default:
		return tree.Accept(b)
	}
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
			ln := s.(*LineNode)
			label = ln.Label
			statement = ln.Statement
		}
	}
	// Handle optional end-of-line comment
	if n := len(ctx.AllComment()); n > 0 {
		comment = ctx.Comment(n - 1).GetText()
	}
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
	// The directive name is the first token only; GetText() on the whole
	// context would concatenate the parameters into the name (".ORG0x1000").
	name := ctx.GetStart().GetText()
	params := []OperandNode{}
	// The identifier comes before the immediate so .EQU params keep their
	// source order: .EQU NAME, value -> [NAME, value].
	if ctx.IDENTIFIER() != nil {
		params = append(params, &IdentifierNode{
			Name:   ctx.IDENTIFIER().GetText(),
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine(),
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn(),
		})
	}
	if ctx.Immediate() != nil {
		params = append(params, b.Visit(ctx.Immediate()).(OperandNode))
	}
	if ctx.DataList() != nil {
		params = append(params, b.Visit(ctx.DataList()).([]OperandNode)...)
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

// VisitLabelledDirective handles lines of the form `name: .DW 10`, which the
// grammar parses as a single labelledDirective rule. It returns a *LineNode
// carrying both the label and the directive so neither is lost.
func (b *ASTBuilder) VisitLabelledDirective(ctx *parser.LabelledDirectiveContext) interface{} {
	var label *LabelNode
	if ctx.IDENTIFIER() != nil {
		label = &LabelNode{
			Name:   ctx.IDENTIFIER().GetText(),
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine(),
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn(),
		}
	}
	var statement StatementNode
	if ctx.Directive() != nil {
		if s := b.Visit(ctx.Directive()); s != nil {
			statement = s.(*DirectiveNode)
		}
	}
	return &LineNode{
		Label:     label,
		Statement: statement,
		Line:      ctx.GetStart().GetLine(),
		Column:    ctx.GetStart().GetColumn(),
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
	// Handle binary operations (ident+imm, ident+ident, imm+ident, imm+imm).
	// Walk the children in source order so left and right end up on the
	// correct sides regardless of which grammar alternative matched.
	if ctx.PLUS() != nil {
		var parts []OperandNode
		for _, child := range ctx.GetChildren() {
			switch c := child.(type) {
			case *parser.ImmediateContext:
				parts = append(parts, b.VisitImmediate(c).(OperandNode))
			case antlr.TerminalNode:
				if c.GetText() != "+" {
					parts = append(parts, &IdentifierNode{
						Name:   c.GetText(),
						Line:   c.GetSymbol().GetLine(),
						Column: c.GetSymbol().GetColumn(),
					})
				}
			}
		}
		if len(parts) != 2 {
			return nil
		}
		return &BinaryOpNode{
			Left:   parts[0],
			Op:     "+",
			Right:  parts[1],
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
