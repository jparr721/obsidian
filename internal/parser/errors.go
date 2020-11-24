package parser

import (
	"fmt"

	"github.com/jparr721/obsidian/internal/tokens"
)

// ParseError represents an error during the parsing process
type ParseError struct {
	token   tokens.Token
	message string
	handled bool
}

func newParseError(token tokens.Token, message string) *ParseError {
	return &ParseError{token, message, false}
}

func (p *ParseError) Error() string {
	var pos string

	if p.token.Variant == tokens.TokenEOF {
		pos = "at end"
	} else {
		pos = "at '" + p.token.Lexeme + "'"
	}

	return fmt.Sprintf("ParseError: [line %d] Error %s: %s\n", p.token.Line, pos, p.message)
}

// ReportParseError geneerates the proper error formatting from a given error type
func ReportParseError(e *ParseError) {
	fmt.Println(e.Error())
	e.handled = true
}
