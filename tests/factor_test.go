package tests

import (
	"fmt"
	"symgolic/comparison"
	"symgolic/conversion"
	"symgolic/evaluation"
	"symgolic/generic"
	"symgolic/parsing"
	"testing"
)

type TermFactorTestData struct {
	Expression string

	Factors []evaluation.TermFactor
}

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

		result := evaluation.GetIsolatedFactors(original.GetRoot(), &original)

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

type TestGetTermFactorsData struct {
	Expression string

	Factors [][]string
}

func TestGetTermFactors(t *testing.T) {

	data := []TestGetTermFactorsData{
		{Expression: "a^3", Factors: [][]string{{"a", "a*a"}, {"a*a", "a"}}},
		{Expression: "3*x*y", Factors: [][]string{{"3", "x*y"}, {"x*y", "3"}, {"3*x", "y"}, {"y", "3*x"}, {"x*y", "3"}, {"3", "x*y"}, {"3*x*y", "1"}, {"1", "3*x*y"}}},
		{Expression: "3*x*(y^2)", Factors: [][]string{
			{"3", "x*y*y"}, {"x*y*y", "3"},
			{"3*x", "y*y"}, {"y*y", "3*x"},
			{"3*y", "x*y"}, {"x*y", "3*y"},
			{"3*x*y", "y"}, {"y", "3*x*y"},
			{"3*y*y", "x"}, {"x", "3*y*y"},
			{"3*x*y*y", "1"}, {"1", "3*x*y*y"}}},
	}

	for _, input := range data {

		original := parsing.ParseExpression(input.Expression)

		actual := evaluation.GetTermFactors(original.GetRoot(), &original)

		for _, expectedTermFactor := range input.Factors {

			factor := parsing.ParseExpression(expectedTermFactor[0])

			counterPart := parsing.ParseExpression(expectedTermFactor[1])

			if !ContainsTermFactor(evaluation.TermFactor{Factor: factor, CounterPart: counterPart}, actual) {

				t.Fatalf("Factor " + factor.ToString() + " not found in expected values")
			}
		}
	}
}

func TestGetCommonFactors(t *testing.T) {

	data := map[string][]string{
		"(2*(x^2))+(6*x)":           {"1", "2", "x", "2*x"},
		"(3*(x^2))+(6*x)":           {"1", "2", "3", "x", "2*x", "3*x"},
		"(8*x)+(16*x*y)+(24*(x^2))": {"1", "2", "4", "8", "x", "2*x", "4*x", "8*x"},
		"(2*(x^2))+(8*x)+(3*x)+12":  {"1"},
	}
	for input, output := range data {

		original := parsing.ParseExpression(input)

		actual := evaluation.GetCommonFactors(original.GetRoot(), &original)

		expected := conversion.ConvertBulkStringToExpression(output)

		for _, value := range actual {

			if !ContainsExpression(value.Factor, expected) {

				t.Fatalf("Factor " + value.Factor.ToString() + " not found in expected values")
			}
		}
	}
}

func TestGetFactorsByGrouping(t *testing.T) {

	data := map[string][]string{
		// "(2*(x^2))+(8*x)+(3*x)+12": {},
		// "(2*(x^2))+(11*x)+12": {},
		// "(2*(x^2))+(6*x)+12": {},
		// "(3*(x^2))+(6*x)+(4*x)+8":       {},
		// "(3*(x^2))-(6*x)-(4*x)+8":       {}, // negative numbers currently not working
		// "(2*(x^3))+(10*(x^2))+(3*x)+15": {},
		"(x^5)+(x^4)+(x^3)+(x^2)+x+1": {},
	}
	for input, _ := range data {

		original := parsing.ParseExpression(input)

		actual := evaluation.GetFactorsByGroupings(original.GetRoot(), &original)

		fmt.Println(actual)
	}
}

func TestEvaluateFactorisation(t *testing.T) {

	data := map[string]string{
		// "(a^3)+(b^3)+(3*a*b*(a+b))": "(a+b)^3",
		// "(a^2)+(2*a*b)+(b^2)": "(a+b)^2",
		"(a^2)-(b^2)": "(a+b)*(a-b)",
		// "(2*(x^2))+(6*x)": "2*x",
		// "(3*(x^2))+(6*x)":           "3*x",
		// "(8*x)+(16*x*y)+(24*(x^2))": "8*x",
		// "(2*(x^2))+(8*x)+(3*x)+12": "1",
	}

	for input, output := range data {

		original := parsing.ParseExpression(input)

		expected := parsing.ParseExpression(output)

		evaluation.EvaluateAndReplace(original.GetRoot(), &original, evaluation.Factor)

		if !comparison.IsEqual(original, expected) {

			err := "Expected " + expected.ToString() + " but instead got " + original.ToString()

			t.Fatalf(err)
		}
	}
}
