package cmd

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

// CustomErrorListener implements antlr.ErrorListener to provide custom error handling.
type CustomErrorListener struct {
	*antlr.DefaultErrorListener
	Errors *ErrorManager
}

// NewCustomErrorListener creates a new CustomErrorListener.
func NewCustomErrorListener(errors *ErrorManager) *CustomErrorListener {
	return &CustomErrorListener{
		DefaultErrorListener: antlr.NewDefaultErrorListener(),
		Errors:               errors,
	}
}

// SyntaxError is called by ANTLR when a syntax error is detected.
func (l *CustomErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// Format a more user-friendly error message
	l.Errors.Errors = append(l.Errors.Errors, fmt.Errorf("syntax error at line %d:%d: %s", line, column, msg))
}
