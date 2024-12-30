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
		{
			input: "5",
			expected: 5,
		},
		{
			input: "10",
			expected: 10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {

}

func TestEvalNullExpression(t *testing.T) {

}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParserProgram()

	return Eval(program)
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
