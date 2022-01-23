package tests

import (
	"symgolic/generic"
	. "symgolic/language/components"
	"symgolic/language/interpretation"
	"symgolic/language/parsing"
)

type BinaryTestData struct {
	Input, Expected string
}

type TernaryTestData struct {
	InputA, InputB, Expected string
}

type QuarternaryTestData struct {
	InputA, InputB, InputC, Expected string
}

func TestBinaryDataOverManipulation(testData []BinaryTestData, manipulate interpretation.Manipulation) (bool, string) {

	for _, data := range testData {

		input := parsing.ParseExpression(data.Input)

		expected := parsing.ParseExpression(data.Expected)

		output := manipulate(From(input))

		if !interpretation.IsEqual(expected, output) {

			return false, "Expected: " + expected.ToString() + ", output: " + input.ToString()
		}
	}
	return true, "Success"
}

func TestBinaryDataOverEvaluation(testData []BinaryTestData, evaluate interpretation.Evaluation) (bool, string) {

	for _, data := range testData {

		input := parsing.ParseExpression(data.Input)

		expected := parsing.ParseExpression(data.Expected)

		interpretation.EvaluateAndReplace(From(input), evaluate)

		if !interpretation.IsEqual(expected, input) {

			return false, "Expected: " + expected.ToString() + ", output: " + input.ToString()
		}
	}
	return true, "Success"
}

func TestTernaryDataOverManipulationAgainst(testData []TernaryTestData, manipulate interpretation.ManipulationAgainst) (bool, string) {

	for _, data := range testData {

		inputA := parsing.ParseExpression(data.InputA)

		inputB := parsing.ParseExpression(data.InputB)

		expected := parsing.ParseExpression(data.Expected)

		output := manipulate(From(inputA), From(inputB))

		if !interpretation.IsEqual(expected, output) {

			return false, "Expected: " + expected.ToString() + ", output: " + output.ToString()
		}
	}
	return true, "Success"
}

func TestTernaryDataOverManipulationForMany(testData []TernaryTestData, manipulate interpretation.ManipulationForMany) (bool, string) {

	for _, data := range testData {

		inputA := parsing.ParseExpression(data.InputA)

		inputB := parsing.ParseExpression(data.InputB)

		expected := parsing.ParseExpression(data.Expected)

		output := manipulate(From(inputA), From(inputB))

		if !interpretation.IsEqual(expected, output) {

			return false, "Expected: " + expected.ToString() + ", output: " + output.ToString()
		}
	}
	return true, "Success"
}

type MatchFunction func(interface{}, interface{}) bool

func MatchUnorderedArray_ForExpression(arrA, arrB []Expression) bool {

	if len(arrA) != len(arrB) {

		return false

	} else {

		visited := make([]int, 0)

		for i := 0; i < len(arrA); i++ {

			found := false

			for j := 0; j < len(arrB); j++ {

				if generic.Contains(j, visited) {

					continue

				} else {

					if interpretation.IsEqual(arrA[i], arrB[j]) {

						found = true

						visited = append(visited, j)

						break
					}
				}
			}
			if !found {

				return false
			}
		}
		return true
	}
}

func ContainsExpression(value Expression, arr []Expression) bool {

	for _, compared := range arr {

		if interpretation.IsEqual(value, compared) {

			return true
		}
	}
	return false
}

func ContainsTermFactor(value interpretation.TermFactor, arr []interpretation.TermFactor) bool {

	for _, compared := range arr {

		if interpretation.IsEqual(value.Factor, compared.Factor) && interpretation.IsEqual(value.CounterPart, compared.CounterPart) {

			return true
		}
	}
	return false
}

func ConvertBulkStringToExpression(values []string) []Expression {

	expressions := make([]Expression, 0)

	for _, value := range values {

		expressions = append(expressions, parsing.ParseExpression(value))
	}
	return expressions
}
