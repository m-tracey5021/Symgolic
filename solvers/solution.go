package solvers

import (
	. "symgolic/comparison"
	. "symgolic/symbols"
)

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
