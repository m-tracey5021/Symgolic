package parsing

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	. "symgolic/symbols"
)

type ParseType int

const (
	Math = iota

	Logic

	NaturalLanguage
)

type ParseState int

const (
	NoneParsed = iota

	SubexpressionParsed

	FunctionParsed

	SetParsed

	VectorParsed
)

type Parser struct {
	program Program

	currentExpression Expression

	tokens []Symbol

	currentToken int

	currentState int
}

func NewParser() Parser {

	parser := Parser{NewProgram(), NewEmptyExpression(), make([]Symbol, 0), 0, NoneParsed}

	return parser
}

func ParseProgramFromFile(path string) Program {

	parser := NewParser()

	file, err := os.Open(path)

	if err != nil {

		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := scanner.Text()

		fmt.Print("Line to parse: ")
		fmt.Println(line)

		if line != "" {

			parser.tokens = lex(line)

			parser.currentToken = 0

			parser.equation()

			parser.program.AddExpression(parser.currentExpression)

			parser.currentExpression = NewEmptyExpression()
		}
	}
	if err := scanner.Err(); err != nil {

		panic(err)
	}
	return parser.program
}

func ParseProgramFromString(text string) Program {

	parser := NewParser()

	parser.tokens = lex(text)

	parser.lines()

	return parser.program
}

func ParseExpression(text string) Expression {

	parser := NewParser()

	parser.tokens = lex(text)

	parser.equation()

	return parser.currentExpression
}

func (p *Parser) setState(symbol SymbolType) {

	if symbol == SubExpressionOpen && p.currentState != FunctionParsed {

		p.currentState = SubexpressionParsed

	} else if symbol == Set {

		p.currentState = SetParsed

	} else if symbol == Vector {

		p.currentState = VectorParsed

	}
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

	p.setState(p.tokens[p.currentToken].SymbolType)

	if p.tokens[p.currentToken].IsEnclosingOperation() {

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
	// set := p.set()

	// if set != -1 {

	// 	p.addAuxiliaries(set, auxiliaries)

	// 	return set
	// }
	// list := p.list()

	// if list != -1 {

	// 	p.addAuxiliaries(list, auxiliaries)

	// 	return list
	// }
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

	if p.tokens[p.currentToken].ClosesExpressionScope() {

		p.currentToken++

		return true

	} else if p.tokens[p.currentToken].SymbolType == NewLine || p.tokens[p.currentToken].SymbolType == EndOfFile {

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

		p.currentExpression.SetRootByIndex(lhs)
	}
}

func (p *Parser) function() int {

	if p.tokens[p.currentToken].SymbolType == Function {

		p.currentState = FunctionParsed

		function := p.addNode()

		p.currentToken++

		params := p.expression()

		p.linkChild(function, params)

		return function
	}
	return -1
}

// func (p *Parser) function() int {

// 	if p.tokens[p.currentToken].SymbolType == Function {

// 		p.currentState = FunctionParsed

// 		functionCall := p.addNode()

// 		p.currentToken++

// 		params := p.expression()

// 		p.linkChild(functionCall, params)

// 		return functionCall
// 	}
// 	return -1
// }

// func (p *Parser) set() int {

// 	if p.tokens[p.currentToken].SymbolType == Set {

// 		parent := p.addNode()

// 		p.currentToken++

// 		elements := make([]int, 0)

// 		elements = p.parameters(elements)

// 		p.linkChildren(parent, elements)

// 		return parent
// 	}
// 	return -1
// }

// func (p *Parser) list() int {

// 	if p.tokens[p.currentToken].IsEnclosingOperation() {

// 		parent := p.addNode()

// 		p.currentToken++

// 		params := make([]int, 0)

// 		params = p.parameters(params)

// 		p.linkChildren(parent, params)

// 		return parent
// 	}
// 	return -1
// }

// func (p *Parser) parameters(params []int) []int {

// 	if !p.close() {

// 		auxiliaries := p.auxiliary(make([]Symbol, 0))

// 		param := p.expression()

// 		p.addAuxiliaries(param, auxiliaries)

// 		params = append(params, param)

// 		params = p.parameters(params)

// 		return params

// 	} else {

// 		return params
// 	}
// }

func (p *Parser) lines() {

	p.equation()

	p.program.AddExpression(p.currentExpression)

	p.currentExpression = NewEmptyExpression()

	if p.tokens[p.currentToken].SymbolType == NewLine {

		p.currentToken++

		p.lines()

	} else if p.tokens[p.currentToken].SymbolType == EndOfFile {

		return

	} else {

		panic(errors.New("unrecognised symbol"))
	}
}

func (p *Parser) addNode() int {

	if p.tokens[p.currentToken].SymbolType == Iteration {

		var enclose int

		if p.currentState == SubexpressionParsed {

			enclose = p.currentExpression.AddToMap(Symbol{NaryTuple, -1, "(...)"})

		} else if p.currentState == FunctionParsed {

			enclose = p.currentExpression.AddToMap(Symbol{FunctionParameters, -1, "(...)"})

		} else if p.currentState == SetParsed {

			enclose = p.currentExpression.AddToMap(Symbol{Set, -1, "{...}"})

		} else if p.currentState == VectorParsed {

			enclose = p.currentExpression.AddToMap(Symbol{Vector, -1, "[...]"})

		} else {

			panic(errors.New("parse not started properly"))
		}
		return enclose

	} else {

		return p.currentExpression.AddToMap(p.tokens[p.currentToken])
	}

}

func (p *Parser) addAuxiliaries(index int, auxiliaries []Symbol) {

	p.currentExpression.InsertAuxiliariesAt(index, auxiliaries)
}

func (p *Parser) linkChild(parent int, child int) {

	p.currentExpression.SetParent(parent, child)
}

func (p *Parser) linkChildren(parent int, children []int) {

	for _, child := range children {

		p.currentExpression.SetParent(parent, child)
	}
}

func (p *Parser) completeEquation(parent int, lhs int, rhs int) {

	p.currentExpression.SetParent(parent, lhs)

	p.currentExpression.SetParent(parent, rhs)

	p.currentExpression.SetRootByIndex(parent)
}
