package tokens

// TokenType represents a token by enum position
type TokenType int

// Token represents a token value in a file
type Token struct {
	Variant TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

// NewToken creates a new token from a given input sequence of values
func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		tokenType,
		lexeme,
		literal,
		line,
	}
}

const (
	// TokenOsquiggle represents a Left squiggle bracket
	TokenOsquiggle TokenType = iota

	// TokenCsquiggle Represents A Right squiggle Bracket
	TokenCsquiggle

	// TokenOparen Represents A Left Parenthesis
	TokenOparen

	// TokenCparen Represents A Right Parenthesis
	TokenCparen

	// TokenComma Represents A Comma
	TokenComma

	// TokenSemi Represents A ;
	TokenSemi

	// TokenDot Represents A ,
	TokenDot

	// TokenPlus Represents A + Sign
	TokenPlus

	// TokenMinus Represents A - Sign
	TokenMinus

	// TokenStar Represents A * Symbol
	TokenStar

	// TokenSlash Represents A / Symbol
	TokenSlash

	// TokenModulo Represents A % Operator
	TokenModulo

	// TokenBang Represents A ! Symbol
	TokenBang

	// TokenBangEqual Represents A != Symbol
	TokenBangEqual

	// TokenEqual Represents A = Symbol
	TokenEqual

	// TokenEqualEqual Represents A == Symbol
	TokenEqualEqual

	// TokenGreater Represents A > Symbol
	TokenGreater

	// TokenGreaterEqual Represents A >= Symbol
	TokenGreaterEqual

	// TokenLess Represents A < Symbol
	TokenLess

	// TokenLessEqual Represents A <= Symbol
	TokenLessEqual

	// TokenIdentifier Represents A Name
	TokenIdentifier

	// TokenString Represents A String Value
	TokenString

	// TokenNumber Represents A Numeric Value
	TokenNumber

	// TokenAnd Represents A And Keyword
	TokenAnd

	// TokenClass Represents A Class Keyword
	TokenClass

	// TokenElse Represents The Else Keyword
	TokenElse

	// TokenFalse Represents The False Keyword
	TokenFalse

	// TokenFun Represents The Fun Keyword
	TokenFun

	// TokenFor Represents A The For Keyword
	TokenFor

	// TokenIf Represents The If Keyword
	TokenIf

	// TokenNil Represents The Nil Keyword
	TokenNil

	// TokenOr Represents The Or Keyword
	TokenOr

	// TokenPrint Represents The Print Keyword
	TokenPrint

	// TokenReturn Represents The Return Keyword
	TokenReturn

	// TokenSuper Represents The Super Keyword
	TokenSuper

	// TokenThis Represents The This Keyword
	TokenThis

	// TokenTrue Represents The True Keyword
	TokenTrue

	// TokenVar Represents The Var Keyword
	TokenVar

	// TokenWhile Represents The While Keyword
	TokenWhile

	// TokenBreak Represents The Break Keyword
	TokenBreak

	// TokenEOF Represents The End Of File
	TokenEOF

	// TokenUnknown Represents An Unknown, Erroneous Character
	TokenUnknown
)

// Keywords represents all of the keyword types
var Keywords = map[string]TokenType{
	"and":    TokenAnd,
	"class":  TokenClass,
	"else":   TokenElse,
	"false":  TokenFalse,
	"for":    TokenFor,
	"fun":    TokenFun,
	"if":     TokenIf,
	"nil":    TokenNil,
	"or":     TokenOr,
	"print":  TokenPrint,
	"return": TokenReturn,
	"super":  TokenSuper,
	"this":   TokenThis,
	"true":   TokenTrue,
	"var":    TokenVar,
	"while":  TokenWhile,
	"break":  TokenBreak,
}
