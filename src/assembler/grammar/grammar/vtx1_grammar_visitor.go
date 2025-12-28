// Code generated from grammar/vtx1_grammar.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // vtx1_grammar

import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by vtx1_grammarParser.
type vtx1_grammarVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by vtx1_grammarParser#program.
	VisitProgram(ctx *ProgramContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#line.
	VisitLine(ctx *LineContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#blankLine.
	VisitBlankLine(ctx *BlankLineContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#label.
	VisitLabel(ctx *LabelContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#comment.
	VisitComment(ctx *CommentContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#instruction.
	VisitInstruction(ctx *InstructionContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#vliwInstruction.
	VisitVliwInstruction(ctx *VliwInstructionContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#mnemonic.
	VisitMnemonic(ctx *MnemonicContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#operand.
	VisitOperand(ctx *OperandContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#register.
	VisitRegister(ctx *RegisterContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#memoryOperand.
	VisitMemoryOperand(ctx *MemoryOperandContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#baseRegister.
	VisitBaseRegister(ctx *BaseRegisterContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#indexRegister.
	VisitIndexRegister(ctx *IndexRegisterContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#offsetImmediate.
	VisitOffsetImmediate(ctx *OffsetImmediateContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#immediate.
	VisitImmediate(ctx *ImmediateContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#directive.
	VisitDirective(ctx *DirectiveContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#dataList.
	VisitDataList(ctx *DataListContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#macroDefinition.
	VisitMacroDefinition(ctx *MacroDefinitionContext) interface{}

	// Visit a parse tree produced by vtx1_grammarParser#macroBody.
	VisitMacroBody(ctx *MacroBodyContext) interface{}
}
