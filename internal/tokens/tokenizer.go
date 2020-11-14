package tokens

import (
	"fmt"
	"regexp"
	"strconv"
)

// Tokenizer represents the PL tokenizer
type Tokenizer struct {
	src     string
	Tokens  []Token
	start   int
	current int
	line    int
}

// NewTokenizer creates a new tokenizer from a source string of values
func NewTokenizer(src string) *Tokenizer {
	return &Tokenizer{
		src:     src,
		Tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

// ScanTokens scans the internal token array and reports errors, or the scanned tokens if successful
func (t *Tokenizer) ScanTokens() ([]Token, *TokenizerError) {
	for !t.end() {
		t.start = t.current

		err := t.scanToken()

		if err != nil {
			return nil, err
		}

	}

	t.Tokens = append(t.Tokens, NewToken(TokenEOF, "", nil, t.line))
	return t.Tokens, nil
}

func (t *Tokenizer) scanToken() *TokenizerError {
	c := t.next()
	switch c {
	case "(":
		t.addToken(TokenOparen, nil)
		break
	case ")":
		t.addToken(TokenCparen, nil)
		break
	case "{":
		t.addToken(TokenOsquiggle, nil)
		break
	case "}":
		t.addToken(TokenCsquiggle, nil)
		break
	case ",":
		t.addToken(TokenComma, nil)
		break
	case ".":
		t.addToken(TokenDot, nil)
		break
	case "-":
		t.addToken(TokenMinus, nil)
		break
	case "+":
		t.addToken(TokenPlus, nil)
		break
	case ";":
		t.addToken(TokenSemi, nil)
		break
	case "*":
		t.addToken(TokenStar, nil)
		break
	case "!":
		if t.match("=") {
			t.addToken(TokenBangEqual, nil)
		} else {
			t.addToken(TokenBang, nil)
		}
		break
	case "=":
		if t.match("=") {
			t.addToken(TokenEqualEqual, nil)
		} else {
			t.addToken(TokenEqual, nil)
		}
		break
	case "<":
		if t.match("=") {
			t.addToken(TokenLessEqual, nil)
		} else {
			t.addToken(TokenLess, nil)
		}
		break
	case ">":
		if t.match("=") {
			t.addToken(TokenGreaterEqual, nil)
		} else {
			t.addToken(TokenGreater, nil)
		}
		break
	case "/":
		if t.match("/") {
			for t.peek() != "\n" && !t.end() {
				t.next()
			}
		} else {
			t.addToken(TokenSlash, nil)
		}
		break
	case " ":
	case "":
	case "\r":
	case "\t":
		break
	case "\n":
		t.line++
		break
	case "\"":
		t.parseString()
		break
	default:
		if t.isDigit(c) {
			t.parseNumber()
		} else if t.isAlphaOrUnderscore(c) {
			t.parseIdentifier()
		} else {
			return newTokenizerError(t.line, fmt.Sprintf("Unexpected character: %s", c))
		}
		break
	}

	return nil
}

func (t *Tokenizer) isAlphaOrUnderscore(c string) bool {
	reg, err := regexp.Compile("^[a-zA-Z0-9_]+")

	if err != nil {
		newTokenizerError(t.line, err.Error())
		return false
	}

	return reg.MatchString(c)
}

func (t *Tokenizer) isDigit(c string) bool {
	if _, err := strconv.Atoi(c); err == nil {
		return true
	}
	return false
}

func (t *Tokenizer) parseIdentifier() {
	for t.isAlphaOrUnderscore(t.peek()) {
		t.next()
	}

	text := t.src[t.start:t.current]
	TokenType, ok := Keywords[text]

	if !ok {
		t.addToken(TokenIdentifier, nil)
		return
	}

	t.addToken(TokenType, nil)
}

func (t *Tokenizer) parseNumber() {
	for t.isDigit(t.peek()) {
		t.next()
	}

	// Check for decimal
	if t.peek() == "." && t.isDigit(t.peekNext()) {
		// eat the decimal
		t.next()

		for t.isDigit(t.peek()) {
			t.next()
		}
	}

	value, err := strconv.ParseFloat(t.src[t.start:t.current], 64)

	if err != nil {
		newTokenizerError(t.line, err.Error())
	}

	t.addToken(TokenNumber, value)
}

func (t *Tokenizer) parseString() {
	for t.peek() != "\"" && !t.end() {
		if t.peek() == "\\n" {
			t.line++
		}
		t.next()
	}

	if t.end() {
		newTokenizerError(t.line, "Unexpected string.")
		return
	}

	t.next()

	// Trim quotes
	value := t.src[t.start+1 : t.current-1]
	t.addToken(TokenString, value)
}

func (t *Tokenizer) peek() string {
	if t.end() {
		return "\\0"
	}

	return string(t.src[t.current])
}

func (t *Tokenizer) peekNext() string {
	if t.end() {
		return "\\0"
	}

	return string(t.src[t.current+1])
}

func (t *Tokenizer) match(expected string) bool {
	if t.end() {
		return false
	}

	if string(t.src[t.current]) != expected {
		return false
	}

	t.current++
	return true
}

func (t *Tokenizer) addToken(TokenType TokenType, literal interface{}) {
	text := t.src[t.start:t.current]
	t.Tokens = append(t.Tokens, NewToken(TokenType, text, literal, t.line))
}

func (t *Tokenizer) next() string {
	t.current++
	return string(t.src[t.current-1])
}

func (t *Tokenizer) end() bool {
	return t.current >= len(t.src)
}
