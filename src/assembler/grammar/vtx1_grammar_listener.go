// Code generated from grammar/vtx1_grammar.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // vtx1_grammar

import "github.com/antlr4-go/antlr/v4"

// vtx1_grammarListener is a complete listener for a parse tree produced by vtx1_grammarParser.
type vtx1_grammarListener interface {
	antlr.ParseTreeListener

	// EnterProgram is called when entering the program production.
	EnterProgram(c *ProgramContext)

	// EnterLine is called when entering the line production.
	EnterLine(c *LineContext)

	// EnterLabel is called when entering the label production.
	EnterLabel(c *LabelContext)

	// EnterComment is called when entering the comment production.
	EnterComment(c *CommentContext)

	// EnterInstruction is called when entering the instruction production.
	EnterInstruction(c *InstructionContext)

	// EnterVliwInstruction is called when entering the vliwInstruction production.
	EnterVliwInstruction(c *VliwInstructionContext)

	// EnterMnemonic is called when entering the mnemonic production.
	EnterMnemonic(c *MnemonicContext)

	// EnterOperand is called when entering the operand production.
	EnterOperand(c *OperandContext)

	// EnterRegister is called when entering the register production.
	EnterRegister(c *RegisterContext)

	// EnterMemoryOperand is called when entering the memoryOperand production.
	EnterMemoryOperand(c *MemoryOperandContext)

	// EnterBaseRegister is called when entering the baseRegister production.
	EnterBaseRegister(c *BaseRegisterContext)

	// EnterIndexRegister is called when entering the indexRegister production.
	EnterIndexRegister(c *IndexRegisterContext)

	// EnterOffsetImmediate is called when entering the offsetImmediate production.
	EnterOffsetImmediate(c *OffsetImmediateContext)

	// EnterImmediate is called when entering the immediate production.
	EnterImmediate(c *ImmediateContext)

	// EnterDirective is called when entering the directive production.
	EnterDirective(c *DirectiveContext)

	// EnterDataList is called when entering the dataList production.
	EnterDataList(c *DataListContext)

	// ExitProgram is called when exiting the program production.
	ExitProgram(c *ProgramContext)

	// ExitLine is called when exiting the line production.
	ExitLine(c *LineContext)

	// ExitLabel is called when exiting the label production.
	ExitLabel(c *LabelContext)

	// ExitComment is called when exiting the comment production.
	ExitComment(c *CommentContext)

	// ExitInstruction is called when exiting the instruction production.
	ExitInstruction(c *InstructionContext)

	// ExitVliwInstruction is called when exiting the vliwInstruction production.
	ExitVliwInstruction(c *VliwInstructionContext)

	// ExitMnemonic is called when exiting the mnemonic production.
	ExitMnemonic(c *MnemonicContext)

	// ExitOperand is called when exiting the operand production.
	ExitOperand(c *OperandContext)

	// ExitRegister is called when exiting the register production.
	ExitRegister(c *RegisterContext)

	// ExitMemoryOperand is called when exiting the memoryOperand production.
	ExitMemoryOperand(c *MemoryOperandContext)

	// ExitBaseRegister is called when exiting the baseRegister production.
	ExitBaseRegister(c *BaseRegisterContext)

	// ExitIndexRegister is called when exiting the indexRegister production.
	ExitIndexRegister(c *IndexRegisterContext)

	// ExitOffsetImmediate is called when exiting the offsetImmediate production.
	ExitOffsetImmediate(c *OffsetImmediateContext)

	// ExitImmediate is called when exiting the immediate production.
	ExitImmediate(c *ImmediateContext)

	// ExitDirective is called when exiting the directive production.
	ExitDirective(c *DirectiveContext)

	// ExitDataList is called when exiting the dataList production.
	ExitDataList(c *DataListContext)
}
