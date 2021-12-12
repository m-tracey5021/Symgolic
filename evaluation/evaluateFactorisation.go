package evaluation

import (
	"symgolic/comparison"
	"symgolic/conversion"
	. "symgolic/symbols"
)

type CommonFactor struct {
	Factor Expression

	CounterParts []Expression
}

type TermFactor struct {
	Factor Expression

	CounterPart Expression
}

func EvaluateFactorisation(index int, expression *Expression) (bool, Expression) {

	commonFactorsGroups := GetCommonFactors(index, expression)

	factoredExpressions := make([]Expression, 0)

	for i := 0; i < len(commonFactorsGroups); i++ {

		factoredRoot, factored := NewExpression(NewOperation(Multiplication))

		factored.AppendExpression(factoredRoot, commonFactorsGroups[i].Factor, false)

		add := factored.AppendNode(factoredRoot, NewOperation(Addition))

		for _, counterPart := range commonFactorsGroups[i].CounterParts {

			factored.AppendExpression(add, counterPart, false)
		}
		factoredExpressions = append(factoredExpressions, factored)
	}

	// add the rest of the non factors

	if len(factoredExpressions) == 0 {

		return false, *expression

	} else {

		resultRoot, result := NewExpression(NewOperation(Vector))

		for _, factored := range factoredExpressions {

			result.AppendExpression(resultRoot, factored, false)
		}

		return true, result
	}
}

func GetCommonFactors(index int, expression *Expression) []CommonFactor {

	commonFactors := make([]CommonFactor, 0)

	if expression.IsSummation(index) {

		termFactorGroups := make([][]TermFactor, 0)

		for _, term := range expression.GetChildren(index) {

			termFactors := GetTermFactors(term, expression)

			termFactorGroups = append(termFactorGroups, termFactors)
		}
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
				if instances == instancesReq && !IsDuplicatedInCommonFactors(commonFactors, factor.Factor) {

					combo := CommonFactor{Factor: factor.Factor, CounterParts: counterParts}

					commonFactors = append(commonFactors, combo)
				}
			}
		}
	}
	return commonFactors
}

func GetTermFactors(index int, expression *Expression) []TermFactor {

	copy := expression.CopySubtree(index)

	EvaluateAndReplace(copy.GetRoot(), &copy, EvaluateExponentExpansion)

	isolatedFactors := GetIsolatedFactors(copy.GetRoot(), &copy)

	factorGroups := GenerateFactorGroups(isolatedFactors, make([]Expression, 0), make([][]Expression, 0), 0)

	return SelectCompatibleFactors(copy, factorGroups)
}

func SelectCompatibleFactors(target Expression, factorGroups [][]Expression) []TermFactor {

	termFactors := make([]TermFactor, 0)

	for i := 0; i < len(factorGroups); i++ {

		currentFactor := Multiply(factorGroups[i]...)

		if IsDuplicatedInTermFactors(termFactors, currentFactor) {

			continue
		}
		for j := 0; j < len(factorGroups); j++ {

			if i == j {

				continue
			}
			comparedFactor := Multiply(factorGroups[j]...)

			result := Multiply(currentFactor, comparedFactor)

			if comparison.IsEqual(target, result) {

				termFactors = append(termFactors, TermFactor{Factor: currentFactor, CounterPart: comparedFactor})

				break
			}
		}
	}
	return termFactors
}

func GetIsolatedFactors(index int, expression *Expression) []Expression {

	factors := make([]Expression, 0)

	if expression.IsConstant(index) {

		value := expression.GetNode(index).NumericValue

		constantFactors := FindFactors(value)

		factors = append(factors, conversion.ConvertBulkIntToExpression(constantFactors)...)

	} else if expression.IsMultiplication(index) {

		for _, child := range expression.GetChildren(index) {

			innerValue := expression.GetNode(child).NumericValue

			if innerValue != -1 {

				innerConstantFactors := FindFactors(innerValue)

				factors = append(factors, conversion.ConvertBulkIntToExpression(innerConstantFactors)...)

			} else {

				factors = append(factors, expression.CopySubtree(child))
			}
		}

	} else {

		factors = append(factors, expression.CopySubtree(index))
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

func IsDuplicatedInCommonFactors(commonFactors []CommonFactor, factorToAdd Expression) bool {

	for _, commonFactor := range commonFactors {

		if comparison.IsEqual(commonFactor.Factor, factorToAdd) {

			return true
		}
	}
	return false
}

func IsDuplicatedInTermFactors(termFactors []TermFactor, factorToAdd Expression) bool {

	for _, termFactor := range termFactors {

		if comparison.IsEqual(termFactor.Factor, factorToAdd) {

			return true
		}
	}
	return false
}
