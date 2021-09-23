package evaluation

import (
	. "symgolic/symbols"
)

func IsEqual(index, indexInOther int, expression, other *Expression) bool {

	if expression.GetAlphaValueByIndex(index) == other.GetAlphaValueByIndex(indexInOther) {

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

func IsEqualByBase(index, indexInOther int, expression, other *Expression) bool {

	if expression.GetAlphaValueByIndex(index) == other.GetAlphaValueByIndex(indexInOther) {

		if expression.IsExponent(index) && other.IsExponent(indexInOther) {

			if !IsEqual(expression.GetChildAtBreadth(index, 0), other.GetChildAtBreadth(indexInOther, 0), expression, other) {

				return false

			} else {

				return true
			}

		} else {

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
		}

	} else {

		return false
	}
}
