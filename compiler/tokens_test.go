package compiler

import "testing"

type TokensTest struct {
	Name     string
	Kind     TokenKind
	Input    string
	Expected bool
}

func TestRegularExpressionMatching(test *testing.T) {
	allTests := []TokensTest{
		{
			Name:     "OSQUARE",
			Kind:     OSQUARE,
			Input:    "[",
			Expected: true,
		},
		{
			Name:     "CSQUARE",
			Kind:     CSQUARE,
			Input:    "]",
			Expected: true,
		},
		{
			Name:     "ALL",
			Kind:     ALL,
			Input:    "ALL",
			Expected: true,
		},
		{
			Name:     "SEMI",
			Kind:     SEMI,
			Input:    ";",
			Expected: true,
		},
		{
			Name:     "VARIABLE",
			Kind:     VARIABLE,
			Input:    "x",
			Expected: true,
		},
		{
			Name:     "VARIABLE",
			Kind:     VARIABLE,
			Input:    "x1",
			Expected: true,
		},
		{
			Name:     "VARIABLE",
			Kind:     VARIABLE,
			Input:    "1x",
			Expected: false,
		},
		{
			Name:     "OPAREN",
			Kind:     OPAREN,
			Input:    "(",
			Expected: true,
		},
		{
			Name:     "CPAREN",
			Kind:     CPAREN,
			Input:    ")",
			Expected: true,
		},
		{
			Name:     "COMMA",
			Kind:     COMMA,
			Input:    ",",
			Expected: true,
		},
		{
			Name:     "WHERE",
			Kind:     WHERE,
			Input:    "WHERE",
			Expected: true,
		},
		{
			Name:     "PLUS",
			Kind:     PLUS,
			Input:    "+",
			Expected: true,
		},
		{
			Name:     "MINUS",
			Kind:     MINUS,
			Input:    "-",
			Expected: true,
		},
		{
			Name:     "MULTIPLY",
			Kind:     MULTIPLY,
			Input:    "*",
			Expected: true,
		},
		{
			Name:     "DIVIDE",
			Kind:     DIVIDE,
			Input:    "/",
			Expected: true,
		},
		{
			Name:     "MODULO",
			Kind:     MODULO,
			Input:    "%",
			Expected: true,
		},
		{
			Name:     "GREATERTHAN",
			Kind:     GREATERTHAN,
			Input:    ">",
			Expected: true,
		},
		{
			Name:     "LESSTHAN",
			Kind:     LESSTHAN,
			Input:    "<",
			Expected: true,
		},
		{
			Name:     "AND",
			Kind:     AND,
			Input:    "AND",
			Expected: true,
		},
		{
			Name:     "OR",
			Kind:     OR,
			Input:    "OR",
			Expected: true,
		},
		{
			Name:     "EQUALS",
			Kind:     EQUALS,
			Input:    "=",
			Expected: true,
		},
		{
			Name:     "POW",
			Kind:     POW,
			Input:    "^",
			Expected: true,
		},
		{
			Name:     "DO",
			Kind:     DO,
			Input:    "do",
			Expected: true,
		},
		{
			Name:     "END",
			Kind:     END,
			Input:    "end",
			Expected: true,
		},
		{
			Name:     "IF",
			Kind:     IF,
			Input:    "if",
			Expected: true,
		},
		{
			Name:     "ELSE",
			Kind:     ELSE,
			Input:    "else",
			Expected: true,
		},
		{
			Name:     "QUOTE",
			Kind:     QUOTE,
			Input:    "\"",
			Expected: true,
		},
		{
			Name:     "FUNCTION_GOOD",
			Kind:     FUNCTION,
			Input:    "Walk",
			Expected: true,
		},
		{
			Name:     "FUNCTION_BAD",
			Kind:     FUNCTION,
			Input:    "walk",
			Expected: false,
		},
		{
			Name:     "STRINGVALUE",
			Kind:     STRINGVALUE,
			Input:    "\"foobar\"",
			Expected: true,
		},
		{
			Name:     "STRINGVALUE_FAIL",
			Kind:     STRINGVALUE,
			Input:    "foobar",
			Expected: false,
		},
		{
			Name:     "NUMERICVALUE",
			Kind:     NUMERICVALUE,
			Input:    "12345",
			Expected: true,
		},
		{
			Name:     "NUMERICVALUE_FAIL",
			Kind:     NUMERICVALUE,
			Input:    "\"foobar\"",
			Expected: false,
		},
	}

	for _, regexTest := range allTests {
		test.Logf("TestRegularExpressionMatching -- %s", regexTest.Name)
		re := regexTest.Kind.Regex()
		if re.MatchString(regexTest.Input) != regexTest.Expected {
			test.Logf("failing test: %s", regexTest.Name)
			test.Logf("Input: %s", regexTest.Input)
			test.Logf("Expected: %t", regexTest.Expected)
			test.Fail()
		}
	}
}
