package parser

import (
	"errors"
	"gcfg"
	"gcfg/lexer"
	"strconv"
)

var ErrNotSimple = errors.New("value is not simple")

type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{
		l: l,
	}
}

func (p *Parser) NextToken() error {
	tok, err := p.l.NextToken()
	if err != nil {
		return err
	}

	p.curToken = p.peekToken
	p.peekToken = tok
	return nil
}

func (p *Parser) ParseFile() (map[string]any, error) {
	err := p.NextToken()
	if err != nil {
		return nil, err
	}
	err = p.NextToken()
	if err != nil {
		return nil, err
	}

	fileMap := make(map[string]any)

	for p.peekToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.IDENT {
			if p.peekToken.Type == lexer.ASSIGN {
				name := p.curToken.Literal
				value, err := p.parseAssign()
				if err != nil {
					return nil, err
				}

				fileMap[name] = value
			} else if p.peekToken.Type == lexer.LBRACE {
				name := p.curToken.Literal
				value, err := p.parseSection(false)
				if err != nil {
					return nil, err
				}

				fileMap[name] = value
			}
		} else if p.curToken.Type == lexer.LBRACKET && p.peekToken.Type == lexer.IDENT {
			err = p.NextToken() // advance past lbracket
			if err != nil {
				return nil, err
			}
			name := p.curToken.Literal
			value, err := p.parseSection(true)
			if err != nil {
				return nil, err
			}

			arr, exists := fileMap[name]

			if exists {
				fileMap[name] = append(arr.([]map[string]any), value)
			} else {
				fileMap[name] = []map[string]any{value}
			}
		}

		err = p.NextToken()
		if err != nil {
			return nil, err
		}
	}

	return fileMap, nil
}

func (p *Parser) parseSection(arrSection bool) (map[string]any, error) {
	err := p.NextToken() // advance past lbrace
	if err != nil {
		return nil, err
	}

	if arrSection {
		if p.curToken.Type != lexer.RBRACKET {
			return nil, errors.New("expected closing ] for array section")
		}
		err := p.NextToken() // advance past ]
		if err != nil {
			return nil, err
		}
	}

	err = p.NextToken() // advance to first
	if err != nil {
		return nil, err
	}

	sectionMap := make(map[string]any)

	for p.curToken.Type != lexer.RBRACE {
		if p.curToken.Type == lexer.IDENT && p.peekToken.Type == lexer.ASSIGN {
			name := p.curToken.Literal
			value, err := p.parseAssign()
			if err != nil {
				return nil, err
			}

			sectionMap[name] = value
		} else {
			return nil, errors.New("something other than assignments found in section")
		}

		err = p.NextToken()
		if err != nil {
			return nil, err
		}
	}

	return sectionMap, nil
}

func (p *Parser) parseAssign() (any, error) {
	err := p.NextToken() // advance past =
	if err != nil {
		return nil, err
	}
	err = p.NextToken() // advance to value
	if err != nil {
		return nil, err
	}

	return p.parseValue()
}

func (p *Parser) parseSimpleValue() (any, error) {
	var val any

	switch p.curToken.Type {
	case lexer.INT:
		value, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			return nil, err
		}
		val = value
	case lexer.FLOAT:
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			return nil, err
		}
		val = value
	case lexer.STRING:
		val = p.curToken.Literal
	case lexer.TRUE:
		val = true
	case lexer.FALSE:
		val = false
	case lexer.NULL:
		val = nil
	default:
		return nil, ErrNotSimple
	}

	return val, nil
}

func (p *Parser) parsePair() (any, error) {
	err := p.NextToken() // advance past lparen
	if err != nil {
		return nil, err
	}

	first, err := p.parseSimpleValue()
	if err != nil {
		return nil, err
	}

	err = p.NextToken()
	if p.curToken.Type != lexer.COMMA {
		return nil, errors.New("expected comma after value in pair")
	}
	if err != nil {
		return nil, err
	}

	err = p.NextToken() // advance past comma
	if err != nil {
		return nil, err
	}

	second, err := p.parseSimpleValue()
	if err != nil {
		return nil, err
	}

	err = p.NextToken()
	if err != nil {
		return nil, err
	}
	if p.curToken.Type != lexer.RPAREN {
		return nil, errors.New("expected rparen after second value in pair")
	}

	return gcfg.Pair[any, any]{
		First:  first,
		Second: second,
	}, nil
}

func (p *Parser) parseArray() (any, error) {
	err := p.NextToken() // advance past lbracket
	if err != nil {
		return nil, err
	}

	first, err := p.parseSimpleValue()
	if err != nil {
		return nil, err
	}
	firstType := p.curToken.Type

	err = p.NextToken()
	if p.curToken.Type != lexer.COMMA {
		return nil, errors.New("expected comma after value in array")
	}
	if err != nil {
		return nil, err
	}

	arr := []any{first}

	for p.curToken.Type != lexer.RBRACKET {
		err = p.NextToken()
		if err != nil {
			return nil, err
		}

		val, err := p.parseSimpleValue()
		if p.curToken.Type != lexer.COMMA {
			return nil, errors.New("expected comma after value in array")
		}
		if p.curToken.Type != firstType {
			return nil, errors.New("arrays must be of single type")
		}
		if err != nil {
			return nil, err
		}
		arr = append(arr, val)

		err = p.NextToken()
		if err != nil {
			return nil, err
		}
	}

	return arr, nil
}

func (p *Parser) parseValue() (any, error) {
	simple, err := p.parseSimpleValue()
	if err != nil && !errors.Is(err, ErrNotSimple) {
		return nil, err
	}

	if errors.Is(err, ErrNotSimple) {
		switch p.curToken.Type {
		case lexer.LPAREN:
			return p.parsePair()
		case lexer.LBRACKET:
			return p.parseArray()
		default:
			return nil, errors.New("invalid value")
		}
	} else {
		return simple, nil
	}
}
