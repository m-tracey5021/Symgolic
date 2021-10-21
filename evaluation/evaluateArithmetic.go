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

		value := expression.GetNumericValueByIndex(child)

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

				continue
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

			EvaluateAndReplace(root, &result, RemoveMultiplicationByOne)
		}
		return change, result
	}
}

func RemoveMultiplicationByOne(index int, expression *Expression) (bool, Expression) {

	if expression.IsMultiplication(index) {

		removed := false

		children := expression.GetChildren(index)

		for i := 0; i < len(children); i++ {

			if expression.GetNumericValueByIndex(children[i]) == 1 {

				children = append(children[0:i], children[:i+1]...)

				removed = true
			}
		}
		if removed {

			if len(children) == 1 {

				return true, expression.CopySubtree(children[0])

			} else if len(children) > 1 {

				mulRoot, mul := NewExpressionWithRoot(Symbol{Multiplication, -1, "*"})

				mul.AppendBulkSubtreesFrom(mulRoot, children, *expression)

				return true, mul

			} else {

				panic(errors.New("Children has no length"))
			}

		} else {

			return false, *expression
		}

	} else {

		return false, *expression
	}
}

func MultiplyMany(operands []Expression) Expression {

	mulRoot, mul := NewExpressionWithRoot(Symbol{Multiplication, -1, "*"})

	for _, operand := range operands {

		mul.AppendExpression(mulRoot, operand, false)
	}
	EvaluateAndReplace(mulRoot, &mul, EvaluateConstants)

	return mul
}

func MultiplyTwo(operandA, operandB Expression) Expression {

	mulRoot, mul := NewExpressionWithRoot(Symbol{Multiplication, -1, "*"})

	mul.AppendExpression(mulRoot, operandA, false)

	mul.AppendExpression(mulRoot, operandB, false)

	EvaluateAndReplace(mulRoot, &mul, EvaluateConstants)

	return mul
}
