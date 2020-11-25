#ifndef __TOKENS_H__
#define __TOKENS_H__

typedef unsigned int token_t;
typedef union {
	int i;
	double d;
	char *c;
} literal_t;

struct token {
	token_t type;
	char *lexeme;
	int line;
	literal_t literal;
};

enum tokens {
	// Osquiggle represents a Left squiggle bracket
	OSQUIGGLE = 0x00,

	// Csquiggle Represents A Right squiggle Bracket
	CSQUIGGLE,

	// Oparen Represents A Left Parenthesis
	OPAREN,

	// Cparen Represents A Right Parenthesis
	CPAREN,

	// Comma Represents A Comma
	COMMA,

	// Semi Represents A ;
	SEMI,

	// Dot Represents A .
	DOT,

	// Plus Represents A + Sign
	PLUS,

	// Minus Represents A - Sign
	MINUS,

	// Star Represents A * Symbol
	STAR,

	// Slash Represents A / Symbol
	SLASH,

	// Modulo Represents A % Operator
	MODULO,

	// Bang Represents A ! Symbol
	BANG,

	// BangEqual Represents A != Symbol
	BANG_EQUAL,

	// Equal Represents A = Symbol
	EQUAL,

	// EqualEqual Represents A == Symbol
	EQUAL_EQUAL,

	// Greater Represents A > Symbol
	GREATER,

	// GreaterEqual Represents A >= Symbol
	GREATER_EQUAL,

	// Less Represents A < Symbol
	LESS,

	// LessEqual Represents A <= Symbol
	LESS_EQUAL,

	// Identifier Represents A Name
	IDENTIFIER,

	// String Represents A String Value
	STRING,

	// Number Represents A Numeric Value
	NUMBER,

	// And Represents A And Keyword
	AND,

	// Class Represents A Class Keyword
	CLASS,

	// Else Represents The Else Keyword
	ELSE,

	// False Represents The False Keyword
	FALSE,

	// Fun Represents The Fun Keyword
	FUN,

	// For Represents A The For Keyword
	FOR,

	// If Represents The If Keyword
	IF,

	// Nil Represents The Nil Keyword
	NIL,

	// Or Represents The Or Keyword
	OR,

	// Print Represents The Print Keyword
	PRINT,

	// Return Represents The Return Keyword
	RETURN,

	// Super Represents The Super Keyword
	SUPER,

	// This Represents The This Keyword
	THIS,

	// True Represents The True Keyword
	TRUE,

	// Var Represents The Var Keyword
	VAR,

	// While Represents The While Keyword
	WHILE,

	// Break Represents The Break Keyword
	BREAK,

	// EOF Represents The End Of File
	EOF
};

#endif /* TOKENS_H */
