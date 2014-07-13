package trollan

import "bufio"
import "fmt"

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

func (l *Lexer) NextToken() (tok *Token, err error) {
	tok = new(Token)
	strTok := make([]byte, 30)
	var b byte
	var foundSpace bool

	for {
		b, err = l.ReadByte ()
		if err != nil {
			
		}

		if foundSpace {
			//l.pos.Offset++
			break
		}

		if (b >= 0x41 && b <= 0x5a) ||
			(b >= 61 && b <= 0x7a) ||
			(b == 0x5f) {
			strTok = append(strTok, b)
		} else if b == 0x20 || b == 0x0 {
			foundSpace = true
		}

		fmt.Sprintf("%c\n",b)
		
		l.pos.Offset++
	}

	tok.Pos = l.pos
	tok.Val = string(strTok)
	tok.Type = ID
	return
}

func (t *Token) StringVal() string {
	return t.Val.(string)
}

func (t *Token) DoubleVal() float64 {
	return t.Val.(float64)
}
