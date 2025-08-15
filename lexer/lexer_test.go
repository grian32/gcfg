package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
[]
()
{}
=,
1.23
123
"hello"
	`

	expectedTokenTypes := []Token{
		newToken(LBRACKET, "["),
		newToken(RBRACKET, "]"),
		newToken(LPAREN, "("),
		newToken(RPAREN, ")"),
		newToken(LBRACE, "{"),
		newToken(RBRACE, "}"),
		newToken(ASSIGN, "="),
		newToken(COMMA, ","),
		newToken(FLOAT, "1.23"),
		newToken(INT, "123"),
		newToken(STRING, "hello"),
		newToken(BOOL, "true"),
		newToken(BOOL, "false"),
		newToken(IDENT, "foo"),
		newToken(NULL, "nil"),
		newToken(EOF, ""),
	}

	l := New([]byte(input))

	for _, tt := range expectedTokenTypes {
		token, err := l.NextToken()

		if token != tt || err != nil {
			t.Errorf("NextToken=%v, %v, wanted match for %v", token, err, tt)
		}
	}
}

func newToken(tokenType TokenType, lit string) Token {
	return Token{Type: tokenType, Literal: lit}
}

func TestBadInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "MultipleDots",
			input: "1.2.3",
		},
		{
			name:  "EndInDot",
			input: "123.",
		},
		{
			name:  "MalformedString",
			input: `"hey`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New([]byte(tt.input))
			_, err := l.NextToken()

			if err == nil {
				t.Errorf("NextToken expected error but got nil")
			}
		})
	}
}
