// Code generated from grammar/vtx1_grammar.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // vtx1_grammar

import "github.com/antlr4-go/antlr/v4"

// Basevtx1_grammarListener is a complete listener for a parse tree produced by vtx1_grammarParser.
type Basevtx1_grammarListener struct{}

var _ vtx1_grammarListener = &Basevtx1_grammarListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *Basevtx1_grammarListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *Basevtx1_grammarListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *Basevtx1_grammarListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *Basevtx1_grammarListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProgram is called when production program is entered.
func (s *Basevtx1_grammarListener) EnterProgram(ctx *ProgramContext) {}

// ExitProgram is called when production program is exited.
func (s *Basevtx1_grammarListener) ExitProgram(ctx *ProgramContext) {}

// EnterLine is called when production line is entered.
func (s *Basevtx1_grammarListener) EnterLine(ctx *LineContext) {}

// ExitLine is called when production line is exited.
func (s *Basevtx1_grammarListener) ExitLine(ctx *LineContext) {}

// EnterBlankLine is called when production blankLine is entered.
func (s *Basevtx1_grammarListener) EnterBlankLine(ctx *BlankLineContext) {}

// ExitBlankLine is called when production blankLine is exited.
func (s *Basevtx1_grammarListener) ExitBlankLine(ctx *BlankLineContext) {}

// EnterLabel is called when production label is entered.
func (s *Basevtx1_grammarListener) EnterLabel(ctx *LabelContext) {}

// ExitLabel is called when production label is exited.
func (s *Basevtx1_grammarListener) ExitLabel(ctx *LabelContext) {}

// EnterComment is called when production comment is entered.
func (s *Basevtx1_grammarListener) EnterComment(ctx *CommentContext) {}

// ExitComment is called when production comment is exited.
func (s *Basevtx1_grammarListener) ExitComment(ctx *CommentContext) {}

// EnterInstruction is called when production instruction is entered.
func (s *Basevtx1_grammarListener) EnterInstruction(ctx *InstructionContext) {}

// ExitInstruction is called when production instruction is exited.
func (s *Basevtx1_grammarListener) ExitInstruction(ctx *InstructionContext) {}

// EnterVliwInstruction is called when production vliwInstruction is entered.
func (s *Basevtx1_grammarListener) EnterVliwInstruction(ctx *VliwInstructionContext) {}

// ExitVliwInstruction is called when production vliwInstruction is exited.
func (s *Basevtx1_grammarListener) ExitVliwInstruction(ctx *VliwInstructionContext) {}

// EnterMnemonic is called when production mnemonic is entered.
func (s *Basevtx1_grammarListener) EnterMnemonic(ctx *MnemonicContext) {}

// ExitMnemonic is called when production mnemonic is exited.
func (s *Basevtx1_grammarListener) ExitMnemonic(ctx *MnemonicContext) {}

// EnterOperand is called when production operand is entered.
func (s *Basevtx1_grammarListener) EnterOperand(ctx *OperandContext) {}

// ExitOperand is called when production operand is exited.
func (s *Basevtx1_grammarListener) ExitOperand(ctx *OperandContext) {}

// EnterRegister is called when production register is entered.
func (s *Basevtx1_grammarListener) EnterRegister(ctx *RegisterContext) {}

// ExitRegister is called when production register is exited.
func (s *Basevtx1_grammarListener) ExitRegister(ctx *RegisterContext) {}

// EnterMemoryOperand is called when production memoryOperand is entered.
func (s *Basevtx1_grammarListener) EnterMemoryOperand(ctx *MemoryOperandContext) {}

// ExitMemoryOperand is called when production memoryOperand is exited.
func (s *Basevtx1_grammarListener) ExitMemoryOperand(ctx *MemoryOperandContext) {}

// EnterBaseRegister is called when production baseRegister is entered.
func (s *Basevtx1_grammarListener) EnterBaseRegister(ctx *BaseRegisterContext) {}

// ExitBaseRegister is called when production baseRegister is exited.
func (s *Basevtx1_grammarListener) ExitBaseRegister(ctx *BaseRegisterContext) {}

// EnterIndexRegister is called when production indexRegister is entered.
func (s *Basevtx1_grammarListener) EnterIndexRegister(ctx *IndexRegisterContext) {}

// ExitIndexRegister is called when production indexRegister is exited.
func (s *Basevtx1_grammarListener) ExitIndexRegister(ctx *IndexRegisterContext) {}

// EnterOffsetImmediate is called when production offsetImmediate is entered.
func (s *Basevtx1_grammarListener) EnterOffsetImmediate(ctx *OffsetImmediateContext) {}

// ExitOffsetImmediate is called when production offsetImmediate is exited.
func (s *Basevtx1_grammarListener) ExitOffsetImmediate(ctx *OffsetImmediateContext) {}

// EnterImmediate is called when production immediate is entered.
func (s *Basevtx1_grammarListener) EnterImmediate(ctx *ImmediateContext) {}

// ExitImmediate is called when production immediate is exited.
func (s *Basevtx1_grammarListener) ExitImmediate(ctx *ImmediateContext) {}

// EnterDirective is called when production directive is entered.
func (s *Basevtx1_grammarListener) EnterDirective(ctx *DirectiveContext) {}

// ExitDirective is called when production directive is exited.
func (s *Basevtx1_grammarListener) ExitDirective(ctx *DirectiveContext) {}

// EnterDataList is called when production dataList is entered.
func (s *Basevtx1_grammarListener) EnterDataList(ctx *DataListContext) {}

// ExitDataList is called when production dataList is exited.
func (s *Basevtx1_grammarListener) ExitDataList(ctx *DataListContext) {}

// EnterMacroDefinition is called when production macroDefinition is entered.
func (s *Basevtx1_grammarListener) EnterMacroDefinition(ctx *MacroDefinitionContext) {}

// ExitMacroDefinition is called when production macroDefinition is exited.
func (s *Basevtx1_grammarListener) ExitMacroDefinition(ctx *MacroDefinitionContext) {}

// EnterMacroBody is called when production macroBody is entered.
func (s *Basevtx1_grammarListener) EnterMacroBody(ctx *MacroBodyContext) {}

// ExitMacroBody is called when production macroBody is exited.
func (s *Basevtx1_grammarListener) ExitMacroBody(ctx *MacroBodyContext) {}
