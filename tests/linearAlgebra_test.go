package tests

import (
	"symgolic/comparison"
	"symgolic/evaluation/linearAlgebra"
	"symgolic/parsing"
	"testing"
)

type ArithmeticTestData struct {
	InputA, InputB, Output string
}

func TestDotProduct(t *testing.T) {

	data := []ArithmeticTestData{

		{InputA: "[1, 2, 3]", InputB: "[4, 5, 6]", Output: "32"},
		{InputA: "[1, x, 3]", InputB: "[4, 5, 6]", Output: "22+(5*x)"},
	}
	for _, input := range data {

		a := parsing.ParseExpression(input.InputA)

		b := parsing.ParseExpression(input.InputB)

		expected := parsing.ParseExpression(input.Output)

		_, product := linearAlgebra.DotProduct(a.GetRoot(), b.GetRoot(), &a, &b)

		if !comparison.IsEqual(expected, product) {

			err := "expected " + expected.ToString() + " but got " + product.ToString()

			t.Fatalf(err)
		}
	}
}

func TestCrossProduct(t *testing.T) {

	data := []ArithmeticTestData{

		{InputA: "[2, 3, 4]", InputB: "[5, 6, 7]", Output: "[-3, 6, -3]"},
		{InputA: "[1, x, 3]", InputB: "[4, 5, 6]", Output: "[-15+(6*x), 6, 5-(4*x)]"},
	}
	for _, input := range data {

		a := parsing.ParseExpression(input.InputA)

		b := parsing.ParseExpression(input.InputB)

		expected := parsing.ParseExpression(input.Output)

		_, product := linearAlgebra.CrossProduct(a.GetRoot(), b.GetRoot(), &a, &b)

		if !comparison.IsEqual(expected, product) {

			err := "expected " + expected.ToString() + " but got " + product.ToString()

			t.Fatalf(err)
		}
	}
}
