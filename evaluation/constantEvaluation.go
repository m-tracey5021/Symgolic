package evaluation

import (
	. "symgolic/symbols"
)

func EvaluateConstants(expression *Expression) Expression {

	root := expression.GetRoot()

	SearchAndReplace(root, expression, performEvaluateConstants)

	return *expression
}

func performEvaluateConstants(index int, expression *Expression) (bool, Expression) {

	return false, *expression
}
