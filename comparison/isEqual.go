package comparison

import (
	"symgolic/generic"
	. "symgolic/symbols"
)

type EqualityComparer func(int, int, *Expression, *Expression) bool

func IsEqualAt(index, indexInOther int, expression, other *Expression) bool {

	if expression.GetNode(index).AlphaValue == other.GetNode(indexInOther).AlphaValue {

		children := expression.GetChildren(index)

		otherChildren := other.GetChildren(indexInOther)

		if len(children) != len(otherChildren) {

			return false

		} else {

			if expression.IsCommutative(index) {

				return IsEqualByCommutation(index, indexInOther, expression, other, children, otherChildren, IsEqualAt)
			}
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

func AreEqual(expressions ...Expression) bool {

	if len(expressions) == 1 {

		return true

	} else {

		for i := 1; i < len(expressions); i++ {

			if !IsEqualAt(expressions[0].GetRoot(), expressions[i].GetRoot(), &expressions[0], &expressions[i]) {

				return false
			}
		}
		return true
	}
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

				if expression.IsCommutative(index) {

					return IsEqualByCommutation(index, indexInOther, expression, other, children, otherChildren, IsEqualByBaseAt)
				}
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

func IsEqualByBase(expression, other Expression) bool {

	return IsEqualByBaseAt(expression.GetRoot(), other.GetRoot(), &expression, &other)
}

func IsEqualByForm(form, compared Expression) (bool, map[string]Expression) {

	variableMap := make(map[string]Expression)

	return IsEqualByFormAt(form.GetRoot(), compared.GetRoot(), &form, &compared, variableMap), variableMap
}

func IsEqualByFormAt(formIndex, comparedIndex int, form, compared *Expression, varMap map[string]Expression) bool {

	if form.IsOperation(formIndex) && compared.IsOperation(comparedIndex) {

		if form.GetNode(formIndex).AlphaValue == compared.GetNode(comparedIndex).AlphaValue {

			children := form.GetChildren(formIndex)

			comparedChildren := compared.GetChildren(comparedIndex)

			if len(children) != len(comparedChildren) {

				return false

			} else {

				if form.IsCommutative(formIndex) { // check if is equal with commutation

					matches := 0

					visited := make([]int, 0)

					for _, child := range children {

						for j, comparedChild := range comparedChildren {

							if generic.Contains(j, visited) {

								continue
							}
							if form.IsVariable(child) {

								if CheckVariableMap(form, compared, child, comparedChild, varMap) {

									matches++

									visited = append(visited, j)

									break
								}

							} else {

								if form.GetNode(child).AlphaValue == compared.GetNode(comparedChild).AlphaValue {

									matches++

									visited = append(visited, j)

									if !IsEqualByFormAt(child, comparedChild, form, compared, varMap) {

										return false
									}
									break
								}
							}
						}
					}
					return matches == len(children)

				} else {

					for i := 0; i < len(children); i++ {

						if !IsEqualByFormAt(children[i], comparedChildren[i], form, compared, varMap) {

							return false
						}
					}
					return true
				}

			}

		} else {

			return false
		}

	} else if form.IsConstant(formIndex) && compared.IsConstant(comparedIndex) {

		return form.GetNode(formIndex).NumericValue == compared.GetNode(comparedIndex).NumericValue

	} else if form.IsVariable(formIndex) {

		return CheckVariableMap(form, compared, formIndex, comparedIndex, varMap)

	} else {

		return false
	}
}

func CheckVariableMap(form, compared *Expression, formIndex, comparedIndex int, varMap map[string]Expression) bool {

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
}

func IsEqualByCommutation(index, indexInOther int, expression, other *Expression, children, comparedChildren []int, isEqual EqualityComparer) bool {

	matches := 0

	visited := make([]int, 0)

	for _, child := range children {

		for j, comparedChild := range comparedChildren {

			if generic.Contains(j, visited) {

				continue
			}
			if isEqual(child, comparedChild, expression, other) {

				matches++

				visited = append(visited, j)

				break
			}
		}
	}
	return matches == len(children)
}