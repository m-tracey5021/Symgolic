package interpretation

import (
	"symgolic/generic"
	. "symgolic/language/components"
	. "symgolic/language/interpretation/conversion"
)

type CommonFactor struct {
	Factor Expression

	CounterParts []Expression
}

type TermFactor struct {
	Factor Expression

	CounterPart Expression
}

func Factor(target ExpressionIndex) (bool, Expression) {

	if target.Expression.IsSummation(target.Index) {

		commonFactors := GetCommonFactors(target)

		commonFactors = append(commonFactors, GetFactorsByGroupings(target)...)

		factoredExpressions := GetFactoredExpressions(commonFactors)

		if len(factoredExpressions) == 0 {

			return false, target.Expression

		} else {

			resultRoot, result := NewExpression(NewOperation(NaryTuple))

			for _, factored := range factoredExpressions {

				result.AppendExpression(resultRoot, factored, false)
			}

			return true, result
		}

	} else {

		return false, target.Expression
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

func GetCommonFactors(target ExpressionIndex) []CommonFactor {

	commonFactors := make([]CommonFactor, 0)

	if target.Expression.IsSummation(target.Index) {

		termFactorGroups := make([][]TermFactor, 0)

		for _, term := range target.Expression.GetChildren(target.Index) {

			termFactors := GetTermFactors(target.At(term))

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

							if IsEqual(factor.Factor, otherFactor.Factor) {

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

func GetFactorsByGroupings(target ExpressionIndex) []CommonFactor {

	children := target.Expression.GetChildren(target.Index)

	groupings := generic.GenerateSubArrayGroups(children)

	// factoredExpressions := make([]Expression, 0)

	commonFactors := make([]CommonFactor, 0)

	for _, grouping := range groupings {

		commonFactorsPerGrouping := make([][]CommonFactor, 0)

		for _, group := range grouping {

			subRoot, subExpression := NewExpression(NewOperation(Addition))

			subExpression.AppendBulkSubtreesFrom(subRoot, group, target.Expression)

			commonFactorsPerGrouping = append(commonFactorsPerGrouping, GetCommonFactors(From(subExpression)))
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
			if AreEqual(counterParts...) {

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

func GetTermFactors(target ExpressionIndex) []TermFactor {

	expanded := target.Expression.CopySubtree(target.Index)

	EvaluateAndReplace(From(expanded), ExpandExponents)

	isolatedFactors := GetIsolatedFactors(From(expanded))

	factorGroups := GenerateFactorGroups(isolatedFactors, make([]Expression, 0), make([][]Expression, 0), 0)

	return SelectCompatibleFactors(expanded, factorGroups)
}

func SelectCompatibleFactors(target Expression, factorGroups [][]Expression) []TermFactor {

	termFactors := make([]TermFactor, 0)

	for i := 0; i < len(factorGroups); i++ {

		currentFactor := Multiply(FromMany(factorGroups[i])...)

		if IsDuplicatedInTermFactors(termFactors, currentFactor) {

			continue
		}
		for j := 0; j < len(factorGroups); j++ {

			if i == j {

				continue
			}
			comparedFactor := Multiply(FromMany(factorGroups[j])...)

			result := Multiply(From(currentFactor), From(comparedFactor))

			if IsEqual(target, result) {

				termFactors = append(termFactors, TermFactor{Factor: currentFactor, CounterPart: comparedFactor})

				break
			}
		}
	}
	return termFactors
}

func GetIsolatedFactors(target ExpressionIndex) []Expression {

	factors := make([]Expression, 0)

	if target.Expression.IsConstant(target.Index) {

		value := target.Expression.GetNode(target.Index).NumericValue

		constantFactors := FindFactors(value)

		factors = append(factors, ConvertBulkIntToExpression(constantFactors)...)

	} else if target.Expression.IsMultiplication(target.Index) {

		oneIncluded := false

		for _, child := range target.Expression.GetChildren(target.Index) {

			innerValue := target.Expression.GetNode(child).NumericValue

			if innerValue != -1 {

				innerConstantFactors := FindFactors(innerValue)

				factors = append(factors, ConvertBulkIntToExpression(innerConstantFactors)...)

				oneIncluded = true

			} else {

				factors = append(factors, target.Expression.CopySubtree(child))
			}
		}
		if !oneIncluded {

			_, one := NewExpression(NewConstant(1))

			factors = append(factors, one)
		}

	} else {

		factors = append(factors, target.Expression.CopySubtree(target.Index))

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

		if IsEqual(commonFactor.Factor, factorToAdd) {

			return true
		}
	}
	return false
}

func IsDuplicatedInTermFactors(termFactors []TermFactor, factorToAdd Expression) bool {

	for _, termFactor := range termFactors {

		if IsEqual(termFactor.Factor, factorToAdd) {

			return true
		}
	}
	return false
}
