package parsing

import (
	"errors"
	. "symgolic/symbols"
	"unicode"
)

type ExpectedType int

const (
	Operation = iota

	Atom

	OpenBracket

	CloseBrakcet
)

func _getTokenType(character string) int {

	if character == "+" || character == "-" {

		return Addition

	} else if character == "-" {

		return Negation

	} else if character == "*" {

		return Multipliaction

	} else if character == "/" {

		return Division

	} else if character == "^" {

		return Exponent

	} else if character == "v" {

		return Radical

	} else {

		return None

	}
}

func _lex(text string) ([]Token, error) {

	var terminators map[rune]bool = map[rune]bool{

		'+': true,

		'-': true,

		'*': true,

		'/': true,

		'^': true,

		'v': true,

		'(': true,

		')': true,
	}

	var tokens []Token

	var characters []rune = []rune(text)

	for i := 0; i < len(characters); i++ {

		var charAt rune = characters[i]

		var exists bool = terminators[charAt]

		if exists || unicode.IsLetter(charAt) {

			var tokenAt string = text[i : i+1]

			var tokenType = _getTokenType(tokenAt)

			tokens = append(tokens, Token{tokenType: tokenType, value: tokenAt})

		} else if unicode.IsDigit(charAt) {

			var j int = i + 1

			for unicode.IsDigit(characters[j]) {

				j++

			}
			var tokenAt string = text[i : i+1]

			var tokenType = _getTokenType(tokenAt)

			tokens = append(tokens, Token{tokenType: tokenType, value: tokenAt})

		} else {

			return tokens, errors.New("text is not correctly fomatted")

		}
	}

	return tokens, nil
}

// after +, -,

type Parser struct {
	expression Expression

	tokens []Token

	currentToken int
}

func ParseExpression(text string) (Expression, error) {

	var expression Expression = Expression{}

	tokens, err := _lex(text)

	if err == nil {

		_expression(expression, tokens, 0, 0)
	}

	return expression, nil
}

func _expression(expression Expression, tokens []Token, tokenIndex int, parentIndex int) int {

	tokenIndex, sign := _auxillary(expression, tokens, tokenIndex, parentIndex)

	tokenIndex = _left(expression, tokens, tokenIndex, parentIndex)

	tokenIndex = _operator(expression, tokens, tokenIndex, parentIndex)

	tokenIndex = _right(expression, tokens, tokenIndex, parentIndex)

	return tokenIndex
}

func _auxillary(expression Expression, tokens []Token, tokenIndex int, parentIndex int) (int, bool) {

	if tokens[tokenIndex].tokenType == Negation {

		return tokenIndex + 1, false

	} else {

		return tokenIndex, true
	}
}

func _left(expression Expression, tokens []Token, tokenIndex int, parentIndex int) int {

	tokenIndex, sign := _auxillary(expression, tokens, tokenIndex, parentIndex)

	if tokens[tokenIndex].tokenType == Variable || tokens[tokenIndex].tokenType == Constant {

		// append atom to tree
		return tokenIndex + 1

	} else {

		return _expression(expression, tokens, tokenIndex, parentIndex)
	}
}

func _operator(expression Expression, tokens []Token, tokenIndex int, parentIndex int) int {

	if tokens[tokenIndex].tokenType == Addition ||
		tokens[tokenIndex].tokenType == Multipliaction ||
		tokens[tokenIndex].tokenType == Division ||
		tokens[tokenIndex].tokenType == Exponent ||
		tokens[tokenIndex].tokenType == Radical {

		// append operator to tree
		return tokenIndex + 1

	} else {

		return tokenIndex
	}
}

func _right(expression Expression, tokens []Token, tokenIndex int, parentIndex int) int {

	tokenIndex, sign := _auxillary(expression, tokens, tokenIndex, parentIndex)

	if tokens[tokenIndex].tokenType == Variable || tokens[tokenIndex].tokenType == Constant {

		// append atom to tree
		return tokenIndex + 1

	} else {

		tokenIndex = _expression(expression, tokens, tokenIndex, parentIndex)

		tokenIndex = _operator(expression, tokens, tokenIndex, parentIndex) // this allows for chaining of particular operators like + and *

		return _right(expression, tokens, tokenIndex, parentIndex)
	}
}
