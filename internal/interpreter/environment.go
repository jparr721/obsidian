package interpreter

import (
	"fmt"

	"github.com/jparr721/obsidian/internal/tokens"
)

type environment struct {
	enclosing *environment
	values    map[string]interface{}
}

func newEnvironment(enclosing *environment) *environment {
	return &environment{enclosing, make(map[string]interface{})}
}

// TODO(@jparr721) - Make values immutable
func (e *environment) define(name string, value interface{}) {
	if value != nil {
		e.values[name] = value
		return
	}

	e.values[name] = nil
}

func (e *environment) get(name tokens.Token) (interface{}, *RuntimeError) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, newRuntimeError(name, fmt.Sprintf("Undefined variable '%s'", name.Lexeme))
}

func (e *environment) assign(name tokens.Token, value interface{}) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}

	return newRuntimeError(name, fmt.Sprintf("Undefined variable '%s'", name.Lexeme))
}
