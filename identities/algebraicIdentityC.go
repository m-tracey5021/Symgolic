package identities

import (
	. "symgolic/language/components"
	. "symgolic/language/parsing"
)

// (a^2)-(b^2) = (a+b)*(a-b)

type AlgebraicIdentityC struct {
	A Expression

	B Expression

	Direction Direction

	IdentityRequisites []IdentityRequisite
}

func NewAlgebraicIdentityC(expression *Expression) AlgebraicIdentityC {

	identityRequisites := []IdentityRequisite{
		{

			Form: "(a+b)*(a-b)",

			Direction: Forwards,

			AlternateForms: []AlternateForm{

				{
					Form: "c*(a-b)", // where c = a + b

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
				{
					Form: "(a+b)*c", // where c = a - b

					Conditions: []FormCondition{

						{

							Target: expression.CopySubtree(expression.GetNodeByPath([]int{1})),

							EqualTo: ParseExpression("a-b"),

							Instances: [][]int{

								{1},
							},
						},
					},
				},
				{
					Form: "c*d", // where c = a + b and d = a - b

					Conditions: []FormCondition{

						{

							Target: expression.CopySubtree(expression.GetNodeByPath([]int{0})),

							EqualTo: ParseExpression("a+b"),

							Instances: [][]int{

								{0},
							},
						},
						{

							Target: expression.CopySubtree(expression.GetNodeByPath([]int{1})),

							EqualTo: ParseExpression("a-b"),

							Instances: [][]int{

								{1},
							},
						},
					},
				},
			},
		},
		{
			Form: "(a^2)-(b^2)",

			Direction: Backwards,

			AlternateForms: []AlternateForm{
				{
					Form: "c-(b^2)", // where c is constant and c = a^2

					Conditions: []FormCondition{

						{

							Target: expression.CopySubtree(expression.GetNodeByPath([]int{0})),

							EqualTo: ParseExpression("a^2"),

							Instances: [][]int{

								{0},
							},
						},
					},
				},
				{
					Form: "(a^2)-c", // where c is constant and c = b^2

					Conditions: []FormCondition{

						{

							Target: expression.CopySubtree(expression.GetNodeByPath([]int{1})),

							EqualTo: ParseExpression("b^2"),

							Instances: [][]int{

								{1},
							},
						},
					},
				},
				{
					Form: "c-d", // where c is constant and c = a^2 and d = b^2

					Conditions: []FormCondition{

						{

							Target: expression.CopySubtree(expression.GetNodeByPath([]int{0})),

							EqualTo: ParseExpression("a^2"),

							Instances: [][]int{

								{0},
							},
						},
						{

							Target: expression.CopySubtree(expression.GetNodeByPath([]int{1})),

							EqualTo: ParseExpression("b^2"),

							Instances: [][]int{

								{1},
							},
						},
					},
				},
			},
		},
	}
	return AlgebraicIdentityC{IdentityRequisites: identityRequisites}
}

func (a *AlgebraicIdentityC) AssignVariables(variableMap map[string]Expression, direction Direction) {

	a.A = variableMap["a"]

	a.B = variableMap["b"]

	a.Direction = direction
}

func (a *AlgebraicIdentityC) ApplyForwards(index int, expression *Expression) Expression {

	addRoot, add := NewExpression(NewOperation(Addition))

	expA := add.AppendNode(addRoot, NewOperation(Exponent))

	expB := add.AppendNode(addRoot, NewOperation(Exponent))

	add.AppendExpression(expA, a.A, false)

	add.AppendNode(expA, NewConstant(2))

	add.AppendExpression(expB, a.B, false)

	add.AppendNode(expB, NewConstant(2))

	return add
}

func (a *AlgebraicIdentityC) ApplyBackwards(index int, expression *Expression) Expression { // "(a^2)-(b^2)"

	mulRoot, mul := NewExpression(NewOperation(Multiplication))

	addA := mul.AppendNode(mulRoot, NewOperation(Addition))

	addB := mul.AppendNode(mulRoot, NewOperation(Addition))

	mul.AppendExpression(addA, a.A, false)

	mul.AppendExpression(addA, a.B, false)

	mul.AppendExpression(addB, a.A, true)

	b2 := mul.AppendExpression(addB, a.B, true)

	mul.AppendAuxiliariesAt(b2, []Symbol{NewOperation(Subtraction)})

	return mul
}

func (a *AlgebraicIdentityC) GetRequisites() []IdentityRequisite {

	return a.IdentityRequisites
}

func (a *AlgebraicIdentityC) GetDirection() Direction {

	return a.Direction
}
