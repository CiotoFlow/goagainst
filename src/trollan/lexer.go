package trollan

import "io"

type Pos struct {
	Offset int
	Row int
	Col int
}

type Lexer struct {
	io.Reader
	Pos
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{r, Pos{0,0,0}}
}

type Token struct {
	pos int
}

func (l *Lexer) NextToken() Token {
	l.Offset++
	return Token { 1 }
}