package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type tokenizer struct {
	src     string
	tokens  []token
	start   int
	current int
	line    int
}

func newTokenizer(src string) *tokenizer {
	return &tokenizer{
		src:     src,
		tokens:  make([]token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (t *tokenizer) scanTokens() []token {
	for !t.end() {
		t.start = t.current

		// Prevents the last space from getting parsed
		if t.current+1 >= len(t.src) {
			break
		}

		t.scanToken()
	}

	t.tokens = append(t.tokens, newToken(EOF, "", nil, t.line))
	return t.tokens
}

func (t *tokenizer) scanToken() {
	c := t.next()
	switch c {
	case "(":
		t.addToken(OPAREN, nil)
		break
	case ")":
		t.addToken(CPAREN, nil)
		break
	case "{":
		t.addToken(OSQUIGGLE, nil)
		break
	case "}":
		t.addToken(CSQUIGGLE, nil)
		break
	case ",":
		t.addToken(COMMA, nil)
		break
	case ".":
		t.addToken(DOT, nil)
		break
	case "-":
		t.addToken(MINUS, nil)
		break
	case "+":
		t.addToken(PLUS, nil)
		break
	case ";":
		t.addToken(SEMI, nil)
		break
	case "*":
		t.addToken(STAR, nil)
		break
	case "!":
		if t.match("=") {
			t.addToken(BANGEQUAL, nil)
		} else {
			t.addToken(BANG, nil)
		}
		break
	case "=":
		if t.match("=") {
			t.addToken(EQUALEQUAL, nil)
		} else {
			t.addToken(EQUAL, nil)
		}
		break
	case "<":
		if t.match("=") {
			t.addToken(LESSEQUAL, nil)
		} else {
			t.addToken(LESS, nil)
		}
		break
	case ">":
		if t.match("=") {
			t.addToken(GREATEREQUAL, nil)
		} else {
			t.addToken(GREATER, nil)
		}
		break
	case "/":
		if t.match("/") {
			for t.peek() != "\\n" && !t.end() {
				t.next()
			}
		} else {
			t.addToken(SLASH, nil)
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
			lineError(t.line, fmt.Sprintf("Unexpected character: %s", c))
		}
		break
	}
}

func (t *tokenizer) isAlphaOrUnderscore(c string) bool {
	reg, err := regexp.Compile("^[a-zA-Z0-9_]+")

	if err != nil {
		lineError(t.line, err.Error())
		return false
	}

	return reg.MatchString(c)
}

func (t *tokenizer) isDigit(c string) bool {
	if _, err := strconv.Atoi(c); err == nil {
		return true
	}
	return false
}

func (t *tokenizer) parseIdentifier() {
	for t.isAlphaOrUnderscore(t.peek()) {
		t.next()
	}

	text := t.src[t.start:t.current]
	tokenType, ok := keywords[text]

	if !ok {
		t.addToken(IDENTIFIER, nil)
		return
	}

	t.addToken(tokenType, nil)
}

func (t *tokenizer) parseNumber() {
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
		lineError(t.line, err.Error())
	}

	t.addToken(NUMBER, value)
}

func (t *tokenizer) parseString() {
	for t.peek() != "\"" && !t.end() {
		if t.peek() == "\\n" {
			t.line++
		}
		t.next()
	}

	if t.end() {
		lineError(t.line, "Unexpected string.")
		return
	}

	t.next()

	// Trim quotes
	value := t.src[t.start+1 : t.current-1]
	t.addToken(STRING, value)
}

func (t *tokenizer) peek() string {
	if t.end() {
		return "\\0"
	}

	return string(t.src[t.current])
}

func (t *tokenizer) peekNext() string {
	if t.end() {
		return "\\0"
	}

	return string(t.src[t.current+1])
}

func (t *tokenizer) match(expected string) bool {
	if t.end() {
		return false
	}

	if string(t.src[t.current]) != expected {
		return false
	}

	t.current++
	return true
}

func (t *tokenizer) addToken(tokenType tokenType, literal interface{}) {
	text := t.src[t.start:t.current]
	t.tokens = append(t.tokens, newToken(tokenType, text, literal, t.line))
}

func (t *tokenizer) next() string {
	t.current++
	return string(t.src[t.current-1])
}

func (t *tokenizer) end() bool {
	return t.current >= len(t.src)
}
