package algebra

import (
	"math"
	"symgolic/language/components"
	"symgolic/language/interpretation"
)

func VectorAdd(indexA, indexB int, expressionA, expressionB *components.Expression) (bool, components.Expression) {

	childrenA := expressionA.GetChildren(indexA)

	childrenB := expressionB.GetChildren(indexB)

	if len(childrenA) == len(childrenB) {

		root, result := components.NewExpression(components.NewOperation(components.Vector))

		for i := 0; i < len(childrenA); i++ {

			value := expressionA.GetNode(childrenA[i]).NumericValue + expressionB.GetNode(childrenB[i]).NumericValue

			result.AppendNode(root, components.NewConstant(value))
		}
		return true, result

	} else {

		return false, components.NewEmptyExpression()
	}
}

func Scale(index int, expression, scalar *components.Expression) (bool, components.Expression) {

	if expression.IsVector(index) {

		for _, child := range expression.GetChildren(index) {

			replacement := interpretation.Multiply(*scalar, expression.CopySubtree(child))

			expression.ReplaceNodeCascade(child, replacement)
		}
		return true, *expression

	} else {

		return false, components.NewEmptyExpression()
	}
}

func Magnitude(index int, expression components.Expression) int {

	if expression.IsVector(index) {

		total := 0

		for _, child := range expression.GetChildren(index) {

			value := expression.GetNode(child).NumericValue

			if value != -1 {

				total += value * value

			} else {

				return -1
			}
		}
		return int(math.Sqrt(float64(total)))

	} else {

		return -1
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

				nthTotal := interpretation.Multiply(expressionA.CopySubtree(children[i]), expressionB.CopySubtree(otherChildren[i]))

				result.AppendExpression(root, nthTotal, false)
			}
			interpretation.EvaluateAndReplace(root, &result, interpretation.ApplyArithmetic)

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

			C1 := interpretation.Subtract(interpretation.Multiply(expressionA.CopySubtree(childrenA[1]), expressionB.CopySubtree(childrenB[2])), interpretation.Multiply(expressionA.CopySubtree(childrenA[2]), expressionB.CopySubtree(childrenB[1])))

			C2 := interpretation.Subtract(interpretation.Multiply(expressionA.CopySubtree(childrenA[2]), expressionB.CopySubtree(childrenB[0])), interpretation.Multiply(expressionA.CopySubtree(childrenA[0]), expressionB.CopySubtree(childrenB[2])))

			C3 := interpretation.Subtract(interpretation.Multiply(expressionA.CopySubtree(childrenA[0]), expressionB.CopySubtree(childrenB[1])), interpretation.Multiply(expressionA.CopySubtree(childrenA[1]), expressionB.CopySubtree(childrenB[0])))

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

func IsSquareMatrix(index int, expression components.Expression) bool {

	if expression.IsNaryTuple(index) {

		children := expression.GetChildren(index)

		cols := len(children)

		for _, child := range children {

			if !expression.IsVector(child) || len(expression.GetChildren(child)) != cols {

				return false
			}
		}
		return true

	} else {

		return false
	}
}

func FindDeterminant(index int, expression components.Expression) int {

	if IsSquareMatrix(index, expression) {

	} else {

		return -1
	}
}

func IsLinearCombination(indexA int, target components.Expression, others ...components.Expression) bool {

	// will have to factor here somehow

	return false
}

func IsLinearlyDependent(expressions map[int]components.Expression) bool {

}
