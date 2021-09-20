package interpretation

import (
	. "symgolic/evaluation"
	. "symgolic/symbols"
)

func InterpretExpression(expression *Expression) {

	root := expression.GetRoot()

	SearchFunctions(root, expression)
}

func SearchFunctions(index int, expression *Expression) {

	for _, child := range expression.GetChildren(index) {

		SearchFunctions(child, expression)
	}
	if expression.IsFunction(index) {

		InvokeFunction(expression.GetNodeByIndex(index).AlphaValue, index, expression)
	}
}

func InvokeFunction(command string, index int, expression *Expression) {

	functions := map[string]Evaluation{

		"ec": EvaluateConstants,
	}
	EvaluateAndReplace(expression.GetRoot(), expression, functions[command])
}
