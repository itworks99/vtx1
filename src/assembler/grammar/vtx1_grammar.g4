grammar vtx1_grammar;

options {
    language=Go;
}

// === Parser Rules ===

program
    : (line | blankLine)* (END_DIRECTIVE (EOL)?)? EOF
    ;

line
    : (label | instruction | directive | labelledDirective | vliwInstruction | macroDefinition | comment) (comment)? EOL
    ;

blankLine
    : EOL
    ;

label
    : IDENTIFIER COLON
    ;

comment
    : COMMENT
    ;

instruction
    : mnemonic (operand (COMMA operand)*)?
    ;

vliwInstruction
    : LSQUARE instruction RSQUARE
      (LSQUARE instruction RSQUARE)?
      (LSQUARE instruction RSQUARE)?
    ;

mnemonic
    : ALU_OP
    | MEM_OP
    | CTRL_OP
    | VEC_OP
    | FP_OP
    | SYS_OP
    | COMPLEX_OP
    | COMPLEX_VEC
    | COMPLEX_MEM
    | COMPLEX_SYS
    ;

operand
    : register
    | immediate
    | memoryOperand
    | IDENTIFIER
    | IDENTIFIER PLUS immediate
    | IDENTIFIER PLUS IDENTIFIER
    | immediate PLUS IDENTIFIER
    | immediate PLUS immediate
    ;

register
    : TB_REG
    | GPR
    | SPECIAL_REG
    | VECTOR_REG
    | FP_REG
    ;

memoryOperand
    : LSQUARE baseRegister (PLUS (indexRegister | offsetImmediate))? RSQUARE
    ;

baseRegister
    : TB_REG
    | GPR
    ;

indexRegister
    : GPR
    ;

offsetImmediate
    : immediate
    ;

immediate
    : DECIMAL
    | HEXADECIMAL
    | BINARY
    | TERNARY
    ;

directive
    : ORG_DIRECTIVE immediate
    | DATA_DIRECTIVE dataList
    | EQU_DIRECTIVE IDENTIFIER COMMA immediate
    | INCLUDE_DIRECTIVE STRING
    | SECTION_DIRECTIVE IDENTIFIER
    | ALIGN_DIRECTIVE immediate
    | SPACE_DIRECTIVE immediate
    ;

dataList
    : dataItem (COMMA dataItem)*
    ;

dataItem
    : immediate
    | STRING
    ;

// Macro definition
macroDefinition
    : MACRO_DIRECTIVE IDENTIFIER (IDENTIFIER)* EOL
      macroBody
      ENDM_DIRECTIVE EOL
    ;

macroBody
    : line+
    ;

// === Lexer Rules ===

WHITESPACE      : [ \t]+ -> channel(HIDDEN);
COMMENT         : ';' ~[\r\n]*;
EOL             : ('\r'? '\n')+;

// Mnemonics - ALU operations
ALU_OP          : 'ADD' | 'SUB' | 'MUL' | 'AND' | 'OR' | 'NOT' | 'XOR'
                | 'SHL' | 'SHR' | 'ROL' | 'ROR' | 'CMP' | 'TEST' | 'INC'
                | 'DEC' | 'NEG';

// Mnemonics - Memory operations
MEM_OP          : 'LD' | 'ST' | 'VLD' | 'VST' | 'FLD' | 'FST' | 'LEA'
                | 'PUSH' | 'POP';

// Mnemonics - Control operations
CTRL_OP         : 'JMP' | 'JAL' | 'JR' | 'JALR' | 'BEQ' | 'BNE' | 'BLT'
                | 'BGE' | 'BLTU' | 'BGEU' | 'BGT' | 'BLE' | 'CALL' | 'RET';

// Mnemonics - Vector operations
VEC_OP          : 'VADD' | 'VSUB' | 'VMUL' | 'VAND' | 'VOR' | 'VNOT'
                | 'VSHL' | 'VSHR';

// Mnemonics - Floating point operations
FP_OP           : 'FADD' | 'FSUB' | 'FMUL' | 'FCMP' | 'FMOV' | 'FNEG';

// Mnemonics - System operations
SYS_OP          : 'NOP' | 'WFI';

// Mnemonics - Complex operations (microcode)
COMPLEX_OP      : 'DIV' | 'MOD' | 'UDIV' | 'UMOD' | 'SQRT' | 'ABS'
                | 'SIN' | 'COS' | 'TAN' | 'ASIN' | 'ACOS' | 'ATAN'
                | 'EXP' | 'LOG';

// Mnemonics - Complex vector operations
COMPLEX_VEC     : 'VDOT' | 'VREDUCE' | 'VMAX' | 'VMIN' | 'VSUM' | 'VPERM';

// Mnemonics - Complex memory operations
COMPLEX_MEM     : 'CACHE' | 'FLUSH' | 'MEMBAR';

// Mnemonics - Complex system operations
COMPLEX_SYS     : 'SYSCALL' | 'BREAK' | 'HALT';

// Registers
TB_REG          : 'TB';
GPR             : 'T'[0-6];
SPECIAL_REG     : 'TA' | 'TC' | 'TS' | 'TI';
VECTOR_REG      : 'VA' | 'VT' | 'VB';
FP_REG          : 'FA' | 'FT' | 'FB';

// Literals
DECIMAL         : [0-9]+;
HEXADECIMAL     : '0x' [0-9a-fA-F]+;
BINARY          : '0b' [01]+;
TERNARY         : '0t' [+\-0]+;
STRING          : '"' ~["]* '"';

// Structural
COLON           : ':';
COMMA           : ',';
PLUS            : '+';
LSQUARE         : '[';
RSQUARE         : ']';

// Directives
ORG_DIRECTIVE   : '.ORG';
DATA_DIRECTIVE  : '.DB' | '.DW' | '.DT';
EQU_DIRECTIVE   : '.EQU';
INCLUDE_DIRECTIVE : '.INCLUDE';
SECTION_DIRECTIVE : '.SECTION';
ALIGN_DIRECTIVE : '.ALIGN';
SPACE_DIRECTIVE : '.SPACE';
END_DIRECTIVE   : 'END';
IDENTIFIER      : [a-zA-Z_] [a-zA-Z0-9_]*;

// Add lexer rules for macros
MACRO_DIRECTIVE : '.MACRO';
ENDM_DIRECTIVE  : '.ENDM';

// Add this rule after label:
labelledDirective
    : IDENTIFIER (COLON)? directive
    ;
