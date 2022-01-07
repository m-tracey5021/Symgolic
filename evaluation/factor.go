package evaluation

import (
	"symgolic/comparison"
	"symgolic/generic"
	"symgolic/parsing"
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

func Factor(index int, expression *Expression) (bool, Expression) {

	if expression.IsSummation(index) {

		commonFactors := GetCommonFactors(index, expression)

		commonFactors = append(commonFactors, GetFactorsByGroupings(index, expression)...)

		factoredExpressions := GetFactoredExpressions(commonFactors)

		if len(factoredExpressions) == 0 {

			return false, *expression

		} else {

			resultRoot, result := NewExpression(NewOperation(NaryTuple))

			for _, factored := range factoredExpressions {

				result.AppendExpression(resultRoot, factored, false)
			}

			return true, result
		}

	} else {

		return false, *expression
	}
}

func GetFactoredExpressions(commonFactors []CommonFactor) []Expression {

	factoredExpressions := make([]Expression, 0)

	for i := 0; i < len(commonFactors); i++ {

		factoredRoot, factored := NewExpression(NewOperation(Multiplication))

		factored.AppendExpression(factoredRoot, commonFactors[i].Factor, false)

		add := factored.AppendNode(factoredRoot, NewOperation(Addition))

		for _, counterPart := range commonFactors[i].CounterParts {

			factored.AppendExpression(add, counterPart, false)
		}
		factoredExpressions = append(factoredExpressions, factored)
	}
	return factoredExpressions
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

func GetFactorsByGroupings(index int, expression *Expression) []CommonFactor {

	children := expression.GetChildren(index)

	groupings := generic.GenerateSubArrayGroups(children)

	// factoredExpressions := make([]Expression, 0)

	commonFactors := make([]CommonFactor, 0)

	for _, grouping := range groupings {

		commonFactorsPerGrouping := make([][]CommonFactor, 0)

		for _, group := range grouping {

			subRoot, subExpression := NewExpression(NewOperation(Addition))

			subExpression.AppendBulkSubtreesFrom(subRoot, group, *expression)

			commonFactorsPerGrouping = append(commonFactorsPerGrouping, GetCommonFactors(subRoot, &subExpression))
		}
		rows := GenerateGroupFactorRows(commonFactorsPerGrouping, make([][]CommonFactor, 0), make([]int, len(commonFactorsPerGrouping)), 0)

		for _, row := range rows {

			factors := make([]Expression, 0)

			counterParts := make([]Expression, 0)

			for _, commonFactorPerGroup := range row {

				factors = append(factors, commonFactorPerGroup.Factor)

				root, counterPart := NewExpression(NewOperation(Addition))

				counterPart.AppendBulkExpressions(root, commonFactorPerGroup.CounterParts)

				counterParts = append(counterParts, counterPart)
			}
			if comparison.AreEqual(counterParts...) {

				commonFactors = append(commonFactors, CommonFactor{Factor: counterParts[0], CounterParts: factors})
			}
		}
		// factoredExpressions = append(factoredExpressions, GetFactoredExpressions(commonFactors)...)
	}
	return commonFactors
}

func GenerateGroupFactorRows(matrix [][]CommonFactor, combinations [][]CommonFactor, rowIndexes []int, currentColumn int) [][]CommonFactor {

	for rowNumber := range matrix[currentColumn] {

		rowIndexes[currentColumn] = rowNumber

		if currentColumn == len(matrix)-1 { // if its the last column

			comboPerLine := make([]CommonFactor, 0)

			for colNumber, rowNumber := range rowIndexes {

				comboPerLine = append(comboPerLine, matrix[colNumber][rowNumber])
			}
			combinations = append(combinations, comboPerLine)

		} else {

			combinations = GenerateGroupFactorRows(matrix, combinations, rowIndexes, currentColumn+1)
		}
	}
	return combinations
}

func GetTermFactors(index int, expression *Expression) []TermFactor {

	expanded := expression.CopySubtree(index)

	EvaluateAndReplace(expanded.GetRoot(), &expanded, ExpandExponents)

	isolatedFactors := GetIsolatedFactors(expanded.GetRoot(), &expanded)

	factorGroups := GenerateFactorGroups(isolatedFactors, make([]Expression, 0), make([][]Expression, 0), 0)

	return SelectCompatibleFactors(expanded, factorGroups)
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

		factors = append(factors, parsing.ConvertBulkIntToExpression(constantFactors)...)

	} else if expression.IsMultiplication(index) {

		oneIncluded := false

		for _, child := range expression.GetChildren(index) {

			innerValue := expression.GetNode(child).NumericValue

			if innerValue != -1 {

				innerConstantFactors := FindFactors(innerValue)

				factors = append(factors, parsing.ConvertBulkIntToExpression(innerConstantFactors)...)

				oneIncluded = true

			} else {

				factors = append(factors, expression.CopySubtree(child))
			}
		}
		if !oneIncluded {

			_, one := NewExpression(NewConstant(1))

			factors = append(factors, one)
		}

	} else {

		factors = append(factors, expression.CopySubtree(index))

		_, one := NewExpression(NewConstant(1))

		factors = append(factors, one)
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
