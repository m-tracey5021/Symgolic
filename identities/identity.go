package identities

import (
	. "symgolic/comparison"
	"symgolic/evaluation"
	"symgolic/parsing"
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

	} else {

		return false
	}
}

func MergeMaps(A map[string]int, B map[string]int) map[string]int {

	merged := make(map[string]int)

	for variable, constant := range A {

		merged[variable] = constant
	}
	for variable, constant := range B {

		merged[variable] = constant
	}
	return merged
}

func MergeMultipleMaps(merged []map[string]int, toMerge map[string]int) []map[string]int {

	if len(merged) == 0 {

		return []map[string]int{toMerge}
	}
	for _, merge := range merged {

		for key, value := range toMerge {

			merge[key] = value
		}
	}
	return merged
}

func FindVariablesWhere(index int, expression *Expression, target int) []map[string]int {

	symbolType := expression.GetSymbolTypeByIndex(index)

	// somehow work backwards to get variables of a certain structure which can equal the target

	operands := make([][]int, 0)

	if symbolType == Addition {

		operands = evaluation.FindAdditives(target)

	} else if symbolType == Multiplication {

		operands = evaluation.FindFactors(target)

	} else if symbolType == Division {

		operands = evaluation.FindDividends(target, 5)

	} else {

		return make([]map[string]int, 0)
	}
	variableMaps := make([]map[string]int, 0)

	for _, operandGroup := range operands {

		children := expression.GetChildren(index)

		currentMap := make(map[string]int)

		lowerMaps := make([]map[string]int, 0)

		if len(operandGroup) == len(children) {

			for i := 0; i < len(operandGroup); i++ {

				if expression.IsOperation(children[i]) {

					lowerMaps = append(lowerMaps, FindVariablesWhere(children[i], expression, operandGroup[i])...) // need to merge smaller maps further down

				} else {

					currentMap[expression.GetAlphaValueByIndex(children[i])] = operandGroup[i]
				}

			}

		}
		totalMaps := MergeMultipleMaps(lowerMaps, currentMap)

		variableMaps = append(variableMaps, totalMaps...)
	}
	return variableMaps
}

func SelectCompatibleMappings(A, B ConstantMapByForm) []map[string]int {

	// return all combinations of possible mappings where there are no contradictory entries

	totalPossibleMappings := make([]map[string]int, 0)

	for i := 0; i < len(A.PossibleMappings); i++ {

		for j := 0; j < len(B.PossibleMappings); j++ {

			compatible := true

			for variable, constant := range A.PossibleMappings[i] {

				comparedConstant, exists := B.PossibleMappings[j][variable]

				if exists && comparedConstant != constant {

					compatible = false
				}
			}
			if compatible {

				merged := MergeMaps(A.PossibleMappings[i], B.PossibleMappings[j])

				totalPossibleMappings = append(totalPossibleMappings, merged)
			}
		}
	}
	return totalPossibleMappings
}

func GenerateCompatibleConstantMaps(values map[int]string) {

	constantMaps := make([]ConstantMapByForm, 0)

	for value, form := range values {

		parsed := parsing.ParseExpression(form)

		constantMaps = append(constantMaps,

			ConstantMapByForm{

				Value: value,

				Form: form,

				PossibleMappings: FindVariablesWhere(parsed.GetRoot(), &parsed, value),
			},
		)
	}
	for i := 0; i < len(constantMaps); i++ {

		for j := i + 1; j < len(constantMaps); j++ {

		}
	}
}

func GenerateConstantMapCombos(constantMaps []ConstantMapByForm) [][]map[string]int {

	indexes := make([]int, len(constantMaps))

	for i := 0; i < len(indexes); i++ {

		indexes[i] = 0
	}
	// indexToIncrement := 0

	// totalMapCombos := make([][]map[string]int, 0)

	// for _, constantMap := range constantMaps {

	// 	for i := 0; i < len(constantMap.PossibleMappings); i++ {

	// 	}
	// }
	return nil
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
