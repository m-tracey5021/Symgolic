package identities

import (
	. "symgolic/language/components"
	. "symgolic/language/parsing"
)

type AlgebraicIdentityA struct {
	A Expression

	B Expression

	Direction Direction

	IdentityRequisites []IdentityRequisite
}

func NewAlgebraicIdentityA(expression *Expression) AlgebraicIdentityA {

	identityRequisites := []IdentityRequisite{

		{
			Form: "(a^2)+(2*a*b)+(b^2)",

			Direction: Forwards,

			AlternateForms: []AlternateForm{

				{
					Form: "(a^2)+(c*a)+(b^2)", // where c = 2 * b

					Conditions: []FormCondition{

						{
							Target: expression.CopySubtree(expression.GetNodeByPath([]int{1, 0})),

							EqualTo: ParseExpression("2*b"),

							Instances: [][]int{

								{1, 0},
							},
						},
					},
				},
				{
					Form: "(a^2)+(c*b)+(b^2)", // where c = 2 * a

					Conditions: []FormCondition{

						{
							Target: expression.CopySubtree(expression.GetNodeByPath([]int{1, 0})),

							EqualTo: ParseExpression("2*a"),

							Instances: [][]int{

								{1, 0},
							},
						},
					},
				},
				{
					Form: "(a^2)+c+(b^2)", // where c = 2 * a * b

					Conditions: []FormCondition{

						{
							Target: expression.CopySubtree(expression.GetNodeByPath([]int{1})),

							EqualTo: ParseExpression("2*a*b"),

							Instances: [][]int{

								{1},
							},
						},
					},
				},
			},
		},
		{

			Form: "(a+b)^2",

			Direction: Backwards,

			AlternateForms: []AlternateForm{

				{
					Form: "c^2", // where c = a + b

					Conditions: []FormCondition{

						{
							Target: expression.CopySubtree(expression.GetNodeByPath([]int{0})),

							EqualTo: ParseExpression("a+b"),

							Instances: [][]int{

								{0},
							},
						},
					},
				},
			},
		},
	}
	return AlgebraicIdentityA{IdentityRequisites: identityRequisites}
}

func (a *AlgebraicIdentityA) AssignVariables(variableMap map[string]Expression, direction Direction) {

	a.A = variableMap["a"]

	a.B = variableMap["b"]

	a.Direction = direction
}

func (a *AlgebraicIdentityA) ApplyForwards(index int, expression *Expression) Expression {

	exponentRoot, exponent := NewExpression(NewOperation(Exponent))

	add := exponent.AppendNode(exponentRoot, NewOperation(Addition))

	sumOperands := []Expression{a.A, a.B}

	exponent.AppendBulkExpressions(add, sumOperands)

	exponent.AppendNode(exponentRoot, NewConstant(2))

	return exponent

}

func (a *AlgebraicIdentityA) ApplyBackwards(index int, expression *Expression) Expression {

	sumRoot, sum := NewExpression(NewOperation(Addition))

	exponentA := sum.AppendNode(sumRoot, NewOperation(Exponent))

	mul := sum.AppendNode(sumRoot, NewOperation(Multiplication))

	exponentB := sum.AppendNode(sumRoot, NewOperation(Exponent))

	sum.AppendExpression(exponentA, a.A, false)

	sum.AppendNode(exponentA, NewConstant(2))

	sum.AppendNode(mul, NewConstant(2))

	sum.AppendExpression(mul, a.A, false)

	sum.AppendExpression(mul, a.B, false)

	sum.AppendExpression(exponentB, a.B, false)

	sum.AppendNode(exponentB, NewConstant(2))

	return sum

}

func (a *AlgebraicIdentityA) GetRequisites() []IdentityRequisite {

	return a.IdentityRequisites
}

func (a *AlgebraicIdentityA) GetDirection() Direction {

	return a.Direction
}
