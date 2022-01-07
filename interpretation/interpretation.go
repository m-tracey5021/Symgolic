package interpretation

// import (
// 	"errors"
// 	. "symgolic/evaluation"
// 	. "symgolic/evaluation/linearAlgebra"
// 	. "symgolic/search"
// 	. "symgolic/symbols"
// )

// func SubstituteVariableDefs(program *Program) {

// 	for i, expression := range program.Expressions {

// 		root := expression.GetRoot()

// 		if expression.IsAssignment(root) {

// 			lhs := expression.GetChildAtBreadth(root, 0)

// 			rhs := expression.GetChildAtBreadth(root, 1)

// 			if expression.IsVariable(lhs) {

// 				variable := expression.GetNode(lhs).AlphaValue

// 				definition := expression.CopySubtree(rhs)

// 				for j, search := range program.Expressions {

// 					if i == j {

// 						continue
// 					}
// 					SubstituteVariableDefFor(search.GetRoot(), &search, &definition, variable)

// 					program.Expressions[j] = search
// 				}

// 			} else {

// 				continue
// 			}
// 			// do something else if one function is defined in terms of another
// 		}
// 	}
// }

// func SubstituteFunctionDefs(program *Program) {

// 	for i, expression := range program.Expressions {

// 		root := expression.GetRoot()

// 		if expression.IsAssignment(root) {

// 			lhs := expression.GetChildAtBreadth(root, 0)

// 			rhs := expression.GetChildAtBreadth(root, 1)

// 			if expression.IsFunction(lhs) && !expression.IsFunction(rhs) {

// 				functionName, paramMap, definition := MapFunctionDefParams(lhs, &expression)

// 				program.FunctionDefs[functionName] = true

// 				for j, search := range program.Expressions {

// 					if i == j {

// 						continue
// 					}
// 					SubstituteFunctionDefFor(search.GetRoot(), &search, functionName, paramMap, definition)

// 					program.Expressions[j] = search
// 				}

// 			} else if expression.IsFunction(lhs) && expression.IsFunction(rhs) {

// 				panic(errors.New("function defined in terms of another"))

// 			} else {

// 				continue
// 			}
// 			// do something else if one function is defined in terms of another
// 		}
// 	}
// }

// func MapFunctionDefParams(index int, expression *Expression) (string, map[int][]int, *Expression) {

// 	paramMap := make(map[int][]int)

// 	if expression.IsFunctionDef(index) {

// 		definition := expression.CopySubtree(expression.GetSiblings(index)[0])

// 		for i, child := range expression.GetChildren(index) {

// 			paramMap[i] = SearchForInstancesOf(child, definition.GetRoot(), *expression, definition, make([]int, 0))
// 		}
// 		return expression.GetNode(index).AlphaValue, paramMap, &definition

// 	} else {

// 		return "", paramMap, nil
// 	}
// }

// func SubstituteVariableDefFor(index int, expression, definition *Expression, variable string) {

// 	for _, child := range expression.GetChildren(index) {

// 		SubstituteVariableDefFor(child, expression, definition, variable)
// 	}
// 	if expression.IsVariable(index) && expression.GetNode(index).AlphaValue == variable {

// 		expression.ReplaceNodeCascade(index, *definition)
// 	}
// }

// func SubstituteFunctionDefFor(index int, expression *Expression, functionName string, paramMap map[int][]int, definition *Expression) {

// 	for _, child := range expression.GetChildren(index) {

// 		SubstituteFunctionDefFor(child, expression, functionName, paramMap, definition)
// 	}
// 	if expression.IsFunctionCall(index) && expression.GetNode(index).AlphaValue == functionName {

// 		ApplyFunctionParams(expression, index, paramMap, *definition)
// 	}
// }

// func ApplyFunctionParams(applyTo *Expression, functionCall int, paramMap map[int][]int, definition Expression) {

// 	params := applyTo.GetChildren(functionCall)

// 	for paramIndex, instances := range paramMap {

// 		for _, instance := range instances {

// 			definition.ReplaceNodeCascade(instance, applyTo.CopySubtree(params[paramIndex]))
// 		}
// 	}
// 	applyTo.ReplaceNodeCascade(functionCall, definition)
// }

// func InterpretProgram(program *Program) []Expression {

// 	results := make([]Expression, 0)

// 	SubstituteFunctionDefs(program)

// 	SubstituteVariableDefs(program)

// 	for _, expression := range program.Expressions {

// 		EvaluateAndReplace(expression.GetRoot(), &expression, ApplyArithmetic)

// 		InterpretExpression(&expression, program)

// 		results = append(results, expression)
// 	}
// 	return results
// }

// func InterpretExpression(expression *Expression, program *Program) {

// 	root := expression.GetRoot()

// 	SearchFunctions(root, expression, program)
// }

// func SearchFunctions(index int, expression *Expression, program *Program) {

// 	for _, child := range expression.GetChildren(index) {

// 		SearchFunctions(child, expression, program)
// 	}
// 	if expression.IsFunctionCall(index) {

// 		functionName := expression.GetNode(index).AlphaValue

// 		_, functionDefined := program.FunctionDefs[functionName]

// 		if !functionDefined {

// 			output := InvokePredefinedFunction(functionName, index, expression)
// 		}

// 	}
// }

// func InvokePredefinedFunction(command string, index int, expression *Expression) Expression {

// 	arguments := expression.GetChildren(index)

// 	input := make([]Expression, 0)

// 	for _, arg := range arguments {

// 		input = append(input, expression.CopySubtree(arg))
// 	}
// 	evaluationFunctions := map[string]Evaluation{

// 		"applyarithmetic": ApplyArithmetic,

// 		"cancel": Cancel,

// 		"distribute": Distribute,

// 		"sumliketerms": SumLikeTerms,

// 		"expandexponents": ExpandExponents,

// 		"factor": Factor,
// 	}

// 	evaluationAgainstFunctions := map[string]EvaluationAgainst{

// 		"dot": DotProduct,

// 		"cross": CrossProduct,
// 	}
// 	call, exists := evaluationFunctions[command]

// 	if exists {

// 		EvaluateAndReplace(input[0].GetRoot(), &input[0], call)

// 		return input[0] // everything in input is copied so just return the value modified in place

// 	} else {

// 		call, exists := evaluationAgainstFunctions[command]

// 		if exists {

// 			_, output := call(input[0].GetRoot(), input[1].GetRoot(), &input[0], &input[1])

// 			return output

// 		} else {

// 			panic(errors.New("function " + command + " is not defined"))
// 		}
// 	}
// }
