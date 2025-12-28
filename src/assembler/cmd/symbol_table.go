package cmd

import (
	"fmt"
)

// Symbol represents an entry in the symbol table, such as a label or a constant.
type Symbol struct {
	Name    string
	Address uint32
	Defined bool
	// Location of the symbol's definition
	File   string
	Line   int
	Column int
	// We can add more info here later, like the file and line number of definition.
}

// SymbolTable manages all symbols for the assembler.
type SymbolTable struct {
	symbols  map[string]*Symbol
	warnings *[]error // Pointer to ErrorManager.Warnings for reporting
}

// NewSymbolTable creates and returns a new SymbolTable.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]*Symbol),
	}
}

// AttachWarnings allows the symbol table to report warnings to ErrorManager.
func (st *SymbolTable) AttachWarnings(w *[]error) {
	st.warnings = w
}

// Define adds a new symbol to the table or updates its address if already present.
// It marks the symbol as defined and records its location.
func (st *SymbolTable) Define(name string, address uint32, file string, line, column int) (*Symbol, error) {
	if s, exists := st.symbols[name]; exists {
		if s.Defined {
			// Report as error, but also add a warning for shadowing
			if st.warnings != nil {
				*st.warnings = append(*st.warnings, fmt.Errorf("warning: label '%s' redefined at %s:%d:%d (original at %s:%d:%d)", name, file, line, column, s.File, s.Line, s.Column))
			}
			return nil, fmt.Errorf("duplicate definition of symbol '%s' at %s:%d:%d (originally defined at %s:%d:%d)", name, file, line, column, s.File, s.Line, s.Column)
		}
		s.Address = address
		s.Defined = true
		s.File = file
		s.Line = line
		s.Column = column
		return s, nil
	}

	newSymbol := &Symbol{
		Name:    name,
		Address: address,
		Defined: true,
		File:    file,
		Line:    line,
		Column:  column,
	}
	st.symbols[name] = newSymbol
	return newSymbol, nil
}

// Reference looks up a symbol. If it doesn't exist, it creates an undefined entry
// to be resolved later. This is key for handling forward references.
func (st *SymbolTable) Reference(name string) *Symbol {
	if s, exists := st.symbols[name]; exists {
		return s
	}

	newSymbol := &Symbol{
		Name:    name,
		Defined: false,
	}
	st.symbols[name] = newSymbol
	return newSymbol
}

// Lookup finds a symbol in the table. It returns the symbol and a boolean
// indicating whether it was found.
func (st *SymbolTable) Lookup(name string) (*Symbol, bool) {
	s, exists := st.symbols[name]
	return s, exists
}

// AllUndefined returns a slice of all symbols that were referenced but never defined.
func (st *SymbolTable) AllUndefined() []*Symbol {
	var undefined []*Symbol
	for _, s := range st.symbols {
		if !s.Defined {
			undefined = append(undefined, s)
		}
	}
	return undefined
}

// UnusedLabels returns a slice of all labels that were defined but never referenced.
func (st *SymbolTable) UnusedLabels() []*Symbol {
	var unused []*Symbol
	for _, s := range st.symbols {
		if s.Defined && s.Address != 0 && !st.isReferenced(s.Name) {
			unused = append(unused, s)
		}
	}
	return unused
}

// isReferenced checks if a symbol was referenced (stub: always false for now)
func (st *SymbolTable) isReferenced(name string) bool {
	// TODO: Implement real reference tracking
	return false
}
