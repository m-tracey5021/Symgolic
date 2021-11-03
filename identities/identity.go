package identities

import (
	. "symgolic/comparison"
	. "symgolic/parsing"
	. "symgolic/symbols"
)

type Identity interface {
	AssignVariables()

	ApplyForwards()

	ApplyBackwards()

	Run()
}

type Assignment func(map[string]Expression, Direction)

type IdentityRequisite struct {
	Form string

	Direction Direction

	ConstantChecks []ConstantCheck
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

func Identify(index int, expression *Expression, identityRequisites []IdentityRequisite, assignment Assignment) bool {

	for _, requisite := range identityRequisites {

		form := ParseExpression(requisite.Form)

		formApplies, variableMap := IsEqualByForm(form, *expression)

		if formApplies {

			if len(requisite.ConstantChecks) != 0 {

				for _, check := range requisite.ConstantChecks {

					if CheckConstantValue(check.Values, check.Target, check.Operation, expression) {

						// assign indexes to struct

						assignment(variableMap, requisite.Direction)

						return true
					}
				}

			} else {

				assignment(variableMap, requisite.Direction)

				return true
			}

		} else {

			return false
		}
	}
	return false
}

func CheckConstantValue(indices []int, targetIndex int, operation SymbolType, expression *Expression) bool {

	values := make([]int, 0)

	target := expression.GetNumericValueByIndex(targetIndex)

	for _, index := range indices {

		value := expression.GetNumericValueByIndex(index)

		if value == -1 {

			return false

		} else {

			values = append(values, value)
		}
	}
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

	} else {

		return false
	}
}
