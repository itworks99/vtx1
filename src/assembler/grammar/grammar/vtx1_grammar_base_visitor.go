// Code generated from grammar/vtx1_grammar.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // vtx1_grammar

import "github.com/antlr4-go/antlr/v4"

type Basevtx1_grammarVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *Basevtx1_grammarVisitor) VisitProgram(ctx *ProgramContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitLine(ctx *LineContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitBlankLine(ctx *BlankLineContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitLabel(ctx *LabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitComment(ctx *CommentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitInstruction(ctx *InstructionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitVliwInstruction(ctx *VliwInstructionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitMnemonic(ctx *MnemonicContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitOperand(ctx *OperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitRegister(ctx *RegisterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitMemoryOperand(ctx *MemoryOperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitBaseRegister(ctx *BaseRegisterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitIndexRegister(ctx *IndexRegisterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitOffsetImmediate(ctx *OffsetImmediateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitImmediate(ctx *ImmediateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitDirective(ctx *DirectiveContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitDataList(ctx *DataListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitMacroDefinition(ctx *MacroDefinitionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *Basevtx1_grammarVisitor) VisitMacroBody(ctx *MacroBodyContext) interface{} {
	return v.VisitChildren(ctx)
}
