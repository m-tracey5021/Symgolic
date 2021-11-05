package identities

import (
	. "symgolic/comparison"
	. "symgolic/parsing"
	. "symgolic/symbols"
)

type IIdentity interface {
	AssignVariables(variableMap map[string]Expression, direction Direction)

	ApplyForwards(index int, expression *Expression) Expression

	ApplyBackwards(index int, expression *Expression) Expression

	GetRequisites() []IdentityRequisite

	GetDirection() Direction
}

// type IdentityBase struct {
// 	Direction Direction

// 	IdentityRequisites []IdentityRequisite
// }

// type Assignment func(map[string]Expression, Direction)

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

// func Identify(index int, expression *Expression, identityRequisites []IdentityRequisite, assignment Assignment) bool {

// 	for _, requisite := range identityRequisites {

// 		form := ParseExpression(requisite.Form)

// 		formApplies, variableMap := IsEqualByForm(form, *expression)

// 		if formApplies {

// 			if len(requisite.ConstantChecks) != 0 {

// 				for _, check := range requisite.ConstantChecks {

// 					if CheckConstantValue(check.Values, check.Target, check.Operation, expression) {

// 						// assign indexes to struct

// 						assignment(variableMap, requisite.Direction)

// 						return true
// 					}
// 				}

// 			} else {

// 				assignment(variableMap, requisite.Direction)

// 				return true
// 			}

// 		} else {

// 			return false
// 		}
// 	}
// 	return false
// }

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

	} else {

		return false
	}
}

func Identify(index int, expression *Expression, identity IIdentity) bool {

	for _, requisite := range identity.GetRequisites() {

		form := ParseExpression(requisite.Form)

		formApplies, variableMap := IsEqualByForm(form, *expression)

		if formApplies {

			if len(requisite.ConstantChecks) != 0 {

				for _, check := range requisite.ConstantChecks {

					if CheckConstantValue(check.Values, check.Target, check.Operation, expression) {

						// assign indexes to struct

						identity.AssignVariables(variableMap, requisite.Direction)

						return true
					}
				}

			} else {

				identity.AssignVariables(variableMap, requisite.Direction)

				return true
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
