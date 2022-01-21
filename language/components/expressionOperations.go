package components

// Generating and Adding Nodes

func (e *Expression) GenerateId() int {

	i := 0

	_, exists := e.treeMap[i]

	for exists {

		i++

		_, exists = e.treeMap[i]
	}
	return i
}

func (e *Expression) AddToMap(node Symbol) int {

	id := e.GenerateId()

	e.treeMap[id] = node

	e.childMap[id] = make([]int, 0)

	return id
}

func (e *Expression) AddToMapWithAux(node Symbol, auxiliaries []Symbol) int {

	id := e.GenerateId()

	e.treeMap[id] = node

	e.childMap[id] = make([]int, 0)

	e.auxMap[id] = auxiliaries

	return id
}

func (e *Expression) AppendNode(parent int, child Symbol) int {

	index := e.AddToMap(child)

	var childIndex int = len(e.treeMap) - 1

	e.parentMap[childIndex] = parent

	e.childMap[parent] = append(e.childMap[parent], childIndex)

	e.updateDisplay()

	return index
}

func (e *Expression) AppendNodeWithAux(parent int, child Symbol, childAux []Symbol) int {

	index := e.AddToMapWithAux(child, childAux)

	var childIndex int = len(e.treeMap) - 1

	e.parentMap[childIndex] = parent

	e.childMap[parent] = append(e.childMap[parent], childIndex)

	e.updateDisplay()

	return index
}

// Auxiliaries

func (e *Expression) AppendAuxiliariesAt(index int, auxiliaries []Symbol) {

	for i := 0; i < len(auxiliaries); i++ {

		e.auxMap[index] = append(e.auxMap[index], auxiliaries[i])
	}
	e.updateDisplay()
}

func (e *Expression) InsertAuxiliariesAt(index int, auxiliaries []Symbol) {

	currentAux := e.auxMap[index]

	if len(currentAux) != 0 {

		auxiliaries = append(auxiliaries, currentAux...)

		e.auxMap[index] = auxiliaries

	} else {

		e.auxMap[index] = append(e.auxMap[index], auxiliaries...)
	}
	e.updateDisplay()
}

func (e *Expression) RemoveAuxiliariesAt(index, auxIndex int) {

	currentAux := e.auxMap[index]

	if len(currentAux) != 0 {

		currentAux = append(currentAux[:auxIndex], currentAux[auxIndex+1:]...)

		e.auxMap[index] = currentAux

	}
	e.updateDisplay()
}

// Appending Expressions

func (e *Expression) AppendExpression(parent int, expression Expression, copy bool) int {

	var result int

	root := expression.GetRoot()

	if copy {

		expression = expression.CopyTree()

		root = expression.GetRoot()
	}
	if (expression.IsMultiplication(root) && e.IsMultiplication(parent)) ||
		(expression.IsSummation(root) && e.IsSummation(parent)) {

		for _, child := range expression.GetChildren(root) {

			e.AppendSubtreeFrom(parent, child, expression)
		}
		result = parent

	} else {

		result = e.AppendExpressionRecurse(parent, expression, root)
	}
	e.updateDisplay()

	return result
}

func (e *Expression) AppendExpressionRecurse(parent int, expression Expression, transferIndex int) int {

	transfer := expression.GetNode(transferIndex)

	transferAux := expression.GetAuxiliaries(transferIndex)

	index := e.AddToMapWithAux(*transfer, transferAux)

	e.parentMap[index] = parent

	e.childMap[parent] = append(e.childMap[parent], index)

	for _, child := range expression.GetChildren(transferIndex) {

		e.AppendExpressionRecurse(index, expression, child)
	}
	return index
}

func (e *Expression) AppendSubtree(parent int, child int) {

	copy := e.CopySubtree(child)

	e.AppendExpression(parent, copy, false)
}

func (e *Expression) AppendSubtreeFrom(parent int, child int, source Expression) {

	copy := source.CopySubtree(child)

	e.AppendExpression(parent, copy, false)
}

func (e *Expression) AppendBulkSubtrees(parent int, children []int) {

	for _, child := range children {

		e.AppendSubtree(parent, child)
	}
}

func (e *Expression) AppendBulkSubtreesFrom(parent int, children []int, source Expression) {

	for _, child := range children {

		e.AppendSubtreeFrom(parent, child, source)
	}
}

func (e *Expression) AppendBulkExpressions(parent int, children []Expression) {

	for _, child := range children {

		e.AppendExpression(parent, child, false)
	}
}

// Inserting, Replacing and Removing

func (e *Expression) InsertExpression(parent, index int, expression Expression) {

	children := e.GetChildren(parent)

	inserted := make([]int, 0)

	insert := func() int {

		insertedNode := e.AddToMap(*expression.GetNode(expression.GetRoot()))

		e.parentMap[insertedNode] = parent

		e.AppendBulkSubtreesFrom(insertedNode, expression.GetChildren(expression.GetRoot()), expression)

		return insertedNode
	}

	if index == len(children) {

		insertedNode := insert()

		e.childMap[parent] = append(e.childMap[parent], insertedNode)

	} else {

		for i, child := range children {

			if i == index {

				insertedNode := insert()

				inserted = append(inserted, insertedNode)

				inserted = append(inserted, child)

			} else {

				inserted = append(inserted, child)
			}
		}
		e.childMap[parent] = inserted
	}
	e.updateDisplay()
}

func (e *Expression) ReplaceNode(index int, symbol Symbol) {

	e.treeMap[index] = symbol

	e.updateDisplay()
}

func (e *Expression) ReplaceNodeCascade(index int, expression Expression) {

	parent := e.GetParent(index)

	if parent == -1 {

		e.SetExpressionAsRoot(expression)

	} else if parent >= 0 && len(e.treeMap) != 0 {

		otherRoot := expression.GetRoot()

		if (e.IsSummation(parent) && expression.IsSummation(otherRoot)) ||
			(e.IsMultiplication(parent) && expression.IsMultiplication(otherRoot)) {

			indexAsChild := e.GetIndexAsChild(index)

			e.RemoveNode(index, true)

			for _, otherChild := range expression.GetChildren(otherRoot) {

				e.InsertExpression(parent, indexAsChild, expression.CopySubtree(otherChild))

				indexAsChild++
			}

		} else {

			for len(e.GetChildren(index)) != 0 {

				e.RemoveNode(e.GetChildAtBreadth(index, 0), true)
			}
			e.treeMap[index] = *expression.GetNode(otherRoot)

			e.AppendBulkSubtreesFrom(index, expression.GetChildren(otherRoot), expression)
		}
	}
	e.updateDisplay()
}

func (e *Expression) RemoveNode(index int, startIndex bool) {

	for _, child := range e.GetChildren(index) {

		e.RemoveNode(child, false)
	}
	if startIndex && index != e.GetRoot() {

		parent := e.GetParent(index)

		if parent != -1 {

			i := e.GetIndexAsChild(index)

			e.childMap[parent] = append(e.childMap[parent][:i], e.childMap[parent][i+1:]...)
		}
	}
	delete(e.treeMap, index)

	delete(e.parentMap, index)

	delete(e.childMap, index)

	e.updateDisplay()
}
