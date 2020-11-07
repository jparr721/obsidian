package main

// tokenType represents a token by enum position
type tokenType int

// token represents a token value in a file
type token struct {
	variant tokenType
	lexeme  string
	literal interface{}
	line    int
}

// newtoken creates a new token from a given input sequence of values
func newToken(tokenType tokenType, lexeme string, literal interface{}, line int) token {
	return token{
		tokenType,
		lexeme,
		literal,
		line,
	}
}

const (
	// OSQUIGGLE represents a left square bracket
	OSQUIGGLE tokenType = iota

	// CSQUIGGLE represents a right square bracket
	CSQUIGGLE

	// OPAREN represents a left parenthesis
	OPAREN

	// CPAREN represents a right parenthesis
	CPAREN

	// COMMA represents a comma
	COMMA

	// SEMI represents a ;
	SEMI

	// DOT represents a ,
	DOT

	// PLUS represents a + sign
	PLUS

	// MINUS represents a - sign
	MINUS

	// STAR represents a * symbol
	STAR

	// SLASH represents a / symbol
	SLASH

	// MODULO represents a % operator
	MODULO

	// BANG represents a ! symbol
	BANG

	// BANGEQUAL represents a != symbol
	BANGEQUAL

	// EQUAL represents a = symbol
	EQUAL

	// EQUALEQUAL represents a == symbol
	EQUALEQUAL

	// GREATER represents a > symbol
	GREATER

	// GREATEREQUAL represents a >= symbol
	GREATEREQUAL

	// LESS represents a < symbol
	LESS

	// LESSEQUAL represents a <= symbol
	LESSEQUAL

	// IDENTIFIER represents a name
	IDENTIFIER

	// STRING represents a string value
	STRING

	// NUMBER represents a numeric value
	NUMBER

	// AND represents a and keyword
	AND

	// CLASS represents a class keyword
	CLASS

	// ELSE represents the else keyword
	ELSE

	// FALSE represents the false keyword
	FALSE

	// FUN represents the fun keyword
	FUN

	// FOR represents a the for keyword
	FOR

	// IF represents the if keyword
	IF

	// NIL represents the nil keyword
	NIL

	// OR represents the or keyword
	OR

	// PRINT represents the print keyword
	PRINT

	// RETURN represents the return keyword
	RETURN

	// SUPER represents the super keyword
	SUPER

	// THIS represents the this keyword
	THIS

	// TRUE represents the true keyword
	TRUE

	// VAR represents the var keyword
	VAR

	// WHILE represents the while keyword
	WHILE

	// BREAK represents the break keyword
	BREAK

	// EOF represents the end of file
	EOF

	// UNKNOWN represents an unknown, erroneous character
	UNKNOWN
)

// Keywords represents all of the keyword types
var keywords = map[string]tokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
	"break":  BREAK,
}
