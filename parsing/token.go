package parsing

type TokenType int

const (
	Addition = iota

	Negation

	Multipliaction

	Division

	Exponent

	Radical

	Variable

	Constant

	None
)

type Token struct {
	tokenType int

	value string
}
