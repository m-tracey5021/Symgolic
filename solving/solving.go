package solving

import (
	. "symgolic/evaluation"
	. "symgolic/generation"
	. "symgolic/symbols"
)

type SolutionSet struct {
	Mapping map[string]int
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

				if value != otherValue {

					return false, SolutionSet{} // not compatible
				}

			} else {

				merged.Mapping[variable] = value
			}
		}
	}
	return true, merged
}

func SolveByConstantValue(target, index int, expression *Expression) []SolutionSet {

	symbolType := expression.GetSymbolTypeByIndex(index)

	// somehow work backwards to get variables of a certain structure which can equal the target

	operands := make([][]int, 0)

	if symbolType == Addition {

		operands = FindAdditives(target)

	} else if symbolType == Multiplication {

		operands = FindFactors(target)

	} else if symbolType == Division {

		operands = FindDividends(target, 5)

	} else {

		return make([]SolutionSet, 0)
	}
	solutions := make([]SolutionSet, 0)

	for _, operandGroup := range operands {

		children := expression.GetChildren(index)

		if len(operandGroup) == len(children) {

			// need to rearrange this list operandGroup
			// at this point and go through each rearranged list

			operandCombinations := GeneratePermutationsOfArray(operandGroup)

			for _, operandCombination := range operandCombinations {

				currentSolution := SolutionSet{}

				lowerSolutions := make([]SolutionSet, 0)

				for i := 0; i < len(operandCombination); i++ {

					if expression.IsOperation(children[i]) {

						lowerSolutions = append(lowerSolutions, SolveByConstantValue(operandCombination[i], children[i], expression)...) // need to merge smaller maps further down

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
