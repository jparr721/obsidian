package main

import (
	"fmt"
	"reflect"
	"strconv"
)

type interpreter struct{}

func (i *interpreter) interpret(e expr) {
	valueOrError := i.evaluate(e)

	switch valueOrError.(type) {
	case *runtimeError:
		reportRuntimeError(valueOrError.(*runtimeError))
	case error:
		// TODO(jparr721) - Do a core-dump for finding out the cause of the bug if this ever somehow reaches this line of code
		panic("Internal compiler error, something seems to have slipped through the cracks. Please file a github issue.")
	default:
		fmt.Println(i.stringify(valueOrError))
	}
}

func (i *interpreter) stringify(evaluated interface{}) string {
	if reflect.TypeOf(evaluated).String() == "float64" {
		return strconv.FormatFloat(evaluated.(float64), 'f', -1, 64)
	}

	return fmt.Sprintf("%v", evaluated)
}

func (i *interpreter) visitLiteralExpr(e *literalExpr) interface{} {
	return e.value
}

func (i *interpreter) visitGroupingExpr(e *groupingExpr) interface{} {
	return i.evaluate(e)
}

func (i *interpreter) evaluate(e expr) interface{} {
	return e.accept(i)
}

func (i *interpreter) visitBinaryExpr(e *binaryExpr) interface{} {
	left := i.evaluate(e.left)
	right := i.evaluate(e.right)

	err := i.checkBinaryNumberOperands(e.operator, left, right)

	if err != nil {
		return err
	}

	switch e.operator.variant {
	case MINUS:
		return left.(float64) - right.(float64)
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case PLUS:
		if reflect.TypeOf(left).String() == "string" && reflect.TypeOf(right).String() == "string" {
			return left.(string) + right.(string)
		}

		if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
			return left.(float64) + right.(float64)
		}
		fmt.Println(fmt.Errorf("Operator '%s' requires two strings or two numbers", "+"))
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

	return newRuntimeError(operator, "Operands must be two numbers or two strings.")
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
