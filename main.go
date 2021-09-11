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
		"Fn f(x)=x+1",
		"Fn f(x,y)=x+y",
		"f(x,y)+z",
	}

	for _, expression := range expressions {

		result := parsing.ParseExpression(expression, parsing.Math)

		fmt.Println(result)
	}

}
