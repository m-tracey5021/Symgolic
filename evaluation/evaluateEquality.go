package evaluation

import (
	. "symgolic/symbols"
)

func IsEqual(index, indexInOther int, expression, other *Expression) bool {

	if expression.GetAlphaValuebyIndex(index) == other.GetAlphaValuebyIndex(indexInOther) {

		children := expression.GetChildren(index)

		otherChildren := other.GetChildren(indexInOther)

		if len(children) != len(otherChildren) {

			return false

		} else {

			for i := 0; i < len(children); i++ {

				if !IsEqual(children[i], otherChildren[i], expression, other) {

					return false
				}
			}
			return true
		}

	} else {

		return false
	}
}
