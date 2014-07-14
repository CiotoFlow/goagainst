package main

import (
	"fmt"
	"strings"
	"bufio"
	"./trollan"
	"./irc"
)

func testIRC() {
	client := irc.NewIRC("goagainst", "ai.irc.mufhd0.net:9999",
			     "#ciotoflow", true)

	err := client.Connect()
	if (err != nil) {
		fmt.Println(err)
		return
	}

	err = client.Loop()
	if (err != nil) {
		fmt.Println(err)
		client.Disconnect()
		return
	}
}

func main() {
	buf := strings.NewReader("test_test test123")
	l := trollan.NewLexer (bufio.NewReader(buf))
	tok, _ := l.NextToken()
	fmt.Println(l)
	fmt.Println(tok)
	tok, _ = l.NextToken()
	fmt.Println(l)
	fmt.Println(tok)

	testIRC()
}
