package evaluation

import (
	"math"
	"symgolic/generic"
	. "symgolic/symbols"
)

func ApplyArithmetic(index int, expression *Expression) (bool, Expression) {

	total := 0

	change := false

	duplicated := make([]Expression, 0)

	for _, child := range expression.GetChildren(index) {

		value := expression.GetNode(child).NumericValue

		if value != -1 {

			aux := expression.GetAuxiliaries(child)

			if len(aux) != 0 {

				if aux[0].SymbolType == Subtraction {

					value = value - (value * 2)
				}
			}
			if expression.IsSummation(index) {

				if !change {

					total = value

					change = true

				} else {

					total += value
				}

			} else if expression.IsMultiplication(index) {

				if value != 1 {

					if !change {

						total = value

						change = true

					} else {

						total *= value
					}
				}

			} else if expression.IsDivision(index) {

				if !change {

					total = value

					change = true

				} else {

					if total%value == 0 {

						total /= value
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

		result := NewEmptyExpression()

		negated := false

		if total < 0 {

			total = total + (total * -2)

			negated = true
		}

		if len(duplicated) == 0 {

			result.SetRoot(NewConstant(total))

			if negated {

				result = Negate(result)
			}

		} else {

			newParent := expression.GetNode(index).Copy()

			root := result.SetRoot(newParent)

			constant := result.AppendNode(root, NewConstant(total))

			if negated {

				result.InsertAuxiliariesAt(constant, []Symbol{NewOperation(Subtraction)})
			}
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

			if expression.GetNode(children[i]).NumericValue == 1 {

				children = append(children[0:i], children[:i+1]...)

				removed = true
			}
		}
		if removed {

			if len(children) == 1 {

				return true, expression.CopySubtree(children[0])

			} else if len(children) > 1 {

				mulRoot, mul := NewExpression(NewOperation(Multiplication))

				mul.AppendBulkSubtreesFrom(mulRoot, children, *expression)

				return true, mul

			} else {

				panic("Children has no length")
			}

		} else {

			return false, *expression
		}

	} else {

		return false, *expression
	}
}

func Negate(operand Expression) Expression {

	copy := operand.CopyTree()

	root := copy.GetRoot()

	copy.InsertAuxiliariesAt(root, []Symbol{NewOperation(Subtraction)})

	return copy
}

func Subtract(operands ...Expression) Expression {

	if len(operands) == 1 {

		return operands[0]
	}
	sumRoot, sum := NewExpression(NewOperation(Addition))

	sum.AppendExpression(sumRoot, operands[0], false)

	for i := 1; i < len(operands); i++ {

		sum.AppendExpression(sumRoot, Negate(operands[i]), false)
	}
	EvaluateAndReplace(sumRoot, &sum, ApplyArithmetic)

	return sum
}

func Multiply(operands ...Expression) Expression {

	if len(operands) == 1 {

		return operands[0]
	}
	mulRoot, mul := NewExpression(NewOperation(Multiplication))

	for _, operand := range operands {

		mul.AppendExpression(mulRoot, operand, false)
	}
	EvaluateAndReplace(mulRoot, &mul, ApplyArithmetic)

	return mul
}

func FindAdditives(value int) []int {

	additives := make([]int, 0)

	for i := 1; i <= value; i++ {

		if value-i != 0 {

			additives = append(additives, value-i)
		}
	}
	if value%2 == 0 {

		additives = append(additives, value/2)
	}
	return additives
}

func FindFactors(value int) []int {

	factors := make([]int, 0)

	for i := 1; i <= value; i++ {

		if value%i == 0 {

			factors = append(factors, i)
		}
	}
	return factors
}

func FindDividends(value, limit int) [][]int {

	dividends := make([][]int, 0)

	for i := 0; i <= limit; i++ {

		dividend := []int{value * i, i}

		dividends = append(dividends, dividend)
	}
	return dividends
}

func FindRoots(value int) [][]int {

	roots := make([][]int, 0)

	done := false

	for i := 1; i < value; i++ {

		total := i

		for j := 1; j < value; j++ {

			total *= i

			if total == value {

				if j == 2 {

					done = true
				}

				roots = append(roots, []int{i, j + 1})

				break

			} else if total > value {

				break
			}
		}
		if done {

			break
		}
	}
	return roots
}

func FindDegree(index int, expression *Expression) int {

	if expression.IsSummation(index) {

		largest := 1

		for _, child := range expression.GetChildren(index) {

			if expression.IsExponent(child) {

				value := expression.GetNode(expression.GetChildAtBreadth(child, 1)).NumericValue

				if value > largest {

					largest = value
				}
			}
		}
		return largest

	} else if expression.IsExponent(index) {

		value := expression.GetNode(expression.GetChildAtBreadth(index, 1)).NumericValue

		return value

	} else {

		return 1
	}
}

func FindAllOperands(value int, operation SymbolType) []int {

	operands := make([]int, 0)

	switch operation {

	case Addition:

		operands = FindAdditives(value)

	case Multiplication:

		operands = FindFactors(value)
	}
	if len(operands) == 1 && (operation == Addition || operation == Multiplication) {

		return make([]int, 0)

	} else {

		totalOperands := make([]int, 0)

		for _, operand := range operands {

			if operand != value {

				innerOperands := FindAllOperands(operand, operation)

				for _, inner := range innerOperands {

					if inner != 1 && inner != operand {

						totalOperands = append(totalOperands, inner)
					}
				}
			}
		}
		totalOperands = append(totalOperands, operands...)

		return totalOperands
	}
}

func GeneratePossibleOperandCombinationsForValue(value, limit int, operation SymbolType) [][]int {

	operandGroups := generic.GenerateSubArrays(FindAllOperands(value, operation), limit)

	operandGroupsNoDuplicates := make([][]int, 0)

	for _, operandGroup := range operandGroups {

		duplicate := false

		for _, operandGroupCompared := range operandGroupsNoDuplicates {

			if len(operandGroup) == len(operandGroupCompared) {

				count := 0

				for i := 0; i < len(operandGroup); i++ {

					if operandGroup[i] == operandGroupCompared[i] {

						count++

					} else {

						break
					}
				}
				duplicate = count == len(operandGroup)

				if duplicate {

					break
				}
			}
		}
		if !duplicate {

			operandGroupsNoDuplicates = append(operandGroupsNoDuplicates, operandGroup)
		}
	}
	return VerifySubArrays(operandGroupsNoDuplicates, value, operation)
}

func VerifySubArrays(subarrays [][]int, target int, operation SymbolType) [][]int {

	verified := make([][]int, 0)

	if operation == Addition {

		for _, subarray := range subarrays {

			total := 0

			for _, value := range subarray {

				total += value
			}
			if total == target {

				verified = append(verified, subarray)
			}
		}

	} else if operation == Multiplication {

		for _, subarray := range subarrays {

			total := 1

			for _, value := range subarray {

				total *= value
			}
			if total == target {

				verified = append(verified, subarray)
			}
		}
	}
	return verified
}
