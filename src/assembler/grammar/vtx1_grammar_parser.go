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
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "'TB'", "",
		"", "", "", "", "", "", "", "", "':'", "','", "'+'", "'['", "']'", "'.ORG'",
		"", "'.EQU'", "'.INCLUDE'", "'.SECTION'", "'.ALIGN'", "'.SPACE'", "'END'",
		"", "'.MACRO'", "'.ENDM'",
	}
	staticData.SymbolicNames = []string{
		"", "WHITESPACE", "COMMENT", "EOL", "ALU_OP", "MEM_OP", "CTRL_OP", "VEC_OP",
		"FP_OP", "SYS_OP", "COMPLEX_OP", "COMPLEX_VEC", "COMPLEX_MEM", "COMPLEX_SYS",
		"TB_REG", "GPR", "SPECIAL_REG", "VECTOR_REG", "FP_REG", "DECIMAL", "HEXADECIMAL",
		"BINARY", "TERNARY", "STRING", "COLON", "COMMA", "PLUS", "LSQUARE",
		"RSQUARE", "ORG_DIRECTIVE", "DATA_DIRECTIVE", "EQU_DIRECTIVE", "INCLUDE_DIRECTIVE",
		"SECTION_DIRECTIVE", "ALIGN_DIRECTIVE", "SPACE_DIRECTIVE", "END_DIRECTIVE",
		"IDENTIFIER", "MACRO_DIRECTIVE", "ENDM_DIRECTIVE",
	}
	staticData.RuleNames = []string{
		"program", "line", "blankLine", "label", "comment", "instruction", "vliwInstruction",
		"mnemonic", "operand", "register", "memoryOperand", "baseRegister",
		"indexRegister", "offsetImmediate", "immediate", "directive", "dataList",
		"dataItem", "macroDefinition", "macroBody", "labelledDirective",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 39, 202, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 1,
		0, 1, 0, 5, 0, 45, 8, 0, 10, 0, 12, 0, 48, 9, 0, 1, 0, 1, 0, 3, 0, 52,
		8, 0, 3, 0, 54, 8, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 3, 1, 65, 8, 1, 1, 1, 3, 1, 68, 8, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3,
		1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 5, 5, 83, 8, 5, 10, 5,
		12, 5, 86, 9, 5, 3, 5, 88, 8, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1,
		6, 3, 6, 97, 8, 6, 1, 6, 1, 6, 1, 6, 1, 6, 3, 6, 103, 8, 6, 1, 7, 1, 7,
		1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8,
		1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 125, 8, 8, 1, 9, 1, 9, 1, 10,
		1, 10, 1, 10, 1, 10, 1, 10, 3, 10, 134, 8, 10, 3, 10, 136, 8, 10, 1, 10,
		1, 10, 1, 11, 1, 11, 1, 12, 1, 12, 1, 13, 1, 13, 1, 14, 1, 14, 1, 15, 1,
		15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15,
		1, 15, 1, 15, 1, 15, 1, 15, 3, 15, 164, 8, 15, 1, 16, 1, 16, 1, 16, 5,
		16, 169, 8, 16, 10, 16, 12, 16, 172, 9, 16, 1, 17, 1, 17, 3, 17, 176, 8,
		17, 1, 18, 1, 18, 1, 18, 5, 18, 181, 8, 18, 10, 18, 12, 18, 184, 9, 18,
		1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 19, 4, 19, 192, 8, 19, 11, 19, 12,
		19, 193, 1, 20, 1, 20, 3, 20, 198, 8, 20, 1, 20, 1, 20, 1, 20, 0, 0, 21,
		0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36,
		38, 40, 0, 4, 1, 0, 4, 13, 1, 0, 14, 18, 1, 0, 14, 15, 1, 0, 19, 22, 215,
		0, 46, 1, 0, 0, 0, 2, 64, 1, 0, 0, 0, 4, 71, 1, 0, 0, 0, 6, 73, 1, 0, 0,
		0, 8, 76, 1, 0, 0, 0, 10, 78, 1, 0, 0, 0, 12, 89, 1, 0, 0, 0, 14, 104,
		1, 0, 0, 0, 16, 124, 1, 0, 0, 0, 18, 126, 1, 0, 0, 0, 20, 128, 1, 0, 0,
		0, 22, 139, 1, 0, 0, 0, 24, 141, 1, 0, 0, 0, 26, 143, 1, 0, 0, 0, 28, 145,
		1, 0, 0, 0, 30, 163, 1, 0, 0, 0, 32, 165, 1, 0, 0, 0, 34, 175, 1, 0, 0,
		0, 36, 177, 1, 0, 0, 0, 38, 191, 1, 0, 0, 0, 40, 195, 1, 0, 0, 0, 42, 45,
		3, 2, 1, 0, 43, 45, 3, 4, 2, 0, 44, 42, 1, 0, 0, 0, 44, 43, 1, 0, 0, 0,
		45, 48, 1, 0, 0, 0, 46, 44, 1, 0, 0, 0, 46, 47, 1, 0, 0, 0, 47, 53, 1,
		0, 0, 0, 48, 46, 1, 0, 0, 0, 49, 51, 5, 36, 0, 0, 50, 52, 5, 3, 0, 0, 51,
		50, 1, 0, 0, 0, 51, 52, 1, 0, 0, 0, 52, 54, 1, 0, 0, 0, 53, 49, 1, 0, 0,
		0, 53, 54, 1, 0, 0, 0, 54, 55, 1, 0, 0, 0, 55, 56, 5, 0, 0, 1, 56, 1, 1,
		0, 0, 0, 57, 65, 3, 6, 3, 0, 58, 65, 3, 10, 5, 0, 59, 65, 3, 30, 15, 0,
		60, 65, 3, 40, 20, 0, 61, 65, 3, 12, 6, 0, 62, 65, 3, 36, 18, 0, 63, 65,
		3, 8, 4, 0, 64, 57, 1, 0, 0, 0, 64, 58, 1, 0, 0, 0, 64, 59, 1, 0, 0, 0,
		64, 60, 1, 0, 0, 0, 64, 61, 1, 0, 0, 0, 64, 62, 1, 0, 0, 0, 64, 63, 1,
		0, 0, 0, 65, 67, 1, 0, 0, 0, 66, 68, 3, 8, 4, 0, 67, 66, 1, 0, 0, 0, 67,
		68, 1, 0, 0, 0, 68, 69, 1, 0, 0, 0, 69, 70, 5, 3, 0, 0, 70, 3, 1, 0, 0,
		0, 71, 72, 5, 3, 0, 0, 72, 5, 1, 0, 0, 0, 73, 74, 5, 37, 0, 0, 74, 75,
		5, 24, 0, 0, 75, 7, 1, 0, 0, 0, 76, 77, 5, 2, 0, 0, 77, 9, 1, 0, 0, 0,
		78, 87, 3, 14, 7, 0, 79, 84, 3, 16, 8, 0, 80, 81, 5, 25, 0, 0, 81, 83,
		3, 16, 8, 0, 82, 80, 1, 0, 0, 0, 83, 86, 1, 0, 0, 0, 84, 82, 1, 0, 0, 0,
		84, 85, 1, 0, 0, 0, 85, 88, 1, 0, 0, 0, 86, 84, 1, 0, 0, 0, 87, 79, 1,
		0, 0, 0, 87, 88, 1, 0, 0, 0, 88, 11, 1, 0, 0, 0, 89, 90, 5, 27, 0, 0, 90,
		91, 3, 10, 5, 0, 91, 96, 5, 28, 0, 0, 92, 93, 5, 27, 0, 0, 93, 94, 3, 10,
		5, 0, 94, 95, 5, 28, 0, 0, 95, 97, 1, 0, 0, 0, 96, 92, 1, 0, 0, 0, 96,
		97, 1, 0, 0, 0, 97, 102, 1, 0, 0, 0, 98, 99, 5, 27, 0, 0, 99, 100, 3, 10,
		5, 0, 100, 101, 5, 28, 0, 0, 101, 103, 1, 0, 0, 0, 102, 98, 1, 0, 0, 0,
		102, 103, 1, 0, 0, 0, 103, 13, 1, 0, 0, 0, 104, 105, 7, 0, 0, 0, 105, 15,
		1, 0, 0, 0, 106, 125, 3, 18, 9, 0, 107, 125, 3, 28, 14, 0, 108, 125, 3,
		20, 10, 0, 109, 125, 5, 37, 0, 0, 110, 111, 5, 37, 0, 0, 111, 112, 5, 26,
		0, 0, 112, 125, 3, 28, 14, 0, 113, 114, 5, 37, 0, 0, 114, 115, 5, 26, 0,
		0, 115, 125, 5, 37, 0, 0, 116, 117, 3, 28, 14, 0, 117, 118, 5, 26, 0, 0,
		118, 119, 5, 37, 0, 0, 119, 125, 1, 0, 0, 0, 120, 121, 3, 28, 14, 0, 121,
		122, 5, 26, 0, 0, 122, 123, 3, 28, 14, 0, 123, 125, 1, 0, 0, 0, 124, 106,
		1, 0, 0, 0, 124, 107, 1, 0, 0, 0, 124, 108, 1, 0, 0, 0, 124, 109, 1, 0,
		0, 0, 124, 110, 1, 0, 0, 0, 124, 113, 1, 0, 0, 0, 124, 116, 1, 0, 0, 0,
		124, 120, 1, 0, 0, 0, 125, 17, 1, 0, 0, 0, 126, 127, 7, 1, 0, 0, 127, 19,
		1, 0, 0, 0, 128, 129, 5, 27, 0, 0, 129, 135, 3, 22, 11, 0, 130, 133, 5,
		26, 0, 0, 131, 134, 3, 24, 12, 0, 132, 134, 3, 26, 13, 0, 133, 131, 1,
		0, 0, 0, 133, 132, 1, 0, 0, 0, 134, 136, 1, 0, 0, 0, 135, 130, 1, 0, 0,
		0, 135, 136, 1, 0, 0, 0, 136, 137, 1, 0, 0, 0, 137, 138, 5, 28, 0, 0, 138,
		21, 1, 0, 0, 0, 139, 140, 7, 2, 0, 0, 140, 23, 1, 0, 0, 0, 141, 142, 5,
		15, 0, 0, 142, 25, 1, 0, 0, 0, 143, 144, 3, 28, 14, 0, 144, 27, 1, 0, 0,
		0, 145, 146, 7, 3, 0, 0, 146, 29, 1, 0, 0, 0, 147, 148, 5, 29, 0, 0, 148,
		164, 3, 28, 14, 0, 149, 150, 5, 30, 0, 0, 150, 164, 3, 32, 16, 0, 151,
		152, 5, 31, 0, 0, 152, 153, 5, 37, 0, 0, 153, 154, 5, 25, 0, 0, 154, 164,
		3, 28, 14, 0, 155, 156, 5, 32, 0, 0, 156, 164, 5, 23, 0, 0, 157, 158, 5,
		33, 0, 0, 158, 164, 5, 37, 0, 0, 159, 160, 5, 34, 0, 0, 160, 164, 3, 28,
		14, 0, 161, 162, 5, 35, 0, 0, 162, 164, 3, 28, 14, 0, 163, 147, 1, 0, 0,
		0, 163, 149, 1, 0, 0, 0, 163, 151, 1, 0, 0, 0, 163, 155, 1, 0, 0, 0, 163,
		157, 1, 0, 0, 0, 163, 159, 1, 0, 0, 0, 163, 161, 1, 0, 0, 0, 164, 31, 1,
		0, 0, 0, 165, 170, 3, 34, 17, 0, 166, 167, 5, 25, 0, 0, 167, 169, 3, 34,
		17, 0, 168, 166, 1, 0, 0, 0, 169, 172, 1, 0, 0, 0, 170, 168, 1, 0, 0, 0,
		170, 171, 1, 0, 0, 0, 171, 33, 1, 0, 0, 0, 172, 170, 1, 0, 0, 0, 173, 176,
		3, 28, 14, 0, 174, 176, 5, 23, 0, 0, 175, 173, 1, 0, 0, 0, 175, 174, 1,
		0, 0, 0, 176, 35, 1, 0, 0, 0, 177, 178, 5, 38, 0, 0, 178, 182, 5, 37, 0,
		0, 179, 181, 5, 37, 0, 0, 180, 179, 1, 0, 0, 0, 181, 184, 1, 0, 0, 0, 182,
		180, 1, 0, 0, 0, 182, 183, 1, 0, 0, 0, 183, 185, 1, 0, 0, 0, 184, 182,
		1, 0, 0, 0, 185, 186, 5, 3, 0, 0, 186, 187, 3, 38, 19, 0, 187, 188, 5,
		39, 0, 0, 188, 189, 5, 3, 0, 0, 189, 37, 1, 0, 0, 0, 190, 192, 3, 2, 1,
		0, 191, 190, 1, 0, 0, 0, 192, 193, 1, 0, 0, 0, 193, 191, 1, 0, 0, 0, 193,
		194, 1, 0, 0, 0, 194, 39, 1, 0, 0, 0, 195, 197, 5, 37, 0, 0, 196, 198,
		5, 24, 0, 0, 197, 196, 1, 0, 0, 0, 197, 198, 1, 0, 0, 0, 198, 199, 1, 0,
		0, 0, 199, 200, 3, 30, 15, 0, 200, 41, 1, 0, 0, 0, 19, 44, 46, 51, 53,
		64, 67, 84, 87, 96, 102, 124, 133, 135, 163, 170, 175, 182, 193, 197,
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

	// Getter signatures
	EOF() antlr.TerminalNode
	AllLine() []ILineContext
	Line(i int) ILineContext
	AllBlankLine() []IBlankLineContext
	BlankLine(i int) IBlankLineContext
	END_DIRECTIVE() antlr.TerminalNode
	EOL() antlr.TerminalNode

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

func (s *ProgramContext) AllBlankLine() []IBlankLineContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IBlankLineContext); ok {
			len++
		}
	}

	tst := make([]IBlankLineContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IBlankLineContext); ok {
			tst[i] = t.(IBlankLineContext)
			i++
		}
	}

	return tst
}

func (s *ProgramContext) BlankLine(i int) IBlankLineContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBlankLineContext); ok {
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
	p.SetState(46)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&480633700348) != 0 {
		p.SetState(44)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

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
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(48)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(53)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserEND_DIRECTIVE {
		{
			p.SetState(49)
			p.Match(vtx1_grammarParserEND_DIRECTIVE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(51)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == vtx1_grammarParserEOL {
			{
				p.SetState(50)
				p.Match(vtx1_grammarParserEOL)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	}
	{
		p.SetState(55)
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
	LabelledDirective() ILabelledDirectiveContext
	VliwInstruction() IVliwInstructionContext
	MacroDefinition() IMacroDefinitionContext
	AllComment() []ICommentContext
	Comment(i int) ICommentContext

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

func (s *LineContext) LabelledDirective() ILabelledDirectiveContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILabelledDirectiveContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILabelledDirectiveContext)
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

func (s *LineContext) MacroDefinition() IMacroDefinitionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMacroDefinitionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMacroDefinitionContext)
}

func (s *LineContext) AllComment() []ICommentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICommentContext); ok {
			len++
		}
	}

	tst := make([]ICommentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICommentContext); ok {
			tst[i] = t.(ICommentContext)
			i++
		}
	}

	return tst
}

func (s *LineContext) Comment(i int) ICommentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICommentContext); ok {
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
	p.SetState(64)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext()) {
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

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.SetState(67)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
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

// IBlankLineContext is an interface to support dynamic dispatch.
type IBlankLineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOL() antlr.TerminalNode

	// IsBlankLineContext differentiates from other interfaces.
	IsBlankLineContext()
}

type BlankLineContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBlankLineContext() *BlankLineContext {
	var p = new(BlankLineContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_blankLine
	return p
}

func InitEmptyBlankLineContext(p *BlankLineContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_blankLine
}

func (*BlankLineContext) IsBlankLineContext() {}

func NewBlankLineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BlankLineContext {
	var p = new(BlankLineContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *BlankLineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterBlankLine(s)
	}
}

func (s *BlankLineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitBlankLine(s)
	}
}

func (s *BlankLineContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitBlankLine(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) BlankLine() (localctx IBlankLineContext) {
	localctx = NewBlankLineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, vtx1_grammarParserRULE_blankLine)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(71)
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
	p.EnterRule(localctx, 6, vtx1_grammarParserRULE_label)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(73)
		p.Match(vtx1_grammarParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(74)
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
	p.EnterRule(localctx, 8, vtx1_grammarParserRULE_comment)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(76)
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
	p.EnterRule(localctx, 10, vtx1_grammarParserRULE_instruction)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(78)
		p.Mnemonic()
	}
	p.SetState(87)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&137581543424) != 0 {
		{
			p.SetState(79)
			p.Operand()
		}
		p.SetState(84)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == vtx1_grammarParserCOMMA {
			{
				p.SetState(80)
				p.Match(vtx1_grammarParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(81)
				p.Operand()
			}

			p.SetState(86)
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
	p.EnterRule(localctx, 12, vtx1_grammarParserRULE_vliwInstruction)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(89)
		p.Match(vtx1_grammarParserLSQUARE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(90)
		p.Instruction()
	}
	{
		p.SetState(91)
		p.Match(vtx1_grammarParserRSQUARE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(96)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(92)
			p.Match(vtx1_grammarParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(93)
			p.Instruction()
		}
		{
			p.SetState(94)
			p.Match(vtx1_grammarParserRSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(102)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserLSQUARE {
		{
			p.SetState(98)
			p.Match(vtx1_grammarParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(99)
			p.Instruction()
		}
		{
			p.SetState(100)
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
	p.EnterRule(localctx, 14, vtx1_grammarParserRULE_mnemonic)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(104)
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
	AllImmediate() []IImmediateContext
	Immediate(i int) IImmediateContext
	MemoryOperand() IMemoryOperandContext
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	PLUS() antlr.TerminalNode

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

func (s *OperandContext) AllImmediate() []IImmediateContext {
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

func (s *OperandContext) Immediate(i int) IImmediateContext {
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
	p.EnterRule(localctx, 16, vtx1_grammarParserRULE_operand)
	p.SetState(124)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 10, p.GetParserRuleContext()) {
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(110)
			p.Match(vtx1_grammarParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(111)
			p.Match(vtx1_grammarParserPLUS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(114)
			p.Match(vtx1_grammarParserPLUS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(115)
			p.Match(vtx1_grammarParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(120)
			p.Immediate()
		}
		{
			p.SetState(121)
			p.Match(vtx1_grammarParserPLUS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(122)
			p.Immediate()
		}

	case antlr.ATNInvalidAltNumber:
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
	TB_REG() antlr.TerminalNode
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
	p.EnterRule(localctx, 18, vtx1_grammarParserRULE_register)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(126)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&507904) != 0) {
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
	p.EnterRule(localctx, 20, vtx1_grammarParserRULE_memoryOperand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(128)
		p.Match(vtx1_grammarParserLSQUARE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(129)
		p.BaseRegister()
	}
	p.SetState(135)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserPLUS {
		{
			p.SetState(130)
			p.Match(vtx1_grammarParserPLUS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(133)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

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
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

	}
	{
		p.SetState(137)
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
	p.EnterRule(localctx, 22, vtx1_grammarParserRULE_baseRegister)
	var _la int

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
	p.EnterRule(localctx, 24, vtx1_grammarParserRULE_indexRegister)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(141)
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
	p.EnterRule(localctx, 26, vtx1_grammarParserRULE_offsetImmediate)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(143)
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
	p.EnterRule(localctx, 28, vtx1_grammarParserRULE_immediate)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(145)
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
	p.EnterRule(localctx, 30, vtx1_grammarParserRULE_directive)
	p.SetState(163)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case vtx1_grammarParserORG_DIRECTIVE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(147)
			p.Match(vtx1_grammarParserORG_DIRECTIVE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(152)
			p.Match(vtx1_grammarParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(153)
			p.Match(vtx1_grammarParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(156)
			p.Match(vtx1_grammarParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case vtx1_grammarParserSECTION_DIRECTIVE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(157)
			p.Match(vtx1_grammarParserSECTION_DIRECTIVE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(158)
			p.Match(vtx1_grammarParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case vtx1_grammarParserALIGN_DIRECTIVE:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(159)
			p.Match(vtx1_grammarParserALIGN_DIRECTIVE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(162)
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
	AllDataItem() []IDataItemContext
	DataItem(i int) IDataItemContext
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

func (s *DataListContext) AllDataItem() []IDataItemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDataItemContext); ok {
			len++
		}
	}

	tst := make([]IDataItemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDataItemContext); ok {
			tst[i] = t.(IDataItemContext)
			i++
		}
	}

	return tst
}

func (s *DataListContext) DataItem(i int) IDataItemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDataItemContext); ok {
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
	p.EnterRule(localctx, 32, vtx1_grammarParserRULE_dataList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(165)
		p.DataItem()
	}
	p.SetState(170)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == vtx1_grammarParserCOMMA {
		{
			p.SetState(166)
			p.Match(vtx1_grammarParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(167)
			p.DataItem()
		}

		p.SetState(172)
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

// IDataItemContext is an interface to support dynamic dispatch.
type IDataItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Immediate() IImmediateContext
	STRING() antlr.TerminalNode

	// IsDataItemContext differentiates from other interfaces.
	IsDataItemContext()
}

type DataItemContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDataItemContext() *DataItemContext {
	var p = new(DataItemContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_dataItem
	return p
}

func InitEmptyDataItemContext(p *DataItemContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_dataItem
}

func (*DataItemContext) IsDataItemContext() {}

func NewDataItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DataItemContext {
	var p = new(DataItemContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_dataItem

	return p
}

func (s *DataItemContext) GetParser() antlr.Parser { return s.parser }

func (s *DataItemContext) Immediate() IImmediateContext {
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

func (s *DataItemContext) STRING() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserSTRING, 0)
}

func (s *DataItemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DataItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DataItemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterDataItem(s)
	}
}

func (s *DataItemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitDataItem(s)
	}
}

func (s *DataItemContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitDataItem(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) DataItem() (localctx IDataItemContext) {
	localctx = NewDataItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, vtx1_grammarParserRULE_dataItem)
	p.SetState(175)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

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

// IMacroDefinitionContext is an interface to support dynamic dispatch.
type IMacroDefinitionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MACRO_DIRECTIVE() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AllEOL() []antlr.TerminalNode
	EOL(i int) antlr.TerminalNode
	MacroBody() IMacroBodyContext
	ENDM_DIRECTIVE() antlr.TerminalNode

	// IsMacroDefinitionContext differentiates from other interfaces.
	IsMacroDefinitionContext()
}

type MacroDefinitionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMacroDefinitionContext() *MacroDefinitionContext {
	var p = new(MacroDefinitionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_macroDefinition
	return p
}

func InitEmptyMacroDefinitionContext(p *MacroDefinitionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_macroDefinition
}

func (*MacroDefinitionContext) IsMacroDefinitionContext() {}

func NewMacroDefinitionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MacroDefinitionContext {
	var p = new(MacroDefinitionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMacroBodyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

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

func (s *MacroDefinitionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterMacroDefinition(s)
	}
}

func (s *MacroDefinitionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitMacroDefinition(s)
	}
}

func (s *MacroDefinitionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitMacroDefinition(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) MacroDefinition() (localctx IMacroDefinitionContext) {
	localctx = NewMacroDefinitionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, vtx1_grammarParserRULE_macroDefinition)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(177)
		p.Match(vtx1_grammarParserMACRO_DIRECTIVE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(178)
		p.Match(vtx1_grammarParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(182)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == vtx1_grammarParserIDENTIFIER {
		{
			p.SetState(179)
			p.Match(vtx1_grammarParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(184)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(185)
		p.Match(vtx1_grammarParserEOL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(186)
		p.MacroBody()
	}
	{
		p.SetState(187)
		p.Match(vtx1_grammarParserENDM_DIRECTIVE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(188)
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

// IMacroBodyContext is an interface to support dynamic dispatch.
type IMacroBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllLine() []ILineContext
	Line(i int) ILineContext

	// IsMacroBodyContext differentiates from other interfaces.
	IsMacroBodyContext()
}

type MacroBodyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMacroBodyContext() *MacroBodyContext {
	var p = new(MacroBodyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_macroBody
	return p
}

func InitEmptyMacroBodyContext(p *MacroBodyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_macroBody
}

func (*MacroBodyContext) IsMacroBodyContext() {}

func NewMacroBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MacroBodyContext {
	var p = new(MacroBodyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_macroBody

	return p
}

func (s *MacroBodyContext) GetParser() antlr.Parser { return s.parser }

func (s *MacroBodyContext) AllLine() []ILineContext {
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

func (s *MacroBodyContext) Line(i int) ILineContext {
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

func (s *MacroBodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MacroBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MacroBodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterMacroBody(s)
	}
}

func (s *MacroBodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitMacroBody(s)
	}
}

func (s *MacroBodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitMacroBody(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) MacroBody() (localctx IMacroBodyContext) {
	localctx = NewMacroBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, vtx1_grammarParserRULE_macroBody)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(191)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&480633700340) != 0) {
		{
			p.SetState(190)
			p.Line()
		}

		p.SetState(193)
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

// ILabelledDirectiveContext is an interface to support dynamic dispatch.
type ILabelledDirectiveContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	Directive() IDirectiveContext
	COLON() antlr.TerminalNode

	// IsLabelledDirectiveContext differentiates from other interfaces.
	IsLabelledDirectiveContext()
}

type LabelledDirectiveContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLabelledDirectiveContext() *LabelledDirectiveContext {
	var p = new(LabelledDirectiveContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_labelledDirective
	return p
}

func InitEmptyLabelledDirectiveContext(p *LabelledDirectiveContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = vtx1_grammarParserRULE_labelledDirective
}

func (*LabelledDirectiveContext) IsLabelledDirectiveContext() {}

func NewLabelledDirectiveContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LabelledDirectiveContext {
	var p = new(LabelledDirectiveContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = vtx1_grammarParserRULE_labelledDirective

	return p
}

func (s *LabelledDirectiveContext) GetParser() antlr.Parser { return s.parser }

func (s *LabelledDirectiveContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserIDENTIFIER, 0)
}

func (s *LabelledDirectiveContext) Directive() IDirectiveContext {
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

func (s *LabelledDirectiveContext) COLON() antlr.TerminalNode {
	return s.GetToken(vtx1_grammarParserCOLON, 0)
}

func (s *LabelledDirectiveContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LabelledDirectiveContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LabelledDirectiveContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.EnterLabelledDirective(s)
	}
}

func (s *LabelledDirectiveContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(vtx1_grammarListener); ok {
		listenerT.ExitLabelledDirective(s)
	}
}

func (s *LabelledDirectiveContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case vtx1_grammarVisitor:
		return t.VisitLabelledDirective(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *vtx1_grammarParser) LabelledDirective() (localctx ILabelledDirectiveContext) {
	localctx = NewLabelledDirectiveContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, vtx1_grammarParserRULE_labelledDirective)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(195)
		p.Match(vtx1_grammarParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(197)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == vtx1_grammarParserCOLON {
		{
			p.SetState(196)
			p.Match(vtx1_grammarParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(199)
		p.Directive()
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
