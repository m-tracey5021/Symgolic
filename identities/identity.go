package identities

import (
	. "symgolic/comparison"
	. "symgolic/parsing"
	. "symgolic/solvers"
	. "symgolic/symbols"
)

type IIdentity interface {
	AssignVariables(variableMap map[string]Expression, direction Direction)

	ApplyForwards(index int, expression *Expression) Expression

	ApplyBackwards(index int, expression *Expression) Expression

	GetRequisites() []IdentityRequisite

	GetDirection() Direction
}

type IdentityRequisite struct {
	Form string

	Direction Direction

	ConstantChecks []ConstantCheck

	Expansions SolutionContext

	AlternateForms []AlternateForm
}

type Direction int

const (
	Forwards = iota

	Backwards
)

type ConstantCheck struct {
	Values []int

	Target int

	Operation SymbolType
}

type AlternateForm struct {
	Form string

	UnknownValues map[int]string

	Replacements []ReplacementCommand

	Conditions []FormCondition // map of indexes to forms, the variable at that index must be equal to this form
}

type FormCondition struct {
	Target Expression

	EqualTo Expression

	Instances [][]int
}

type ReplacementCommand struct {
	Indexes []int

	ReplacementForm string
}

func NewConstantCheck() ConstantCheck {

	return ConstantCheck{Values: make([]int, 0)}
}

type ConstantMapByForm struct {
	Value int

	Form string

	PossibleMappings []map[string]int
}

func CheckConstantValue(values []int, target int, operation SymbolType, expression *Expression) bool {

	if operation == Addition {

		total := 0

		for _, value := range values {

			total += value
		}
		return total == target

	} else if operation == Multiplication {

		total := 1

		for _, value := range values {

			total *= value
		}
		return total == target

	} else if operation == Equality {

		return values[0] == target

	} else {

		return false
	}
}

func GenerateExpansions(unknownValues map[int]string, knownValues map[string]Expression) []ConstantCheck {

	constantChecks := make([]ConstantCheck, 0)

	// solutionContext := SolveForMultipleConstantValues(unknownValues)

	// knownSolutionSet := SolutionSet{Mapping: knownValues}

	return constantChecks
}

func GenerateConstantCheckForExpandedForm(expandedForm string, valueIndexes []int, targetIndex int, operation SymbolType) ConstantCheck {

	check := NewConstantCheck()

	expanded := ParseExpression(expandedForm)

	for _, valueIndex := range valueIndexes {

		check.Values = append(check.Values, expanded.GetNumericValueByIndex(valueIndex))
	}
	check.Target = expanded.GetNumericValueByIndex(targetIndex)

	check.Operation = operation

	return check
}

// func Identify(index int, expression *Expression, identity IIdentity) bool {

// 	for _, requisite := range identity.GetRequisites() {

// 		form := ParseExpression(requisite.Form)

// 		formApplies, variableMap := IsEqualByForm(form, *expression)

// 		if formApplies {

// 			if len(requisite.ConstantChecks) != 0 {

// 				for _, check := range requisite.ConstantChecks {

// 					if CheckConstantValue(check.Values, check.Target, check.Operation, expression) {

// 						// assign indexes to struct

// 						identity.AssignVariables(variableMap, requisite.Direction)

// 						return true
// 					}
// 				}

// 			} else {

// 				identity.AssignVariables(variableMap, requisite.Direction)

// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

func Identify(index int, expression *Expression, identity IIdentity) bool {

	for _, requisite := range identity.GetRequisites() {

		form := ParseExpression(requisite.Form)

		formApplies, variableMap := IsEqualByForm(form, *expression)

		if formApplies {

			identity.AssignVariables(variableMap, requisite.Direction)

			return true

		} else {

			if len(requisite.AlternateForms) != 0 {

				for _, alternative := range requisite.AlternateForms {

					unknownValues := make([]SolveRequest, 0)

					for _, condition := range alternative.Conditions {

						unknownValues = append(unknownValues, SolveRequest{Value: condition.Target, Given: condition.EqualTo})
					}
					solutionContext := SolveForMultipleConstantValues(unknownValues)

					for _, solution := range solutionContext.SolutionsOverValues {

						copy := expression.CopyTree()

						for _, condition := range alternative.Conditions {

							replacement := condition.EqualTo.CopyTree()

							SubstituteSolutionSet(replacement.GetRoot(), &replacement, solution)

							for _, instance := range condition.Instances {

								index := expression.GetChildByPath(instance)

								copy.ReplaceNodeCascade(index, replacement.CopyTree())
							}
						}
						formApplies, variableMap := IsEqualByForm(form, copy)

						if formApplies {

							identity.AssignVariables(variableMap, requisite.Direction)

							return true
						}
					}
				}
			}
		}
	}
	return false
}

func Run(index int, expression *Expression, identity IIdentity) (bool, Expression) {

	if Identify(index, expression, identity) {

		if identity.GetDirection() == Forwards {

			return true, identity.ApplyForwards(index, expression)

		} else {

			return true, identity.ApplyBackwards(index, expression)
		}

	} else {

		return false, *expression
	}
}
