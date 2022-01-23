package interpretation

import (
	"math"
	"symgolic/generic"
	. "symgolic/language/components"
)

func ApplyArithmetic(target ExpressionIndex) (bool, Expression) {

	total := 0

	change := false

	duplicated := make([]Expression, 0)

	for _, child := range target.Expression.GetChildren(target.Index) {

		value := target.Expression.GetNode(child).NumericValue

		if value != -1 { // child is a number

			aux := target.Expression.GetAuxiliaries(child)

			if len(aux) != 0 {

				if aux[0].SymbolType == Subtraction {

					value = value - (value * 2)
				}
			}
			if target.Expression.IsSummation(target.Index) {

				if !change {

					total = value

					change = true

				} else {

					total += value
				}

			} else if target.Expression.IsMultiplication(target.Index) {

				if value != 1 {

					if !change {

						total = value

						change = true

					} else {

						total *= value
					}
				}

			} else if target.Expression.IsDivision(target.Index) {

				if !change {

					total = value

					change = true

				} else {

					if total%value == 0 {

						total /= value

					} else {

						change = false
					}
				}

			} else if target.Expression.IsExponent((target.Index)) {

				if !change {

					total = value

					change = true

				} else {

					total = int(math.Pow(float64(total), float64(value)))
				}
			} else {

				continue
			}

		} else { // child is a data structure

			duplicated = append(duplicated, target.Expression.CopySubtree(child))
		}
	}
	if !change {

		return change, target.Expression

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

				Negate(target.At(target.Expression.GetRoot()))
			}

		} else {

			newParent := target.Expression.GetNode(target.Index).Copy()

			root := result.SetRoot(newParent)

			constant := result.AppendNode(root, NewConstant(total))

			if negated {

				result.InsertAuxiliariesAt(constant, []Symbol{NewOperation(Subtraction)})
			}
			result.AppendBulkExpressions(root, duplicated)

			EvaluateAndReplace(From(result), RemoveMultiplicationByOne)
		}
		return change, result
	}
}

func RemoveMultiplicationByOne(target ExpressionIndex) (bool, Expression) {

	if target.Expression.IsMultiplication(target.Index) {

		removed := false

		children := target.Expression.GetChildren(target.Index)

		for i := 0; i < len(children); i++ {

			if target.Expression.GetNode(children[i]).NumericValue == 1 {

				children = append(children[0:i], children[:i+1]...)

				removed = true
			}
		}
		if removed {

			if len(children) == 1 {

				return true, target.Expression.CopySubtree(children[0])

			} else if len(children) > 1 {

				mulRoot, mul := NewExpression(NewOperation(Multiplication))

				mul.AppendBulkSubtreesFrom(mulRoot, children, target.Expression)

				return true, mul

			} else {

				panic("Children has no length")
			}

		} else {

			return false, target.Expression
		}

	} else {

		return false, target.Expression
	}
}

func EvaluateArithmetic(target ExpressionIndex) (bool, Expression) {

	symbolType := target.Expression.GetNode(target.Index).SymbolType

	operands := make([]ExpressionIndex, 0)

	for _, child := range target.Expression.GetChildren(target.Index) {

		operands = append(operands, target.At(child))
	}
	switch symbolType {

	case Addition:

		return true, Add(operands...)

	case Multiplication:

		return true, Multiply(operands...)

	case Division:

		return true, Divide(operands[0], operands[1])
	}
	return false, NewEmptyExpression()
}

type Pairing struct {
	First, Second SymbolType
}

type Operation struct {
	Call func(ExpressionIndex, ExpressionIndex, bool) Expression

	Reverse bool
}

// Comparer

func GetOperativeCall(a, b, operation SymbolType) Operation {

	var pairing map[Pairing]Operation

	switch operation {

	case Addition:

		pairing = map[Pairing]Operation{

			{First: Addition, Second: Addition}:             {Call: SSS, Reverse: false},
			{First: Addition, Second: Multiplication}:       {Call: SSM, Reverse: false},
			{First: Addition, Second: Division}:             {Call: SSD, Reverse: false},
			{First: Multiplication, Second: Addition}:       {Call: SSM, Reverse: true},
			{First: Multiplication, Second: Multiplication}: {Call: MSM, Reverse: false},
			{First: Multiplication, Second: Division}:       {Call: MSD, Reverse: false},
			{First: Division, Second: Addition}:             {Call: SSD, Reverse: true},
			{First: Division, Second: Multiplication}:       {Call: MSD, Reverse: true},
			{First: Division, Second: Division}:             {Call: DSD, Reverse: false},
		}

	case Multiplication:

		pairing = map[Pairing]Operation{}

	case Division:

		pairing = map[Pairing]Operation{}
	}
	match, exists := pairing[Pairing{First: a, Second: b}]

	if exists {

		return match
	}
	panic("No matching function")
}

// Arithmetic base

func Negate(target ExpressionIndex) {

	negation := make([]Symbol, 0)

	negation = append(negation, NewOperation(Subtraction))

	target.Expression.InsertAuxiliariesAt(target.Index, negation)
}

func Add(operands ...ExpressionIndex) Expression {

	cumulative := operands[0].Expression

	cumulativeIndex := operands[0].Index

	for i, operand := range operands {

		if i == 0 {

			continue
		}
		operation := GetOperativeCall(cumulative.GetNode(cumulativeIndex).SymbolType, operand.Expression.GetNode(operand.Index).SymbolType, Addition)

		cumulative = operation.Call(From(cumulative), operand, operation.Reverse)

		cumulativeIndex = cumulative.GetRoot()
	}
	return cumulative
}

func Subtract(a, b ExpressionIndex) Expression {

	operation := GetOperativeCall(a.Expression.GetNode(a.Index).SymbolType, b.Expression.GetNode(b.Index).SymbolType, Addition)

	Negate(b)

	return operation.Call(a, b, operation.Reverse)
}

func Multiply(operands ...ExpressionIndex) Expression {

	cumulative := operands[0].Expression

	cumulativeIndex := operands[0].Index

	for i, operand := range operands {

		if i == 0 {

			continue
		}

		operation := GetOperativeCall(cumulative.GetNode(cumulativeIndex).SymbolType, operand.Expression.GetNode(operand.Index).SymbolType, Multiplication)

		cumulative = operation.Call(From(cumulative), operand, operation.Reverse)

		cumulativeIndex = cumulative.GetRoot()
	}
	return cumulative
}

func Divide(a, b ExpressionIndex) Expression {

	operation := GetOperativeCall(a.Expression.GetNode(a.Index).SymbolType, b.Expression.GetNode(b.Index).SymbolType, Division)

	return operation.Call(a, b, operation.Reverse)
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
