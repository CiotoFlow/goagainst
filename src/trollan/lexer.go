package trollan

import "bufio"
//import "fmt"

type Pos struct {
	Offset int
	Row int
	Col int
}

type Lexer struct {
	*bufio.Reader
	pos Pos
	ahead []byte
}

func NewLexer(r *bufio.Reader) *Lexer {
	return &Lexer{r, Pos{0,0,0}, make([]byte, 0)}
}

type TokenType int;

const (
	EOF TokenType = iota
	ID
	STR
	FLOAT
)

type Token struct {
	Pos
	Val interface{}
	Type TokenType
}

func (l *Lexer) skipSpaces() {
	for {
		b, err := l.nextByte ()
		if err != nil {
			return
		}

		if b == 0 {
			break
		}

		if b != 0x20 {
			l.ahead = append (l.ahead, b)
			break
		}
	}
}

func (l *Lexer) nextByte() (b byte, err error) {
	if len(l.ahead) > 0 {
		b = l.ahead[0]
		err = nil
		l.ahead = l.ahead[1:]
	} else {
		b, err = l.ReadByte()
	}
	return
}

func isAlpha(b byte) bool {
	return (b >= 0x41 && b <= 0x5a) ||
			(b >= 61 && b <= 0x7a) ||
			(b == 0x5f)
}

func (l *Lexer) NextToken() (tok *Token, err error) {
	tok = new(Token)
	var b byte

	l.skipSpaces()

	b, err = l.nextByte ()
	if err != nil {
		return
	}

	if isAlpha(b) {
		strTok := make([]byte, 30)
		tok.Pos = l.pos
		for {
			strTok = append(strTok, b);
			l.pos.Offset++

			b, err = l.nextByte()
			if b == 0 || err != nil {
				break
			} else if !isAlpha(b) {
				l.ahead = append(l.ahead, b)
				break
			}
		}
		tok.Type = ID
		tok.Val = string(strTok)
	}

	return
}

func (t *Token) StringVal() string {
	return t.Val.(string)
}

func (t *Token) DoubleVal() float64 {
	return t.Val.(float64)
}
