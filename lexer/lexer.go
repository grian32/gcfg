package lexer

type Lexer struct {
	input   []byte
	readPos int
	prevPos int
	ch      byte
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
	return Token{}
}
