package parsing

import (
	. "symgolic/symbols"
	"unicode"
)

type ParseType int

const (
	Math = iota

	Logic

	NaturalLanguage
)

func getOperatorSymbolType(character string, parseType int) int {

	if character == "=" {

		return Equality

	} else if character == ">" {

		return GreaterThan

	} else if character == "<" {

		return LessThan

	} else if character == ">=" {

		return GreaterThanOrEqualTo

	} else if character == "<=" {

		return LessThanOrEqualTo

	} else if character == "(" {

		return Open

	} else if character == ")" {

		return Close

	} else if character == "+" {

		return Addition

	} else if character == "-" {

		return Subtraction

	} else if character == "*" {

		return Multiplication

	} else if character == "/" {

		return Division

	} else if character == "^" {

		if parseType == Math {

			return Exponent

		} else {

			return And
		}
	} else if character == "v" {

		if parseType == Math {

			return Radical

		} else {

			return Or
		}
	} else if character == "->" {

		return If

	} else if character == "<>" {

		return Iff

	} else if character == "~" {

		return Negation

	} else if character == "!" {

		return Necessity

	} else if character == "?" {

		return Possibility

	} else if character == "A" {

		return Universal

	} else if character == "E" {

		return Existential

	} else {

		return None

	}
}

func lex(text string, parseType int) []Symbol {

	var tokens map[string]SymbolType = map[string]SymbolType{

		"=": Equality,

		">": GreaterThan,

		"<": LessThan,

		">=": GreaterThanOrEqualTo,

		"<=": LessThanOrEqualTo,

		"(": Open,

		")": Close,

		"+": Addition,

		"-": Subtraction,

		"*": Multiplication,

		"/": Division,

		"^": Exponent,

		"v": Radical,

		"&": And,

		"|": Or,

		"->": If,

		"<>": Iff,

		"~": Negation,

		"!": Necessity,

		"?": Possibility,

		"A": Universal,

		"E": Existential,
	}
	var symbols []Symbol

	var characters []rune = []rune(text)

	for i := 0; i < len(characters); i++ {

		characterAt := characters[i]

		if unicode.IsLetter(characterAt) && characterAt != 'v' && characterAt != 'A' && characterAt != 'E' {

			symbols = append(symbols, Symbol{Variable, -1, text[i : i+1]})

		} else if unicode.IsDigit(characterAt) {

			if i+1 < len(text) {

				var j int = i + 1

				for unicode.IsDigit(characters[j]) {

					j++
				}
				symbols = append(symbols, Symbol{Constant, 0, text[i:j]})

			} else {

				symbols = append(symbols, Symbol{Constant, 0, text[i : i+1]})
			}

		} else {

			simpleToken := text[i : i+1]

			simple, simpleExists := tokens[simpleToken]

			if simpleExists {

				if i+2 <= len(characters) {

					compoundToken := text[i : i+2]

					compound, compoundExists := tokens[compoundToken]

					if compoundExists {

						symbols = append(symbols, Symbol{compound, -1, compoundToken})

						i++

					} else {

						symbols = append(symbols, Symbol{simple, -1, simpleToken})
					}
				} else {

					symbols = append(symbols, Symbol{simple, -1, simpleToken})
				}
			} else {

				continue
			}
		}
		if i == len(characters)-1 {

			symbols = append(symbols, Symbol{None, -1, ""})
		}
	}
	return symbols
}

// after +, -,

type Parser struct {
	parsed Expression

	tokens []Symbol

	currentToken int
}

func ParseExpression(text string, parseType int) (Expression, error) {

	var expression Expression = NewExpression()

	var parser Parser = Parser{expression, lex(text, parseType), 0}

	parser.expression()

	return parser.parsed, nil
}

func (p *Parser) auxillary() bool {

	if p.tokens[p.currentToken].SymbolType == Subtraction {

		p.currentToken++

		return false

	} else {

		return true
	}
}

func (p *Parser) open() {

	if p.tokens[p.currentToken].SymbolType == Open {

		p.currentToken++
	}
}

func (p *Parser) atom(sign bool) int {

	if p.tokens[p.currentToken].SymbolType == Variable || p.tokens[p.currentToken].SymbolType == Constant {

		child := p.addNode(sign)

		p.currentToken++

		return child
	}
	return -1
}

func (p *Parser) operand() int {

	sign := p.auxillary()

	child := p.atom(sign)

	if child == -1 {

		child = p.subExpression(sign)

		return child

	} else {

		return child
	}
}

func (p *Parser) operator(sign bool) int {

	if p.tokens[p.currentToken].IsOperation() {

		parent := p.addNode(sign)

		p.currentToken++

		return parent
	}
	return -1
}

func (p *Parser) operandChain(children []int) []int {

	if p.close() {

		return children

	} else if p.tokens[p.currentToken].IsOperation() {

		p.currentToken++

		children = p.operandChain(children)

		return children

	} else {

		children = append(children, p.operand())

		children = p.operandChain(children)

		return children
	}
}

func (p *Parser) operands(sign bool, parent int, children []int) (int, []int) {

	if p.close() { // used directly after first operand, i.e. the lhs, so it atomic will close

		return parent, children

	} else {

		if parent != -1 {

			parent = p.operator(sign)
		}

		children = append(children, p.operand())

		parent, children = p.operands(sign, parent, children)

		return parent, children
	}
}

func (p *Parser) close() bool {

	if p.tokens[p.currentToken].SymbolType == Close {

		p.currentToken++

		return true

	} else if p.tokens[p.currentToken].SymbolType == None {

		return true

	} else {

		return false
	}
}

func (p *Parser) subExpression(expressionSign bool) int {

	lhs := make([]int, 0)

	rhs := make([]int, 0)

	p.open()

	lhsChild := p.operand()

	lhs = append(lhs, lhsChild)

	parent := p.operator(expressionSign)

	rhsChild := p.operand()

	rhs = append(rhs, rhsChild)

	rhs = p.operandChain(rhs)

	p.complete(parent, append(lhs, rhs...)) // return children[0] if parent is -1

	return parent
}

func (p *Parser) expression() {

	sign := p.auxillary()

	subExpressions := make([]int, 0)

	first := p.subExpression(sign)

	if p.close() {

		p.parsed.SetRootByIndex(first)

	} else {

		chain := p.operator(true)

		subExpressions = append(subExpressions, first)

		subExpressions = p.operandChain(subExpressions)

		p.complete(chain, subExpressions)

		p.parsed.SetRootByIndex(chain)
	}
}

func (p *Parser) addNode(sign bool) int {

	return p.parsed.AddToMap(p.tokens[p.currentToken], sign)
}

func (p *Parser) complete(parent int, children []int) {

	for _, child := range children {

		p.parsed.SetParent(parent, child)
	}
}
