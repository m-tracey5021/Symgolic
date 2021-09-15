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

func lex(text string, parseType int) []Symbol {

	var symbols []Symbol

	var characters []rune = []rune(text)

	for i := 0; i < len(characters); i++ {

		characterAt := characters[i]

		if i == 0 {

			symbols = append(symbols, Symbol{Open, -1, "expression start"})
		}
		symbolType, val, symbol, end := lexOperand(text, characters, i, true) // gets name, variable or constant

		if symbolType != None {

			symbols = append(symbols, Symbol{symbolType, val, symbol})

			i = end - 1

		} else { // check for operators, opens and closes

			if characterAt == ' ' {

				continue

			} else {

				operatorType, val, operator, end := lexOperator(text, characters, i)

				symbol := Symbol{operatorType, val, operator}

				if symbol.SymbolType == SetElement {

					symbols = append(symbols, Symbol{Close, -1, "parameter end"})

					symbols = append(symbols, Symbol{Open, -1, "parameter start"})

				} else if symbol.SymbolType == Set {

					symbols = append(symbols, symbol)

					symbols = append(symbols, Symbol{Open, -1, "parameter started"})

				} else if symbol.SymbolType == SetClose {

					symbols = append(symbols, Symbol{Close, -1, "parameter end"})

					symbols = append(symbols, symbol)

				} else if symbol.IsComparison() {

					symbols = append(symbols, Symbol{Close, -1, "expression end"})

					symbols = append(symbols, symbol)

					symbols = append(symbols, Symbol{Open, -1, "expression start"})

				} else {

					symbols = append(symbols, symbol)
				}
				i = end
			}
		}
		if i == len(characters)-1 {

			symbols = append(symbols, Symbol{Close, -1, "expression end"})

			symbols = append(symbols, Symbol{None, -1, "EOF"})
		}
	}
	return symbols
}

func lexOperator(text string, characters []rune, index int) (SymbolType, int, string, int) {

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

		"{": Set,

		"}": SetClose,

		",": SetElement,

		"u": Union,

		"n": Intersection,

		"c": Subset,

		"c=": ProperSubset,
	}

	// var tokens map[rune]SymbolType = map[rune]SymbolType{

	// 	'=': Equality,

	// 	'>': GreaterThan,

	// 	'<': LessThan,

	// 	'\u2265': GreaterThanOrEqualTo,

	// 	'\u2264': LessThanOrEqualTo,

	// 	'(': Open,

	// 	')': Close,

	// 	'+': Addition,

	// 	'-': Subtraction,

	// 	'*': Multiplication,

	// 	'/': Division,

	// 	'^': Exponent,

	// 	'v': Radical,

	// 	'\u2227': And,

	// 	'\u2228': Or,

	// 	'\u2192': If,

	// 	'\u2261': Iff,

	// 	'~': Negation,

	// 	'!': Necessity,

	// 	'?': Possibility,

	// 	'\u2200': Universal,

	// 	'\u2203': Existential,

	// 	'{': Set,

	// 	'}': SetClose,

	// 	',': SetElement,

	// 	'\u222a': Union,

	// 	'\u2229': Intersection,
	// }

	simpleToken := text[index : index+1]

	simple, simpleExists := tokens[simpleToken]

	if index+2 <= len(characters) {

		compoundToken := text[index : index+2]

		compound, compoundExists := tokens[compoundToken]

		if compoundExists {

			return compound, -1, compoundToken, index + 1

		} else {

			if simpleExists {

				return simple, -1, simpleToken, index

			} else {

				return None, -1, "", index
			}
		}

	} else {

		if simpleExists {

			return simple, -1, simpleToken, index

		} else {

			return None, -1, "", index
		}
	}
}

func lexOperand(text string, characters []rune, index int, predefined bool) (SymbolType, int, string, int) {

	symbolType, value, symbol, end := lexWord(text, characters, index, predefined)

	if symbolType == None {

		return lexNumber(text, characters, index)

	} else {

		return symbolType, value, symbol, end
	}
}

func lexWord(text string, characters []rune, index int, predefined bool) (SymbolType, int, string, int) {

	end := index

	for i := index; i < len(characters); i++ {

		if unicode.IsLetter(characters[end]) &&
			characters[end] != 'v' &&
			characters[end] != 'A' &&
			characters[end] != 'E' &&
			characters[end] != 'u' &&
			characters[end] != 'n' &&
			characters[end] != 'c' {

			end++

		} else {

			break
		}
	}
	if end == index {

		return None, -1, "", end

	} else {

		if end < len(characters) {

			if characters[end] == '{' {

				return Function, -1, text[index:end], end

			} else {

				return Variable, -1, text[index:end], end
			}
		} else {

			return Variable, -1, text[index:end], end
		}
	}

}

func lexNumber(text string, characters []rune, index int) (SymbolType, int, string, int) {

	end := index

	for i := index; i < len(characters); i++ {

		if unicode.IsDigit(characters[end]) {

			end++

		} else {

			break
		}
	}
	if end == index {

		return None, -1, "", end

	} else {

		return Constant, 0, text[index:end], end
	}
}

type Parser struct {
	parsed Expression

	tokens []Symbol

	currentToken int
}

func NewParser(text string, parseType int) Parser {

	parser := Parser{NewExpression(), lex(text, parseType), 0}

	return parser
}

func ParseExpression(text string, parseType int) Expression {

	parser := NewParser(text, parseType)

	parser.equation()

	return parser.parsed
}

func (p *Parser) auxiliary(auxiliaries []Symbol) []Symbol {

	if p.tokens[p.currentToken].IsAuxiliary() {

		auxiliaries = append(auxiliaries, p.tokens[p.currentToken])

		p.currentToken++

		auxiliaries = p.auxiliary(auxiliaries)

		return auxiliaries

	} else {

		return auxiliaries
	}
}

func (p *Parser) open() bool {

	if p.tokens[p.currentToken].SymbolType == Open || p.tokens[p.currentToken].SymbolType == Set {

		p.currentToken++

		return true

	} else {

		return false
	}
}

func (p *Parser) atom() int {

	if p.tokens[p.currentToken].SymbolType == Variable || p.tokens[p.currentToken].SymbolType == Constant {

		child := p.addNode()

		p.currentToken++

		return child
	}
	return -1
}

func (p *Parser) operand() int {

	auxiliaries := p.auxiliary(make([]Symbol, 0))

	// these types all need to be mutually exclusive

	atom := p.atom()

	if atom != -1 {

		p.addAuxiliaries(atom, auxiliaries)

		return atom
	}
	set := p.set()

	if set != -1 {

		p.addAuxiliaries(set, auxiliaries)

		return set
	}
	function := p.function()

	if function != -1 {

		p.addAuxiliaries(function, auxiliaries)

		return function
	}
	subExpression := p.expression()

	if subExpression != -1 {

		p.addAuxiliaries(subExpression, auxiliaries)

		return subExpression
	}
	return -1

}

func (p *Parser) operands(parent int, children []int) (int, []int) {

	if p.close() { // used directly after first operand, i.e. the lhs, so if atomic will close

		return parent, children

	} else {

		if p.tokens[p.currentToken].IsOperation() {

			if parent == -1 {

				parent = p.addNode()
			}
			p.currentToken++
		}

		children = append(children, p.operand())

		parent, children = p.operands(parent, children)

		return parent, children
	}
}

func (p *Parser) close() bool {

	if p.tokens[p.currentToken].SymbolType == Close || p.tokens[p.currentToken].SymbolType == SetClose || p.tokens[p.currentToken].IsComparison() {

		p.currentToken++

		return true

	} else if p.tokens[p.currentToken].SymbolType == None {

		return true

	} else {

		return false
	}
}

func (p *Parser) expression() int {

	subExpressions := make([]int, 0)

	p.open() // optional open for negatives and subexpressions

	first := p.operand()

	parent, children := p.operands(-1, append(subExpressions, first))

	if parent == -1 {

		return first // return an atom

	} else {

		p.linkChildren(parent, children)

		return parent // return a subexpression
	}
}

func (p *Parser) equation() {

	lhsAux := p.auxiliary(make([]Symbol, 0))

	lhs := p.expression()

	p.addAuxiliaries(lhs, lhsAux)

	if p.tokens[p.currentToken].IsComparison() {

		comparison := p.addNode()

		p.currentToken++

		rhsAux := p.auxiliary(make([]Symbol, 0))

		rhs := p.expression()

		p.addAuxiliaries(rhs, rhsAux)

		p.completeEquation(comparison, lhs, rhs)

	} else {

		p.parsed.SetRootByIndex(lhs)
	}
}

func (p *Parser) function() int {

	if p.tokens[p.currentToken].SymbolType == Function {

		functionCall := p.addNode()

		p.currentToken++

		set := p.set()

		p.linkChild(functionCall, set)

		return functionCall
	}
	return -1
}

func (p *Parser) set() int {

	if p.tokens[p.currentToken].SymbolType == Set {

		parent := p.addNode()

		p.currentToken++

		elements := make([]int, 0)

		elements = p.setElements(elements)

		p.linkChildren(parent, elements)

		return parent
	}
	return -1
}

func (p *Parser) setElements(elements []int) []int {

	if !p.close() {

		auxiliaries := p.auxiliary(make([]Symbol, 0))

		element := p.expression()

		p.addAuxiliaries(element, auxiliaries)

		elements = append(elements, element)

		elements = p.setElements(elements)

		return elements

	} else {

		return elements
	}
}

func (p *Parser) addNode() int {

	return p.parsed.AddToMap(p.tokens[p.currentToken])
}

func (p *Parser) addAuxiliaries(index int, auxiliaries []Symbol) {

	p.parsed.InsertAuxilliariesAt(index, auxiliaries)
}

func (p *Parser) linkChild(parent int, child int) {

	p.parsed.SetParent(parent, child)
}

func (p *Parser) linkChildren(parent int, children []int) {

	for _, child := range children {

		p.parsed.SetParent(parent, child)
	}
}

func (p *Parser) completeEquation(parent int, lhs int, rhs int) {

	p.parsed.SetParent(parent, lhs)

	p.parsed.SetParent(parent, rhs)

	p.parsed.SetRootByIndex(parent)
}
