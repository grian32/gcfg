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
			} else if p.peekToken.Type == lexer.LBRACE {
				name := p.curToken.Literal
				value, err := p.parseSection()
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

func (p *Parser) parseSection() (map[string]any, error) {
	err := p.NextToken() // advance past lbrace
	if err != nil {
		return nil, err
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

func (p *Parser) parseValue() (any, error) {
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
		return nil, errors.New("non accepted value")
	}

	return val, nil
}
