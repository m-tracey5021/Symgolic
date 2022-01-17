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

func IsTwoByTwoMatrix(index int, expression components.Expression) bool {

	if expression.IsNaryTuple(index) {

		children := expression.GetChildren(index)

		cols := len(children)

		if cols != 2 {

			return false
		}
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

func FindDeterminant(index int, expression components.Expression) components.Expression {

	if expression.IsNaryTuple(index) {

		children := expression.GetChildren(index)

		cols := len(children)

		if cols != 2 {

			// recurse into new matrix

			sumRoot, sum := components.NewExpression(components.NewOperation(components.Addition))

			sign := true

			for i, col := range children {

				rows := expression.GetChildren(col)

				if !expression.IsVector(col) || cols != len(rows) {

					return components.NewEmptyExpression()
				}
				target := rows[0]

				matrixRoot, subMatrix := components.NewExpression(components.NewOperation(components.NaryTuple))

				for j, comparedCol := range children {

					if i == j {

						continue
					}
					root, vector := components.NewExpression(components.NewOperation(components.Vector))

					for k, rowEntry := range expression.GetChildren(comparedCol) {

						if k != 0 {

							vector.AppendSubtreeFrom(root, rowEntry, expression)
						}
					}
					subMatrix.AppendExpression(matrixRoot, vector, false)
				}
				subDeterminant := FindDeterminant(matrixRoot, subMatrix)

				subMatrixResult := interpretation.Multiply(subDeterminant, expression.CopySubtree(target))

				if !sign {

					subMatrixResult = interpretation.Negate(subMatrixResult)
				}
				sign = !sign

				sum.AppendExpression(sumRoot, subMatrixResult, false)
			}
			interpretation.EvaluateAndReplace(sumRoot, &sum, interpretation.ApplyArithmetic)

			return sum

		} else {

			if !(len(expression.GetChildren(children[0])) == 2 && len(expression.GetChildren(children[1])) == 2) {

				return components.NewEmptyExpression()
			}
			mulA := interpretation.Multiply(expression.CopySubtree(expression.GetChildren(children[0])[0]), expression.CopySubtree(expression.GetChildren(children[1])[1]))

			mulB := interpretation.Multiply(expression.CopySubtree(expression.GetChildren(children[0])[1]), expression.CopySubtree(expression.GetChildren(children[1])[0]))

			mulB = interpretation.Negate(mulB)

			result := interpretation.Add(mulA, mulB)

			return result
		}

	} else {

		return components.NewEmptyExpression()
	}
}

func IsLinearCombination(indexA, indexB, targetIndex int, expressionA, expressionB, target components.Expression) bool {

	// set the vectors up as an augmented matrix, then row reduce, if there is a row where 0, 0, 0, 0, c where c != 0 then the system is inconsistent
	// and the target vector is not a linear combination of the others

	return false
}

func IsLinearlyDependentMatrix(index int, expression components.Expression) bool {

	determinant := FindDeterminant(index, expression)

	_, zero := components.NewExpression(components.NewConstant(0))

	return !interpretation.IsEqual(determinant, zero)
}
