package codegenerator

import (
	"fmt"
	"liggi-go-jack-compiler/parser"
	"liggi-go-jack-compiler/token"
	"liggi-go-jack-compiler/tokeniser"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func element(tag string, elements ...token.Node) *token.Element {
	return &token.Element{
		Tag:      tag,
		Children: elements,
	}
}

func tk(tokenType, value string) *token.Token {
	return &token.Token{
		TokenType: tokenType,
		Value:     value,
	}
}

func keyword(value string) *token.Token {
	return tk("keyword", value)
}

func identifier(value string) *token.Token {
	return tk("identifier", value)
}

func symbol(value rune) *token.Token {
	return tk("symbol", string(value))
}

func class(elements ...token.Node) *token.Element {
	return element("class", elements...)
}

func subroutineDec(elements ...token.Node) *token.Element {
	return element("subroutineDec", elements...)
}

func parameterList(elements ...token.Node) *token.Element {
	return element("parameterList", elements...)
}

func subroutineBody(elements ...token.Node) *token.Element {
	return element("subroutineBody", elements...)
}

func varDec(elements ...token.Node) *token.Element {
	return element("varDec", elements...)
}

func statements(elements ...token.Node) *token.Element {
	return element("statements", elements...)
}

func letStatement(elements ...token.Node) *token.Element {
	return element("letStatement", elements...)
}

func doStatement(elements ...token.Node) *token.Element {
	return element("doStatement", elements...)
}

func ifStatement(elements ...token.Node) *token.Element {
	return element("ifStatement", elements...)
}

func whileStatement(elements ...token.Node) *token.Element {
	return element("whileStatement", elements...)
}

func returnStatement(elements ...token.Node) *token.Element {
	return element("returnStatement", elements...)
}

func expressionList(elements ...token.Node) *token.Element {
	return element("expressionList", elements...)
}

func expression(elements ...token.Node) *token.Element {
	return element("expression", elements...)
}

func term(elements ...token.Node) *token.Element {
	return element("term", elements...)
}

func stringConstant(value string) *token.Token {
	return tk("stringConstant", value)
}

func integerConstant(value int) *token.Token {
	return tk("integerConstant", fmt.Sprintf("%d", value))
}

func classVarDec(elements ...token.Node) *token.Element {
	return element("classVarDec", elements...)
}

func TestCodeGenerator_Seven(t *testing.T) {
	input := `
		class Main {
			function void main() {
				 do Output.printInt(1 + (2 * 3));
				 return;
			}
	 	}`

	stringReader := strings.NewReader(input)
	tokeniser := tokeniser.NewTokeniser(stringReader)
	tokens, _ := tokeniser.Tokenise()

	parser := parser.NewParser(tokens)
	syntax, _ := parser.Parse()

	codeGenerator := NewCodeGenerator(syntax)

	expected := `function Main.main 0
		push constant 1
		push constant 2
		push constant 3
		call Math.multiply 2
		add
		call Output.printInt 1
		pop temp 0
		push constant 0
		return`
	expected = strings.ReplaceAll(expected, "\t", "")

	generated := codeGenerator.Generate()

	diff := cmp.Diff(generated, expected)
	if diff != "" {
		t.Errorf("Diff: %v", diff)
	}
}

func TestCodeGenerator_SquareMain(t *testing.T) {
	input := `
		class Main {
			function void main() {
					var SquareGame game;
					let game = SquareGame.new();
					do game.run();
					do game.dispose();
					return;
			}
		}`

	stringReader := strings.NewReader(input)
	tokeniser := tokeniser.NewTokeniser(stringReader)
	tokens, _ := tokeniser.Tokenise()

	parser := parser.NewParser(tokens)
	syntax, _ := parser.Parse()

	codeGenerator := NewCodeGenerator(syntax)

	expected := `function Main.main 1
	call SquareGame.new 0
	pop local 0
	push local 0
	call SquareGame.run 1
	pop temp 0
	push local 0
	call SquareGame.dispose 1
	pop temp 0
	push constant 0
	return`
	expected = strings.ReplaceAll(expected, "\t", "")

	generated := codeGenerator.Generate()

	diff := cmp.Diff(generated, expected)
	if diff != "" {
		t.Errorf("Diff: %v", diff)
	}
}

func TestCodeGenerator_ConvertToBin(t *testing.T) {
	input := `
		class Main {
			/**
			* Initializes RAM[8001]..RAM[8016] to -1,
			* and converts the value in RAM[8000] to binary.
			*/
			function void main() {
				var int value;
					do Main.fillMemory(8001, 16, -1); // sets RAM[8001]..RAM[8016] to -1
					let value = Memory.peek(8000);    // reads a value from RAM[8000]
					do Main.convert(value);           // performs the conversion
					return;
			}
			
			/** Converts the given decimal value to binary, and puts 
			*  the resulting bits in RAM[8001]..RAM[8016]. */
			function void convert(int value) {
				var int mask, position;
				var boolean loop;
				
				let loop = true;
				while (loop) {
						let position = position + 1;
						let mask = Main.nextMask(mask);
				
						if (~(position > 16)) {
				
								if (~((value & mask) = 0)) {
										do Memory.poke(8000 + position, 1);
									}
								else {
										do Memory.poke(8000 + position, 0);
									}    
						}
						else {
								let loop = false;
						}
				}
				return;
			}
	
			/** Returns the next mask (the mask that should follow the given mask). */
			function int nextMask(int mask) {
				if (mask = 0) {
						return 1;
				}
				else {
				return mask * 2;
				}
			}
			
			/** Fills 'length' consecutive memory locations with 'value',
				* starting at 'startAddress'. */
			function void fillMemory(int startAddress, int length, int value) {
					while (length > 0) {
							do Memory.poke(startAddress, value);
							let length = length - 1;
							let startAddress = startAddress + 1;
					}
					return;
			}
		}`

	stringReader := strings.NewReader(input)
	tokeniser := tokeniser.NewTokeniser(stringReader)
	tokens, tokeniserErr := tokeniser.Tokenise()
	if tokeniserErr != nil {
		t.Errorf("Tokeniser error: %v", tokeniserErr)
	}

	parser := parser.NewParser(tokens)
	syntax, parserErr := parser.Parse()
	if parserErr != nil {
		t.Errorf("Parser error: %v", parserErr)
	}

	codeGenerator := NewCodeGenerator(syntax)

	expected := `function Main.main 1
		push constant 8001
		push constant 16
		push constant 1
		neg
		call Main.fillMemory 3
		pop temp 0
		push constant 8000
		call Memory.peek 1
		pop local 0
		push local 0
		call Main.convert 1
		pop temp 0
		push constant 0
		return
		function Main.convert 3
		push constant 0
		not
		pop local 2
		label WHILE_EXP0
		push local 2
		not
		if-goto WHILE_END0
		push local 1
		push constant 1
		add
		pop local 1
		push local 0
		call Main.nextMask 1
		pop local 0
		push local 1
		push constant 16
		gt
		not
		if-goto IF_TRUE0
		goto IF_FALSE0
		label IF_TRUE0
		push argument 0
		push local 0
		and
		push constant 0
		eq
		not
		if-goto IF_TRUE1
		goto IF_FALSE1
		label IF_TRUE1
		push constant 8000
		push local 1
		add
		push constant 1
		call Memory.poke 2
		pop temp 0
		goto IF_END1
		label IF_FALSE1
		push constant 8000
		push local 1
		add
		push constant 0
		call Memory.poke 2
		pop temp 0
		label IF_END1
		goto IF_END0
		label IF_FALSE0
		push constant 0
		pop local 2
		label IF_END0
		goto WHILE_EXP0
		label WHILE_END0
		push constant 0
		return
		function Main.nextMask 0
		push argument 0
		push constant 0
		eq
		if-goto IF_TRUE0
		goto IF_FALSE0
		label IF_TRUE0
		push constant 1
		return
		goto IF_END0
		label IF_FALSE0
		push argument 0
		push constant 2
		call Math.multiply 2
		return
		label IF_END0
		function Main.fillMemory 0
		label WHILE_EXP0
		push argument 1
		push constant 0
		gt
		not
		if-goto WHILE_END0
		push argument 0
		push argument 2
		call Memory.poke 2
		pop temp 0
		push argument 1
		push constant 1
		sub
		pop argument 1
		push argument 0
		push constant 1
		add
		pop argument 0
		goto WHILE_EXP0
		label WHILE_END0
		push constant 0
		return`
	expected = strings.ReplaceAll(expected, "\t", "")

	generated := codeGenerator.Generate()

	fmt.Println(generated)

	expectedLines := strings.Split(expected, "\n")
	actualLines := strings.Split(generated, "\n")

	diff := cmp.Diff(actualLines, expectedLines)
	if diff != "" {
		t.Errorf("Diff: %v", diff)
	}
}
