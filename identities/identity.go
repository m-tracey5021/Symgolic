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

func FindConstantMapByForm(index int, expression *Expression, target int) []map[string]int {

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

			// need to rearrange this list operandGroup
			// at this point and go through each rearranged list

			for i := 0; i < len(operandGroup); i++ {

				if expression.IsOperation(children[i]) {

					lowerMaps = append(lowerMaps, FindConstantMapByForm(children[i], expression, operandGroup[i])...) // need to merge smaller maps further down

				} else {

					currentMap[expression.GetAlphaValueByIndex(children[i])] = operandGroup[i]
				}
			}
		}
		if len(lowerMaps) != 0 || len(currentMap) != 0 {

			totalMaps := MergeMultipleMaps(lowerMaps, currentMap)

			variableMaps = append(variableMaps, totalMaps...)
		}
	}
	return variableMaps
}

func GenerateConstantMapCombinations(constantMaps []ConstantMapByForm) [][]map[string]int {

	return GenerateConstantMapCombinationsRecurse(constantMaps, make([][]map[string]int, 0), make([]int, len(constantMaps)), 0)
}

func GenerateConstantMapCombinationsRecurse(constantMaps []ConstantMapByForm, combinations [][]map[string]int, rowIndexes []int, currentColumn int) [][]map[string]int {

	for rowNumber := range constantMaps[currentColumn].PossibleMappings {

		rowIndexes[currentColumn] = rowNumber

		if currentColumn == len(constantMaps)-1 { // if its the last column

			comboPerLine := make([]map[string]int, 0)

			for colNumber, rowNumber := range rowIndexes {

				comboPerLine = append(comboPerLine, constantMaps[colNumber].PossibleMappings[rowNumber])
			}
			combinations = append(combinations, comboPerLine)

		} else {

			combinations = GenerateConstantMapCombinationsRecurse(constantMaps, combinations, rowIndexes, currentColumn+1)
		}
	}
	return combinations
}

func ConstantMapCombinationIsCompatible(constantMaps []map[string]int) bool {

	if len(constantMaps) != 1 {

		initial := constantMaps[0]

		for i := 1; i < len(constantMaps); i++ {

			for variable, value := range initial {

				comparedValue, exists := constantMaps[i][variable]

				if exists {

					if value != comparedValue {

						return false
					}
				}
			}
		}
		return true

	} else {

		return true
	}
}

func GenerateCompatibleConstantMapsForValues(values map[int]string) [][]map[string]int {

	constantMaps := make([]ConstantMapByForm, 0)

	for value, form := range values {

		parsed := parsing.ParseExpression(form)

		constantMaps = append(constantMaps,

			ConstantMapByForm{

				Value: value,

				Form: form,

				PossibleMappings: FindConstantMapByForm(parsed.GetRoot(), &parsed, value),
			},
		)
	}
	combinations := GenerateConstantMapCombinations(constantMaps)

	compatibleCombinations := make([][]map[string]int, 0)

	for _, combination := range combinations {

		if ConstantMapCombinationIsCompatible(combination) {

			compatibleCombinations = append(compatibleCombinations, combination)
		}
	}
	return compatibleCombinations
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
