package algebra

import (
	"math"
	. "symgolic/language/components"
	"symgolic/language/interpretation"
)

func VectorAdd(indexA, indexB int, expressionA, expressionB *Expression) (bool, Expression) {

	childrenA := expressionA.GetChildren(indexA)

	childrenB := expressionB.GetChildren(indexB)

	if len(childrenA) == len(childrenB) {

		root, result := NewExpression(NewOperation(Vector))

		for i := 0; i < len(childrenA); i++ {

			value := expressionA.GetNode(childrenA[i]).NumericValue + expressionB.GetNode(childrenB[i]).NumericValue

			result.AppendNode(root, NewConstant(value))
		}
		return true, result

	} else {

		return false, NewEmptyExpression()
	}
}

func Scale(index int, expression, scalar *Expression) (bool, Expression) {

	if expression.IsVector(index) {

		for _, child := range expression.GetChildren(index) {

			replacement := interpretation.Multiply(From(*scalar), From(expression.CopySubtree(child)))

			expression.ReplaceNodeCascade(child, replacement)
		}
		return true, *expression

	} else {

		return false, NewEmptyExpression()
	}
}

func Magnitude(index int, expression Expression) int {

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

func DotProduct(a, b ExpressionIndex) (bool, Expression) {

	if a.Expression.IsVector(a.Index) && b.Expression.IsVector(b.Index) {

		children := a.Expression.GetChildren(a.Index)

		otherChildren := b.Expression.GetChildren(b.Index)

		if len(children) != len(otherChildren) {

			return false, NewEmptyExpression()

		} else {

			root, result := NewExpression(NewOperation(Addition))

			for i := 0; i < len(children); i++ {

				nthTotal := interpretation.Multiply(a.At(children[i]), b.At(otherChildren[i]))

				result.AppendExpression(root, nthTotal, false)
			}
			interpretation.EvaluateAndReplace(From(result), interpretation.ApplyArithmetic)

			return true, result
		}

	} else {

		return false, NewEmptyExpression()
	}
}

func CrossProduct(a, b ExpressionIndex) (bool, Expression) {

	if a.Expression.IsVector(a.Index) && b.Expression.IsVector(b.Index) {

		childrenA := a.Expression.GetChildren(a.Index)

		childrenB := b.Expression.GetChildren(b.Index)

		if len(childrenA) != 3 && len(childrenB) != 3 { // cross product only works in 3rd and 7th dimension

			return false, NewEmptyExpression()

		} else {

			x := interpretation.Multiply(a.At(childrenA[1]), b.At(childrenB[2]))

			y := interpretation.Multiply(a.At(childrenA[2]), b.At(childrenB[1]))

			C1 := interpretation.Subtract(From(x), From(y))

			x = interpretation.Multiply(a.At(childrenA[2]), b.At(childrenB[0]))

			y = interpretation.Multiply(a.At(childrenA[0]), b.At(childrenB[2]))

			C2 := interpretation.Subtract(From(x), From(y))

			x = interpretation.Multiply(a.At(childrenA[0]), b.At(childrenB[1]))

			y = interpretation.Multiply(a.At(childrenA[1]), b.At(childrenB[0]))

			C3 := interpretation.Subtract(From(x), From(y))

			root, result := NewExpression(NewOperation(Vector))

			result.AppendExpression(root, C1, false)

			result.AppendExpression(root, C2, false)

			result.AppendExpression(root, C3, false)

			return true, result
		}

	} else {

		return false, NewEmptyExpression()
	}
}
