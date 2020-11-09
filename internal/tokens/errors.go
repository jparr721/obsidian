package tokens

import "fmt"

// TokenizerError is an error that occurs during the tokenization process
type TokenizerError struct {
	line    int
	message string
}

func newTokenizerError(line int, message string) *TokenizerError {
	return &TokenizerError{line, message}
}

func (l *TokenizerError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s\n", l.line, l.message)
}

// TODO(@jparr721) - Add global "did error" state.
func reportTokenizerError(l *TokenizerError) {
	fmt.Println(l.Error())
}
