package main

import (
	"fmt"
	"strings"
)

// astPrinter is a printer module for visually representing the AST
type astPrinter struct{}

// newastPrinter returns a nee pointer to the ast printer module
func newAstPrinter() *astPrinter {
	return &astPrinter{}
}

// Print prints the existing AST from the base expr statement
func (a *astPrinter) printAst(expr expr) string {
	return expr.accept(a).(string)
}

func (a *astPrinter) parenthesize(name string, exprs ...interface{}) string {
	builder := strings.Builder{}

	builder.WriteString("(")
	builder.WriteString(name)
	for _, ex := range exprs {
		builder.WriteString(" ")
		builder.WriteString(fmt.Sprintf("%v", ex.(expr).accept(a)))
	}
	builder.WriteString(")")

	return builder.String()
}

func (a *astPrinter) visitbinaryExpr(b *binaryExpr) interface{} {
	return a.parenthesize(b.operator.lexeme, b.Left, b.Right)
}

func (a *astPrinter) visitgroupingExpr(g *groupingExpr) interface{} {
	return a.parenthesize("group", g.expression)
}

func (a *astPrinter) visitliteralExpr(l *literalExpr) interface{} {
	if l.value == nil {
		return "nil"
	}
	return l.value
}

func (a *astPrinter) visitunaryExpr(u *unaryExpr) interface{} {
	return a.parenthesize(u.operator.lexeme, u.right)
}
