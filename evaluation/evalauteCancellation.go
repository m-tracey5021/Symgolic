package evaluation

import (
	. "symgolic/symbols"
)

func EvaluateCancellation(index int, expression *Expression) (bool, Expression) {

	if expression.IsDivision(index) {

		exponents := make([]Expression, 0)

		num := expression.GetChildAtBreadth(index, 0)

		denom := expression.GetChildAtBreadth(index, 1)

		cancelledNums, cancelledDenoms, cont := InitCancelled(expression, num, denom)

		if cont {

			for i := 0; i < len(cancelledNums); i++ {

				// removed := false

				for j := 0; j < len(cancelledDenoms); j++ {

					if IsEqualByBase(cancelledNums[i], cancelledDenoms[j], expression, expression) {

						change, subtracted := SubtractExponents(expression, cancelledNums[i], cancelledDenoms[j])

						if change {

							exponents = append(exponents, subtracted)
						}

						cancelledNums = append(cancelledNums[:i], cancelledNums[i+1:]...)

						cancelledDenoms = append(cancelledDenoms[:j], cancelledDenoms[j+1:]...)

						// removed = true

						// i--

						// j = -1

						i = i - 1

						break
					}
				}
				// if removed {

				// 	i = i - 1
				// }
			}
			finalNums := DuplicateCancelled(expression, cancelledNums)

			finalNums = append(finalNums, exponents...)

			finalDenoms := DuplicateCancelled(expression, cancelledDenoms)

			result := CreateExpressionFromTerms(finalNums, finalDenoms)

			return true, result

		} else {

			return false, *expression
		}

	} else {

		return false, *expression
	}
}

func InitCancelled(expression *Expression, num, denom int) ([]int, []int, bool) {

	cancelledNums := make([]int, 0)

	cancelledDenoms := make([]int, 0)

	if expression.IsMultiplication(num) && expression.IsMultiplication(denom) {

		cancelledNums = append(cancelledNums, expression.GetChildren(num)...)

		cancelledDenoms = append(cancelledDenoms, expression.GetChildren(denom)...)

	} else if !expression.IsMultiplication(num) && expression.IsMultiplication(denom) {

		cancelledNums = append(cancelledNums, num)

		cancelledDenoms = append(cancelledDenoms, expression.GetChildren(denom)...)

	} else if expression.IsMultiplication(num) && !expression.IsMultiplication(denom) {

		cancelledNums = append(cancelledNums, expression.GetChildren(num)...)

		cancelledDenoms = append(cancelledDenoms, denom)

	} else {

		return cancelledNums, cancelledDenoms, false
	}
	return cancelledNums, cancelledDenoms, true
}

func SubtractExponents(expression *Expression, i, j int) (bool, Expression) {

	if expression.IsExponent(i) && expression.IsExponent(j) {

		lhs := expression.CopySubtree(expression.GetChildAtBreadth(i, 1))

		rhs := expression.CopySubtree(expression.GetChildAtBreadth(j, 1))

		sub := lhs.Subtract(rhs)

		exp := NewExpression()

		root := exp.SetRoot(Symbol{Exponent, -1, "^"})

		// exp.AppendExpression(root, expression.CopySubtree(expression.GetChildAtBreadth(i, 0)), false)

		exp.AppendSubtreeFrom(root, expression.GetChildAtBreadth(i, 0), *expression)

		exp.AppendExpression(root, sub, false)

		return true, exp

	} else {

		return false, *expression
	}
}

func DuplicateCancelled(expression *Expression, cancelled []int) []Expression {

	duplicated := make([]Expression, 0)

	for _, cancellation := range cancelled {

		copy := expression.CopySubtree(cancellation)

		duplicated = append(duplicated, copy)
	}
	return duplicated
}

func CreateExpressionFromTerms(nums, denoms []Expression) Expression {

	result := NewExpression()

	if len(nums) == 0 && len(denoms) == 0 {

		result.SetRoot(Symbol{Constant, 1, "1"})

	} else if len(nums) == 1 && len(denoms) == 0 {

		result.SetExpressionAsRoot(nums[0])

	} else if len(nums) == 0 && len(denoms) == 1 {

		result.SetExpressionAsRoot(denoms[0])

	} else if len(nums) > 1 && len(denoms) == 0 {

		root := result.SetRoot(Symbol{Multiplication, -1, "*"})

		result.AppendBulkExpressions(root, nums)

	} else if len(nums) == 0 && len(denoms) > 1 {

		root := result.SetRoot(Symbol{Multiplication, -1, "*"})

		result.AppendBulkExpressions(root, denoms)

	} else {

		div := Symbol{Division, -1, "/"}

		root := result.SetRoot(div)

		numMul := Symbol{Multiplication, -1, "*"}

		denomMul := Symbol{Multiplication, -1, "*"}

		if len(nums) == 1 && len(denoms) == 1 {

			result.AppendExpression(root, nums[0], false)

			result.AppendExpression(root, denoms[0], false)

		} else if len(nums) > 1 && len(denoms) == 1 {

			numIndex := result.AppendNode(root, numMul)

			result.AppendBulkExpressions(numIndex, nums)

			result.AppendExpression(root, denoms[0], false)

		} else if len(nums) == 1 && len(denoms) > 1 {

			result.AppendExpression(root, nums[0], false)

			denomIndex := result.AppendNode(root, denomMul)

			result.AppendBulkExpressions(denomIndex, denoms)

		} else {

			numIndex := result.AppendNode(root, numMul)

			denomIndex := result.AppendNode(root, denomMul)

			result.AppendBulkExpressions(numIndex, nums)

			result.AppendBulkExpressions(denomIndex, denoms)
		}
	}
	return result
}
