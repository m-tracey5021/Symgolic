package identities

import (
	. "symgolic/comparison"
	. "symgolic/parsing"
	. "symgolic/symbols"
)

// (a+b)^2 = (a^2)+(2*a*b)+(b^2) || (2^2)+12+(3^2) || (a^2)+(6*a)+(3^2) || (2^2)+(4*b)+(b^2)

type AlgebraicIdentityA struct {
	Result bool

	Stage int

	Structure map[int]Stage

	A int

	B int

	Coeff int
}

func NewAlgebraicIdentityA() AlgebraicIdentityA {

	structure := map[int]Stage{

		0: Stage{SymbolType: Addition, Recurse: true},

		1: Stage{SymbolType: Exponent, Recurse: true},

		2: Stage{SymbolType: Variable, Recurse: false},

		3: Stage{SymbolType: Constant, Recurse: false},

		4: Stage{SymbolType: Multiplication, Recurse: true},

		5: Stage{SymbolType: Constant, Recurse: false},

		6: Stage{SymbolType: Variable, Recurse: false},

		7: Stage{SymbolType: Variable, Recurse: false},

		8: Stage{SymbolType: Exponent, Recurse: true},

		9: Stage{SymbolType: Variable, Recurse: false},

		10: Stage{SymbolType: Constant, Recurse: false},
	}
	return AlgebraicIdentityA{Result: false, Stage: 0, Structure: structure}
}

func (a *AlgebraicIdentityA) Pass(node Symbol) bool {

	expected := a.Structure[a.Stage].SymbolType

	if expected == Variable {

		if node.SymbolType != Constant {

			return true

		} else {

			return false
		}

	} else {

		return node.SymbolType == expected
	}
}

func (a *AlgebraicIdentityA) Failure() {

	a.Result = false
}

// func (a *AlgebraicIdentityA) Identify(index int, expression *Expression) bool {

// 	node := expression.GetNodeByIndex(index)

// 	pass := a.Pass(*node)

// 	if pass {

// 		if a.Stage == 2 {

// 			a.A = index

// 		} else if a.Stage == 5 {

// 			// check that coeff is two times others

// 			if node.NumericValue != -1 {

// 				a.Coeff = node.NumericValue

// 			} else {

// 				return false
// 			}

// 		} else if a.Stage == 6 {

// 			if !IsEqualAt(a.A, index, expression, expression) {

// 				return false
// 			}

// 		} else if a.Stage == 7 {

// 			a.B = index

// 		} else if a.Stage == 9 {

// 			if !IsEqualAt(a.B, index, expression, expression) {

// 				return false
// 			}
// 		}
// 		recurse := a.Structure[a.Stage].Recurse

// 		a.Stage++

// 		if recurse {

// 			for _, child := range expression.GetChildren(index) {

// 				result := a.Identify(child, expression)

// 				if !result {

// 					return false
// 				}
// 			}
// 		}
// 		return true

// 	} else {

// 		return false
// 	}
// }

func (a *AlgebraicIdentityA) Identify(index int, expression *Expression) bool {

	root := expression.GetRoot()

	formA := ParseExpression("(a^2)+(2*a*b)+(b^2)") // standard form

	formAApplies := IsEqualByFormAt(formA.GetRoot(), index, &formA, expression, make(map[string]Expression))

	if formAApplies {

		a.A = expression.GetChildByPath(root, []int{0, 0})

		a.B = expression.GetChildByPath(root, []int{2, 0})

		return true
	}
	formB := ParseExpression("(a^2)+c+(b^2)") // where c is 2*a*b and a and b are constant

	formBApplies := IsEqualByFormAt(formB.GetRoot(), index, &formB, expression, make(map[string]Expression))

	if formBApplies {

		A := expression.GetChildByPath(root, []int{0, 0})

		B := expression.GetChildByPath(root, []int{2, 0})

		C := expression.GetChildByPath(root, []int{1})

		if expression.IsConstant(A) && expression.IsConstant(B) && expression.IsConstant(C) {

			mul := expression.GetNumericValueByIndex(A) * expression.GetNumericValueByIndex(B) * 2

			if expression.GetNumericValueByIndex(C) == mul {

				a.A = A

				a.B = B

				return true
			}
		}
	}
	formC := ParseExpression("(a^2)+(c*a)+(b^2)") // where c is 2*b and b is constant

	formCApplies := IsEqualByFormAt(formC.GetRoot(), index, &formC, expression, make(map[string]Expression))

	if formCApplies {

		B := expression.GetChildByPath(root, []int{2, 0})

		C := expression.GetChildByPath(root, []int{1, 0})

		if expression.IsConstant(B) && expression.IsConstant(C) {

			mul := expression.GetNumericValueByIndex(B) * 2

			if expression.GetNumericValueByIndex(C) == mul {

				a.A = expression.GetChildByPath(root, []int{0, 0})

				a.B = B

				return true
			}
		}
	}
	formD := ParseExpression("(a^2)+(c*b)+(b^2)") // where c is 2*a and a is constant

	formDApplies := IsEqualByFormAt(formD.GetRoot(), index, &formD, expression, make(map[string]Expression))

	if formDApplies {

		A := expression.GetChildByPath(root, []int{0, 0})

		C := expression.GetChildByPath(root, []int{1, 0})

		if expression.IsConstant(A) && expression.IsConstant(C) {

			mul := expression.GetNumericValueByIndex(A) * 2

			if expression.GetNumericValueByIndex(C) == mul {

				a.A = A

				a.B = expression.GetChildByPath(root, []int{2, 0})

				return true
			}
		}
	}
	return false
}

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
