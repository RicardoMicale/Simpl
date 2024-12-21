package lexer

import (
	"language/token"
)

type Lexer struct {
	input string
	position int //	current position in input, points t current char
	readPosition int //	current reading position in input, after current char
	ch byte //	current char under examination
}

func New(input string) *Lexer {
	//	assigns the Lexer memory address value to the l variable
	l := &Lexer{ input: input }
	//	Reads the first character of the input
	l.readChar()
	//	returns the new Lexer value
	return l
}

func (l *Lexer) readChar() {
	//	If the position of the next character is the end or after the end of the input
	if l.readPosition >= len(l.input) {
		//	set the character under examination as 0
		l.ch = 0
	} else {
		//	else it re-assigns the character to the next character in the input
		l.ch = l.input[l.readPosition]
	}
	//	advances the current position stored in the Lexer to the next one
	l.position = l.readPosition
	//	advances the next position stored in the Lexer to the next one
	l.readPosition += 1
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	//	creates a new token using a token type defined in the token folder and the character parsed from byte to string would be the literal
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	//	checks the byte value of the character and compares it to the byte values to determine if its a letter
	//	it uses the ascii values as the byte value
	//	incorporates the underscore (_)
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIdentifier() string {
	//	assigns the current position to a variable to reference it later
	position := l.position
	for isLetter(l.ch) {
		//	acts like a while loop
		//	while the character is a letter, read the character
		l.readChar()
	}

	//	returns the input on the next non-letter position
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	//	checks the byte value of the character being read and compares it to the different whitespace options
	//	as long as there is whitespace, it will read the character and do nothing but advance the input position pointers
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	//	checks the byte value and compares it in an interval of 0 to 9 to check if the character is a number
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	//	assigns the current position to a value to use it later on
	position := l.position

	//	works like a while loop
	for isDigit(l.ch) {
		//	as long as the character is a number it reads the character and advances the position
		l.readChar()
	}

	//	returns the next non-number position
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	//	checks if the next position is at the end or after the end of the input
	if l.readPosition >= len(l.input) {
		//	if it is, return 0
		return 0
	} else {
		//	else it returns the next character
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	//	creates a token variable uninitialized
	var tok token.Token

	//	skips the whitespace
	l.skipWhitespace()

	//	switches according to the character encountered
	switch l.ch {
	case '=':
		//	check for equals (==)
		if l.peekChar() == '=' {
			//	if the next character is an equals sign (=)
			//	stores the character in a variable to access it later
			ch := l.ch
			//	reads the character and advances position
			l.readChar()
			//	assign an EQUALS token using the previously stored character and the current character
			tok = token.Token{ Type: token.EQUALS, Literal: string(ch) + string(l.ch) }
		} else {
			//	creates the token using the appropriate token type and the current character
			tok = newToken(token.ASSIGN, l.ch)
		}
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
	case '>':
		//	checks for greater than or equal (>=)
		if l.peekChar() == '=' {
			//	if the next character is an equals sign (=)
			//	stores the character in a variable to access it later
			ch := l.ch
			//	reads the character and advances position
			l.readChar()
			//	assign a GREATER_THAN_OR_EQUAL token using the previously stored character and the current character
			tok = token.Token{ Type: token.GREATER_THAN_OR_EQUAL, Literal: string(ch) + string(l.ch) }
		} else {
			tok = newToken(token.GREATER_THAN, l.ch)
		}
	case '<':
		//	checks for less than or equal (<=)
		if l.peekChar() == '=' {
			//	if the next character is an equals sign (=)
			//	stores the character in a variable to access it later
			ch := l.ch
			//	reads the character and advances position
			l.readChar()
			//	assign a LESS_THAN_OR_EQUAL token using the previously stored character and the current character
			tok = token.Token{ Type: token.LESS_THAN_OR_EQUAL, Literal: string(ch) + string(l.ch) }
		} else {
			tok = newToken(token.LESS_THAN, l.ch)
		}
	case '!':
		//	checks for a not equals (!=)
		if l.peekChar() == '=' {
			//	if the next character is an equals sign (=)
			//	stores the character in a variable to access it later
			ch := l.ch
			//	reads the character and advances position
			l.readChar()
			//	assign a NOT_EQUALS token using the previously stored character and the current character

			tok = token.Token{ Type: token.NOT_EQUALS, Literal: string(ch) + string(l.ch) }
		} else {
			tok = newToken(token.NOT, l.ch)
		}
	case '[':
		tok = newToken(token.L_BRACK, l.ch)
	case ']':
		tok = newToken(token.R_BRACK, l.ch)
	case 0:
		//	when a zero is found, it means it is the end of the file
		//	the token literal is an empty string and the type is an End Of File (EOF)
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		//	if the token is a letter
		if isLetter(l.ch) {
			//	the literal is the identifier (variables and keywords)
			tok.Literal = l.readIdentifier()
			//	the type is looked in the LookupIdent function of the token file
			tok.Type = token.LookupIdent(tok.Literal)
			//	returns the token
			return tok
		} else if isDigit(l.ch) {
			//	if the token is a number it assigns an INT token type
			//	TODO CHECK FOR DOUBLES
			tok.Type = token.INT
			//	reads the number and assigns it as the literal
			tok.Literal = l.readNumber()
			//	returns the token
			return tok
		} else {
			//	if it is not a number or a letter, it is classified as an illegal type
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	//	advances position and return the token
	l.readChar()
	return tok
}
