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

	AlternateForms []AlternateForm
}

type Direction int

const (
	Forwards = iota

	Backwards
)

type FormCondition struct {
	Target Expression

	EqualTo Expression

	Instances [][]int
}

type AlternateForm struct {
	Form string

	Conditions []FormCondition // map of indexes to forms, the variable at that index must be equal to this form
}

func GenerateAlternateForms(index int, form, expression *Expression) {

	for _, child := range form.GetChildren(index) {

		GenerateAlternateForms(child, form, expression)
	}
	if form.IsSummation(index) {

		condition := FormCondition{

			Target: expression.CopySubtree(index),

			EqualTo: form.CopySubtree(index),

			Instances: [][]int{

				{index},
			},
		}
	}
}

func GetSolutionContextForConditions(alternative AlternateForm) SolutionContext {

	unknownValues := make([]SolveRequest, 0)

	for _, condition := range alternative.Conditions {

		unknownValues = append(unknownValues, SolveRequest{Value: condition.Target, Given: condition.EqualTo})
	}
	return SolveForMultipleConstantValues(unknownValues)

}

func ApplyConditions(form, expression Expression, alternative AlternateForm, solution SolutionSet) (bool, map[string]Expression) {

	expanded := expression.CopyTree()

	for _, condition := range alternative.Conditions {

		replacement := condition.EqualTo.CopyTree()

		SubstituteSolutionSet(replacement.GetRoot(), &replacement, solution)

		for _, instance := range condition.Instances {

			index := expanded.GetNodeByPath(instance)

			expanded.ReplaceNodeCascade(index, replacement.CopyTree())
		}
	}
	return IsEqualByForm(form, expanded)
}

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

					solutionContext := GetSolutionContextForConditions(alternative)

					for _, solution := range solutionContext.SolutionsOverValues {

						formApplies, variableMap := ApplyConditions(form, *expression, alternative, solution)

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
