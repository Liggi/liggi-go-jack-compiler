package parser

import (
	"fmt"
	"liggi-go-jack-compiler/token"
	"reflect"
)

type Token = token.Token
type TokenMatchable = token.TokenMatchable
type Node = token.Node
type Element = token.Element
type PossibleTokens = token.PossibleTokens

type Parser struct {
	tokens []Token
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Scan() bool {
	return len(p.tokens) > 0
}

func (p *Parser) Next() Token {
	token := p.tokens[0]
	p.tokens = p.tokens[1:]

	return token
}

func (p *Parser) Peek() Token {
	if len(p.tokens) == 0 {
		return Token{}
	}

	return p.tokens[0]
}

func (p *Parser) Expect(expected TokenMatchable) (*Token, error) {
	next := p.Next()

	if expected.Match(next) {
		return &next, nil
	}

	return nil, fmt.Errorf("expected %s, encountered %s", expected, next)
}

func (p *Parser) ExpectMaybe(expected TokenMatchable) (*Token, error) {
	next := p.Peek()

	if expected.Match(next) {
		return p.Expect(expected)
	}

	return nil, nil
}

func (p *Parser) ExpectSequence(expected []TokenMatchable) ([]Token, error) {
	var parsedTokens []Token

	for _, token := range expected {
		t, err := p.Expect(token)
		if err != nil {
			return nil, err
		}
		parsedTokens = append(parsedTokens, *t)
	}

	return parsedTokens, nil
}

func (p *Parser) Parse(tokens []Token) ([]Node, error) {
	parsed := []Node{}

	for p.Scan() {
		token := p.Next()

		switch token.Value {
		case "class":
			parsedClass, err := p.parseClass(token)
			if err != nil {
				return nil, err
			}
			parsed = append(parsed, parsedClass)

		case "let":
			parsedLet, err := p.parseLet(token)
			if err != nil {
				return nil, err
			}
			parsed = append(parsed, parsedLet)

		case "do":
			parsedDo, err := p.parseDo(token)
			if err != nil {
				return nil, err
			}
			parsed = append(parsed, parsedDo)

		case "if":
			parsedIf, err := p.parseIf(token)
			if err != nil {
				return nil, err
			}
			parsed = append(parsed, parsedIf)

		default:
			return nil, fmt.Errorf("unexpected token: %s", token.Value)
		}
	}

	return parsed, nil
}

func (p *Parser) ParseUntil(terminator TokenMatchable) ([]Node, error) {
	tokens := []Node{}

	for !terminator.Match(p.Peek()) {
		next := p.Next()

		switch next.Value {
		case "field", "static":
			parsed, err := p.parseClassVar(next)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, parsed)

		case "method", "function", "constructor":
			parsed, err := p.parseSubroutine(next)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, parsed)

		case "var":
			parsed, err := p.parseVar(next)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, parsed)

		default:
			if token.AnyStatement().Match(next) {
				parsed, err := p.parseStatement(next)
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, parsed)
			} else {
				return nil, fmt.Errorf("unexpected token: %s", next.Value)
			}
		}
	}

	if len(tokens) == 0 {
		return nil, nil
	}

	return tokens, nil
}

func (p *Parser) ParseRepeatedSequenceUntil(sequence []TokenMatchable, terminator Token) ([]Node, error) {
	var parsedTokens []Node
	seqIndex := 0

	// If the very next token is the terminator, there's nothing to parse
	if p.Peek() == terminator {
		return nil, nil
	}

	for {
		// Match each token in the sequence one by one
		for seqIndex < len(sequence) {
			token, err := p.Expect(sequence[seqIndex])
			if err != nil {
				return nil, err
			}

			parsedTokens = append(parsedTokens, token)
			seqIndex++
		}

		// We've matched the entire sequence, so reset the index
		seqIndex = 0

		// At the end, check if the next token is the terminator
		// and break if it is
		if p.Peek() == terminator {
			break
		}
	}

	return parsedTokens, nil
}

func (p *Parser) parseClass(initial Token) (Node, error) {
	openingTokens, err := p.ExpectSequence([]TokenMatchable{
		token.AnyIdentifier(),
		token.Symbol('{'),
	})
	if err != nil {
		return &Element{}, err
	}

	declarations, err := p.ParseUntil(token.OneOf(token.Symbol('}'), token.AnyStatement()))
	if err != nil {
		return &Element{}, err
	}

	closingBracket, err := p.Expect(token.Symbol('}'))
	if err != nil {
		return &Element{}, err
	}

	return &Element{
		Tag: "class",
		Children: combineNodeSlices(
			[]Node{&initial, &openingTokens[0], &openingTokens[1]},
			declarations,
			[]Node{closingBracket},
		),
	}, nil
}

func (p *Parser) parseExpressionList() (Node, error) {
	var expressionList []Node

	next := p.Peek()

	// The expression list might be empty, so just return immediately if it is
	if token.Symbol(')').Match(next) {
		return &Element{
			Tag: "expressionList",
		}, nil
	}

	// Otherwise, keep parsing expressions until we run out of commas
	for {
		expression, err := p.parseExpression()
		if err != nil {
			return &Element{}, err
		}

		expressionList = append(expressionList, expression)

		comma, err := p.ExpectMaybe(token.Symbol(','))
		if err != nil {
			return &Element{}, err
		}

		if comma == nil {
			break
		} else {
			expressionList = append(expressionList, comma)
		}
	}

	return &Element{
		Tag:      "expressionList",
		Children: expressionList,
	}, nil
}

func (p *Parser) parseTerm() (Node, error) {
	next := p.Peek()

	// Identifier
	if next.Match(token.AnyIdentifier()) {
		identifier, err := p.Expect(token.AnyIdentifier())
		if err != nil {
			return &Element{}, err
		}

		// Subroutine call
		if token.OneOf(
			token.Symbol('('),
			token.Symbol('.'),
		).Match(p.Peek()) {
			subsequentIdentifiers, err := p.ParseRepeatedSequenceUntil([]TokenMatchable{
				token.Symbol('.'),
				token.AnyIdentifier(),
			}, token.Symbol('('))
			if err != nil {
				return &Element{}, err
			}

			openingBracket, err := p.Expect(token.Symbol('('))
			if err != nil {
				return &Element{}, err
			}

			expressionList, err := p.parseExpressionList()
			if err != nil {
				return &Element{}, err
			}

			closingBracket, err := p.Expect(token.Symbol(')'))
			if err != nil {
				return &Element{}, err
			}

			return &Element{
				Tag: "term",
				Children: combineNodeSlices(
					[]Node{identifier},
					subsequentIdentifiers,
					[]Node{openingBracket, expressionList, closingBracket},
				),
			}, nil
		} else if token.Symbol('[').Match(p.Peek()) {
			openBracket, err := p.Expect(token.Symbol('['))
			if err != nil {
				return &Element{}, err
			}

			expression, err := p.parseExpression()
			if err != nil {
				return &Element{}, err
			}

			closeBracket, err := p.Expect(token.Symbol(']'))
			if err != nil {
				return &Element{}, err
			}

			return &Element{
				Tag: "term",
				Children: []Node{
					identifier,
					openBracket,
					expression,
					closeBracket,
				},
			}, nil
		} else {
			// Just the identifier
			return &Element{
				Tag:      "term",
				Children: []Node{identifier},
			}, nil
		}
	}

	// Unary operation
	if token.AnyUnaryOperation().Match(next) {
		unary, err := p.Expect(token.AnyUnaryOperation())
		if err != nil {
			return &Element{}, err
		}

		term, err := p.parseTerm()
		if err != nil {
			return &Element{}, err
		}

		return &Element{
			Tag:      "term",
			Children: []Node{unary, term},
		}, nil
	}

	// Expression
	if token.Symbol('(').Match(next) {
		openingBracket, err := p.Expect(token.Symbol('('))
		if err != nil {
			return &Element{}, err
		}

		expression, err := p.parseExpression()
		if err != nil {
			return &Element{}, err
		}

		closingBracket, err := p.Expect(token.Symbol(')'))
		if err != nil {
			return &Element{}, err
		}

		return &Element{
			Tag:      "term",
			Children: []Node{openingBracket, expression, closingBracket},
		}, nil
	}

	// Ok, final possibility, it's a keyword constant, an integer constant or a string constant
	if token.AnyConstant().Match(next) {
		constant, err := p.Expect(token.AnyConstant())
		if err != nil {
			return &Element{}, err
		}

		return &Element{
			Tag:      "term",
			Children: []Node{constant},
		}, nil
	}

	return &Element{}, fmt.Errorf("unexpected token %s", next)
}

func (p *Parser) parseExpression() (Node, error) {
	term, err := p.parseTerm()
	if err != nil {
		return &Element{}, err
	}

	if !token.AnyOperation().Match(p.Peek()) {
		return &token.Element{
			Tag:      "expression",
			Children: []Node{term},
		}, nil
	}

	operation, err := p.Expect(token.AnyOperation())
	if err != nil {
		return &Element{}, err
	}

	secondTerm, err := p.parseTerm()
	if err != nil {
		return &Element{}, err
	}

	return &Element{
		Tag: "expression",
		Children: []Node{
			term,
			operation,
			secondTerm,
		},
	}, nil
}

func (p *Parser) parseWhile(initial Token) (Node, error) {
	opening, err := p.Expect(token.Symbol('('))
	if err != nil {
		return &Element{}, err
	}

	expression, err := p.parseExpression()
	if err != nil {
		return &Element{}, err
	}

	closing, err := p.Expect(token.Symbol(')'))
	if err != nil {
		return &Element{}, err
	}

	openingBracket, err := p.Expect(token.Symbol('{'))
	if err != nil {
		return &Element{}, err
	}

	statements, err := p.ParseUntil(token.Symbol('}'))
	if err != nil {
		return &Element{}, err
	}

	closingBracket, err := p.Expect(token.Symbol('}'))
	if err != nil {
		return &Element{}, err
	}

	return &Element{
		Tag: "whileStatement",
		Children: combineNodeSlices(
			[]Node{&initial, opening, expression, closing, openingBracket},
			[]Node{
				&Element{
					Tag:      "statements",
					Children: statements,
				},
			},
			[]Node{closingBracket},
		),
	}, nil
}

func (p *Parser) parseLet(initial Token) (Node, error) {
	identifier, err := p.Expect(token.AnyIdentifier())
	if err != nil {
		return &Element{}, err
	}

	opening := []Node{&initial, identifier}

	// Could be an `identifier[expression]`, so handle that case
	openSquareBracket, err := p.ExpectMaybe(token.Symbol('['))
	if err != nil {
		return &Element{}, err
	}

	if openSquareBracket != nil {
		expression, err := p.parseExpression()
		if err != nil {
			return &Element{}, err
		}

		closeSquareBracket, err := p.Expect(token.Symbol(']'))
		if err != nil {
			return &Element{}, err
		}

		opening = append(opening, openSquareBracket, expression, closeSquareBracket)
	}

	assignment, err := p.Expect(token.Symbol('='))
	if err != nil {
		return &Element{}, err
	}

	expression, err := p.parseExpression()
	if err != nil {
		return &Element{}, err
	}

	endOfLine, err := p.Expect(token.Symbol(';'))
	if err != nil {
		return &Element{}, err
	}

	return &Element{
		Tag: "letStatement",
		Children: combineNodeSlices(
			opening,
			[]Node{assignment, expression, endOfLine},
		),
	}, nil
}

func (p *Parser) parseStatement(initial Token) (Node, error) {
	switch initial.Value {
	case "let":
		return p.parseLet(initial)
	case "do":
		return p.parseDo(initial)
	case "if":
		return p.parseIf(initial)
	case "while":
		return p.parseWhile(initial)
	case "return":
		return p.parseReturn(initial)
	default:
		return nil, fmt.Errorf("unexpected token: %s", initial.Value)
	}
}

func (p *Parser) parseStatementsUntil(terminator Token) ([]Node, error) {
	var statements []Node

	for p.Peek() != terminator {
		token := p.Next()

		statement, err := p.parseStatement(token)
		if err != nil {
			return nil, err
		}

		statements = append(statements, statement)
	}

	return statements, nil
}

func (p *Parser) parseDo(initial Token) (Node, error) {
	firstIdentifier, err := p.Expect(token.AnyIdentifier())
	if err != nil {
		return &Element{}, err
	}

	subsequentIdentifiers, err := p.ParseRepeatedSequenceUntil([]TokenMatchable{
		token.Symbol('.'),
		token.AnyIdentifier(),
	}, token.Symbol('('))
	if err != nil {
		return &Element{}, err
	}

	openingBracket, err := p.Expect(token.Symbol('('))
	if err != nil {
		return &Element{}, err
	}

	expressionList, err := p.parseExpressionList()
	if err != nil {
		return &Element{}, err
	}

	closing, err := p.ExpectSequence([]TokenMatchable{
		token.Symbol(')'),
		token.Symbol(';'),
	})
	if err != nil {
		return &Element{}, err
	}

	return &Element{
		Tag: "doStatement",
		Children: combineNodeSlices(
			[]Node{&initial, firstIdentifier},
			subsequentIdentifiers,
			[]Node{
				openingBracket,
				expressionList,
				&closing[0],
				&closing[1],
			},
		),
	}, nil
}

func (p *Parser) parseReturn(initial Token) (Node, error) {
	next := p.Peek()

	if next == token.Symbol(';') {
		endOfLine, err := p.Expect(token.Symbol(';'))
		if err != nil {
			return &Element{}, err
		}

		return &Element{
			Tag: "returnStatement",
			Children: []Node{
				&initial,
				endOfLine,
			},
		}, nil
	}

	returnExpression, err := p.parseExpression()
	if err != nil {
		return &Element{}, err
	}

	endOfLine, err := p.Expect(token.Symbol(';'))
	if err != nil {
		return &Element{}, err
	}

	return &Element{
		Tag: "returnStatement",
		Children: []Node{
			&initial,
			returnExpression,
			endOfLine,
		},
	}, nil
}

func (p *Parser) parseIf(initial Token) (Node, error) {
	openBracket, err := p.Expect(token.Symbol('('))
	if err != nil {
		return &Element{}, err
	}

	expression, err := p.parseExpression()
	if err != nil {
		return &Element{}, err
	}

	openFirstStatements, err := p.ExpectSequence([]TokenMatchable{
		token.Symbol(')'),
		token.Symbol('{'),
	})
	if err != nil {
		return &Element{}, err
	}

	// Parse all the statements
	statements, err := p.parseStatementsUntil(token.Symbol('}'))
	if err != nil {
		return &Element{}, err
	}

	closingBracket, err := p.Expect(token.Symbol('}'))
	if err != nil {
		return &Element{}, err
	}

	ifNodes := combineNodeSlices(
		[]Node{&initial, openBracket, expression, &openFirstStatements[0], &openFirstStatements[1]},
		[]Node{
			&Element{
				Tag:      "statements",
				Children: statements,
			},
		},
		[]Node{closingBracket},
	)

	// Check if there's an `else` statement
	elseToken, err := p.ExpectMaybe(token.Keyword("else"))
	if err != nil {
		return &Element{}, err
	}

	if elseToken != nil {
		openSecondStatements, err := p.Expect(token.Symbol('{'))
		if err != nil {
			return &Element{}, err
		}

		// Parse all the statements
		elseStatements, err := p.parseStatementsUntil(token.Symbol('}'))
		if err != nil {
			return &Element{}, err
		}

		closingBracket, err := p.Expect(token.Symbol('}'))
		if err != nil {
			return &Element{}, err
		}

		ifNodes = combineNodeSlices(
			ifNodes,
			[]Node{elseToken, openSecondStatements},
			[]Node{
				&Element{
					Tag:      "statements",
					Children: elseStatements,
				},
			},
			[]Node{closingBracket},
		)
	}

	return &Element{
		Tag:      "ifStatement",
		Children: ifNodes,
	}, nil
}

func (p *Parser) parseClassVar(initial Token) (Node, error) {
	return p.parseVar(initial, "classVarDec")
}

func (p *Parser) parseVar(initial Token, varType ...string) (Node, error) {
	var declarationWrapper string

	if len(varType) > 0 {
		declarationWrapper = varType[0]
	} else {
		declarationWrapper = "varDec"
	}

	tokens, err := p.ExpectSequence([]TokenMatchable{
		token.ValidType(),
		token.AnyIdentifier(),
	})
	if err != nil {
		return &Element{}, err
	}

	subsequentIdentifiers, err := p.ParseRepeatedSequenceUntil([]TokenMatchable{
		token.Symbol(','),
		token.AnyIdentifier(),
	}, token.Symbol(';'))
	if err != nil {
		return &Element{}, err
	}

	endOfLine, err := p.Expect(token.Symbol(';'))
	if err != nil {
		return &Element{}, err
	}

	return &Element{
		Tag: declarationWrapper,
		Children: combineNodeSlices(
			[]Node{&initial, &tokens[0], &tokens[1]},
			subsequentIdentifiers,
			[]Node{endOfLine},
		),
	}, nil
}

func (p *Parser) parseSubroutine(initial Token) (Node, error) {
	opening, err := p.ExpectSequence([]TokenMatchable{
		token.ValidType(),
		token.AnyIdentifier(),
		token.Symbol('('),
	})
	if err != nil {
		return &Element{}, err
	}

	parameterList := &Element{
		Tag: "parameterList",
	}

	firstParameterType, err := p.ExpectMaybe(token.ValidType())
	if err != nil {
		return &Element{}, err
	}

	if firstParameterType != nil {
		firstIdentifier, err := p.Expect(token.AnyIdentifier())
		if err != nil {
			return &Element{}, err
		}

		subsequentParameters, err := p.ParseRepeatedSequenceUntil([]TokenMatchable{
			token.Symbol(','),
			token.ValidType(),
			token.AnyIdentifier(),
		}, token.Symbol(')'))
		if err != nil {
			return &Element{}, err
		}

		parameterList.Children = combineNodeSlices(
			[]Node{firstParameterType, firstIdentifier},
			subsequentParameters,
		)
	}

	predeclarations, err := p.ExpectSequence([]TokenMatchable{
		token.Symbol(')'),
		token.Symbol('{'),
	})
	if err != nil {
		return &Element{}, err
	}

	declarations, err := p.ParseUntil(token.OneOf(token.Symbol('}'), token.AnyStatement()))
	if err != nil {
		return &Element{}, err
	}

	statements, err := p.ParseUntil(token.Symbol('}'))
	if err != nil {
		return &Element{}, err
	}

	closing, err := p.ExpectSequence([]TokenMatchable{
		token.Symbol('}'),
	})
	if err != nil {
		return &Element{}, err
	}

	return &Element{
		Tag: "subroutineDec",
		Children: []Node{
			&initial,
			&opening[0],
			&opening[1],
			&opening[2],
			parameterList,
			&predeclarations[0],
			&Element{
				Tag: "subroutineBody",
				Children: combineNodeSlices(
					[]Node{&predeclarations[1]},
					declarations,
					[]Node{
						&Element{
							Tag:      "statements",
							Children: statements,
						},
						&closing[0]},
				),
			},
		},
	}, nil
}

func removeNilNodes(slice []Node) []Node {
	nonNilNodes := make([]Node, 0, len(slice))
	for _, node := range slice {
		nodeValue := reflect.ValueOf(node)
		isNil := false

		// Check if the node is a pointer and if it is nil
		if nodeValue.Kind() == reflect.Ptr {
			isNil = nodeValue.IsNil()
		}

		if !isNil {
			nonNilNodes = append(nonNilNodes, node)
		}
	}
	if len(nonNilNodes) == 0 {
		return nil
	}
	return nonNilNodes
}

func combineNodeSlices(slices ...[]Node) []Node {
	result := []Node{}
	for _, slice := range slices {
		nonNilSlice := removeNilNodes(slice)
		result = append(result, nonNilSlice...)
	}

	return result
}
