package identities

import (
	. "symgolic/language/components"
	. "symgolic/language/parsing"
)

type AlgebraicIdentityE struct {
	A Expression

	B Expression

	C Expression

	Direction Direction

	IdentityRequisites []IdentityRequisite
}

func NewAlgebraicIdentityE(expression *Expression) AlgebraicIdentityE {

	identityRequisites := []IdentityRequisite{

		{
			Form: "(a^2)+(b^2)+(c^2)+(2*a*b)+(2*b*c)+(2*a*c)",

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
	return AlgebraicIdentityE{IdentityRequisites: identityRequisites}
}

func (a *AlgebraicIdentityE) AssignVariables(variableMap map[string]Expression, direction Direction) {

	a.A = variableMap["a"]

	a.B = variableMap["b"]

	a.C = variableMap["c"]

	a.Direction = direction
}

func (a *AlgebraicIdentityE) ApplyForwards(index int, expression *Expression) Expression {

	exponentRoot, exponent := NewExpression(NewOperation(Exponent))

	add := exponent.AppendNode(exponentRoot, NewOperation(Addition))

	sumOperands := []Expression{a.A, a.B, a.C}

	exponent.AppendBulkExpressions(add, sumOperands)

	exponent.AppendNode(exponentRoot, NewConstant(2))

	return exponent

}

func (a *AlgebraicIdentityE) ApplyBackwards(index int, expression *Expression) Expression {

	sumRoot, sum := NewExpression(NewOperation(Addition))

	exponentA := sum.AppendNode(sumRoot, NewOperation(Exponent))

	sum.AppendExpression(exponentA, a.A, false)

	sum.AppendNode(exponentA, NewConstant(2))

	exponentB := sum.AppendNode(sumRoot, NewOperation(Exponent))

	sum.AppendExpression(exponentB, a.B, false)

	sum.AppendNode(exponentB, NewConstant(2))

	exponentC := sum.AppendNode(sumRoot, NewOperation(Exponent))

	sum.AppendExpression(exponentC, a.C, false)

	sum.AppendNode(exponentC, NewConstant(2))

	mulA := sum.AppendNode(sumRoot, NewOperation(Multiplication))

	sum.AppendNode(mulA, NewConstant(2))

	sum.AppendExpression(mulA, a.A, false)

	sum.AppendExpression(mulA, a.B, false)

	mulB := sum.AppendNode(sumRoot, NewOperation(Multiplication))

	sum.AppendNode(mulB, NewConstant(2))

	sum.AppendExpression(mulB, a.B, false)

	sum.AppendExpression(mulB, a.C, false)

	mulC := sum.AppendNode(sumRoot, NewOperation(Multiplication))

	sum.AppendNode(mulC, NewConstant(2))

	sum.AppendExpression(mulC, a.C, false)

	sum.AppendExpression(mulC, a.A, false)

	return sum

}

func (a *AlgebraicIdentityE) GetRequisites() []IdentityRequisite {

	return a.IdentityRequisites
}

func (a *AlgebraicIdentityE) GetDirection() Direction {

	return a.Direction
}
