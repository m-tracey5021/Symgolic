package conversion

import (
	"symgolic/symbols"
)

func ConvertIntToExpression(value int) symbols.Expression {

	_, expression := symbols.NewExpression(symbols.NewConstant(value))

	return expression
}

func ConvertBulkIntToExpression(values []int) []symbols.Expression {

	expressions := make([]symbols.Expression, 0)

	for _, value := range values {

		expressions = append(expressions, ConvertIntToExpression(value))
	}
	return expressions
}
