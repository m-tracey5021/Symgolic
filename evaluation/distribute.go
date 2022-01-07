package evaluation

import (
	. "symgolic/symbols"
)

func Distribute(index int, expression *Expression) (bool, Expression) {

	if expression.IsMultiplication(index) {

		result := NewEmptyExpression()

		add := Symbol{Addition, -1, "+"}

		root := result.SetRoot(add)

		multiplications := DistributeAcross(expression, expression.GetChildren(index), 0, make(map[int]int))

		for _, multiplication := range multiplications {

			result.AppendExpression(root, multiplication, false)
		}
		return true, result

	} else {

		return false, *expression
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

			values := make([]int, 0)

			for _, value := range sumMap {

				values = append(values, value)
			}
			multiplications = append(multiplications, expression.Multiply(values))
		}
	}
	return multiplications
}
