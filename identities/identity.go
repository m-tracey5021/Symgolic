package identities

import (
	. "symgolic/symbols"
)

type Identity interface {
	Identify()

	Apply()

	Run()
}

type ConstantCheck struct {
	Values []int

	Target int

	Operation SymbolType
}

type IdentityRequisite struct {
	Form string

	ConstantChecks []ConstantCheck
}

func CheckConstantValue(indices []int, targetIndex int, operation SymbolType, expression *Expression) bool {

	values := make([]int, 0)

	target := expression.GetNumericValueByIndex(targetIndex)

	for _, index := range indices {

		value := expression.GetNumericValueByIndex(index)

		if value == -1 {

			return false

		} else {

			values = append(values, value)
		}
	}
	if operation == Addition {

		total := 0

		for _, value := range values {

			total += value
		}
		return total == target

	} else if operation == Multiplication {

		total := 1

		for _, value := range values {

			total *= value
		}
		return total == target

	} else {

		return false
	}
}
