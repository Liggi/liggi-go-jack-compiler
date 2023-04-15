package codegenerator

import (
	"fmt"
	"liggi-go-jack-compiler/token"
	"log"
)

type CodeGenerator struct {
	code                  []token.Node
	classSymbolTable      *SymbolTable
	subroutineSymbolTable *SymbolTable
	whileStatementCount   int
	ifStatementCount      int
}

type SymbolTableRow struct {
	Name  string
	Type  string
	Kind  string
	Index int
}

type SymbolTable struct {
	Rows []SymbolTableRow
}

func NewCodeGenerator(code []token.Node) *CodeGenerator {
	return &CodeGenerator{
		code:                  code,
		classSymbolTable:      NewSymbolTable(),
		subroutineSymbolTable: NewSymbolTable(),
	}
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Rows: []SymbolTableRow{},
	}
}

func (s *SymbolTable) Add(name, kind, typ string) {
	s.Rows = append(s.Rows, SymbolTableRow{
		Name:  name,
		Type:  typ,
		Kind:  kind,
		Index: s.Count(kind),
	})
}

func (s *SymbolTable) Get(name string) *SymbolTableRow {
	for _, row := range s.Rows {
		if row.Name == name {
			return &row
		}
	}

	return nil
}

func (s *SymbolTable) Count(kind string) int {
	count := 0
	for _, row := range s.Rows {
		if row.Kind == kind {
			count++
		}
	}

	return count
}

func (s *SymbolTable) Clear() {
	s.Rows = []SymbolTableRow{}
}

func findElement(element *token.Element, tag string) *token.Element {
	for _, child := range element.Children {
		el, ok := child.(*token.Element)
		if !ok {
			continue
		}

		if el.Tag == tag {
			return el
		} else {
			found := findElement(el, tag)
			if found != nil {
				return found
			}
		}
	}

	return nil
}

func findChildElement(element *token.Element, tag string) *token.Element {
	for _, child := range element.Children {
		el, ok := child.(*token.Element)
		if !ok {
			continue
		}

		if el.Tag == tag {
			return el
		}
	}

	return nil
}

func getAllChildElements(element *token.Element, tag string) []*token.Element {
	var elements []*token.Element

	for _, child := range element.Children {
		el, ok := child.(*token.Element)
		if !ok {
			continue
		}

		if el.Tag == tag {
			elements = append(elements, el)
		} else {
			found := getAllChildElements(el, tag)
			if found != nil {
				elements = append(elements, found...)
			}
		}
	}

	return elements
}

func findToken(element *token.Element, tokenType string, value ...string) *token.Token {
	valueToFind := ""

	if len(value) > 0 {
		valueToFind = value[0]
	}

	for _, child := range element.Children {
		switch child := child.(type) {
		case *token.Token:
			if child.TokenType == tokenType && (valueToFind == "" || child.Value == valueToFind) {
				return child
			}
		}
	}

	return nil
}

func getChildTokens(element *token.Element, tokenType string, value ...string) []*token.Token {
	valueToFind := ""

	if len(value) > 0 {
		valueToFind = value[0]
	}

	var tokens []*token.Token

	for _, child := range element.Children {
		switch child := child.(type) {
		case *token.Token:
			if child.TokenType == tokenType && (valueToFind == "" || child.Value == valueToFind) {
				tokens = append(tokens, child)
			}
		}
	}

	return tokens
}

func getChildElements(element *token.Element, tag ...string) []token.Element {
	var children []token.Element

	for _, child := range element.Children {
		switch child := child.(type) {
		case *token.Element:
			if len(tag) == 0 || child.Tag == tag[0] {
				children = append(children, *child)
			}
		}
	}

	return children
}

func getVarDefs(v *token.Element) ([]string, string) {
	varNames := []string{}
	// May have multiple identifiers, let's get them all
	identifiers := getChildTokens(v, "identifier")
	typ := v.Children[1].(*token.Token).Value

	for _, identifier := range identifiers {
		varNames = append(varNames, identifier.Value)
	}

	return varNames, typ
}

func (c *CodeGenerator) compileTerm(term *token.Element) string {
	t := term.Children[0].(*token.Token)

	termHasExpression := findChildElement(term, "expression") != nil

	if termHasExpression {
		return c.compileExpression(term.Children[1].(*token.Element))
	}

	switch t.TokenType {
	case "integerConstant":
		return "push constant " + t.Value + "\n"
	case "keyword":
		switch t.Value {
		case "true":
			return "push constant 0\nnot\n"
		case "false":
			return "push constant 0\n"
		}
	case "identifier":
		// Might be a function call
		openingBracket := findToken(term, "symbol", "(")
		if openingBracket != nil {
			// Get number of expressions
			expressions := findElement(term, "expressionList")
			args := len(expressions.Children)
			// compile all the expressions!
			code, _ := c.compileExpressionList(*expressions)

			code += fmt.Sprintf("call %s.%s %d\n", t.Value, term.Children[2].(*token.Token).Value, args)

			return code
		} else {
			ident := c.findSymbol(t.Value)
			return "push " + ident.Kind + " " + fmt.Sprintf("%d", ident.Index) + "\n"
		}

	}

	return ""
}

func (c *CodeGenerator) compileExpression(expression *token.Element) string {
	var expressionCode string

	terms := getChildElements(expression, "term")

	for _, term := range terms {
		// Is the term another expression?
		termExpression := findChildElement(&term, "expression")

		// Is the first part of the term an op?
		if term.Children[0].(*token.Token).Value == "-" {
			// It's not always a term here, could also be an expression too!
			expressionCode += c.compileTerm(term.Children[1].(*token.Element))

			expressionCode += "neg\n"
		}

		if term.Children[0].(*token.Token).Value == "~" {
			// It's not always a term here, could also be an expression too!
			expressionCode += c.compileTerm(term.Children[1].(*token.Element))

			expressionCode += "not\n"
		}

		if termExpression == nil {
			expressionCode += c.compileTerm(&term)
		} else {
			expressionCode += c.compileExpression(termExpression)
		}
	}

	op := findToken(expression, "symbol")

	if op != nil {
		if op.Value == "+" {
			expressionCode += "add\n"
		}

		if op.Value == "-" {
			expressionCode += "sub\n"
		}

		if op.Value == "*" {
			expressionCode += "call Math.multiply 2\n"
		}

		if op.Value == ">" {
			expressionCode += "gt\n"
		}

		if op.Value == "&" {
			expressionCode += "and\n"
		}

		if op.Value == "=" {
			expressionCode += "eq\n"
		}
	}

	return expressionCode
}

func (c *CodeGenerator) compileStatement(statement token.Element) string {
	var code string

	switch statement.Tag {
	case "letStatement":
		code += c.compileLetStatement(statement)
	case "doStatement":
		code += c.compileDoStatement(statement)
	case "returnStatement":
		code += c.compileReturnStatement(statement)
	case "whileStatement":
		code += c.compileWhileStatement(statement)
	case "ifStatement":
		code += c.compileIfStatement(statement)
	}

	return code
}

func (c *CodeGenerator) compileIfStatement(ifStatement token.Element) string {
	var code string

	expression := findChildElement(&ifStatement, "expression")
	statementsContainer := findChildElement(&ifStatement, "statements")
	statements := getChildElements(statementsContainer)

	startLabel := fmt.Sprintf("IF_TRUE%d", c.ifStatementCount)
	endLabel := fmt.Sprintf("IF_END%d", c.ifStatementCount)
	elseLabel := fmt.Sprintf("IF_FALSE%d", c.ifStatementCount)

	c.ifStatementCount++

	code += c.compileExpression(expression)

	elseKeyword := getChildTokens(&ifStatement, "keyword", "else")

	code += fmt.Sprintf("if-goto %s\n", startLabel)
	if len(elseKeyword) > 0 {
		code += fmt.Sprintf("goto %s\n", elseLabel)
	} else {
		code += fmt.Sprintf("goto %s\n", endLabel)
	}

	code += fmt.Sprintf("label %s\n", startLabel)
	for _, statement := range statements {
		code += c.compileStatement(statement)
	}

	code += fmt.Sprintf("goto %s\n", endLabel)

	if len(elseKeyword) > 0 {
		code += fmt.Sprintf("label %s\n", elseLabel)

		elseStatementsContainer := getChildElements(&ifStatement, "statements")[1]
		elseStatements := getChildElements(&elseStatementsContainer)
		for _, statement := range elseStatements {
			code += c.compileStatement(statement)
		}
	}

	code += fmt.Sprintf("label %s\n", endLabel)

	return code
}

func (c *CodeGenerator) compileWhileStatement(whileStatement token.Element) string {
	var code string

	expression := findChildElement(&whileStatement, "expression")
	statementsContainer := findChildElement(&whileStatement, "statements")
	statements := getChildElements(statementsContainer)

	startLabel := fmt.Sprintf("WHILE_EXP%d", c.whileStatementCount)
	endLabel := fmt.Sprintf("WHILE_END%d", c.whileStatementCount)

	c.whileStatementCount++

	code = fmt.Sprintf("label %s\n", startLabel)
	code += c.compileExpression(expression)
	code += "not\n"
	code += fmt.Sprintf("if-goto %s\n", endLabel)
	for _, statement := range statements {
		code += c.compileStatement(statement)
	}

	code += fmt.Sprintf("goto %s\n", startLabel)
	code += fmt.Sprintf("label %s\n", endLabel)

	return code
}

func (c *CodeGenerator) compileLetStatement(letStatement token.Element) string {
	var code string

	identifier := findToken(&letStatement, "identifier")
	if identifier == nil {
		log.Fatalf("Expected identifier")
	}

	expression := findElement(&letStatement, "expression")
	if expression == nil {
		log.Fatalf("Expected expression")
	}

	code += c.compileExpression(expression)

	symbol := c.findSymbol(identifier.Value)

	if symbol == nil {
		log.Fatalf("Symbol not found")
	}

	code += fmt.Sprintf("pop %s %d\n", symbol.Kind, symbol.Index)

	return code
}

func getDoStatement(doStatement token.Element) (string, string) {
	// Get the class or object name
	obj := doStatement.Children[1].(*token.Token).Value

	// Get the subroutine name
	subroutineName := doStatement.Children[3].(*token.Token).Value

	return obj, subroutineName
}

func (c *CodeGenerator) findSymbol(obj string) *SymbolTableRow {
	// Find the symbol in the symbol table
	symbol := c.subroutineSymbolTable.Get(obj)
	if symbol == nil {
		symbol = c.classSymbolTable.Get(obj)
	}

	return symbol
}

func (c *CodeGenerator) compileExpressionList(expressionList token.Element) (string, int) {
	var code string

	expressions := getChildElements(&expressionList, "expression")

	for _, expression := range expressions {
		code += c.compileExpression(&expression)
	}

	return code, len(expressions)
}

func (c *CodeGenerator) compileDoStatement(doStatement token.Element) string {
	code, argCount := c.compileExpressionList(*findElement(&doStatement, "expressionList"))

	qualifier, subroutine := getDoStatement(doStatement)

	instanceVar := c.findSymbol(qualifier)

	if instanceVar != nil {
		code += fmt.Sprintf("push %s %d\n", instanceVar.Kind, instanceVar.Index)
		qualifier = instanceVar.Type
		argCount++
	}

	code += fmt.Sprintf("call %s.%s %d\n", qualifier, subroutine, argCount)

	// Do statements don't have a return value, so just dump it
	code += "pop temp 0\n"

	return code
}

func (c *CodeGenerator) compileReturnStatement(returnStatement token.Element) string {
	var code string

	expression := findElement(&returnStatement, "expression")

	if expression == nil {
		// Handle an empty return
		code += "push constant 0\n"
	} else {
		code += c.compileExpression(expression)
	}

	code += "return\n"

	return code
}

func (c *CodeGenerator) compileVarDecs(varDecs []*token.Element) int {
	var count int

	for _, varDec := range varDecs {
		// Get all the identifiers
		typ := varDec.Children[1].(*token.Token).Value
		identifiers := getChildTokens(varDec, "identifier")

		for _, identifier := range identifiers {
			// Add the identifier to the symbol table
			count++
			c.subroutineSymbolTable.Add(identifier.Value, typ, "local")
		}
	}

	return count
}

func (c *CodeGenerator) initialiseParameterList(dec *token.Element, symbolTable *SymbolTable) {
	// Get the parameter list
	parameterList := findElement(dec, "parameterList")

	// Get all the identifiers
	identifiers := getChildTokens(parameterList, "identifier")

	// Get all the types
	types := getChildTokens(parameterList, "keyword")

	// Add all the parameters to the symbol table
	for i, identifier := range identifiers {
		symbolTable.Add(identifier.Value, "argument", types[i].Value)
	}
}

func (c *CodeGenerator) compileSubroutine(class, dec *token.Element) string {
	numLocalVars := c.initSymbolTable(
		c.subroutineSymbolTable,
		getAllChildElements(dec, "varDec"),
		"local",
	)

	c.initialiseParameterList(dec, c.subroutineSymbolTable)

	c.whileStatementCount = 0
	c.ifStatementCount = 0

	className := findToken(class, "identifier").Value
	funcName := findToken(dec, "identifier").Value

	code := fmt.Sprintf("function %s.%s %d\n", className, funcName, numLocalVars)

	statementsContainer := findElement(dec, "statements")
	statements := getChildElements(statementsContainer)
	for _, statement := range statements {
		code += c.compileStatement(statement)
	}

	return code
}

func (c *CodeGenerator) compileClass(class *token.Element) string {
	var code string

	c.initSymbolTable(
		c.classSymbolTable,
		getAllChildElements(class, "classVarDec"),
		"field",
	)

	subroutineDecs := getChildElements(class, "subroutineDec")
	for _, subroutineDec := range subroutineDecs {
		code += c.compileSubroutine(class, &subroutineDec)
	}

	return code
}

func (c *CodeGenerator) initSymbolTable(s *SymbolTable, vars []*token.Element, kind string) int {
	count := 0
	s.Clear()

	for _, v := range vars {
		idents, typ := getVarDefs(v)
		for _, i := range idents {
			count++
			s.Add(i, kind, typ)
		}
	}

	return count
}

func (c *CodeGenerator) Generate() string {
	var code string

	for _, child := range c.code {
		switch child := child.(type) {
		case *token.Element:
			element := child
			switch element.Tag {
			case "class":
				code += c.compileClass(child)
			}
		}
	}

	// Strip any trailing \n if there is one
	if len(code) > 0 && code[len(code)-1] == '\n' {
		code = code[:len(code)-1]
	}

	return code
}
