package main

import "fmt"

type runtimeError struct {
	token   token
	message string
}

func newRuntimeError(token token, message string) *runtimeError {
	return &runtimeError{token, message}
}

func (r *runtimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", r.message, r.token.line)
}

func reportRuntimeError(e *runtimeError) {
	fmt.Println(e.Error())
	hadRuntimeError = true
}

func lineError(line int, message string) {
	reportError(line, "", message)
}

func parseError(t token, m string) {
	if t.variant == EOF {
		reportError(t.line, " at end", m)
	} else {
		reportError(t.line, " at "+t.lexeme+"'", m)
	}
}

func reportError(line int, where, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
	hadParseError = true
}
