package lexer

type Lexer struct {
	input   []byte
	readPos int
	prevPos int
	ch      byte
}

var singleCharTokens = map[byte]TokenType{
	'[': LBRACKET,
	']': RBRACKET,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	'=': ASSIGN,
	',': COMMA,
}

func New(input []byte) *Lexer {
	l := &Lexer{input: input}
	l.advance()
	return l
}

func (l *Lexer) advance() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}

	l.prevPos = l.readPos
	l.readPos += 1
}

func (l *Lexer) peek() byte {
	if l.readPos >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPos]
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	// TODO: consider not doing this.. maybe just for indentation
	l.skipWhitespace()

	singleTok, exists := singleCharTokens[l.ch]
	if exists {
		tok = newSingleToken(singleTok, l.ch)
		l.advance()
		return tok
	}

	if l.ch == 0 {
		tok.Literal = ""
		tok.Type = EOF
	} else if l.ch == '"' {
		// strings
	} else if IsDigit(l.ch) {
		// ints and floats
	} else {
		// identifiers, booleans, null
	}

	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.advance()
	}
}

func newSingleToken(tokType TokenType, ch byte) Token {
	return Token{Type: tokType, Literal: string(ch)}
}
