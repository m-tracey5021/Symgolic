package evaluation

import (
	. "symgolic/symbols"
)

type Evaluation func(int, *Expression) (bool, Expression)

type IndeterminateEvaluation func(int, *Expression) (bool, []Expression)

func EvaluateAndReplace(index int, expression *Expression, evalFunc Evaluation) {

	for _, child := range expression.GetChildren(index) {

		EvaluateAndReplace(child, expression, evalFunc)
	}
	change, result := evalFunc(index, expression)

	if change {

		expression.ReplaceNodeCascade(index, result)
	}
}
