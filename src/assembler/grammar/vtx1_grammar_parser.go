// Code generated from grammar/vtx1_grammar.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // vtx1_grammar

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type vtx1_grammarParser struct {
	*antlr.BaseParser
}

var Vtx1_grammarParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func vtx1_grammarParserInit() {
	staticData := &Vtx1_grammarParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "'TB'",
		"", "", "", "", "", "", "", "", "':'", "','", "'+'", "'['", "']'", "'.ORG'",
		"", "'.EQU'", "'.INCLUDE'", "'.SECTION'", "'.ALIGN'", "'.SPACE'",
	}
	staticData.SymbolicNames = []string{
		"", "WHITESPACE", "COMMENT", "EOL", "ALU_OP", "MEM_OP", "CTRL_OP", "VEC_OP",
		"FP_OP", "SYS_OP", "COMPLEX_OP", "COMPLEX_VEC", "COMPLEX_MEM", "COMPLEX_SYS",
		"GPR", "TB_REG", "SPECIAL_REG", "VECTOR_REG", "FP_REG", "DECIMAL", "HEXADECIMAL",
		"BINARY", "TERNARY", "STRING", "COLON", "COMMA", "PLUS", "LSQUARE",
		"RSQUARE", "ORG_DIRECTIVE", "DATA_DIRECTIVE", "EQU_DIRECTIVE", "INCLUDE_DIRECTIVE",
		"SECTION_DIRECTIVE", "ALIGN_DIRECTIVE", "SPACE_DIRECTIVE", "IDENTIFIER",
	}
	staticData.RuleNames = []string{
		"program", "line", "label", "comment", "instruction", "vliwInstruction",
		"mnemonic", "operand", "register", "memoryOperand", "baseRegister",
		"indexRegister", "offsetImmediate", "immediate", "directive", "dataList",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 36, 140, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		1, 0, 5, 0, 34, 8, 0, 10, 0, 12, 0, 37, 9, 0, 1, 0, 1, 0, 1, 1, 3, 1, 42,
		8, 1, 1, 1, 1, 1, 1, 1, 3, 1, 47, 8, 1, 1, 1, 3, 1, 50, 8, 1, 1, 1, 1,
		1, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 5, 4, 63, 8, 4,
		10, 4, 12, 4, 66, 9, 4, 3, 4, 68, 8, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1,
		5, 1, 5, 3, 5, 77, 8, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 83, 8, 5, 1, 6,
		1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 3, 7, 91, 8, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1,
		9, 1, 9, 1, 9, 3, 9, 100, 8, 9, 3, 9, 102, 8, 9, 1, 9, 1, 9, 1, 10, 1,
		10, 1, 11, 1, 11, 1, 12, 1, 12, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 1, 14,
		1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1,
		14, 1, 14, 3, 14, 130, 8, 14, 1, 15, 1, 15, 1, 15, 5, 15, 135, 8, 15, 10,
		15, 12, 15, 138, 9, 15, 1, 15, 0, 0, 16, 0, 2, 4, 6, 8, 10, 12, 14, 16,
		18, 20, 22, 24, 26, 28, 30, 0, 4, 1, 0, 4, 13, 2, 0, 14, 14, 16, 18, 1,
		0, 14, 15, 1, 0, 19, 22, 145, 0, 35, 1, 0, 0, 0, 2, 41, 1, 0, 0, 0, 4,
		53, 1, 0, 0, 0, 6, 56, 1, 0, 0, 0, 8, 58, 1, 0, 0, 0, 10, 69, 1, 0, 0,
		0, 12, 84, 1, 0, 0, 0, 14, 90, 1, 0, 0, 0, 16, 92, 1, 0, 0, 0, 18, 94,
		1, 0, 0, 0, 20, 105, 1, 0, 0, 0, 22, 107, 1, 0, 0, 0, 24, 109, 1, 0, 0,
		0, 26, 111, 1, 0, 0, 0, 28, 129, 1, 0, 0, 0, 30, 131, 1, 0, 0, 0, 32, 34,
		3, 2, 1, 0, 33, 32, 1, 0, 0, 0, 34, 37, 1, 0, 0, 0, 35, 33, 1, 0, 0, 0,
		35, 36, 1, 0, 0, 0, 36, 38, 1, 0, 0, 0, 37, 35, 1, 0, 0, 0, 38, 39, 5,
		0, 0, 1, 39, 1, 1, 0, 0, 0, 40, 42, 3, 4, 2, 0, 41, 40, 1, 0, 0, 0, 41,
		42, 1, 0, 0, 0, 42, 46, 1, 0, 0, 0, 43, 47, 3, 8, 4, 0, 44, 47, 3, 28,
		14, 0, 45, 47, 3, 10, 5, 0, 46, 43, 1, 0, 0, 0, 46, 44, 1, 0, 0, 0, 46,
		45, 1, 0, 0, 0, 46, 47, 1, 0, 0, 0, 47, 49, 1, 0, 0, 0, 48, 50, 3, 6, 3,
		0, 49, 48, 1, 0, 0, 0, 49, 50, 1, 0, 0, 0, 50, 51, 1, 0, 0, 0, 51, 52,
		5, 3, 0, 0, 52, 3, 1, 0, 0, 0, 53, 54, 5, 36, 0, 0, 54, 55, 5, 24, 0, 0,
		55, 5, 1, 0, 0, 0, 56, 57, 5, 2, 0, 0, 57, 7, 1, 0, 0, 0, 58, 67, 3, 12,
		6, 0, 59, 64, 3, 14, 7, 0, 60, 61, 5, 25, 0, 0, 61, 63, 3, 14, 7, 0, 62,
		60, 1, 0, 0, 0, 63, 66, 1, 0, 0, 0, 64, 62, 1, 0, 0, 0, 64, 65, 1, 0, 0,
		0, 65, 68, 1, 0, 0, 0, 66, 64, 1, 0, 0, 0, 67, 59, 1, 0, 0, 0, 67, 68,
		1, 0, 0, 0, 68, 9, 1, 0, 0, 0, 69, 70, 5, 27, 0, 0, 70, 71, 3, 8, 4, 0,
		71, 76, 5, 28, 0, 0, 72, 73, 5, 27, 0, 0, 73, 74, 3, 8, 4, 0, 74, 75, 5,
		28, 0, 0, 75, 77, 1, 0, 0, 0, 76, 72, 1, 0, 0, 0, 76, 77, 1, 0, 0, 0, 77,
		82, 1, 0, 0, 0, 78, 79, 5, 27, 0, 0, 79, 80, 3, 8, 4, 0, 80, 81, 5, 28,
		0, 0, 81, 83, 1, 0, 0, 0, 82, 78, 1, 0, 0, 0, 82, 83, 1, 0, 0, 0, 83, 11,
		1, 0, 0, 0, 84, 85, 7, 0, 0, 0, 85, 13, 1, 0, 0, 0, 86, 91, 3, 16, 8, 0,
		87, 91, 3, 26, 13, 0, 88, 91, 3, 18, 9, 0, 89, 91, 5, 36, 0, 0, 90, 86,
		1, 0, 0, 0, 90, 87, 1, 0, 0, 0, 90, 88, 1, 0, 0, 0, 90, 89, 1, 0, 0, 0,
		91, 15, 1, 0, 0, 0, 92, 93, 7, 1, 0, 0, 93, 17, 1, 0, 0, 0, 94, 95, 5,
		27, 0, 0, 95, 101, 3, 20, 10, 0, 96, 99, 5, 26, 0, 0, 97, 100, 3, 22, 11,
		0, 98, 100, 3, 24, 12, 0, 99, 97, 1, 0, 0, 0, 99, 98, 1, 0, 0, 0, 100,
		102, 1, 0, 0, 0, 101, 96, 1, 0, 0, 0, 101, 102, 1, 0, 0, 0, 102, 103, 1,
		0, 0, 0, 103, 104, 5, 28, 0, 0, 104, 19, 1, 0, 0, 0, 105, 106, 7, 2, 0,
		0, 106, 21, 1, 0, 0, 0, 107, 108, 5, 14, 0, 0, 108, 23, 1, 0, 0, 0, 109,
		110, 3, 26, 13, 0, 110, 25, 1, 0, 0, 0, 111, 112, 7, 3, 0, 0, 112, 27,
		1, 0, 0, 0, 113, 114, 5, 29, 0, 0, 114, 130, 3, 26, 13, 0, 115, 116, 5,
		30, 0, 0, 116, 130, 3, 30, 15, 0, 117, 118, 5, 31, 0, 0, 118, 119, 5, 36,
		0, 0, 119, 120, 5, 25, 0, 0, 120, 130, 3, 26, 13, 0, 121, 122, 5, 32, 0,
		0, 122, 130, 5, 23, 0, 0, 123, 124, 5, 33, 0, 0, 124, 130, 5, 36, 0, 0,
		125, 126, 5, 34, 0, 0, 126, 130, 3, 26, 13, 0, 127, 128, 5, 35, 0, 0, 128,
		130, 3, 26, 13, 0, 129, 113, 1, 0, 0, 0, 129, 115, 1, 0, 0, 0, 129, 117,
		1, 0, 0, 0, 129, 121, 1, 0, 0, 0, 129, 123, 1, 0, 0, 0, 129, 125, 1, 0,
		0, 0, 129, 127, 1, 0, 0, 0, 130, 29, 1, 0, 0, 0, 131, 136, 3, 26, 13, 0,
		132, 133, 5, 25, 0, 0, 133, 135, 3, 26, 13, 0, 134, 132, 1, 0, 0, 0, 135,
		138, 1, 0, 0, 0, 136, 134, 1, 0, 0, 0, 136, 137, 1, 0, 0, 0, 137, 31, 1,
		0, 0, 0, 138, 136, 1, 0, 0, 0, 13, 35, 41, 46, 49, 64, 67, 76, 82, 90,
		99, 101, 129, 136,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// vtx1_grammarParserInit initializes any static state used to implement vtx1_grammarParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// Newvtx1_grammarParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func Vtx1_grammarParserInit() {
	staticData := &Vtx1_grammarParserStaticData
	staticData.once.Do(vtx1_grammarParserInit)
}

// Newvtx1_grammarParser produces a new parser instance for the optional input antlr.TokenStream.
func Newvtx1_grammarParser(input antlr.TokenStream) *vtx1_grammarParser {
	Vtx1_grammarParserInit()
	this := new(vtx1_grammarParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &Vtx1_grammarParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
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

	// Getter signatures
	EOF() antlr.TerminalNode
	AllLine() []ILineContext
	Line(i int) ILineContext

	// IsProgramContext differentiates from other interfaces.
	IsProgramContext()
}

type ProgramContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgramContext() *ProgramContext {
	var p = new(ProgramContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_program
	return p
}

func InitEmptyProgramContext(p *ProgramContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_program
}

func (*ProgramContext) IsProgramContext() {}

func NewProgramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramContext {
	var p = new(ProgramContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_program

	return p
}

func (s *ProgramContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramContext) EOF() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserEOF, 0)
}

func (s *ProgramContext) AllLine() []ILineContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILineContext); ok {
			len++
		}
	}

	tst := make([]ILineContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILineContext); ok {
			tst[i] = t.(ILineContext)
			i++
		}
	}

	return tst
}

func (s *ProgramContext) Line(i int) ILineContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILineContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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

func (s *ProgramContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterProgram(s)
	}
}

func (s *ProgramContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitProgram(s)
	}
}

func (s *ProgramContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitProgram(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) Program() (localctx IProgramContext) {
	localctx = NewProgramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, vtx1_grammarParserRULE_program)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(35)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&137036316668) != 0 {
		{
			p.SetState(32)
			p.Line()
		}

		p.SetState(37)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(38)
		p.Match(vtx1_grammarParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILineContext is an interface to support dynamic dispatch.
type ILineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOL() antlr.TerminalNode
	Label() ILabelContext
	Instruction() IInstructionContext
	Directive() IDirectiveContext
	VliwInstruction() IVliwInstructionContext
	Comment() ICommentContext

	// IsLineContext differentiates from other interfaces.
	IsLineContext()
}

type LineContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLineContext() *LineContext {
	var p = new(LineContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_line
	return p
}

func InitEmptyLineContext(p *LineContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_line
}

func (*LineContext) IsLineContext() {}

func NewLineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LineContext {
	var p = new(LineContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_line

	return p
}

func (s *LineContext) GetParser() antlr.Parser { return s.parser }

func (s *LineContext) EOL() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserEOL, 0)
}

func (s *LineContext) Label() ILabelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILabelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILabelContext)
}

func (s *LineContext) Instruction() IInstructionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInstructionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInstructionContext)
}

func (s *LineContext) Directive() IDirectiveContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDirectiveContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDirectiveContext)
}

func (s *LineContext) VliwInstruction() IVliwInstructionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVliwInstructionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVliwInstructionContext)
}

func (s *LineContext) Comment() ICommentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICommentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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

func (s *LineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterLine(s)
	}
}

func (s *LineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitLine(s)
	}
}

func (s *LineContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitLine(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) Line() (localctx ILineContext) {
	localctx = NewLineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, vtx1_grammarParserRULE_line)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(41)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserIDENTIFIER {
		{
			p.SetState(40)
			p.Label()
		}

	}
	p.SetState(46)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
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
	if p.HasError() {
		goto errorExit
	}
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
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILabelContext is an interface to support dynamic dispatch.
type ILabelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode

	// IsLabelContext differentiates from other interfaces.
	IsLabelContext()
}

type LabelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLabelContext() *LabelContext {
	var p = new(LabelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_label
	return p
}

func InitEmptyLabelContext(p *LabelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_label
}

func (*LabelContext) IsLabelContext() {}

func NewLabelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LabelContext {
	var p = new(LabelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *LabelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterLabel(s)
	}
}

func (s *LabelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitLabel(s)
	}
}

func (s *LabelContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitLabel(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) Label() (localctx ILabelContext) {
	localctx = NewLabelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, vtx1_grammarParserRULE_label)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(53)
		p.Match(vtx1_grammarParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(54)
		p.Match(vtx1_grammarParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICommentContext is an interface to support dynamic dispatch.
type ICommentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COMMENT() antlr.TerminalNode

	// IsCommentContext differentiates from other interfaces.
	IsCommentContext()
}

type CommentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommentContext() *CommentContext {
	var p = new(CommentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_comment
	return p
}

func InitEmptyCommentContext(p *CommentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_comment
}

func (*CommentContext) IsCommentContext() {}

func NewCommentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommentContext {
	var p = new(CommentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *CommentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterComment(s)
	}
}

func (s *CommentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitComment(s)
	}
}

func (s *CommentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitComment(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) Comment() (localctx ICommentContext) {
	localctx = NewCommentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, vtx1_grammarParserRULE_comment)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(56)
		p.Match(vtx1_grammarParserCOMMENT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IInstructionContext is an interface to support dynamic dispatch.
type IInstructionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Mnemonic() IMnemonicContext
	AllOperand() []IOperandContext
	Operand(i int) IOperandContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsInstructionContext differentiates from other interfaces.
	IsInstructionContext()
}

type InstructionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInstructionContext() *InstructionContext {
	var p = new(InstructionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_instruction
	return p
}

func InitEmptyInstructionContext(p *InstructionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_instruction
}

func (*InstructionContext) IsInstructionContext() {}

func NewInstructionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InstructionContext {
	var p = new(InstructionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_instruction

	return p
}

func (s *InstructionContext) GetParser() antlr.Parser { return s.parser }

func (s *InstructionContext) Mnemonic() IMnemonicContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMnemonicContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMnemonicContext)
}

func (s *InstructionContext) AllOperand() []IOperandContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOperandContext); ok {
			len++
		}
	}

	tst := make([]IOperandContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOperandContext); ok {
			tst[i] = t.(IOperandContext)
			i++
		}
	}

	return tst
}

func (s *InstructionContext) Operand(i int) IOperandContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperandContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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

func (s *InstructionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterInstruction(s)
	}
}

func (s *InstructionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitInstruction(s)
	}
}

func (s *InstructionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitInstruction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) Instruction() (localctx IInstructionContext) {
	localctx = NewInstructionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, vtx1_grammarParserRULE_instruction)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(58)
		p.Mnemonic()
	}
	p.SetState(67)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68862033920) != 0 {
		{
			p.SetState(59)
			p.Operand()
		}
		p.SetState(64)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == vtx1_grammarParserCOMMA {
			{
				p.SetState(60)
				p.Match(vtx1_grammarParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(61)
				p.Operand()
			}

			p.SetState(66)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IVliwInstructionContext is an interface to support dynamic dispatch.
type IVliwInstructionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllLSQUARE() []antlr.TerminalNode
	LSQUARE(i int) antlr.TerminalNode
	AllInstruction() []IInstructionContext
	Instruction(i int) IInstructionContext
	AllRSQUARE() []antlr.TerminalNode
	RSQUARE(i int) antlr.TerminalNode

	// IsVliwInstructionContext differentiates from other interfaces.
	IsVliwInstructionContext()
}

type VliwInstructionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVliwInstructionContext() *VliwInstructionContext {
	var p = new(VliwInstructionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_vliwInstruction
	return p
}

func InitEmptyVliwInstructionContext(p *VliwInstructionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_vliwInstruction
}

func (*VliwInstructionContext) IsVliwInstructionContext() {}

func NewVliwInstructionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VliwInstructionContext {
	var p = new(VliwInstructionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IInstructionContext); ok {
			len++
		}
	}

	tst := make([]IInstructionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IInstructionContext); ok {
			tst[i] = t.(IInstructionContext)
			i++
		}
	}

	return tst
}

func (s *VliwInstructionContext) Instruction(i int) IInstructionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInstructionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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

func (s *VliwInstructionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterVliwInstruction(s)
	}
}

func (s *VliwInstructionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitVliwInstruction(s)
	}
}

func (s *VliwInstructionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitVliwInstruction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) VliwInstruction() (localctx IVliwInstructionContext) {
	localctx = NewVliwInstructionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, vtx1_grammarParserRULE_vliwInstruction)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(69)
		p.Match(vtx1_grammarParserLSQUARE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(70)
		p.Instruction()
	}
	{
		p.SetState(71)
		p.Match(vtx1_grammarParserRSQUARE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(76)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(72)
			p.Match(vtx1_grammarParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(73)
			p.Instruction()
		}
		{
			p.SetState(74)
			p.Match(vtx1_grammarParserRSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(82)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserLSQUARE {
		{
			p.SetState(78)
			p.Match(vtx1_grammarParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(79)
			p.Instruction()
		}
		{
			p.SetState(80)
			p.Match(vtx1_grammarParserRSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMnemonicContext is an interface to support dynamic dispatch.
type IMnemonicContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ALU_OP() antlr.TerminalNode
	MEM_OP() antlr.TerminalNode
	CTRL_OP() antlr.TerminalNode
	VEC_OP() antlr.TerminalNode
	FP_OP() antlr.TerminalNode
	SYS_OP() antlr.TerminalNode
	COMPLEX_OP() antlr.TerminalNode
	COMPLEX_VEC() antlr.TerminalNode
	COMPLEX_MEM() antlr.TerminalNode
	COMPLEX_SYS() antlr.TerminalNode

	// IsMnemonicContext differentiates from other interfaces.
	IsMnemonicContext()
}

type MnemonicContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMnemonicContext() *MnemonicContext {
	var p = new(MnemonicContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_mnemonic
	return p
}

func InitEmptyMnemonicContext(p *MnemonicContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_mnemonic
}

func (*MnemonicContext) IsMnemonicContext() {}

func NewMnemonicContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MnemonicContext {
	var p = new(MnemonicContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *MnemonicContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterMnemonic(s)
	}
}

func (s *MnemonicContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitMnemonic(s)
	}
}

func (s *MnemonicContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitMnemonic(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) Mnemonic() (localctx IMnemonicContext) {
	localctx = NewMnemonicContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, vtx1_grammarParserRULE_mnemonic)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(84)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&16368) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOperandContext is an interface to support dynamic dispatch.
type IOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Register() IRegisterContext
	Immediate() IImmediateContext
	MemoryOperand() IMemoryOperandContext
	IDENTIFIER() antlr.TerminalNode

	// IsOperandContext differentiates from other interfaces.
	IsOperandContext()
}

type OperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperandContext() *OperandContext {
	var p = new(OperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_operand
	return p
}

func InitEmptyOperandContext(p *OperandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_operand
}

func (*OperandContext) IsOperandContext() {}

func NewOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperandContext {
	var p = new(OperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_operand

	return p
}

func (s *OperandContext) GetParser() antlr.Parser { return s.parser }

func (s *OperandContext) Register() IRegisterContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRegisterContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRegisterContext)
}

func (s *OperandContext) Immediate() IImmediateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImmediateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImmediateContext)
}

func (s *OperandContext) MemoryOperand() IMemoryOperandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemoryOperandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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

func (s *OperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterOperand(s)
	}
}

func (s *OperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitOperand(s)
	}
}

func (s *OperandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitOperand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) Operand() (localctx IOperandContext) {
	localctx = NewOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, vtx1_grammarParserRULE_operand)
	p.SetState(90)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRegisterContext is an interface to support dynamic dispatch.
type IRegisterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	GPR() antlr.TerminalNode
	SPECIAL_REG() antlr.TerminalNode
	VECTOR_REG() antlr.TerminalNode
	FP_REG() antlr.TerminalNode

	// IsRegisterContext differentiates from other interfaces.
	IsRegisterContext()
}

type RegisterContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRegisterContext() *RegisterContext {
	var p = new(RegisterContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_register
	return p
}

func InitEmptyRegisterContext(p *RegisterContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_register
}

func (*RegisterContext) IsRegisterContext() {}

func NewRegisterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RegisterContext {
	var p = new(RegisterContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *RegisterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterRegister(s)
	}
}

func (s *RegisterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitRegister(s)
	}
}

func (s *RegisterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitRegister(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) Register() (localctx IRegisterContext) {
	localctx = NewRegisterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, vtx1_grammarParserRULE_register)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(92)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&475136) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMemoryOperandContext is an interface to support dynamic dispatch.
type IMemoryOperandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LSQUARE() antlr.TerminalNode
	BaseRegister() IBaseRegisterContext
	RSQUARE() antlr.TerminalNode
	PLUS() antlr.TerminalNode
	IndexRegister() IIndexRegisterContext
	OffsetImmediate() IOffsetImmediateContext

	// IsMemoryOperandContext differentiates from other interfaces.
	IsMemoryOperandContext()
}

type MemoryOperandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemoryOperandContext() *MemoryOperandContext {
	var p = new(MemoryOperandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_memoryOperand
	return p
}

func InitEmptyMemoryOperandContext(p *MemoryOperandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_memoryOperand
}

func (*MemoryOperandContext) IsMemoryOperandContext() {}

func NewMemoryOperandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemoryOperandContext {
	var p = new(MemoryOperandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_memoryOperand

	return p
}

func (s *MemoryOperandContext) GetParser() antlr.Parser { return s.parser }

func (s *MemoryOperandContext) LSQUARE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserLSQUARE, 0)
}

func (s *MemoryOperandContext) BaseRegister() IBaseRegisterContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBaseRegisterContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIndexRegisterContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIndexRegisterContext)
}

func (s *MemoryOperandContext) OffsetImmediate() IOffsetImmediateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOffsetImmediateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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

func (s *MemoryOperandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterMemoryOperand(s)
	}
}

func (s *MemoryOperandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitMemoryOperand(s)
	}
}

func (s *MemoryOperandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitMemoryOperand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) MemoryOperand() (localctx IMemoryOperandContext) {
	localctx = NewMemoryOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, vtx1_grammarParserRULE_memoryOperand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(94)
		p.Match(vtx1_grammarParserLSQUARE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(95)
		p.BaseRegister()
	}
	p.SetState(101)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserPLUS {
		{
			p.SetState(96)
			p.Match(vtx1_grammarParserPLUS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(99)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

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
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

	}
	{
		p.SetState(103)
		p.Match(vtx1_grammarParserRSQUARE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBaseRegisterContext is an interface to support dynamic dispatch.
type IBaseRegisterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TB_REG() antlr.TerminalNode
	GPR() antlr.TerminalNode

	// IsBaseRegisterContext differentiates from other interfaces.
	IsBaseRegisterContext()
}

type BaseRegisterContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBaseRegisterContext() *BaseRegisterContext {
	var p = new(BaseRegisterContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_baseRegister
	return p
}

func InitEmptyBaseRegisterContext(p *BaseRegisterContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_baseRegister
}

func (*BaseRegisterContext) IsBaseRegisterContext() {}

func NewBaseRegisterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BaseRegisterContext {
	var p = new(BaseRegisterContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *BaseRegisterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterBaseRegister(s)
	}
}

func (s *BaseRegisterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitBaseRegister(s)
	}
}

func (s *BaseRegisterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitBaseRegister(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) BaseRegister() (localctx IBaseRegisterContext) {
	localctx = NewBaseRegisterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, vtx1_grammarParserRULE_baseRegister)
	var _la int

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

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIndexRegisterContext is an interface to support dynamic dispatch.
type IIndexRegisterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	GPR() antlr.TerminalNode

	// IsIndexRegisterContext differentiates from other interfaces.
	IsIndexRegisterContext()
}

type IndexRegisterContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIndexRegisterContext() *IndexRegisterContext {
	var p = new(IndexRegisterContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_indexRegister
	return p
}

func InitEmptyIndexRegisterContext(p *IndexRegisterContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_indexRegister
}

func (*IndexRegisterContext) IsIndexRegisterContext() {}

func NewIndexRegisterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IndexRegisterContext {
	var p = new(IndexRegisterContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *IndexRegisterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterIndexRegister(s)
	}
}

func (s *IndexRegisterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitIndexRegister(s)
	}
}

func (s *IndexRegisterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitIndexRegister(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) IndexRegister() (localctx IIndexRegisterContext) {
	localctx = NewIndexRegisterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, vtx1_grammarParserRULE_indexRegister)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(107)
		p.Match(vtx1_grammarParserGPR)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOffsetImmediateContext is an interface to support dynamic dispatch.
type IOffsetImmediateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Immediate() IImmediateContext

	// IsOffsetImmediateContext differentiates from other interfaces.
	IsOffsetImmediateContext()
}

type OffsetImmediateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOffsetImmediateContext() *OffsetImmediateContext {
	var p = new(OffsetImmediateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_offsetImmediate
	return p
}

func InitEmptyOffsetImmediateContext(p *OffsetImmediateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_offsetImmediate
}

func (*OffsetImmediateContext) IsOffsetImmediateContext() {}

func NewOffsetImmediateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OffsetImmediateContext {
	var p = new(OffsetImmediateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_offsetImmediate

	return p
}

func (s *OffsetImmediateContext) GetParser() antlr.Parser { return s.parser }

func (s *OffsetImmediateContext) Immediate() IImmediateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImmediateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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

func (s *OffsetImmediateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterOffsetImmediate(s)
	}
}

func (s *OffsetImmediateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitOffsetImmediate(s)
	}
}

func (s *OffsetImmediateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitOffsetImmediate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) OffsetImmediate() (localctx IOffsetImmediateContext) {
	localctx = NewOffsetImmediateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, vtx1_grammarParserRULE_offsetImmediate)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(109)
		p.Immediate()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IImmediateContext is an interface to support dynamic dispatch.
type IImmediateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DECIMAL() antlr.TerminalNode
	HEXADECIMAL() antlr.TerminalNode
	BINARY() antlr.TerminalNode
	TERNARY() antlr.TerminalNode

	// IsImmediateContext differentiates from other interfaces.
	IsImmediateContext()
}

type ImmediateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImmediateContext() *ImmediateContext {
	var p = new(ImmediateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_immediate
	return p
}

func InitEmptyImmediateContext(p *ImmediateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_immediate
}

func (*ImmediateContext) IsImmediateContext() {}

func NewImmediateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImmediateContext {
	var p = new(ImmediateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *ImmediateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterImmediate(s)
	}
}

func (s *ImmediateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitImmediate(s)
	}
}

func (s *ImmediateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitImmediate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) Immediate() (localctx IImmediateContext) {
	localctx = NewImmediateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, vtx1_grammarParserRULE_immediate)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(111)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&7864320) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDirectiveContext is an interface to support dynamic dispatch.
type IDirectiveContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ORG_DIRECTIVE() antlr.TerminalNode
	Immediate() IImmediateContext
	DATA_DIRECTIVE() antlr.TerminalNode
	DataList() IDataListContext
	EQU_DIRECTIVE() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	COMMA() antlr.TerminalNode
	INCLUDE_DIRECTIVE() antlr.TerminalNode
	STRING() antlr.TerminalNode
	SECTION_DIRECTIVE() antlr.TerminalNode
	ALIGN_DIRECTIVE() antlr.TerminalNode
	SPACE_DIRECTIVE() antlr.TerminalNode

	// IsDirectiveContext differentiates from other interfaces.
	IsDirectiveContext()
}

type DirectiveContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDirectiveContext() *DirectiveContext {
	var p = new(DirectiveContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_directive
	return p
}

func InitEmptyDirectiveContext(p *DirectiveContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_directive
}

func (*DirectiveContext) IsDirectiveContext() {}

func NewDirectiveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DirectiveContext {
	var p = new(DirectiveContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_directive

	return p
}

func (s *DirectiveContext) GetParser() antlr.Parser { return s.parser }

func (s *DirectiveContext) ORG_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserORG_DIRECTIVE, 0)
}

func (s *DirectiveContext) Immediate() IImmediateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImmediateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImmediateContext)
}

func (s *DirectiveContext) DATA_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserDATA_DIRECTIVE, 0)
}

func (s *DirectiveContext) DataList() IDataListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDataListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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

func (s *DirectiveContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterDirective(s)
	}
}

func (s *DirectiveContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitDirective(s)
	}
}

func (s *DirectiveContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitDirective(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) Directive() (localctx IDirectiveContext) {
	localctx = NewDirectiveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, vtx1_grammarParserRULE_directive)
	p.SetState(129)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case vtx1_grammarParserORG_DIRECTIVE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(113)
			p.Match(vtx1_grammarParserORG_DIRECTIVE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(118)
			p.Match(vtx1_grammarParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(119)
			p.Match(vtx1_grammarParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(122)
			p.Match(vtx1_grammarParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case vtx1_grammarParserSECTION_DIRECTIVE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(123)
			p.Match(vtx1_grammarParserSECTION_DIRECTIVE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(124)
			p.Match(vtx1_grammarParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case vtx1_grammarParserALIGN_DIRECTIVE:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(125)
			p.Match(vtx1_grammarParserALIGN_DIRECTIVE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(128)
			p.Immediate()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDataListContext is an interface to support dynamic dispatch.
type IDataListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllImmediate() []IImmediateContext
	Immediate(i int) IImmediateContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsDataListContext differentiates from other interfaces.
	IsDataListContext()
}

type DataListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDataListContext() *DataListContext {
	var p = new(DataListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_dataList
	return p
}

func InitEmptyDataListContext(p *DataListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_dataList
}

func (*DataListContext) IsDataListContext() {}

func NewDataListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DataListContext {
	var p = new(DataListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_dataList

	return p
}

func (s *DataListContext) GetParser() antlr.Parser { return s.parser }

func (s *DataListContext) AllImmediate() []IImmediateContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IImmediateContext); ok {
			len++
		}
	}

	tst := make([]IImmediateContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IImmediateContext); ok {
			tst[i] = t.(IImmediateContext)
			i++
		}
	}

	return tst
}

func (s *DataListContext) Immediate(i int) IImmediateContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImmediateContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

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

func (s *DataListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterDataList(s)
	}
}

func (s *DataListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitDataList(s)
	}
}

func (s *DataListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitDataList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) DataList() (localctx IDataListContext) {
	localctx = NewDataListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, vtx1_grammarParserRULE_dataList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(131)
		p.Immediate()
	}
	p.SetState(136)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == vtx1_grammarParserCOMMA {
		{
			p.SetState(132)
			p.Match(vtx1_grammarParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(133)
			p.Immediate()
		}

		p.SetState(138)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
