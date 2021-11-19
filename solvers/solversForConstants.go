package solvers

import (
	. "symgolic/evaluation"
	. "symgolic/symbols"
)

func SolveForConstantValue(index int, target, expression *Expression) (bool, []SolutionSet) {

	symbolType := expression.GetSymbolTypeByIndex(index)

	operands := make([][]int, 0)

	if symbolType == Addition {

		operands = FindAdditives(target.GetNode(target.GetRoot()).NumericValue)

	} else if symbolType == Multiplication {

		operands = FindFactors(target.GetNode(target.GetRoot()).NumericValue)

	} else if symbolType == Division {

		operands = FindDividends(target.GetNode(target.GetRoot()).NumericValue, 5)

	} else if symbolType == Variable && expression.GetParent(index) == -1 {

		return true, []SolutionSet{

			{
				Mapping: map[string]Expression{

					expression.GetNode(index).AlphaValue: *target,
				},
			},
		}

	} else if symbolType == Constant {

		if target.GetNode(target.GetRoot()).NumericValue != expression.GetNode(index).NumericValue {

			return false, make([]SolutionSet, 0)

		} else {

			return true, make([]SolutionSet, 0)
		}

	} else {

		return true, make([]SolutionSet, 0)
	}
	solutions := make([]SolutionSet, 0)

	for _, operandGroup := range operands {

		children := expression.GetChildren(index)

		if len(operandGroup) == len(children) {

			operandGroupAsExpression := ConvertIntToExpression(operandGroup)

			operandCombinations := Expression_GeneratePermutationsOfArray(operandGroupAsExpression)

			for _, operandCombination := range operandCombinations {

				currentSolution := NewSolutionSet()

				lowerSolutions := make([]SolutionSet, 0)

				solutionExists := true

				for i := 0; i < len(operandCombination); i++ {

					if expression.IsOperation(children[i]) || expression.IsConstant(children[i]) {

						solutionExistsForChild, solutionsForChild := SolveForConstantValue(children[i], &operandCombination[i], expression)

						if solutionExistsForChild {

							lowerSolutions = append(lowerSolutions, solutionsForChild...) // need to merge smaller maps further down

						} else {

							solutionExists = false

							break
						}

					} else {

						currentSolution.Mapping[expression.GetNode(children[i]).AlphaValue] = operandCombination[i]
					}
				}
				if (len(lowerSolutions) != 0 || len(currentSolution.Mapping) != 0) && solutionExists {

					totalSolutions := MergeMultipleSolutionsOneToMany(lowerSolutions, currentSolution)

					solutions = append(solutions, totalSolutions...)
				}
			}
		}
	}
	return len(solutions) != 0, solutions
}

func SolveForMultipleConstantValues(values []SolveRequest) SolutionContext {

	solutionsForValues := make([]SolutionFor, 0)

	for _, request := range values {

		_, solutions := SolveForConstantValue(request.Given.GetRoot(), &request.Value, &request.Given)

		solutionsForValues = append(solutionsForValues, SolutionFor{

			Value: request.Value,

			Given: request.Given,

			Solutions: solutions,
		})
	}
	return GenerateCompatibleSolutionContext(solutionsForValues)
}
