package cmd

// AST represents the root of the abstract syntax tree
// for a VTX1 assembly program.
type AST struct {
	Program *ProgramNode
}

// ProgramNode represents the entire program (list of lines)
type ProgramNode struct {
	Lines []*LineNode
}

// LineNode represents a single line in the assembly source
// (may contain label, instruction/directive, comment)
type LineNode struct {
	Label     *LabelNode
	Statement StatementNode // InstructionNode, DirectiveNode, or VLIWInstructionNode
	Comment   string
	Line      int
	Column    int
}

// StatementNode is an interface for all statement types
// (Instruction, Directive, VLIWInstruction)
type StatementNode interface {
	isStatement()
}

// LabelNode represents a label definition
// e.g., main:
type LabelNode struct {
	Name   string
	Line   int
	Column int
}

// InstructionNode represents a single instruction
// e.g., ADD T0, T1, T2
type InstructionNode struct {
	Mnemonic string
	Operands []OperandNode
	Line     int
	Column   int
}

func (InstructionNode) isStatement() {}

// DirectiveNode represents an assembler directive
// e.g., .ORG 0x1000
type DirectiveNode struct {
	Name   string
	Params []OperandNode
	Line   int
	Column int
}

func (DirectiveNode) isStatement() {}

// VLIWInstructionNode represents a VLIW instruction group
// e.g., [ADD T0, T1, T2] [SUB T3, T4, T5]
type VLIWInstructionNode struct {
	Instructions []*InstructionNode
	Line         int
	Column       int
}

func (VLIWInstructionNode) isStatement() {}

// OperandNode is an interface for all operand types
type OperandNode interface {
	isOperand()
}

// RegisterNode represents a register operand
type RegisterNode struct {
	Name   string
	Line   int
	Column int
}

func (RegisterNode) isOperand() {}

// ImmediateNode represents an immediate value operand
type ImmediateNode struct {
	Value  string // Keep as string for now (can parse to int/ternary later)
	Line   int
	Column int
}

func (ImmediateNode) isOperand() {}

// MemoryOperandNode represents a memory operand
// e.g., [T0+T1] or [T0+0x10]
type MemoryOperandNode struct {
	Base   string // Register name
	Index  string // Optional index register
	Offset string // Optional offset (immediate)
	Line   int
	Column int
}

func (MemoryOperandNode) isOperand() {}

// IdentifierNode represents a symbol or label reference as an operand
type IdentifierNode struct {
	Name   string
	Line   int
	Column int
}

func (IdentifierNode) isOperand() {}

// Add BinaryOpNode definition here for use by ASTBuilder and codegen
type BinaryOpNode struct {
	Left   OperandNode
	Op     string
	Right  OperandNode
	Line   int
	Column int
}
