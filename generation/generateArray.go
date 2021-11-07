package generation

import "fmt"

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

func GeneratePermutationsOfArray(arr []int) [][]int {

	length := len(arr)

	return GeneratePermutationsOfArrayRecurse(arr, make([]int, length), 0, length-1, 0, length, make([][]int, 0))
}

func GeneratePermutationsOfArrayRecurse(arr, currentCombination []int, start, end, index, length int, combinations [][]int) [][]int {

	if index == length {

		combinations = append(combinations, currentCombination)

		fmt.Println(currentCombination)

		return combinations
	}
	for i := start; (i <= end) && (end-i+1 >= length-index); i++ {

		currentCombination[index] = arr[i]

		newStart := i + 1

		newIndex := index + 1

		combinations = GeneratePermutationsOfArrayRecurse(arr, currentCombination, newStart, end, newIndex, length, combinations)
	}
	return combinations
}
