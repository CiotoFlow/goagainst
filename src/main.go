package main

import (
	"fmt"
	"strings"
	"bufio"
	"./trollan"
)

func main() {
	buf := strings.NewReader("test test")
	l := trollan.NewLexer (bufio.NewReader(buf))
	tok, _ := l.NextToken()
	fmt.Println(l)
	fmt.Println(tok)
}