package symbols

type Program struct {
	Expressions []Expression
}

func NewProgram() Program {

	Expressions := make([]Expression, 0)

	return Program{Expressions}
}

func (p *Program) AddExpression(expression Expression) {

	p.Expressions = append(p.Expressions, expression)
}
