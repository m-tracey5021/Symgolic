package solvers

import (
	. "symgolic/comparison"
	. "symgolic/symbols"
)

type SolveRequest struct {
	Value Expression

	Given Expression
}

type SolutionSet struct {
	Mapping map[string]Expression
}

type SolutionFor struct {
	Value Expression

	Given Expression

	Solutions []SolutionSet
}

type SolutionContext struct {
	Expression Expression

	SolutionsForValues []SolutionFor

	SolutionsOverValues []SolutionSet
}

func NewSolutionSet() SolutionSet {

	return SolutionSet{Mapping: make(map[string]Expression, 0)}
}

func NewSolutionFor() SolutionFor {

	return SolutionFor{Solutions: make([]SolutionSet, 0)}
}

func NewSolutionContext() SolutionContext {

	return SolutionContext{SolutionsForValues: make([]SolutionFor, 0)}
}

func (s *SolutionFor) AppendSolutionSets(solutionSets []SolutionSet) {

	for _, compared := range solutionSets {

		exists := false

		for _, nthSolutionSet := range s.Solutions {

			if SolutionsAreEqual(nthSolutionSet, compared) {

				exists = true

				break
			}
		}
		if exists {

			continue

		} else {

			s.Solutions = append(s.Solutions, compared)
		}
	}
}

func (s *SolutionContext) AppendSolutionFor(solution SolutionFor) {

	for i, nthSolution := range s.SolutionsForValues {

		if IsEqual(solution.Value, nthSolution.Value) {

			nthSolution.AppendSolutionSets(solution.Solutions)

			s.SolutionsForValues[i] = nthSolution

			return
		}
	}
	s.SolutionsForValues = append(s.SolutionsForValues, solution)
}

func (s *SolutionContext) GetValuesFor(variable string) []Expression {

	values := make([]Expression, 0)

	for _, solutionFor := range s.SolutionsForValues {

		for _, solutionSet := range solutionFor.Solutions {

			value, exists := solutionSet.Mapping[variable]

			if exists {

				add := true

				for _, preexisting := range values {

					if IsEqual(preexisting, value) {

						add = false

						break
					}
				}
				if add {

					values = append(values, value)
				}
			}
		}
	}
	return values
}

func MergeSolutions(A SolutionSet, B SolutionSet) SolutionSet {

	merged := NewSolutionSet()

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

	merged := NewSolutionSet()

	for _, solution := range toMerge {

		for variable, value := range solution.Mapping {

			otherValue, exists := merged.Mapping[variable]

			if exists {

				if !IsEqual(value, otherValue) {

					return false, merged // not compatible
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

func SubstituteSolutionSet(index int, expression *Expression, solution SolutionSet) {

	if expression.IsOperation(index) {

		for _, child := range expression.GetChildren(index) {

			SubstituteSolutionSet(child, expression, solution)
		}

	} else {

		value, exists := solution.Mapping[expression.GetNode(index).AlphaValue]

		if exists {

			expression.ReplaceNodeCascade(index, value)
		}
	}
}

func GenerateCompatibleSolutionContext(solutionsForValues []SolutionFor) SolutionContext {

	return GenerateCompatibleSolutionContextRecurse(solutionsForValues, NewSolutionContext(), make([]int, len(solutionsForValues)), 0)
}

func GenerateCompatibleSolutionContextRecurse(solutionsForValues []SolutionFor, context SolutionContext, rowIndexes []int, currentColumn int) SolutionContext {

	for rowNumber := range solutionsForValues[currentColumn].Solutions {

		rowIndexes[currentColumn] = rowNumber

		if currentColumn == len(solutionsForValues)-1 { // if its the last column

			lineCombination := make([]SolutionSet, 0)

			for colNumber, rowNumber := range rowIndexes {

				lineCombination = append(lineCombination, solutionsForValues[colNumber].Solutions[rowNumber])
			}
			isCompatible, solutionOverValues := MergeMultipleSolutionsManyToOne(lineCombination)

			if isCompatible {

				context.SolutionsOverValues = append(context.SolutionsOverValues, solutionOverValues)

				for colNumber, rowNumber := range rowIndexes {

					solutionToAppend := SolutionFor{

						Value: solutionsForValues[colNumber].Value,

						Given: solutionsForValues[colNumber].Given,

						Solutions: make([]SolutionSet, 0),
					}
					solutionToAppend.Solutions = append(solutionToAppend.Solutions, solutionsForValues[colNumber].Solutions[rowNumber])

					context.AppendSolutionFor(solutionToAppend)
				}
			}

		} else {

			context = GenerateCompatibleSolutionContextRecurse(solutionsForValues, context, rowIndexes, currentColumn+1)
		}
	}
	return context
}
