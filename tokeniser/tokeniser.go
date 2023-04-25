package tokeniser

import (
	"bufio"
	"fmt"
	"io"
	"liggi-go-jack-compiler/token"
	"regexp"
)

type Tokeniser struct {
	scanner      *bufio.Scanner
	token_buffer rune
}

type Token = token.Token

func NewTokeniser(r io.Reader) *Tokeniser {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	return &Tokeniser{scanner: scanner}
}

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

func isDigit(char rune) bool {
	return regexp.MustCompile(`\d`).MatchString(string(char))
}

func isValidIdentifier(char rune) bool {
	return regexp.MustCompile(`[a-zA-Z]`).MatchString(string(char))
}

func isKeyword(identifier string) bool {
	keywords := []string{
		"class", "function", "void", "return", "do", "let", "var", "int", "while", "field", "constructor", "this", "method", "true", "false", "if", "else", "boolean", "null",
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
