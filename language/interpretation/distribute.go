package interpretation

import (
	. "symgolic/language/components"
)

func Distribute(target ExpressionIndex) (bool, Expression) {

	if target.Expression.IsMultiplication(target.Index) {

		result := NewEmptyExpression()

		add := Symbol{Addition, -1, "+"}

		root := result.SetRoot(add)

		multiplications := DistributeAcross(&target.Expression, target.Expression.GetChildren(target.Index), 0, make(map[int]int))

		for _, multiplication := range multiplications {

			result.AppendExpression(root, multiplication, false)
		}
		return true, result

	} else {

		return false, target.Expression
	}
}

func DistributeAcross(expression *Expression, symbols []int, currentIndex int, sumMap map[int]int) []Expression {

	multiplications := make([]Expression, 0)

	children := make([]int, 0)

	if expression.IsSummation(symbols[currentIndex]) {

		children = expression.GetChildren(symbols[currentIndex])

	} else {

		children = append(children, symbols[currentIndex])
	}
	for i := 0; i < len(children); i++ {

		sumMap[symbols[currentIndex]] = children[i]

		if currentIndex != len(symbols)-1 {

			multiplications = append(multiplications, DistributeAcross(expression, symbols, currentIndex+1, sumMap)...)

		} else {

			values := make([]ExpressionIndex, 0)

			for _, value := range sumMap {

				values = append(values, ExpressionIndex{Expression: *expression, Index: value})
			}
			multiplications = append(multiplications, Multiply(values...))
		}
	}
	return multiplications
}
