package main

import (
	"fmt"
)

type parser struct {
  tokens  []token
  current int
}

func newParser(tokens []token) *parser {
  return &parser{
    tokens: tokens,
    current: 0,
  }
}

func (p *parser) parse() expr {
  return p.expression()
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
  return p.tokens[p.current - 1]
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
func (p *parser) primary() (expr, error) {
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
    p.consume(CPAREN, "expected ')' after expression.")
    return newGroupingExpr(expr), nil
  }

  return nil, fmt.Errorf("Expected expression: %v", p.peek())
}

func (p *parser) consume(t tokenType, errorMsg string) token {
  if p.check(t) {
    return p.next()
  }

  parseError(p.peek(), errorMsg)
  panic(fmt.Sprintf("found: %v, message: %s", p.peek(), errorMsg))
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
