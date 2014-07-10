package trollan

import "bufio"

type Pos struct {
	Offset int
	Row int
	Col int
}

type Lexer struct {
	*bufio.Reader
	pos Pos
}

func NewLexer(r *bufio.Reader) *Lexer {
	return &Lexer{r, Pos{0,0,0}}
}

type Token struct {
	Pos
	Val interface{}
}

func (l *Lexer) NextToken() (tok *Token, err error) {
	tok = new(Token)
	b, err := l.ReadByte ()
	if err != nil {
		return
	}
	
	l.pos.Offset++
	tok.Pos = l.pos
	tok.Val = string(b)
	return
}

func (t *Token) StringVal() string {
	return t.Val.(string)
}

func (t *Token) DoubleVal() float64 {
	return t.Val.(float64)
}