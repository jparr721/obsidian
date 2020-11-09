package expression

import "github.com/jparr721/obsidian/internal/tokens"

type Visitor interface {
	VisitBinaryExpression(*BinaryExpression) (interface{}, error)
	VisitGroupingExpression(*GroupingExpression) (interface{}, error)
	VisitLiteralExpression(*LiteralExpression) (interface{}, error)
	VisitUnaryExpression(*UnaryExpression) (interface{}, error)
	VisitVariableExpression(*VariableExpression) (interface{}, error)
	VisitAssignExpression(*AssignExpression) (interface{}, error)
	VisitLogicalExpression(*LogicalExpression) (interface{}, error)
}

type Expression interface {
	Accept(Visitor) (interface{}, error)
}

type BinaryExpression struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

func (b *BinaryExpression) Accept(v Visitor) (interface{}, error) {
	return v.VisitBinaryExpression(b)
}

func NewBinaryExpression(left, right Expression, operator tokens.Token) *BinaryExpression {
	return &BinaryExpression{
		left,
		operator,
		right,
	}
}

type GroupingExpression struct {
	Expression interface{}
}

func (g *GroupingExpression) Accept(v Visitor) (interface{}, error) {
	return v.VisitGroupingExpression(g)
}

func NewGroupingExpression(Expression interface{}) *GroupingExpression {
	return &GroupingExpression{
		Expression,
	}
}

type LiteralExpression struct {
	Value interface{}
}

func (l *LiteralExpression) Accept(v Visitor) (interface{}, error) {
	return v.VisitLiteralExpression(l)
}

func NewLiteralExpression(value interface{}) *LiteralExpression {
	return &LiteralExpression{
		value,
	}
}

type UnaryExpression struct {
	Operator tokens.Token
	Right    interface{}
}

func (u *UnaryExpression) Accept(v Visitor) (interface{}, error) {
	return v.VisitUnaryExpression(u)
}

func NewUnaryExpression(operator tokens.Token, right interface{}) *UnaryExpression {
	return &UnaryExpression{
		operator,
		right,
	}
}

type VariableExpression struct {
	Name tokens.Token
}

func (ve *VariableExpression) Accept(v Visitor) (interface{}, error) {
	return v.VisitVariableExpression(ve)
}

func NewVariableExpression(name tokens.Token) *VariableExpression {
	return &VariableExpression{name}
}

type AssignExpression struct {
	Name  tokens.Token
	Value Expression
}

func (a *AssignExpression) Accept(v Visitor) (interface{}, error) {
	return v.VisitAssignExpression(a)
}

func NewAssignExpression(name tokens.Token, value Expression) *AssignExpression {
	return &AssignExpression{name, value}
}

type LogicalExpression struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

func (l *LogicalExpression) Accept(v Visitor) (interface{}, error) {
	return v.VisitLogicalExpression(l)
}

func NewLogicalExpression(left, right Expression, operator tokens.Token) *LogicalExpression {
	return &LogicalExpression{left, operator, right}
}
