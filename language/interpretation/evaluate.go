package interpretation

import (
	. "symgolic/language/components"
)

// produces one result

type Evaluation func(ExpressionIndex) (bool, Expression)

type EvaluationInPlace func(ExpressionIndex)

type EvaluationAgainst func(ExpressionIndex, ExpressionIndex) (bool, Expression)

type EvaluationForMany func(...ExpressionIndex) (bool, Expression)

// produces many results

type IndeterminateEvaluation func(ExpressionIndex) (bool, []Expression)

type IndeterminateEvaluationInPlace func(ExpressionIndex)

type IndeterminateEvaluationAgainst func(ExpressionIndex, ExpressionIndex) (bool, []Expression)

type IndeterminateEvaluationForMany func(...ExpressionIndex) (bool, []Expression)

func EvaluateAndReplace(target ExpressionIndex, evalFunc Evaluation) {

	for _, child := range target.Expression.GetChildren(target.Index) {

		EvaluateAndReplace(target.At(child), evalFunc)
	}
	change, result := evalFunc(target)

	if change {

		target.Expression.ReplaceNodeCascade(target.Index, result)
	}
}
