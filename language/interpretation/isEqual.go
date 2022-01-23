package interpretation

import (
	"symgolic/generic"
	. "symgolic/language/components"
)

type EqualityComparer func(ExpressionIndex, ExpressionIndex) bool

func IsEqualAt(a, b ExpressionIndex) bool {

	if a.Expression.GetNode(a.Index).AlphaValue == b.Expression.GetNode(b.Index).AlphaValue {

		children := a.Expression.GetChildren(a.Index)

		otherChildren := b.Expression.GetChildren(b.Index)

		if len(children) != len(otherChildren) {

			return false

		} else {

			if a.Expression.IsCommutative(a.Index) {

				return IsEqualByCommutation(a, b, children, otherChildren, IsEqualAt)
			}
			for i := 0; i < len(children); i++ {

				if !IsEqualAt(a.At(children[i]), b.At(otherChildren[i])) {

					return false
				}
			}
			return true
		}

	} else {

		return false
	}
}

func IsEqual(a, b Expression) bool {

	return IsEqualAt(From(a), From(b))
}

func AreEqual(expressions ...Expression) bool {

	if len(expressions) == 1 {

		return true

	} else {

		for i := 1; i < len(expressions); i++ {

			if !IsEqualAt(From(expressions[0]), From(expressions[i])) {

				return false
			}
		}
		return true
	}
}

func IsEqualByBaseAt(a, b ExpressionIndex) bool {

	if a.Expression.GetNode(a.Index).AlphaValue == b.Expression.GetNode(b.Index).AlphaValue {

		if a.Expression.IsExponent(a.Index) && b.Expression.IsExponent(b.Index) {

			if !IsEqualAt(a.At(a.Expression.GetChildAtBreadth(a.Index, 0)), b.At(b.Expression.GetChildAtBreadth(b.Index, 0))) {

				return false

			} else {

				return true
			}

		} else {

			children := a.Expression.GetChildren(a.Index)

			otherChildren := b.Expression.GetChildren(b.Index)

			if len(children) != len(otherChildren) {

				return false

			} else {

				if a.Expression.IsCommutative(a.Index) {

					return IsEqualByCommutation(a, b, children, otherChildren, IsEqualByBaseAt)
				}
				for i := 0; i < len(children); i++ {

					if !IsEqualAt(a.At(children[i]), b.At(otherChildren[i])) {

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

func IsEqualByBase(a, b Expression) bool {

	return IsEqualByBaseAt(From(a), From(b))
}

func IsEqualByFormAt(form, compared ExpressionIndex, varMap map[string]Expression) bool {

	if form.Expression.IsOperation(form.Index) && compared.Expression.IsOperation(compared.Index) {

		if form.Expression.GetNode(form.Index).AlphaValue == compared.Expression.GetNode(compared.Index).AlphaValue {

			children := form.Expression.GetChildren(form.Index)

			comparedChildren := compared.Expression.GetChildren(compared.Index)

			if len(children) != len(comparedChildren) {

				return false

			} else {

				if form.Expression.IsCommutative(form.Index) { // check if is equal with commutation

					matches := 0

					visited := make([]int, 0)

					for _, child := range children {

						for j, comparedChild := range comparedChildren {

							if generic.Contains(j, visited) {

								continue
							}
							if form.Expression.IsVariable(child) {

								if CheckVariableMap(form.At(child), compared.At(comparedChild), varMap) {

									matches++

									visited = append(visited, j)

									break
								}

							} else {

								if form.Expression.GetNode(child).AlphaValue == compared.Expression.GetNode(comparedChild).AlphaValue {

									matches++

									visited = append(visited, j)

									if !IsEqualByFormAt(form.At(child), compared.At(comparedChild), varMap) {

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

						if !IsEqualByFormAt(form.At(children[i]), compared.At(comparedChildren[i]), varMap) {

							return false
						}
					}
					return true
				}

			}

		} else {

			return false
		}

	} else if form.Expression.IsConstant(form.Index) && compared.Expression.IsConstant(compared.Index) {

		return form.Expression.GetNode(form.Index).NumericValue == compared.Expression.GetNode(compared.Index).NumericValue

	} else if form.Expression.IsVariable(form.Index) {

		return CheckVariableMap(form, compared, varMap)

	} else {

		return false
	}
}

func IsEqualByForm(form, compared Expression) (bool, map[string]Expression) {

	variableMap := make(map[string]Expression)

	return IsEqualByFormAt(From(form), From(compared), variableMap), variableMap
}

func CheckVariableMap(form, compared ExpressionIndex, varMap map[string]Expression) bool {

	variable := form.Expression.GetNode(form.Index).AlphaValue

	value, exists := varMap[variable]

	if exists {

		if !IsEqualAt(From(value), compared) {

			return false
		}

	} else {

		varMap[variable] = compared.Expression.CopySubtree(compared.Index)
	}
	return true
}

func IsEqualByCommutation(a, b ExpressionIndex, children, comparedChildren []int, isEqual EqualityComparer) bool {

	matches := 0

	visited := make([]int, 0)

	for _, child := range children {

		for j, comparedChild := range comparedChildren {

			if generic.Contains(j, visited) {

				continue
			}
			if isEqual(a.At(child), b.At(comparedChild)) {

				matches++

				visited = append(visited, j)

				break
			}
		}
	}
	return matches == len(children)
}
