package parser

import (
	"github.com/jparr721/obsidian/internal/expression"
	"github.com/jparr721/obsidian/internal/statement"
	"github.com/jparr721/obsidian/internal/tokens"
)

const (
	inLoopStatement    = true
	notInLoopStatement = false
)

type parser struct {
	tokens  []tokens.Token
	current int
}

func NewParser(tokens []tokens.Token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *parser) Parse() ([]statement.Statement, *ParseError) {
	statements := make([]statement.Statement, 0)

	for !p.end() {
		statement, err := p.declaration(notInLoopStatement)

		// Stop parsing - report the error
		if err != nil {
			ReportParseError(err)
			return nil, err
		}

		statements = append(statements, statement)
	}

	return statements, nil
}

func (p *parser) check(tType tokens.TokenType) bool {
	if p.end() {
		return false
	}

	return p.peek().Variant == tType
}

func (p *parser) end() bool {
	return p.peek().Variant == tokens.TokenEOF
}

func (p *parser) next() tokens.Token {
	if !p.end() {
		p.current++
	}

	return p.prev()
}

func (p *parser) peek() tokens.Token {
	return p.tokens[p.current]
}

func (p *parser) prev() tokens.Token {
	return p.tokens[p.current-1]
}

func (p *parser) match(tTypes ...tokens.TokenType) bool {
	for _, t := range tTypes {
		if p.check(t) {
			p.next()
			return true
		}
	}

	return false
}

// declaration -> varDecl | statement;
func (p *parser) declaration(inLoop bool) (statement.Statement, *ParseError) {
	if p.match(tokens.TokenVar) {
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

func (p *parser) varDeclaration() (statement.Statement, *ParseError) {
	name, err := p.consume(tokens.TokenIdentifier, "Expected variable name.")

	if err != nil {
		return nil, err
	}

	var initializer expression.Expression

	if p.match(tokens.TokenEqual) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(tokens.TokenSemi, "Expected ';' after variable declaration")
	variableStatement := statement.NewVariableStatement(name, initializer)
	return variableStatement, nil
}

// statement -> statement.StatementStatement | printStatement;
func (p *parser) statement(inLoop bool) (statement.Statement, *ParseError) {
	if p.match(tokens.TokenFor) {
		return p.forStatement()
	}

	if p.match(tokens.TokenIf) {
		return p.ifStatement()
	}

	if p.match(tokens.TokenPrint) {
		return p.printStatement()
	}

	if p.match(tokens.TokenWhile) {
		return p.whileStatement()
	}

	if p.match(tokens.TokenBreak) {
		return p.breakStatement(inLoop)
	}

	if p.match(tokens.TokenOsquiggle) {
		statements, err := p.block(inLoop)

		if err != nil {
			return nil, err
		}

		return statement.NewBlockStatement(statements), nil
	}

	return p.expressionStatement()
}

// for -> "for" "(" (varDecl | exprStatement | ";") expression?";" expression?";" statement ;
func (p *parser) forStatement() (statement.Statement, *ParseError) {
	p.consume(tokens.TokenOparen, "Expected '(' after 'for'")
	var err *ParseError

	// first clause in the for loop
	var initializer statement.Statement
	if p.match(tokens.TokenSemi) {
		initializer = nil
	} else if p.match(tokens.TokenVar) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	// Second clause in the for loop
	var condition expression.Expression
	if !p.check(tokens.TokenSemi) {
		condition, err = p.expression()

		if err != nil {
			return nil, err
		}

	}

	p.consume(tokens.TokenSemi, "Expected ';' after loop condition")

	// Third clause in the for loop
	var increment expression.Expression
	if !p.check(tokens.TokenCparen) {
		increment, err = p.expression()

		if err != nil {
			return nil, err
		}

	}

	_, err = p.consume(tokens.TokenCparen, "Expected ')' after for clauses")

	if err != nil {
		return nil, err
	}

	body, err := p.statement(inLoopStatement)

	if err != nil {
		return nil, err
	}

	// Append the increment to the bottom of the block if set
	if increment != nil {
		body.(*statement.BlockStatement).Statements = append(body.(*statement.BlockStatement).Statements, statement.NewExpressionStatement(increment))
	}

	if condition == nil {
		condition = expression.NewLiteralExpression(true)
	}

	body = statement.NewWhileStatement(condition, body)

	if initializer != nil {
		block := []statement.Statement{
			initializer,
			body,
		}
		body = statement.NewBlockStatement(block)
	}

	return body, nil
}

// break -> "break";
func (p *parser) breakStatement(inLoop bool) (statement.Statement, *ParseError) {
	if inLoop {
		b := p.prev()
		_, err := p.consume(tokens.TokenSemi, "Expected ';' after break statement")

		if err != nil {
			return nil, err
		}

		return statement.NewBreakStatement(b), nil
	}

	return nil, newParseError(p.prev(), "Expected 'break' inside of while or for loop")
}

// while -> "while" "(" expression ")" statement;
func (p *parser) whileStatement() (statement.Statement, *ParseError) {
	p.consume(tokens.TokenOparen, "Expected '(' after 'while'")
	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(tokens.TokenCparen, "Expected ')' after condition")

	body, err := p.statement(inLoopStatement)

	if err != nil {
		return nil, err
	}

	return statement.NewWhileStatement(condition, body), nil
}

// if -> "if" "(" expression ")" statement ("else" statement)?;
func (p *parser) ifStatement() (statement.Statement, *ParseError) {
	p.consume(tokens.TokenOparen, "Expected '(' after 'if'")
	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(tokens.TokenCparen, "Expected ')' after if condition")

	thenBranch, err := p.statement(notInLoopStatement)

	if err != nil {
		return nil, err
	}

	var elseBranch statement.Statement
	if p.match(tokens.TokenElse) {
		elseBranch, err = p.statement(notInLoopStatement)

		if err != nil {
			return nil, err
		}
	}

	return statement.NewIfStatement(condition, thenBranch, elseBranch), nil
}

// statement.StatementStatement -> statement.Statementession ";";
func (p *parser) expressionStatement() (statement.Statement, *ParseError) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(tokens.TokenSemi, "Expected a ';' after expression.")
	if err != nil {
		return nil, err
	}

	expressionStatement := statement.NewExpressionStatement(value)

	return expressionStatement, nil
}

// printStatement -> "print" statement.Statementession ";";
func (p *parser) printStatement() (statement.Statement, *ParseError) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(tokens.TokenSemi, "Expected a ';' after value.")
	if err != nil {
		return nil, err
	}

	printStatement := statement.NewPrintStatement(value)

	return printStatement, nil
}

// block -> "{" declaration* "}";
func (p *parser) block(inLoop bool) ([]statement.Statement, *ParseError) {
	statements := make([]statement.Statement, 0)

	for !p.check(tokens.TokenCsquiggle) && !p.end() {
		statement, err := p.declaration(inLoop)

		if err != nil {
			return nil, err
		}

		statements = append(statements, statement)
	}

	_, err := p.consume(tokens.TokenCsquiggle, "Expected '}' after block statement")

	if err != nil {
		return nil, err
	}

	return statements, nil
}

// assignment -> tokens.TokenIdentifier "=" assignment | equality;
func (p *parser) assignment() (expression.Expression, *ParseError) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(tokens.TokenEqual) {
		equals := p.prev()
		value, err := p.assignment()

		if err != nil {
			return nil, err
		}

		switch expr.(type) {
		case *expression.VariableExpression:
			name := expr.(*expression.VariableExpression).Name
			return expression.NewAssignExpression(name, value), nil
		default:
			return nil, newParseError(equals, "Invalid assignment target")
		}
	}

	return expr, nil
}

// represents an expression statement
func (p *parser) expression() (expression.Expression, *ParseError) {
	return p.assignment()
}

// logical or -> logical and ("or" logical and)*
func (p *parser) or() (expression.Expression, *ParseError) {
	expr, err := p.and()

	if err != nil {
		return nil, err
	}

	for p.match(tokens.TokenOr) {
		operator := p.prev()
		right, err := p.and()

		if err != nil {
			return nil, err
		}

		expr = expression.NewLogicalExpression(expr, right, operator)
	}

	return expr, nil
}

func (p *parser) and() (expression.Expression, *ParseError) {
	expr, err := p.equality()

	if err != nil {
		return nil, err
	}

	for p.match(tokens.TokenAnd) {
		operator := p.prev()
		right, err := p.equality()

		if err != nil {
			return nil, err
		}

		expr = expression.NewLogicalExpression(expr, right, operator)
	}

	return expr, nil
}

// equality -> comparison ( ( "!=" | "==" ) comparison )*;
func (p *parser) equality() (expression.Expression, *ParseError) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.TokenBangEqual, tokens.TokenEqualEqual) {
		operator := p.prev()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinaryExpression(expr, right, operator)
	}

	return expr, nil
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )*;
func (p *parser) comparison() (expression.Expression, *ParseError) {
	expr, err := p.term()

	if err != nil {
		return nil, err
	}

	for p.match(tokens.TokenGreater, tokens.TokenGreaterEqual, tokens.TokenLess, tokens.TokenLessEqual) {
		operator := p.prev()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = expression.NewBinaryExpression(expr, right, operator)
	}

	return expr, nil
}

// term -> ( term ( "+" | "-" ) )*;
func (p *parser) term() (expression.Expression, *ParseError) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.TokenMinus, tokens.TokenPlus) {
		operator := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = expression.NewBinaryExpression(expr, right, operator)
	}

	return expr, nil
}

// factor -> ( term ( "/" | "*" ) term )*;
func (p *parser) factor() (expression.Expression, *ParseError) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.TokenSlash, tokens.TokenStar) {
		operator := p.prev()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		expr = expression.NewBinaryExpression(expr, right, operator)
	}

	return expr, nil
}

// unary -> ( "!" | "-" ) unary | primary;
func (p *parser) unary() (expression.Expression, *ParseError) {
	if p.match(tokens.TokenBang, tokens.TokenMinus) {
		operator := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return expression.NewUnaryExpression(operator, right), nil
	}

	return p.primary()
}

// primary -> tokens.TokenNumber | tokens.TokenString | "true" | "false" | "nil" | "(" expression ")";
func (p *parser) primary() (expression.Expression, *ParseError) {
	if p.match(tokens.TokenFalse) {
		return expression.NewLiteralExpression(false), nil
	}

	if p.match(tokens.TokenTrue) {
		return expression.NewLiteralExpression(true), nil
	}

	if p.match(tokens.TokenNil) {
		return expression.NewLiteralExpression(nil), nil
	}

	if p.match(tokens.TokenNumber, tokens.TokenString) {
		return expression.NewLiteralExpression(p.prev().Literal), nil
	}

	if p.match(tokens.TokenOparen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(tokens.TokenCparen, "expected ')' after expression.")
		if err != nil {
			return nil, err
		}

		return expression.NewGroupingExpression(expr), nil
	}

	if p.match(tokens.TokenIdentifier) {
		return expression.NewVariableExpression(p.prev()), nil
	}

	return nil, newParseError(p.peek(), "expected expression")
}

func (p *parser) consume(t tokens.TokenType, errorMsg string) (tokens.Token, *ParseError) {
	if p.check(t) {
		return p.next(), nil
	}

	return tokens.Token{}, newParseError(p.peek(), errorMsg)
}

// Synchronize when we've caught an error to the next valid tokens.Token
func (p *parser) synchronize() {
	p.next()

	for !p.end() {
		if p.prev().Variant == tokens.TokenSemi {
			return
		}

		switch p.peek().Variant {
		case tokens.TokenClass:
		case tokens.TokenFun:
		case tokens.TokenVar:
		case tokens.TokenFor:
		case tokens.TokenIf:
		case tokens.TokenWhile:
		case tokens.TokenPrint:
		case tokens.TokenReturn:
			return
		}

		p.next()
	}
}
