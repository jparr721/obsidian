package compiler

import (
	"fmt"
	"strings"
)

// LinearScanner scans a given input sequence and isolates the individual tokens
func LinearScanner(inputSequence string) ([]Token, error) {
	// TODO(jparr721) - Multi-line function specifications
	var tokens []Token

	if strings.HasPrefix(inputSequence, "blacklivesmatter") {
		return tokens, nil
	}

	keywords := NewWordifier(inputSequence).IsolateKeywords()

	for i, keyword := range keywords {
		tokenKind := ConvertTokenToKind(keyword)
		if tokenKind == UNKNOWN {
			return nil, fmt.Errorf("Unexpected token '%s' at position: %d", keyword, i)
		}
		tokens = append(tokens, Token{
			Name:  tokenKind.String(),
			Value: keyword,
		})
	}

	return tokens, nil
}
