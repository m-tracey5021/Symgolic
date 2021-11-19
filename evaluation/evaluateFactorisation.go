package evaluation

import (
	"strconv"
	"symgolic/comparison"
	. "symgolic/symbols"
)

type CommonFactorCombination struct {
	CommonFactor Expression

	CounterParts []Expression
}

type TermFactorCombination struct {
	Factor Expression

	CounterPart Expression
}

func EvaluateFactorisation(index int, expression *Expression) (bool, Expression) {

	if expression.IsSummation(index) {

		// factorGroups := make([][]Expression, 0)

		// counterPartMaps := make([]map[int]int, 0)

		termFactorGroups := make([][]TermFactorCombination, 0)

		for _, term := range expression.GetChildren(index) {

			termFactors := GetTermFactors(term, expression)

			termFactorGroups = append(termFactorGroups, termFactors)

			// factorGroups = append(factorGroups, termFactors)

			// counterPartMaps = append(counterPartMaps, termCounterParts)
		}

		commonFactorsGroups := GetCommonFactors(expression, termFactorGroups)

		factoredExpressions := make([]Expression, 0)

		for i := 0; i < len(commonFactorsGroups); i++ {

			factoredRoot, factored := NewExpression(Symbol{Multiplication, -1, "*"})

			factored.AppendExpression(factoredRoot, commonFactorsGroups[i].CommonFactor, false)

			add := factored.AppendNode(factoredRoot, Symbol{Addition, -1, "+"})

			for _, counterPart := range commonFactorsGroups[i].CounterParts {

				factored.AppendExpression(add, counterPart, false)
			}
			factoredExpressions = append(factoredExpressions, factored)
		}

		// add the rest of the non factors

		resultRoot, result := NewExpression(Symbol{Vector, -1, "[...]"})

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

func GetCommonFactors(expression *Expression, termFactorGroups [][]TermFactorCombination) []CommonFactorCombination {

	// commonFactors := make([]Expression, 0)

	// counterPartFactors := make([][]Expression, 0)

	commonFactorCombos := make([]CommonFactorCombination, 0)

	instancesReq := len(termFactorGroups)

	for i, group := range termFactorGroups {

		for _, factor := range group {

			instances := 1

			counterParts := make([]Expression, 0)

			counterParts = append(counterParts, factor.CounterPart)

			for k, otherGroup := range termFactorGroups {

				if k == i {

					continue

				} else {

					for _, otherFactor := range otherGroup {

						if comparison.IsEqual(factor.Factor, otherFactor.Factor) {

							instances++

							counterParts = append(counterParts, otherFactor.CounterPart)

							continue
						}
					}
				}
			}
			if instances == instancesReq && !IsDuplicatedInCommonFactors(commonFactorCombos, factor.Factor) {

				combo := CommonFactorCombination{CommonFactor: factor.Factor, CounterParts: counterParts}

				commonFactorCombos = append(commonFactorCombos, combo)

				// commonFactors = append(commonFactors, factor)

				// counterPartFactors = append(counterPartFactors, counterParts)
			}
		}
	}
	return commonFactorCombos
}

func GetTermFactors(index int, expression *Expression) []TermFactorCombination {

	copy := expression.CopyTree()

	// expand exponents

	EvaluateAndReplace(index, &copy, EvaluateExponentExpansion)

	// get factors of constant add to list

	largestConstantFactor, isolatedFactors := GetIsolatedFactors(index, expression)

	// get all sublists of isolatedFactors

	factorGroups := GenerateFactorGroups(isolatedFactors, make([]Expression, 0), make([][]Expression, 0), 0)

	factors := make([]Expression, 0)

	termFactorCombos := make([]TermFactorCombination, 0)

	// iterate through all sublists, times each tuple together

	for _, group := range factorGroups {

		var factorToAdd Expression

		if len(group) > 1 {

			factorToAdd = MultiplyMany(group)

		} else {

			factorToAdd = group[0]
		}
		// make sure constants are not times together to be bigger than the initial value

		if !IsDuplicated(factors, factorToAdd) && !ExceedsLargestConstantFactor(largestConstantFactor, factorToAdd) {

			factors = append(factors, factorToAdd)

			// add to final factors if it equals the target
		}
	}

	// list counterparts for each factor

	for i := 0; i < len(factors); i++ {

		for j := 0; j < len(factors); j++ {

			if i == j {

				continue
			}
			mul := MultiplyTwo(factors[i], factors[j])

			mulRoot := mul.GetRoot()

			if comparison.IsEqualAt(index, mulRoot, expression, &mul) {

				// counterParts[i] = j

				termFactorCombos = append(termFactorCombos, TermFactorCombination{Factor: factors[i], CounterPart: factors[j]})
			}
		}
	}
	return termFactorCombos
}

func IsDuplicated(factors []Expression, factorToAdd Expression) bool {

	for _, factor := range factors {

		if comparison.IsEqual(factor, factorToAdd) {

			return true
		}
	}
	return false
}

func IsDuplicatedInCommonFactors(commonFactors []CommonFactorCombination, factorToAdd Expression) bool {

	for _, commonFactor := range commonFactors {

		if comparison.IsEqual(commonFactor.CommonFactor, factorToAdd) {

			return true
		}
	}
	return false
}

func IsDuplicatedInTermFactors(termFactors []TermFactorCombination, factorToAdd Expression) bool {

	for _, termFactor := range termFactors {

		if comparison.IsEqual(termFactor.Factor, factorToAdd) {

			return true
		}
	}
	return false
}

func ExceedsLargestConstantFactor(largestFactor int, compared Expression) bool {

	root := compared.GetRoot()

	if compared.IsConstant(root) {

		if compared.GetNode(root).NumericValue > largestFactor {

			return true

		} else {

			return false
		}

	} else if compared.IsMultiplication(root) {

		for _, child := range compared.GetChildren(root) {

			if compared.GetNode(child).NumericValue > largestFactor {

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

	_, expression := NewExpression(Symbol{Constant, value, strconv.Itoa(value)})

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

		value = expression.GetNode(index).NumericValue

		constantFactors := GetConstantFactors(index)

		factors = append(factors, IntegerFactorsToConstantExpression(constantFactors)...)

	} else if expression.IsMultiplication(index) {

		for _, child := range expression.GetChildren(index) {

			innerValue := expression.GetNode(child).NumericValue

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
