package symbols

type SymbolType int

const (
	Equality = iota

	GreaterThan

	LessThan

	GreaterThanOrEqualTo

	LessThanOrEqualTo

	Open

	Close

	Addition

	Subtraction

	Multiplication

	Division

	Exponent

	Radical

	Variable

	Constant

	And

	Or

	If

	Iff

	Negation

	Necessity

	Possibility

	Universal

	Existential

	None
)

type Symbol struct {
	SymbolType SymbolType

	NumericValue int

	CharacterValue string
}

func (s *Symbol) Copy() Symbol {

	return Symbol{s.SymbolType, s.NumericValue, s.CharacterValue}
}
