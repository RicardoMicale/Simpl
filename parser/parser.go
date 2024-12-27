package parser

import (
	"fmt"
	"language/ast"
	"language/lexer"
	"language/token"
	"strconv"
)

//	This determines the order of priority of an operation
//	going from the lower priority (iota) to the highest priority (function calls)
const (
	_ int = iota
	LOWEST
	EQUALS //	==
	LESS_GREATER //	< or >
	SUM //	+
	PRODUCT //	*
	PREFIX //	-X or !X
	CALL //	myFunc(X)
)

var precedences = map[token.TokenType]int{
	token.EQUALS: EQUALS,
	token.NOT_EQUALS: EQUALS,
	token.LESS_THAN: LESS_GREATER,
	token.GREATER_THAN: LESS_GREATER,
	token.PLUS: SUM,
	token.MINUS: SUM,
	token.DIVIDE: PRODUCT,
	token.MULTIPLY: PRODUCT,
}

type (
	prefixParseFunc func() ast.Expression
	infixParseFunc func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer
	errors []string

	currentToken token.Token //	current token being read
	peekToken token.Token //	next token being peeked

	prefixParseFuncs map[token.TokenType]prefixParseFunc
	infixParseFuncs map[token.TokenType]infixParseFunc
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{ l: l, errors: []string{} }

	//	read two tokens to set both currentToken and peerToken
	p.nextToken()
	p.nextToken()

	//	makes prefix functions
	p.prefixParseFuncs = make(map[token.TokenType]prefixParseFunc)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.L_PAREN, p.parseGroupedExpression)

	//	makes infix functions
	p.infixParseFuncs = make(map[token.TokenType]infixParseFunc)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.DIVIDE, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLY, p.parseInfixExpression)
	p.registerInfix(token.EQUALS, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUALS, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN, p.parseInfixExpression)

	//	returns the parser
	return p
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.R_PAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token: p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token: p.currentToken,
		Operator: p.currentToken.Literal,
		Left: left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{ Token: p.currentToken, Value: p.currentTokenIs(token.TRUE) }
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{ Token: p.currentToken, Value: p.currentToken.Literal }
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	//	the next token is the one being peeked at
	p.currentToken = p.peekToken
	//	then the peeked at token is the next one
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	//	prints and stores an error message
	message := fmt.Sprintf("Expected token to be %s but received %s", t, p.peekToken.Type)
	//	adds it to the error array of the parser
	p.errors = append(p.errors, message)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseConstStatement() *ast.ConstStatement {
	//	creates the statement object and assigns its memory address to a variable
	statement := &ast.ConstStatement{ Token: p.currentToken}
	//	helper variable with data type tokens
	dataTypes := []token.TokenType{token.INT, token.STRING, token.DOUBLE, token.BOOL}
	//	used to flag if the data type is not specified
	dataTypeFlag := true
	//	checks if the next token is a data type token
	for _, dataType := range dataTypes  {
		if p.expectPeek(dataType) {
			dataTypeFlag = false
			break //	once a valid data type is found, break the loop
		}
	}
	//	if the flag is true, it did not find any data type specified
	//	therefore it is not a valid statement
	if dataTypeFlag {
		return nil
	}

	//	checks if the next token is not a variable name
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	//	defines the Name attribute of the statement as an identifier
	statement.Name = &ast.Identifier{ Token: p.currentToken, Value: p.currentToken.Literal }

	//	checks if the next token after the variable name is an assign token
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//	TODO Skipping expressions until we find a semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	//	creates the statement object and assigns its memory address to a variable
	statement := &ast.VarStatement{ Token: p.currentToken}
	//	helper variable with data type tokens
	dataTypes := []token.TokenType{token.INT, token.STRING, token.DOUBLE, token.BOOL}
	//	used to flag if the data type is not specified
	dataTypeFlag := true
	//	checks if the next token is a data type token
	for _, dataType := range dataTypes  {
		if p.expectPeek(dataType) {
			dataTypeFlag = false
			break //	once a valid data type is found, break the loop
		}
	}
	//	if the flag is true, it did not find any data type specified
	//	therefore it is not a valid statement
	if dataTypeFlag {
		return nil
	}

	//	checks if the next token is not a variable name
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	//	defines the Name attribute of the statement as an identifier
	statement.Name = &ast.Identifier{ Token: p.currentToken, Value: p.currentToken.Literal }

	//	checks if the next token after the variable name is an assign token
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//	TODO Skipping expressions until we find a semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	//	creates a statement variable with the current return token
	statement := &ast.ReturnStatement{ Token: p.currentToken }

	//	goes to the next token
	p.nextToken()

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) noPrefixParseError(t token.TokenType) {
	message := fmt.Sprintf("No prefix function for %s found", t)
	p.errors = append(p.errors, message)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	//	gets the function corresponding to the token type being parsed
	prefix := p.prefixParseFuncs[p.currentToken.Type]

	if prefix == nil {
		p.noPrefixParseError(p.currentToken.Type)
		return nil
	}
	//	calls the prefixFunc found
	leftExpression := prefix()

	//	loops the expression until it finds a semicolon
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		//	creates an infix variable assigned to the corresponding next token function on the infix callbacks
		infix := p.infixParseFuncs[p.peekToken.Type]

		//	if there is no infix, returns the left expression
		if infix == nil {
			return leftExpression
		}

		//	advances the parser
		p.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{ Token: p.currentToken }

	statement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	//	assigns the literal variable to an IntegerLiteral
	literal := &ast.IntegerLiteral{ Token: p.currentToken }
	//	uses the string converter library to parse the literal from a string to an integer
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)

	if err != nil {
		//	creates an error message and appends it to the parser error list
		message := fmt.Sprintf("Could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, message)
		return nil
	}

	//	else it assigns the parsed integer value to the literal aValue attribute and returns the literal
	literal.Value = value

	return literal
}

func (p *Parser) parseStatement() ast.Statement {
	//	switches depending on the token type and parses that specific type
	switch p.currentToken.Type {
	case token.CONST:
		return p.parseConstStatement()
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFunc) {
	//	registers the prefix parser function a the token type
	p.prefixParseFuncs[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFunc) {
	//	registers the infix parser function a the token type
	p.infixParseFuncs[tokenType] = fn
}

func (p *Parser) ParserProgram() *ast.Program {
	//	create a program object
	program := &ast.Program{}
	//	initialize the statements array with an empty one
	program.Statements = []ast.Statement{}

	//	loop through the tokens while the current token is not an EOF type
	for !p.currentTokenIs(token.EOF) {
		//	gets the parsed statement
		statement := p.parseStatement()

		//	if the statement is not nil, append it to the statement array
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		//	goes to the next token of the input
		p.nextToken()
	}

	return program
}
