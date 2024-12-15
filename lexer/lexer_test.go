package lexer

import (
	"testing"

	"language/token"
)

type Lexer struct {
	input string
	position int //	current position in input, points t current char
	readPosition int //	current reading position in input, after current char
	ch byte //	current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{ input: input }
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.MULTIPLY, l.ch)
	case '/':
		tok = newToken(token.DIVIDE, l.ch)
	case '%':
		tok = newToken(token.MODULO, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.L_PAREN, l.ch)
	case ')':
		tok = newToken(token.R_PAREN, l.ch)
	case '{':
		tok = newToken(token.L_BRACE, l.ch)
	case '}':
		tok = newToken(token.R_BRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func TestNextToken(t *testing.T) {
	input := `
	const int multiplier = 3

	func add(int x, int y) {
		return (x + y) * multiplier
	}

	add(2, 4)
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
		{token.IDENTIFIER, "add"},
		{token.L_PAREN, "("},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "4"},
		{token.R_PAREN, ")"},
		// {token.MINUS, "-"},
		// {token.MULTIPLY, "*"},
		// {token.DIVIDE, "/"},
		// {token.MODULO, "%"},
		// {token.EXACT_DIVISION, "//"},
		// {token.POWER, "**"},
		// {token.L_PAREN, "("},
		// {token.R_PAREN, ")"},
		// {token.L_BRACE, "{"},
		// {token.R_BRACE, "}"},
		// {token.COMMA, ","},
		// {token.SEMICOLON, ";"},
		// {token.COLON, ":"},
		// {token.AND, "&&"},
		// {token.OR, "||"},
		// {token.EQUALS, "=="},
		// {token.NOT_EQUALS, "!="},
		// {token.NOT, "!"},
		// {token.GREATER_THAN, ">"},
		// {token.GREATER_THAN_OR_EQUAL, ">="},
		// {token.LESS_THAN, "<"},
		// {token.LESS_THAN_OR_EQUAL, "<="},
		// {token.RETURN, "return"},
		// {token.ARRAY, "array"},
		// {token.STRING, "string"},
		// {token.BOOL, "bool"},
		// {token.DOUBLE, "double"},
		// {token.VAR, "var"},
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
