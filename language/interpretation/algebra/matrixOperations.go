package algebra

import (
	"symgolic/language/components"
	"symgolic/language/interpretation"
)

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

func Rref(index int, expression components.Expression) components.Expression { // reduced row echelon form

	// =============================

	// two loops, one to iterate through rows, one inner loop to iterate through columns

	// if matrix consists of n rows

	// column iterations = length of columns - 1

	// enter loop, in reverse, take nth row first

	// enter inner loop

	// take that rows iterated entry as a target, find entry in the same column in another row that is not n, say i, and find x where i[0] * x = n[0]

	// times every entry in row i by x

	// add row i to row n

	// decrement column iterations by 1

	if expression.IsNaryTuple(index) {

		rows := expression.GetChildren(index)

		colIterations := len(expression.GetChildren(rows[0])) - 1

		for i := len(rows) - 1; i >= 0; i-- { // iterate through rows

			row := expression.GetChildren(rows[i])

			for col := 0; col < colIterations; col++ { // iterate through cols

				_, value := components.NewExpression(components.NewConstant(expression.GetNode(row[col]).NumericValue))

				for j := 0; j < len(rows); j++ { // find other row which matches

					if i == j {

						continue
					}
					_, compared := components.NewExpression(components.NewConstant(expression.GetNode(expression.GetChildren(rows[j])[col]).NumericValue))

					// _, scalar := components.NewExpression(components.NewConstant(value / compared))

					// value = x * compared, find x

					scalar := interpretation.Divide(value, compared)

					toScale := expression.CopySubtree(rows[j])

					toSum := expression.CopySubtree(rows[i])

					_, scaled := Scale(toScale.GetRoot(), &toScale, &scalar)

					_, summed := VectorAdd(toSum.GetRoot(), scaled.GetRoot(), &toSum, &scaled)

					expression.ReplaceNodeCascade(rows[j], scaled)

					expression.ReplaceNodeCascade(rows[i], summed)
				}
			}
			colIterations--
		}
		return expression
	}

	return components.NewEmptyExpression()
}
