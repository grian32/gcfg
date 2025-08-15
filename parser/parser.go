package parser

import (
	"errors"
	"gcfg/lexer"
	"strconv"
)

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
			}
		}

		err = p.NextToken()
		if err != nil {
			return nil, err
		}
	}

	return fileMap, nil
}

func (p *Parser) parseAssign() (any, error) {
	err := p.NextToken()
	if err != nil {
		return nil, err
	}

	var val any

	switch p.peekToken.Type {
	case lexer.INT:
		value, err := strconv.Atoi(p.peekToken.Literal)
		if err != nil {
			return nil, err
		}
		val = value
	case lexer.FLOAT:
		value, err := strconv.ParseFloat(p.peekToken.Literal, 64)
		if err != nil {
			return nil, err
		}
		val = value
	case lexer.STRING:
		val = p.peekToken.Literal
	case lexer.TRUE:
		val = true
	case lexer.FALSE:
		val = false
	case lexer.NULL:
		val = nil
	default:
		return nil, errors.New("")
	}

	err = p.NextToken()
	if err != nil {
		return nil, err
	}

	return val, nil
}
