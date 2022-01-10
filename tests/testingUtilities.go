package tests

import (
	"symgolic/generic"
	"symgolic/language/components"
	"symgolic/language/interpretation"
	"symgolic/language/parsing"
)

type MatchFunction func(interface{}, interface{}) bool

func MatchUnorderedArray_ForExpression(arrA, arrB []components.Expression) bool {

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

func ContainsExpression(value components.Expression, arr []components.Expression) bool {

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

func ConvertBulkStringToExpression(values []string) []components.Expression {

	expressions := make([]components.Expression, 0)

	for _, value := range values {

		expressions = append(expressions, parsing.ParseExpression(value))
	}
	return expressions
}
