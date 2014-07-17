package trollan

import "testing"
import "bufio"
import "strings"

func TestSimple(t *testing.T) {
	buf := strings.NewReader("test_test test123 1234.12")
	l := NewLexer (bufio.NewReader (buf))

	tok, err := l.NextToken()
	if err != nil { t.Errorf(err.Error()) }
	if tok.Type != TOK_ID && tok.Val != "test_test" { t.Errorf("%s", tok) }
	
	tok, err = l.NextToken()
	if err != nil { t.Errorf(err.Error()) }
	if tok.Type != TOK_ID && tok.Val != "test123" { t.Errorf("%s", tok) }

	tok, err = l.NextToken()
	if err != nil { t.Errorf(err.Error()) }
	if tok.Type != TOK_FLOAT && tok.Val != "1234.12" { t.Errorf("%s", tok) }
}
