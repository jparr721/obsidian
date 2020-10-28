package main

import "fmt"

type environment struct {
	enclosing *environment
	values    map[string]interface{}
}

func newEnvironment(enclosing *environment) *environment {
	var values map[string]interface{}
	return &environment{enclosing, values}
}

// TODO(@jparr721) - Make values immutable
func (e *environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *environment) get(name token) (interface{}, *runtimeError) {
	if value, ok := e.values[name.lexeme]; ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, newRuntimeError(name, fmt.Sprintf("Undefined variable '%s'"))
}

func (e *environment) assign(name token, value interface{}) *runtimeError {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}

	return newRuntimeError(name, fmt.Sprintf("Undefined variable '%s'"))
}
