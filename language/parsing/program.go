package parsing

import (
	"errors"
	. "symgolic/language/components"
	. "symgolic/language/interpretation"
	. "symgolic/language/interpretation/algebra"
)

type Program struct {
	Input, Output []Expression

	FunctionDefs map[string]bool
}

func NewProgram() Program {

	return Program{Input: make([]Expression, 0), Output: make([]Expression, 0), FunctionDefs: make(map[string]bool)}
}

func (p *Program) AddExpression(expression Expression) {

	p.Input = append(p.Input, expression)
}

func (p *Program) SubstituteVariableDefs() {

	for i, expression := range p.Input {

		root := expression.GetRoot()

		if expression.IsAssignment(root) {

			lhs := expression.GetChildAtBreadth(root, 0)

			rhs := expression.GetChildAtBreadth(root, 1)

			if expression.IsVariable(lhs) {

				variable := expression.GetNode(lhs).AlphaValue

				definition := expression.CopySubtree(rhs)

				for j, search := range p.Input {

					if i == j {

						continue
					}
					SubstituteVariableDefFor(search.GetRoot(), &search, &definition, variable)

					p.Input[j] = search
				}

			} else {

				continue
			}
			// do something else if one function is defined in terms of another
		}
	}
}

func (p *Program) SubstituteFunctionDefs() {

	for i, expression := range p.Input {

		root := expression.GetRoot()

		if expression.IsAssignment(root) {

			lhs := expression.GetChildAtBreadth(root, 0)

			rhs := expression.GetChildAtBreadth(root, 1)

			if expression.IsFunction(lhs) && !expression.IsFunction(rhs) {

				functionName, paramMap, definition := MapFunctionDefParams(From(expression).At(lhs))

				p.FunctionDefs[functionName] = true

				for j, search := range p.Input {

					if i == j {

						continue
					}
					SubstituteFunctionDefFor(search.GetRoot(), &search, functionName, paramMap, definition)

					p.Input[j] = search
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

func SubstituteVariableDefFor(index int, expression, definition *Expression, variable string) {

	for _, child := range expression.GetChildren(index) {

		SubstituteVariableDefFor(child, expression, definition, variable)
	}
	if expression.IsVariable(index) && expression.GetNode(index).AlphaValue == variable {

		expression.ReplaceNodeCascade(index, *definition)
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

func MapFunctionDefParams(target ExpressionIndex) (string, map[int][]int, *Expression) {

	paramMap := make(map[int][]int)

	if target.Expression.IsFunctionDef(target.Index) {

		definition := target.Expression.CopySubtree(target.Expression.GetSiblings(target.Index)[0])

		for i, child := range target.Expression.GetChildren(target.Index) {

			paramMap[i] = SearchForInstancesOf(target.At(child), From(definition), make([]int, 0))
		}
		return target.Expression.GetNode(target.Index).AlphaValue, paramMap, &definition

	} else {

		return "", paramMap, nil
	}
}

func ApplyFunctionParams(applyTo *Expression, functionCall int, paramMap map[int][]int, definition Expression) {

	params := applyTo.GetChildren(functionCall)

	for paramIndex, instances := range paramMap {

		for _, instance := range instances {

			definition.ReplaceNodeCascade(instance, applyTo.CopySubtree(params[paramIndex]))
		}
	}
	applyTo.ReplaceNodeCascade(functionCall, definition)
}

func (p *Program) Interpret() []Expression {

	results := make([]Expression, 0)

	p.SubstituteFunctionDefs()

	p.SubstituteVariableDefs()

	for _, expression := range p.Input {

		EvaluateAndReplace(From(expression), ApplyArithmetic)

		p.InterpretExpression(&expression)

		results = append(results, expression)
	}
	return results
}

func (p *Program) InterpretExpression(expression *Expression) {

	p.SearchFunctions(From(*expression))
}

func (p *Program) SearchFunctions(target ExpressionIndex) {

	for _, child := range target.Expression.GetChildren(target.Index) {

		p.SearchFunctions(target.At(child))
	}
	if target.Expression.IsFunctionCall(target.Index) {

		functionName := target.Expression.GetNode(target.Index).AlphaValue

		_, functionDefined := p.FunctionDefs[functionName]

		if !functionDefined {

			p.InvokePredefinedFunction(functionName, target.Index, &target.Expression)
		}

	}
}

func (p *Program) InvokePredefinedFunction(command string, index int, expression *Expression) {

	arguments := expression.GetChildren(index)

	input := make([]Expression, 0)

	for _, arg := range arguments {

		input = append(input, expression.CopySubtree(arg))
	}
	evaluationFunctions := map[string]Evaluation{

		"applyarithmetic": ApplyArithmetic,

		"cancel": Cancel,

		"distribute": Distribute,

		"sumliketerms": SumLikeTerms,

		"expandexponents": ExpandExponents,

		"factor": Factor,
	}
	evaluationAgainstFunctions := map[string]EvaluationAgainst{

		"dot": DotProduct,

		"cross": CrossProduct,
	}
	call, exists := evaluationFunctions[command]

	if exists {

		EvaluateAndReplace(From(input[0]), call)

		p.Output = append(p.Output, input[0]) // everything in input is copied so just return the value modified in place

	} else {

		call, exists := evaluationAgainstFunctions[command]

		if exists {

			_, output := call(From(input[0]), From(input[1]))

			p.Output = append(p.Output, output)

		} else {

			panic(errors.New("function " + command + " is not defined"))
		}
	}
}
