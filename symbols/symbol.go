package symbols

type SymbolType int

const (
	Equality = iota

	GreaterThan = 1

	LessThan = 2

	GreaterThanOrEqualTo = 3

	LessThanOrEqualTo = 4

	ExpressionOpen = 5

	ExpressionClose = 6

	SubExpressionOpen = 7

	SubExpressionClose = 8

	Addition = 9

	Subtraction = 10

	Multiplication = 11

	Division = 12

	Exponent = 13

	Radical = 14

	Iteration = 15

	Variable = 16

	Constant = 17

	And = 18

	Or = 19

	If = 20

	Iff = 21

	Negation = 22

	Necessity = 23

	Possibility = 24

	Universal = 25

	Existential = 26

	Function = 27

	NaryTuple = 28

	Set = 29

	SetClose = 30

	Vector = 31

	VectorClose = 32

	Union = 33

	Intersection = 34

	Subset = 35

	ProperSubset = 36

	NewLine = 37

	EndOfFile = 38

	None = 39
)

type Symbol struct {
	SymbolType SymbolType

	NumericValue int

	AlphaValue string
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
		s.SymbolType == Iteration ||
		s.SymbolType == Union ||
		s.SymbolType == Intersection {

		return true

	} else {

		return false
	}
}

func (s *Symbol) IsEnclosingOperation() bool {
	if s.SymbolType == ExpressionOpen ||
		s.SymbolType == SubExpressionOpen ||
		s.SymbolType == Set ||
		s.SymbolType == Vector {

		return true

	} else {

		return false
	}
}

func (s *Symbol) ClosesExpressionScope() bool {
	if s.SymbolType == ExpressionClose ||
		s.SymbolType == SubExpressionClose ||
		s.SymbolType == SetClose ||
		s.SymbolType == VectorClose ||
		// s.SymbolType == ParameterClose ||
		s.IsComparison() {

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

	return Symbol{s.SymbolType, s.NumericValue, s.AlphaValue}
}
