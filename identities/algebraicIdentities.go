package identities

import (
	// . "symgolic/solvers"
	. "symgolic/parsing"
	. "symgolic/symbols"
)

// (a+b)^2 = (a^2)+(2*a*b)+(b^2) || (2^2)+12+(3^2) || (a^2)+(6*a)+(3^2) || (2^2)+(4*b)+(b^2)

type AlgebraicIdentityA struct {
	A Expression

	B Expression

	Direction Direction

	IdentityRequisites []IdentityRequisite
}

func NewAlgebraicIdentityA(expression *Expression) AlgebraicIdentityA {

	identityRequisites := []IdentityRequisite{

		IdentityRequisite{Form: "(a^2)+(2*a*b)+(b^2)", Direction: Forwards},

		IdentityRequisite{

			Form: "(a^2)+c+(b^2)",

			Direction: Forwards,

			ConstantChecks: []ConstantCheck{

				ConstantCheck{

					Values: []int{

						2,

						expression.GetNumericValueByPath([]int{0, 0}),

						expression.GetNumericValueByPath([]int{2, 0}),
					},
					Target: expression.GetNumericValueByPath([]int{1}),

					Operation: Multiplication,
				},
			},
		},
		IdentityRequisite{

			Form: "(a^2)+(c*a)+(b^2)",

			Direction: Forwards,

			ConstantChecks: []ConstantCheck{

				ConstantCheck{

					Values: []int{

						2,

						expression.GetNumericValueByPath([]int{2, 0}),
					},
					Target: expression.GetNumericValueByPath([]int{1, 0}),

					Operation: Multiplication,
				},
			},
		},
		IdentityRequisite{

			Form: "(a^2)+(c*b)+(b^2)",

			Direction: Forwards,

			ConstantChecks: []ConstantCheck{

				ConstantCheck{

					Values: []int{

						2,

						expression.GetNumericValueByPath([]int{0, 0}),
					},
					Target: expression.GetNumericValueByPath([]int{1, 0}),

					Operation: Multiplication,
				},
			},
		},
		IdentityRequisite{Form: "(a+b)^2", Direction: Backwards},
	}
	return AlgebraicIdentityA{IdentityRequisites: identityRequisites}
}

func (a *AlgebraicIdentityA) AssignVariables(variableMap map[string]Expression, direction Direction) {

	a.A = variableMap["a"]

	a.B = variableMap["b"]

	a.Direction = direction
}

// func (a *AlgebraicIdentityA) Identify(index int, expression *Expression) bool {

// 	return Identify(index, expression, a.IdentityRequisites, a.AssignVariables)
// }

func (a *AlgebraicIdentityA) ApplyForwards(index int, expression *Expression) Expression {

	exponentRoot, exponent := NewExpressionWithRoot(Symbol{Exponent, -1, "^"})

	add := exponent.AppendNode(exponentRoot, Symbol{Addition, -1, "+"})

	sumOperands := []Expression{a.A, a.B}

	exponent.AppendBulkExpressions(add, sumOperands)

	exponent.AppendNode(exponentRoot, Symbol{Constant, 2, "2"})

	return exponent

}

func (a *AlgebraicIdentityA) ApplyBackwards(index int, expression *Expression) Expression {

	sumRoot, sum := NewExpressionWithRoot(Symbol{Addition, -1, "+"})

	exponentA := sum.AppendNode(sumRoot, Symbol{Exponent, -1, "^"})

	mul := sum.AppendNode(sumRoot, Symbol{Multiplication, -1, "*"})

	exponentB := sum.AppendNode(sumRoot, Symbol{Exponent, -1, "^"})

	sum.AppendExpression(exponentA, a.A, false)

	sum.AppendNode(exponentA, Symbol{Constant, 2, "2"})

	sum.AppendNode(mul, Symbol{Constant, 2, "2"})

	sum.AppendExpression(mul, a.A, false)

	sum.AppendExpression(mul, a.B, false)

	sum.AppendExpression(exponentB, a.B, false)

	sum.AppendNode(exponentB, Symbol{Constant, 2, "2"})

	return sum

}

func (a *AlgebraicIdentityA) GetRequisites() []IdentityRequisite {

	return a.IdentityRequisites
}

func (a *AlgebraicIdentityA) GetDirection() Direction {

	return a.Direction
}

// (a-b)^2=(a^2)-(2*a*b)+(b^2)

// (a^2)-(b^2)=(a+b)*(a-b)

// (x+a)*(x+b) = (x^2)+((a+b)*x)+(a*b)

type AlgebraicIdentityD struct {
	A Expression

	B Expression

	X Expression

	Direction Direction

	IdentityRequisites []IdentityRequisite
}

func NewAlgebraicIdentityD(expression *Expression) AlgebraicIdentityD {

	identityRequisites := []IdentityRequisite{

		IdentityRequisite{

			Form: "(x^2)+((a+b)*x)+(a*b)",

			Direction: Forwards,

			AlternateForms: []AlternateForm{

				{
					Form: "(x^2)+(c*x)+(a*b)", // where c is constant and c = a + b

					Conditions: []FormCondition{

						FormCondition{

							Target: expression.CopySubtree(expression.GetChildByPath([]int{1, 0})),

							EqualTo: ParseExpression("a+b"),

							Instances: [][]int{

								{1, 0},
							},
						},
					},
				},
				{
					Form: "(x^2)+((a+b)*x)+c", // where c = a * b

					Conditions: []FormCondition{

						FormCondition{

							Target: expression.CopySubtree(expression.GetChildByPath([]int{2})),

							EqualTo: ParseExpression("a*b"),

							Instances: [][]int{

								{2},
							},
						},
					},
				},
				{
					Form: "(x^2)+(c*x)+d", // where there are variables y and z where y + z = c, and y * z = d

					Conditions: []FormCondition{

						FormCondition{

							Target: expression.CopySubtree(expression.GetChildByPath([]int{1, 0})),

							EqualTo: ParseExpression("a+b"),

							Instances: [][]int{

								{1, 0},
							},
						},
						FormCondition{

							Target: expression.CopySubtree(expression.GetChildByPath([]int{2})),

							EqualTo: ParseExpression("a*b"),

							Instances: [][]int{

								{2},
							},
						},
					},
				},
				{
					Form: "(x^2)+c+d", // where there are variables j, k and l where (j + k) * l = c, j * k = d and l = x

					Conditions: []FormCondition{

						FormCondition{

							Target: expression.CopySubtree(expression.GetChildByPath([]int{1})),

							EqualTo: ParseExpression("(a+b)*x"),

							Instances: [][]int{

								{1},
							},
						},
						FormCondition{

							Target: expression.CopySubtree(expression.GetChildByPath([]int{2})),

							EqualTo: ParseExpression("a*b"),

							Instances: [][]int{

								{2},
							},
						},
					},
				},
			},
		},
		IdentityRequisite{

			Form: "(x+a)*(x+b)",

			Direction: Backwards,

			AlternateForms: []AlternateForm{

				{
					Form: "(x+a)*c", // where c = x + b

					Conditions: []FormCondition{

						FormCondition{

							Target: expression.CopySubtree(expression.GetChildByPath([]int{1})),

							EqualTo: ParseExpression("x+b"),

							Instances: [][]int{

								{1},
							},
						},
					},
				},
				{
					Form: "c*(x+b)", // where c = x + a

					Conditions: []FormCondition{

						FormCondition{

							Target: expression.CopySubtree(expression.GetChildByPath([]int{0})),

							EqualTo: ParseExpression("x+a"),

							Instances: [][]int{

								{0},
							},
						},
					},
				},
				{
					Form: "c*d", // where c = x + a and d = x + b

					Conditions: []FormCondition{

						FormCondition{

							Target: expression.CopySubtree(expression.GetChildByPath([]int{0})),

							EqualTo: ParseExpression("x+a"),

							Instances: [][]int{

								{0},
							},
						},
						FormCondition{

							Target: expression.CopySubtree(expression.GetChildByPath([]int{1})),

							EqualTo: ParseExpression("x+b"),

							Instances: [][]int{

								{1},
							},
						},
					},
				},
			},
		},
	}
	return AlgebraicIdentityD{IdentityRequisites: identityRequisites}
}

func (a *AlgebraicIdentityD) AssignVariables(variableMap map[string]Expression, direction Direction) {

	a.A = variableMap["a"]

	a.B = variableMap["b"]

	a.X = variableMap["x"]

	a.Direction = direction
}

func (a *AlgebraicIdentityD) ApplyForwards(index int, expression *Expression) Expression {

	mulRoot, mul := NewExpressionWithRoot(Symbol{Multiplication, -1, "*"})

	addA := mul.AppendNode(mulRoot, Symbol{Addition, -1, "+"})

	addB := mul.AppendNode(mulRoot, Symbol{Addition, -1, "+"})

	mul.AppendExpression(addA, a.X, false)

	mul.AppendExpression(addA, a.A, false)

	mul.AppendExpression(addB, a.X, false)

	mul.AppendExpression(addB, a.B, false)

	return mul

}

func (a *AlgebraicIdentityD) ApplyBackwards(index int, expression *Expression) Expression { // "(x^2)+((a+b)*x)+(a*b)"

	resultRoot, result := NewExpressionWithRoot(Symbol{Addition, -1, "+"})

	exp := result.AppendNode(resultRoot, Symbol{Exponent, -1, "^"})

	result.AppendExpression(exp, a.X, false)

	result.AppendNode(exp, Symbol{Constant, 2, "2"})

	mulA := result.AppendNode(resultRoot, Symbol{Multiplication, -1, "*"})

	innerresult := result.AppendNode(mulA, Symbol{Addition, -1, "+"})

	result.AppendExpression(innerresult, a.A, false)

	result.AppendExpression(innerresult, a.B, false)

	result.AppendExpression(mulA, a.X, false)

	mulB := result.AppendNode(resultRoot, Symbol{Multiplication, -1, "*"})

	result.AppendExpression(mulB, a.A, false)

	result.AppendExpression(mulB, a.B, false)

	return result
}

func (a *AlgebraicIdentityD) GetRequisites() []IdentityRequisite {

	return a.IdentityRequisites
}

func (a *AlgebraicIdentityD) GetDirection() Direction {

	return a.Direction
}
