package components

// Config

type Pairing struct {
	First, Second SymbolType
}

type Operation struct {
	Call func(Expression, Expression, bool) Expression

	Reverse bool
}

// Comparer

func Compare(a, b, operation SymbolType) Operation {

	var pairing map[Pairing]Operation

	switch operation {

	case Addition:

		pairing = map[Pairing]Operation{

			{First: Addition, Second: Addition}:             {Call: SSS, Reverse: false},
			{First: Addition, Second: Multiplication}:       {Call: SSM, Reverse: false},
			{First: Addition, Second: Division}:             {Call: SSD, Reverse: false},
			{First: Multiplication, Second: Addition}:       {Call: SSM, Reverse: true},
			{First: Multiplication, Second: Multiplication}: {Call: MSM, Reverse: false},
			{First: Multiplication, Second: Division}:       {Call: MSD, Reverse: false},
			{First: Division, Second: Addition}:             {Call: SSD, Reverse: true},
			{First: Division, Second: Multiplication}:       {Call: MSD, Reverse: true},
			{First: Division, Second: Division}:             {Call: DSD, Reverse: false},
		}

	case Multiplication:

		pairing = map[Pairing]Operation{}

	case Division:

		pairing = map[Pairing]Operation{}
	}
	match, exists := pairing[Pairing{First: a, Second: b}]

	if exists {

		return match
	}
	panic("No matching function")
}

// Arithmetic base

func Negate(target ExpressionIndex) {

	negation := make([]Symbol, 0)

	negation = append(negation, NewOperation(Subtraction))

	target.Expression.InsertAuxiliariesAt(target.Index, negation)
}

func Add(operands ...ExpressionIndex) Expression {

	cumulative := operands[0].Expression

	cumulativeIndex := operands[0].Index

	for _, operand := range operands {

		operation := Compare(cumulative.GetNode(cumulativeIndex).SymbolType, operand.Expression.GetNode(operand.Index).SymbolType, Addition)

		cumulative = operation.Call(cumulative, operand.Expression, operation.Reverse)

		cumulativeIndex = cumulative.GetRoot()
	}
	return cumulative
}

func Subtract(a, b ExpressionIndex) Expression {

	operation := Compare(a.Expression.GetNode(a.Index).SymbolType, b.Expression.GetNode(b.Index).SymbolType, Addition)

	Negate(ExpressionIndex{Expression: b.Expression, Index: b.Expression.GetRoot()})

	return operation.Call(a.Expression, b.Expression, operation.Reverse)
}

func Multiply(operands ...ExpressionIndex) Expression {

	cumulative := operands[0].Expression

	cumulativeIndex := operands[0].Index

	for _, operand := range operands {

		operation := Compare(cumulative.GetNode(cumulativeIndex).SymbolType, operand.Expression.GetNode(operand.Index).SymbolType, Multiplication)

		cumulative = operation.Call(cumulative, operand.Expression, operation.Reverse)

		cumulativeIndex = cumulative.GetRoot()
	}
	return cumulative
}

func Divide(a, b ExpressionIndex) Expression {

	operation := Compare(a.Expression.GetNode(a.Index).SymbolType, b.Expression.GetNode(b.Index).SymbolType, Division)

	return operation.Call(a.Expression, b.Expression, operation.Reverse)
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
