package main

import (
	"fmt"
	"symgolic/parsing"
)

func main() {

	var text string = "x+y+(3*9)"

	result, err := parsing.ParseExpression(text, parsing.Math)

	if err == nil {

		fmt.Println(result)
	}
}
