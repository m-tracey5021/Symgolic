package evaluation

import (
	"strconv"
	"symgolic/comparison"
	. "symgolic/symbols"
)

func EvaluateFactorisation(index int, expression *Expression) (bool, Expression) {

	if expression.IsSummation(index) {

		factorGroups := make([][]Expression, 0)

		counterPartMaps := make([]map[int]int, 0)

		for _, term := range expression.GetChildren(index) {

			termFactors, termCounterParts := GetTermFactors(term, expression)

			factorGroups = append(factorGroups, termFactors)

			counterPartMaps = append(counterPartMaps, termCounterParts)
		}

		commonFactors, counterParts := GetCommonFactors(expression, factorGroups, counterPartMaps)

		factoredExpressions := make([]Expression, 0)

		for i := 0; i < len(commonFactors); i++ {

			factoredRoot, factored := NewExpressionWithRoot(Symbol{Multiplication, -1, "*"})

			factored.AppendExpression(factoredRoot, commonFactors[i], false)

			add := factored.AppendNode(factoredRoot, Symbol{Addition, -1, "+"})

			for _, counterPart := range counterParts[i] {

				factored.AppendExpression(add, counterPart, false)
			}
		}

		// add the rest of the non factors

		resultRoot, result := NewExpressionWithRoot(Symbol{Vector, -1, "[...]"})

		for _, factored := range factoredExpressions {

			result.AppendExpression(resultRoot, factored, false)
		}

		return true, result

	} else {

		return false, *expression
	}
}

func GetHighestCommonFactor(commonFactors []Expression) Expression {

	largest := 0

	largestIndex := -1

	for i, factor := range commonFactors {

		root := factor.GetRoot()

		value := 1

		if factor.IsMultiplication(root) {

			value, _ = GetTerms(root, &factor)

		} else if factor.IsConstant(root) {

			value = factor.GetNumericValueByIndex(root)
		}
		if value > largest {

			largest = value

			largestIndex = i
		}
	}
	return commonFactors[largestIndex]
}

// func GetCommonFactors(factorGroups [][]Expression, counterPartGroups map[int]int) []Expression {

// 	commonFactors := make([]Expression, 0)

// 	instancesReq := len(factorGroups)

// 	allFactors := make([]Expression, 0)

// 	visited := make([]int, 0)

// 	for _, factorGroup := range factorGroups {

// 		allFactors = append(allFactors, factorGroup...)
// 	}
// 	for i := 0; i < len(allFactors); i++ {

// 		for _, index := range visited {

// 			if i == index {

// 				continue
// 			}
// 		}
// 		instances := 1

// 		for j := i + 1; j < len(allFactors); j++ {

// 			for _, index := range visited {

// 				if i == index {

// 					continue
// 				}
// 			}
// 			if comparison.IsEqualByRoot(allFactors[i], allFactors[j]) {

// 				instances++
// 			}
// 		}
// 		if instances == instancesReq {

// 			commonFactors = append(commonFactors, allFactors[i])
// 		}
// 	}
// 	return commonFactors
// }

func GetCommonFactors(expression *Expression, factorGroups [][]Expression, counterPartMaps []map[int]int) ([]Expression, [][]Expression) {

	commonFactors := make([]Expression, 0)

	counterPartFactors := make([][]Expression, 0)

	instancesReq := len(factorGroups)

	for i, group := range factorGroups {

		for j, factor := range group {

			instances := 1

			counterPartIndexes := make([]int, 0)

			counterPartIndexes = append(counterPartIndexes, counterPartMaps[i][j])

			for k, otherGroup := range factorGroups {

				if k == i {

					continue

				} else {

					for l, otherFactor := range otherGroup {

						if comparison.IsEqualByRoot(factor, otherFactor) {

							instances++

							counterPartIndexes = append(counterPartIndexes, counterPartMaps[k][l])

							continue
						}
					}
				}
			}
			if instances == instancesReq {

				commonFactors = append(commonFactors, factor)

				counterPartsPerCommonFactor := make([]Expression, 0)

				for _, index := range counterPartIndexes {

					counterPartCopy := expression.CopySubtree(index)

					counterPartsPerCommonFactor = append(counterPartsPerCommonFactor, counterPartCopy)
				}

				counterPartFactors = append(counterPartFactors, counterPartsPerCommonFactor)
			}
		}
	}
	return commonFactors, counterPartFactors
}

func GetTermFactors(index int, expression *Expression) ([]Expression, map[int]int) {

	isolatedFactors := make([]Expression, 0)

	copy := expression.CopyTree()

	// get factors of constant add to list

	constantFactors := make([]int, 0)

	if copy.IsMultiplication(index) {

		for _, child := range copy.GetChildren(index) {

			constantFactors = append(constantFactors, GetConstantFactors(child, &copy)...)
		}

	} else {

		constantFactors = GetConstantFactors(index, &copy)
	}

	for _, constantFactor := range constantFactors {

		_, toAdd := NewExpressionWithRoot(Symbol{Constant, constantFactor, strconv.Itoa(constantFactor)})

		isolatedFactors = append(isolatedFactors, toAdd)
	}
	// expand exponents and add to list

	EvaluateAndReplace(index, &copy, EvaluateExponentExpansion)

	for _, child := range copy.GetChildren(index) {

		if !copy.IsConstant(child) {

			isolatedFactors = append(isolatedFactors, copy.CopySubtree(child))
		}
	}
	// get all sublists of isolatedFactors

	factorGroups := GenerateFactorGroups(isolatedFactors, make([]Expression, 0), make([][]Expression, 0), 0)

	factors := make([]Expression, 0)

	// iterate through all sublists, times each tuple together

	for _, group := range factorGroups {

		if len(group) > 1 {

			mul := MultiplyMany(group)

			factors = append(factors, mul)

		} else {

			factors = append(factors, group[0])
		}
		// add to final factors if it equals the target
	}

	// list counterparts for each factor

	counterParts := make(map[int]int)

	for i := 0; i < len(factors); i++ {

		for j := i + 1; j < len(factors); j++ {

			mul := MultiplyTwo(factors[i], factors[j])

			mulRoot := mul.GetRoot()

			if comparison.IsEqual(index, mulRoot, expression, &mul) {

				counterParts[i] = j
			}
		}
	}
	return factors, counterParts
}

func GetConstantFactors(index int, expression *Expression) []int {

	factors := make([]int, 0)

	if expression.IsConstant(index) {

		value := expression.GetNumericValueByIndex(index)

		for i := 1; i <= value; i++ {

			if value%i == 0 {

				factors = append(factors, i)
			}
		}
	}
	return factors
}

func GenerateFactorGroups(factors, output []Expression, factorGroups [][]Expression, index int) [][]Expression {

	if index == len(factors) {

		if len(output) != 0 {

			factorGroups = append(factorGroups, output)
		}
		return factorGroups
	}

	factorGroups = GenerateFactorGroups(factors, output, factorGroups, index+1)

	output = append(output, factors[index])

	factorGroups = GenerateFactorGroups(factors, output, factorGroups, index+1)

	return factorGroups
}
