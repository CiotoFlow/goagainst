package trollan

import "errors"

type Parser struct {
	l *Lexer
	cur Token
	ahead []Token
}

func NewParser(l *Lexer) *Parser {
	return &Parser{l, Token{}, make([]Token, 0)}
}

func (p *Parser) pushAhead(t Token) {
	p.ahead = append(p.ahead, t)
}

func (p *Parser) next() (Token, error) {
	cur, err := p.l.NextToken()
	p.cur = cur
	return cur, err
}

func (p *Parser) ParseExpr() (Expr, error) {
	p.next()
	return p.parseAccess()
}

func (p *Parser) parseAccess() (Expr, error) {
	if p.cur.Type == TOK_ID {
		return &Access{nil, p.cur.Val.(string)}, nil
	}
	return nil, errors.New("Expected ID")
}

