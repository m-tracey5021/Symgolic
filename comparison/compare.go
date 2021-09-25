package comparison

import (
	. "symgolic/symbols"
)

type Comparison func(int, int, Expression, Expression) bool

func Compare(index, indexInOther int, expression, other Expression, compareFunc Comparison) bool {

	return compareFunc(index, indexInOther, expression, other)
}
