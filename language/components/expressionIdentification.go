package components

// Identifiers

func (e *Expression) GetSymbolTypeByIndex(index int) SymbolType {

	if index < 0 {

		return None

	} else {

		return e.GetNode(index).SymbolType
	}
}

func (e *Expression) IsEmpty() bool {

	return len(e.treeMap) == 0
}

func (e *Expression) IsAssignment(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Assignment {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsEquality(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Equality {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsOperation(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Addition ||
		symbolType == Multiplication ||
		symbolType == Division ||
		symbolType == Exponent ||
		symbolType == Radical {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsCommutative(index int) bool {

	return e.IsSummation(index) || e.IsMultiplication(index)
}

func (e *Expression) IsSummation(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Addition {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsMultiplication(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Multiplication {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsDivision(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Division {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsExponent(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Exponent {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsRadical(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Radical {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsAtomic(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Variable || symbolType == Constant {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsVariable(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Variable {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsConstant(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Constant {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsFunction(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Function {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsFunctionCall(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Function {

		parent := e.GetParent(index)

		if e.GetChildren(parent)[1] == index {

			return true

		} else {

			return false
		}

	} else {

		return false
	}
}

func (e *Expression) IsFunctionDef(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Function {

		parent := e.GetParent(index)

		if e.IsAssignment(parent) && e.GetChildren(parent)[0] == index {

			return true

		} else {

			return false
		}

	} else {

		return false
	}
}

func (e *Expression) IsNaryTuple(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == NaryTuple {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsSet(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Set {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsVector(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Vector {

		return true

	} else {

		return false
	}
}

func (e *Expression) IsMatrix(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == NaryTuple {

		for _, child := range e.GetChildren(index) {

			if e.GetSymbolTypeByIndex(child) != Vector {

				return false
			}
		}
		return true

	} else {

		return false
	}
}

func (e *Expression) IsAugmentedMatrix(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	augmentedFound := false

	if symbolType == NaryTuple {

		for _, child := range e.GetChildren(index) {

			if e.GetSymbolTypeByIndex(child) != Vector {

				return false

			} else {

				aux := e.GetAuxiliaries(child)

				if len(aux) != 0 {

					if aux[0].SymbolType == Augmented {

						augmentedFound = !augmentedFound
					}
				}
			}
		}
		return augmentedFound

	} else {

		return false
	}
}
