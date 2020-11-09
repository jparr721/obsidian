package statement

import (
	"github.com/jparr721/obsidian/internal/expression"
	"github.com/jparr721/obsidian/internal/tokens"
)

type Visitor interface {
	VisitExpressionStatement(*ExpressionStatement) (interface{}, error)
	VisitPrintStatement(*PrintStatement) (interface{}, error)
	VisitVariableStatement(*VariableStatement) (interface{}, error)
	VisitBlockStatement(*BlockStatement) (interface{}, error)
	VisitIfStatement(*IfStatement) (interface{}, error)
	VisitWhileStatement(*WhileStatement) (interface{}, error)
	VisitBreakStatement(*BreakStatement) (interface{}, error)
}

// Statement represents
type Statement interface {
	Accept(Visitor) (interface{}, error)
}

// ExpressionStatement is a statement that represents an expression
type ExpressionStatement struct {
	Expression expression.Expression
}

// NewExpressionStatement
func NewExpressionStatement(e expression.Expression) *ExpressionStatement {
	return &ExpressionStatement{e}
}

func (e *ExpressionStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitExpressionStatement(e)
}

// Print Statement
type PrintStatement struct {
	Expression expression.Expression
}

func NewPrintStatement(expression expression.Expression) *PrintStatement {
	return &PrintStatement{expression}
}

func (p *PrintStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitPrintStatement(p)
}

// Variable Statement
type VariableStatement struct {
	Name        tokens.Token
	Initializer expression.Expression
}

func NewVariableStatement(name tokens.Token, initializer expression.Expression) *VariableStatement {
	return &VariableStatement{name, initializer}
}

func (vs *VariableStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitVariableStatement(vs)
}

// Block Statement
type BlockStatement struct {
	Statements []Statement
}

func NewBlockStatement(statements []Statement) *BlockStatement {
	return &BlockStatement{statements}
}

func (b *BlockStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitBlockStatement(b)
}

// If Statement
type IfStatement struct {
	Condition  expression.Expression
	ThenBranch Statement
	ElseBranch Statement
}

func NewIfStatement(condition expression.Expression, thenBranch, elseBranch Statement) *IfStatement {
	return &IfStatement{condition, thenBranch, elseBranch}
}

func (i *IfStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitIfStatement(i)
}

// While Statement
type WhileStatement struct {
	Condition expression.Expression
	Body      Statement
}

func NewWhileStatement(condition expression.Expression, body Statement) *WhileStatement {
	return &WhileStatement{condition, body}
}

func (w *WhileStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitWhileStatement(w)
}

// Break statement
type BreakStatement struct {
	Instance tokens.Token
}

func NewBreakStatement(instance tokens.Token) *BreakStatement {
	return &BreakStatement{instance}
}

func (b *BreakStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitBreakStatement(b)
}
