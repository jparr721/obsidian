package main

const (
	inLoopStatement    = true
	notInLoopStatement = false
)

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

func (p *parser) parse() ([]stmt, *parseError) {
	statements := make([]stmt, 0)

	for !p.end() {
		stmt, err := p.declaration(notInLoopStatement)

		// Stop parsing - report the error
		if err != nil {
			reportParseError(err)
			return nil, err
		}

		statements = append(statements, stmt)
	}

	return statements, nil
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
func (p *parser) declaration(inLoop bool) (stmt, *parseError) {
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

	return p.statement(inLoop)
}

func (p *parser) varDeclaration() (stmt, *parseError) {
	name, err := p.consume(IDENTIFIER, "Expected variable name.")

	if err != nil {
		return nil, err
	}

	var initializer expr

	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(SEMI, "Expected ';' after variable declaration")
	variableStmt := newVariableStmt(name, initializer)
	return variableStmt, nil
}

// statement -> stmtStmt | printStmt;
func (p *parser) statement(inLoop bool) (stmt, *parseError) {
	if p.match(FOR) {
		return p.forStatement()
	}

	if p.match(IF) {
		return p.ifStatement()
	}

	if p.match(PRINT) {
		return p.printStmt()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	if p.match(BREAK) {
		return p.breakStatement(inLoop)
	}

	if p.match(OSQUIGGLE) {
		statements, err := p.block(inLoop)

		if err != nil {
			return nil, err
		}

		return newBlockStmt(statements), nil
	}

	return p.expressionStmt()
}

// for -> "for" "(" (varDecl | exprStmt | ";") expression?";" expression?";" statement ;
func (p *parser) forStatement() (stmt, *parseError) {
	p.consume(OPAREN, "Expected '(' after 'for'")
	var err *parseError

	// first clause in the for loop
	var initializer stmt
	if p.match(SEMI) {
		initializer = nil
	} else if p.match(VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStmt()
		if err != nil {
			return nil, err
		}
	}

	// Second clause in the for loop
	var condition expr
	if !p.check(SEMI) {
		condition, err = p.expression()

		if err != nil {
			return nil, err
		}

	}

	p.consume(SEMI, "Expected ';' after loop condition")

	// Third clause in the for loop
	var increment expr
	if !p.check(CPAREN) {
		increment, err = p.expression()

		if err != nil {
			return nil, err
		}

	}

	_, err = p.consume(CPAREN, "Expected ')' after for clauses")

	if err != nil {
		return nil, err
	}

	body, err := p.statement(inLoopStatement)

	if err != nil {
		return nil, err
	}

	// Append the increment to the bottom of the block if set
	if increment != nil {
		body.(*blockStmt).statements = append(body.(*blockStmt).statements, newExpressionStmt(increment))
	}

	if condition == nil {
		condition = newLiteralExpr(true)
	}

	body = newWhileStmt(condition, body)

	if initializer != nil {
		block := []stmt{
			initializer,
			body,
		}
		body = newBlockStmt(block)
	}

	return body, nil
}

// break -> "break";
func (p *parser) breakStatement(inLoop bool) (stmt, *parseError) {
	if inLoop {
		b := p.prev()
		_, err := p.consume(SEMI, "Expected ';' after break statement")

		if err != nil {
			return nil, err
		}

		return newBreakStmt(b), nil
	}

	return nil, newParseError(p.prev(), "Expected 'break' inside of while or for loop")
}

// while -> "while" "(" expression ")" statement;
func (p *parser) whileStatement() (stmt, *parseError) {
	p.consume(OPAREN, "Expected '(' after 'while'")
	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(CPAREN, "Expected ')' after condition")

	body, err := p.statement(inLoopStatement)

	if err != nil {
		return nil, err
	}

	return newWhileStmt(condition, body), nil
}

// if -> "if" "(" expression ")" statement ("else" statement)?;
func (p *parser) ifStatement() (stmt, *parseError) {
	p.consume(OPAREN, "Expected '(' after 'if'")
	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(CPAREN, "Expected ')' after if condition")

	thenBranch, err := p.statement(notInLoopStatement)

	if err != nil {
		return nil, err
	}

	var elseBranch stmt
	if p.match(ELSE) {
		elseBranch, err = p.statement(notInLoopStatement)

		if err != nil {
			return nil, err
		}
	}

	return newIfStmt(condition, thenBranch, elseBranch), nil
}

// stmtStmt -> stmtession ";";
func (p *parser) expressionStmt() (stmt, *parseError) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(SEMI, "Expected a ';' after expression.")
	if err != nil {
		return nil, err
	}

	expressionStmt := newExpressionStmt(value)

	return expressionStmt, nil
}

// printStmt -> "print" stmtession ";";
func (p *parser) printStmt() (stmt, *parseError) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(SEMI, "Expected a ';' after value.")
	if err != nil {
		return nil, err
	}

	printStmt := newPrintStmt(value)

	return printStmt, nil
}

// block -> "{" declaration* "}";
func (p *parser) block(inLoop bool) ([]stmt, *parseError) {
	statements := make([]stmt, 0)

	for !p.check(CSQUIGGLE) && !p.end() {
		statement, err := p.declaration(inLoop)

		if err != nil {
			return nil, err
		}

		statements = append(statements, statement)
	}

	_, err := p.consume(CSQUIGGLE, "Expected '}' after block statement")

	if err != nil {
		return nil, err
	}

	return statements, nil
}

// assignment -> IDENTIFIER "=" assignment | equality;
func (p *parser) assignment() (expr, *parseError) {
	expression, err := p.or()
	if err != nil {
		return nil, err
	}

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
func (p *parser) expression() (expr, *parseError) {
	return p.assignment()
}

// logical or -> logical and ("or" logical and)*
func (p *parser) or() (expr, *parseError) {
	expression, err := p.and()

	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.prev()
		right, err := p.and()

		if err != nil {
			return nil, err
		}

		expression = newLogicalExpr(expression, right, operator)
	}

	return expression, nil
}

func (p *parser) and() (expr, *parseError) {
	expression, err := p.equality()

	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.prev()
		right, err := p.equality()

		if err != nil {
			return nil, err
		}

		expression = newLogicalExpr(expression, right, operator)
	}

	return expression, nil
}

// equality -> comparison ( ( "!=" | "==" ) comparison )*;
func (p *parser) equality() (expr, *parseError) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANGEQUAL, EQUALEQUAL) {
		operator := p.prev()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = newBinaryExpr(expr, right, operator)
	}

	return expr, nil
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )*;
func (p *parser) comparison() (expr, *parseError) {
	expr, err := p.term()

	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATEREQUAL, LESS, LESSEQUAL) {
		operator := p.prev()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = newBinaryExpr(expr, right, operator)
	}

	return expr, nil
}

// term -> ( term ( "+" | "-" ) )*;
func (p *parser) term() (expr, *parseError) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = newBinaryExpr(expr, right, operator)
	}

	return expr, nil
}

// factor -> ( term ( "/" | "*" ) term )*;
func (p *parser) factor() (expr, *parseError) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.prev()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		expr = newBinaryExpr(expr, right, operator)
	}

	return expr, nil
}

// unary -> ( "!" | "-" ) unary | primary;
func (p *parser) unary() (expr, *parseError) {
	if p.match(BANG, MINUS) {
		operator := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return newUnaryExpr(operator, right), nil
	}

	return p.primary()
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
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(CPAREN, "expected ')' after expression.")
		if err != nil {
			return nil, err
		}

		return newGroupingExpr(expr), nil
	}

	if p.match(IDENTIFIER) {
		return newVariableExpr(p.prev()), nil
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
