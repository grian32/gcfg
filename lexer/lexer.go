package lexer

import "errors"

type Lexer struct {
	input   []byte
	readPos int
	pos     int
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

	l.pos = l.readPos
	l.readPos += 1
}

func (l *Lexer) peek() byte {
	if l.readPos >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPos]
	}
}

func (l *Lexer) NextToken() (Token, error) {
	// TODO: consider not doing this.. maybe just for indentation
	l.skipWhitespace()

	singleTok, exists := singleCharTokens[l.ch]
	if exists {
		tok := newSingleToken(singleTok, l.ch)
		l.advance()
		return tok, nil
	}

	if l.ch == 0 {
		l.advance()
		return Token{Type: EOF, Literal: ""}, nil
	} else if l.ch == '"' {
		return l.readString()
	} else if IsDigit(l.ch) {
		return l.readNumber()
	} else {
		//return l.readIdent()
		return Token{}, nil
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.advance()
	}
}

func (l *Lexer) readString() (Token, error) {
	startPos := l.pos + 1

	l.advance() // go past "

	for l.ch != '"' {
		l.advance()
	}

	l.advance()

	return Token{Type: STRING, Literal: string(l.input[startPos : l.pos-1])}, nil
}

func (l *Lexer) readNumber() (Token, error) {
	startPos := l.pos
	tokType := INT
	float := false

	for IsDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			if float {
				return Token{}, errors.New("multiple dots not allowed in number")
			}
			float = true
			tokType = FLOAT
		}
		l.advance()
	}

	literal := string(l.input[startPos:l.pos])

	if literal[len(literal)-1] == '.' {
		return Token{}, errors.New("numbers not allowed to end in dot")
	}

	l.advance()

	return Token{Type: tokType, Literal: literal}, nil
}

func newSingleToken(tokType TokenType, ch byte) Token {
	return Token{Type: tokType, Literal: string(ch)}
}
