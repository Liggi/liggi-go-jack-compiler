package token

import "fmt"

type TokenMatchable interface {
	Match(token Token) bool
}

type Token struct {
	TokenType string
	Value     string
}

type PossibleTokens struct {
	Tokens []TokenMatchable
}

type Node interface {
	isNode()
}

type Element struct {
	Tag      string
	Children []Node
}

func (t Token) isNode() {}
func (t Token) Match(other Token) bool {
	return t.TokenType == other.TokenType && (other.Value == "" || t.Value == "" || t.Value == other.Value)
}

func Symbol(symbol rune) Token {
	return Token{
		TokenType: "symbol",
		Value:     string(symbol),
	}
}

func Keyword(keyword string) Token {
	return Token{
		TokenType: "keyword",
		Value:     keyword,
	}
}

func Identifier(identifier string) Token {
	return Token{
		TokenType: "identifier",
		Value:     identifier,
	}
}

func StringConstant(value string) Token {
	return Token{
		TokenType: "stringConstant",
		Value:     value,
	}
}

func AnyIdentifier() Token {
	return Token{
		TokenType: "identifier",
	}
}

func AnyStatement() PossibleTokens {
	return OneOf(
		Token{TokenType: "keyword", Value: "let"},
		Token{TokenType: "keyword", Value: "return"},
		Token{TokenType: "keyword", Value: "do"},
		Token{TokenType: "keyword", Value: "if"},
		Token{TokenType: "keyword", Value: "while"},
	)
}

func OneOf(tokens ...TokenMatchable) PossibleTokens {
	return PossibleTokens{Tokens: tokens}
}

func IntegerConstant(value int) Token {
	return Token{
		TokenType: "integerConstant",
		Value:     fmt.Sprintf("%d", value),
	}
}

func AnyStringConstant() Token {
	return Token{
		TokenType: "stringConstant",
	}
}

func AnyIntegerConstant() Token {
	return Token{
		TokenType: "integerConstant",
	}
}

func ValidType() PossibleTokens {
	return OneOf(
		Token{TokenType: "keyword", Value: "int"},
		Token{TokenType: "keyword", Value: "void"},
		Token{TokenType: "keyword", Value: "boolean"},
		Token{TokenType: "keyword", Value: "char"},
		Token{TokenType: "identifier"},
	)
}

func ValidTerm() PossibleTokens {
	return OneOf(
		AnyPrimitive(),
		Token{TokenType: "identifier"},
		Token{TokenType: "keyword", Value: "this"},
	)
}

func AnyPrimitive() PossibleTokens {
	return OneOf(
		Token{TokenType: "integerConstant"},
		Token{TokenType: "stringConstant"},
		Token{TokenType: "keyword", Value: "true"},
		Token{TokenType: "keyword", Value: "false"},
		Token{TokenType: "keyword", Value: "null"},
	)
}

func AnyConstant() PossibleTokens {
	return OneOf(
		AnyIntegerConstant(),
		AnyStringConstant(),
		AnyKeywordConstant(),
	)
}

func AnyKeywordConstant() PossibleTokens {
	return OneOf(
		Token{TokenType: "keyword", Value: "true"},
		Token{TokenType: "keyword", Value: "false"},
		Token{TokenType: "keyword", Value: "null"},
		Token{TokenType: "keyword", Value: "this"},
	)
}

func AnyOperation() PossibleTokens {
	return OneOf(
		Token{TokenType: "symbol", Value: "+"},
		Token{TokenType: "symbol", Value: "-"},
		Token{TokenType: "symbol", Value: "*"},
		Token{TokenType: "symbol", Value: "/"},
		Token{TokenType: "symbol", Value: "&"},
		Token{TokenType: "symbol", Value: "|"},
		Token{TokenType: "symbol", Value: "<"},
		Token{TokenType: "symbol", Value: ">"},
		Token{TokenType: "symbol", Value: "="},
	)
}

func AnyUnaryOperation() PossibleTokens {
	return OneOf(
		Token{TokenType: "symbol", Value: "-"},
		Token{TokenType: "symbol", Value: "~"},
	)
}

func Empty() Token {
	return Token{}
}

func (e *Element) isNode() {}

func (pt PossibleTokens) Match(other Token) bool {
	for _, token := range pt.Tokens {
		if token.Match(other) {
			return true
		}
	}
	return false
}
