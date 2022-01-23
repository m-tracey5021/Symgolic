package interpretation

import (
	"errors"
	. "symgolic/language/components"
)

func ExpandExponents(target ExpressionIndex) (bool, Expression) {

	if target.Expression.IsExponent(target.Index) {

		result := target.Expression

		change := false

		power := target.Expression.GetChildAtBreadth(target.Index, 1)

		if !IsAtomicExponent(power, result) {

			index := target.Expression.GetChildAtBreadth(target.Index, 0)

			if target.Expression.IsSummation(power) {

				change, result = ExpandSummation(index, power, &target.Expression)

			} else if target.Expression.IsMultiplication(power) {

				change, result = ExpandMultiplication(index, power, &target.Expression)

			} else if target.Expression.IsDivision(power) {

				change, result = ExpandDivision(index, power, &target.Expression)

			} else if target.Expression.IsExponent(power) {

				change, result = ExpandExponents(From(result).At(power))

			} else if target.Expression.IsConstant(power) {

				change, result = ExpandConstant(index, power, &target.Expression)

			} else {

				panic(errors.New("symbol doesnt make sense as exponent"))
			}
			for _, child := range result.GetChildren(result.GetRoot()) {

				innerChange, innerResult := ExpandExponents(From(result).At(child))

				if innerChange {

					result.ReplaceNodeCascade(child, innerResult)
				}
			}
		}
		return change, result

	} else {

		return false, target.Expression
	}
}

func IsAtomicExponent(index int, expression Expression) bool {

	if expression.IsVariable(index) {

		return true

	} else if expression.IsMultiplication(index) {

		coefficient, _ := GetTerms(From(expression).At(index))

		if coefficient > 1 {

			return false

		} else {

			return true
		}

	} else if expression.IsDivision(index) {

		num := expression.GetChildAtBreadth(index, 0)

		numVal := expression.GetNode(num).NumericValue

		if numVal <= 1 {

			return true

		} else {

			return false
		}

	} else {

		return false
	}
}

func ExpandSummation(target, power int, expression *Expression) (bool, Expression) {

	resultRoot, result := NewExpression(NewOperation(Multiplication))

	for _, child := range expression.GetChildren(power) {

		root, operand := NewExpression(NewOperation(Exponent))

		operand.AppendSubtreeFrom(root, target, *expression)

		operand.AppendSubtreeFrom(root, child, *expression)

		result.AppendExpression(resultRoot, operand, false)
	}
	return true, result
}

func ExpandMultiplication(target, power int, expression *Expression) (bool, Expression) {

	resultRoot, result := NewExpression(NewOperation(Multiplication))

	coefficient, terms := GetTerms(From(*expression).At(power))

	if coefficient != 1 {

		duplicatedPower := NewEmptyExpression()

		if len(terms) > 1 {

			duplicatedPower.SetRoot(NewOperation(Multiplication))

			root := duplicatedPower.GetRoot()

			for _, term := range terms {

				duplicatedPower.AppendSubtreeFrom(root, term, *expression)
			}

		} else if len(terms) == 1 {

			duplicatedPower.SetExpressionAsRoot(expression.CopySubtree(terms[0]))

		} else {

			panic(errors.New("not a multiplication"))
		}
		for i := 0; i < coefficient; i++ {

			expRoot, exp := NewExpression(NewOperation(Exponent))

			exp.AppendSubtreeFrom(expRoot, target, *expression)

			exp.AppendExpression(expRoot, duplicatedPower, true)

			result.AppendExpression(resultRoot, exp, false)
		}
		return true, result

	} else {

		return false, *expression
	}
}

func ExpandDivision(target, power int, expression *Expression) (bool, Expression) {

	num := expression.GetChildAtBreadth(power, 0)

	denom := expression.GetChildAtBreadth(power, 1)

	numVal := expression.GetNode(num).NumericValue

	if numVal > 1 {

		resultRoot, result := NewExpression(NewOperation(Multiplication))

		for i := 0; i < numVal; i++ {

			root, duplicatedPower := NewExpression(NewOperation(Division))

			duplicatedPower.AppendNode(root, NewConstant(1))

			duplicatedPower.AppendSubtreeFrom(root, denom, *expression)

			expRoot, exp := NewExpression(NewOperation(Exponent))

			exp.AppendSubtreeFrom(expRoot, target, *expression)

			exp.AppendExpression(expRoot, duplicatedPower, false)

			result.AppendExpression(resultRoot, exp, false)
		}
		return true, result

	} else {

		return false, *expression
	}
}

func ExpandConstant(target, power int, expression *Expression) (bool, Expression) {

	value := expression.GetNode(power).NumericValue

	if value > 1 {

		resultRoot, result := NewExpression(NewOperation(Multiplication))

		for i := 0; i < value; i++ {

			result.AppendSubtreeFrom(resultRoot, target, *expression)
		}
		return true, result

	} else if value == 1 {

		self := expression.CopySubtree(target)

		return true, self

	} else {

		_, one := NewExpression(NewConstant(1))

		return true, one
	}
}
