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

		"Fn": FunctionDef,
	}
	var symbols []Symbol

	var characters []rune = []rune(text)

	functionDefInProgress := false

	for i := 0; i < len(characters); i++ {

		characterAt := characters[i]

		if functionDefInProgress {

			if unicode.IsLetter(characterAt) || unicode.IsDigit(characterAt) {

				symbolType, val, symbol, end := lexWord(text, characters, i, false)

				if symbolType != None {

					symbols = append(symbols, Symbol{symbolType, val, symbol})

					i = end
				}
			} else if characterAt == ' ' {

				continue

			} else {

				panic(errors.New("no function name supplied"))
			}

		} else {

			if unicode.IsLetter(characterAt) || unicode.IsDigit(characterAt) {

				symbolType, val, symbol, end := lexOperand(text, characters, i, true)

				if symbolType != None {

					symbols = append(symbols, Symbol{symbolType, val, symbol})

					i = end
				}

			} else if characterAt == ' ' || characterAt == ',' {

				continue

			} else {

				simpleToken := text[i : i+1]

				simple, simpleExists := tokens[simpleToken]

				if i+2 <= len(characters) {

					compoundToken := text[i : i+2]

					compound, compoundExists := tokens[compoundToken]

					if compoundExists {

						if compoundToken == "Fn" {

							functionDefInProgress = true

						} else {

							functionDefInProgress = false
						}
						symbols = append(symbols, Symbol{compound, -1, compoundToken})

						i++

					} else {

						if simpleExists {

							symbols = append(symbols, Symbol{simple, -1, simpleToken})

						} else {

							continue
						}
					}

				} else {

					if simpleExists {

						symbols = append(symbols, Symbol{simple, -1, simpleToken})

					} else {

						continue
					}
				}
			}
		}

		if i == len(characters)-1 {

			symbols = append(symbols, Symbol{None, -1, "EOF"})
		}
	}
	return symbols
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

	for unicode.IsLetter(characters[end]) &&
		characters[end] != 'v' &&
		characters[end] != 'A' &&
		characters[end] != 'E' &&
		characters[end] != 'F' {

		end++
	}
	if end == index {

		return None, -1, "", end

	} else {

		if end+1 < len(characters) {

			if characters[end+1] == '(' {

				if predefined {

					return FunctionCall, -1, text[index:end], end

				} else {

					return Function, -1, text[index:end], end
				}

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

	for unicode.IsDigit(characters[end]) {

		end++
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

// func (p *Parser) auxillary() bool {

// 	if p.tokens[p.currentToken].SymbolType == Subtraction {

// 		p.currentToken++

// 		return false

// 	} else {

// 		return true
// 	}
// }

func (p *Parser) auxillary(auxillaries []Symbol) []Symbol {

	if p.tokens[p.currentToken].IsAuxillary() {

		auxillaries = append(auxillaries, p.tokens[p.currentToken])

		p.currentToken++

		auxillaries = p.auxillary(auxillaries)

		return auxillaries

	} else {

		return auxillaries
	}
}

func (p *Parser) open() {

	if p.tokens[p.currentToken].SymbolType == Open {

		p.currentToken++
	}
}

func (p *Parser) atom(auxillaries []Symbol) int {

	if p.tokens[p.currentToken].SymbolType == Variable || p.tokens[p.currentToken].SymbolType == Constant {

		child := p.addNode(auxillaries)

		p.currentToken++

		return child
	}
	return -1
}

func (p *Parser) operand() int {

	auxillaries := p.auxillary(make([]Symbol, 0))

	child := p.atom(auxillaries)

	if child == -1 {

		child = p.functionCall()

		if child == -1 {

			child = p.subExpression(auxillaries)

			return child

		} else {

			return child
		}

	} else {

		return child
	}
}

func (p *Parser) operator(auxillaries []Symbol) int {

	if p.tokens[p.currentToken].IsOperation() {

		parent := p.addNode(auxillaries)

		p.currentToken++

		return parent
	}
	return -1
}

func (p *Parser) comparison() int {

	if p.tokens[p.currentToken].IsComparison() {

		parent := p.addNode(make([]Symbol, 0))

		p.currentToken++

		return parent
	}
	return -1
}

// func (p *Parser) operandChain(children []int) []int {

// 	if p.close() {

// 		return children

// 	} else if p.tokens[p.currentToken].IsOperation() {

// 		p.currentToken++

// 		children = p.operandChain(children)

// 		return children

// 	} else {

// 		children = append(children, p.operand())

// 		children = p.operandChain(children)

// 		return children
// 	}
// }

func (p *Parser) operands(auxillaries []Symbol, parent int, children []int) (int, []int) {

	if p.close() { // used directly after first operand, i.e. the lhs, so if atomic will close

		return parent, children

	} else {

		if parent == -1 {

			parent = p.operator(auxillaries)
		}

		children = append(children, p.operand())

		parent, children = p.operands(auxillaries, parent, children)

		return parent, children
	}
}

func (p *Parser) close() bool {

	if p.tokens[p.currentToken].SymbolType == Close {

		p.currentToken++

		return true

	} else if p.tokens[p.currentToken].IsComparison() || p.tokens[p.currentToken].SymbolType == None {

		return true

	} else {

		return false
	}
}

func (p *Parser) subExpression(expressionAuxillaries []Symbol) int {

	children := make([]int, 0)

	p.open()

	lhsChild := p.operand()

	children = append(children, lhsChild)

	parent, children := p.operands(expressionAuxillaries, -1, children)

	if parent == -1 {

		p.parsed.InsertAuxilliariesAt(children[0], expressionAuxillaries)

		return children[0]

	} else {

		p.complete(parent, children)

		return parent
	}
}

func (p *Parser) expression() int {

	auxillaries := p.auxillary(make([]Symbol, 0))

	subExpressions := make([]int, 0)

	first := p.subExpression(auxillaries)

	parent, children := p.operands(make([]Symbol, 0), -1, append(subExpressions, first))

	if parent == -1 {

		// p.parsed.SetRootByIndex(first)

		return first // return an atom

	} else {

		p.complete(parent, children)

		// p.parsed.SetRootByIndex(parent)

		return parent // return a subexpression
	}
}

func (p *Parser) equation() {

	lhs := p.functionDef()

	if lhs == -1 {

		lhs = p.expression()
	}

	comparison := p.comparison()

	if comparison == -1 {

		p.parsed.SetRootByIndex(lhs)

	} else {

		p.completeEquation(comparison, lhs, p.expression())
	}
}

func (p *Parser) functionDef() int {

	if p.tokens[p.currentToken].SymbolType == FunctionDef {

		p.currentToken++

		if p.tokens[p.currentToken].SymbolType == Function {

			functionDef := p.addNode(make([]Symbol, 0))

			p.currentToken++

			params := p.functionParams()

			p.complete(functionDef, params)

			return functionDef
		}
		return -1
	}
	return -1
}

func (p *Parser) functionParams() []int {

	p.open()

	params := make([]int, 0)

	params = append(params, p.atom(make([]Symbol, 0)))

	params = p.additionalParams(params)

	return params

}

func (p *Parser) additionalParams(params []int) []int {

	if !p.close() {

		params = append(params, p.atom(make([]Symbol, 0)))

		params = p.additionalParams(params)

		return params

	} else {

		return params
	}
}

func (p *Parser) functionCall() int {

	if p.tokens[p.currentToken].SymbolType == FunctionCall {

		functionCall := p.addNode(make([]Symbol, 0))

		p.currentToken++

		params := p.functionCallParams()

		p.complete(functionCall, params)

		return functionCall
	}
	return -1
}

func (p *Parser) functionCallParams() []int {

	p.open()

	params := make([]int, 0)

	params = append(params, p.expression())

	params = p.additionalFunctionCallParams(params)

	return params
}

func (p *Parser) additionalFunctionCallParams(params []int) []int {

	if !p.close() {

		params = append(params, p.expression())

		params = p.additionalParams(params)

		return params

	} else {

		return params
	}
}

func (p *Parser) addNode(auxillaries []Symbol) int {

	return p.parsed.AddToMap(p.tokens[p.currentToken], auxillaries)
}

func (p *Parser) complete(parent int, children []int) {

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
