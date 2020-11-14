package interpreter

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/jparr721/obsidian/internal/expression"
	"github.com/jparr721/obsidian/internal/statement"
	"github.com/jparr721/obsidian/internal/tokens"
)

//TODO(@jparr721) - This _really_ needs to use Visitors with error propagation.

type interpreter struct {
	environment *environment
}

func NewInterpreter() *interpreter {
	return &interpreter{NewEnvironment(nil)}
}

func (i *interpreter) Interpret(statements []statement.Statement) error {
	for _, statement := range statements {
		_, err := i.execute(statement)

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *interpreter) executeBlock(statements []statement.Statement, environment *environment) error {
	previous := i.environment

	// Lift into new environment context
	i.environment = environment

	// Run everything in this scope
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			return err
		}
	}

	// Hand the environment back
	i.environment = previous
	return nil
}

func (i *interpreter) execute(s statement.Statement) (interface{}, error) {
	return s.Accept(i)
}

func (i *interpreter) stringify(evaluated interface{}) string {
	if evaluated == nil {
		return "nil"
	}

	if reflect.TypeOf(evaluated).String() == "float64" {
		return strconv.FormatFloat(evaluated.(float64), 'f', -1, 64)
	}

	return fmt.Sprintf("%v", evaluated)
}

func (i *interpreter) evaluate(e expression.Expression) (interface{}, error) {
	return e.Accept(i)
}

func (i *interpreter) VisitBreakStatement(s *statement.BreakStatement) (interface{}, error) {
	return nil, newRuntimeError(s.Instance, "'break' found outside of loop statement, please file an issue")
}

func (i *interpreter) VisitWhileStatement(s *statement.WhileStatement) (interface{}, error) {
	for {
		cond, err := i.evaluate(s.Condition)

		if err != nil {
			return nil, err
		}

		if !i.isTruthy(cond) {
			break
		}

		_, err = i.execute(s.Body)

		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *interpreter) VisitIfStatement(s *statement.IfStatement) (interface{}, error) {
	cond, err := i.evaluate(s.Condition)

	if err != nil {
		return nil, err
	}

	if i.isTruthy(cond) {
		i.execute(s.ThenBranch)
	} else if s.ElseBranch != nil {
		i.execute(s.ElseBranch)
	}

	return nil, nil
}

func (i *interpreter) VisitBlockStatement(s *statement.BlockStatement) (interface{}, error) {
	return nil, i.executeBlock(s.Statements, NewEnvironment(i.environment))
}

func (i *interpreter) VisitVariableStatement(s *statement.VariableStatement) (interface{}, error) {
	var value interface{} = nil
	var err error

	if s.Initializer != nil {
		value, err = i.evaluate(s.Initializer)

		if err != nil {
			return nil, err
		}

	}

	i.environment.define(s.Name.Lexeme, value)

	return nil, nil
}

func (i *interpreter) VisitExpressionStatement(s *statement.ExpressionStatement) (interface{}, error) {
	_, err := i.evaluate(s.Expression)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *interpreter) VisitPrintStatement(s *statement.PrintStatement) (interface{}, error) {
	value, err := i.evaluate(s.Expression)

	if err != nil {
		return nil, err
	}

	fmt.Println(i.stringify(value))

	return nil, nil
}

func (i *interpreter) VisitLogicalExpression(a *expression.LogicalExpression) (interface{}, error) {
	left, err := i.evaluate(a.Left)

	if err != nil {
		return nil, err
	}

	if a.Operator.Variant == tokens.TokenOr {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(a.Right)
}

func (i *interpreter) VisitAssignExpression(a *expression.AssignExpression) (interface{}, error) {
	value, err := i.evaluate(a.Value)

	if err != nil {
		return nil, err
	}

	err = i.environment.assign(a.Name, value)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *interpreter) VisitVariableExpression(e *expression.VariableExpression) (interface{}, error) {
	value, err := i.environment.get(e.Name)

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (i *interpreter) VisitLiteralExpression(e *expression.LiteralExpression) (interface{}, error) {
	return e.Value, nil
}

func (i *interpreter) VisitGroupingExpression(e *expression.GroupingExpression) (interface{}, error) {
	return i.evaluate(e)
}

func (i *interpreter) VisitBinaryExpression(e *expression.BinaryExpression) (interface{}, error) {
	left, err := i.evaluate(e.Left)

	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(e.Right)

	if err != nil {
		return nil, err
	}

	if e.Operator.Variant != tokens.TokenPlus {
		err := i.checkBinaryNumberOperands(e.Operator, left, right)

		if err != nil {
			return nil, err
		}
	}

	switch e.Operator.Variant {
	case tokens.TokenMinus:
		return left.(float64) - right.(float64), nil
	case tokens.TokenSlash:
		// Handle divide by zero
		if right.(float64) == 0 {
			return nil, newRuntimeError(e.Operator, "error! attempted to divide by zero")
		}

		return left.(float64) / right.(float64), nil
	case tokens.TokenStar:
		return left.(float64) * right.(float64), nil
	case tokens.TokenPlus:
		if reflect.TypeOf(left).String() == "string" && reflect.TypeOf(right).String() != "string" {
			//TODO(@jparr721) - We should add custom checks for various types
			// Rough conversion of right value
			rstr := fmt.Sprintf("%v", right)
			return left.(string) + rstr, nil
		}

		if reflect.TypeOf(left).String() == "string" && reflect.TypeOf(right).String() == "string" {
			return left.(string) + right.(string), nil
		}

		if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
			return left.(float64) + right.(float64), nil
		}

		return nil, newRuntimeError(e.Operator, "Operator requires two strings or two numbers.")
	case tokens.TokenGreater:
		return left.(float64) > right.(float64), nil
	case tokens.TokenGreaterEqual:
		return left.(float64) >= right.(float64), nil
	case tokens.TokenLess:
		return left.(float64) < right.(float64), nil
	case tokens.TokenLessEqual:
		return left.(float64) <= right.(float64), nil
	case tokens.TokenBangEqual:
		return !i.isEqual(left, right), nil
	case tokens.TokenEqualEqual:
		return i.isEqual(left, right), nil
	}

	// unreachable
	return nil, nil
}

func (i *interpreter) VisitUnaryExpression(e *expression.UnaryExpression) (interface{}, error) {
	right, err := i.evaluate(e.Right.(expression.Expression))

	if err != nil {
		return nil, err
	}

	err = i.checkInfixNumberOperand(e.Operator, right)

	if err != nil {
		return nil, err
	}

	switch e.Operator.Variant {
	case tokens.TokenMinus:
		return -right.(float64), nil
	case tokens.TokenBang:
		return !i.isTruthy(right), nil
	}

	// unreachable
	return nil, nil
}

func (i *interpreter) isEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func (i *interpreter) checkBinaryNumberOperands(operator tokens.Token, left, right interface{}) *RuntimeError {
	if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
		return nil
	}

	return newRuntimeError(operator, "Operands must be two numbers.")
}

func (i *interpreter) checkInfixNumberOperand(operator tokens.Token, operand interface{}) *RuntimeError {
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
