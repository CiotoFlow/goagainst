package trollan

import "fmt"

type Expr interface {
	isExpr()
	fmt.Stringer
}

/* ACCESS */
type Access struct {
	inner *Access
	name string
}

func (a *Access) isExpr() {}

func (a *Access) String() string {
	if a == nil {
		return ""
	}
	
	if (a.inner != nil) {
		return (*a.inner).String()+"."+a.name
	}
	
	return a.name
}
