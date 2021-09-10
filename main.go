package main

import (
	"fmt"
	"symgolic/parsing"
)

func main() {

	var expressions []string = []string{

		"x+y",
		"-(x+y)",
		"-(x+y)+3",
		"-(x+y)+-3",
		"-(x+y)+-(3+4)",
		"",
		"",
		"",
		"",
		"",
	}

	for _, expression := range expressions {

		result, err := parsing.ParseExpression(expression, parsing.Math)

		if err == nil {

			fmt.Println(result)
		}
	}

}
