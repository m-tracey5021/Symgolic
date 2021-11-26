package interpretation

import (
	"errors"
	. "symgolic/evaluation"
	. "symgolic/search"
	. "symgolic/symbols"
)

func SubstituteFunctionDefs(program *Program) {

	for _, expression := range program.Expressions {

		root := expression.GetRoot()

		if expression.IsEquality(root) {

			lhs := expression.GetChildAtBreadth(root, 0)

			rhs := expression.GetChildAtBreadth(root, 1)

			if expression.IsFunction(lhs) && !expression.IsFunction(rhs) {

				functionName, paramMap, definition := MapFunctionDefParams(lhs, &expression)

				for i, search := range program.Expressions {

					SubstituteFunctionDefFor(search.GetRoot(), &search, functionName, paramMap, definition)

					program.Expressions[i] = search
				}

			} else if !expression.IsFunction(lhs) && expression.IsFunction(rhs) {

				functionName, paramMap, definition := MapFunctionDefParams(rhs, &expression)

				for _, search := range program.Expressions {

					SubstituteFunctionDefFor(search.GetRoot(), &search, functionName, paramMap, definition)
				}

			} else if expression.IsFunction(lhs) && expression.IsFunction(rhs) {

				panic(errors.New("function defined in terms of another"))

			} else {

				continue
			}
			// do something else if one function is defined in terms of another
		}
	}
}

func MapFunctionDefParams(index int, expression *Expression) (string, map[int][]int, *Expression) {

	paramMap := make(map[int][]int)

	if expression.IsFunctionDef(index) {

		definition := expression.CopySubtree(expression.GetSiblings(index)[0])

		for i, child := range expression.GetChildren(index) {

			paramMap[i] = SearchForInstancesOf(child, definition.GetRoot(), *expression, definition, make([]int, 0))
		}
		return expression.GetNode(index).AlphaValue, paramMap, &definition

	} else {

		return "", paramMap, nil
	}
}

func SubstituteFunctionDefFor(index int, expression *Expression, functionName string, paramMap map[int][]int, definition *Expression) {

	for _, child := range expression.GetChildren(index) {

		SubstituteFunctionDefFor(child, expression, functionName, paramMap, definition)
	}
	if expression.IsFunctionCall(index) && expression.GetNode(index).AlphaValue == functionName {

		ApplyFunctionParams(expression, index, paramMap, *definition)
	}
}

func ApplyFunctionParams(applyTo *Expression, functionCall int, paramMap map[int][]int, definition Expression) {

	params := applyTo.GetChildren(functionCall)

	for paramIndex, instances := range paramMap {

		for _, instance := range instances {

			definition.ReplaceNodeCascade(instance, applyTo.CopySubtree(params[paramIndex]), NewEmptyExpression())
		}
	}
	applyTo.ReplaceNodeCascade(functionCall, definition, NewEmptyExpression())
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

			InvokeFunction(expression.GetNode(index).AlphaValue, index, expression)
		}
	}
}

func InvokeFunction(command string, index int, expression *Expression) {

	functions := map[string]Evaluation{

		"ec": EvaluateConstants,

		"cancel": EvaluateCancellation,

		"distribute": EvaluateDistribution,

		"sumliketerms": EvaluateLikeTerms,

		"expandexponents": EvaluateExponentExpansion,

		"factor": EvaluateFactorisation,
	}
	call, exists := functions[command]

	if exists {

		EvaluateAndReplace(expression.GetRoot(), expression, call)

	} else {

		panic(errors.New("function " + command + " is not defined"))
	}
}
