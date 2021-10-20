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

// func GetHighestCommonFactor(commonFactors []Expression) Expression {

// 	largest := 0

// 	largestIndex := -1

// 	for i, factor := range commonFactors {

// 		root := factor.GetRoot()

// 		value := 1

// 		if factor.IsMultiplication(root) {

// 			value, _ = GetTerms(root, &factor)

// 		} else if factor.IsConstant(root) {

// 			value = factor.GetNumericValueByIndex(root)
// 		}
// 		if value > largest {

// 			largest = value

// 			largestIndex = i
// 		}
// 	}
// 	return commonFactors[largestIndex]
// }

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

			counterParts := make([]Expression, 0)

			counterParts = append(counterParts, group[counterPartMaps[i][j]].CopyTree())

			for k, otherGroup := range factorGroups {

				if k == i {

					continue

				} else {

					for l, otherFactor := range otherGroup {

						if comparison.IsEqualByRoot(factor, otherFactor) {

							instances++

							counterParts = append(counterParts, otherGroup[counterPartMaps[k][l]].CopyTree())

							continue
						}
					}
				}
			}
			if instances == instancesReq {

				commonFactors = append(commonFactors, factor)

				counterPartFactors = append(counterPartFactors, counterParts)
			}
		}
	}
	return commonFactors, counterPartFactors
}

func GetTermFactors(index int, expression *Expression) ([]Expression, map[int]int) {

	isolatedFactors := make([]Expression, 0)

	copy := expression.CopyTree()

	// expand exponents

	EvaluateAndReplace(index, &copy, EvaluateExponentExpansion)

	// get factors of constant add to list

	largestConstantFactor, isolatedFactors := GetIsolatedFactors(index, expression)

	// get all sublists of isolatedFactors

	factorGroups := GenerateFactorGroups(isolatedFactors, make([]Expression, 0), make([][]Expression, 0), 0)

	factors := make([]Expression, 0)

	// iterate through all sublists, times each tuple together

	for _, group := range factorGroups {

		var factorToAdd Expression

		if len(group) > 1 {

			factorToAdd = MultiplyMany(group)

		} else {

			factorToAdd = group[0]
		}
		// make sure constants are not times together to be bigger than the initial value

		if !ExceedsLargestConstantFactor(largestConstantFactor, factorToAdd) {

			exists := false

			for _, factor := range factors {

				if comparison.IsEqualByRoot(factor, factorToAdd) {

					exists = true

					break
				}
			}
			if !exists {

				factors = append(factors, factorToAdd)
			}
			// add to final factors if it equals the target
		}
	}

	// list counterparts for each factor

	counterParts := make(map[int]int)

	for i := 0; i < len(factors); i++ {

		for j := i + 1; j < len(factors); j++ {

			mul := MultiplyTwo(factors[i], factors[j])

			mulRoot := mul.GetRoot()

			if comparison.IsEqualAt(index, mulRoot, expression, &mul) {

				counterParts[i] = j
			}
		}
	}
	return factors, counterParts
}

func ExceedsLargestConstantFactor(largestFactor int, compared Expression) bool {

	root := compared.GetRoot()

	if compared.IsConstant(root) {

		if compared.GetNumericValueByIndex(root) > largestFactor {

			return true

		} else {

			return false
		}

	} else if compared.IsMultiplication(root) {

		for _, child := range compared.GetChildren(root) {

			if compared.GetNumericValueByIndex(child) > largestFactor {

				return true

			} else {

				return false
			}
		}
		return false

	} else {

		return false
	}
}

func IntegerToConstantExpression(value int) Expression {

	_, expression := NewExpressionWithRoot(Symbol{Constant, value, strconv.Itoa(value)})

	return expression
}

func IntegerFactorsToConstantExpression(integerFactors []int) []Expression {

	factors := make([]Expression, 0)

	for _, integer := range integerFactors {

		factors = append(factors, IntegerToConstantExpression(integer))
	}
	return factors
}

func GetConstantFactors(value int) []int {

	factors := make([]int, 0)

	for i := 1; i <= value; i++ {

		if value%i == 0 {

			factors = append(factors, i)
		}
	}
	return factors
}

func GetIsolatedFactors(index int, expression *Expression) (int, []Expression) {

	value := -1

	factors := make([]Expression, 0)

	if expression.IsConstant(index) {

		value = expression.GetNumericValueByIndex(index)

		constantFactors := GetConstantFactors(index)

		factors = append(factors, IntegerFactorsToConstantExpression(constantFactors)...)

	} else if expression.IsMultiplication(index) {

		for _, child := range expression.GetChildren(index) {

			innerValue := expression.GetNumericValueByIndex(child)

			if innerValue != -1 {

				innerConstantFactors := GetConstantFactors(innerValue)

				factors = append(factors, IntegerFactorsToConstantExpression(innerConstantFactors)...)

				if innerValue > value {

					value = innerValue
				}

			} else {

				factors = append(factors, expression.CopySubtree(child))
			}

		}

	} else {

		factors = append(factors, expression.CopySubtree(index))
	}
	return value, factors
}

// func AppendToConstantFactors(child, largestConstant int, copy *Expression, constantFactors []int) (int, []int) {

// 	constant, constantFactorsToAdd := GetConstantFactors(child, copy)

// 	constantFactors = append(constantFactors, constantFactorsToAdd...)

// 	if constant > largestConstant {

// 		largestConstant = constant
// 	}
// 	return largestConstant, constantFactors
// }

func GenerateFactorGroups(factors, output []Expression, factorGroups [][]Expression, index int) [][]Expression {

	if index == len(factors) {

		if len(output) != 0 {

			// exists := false

			// for _, group := range factorGroups {

			// 	if len(group) == len(output) {

			// 		match := true

			// 		for i := 0; i < len(group); i++ {

			// 			if !comparison.IsEqualByRoot(group[i], output[i]) {

			// 				match = false
			// 			}
			// 		}
			// 		if match {

			// 			exists = true

			// 			break
			// 		}
			// 	}
			// }
			// if !exists {

			// 	factorGroups = append(factorGroups, output)
			// }
			factorGroups = append(factorGroups, output)
		}
		return factorGroups
	}

	factorGroups = GenerateFactorGroups(factors, output, factorGroups, index+1)

	output = append(output, factors[index])

	factorGroups = GenerateFactorGroups(factors, output, factorGroups, index+1)

	return factorGroups
}
