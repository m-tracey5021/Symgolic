package generation

import (
	. "symgolic/symbols"
)

// make these generic once it is out

func GenerateCombinations(terms [][]int, limit int) [][]int { // needs to be [][]interface{}

	return GenerateCombinationsRecurse(terms, make([][]int, 0), make([]int, 0), limit)
}

func GenerateCombinationsRecurse(terms, combinations [][]int, accumulated []int, limit int) [][]int {

	last := len(terms) == 1

	n := len(terms[0])

	for i := 0; i < n; i++ {

		accumulated = append(accumulated, terms[0][i])

		item := accumulated

		if last {

			if len(item) == limit {

				combinations = append(combinations, item)
			}

		} else {

			combinations = GenerateCombinationsRecurse(terms[1:], combinations, item, limit)
		}
	}
	return combinations
}

func GenerateCombinationsByRow(matrix [][]int) [][]int {

	return GenerateCombinationsByRowRecurse(matrix, make([][]int, 0), make([]int, len(matrix)), 0)
}

func GenerateCombinationsByRowRecurse(matrix [][]int, combinations [][]int, rowIndexes []int, currentColumn int) [][]int {

	for rowNumber := range matrix[currentColumn] {

		rowIndexes[currentColumn] = rowNumber

		if currentColumn == len(matrix)-1 { // if its the last column

			comboPerLine := make([]int, 0)

			for colNumber, rowNumber := range rowIndexes {

				comboPerLine = append(comboPerLine, matrix[colNumber][rowNumber])
			}
			combinations = append(combinations, comboPerLine)

		} else {

			combinations = GenerateCombinationsByRowRecurse(matrix, combinations, rowIndexes, currentColumn+1)
		}
	}
	return combinations
}

func GeneratePermutationsOfArray(arr []Expression) [][]Expression {

	combinations := GeneratePermutationsOfArrayRecurse(arr, make([]Expression, 0), make([][]Expression, 0))

	return combinations
}

func GeneratePermutationsOfArrayRecurse(arr, currentCombination []Expression, combinations [][]Expression) [][]Expression {

	if len(arr) != 0 {

		for i := 0; i < len(arr); i++ {

			element := arr[i]

			nextCombination := append(currentCombination, element)

			remaining := make([]Expression, 0)

			remaining = append(remaining, arr[i+1:]...)

			remaining = append(remaining, arr[:i]...)

			combinations = GeneratePermutationsOfArrayRecurse(remaining, nextCombination, combinations)
		}
		return combinations

	} else {

		combinations = append(combinations, currentCombination)

		return combinations
	}
}
