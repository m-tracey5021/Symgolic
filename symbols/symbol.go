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

	Function

	Set

	SetClose

	SetElement

	Union

	Intersection

	Subset

	ProperSubset

	None
)

type Symbol struct {
	SymbolType SymbolType

	NumericValue int

	CharacterValue string
}

func (s *Symbol) IsAuxiliary() bool {

	if s.SymbolType == Subtraction ||
		s.SymbolType == Negation ||
		s.SymbolType == Necessity ||
		s.SymbolType == Possibility {

		return true

	} else {

		return false
	}
}

func (s *Symbol) IsOperation() bool {
	if s.SymbolType == Addition ||
		s.SymbolType == Multiplication ||
		s.SymbolType == Division ||
		s.SymbolType == Exponent ||
		s.SymbolType == Radical ||
		s.SymbolType == Union ||
		s.SymbolType == Intersection {

		return true

	} else {

		return false
	}
}

func (s *Symbol) IsComparison() bool {

	if s.SymbolType == Equality ||
		s.SymbolType == GreaterThan ||
		s.SymbolType == LessThan ||
		s.SymbolType == GreaterThanOrEqualTo ||
		s.SymbolType == LessThanOrEqualTo {

		return true

	} else {

		return false
	}
}

func (s *Symbol) Copy() Symbol {

	return Symbol{s.SymbolType, s.NumericValue, s.CharacterValue}
}
