package parser

import (
	"fmt"
	"language/ast"
	"language/lexer"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	//	gets the error array from the parser
	errors := p.Errors()
	//	checks if there are any
	if len(errors) == 0 {
		return
	}

	//	notifies the amount of errors found
	t.Errorf("Parser has %d errors", len(errors))
	//	loops and prints the errors
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}

func TestConstStatements(t *testing.T) {
	input := `
		const int x = 5;
		const int y = 6;
		const int z = 12412;
	`

	//	creates a lexer object
	l := lexer.New(input)
	//	creates a parser object
	p := New(l)

	//	creates a program
	program := p.ParserProgram()
	checkParserErrors(t, p)

	//	checks if the program is nil
	if program == nil {
		t.Fatalf("ParserProgram returned nil")
	}

	// checks if the input has the expected amount of statements in the program's statements array
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not have the expected amount of statements. Expected 3, got %d", len(program.Statements))
	}

	//	declares the tests to try
	tests := []struct {
		input string
		expectedIdentifier string
		expectedValue interface{}
	} {
		{
			"const int x = 5;", "x", 5,
		},
		{
			"const bool y = true;", "y", true,
		},
		{
			"const bool z = y;", "z", "y",
		},
	}

	//	loops through the tests
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParserProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not have %d elements, received %d",
				1,
				len(program.Statements),
			)
		}

		statement := program.Statements[0]

		if !testConstStatements(t, statement, tt.expectedIdentifier) {
			return
		}

		value := statement.(*ast.ConstStatement).Value

		if !testLiteralExpression(t, value, tt.expectedValue) {
			return
		}
	}
}

func TestVarStatements(t *testing.T) {
	input := `
		var int x = 5;
		var int y = 6;
		var int z = 12412;
	`

	//	creates a lexer object
	l := lexer.New(input)
	//	creates a parser object
	p := New(l)

	//	creates a program
	program := p.ParserProgram()
	checkParserErrors(t, p)

	//	checks if the program is nil
	if program == nil {
		t.Fatalf("ParserProgram returned nil")
	}

	// checks if the input has the expected amount of statements in the program's statements array
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not have the expected amount of statements. Expected 3, got %d", len(program.Statements))
	}

	//	declares the tests to try
	tests := []struct {
		input string
		expectedIdentifier string
		expectedType string
		expectedValue interface{}
	} {
		{
			"var int x = 5;", "x", "int", 5,
		},
		{
			"var bool y = true;", "y", "bool", true,
		},
		{
			"var bool z = y;", "z", "bool", "y",
		},
	}

	//	loops through the tests
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParserProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not have %d elements, received %d",
				1,
				len(program.Statements),
			)
		}

		statement := program.Statements[0]

		if !testVarStatements(t, statement, tt.expectedIdentifier) {
			return
		}

		value := statement.(*ast.VarStatement).Value

		if !testLiteralExpression(t, value, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
		return 23;
		return 3;
		return 44;
	`
	//	create lexer and parser objects
	l := lexer.New(input)
	p := New(l)

	//	creates program object
	program := p.ParserProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParserProgram returned nil")
	}

	// checks if the input has the expected amount of statements in the program's statements array
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not have the expected amount of statements. Expected 3, got %d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("statement not *ast.ReturnStatement, got %T", returnStatement)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf(
				"returnStatement.Token.Literal not 'return', got %s",
				returnStatement.TokenLiteral(),
			)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements. Got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got %T",
			program.Statements[0],
		)
	}

	identifier, ok := statement.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("statement.Expression is not Identifier. Got %T", statement.Expression)
	}

	if identifier.Value != "foobar" {
		t.Errorf("identifier.Value not %s, got %s", "foobar", identifier.Value)
	}

	if identifier.TokenLiteral() != "foobar" {
		t.Errorf("identifier.TokenLiteral not %s, got %s", "foobar", identifier.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements. got %d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statement is not ast.ExpressionStatement. Got %T", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Statement is not ast.IntegerLiteral, got %T", statement.Expression)
	}

	if literal.Value != 5 {
		t.Fatalf("literal.Value not %d, got %d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Fatalf("literal.TokenLiteral not %s, got %s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input string
		operator string
		value interface{}
	}{
		{ "!5", "!", 5 },
		{ "-15", "-", 15 },
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParserProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. Got %d",
				1,
				len(program.Statements),
			)
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf(
				"program.Statements is not of type ExpressionStatement, got %T",
				program.Statements[0],
			)
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("statement is not an ast.PrefixExpression, got %T", statement.Expression)
		}

		if expression.Operator != tt.operator {
			t.Fatalf("expression.Operator is not %s, got %s", tt.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input string
		leftValue interface{}
		operator string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		// {"5 >= 5", 5, ">=", 5},
		// {"5 <= 5", 5, "<=", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	//	loops through the infix tests above
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParserProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d elements, got %d",
				1,
				len(program.Statements),
			)
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		//	checks that the statement is of type Statement
		if !ok {
			t.Fatalf(
				"program.Stataements[0] is not ast.ExpressionStatement, got %T",
				program.Statements[0],
			)
		}

		if !testInfixExpression(t, statement.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}

		expression, ok := statement.Expression.(*ast.InfixExpression)
		//	checks that the expression is of type InfixExpression
		if !ok {
			t.Fatalf(
				"statement.Expression is not of type ast.InfixExpression, got %T",
				statement.Expression,
			)
		}
		//	checks that the integer literal contained in the left side of the expression is correct
		if !testLiteralExpression(t, expression.Left, tt.leftValue) {
			return
		}
		//	checks that the operator is the correct one
		if expression.Operator != tt.operator {
			t.Fatalf("expression.Operator is not %s, got %s", tt.operator, expression.Operator)
		}
		//	checks that the integer literal contained in the right side of the expression is correct
		if !testLiteralExpression(t, expression.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input string
		expected string
	} {
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParserProgram()
		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Fatalf("expected: %q, received: %q", tt.expected, actual)
		}
	}
}

func TestBooleanParsing(t *testing.T) {
	input := "false;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected %d statement(s), got %d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement, got %T",
			program.Statements[0],
		)
	}

	boolean, ok := statement.Expression.(*ast.Boolean)

	if !ok {
		t.Fatalf("statement.Expression is not ast.Boolean, got %T", statement.Expression)
	}

	if boolean.Value != false {
		t.Fatalf("boolean.Value not %t, got %t", true, boolean.Value)
	}

	if boolean.TokenLiteral() != "false" {
		t.Fatalf("boolean.TokenLiteral not %s, got %s", "true", boolean.TokenLiteral())
	}
}

func TestIfStatement(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d elements. Got %d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. Got %T",
			program.Statements[0],
		)
	}

	expression, ok := statement.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("statement.Expression not ast.IfExpression, got %T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement, got %d\n", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"consequence statements not ast.ExpressionStatement, got %T",
			expression.Consequence.Statements[0].(*ast.ExpressionStatement),
		)
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("expression.Alternative.Statements is not nil, got %+v", expression.Alternative)
	}
}

func TestIfElseStatement(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d elements. Got %d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. Got %T",
			program.Statements[0],
		)
	}

	expression, ok := statement.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("statement.Expression not ast.IfExpression, got %T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement, got %d\n", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"consequence statements not ast.ExpressionStatement, got %T",
			expression.Consequence.Statements[0].(*ast.ExpressionStatement),
		)
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
}

func TestFuncExpression(t *testing.T) {
	input := `func(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Statements does not contain %d elements, received %d",
			1,
			len(program.Statements),
		)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement, received %T",
			program.Statements[0],
		)
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("statement.Expression is not FunctionLiteral, got %T", statement.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function parameters wrong, expected 2, got %d", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf(
			"function.Body.Statements has not 1 statement. Got %d",
			len(function.Body.Statements),
		)
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function body statement is not an expression, got %T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestFunctionParameters(t *testing.T) {
	tests := []struct{
		input string
		expectedParams []string
	}{
		{
			input: "func() {}",
			expectedParams: []string{},
		},
		{
			input: "func(x) {}",
			expectedParams: []string{"x"},
		},
		{
			input: "func(x, y, z) {}",
			expectedParams: []string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParserProgram()
		checkParserErrors(t, p)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Fatalf(
				"function.Params returned %d elements, expected %d",
				len(function.Parameters),
				len(tt.expectedParams),
			)
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Statements does not contain %d elements, received %d",
			1,
			len(program.Statements),
		)
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"program.Statements is not of type ast.ExpressionStatement, got %T",
			program.Statements[0],
		)
	}

	expression, ok := statement.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("statement.Expression is not ast.CallExpression, got %T", statement.Expression)
	}

	if !testIdentifier(t, expression.Function, "add") {
		return
	}

	if len(expression.Arguments) != 3 {
		t.Fatalf("got %d arguments, expected 3", len(expression.Arguments))
	}

	testLiteralExpression(t, expression.Arguments[0], 1)
	testInfixExpression(t, expression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expression.Arguments[2], 4, "+", 5)
}

func TestCallArguments(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParserProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func testConstStatements(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "const" {
		t.Errorf("statement.TokenLiteral() not const. got %s", statement.TokenLiteral())
		return false
	}

	constStatement, ok := statement.(*ast.ConstStatement)

	if !ok {
		t.Errorf("statement is not ast.Statement, got %T", statement)
		return false
	}

	if constStatement.Name.Value != name {
		t.Errorf("constStatement.Name.Value not %s. got %s", name, constStatement.Name.Value)
		return false
	}

	if constStatement.Name.TokenLiteral() != name {
		t.Errorf("constStatement.Name not %s. got %s", name, constStatement.Name)
		return false
	}

	return true
}

func testVarStatements(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "var" {
		t.Errorf("statement.TokenLiteral() not var. got %s", statement.TokenLiteral())
		return false
	}

	varStatement, ok := statement.(*ast.VarStatement)

	if !ok {
		t.Errorf("statement is not ast.Statement, got %T", statement)
		return false
	}

	if varStatement.Name.Value != name {
		t.Errorf("varStatement.Name.Value not %s. got %s", name, varStatement.Name.Value)
		return false
	}

	if varStatement.Name.TokenLiteral() != name {
		t.Errorf("varStatement.Name not %s. got %s", name, varStatement.Name)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il is not an *ast.IntegerLiteral, got %T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value is not %d, got %d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() not %d, got %s", value, integer.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("exp no ast.Expression, got %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %q, got %q", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() is not %q, got %q", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)

	if !ok {
		t.Errorf("exp is not *ast.Boolean, got %T", exp)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value is not %t, got %t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral is not %t, got %s", value, boolean.TokenLiteral())
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled, got %T", exp)
	return false
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	operatorExp, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp not ast.InfixExpression, got %T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, operatorExp.Left, left) {
		return false
	}

	if operatorExp.Operator != operator {
		t.Errorf("exp.Operator not %q, got %q", operator, operatorExp.Operator)
		return false
	}

	if !testLiteralExpression(t, operatorExp.Right, right) {
		return false
	}

	return true
}
