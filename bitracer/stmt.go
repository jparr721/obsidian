package main

type stmtVisitor interface {
	visitExpressionStmt(*expressionStmt) interface{}
	visitPrintStmt(*printStmt) interface{}
}

type stmt interface {
	accept(stmtVisitor) interface{}
}

type expressionStmt struct{
	expression expr
}

func newExpressionStatement(e expr) *expressionStmt {
	return &expressionStmt{e}
}

func (e *expressionStmt) accept(v stmtVisitor) interface{} {
	return v.visitExpressionStmt(e)
}

type printStmt struct {
	expression expr
}

func newPrintStmt(expression expr) *printStmt {
	return &printStmt{expression}
}

func (p *printStmt) accept(v stmtVisitor) interface{} {
	return v.visitPrintStmt(p)
}
