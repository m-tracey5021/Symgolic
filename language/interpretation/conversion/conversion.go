package conversion

import (
	. "symgolic/language/components"
)

func ConvertIntToExpression(value int) Expression {

	_, expression := NewExpression(NewConstant(value))

	return expression
}

func ConvertBulkIntToExpression(values []int) []Expression {

	expressions := make([]Expression, 0)

	for _, value := range values {

		expressions = append(expressions, ConvertIntToExpression(value))
	}
	return expressions
}
