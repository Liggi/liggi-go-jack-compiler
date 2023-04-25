package codegenerator

import (
	"fmt"
	"liggi-go-jack-compiler/token"
)

type CodeGenerator struct {
	code                  []token.Node
	classSymbolTable      *SymbolTable
	subroutineSymbolTable *SymbolTable
	whileStatementCount   int
	ifStatementCount      int
	className             string
}

type Symbol struct {
	Name  string
	Type  string
	Kind  string
	Index int
}

type SymbolTable struct {
	Entries []Symbol
}

type SubroutineType int

const (
	Function SubroutineType = iota
	Method
	Constructor
)

type SubroutineDec struct {
	name           string
	returnType     string
	subroutineType SubroutineType
}

type Statement struct {
	typ string
}

func NewCodeGenerator(code []token.Node) *CodeGenerator {
	return &CodeGenerator{
		code:                  code,
		classSymbolTable:      NewSymbolTable(),
		subroutineSymbolTable: NewSymbolTable(),
	}
}

var CharacterMap = map[rune]int{
	' ':  32,
	'!':  33,
	'"':  34,
	'#':  35,
	'$':  36,
	'%':  37,
	'&':  38,
	'\'': 39,
	'(':  40,
	')':  41,
	'*':  42,
	'+':  43,
	',':  44,
	'-':  45,
	'.':  46,
	'/':  47,
	'0':  48,
	'1':  49,
	'2':  50,
	'3':  51,
	'4':  52,
	'5':  53,
	'6':  54,
	'7':  55,
	'8':  56,
	'9':  57,
	':':  58,
	';':  59,
	'<':  60,
	'=':  61,
	'>':  62,
	'?':  63,
	'@':  64,
	'A':  65,
	'B':  66,
	'C':  67,
	'D':  68,
	'E':  69,
	'F':  70,
	'G':  71,
	'H':  72,
	'I':  73,
	'J':  74,
	'K':  75,
	'L':  76,
	'M':  77,
	'N':  78,
	'O':  79,
	'P':  80,
	'Q':  81,
	'R':  82,
	'S':  83,
	'T':  84,
	'U':  85,
	'V':  86,
	'W':  87,
	'X':  88,
	'Y':  89,
	'Z':  90,
	'[':  91,
	'\\': 92,
	']':  93,
	'^':  94,
	'_':  95,
	'`':  96,
	'a':  97,
	'b':  98,
	'c':  99,
	'd':  100,
	'e':  101,
	'f':  102,
	'g':  103,
	'h':  104,
	'i':  105,
	'j':  106,
	'k':  107,
	'l':  108,
	'm':  109,
	'n':  110,
	'o':  111,
	'p':  112,
	'q':  113,
	'r':  114,
	's':  115,
	't':  116,
	'u':  117,
	'v':  118,
	'w':  119,
	'x':  120,
	'y':  121,
	'z':  122,
	'{':  123,
	'|':  124,
	'}':  125,
	'~':  126,
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Entries: []Symbol{},
	}
}

func (s *SymbolTable) Add(name, kind, typ string) {
	s.Entries = append(s.Entries, Symbol{
		Name:  name,
		Type:  typ,
		Kind:  kind,
		Index: s.Count(kind),
	})
}

func (s *SymbolTable) Get(name string) Symbol {
	for _, row := range s.Entries {
		if row.Name == name {
			return row
		}
	}

	return Symbol{}
}

func (s *SymbolTable) Count(kind string) int {
	count := 0
	for _, row := range s.Entries {
		if row.Kind == kind {
			count++
		}
	}

	return count
}

func (s *SymbolTable) Clear() {
	s.Entries = []Symbol{}
}

func (s *Symbol) Push() string {
	segment := s.Kind

	if s.Kind == "field" {
		segment = "this"
	}

	return fmt.Sprintf("push %s %d\n", segment, s.Index)
}

func (s *Symbol) Pop() string {
	segment := s.Kind

	if s.Kind == "field" {
		segment = "this"
	}

	return fmt.Sprintf("pop %s %d\n", segment, s.Index)
}

func SubroutineDecFromSyntax(syntax *token.Element) (*SubroutineDec, error) {
	var subroutineType SubroutineType
	defKeyword := syntax.Children[0].(*token.Token).Value

	switch defKeyword {
	case "function":
		subroutineType = Function
	case "method":
		subroutineType = Method
	case "constructor":
		subroutineType = Constructor
	default:
		return nil, fmt.Errorf("unknown subroutine type: %s", defKeyword)
	}

	name, err := syntax.ChildAsToken(2)
	if err != nil {
		return nil, err
	}

	returnType, err := syntax.ChildAsToken(1)
	if err != nil {
		return nil, err
	}

	return &SubroutineDec{
		name:           name.Value,
		returnType:     returnType.Value,
		subroutineType: subroutineType,
	}, nil
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

func getVarDefs(v *token.Element) ([]string, string, error) {
	varNames := []string{}

	typ, err := v.ChildAsToken(1)
	if err != nil {
		return nil, "", err
	}

	for i := 2; i < len(v.Children); i++ {
		if identifier, ok := v.Children[i].(*token.Token); ok {
			// Check if the token type is "identifier" before appending
			if identifier.TokenType == "identifier" {
				varNames = append(varNames, identifier.Value)
			}
		} else {
			return nil, "", fmt.Errorf("unexpected non-token child at index %d", i)
		}
	}

	return varNames, typ.Value, nil
}

func isUnaryOperation(op string) bool {
	return op == "-" || op == "~"
}

func unaryOpToCode(op string) (string, error) {
	switch op {
	case "-":
		return "neg", nil
	case "~":
		return "not", nil
	default:
		return "", fmt.Errorf("unknown unary operator: %s", op)
	}
}

func opToCode(op string) (string, error) {
	switch op {
	case "+":
		return "add", nil
	case "-":
		return "sub", nil
	case "*":
		return "call Math.multiply 2", nil
	case "/":
		return "call Math.divide 2", nil
	case "&":
		return "and", nil
	case "|":
		return "or", nil
	case "<":
		return "lt", nil
	case ">":
		return "gt", nil
	case "=":
		return "eq", nil
	default:
		return "", fmt.Errorf("unknown operator: %s", op)
	}
}

func (c *CodeGenerator) compileExpressionTerm(term *token.Element) (string, error) {
	exp, err := term.ChildAsElement(1)
	if err != nil {
		return "", err
	}

	compiledExpression, err := c.compileExpression(exp)
	if err != nil {
		return "", err
	}

	return compiledExpression, nil
}

func (c *CodeGenerator) compileFunctionCallTerm(term *token.Element) (string, error) {
	var code string
	var numArgs int

	qualifierToken, err := term.ChildAsToken(0)
	if err != nil {
		return "", err
	}

	qualifier := qualifierToken.Value
	subroutineName, err := term.ChildAsToken(2)
	if err != nil {
		return "", err
	}

	instanceVar := c.findSymbol(qualifier)

	if (instanceVar != Symbol{}) {
		code += instanceVar.Push()

		qualifier = instanceVar.Type
		numArgs = 1
	}

	if qualifier == "" {
		code += "push pointer 0\n"

		qualifier = c.className
		numArgs = 1
	}

	expressionList := term.FindElement("expressionList")
	expressions := expressionList.AllChildElementsByTag("expression")

	args := len(expressions)
	compiledExpressionList, _, err := c.compileExpressionList(*expressionList)
	if err != nil {
		return "", err
	}

	code += compiledExpressionList

	numArgs += args

	code += fmt.Sprintf("call %s.%s %d\n", qualifier, subroutineName.Value, numArgs)

	return code, nil
}

func (c *CodeGenerator) compileString(s string) (string, error) {
	code := fmt.Sprintf("push constant %d\n", len(s))
	code += "call String.new 1\n"

	for _, char := range s {
		charCode, ok := CharacterMap[char]
		if !ok {
			return "", fmt.Errorf("unknown character: %s", string(char))
		}

		code += fmt.Sprintf("push constant %d\n", charCode)
		code += "call String.appendChar 2\n"
	}

	return code, nil
}

func (c *CodeGenerator) compileTerm(term *token.Element) (string, error) {
	// If the term has an expression, we compile the expression, then return
	if findChildElement(term, "expression") != nil {
		return c.compileExpressionTerm(term)
	}

	// If the expression is a function call, we compile the expression list
	// call the function, then return
	if term.FindElement("expressionList") != nil {
		return c.compileFunctionCallTerm(term)
	}

	token, err := term.ChildAsToken(0)
	if err != nil {
		return "", err
	}

	// Otherwise, it's just a simple push operation
	switch token.TokenType {
	case "integerConstant":
		return "push constant " + token.Value + "\n", nil

	case "stringConstant":
		return c.compileString(token.Value)

	case "keyword":
		switch token.Value {
		case "true":
			return "push constant 0\nnot\n", nil
		case "false":
			return "push constant 0\n", nil
		case "this":
			return "push pointer 0\n", nil
		case "null":
			return "push constant 0\n", nil
		}

	case "identifier":
		ident := c.findSymbol(token.Value)
		return ident.Push(), nil
	}

	return "", nil
}

func (c *CodeGenerator) compileExpression(expression *token.Element) (string, error) {
	var expressionCode string

	terms := expression.AllChildElementsByTag("term")

	for _, term := range terms {
		// Is the term an array access?
		if term.FindChildToken("symbol", "[") != nil {
			identifier, err := term.ChildAsToken(0)
			if err != nil {
				return "", err
			}

			ident := c.findSymbol(identifier.Value)

			arrayAccessExpression := term.FindElement("expression")
			compiledExpression, err := c.compileExpression(arrayAccessExpression)
			if err != nil {
				return "", err
			}

			expressionCode += compiledExpression
			expressionCode += ident.Push()
			expressionCode += "add\n"
			expressionCode += "pop pointer 1\n"
			expressionCode += "push that 0\n"

			continue
		}

		// Is the term another expression?
		termExpression := term.FindChildElement("expression")
		if termExpression != nil {
			compiledExpression, err := c.compileExpression(termExpression)
			if err != nil {
				return "", err
			}

			expressionCode += compiledExpression

			continue
		}

		// Is it a unary operation?
		unaryOp, _ := term.ChildAsToken(0)
		if isUnaryOperation(unaryOp.Value) {
			termToCompile, err := term.ChildAsElement(1)
			if err != nil {
				return "", err
			}

			compiledTerm, err := c.compileTerm(termToCompile)
			if err != nil {
				return "", err
			}

			expressionCode += compiledTerm
			op, err := unaryOpToCode(unaryOp.Value)
			if err != nil {
				return "", err
			}

			expressionCode += op + "\n"

			continue
		}

		// Otherwise, it's just a regular term
		compiledTerm, err := c.compileTerm(term)
		if err != nil {
			return "", err
		}

		expressionCode += compiledTerm
	}

	// Handle the operation, if there is one!
	op := expression.FindChildToken("symbol")
	if op != nil {
		operation, err := opToCode(op.Value)
		if err != nil {
			return "", err
		}

		expressionCode += operation + "\n"
	}

	return expressionCode, nil
}

func (c *CodeGenerator) compileStatement(statement token.Element) (string, error) {
	var code string

	switch statement.Tag {
	case "letStatement":
		compiledLet, err := c.compileLetStatement(statement)
		if err != nil {
			return "", fmt.Errorf("failed compiling let statement: %w", err)
		}

		code += compiledLet
	case "doStatement":
		compiledDo, err := c.compileDoStatement(statement)
		if err != nil {
			return "", fmt.Errorf("failed compiling do statement: %w", err)
		}

		code += compiledDo
	case "returnStatement":
		compiledReturn, err := c.compileReturnStatement(statement)
		if err != nil {
			return "", fmt.Errorf("failed compiling return statement: %w", err)
		}

		code += compiledReturn
	case "whileStatement":
		compiledWhile, err := c.compileWhileStatement(statement)
		if err != nil {
			return "", fmt.Errorf("failed compiling while statement: %w", err)
		}

		code += compiledWhile
	case "ifStatement":
		compiledIf, err := c.compileIfStatement(statement)
		if err != nil {
			return "", fmt.Errorf("failed compiling if statement: %w", err)
		}

		code += compiledIf
	}

	return code, nil
}

func (c *CodeGenerator) compileIfStatement(ifStatement token.Element) (string, error) {
	var code string

	expression := findChildElement(&ifStatement, "expression")
	if expression == nil {
		return "", fmt.Errorf("failed to find expression")
	}

	statementsContainer := findChildElement(&ifStatement, "statements")
	if statementsContainer == nil {
		return "", fmt.Errorf("failed to find statements")
	}
	statements := statementsContainer.AllChildElements()

	startLabel := fmt.Sprintf("IF_TRUE%d", c.ifStatementCount)
	endLabel := fmt.Sprintf("IF_END%d", c.ifStatementCount)
	elseLabel := fmt.Sprintf("IF_FALSE%d", c.ifStatementCount)

	c.ifStatementCount++

	compiledExpression, err := c.compileExpression(expression)
	if err != nil {
		return "", err
	}
	code += compiledExpression

	elseKeyword := ifStatement.FindChildToken("keyword", "else")

	if elseKeyword == nil {
		code += fmt.Sprintf("if-goto %s\n", startLabel)
		code += fmt.Sprintf("goto %s\n", elseLabel) // This is confusing

		code += fmt.Sprintf("label %s\n", startLabel)

		for _, statement := range statements {
			compiledStatement, err := c.compileStatement(*statement)
			if err != nil {
				return "", err
			}

			code += compiledStatement
		}

		code += fmt.Sprintf("label %s\n", elseLabel)
	} else {
		code += fmt.Sprintf("if-goto %s\n", startLabel)
		code += fmt.Sprintf("goto %s\n", elseLabel)

		code += fmt.Sprintf("label %s\n", startLabel)
		for _, statement := range statements {
			compiledStatement, err := c.compileStatement(*statement)
			if err != nil {
				return "", err
			}

			code += compiledStatement
		}

		code += fmt.Sprintf("goto %s\n", endLabel)

		code += fmt.Sprintf("label %s\n", elseLabel)

		// TODO: Bad, make good
		statementsBlocks := ifStatement.AllChildElementsByTag("statements")
		if len(statementsBlocks) < 2 {
			return "", fmt.Errorf("failed to find else statements, despite finding else keyword")
		}
		elseStatementsContainer := statementsBlocks[1]

		elseStatements := elseStatementsContainer.AllChildElements()
		for _, statement := range elseStatements {
			compiledStatement, err := c.compileStatement(*statement)
			if err != nil {
				return "", err
			}

			code += compiledStatement
		}

		code += fmt.Sprintf("label %s\n", endLabel)
	}

	return code, nil
}

func (c *CodeGenerator) compileWhileStatement(whileStatement token.Element) (string, error) {
	var code string

	expression := findChildElement(&whileStatement, "expression")
	if expression == nil {
		return "", fmt.Errorf("failed to find expression")
	}

	statementsContainer := findChildElement(&whileStatement, "statements")
	if statementsContainer == nil {
		return "", fmt.Errorf("failed to find statements")
	}

	statements := statementsContainer.AllChildElements()

	startLabel := fmt.Sprintf("WHILE_EXP%d", c.whileStatementCount)
	endLabel := fmt.Sprintf("WHILE_END%d", c.whileStatementCount)

	c.whileStatementCount++

	code = fmt.Sprintf("label %s\n", startLabel)

	compiledExpression, err := c.compileExpression(expression)
	if err != nil {
		return "", err
	}
	code += compiledExpression
	code += "not\n"
	code += fmt.Sprintf("if-goto %s\n", endLabel)
	for _, statement := range statements {
		compiledStatement, err := c.compileStatement(*statement)
		if err != nil {
			return "", err
		}

		code += compiledStatement
	}

	code += fmt.Sprintf("goto %s\n", startLabel)
	code += fmt.Sprintf("label %s\n", endLabel)

	return code, nil
}

func (c *CodeGenerator) compileLetStatement(letStatement token.Element) (string, error) {
	var code string

	hasLeftHandArrayAssignment := letStatement.FindChildToken("symbol", "[") != nil

	identifier := letStatement.FindChildToken("identifier")
	if identifier == nil {
		return "", fmt.Errorf("expected identifier")
	}

	var assignmentExpression *token.Element

	if hasLeftHandArrayAssignment {
		assignmentExpression = letStatement.AllChildElementsByTag("expression")[1]
	} else {
		assignmentExpression = letStatement.FindChildElement("expression")
	}

	if assignmentExpression == nil {
		return "", fmt.Errorf("expected assignment expression")
	}

	compiledAssignmentExpression, err := c.compileExpression(assignmentExpression)
	if err != nil {
		return "", err
	}

	symbol := c.findSymbol(identifier.Value)
	if (symbol == Symbol{}) {
		return "", fmt.Errorf("symbol (%s) not found", identifier.Value)
	}

	if hasLeftHandArrayAssignment {

		arrayExpression := letStatement.AllChildElementsByTag("expression")[0]
		compiledArrayExpression, err := c.compileExpression(arrayExpression)
		if err != nil {
			return "", fmt.Errorf("failed compiling array assignment expression: %w", err)
		}

		code += compiledArrayExpression
		code += symbol.Push()
		code += "add\n"
	}

	code += compiledAssignmentExpression

	if hasLeftHandArrayAssignment {
		code += "pop temp 0\n"

		code += "pop pointer 1\n"
		code += "push temp 0\n"
		code += "pop that 0\n"

		return code, nil
	}

	code += symbol.Pop()

	return code, nil
}

func getDoStatement(doStatement token.Element) (string, string, error) {
	hasQualifier := doStatement.FindChildToken("symbol", ".") != nil

	if !hasQualifier {
		subroutineName, err := doStatement.ChildAsToken(1)
		if err != nil {
			return "", "", err
		}

		return "", subroutineName.Value, nil
	}

	qualifier, err := doStatement.ChildAsToken(1)
	if err != nil {
		return "", "", err
	}

	subroutineName, err := doStatement.ChildAsToken(3)
	if err != nil {
		return "", "", err
	}

	return qualifier.Value, subroutineName.Value, nil
}

func (c *CodeGenerator) findSymbol(obj string) Symbol {
	// Find the symbol in the symbol table
	symbol := c.subroutineSymbolTable.Get(obj)
	if (symbol == Symbol{}) {
		symbol = c.classSymbolTable.Get(obj)
	}

	return symbol
}

func (c *CodeGenerator) compileExpressionList(expressionList token.Element) (string, int, error) {
	expressions := expressionList.AllChildElementsByTag("expression")
	if len(expressions) == 0 {
		return "", 0, nil
	}

	var code string

	for _, expression := range expressions {
		compiledExpression, err := c.compileExpression(expression)
		if err != nil {
			return "", 0, err
		}

		code += compiledExpression
	}

	return code, len(expressions), nil
}

func (c *CodeGenerator) compileDoStatement(doStatement token.Element) (string, error) {
	var code string
	var numArgs int

	qualifier, subroutine, err := getDoStatement(doStatement)
	if err != nil {
		return "", err
	}

	instanceVar := c.findSymbol(qualifier)

	if (instanceVar != Symbol{}) {
		code += instanceVar.Push()

		qualifier = instanceVar.Type
		numArgs = 1
	}

	if qualifier == "" {
		code += "push pointer 0\n"

		qualifier = c.className
		numArgs = 1
	}

	expressionList := doStatement.FindChildElement("expressionList")

	compiledExpressionList, argCount, err := c.compileExpressionList(*expressionList)
	if err != nil {
		return "", err
	}

	code += compiledExpressionList

	numArgs += argCount

	code += fmt.Sprintf("call %s.%s %d\n", qualifier, subroutine, numArgs)

	// Do statements don't have a return value, so just dump it
	code += "pop temp 0\n"

	return code, nil
}

func (c *CodeGenerator) compileReturnStatement(returnStatement token.Element) (string, error) {
	var code string

	expression := returnStatement.FindElement("expression")

	if expression == nil {
		// Handle an empty return
		code += "push constant 0\n"
	} else {
		compiledExpression, err := c.compileExpression(expression)
		if err != nil {
			return "", err
		}

		code += compiledExpression
	}

	code += "return\n"

	return code, nil
}

func (c *CodeGenerator) initialiseParameterList(dec *token.Element, symbolTable *SymbolTable, isMethod bool) {
	// Get the parameter list
	parameterList := dec.FindElement("parameterList")

	typesAndIdentifiers := []token.Token{}

	for _, child := range parameterList.AllChildTokens() {
		if child.TokenType == "keyword" || child.TokenType == "identifier" {
			typesAndIdentifiers = append(typesAndIdentifiers, *child)
		}
	}

	if isMethod {
		symbolTable.Add("__placeholder_for_this__", "argument", c.className)
	}

	for i := 0; i < len(typesAndIdentifiers); i += 2 {
		symbolTable.Add(typesAndIdentifiers[i+1].Value, "argument", typesAndIdentifiers[i].Value)
	}
}

func (c *CodeGenerator) compileSubroutine(class, dec *token.Element) (string, error) {
	numLocalVars, err := c.initSymbolTable(
		c.subroutineSymbolTable,
		dec.FindChildElement("subroutineBody").AllChildElementsByTag("varDec"),
	)
	if err != nil {
		return "", fmt.Errorf("error initialising local symbol table for subroutine: %w", err)
	}

	subroutine, err := SubroutineDecFromSyntax(dec)
	if err != nil {
		return "", err
	}

	c.initialiseParameterList(dec, c.subroutineSymbolTable, subroutine.subroutineType == Method)

	c.whileStatementCount = 0
	c.ifStatementCount = 0

	className := class.FindChildToken("identifier").Value
	funcName := subroutine.name

	code := fmt.Sprintf("function %s.%s %d\n", className, funcName, numLocalVars)

	if subroutine.subroutineType == Constructor {
		// If it's a constructor, allocate memory for the object
		code += fmt.Sprintf("push constant %d\n", c.classSymbolTable.Count("field"))
		code += "call Memory.alloc 1\n"
		code += "pop pointer 0\n"
	} else if subroutine.subroutineType == Method {
		// If it's a method, set the first argument to the object
		code += "push argument 0\n"
		code += "pop pointer 0\n"
	}

	statementsContainer := dec.FindElement("statements")
	statements := statementsContainer.AllChildElements()
	for _, statement := range statements {
		compiledStatement, err := c.compileStatement(*statement)
		if err != nil {
			return "", fmt.Errorf("error compiling subroutine (%s): %w", funcName, err)
		}

		code += compiledStatement
	}

	return code, nil
}

func (c *CodeGenerator) compileClass(class *token.Element) (string, error) {
	var code string

	c.className = class.FindChildToken("identifier").Value

	_, err := c.initSymbolTable(
		c.classSymbolTable,
		class.AllChildElementsByTag("classVarDec"),
	)
	if err != nil {
		return "", fmt.Errorf("error initialising class symbol table: %w", err)
	}

	subroutineDecs := class.AllChildElementsByTag("subroutineDec")
	for _, subroutineDec := range subroutineDecs {
		compiledSubroutine, err := c.compileSubroutine(class, subroutineDec)
		if err != nil {
			return "", err
		}
		code += compiledSubroutine
	}

	return code, nil
}

func (c *CodeGenerator) initSymbolTable(s *SymbolTable, vars []*token.Element) (int, error) {
	count := 0
	s.Clear()

	for _, v := range vars {
		// Get the kind of variable
		var kind string
		decLabel, err := v.ChildAsToken(0)
		if err != nil {
			return 0, fmt.Errorf("no valid kind found for variable: %w", err)
		}

		if decLabel.Value == "var" {
			kind = "local"
		} else {
			kind = decLabel.Value
		}

		idents, typ, err := getVarDefs(v)
		if err != nil {
			return 0, err
		}

		for _, i := range idents {
			count++
			s.Add(i, kind, typ)
		}
	}

	return count, nil
}

func (c *CodeGenerator) Generate() (string, error) {
	var code string

	for _, child := range c.code {
		switch child := child.(type) {
		case *token.Element:
			element := child
			switch element.Tag {
			case "class":
				compiledClass, err := c.compileClass(child)
				if err != nil {
					return "", fmt.Errorf("error compiling class: %w", err)
				}

				code += compiledClass
			}
		}
	}

	return code, nil
}
