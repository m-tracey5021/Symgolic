package evaluation

import (
	"errors"
	"math"
	"strconv"
	. "symgolic/symbols"
)

func EvaluateConstants(index int, expression *Expression) (bool, Expression) {

	total := 0

	change := false

	duplicated := make([]Expression, 0)

	for _, child := range expression.GetChildren(index) {

		value := expression.GetNumericValuebyIndex(child)

		if value != -1 {

			if expression.IsSummation(index) {

				if !change {

					total = value

					change = true

				} else {

					total += value
				}

			} else if expression.IsMultiplication(index) {

				if !change {

					total = value

					change = true

				} else {

					total *= value
				}

			} else if expression.IsDivision(index) {

				if !change {

					total = value

					change = true

				} else {

					if total%value == 0 {

						total *= value
					}
				}

			} else if expression.IsExponent((index)) {

				if !change {

					total = value

					change = true

				} else {

					total = int(math.Pow(float64(total), float64(value)))
				}
			} else {

				panic(errors.New("not an arithmetic operation"))
			}

		} else {

			duplicated = append(duplicated, expression.CopySubtree(child))
		}
	}
	if !change {

		return change, *expression

	} else {

		result := NewExpression()

		if len(duplicated) == 0 {

			result.SetRoot(Symbol{Constant, total, strconv.Itoa(total)})

		} else {

			newParent := expression.GetNodeByIndex(index).Copy()

			root := result.SetRoot(newParent)

			result.AppendNode(root, Symbol{Constant, total, strconv.Itoa(total)})

			result.AppendBulkExpressions(root, duplicated)
		}
		return change, result
	}
}
