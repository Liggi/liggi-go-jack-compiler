package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

type Tokeniser struct {
	scanner      *bufio.Scanner
	token_buffer rune
}

type Parser struct {
	tokens []Token
}

type Node interface {
	isNode()
}

type TokenMatchable interface {
	Match(token Token) bool
}

type Token struct {
	TokenType string
	Value     string
}

type PossibleTokens struct {
	Tokens []Token
}

type Element struct {
	Tag      string
	Children []Node
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No file specified")
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	tokeniser := NewTokeniser(file)
	tokens, err := tokeniser.Tokenise()
	if err != nil {
		log.Fatal(err)
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func NewTokeniser(r io.Reader) *Tokeniser {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	return &Tokeniser{scanner: scanner}
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (t *Token) isNode()   {}
func (e *Element) isNode() {}

func (t *Tokeniser) Scan() bool {
	if t.token_buffer != 0 {
		return true
	}

	return t.scanner.Scan()
}

func (t *Tokeniser) Text() string {
	if t.token_buffer != 0 {
		token := string(t.token_buffer)
		t.token_buffer = 0
		return string(token)
	}

	return t.scanner.Text()
}

func (t *Tokeniser) Peek() rune {
	if t.token_buffer == 0 {
		if !t.scanner.Scan() {
			return 0
		}

		t.token_buffer = rune(t.Text()[0])
	}

	return t.token_buffer
}

func isDigit(char rune) bool {
	return regexp.MustCompile(`\d`).MatchString(string(char))
}

func isValidIdentifier(char rune) bool {
	return regexp.MustCompile(`[a-zA-Z]`).MatchString(string(char))
}

func isKeyword(identifier string) bool {
	keywords := []string{
		"class", "function", "void", "return", "do", "let", "var", "int", "while", "field", "constructor", "this", "method", "true", "false", "if",
	}

	for _, keyword := range keywords {
		if keyword == identifier {
			return true
		}
	}

	return false
}

func isSymbol(char rune) bool {
	symbols := []rune{
		'{', '}', '(', ')', '[', ']', '.', ',', ';', '+', '-', '*', '/', '&', '|', '<', '>', '=', '~',
	}

	for _, symbol := range symbols {
		if symbol == char {
			return true
		}
	}

	return false
}

func isWhitespace(char rune) bool {
	return regexp.MustCompile(`\s`).MatchString(string(char))
}

func isSingleLineComment(char rune, nextChar rune) bool {
	return char == '/' && nextChar == '/'
}

func isMultiLineComment(char rune, nextChar rune) bool {
	return char == '/' && nextChar == '*'
}

func (t *Tokeniser) ScanUntilNot(exp string) string {
	regex := regexp.MustCompile(exp)

	output := ""

	// If the first character we look at stops matching, we don't need to scan any further
	immediateNextChar := string(t.Peek())
	if !regex.MatchString(immediateNextChar) {
		return ""
	}

	for t.Scan() {
		currentChar := t.Text()
		nextChar := string(t.Peek())

		output += currentChar

		if !regex.MatchString(nextChar) {
			break
		}
	}

	return output
}

func (t *Tokeniser) ScanUntil(exp string, matches bool, charsToMatch ...int) string {
	var scanAheadBy int
	if len(charsToMatch) == 0 {
		scanAheadBy = 1
	} else {
		scanAheadBy = charsToMatch[0]
	}

	regex := regexp.MustCompile(exp)

	output := ""

	for t.Scan() {
		firstChar := t.Text()
		secondChar := string(t.Peek())

		if scanAheadBy == 1 && regex.MatchString(firstChar) == matches {
			break
		}

		if scanAheadBy == 2 && regex.MatchString(firstChar+secondChar) == matches {
			t.Text()
			break
		}

		output += firstChar
	}

	return output
}

func (t *Tokeniser) Tokenise() ([]Token, error) {
	tokens := []Token{}

	for t.Scan() {
		text := t.Text()
		char := rune(text[0])

		switch {
		case isWhitespace(char):
			continue
		case isSingleLineComment(char, t.Peek()):
			t.ScanUntil(`\n`, true)
			continue
		case isMultiLineComment(char, t.Peek()):
			t.ScanUntil(`\*/`, true, 2)
			continue
		case isSymbol(char):
			token := Token{TokenType: "symbol", Value: string(char)}
			tokens = append(tokens, token)
		case isDigit(char):
			integer_const := string(char) + t.ScanUntilNot(`\d`)
			token := Token{TokenType: "integerConstant", Value: integer_const}

			tokens = append(tokens, token)
		case char == '"':
			string_const := t.ScanUntil(`"`, true)
			token := Token{TokenType: "stringConstant", Value: string_const}

			tokens = append(tokens, token)
		case isValidIdentifier(char):
			str := string(char) + t.ScanUntilNot(`[a-zA-Z0-9_]`)

			if isKeyword(str) {
				token := Token{TokenType: "keyword", Value: str}
				tokens = append(tokens, token)
			} else {
				token := Token{TokenType: "identifier", Value: str}
				tokens = append(tokens, token)
			}
		default:
			// Return an error
			return nil, fmt.Errorf("unrecognised character: %s", string(char))
		}
	}

	return tokens, nil
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
	return p.tokens[0]
}

func (p *Parser) Expect(expected TokenMatchable) (*Token, error) {
	token := p.Next()

	if expected.Match(token) {
		return &token, nil
	}

	return &Token{}, fmt.Errorf("expected %s, encountered %s", expected, token)
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

func (p *Parser) ExpectValidType() (*Token, error) {
	token := p.Next()

	if token.TokenType == "identifier" {
		return &token, nil
	}

	if token.TokenType == "keyword" && (token.Value == "int" || token.Value == "char" || token.Value == "boolean" || token.Value == "void") {
		return &token, nil
	}

	return nil, fmt.Errorf("expected valid type, encountered %s", token.TokenType)
}

func (p *Parser) ParseUntil(terminator Token) ([]Node, error) {
	tokens := []Node{}

	for p.Peek() != terminator {
		next := p.Next()

		if (next == Token{TokenType: "keyword", Value: "field"}) {
			parsed, err := p.parseClassVar(next)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, parsed)
		} else if (next == Token{TokenType: "keyword", Value: "method"}) {
			parsed, err := p.parseSubroutine(next)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, parsed)
		}
	}

	return tokens, nil
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

		default:
			return nil, fmt.Errorf("unexpected token: %s", token.Value)
		}
	}

	return parsed, nil
}

func (p *Parser) parseClass(initial Token) (Node, error) {
	openingTokens, err := p.ExpectSequence([]TokenMatchable{
		Token{TokenType: "identifier"},
		Token{TokenType: "symbol", Value: "{"},
	})
	if err != nil {
		return &Element{}, err
	}

	declarations, err := p.ParseUntil(Token{TokenType: "symbol", Value: "}"})
	if err != nil {
		return &Element{}, err
	}

	closingBracket, err := p.Expect(Token{TokenType: "symbol", Value: "}"})
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

func (p *Parser) parseLet(initial Token) (Node, error) {
	tokens, err := p.ExpectSequence([]TokenMatchable{
		Token{TokenType: "identifier"},
		Token{TokenType: "symbol", Value: "="},
		oneOf(
			Token{TokenType: "integerConstant"},
			Token{TokenType: "stringConstant"},
			Token{TokenType: "keyword", Value: "true"},
			Token{TokenType: "keyword", Value: "false"},
		),
		Token{TokenType: "symbol", Value: ";"},
	})
	if err != nil {
		return &Element{}, err
	}

	return &Element{
		Tag: "letStatement",
		Children: []Node{
			&initial,
			&tokens[0],
			&tokens[1],
			&Element{
				Tag: "expression",
				Children: []Node{
					&Element{
						Tag: "term",
						Children: []Node{
							&tokens[2],
						},
					},
				},
			},
			&tokens[3],
		},
	}, nil
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

		seqIndex = 0

		// At the end, check if the next token is the terminator
		// and break if it is
		if p.Peek() == terminator {
			break
		}
	}

	return parsedTokens, nil
}

func (p *Parser) parseDo(initial Token) (Node, error) {
	firstIdentifier, err := p.Expect(Token{TokenType: "identifier"})
	if err != nil {
		return &Element{}, err
	}

	subsequentIdentifiers, err := p.ParseRepeatedSequenceUntil([]TokenMatchable{
		Token{TokenType: "symbol", Value: "."},
		Token{TokenType: "identifier"},
	}, Token{TokenType: "symbol", Value: "("})
	if err != nil {
		return &Element{}, err
	}

	openingBracket, err := p.Expect(Token{TokenType: "symbol", Value: "("})
	if err != nil {
		return &Element{}, err
	}

	firstArgument, err := p.ExpectMaybe(oneOf(
		Token{TokenType: "integerConstant"},
		Token{TokenType: "stringConstant"},
		Token{TokenType: "keyword", Value: "true"},
		Token{TokenType: "keyword", Value: "false"},
		Token{TokenType: "identifier"},
	))
	if err != nil {
		return &Element{}, err
	}

	arguments, err := p.ParseRepeatedSequenceUntil([]TokenMatchable{
		Token{TokenType: "symbol", Value: ","},
		oneOf(
			Token{TokenType: "integerConstant"},
			Token{TokenType: "stringConstant"},
			Token{TokenType: "keyword", Value: "true"},
			Token{TokenType: "keyword", Value: "false"},
			Token{TokenType: "identifier"},
		),
	}, Token{TokenType: "symbol", Value: ")"})
	if err != nil {
		return &Element{}, err
	}

	tokens, err := p.ExpectSequence([]TokenMatchable{
		Token{TokenType: "symbol", Value: ")"},
		Token{TokenType: "symbol", Value: ";"},
	})
	if err != nil {
		return &Element{}, err
	}

	var expressionListChildren []Node
	if firstArgument != nil {
		expressionListChildren = combineNodeSlices(
			[]Node{firstArgument},
			arguments,
		)
	}

	return &Element{
		Tag: "doStatement",
		Children: combineNodeSlices(
			[]Node{&initial, firstIdentifier},
			subsequentIdentifiers,
			[]Node{
				openingBracket,
				&Element{
					Tag:      "expressionList",
					Children: expressionListChildren,
				},
				&tokens[0],
				&tokens[1],
			},
		),
	}, nil
}

func (p *Parser) parseClassVar(initial Token) (Node, error) {
	tokens, err := p.ExpectSequence([]TokenMatchable{
		validType(),
		Token{TokenType: "identifier"},
		Token{TokenType: "symbol", Value: ";"},
	})
	if err != nil {
		return &Element{}, err
	}

	// MISSING PIECE: You can declare multiple variables on the same line, comma separated
	// So I'll need to add some parsing logic to handle that possibility
	// While I'm encounting commas, keep parsing for identifiers

	return &Element{
		Tag: "classVarDec",
		Children: []Node{
			&initial,
			&tokens[0],
			&tokens[1],
			&tokens[2],
		},
	}, nil
}

func (p *Parser) parseSubroutine(initial Token) (Node, error) {
	tokens, err := p.ExpectSequence([]TokenMatchable{
		validType(),
		Token{TokenType: "identifier"},
		Token{TokenType: "symbol", Value: "("},
		Token{TokenType: "symbol", Value: ")"},
		Token{TokenType: "symbol", Value: "{"},
		Token{TokenType: "keyword", Value: "return"},
		Token{TokenType: "integerConstant", Value: "5"},
		Token{TokenType: "symbol", Value: ";"},
		Token{TokenType: "symbol", Value: "}"},
	})
	if err != nil {
		return &Element{}, err
	}

	return &Element{
		Tag: "subroutineDec",
		Children: []Node{
			&initial,
			&tokens[0],
			&tokens[1],
			&tokens[2],
			&Element{
				Tag: "parameterList",
			},
			&tokens[3],
			&Element{
				Tag: "subroutineBody",
				Children: []Node{
					&tokens[4],
					&Element{
						Tag: "statements",
						Children: []Node{
							&Element{
								Tag: "returnStatement",
								Children: []Node{
									&tokens[5],
									&tokens[6],
									&tokens[7],
								},
							},
						},
					},
					&tokens[8],
				},
			},
		},
	}, nil
}

func (t Token) Match(other Token) bool {
	return t.TokenType == other.TokenType && (other.Value != "" || t.Value == other.Value)
}

func (pt PossibleTokens) Match(other Token) bool {
	for _, token := range pt.Tokens {
		if token.Match(other) {
			return true
		}
	}
	return false
}

func removeNilNodes(slice []Node) []Node {
	nonNilNodes := make([]Node, 0, len(slice))
	for _, node := range slice {
		if node != nil {
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

func oneOf(tokens ...Token) PossibleTokens {
	return PossibleTokens{Tokens: tokens}
}

func validType() PossibleTokens {
	return oneOf(
		Token{TokenType: "keyword", Value: "int"},
		Token{TokenType: "keyword", Value: "char"},
		Token{TokenType: "keyword", Value: "boolean"},
		Token{TokenType: "identifier"},
	)
}
