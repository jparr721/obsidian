package compiler

import "regexp"

// TokenKind represents the type of token in a given operation sequence
type TokenKind int

// Token represents a parsed token value from an input operation sequence
type Token struct {
	Name  string
	Value interface{}
}

const (
	// UNKNOWN is the unknown token, this will throw on compile
	UNKNOWN TokenKind = iota

	// OSQUARE represents a left square bracket
	OSQUARE

	// CSQUARE represents a right square bracket
	CSQUARE

	// ALL represents a for-all loop
	ALL

	// SEMI represents a semicolon
	SEMI

	// VARIABLE represents a variable
	VARIABLE

	// OPAREN represents a left parenthesis
	OPAREN

	// CPAREN represents a right parenthesis
	CPAREN

	// COMMA represents a comma
	COMMA

	// WHERE represents a where statement
	WHERE

	// PLUS represents a plus sign
	PLUS

	// MINUS represents a minus sign
	MINUS

	// MULTIPLY represents a multiplication symbol
	MULTIPLY

	// DIVIDE represents a division symbol
	DIVIDE

	// MODULO represents a modulus operator
	MODULO

	// GREATERTHAN represents a greater than operator
	GREATERTHAN

	// LESSTHAN represents a less than operator
	LESSTHAN

	// AND represents the `and` keyword
	AND

	// OR represents the `or` keyword
	OR

	// EQUALS represents an equals sign
	EQUALS

	// POW represents raising a value to a power
	POW
)

// BinaryOps are binary operations that runtime methods can perform
var BinaryOps = []TokenKind{
	EQUALS,
	AND,
	OR,
	PLUS,
	MINUS,
	MULTIPLY,
	DIVIDE,
	MODULO,
	GREATERTHAN,
	LESSTHAN,
}

func (kind TokenKind) String() string {
	switch kind {

	case OSQUARE:
		return "OSQUARE"
	case CSQUARE:
		return "CSQUARE"
	case ALL:
		return "ALL"
	case SEMI:
		return "SEMI"
	case VARIABLE:
		return "VARIABLE"
	case OPAREN:
		return "OPAREN"
	case CPAREN:
		return "CPAREN"
	case COMMA:
		return "COMMA"
	case WHERE:
		return "WHERE"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case MODULO:
		return "MODULO"
	case GREATERTHAN:
		return "GREATERTHAN"
	case LESSTHAN:
		return "LESSTHAN"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case EQUALS:
		return "EQUALS"
	case POW:
		return "POW"
	}

	return "UNKNOWN"
}

func (kind TokenKind) Regex() *regexp.Regexp {
	switch kind {

	case OSQUARE:
		return regexp.MustCompile("(?P<OSQUARE>\\A(\\[))")
	case CSQUARE:
		return regexp.MustCompile("(?P<CSQUARE>\\A(\\]))")
	case ALL:
		return regexp.MustCompile("(?P<ALL>\\A(ALL))")
	case SEMI:
		return regexp.MustCompile("(?P<SEMI>\\A(;))")
	case VARIABLE:
		return regexp.MustCompile("(?P<VARIABLE>\\A(^[a-zA-Z]\\w*$))")
	case OPAREN:
		return regexp.MustCompile("(?P<OPAREN>\\A(\\())")
	case CPAREN:
		return regexp.MustCompile("(?P<CPAREN>\\A(\\)))")
	case COMMA:
		return regexp.MustCompile("(?P<COMMA>\\A(,))")
	case WHERE:
		return regexp.MustCompile("(?P<WHERE>\\A(WHERE))")
	case PLUS:
		return regexp.MustCompile("(?P<PLUS>\\A(\\+))")
	case MINUS:
		return regexp.MustCompile("(?P<MINUS>\\A(-))")
	case MULTIPLY:
		return regexp.MustCompile("(?P<MULTIPLY>\\A(\\*))")
	case DIVIDE:
		return regexp.MustCompile("(?P<DIVIDE>\\A(\\/))")
	case MODULO:
		return regexp.MustCompile("(?P<MODULO>\\A(\\%))")
	case GREATERTHAN:
		return regexp.MustCompile("(?P<GREATERTHAN>\\A(\\>))")
	case LESSTHAN:
		return regexp.MustCompile("(?P<LESSTHAN>\\A(\\<))")
	case AND:
		return regexp.MustCompile("(?P<AND>\\A(AND))")
	case OR:
		return regexp.MustCompile("(?P<OR>\\A(OR))")
	case EQUALS:
		return regexp.MustCompile("(?P<EQUALS>\\A(=))")
	case POW:
		return regexp.MustCompile("(?P<POW>\\A(\\^))")
	}

	return regexp.MustCompile("(?<UNKNOWN>\\A(UNKNOWN))")
}
