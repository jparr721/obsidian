package main

type expressionVisitor interface {
	visitBinaryExpr(*binaryExpr) interface{}
	visitGroupingExpr(*groupingExpr) interface{}
	visitLiteralExpr(*literalExpr) interface{}
	visitUnaryExpr(*unaryExpr) interface{}
	visitVariableExpr(*variableExpr) interface{}
	visitAssignExpr(*assignExpr) interface{}
	visitLogicalExpr(*logicalExpr) interface{}
}

type expr interface {
	accept(expressionVisitor) interface{}
}

type binaryExpr struct {
	left     expr
	operator token
	right    expr
}

func (b *binaryExpr) accept(v expressionVisitor) interface{} {
	return v.visitBinaryExpr(b)
}

func newBinaryExpr(left, right expr, operator token) *binaryExpr {
	return &binaryExpr{
		left,
		operator,
		right,
	}
}

type groupingExpr struct {
	expression interface{}
}

func (g *groupingExpr) accept(v expressionVisitor) interface{} {
	return v.visitGroupingExpr(g)
}

func newGroupingExpr(expression interface{}) *groupingExpr {
	return &groupingExpr{
		expression,
	}
}

type literalExpr struct {
	value interface{}
}

func (l *literalExpr) accept(v expressionVisitor) interface{} {
	return v.visitLiteralExpr(l)
}

func newLiteralExpr(value interface{}) *literalExpr {
	return &literalExpr{
		value,
	}
}

type unaryExpr struct {
	operator token
	right    interface{}
}

func (u *unaryExpr) accept(v expressionVisitor) interface{} {
	return v.visitUnaryExpr(u)
}

func newUnaryExpr(operator token, right interface{}) *unaryExpr {
	return &unaryExpr{
		operator,
		right,
	}
}

type variableExpr struct {
	name token
}

func (ve *variableExpr) accept(v expressionVisitor) interface{} {
	return v.visitVariableExpr(ve)
}

func newVariableExpr(name token) *variableExpr {
	return &variableExpr{name}
}

type assignExpr struct {
	name  token
	value expr
}

func (a *assignExpr) accept(v expressionVisitor) interface{} {
	return v.visitAssignExpr(a)
}

func newAssignExpr(name token, value expr) *assignExpr {
	return &assignExpr{name, value}
}

type logicalExpr struct {
	left     expr
	operator token
	right    expr
}

func (l *logicalExpr) accept(v expressionVisitor) interface{} {
	return v.visitLogicalExpr(l)
}

func newLogicalExpr(left, right expr, operator token) *logicalExpr {
	return &logicalExpr{left, operator, right}
}
