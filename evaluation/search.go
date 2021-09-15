package evaluation

import (
	. "symgolic/symbols"
)

type search func(int, *Expression) (bool, Expression)

func SearchAndReplace(index int, expression *Expression, replaceFunc search) {

	for _, child := range expression.GetChildren(index) {

		SearchAndReplace(child, expression, replaceFunc)
	}
	change, result := replaceFunc(index, expression)

	if change {

		expression.ReplaceNodeCascade(index, result)
	}
}
