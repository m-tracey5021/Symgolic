package comparison

import (
	. "symgolic/symbols"
)

func IsEqualAt(index, indexInOther int, expression, other *Expression) bool {

	if expression.GetAlphaValueByIndex(index) == other.GetAlphaValueByIndex(indexInOther) {

		children := expression.GetChildren(index)

		otherChildren := other.GetChildren(indexInOther)

		if len(children) != len(otherChildren) {

			return false

		} else {

			for i := 0; i < len(children); i++ {

				if !IsEqualAt(children[i], otherChildren[i], expression, other) {

					return false
				}
			}
			return true
		}

	} else {

		return false
	}
}

func IsEqualAtBreadthFirst(index, indexInOther int, expression, other *Expression) bool {

	if expression.GetAlphaValueByIndex(index) == other.GetAlphaValueByIndex(indexInOther) {

		children := expression.GetChildren(index)

		otherChildren := other.GetChildren(indexInOther)

		if len(children) != len(otherChildren) {

			return false

		} else {

			for i := 0; i < len(children); i++ {

				if !IsEqualAt(children[i], otherChildren[i], expression, other) {

					return false
				}
			}
			return true
		}

	} else {

		return false
	}
}

func IsEqualByRoot(expression, other Expression) bool {

	return IsEqualAt(expression.GetRoot(), other.GetRoot(), &expression, &other)
}

func IsEqualByBaseAt(index, indexInOther int, expression, other *Expression) bool {

	if expression.GetAlphaValueByIndex(index) == other.GetAlphaValueByIndex(indexInOther) {

		if expression.IsExponent(index) && other.IsExponent(indexInOther) {

			if !IsEqualAt(expression.GetChildAtBreadth(index, 0), other.GetChildAtBreadth(indexInOther, 0), expression, other) {

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

					if !IsEqualAt(children[i], otherChildren[i], expression, other) {

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
