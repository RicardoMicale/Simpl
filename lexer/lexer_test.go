package lexer

import (
	"language/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
	const int multiplier = 3

	func int add(int x, int y) {
		return (x + y) * multiplier
	}

	const int sum = add(2, 4)

	if (sum > 4) {
		return true
	} else {
		return false
	}

	10 == 9
	10 != 9
	10 >= 9
	10 <= 9

	"foobar"
	"foo bar"
	[1, 2];
	{"foo": "bar"};
	for (i < 10) {
		x = x + i
	};
	/#	comments #/
	multiplier = 4;
	`

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	}{
		{token.CONST, "const"},
		{token.INT, "int"},
		{token.IDENTIFIER, "multiplier"},
		{token.ASSIGN, "="},
		{token.INT, "3"},
		{token.FUNCTION, "func"},
		{token.INT, "int"},
		{token.IDENTIFIER, "add"},
		{token.L_PAREN, "("},
		{token.INT, "int"},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.INT, "int"},
		{token.IDENTIFIER, "y"},
		{token.R_PAREN, ")"},
		{token.L_BRACE, "{"},
		{token.RETURN, "return"},
		{token.L_PAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.R_PAREN, ")"},
		{token.MULTIPLY, "*"},
		{token.IDENTIFIER, "multiplier"},
		{token.R_BRACE, "}"},
		{token.CONST, "const"},
		{token.INT, "int"},
		{token.IDENTIFIER, "sum"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.L_PAREN, "("},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "4"},
		{token.R_PAREN, ")"},
		{token.IF, "if"},
		{token.L_PAREN, "("},
		{token.IDENTIFIER, "sum"},
		{token.GREATER_THAN, ">"},
		{token.INT, "4"},
		{token.R_PAREN, ")"},
		{token.L_BRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.R_BRACE, "}"},
		{token.ELSE, "else"},
		{token.L_BRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.R_BRACE, "}"},
		{token.INT, "10"},
		{token.EQUALS, "=="},
		{token.INT, "9"},
		{token.INT, "10"},
		{token.NOT_EQUALS, "!="},
		{token.INT, "9"},
		{token.INT, "10"},
		{token.GREATER_THAN_OR_EQUAL, ">="},
		{token.INT, "9"},
		{token.INT, "10"},
		{token.LESS_THAN_OR_EQUAL, "<="},
		{token.INT, "9"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.L_BRACK, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.R_BRACK, "]"},
		{token.SEMICOLON, ";"},
		{token.L_BRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.R_BRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.FOR, "for"},
		{token.L_PAREN, "("},
		{token.IDENTIFIER, "i"},
		{token.LESS_THAN, "<"},
		{token.INT, "10"},
		{token.R_PAREN, ")"},
		{token.L_BRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "i"},
		{token.R_BRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.DIVIDE, "/"},
		{token.IDENTIFIER, "multiplier"},
		{token.ASSIGN, "="},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. Expected=%q, got=%q at=%q", i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. Expected=%q, got=%q at=%q", i, tt.expectedLiteral, tok.Literal, tok.Literal)
		}
	}
}
