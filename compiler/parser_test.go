package compiler

import "testing"

type ParserTest struct {
	Name     string
	Input    string
	Expected map[int]interface{}
}

func TestParsingEfficacy(test *testing.T) {
	allTests := []ParserTest{
		{
			Name:  "Basic Function Parsing",
			Input: "Function x, y do x = \"foobar\" end",
			Expected: map[int]interface{}{
				0: AstNode{
					Name: "FunctionAstNode",
					Value: FunctionAstNode{
						Name: "Function",
						Parameters: []InputParameter{
							{
								Name: "x",
							},
							{
								Name: "y",
							},
						},
					},
				},
			},
		},
	}

  for _, parserTest := range allTests {
    test.Logf("TestParsingEfficacy -- %s", parserTest.Name)
    parser := NewParser()
  }
}
