package identities

import (
	. "symgolic/parsing"
	. "symgolic/symbols"
)

// (a^2)-(b^2) = (a+b)*(a-b)

type DifferenceOfTwoSquares struct {
	A Expression

	B Expression

	Direction Direction

	IdentityRequisites []IdentityRequisite
}

func NewDifferenceOfTwoSquares(expression *Expression) DifferenceOfTwoSquares {

	identityRequisites := []IdentityRequisite{

		{

			Form: "(a^2)-(b^2)",

			Direction: Forwards,

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
		{

			Form: "(a+b)*(a-b)",

			Direction: Backwards,

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
	}
	return DifferenceOfTwoSquares{IdentityRequisites: identityRequisites}
}

func (d *DifferenceOfTwoSquares) AssignVariables(variableMap map[string]Expression, direction Direction) {

	d.A = variableMap["a"]

	d.B = variableMap["b"]

	d.Direction = direction
}

func (d *DifferenceOfTwoSquares) ApplyForwards(index int, expression *Expression) Expression {

	mulRoot, mul := NewExpression(NewOperation(Multiplication))

	addA := mul.AppendNode(mulRoot, NewOperation(Addition))

	addB := mul.AppendNode(mulRoot, NewOperation(Addition))

	mul.AppendExpression(addA, d.A, false)

	mul.AppendExpression(addA, d.B, false)

	mul.AppendExpression(addB, d.A, true)

	b2 := mul.AppendExpression(addB, d.B, true)

	mul.AppendAuxiliariesAt(b2, []Symbol{NewOperation(Subtraction)})

	return mul
}

func (d *DifferenceOfTwoSquares) ApplyBackwards(index int, expression *Expression) Expression { // "(a^2)-(b^2)"

	addRoot, add := NewExpression(NewOperation(Addition))

	expA := add.AppendNode(addRoot, NewOperation(Exponent))

	expB := add.AppendNode(addRoot, NewOperation(Exponent))

	add.AppendExpression(expA, d.A, false)

	add.AppendNode(expA, NewConstant(2))

	add.AppendExpression(expB, d.B, false)

	add.AppendNode(expB, NewConstant(2))

	return add
}

func (d *DifferenceOfTwoSquares) GetRequisites() []IdentityRequisite {

	return d.IdentityRequisites
}

func (d *DifferenceOfTwoSquares) GetDirection() Direction {

	return d.Direction
}
