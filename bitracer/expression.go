package main

type visitor interface {
	visitbinaryExpr(*binaryExpr) interface{}
	visitgroupingExpr(*groupingExpr) interface{}
	visitliteralExpr(*literalExpr) interface{}
	visitunaryExpr(*unaryExpr) interface{}
}

// Expr repsents an interface underneath each expression statement
type expr interface {
	accept(visitor) interface{}
}

// binaryExpr is a recursive data structure representing a syntax tree
type binaryExpr struct {
	Left     expr
	operator token
	Right    expr
}

func (b *binaryExpr) accept(v visitor) interface{} {
	return v.visitbinaryExpr(b)
}

// newbinaryExpr creates a new binary expression given the parameters
func newBinaryExpr(left, right expr, operator token) *binaryExpr {
	return &binaryExpr{
		left,
		operator,
		right,
	}
}

// groupingExpr is a recursive data structure representing a syntax tree
type groupingExpr struct {
	expression interface{}
}

func (g *groupingExpr) accept(v visitor) interface{} {
	return v.visitgroupingExpr(g)
}

// newgroupingExpr creates a new grouping expression given the parameters
func newGroupingExpr(expression interface{}) *groupingExpr {
	return &groupingExpr{
		expression,
	}
}

// literalExpr is a recursive data structure representing a syntax tree
type literalExpr struct {
	value interface{}
}

func (l *literalExpr) accept(v visitor) interface{} {
	return v.visitliteralExpr(l)
}

// newliteralExpr creates a new literal value given the parameters
func newLiteralExpr(value interface{}) *literalExpr {
	return &literalExpr{
		value,
	}
}

// unaryExpr is a recursive data structure representing a syntax tree
type unaryExpr struct {
	operator token
	right    interface{}
}

func (u *unaryExpr) accept(v visitor) interface{} {
	return v.visitunaryExpr(u)
}

// newunaryExpr creates a new unary value given the parameters
func newUnaryExpr(operator token, right interface{}) *unaryExpr {
	return &unaryExpr{
		operator,
		right,
	}
}
