package main

type parser struct {
	tokens  []token
	current int
}

func newParser(tokens []token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *parser) parse() []stmt {
	statements := make([]stmt, 0)

	for !p.end() {
		stmt, err := p.statement()

		// Stop parsing - report the error
		if err != nil {
			reportParseError(err)
			break
		}

		statements = append(statements, stmt)
	}

	return statements
}

func (p *parser) check(tType tokenType) bool {
	if p.end() {
		return false
	}

	return p.peek().variant == tType
}

func (p *parser) end() bool {
	return p.peek().variant == EOF
}

func (p *parser) next() token {
	if !p.end() {
		p.current++
	}

	return p.prev()
}

func (p *parser) peek() token {
	return p.tokens[p.current]
}

func (p *parser) prev() token {
	return p.tokens[p.current-1]
}

func (p *parser) match(tTypes ...tokenType) bool {
	for _, t := range tTypes {
		if p.check(t) {
			p.next()
			return true
		}
	}

	return false
}

// declaration -> varDecl | statement;
func (p *parser) declaration() (stmt, *parseError) {
	if p.match(VAR) {
		value, err := p.varDeclaration()

		if err != nil {
			// unwind errors to next valid expression.
			p.synchronize()
			// Bubble up error regardless.
			return nil, err
		}

		return value, nil
	}

	return p.statement()
}

func (p *parser) varDeclaration() (stmt, *parseError) {
	name, err := p.consume(IDENTIFIER, "Expected variable name.")

	if err != nil {
		return nil, err
	}

	var initializer expr

	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMI, "Expected ';' after variable declaration")
	variableStmt := newVariableStmt(name, initializer)
	return variableStmt, nil
}

// statement -> stmtStmt | printStmt;
func (p *parser) statement() (stmt, *parseError) {
	if p.match(PRINT) {
		return p.printStmt()
	}

	if p.match(OSQUIGGLE) {
		statements, err := p.block()

		if err != nil {
			return nil, err
		}

		return newBlockStmt(statements), nil
	}

	return p.expressionStmt()
}

// stmtStmt -> stmtession ";";
func (p *parser) expressionStmt() (stmt, *parseError) {
	value := p.expression()

	_, err := p.consume(SEMI, "Expected a ';' after expression.")
	if err != nil {
		return nil, err
	}

	expressionStmt := newExpressionStatement(value)

	return expressionStmt, nil
}

// printStmt -> "print" stmtession ";";
func (p *parser) printStmt() (stmt, *parseError) {
	value := p.expression()

	_, err := p.consume(SEMI, "Expected a ';' after value.")
	if err != nil {
		return nil, err
	}

	printStmt := newPrintStmt(value)

	return printStmt, nil
}

// block -> "{" declaration* "}";
func (p *parser) block() ([]stmt, *parseError) {
	statements := make([]stmt, 0)

	for !p.check(CSQUIGGLE) && !p.end() {
		declaration, err := p.declaration()

		if err != nil {
			return nil, err
		}

		statements = append(statements, declaration)
	}

	_, err := p.consume(CSQUIGGLE, "Expected '}' after block statement")

	if err != nil {
		return nil, err
	}

	return statements, nil
}

// assignment -> IDENTIFIER "=" assignment | equality;
func (p *parser) assignment() (expr, *parseError) {
	expression := p.equality()

	if p.match(EQUAL) {
		equals := p.prev()
		value, err := p.assignment()

		if err != nil {
			return nil, err
		}

		switch expression.(type) {
		case *variableExpr:
			name := expression.(*variableExpr).name
			return newAssignExpr(name, value), nil
		default:
			return nil, newParseError(equals, "Invalid assignment target")
		}
	}

	return expression, nil
}

// represents an expression statement
func (p *parser) expression() expr {
	return p.equality()
}

// equality -> comparison ( ( "!=" | "==" ) comparison )*;
func (p *parser) equality() expr {
	expr := p.comparison()

	for p.match(BANGEQUAL, EQUALEQUAL) {
		operator := p.prev()
		right := p.comparison()
		expr = newBinaryExpr(expr, right, operator)
	}

	return expr
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )*;
func (p *parser) comparison() expr {
	expr := p.term()

	for p.match(GREATER, GREATEREQUAL, LESS, LESSEQUAL) {
		operator := p.prev()
		right := p.term()
		expr = newBinaryExpr(expr, right, operator)
	}

	return expr
}

// term -> ( term ( "+" | "-" ) )*;
func (p *parser) term() expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.prev()
		right := p.unary()
		expr = newBinaryExpr(expr, right, operator)
	}

	return expr
}

// factor -> ( term ( "/" | "*" ) term )*;
func (p *parser) factor() expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.prev()
		right := p.unary()
		expr = newBinaryExpr(expr, right, operator)
	}

	return expr
}

// unary -> ( "!" | "-" ) unary | primary;
func (p *parser) unary() expr {
	if p.match(BANG, MINUS) {
		operator := p.prev()
		right := p.unary()
		return newUnaryExpr(operator, right)
	}

	parsed, err := p.primary()

	if err != nil {
		return nil
	}

	return parsed
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")";
func (p *parser) primary() (expr, *parseError) {
	if p.match(FALSE) {
		return newLiteralExpr(false), nil
	}

	if p.match(TRUE) {
		return newLiteralExpr(true), nil
	}

	if p.match(NIL) {
		return newLiteralExpr(nil), nil
	}

	if p.match(NUMBER, STRING) {
		return newLiteralExpr(p.prev().literal), nil
	}

	if p.match(OPAREN) {
		expr := p.expression()
		_, err := p.consume(CPAREN, "expected ')' after expression.")

		if err != nil {
			return nil, err
		}

		return newGroupingExpr(expr), nil
	}

	return nil, newParseError(p.peek(), "expected expression")
}

func (p *parser) consume(t tokenType, errorMsg string) (token, *parseError) {
	if p.check(t) {
		return p.next(), nil
	}

	return token{}, newParseError(p.peek(), errorMsg)
}

// Synchronize when we've caught an error to the next valid token
func (p *parser) synchronize() {
	p.next()

	for !p.end() {
		if p.prev().variant == SEMI {
			return
		}

		switch p.peek().variant {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.next()
	}
}
