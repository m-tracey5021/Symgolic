package symbols

import "strconv"

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

	FunctionParameters = 28

	NaryTuple = 29

	Set = 30

	SetClose = 31

	Vector = 32

	VectorClose = 33

	Union = 34

	Intersection = 35

	Subset = 36

	ProperSubset = 37

	Assignment = 38

	NewLine = 39

	EndOfFile = 40

	None = 41
)

type Symbol struct {
	SymbolType SymbolType

	NumericValue int

	AlphaValue string
}

func NewOperation(symbolType SymbolType) Symbol {

	var alpha string

	switch symbolType {

	case Addition:

		alpha = "+"

	case Subtraction:

		alpha = "-"

	case Multiplication:

		alpha = "*"

	case Division:

		alpha = "/"

	case Exponent:

		alpha = "^"

	case Radical:

		alpha = "v"

	case NaryTuple:

		alpha = "(...)"

	case Set:

		alpha = "{...}"

	case Vector:

		alpha = "[...]"
	}
	return Symbol{SymbolType: symbolType, NumericValue: -1, AlphaValue: alpha}
}

func NewVariable(variable string) Symbol {

	return Symbol{SymbolType: Variable, NumericValue: -1, AlphaValue: variable}
}

func NewConstant(value int) Symbol {

	return Symbol{SymbolType: Constant, NumericValue: value, AlphaValue: strconv.Itoa(value)}
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
		s.SymbolType == Subtraction ||
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
		s.SymbolType == Assignment ||
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

	return Symbol{SymbolType: s.SymbolType, NumericValue: s.NumericValue, AlphaValue: s.AlphaValue}
}
