package parsing

import (
	"strconv"
	. "symgolic/language/components"
	"unicode"
)

func lex(text string) []Symbol {

	var symbols []Symbol

	var characters []rune = []rune(text)

	for i := 0; i < len(characters); i++ {

		characterAt := characters[i]

		if i == 0 {

			symbols = append(symbols, Symbol{SymbolType: ExpressionOpen, NumericValue: -1, AlphaValue: "expression start"})
		}
		symbolType, val, symbol, end := lexOperand(text, characters, i, true) // gets name, variable or constant

		if symbolType != None {

			symbols = append(symbols, Symbol{symbolType, val, symbol})

			i = end - 1

		} else { // check for operators, ExpressionOpens and ExpressionCloses

			if characterAt == ' ' {

				continue

			} else if characterAt == '\n' {

				symbols = append(symbols, Symbol{ExpressionClose, -1, "expression end"})

				symbols = append(symbols, Symbol{NewLine, -1, "EOL"})

			} else {

				operatorType, val, operator, end := lexOperator(text, characters, i)

				symbol := Symbol{operatorType, val, operator}

				if symbol.SymbolType == Iteration {

					symbols = append(symbols, Symbol{ExpressionClose, -1, "parameter end"})

					symbols = append(symbols, symbol)

					symbols = append(symbols, Symbol{ExpressionOpen, -1, "parameter start"})

				} else if symbol.SymbolType == SubExpressionOpen || symbol.SymbolType == Set || symbol.SymbolType == Vector {

					symbols = append(symbols, symbol)

					symbols = append(symbols, Symbol{ExpressionOpen, -1, "parameter started"})

				} else if symbol.SymbolType == SubExpressionClose || symbol.SymbolType == SetClose || symbol.SymbolType == VectorClose {

					symbols = append(symbols, Symbol{ExpressionClose, -1, "parameter end"})

					symbols = append(symbols, symbol)

				} else if symbol.IsComparison() {

					symbols = append(symbols, Symbol{ExpressionClose, -1, "expression end"})

					symbols = append(symbols, symbol)

					symbols = append(symbols, Symbol{ExpressionOpen, -1, "expression start"})

				} else {

					symbols = append(symbols, symbol)
				}
				i = end
			}
		}
		if i == len(characters)-1 {

			symbols = append(symbols, Symbol{ExpressionClose, -1, "expression end"})

			symbols = append(symbols, Symbol{EndOfFile, -1, "EOF"})
		}
	}
	return symbols
}

func getSpecialString(symbol SymbolType) string {

	specialSymbols := map[SymbolType]string{

		Iteration: "(...)",

		Set: "{}",

		Vector: "[]",
	}
	special, exists := specialSymbols[symbol]

	if exists {

		return special

	} else {

		return ""
	}
}

func lexOperator(text string, characters []rune, index int) (SymbolType, int, string, int) {

	var tokens map[string]SymbolType = map[string]SymbolType{

		":=": Assignment,

		"=": Equality,

		">": GreaterThan,

		"<": LessThan,

		">=": GreaterThanOrEqualTo,

		"<=": LessThanOrEqualTo,

		"(": SubExpressionOpen,

		")": SubExpressionClose,

		"+": Addition,

		"-": Subtraction,

		"*": Multiplication,

		"/": Division,

		"^": Exponent,

		"v": Radical,

		",": Iteration,

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

		"[": Vector,

		"]": VectorClose,

		"||": Augmented,

		"U": Union,

		"N": Intersection,

		"C": Subset,

		"C=": ProperSubset,
	}

	// var tokens map[rune]SymbolType = map[rune]SymbolType{

	// 	'=': Equality,

	// 	'>': GreaterThan,

	// 	'<': LessThan,

	// 	'\u2265': GreaterThanOrEqualTo,

	// 	'\u2264': LessThanOrEqualTo,

	// 	'(': ExpressionOpen,

	// 	')': ExpressionClose,

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

	// 	'}': SetExpressionClose,

	// 	',': SetElement,

	// 	'\u222a': Union,

	// 	'\u2229': Intersection,
	// }

	simpleToken := text[index : index+1]

	simple, simpleExists := tokens[simpleToken]

	// special := getSpecialString(simple)

	// if special != "" {

	// 	simpleToken = special
	// }
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

func isReserved(word string) bool {

	reserved := []string{

		"v",
		"A",
		"E",
		"U",
		"N",
		"C",
		"C=",
	}
	for _, str := range reserved {

		if word == str {

			return true
		}
	}
	return false
}

func lexWord(text string, characters []rune, index int, predefined bool) (SymbolType, int, string, int) {

	end := index

	for i := index; i < len(characters); i++ {

		if unicode.IsLetter(characters[end]) {

			end++

		} else {

			break
		}
	}
	if end == index {

		return None, -1, "", end

	} else {

		if end < len(characters) {

			word := text[index:end]

			if !isReserved(word) {

				if characters[end] == '(' {

					return Function, -1, text[index:end], end

				} else {

					return Variable, -1, text[index:end], end
				}

			} else {

				return None, -1, "", end
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

		number := text[index:end]

		numberVal, _ := strconv.Atoi(number)

		return Constant, numberVal, number, end
	}
}
