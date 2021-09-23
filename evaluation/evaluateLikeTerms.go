package evaluation

import (
	"strconv"
	. "symgolic/symbols"
)

func EvaluateLikeTerms(index int, expression *Expression) (bool, Expression) {

	result := NewExpression()

	totalTerms := make([]Expression, 0)

	if expression.IsSummation(index) {

		children := expression.GetChildren(index)

		visited := make(map[int]bool)

		for i := 0; i < len(children); i++ {

			totalSum := 0

			firstVisited := visited[i]

			if !firstVisited && expression.IsMultiplication(children[i]) {

				compared := children[i]

				visited[i] = true

				coeff, terms := GetTerms(children[i], expression)

				totalSum += coeff

				for j := i + 1; j < len(children); j++ {

					secondVisited := visited[j]

					if !secondVisited {

						if IsLikeTerm(compared, children[j], expression) {

							visited[j] = true

							coeff, _ := GetTerms(children[j], expression)

							totalSum += coeff
						}
					}
				}
				summed := NewExpression()

				mul := Symbol{Multiplication, -1, "*"}

				root := summed.SetRoot(mul)

				if totalSum > 1 {

					summed.AppendNode(root, Symbol{Constant, totalSum, strconv.Itoa(totalSum)})
				}
				summed.AppendBulkSubtreesFrom(root, terms, *expression)

				totalTerms = append(totalTerms, summed)

			} else {

				copy := expression.CopySubtree(i)

				totalTerms = append(totalTerms, copy)
			}
		}
		if len(totalTerms) > 1 {

			add := Symbol{Addition, -1, "+"}

			root := result.SetRoot(add)

			result.AppendBulkExpressions(root, totalTerms)

		} else {

			result.SetExpressionAsRoot(totalTerms[0])
		}
		return true, result

	} else {

		return false, *expression
	}

}

func IsLikeTerm(first, second int, expression *Expression) bool {

	if expression.IsMultiplication(first) && expression.IsMultiplication(second) {

		_, firstTerms := GetTerms(first, expression)

		_, secondTerms := GetTerms(first, expression)

		if len(firstTerms) != len(secondTerms) {

			return false

		} else {

			for i := 0; i < len(firstTerms); i++ {

				if !IsEqual(firstTerms[i], secondTerms[i], expression, expression) {

					return false
				}
			}
			return true
		}

	} else {

		return false
	}
}

func GetTerms(index int, expression *Expression) (int, []int) {

	terms := make([]int, 0)

	if expression.IsMultiplication(index) {

		coefficient := 1

		for _, child := range expression.GetChildren(index) {

			if expression.IsConstant(child) {

				if coefficient == 1 {

					coefficient = expression.GetNumericValueByIndex(child)

				} else {

					coefficient *= expression.GetNumericValueByIndex(child)
				}

			} else {

				terms = append(terms, child)
			}
		}
		return coefficient, terms

	} else {

		return 1, terms
	}
}
