package tests

import (
	"fmt"
	"symgolic/interpretation"
	"symgolic/parsing"
	"testing"
)

func TestInterpret(t *testing.T) {

	data := []string{

		// "x=2+3\ny=4+1",
		// "f(x)=2+3\ny=f(x)*8",
		// "f(x)=x+3\nf(x)*8",
		// "f(x)=x+1\nf(y)",
		// "f(x)=x+1\ng(y)=f(y)+2",
		// "f(x)=x+1\ng(y)=f(x)+2",
		// "g(y)=f(x)\nf(x)=x+1",
		// "f(x)=expandexponents(x^((2*y)+(3*z)))\nf(x)*8",
		// "f(x):=(x+2)*3\nz:=f(y)\na:=z+2",
		"a:=distribute(2*(x+y))",
		// "f(x)=(x+2)*3\ny=5\nz=f(y)",
	}
	for _, program := range data {

		parsed := parsing.ParseProgramFromString(program)

		results := interpretation.InterpretProgram(&parsed)

		fmt.Println(results)
	}
}
