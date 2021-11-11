package solvers

import (
	"strconv"
	. "symgolic/comparison"
	. "symgolic/evaluation"
	. "symgolic/parsing"
	. "symgolic/symbols"
)

type SolutionSet struct {
	Mapping map[string]Expression
}

type SolutionFor struct {
	Value Expression

	Expression Expression

	Solutions []SolutionSet
}

func MergeSolutions(A SolutionSet, B SolutionSet) SolutionSet {

	merged := SolutionSet{}

	for variable, constant := range A.Mapping {

		merged.Mapping[variable] = constant
	}
	for variable, constant := range B.Mapping {

		merged.Mapping[variable] = constant
	}
	return merged
}

func MergeMultipleSolutionsOneToMany(merged []SolutionSet, toMerge SolutionSet) []SolutionSet { // maybe check for compatibility here too

	if len(merged) == 0 {

		return []SolutionSet{toMerge}
	}
	for _, merge := range merged {

		for key, value := range toMerge.Mapping {

			merge.Mapping[key] = value
		}
	}
	return merged
}

func MergeMultipleSolutionsManyToOne(toMerge []SolutionSet) (bool, SolutionSet) {

	merged := SolutionSet{}

	for _, solution := range toMerge {

		for variable, value := range solution.Mapping {

			otherValue, exists := merged.Mapping[variable]

			if exists {

				if !IsEqual(value, otherValue) {

					return false, SolutionSet{} // not compatible
				}

			} else {

				merged.Mapping[variable] = value
			}
		}
	}
	return true, merged
}

func SolutionsAreEqual(solutionA, solutionB SolutionSet) bool {

	for variable, value := range solutionA.Mapping {

		comparedValue, exists := solutionB.Mapping[variable]

		if exists {

			if !IsEqual(value, comparedValue) {

				return false
			}

		} else {

			return false
		}
	}
	return true
}

func SolutionIsDuplicated(solutions []SolutionSet, target SolutionSet) bool {

	for _, solution := range solutions {

		if SolutionsAreEqual(solution, target) {

			return true
		}
	}
	return false
}

func SolveByConstantValue() {

}

func SolveForConstantValue(index int, target, expression *Expression) []SolutionSet {

	symbolType := expression.GetSymbolTypeByIndex(index)

	operands := make([][]int, 0)

	if symbolType == Addition {

		operands = FindAdditives(target.GetNumericValueByIndex(target.GetRoot()))

	} else if symbolType == Multiplication {

		operands = FindFactors(target.GetNumericValueByIndex(target.GetRoot()))

	} else if symbolType == Division {

		operands = FindDividends(target.GetNumericValueByIndex(target.GetRoot()), 5)

	} else {

		return make([]SolutionSet, 0)
	}
	solutions := make([]SolutionSet, 0)

	for _, operandGroup := range operands {

		children := expression.GetChildren(index)

		if len(operandGroup) == len(children) {

			operandGroupAsExpression := ConvertIntToExpression(operandGroup)

			operandCombinations := Expression_GeneratePermutationsOfArray(operandGroupAsExpression)

			for _, operandCombination := range operandCombinations {

				currentSolution := SolutionSet{}

				lowerSolutions := make([]SolutionSet, 0)

				for i := 0; i < len(operandCombination); i++ {

					if expression.IsOperation(children[i]) {

						lowerSolutions = append(lowerSolutions, SolveForConstantValue(children[i], &operandCombination[i], expression)...) // need to merge smaller maps further down

					} else {

						currentSolution.Mapping[expression.GetAlphaValueByIndex(children[i])] = operandCombination[i]
					}
				}
				if len(lowerSolutions) != 0 || len(currentSolution.Mapping) != 0 {

					totalSolutions := MergeMultipleSolutionsOneToMany(lowerSolutions, currentSolution)

					solutions = append(solutions, totalSolutions...)
				}
			}
		}
	}
	return solutions
}

func SolveForMultipleConstantValues(values map[int]string) []SolutionFor {

	solutionsMatrix := make([][]SolutionSet, 0)

	for value, form := range values {

		parsed := ParseExpression(form)

		_, constantAsExpression := NewExpressionWithRoot(Symbol{Constant, value, strconv.Itoa(value)})

		solutionsMatrix = append(solutionsMatrix, SolveForConstantValue(parsed.GetRoot(), &constantAsExpression, &parsed))
	}
	combinations := SolutionSet_GenerateCombinationsByRow(solutionsMatrix)

	compatibleSolutions := make([]SolutionSet, 0)

	for _, combination := range combinations {

		isCompatible, merge := MergeMultipleSolutionsManyToOne(combination)

		if isCompatible && !SolutionIsDuplicated(compatibleSolutions, merge) {

			compatibleSolutions = append(compatibleSolutions, merge)
		}
	}
	// return solutions per initial value somehow
}

func SubstituteSolutionSet(index int, expression *Expression, solution SolutionSet) {

	if expression.IsOperation(index) {

		for _, child := range expression.GetChildren(index) {

			SubstituteSolutionSet(child, expression, solution)
		}

	} else {

		value, exists := solution.Mapping[expression.GetAlphaValueByIndex(index)]

		if exists {

			expression.ReplaceNodeCascade(index, value)
		}
	}
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

func SolutionSet_GenerateCombinationsByRow(matrix [][]SolutionSet) [][]SolutionSet {

	return SolutionSet_GenerateCombinationsByRowRecurse(matrix, make([][]SolutionSet, 0), make([]int, len(matrix)), 0)
}

func SolutionSet_GenerateCombinationsByRowRecurse(matrix [][]SolutionSet, combinations [][]SolutionSet, rowIndexes []int, currentColumn int) [][]SolutionSet {

	for rowNumber := range matrix[currentColumn] {

		rowIndexes[currentColumn] = rowNumber

		if currentColumn == len(matrix)-1 { // if its the last column

			comboPerLine := make([]SolutionSet, 0)

			for colNumber, rowNumber := range rowIndexes {

				comboPerLine = append(comboPerLine, matrix[colNumber][rowNumber])
			}
			combinations = append(combinations, comboPerLine)

		} else {

			combinations = SolutionSet_GenerateCombinationsByRowRecurse(matrix, combinations, rowIndexes, currentColumn+1)
		}
	}
	return combinations
}