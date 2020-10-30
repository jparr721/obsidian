package main

import (
	"fmt"
	"reflect"
	"strconv"
)

type interpreter struct {
	environment *environment
}

func newInterpreter() *interpreter {
	return &interpreter{newEnvironment(nil)}
}

func (i *interpreter) interpret(statements []stmt) {
	for _, statement := range statements {
		result := i.execute(statement)

		switch result.(type) {
		case *runtimeError:
			reportRuntimeError(result.(*runtimeError))
		case error:
			coreDump(result)
			panic("Internal compiler error, something seems to have slipped through the cracks. Please file a github issue with the core dump")
		}
	}
}

func (i *interpreter) executeBlock(statements []stmt, environment *environment) {
  previous := i.environment

  // Life into new environment context
  i.environment = environment

  // Run everything in this scope
  for _, statement := range statements {
    i.execute(statement)
  }

  // Hand the environment back
  i.environment = previous
}

func (i *interpreter) execute(s stmt) interface{} {
	return s.accept(i)
}

func (i *interpreter) stringify(evaluated interface{}) string {
	if reflect.TypeOf(evaluated).String() == "float64" {
		return strconv.FormatFloat(evaluated.(float64), 'f', -1, 64)
	}

	return fmt.Sprintf("%v", evaluated)
}

func (i *interpreter) evaluate(e expr) interface{} {
	return e.accept(i)
}

func (i *interpreter) visitBlockStmt(s *blockStmt) interface{} {
  i.executeBlock(s.statements, newEnvironment(i.environment))
  return nil
}

func (i *interpreter) visitVariableStmt(s *variableStmt) interface{} {
	var value interface{}
	if s.initializer != nil {
		value = i.evaluate(s.initializer)
	}

	i.environment.define(s.name.lexeme, value)
	return nil
}

func (i *interpreter) visitExpressionStmt(s *expressionStmt) interface{} {
	i.evaluate(s.expression)
	return nil
}

func (i *interpreter) visitPrintStmt(s *printStmt) interface{} {
	value := i.evaluate(s.expression)
	fmt.Println(i.stringify(value))
	return nil
}

func (i *interpreter) visitAssignExpr(a *assignExpr) interface{} {
	value := i.evaluate(a.value)

	err := i.environment.assign(a.name, value)
	if err != nil {
		return err
	}
	return nil
}

func (i *interpreter) visitVariableExpr(e *variableExpr) interface{} {
	value, err := i.environment.get(e.name)

	if err != nil {
		return err
	}

	return value
}

func (i *interpreter) visitLiteralExpr(e *literalExpr) interface{} {
	return e.value
}

func (i *interpreter) visitGroupingExpr(e *groupingExpr) interface{} {
	return i.evaluate(e)
}

func (i *interpreter) visitBinaryExpr(e *binaryExpr) interface{} {
	left := i.evaluate(e.left)
	right := i.evaluate(e.right)

	if e.operator.variant != PLUS {
		err := i.checkBinaryNumberOperands(e.operator, left, right)

		if err != nil {
			return err
		}
	}

	switch e.operator.variant {
	case MINUS:
		return left.(float64) - right.(float64)
	case SLASH:
		// Handle divide by zero
		if right.(float64) == 0 {
			return newRuntimeError(e.operator, "error! attempted to divide by zero")
		}

		return left.(float64) / right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case PLUS:
		if reflect.TypeOf(left).String() == "string" && reflect.TypeOf(right).String() != "string" {
			//TODO(@jparr721) - We should add custom checks for various types
			// Rough conversion of right value
			rstr := fmt.Sprintf("%v", right)
			return left.(string) + rstr
		}

		if reflect.TypeOf(left).String() == "string" && reflect.TypeOf(right).String() == "string" {
			return left.(string) + right.(string)
		}

		if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
			return left.(float64) + right.(float64)
		}

		return newRuntimeError(e.operator, "Operator requires two strings or two numbers.")
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATEREQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESSEQUAL:
		return left.(float64) <= right.(float64)
	case BANGEQUAL:
		return !i.isEqual(left, right)
	case EQUALEQUAL:
		return i.isEqual(left, right)
	}

	return nil
}

func (i *interpreter) visitUnaryExpr(e *unaryExpr) interface{} {
	right := i.evaluate(e.right.(expr))

	err := i.checkInfixNumberOperand(e.operator, right)

	if err != nil {
		return err
	}

	switch e.operator.variant {
	case MINUS:
		return -right.(float64)
	case BANG:
		return !i.isTruthy(right)
	}

	return nil
}

func (i *interpreter) isEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func (i *interpreter) checkBinaryNumberOperands(operator token, left, right interface{}) *runtimeError {
	if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
		return nil
	}

	return newRuntimeError(operator, "Operands must be two numbers.")
}

func (i *interpreter) checkInfixNumberOperand(operator token, operand interface{}) *runtimeError {
	if reflect.TypeOf(operand).String() == "float64" {
		return nil
	}

	return newRuntimeError(operator, "Operans must be numbers.")
}

// only nil and false are falsy values, everything else evaluates to truthy
func (i *interpreter) isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}

	switch value.(type) {
	case bool:
		return value.(bool)
	default:
		return true
	}
}
