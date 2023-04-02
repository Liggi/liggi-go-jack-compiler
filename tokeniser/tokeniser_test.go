package tokeniser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testTokeniser(t *testing.T, input string, expected []Token) {
	reader := strings.NewReader(input)
	tokeniser := NewTokeniser(reader)
	tokens, _ := tokeniser.Tokenise()

	for i, token := range tokens {
		if !reflect.DeepEqual(token, expected[i]) {
			diff := cmp.Diff(expected[i], token)
			t.Errorf("Token at index %d did not match expected Value:\n%s", i, diff)
		}
	}
}

func TestTokeniser_EmptyInput(t *testing.T) {
	input := ""
	expected := []Token{}

	testTokeniser(t, input, expected)
}

func TestTokeniser_WhitespaceOnly(t *testing.T) {
	input := "   \n\t  "
	expected := []Token{}

	testTokeniser(t, input, expected)
}

func TestTokeniser_StringConstants(t *testing.T) {
	input := "\"Hello, World\" \"String constants are working as expected\""
	expected := []Token{
		{TokenType: "stringConstant", Value: "Hello, World"},
		{TokenType: "stringConstant", Value: "String constants are working as expected"},
	}

	testTokeniser(t, input, expected)
}

func TestTokeniser_IntegerConstants(t *testing.T) {
	input := "let numbers = 123 + 456 + 789;"
	expected := []Token{
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "numbers"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "integerConstant", Value: "123"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "456"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "789"},
		{TokenType: "symbol", Value: ";"},
	}

	testTokeniser(t, input, expected)
}

func TestTokeniser_SingleLineComments(t *testing.T) {
	input := `let str = "Hello, World!"; // hello world is a test program
	// both of these single line comments should be ignored`
	expected := []Token{
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "str"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "stringConstant", Value: "Hello, World!"},
		{TokenType: "symbol", Value: ";"},
	}

	testTokeniser(t, input, expected)
}

func TestTokeniser_MultiLineComments(t *testing.T) {
	input := `let str = "Hello, World!";
	/* And now we have a large, multi-line comment block and we
	want to make sure that it doesn't get tokenised, so let's run it
	through the tokeniser and see what comes out the other side */`
	expected := []Token{
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "str"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "stringConstant", Value: "Hello, World!"},
		{TokenType: "symbol", Value: ";"},
	}

	testTokeniser(t, input, expected)
}

func TestTokeniser_IdentifiersAndSymbols(t *testing.T) {
	input := `let var1 = 1234;
	let var2 = "string" + 9999;
	return;`
	expected := []Token{
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "var1"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "integerConstant", Value: "1234"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "var2"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "stringConstant", Value: "string"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "9999"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "symbol", Value: ";"},
	}

	testTokeniser(t, input, expected)
}

func TestTokeniser_Array(t *testing.T) {
	input := `// This file is part of www.nand2tetris.org
	// and the book "The Elements of Computing Systems"
	// by Nisan and Schocken, MIT Press.
	// File name: projects/10/ArrayTest/Main.jack
	
	// (identical to projects/09/Average/Main.jack)
	
	/** Computes the average of a sequence of integers. */
	class Main {
		function void main() {
			var Array a;
			var int length;
			var int i, sum;
		
			let length = Keyboard.readInt("HOW MANY NUMBERS? ");
			let a = Array.new(length);
			let i = 0;

			while (i < length) {
					let a[i] = Keyboard.readInt("ENTER THE NEXT NUMBER: ");
					let i = i + 1;
			}

			let i = 0;
			let sum = 0;

			while (i < length) {
					let sum = sum + a[i];
					let i = i + 1;
			}

			do Output.printString("THE AVERAGE IS: ");
			do Output.printInt(sum / length);
			do Output.println();

			return;
		}
	}`
	expected := []Token{
		{TokenType: "keyword", Value: "class"},
		{TokenType: "identifier", Value: "Main"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "function"},
		{TokenType: "keyword", Value: "void"},
		{TokenType: "identifier", Value: "main"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "var"},
		{TokenType: "identifier", Value: "Array"},
		{TokenType: "identifier", Value: "a"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "var"},
		{TokenType: "keyword", Value: "int"},
		{TokenType: "identifier", Value: "length"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "var"},
		{TokenType: "keyword", Value: "int"},
		{TokenType: "identifier", Value: "i"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "sum"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "length"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "Keyboard"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "readInt"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "stringConstant", Value: "HOW MANY NUMBERS? "},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "a"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "Array"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "new"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "length"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "i"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "integerConstant", Value: "0"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "while"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "i"},
		{TokenType: "symbol", Value: "<"},
		{TokenType: "identifier", Value: "length"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "a"},
		{TokenType: "symbol", Value: "["},
		{TokenType: "identifier", Value: "i"},
		{TokenType: "symbol", Value: "]"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "Keyboard"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "readInt"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "stringConstant", Value: "ENTER THE NEXT NUMBER: "},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "i"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "i"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "i"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "integerConstant", Value: "0"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "sum"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "integerConstant", Value: "0"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "while"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "i"},
		{TokenType: "symbol", Value: "<"},
		{TokenType: "identifier", Value: "length"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "sum"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "sum"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "a"},
		{TokenType: "symbol", Value: "["},
		{TokenType: "identifier", Value: "i"},
		{TokenType: "symbol", Value: "]"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "i"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "i"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Output"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "printString"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "stringConstant", Value: "THE AVERAGE IS: "},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Output"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "printInt"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "sum"},
		{TokenType: "symbol", Value: "/"},
		{TokenType: "identifier", Value: "length"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Output"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "println"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "symbol", Value: "}"},
	}

	testTokeniser(t, input, expected)
}

func TestTokeniser_Square(t *testing.T) {
	input := `// This file is part of www.nand2tetris.org
	// and the book "The Elements of Computing Systems"
	// by Nisan and Schocken, MIT Press.
	// File name: projects/10/Square/Square.jack
	
	// (same as projects/09/Square/Square.jack)
	
	/** Implements a graphical square. */
	class Square {
	
		 field int x, y; // screen location of the square's top-left corner
		 field int size; // length of this square, in pixels
	
		 /** Constructs a new square with a given location and size. */
		 constructor Square new(int Ax, int Ay, int Asize) {
				let x = Ax;
				let y = Ay;
				let size = Asize;
				do draw();
				return this;
		 }
	
		 /** Disposes this square. */
		 method void dispose() {
				do Memory.deAlloc(this);
				return;
		 }
	
		 /** Draws the square on the screen. */
		 method void draw() {
				do Screen.setColor(true);
				do Screen.drawRectangle(x, y, x + size, y + size);
				return;
		 }
	
		 /** Erases the square from the screen. */
		 method void erase() {
				do Screen.setColor(false);
				do Screen.drawRectangle(x, y, x + size, y + size);
				return;
		 }
	
			/** Increments the square size by 2 pixels. */
		 method void incSize() {
				if (((y + size) < 254) & ((x + size) < 510)) {
					 do erase();
					 let size = size + 2;
					 do draw();
				}
				return;
		 }
	
		 /** Decrements the square size by 2 pixels. */
		 method void decSize() {
				if (size > 2) {
					 do erase();
					 let size = size - 2;
					 do draw();
				}
				return;
		 }
	
		 /** Moves the square up by 2 pixels. */
		 method void moveUp() {
				if (y > 1) {
					 do Screen.setColor(false);
					 do Screen.drawRectangle(x, (y + size) - 1, x + size, y + size);
					 let y = y - 2;
					 do Screen.setColor(true);
					 do Screen.drawRectangle(x, y, x + size, y + 1);
				}
				return;
		 }
	
		 /** Moves the square down by 2 pixels. */
		 method void moveDown() {
				if ((y + size) < 254) {
					 do Screen.setColor(false);
					 do Screen.drawRectangle(x, y, x + size, y + 1);
					 let y = y + 2;
					 do Screen.setColor(true);
					 do Screen.drawRectangle(x, (y + size) - 1, x + size, y + size);
				}
				return;
		 }
	
		 /** Moves the square left by 2 pixels. */
		 method void moveLeft() {
				if (x > 1) {
					 do Screen.setColor(false);
					 do Screen.drawRectangle((x + size) - 1, y, x + size, y + size);
					 let x = x - 2;
					 do Screen.setColor(true);
					 do Screen.drawRectangle(x, y, x + 1, y + size);
				}
				return;
		 }
	
		 /** Moves the square right by 2 pixels. */
		 method void moveRight() {
				if ((x + size) < 510) {
					 do Screen.setColor(false);
					 do Screen.drawRectangle(x, y, x + 1, y + size);
					 let x = x + 2;
					 do Screen.setColor(true);
					 do Screen.drawRectangle((x + size) - 1, y, x + size, y + size);
				}
				return;
		 }
	}`
	expected := []Token{
		{TokenType: "keyword", Value: "class"},
		{TokenType: "identifier", Value: "Square"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "field"},
		{TokenType: "keyword", Value: "int"},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "field"},
		{TokenType: "keyword", Value: "int"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "constructor"},
		{TokenType: "identifier", Value: "Square"},
		{TokenType: "identifier", Value: "new"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "int"},
		{TokenType: "identifier", Value: "Ax"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "keyword", Value: "int"},
		{TokenType: "identifier", Value: "Ay"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "keyword", Value: "int"},
		{TokenType: "identifier", Value: "Asize"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "Ax"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "Ay"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "Asize"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "draw"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "keyword", Value: "this"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "method"},
		{TokenType: "keyword", Value: "void"},
		{TokenType: "identifier", Value: "dispose"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Memory"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "deAlloc"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "this"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "method"},
		{TokenType: "keyword", Value: "void"},
		{TokenType: "identifier", Value: "draw"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "setColor"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "true"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "drawRectangle"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "method"},
		{TokenType: "keyword", Value: "void"},
		{TokenType: "identifier", Value: "erase"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "setColor"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "false"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "drawRectangle"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "method"},
		{TokenType: "keyword", Value: "void"},
		{TokenType: "identifier", Value: "incSize"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "if"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "<"},
		{TokenType: "integerConstant", Value: "254"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "&"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "<"},
		{TokenType: "integerConstant", Value: "510"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "erase"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "2"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "draw"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "method"},
		{TokenType: "keyword", Value: "void"},
		{TokenType: "identifier", Value: "decSize"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "if"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ">"},
		{TokenType: "integerConstant", Value: "2"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "erase"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: "-"},
		{TokenType: "integerConstant", Value: "2"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "draw"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "method"},
		{TokenType: "keyword", Value: "void"},
		{TokenType: "identifier", Value: "moveUp"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "if"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: ">"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "setColor"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "false"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "drawRectangle"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "-"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "-"},
		{TokenType: "integerConstant", Value: "2"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "setColor"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "true"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "drawRectangle"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "method"},
		{TokenType: "keyword", Value: "void"},
		{TokenType: "identifier", Value: "moveDown"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "if"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "<"},
		{TokenType: "integerConstant", Value: "254"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "setColor"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "false"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "drawRectangle"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "2"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "setColor"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "true"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "drawRectangle"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "-"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "method"},
		{TokenType: "keyword", Value: "void"},
		{TokenType: "identifier", Value: "moveLeft"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "if"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ">"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "setColor"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "false"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "drawRectangle"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "-"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "-"},
		{TokenType: "integerConstant", Value: "2"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "setColor"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "true"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "drawRectangle"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "method"},
		{TokenType: "keyword", Value: "void"},
		{TokenType: "identifier", Value: "moveRight"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "if"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "<"},
		{TokenType: "integerConstant", Value: "510"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "{"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "setColor"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "false"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "drawRectangle"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "let"},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "="},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "integerConstant", Value: "2"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "setColor"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "keyword", Value: "true"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "keyword", Value: "do"},
		{TokenType: "identifier", Value: "Screen"},
		{TokenType: "symbol", Value: "."},
		{TokenType: "identifier", Value: "drawRectangle"},
		{TokenType: "symbol", Value: "("},
		{TokenType: "symbol", Value: "("},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: "-"},
		{TokenType: "integerConstant", Value: "1"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "x"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ","},
		{TokenType: "identifier", Value: "y"},
		{TokenType: "symbol", Value: "+"},
		{TokenType: "identifier", Value: "size"},
		{TokenType: "symbol", Value: ")"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "keyword", Value: "return"},
		{TokenType: "symbol", Value: ";"},
		{TokenType: "symbol", Value: "}"},
		{TokenType: "symbol", Value: "}"},
	}

	testTokeniser(t, input, expected)
}
