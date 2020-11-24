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
	VisitFunctionStatement(*FunctionStatement) (interface{}, error)
	VisitReturnStatement(*ReturnStatement) (interface{}, error)
}

// Statement represents
type Statement interface {
	Accept(Visitor) (interface{}, error)
}

// ExpressionStatement represents an expression statement
type ExpressionStatement struct {
	Expression expression.Expression
}

// NewExpressionStatement creates a new ExpressionStatement
func NewExpressionStatement(e expression.Expression) *ExpressionStatement {
	return &ExpressionStatement{e}
}

// Accept is the method which invokes this type's functionality
func (e *ExpressionStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitExpressionStatement(e)
}

// PrintStatement represents a print statement
type PrintStatement struct {
	Expression expression.Expression
}

// NewPrintStatement creates a new PrintStatement
func NewPrintStatement(expression expression.Expression) *PrintStatement {
	return &PrintStatement{expression}
}

// Accept is the method which invokes this type's functionality
func (p *PrintStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitPrintStatement(p)
}

// VariableStatement represents a variable statement
type VariableStatement struct {
	Name        tokens.Token
	Initializer expression.Expression
}

// NewVariableStatement creates a new VariableStatement
func NewVariableStatement(name tokens.Token, initializer expression.Expression) *VariableStatement {
	return &VariableStatement{name, initializer}
}

// Accept is the method which invokes this type's functionality
func (vs *VariableStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitVariableStatement(vs)
}

// BlockStatement represents a block statement
type BlockStatement struct {
	Statements []Statement
}

// NewBlockStatement creates a new BlockStatement
func NewBlockStatement(statements []Statement) *BlockStatement {
	return &BlockStatement{statements}
}

// Accept is the method which invokes this type's functionality
func (b *BlockStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitBlockStatement(b)
}

// IfStatement represents a if statement
type IfStatement struct {
	Condition  expression.Expression
	ThenBranch Statement
	ElseBranch Statement
}

// NewIfStatement creates a new IfStatement
func NewIfStatement(condition expression.Expression, thenBranch, elseBranch Statement) *IfStatement {
	return &IfStatement{condition, thenBranch, elseBranch}
}

// Accept is the method which invokes this type's functionality
func (i *IfStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitIfStatement(i)
}

// WhileStatement represents a while statement
type WhileStatement struct {
	Condition expression.Expression
	Body      Statement
}

// NewWhileStatement creates a new WhileStatement
func NewWhileStatement(condition expression.Expression, body Statement) *WhileStatement {
	return &WhileStatement{condition, body}
}

// Accept is the method which invokes this type's functionality
func (w *WhileStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitWhileStatement(w)
}

// BreakStatement represents a break statement
type BreakStatement struct {
	Instance tokens.Token
}

// NewBreakStatement creates a new BreakStatement
func NewBreakStatement(instance tokens.Token) *BreakStatement {
	return &BreakStatement{instance}
}

// Accept is the method which invokes this type's functionality
func (b *BreakStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitBreakStatement(b)
}

// FunctionStatement represents a function statement
type FunctionStatement struct {
	Name      tokens.Token
	Arguments []tokens.Token
	Body      []Statement
}

// NewFunctionStatement creates a new FunctionStatement
func NewFunctionStatement(name tokens.Token, arguments []tokens.Token, body []Statement) *FunctionStatement {
	return &FunctionStatement{name, arguments, body}
}

// Accept is the method which invokes this type's functionality
func (f *FunctionStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitFunctionStatement(f)
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Keyword tokens.Token
	Value   expression.Expression
}

// NewReturnStatement creates a new ReturnStatement
func NewReturnStatement(keyword tokens.Token, value expression.Expression) *ReturnStatement {
	return &ReturnStatement{keyword, value}
}

// Accept is the method which invokes this type's functionality
func (r *ReturnStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitReturnStatement(r)
}
