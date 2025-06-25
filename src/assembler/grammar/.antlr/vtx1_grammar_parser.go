// Code generated from /home/itworks/Projects/vtx1/src/assembler/grammar/vtx1_grammar.g4 by ANTLR 4.9.2. DO NOT EDIT.

package parser // vtx1_grammar

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 38, 142,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 3, 2, 7,
	2, 36, 10, 2, 12, 2, 14, 2, 39, 11, 2, 3, 2, 3, 2, 3, 3, 5, 3, 44, 10,
	3, 3, 3, 3, 3, 3, 3, 5, 3, 49, 10, 3, 3, 3, 5, 3, 52, 10, 3, 3, 3, 3, 3,
	3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 7, 6, 65, 10, 6,
	12, 6, 14, 6, 68, 11, 6, 5, 6, 70, 10, 6, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7,
	3, 7, 3, 7, 5, 7, 79, 10, 7, 3, 7, 3, 7, 3, 7, 3, 7, 5, 7, 85, 10, 7, 3,
	8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 9, 5, 9, 93, 10, 9, 3, 10, 3, 10, 3, 11,
	3, 11, 3, 11, 3, 11, 3, 11, 5, 11, 102, 10, 11, 5, 11, 104, 10, 11, 3,
	11, 3, 11, 3, 12, 3, 12, 3, 13, 3, 13, 3, 14, 3, 14, 3, 15, 3, 15, 3, 16,
	3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 3,
	16, 3, 16, 3, 16, 3, 16, 3, 16, 5, 16, 132, 10, 16, 3, 17, 3, 17, 3, 17,
	7, 17, 137, 10, 17, 12, 17, 14, 17, 140, 11, 17, 3, 17, 2, 2, 18, 2, 4,
	6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 2, 6, 3, 2, 6, 15,
	4, 2, 16, 16, 18, 20, 3, 2, 16, 17, 3, 2, 21, 24, 2, 147, 2, 37, 3, 2,
	2, 2, 4, 43, 3, 2, 2, 2, 6, 55, 3, 2, 2, 2, 8, 58, 3, 2, 2, 2, 10, 60,
	3, 2, 2, 2, 12, 71, 3, 2, 2, 2, 14, 86, 3, 2, 2, 2, 16, 92, 3, 2, 2, 2,
	18, 94, 3, 2, 2, 2, 20, 96, 3, 2, 2, 2, 22, 107, 3, 2, 2, 2, 24, 109, 3,
	2, 2, 2, 26, 111, 3, 2, 2, 2, 28, 113, 3, 2, 2, 2, 30, 131, 3, 2, 2, 2,
	32, 133, 3, 2, 2, 2, 34, 36, 5, 4, 3, 2, 35, 34, 3, 2, 2, 2, 36, 39, 3,
	2, 2, 2, 37, 35, 3, 2, 2, 2, 37, 38, 3, 2, 2, 2, 38, 40, 3, 2, 2, 2, 39,
	37, 3, 2, 2, 2, 40, 41, 7, 2, 2, 3, 41, 3, 3, 2, 2, 2, 42, 44, 5, 6, 4,
	2, 43, 42, 3, 2, 2, 2, 43, 44, 3, 2, 2, 2, 44, 48, 3, 2, 2, 2, 45, 49,
	5, 10, 6, 2, 46, 49, 5, 30, 16, 2, 47, 49, 5, 12, 7, 2, 48, 45, 3, 2, 2,
	2, 48, 46, 3, 2, 2, 2, 48, 47, 3, 2, 2, 2, 48, 49, 3, 2, 2, 2, 49, 51,
	3, 2, 2, 2, 50, 52, 5, 8, 5, 2, 51, 50, 3, 2, 2, 2, 51, 52, 3, 2, 2, 2,
	52, 53, 3, 2, 2, 2, 53, 54, 7, 5, 2, 2, 54, 5, 3, 2, 2, 2, 55, 56, 7, 38,
	2, 2, 56, 57, 7, 26, 2, 2, 57, 7, 3, 2, 2, 2, 58, 59, 7, 4, 2, 2, 59, 9,
	3, 2, 2, 2, 60, 69, 5, 14, 8, 2, 61, 66, 5, 16, 9, 2, 62, 63, 7, 27, 2,
	2, 63, 65, 5, 16, 9, 2, 64, 62, 3, 2, 2, 2, 65, 68, 3, 2, 2, 2, 66, 64,
	3, 2, 2, 2, 66, 67, 3, 2, 2, 2, 67, 70, 3, 2, 2, 2, 68, 66, 3, 2, 2, 2,
	69, 61, 3, 2, 2, 2, 69, 70, 3, 2, 2, 2, 70, 11, 3, 2, 2, 2, 71, 72, 7,
	29, 2, 2, 72, 73, 5, 10, 6, 2, 73, 78, 7, 30, 2, 2, 74, 75, 7, 29, 2, 2,
	75, 76, 5, 10, 6, 2, 76, 77, 7, 30, 2, 2, 77, 79, 3, 2, 2, 2, 78, 74, 3,
	2, 2, 2, 78, 79, 3, 2, 2, 2, 79, 84, 3, 2, 2, 2, 80, 81, 7, 29, 2, 2, 81,
	82, 5, 10, 6, 2, 82, 83, 7, 30, 2, 2, 83, 85, 3, 2, 2, 2, 84, 80, 3, 2,
	2, 2, 84, 85, 3, 2, 2, 2, 85, 13, 3, 2, 2, 2, 86, 87, 9, 2, 2, 2, 87, 15,
	3, 2, 2, 2, 88, 93, 5, 18, 10, 2, 89, 93, 5, 28, 15, 2, 90, 93, 5, 20,
	11, 2, 91, 93, 7, 38, 2, 2, 92, 88, 3, 2, 2, 2, 92, 89, 3, 2, 2, 2, 92,
	90, 3, 2, 2, 2, 92, 91, 3, 2, 2, 2, 93, 17, 3, 2, 2, 2, 94, 95, 9, 3, 2,
	2, 95, 19, 3, 2, 2, 2, 96, 97, 7, 29, 2, 2, 97, 103, 5, 22, 12, 2, 98,
	101, 7, 28, 2, 2, 99, 102, 5, 24, 13, 2, 100, 102, 5, 26, 14, 2, 101, 99,
	3, 2, 2, 2, 101, 100, 3, 2, 2, 2, 102, 104, 3, 2, 2, 2, 103, 98, 3, 2,
	2, 2, 103, 104, 3, 2, 2, 2, 104, 105, 3, 2, 2, 2, 105, 106, 7, 30, 2, 2,
	106, 21, 3, 2, 2, 2, 107, 108, 9, 4, 2, 2, 108, 23, 3, 2, 2, 2, 109, 110,
	7, 16, 2, 2, 110, 25, 3, 2, 2, 2, 111, 112, 5, 28, 15, 2, 112, 27, 3, 2,
	2, 2, 113, 114, 9, 5, 2, 2, 114, 29, 3, 2, 2, 2, 115, 116, 7, 31, 2, 2,
	116, 132, 5, 28, 15, 2, 117, 118, 7, 32, 2, 2, 118, 132, 5, 32, 17, 2,
	119, 120, 7, 33, 2, 2, 120, 121, 7, 38, 2, 2, 121, 122, 7, 27, 2, 2, 122,
	132, 5, 28, 15, 2, 123, 124, 7, 34, 2, 2, 124, 132, 7, 25, 2, 2, 125, 126,
	7, 35, 2, 2, 126, 132, 7, 38, 2, 2, 127, 128, 7, 36, 2, 2, 128, 132, 5,
	28, 15, 2, 129, 130, 7, 37, 2, 2, 130, 132, 5, 28, 15, 2, 131, 115, 3,
	2, 2, 2, 131, 117, 3, 2, 2, 2, 131, 119, 3, 2, 2, 2, 131, 123, 3, 2, 2,
	2, 131, 125, 3, 2, 2, 2, 131, 127, 3, 2, 2, 2, 131, 129, 3, 2, 2, 2, 132,
	31, 3, 2, 2, 2, 133, 138, 5, 28, 15, 2, 134, 135, 7, 27, 2, 2, 135, 137,
	5, 28, 15, 2, 136, 134, 3, 2, 2, 2, 137, 140, 3, 2, 2, 2, 138, 136, 3,
	2, 2, 2, 138, 139, 3, 2, 2, 2, 139, 33, 3, 2, 2, 2, 140, 138, 3, 2, 2,
	2, 15, 37, 43, 48, 51, 66, 69, 78, 84, 92, 101, 103, 131, 138,
}
var literalNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "'TB'", "",
	"", "", "", "", "", "", "", "':'", "','", "'+'", "'['", "']'", "'.ORG'",
	"", "'.EQU'", "'.INCLUDE'", "'.SECTION'", "'.ALIGN'", "'.SPACE'",
}
var symbolicNames = []string{
	"", "WHITESPACE", "COMMENT", "EOL", "ALU_OP", "MEM_OP", "CTRL_OP", "VEC_OP",
	"FP_OP", "SYS_OP", "COMPLEX_OP", "COMPLEX_VEC", "COMPLEX_MEM", "COMPLEX_SYS",
	"GPR", "TB_REG", "SPECIAL_REG", "VECTOR_REG", "FP_REG", "DECIMAL", "HEXADECIMAL",
	"BINARY", "TERNARY", "STRING", "COLON", "COMMA", "PLUS", "LSQUARE", "RSQUARE",
	"ORG_DIRECTIVE", "DATA_DIRECTIVE", "EQU_DIRECTIVE", "INCLUDE_DIRECTIVE",
	"SECTION_DIRECTIVE", "ALIGN_DIRECTIVE", "SPACE_DIRECTIVE", "IDENTIFIER",
}

var ruleNames = []string{
	"program", "line", "label", "comment", "instruction", "vliwInstruction",
	"mnemonic", "operand", "register", "memoryOperand", "baseRegister", "indexRegister",
	"offsetImmediate", "immediate", "directive", "dataList",
}

type vtx1_grammarParser struct {
	*antlr.BaseParser
}

// Newvtx1_grammarParser produces a new parser instance for the optional input antlr.TokenStream.
//
// The *vtx1_grammarParser instance produced may be reused by calling the SetInputStream method.
// The initial parser configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func Newvtx1_grammarParser(input antlr.TokenStream) *vtx1_grammarParser {
	this := new(vtx1_grammarParser)
	deserializer := antlr.NewATNDeserializer(nil)
	deserializedATN := deserializer.DeserializeFromUInt16(parserATN)
	decisionToDFA := make([]*antlr.DFA, len(deserializedATN.DecisionToState))
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "vtx1_grammar.g4"

	return this
}

// vtx1_grammarParser tokens.
const (
	vtx1_grammarParserEOF               = antlr.TokenEOF
	vtx1_grammarParserWHITESPACE        = 1
	vtx1_grammarParserCOMMENT           = 2
	vtx1_grammarParserEOL               = 3
	vtx1_grammarParserALU_OP            = 4
	vtx1_grammarParserMEM_OP            = 5
	vtx1_grammarParserCTRL_OP           = 6
	vtx1_grammarParserVEC_OP            = 7
	vtx1_grammarParserFP_OP             = 8
	vtx1_grammarParserSYS_OP            = 9
	vtx1_grammarParserCOMPLEX_OP        = 10
	vtx1_grammarParserCOMPLEX_VEC       = 11
	vtx1_grammarParserCOMPLEX_MEM       = 12
	vtx1_grammarParserCOMPLEX_SYS       = 13
	vtx1_grammarParserGPR               = 14
	vtx1_grammarParserTB_REG            = 15
	vtx1_grammarParserSPECIAL_REG       = 16
	vtx1_grammarParserVECTOR_REG        = 17
	vtx1_grammarParserFP_REG            = 18
	vtx1_grammarParserDECIMAL           = 19
	vtx1_grammarParserHEXADECIMAL       = 20
	vtx1_grammarParserBINARY            = 21
	vtx1_grammarParserTERNARY           = 22
	vtx1_grammarParserSTRING            = 23
	vtx1_grammarParserCOLON             = 24
	vtx1_grammarParserCOMMA             = 25
	vtx1_grammarParserPLUS              = 26
	vtx1_grammarParserLSQUARE           = 27
	vtx1_grammarParserRSQUARE           = 28
	vtx1_grammarParserORG_DIRECTIVE     = 29
	vtx1_grammarParserDATA_DIRECTIVE    = 30
	vtx1_grammarParserEQU_DIRECTIVE     = 31
	vtx1_grammarParserINCLUDE_DIRECTIVE = 32
	vtx1_grammarParserSECTION_DIRECTIVE = 33
	vtx1_grammarParserALIGN_DIRECTIVE   = 34
	vtx1_grammarParserSPACE_DIRECTIVE   = 35
	vtx1_grammarParserIDENTIFIER        = 36
)

// vtx1_grammarParser rules.
const (
	vtx1_grammarParserRULE_program         = 0
	vtx1_grammarParserRULE_line            = 1
	vtx1_grammarParserRULE_label           = 2
	vtx1_grammarParserRULE_comment         = 3
	vtx1_grammarParserRULE_instruction     = 4
	vtx1_grammarParserRULE_vliwInstruction = 5
	vtx1_grammarParserRULE_mnemonic        = 6
	vtx1_grammarParserRULE_operand         = 7
	vtx1_grammarParserRULE_register        = 8
	vtx1_grammarParserRULE_memoryOperand   = 9
	vtx1_grammarParserRULE_baseRegister    = 10
	vtx1_grammarParserRULE_indexRegister   = 11
	vtx1_grammarParserRULE_offsetImmediate = 12
	vtx1_grammarParserRULE_immediate       = 13
	vtx1_grammarParserRULE_directive       = 14
	vtx1_grammarParserRULE_dataList        = 15
)

// IProgramContext is an interface to support dynamic dispatch.
type IProgramContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsProgramContext differentiates from other interfaces.
	IsProgramContext()
}

type ProgramContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgramContext() *ProgramContext {
	var p = new(ProgramContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_program
	return p
}

func (*ProgramContext) IsProgramContext() {}

func NewProgramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramContext {
	var p = new(ProgramContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_program

	return p
}

func (s *ProgramContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramContext) EOF() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserEOF, 0)
}

func (s *ProgramContext) AllLine() []ILineContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ILineContext)(nil)).Elem())
	var tst = make([]ILineContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ILineContext)
		}
	}

	return tst
}

func (s *ProgramContext) Line(i int) ILineContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILineContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ILineContext)
}

func (s *ProgramContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgramContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) Program() (localctx IProgramContext) {
	localctx = NewProgramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, vtx1_grammarParserRULE_program)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(35)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<vtx1_grammarParserCOMMENT)|(1<<vtx1_grammarParserEOL)|(1<<vtx1_grammarParserALU_OP)|(1<<vtx1_grammarParserMEM_OP)|(1<<vtx1_grammarParserCTRL_OP)|(1<<vtx1_grammarParserVEC_OP)|(1<<vtx1_grammarParserFP_OP)|(1<<vtx1_grammarParserSYS_OP)|(1<<vtx1_grammarParserCOMPLEX_OP)|(1<<vtx1_grammarParserCOMPLEX_VEC)|(1<<vtx1_grammarParserCOMPLEX_MEM)|(1<<vtx1_grammarParserCOMPLEX_SYS)|(1<<vtx1_grammarParserLSQUARE)|(1<<vtx1_grammarParserORG_DIRECTIVE)|(1<<vtx1_grammarParserDATA_DIRECTIVE)|(1<<vtx1_grammarParserEQU_DIRECTIVE))) != 0) || (((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(vtx1_grammarParserINCLUDE_DIRECTIVE-32))|(1<<(vtx1_grammarParserSECTION_DIRECTIVE-32))|(1<<(vtx1_grammarParserALIGN_DIRECTIVE-32))|(1<<(vtx1_grammarParserSPACE_DIRECTIVE-32))|(1<<(vtx1_grammarParserIDENTIFIER-32)))) != 0) {
		{
			p.SetState(32)
			p.Line()
		}

		p.SetState(37)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(38)
		p.Match(vtx1_grammarParserEOF)
	}

	return localctx
}

// ILineContext is an interface to support dynamic dispatch.
type ILineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLineContext differentiates from other interfaces.
	IsLineContext()
}

type LineContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLineContext() *LineContext {
	var p = new(LineContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_line
	return p
}

func (*LineContext) IsLineContext() {}

func NewLineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LineContext {
	var p = new(LineContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_line

	return p
}

func (s *LineContext) GetParser() antlr.Parser { return s.parser }

func (s *LineContext) EOL() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserEOL, 0)
}

func (s *LineContext) Label() ILabelContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILabelContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILabelContext)
}

func (s *LineContext) Instruction() IInstructionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInstructionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInstructionContext)
}

func (s *LineContext) Directive() IDirectiveContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDirectiveContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDirectiveContext)
}

func (s *LineContext) VliwInstruction() IVliwInstructionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVliwInstructionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVliwInstructionContext)
}

func (s *LineContext) Comment() ICommentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICommentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICommentContext)
}

func (s *LineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) Line() (localctx ILineContext) {
	localctx = NewLineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, vtx1_grammarParserRULE_line)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(41)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserIDENTIFIER {
		{
			p.SetState(40)
			p.Label()
		}

	}
	p.SetState(46)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case vtx1_grammarParserALU_OP, vtx1_grammarParserMEM_OP, vtx1_grammarParserCTRL_OP, vtx1_grammarParserVEC_OP, vtx1_grammarParserFP_OP, vtx1_grammarParserSYS_OP, vtx1_grammarParserCOMPLEX_OP, vtx1_grammarParserCOMPLEX_VEC, vtx1_grammarParserCOMPLEX_MEM, vtx1_grammarParserCOMPLEX_SYS:
		{
			p.SetState(43)
			p.Instruction()
		}

	case vtx1_grammarParserORG_DIRECTIVE, vtx1_grammarParserDATA_DIRECTIVE, vtx1_grammarParserEQU_DIRECTIVE, vtx1_grammarParserINCLUDE_DIRECTIVE, vtx1_grammarParserSECTION_DIRECTIVE, vtx1_grammarParserALIGN_DIRECTIVE, vtx1_grammarParserSPACE_DIRECTIVE:
		{
			p.SetState(44)
			p.Directive()
		}

	case vtx1_grammarParserLSQUARE:
		{
			p.SetState(45)
			p.VliwInstruction()
		}

	case vtx1_grammarParserCOMMENT, vtx1_grammarParserEOL:

	default:
	}
	p.SetState(49)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserCOMMENT {
		{
			p.SetState(48)
			p.Comment()
		}

	}
	{
		p.SetState(51)
		p.Match(vtx1_grammarParserEOL)
	}

	return localctx
}

// ILabelContext is an interface to support dynamic dispatch.
type ILabelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLabelContext differentiates from other interfaces.
	IsLabelContext()
}

type LabelContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLabelContext() *LabelContext {
	var p = new(LabelContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_label
	return p
}

func (*LabelContext) IsLabelContext() {}

func NewLabelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LabelContext {
	var p = new(LabelContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_label

	return p
}

func (s *LabelContext) GetParser() antlr.Parser { return s.parser }

func (s *LabelContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserIDENTIFIER, 0)
}

func (s *LabelContext) COLON() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCOLON, 0)
}

func (s *LabelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LabelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) Label() (localctx ILabelContext) {
	localctx = NewLabelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, vtx1_grammarParserRULE_label)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(53)
		p.Match(vtx1_grammarParserIDENTIFIER)
	}
	{
		p.SetState(54)
		p.Match(vtx1_grammarParserCOLON)
	}

	return localctx
}

// ICommentContext is an interface to support dynamic dispatch.
type ICommentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCommentContext differentiates from other interfaces.
	IsCommentContext()
}

type CommentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommentContext() *CommentContext {
	var p = new(CommentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_comment
	return p
}

func (*CommentContext) IsCommentContext() {}

func NewCommentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommentContext {
	var p = new(CommentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_comment

	return p
}

func (s *CommentContext) GetParser() antlr.Parser { return s.parser }

func (s *CommentContext) COMMENT() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCOMMENT, 0)
}

func (s *CommentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) Comment() (localctx ICommentContext) {
	localctx = NewCommentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, vtx1_grammarParserRULE_comment)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(56)
		p.Match(vtx1_grammarParserCOMMENT)
	}

	return localctx
}

// IInstructionContext is an interface to support dynamic dispatch.
type IInstructionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInstructionContext differentiates from other interfaces.
	IsInstructionContext()
}

type InstructionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInstructionContext() *InstructionContext {
	var p = new(InstructionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_instruction
	return p
}

func (*InstructionContext) IsInstructionContext() {}

func NewInstructionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InstructionContext {
	var p = new(InstructionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_instruction

	return p
}

func (s *InstructionContext) GetParser() antlr.Parser { return s.parser }

func (s *InstructionContext) Mnemonic() IMnemonicContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMnemonicContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMnemonicContext)
}

func (s *InstructionContext) AllOperand() []IOperandContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IOperandContext)(nil)).Elem())
	var tst = make([]IOperandContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IOperandContext)
		}
	}

	return tst
}

func (s *InstructionContext) Operand(i int) IOperandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOperandContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IOperandContext)
}

func (s *InstructionContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(vtx1_grammarParserCOMMA)
}

func (s *InstructionContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCOMMA, i)
}

func (s *InstructionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InstructionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) Instruction() (localctx IInstructionContext) {
	localctx = NewInstructionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, vtx1_grammarParserRULE_instruction)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(58)
		p.Mnemonic()
	}
	p.SetState(67)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if ((_la-14)&-(0x1f+1)) == 0 && ((1<<uint((_la-14)))&((1<<(vtx1_grammarParserGPR-14))|(1<<(vtx1_grammarParserSPECIAL_REG-14))|(1<<(vtx1_grammarParserVECTOR_REG-14))|(1<<(vtx1_grammarParserFP_REG-14))|(1<<(vtx1_grammarParserDECIMAL-14))|(1<<(vtx1_grammarParserHEXADECIMAL-14))|(1<<(vtx1_grammarParserBINARY-14))|(1<<(vtx1_grammarParserTERNARY-14))|(1<<(vtx1_grammarParserLSQUARE-14))|(1<<(vtx1_grammarParserIDENTIFIER-14)))) != 0 {
		{
			p.SetState(59)
			p.Operand()
		}
		p.SetState(64)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == vtx1_grammarParserCOMMA {
			{
				p.SetState(60)
				p.Match(vtx1_grammarParserCOMMA)
			}
			{
				p.SetState(61)
				p.Operand()
			}

			p.SetState(66)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	}

	return localctx
}

// IVliwInstructionContext is an interface to support dynamic dispatch.
type IVliwInstructionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsVliwInstructionContext differentiates from other interfaces.
	IsVliwInstructionContext()
}

type VliwInstructionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVliwInstructionContext() *VliwInstructionContext {
	var p = new(VliwInstructionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_vliwInstruction
	return p
}

func (*VliwInstructionContext) IsVliwInstructionContext() {}

func NewVliwInstructionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VliwInstructionContext {
	var p = new(VliwInstructionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_vliwInstruction

	return p
}

func (s *VliwInstructionContext) GetParser() antlr.Parser { return s.parser }

func (s *VliwInstructionContext) AllLSQUARE() []antlr.TerminalNode {
	return s.GetTokens(vtx1_grammarParserLSQUARE)
}

func (s *VliwInstructionContext) LSQUARE(i int) antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserLSQUARE, i)
}

func (s *VliwInstructionContext) AllInstruction() []IInstructionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IInstructionContext)(nil)).Elem())
	var tst = make([]IInstructionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IInstructionContext)
		}
	}

	return tst
}

func (s *VliwInstructionContext) Instruction(i int) IInstructionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInstructionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IInstructionContext)
}

func (s *VliwInstructionContext) AllRSQUARE() []antlr.TerminalNode {
	return s.GetTokens(vtx1_grammarParserRSQUARE)
}

func (s *VliwInstructionContext) RSQUARE(i int) antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserRSQUARE, i)
}

func (s *VliwInstructionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VliwInstructionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) VliwInstruction() (localctx IVliwInstructionContext) {
	localctx = NewVliwInstructionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, vtx1_grammarParserRULE_vliwInstruction)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(69)
		p.Match(vtx1_grammarParserLSQUARE)
	}
	{
		p.SetState(70)
		p.Instruction()
	}
	{
		p.SetState(71)
		p.Match(vtx1_grammarParserRSQUARE)
	}
	p.SetState(76)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(72)
			p.Match(vtx1_grammarParserLSQUARE)
		}
		{
			p.SetState(73)
			p.Instruction()
		}
		{
			p.SetState(74)
			p.Match(vtx1_grammarParserRSQUARE)
		}

	}
	p.SetState(82)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserLSQUARE {
		{
			p.SetState(78)
			p.Match(vtx1_grammarParserLSQUARE)
		}
		{
			p.SetState(79)
			p.Instruction()
		}
		{
			p.SetState(80)
			p.Match(vtx1_grammarParserRSQUARE)
		}

	}

	return localctx
}

// IMnemonicContext is an interface to support dynamic dispatch.
type IMnemonicContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMnemonicContext differentiates from other interfaces.
	IsMnemonicContext()
}

type MnemonicContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMnemonicContext() *MnemonicContext {
	var p = new(MnemonicContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_mnemonic
	return p
}

func (*MnemonicContext) IsMnemonicContext() {}

func NewMnemonicContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MnemonicContext {
	var p = new(MnemonicContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_mnemonic

	return p
}

func (s *MnemonicContext) GetParser() antlr.Parser { return s.parser }

func (s *MnemonicContext) ALU_OP() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserALU_OP, 0)
}

func (s *MnemonicContext) MEM_OP() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserMEM_OP, 0)
}

func (s *MnemonicContext) CTRL_OP() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCTRL_OP, 0)
}

func (s *MnemonicContext) VEC_OP() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserVEC_OP, 0)
}

func (s *MnemonicContext) FP_OP() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserFP_OP, 0)
}

func (s *MnemonicContext) SYS_OP() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserSYS_OP, 0)
}

func (s *MnemonicContext) COMPLEX_OP() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCOMPLEX_OP, 0)
}

func (s *MnemonicContext) COMPLEX_VEC() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCOMPLEX_VEC, 0)
}

func (s *MnemonicContext) COMPLEX_MEM() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCOMPLEX_MEM, 0)
}

func (s *MnemonicContext) COMPLEX_SYS() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCOMPLEX_SYS, 0)
}

func (s *MnemonicContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MnemonicContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) Mnemonic() (localctx IMnemonicContext) {
	localctx = NewMnemonicContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, vtx1_grammarParserRULE_mnemonic)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(84)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<vtx1_grammarParserALU_OP)|(1<<vtx1_grammarParserMEM_OP)|(1<<vtx1_grammarParserCTRL_OP)|(1<<vtx1_grammarParserVEC_OP)|(1<<vtx1_grammarParserFP_OP)|(1<<vtx1_grammarParserSYS_OP)|(1<<vtx1_grammarParserCOMPLEX_OP)|(1<<vtx1_grammarParserCOMPLEX_VEC)|(1<<vtx1_grammarParserCOMPLEX_MEM)|(1<<vtx1_grammarParserCOMPLEX_SYS))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IOperandContext is an interface to support dynamic dispatch.
type IOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOperandContext differentiates from other interfaces.
	IsOperandContext()
}

type OperandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperandContext() *OperandContext {
	var p = new(OperandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_operand
	return p
}

func (*OperandContext) IsOperandContext() {}

func NewOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperandContext {
	var p = new(OperandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_operand

	return p
}

func (s *OperandContext) GetParser() antlr.Parser { return s.parser }

func (s *OperandContext) Register() IRegisterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRegisterContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRegisterContext)
}

func (s *OperandContext) Immediate() IImmediateContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImmediateContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImmediateContext)
}

func (s *OperandContext) MemoryOperand() IMemoryOperandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMemoryOperandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMemoryOperandContext)
}

func (s *OperandContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserIDENTIFIER, 0)
}

func (s *OperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) Operand() (localctx IOperandContext) {
	localctx = NewOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, vtx1_grammarParserRULE_operand)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(90)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case vtx1_grammarParserGPR, vtx1_grammarParserSPECIAL_REG, vtx1_grammarParserVECTOR_REG, vtx1_grammarParserFP_REG:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(86)
			p.Register()
		}

	case vtx1_grammarParserDECIMAL, vtx1_grammarParserHEXADECIMAL, vtx1_grammarParserBINARY, vtx1_grammarParserTERNARY:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(87)
			p.Immediate()
		}

	case vtx1_grammarParserLSQUARE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(88)
			p.MemoryOperand()
		}

	case vtx1_grammarParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(89)
			p.Match(vtx1_grammarParserIDENTIFIER)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IRegisterContext is an interface to support dynamic dispatch.
type IRegisterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRegisterContext differentiates from other interfaces.
	IsRegisterContext()
}

type RegisterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRegisterContext() *RegisterContext {
	var p = new(RegisterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_register
	return p
}

func (*RegisterContext) IsRegisterContext() {}

func NewRegisterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RegisterContext {
	var p = new(RegisterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_register

	return p
}

func (s *RegisterContext) GetParser() antlr.Parser { return s.parser }

func (s *RegisterContext) GPR() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserGPR, 0)
}

func (s *RegisterContext) SPECIAL_REG() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserSPECIAL_REG, 0)
}

func (s *RegisterContext) VECTOR_REG() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserVECTOR_REG, 0)
}

func (s *RegisterContext) FP_REG() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserFP_REG, 0)
}

func (s *RegisterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RegisterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) Register() (localctx IRegisterContext) {
	localctx = NewRegisterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, vtx1_grammarParserRULE_register)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(92)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<vtx1_grammarParserGPR)|(1<<vtx1_grammarParserSPECIAL_REG)|(1<<vtx1_grammarParserVECTOR_REG)|(1<<vtx1_grammarParserFP_REG))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IMemoryOperandContext is an interface to support dynamic dispatch.
type IMemoryOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMemoryOperandContext differentiates from other interfaces.
	IsMemoryOperandContext()
}

type MemoryOperandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemoryOperandContext() *MemoryOperandContext {
	var p = new(MemoryOperandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_memoryOperand
	return p
}

func (*MemoryOperandContext) IsMemoryOperandContext() {}

func NewMemoryOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemoryOperandContext {
	var p = new(MemoryOperandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_memoryOperand

	return p
}

func (s *MemoryOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *MemoryOperandContext) LSQUARE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserLSQUARE, 0)
}

func (s *MemoryOperandContext) BaseRegister() IBaseRegisterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBaseRegisterContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBaseRegisterContext)
}

func (s *MemoryOperandContext) RSQUARE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserRSQUARE, 0)
}

func (s *MemoryOperandContext) PLUS() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserPLUS, 0)
}

func (s *MemoryOperandContext) IndexRegister() IIndexRegisterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIndexRegisterContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIndexRegisterContext)
}

func (s *MemoryOperandContext) OffsetImmediate() IOffsetImmediateContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOffsetImmediateContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOffsetImmediateContext)
}

func (s *MemoryOperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemoryOperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) MemoryOperand() (localctx IMemoryOperandContext) {
	localctx = NewMemoryOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, vtx1_grammarParserRULE_memoryOperand)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(94)
		p.Match(vtx1_grammarParserLSQUARE)
	}
	{
		p.SetState(95)
		p.BaseRegister()
	}
	p.SetState(101)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserPLUS {
		{
			p.SetState(96)
			p.Match(vtx1_grammarParserPLUS)
		}
		p.SetState(99)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case vtx1_grammarParserGPR:
			{
				p.SetState(97)
				p.IndexRegister()
			}

		case vtx1_grammarParserDECIMAL, vtx1_grammarParserHEXADECIMAL, vtx1_grammarParserBINARY, vtx1_grammarParserTERNARY:
			{
				p.SetState(98)
				p.OffsetImmediate()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

	}
	{
		p.SetState(103)
		p.Match(vtx1_grammarParserRSQUARE)
	}

	return localctx
}

// IBaseRegisterContext is an interface to support dynamic dispatch.
type IBaseRegisterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBaseRegisterContext differentiates from other interfaces.
	IsBaseRegisterContext()
}

type BaseRegisterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBaseRegisterContext() *BaseRegisterContext {
	var p = new(BaseRegisterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_baseRegister
	return p
}

func (*BaseRegisterContext) IsBaseRegisterContext() {}

func NewBaseRegisterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BaseRegisterContext {
	var p = new(BaseRegisterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_baseRegister

	return p
}

func (s *BaseRegisterContext) GetParser() antlr.Parser { return s.parser }

func (s *BaseRegisterContext) TB_REG() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserTB_REG, 0)
}

func (s *BaseRegisterContext) GPR() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserGPR, 0)
}

func (s *BaseRegisterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BaseRegisterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) BaseRegister() (localctx IBaseRegisterContext) {
	localctx = NewBaseRegisterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, vtx1_grammarParserRULE_baseRegister)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(105)
		_la = p.GetTokenStream().LA(1)

		if !(_la == vtx1_grammarParserGPR || _la == vtx1_grammarParserTB_REG) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IIndexRegisterContext is an interface to support dynamic dispatch.
type IIndexRegisterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIndexRegisterContext differentiates from other interfaces.
	IsIndexRegisterContext()
}

type IndexRegisterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIndexRegisterContext() *IndexRegisterContext {
	var p = new(IndexRegisterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_indexRegister
	return p
}

func (*IndexRegisterContext) IsIndexRegisterContext() {}

func NewIndexRegisterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IndexRegisterContext {
	var p = new(IndexRegisterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_indexRegister

	return p
}

func (s *IndexRegisterContext) GetParser() antlr.Parser { return s.parser }

func (s *IndexRegisterContext) GPR() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserGPR, 0)
}

func (s *IndexRegisterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IndexRegisterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) IndexRegister() (localctx IIndexRegisterContext) {
	localctx = NewIndexRegisterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, vtx1_grammarParserRULE_indexRegister)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(107)
		p.Match(vtx1_grammarParserGPR)
	}

	return localctx
}

// IOffsetImmediateContext is an interface to support dynamic dispatch.
type IOffsetImmediateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOffsetImmediateContext differentiates from other interfaces.
	IsOffsetImmediateContext()
}

type OffsetImmediateContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOffsetImmediateContext() *OffsetImmediateContext {
	var p = new(OffsetImmediateContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_offsetImmediate
	return p
}

func (*OffsetImmediateContext) IsOffsetImmediateContext() {}

func NewOffsetImmediateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OffsetImmediateContext {
	var p = new(OffsetImmediateContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_offsetImmediate

	return p
}

func (s *OffsetImmediateContext) GetParser() antlr.Parser { return s.parser }

func (s *OffsetImmediateContext) Immediate() IImmediateContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImmediateContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImmediateContext)
}

func (s *OffsetImmediateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OffsetImmediateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) OffsetImmediate() (localctx IOffsetImmediateContext) {
	localctx = NewOffsetImmediateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, vtx1_grammarParserRULE_offsetImmediate)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(109)
		p.Immediate()
	}

	return localctx
}

// IImmediateContext is an interface to support dynamic dispatch.
type IImmediateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImmediateContext differentiates from other interfaces.
	IsImmediateContext()
}

type ImmediateContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImmediateContext() *ImmediateContext {
	var p = new(ImmediateContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_immediate
	return p
}

func (*ImmediateContext) IsImmediateContext() {}

func NewImmediateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImmediateContext {
	var p = new(ImmediateContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_immediate

	return p
}

func (s *ImmediateContext) GetParser() antlr.Parser { return s.parser }

func (s *ImmediateContext) DECIMAL() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserDECIMAL, 0)
}

func (s *ImmediateContext) HEXADECIMAL() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserHEXADECIMAL, 0)
}

func (s *ImmediateContext) BINARY() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserBINARY, 0)
}

func (s *ImmediateContext) TERNARY() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserTERNARY, 0)
}

func (s *ImmediateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImmediateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) Immediate() (localctx IImmediateContext) {
	localctx = NewImmediateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, vtx1_grammarParserRULE_immediate)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(111)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<vtx1_grammarParserDECIMAL)|(1<<vtx1_grammarParserHEXADECIMAL)|(1<<vtx1_grammarParserBINARY)|(1<<vtx1_grammarParserTERNARY))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IDirectiveContext is an interface to support dynamic dispatch.
type IDirectiveContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDirectiveContext differentiates from other interfaces.
	IsDirectiveContext()
}

type DirectiveContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDirectiveContext() *DirectiveContext {
	var p = new(DirectiveContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_directive
	return p
}

func (*DirectiveContext) IsDirectiveContext() {}

func NewDirectiveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DirectiveContext {
	var p = new(DirectiveContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_directive

	return p
}

func (s *DirectiveContext) GetParser() antlr.Parser { return s.parser }

func (s *DirectiveContext) ORG_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserORG_DIRECTIVE, 0)
}

func (s *DirectiveContext) Immediate() IImmediateContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImmediateContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImmediateContext)
}

func (s *DirectiveContext) DATA_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserDATA_DIRECTIVE, 0)
}

func (s *DirectiveContext) DataList() IDataListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDataListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDataListContext)
}

func (s *DirectiveContext) EQU_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserEQU_DIRECTIVE, 0)
}

func (s *DirectiveContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserIDENTIFIER, 0)
}

func (s *DirectiveContext) COMMA() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCOMMA, 0)
}

func (s *DirectiveContext) INCLUDE_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserINCLUDE_DIRECTIVE, 0)
}

func (s *DirectiveContext) STRING() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserSTRING, 0)
}

func (s *DirectiveContext) SECTION_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserSECTION_DIRECTIVE, 0)
}

func (s *DirectiveContext) ALIGN_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserALIGN_DIRECTIVE, 0)
}

func (s *DirectiveContext) SPACE_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserSPACE_DIRECTIVE, 0)
}

func (s *DirectiveContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DirectiveContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) Directive() (localctx IDirectiveContext) {
	localctx = NewDirectiveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, vtx1_grammarParserRULE_directive)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(129)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case vtx1_grammarParserORG_DIRECTIVE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(113)
			p.Match(vtx1_grammarParserORG_DIRECTIVE)
		}
		{
			p.SetState(114)
			p.Immediate()
		}

	case vtx1_grammarParserDATA_DIRECTIVE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(115)
			p.Match(vtx1_grammarParserDATA_DIRECTIVE)
		}
		{
			p.SetState(116)
			p.DataList()
		}

	case vtx1_grammarParserEQU_DIRECTIVE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(117)
			p.Match(vtx1_grammarParserEQU_DIRECTIVE)
		}
		{
			p.SetState(118)
			p.Match(vtx1_grammarParserIDENTIFIER)
		}
		{
			p.SetState(119)
			p.Match(vtx1_grammarParserCOMMA)
		}
		{
			p.SetState(120)
			p.Immediate()
		}

	case vtx1_grammarParserINCLUDE_DIRECTIVE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(121)
			p.Match(vtx1_grammarParserINCLUDE_DIRECTIVE)
		}
		{
			p.SetState(122)
			p.Match(vtx1_grammarParserSTRING)
		}

	case vtx1_grammarParserSECTION_DIRECTIVE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(123)
			p.Match(vtx1_grammarParserSECTION_DIRECTIVE)
		}
		{
			p.SetState(124)
			p.Match(vtx1_grammarParserIDENTIFIER)
		}

	case vtx1_grammarParserALIGN_DIRECTIVE:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(125)
			p.Match(vtx1_grammarParserALIGN_DIRECTIVE)
		}
		{
			p.SetState(126)
			p.Immediate()
		}

	case vtx1_grammarParserSPACE_DIRECTIVE:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(127)
			p.Match(vtx1_grammarParserSPACE_DIRECTIVE)
		}
		{
			p.SetState(128)
			p.Immediate()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IDataListContext is an interface to support dynamic dispatch.
type IDataListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDataListContext differentiates from other interfaces.
	IsDataListContext()
}

type DataListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDataListContext() *DataListContext {
	var p = new(DataListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_dataList
	return p
}

func (*DataListContext) IsDataListContext() {}

func NewDataListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DataListContext {
	var p = new(DataListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_dataList

	return p
}

func (s *DataListContext) GetParser() antlr.Parser { return s.parser }

func (s *DataListContext) AllImmediate() []IImmediateContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IImmediateContext)(nil)).Elem())
	var tst = make([]IImmediateContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IImmediateContext)
		}
	}

	return tst
}

func (s *DataListContext) Immediate(i int) IImmediateContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImmediateContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IImmediateContext)
}

func (s *DataListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(vtx1_grammarParserCOMMA)
}

func (s *DataListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCOMMA, i)
}

func (s *DataListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DataListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) DataList() (localctx IDataListContext) {
	localctx = NewDataListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, vtx1_grammarParserRULE_dataList)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(131)
		p.Immediate()
	}
	p.SetState(136)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == vtx1_grammarParserCOMMA {
		{
			p.SetState(132)
			p.Match(vtx1_grammarParserCOMMA)
		}
		{
			p.SetState(133)
			p.Immediate()
		}

		p.SetState(138)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}
