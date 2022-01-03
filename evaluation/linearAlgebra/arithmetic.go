package linearAlgebra

import (
	"symgolic/evaluation"
	"symgolic/symbols"
)

func Scale(index, scalar int, expression *symbols.Expression) (bool, symbols.Expression) {

	if expression.IsVector(index) {

		_, scalarExpression := symbols.NewExpression(symbols.NewConstant(scalar))

		for _, child := range expression.GetChildren(index) {

			replacement := evaluation.Multiply(scalarExpression, expression.CopySubtree(child))

			expression.ReplaceNodeCascade(child, replacement)
		}
		return true, *expression

	} else {

		return false, symbols.NewEmptyExpression()
	}
}

func DotProduct(index, indexInOther int, expression, other *symbols.Expression) (bool, symbols.Expression) {

	if expression.IsVector(index) && other.IsVector(index) {

		children := expression.GetChildren(index)

		otherChildren := other.GetChildren(indexInOther)

		if len(children) != len(otherChildren) {

			return false, symbols.NewEmptyExpression()

		} else {

			root, result := symbols.NewExpression(symbols.NewOperation(symbols.Addition))

			for i := 0; i < len(children); i++ {

				nthTotal := evaluation.Multiply(expression.CopySubtree(children[i]), other.CopySubtree(otherChildren[i]))

				result.AppendExpression(root, nthTotal, false)
			}
			evaluation.EvaluateAndReplace(root, &result, evaluation.EvaluateConstants)

			return true, result
		}

	} else {

		return false, symbols.NewEmptyExpression()
	}
}

func CrossProduct(indexA, indexB int, expressionA, expressionB *symbols.Expression) (bool, symbols.Expression) {

	if expressionA.IsVector(indexA) && expressionB.IsVector(indexB) {

		childrenA := expressionA.GetChildren(indexA)

		childrenB := expressionB.GetChildren(indexB)

		if len(childrenA) != 3 && len(childrenB) != 3 { // cross product only works in 3rd and 7th dimension

			return false, symbols.NewEmptyExpression()

		} else {

			C1 := evaluation.Subtract(evaluation.Multiply(expressionA.CopySubtree(childrenA[1]), expressionB.CopySubtree(childrenB[2])), evaluation.Multiply(expressionA.CopySubtree(childrenA[2]), expressionB.CopySubtree(childrenB[1])))

			C2 := evaluation.Subtract(evaluation.Multiply(expressionA.CopySubtree(childrenA[2]), expressionB.CopySubtree(childrenB[0])), evaluation.Multiply(expressionA.CopySubtree(childrenA[0]), expressionB.CopySubtree(childrenB[2])))

			C3 := evaluation.Subtract(evaluation.Multiply(expressionA.CopySubtree(childrenA[0]), expressionB.CopySubtree(childrenB[1])), evaluation.Multiply(expressionA.CopySubtree(childrenA[1]), expressionB.CopySubtree(childrenB[0])))

			root, result := symbols.NewExpression(symbols.NewOperation(symbols.Vector))

			result.AppendExpression(root, C1, false)

			result.AppendExpression(root, C2, false)

			result.AppendExpression(root, C3, false)

			return true, result
		}

	} else {

		return false, symbols.NewEmptyExpression()
	}
}
