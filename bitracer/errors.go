package main

import (
	"fmt"
	"os"
	"time"
)

type runtimeError struct {
	token   token
	message string
}

func newRuntimeError(token token, message string) *runtimeError {
	return &runtimeError{token, message}
}

func (r *runtimeError) Error() string {
	return fmt.Sprintf("[line %d] %s", r.token.line, r.message)
}

func reportRuntimeError(e *runtimeError) {
	fmt.Println(e.Error())
	hadRuntimeError = true
}

type parseError struct {
	token   token
	message string
}

func newParseError(token token, message string) *parseError {
	return &parseError{token, message}
}

func (p *parseError) Error() string {
	var pos string

	if p.token.variant == EOF {
		pos = "at end"
	} else {
		pos = "at '" + p.token.lexeme + "'"
	}

	return fmt.Sprintf("[line %d] Error %s: %s\n", p.token.line, pos, p.message)
}

func reportParseError(e *parseError) {
	fmt.Println(e.Error())
	hadParseError = true
}

func lineError(line int, message string) {
	reportError(line, "", message)
}

func reportError(line int, where, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
	hadParseError = true
}

// Rudimentary core dump function. This should include more useful intel later, but it should do for now.
// TODO(@jparr721) - Consider this: https://stackoverflow.com/questions/52103182/how-to-get-the-stacktrace-of-a-panic-and-store-as-a-variable
func coreDump(metadata interface{}) {
	fileName := fmt.Sprintf("core_dump_%s.log", time.Now().Format(time.RFC3339))
	metadataStr := fmt.Sprintf("%v", metadata)

	f, err := os.Create(fileName)
	_, err = f.WriteString(metadataStr)
	err = f.Close()
	if err != nil {
		fmt.Printf("Last ditch, failed to dump core to file, puking here. Include this with your issue please")
		fmt.Println(metadataStr)
		return
	}
}
