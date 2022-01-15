package generic

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

	combinations := GeneratePermutationsOfArrayRecurse(arr, make([]int, 0), make([][]int, 0))

	return combinations
}

func GeneratePermutationsOfArrayRecurse(arr, currentCombination []int, combinations [][]int) [][]int {

	if len(arr) != 0 {

		for i := 0; i < len(arr); i++ {

			element := arr[i]

			nextCombination := append(currentCombination, element)

			remaining := make([]int, 0)

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

func GenerateSubArrays(array []int, size int) [][]int {

	return GenerateSubArraysRecurse(array, make([]int, 0), make([][]int, 0), 0, size)
}

func GenerateSubArraysRecurse(array, output []int, subarrays [][]int, index, size int) [][]int {

	if index == len(array) {

		if len(output) != 0 {

			subarrays = append(subarrays, output)
		}
		return subarrays
	}
	subarrays = GenerateSubArraysRecurse(array, output, subarrays, index+1, size)

	if len(output) != size {

		output = append(output, array[index])

	} else {

		return subarrays
	}
	subarrays = GenerateSubArraysRecurse(array, output, subarrays, index+1, size)

	return subarrays
}

func GenerateSubArrayGroups(array []int) [][][]int {

	output := make([][][]int, 0)

	length := len(array)

	for i := 1; i <= length; i++ {

		currentGrouping := make([][]int, 0)

		for j := 0; j < i; j++ {

			currentGrouping = append(currentGrouping, make([]int, 0))
		}
		groups := GenerateSubArrayGroupsRecurse(0, length, i, 0, array, currentGrouping, make([][][]int, 0))

		output = append(output, groups...)
	}
	return output
}

func GenerateSubArrayGroupsRecurse(index, length, groupCountTarget, groupCount int, input []int, currentGrouping [][]int, output [][][]int) [][][]int {

	if index >= length {

		if groupCount == groupCountTarget {

			completeGroup := make([][]int, 0)

			for _, group := range currentGrouping {

				completeGroupValue := make([]int, 0)

				completeGroupValue = append(completeGroupValue, group...)

				completeGroup = append(completeGroup, completeGroupValue)
			}

			output = append(output, completeGroup)
		}
		return output
	}
	for i := 0; i < groupCountTarget; i++ {

		if len(currentGrouping[i]) > 0 {

			currentGrouping[i] = append(currentGrouping[i], input[index]) // push

			output = GenerateSubArrayGroupsRecurse(index+1, length, groupCountTarget, groupCount, input, currentGrouping, output)

			currentGrouping[i] = currentGrouping[i][:len(currentGrouping[i])-1] // pop

		} else {

			currentGrouping[i] = append(currentGrouping[i], input[index]) // push

			output = GenerateSubArrayGroupsRecurse(index+1, length, groupCountTarget, groupCount+1, input, currentGrouping, output)

			currentGrouping[i] = currentGrouping[i][:len(currentGrouping[i])-1] // pop

			break
		}
	}
	return output
}

func Contains(value int, arr []int) bool {

	for _, compared := range arr {

		if value == compared {

			return true
		}
	}
	return false
}

func MatchOrderedArray(arrA, arrB []int) bool {

	if len(arrA) != len(arrB) {

		return false

	} else {

		for i := 0; i < len(arrA); i++ {

			if arrA[i] != arrB[i] {

				return false
			}
		}
		return true
	}
}

func MatchUnorderedArray(arrA, arrB []int) bool {

	if len(arrA) != len(arrB) {

		return false

	} else {

		visited := make([]int, 0)

		for i := 0; i < len(arrA); i++ {

			found := false

			for j := 0; j < len(arrB); j++ {

				if Contains(j, visited) {

					continue

				} else {

					if arrA[i] == arrB[j] {

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

func RemoveDuplicates(arr []int) []int {

	keys := make(map[int]bool)

	list := make([]int, 0)

	for _, item := range arr {

		if _, value := keys[item]; !value {

			keys[item] = true

			list = append(list, item)
		}
	}
	return list
}

// ========== PER TYPE =================
