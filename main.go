package main

import (
	"fmt"
	"os"
	"strconv"
	"symgolic/parsing"
	"symgolic/symbols"
)

func main() {

	files := os.Args[1:]

	// fmt.Println(files)

	var expressions []string = []string{

		// "2+3",
		// "2+(4+5)",
		// "2+(4*5)",
		// "2+(4*x)",
		// "2+3+(4*x)",
		// "x+y",
		// "-(x+y)",
		// "-(x+y)+3",
		// "-(x+y)+-3",
		// "-(x+y)+-(3+4)",
		// "x=y",
		// "x=y+z",
		// "y=9+1",
		// "x+y=a+b",
		// "x+y=a+(b*d)",
		// "x+(y*z)=a+(b*d)",
		// "x+y+z",
		// "x+((3*y)+(2/3))",
		"f(x)=x+1",
		"f(x,y)=x+y",
		"f(x,y)+z",
		"1+{4,5}",
		"1+-{4,5}",
		"{2,3}u{4,5}",
		"{2,3}n{4,5}",
		"[1,2,6,5]",
		"[1,2,6,5]+[2,x]",
		"(1,2,6,5)+(2,x)",
	}

	// var programs []string = []string{

	// 	"x=2+3\ny=4+1",
	// }

	for _, expression := range expressions {

		result := parsing.ParseExpression(expression)

		// interpretation.InvokeFunction("ec", result.GetRoot(), &result)

		printTreeInfo(expression, result)
	}

	// for _, program := range programs {

	// 	fmt.Println("Program: ", program)
	// 	fmt.Println()

	// 	result := parsing.ParseProgramFromString(program)

	// 	for i, expression := range result.Expressions {

	// 		programLine := "From program line " + strconv.Itoa(i)

	// 		printTreeInfo(programLine, expression)
	// 	}

	// }

	if len(files) != 0 {

		fmt.Println("File: ")
		fmt.Println()
		fmt.Println(files[0])

		fromFile := parsing.ParseProgramFromFile(files[0])

		for i, expression := range fromFile.Expressions {

			programLine := "From file line " + strconv.Itoa(i)

			printTreeInfo(programLine, expression)
		}

	} else {

		fmt.Println("no file supplied")
	}

}

func printTreeInfo(original string, parsed symbols.Expression) {

	fmt.Println("==========================")
	fmt.Println()
	fmt.Println("Original expression: ", original)
	fmt.Println()
	fmt.Println("Maps to: ", parsed)
	fmt.Println()
	fmt.Println("Pretty print: ", parsed.ToString())
	fmt.Println()
	fmt.Println("vvvv Indented print vvvv")

	parsed.PrintTree(parsed.GetRoot(), 2, 0)

	fmt.Println("^^^^ Indented print ^^^^")

	fmt.Println("==========================")
	fmt.Println()
}
