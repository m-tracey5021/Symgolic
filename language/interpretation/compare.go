package interpretation

import (
	. "symgolic/language/components"
)

type Comparison func(int, int, Expression, Expression) bool

func Compare(index, indexInOther int, expression, other Expression, compareFunc Comparison) bool {

	return compareFunc(index, indexInOther, expression, other)
}
