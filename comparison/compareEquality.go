package comparison

import (
	. "symgolic/symbols"
)

func IsEqualAt(index, indexInOther int, expression, other *Expression) bool {

	if expression.GetNode(index).AlphaValue == other.GetNode(indexInOther).AlphaValue {

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

func IsEqual(expression, other Expression) bool {

	return IsEqualAt(expression.GetRoot(), other.GetRoot(), &expression, &other)
}

func IsEqualByBaseAt(index, indexInOther int, expression, other *Expression) bool {

	if expression.GetNode(index).AlphaValue == other.GetNode(indexInOther).AlphaValue {

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

func IsEqualByFormAt(formIndex, comparedIndex int, form, compared *Expression, varMap map[string]Expression) bool {

	if form.IsOperation(formIndex) && compared.IsOperation(comparedIndex) {

		if form.GetNode(formIndex).AlphaValue == compared.GetNode(comparedIndex).AlphaValue {

			children := form.GetChildren(formIndex)

			comparedChildren := compared.GetChildren(comparedIndex)

			if len(children) != len(comparedChildren) {

				return false

			} else {

				for i := 0; i < len(children); i++ {

					if !IsEqualByFormAt(children[i], comparedChildren[i], form, compared, varMap) {

						return false
					}
				}
				return true
			}

		} else {

			return false
		}

	} else if form.IsConstant(formIndex) && compared.IsConstant(comparedIndex) {

		return form.GetNode(formIndex).NumericValue == compared.GetNode(comparedIndex).NumericValue

	} else if form.IsVariable(formIndex) {

		variable := form.GetNode(formIndex).AlphaValue

		value, exists := varMap[variable]

		if exists {

			if !IsEqualAt(value.GetRoot(), comparedIndex, &value, compared) {

				return false
			}

		} else {

			varMap[variable] = compared.CopySubtree(comparedIndex)
		}
		return true

	} else {

		return false
	}
}

func IsEqualByForm(form, compared Expression) (bool, map[string]Expression) {

	variableMap := make(map[string]Expression)

	return IsEqualByFormAt(form.GetRoot(), compared.GetRoot(), &form, &compared, variableMap), variableMap
}
