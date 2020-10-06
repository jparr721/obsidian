package compiler

import (
	"fmt"
	"strings"
)

// LinearScanner scans a given input sequence and isolates the individual tokens
func LinearScanner(inputSequence []string) (map[int][]Token, error) {
	var tokens map[int][]Token

  for _, line := range inputSequence {
    if strings.HasPrefix(line, "#") {
      continue
    }

    keywords := NewWordifier(line).IsolateKeywords()
    tkns := make([]Token, 0)

    for i, keyword := range keywords {
      tokenKind := ConvertTokenToKind(keyword)
      if tokenKind == UNKNOWN {
        return nil, fmt.Errorf("Unexpected token '%s' at position: %d", keyword, i)
      }

      tkns = append(tkns, Token{
        Name:  tokenKind.String(),
        Value: keyword,
      })
    }
  }

    return tokens, nil
}
