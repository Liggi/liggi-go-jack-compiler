package parser

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParser_LetStatement(t *testing.T) {
	tokens := []Token{
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "a"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "integerConstant", Value: "5"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "b"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "stringConstant", Value: "Hello, world!"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "c"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "keyword", Value: "true"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "d"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "keyword", Value: "false"},
		{TokenType: "symbol", Value: ";"},
	}
	parser := NewParser(tokens)

	expected := []Node{
		&Element{
			Tag: "letStatement",
			Children: []Node{
				&Token{
					TokenType: "keyword",
					Value:     "let",
				},
				&Token{
					TokenType: "identifier",
					Value:     "a",
				},
				&Token{
					TokenType: "symbol",
					Value:     "=",
				},
				&Element{
					Tag: "expression",
					Children: []Node{
						&Element{
							Tag: "term",
							Children: []Node{
								&Token{
									TokenType: "integerConstant",
									Value:     "5",
								},
							},
						},
					},
				},
				&Token{
					TokenType: "symbol",
					Value:     ";",
				},
			},
		},
		&Element{
			Tag: "letStatement",
			Children: []Node{
				&Token{
					TokenType: "keyword",
					Value:     "let",
				},
				&Token{
					TokenType: "identifier",
					Value:     "b",
				},
				&Token{
					TokenType: "symbol",
					Value:     "=",
				},
				&Element{
					Tag: "expression",
					Children: []Node{
						&Element{
							Tag: "term",
							Children: []Node{
								&Token{
									TokenType: "stringConstant",
									Value:     "Hello, world!",
								},
							},
						},
					},
				},
				&Token{
					TokenType: "symbol",
					Value:     ";",
				},
			},
		},
		&Element{
			Tag: "letStatement",
			Children: []Node{
				&Token{
					TokenType: "keyword",
					Value:     "let",
				},
				&Token{
					TokenType: "identifier",
					Value:     "c",
				},
				&Token{
					TokenType: "symbol",
					Value:     "=",
				},
				&Element{
					Tag: "expression",
					Children: []Node{
						&Element{
							Tag: "term",
							Children: []Node{
								&Token{
									TokenType: "keyword",
									Value:     "true",
								},
							},
						},
					},
				},
				&Token{
					TokenType: "symbol",
					Value:     ";",
				},
			},
		},
		&Element{
			Tag: "letStatement",
			Children: []Node{
				&Token{
					TokenType: "keyword",
					Value:     "let",
				},
				&Token{
					TokenType: "identifier",
					Value:     "d",
				},
				&Token{
					TokenType: "symbol",
					Value:     "=",
				},
				&Element{
					Tag: "expression",
					Children: []Node{
						&Element{
							Tag: "term",
							Children: []Node{
								&Token{
									TokenType: "keyword",
									Value:     "false",
								},
							},
						},
					},
				},
				&Token{
					TokenType: "symbol",
					Value:     ";",
				},
			},
		},
	}
	parsed, _ := parser.Parse(tokens)

	if !reflect.DeepEqual(parsed, expected) {
		diff := cmp.Diff(parsed, expected)
		t.Errorf("Diff %v", diff)
	}
}

func TestParser_DoStatementNoArgs(t *testing.T) {
	tokens := []Token{
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "my_function"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "my_object"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "my_function"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "my_object"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "my_property"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "my_function"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
	}
	parser := NewParser(tokens)

	expected := []Node{
		&Element{
			Tag: "doStatement",
			Children: []Node{
				&Token{
					TokenType: "keyword",
					Value:     "do",
				},
				&Token{
					TokenType: "identifier",
					Value:     "my_function",
				},
				&Token{
					TokenType: "symbol",
					Value:     "(",
				},
				&Element{
					Tag: "expressionList",
				},
				&Token{
					TokenType: "symbol",
					Value:     ")",
				},
				&Token{
					TokenType: "symbol",
					Value:     ";",
				},
			},
		},
		&Element{
			Tag: "doStatement",
			Children: []Node{
				&Token{
					TokenType: "keyword",
					Value:     "do",
				},
				&Token{
					TokenType: "identifier",
					Value:     "my_object",
				},
				&Token{
					TokenType: "symbol",
					Value:     ".",
				},
				&Token{
					TokenType: "identifier",
					Value:     "my_function",
				},
				&Token{
					TokenType: "symbol",
					Value:     "(",
				},
				&Element{
					Tag: "expressionList",
				},
				&Token{
					TokenType: "symbol",
					Value:     ")",
				},
				&Token{
					TokenType: "symbol",
					Value:     ";",
				},
			},
		},
		&Element{
			Tag: "doStatement",
			Children: []Node{
				&Token{
					TokenType: "keyword",
					Value:     "do",
				},
				&Token{
					TokenType: "identifier",
					Value:     "my_object",
				},
				&Token{
					TokenType: "symbol",
					Value:     ".",
				},
				&Token{
					TokenType: "identifier",
					Value:     "my_property",
				},
				&Token{
					TokenType: "symbol",
					Value:     ".",
				},
				&Token{
					TokenType: "identifier",
					Value:     "my_function",
				},
				&Token{
					TokenType: "symbol",
					Value:     "(",
				},
				&Element{
					Tag: "expressionList",
				},
				&Token{
					TokenType: "symbol",
					Value:     ")",
				},
				&Token{
					TokenType: "symbol",
					Value:     ";",
				},
			},
		},
	}
	parsed, _ := parser.Parse(tokens)

	if !reflect.DeepEqual(parsed, expected) {
		diff := cmp.Diff(parsed, expected)
		t.Errorf("Diff %v", diff)
	}
}

func TestParser_DoStatementWithArgs(t *testing.T) {
	tokens := []Token{
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "my_function"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "my_function_with_args"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "false"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "keyword", Value: "true"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "integerConstant", Value: "555"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "stringConstant", Value: "hello"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
	}
	parser := NewParser(tokens)

	expected := []Node{
		&Element{
			Tag: "doStatement",
			Children: []Node{
				&Token{
					TokenType: "keyword",
					Value:     "do",
				},
				&Token{
					TokenType: "identifier",
					Value:     "my_function",
				},
				&Token{
					TokenType: "symbol",
					Value:     "(",
				},
				&Element{
					Tag: "expressionList",
					Children: []Node{
						&Token{
							TokenType: "identifier",
							Value:     "x",
						},
					},
				},
				&Token{
					TokenType: "symbol",
					Value:     ")",
				},
				&Token{
					TokenType: "symbol",
					Value:     ";",
				},
			},
		},
		&Element{
			Tag: "doStatement",
			Children: []Node{
				&Token{
					TokenType: "keyword",
					Value:     "do",
				},
				&Token{
					TokenType: "identifier",
					Value:     "my_function_with_args",
				},
				&Token{
					TokenType: "symbol",
					Value:     "(",
				},
				&Element{
					Tag: "expressionList",
					Children: []Node{
						&Token{
							TokenType: "keyword",
							Value:     "false",
						},
						&Token{
							TokenType: "symbol",
							Value:     ",",
						},
						&Token{
							TokenType: "keyword",
							Value:     "true",
						},
						&Token{
							TokenType: "symbol",
							Value:     ",",
						},
						&Token{
							TokenType: "integerConstant",
							Value:     "555",
						},
						&Token{
							TokenType: "symbol",
							Value:     ",",
						},
						&Token{
							TokenType: "stringConstant",
							Value:     "hello",
						},
					},
				},
				&Token{
					TokenType: "symbol",
					Value:     ")",
				},
				&Token{
					TokenType: "symbol",
					Value:     ";",
				},
			},
		},
	}
	parsed, _ := parser.Parse(tokens)

	if !reflect.DeepEqual(parsed, expected) {
		diff := cmp.Diff(parsed, expected)
		t.Errorf("Diff %v", diff)
	}
}

func TestParser_ClassDeclaration(t *testing.T) {
	tokens := []Token{
		{TokenType: "keyword", Value: "class"},
		{TokenType: "identifier", Value: "Main"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "field"},
		{TokenType: "keyword", Value: "int"},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "method"},
		{TokenType: "keyword", Value: "void"},
		{TokenType: "identifier", Value: "example"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "integerConstant", Value: "5"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "symbol", Value: "}"},
	}
	parser := NewParser(tokens)

	expected := []Node{
		&Element{
			Tag: "class",
			Children: []Node{
				&Token{
					TokenType: "keyword",
					Value:     "class",
				},
				&Token{
					TokenType: "identifier",
					Value:     "Main",
				},
				&Token{
					TokenType: "symbol",
					Value:     "{",
				},
				&Element{
					Tag: "classVarDec",
					Children: []Node{
						&Token{
							TokenType: "keyword",
							Value:     "field",
						},
						&Token{
							TokenType: "keyword",
							Value:     "int",
						},
						&Token{
							TokenType: "identifier",
							Value:     "x",
						},
						&Token{
							TokenType: "symbol",
							Value:     ";",
						},
					},
				},
				&Element{
					Tag: "subroutineDec",
					Children: []Node{
						&Token{
							TokenType: "keyword",
							Value:     "method",
						},
						&Token{
							TokenType: "keyword",
							Value:     "void",
						},
						&Token{
							TokenType: "identifier",
							Value:     "example",
						},
						&Token{
							TokenType: "symbol",
							Value:     "(",
						},
						&Element{
							Tag: "parameterList",
						},
						&Token{
							TokenType: "symbol",
							Value:     ")",
						},
						&Element{
							Tag: "subroutineBody",
							Children: []Node{
								&Token{
									TokenType: "symbol",
									Value:     "{",
								},
								&Element{
									Tag: "statements",
									Children: []Node{
										&Element{
											Tag: "returnStatement",
											Children: []Node{
												&Token{
													TokenType: "keyword",
													Value:     "return",
												},
												&Token{
													TokenType: "integerConstant",
													Value:     "5",
												},
												&Token{
													TokenType: "symbol",
													Value:     ";",
												},
											},
										},
									},
								},
								&Token{
									TokenType: "symbol",
									Value:     "}",
								},
							},
						},
					},
				},
				&Token{
					TokenType: "symbol",
					Value:     "}",
				},
			},
		},
	}
	parsed, _ := parser.Parse(tokens)

	if !reflect.DeepEqual(parsed, expected) {
		diff := cmp.Diff(parsed, expected)
		t.Errorf("Diff %v", diff)
	}
}
