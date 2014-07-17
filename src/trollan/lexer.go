package trollan

import "bufio"
import "unicode"
import "io"
//import "fmt"

type Pos struct {
	Offset int
	Row int
	Col int
}

type Lexer struct {
	*bufio.Reader
	pos Pos
	ahead []rune
}

func NewLexer(r *bufio.Reader) *Lexer {
	return &Lexer{r, Pos{0,0,0}, make([]rune, 0)}
}

type TokenType int;

const (
	TOK_EOF TokenType = iota
	TOK_ID
	TOK_STR
	TOK_FLOAT
)

type Token struct {
	Pos
	Val interface{}
	Type TokenType
}

func (l *Lexer) skipSpaces() {
	for {
		b, err := l.nextRune ()
		if err != nil {
			return
		}

		if b == 0 {
			break
		}

		if !unicode.IsSpace (b) {
			l.ahead = append (l.ahead, b)
			break
		}
	}
}

func (l *Lexer) nextRune() (b rune, err error) {
	if len(l.ahead) > 0 {
		b = l.ahead[0]
		err = nil
		l.ahead = l.ahead[1:]
	} else {
		b, _, err = l.ReadRune()
	}
	return
}

func (l *Lexer) NextToken() (tok *Token, err error) {
	tok, err = l.nextTokenReal()
	if err == io.EOF {
		err = nil
	}
	return
}

func (l *Lexer) nextTokenReal() (tok *Token, err error) {
	tok = new(Token)
	tok.Type = TOK_EOF
	var b rune

	l.skipSpaces()

	b, err = l.nextRune ()
	if err != nil {
		return
	}

	if unicode.IsLetter(b) {
		strTok := make([]rune, 30)
		tok.Pos = l.pos
		for {
			strTok = append(strTok, b);
			l.pos.Offset++

			b, err = l.nextRune()
			if b == 0 || err != nil {
				break
			} else if !(unicode.IsLetter(b) || b == '_' || unicode.IsDigit(b)) {
				l.ahead = append(l.ahead, b)
				break
			}
		}
		tok.Type = TOK_ID
		tok.Val = string(strTok)
	}

	if unicode.IsDigit(b) {
		strTok := make([]rune, 30)
		tok.Pos = l.pos
		var foundDot bool

		for {

			strTok = append(strTok, b)
			l.pos.Offset++

			b, err = l.nextRune()
			if b == 0 || err != nil {
				break
			} else if b == '.' {
				if foundDot {
					l.ahead = append(l.ahead, b)
					break
				} else {
					foundDot = true
				}
			} else if !unicode.IsDigit(b) {
				l.ahead = append(l.ahead, b)
				break;
			}
			tok.Type = TOK_FLOAT
			tok.Val = string(strTok)
		}
		
	}

	return
}

func (t *Token) StringVal() string {
	return t.Val.(string)
}

func (t *Token) DoubleVal() float64 {
	return t.Val.(float64)
}
