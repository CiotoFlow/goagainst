package main

import (
	"fmt"
	"./trollan"
)

func main() {
	tok := trollan.NextToken()
	fmt.Println(tok)
}