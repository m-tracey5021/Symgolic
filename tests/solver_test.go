package tests

import (
	"fmt"
	"symgolic/solvers"
	"testing"
)

func TestSolveForMultipleConstantValues(t *testing.T) {

	valuesA := map[int]string{

		6: "a+b",
		9: "a*b",
	}
	valuesB := map[int]string{

		12: "(a+b)*c", // a = 1, b = 2, c = 4
		2:  "a*b",     //
	}
	resultA := solvers.SolveForMultipleConstantValues(valuesA)

	resultB := solvers.SolveForMultipleConstantValues(valuesB)

	valuesForA := resultB.GetValuesFor("a")

	valuesForB := resultB.GetValuesFor("b")

	valuesForC := resultB.GetValuesFor("c")

	fmt.Println(valuesForA)

	fmt.Println(valuesForB)

	fmt.Println(valuesForC)

	fmt.Println(resultA)

	fmt.Println(resultB)
}
