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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 41, 204,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 3, 2, 3, 2,
	7, 2, 47, 10, 2, 12, 2, 14, 2, 50, 11, 2, 3, 2, 3, 2, 5, 2, 54, 10, 2,
	5, 2, 56, 10, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	5, 3, 67, 10, 3, 3, 3, 5, 3, 70, 10, 3, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5, 3,
	5, 3, 5, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 3, 7, 7, 7, 85, 10, 7, 12, 7, 14,
	7, 88, 11, 7, 5, 7, 90, 10, 7, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8,
	5, 8, 99, 10, 8, 3, 8, 3, 8, 3, 8, 3, 8, 5, 8, 105, 10, 8, 3, 9, 3, 9,
	3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3,
	10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 5, 10, 127, 10, 10,
	3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 5, 12, 136, 10, 12, 5,
	12, 138, 10, 12, 3, 12, 3, 12, 3, 13, 3, 13, 3, 14, 3, 14, 3, 15, 3, 15,
	3, 16, 3, 16, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3,
	17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 5, 17, 166, 10, 17,
	3, 18, 3, 18, 3, 18, 7, 18, 171, 10, 18, 12, 18, 14, 18, 174, 11, 18, 3,
	19, 3, 19, 5, 19, 178, 10, 19, 3, 20, 3, 20, 3, 20, 7, 20, 183, 10, 20,
	12, 20, 14, 20, 186, 11, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 21,
	6, 21, 194, 10, 21, 13, 21, 14, 21, 195, 3, 22, 3, 22, 5, 22, 200, 10,
	22, 3, 22, 3, 22, 3, 22, 2, 2, 23, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20,
	22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 2, 6, 3, 2, 6, 15, 3, 2, 16,
	20, 3, 2, 16, 17, 3, 2, 21, 24, 2, 217, 2, 48, 3, 2, 2, 2, 4, 66, 3, 2,
	2, 2, 6, 73, 3, 2, 2, 2, 8, 75, 3, 2, 2, 2, 10, 78, 3, 2, 2, 2, 12, 80,
	3, 2, 2, 2, 14, 91, 3, 2, 2, 2, 16, 106, 3, 2, 2, 2, 18, 126, 3, 2, 2,
	2, 20, 128, 3, 2, 2, 2, 22, 130, 3, 2, 2, 2, 24, 141, 3, 2, 2, 2, 26, 143,
	3, 2, 2, 2, 28, 145, 3, 2, 2, 2, 30, 147, 3, 2, 2, 2, 32, 165, 3, 2, 2,
	2, 34, 167, 3, 2, 2, 2, 36, 177, 3, 2, 2, 2, 38, 179, 3, 2, 2, 2, 40, 193,
	3, 2, 2, 2, 42, 197, 3, 2, 2, 2, 44, 47, 5, 4, 3, 2, 45, 47, 5, 6, 4, 2,
	46, 44, 3, 2, 2, 2, 46, 45, 3, 2, 2, 2, 47, 50, 3, 2, 2, 2, 48, 46, 3,
	2, 2, 2, 48, 49, 3, 2, 2, 2, 49, 55, 3, 2, 2, 2, 50, 48, 3, 2, 2, 2, 51,
	53, 7, 38, 2, 2, 52, 54, 7, 5, 2, 2, 53, 52, 3, 2, 2, 2, 53, 54, 3, 2,
	2, 2, 54, 56, 3, 2, 2, 2, 55, 51, 3, 2, 2, 2, 55, 56, 3, 2, 2, 2, 56, 57,
	3, 2, 2, 2, 57, 58, 7, 2, 2, 3, 58, 3, 3, 2, 2, 2, 59, 67, 5, 8, 5, 2,
	60, 67, 5, 12, 7, 2, 61, 67, 5, 32, 17, 2, 62, 67, 5, 42, 22, 2, 63, 67,
	5, 14, 8, 2, 64, 67, 5, 38, 20, 2, 65, 67, 5, 10, 6, 2, 66, 59, 3, 2, 2,
	2, 66, 60, 3, 2, 2, 2, 66, 61, 3, 2, 2, 2, 66, 62, 3, 2, 2, 2, 66, 63,
	3, 2, 2, 2, 66, 64, 3, 2, 2, 2, 66, 65, 3, 2, 2, 2, 67, 69, 3, 2, 2, 2,
	68, 70, 5, 10, 6, 2, 69, 68, 3, 2, 2, 2, 69, 70, 3, 2, 2, 2, 70, 71, 3,
	2, 2, 2, 71, 72, 7, 5, 2, 2, 72, 5, 3, 2, 2, 2, 73, 74, 7, 5, 2, 2, 74,
	7, 3, 2, 2, 2, 75, 76, 7, 39, 2, 2, 76, 77, 7, 26, 2, 2, 77, 9, 3, 2, 2,
	2, 78, 79, 7, 4, 2, 2, 79, 11, 3, 2, 2, 2, 80, 89, 5, 16, 9, 2, 81, 86,
	5, 18, 10, 2, 82, 83, 7, 27, 2, 2, 83, 85, 5, 18, 10, 2, 84, 82, 3, 2,
	2, 2, 85, 88, 3, 2, 2, 2, 86, 84, 3, 2, 2, 2, 86, 87, 3, 2, 2, 2, 87, 90,
	3, 2, 2, 2, 88, 86, 3, 2, 2, 2, 89, 81, 3, 2, 2, 2, 89, 90, 3, 2, 2, 2,
	90, 13, 3, 2, 2, 2, 91, 92, 7, 29, 2, 2, 92, 93, 5, 12, 7, 2, 93, 98, 7,
	30, 2, 2, 94, 95, 7, 29, 2, 2, 95, 96, 5, 12, 7, 2, 96, 97, 7, 30, 2, 2,
	97, 99, 3, 2, 2, 2, 98, 94, 3, 2, 2, 2, 98, 99, 3, 2, 2, 2, 99, 104, 3,
	2, 2, 2, 100, 101, 7, 29, 2, 2, 101, 102, 5, 12, 7, 2, 102, 103, 7, 30,
	2, 2, 103, 105, 3, 2, 2, 2, 104, 100, 3, 2, 2, 2, 104, 105, 3, 2, 2, 2,
	105, 15, 3, 2, 2, 2, 106, 107, 9, 2, 2, 2, 107, 17, 3, 2, 2, 2, 108, 127,
	5, 20, 11, 2, 109, 127, 5, 30, 16, 2, 110, 127, 5, 22, 12, 2, 111, 127,
	7, 39, 2, 2, 112, 113, 7, 39, 2, 2, 113, 114, 7, 28, 2, 2, 114, 127, 5,
	30, 16, 2, 115, 116, 7, 39, 2, 2, 116, 117, 7, 28, 2, 2, 117, 127, 7, 39,
	2, 2, 118, 119, 5, 30, 16, 2, 119, 120, 7, 28, 2, 2, 120, 121, 7, 39, 2,
	2, 121, 127, 3, 2, 2, 2, 122, 123, 5, 30, 16, 2, 123, 124, 7, 28, 2, 2,
	124, 125, 5, 30, 16, 2, 125, 127, 3, 2, 2, 2, 126, 108, 3, 2, 2, 2, 126,
	109, 3, 2, 2, 2, 126, 110, 3, 2, 2, 2, 126, 111, 3, 2, 2, 2, 126, 112,
	3, 2, 2, 2, 126, 115, 3, 2, 2, 2, 126, 118, 3, 2, 2, 2, 126, 122, 3, 2,
	2, 2, 127, 19, 3, 2, 2, 2, 128, 129, 9, 3, 2, 2, 129, 21, 3, 2, 2, 2, 130,
	131, 7, 29, 2, 2, 131, 137, 5, 24, 13, 2, 132, 135, 7, 28, 2, 2, 133, 136,
	5, 26, 14, 2, 134, 136, 5, 28, 15, 2, 135, 133, 3, 2, 2, 2, 135, 134, 3,
	2, 2, 2, 136, 138, 3, 2, 2, 2, 137, 132, 3, 2, 2, 2, 137, 138, 3, 2, 2,
	2, 138, 139, 3, 2, 2, 2, 139, 140, 7, 30, 2, 2, 140, 23, 3, 2, 2, 2, 141,
	142, 9, 4, 2, 2, 142, 25, 3, 2, 2, 2, 143, 144, 7, 17, 2, 2, 144, 27, 3,
	2, 2, 2, 145, 146, 5, 30, 16, 2, 146, 29, 3, 2, 2, 2, 147, 148, 9, 5, 2,
	2, 148, 31, 3, 2, 2, 2, 149, 150, 7, 31, 2, 2, 150, 166, 5, 30, 16, 2,
	151, 152, 7, 32, 2, 2, 152, 166, 5, 34, 18, 2, 153, 154, 7, 33, 2, 2, 154,
	155, 7, 39, 2, 2, 155, 156, 7, 27, 2, 2, 156, 166, 5, 30, 16, 2, 157, 158,
	7, 34, 2, 2, 158, 166, 7, 25, 2, 2, 159, 160, 7, 35, 2, 2, 160, 166, 7,
	39, 2, 2, 161, 162, 7, 36, 2, 2, 162, 166, 5, 30, 16, 2, 163, 164, 7, 37,
	2, 2, 164, 166, 5, 30, 16, 2, 165, 149, 3, 2, 2, 2, 165, 151, 3, 2, 2,
	2, 165, 153, 3, 2, 2, 2, 165, 157, 3, 2, 2, 2, 165, 159, 3, 2, 2, 2, 165,
	161, 3, 2, 2, 2, 165, 163, 3, 2, 2, 2, 166, 33, 3, 2, 2, 2, 167, 172, 5,
	36, 19, 2, 168, 169, 7, 27, 2, 2, 169, 171, 5, 36, 19, 2, 170, 168, 3,
	2, 2, 2, 171, 174, 3, 2, 2, 2, 172, 170, 3, 2, 2, 2, 172, 173, 3, 2, 2,
	2, 173, 35, 3, 2, 2, 2, 174, 172, 3, 2, 2, 2, 175, 178, 5, 30, 16, 2, 176,
	178, 7, 25, 2, 2, 177, 175, 3, 2, 2, 2, 177, 176, 3, 2, 2, 2, 178, 37,
	3, 2, 2, 2, 179, 180, 7, 40, 2, 2, 180, 184, 7, 39, 2, 2, 181, 183, 7,
	39, 2, 2, 182, 181, 3, 2, 2, 2, 183, 186, 3, 2, 2, 2, 184, 182, 3, 2, 2,
	2, 184, 185, 3, 2, 2, 2, 185, 187, 3, 2, 2, 2, 186, 184, 3, 2, 2, 2, 187,
	188, 7, 5, 2, 2, 188, 189, 5, 40, 21, 2, 189, 190, 7, 41, 2, 2, 190, 191,
	7, 5, 2, 2, 191, 39, 3, 2, 2, 2, 192, 194, 5, 4, 3, 2, 193, 192, 3, 2,
	2, 2, 194, 195, 3, 2, 2, 2, 195, 193, 3, 2, 2, 2, 195, 196, 3, 2, 2, 2,
	196, 41, 3, 2, 2, 2, 197, 199, 7, 39, 2, 2, 198, 200, 7, 26, 2, 2, 199,
	198, 3, 2, 2, 2, 199, 200, 3, 2, 2, 2, 200, 201, 3, 2, 2, 2, 201, 202,
	5, 32, 17, 2, 202, 43, 3, 2, 2, 2, 21, 46, 48, 53, 55, 66, 69, 86, 89,
	98, 104, 126, 135, 137, 165, 172, 177, 184, 195, 199,
}
var literalNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "'TB'", "", "",
	"", "", "", "", "", "", "", "':'", "','", "'+'", "'['", "']'", "'.ORG'",
	"", "'.EQU'", "'.INCLUDE'", "'.SECTION'", "'.ALIGN'", "'.SPACE'", "'END'",
	"", "'.MACRO'", "'.ENDM'",
}
var symbolicNames = []string{
	"", "WHITESPACE", "COMMENT", "EOL", "ALU_OP", "MEM_OP", "CTRL_OP", "VEC_OP",
	"FP_OP", "SYS_OP", "COMPLEX_OP", "COMPLEX_VEC", "COMPLEX_MEM", "COMPLEX_SYS",
	"TB_REG", "GPR", "SPECIAL_REG", "VECTOR_REG", "FP_REG", "DECIMAL", "HEXADECIMAL",
	"BINARY", "TERNARY", "STRING", "COLON", "COMMA", "PLUS", "LSQUARE", "RSQUARE",
	"ORG_DIRECTIVE", "DATA_DIRECTIVE", "EQU_DIRECTIVE", "INCLUDE_DIRECTIVE",
	"SECTION_DIRECTIVE", "ALIGN_DIRECTIVE", "SPACE_DIRECTIVE", "END_DIRECTIVE",
	"IDENTIFIER", "MACRO_DIRECTIVE", "ENDM_DIRECTIVE",
}

var ruleNames = []string{
	"program", "line", "blankLine", "label", "comment", "instruction", "vliwInstruction",
	"mnemonic", "operand", "register", "memoryOperand", "baseRegister", "indexRegister",
	"offsetImmediate", "immediate", "directive", "dataList", "dataItem", "macroDefinition",
	"macroBody", "labelledDirective",
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
	vtx1_grammarParserTB_REG            = 14
	vtx1_grammarParserGPR               = 15
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
	vtx1_grammarParserEND_DIRECTIVE     = 36
	vtx1_grammarParserIDENTIFIER        = 37
	vtx1_grammarParserMACRO_DIRECTIVE   = 38
	vtx1_grammarParserENDM_DIRECTIVE    = 39
)

// vtx1_grammarParser rules.
const (
	vtx1_grammarParserRULE_program           = 0
	vtx1_grammarParserRULE_line              = 1
	vtx1_grammarParserRULE_blankLine         = 2
	vtx1_grammarParserRULE_label             = 3
	vtx1_grammarParserRULE_comment           = 4
	vtx1_grammarParserRULE_instruction       = 5
	vtx1_grammarParserRULE_vliwInstruction   = 6
	vtx1_grammarParserRULE_mnemonic          = 7
	vtx1_grammarParserRULE_operand           = 8
	vtx1_grammarParserRULE_register          = 9
	vtx1_grammarParserRULE_memoryOperand     = 10
	vtx1_grammarParserRULE_baseRegister      = 11
	vtx1_grammarParserRULE_indexRegister     = 12
	vtx1_grammarParserRULE_offsetImmediate   = 13
	vtx1_grammarParserRULE_immediate         = 14
	vtx1_grammarParserRULE_directive         = 15
	vtx1_grammarParserRULE_dataList          = 16
	vtx1_grammarParserRULE_dataItem          = 17
	vtx1_grammarParserRULE_macroDefinition   = 18
	vtx1_grammarParserRULE_macroBody         = 19
	vtx1_grammarParserRULE_labelledDirective = 20
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

func (s *ProgramContext) AllBlankLine() []IBlankLineContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IBlankLineContext)(nil)).Elem())
	var tst = make([]IBlankLineContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IBlankLineContext)
		}
	}

	return tst
}

func (s *ProgramContext) BlankLine(i int) IBlankLineContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBlankLineContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IBlankLineContext)
}

func (s *ProgramContext) END_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserEND_DIRECTIVE, 0)
}

func (s *ProgramContext) EOL() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserEOL, 0)
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
	p.SetState(46)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<vtx1_grammarParserCOMMENT)|(1<<vtx1_grammarParserEOL)|(1<<vtx1_grammarParserALU_OP)|(1<<vtx1_grammarParserMEM_OP)|(1<<vtx1_grammarParserCTRL_OP)|(1<<vtx1_grammarParserVEC_OP)|(1<<vtx1_grammarParserFP_OP)|(1<<vtx1_grammarParserSYS_OP)|(1<<vtx1_grammarParserCOMPLEX_OP)|(1<<vtx1_grammarParserCOMPLEX_VEC)|(1<<vtx1_grammarParserCOMPLEX_MEM)|(1<<vtx1_grammarParserCOMPLEX_SYS)|(1<<vtx1_grammarParserLSQUARE)|(1<<vtx1_grammarParserORG_DIRECTIVE)|(1<<vtx1_grammarParserDATA_DIRECTIVE)|(1<<vtx1_grammarParserEQU_DIRECTIVE))) != 0) || (((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(vtx1_grammarParserINCLUDE_DIRECTIVE-32))|(1<<(vtx1_grammarParserSECTION_DIRECTIVE-32))|(1<<(vtx1_grammarParserALIGN_DIRECTIVE-32))|(1<<(vtx1_grammarParserSPACE_DIRECTIVE-32))|(1<<(vtx1_grammarParserIDENTIFIER-32))|(1<<(vtx1_grammarParserMACRO_DIRECTIVE-32)))) != 0) {
		p.SetState(44)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case vtx1_grammarParserCOMMENT, vtx1_grammarParserALU_OP, vtx1_grammarParserMEM_OP, vtx1_grammarParserCTRL_OP, vtx1_grammarParserVEC_OP, vtx1_grammarParserFP_OP, vtx1_grammarParserSYS_OP, vtx1_grammarParserCOMPLEX_OP, vtx1_grammarParserCOMPLEX_VEC, vtx1_grammarParserCOMPLEX_MEM, vtx1_grammarParserCOMPLEX_SYS, vtx1_grammarParserLSQUARE, vtx1_grammarParserORG_DIRECTIVE, vtx1_grammarParserDATA_DIRECTIVE, vtx1_grammarParserEQU_DIRECTIVE, vtx1_grammarParserINCLUDE_DIRECTIVE, vtx1_grammarParserSECTION_DIRECTIVE, vtx1_grammarParserALIGN_DIRECTIVE, vtx1_grammarParserSPACE_DIRECTIVE, vtx1_grammarParserIDENTIFIER, vtx1_grammarParserMACRO_DIRECTIVE:
			{
				p.SetState(42)
				p.Line()
			}

		case vtx1_grammarParserEOL:
			{
				p.SetState(43)
				p.BlankLine()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(48)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(53)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserEND_DIRECTIVE {
		{
			p.SetState(49)
			p.Match(vtx1_grammarParserEND_DIRECTIVE)
		}
		p.SetState(51)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == vtx1_grammarParserEOL {
			{
				p.SetState(50)
				p.Match(vtx1_grammarParserEOL)
			}

		}

	}
	{
		p.SetState(55)
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

func (s *LineContext) LabelledDirective() ILabelledDirectiveContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILabelledDirectiveContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILabelledDirectiveContext)
}

func (s *LineContext) VliwInstruction() IVliwInstructionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVliwInstructionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVliwInstructionContext)
}

func (s *LineContext) MacroDefinition() IMacroDefinitionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMacroDefinitionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMacroDefinitionContext)
}

func (s *LineContext) AllComment() []ICommentContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ICommentContext)(nil)).Elem())
	var tst = make([]ICommentContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ICommentContext)
		}
	}

	return tst
}

func (s *LineContext) Comment(i int) ICommentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICommentContext)(nil)).Elem(), i)

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
	p.SetState(64)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(57)
			p.Label()
		}

	case 2:
		{
			p.SetState(58)
			p.Instruction()
		}

	case 3:
		{
			p.SetState(59)
			p.Directive()
		}

	case 4:
		{
			p.SetState(60)
			p.LabelledDirective()
		}

	case 5:
		{
			p.SetState(61)
			p.VliwInstruction()
		}

	case 6:
		{
			p.SetState(62)
			p.MacroDefinition()
		}

	case 7:
		{
			p.SetState(63)
			p.Comment()
		}

	}
	p.SetState(67)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserCOMMENT {
		{
			p.SetState(66)
			p.Comment()
		}

	}
	{
		p.SetState(69)
		p.Match(vtx1_grammarParserEOL)
	}

	return localctx
}

// IBlankLineContext is an interface to support dynamic dispatch.
type IBlankLineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBlankLineContext differentiates from other interfaces.
	IsBlankLineContext()
}

type BlankLineContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBlankLineContext() *BlankLineContext {
	var p = new(BlankLineContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_blankLine
	return p
}

func (*BlankLineContext) IsBlankLineContext() {}

func NewBlankLineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BlankLineContext {
	var p = new(BlankLineContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_blankLine

	return p
}

func (s *BlankLineContext) GetParser() antlr.Parser { return s.parser }

func (s *BlankLineContext) EOL() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserEOL, 0)
}

func (s *BlankLineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BlankLineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) BlankLine() (localctx IBlankLineContext) {
	localctx = NewBlankLineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, vtx1_grammarParserRULE_blankLine)

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
		p.SetState(71)
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
	p.EnterRule(localctx, 6, vtx1_grammarParserRULE_label)

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
		p.SetState(73)
		p.Match(vtx1_grammarParserIDENTIFIER)
	}
	{
		p.SetState(74)
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
	p.EnterRule(localctx, 8, vtx1_grammarParserRULE_comment)

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
		p.SetState(76)
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
	p.EnterRule(localctx, 10, vtx1_grammarParserRULE_instruction)
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
		p.SetState(78)
		p.Mnemonic()
	}
	p.SetState(87)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if ((_la-14)&-(0x1f+1)) == 0 && ((1<<uint((_la-14)))&((1<<(vtx1_grammarParserTB_REG-14))|(1<<(vtx1_grammarParserGPR-14))|(1<<(vtx1_grammarParserSPECIAL_REG-14))|(1<<(vtx1_grammarParserVECTOR_REG-14))|(1<<(vtx1_grammarParserFP_REG-14))|(1<<(vtx1_grammarParserDECIMAL-14))|(1<<(vtx1_grammarParserHEXADECIMAL-14))|(1<<(vtx1_grammarParserBINARY-14))|(1<<(vtx1_grammarParserTERNARY-14))|(1<<(vtx1_grammarParserLSQUARE-14))|(1<<(vtx1_grammarParserIDENTIFIER-14)))) != 0 {
		{
			p.SetState(79)
			p.Operand()
		}
		p.SetState(84)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == vtx1_grammarParserCOMMA {
			{
				p.SetState(80)
				p.Match(vtx1_grammarParserCOMMA)
			}
			{
				p.SetState(81)
				p.Operand()
			}

			p.SetState(86)
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
	p.EnterRule(localctx, 12, vtx1_grammarParserRULE_vliwInstruction)
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
		p.SetState(89)
		p.Match(vtx1_grammarParserLSQUARE)
	}
	{
		p.SetState(90)
		p.Instruction()
	}
	{
		p.SetState(91)
		p.Match(vtx1_grammarParserRSQUARE)
	}
	p.SetState(96)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(92)
			p.Match(vtx1_grammarParserLSQUARE)
		}
		{
			p.SetState(93)
			p.Instruction()
		}
		{
			p.SetState(94)
			p.Match(vtx1_grammarParserRSQUARE)
		}

	}
	p.SetState(102)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserLSQUARE {
		{
			p.SetState(98)
			p.Match(vtx1_grammarParserLSQUARE)
		}
		{
			p.SetState(99)
			p.Instruction()
		}
		{
			p.SetState(100)
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
	p.EnterRule(localctx, 14, vtx1_grammarParserRULE_mnemonic)
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
		p.SetState(104)
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

func (s *OperandContext) AllImmediate() []IImmediateContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IImmediateContext)(nil)).Elem())
	var tst = make([]IImmediateContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IImmediateContext)
		}
	}

	return tst
}

func (s *OperandContext) Immediate(i int) IImmediateContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImmediateContext)(nil)).Elem(), i)

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

func (s *OperandContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(vtx1_grammarParserIDENTIFIER)
}

func (s *OperandContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserIDENTIFIER, i)
}

func (s *OperandContext) PLUS() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserPLUS, 0)
}

func (s *OperandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) Operand() (localctx IOperandContext) {
	localctx = NewOperandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, vtx1_grammarParserRULE_operand)

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

	p.SetState(124)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(106)
			p.Register()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(107)
			p.Immediate()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(108)
			p.MemoryOperand()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(109)
			p.Match(vtx1_grammarParserIDENTIFIER)
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(110)
			p.Match(vtx1_grammarParserIDENTIFIER)
		}
		{
			p.SetState(111)
			p.Match(vtx1_grammarParserPLUS)
		}
		{
			p.SetState(112)
			p.Immediate()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(113)
			p.Match(vtx1_grammarParserIDENTIFIER)
		}
		{
			p.SetState(114)
			p.Match(vtx1_grammarParserPLUS)
		}
		{
			p.SetState(115)
			p.Match(vtx1_grammarParserIDENTIFIER)
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(116)
			p.Immediate()
		}
		{
			p.SetState(117)
			p.Match(vtx1_grammarParserPLUS)
		}
		{
			p.SetState(118)
			p.Match(vtx1_grammarParserIDENTIFIER)
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(120)
			p.Immediate()
		}
		{
			p.SetState(121)
			p.Match(vtx1_grammarParserPLUS)
		}
		{
			p.SetState(122)
			p.Immediate()
		}

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

func (s *RegisterContext) TB_REG() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserTB_REG, 0)
}

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
	p.EnterRule(localctx, 18, vtx1_grammarParserRULE_register)
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
		p.SetState(126)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<vtx1_grammarParserTB_REG)|(1<<vtx1_grammarParserGPR)|(1<<vtx1_grammarParserSPECIAL_REG)|(1<<vtx1_grammarParserVECTOR_REG)|(1<<vtx1_grammarParserFP_REG))) != 0) {
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
	p.EnterRule(localctx, 20, vtx1_grammarParserRULE_memoryOperand)
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
		p.SetState(128)
		p.Match(vtx1_grammarParserLSQUARE)
	}
	{
		p.SetState(129)
		p.BaseRegister()
	}
	p.SetState(135)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserPLUS {
		{
			p.SetState(130)
			p.Match(vtx1_grammarParserPLUS)
		}
		p.SetState(133)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case vtx1_grammarParserGPR:
			{
				p.SetState(131)
				p.IndexRegister()
			}

		case vtx1_grammarParserDECIMAL, vtx1_grammarParserHEXADECIMAL, vtx1_grammarParserBINARY, vtx1_grammarParserTERNARY:
			{
				p.SetState(132)
				p.OffsetImmediate()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

	}
	{
		p.SetState(137)
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
	p.EnterRule(localctx, 22, vtx1_grammarParserRULE_baseRegister)
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
		p.SetState(139)
		_la = p.GetTokenStream().LA(1)

		if !(_la == vtx1_grammarParserTB_REG || _la == vtx1_grammarParserGPR) {
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
	p.EnterRule(localctx, 24, vtx1_grammarParserRULE_indexRegister)

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
		p.SetState(141)
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
	p.EnterRule(localctx, 26, vtx1_grammarParserRULE_offsetImmediate)

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
		p.SetState(143)
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
	p.EnterRule(localctx, 28, vtx1_grammarParserRULE_immediate)
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
		p.SetState(145)
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
	p.EnterRule(localctx, 30, vtx1_grammarParserRULE_directive)

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

	p.SetState(163)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case vtx1_grammarParserORG_DIRECTIVE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(147)
			p.Match(vtx1_grammarParserORG_DIRECTIVE)
		}
		{
			p.SetState(148)
			p.Immediate()
		}

	case vtx1_grammarParserDATA_DIRECTIVE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(149)
			p.Match(vtx1_grammarParserDATA_DIRECTIVE)
		}
		{
			p.SetState(150)
			p.DataList()
		}

	case vtx1_grammarParserEQU_DIRECTIVE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(151)
			p.Match(vtx1_grammarParserEQU_DIRECTIVE)
		}
		{
			p.SetState(152)
			p.Match(vtx1_grammarParserIDENTIFIER)
		}
		{
			p.SetState(153)
			p.Match(vtx1_grammarParserCOMMA)
		}
		{
			p.SetState(154)
			p.Immediate()
		}

	case vtx1_grammarParserINCLUDE_DIRECTIVE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(155)
			p.Match(vtx1_grammarParserINCLUDE_DIRECTIVE)
		}
		{
			p.SetState(156)
			p.Match(vtx1_grammarParserSTRING)
		}

	case vtx1_grammarParserSECTION_DIRECTIVE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(157)
			p.Match(vtx1_grammarParserSECTION_DIRECTIVE)
		}
		{
			p.SetState(158)
			p.Match(vtx1_grammarParserIDENTIFIER)
		}

	case vtx1_grammarParserALIGN_DIRECTIVE:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(159)
			p.Match(vtx1_grammarParserALIGN_DIRECTIVE)
		}
		{
			p.SetState(160)
			p.Immediate()
		}

	case vtx1_grammarParserSPACE_DIRECTIVE:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(161)
			p.Match(vtx1_grammarParserSPACE_DIRECTIVE)
		}
		{
			p.SetState(162)
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

func (s *DataListContext) AllDataItem() []IDataItemContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IDataItemContext)(nil)).Elem())
	var tst = make([]IDataItemContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IDataItemContext)
		}
	}

	return tst
}

func (s *DataListContext) DataItem(i int) IDataItemContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDataItemContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IDataItemContext)
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
	p.EnterRule(localctx, 32, vtx1_grammarParserRULE_dataList)
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
		p.SetState(165)
		p.DataItem()
	}
	p.SetState(170)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == vtx1_grammarParserCOMMA {
		{
			p.SetState(166)
			p.Match(vtx1_grammarParserCOMMA)
		}
		{
			p.SetState(167)
			p.DataItem()
		}

		p.SetState(172)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IDataItemContext is an interface to support dynamic dispatch.
type IDataItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDataItemContext differentiates from other interfaces.
	IsDataItemContext()
}

type DataItemContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDataItemContext() *DataItemContext {
	var p = new(DataItemContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_dataItem
	return p
}

func (*DataItemContext) IsDataItemContext() {}

func NewDataItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DataItemContext {
	var p = new(DataItemContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_dataItem

	return p
}

func (s *DataItemContext) GetParser() antlr.Parser { return s.parser }

func (s *DataItemContext) Immediate() IImmediateContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImmediateContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImmediateContext)
}

func (s *DataItemContext) STRING() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserSTRING, 0)
}

func (s *DataItemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DataItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) DataItem() (localctx IDataItemContext) {
	localctx = NewDataItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, vtx1_grammarParserRULE_dataItem)

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

	p.SetState(175)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case vtx1_grammarParserDECIMAL, vtx1_grammarParserHEXADECIMAL, vtx1_grammarParserBINARY, vtx1_grammarParserTERNARY:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(173)
			p.Immediate()
		}

	case vtx1_grammarParserSTRING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(174)
			p.Match(vtx1_grammarParserSTRING)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IMacroDefinitionContext is an interface to support dynamic dispatch.
type IMacroDefinitionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMacroDefinitionContext differentiates from other interfaces.
	IsMacroDefinitionContext()
}

type MacroDefinitionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMacroDefinitionContext() *MacroDefinitionContext {
	var p = new(MacroDefinitionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_macroDefinition
	return p
}

func (*MacroDefinitionContext) IsMacroDefinitionContext() {}

func NewMacroDefinitionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MacroDefinitionContext {
	var p = new(MacroDefinitionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_macroDefinition

	return p
}

func (s *MacroDefinitionContext) GetParser() antlr.Parser { return s.parser }

func (s *MacroDefinitionContext) MACRO_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserMACRO_DIRECTIVE, 0)
}

func (s *MacroDefinitionContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(vtx1_grammarParserIDENTIFIER)
}

func (s *MacroDefinitionContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserIDENTIFIER, i)
}

func (s *MacroDefinitionContext) AllEOL() []antlr.TerminalNode {
	return s.GetTokens(vtx1_grammarParserEOL)
}

func (s *MacroDefinitionContext) EOL(i int) antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserEOL, i)
}

func (s *MacroDefinitionContext) MacroBody() IMacroBodyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMacroBodyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMacroBodyContext)
}

func (s *MacroDefinitionContext) ENDM_DIRECTIVE() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserENDM_DIRECTIVE, 0)
}

func (s *MacroDefinitionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MacroDefinitionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) MacroDefinition() (localctx IMacroDefinitionContext) {
	localctx = NewMacroDefinitionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, vtx1_grammarParserRULE_macroDefinition)
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
		p.SetState(177)
		p.Match(vtx1_grammarParserMACRO_DIRECTIVE)
	}
	{
		p.SetState(178)
		p.Match(vtx1_grammarParserIDENTIFIER)
	}
	p.SetState(182)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == vtx1_grammarParserIDENTIFIER {
		{
			p.SetState(179)
			p.Match(vtx1_grammarParserIDENTIFIER)
		}

		p.SetState(184)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(185)
		p.Match(vtx1_grammarParserEOL)
	}
	{
		p.SetState(186)
		p.MacroBody()
	}
	{
		p.SetState(187)
		p.Match(vtx1_grammarParserENDM_DIRECTIVE)
	}
	{
		p.SetState(188)
		p.Match(vtx1_grammarParserEOL)
	}

	return localctx
}

// IMacroBodyContext is an interface to support dynamic dispatch.
type IMacroBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMacroBodyContext differentiates from other interfaces.
	IsMacroBodyContext()
}

type MacroBodyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMacroBodyContext() *MacroBodyContext {
	var p = new(MacroBodyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_macroBody
	return p
}

func (*MacroBodyContext) IsMacroBodyContext() {}

func NewMacroBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MacroBodyContext {
	var p = new(MacroBodyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_macroBody

	return p
}

func (s *MacroBodyContext) GetParser() antlr.Parser { return s.parser }

func (s *MacroBodyContext) AllLine() []ILineContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ILineContext)(nil)).Elem())
	var tst = make([]ILineContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ILineContext)
		}
	}

	return tst
}

func (s *MacroBodyContext) Line(i int) ILineContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILineContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ILineContext)
}

func (s *MacroBodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MacroBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) MacroBody() (localctx IMacroBodyContext) {
	localctx = NewMacroBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, vtx1_grammarParserRULE_macroBody)
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
	p.SetState(191)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<vtx1_grammarParserCOMMENT)|(1<<vtx1_grammarParserALU_OP)|(1<<vtx1_grammarParserMEM_OP)|(1<<vtx1_grammarParserCTRL_OP)|(1<<vtx1_grammarParserVEC_OP)|(1<<vtx1_grammarParserFP_OP)|(1<<vtx1_grammarParserSYS_OP)|(1<<vtx1_grammarParserCOMPLEX_OP)|(1<<vtx1_grammarParserCOMPLEX_VEC)|(1<<vtx1_grammarParserCOMPLEX_MEM)|(1<<vtx1_grammarParserCOMPLEX_SYS)|(1<<vtx1_grammarParserLSQUARE)|(1<<vtx1_grammarParserORG_DIRECTIVE)|(1<<vtx1_grammarParserDATA_DIRECTIVE)|(1<<vtx1_grammarParserEQU_DIRECTIVE))) != 0) || (((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(vtx1_grammarParserINCLUDE_DIRECTIVE-32))|(1<<(vtx1_grammarParserSECTION_DIRECTIVE-32))|(1<<(vtx1_grammarParserALIGN_DIRECTIVE-32))|(1<<(vtx1_grammarParserSPACE_DIRECTIVE-32))|(1<<(vtx1_grammarParserIDENTIFIER-32))|(1<<(vtx1_grammarParserMACRO_DIRECTIVE-32)))) != 0) {
		{
			p.SetState(190)
			p.Line()
		}

		p.SetState(193)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ILabelledDirectiveContext is an interface to support dynamic dispatch.
type ILabelledDirectiveContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLabelledDirectiveContext differentiates from other interfaces.
	IsLabelledDirectiveContext()
}

type LabelledDirectiveContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLabelledDirectiveContext() *LabelledDirectiveContext {
	var p = new(LabelledDirectiveContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_labelledDirective
	return p
}

func (*LabelledDirectiveContext) IsLabelledDirectiveContext() {}

func NewLabelledDirectiveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LabelledDirectiveContext {
	var p = new(LabelledDirectiveContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_labelledDirective

	return p
}

func (s *LabelledDirectiveContext) GetParser() antlr.Parser { return s.parser }

func (s *LabelledDirectiveContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserIDENTIFIER, 0)
}

func (s *LabelledDirectiveContext) Directive() IDirectiveContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDirectiveContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDirectiveContext)
}

func (s *LabelledDirectiveContext) COLON() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCOLON, 0)
}

func (s *LabelledDirectiveContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LabelledDirectiveContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (p *vtx1_grammarParser) LabelledDirective() (localctx ILabelledDirectiveContext) {
	localctx = NewLabelledDirectiveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, vtx1_grammarParserRULE_labelledDirective)
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
		p.SetState(195)
		p.Match(vtx1_grammarParserIDENTIFIER)
	}
	p.SetState(197)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserCOLON {
		{
			p.SetState(196)
			p.Match(vtx1_grammarParserCOLON)
		}

	}
	{
		p.SetState(199)
		p.Directive()
	}

	return localctx
}
