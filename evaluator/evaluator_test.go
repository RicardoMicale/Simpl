package evaluator

import (
	"language/lexer"
	"language/object"
	"language/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct{
		input string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5",-5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct{
		input string
		expected bool
	}{
		{ "true", true },
		{ "false", false },
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{`"Hello" == "Hello"`, true},
		{`"Hello" != "World"`, true},
		{`"Hello" == "World"`, false},
		{`"Hello" != "Hello"`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestNotOperator(t *testing.T) {
	tests := []struct{
		input string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct{
		input string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct{
		input string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
				if (10 > 1) {
					if (10 > 1) {
						return 10;
					}

					return 1;
				}
			`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct{
		input string
		expected string
	}{
		{
			"5 + true;",
			"Type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"Type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"Unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			132
			if (10 > 1) {
			if (10 > 1) {
			return true + false;
			}
			return 1;
			}
			`,
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"Identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"Unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("No error object returned. Got %T (+%v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf("Wrong error message.\nExpected: %q\nGot: %q", tt.expected, errObj.Message)
		}
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct{
		input string
		expected int64
	}{
		{"var int a = 5; a;", 5},
		{"var int a = 5 * 5; a;", 25},
		{"var int a = 5; var int b = a;", 5},
		{"var int a = 5; var int b = a; var int c = a + b + 5; c", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestConstStatements(t *testing.T) {
	tests := []struct{
		input string
		expected int64
	}{
		{"const int a = 5; a;", 5},
		{"const int a = 5 * 5; a;", 25},
		{"const int a = 5; const int b = a;", 5},
		{"const int a = 5; const int b = a; const int c = a + b + 5; c", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not a Function, got %T", evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters (%+v)", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not x, got %q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q, got %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct{
		input string
		expected int64
	}{
		{"const int identity = func(x) { x; }; identity(5);", 5},
		{"const int identity = func(x) { return x; }; identity(5);", 5},
		{"const int multiply = func(x) { x * 2; }; multiply(5);", 10},
		{"const int add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"const int add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
		const int adder = func(x) {
			func(y) { x + y; };
		};

		const int addTwo = adder(2);
		addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello world!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not of type String. got %T", evaluated)
	}

	if str.Value != "Hello world!" {
		t.Errorf("String has wrong value, got %q", str.Value)
	}
}

func TestBuiltInFunctions(t *testing.T) {
	tests := []struct{
		input string
		expected interface{}
	}{
		{`length("")`, 0},
		{`length("four")`, 4},
		{`length("Hello World")`, 11},
		{`length(1)`, "Argument to `length` not supported, got INTEGER"},
		{`length("one", "two")`, "Wrong number of arguments, expected 1, got 2"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)

			if !ok {
				t.Errorf("Object is not error, got %T, (%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Message != expected {
				t.Errorf("Wrong error message, Expected %q, got %q", expected, errObj.Message)
			}
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParserProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not an Integer, got %T", obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Received result %d, expected %d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("obj not of type Boolean, got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("result.Value was %t, expected %t", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL, got %T (%+v)", obj, obj)
		return false
	}

	return true
}
