package identities

import (
	. "symgolic/comparison"
	. "symgolic/parsing"
	. "symgolic/symbols"
)

// (a+b)^2 = (a^2)+(2*a*b)+(b^2) || (2^2)+12+(3^2) || (a^2)+(6*a)+(3^2) || (2^2)+(4*b)+(b^2)

type AlgebraicIdentityA struct {
	A int

	B int

	IdentityRequisites []IdentityRequisite
}

func NewAlgebraicIdentityA(expression *Expression) AlgebraicIdentityA {

	identityRequisites := []IdentityRequisite{

		IdentityRequisite{Form: "(a^2)+(2*a*b)+(b^2)"},

		IdentityRequisite{

			Form: "(a^2)+c+(b^2)",

			ConstantChecks: []ConstantCheck{

				ConstantCheck{

					Values: []int{

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
	}
	return AlgebraicIdentityA{IdentityRequisites: identityRequisites}
}

func (a *AlgebraicIdentityA) Identify(index int, expression *Expression) bool {

	for _, requisite := range a.IdentityRequisites {

		form := ParseExpression(requisite.Form)

		formApplies := IsEqualByForm(form, *expression)

		if formApplies {

			if len(requisite.ConstantChecks) != 0 {

				for _, check := range requisite.ConstantChecks {

					if CheckConstantValue(check.Values, check.Target, check.Operation, expression) {

						// assign indexes to struct

						return true
					}
				}

			} else {

				return true
			}
		}
	}
	return false
}

// func (a *AlgebraicIdentityA) Identify(index int, expression *Expression) bool {

// 	root := expression.GetRoot()

// 	formA := ParseExpression("(a^2)+(2*a*b)+(b^2)") // standard form

// 	formAApplies := IsEqualByFormAt(formA.GetRoot(), index, &formA, expression, make(map[string]Expression))

// 	if formAApplies {

// 		a.A = expression.GetChildByPath(root, []int{0, 0})

// 		a.B = expression.GetChildByPath(root, []int{2, 0})

// 		return true
// 	}
// 	formB := ParseExpression("(a^2)+c+(b^2)") // where c is 2*a*b and a and b are constant

// 	formBApplies := IsEqualByFormAt(formB.GetRoot(), index, &formB, expression, make(map[string]Expression))

// 	if formBApplies {

// 		A := expression.GetChildByPath(root, []int{0, 0})

// 		B := expression.GetChildByPath(root, []int{2, 0})

// 		C := expression.GetChildByPath(root, []int{1})

// 		if CheckConstantValue([]int{2, A, B}, C, Multiplication, expression) {

// 			a.A = A

// 			a.B = B

// 			return true
// 		}
// 	}
// 	formC := ParseExpression("(a^2)+(c*a)+(b^2)") // where c is 2*b and b is constant

// 	formCApplies := IsEqualByFormAt(formC.GetRoot(), index, &formC, expression, make(map[string]Expression))

// 	if formCApplies {

// 		B := expression.GetChildByPath(root, []int{2, 0})

// 		C := expression.GetChildByPath(root, []int{1, 0})

// 		if expression.IsConstant(B) && expression.IsConstant(C) {

// 			mul := expression.GetNumericValueByIndex(B) * 2

// 			if expression.GetNumericValueByIndex(C) == mul {

// 				a.A = expression.GetChildByPath(root, []int{0, 0})

// 				a.B = B

// 				return true
// 			}
// 		}
// 	}
// 	formD := ParseExpression("(a^2)+(c*b)+(b^2)") // where c is 2*a and a is constant

// 	formDApplies := IsEqualByFormAt(formD.GetRoot(), index, &formD, expression, make(map[string]Expression))

// 	if formDApplies {

// 		A := expression.GetChildByPath(root, []int{0, 0})

// 		C := expression.GetChildByPath(root, []int{1, 0})

// 		if expression.IsConstant(A) && expression.IsConstant(C) {

// 			mul := expression.GetNumericValueByIndex(A) * 2

// 			if expression.GetNumericValueByIndex(C) == mul {

// 				a.A = A

// 				a.B = expression.GetChildByPath(root, []int{2, 0})

// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

func (a *AlgebraicIdentityA) Apply(index int, expression *Expression) Expression {

	exponentRoot, exponent := NewExpressionWithRoot(Symbol{Exponent, -1, "^"})

	add := exponent.AppendNode(exponentRoot, Symbol{Addition, -1, "+"})

	sumOperands := []int{a.A, a.B}

	exponent.AppendBulkSubtreesFrom(add, sumOperands, *expression)

	exponent.AppendNode(exponentRoot, Symbol{Constant, 2, "2"})

	return exponent

}

func (a *AlgebraicIdentityA) Run(index int, expression *Expression) (bool, Expression) {

	if a.Identify(index, expression) {

		return true, a.Apply(index, expression)

	} else {

		return false, *expression
	}
}

// (a-b)^2=(a^2)-(2*a*b)+(b^2)

// (a^2)-(b^2)=(a+b)*(a-b)

// (x+a)*(x+b) = (x^2)+((a+b)*x)+(a*b)

type AlgebraicIdentityD struct {
	A int

	B int

	X int
}
