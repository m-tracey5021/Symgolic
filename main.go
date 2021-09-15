package main

import (
	"fmt"
	"symgolic/parsing"
)

func main() {

	var expressions []string = []string{

		// "x+y",
		// "-(x+y)",
		// "-(x+y)+3",
		// "-(x+y)+-3",
		// "-(x+y)+-(3+4)",
		// "x=y",
		// "x=y+z",
		// "x+y=a+b",
		// "x+y=a+(b*c)",
		// "x+(y*z)=a+(b*c)",
		// "x+y+z",
		// "x+((3*y)+(2/3))",
		// "f{x}=x+1",
		// "f{x,y}=x+y",
		// "f{x,y}+z",
		// "1+{4,5}",
		// "1+-{4,5}",
		"{2,3}u{4,5}",
		"{2,3}n{4,5}",
	}

	for _, expression := range expressions {

		result := parsing.ParseExpression(expression, parsing.Math)

		fmt.Print("Expression: ")
		fmt.Print(expression)
		fmt.Print(" maps to -> ")
		fmt.Println(result)
	}

}
