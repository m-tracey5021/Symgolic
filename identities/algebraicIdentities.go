package identities

import (
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
	A int

	B int

	X int
}

func NewAlgebraicIdentityD() {

	identityRequisites := []IdentityRequisite{

		IdentityRequisite{Form: "(x^2)+((a+b)*x)+(a*b)", Direction: Forwards},

		IdentityRequisite{Form: "(x^2)+(c*x)+(a*b)", Direction: Forwards}, // where c = a + b

		IdentityRequisite{Form: "(x^2)+((a+b)*x)+c", Direction: Forwards}, // where c = a * b

		IdentityRequisite{Form: "(x^2)+(c*x)+d", Direction: Forwards}, // where there are variables y and z where y + z = c, and y * z = d

		IdentityRequisite{

			Form: "(x^2)+c+d", // where there are variables j, k and l where (j + k) * l = c, and j * k = d

			Direction: Forwards,

			ConstantChecks: []ConstantCheck{

				ConstantCheck{
					
				}	
			},
			
		}, 

		IdentityRequisite{Form: "(x+a)*(x+b)"},
	}
}