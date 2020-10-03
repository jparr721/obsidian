package compiler

import (
	"container/list"
	"fmt"
)

type AstStack struct {
	stack *list.List
}

func (s *AstStack) Push(node AstNode) {
	s.stack.PushFront(node)
}

func (s *AstStack) Pop() error {
	if s.Size() > 0 {
		v := s.stack.Front()
		s.stack.Remove(v)
	}

	return fmt.Errorf("Stack is empty, bruh")
}

func (s *AstStack) Peek() (AstNode, error) {
  if s.Size() > 0 {
    return s.stack.Front().Value.(AstNode), nil
  }

	return AstNode{}, fmt.Errorf("Stack is empty, bruh")
}

func (s *AstStack) Size() int {
	return s.stack.Len()
}

func (s *AstStack) Empty() bool {
	return s.Size() == 0
}

type InputParameter struct {
	Name  string
	Value interface{}
}

// Instantiation keep track of assignments
type Instantiation struct {
	Name  string
	Value interface{}
}

func ReflectUnderlyingInstantiationType(instantiation Instantiation) {
}

type AstNode struct {
	Name     string
	Value    interface{}
	Children []AstNode
}

type Ast struct {
	Node  AstNode
	Scope AstStack
}

func NewAst() *Ast {
	return &Ast{
		Node: AstNode{
			Name:     "init",
			Children: make([]AstNode, 0),
		},
	}
}

func (a *Ast) Insert(node AstNode) (bool, error) {
  currentLexicalScope := a.Scope.Peek()
}

type FunctionAstNode struct {
	Name       string
	Parameters []InputParameter
}

type Parser struct {
	tokens     []Token
  SyntaxTree *Ast
}

func NewParser(tokens []Token) *Parser {
  return &Parser{
    tokens: tokens,
    SyntaxTree: NewAst(),
  }
}

func (p *Parser) ConsumeFunction() {
	functionName := p.consumeToken(FUNCTION.String())
	var parameters []InputParameter
	for p.peekToken(VARIABLE.String()) {
		parameters = append(parameters, InputParameter{Name: p.consumeToken(VARIABLE.String()).Value.(string)})
	}

	functionAstNode := FunctionAstNode{Name: functionName, Parameters: parameters}
}

func (p *Parser) peekToken(expected string) bool {
	token := p.tokens[0]
	return token.Name == expected
}

func (p *Parser) consumeToken(expected string) Token {
	token := p.tokens[0]

	if token.Name == expected {
		p.tokens = p.tokens[1:]
		return token
	}

	panic(fmt.Sprintf("unexpected token: '%s', expected: '%s'", expected, token.Name))
}

func (p *Parser) ParseDirective() {
	p.consumeToken(ALL.String())
}
