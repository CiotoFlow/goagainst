package main

import (
	"fmt"
	"./trollan"
)

func main() {
	l := trollan.NewLexer (nil)
	cp := new(trollan.Lexer)
	*cp = *l
	l.NextToken()
	fmt.Println(l.Offset)
	fmt.Println(cp.Offset)
}