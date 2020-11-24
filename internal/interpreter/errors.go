package interpreter

import (
	"fmt"

	"github.com/jparr721/obsidian/internal/tokens"
)

// RuntimeError is a type of error that happens during the interpret step of the compiler's execution.
type RuntimeError struct {
	token   tokens.Token
	message string
}

func newRuntimeError(token tokens.Token, message string) *RuntimeError {
	return &RuntimeError{token, message}
}

func (r *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: [line %d] %s", r.token.Line, r.message)
}

// ReturnInterrupt represents the return statement and its value as an error to break nested calls.
type ReturnInterrupt struct {
	value interface{}
}

func newReturnInterrupt(value interface{}) *ReturnInterrupt {
	return &ReturnInterrupt{value}
}

func (r *ReturnInterrupt) Error() string {
	return stringify(r.value)
}
