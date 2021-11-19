package tests

import (
	"fmt"
	. "symgolic/parsing"
	. "symgolic/solvers"
	"testing"
)

func TestSolveForConstantValue(t *testing.T) {

	target := ParseExpression("12")

	expression := ParseExpression("(3+b)*c")

	solutionExists, solutions := SolveForConstantValue(expression.GetRoot(), &target, &expression)

	if solutionExists {

		for _, solution := range solutions {

			fmt.Println(solution)
		}
	}
}

func TestSolveForMultipleConstantValues(t *testing.T) {

	valuesA := []SolveRequest{

		{ParseExpression("6"), ParseExpression("a+b")},

		{ParseExpression("9"), ParseExpression("a*b")},
	}
	valuesB := []SolveRequest{

		{ParseExpression("12"), ParseExpression("(a+b)*c")}, // a = 1, b = 2, c = 4

		{ParseExpression("2"), ParseExpression("a*b")},
	}
	valuesC := []SolveRequest{

		{ParseExpression("12"), ParseExpression("a")},
	}
	resultA := SolveForMultipleConstantValues(valuesA)

	resultB := SolveForMultipleConstantValues(valuesB)

	resultC := SolveForMultipleConstantValues(valuesC)

	valuesForA := resultB.GetValuesFor("a")

	valuesForB := resultB.GetValuesFor("b")

	valuesForC := resultB.GetValuesFor("c")

	fmt.Println(valuesForA)

	fmt.Println(valuesForB)

	fmt.Println(valuesForC)

	fmt.Println(resultA)

	fmt.Println(resultB)

	fmt.Println(resultC)
}
