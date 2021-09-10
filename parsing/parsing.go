package parsing

import (
	"errors"
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

	// var tmpExpression Expression = symbols.NewExpression()

	var parser Parser = Parser{expression, lex(text, parseType), 0}

	parser.topLevelExpression()

	return parser.parsed, nil
}

func (p *Parser) topLevelExpression() {

	subExpressions := make([]int, 0)

	mainAux := p.auxillary()

	first = p.subExpression(mainAux)

	continuation := p.operator(true) // optional from here on

	rhs = append(rhs, p.right(rhs)...)

	p.parsed.SetRootByIndex(p.subExpression(mainAux))

}

func (p *Parser) atom() int {

	sign := true

	if p.tokens[p.currentToken].SymbolType == Subtraction {

		sign = false

		p.currentToken++
	}

	if p.tokens[p.currentToken].SymbolType == Variable || p.tokens[p.currentToken].SymbolType == Constant {

		child := p.addNode(sign)

		p.currentToken++

		return child
	}
	return -1
}

func (p *Parser) subExpression(sign bool) int {

	lhs := make([]int, 0)

	rhs := make([]int, 0)

	// mainAux := p.auxillary()

	p.open()

	lhs = append(lhs, p.left())

	parent := p.operator(sign)

	rhs = append(rhs, p.right(rhs)...)

	p.close(parent, append(lhs, rhs...))

	return parent
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

func (p *Parser) left() int {

	sign := p.auxillary()

	// if p.tokens[p.currentToken].SymbolType == Variable || p.tokens[p.currentToken].SymbolType == Constant {

	// 	child := p.addNode(sign)

	// 	p.currentToken++

	// 	return child

	// } else {

	// 	child := p.subExpression(sign)

	// 	return child
	// }

	atom := p.atom(sign)

	if atom != -1 {

		return atom

	} else {

		return p.subExpression(sign)
	}

}

func (p *Parser) operator(sign bool) int {

	if p.tokens[p.currentToken].SymbolType == Addition ||
		p.tokens[p.currentToken].SymbolType == Multiplication ||
		p.tokens[p.currentToken].SymbolType == Division ||
		p.tokens[p.currentToken].SymbolType == Exponent ||
		p.tokens[p.currentToken].SymbolType == Radical {

		// append operator to tree
		parent := p.addNode(sign)

		p.currentToken++

		return parent
	}
	return -1
}

// func (p *Parser) right(sign bool) int {

// 	if p.tokens[p.currentToken].SymbolType == Variable || p.tokens[p.currentToken].SymbolType == Constant {

// 		child := p.addNode(sign)

// 		p.currentToken++

// 		return child

// 	} else {

// 		child := p.expression(sign)

// 		return child
// 	}
// }

func (p *Parser) right(children []int) []int {

	sign := p.auxillary()

	child := p.atom(sign)

	if child == -1 {

		child = p.subExpression(sign)
	}

	operator := p.operator(sign)

	if p.tokens[p.currentToken].SymbolType == Variable || p.tokens[p.currentToken].SymbolType == Constant {

		children = append(children, p.addNode(sign))

		p.currentToken++

		children = p.right(children)

	} else if p.tokens[p.currentToken].SymbolType == Addition ||
		p.tokens[p.currentToken].SymbolType == Multiplication ||
		p.tokens[p.currentToken].SymbolType == Division ||
		p.tokens[p.currentToken].SymbolType == Exponent ||
		p.tokens[p.currentToken].SymbolType == Radical {

		p.currentToken++

		children = p.right(children)

	} else if p.tokens[p.currentToken].SymbolType == Open {

		children = append(children, p.subExpression(sign))

		children = p.right(children)

	} else if p.tokens[p.currentToken].SymbolType == Close || p.tokens[p.currentToken].SymbolType == None {

		return children

	} else {

		panic(errors.New("unrecognised token"))
	}
	return children
}

func (p *Parser) close(parent int, children []int) {

	if p.tokens[p.currentToken].SymbolType == Close || p.tokens[p.currentToken].SymbolType == None {

		for _, child := range children {

			p.parsed.SetParent(parent, child)
		}
	}
}

func (p *Parser) addNode(sign bool) int {

	return p.parsed.AddToMap(p.tokens[p.currentToken], sign)
}

// func (p *Parser) addChild() int {

// 	return p.parsed.AddToMap(p.tokens[p.currentToken], p.currentSign)
// }

// func (p *Parser) addChildren() {

// 	for _, child := range p.children {

// 		p.parsed.AppendNode(p.parent, child, true)
// 	}
// }

func (p *Parser) complete() {

}
