package main

type stmtVisitor interface {
	visitExpressionStmt(*expressionStmt) interface{}
	visitPrintStmt(*printStmt) interface{}
	visitVariableStmt(*variableStmt) interface{}
	visitBlockStmt(*blockStmt) interface{}
}

type stmt interface {
	accept(stmtVisitor) interface{}
}

// Expression Statement
type expressionStmt struct {
	expression expr
}

func newExpressionStatement(e expr) *expressionStmt {
	return &expressionStmt{e}
}

func (e *expressionStmt) accept(v stmtVisitor) interface{} {
	return v.visitExpressionStmt(e)
}

// Print Statement
type printStmt struct {
	expression expr
}

func newPrintStmt(expression expr) *printStmt {
	return &printStmt{expression}
}

func (p *printStmt) accept(v stmtVisitor) interface{} {
	return v.visitPrintStmt(p)
}

// Variable Statement
type variableStmt struct {
	name        token
	initializer expr
}

func newVariableStmt(name token, initializer expr) *variableStmt {
	return &variableStmt{name, initializer}
}

func (vs *variableStmt) accept(v stmtVisitor) interface{} {
	return v.visitVariableStmt(vs)
}

// Block Statement
type blockStmt struct {
	statements []stmt
}

func newBlockStmt(statements []stmt) *blockStmt {
	return &blockStmt{statements}
}

func (b *blockStmt) accept(v stmtVisitor) interface{} {
	return v.visitBlockStmt(b)
}
