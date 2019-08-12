package token

const (
	ILLEGAL Type = "ILLEGAL"
	EOF     Type = "EOF"

	// Identifiers + literals
	IDENT Type = "IDENT"
	INT   Type = "INT"

	// Operators
	ASSIGN   Type = "="
	PLUS     Type = "+"
	MINUS    Type = "-"
	BANG     Type = "!"
	ASTERISK Type = "*"
	SLASH    Type = "/"
	LT       Type = "<"
	GT       Type = ">"
	EQ       Type = "=="
	NOT_EQ   Type = "!="

	// Delimiters
	COMMA     Type = ","
	SEMICOLON Type = ";"

	LPAREN Type = "("
	RPAREN Type = ")"
	LBRACE Type = "{"
	RBRACE Type = "}"

	// Keywords
	FUNCTION Type = "FUNCTION"
	LET      Type = "LET"
	TRUE     Type = "TRUE"
	FALSE    Type = "FALSE"
	IF       Type = "IF"
	ELSE     Type = "ELSE"
	RETURN   Type = "RETURN"
)

var (
	SingleByteLiteralToType = map[byte]Type{
		0:   EOF,
		'+': PLUS,
		',': COMMA,
		';': SEMICOLON,
		'(': LPAREN,
		')': RPAREN,
		'{': LBRACE,
		'}': RBRACE,
		'-': MINUS,
		'*': ASTERISK,
		'/': SLASH,
		'<': LT,
		'>': GT,
	}

	Keywords = map[string]Type{
		"fn":     FUNCTION,
		"let":    LET,
		"true":   TRUE,
		"false":  FALSE,
		"if":     IF,
		"else":   ELSE,
		"return": RETURN,
	}
)

type (
	Type string

	Token struct {
		Type    Type
		Literal string
	}
)

func (t Token) TokenLiteral() string {
	return t.Literal
}

func LookupIdent(ident string) Type {
	if t, ok := Keywords[ident]; ok {
		return t
	}

	return IDENT
}

func (t Type) Token(literal string) Token {
	return Token{Type: t, Literal: literal}
}
