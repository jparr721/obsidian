package main

// // AstPrinter is a printer module for visually representing the AST
// type AstPrinter struct{}

// // newAstPrinter returns a nee pointer to the ast printer module
// func newAstPrinter() *AstPrinter {
// 	return &AstPrinter{}
// }

// // Print prints the existing AST from the base expression.Expression statement
// func (a *AstPrinter) printAst(e expression.Expression) string {
// 	if e == nil {
// 		panic("no parseable expression.Expressionession found")
// 	}
// 	return e.Accept(a).(string)
// }

// func (a *AstPrinter) parenthesize(name string, exprs ...interface{}) string {
// 	builder := strings.Builder{}

// 	builder.WriteString("(")
// 	builder.WriteString(name)
// 	for _, expr := range exprs {
// 		builder.WriteString(" ")
// 		builder.WriteString(fmt.Sprintf("%v", expr.(expression.Expression).Accept(a)))
// 	}
// 	builder.WriteString(")")

// 	return builder.String()
// }

// func (a *AstPrinter) visitBinaryExpression(b *expression.BinaryExpression) interface{} {
// 	return a.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
// }

// func (a *AstPrinter) visitGroupingExpression(g *expression.GroupingExpression) interface{} {
// 	return a.parenthesize("group", g.Expression)
// }

// func (a *AstPrinter) visitLiteralExpression(l *expression.LiteralExpression) interface{} {
// 	if l.Value == nil {
// 		return "nil"
// 	}
// 	return l.Value
// }

// func (a *AstPrinter) visitUnaryExpression(u *expression.UnaryExpression) interface{} {
// 	return a.parenthesize(u.Operator.Lexeme, u.Right)
// }

// func (a *AstPrinter) visitVariableExpression(v *expression.VariableExpression) interface{} {
// 	return a.parenthesize(v.Name.Lexeme)
// }

// func (a *AstPrinter) visitAssignExpression(ae *expression.AssignExpression) interface{} {
// 	return a.parenthesize(ae.Name.Lexeme, ae.Value)
// }

// func (a *AstPrinter) visitLogicalExpression(l *expression.LogicalExpression) interface{} {
// 	return a.parenthesize(l.Operator.Lexeme, l.Left, l.Right)
// }
