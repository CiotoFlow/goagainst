package trollan

import "bufio"
import "unicode"
import "io"
import "strconv"
/* import "fmt" */

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
	TOK_INT
	TOK_FLOAT
	TOK_OPER
	TOK_PUN
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

func (l *Lexer) NextToken() (tok Token, err error) {
	tok, err = l.nextTokenReal()
	if err == io.EOF {
		err = nil
	}
	return
}

func (l *Lexer) pushAhead(b rune) {
	l.ahead = append(l.ahead, b)
}

func isOper(k rune) bool {
	return k == '+' || k == '-' || k == '*' ||
	  k == '/' || k == '>' || k == '<' ||
	  k == '='
}

func isPun(k rune) bool {
	return k == '[' || k == ']' || k == '(' ||
		k == ')' || k == '{' || k == '}' ||
		k == '.' || k == ';' || k == ':' ||
		k == ','
}

func (l *Lexer) nextTokenReal() (tok Token, err error) {
	tok = Token{}
	tok.Type = TOK_EOF
	var b rune

	l.skipSpaces()

	b, err = l.nextRune ()
	if err != nil {
		return
	}

	if unicode.IsLetter(b) {
		strTok := make([]rune, 0)
		tok.Pos = l.pos
		for {
			strTok = append(strTok, b);
			l.pos.Offset++

			b, err = l.nextRune()
			if b == 0 || err != nil {
				break
			} else if !(unicode.IsLetter(b) || b == '_' || unicode.IsDigit(b)) {
				l.pushAhead(b)
				break
			}
		}
		tok.Type = TOK_ID
		tok.Val = string(strTok)
	} else if unicode.IsDigit(b) {
		strTok := make([]rune, 0)
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
					l.pushAhead(b)
					break
				} else {
					foundDot = true
				}
			} else if !unicode.IsDigit(b) {
				l.pushAhead(b)
				break;
			}
		}

		if foundDot {
			tok.Type = TOK_FLOAT
			tok.Val, err = strconv.ParseFloat(string(strTok), 64)
		} else {
			tok.Type = TOK_INT
			tok.Val, err = strconv.ParseInt(string(strTok), 10, 64)
		}		
	} else if b == '"' {
		strTok := make([]rune, 0)
		tok.Pos = l.pos
		
		for {
			b, err = l.nextRune()
			l.pos.Offset++
			
			if b == 0 || err != nil {
				break
			} else if b == '"' {
				l.nextRune()
				break
			} else {
				strTok = append(strTok, b)
			}

		}
		
		tok.Type = TOK_STR
		tok.Val = string(strTok)
	} else if isOper(b) {
		tok.Pos = l.pos
		tok.Type = TOK_OPER
		tok.Val = string(b)

		b, err = l.nextRune()
		if b == 0 || err != nil {
			return
		} else if isOper(b) {
			tok.Val = tok.Val.(string)+string(b)
		} else {
			l.pushAhead(b)
		}
	} else if isPun(b) {
		tok.Pos = l.pos
		tok.Type = TOK_PUN
		tok.Val = string(b)
	}
	return
}
