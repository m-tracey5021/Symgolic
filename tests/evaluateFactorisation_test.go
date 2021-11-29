package tests

import (
	"fmt"
	"symgolic/comparison"
	"symgolic/evaluation"
	"symgolic/generic"
	"symgolic/parsing"
	"testing"
)

func TestGetIsolatedFactors(t *testing.T) {

	data := map[string][]string{
		"a*a*a":   {"a", "a", "a"},
		"3":       {"3", "1"},
		"3*x":     {"3", "1", "x"},
		"3*(x^2)": {"3", "1", "x^2"},
		"3*x*y":   {"3", "1", "x", "y"},
		"4*x*y":   {"1", "2", "4", "x", "y"},
	}
	for input, output := range data {

		original := parsing.ParseExpression(input)

		_, result := evaluation.GetIsolatedFactors(original.GetRoot(), &original)

		err := "Expected does not match actual"

		if len(result) != len(output) {

			t.Fatalf(err)

		} else {

			visited := make([]int, 0)

			for i := 0; i < len(output); i++ {

				found := false

				for j := 0; j < len(output); j++ {

					if generic.Contains(j, visited) {

						continue

					} else {

						if comparison.IsEqual(result[i], parsing.ParseExpression(output[j])) {

							found = true

							visited = append(visited, j)

							break
						}
					}
				}
				if !found {

					t.Fatalf(err)
				}

			}
		}
	}
}

func TestGetTermFactors(t *testing.T) {

	data := map[string]string{
		"a^3":   "a",
		"3*x*y": "a",
	}

	for input, _ := range data {

		original := parsing.ParseExpression(input)

		// expected := parsing.ParseExpression(output)

		result := evaluation.GetTermFactors(original.GetRoot(), &original)

		fmt.Print(result)
	}
}

func TestGetCommonFactors(t *testing.T) {

}

func TestEvaluateFactorisation(t *testing.T) {

	data := map[string]string{
		"(a^3)+(b^3)+(3*a*b*(a+b))": "a",
	}

	for input, output := range data {

		original := parsing.ParseExpression(input)

		expected := parsing.ParseExpression(output)

		evaluation.EvaluateAndReplace(original.GetRoot(), &original, evaluation.EvaluateFactorisation)

		if !comparison.IsEqual(original, expected) {

			err := "Expected " + expected.ToString() + " but instead got " + original.ToString()

			t.Fatalf(err)
		}
	}
}
