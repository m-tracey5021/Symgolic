package interpretation

import (
	"symgolic/language/components"
)

func Scale(index int, expression, scalar *components.Expression) (bool, components.Expression) {

	if expression.IsVector(index) {

		for _, child := range expression.GetChildren(index) {

			replacement := Multiply(*scalar, expression.CopySubtree(child))

			expression.ReplaceNodeCascade(child, replacement)
		}
		return true, *expression

	} else {

		return false, components.NewEmptyExpression()
	}
}

func DotProduct(indexA, indexB int, expressionA, expressionB *components.Expression) (bool, components.Expression) {

	if expressionA.IsVector(indexA) && expressionB.IsVector(indexB) {

		children := expressionA.GetChildren(indexA)

		otherChildren := expressionB.GetChildren(indexB)

		if len(children) != len(otherChildren) {

			return false, components.NewEmptyExpression()

		} else {

			root, result := components.NewExpression(components.NewOperation(components.Addition))

			for i := 0; i < len(children); i++ {

				nthTotal := Multiply(expressionA.CopySubtree(children[i]), expressionB.CopySubtree(otherChildren[i]))

				result.AppendExpression(root, nthTotal, false)
			}
			EvaluateAndReplace(root, &result, ApplyArithmetic)

			return true, result
		}

	} else {

		return false, components.NewEmptyExpression()
	}
}

func CrossProduct(indexA, indexB int, expressionA, expressionB *components.Expression) (bool, components.Expression) {

	if expressionA.IsVector(indexA) && expressionB.IsVector(indexB) {

		childrenA := expressionA.GetChildren(indexA)

		childrenB := expressionB.GetChildren(indexB)

		if len(childrenA) != 3 && len(childrenB) != 3 { // cross product only works in 3rd and 7th dimension

			return false, components.NewEmptyExpression()

		} else {

			C1 := Subtract(Multiply(expressionA.CopySubtree(childrenA[1]), expressionB.CopySubtree(childrenB[2])), Multiply(expressionA.CopySubtree(childrenA[2]), expressionB.CopySubtree(childrenB[1])))

			C2 := Subtract(Multiply(expressionA.CopySubtree(childrenA[2]), expressionB.CopySubtree(childrenB[0])), Multiply(expressionA.CopySubtree(childrenA[0]), expressionB.CopySubtree(childrenB[2])))

			C3 := Subtract(Multiply(expressionA.CopySubtree(childrenA[0]), expressionB.CopySubtree(childrenB[1])), Multiply(expressionA.CopySubtree(childrenA[1]), expressionB.CopySubtree(childrenB[0])))

			root, result := components.NewExpression(components.NewOperation(components.Vector))

			result.AppendExpression(root, C1, false)

			result.AppendExpression(root, C2, false)

			result.AppendExpression(root, C3, false)

			return true, result
		}

	} else {

		return false, components.NewEmptyExpression()
	}
}
