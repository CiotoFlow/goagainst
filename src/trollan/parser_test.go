package trollan

import "testing"
import "bufio"
import "strings"

func TestParser(t *testing.T) {
	buf := strings.NewReader("test_test test123 1234 1234.12 \"test string\" * ** [")
	l := NewLexer (bufio.NewReader (buf))
	p := NewParser (l)

	expr, err := p.ParseExpr()
	if err != nil { t.Errorf(err.Error()) }

	a, ok := expr.(*Access)
	if !ok { t.Errorf("Not Access: %s", expr) }
	if a.name != "test_test" { t.Errorf("Not test_test access: %s", a) }
}
