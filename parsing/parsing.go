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

		",": SetElement,

		"{": Set,

		"}": SetClose,
	}

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
			characters[end] != 'E' {

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

// after +, -,

type Parser struct {
	parsed Expression

	// functions []FunctionDef

	tokens []Symbol

	currentToken int
}

// type FunctionDef struct {
// 	name string

// 	params []Symbol

// 	definition Expression
// }

func NewParser(text string, parseType int) Parser {

	parser := Parser{NewExpression(), lex(text, parseType), 0}

	return parser
}

func ParseExpression(text string, parseType int) Expression {

	parser := NewParser(text, parseType)

	// parser.expression()

	parser.equation()

	return parser.parsed
}

// func (p *Parser) auxiliary() bool {

// 	if p.tokens[p.currentToken].SymbolType == Subtraction {

// 		p.currentToken++

// 		return false

// 	} else {

// 		return true
// 	}
// }

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

// func (p *Parser) atom(auxiliaries []Symbol) int {

// 	if p.tokens[p.currentToken].SymbolType == Variable || p.tokens[p.currentToken].SymbolType == Constant {

// 		child := p.addNode(auxiliaries)

// 		p.currentToken++

// 		return child
// 	}
// 	return -1
// }

func (p *Parser) atom() int {

	// auxiliaries := make([]Symbol, 0)

	// auxiliaries = p.auxiliary(auxiliaries)

	if p.tokens[p.currentToken].SymbolType == Variable || p.tokens[p.currentToken].SymbolType == Constant {

		child := p.addNode()

		p.currentToken++

		return child
	}
	return -1
}

// func (p *Parser) operand() int {

// 	tmpToken := p.currentToken

// 	child := p.atom()

// 	if child == -1 {

// 		p.currentToken = tmpToken

// 		child = p.set()

// 		if child == -1 {

// 			p.currentToken = tmpToken

// 			child = p.function()

// 			if child == -1 {

// 				p.currentToken = tmpToken

// 				child = p.expression()

// 				return child

// 			} else {

// 				return child
// 			}

// 		} else {

// 			return child
// 		}

// 	} else {

// 		return child
// 	}
// }

func (p *Parser) operand() int {

	tmpToken := p.currentToken

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
	p.currentToken = tmpToken // expression checks for auxiliaries within, so reset pointer

	subExpression := p.expression()

	if subExpression != -1 {

		return subExpression
	}
	return -1

}

// func (p *Parser) operator(auxiliaries []Symbol) int {

// 	if p.tokens[p.currentToken].IsOperation() {

// 		parent := p.addNode(auxiliaries)

// 		p.currentToken++

// 		return parent
// 	}
// 	return -1
// }

// func (p *Parser) comparison() int {

// 	if p.tokens[p.currentToken].IsComparison() {

// 		parent := p.addNode()

// 		p.currentToken++

// 		return parent
// 	}
// 	return -1
// }

// func (p *Parser) operands(auxiliaries []Symbol, parent int, children []int) (int, []int) {

// 	if p.close() { // used directly after first operand, i.e. the lhs, so if atomic will close

// 		return parent, children

// 	} else {

// 		if parent == -1 {

// 			parent = p.operator(auxiliaries)
// 		}

// 		children = append(children, p.operand())

// 		parent, children = p.operands(auxiliaries, parent, children)

// 		return parent, children
// 	}
// }

func (p *Parser) operands(auxiliaries []Symbol, parent int, children []int) (int, []int) {

	if p.close() { // used directly after first operand, i.e. the lhs, so if atomic will close

		return parent, children

	} else {

		if p.tokens[p.currentToken].IsOperation() {

			if parent == -1 {

				parent = p.addNode()

				p.addAuxiliaries(parent, auxiliaries)
			}
			p.currentToken++
		}

		children = append(children, p.operand())

		parent, children = p.operands(auxiliaries, parent, children)

		return parent, children
	}
}

// func (p *Parser) close() bool {

// 	if p.tokens[p.currentToken].SymbolType == Close {

// 		p.currentToken++

// 		return true

// 	} else if p.tokens[p.currentToken].IsComparison() || p.tokens[p.currentToken].SymbolType == None || p.tokens[p.currentToken].SymbolType == ParameterClose {

// 		return true

// 	} else {

// 		return false
// 	}
// }

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

// func (p *Parser) subExpression(expressionauxiliaries []Symbol) int {

// 	children := make([]int, 0)

// 	p.open()

// 	lhsChild := p.operand()

// 	children = append(children, lhsChild)

// 	parent, children := p.operands(expressionauxiliaries, -1, children)

// 	if parent == -1 {

// 		p.parsed.InsertAuxilliariesAt(children[0], expressionauxiliaries)

// 		return children[0]

// 	} else {

// 		p.complete(parent, children)

// 		return parent
// 	}
// }

// func (p *Parser) expression() int {

// 	auxiliaries := p.auxiliary(make([]Symbol, 0))

// 	subExpressions := make([]int, 0)

// 	first := p.subExpression(auxiliaries)

// 	parent, children := p.operands(make([]Symbol, 0), -1, append(subExpressions, first))

// 	if parent == -1 {

// 		return first // return an atom

// 	} else {

// 		p.complete(parent, children)

// 		return parent // return a subexpression
// 	}
// }

func (p *Parser) expression(checkForAux bool) int {

	// p.open() // this matches close #1

	if checkForAux {

		auxiliaries := p.auxiliary(make([]Symbol, 0))
	}

	subExpressions := make([]int, 0)

	p.open() // optional open for negatives and subexpressions

	first := p.operand()

	parent, children := p.operands(make([]Symbol, 0), -1, append(subExpressions, first)) // close #2 in here

	// if optionalOpen {

	// 	p.close() // close that matches the optional open
	// }
	if parent == -1 {

		p.addAuxiliaries(first, auxiliaries)

		return first // return an atom

	} else {

		p.linkChildren(parent, children)

		p.addAuxiliaries(parent, auxiliaries)

		return parent // return a subexpression
	}
}

func (p *Parser) equation() {

	lhs := p.expression()

	if p.tokens[p.currentToken].IsComparison() {

		comparison := p.addNode()

		p.currentToken++

		rhs := p.expression()

		p.completeEquation(comparison, lhs, rhs)

	} else {

		p.parsed.SetRootByIndex(lhs)
	}
	// comparison := p.comparison()

	// if comparison == -1 {

	// 	p.parsed.SetRootByIndex(lhs)

	// } else {

	// 	rhs := p.expression()

	// 	p.completeEquation(comparison, lhs, rhs)
	// }
}

// func (p *Parser) functionDef() int {

// 	p.open()

// 	if p.tokens[p.currentToken].SymbolType == FunctionDef {

// 		p.currentToken++

// 		if p.tokens[p.currentToken].SymbolType == Function {

// 			functionDef := p.addNode(make([]Symbol, 0))

// 			p.currentToken++

// 			p.open()

// 			params := make([]int, 0)

// 			params = p.setElements(params)

// 			p.close()

// 			p.complete(functionDef, params)

// 			return functionDef
// 		}
// 		return -1
// 	}
// 	return -1
// }

// func (p *Parser) functionCall() int {

// 	if p.tokens[p.currentToken].SymbolType == FunctionCall {

// 		functionCall := p.addNode(make([]Symbol, 0))

// 		p.currentToken++

// 		p.open()

// 		params := make([]int, 0)

// 		params = p.functionParams(params)

// 		p.complete(functionCall, params)

// 		return functionCall
// 	}
// 	return -1
// }

// func (p *Parser) functionDef() int {

// 	p.open()

// 	if p.tokens[p.currentToken].SymbolType == FunctionDef {

// 		p.currentToken++

// 		if p.tokens[p.currentToken].SymbolType == Function {

// 			functionDef := p.addNode(make([]Symbol, 0))

// 			p.currentToken++

// 			set := p.set()

// 			p.linkChild(functionDef, set)

// 			return functionDef
// 		}
// 		return -1
// 	}
// 	return -1
// }

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

	// auxiliaries := make([]Symbol, 0)

	// auxiliaries = p.auxiliary(auxiliaries)

	if p.tokens[p.currentToken].SymbolType == Set {

		parent := p.addNode()

		// p.addAuxiliaries(parent, auxiliaries)

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

		elements = append(elements, p.expression())

		elements = p.setElements(elements)

		return elements

	} else {

		return elements
	}
}

// func (p *Parser) functionParams(params []int) []int {

// 	p.open()

// 	param, closeType := p.expression()

// 	params = append(params, param)

// 	if closeType == ParameterClose {

// 		params = p.functionParams(params)

// 		return params

// 	} else {

// 		return params
// 	}

// }

// func (p *Parser) additionalParams(params []int) []int {

// 	if p.paramClose() {

// 		p.currentToken++

// 		params = append(params, p.expression())

// 		params = p.additionalParams(params)

// 		return params

// 	} else {

// 		return params
// 	}
// }

// func (p *Parser) addNode(auxiliaries []Symbol) int {

// 	return p.parsed.AddToMap(p.tokens[p.currentToken], auxiliaries)
// }

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

// func (p *Parser) getFunctionDef(name string) ([]Symbol, Expression) {

// 	for _, function := range p.functions {

// 		if function.name == name {

// 			return function.params, function.definition
// 		}
// 	}
// 	panic(errors.New("function is not defined"))
// }
