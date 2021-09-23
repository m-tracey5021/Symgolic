package evaluation

import (
	"errors"
	. "symgolic/symbols"
)

func EvaluateExponentExpansion(index int, expression *Expression) (bool, Expression) {

	if expression.IsExponent(index) {

		result := *expression

		change := false

		power := expression.GetChildAtBreadth(index, 1)

		if !IsAtomicExponent(power, result) {

			target := expression.GetChildAtBreadth(index, 0)

			if expression.IsSummation(power) {

				change, result = ExpandSummation(target, power, expression)

			} else if expression.IsMultiplication(power) {

				change, result = ExpandMultiplication(target, power, expression)

			} else if expression.IsDivision(power) {

				change, result = ExpandDivision(target, power, expression)

			} else if expression.IsExponent(power) {

				change, result = EvaluateExponentExpansion(power, &result)

			} else {

				panic(errors.New("symbol doesnt make sense as exponent"))
			}
			for _, child := range result.GetChildren(result.GetRoot()) {

				change, result = EvaluateExponentExpansion(child, &result)
			}
		}
		return change, result

	} else {

		return false, *expression
	}
}

func IsAtomicExponent(index int, expression Expression) bool {

	if expression.IsVariable(index) {

		return true

	} else if expression.IsDivision(index) {

		num := expression.GetChildAtBreadth(index, 0)

		numVal := expression.GetNumericValueByIndex(num)

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

	result := NewExpression()

	mul := Symbol{Multiplication, -1, "*"}

	resultRoot := result.SetRoot(mul)

	for _, child := range expression.GetChildren(power) {

		operand := NewExpression()

		exp := Symbol{Exponent, -1, "^"}

		root := operand.SetRoot(exp)

		operand.AppendSubtreeFrom(root, target, *expression)

		operand.AppendSubtreeFrom(root, child, *expression)

		result.AppendExpression(resultRoot, operand, false)
	}
	return true, result
}

func ExpandMultiplication(target, power int, expression *Expression) (bool, Expression) {

	resultRoot, result := NewExpressionWithRoot(Symbol{Multiplication, -1, "*"})

	coefficient, terms := GetTerms(power, expression)

	if coefficient != 1 {

		duplicatedPower := NewExpression()

		if len(terms) > 1 {

			duplicatedPower.SetRoot(Symbol{Multiplication, -1, "*"})

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

			expRoot, exp := NewExpressionWithRoot(Symbol{Exponent, -1, "^"})

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

	numVal := expression.GetNumericValueByIndex(num)

	if numVal > 1 {

		resultRoot, result := NewExpressionWithRoot(Symbol{Multiplication, -1, "*"})

		for i := 0; i < numVal; i++ {

			root, duplicatedPower := NewExpressionWithRoot(Symbol{Division, -1, "/"})

			duplicatedPower.AppendNode(root, Symbol{Constant, 1, "1"})

			duplicatedPower.AppendSubtreeFrom(root, denom, *expression)

			expRoot, exp := NewExpressionWithRoot(Symbol{Exponent, -1, "^"})

			exp.AppendSubtreeFrom(expRoot, target, *expression)

			exp.AppendExpression(expRoot, duplicatedPower, false)

			result.AppendExpression(resultRoot, exp, false)
		}
		return true, result

	} else {

		return false, *expression
	}
}
