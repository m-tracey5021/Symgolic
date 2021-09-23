package interpretation

import (
	"errors"
	. "symgolic/evaluation"
	. "symgolic/symbols"
)

func SubstituteFunctionDefs(program *Program) {

	for _, expression := range program.Expressions {

		root := expression.GetRoot()

		if expression.IsEquality(root) {

			lhs := expression.GetChildAtBreadth(root, 0)

			rhs := expression.GetChildAtBreadth(root, 1)

			if expression.IsFunction(lhs) && !expression.IsFunction(rhs) {

				target := expression.CopySubtree(rhs)

				for _, search := range program.Expressions {

					SubstituteFunctionDefsFor(search.GetRoot(), search, expression.GetAlphaValueByIndex(lhs), target)
				}

			} else if !expression.IsFunction(lhs) && expression.IsFunction(rhs) {

				target := expression.CopySubtree(lhs)

				for _, search := range program.Expressions {

					SubstituteFunctionDefsFor(search.GetRoot(), search, expression.GetAlphaValueByIndex(rhs), target)
				}

			} else {

				continue
			}
			// do something else if one function is defined in terms of another
		}
	}
}

func SubstituteFunctionDefsFor(index int, expression Expression, function string, target Expression) {

	for _, child := range expression.GetChildren(index) {

		SubstituteFunctionDefsFor(child, expression, function, target)
	}
	if expression.IsFunctionCall(index) && expression.GetAlphaValueByIndex(index) == function {

		expression.ReplaceNodeCascade(index, target)
	}
}

func InterpretProgram(program *Program) {

	SubstituteFunctionDefs(program)
}

func InterpretExpression(expression *Expression) {

	root := expression.GetRoot()

	SearchFunctions(root, expression)
}

func SearchFunctions(index int, expression *Expression) {

	for _, child := range expression.GetChildren(index) {

		SearchFunctions(child, expression)
	}
	if expression.IsFunction(index) {

		parent := expression.GetParent(index)

		if !expression.IsEquality(parent) {

			InvokeFunction(expression.GetNodeByIndex(index).AlphaValue, index, expression)
		}
	}
}

func InvokeFunction(command string, index int, expression *Expression) {

	functions := map[string]Evaluation{

		"ec": EvaluateConstants,

		"cancel": EvaluateCancellation,

		"distribute": EvaluateDistribution,

		"sumliketerms": EvaluateLikeTerms,
	}
	call, exists := functions[command]

	if exists {

		EvaluateAndReplace(expression.GetRoot(), expression, call)

	} else {

		panic(errors.New("function " + command + " is not defined"))
	}
}
