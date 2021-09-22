package main

import (
	"fmt"
	"os"
	"strconv"
	"symgolic/interpretation"
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
		// "f(x,y)+z",
		// "1+{4,5}",
		// "1+-{4,5}",
		// "{2,3}u{4,5}",
		// "{2,3}n{4,5}",
		// "[1,2,6,5]",
		// "[1,2,6,5]+[2,x]",
		// "(1,2,6,5)+(2,x)",
		// "ec(1+2+4)",
		// "sumliketerms((2*x)+(3*x))",
		// "distribute((2+y)*(3+x))",
		// "distribute(2*(3+x))",
		// "cancel((2*x*y)/(2*x*y))",
		// "cancel((2*x*y)/x)",
		// "cancel((2*x*y)/(x*y))",
		"cancel((2*x*(3+y))/(2*x*(3+y)))",
		"cancel((2*x*(3+y))/(2*(3+y)))",
	}

	var programs []string = []string{

		"x=2+3\ny=4+1",
		"f(x)=2+3\nf(x)*8",
		"f(x)=x+3\nf(x)*8",
	}

	for _, expression := range expressions {

		result := parsing.ParseExpression(expression)

		printTreeInfo(expression, result)

		interpretation.InterpretExpression(&result)

		fmt.Println("After function invocation")
		fmt.Println()

		printTreeInfo(expression, result)
	}

	for _, program := range programs {

		result := parsing.ParseProgramFromString(program)

		printProgramInfo(program, result)

		interpretation.InterpretProgram(&result)

		fmt.Println("After program interpretation")
		fmt.Println()

		printProgramInfo(program, result)
	}

	if len(files) != 0 {

		for _, program := range files {

			result := parsing.ParseProgramFromString(program)

			printProgramInfo(program, result)

			interpretation.InterpretProgram(&result)

			fmt.Println("After program interpretation")
			fmt.Println()

			printProgramInfo(program, result)
		}

	} else {

		fmt.Println("no files supplied")
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

func printProgramInfo(original string, program symbols.Program) {

	fmt.Println("Program: ")
	fmt.Println(program)
	fmt.Println()

	for i, expression := range program.Expressions {

		programLine := "From program line " + strconv.Itoa(i)

		printTreeInfo(programLine, expression)
	}
}
