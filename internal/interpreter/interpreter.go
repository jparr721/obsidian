package interpreter

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/jparr721/obsidian/internal/expression"
	"github.com/jparr721/obsidian/internal/statement"
	"github.com/jparr721/obsidian/internal/tokens"
)

func stringify(evaluated interface{}) string {
	if evaluated == nil {
		return "nil"
	}

	if reflect.TypeOf(evaluated).String() == "float64" {
		return strconv.FormatFloat(evaluated.(float64), 'f', -1, 64)
	}

	return fmt.Sprintf("%v", evaluated)
}

type Interpreter struct {
	// globals are the global native objects, constants, and functions
	globals      *environment
	environment  *environment
	loopDidBreak bool
}

func NewInterpreter() *Interpreter {
	environment := NewEnvironment(nil)
	globals := NewEnvironment(nil)
	globals.define("clock", new(clockFunction))

	return &Interpreter{globals, environment, false}
}

func (i *Interpreter) Interpret(statements []statement.Statement) error {
	for _, statement := range statements {
    fmt.Println(i.environment)
		_, err := i.execute(statement)

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (i *Interpreter) executeBlock(statements []statement.Statement, environment *environment) error {
	previous := i.environment

	// Lift into new environment context
	i.environment = environment

	// Run everything in this scope
	for _, statement := range statements {
    fmt.Println(reflect.TypeOf(statement))
		if reflect.TypeOf(statement).String() == "*statement.BreakStatement" {
			i.loopDidBreak = true
			break
		}

		_, err := i.execute(statement)
		if err != nil {
			return err
		}
	}

	// Hand the environment back
	i.environment = previous
	return nil
}

func (i *Interpreter) execute(s statement.Statement) (interface{}, error) {
	return s.Accept(i)
}

func (i *Interpreter) evaluate(e expression.Expression) (interface{}, error) {
	return e.Accept(i)
}

func (i *Interpreter) VisitFunctionStatement(s *statement.FunctionStatement) (interface{}, error) {
	function := NewFunction(s)
	i.environment.define(s.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitBreakStatement(s *statement.BreakStatement) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) VisitWhileStatement(s *statement.WhileStatement) (interface{}, error) {
	for {
		if i.loopDidBreak {
			i.loopDidBreak = false
			break
		}

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

func (i *Interpreter) VisitIfStatement(s *statement.IfStatement) (interface{}, error) {
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

func (i *Interpreter) VisitBlockStatement(s *statement.BlockStatement) (interface{}, error) {
	return nil, i.executeBlock(s.Statements, NewEnvironment(i.environment))
}

func (i *Interpreter) VisitVariableStatement(s *statement.VariableStatement) (interface{}, error) {
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

func (i *Interpreter) VisitExpressionStatement(s *statement.ExpressionStatement) (interface{}, error) {
	_, err := i.evaluate(s.Expression)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Interpreter) VisitPrintStatement(s *statement.PrintStatement) (interface{}, error) {
	value, err := i.evaluate(s.Expression)

	if err != nil {
		return nil, err
	}

	fmt.Println(stringify(value))

	return nil, nil
}

func (i *Interpreter) VisitReturnStatement(s *statement.ReturnStatement) (interface{}, error) {
	if s.Value != nil {
		value, err := i.evaluate(s.Value)

		if err != nil {
			return nil, err
		}

		return nil, newReturnInterrupt(value)
	}

	return nil, newReturnInterrupt(nil)
}

func (i *Interpreter) VisitCallExpression(e *expression.CallExpression) (interface{}, error) {
	callee, err := i.evaluate(e.Callee)

	if err != nil {
		return nil, err
	}

	arguments := make([]interface{}, 0)
	for _, argument := range e.Arguments {
		a, err := i.evaluate(argument)

		if err != nil {
			return nil, err
		}

		arguments = append(arguments, a)
	}

	if reflect.TypeOf(callee).String() != "*interpreter.Function" {
		return nil, newRuntimeError(e.Paren, "Only function and class types are callable.")
	}

	function := callee.(Callable)

	if function.Arity() != len(arguments) {
		return nil, newRuntimeError(e.Paren, fmt.Sprintf("Expected %d arguments, but got %d.", function.Arity(), len(arguments)))
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitLogicalExpression(e *expression.LogicalExpression) (interface{}, error) {
	left, err := i.evaluate(e.Left)

	if err != nil {
		return nil, err
	}

	if e.Operator.Variant == tokens.TokenOr {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(e.Right)
}

func (i *Interpreter) VisitAssignExpression(e *expression.AssignExpression) (interface{}, error) {
	value, err := i.evaluate(e.Value)

	if err != nil {
		return nil, err
	}

	err = i.environment.assign(e.Name, value)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Interpreter) VisitVariableExpression(e *expression.VariableExpression) (interface{}, error) {
	value, err := i.environment.get(e.Name)

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (i *Interpreter) VisitLiteralExpression(e *expression.LiteralExpression) (interface{}, error) {
	return e.Value, nil
}

func (i *Interpreter) VisitGroupingExpression(e *expression.GroupingExpression) (interface{}, error) {
	return i.evaluate(e)
}

func (i *Interpreter) VisitBinaryExpression(e *expression.BinaryExpression) (interface{}, error) {
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

func (i *Interpreter) VisitUnaryExpression(e *expression.UnaryExpression) (interface{}, error) {
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

func (i *Interpreter) isEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func (i *Interpreter) checkBinaryNumberOperands(operator tokens.Token, left, right interface{}) *RuntimeError {
	if reflect.TypeOf(left).String() == "float64" && reflect.TypeOf(right).String() == "float64" {
		return nil
	}

	return newRuntimeError(operator, "Operands must be two numbers.")
}

func (i *Interpreter) checkInfixNumberOperand(operator tokens.Token, operand interface{}) *RuntimeError {
	if reflect.TypeOf(operand).String() == "float64" {
		return nil
	}

	return newRuntimeError(operator, "Operans must be numbers.")
}

// only nil and false are falsy values, everything else evaluates to truthy
func (i *Interpreter) isTruthy(value interface{}) bool {
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
