package interpreter

import (
	"fmt"

	"github.com/jparr721/obsidian/internal/statement"
)

// Callable represents a callable type which takes arguments and an interpreter instance
type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error)
}

// Function represents a function callable
type Function struct {
	Declaration *statement.FunctionStatement
}

func NewFunction(declaration *statement.FunctionStatement) *Function {
	return &Function{declaration}
}

func (f *Function) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	environment := NewEnvironment(interpreter.globals)

	for i, arg := range arguments {
		lexeme := f.Declaration.Arguments[i].Lexeme
		environment.define(lexeme, arg)
	}

	err := interpreter.executeBlock(f.Declaration.Body, environment)

	if err != nil {
		switch err.(type) {
		case *RuntimeError:
			return nil, err
		case *ReturnInterrupt:
			return err.(*ReturnInterrupt).value, nil
		}
	}

	return nil, nil
}

func (f *Function) Arity() int {
	return len(f.Declaration.Arguments)
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}
