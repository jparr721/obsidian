package main

type stmtVisitor interface {
	visitExpressionStmt(*expressionStmt) interface{}
	visitPrintStmt(*printStmt) interface{}
	visitVariableStmt(*variableStmt) interface{}
	visitBlockStmt(*blockStmt) interface{}
	visitIfStmt(*ifStmt) interface{}
	visitWhileStmt(*whileStmt) interface{}
	visitBreakStmt(*breakStmt) interface{}
}

type stmt interface {
	accept(stmtVisitor) interface{}
}

// Expression Statement
type expressionStmt struct {
	expression expr
}

func newExpressionStmt(e expr) *expressionStmt {
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

// If Statement
type ifStmt struct {
	condition  expr
	thenBranch stmt
	elseBranch stmt
}

func newIfStmt(condition expr, thenBranch, elseBranch stmt) *ifStmt {
	return &ifStmt{condition, thenBranch, elseBranch}
}

func (i *ifStmt) accept(v stmtVisitor) interface{} {
	return v.visitIfStmt(i)
}

// While Statement
type whileStmt struct {
	condition expr
	body      stmt
}

func newWhileStmt(condition expr, body stmt) *whileStmt {
	return &whileStmt{condition, body}
}

func (w *whileStmt) accept(v stmtVisitor) interface{} {
	return v.visitWhileStmt(w)
}

// Break statement
type breakStmt struct {
	instance token
}

func newBreakStmt(instance token) *breakStmt {
	return &breakStmt{instance}
}

func (b *breakStmt) accept(v stmtVisitor) interface{} {
	return v.visitBreakStmt(b)
}
