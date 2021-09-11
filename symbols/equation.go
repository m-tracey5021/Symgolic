package symbols

type Equation struct {
	lhs Expression

	relation SymbolType

	rhs Expression
}

func NewEquation() Equation {

	var equation Equation = Equation{}

	equation.lhs = NewExpression()

	equation.rhs = NewExpression()

	return equation
}
