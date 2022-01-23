package interpretation

import (
	. "symgolic/language/components"
)

func SumLikeTerms(target ExpressionIndex) (bool, Expression) {

	totalTerms := make([]Expression, 0)

	if target.Expression.IsSummation(target.Index) {

		children := target.Expression.GetChildren(target.Index)

		visited := make(map[int]bool)

		for i := 0; i < len(children); i++ {

			totalSum := 0

			if !visited[i] && target.Expression.IsMultiplication(children[i]) {

				visited[i] = true

				commonTerms := make([]Expression, 0)

				for j := i + 1; j < len(children); j++ {

					if !visited[j] {

						isLikeTerm, coeff, terms := IsLikeTerm(ExpressionIndex{Expression: target.Expression, Index: children[i]}, ExpressionIndex{Expression: target.Expression, Index: children[j]})

						if isLikeTerm {

							visited[j] = true

							totalSum += coeff

							if len(commonTerms) == 0 {

								commonTerms = terms
							}
						}
					}
				}
				root, mul := NewExpression(NewOperation(Multiplication))

				if totalSum > 1 {

					mul.AppendNode(root, NewConstant(totalSum))
				}
				mul.AppendBulkExpressions(root, commonTerms)

				totalTerms = append(totalTerms, mul)

			} else {

				copy := target.Expression.CopySubtree(i)

				totalTerms = append(totalTerms, copy)
			}
		}
		if len(totalTerms) > 1 {

			root, result := NewExpression(NewOperation(Addition))

			result.AppendBulkExpressions(root, totalTerms)

			return true, result

		} else {

			return true, totalTerms[0]
		}

	} else {

		return false, target.Expression
	}

}

func IsLikeTerm(a, b ExpressionIndex) (bool, int, []Expression) {

	if a.Expression.IsMultiplication(a.Index) && b.Expression.IsMultiplication(b.Index) {

		coeffA, termsA := GetTerms(a)

		coeffB, termsB := GetTerms(b)

		if len(termsA) != len(termsB) {

			return false, 0, nil

		} else {

			termsExp := make([]Expression, 0)

			for i := 0; i < len(termsA); i++ {

				if !IsEqualAt(a.At(termsA[i]), b.At(termsB[i])) {

					return false, 0, nil

				} else {

					termsExp = append(termsExp, a.Expression.CopySubtree(termsA[i]))
				}
			}
			return true, coeffA + coeffB, termsExp
		}

	} else {

		return false, 0, nil
	}
}

func GetTerms(target ExpressionIndex) (int, []int) {

	terms := make([]int, 0)

	if target.Expression.IsMultiplication(target.Index) {

		coefficient := 1

		for _, child := range target.Expression.GetChildren(target.Index) {

			if target.Expression.IsConstant(child) {

				if coefficient == 1 {

					coefficient = target.Expression.GetNode(child).NumericValue

				} else {

					coefficient *= target.Expression.GetNode(child).NumericValue
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
