package interpreter

import (
	"testing"

	"github.com/jparr721/obsidian/internal/tokens"
)

type defineTest struct {
	Name     string
	VarName  string
	VarValue interface{}
	Expected interface{}
}

func TestDefine(t *testing.T) {
	tests := []defineTest{
		{
			Name:     "Test define variable with value",
			VarName:  "foo",
			VarValue: "bar",
			Expected: "bar",
		},
		{
			Name:     "Test define variable with no value",
			VarName:  "foo",
			VarValue: "bar",
			Expected: "bar",
		},
	}

	for _, test := range tests {
		e := NewEnvironment(nil)

		t.Logf("Running: %s\n", test.Name)

		e.define(test.VarName, test.VarValue)

		if e.values[test.VarName] != test.VarValue {
			t.Error("defined value did not match test value")
		}
	}
}

type getTest struct {
	Name      string
	Key       tokens.Token
	VarName   string
	VarValue  interface{}
	Enclosing *environment
	Expected  interface{}
}

func buildEnvironment(t *testing.T, k string, v interface{}) *environment {
	e := NewEnvironment(nil)
	e.define(k, v)
	if e.values[k] != v {
		t.Error("Failed configuring test environment")
	}

	return e
}

func TestGet(t *testing.T) {
	tests := []getTest{
		{
			Name:      "Returns expected value in environment",
			Key:       tokens.NewToken(tokens.TokenIdentifier, "foo", nil, 1),
			VarName:   "foo",
			VarValue:  "bar",
			Enclosing: nil,
			Expected:  "bar",
		},
		{
			Name:      "Returns nil value when unset",
			Key:       tokens.NewToken(tokens.TokenIdentifier, "foo", nil, 1),
			VarName:   "foo",
			VarValue:  nil,
			Enclosing: nil,
			Expected:  nil,
		},
		{
			Name:      "Returns error when name and value unset",
			Key:       tokens.NewToken(tokens.TokenIdentifier, "bar", nil, 1),
			VarName:   "",
			VarValue:  nil,
			Enclosing: nil,
			Expected:  newRuntimeError(tokens.NewToken(tokens.TokenIdentifier, "bar", nil, 1), "Undefined variable 'bar'"),
		},
		{
			Name:      "Gets from enclosing when enclosing is not nil",
			Key:       tokens.NewToken(tokens.TokenIdentifier, "foo", nil, 1),
			VarName:   "bar",
			VarValue:  "baz",
			Enclosing: buildEnvironment(t, "foo", "bar"),
			Expected:  "bar",
		},
	}

	for _, test := range tests {
		e := NewEnvironment(test.Enclosing)

		t.Logf("Running: %s\n", test.Name)

		e.define(test.VarName, test.VarValue)
		value, err := e.get(test.Key)

		if err != nil {
			if err.Error() != test.Expected.(error).Error() {
				t.Errorf("returned error: '%v' did not match expected error: '%v'", err, test.Expected)
			}
			return
		}

		if value != test.Expected {
			t.Errorf("defined value: '%v' did not match expected value: '%v'", value, test.Expected)
		}
	}
}

type assignTest struct {
	Name string
	
}