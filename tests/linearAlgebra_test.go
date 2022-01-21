package tests

import (
	"symgolic/language/interpretation"
	"symgolic/language/interpretation/algebra"
	"symgolic/language/parsing"
	"testing"
)

type ScaleTestData struct {
	Input, Output, Scalar string
}

type VectorTestData struct {
	InputA, InputB, Output string
}

type DeterminantTestData struct {
	Input, Output string
}

func TestScale(t *testing.T) {

	data := []ScaleTestData{

		{Input: "[1, 2, 3]", Output: "[3, 6, 9]", Scalar: "3"},
	}
	for _, input := range data {

		a := parsing.ParseExpression(input.Input)

		scalar := parsing.ParseExpression(input.Scalar)

		expected := parsing.ParseExpression(input.Output)

		_, scaled := algebra.Scale(a.GetRoot(), &a, &scalar)

		if !interpretation.IsEqual(expected, scaled) {

			err := "expected " + expected.ToString() + " but got " + scaled.ToString()

			t.Fatalf(err)
		}
	}
}

func TestVectorAdd(t *testing.T) {

	data := []VectorTestData{

		{InputA: "[1, 2, 3]", InputB: "[3, 6, 9]", Output: "[4, 8, 12]"},
	}
	for _, input := range data {

		a := parsing.ParseExpression(input.InputA)

		b := parsing.ParseExpression(input.InputB)

		expected := parsing.ParseExpression(input.Output)

		_, added := algebra.VectorAdd(a.GetRoot(), b.GetRoot(), &a, &b)

		if !interpretation.IsEqual(expected, added) {

			err := "expected " + expected.ToString() + " but got " + added.ToString()

			t.Fatalf(err)
		}
	}
}

func TestDotProduct(t *testing.T) {

	data := []VectorTestData{

		{InputA: "[1, 2, 3]", InputB: "[4, 5, 6]", Output: "32"},
		{InputA: "[1, x, 3]", InputB: "[4, 5, 6]", Output: "22+(5*x)"},
	}
	for _, input := range data {

		a := parsing.ParseExpression(input.InputA)

		b := parsing.ParseExpression(input.InputB)

		expected := parsing.ParseExpression(input.Output)

		_, product := algebra.DotProduct(a.GetRoot(), b.GetRoot(), &a, &b)

		if !interpretation.IsEqual(expected, product) {

			err := "expected " + expected.ToString() + " but got " + product.ToString()

			t.Fatalf(err)
		}
	}
}

func TestCrossProduct(t *testing.T) {

	data := []VectorTestData{

		{InputA: "[2, 3, 4]", InputB: "[5, 6, 7]", Output: "[-3, 6, -3]"},
		{InputA: "[1, x, 3]", InputB: "[4, 5, 6]", Output: "[-15+(6*x), 6, 5-(4*x)]"},
	}
	for _, input := range data {

		a := parsing.ParseExpression(input.InputA)

		b := parsing.ParseExpression(input.InputB)

		expected := parsing.ParseExpression(input.Output)

		_, product := algebra.CrossProduct(a.GetRoot(), b.GetRoot(), &a, &b)

		if !interpretation.IsEqual(expected, product) {

			err := "expected " + expected.ToString() + " but got " + product.ToString()

			t.Fatalf(err)
		}
	}
}

func TestFindDeterminant(t *testing.T) {

	data := []DeterminantTestData{

		{Input: "([2, 3], [1, 4])", Output: "5"},
		{Input: "([1, 2, 3], [4, 5, 6], [7, 8, 9])", Output: "0"},
	}
	for _, input := range data {

		a := parsing.ParseExpression(input.Input)

		expected := parsing.ParseExpression(input.Output)

		determinant := algebra.FindDeterminant(a.GetRoot(), a)

		if !interpretation.IsEqual(expected, determinant) {

			err := "expected " + expected.ToString() + " but got " + determinant.ToString()

			t.Fatalf(err)
		}
	}
}

func TestRref(t *testing.T) {

	data := []DeterminantTestData{

		{Input: "([2, 4], [1, 4])", Output: "([1, 2], [0, 2])"},
	}
	for _, input := range data {

		a := parsing.ParseExpression(input.Input)

		expected := parsing.ParseExpression(input.Output)

		rref := algebra.Rref(a.GetRoot(), a)

		if !interpretation.IsEqual(expected, rref) {

			err := "expected " + expected.ToString() + " but got " + rref.ToString()

			t.Fatalf(err)
		}
	}
}
