package parser

import (
	"fmt"

	"github.com/jparr721/obsidian/internal/expression"
	"github.com/jparr721/obsidian/internal/statement"
	"github.com/jparr721/obsidian/internal/tokens"
)

const (
	inLoopStatement    = true
	notInLoopStatement = false
)

// Parser represents the Obsidian recurisve descent parser
type Parser struct {
	tokens  []tokens.Token
	current int
	inLoop  bool
}

// NewParser creates a new parsing object for a list of tokens
func NewParser(tokens []tokens.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Parse initiates the recurisve descent parser on a supplied list of tokens
// this function returns a ParseError in the event that something fails
func (p *Parser) Parse() ([]statement.Statement, *ParseError) {
	statements := make([]statement.Statement, 0)

	for !p.end() {
		statement, err := p.declaration()

		// Stop parsing - report the error
		if err != nil {
			ReportParseError(err)
			return nil, err
		}

		statements = append(statements, statement)
	}

	return statements, nil
}

func (p *Parser) check(tType tokens.TokenType) bool {
	if p.end() {
		return false
	}

	return p.peek().Variant == tType
}

func (p *Parser) end() bool {
	return p.peek().Variant == tokens.TokenEOF
}

func (p *Parser) next() tokens.Token {
	if !p.end() {
		p.current++
	}

	return p.prev()
}

func (p *Parser) peek() tokens.Token {
	return p.tokens[p.current]
}

func (p *Parser) prev() tokens.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) match(tTypes ...tokens.TokenType) bool {
	for _, t := range tTypes {
		if p.check(t) {
			p.next()
			return true
		}
	}

	return false
}

// declaration -> varDecl | statement;
func (p *Parser) declaration() (statement.Statement, *ParseError) {
	if p.match(tokens.TokenFun) {
		return p.function("function")
	}
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

	return p.statement()
}

// return -> "return" expression? ";";
func (p *Parser) ret() (statement.Statement, *ParseError) {
	var err *ParseError

	keyword := p.prev()
	var value expression.Expression

	if !p.check(tokens.TokenSemi) {
		value, err = p.expression()

		if err != nil {
			return nil, err
		}
	}

	p.consume(tokens.TokenSemi, "Expected ';' after return value")
	return statement.NewReturnStatement(keyword, value), nil
}

// function -> identifier "(" parameters ")" block;
func (p *Parser) function(kind string) (statement.Statement, *ParseError) {
	name, err := p.consume(tokens.TokenIdentifier, fmt.Sprintf("Expected %s name.", kind))

	if err != nil {
		return nil, err
	}

	_, err = p.consume(tokens.TokenOparen, fmt.Sprintf("Expected '(' after %s name.", kind))

	if err != nil {
		return nil, err
	}

	arguments := make([]tokens.Token, 0)

	if !p.check(tokens.TokenCparen) {
		for remainingArgs := true; remainingArgs; remainingArgs = p.match(tokens.TokenComma) {
			if len(arguments) >= 255 {
				return nil, newParseError(p.peek(), "A function cannot have more than 255 arguments.")
			}

			arg, err := p.consume(tokens.TokenIdentifier, "Expected argument name.")

			if err != nil {
				return nil, err
			}

			arguments = append(arguments, arg)
		}
	}
	_, err = p.consume(tokens.TokenCparen, "Expected ')' after argument list.")

	if err != nil {
		return nil, err
	}

	_, err = p.consume(tokens.TokenOsquiggle, fmt.Sprintf("Expected '{' before %s body.", kind))

	if err != nil {
		return nil, err
	}

	body, err := p.block()

	if err != nil {
		return nil, err
	}

	return statement.NewFunctionStatement(name, arguments, body), nil
}

func (p *Parser) varDeclaration() (statement.Statement, *ParseError) {
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
func (p *Parser) statement() (statement.Statement, *ParseError) {
	if p.match(tokens.TokenFor) {
		return p.forStatement()
	}

	if p.match(tokens.TokenIf) {
		return p.ifStatement()
	}

	if p.match(tokens.TokenPrint) {
		return p.printStatement()
	}

	if p.match(tokens.TokenReturn) {
		return p.ret()
	}

	if p.match(tokens.TokenWhile) {
		return p.whileStatement()
	}

	if p.match(tokens.TokenBreak) {
		return p.breakStatement()
	}

	if p.match(tokens.TokenOsquiggle) {
		statements, err := p.block()

		if err != nil {
			return nil, err
		}

		return statement.NewBlockStatement(statements), nil
	}

	return p.expressionStatement()
}

// for -> "for" "(" (varDecl | exprStatement | ";") expression?";" expression?";" statement ;
func (p *Parser) forStatement() (statement.Statement, *ParseError) {
	p.inLoop = true
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

	body, err := p.statement()

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

	p.inLoop = false

	return body, nil
}

// break -> "break";
func (p *Parser) breakStatement() (statement.Statement, *ParseError) {
	if p.inLoop {
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
func (p *Parser) whileStatement() (statement.Statement, *ParseError) {
	p.inLoop = true
	p.consume(tokens.TokenOparen, "Expected '(' after 'while'")
	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(tokens.TokenCparen, "Expected ')' after condition")

	body, err := p.statement()

	if err != nil {
		return nil, err
	}

	p.inLoop = false

	return statement.NewWhileStatement(condition, body), nil
}

// if -> "if" "(" expression ")" statement ("else" statement)?;
func (p *Parser) ifStatement() (statement.Statement, *ParseError) {
	p.consume(tokens.TokenOparen, "Expected '(' after 'if'")
	condition, err := p.expression()

	if err != nil {
		return nil, err
	}

	p.consume(tokens.TokenCparen, "Expected ')' after if condition")

	thenBranch, err := p.statement()

	if err != nil {
		return nil, err
	}

	var elseBranch statement.Statement
	if p.match(tokens.TokenElse) {
		elseBranch, err = p.statement()

		if err != nil {
			return nil, err
		}
	}

	return statement.NewIfStatement(condition, thenBranch, elseBranch), nil
}

// statement.StatementStatement -> statement.Statement ";";
func (p *Parser) expressionStatement() (statement.Statement, *ParseError) {
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

// printStatement -> "print" statement.Statement ";";
func (p *Parser) printStatement() (statement.Statement, *ParseError) {
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
func (p *Parser) block() ([]statement.Statement, *ParseError) {
	statements := make([]statement.Statement, 0)

	for !p.check(tokens.TokenCsquiggle) && !p.end() {
		statement, err := p.declaration()

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
func (p *Parser) assignment() (expression.Expression, *ParseError) {
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
func (p *Parser) expression() (expression.Expression, *ParseError) {
	return p.assignment()
}

// logical or -> logical and ("or" logical and)*
func (p *Parser) or() (expression.Expression, *ParseError) {
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

func (p *Parser) and() (expression.Expression, *ParseError) {
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
func (p *Parser) equality() (expression.Expression, *ParseError) {
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
func (p *Parser) comparison() (expression.Expression, *ParseError) {
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
func (p *Parser) term() (expression.Expression, *ParseError) {
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
func (p *Parser) factor() (expression.Expression, *ParseError) {
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
func (p *Parser) unary() (expression.Expression, *ParseError) {
	if p.match(tokens.TokenBang, tokens.TokenMinus) {
		operator := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return expression.NewUnaryExpression(operator, right), nil
	}

	return p.call()
}

// call -> primary ( "(" arguments? ")" )*;
func (p *Parser) call() (expression.Expression, *ParseError) {
	expr, err := p.primary()

	if err != nil {
		return nil, err
	}

	for {
		if p.match(tokens.TokenOparen) {
			expr, err = p.finishCall(expr)

			if err != nil {
				return nil, err
			}

		} else {
			break
		}
	}

	return expr, nil
}

// primary -> tokens.TokenNumber | tokens.TokenString | "true" | "false" | "nil" | "(" expression ")";
func (p *Parser) primary() (expression.Expression, *ParseError) {
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

		_, err = p.consume(tokens.TokenCparen, "Expected ')' after expression.")
		if err != nil {
			return nil, err
		}

		return expression.NewGroupingExpression(expr), nil
	}

	if p.match(tokens.TokenIdentifier) {
		return expression.NewVariableExpression(p.prev()), nil
	}

	return nil, newParseError(p.peek(), "Expected expression.")
}

func (p *Parser) finishCall(callee expression.Expression) (expression.Expression, *ParseError) {
	arguments := make([]expression.Expression, 0)

	if !p.check(tokens.TokenCparen) {
		for remainingArgs := true; remainingArgs; remainingArgs = p.match(tokens.TokenComma) {
			// Cap arg length to 255.
			if len(arguments) >= 255 {
				return nil, newParseError(p.peek(), "A function cannot have more than 255 arguments.")
			}

			expr, err := p.expression()

			if err != nil {
				return nil, err
			}

			arguments = append(arguments, expr)
		}
	}

	paren, err := p.consume(tokens.TokenCparen, "Expected ')' after function arguments.")

	if err != nil {
		return nil, err
	}

	return expression.NewCallExpression(callee, paren, arguments), nil
}

func (p *Parser) consume(t tokens.TokenType, errorMsg string) (tokens.Token, *ParseError) {
	if p.check(t) {
		return p.next(), nil
	}

	return tokens.Token{}, newParseError(p.peek(), errorMsg)
}

// Synchronize when we've caught an error to the next valid tokens.Token
func (p *Parser) synchronize() {
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
