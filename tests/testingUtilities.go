package tests

import (
	"symgolic/comparison"
	"symgolic/evaluation"
	"symgolic/generic"
	"symgolic/symbols"
)

type MatchFunction func(interface{}, interface{}) bool

func MatchUnorderedArray_ForExpression(arrA, arrB []symbols.Expression) bool {

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

					if comparison.IsEqual(arrA[i], arrB[j]) {

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

func ContainsExpression(value symbols.Expression, arr []symbols.Expression) bool {

	for _, compared := range arr {

		if comparison.IsEqual(value, compared) {

			return true
		}
	}
	return false
}

func ContainsTermFactor(value evaluation.TermFactor, arr []evaluation.TermFactor) bool {

	for _, compared := range arr {

		if comparison.IsEqual(value.Factor, compared.Factor) && comparison.IsEqual(value.CounterPart, compared.CounterPart) {

			return true
		}
	}
	return false
}
