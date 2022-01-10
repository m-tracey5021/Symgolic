package solvers

import (
	. "symgolic/language/components"
)

func SolveByConstantValue() {

}

// ========= RECURSIVE FUNCTIONS THAT SHOULD BE GENERIC ===========

func Expression_GeneratePermutationsOfArray(arr []Expression) [][]Expression {

	combinations := Expression_GeneratePermutationsOfArrayRecurse(arr, make([]Expression, 0), make([][]Expression, 0))

	return combinations
}

func Expression_GeneratePermutationsOfArrayRecurse(arr, currentCombination []Expression, combinations [][]Expression) [][]Expression {

	if len(arr) != 0 {

		for i := 0; i < len(arr); i++ {

			element := arr[i]

			nextCombination := append(currentCombination, element)

			remaining := make([]Expression, 0)

			remaining = append(remaining, arr[i+1:]...)

			remaining = append(remaining, arr[:i]...)

			combinations = Expression_GeneratePermutationsOfArrayRecurse(remaining, nextCombination, combinations)
		}
		return combinations

	} else {

		combinations = append(combinations, currentCombination)

		return combinations
	}
}
