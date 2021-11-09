package identities

import (
	. "symgolic/comparison"
	. "symgolic/evaluation"
	. "symgolic/generation"
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

	} else if operation == Equality {

		return values[0] == target

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

func MergeMultipleMapsOneToMany(merged []map[string]int, toMerge map[string]int) []map[string]int {

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

func MergeMultipleMapsManyToOne(toMerge []map[string]int) (bool, map[string]int) {

	merged := make(map[string]int)

	for _, mapping := range toMerge {

		for variable, value := range mapping {

			otherValue, exists := merged[variable]

			if exists {

				if value != otherValue {

					return false, nil // not compatible
				}

			} else {

				merged[variable] = value
			}
		}
	}
	return true, merged
}

func MappingsAreEqual(mappingA, mappingB map[string]int) bool {

	for variable, value := range mappingA {

		comparedValue, exists := mappingB[variable]

		if exists {

			if value != comparedValue {

				return false
			}

		} else {

			return false
		}
	}
	return true
}

func MappingIsDuplicated(target map[string]int, mappings []map[string]int) bool {

	for _, mapping := range mappings {

		if MappingsAreEqual(mapping, target) {

			return true
		}
	}
	return false
}

func FindConstantMapByForm(index int, expression *Expression, target int) []map[string]int {

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

		return make([]map[string]int, 0)
	}
	variableMaps := make([]map[string]int, 0)

	for _, operandGroup := range operands {

		children := expression.GetChildren(index)

		if len(operandGroup) == len(children) {

			// need to rearrange this list operandGroup
			// at this point and go through each rearranged list

			operandCombinations := GeneratePermutationsOfArray(operandGroup)

			for _, operandCombination := range operandCombinations {

				currentMap := make(map[string]int)

				lowerMaps := make([]map[string]int, 0)

				for i := 0; i < len(operandCombination); i++ {

					if expression.IsOperation(children[i]) {

						lowerMaps = append(lowerMaps, FindConstantMapByForm(children[i], expression, operandCombination[i])...) // need to merge smaller maps further down

					} else {

						currentMap[expression.GetAlphaValueByIndex(children[i])] = operandCombination[i]
					}
				}
				if len(lowerMaps) != 0 || len(currentMap) != 0 {

					totalMaps := MergeMultipleMapsOneToMany(lowerMaps, currentMap)

					variableMaps = append(variableMaps, totalMaps...)
				}
			}
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

func GenerateCompatibleConstantMapsForValues(values map[int]string) []map[string]int {

	constantMaps := make([]ConstantMapByForm, 0)

	for value, form := range values {

		parsed := ParseExpression(form)

		constantMaps = append(constantMaps,

			ConstantMapByForm{

				Value: value,

				Form: form,

				PossibleMappings: FindConstantMapByForm(parsed.GetRoot(), &parsed, value),
			},
		)
	}
	combinations := GenerateConstantMapCombinations(constantMaps)

	compatibleCombinations := make([]map[string]int, 0)

	for _, combination := range combinations {

		compatible, merge := MergeMultipleMapsManyToOne(combination)

		if compatible && !MappingIsDuplicated(merge, compatibleCombinations) {

			compatibleCombinations = append(compatibleCombinations, merge)
		}
	}
	return compatibleCombinations
}

func GenerateConstantChecksForValues(unknownValues, valuesEqualToUnknown map[int]string, knownValues []ConstantCheck) []ConstantCheck {

	constantChecks := make([]ConstantCheck, 0)

	compatibleCombinations := GenerateCompatibleConstantMapsForValues(unknownValues)

	for _, combination := range compatibleCombinations {

		for unknownValue, form := range unknownValues {

			check := ConstantCheck{

				Values: []int{}
			}
		}
	}

	for _, check := range knownValues {

		constantChecks = append(constantChecks, check)
	}
	return constantChecks
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
