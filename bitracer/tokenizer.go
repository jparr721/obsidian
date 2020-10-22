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
		t.addtoken(OPAREN, nil)
	case ")":
		t.addtoken(CPAREN, nil)
	case "{":
		t.addtoken(OSQUIGGLE, nil)
	case "}":
		t.addtoken(CSQUIGGLE, nil)
	case ",":
		t.addtoken(COMMA, nil)
	case ".":
		t.addtoken(DOT, nil)
	case "-":
		t.addtoken(MINUS, nil)
	case "+":
		t.addtoken(PLUS, nil)
	case ";":
		t.addtoken(SEMI, nil)
	case "*":
		t.addtoken(STAR, nil)
	case "!":
		if t.match("=") {
			t.addtoken(BANGEQUAL, nil)
		} else {
			t.addtoken(BANG, nil)
		}
	case "=":
		if t.match("=") {
			t.addtoken(EQUALEQUAL, nil)
		} else {
			t.addtoken(EQUAL, nil)
		}
	case "<":
		if t.match("=") {
			t.addtoken(LESSEQUAL, nil)
		} else {
			t.addtoken(LESS, nil)
		}
	case ">":
		if t.match("=") {
			t.addtoken(GREATEREQUAL, nil)
		} else {
			t.addtoken(GREATER, nil)
		}
	case "/":
		if t.match("/") {
			for t.peek() != "\\n" && !t.end() {
				t.next()
			}
		} else {
			t.addtoken(SLASH, nil)
		}
	case " ":
	case "\\r":
	case "\\t":
		break
	case "\\n":
		t.line++
	case "\"":
		t.parseString()
	default:
		if t.isDigit(c) {
			t.parseNumber()
		} else if t.isAlphaOrUnderscore(c) {
			t.parseIdentifier()
		} else {
			lineError(t.line, fmt.Sprintf("Unexpected character: %s", c))
		}
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
		t.addtoken(IDENTIFIER, nil)
		return
	}

	t.addtoken(tokenType, nil)
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

	t.addtoken(NUMBER, value)
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
	t.addtoken(STRING, value)
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

func (t *tokenizer) addtoken(tokenType tokenType, literal interface{}) {
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
