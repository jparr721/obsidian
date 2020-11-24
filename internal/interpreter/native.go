package interpreter

import (
	"fmt"
	"time"
)

// native.go implement's obsidian's native function interface

type clockFunction struct{}

func (c *clockFunction) String() string {
	return "<native fn: 'clock'>"
}
func (c *clockFunction) Arity() int { return 0 }
func (c *clockFunction) Call(interpreter Interpreter, arguments ...interface{}) (interface{}, error) {
	return time.Now(), nil
}

type printFunction struct{}

func (p *printFunction) String() string {
	return "<native fn: 'print'>"
}

func (p *printFunction) Arity() int { return 254 }
func (p *printFunction) Call(interpreter Interpreter, arguments ...interface{}) (interface{}, error) {
	template := arguments[0].(string)

	fmt.Printf(template, arguments[1:])
	return nil, nil
}
