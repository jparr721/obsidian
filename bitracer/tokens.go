package main

import "fmt"

// tokenType represents a token by enum position
type tokenType int

func (t tokenType) String() string {
	switch t {
	case OSQUIGGLE:
		return "OSQUIGGLE"
	case CSQUIGGLE:
		return "CSQUIGGLE"
	case OPAREN:
		return "OPAREN"
	case CPAREN:
		return "CPAREN"
	case COMMA:
		return "COMMA"
	case SEMI:
		return "SEMI"
	case DOT:
		return "DOT"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case STAR:
		return "STAR"
	case SLASH:
		return "SLASH"
	case MODULO:
		return "MODULO"
	case BANG:
		return "BANG"
	case BANGEQUAL:
		return "BANGEQUAL"
	case EQUAL:
		return "EQUAL"
	case EQUALEQUAL:
		return "EQUALEQUAL"
	case GREATER:
		return "GREATER"
	case GREATEREQUAL:
		return "GREATEREQUAL"
	case LESS:
		return "LESS"
	case LESSEQUAL:
		return "LESSEQUAL"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FUN:
		return "FUN"
	case FOR:
		return "FOR"
	case IF:
		return "IF"
	case NIL:
		return "NIL"
	case OR:
		return "OR"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case TRUE:
		return "TRUE"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	case EOF:
		return "EOF"
	}

	return "UNKNOWN"
}

func tokenFromString(s string) tokenType {
	switch s {
	case "OSQUIGGLE":
		return OSQUIGGLE
	case "CSQUIGGLE":
		return CSQUIGGLE
	case "OPAREN":
		return OPAREN
	case "CPAREN":
		return CPAREN
	case "COMMA":
		return COMMA
	case "SEMI":
		return SEMI
	case "DOT":
		return DOT
	case "PLUS":
		return PLUS
	case "MINUS":
		return MINUS
	case "STAR":
		return STAR
	case "SLASH":
		return SLASH
	case "MODULO":
		return MODULO
	case "BANG":
		return BANG
	case "BANGEQUAL":
		return BANGEQUAL
	case "EQUAL":
		return EQUAL
	case "EQUALEQUAL":
		return EQUALEQUAL
	case "GREATER":
		return GREATER
	case "GREATEREQUAL":
		return GREATEREQUAL
	case "LESS":
		return LESS
	case "LESSEQUAL":
		return LESSEQUAL
	case "IDENTIFIER":
		return IDENTIFIER
	case "STRING":
		return STRING
	case "NUMBER":
		return NUMBER
	case "AND":
		return AND
	case "CLASS":
		return CLASS
	case "ELSE":
		return ELSE
	case "FALSE":
		return FALSE
	case "FUN":
		return FUN
	case "FOR":
		return FOR
	case "IF":
		return IF
	case "NIL":
		return NIL
	case "OR":
		return OR
	case "PRINT":
		return PRINT
	case "RETURN":
		return RETURN
	case "SUPER":
		return SUPER
	case "THIS":
		return THIS
	case "TRUE":
		return TRUE
	case "VAR":
		return VAR
	case "WHILE":
		return WHILE
	case "EOF":
		return EOF
	}
	return UNKNOWN
}

// token represents a token value in a file
type token struct {
	variant   tokenType
	lexeme    string
	literal   interface{}
	line      int
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

func (t token) String() string {
	return fmt.Sprintf("%s %s %v", t.variant.String(), t.lexeme, t.literal)
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
}
