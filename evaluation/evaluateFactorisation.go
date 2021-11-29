package evaluation

import (
	"symgolic/comparison"
	"symgolic/conversion"
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

		termFactorGroups := make([][]TermFactorCombination, 0)

		for _, term := range expression.GetChildren(index) {

			termFactors := GetTermFactors(term, expression)

			termFactorGroups = append(termFactorGroups, termFactors)
		}

		commonFactorsGroups := GetCommonFactors(expression, termFactorGroups)

		factoredExpressions := make([]Expression, 0)

		for i := 0; i < len(commonFactorsGroups); i++ {

			factoredRoot, factored := NewExpression(NewOperation(Multiplication))

			factored.AppendExpression(factoredRoot, commonFactorsGroups[i].CommonFactor, false)

			add := factored.AppendNode(factoredRoot, NewOperation(Addition))

			for _, counterPart := range commonFactorsGroups[i].CounterParts {

				factored.AppendExpression(add, counterPart, false)
			}
			factoredExpressions = append(factoredExpressions, factored)
		}

		// add the rest of the non factors

		resultRoot, result := NewExpression(NewOperation(Vector))

		for _, factored := range factoredExpressions {

			result.AppendExpression(resultRoot, factored, false)
		}

		return true, result

	} else {

		return false, *expression
	}
}

func GetCommonFactors(expression *Expression, termFactorGroups [][]TermFactorCombination) []CommonFactorCombination {

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
			}
		}
	}
	return commonFactorCombos
}

func GetTermFactors(index int, expression *Expression) []TermFactorCombination {

	copy := expression.CopySubtree(index)

	// expand exponents

	EvaluateAndReplace(copy.GetRoot(), &copy, EvaluateExponentExpansion)

	// get factors of constant add to list

	largestConstantFactor, isolatedFactors := GetIsolatedFactors(copy.GetRoot(), &copy)

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

			if comparison.IsEqualAt(copy.GetRoot(), mulRoot, &copy, &mul) {

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

func GetIsolatedFactors(index int, expression *Expression) (int, []Expression) {

	value := -1

	factors := make([]Expression, 0)

	if expression.IsConstant(index) {

		value = expression.GetNode(index).NumericValue

		constantFactors := FindFactors(value)

		factors = append(factors, conversion.ConvertBulkIntToExpression(constantFactors)...)

	} else if expression.IsMultiplication(index) {

		for _, child := range expression.GetChildren(index) {

			innerValue := expression.GetNode(child).NumericValue

			if innerValue != -1 {

				innerConstantFactors := FindFactors(innerValue)

				factors = append(factors, conversion.ConvertBulkIntToExpression(innerConstantFactors)...)

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
