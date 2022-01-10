package interpretation

import (
	. "symgolic/language/components"
)

// evaluates in place

type Evaluation func(int, *Expression) (bool, Expression)

// operates on two expressions

type EvaluationAgainst func(int, int, *Expression, *Expression) (bool, Expression)

// produces many results

type IndeterminateEvaluation func(int, *Expression) (bool, []Expression)

type IndeterminateEvaluationAgainst func(int, int, *Expression, *Expression) (bool, []Expression)

func EvaluateAndReplace(index int, expression *Expression, evalFunc Evaluation) {

	for _, child := range expression.GetChildren(index) {

		EvaluateAndReplace(child, expression, evalFunc)
	}
	change, result := evalFunc(index, expression)

	if change {

		expression.ReplaceNodeCascade(index, result)
	}
}
