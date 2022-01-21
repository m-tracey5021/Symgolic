package components

// Config

type Pairing struct {
	First, Second SymbolType
}

type Operation struct {
	Call func(Expression, Expression, bool) Expression

	Reverse bool
}

var sumPairings = map[Pairing]Operation{

	Pairing{First: Addition, Second: Addition}:             Operation{Call: SSS, Reverse: false},
	Pairing{First: Addition, Second: Multiplication}:       Operation{Call: SSM, Reverse: false},
	Pairing{First: Addition, Second: Division}:             Operation{Call: SSD, Reverse: false},
	Pairing{First: Multiplication, Second: Addition}:       Operation{Call: SSM, Reverse: true},
	Pairing{First: Multiplication, Second: Multiplication}: Operation{Call: MSM, Reverse: false},
	Pairing{First: Multiplication, Second: Division}:       Operation{Call: MSD, Reverse: false},
	Pairing{First: Division, Second: Addition}:             Operation{Call: SSD, Reverse: true},
	Pairing{First: Division, Second: Multiplication}:       Operation{Call: MSD, Reverse: true},
	Pairing{First: Division, Second: Division}:             Operation{Call: DSD, Reverse: false},
}

var mulPairings = map[Pairing]Operation{

	Pairing{First: Addition, Second: Multiplication}: Operation{Call: SSM, Reverse: false},
	Pairing{First: Addition, Second: Division}:       SSD,
	Pairing{First: Division, Second: Multiplication}: MulDiv,
}

var divPairings = map[Pairing]Operation{

	Pairing{First: Addition, Second: Multiplication}: Operation{Call: SSM, Reverse: false},
	Pairing{First: Addition, Second: Division}:       SSD,
	Pairing{First: Division, Second: Multiplication}: MulDiv,
}

func Compare(a, b, operation SymbolType) Operation {

	var pairing map[Pairing]Operation

	switch operation {

	case Addition:

		pairing = sumPairings

	case Multiplication:

		pairing = mulPairings

	case Division:

		pairing = divPairings
	}
	match, exists := pairing[Pairing{First: a, Second: b}]

	if exists {

		return match
	}
	panic("No matching function")
}

// Implementations

func SSS(a, b Expression, reverse bool) Expression { // adding two sums

	if reverse {

		return SSS(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.GetChildren(a.GetRoot()), a)

	result.AppendBulkSubtreesFrom(root, b.GetChildren(b.GetRoot()), b)

	return result
}

func SSM(a, b Expression, reverse bool) Expression { // adding sum and multiplication

	if reverse {

		return SSM(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.GetChildren(a.GetRoot()), a)

	result.AppendSubtreeFrom(root, b.GetRoot(), b)

	return result
}

func SSD(a, b Expression, reverse bool) Expression { // adding sum and division

	if reverse {

		return SSD(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.GetChildren(a.GetRoot()), a)

	result.AppendSubtreeFrom(root, b.GetRoot(), b)

	return result

}

func MSM(a, b Expression, reverse bool) Expression { // adding multiplication and multiplication

	if reverse {

		return MSM(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	coeffA, termA := SeparateTerm(a.GetRoot(), a)

	coeffB, termB := SeparateTerm(b.GetRoot(), b)

	if termB == termA {

		resultCoeff := coeffA + coeffB

		return CreateLikeTerm(resultCoeff, termA)

	} else {

		result.AppendSubtreeFrom(root, a.GetRoot(), a)

		result.AppendSubtreeFrom(root, b.GetRoot(), b)

		return result
	}
}

func MSD(a, b Expression, reverse bool) Expression { // adding multiplication and division

	if reverse {

		return MSD(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.GetChildren(a.GetRoot()), a)

	result.AppendSubtreeFrom(root, b.GetRoot(), b)

	return result
}

func DSD(a, b Expression, reverse bool) Expression { // adding division and division

}

func VSV(a, b Expression) Expression { // adding vector and vector

}

// Arithmetic base

func (e *Expression) Negate() {

	root := e.GetRoot()

	negation := make([]Symbol, 0)

	negation = append(negation, NewOperation(Subtraction))

	e.InsertAuxiliariesAt(root, negation)
}

func (e *Expression) Add(others ...Expression) Expression {

	cumulative := *e

	for _, operand := range others {

		operation := Compare(cumulative.GetNode(cumulative.GetRoot()).SymbolType, operand.GetNode(e.GetRoot()).SymbolType, Addition)

		cumulative = operation.Call(*e, operand)
	}
	return cumulative
}

func (e *Expression) Subtract(other Expression) Expression {

	return Compare(e.GetNode(e.GetRoot()).SymbolType, other.GetNode(e.GetRoot()).SymbolType, Addition)(*e, other)
}

func (e *Expression) Multiply(others ...Expression) Expression {

	cumulative := *e

	for _, operand := range others {

		operation := Compare(cumulative.GetNode(cumulative.GetRoot()).SymbolType, operand.GetNode(e.GetRoot()).SymbolType, Multiplication)

		cumulative = operation(*e, operand)
	}
	return cumulative
}

func (e *Expression) Divide(other Expression) Expression {

	return Compare(e.GetNode(e.GetRoot()).SymbolType, other.GetNode(e.GetRoot()).SymbolType, Division)(*e, other)
}

// Util

func SeparateTerm(index int, expression Expression) (int, string) {

	coeff := 1

	term := ""

	for _, child := range expression.GetChildren(index) {

		node := expression.GetNode(child)

		if node.SymbolType == Variable {

			term += node.AlphaValue

		} else if node.SymbolType == Constant {

			coeff = node.NumericValue
		}
	}
	return coeff, term
}

func CreateLikeTerm(coeff int, term string) Expression {

	_, coeffExp := NewExpression(NewConstant(coeff))

	_, termExp := NewExpression(NewVariable(term))

	if coeff == 0 {

		return coeffExp

	} else if coeff == 1 {

		return termExp

	} else {

		likeRoot, likeTerm := NewExpression(NewOperation(Multiplication))

		likeTerm.AppendExpression(likeRoot, coeffExp, false)

		likeTerm.AppendExpression(likeRoot, termExp, false)

		return likeTerm
	}
}
