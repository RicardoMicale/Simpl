package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF" //	End Of File

	//	Identifiers and literals
	IDENTIFIER = "IDENTIFIER" //	variables
	INT = "INT" //	Integer data type
	STRING = "STRING" //	String data type
	DOUBLE = "DOUBLE" //	Double data type
	BOOL = "BOOL" //	Boolean data type
	ARRAY = "ARRAY" //	Array data type

	//	Operators
	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"
	MULTIPLY = "*"
	DIVIDE = "/"
	MODULO = "%"
	EXACT_DIVISION = "//"
	POWER = "**"

	//	Logical operators
	AND = "&&"
	OR = "||"
	NOT = "!"
	EQUALS = "=="
	NOT_EQUALS = "!="
	LESS_THAN = "<"
	LESS_THAN_OR_EQUAL = "<="
	GREATER_THAN = ">"
	GREATER_THAN_OR_EQUAL = ">="

	//	Delimiters
	COMMA = ","
	SEMICOLON = ";"
	COLON = ":"

	L_PAREN = "("
	R_PAREN = ")"
	L_BRACE = "{"
	R_BRACE = "}"
	L_BRACK = "["
	R_BRACK = "]"

	//	Keywords
	FUNCTION = "FUNC"
	VAR = "VAR"
	CONST = "CONST"
	RETURN = "RETURN"
)

var keywords = map[string]TokenType {
	"func": FUNCTION,
	"const": CONST,
	"var": VAR,
	"return": RETURN,
	"int": INT,
	"string": STRING,
	"double": DOUBLE,
	"bool": BOOL,
	"array": ARRAY,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENTIFIER
}
