package compiler

import (
	"container/list"
	"fmt"
	"reflect"
	"strconv"
)

type AstStack struct {
	stack *list.List
}

func (s *AstStack) Push(node *AstNode) {
	s.stack.PushFront(node)
}

func (s *AstStack) Pop() error {
	if s.Size() > 0 {
		v := s.stack.Front()
		s.stack.Remove(v)
	}

	return fmt.Errorf("Stack is empty, bruh")
}

func (s *AstStack) Peek() (*AstNode, error) {
	if s.Size() > 0 {
		return s.stack.Front().Value.(*AstNode), nil
	}

	return &AstNode{}, fmt.Errorf("Stack is empty, bruh")
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
	Scope *AstStack
}

func NewAst() *Ast {
	return &Ast{
		Node: AstNode{
			Name:     "init",
			Children: make([]AstNode, 0),
		},
	}
}

func (a *Ast) Insert(node interface{}) (bool, error) {
	currentLexicalScope, err := a.Scope.Peek()

	if err != nil {
		return false, err
	}

	astNode := AstNode{
		Name:     a.inferNodeType(node),
		Value:    node,
		Children: make([]AstNode, 0),
	}

	if astNode.Name == "FunctionAstNode" && currentLexicalScope.Name == "FunctionAstNode" {
		panic("cannot nest two function expressions, idiot")
	}

	currentLexicalScope.Children = append(currentLexicalScope.Children, astNode)
	return true, nil
}

func (a *Ast) inferNodeType(node interface{}) string {
	return reflect.TypeOf(node).Name()
}

type FunctionAstNode struct {
	Name       string
	Parameters []InputParameter
}

type VariableAstNode struct {
	Name  string
	Type  string
	Value interface{}
}

type Parser struct {
	tokens     []Token
	SyntaxTree *Ast
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:     tokens,
		SyntaxTree: NewAst(),
	}
}

func (p *Parser) ConsumeFunction() {
	functionName := p.consumeToken(FUNCTION.String())
	var parameters []InputParameter
	for p.peekToken(VARIABLE.String()) {
		parameters = append(parameters, InputParameter{Name: p.consumeToken(VARIABLE.String()).Value.(string)})
	}

	if !p.peekToken(DO.String()) {
		p.consumeToken(DO.String())
	}

	functionAstNode := FunctionAstNode{Name: functionName.Value.(string), Parameters: parameters}
	p.SyntaxTree.Insert(functionAstNode)
}

func (p *Parser) ConsumeVariable() {
	variableName := p.consumeToken(VARIABLE.String()).Value.(string)
	variableType := ""
	p.consumeToken(EQUALS.String())

	var stringVariableValue string
	var numericVariableValue float64
	if p.peekToken(STRINGVALUE.String()) {
		stringVariableValue = p.consumeToken(STRINGVALUE.String()).Value.(string)
		variableType = STRINGVALUE.String()
	} else if p.peekToken(NUMERICVALUE.String()) {
		stringVariableValue = p.consumeToken(NUMERICVALUE.String()).Value.(string)

		var err error
		numericVariableValue, err = strconv.ParseFloat(stringVariableValue, 64)

		if err != nil {
			//TODO(jparr721) - Make this error not fucking terrible to reason about.
			panic(err)
		}

		variableType = NUMERICVALUE.String()
	}

	switch variableType {
	case NUMERICVALUE.String():
		variableAstNode := VariableAstNode{
			Name:  variableName,
			Type:  variableType,
			Value: numericVariableValue,
		}
		p.SyntaxTree.Insert(variableAstNode)
	case STRINGVALUE.String():
		variableAstNode := VariableAstNode{
			Name:  variableName,
			Type:  variableType,
			Value: stringVariableValue,
		}
		p.SyntaxTree.Insert(variableAstNode)
	}
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
