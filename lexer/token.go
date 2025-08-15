package lexer

//go:generate stringer -type=TokenType
type TokenType byte

const (
	LBRACKET TokenType = iota
	RBRACKET

	LPAREN
	RPAREN

	LBRACE
	RBRACE

	ASSIGN
	COMMA

	IDENT
	INT
	FLOAT
	STRING
	TRUE
	FALSE
	NULL

	EOF
)

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	return t.Type.String() + "(" + t.Literal + ")"
}
