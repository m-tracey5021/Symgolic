package symbols

import (
	"errors"
	"reflect"
)

type Expression struct {
	root int

	auxMap map[int][]Symbol

	treeMap map[int]Symbol

	parentMap map[int]int

	childMap map[int][]int

	// reverseTree map[Symbol]int
}

// New

func NewExpression() Expression {

	var expression Expression = Expression{}

	expression.auxMap = make(map[int][]Symbol)

	expression.treeMap = make(map[int]Symbol)

	expression.parentMap = make(map[int]int)

	expression.childMap = make(map[int][]int)

	// result.reverseTree = make(map[Symbol]int)

	return expression
}

// Retrieval

func (e *Expression) GetRoot() int {

	return e.root
}

func (e *Expression) GetValuebyIndex(index int) int {

	return e.treeMap[index].NumericValue
}

func (e *Expression) GetAuxilliariesByIndex(index int) []Symbol {

	return e.auxMap[index]
}

func (e *Expression) GetNodeByIndex(index int) *Symbol {

	node := e.treeMap[index]

	return &node
}

func (e *Expression) GetIndexByNode(node Symbol) int {

	for key, value := range e.treeMap {

		if reflect.DeepEqual(node, value) {

			return key
		}
	}
	return -1
	// return e.reverseTree[node]
}

func (e *Expression) GetParent(index int) int {

	if e.root == index {

		return -1

	} else {

		return e.parentMap[index]
	}
}

func (e *Expression) GetParentAtDepth(index int, depth int) int {

	if e.root == index {

		return -1

	} else {

		var nextParent int = e.parentMap[index]

		for i := 0; i < depth; i++ {

			nextParent = e.parentMap[nextParent]
		}
		return nextParent
	}
}

func (e *Expression) GetChildren(index int) []int {

	return e.childMap[index]
}

func (e *Expression) GetChildAtBreadth(index int, breadth int) int {

	return e.childMap[index][breadth]
}

func (e *Expression) GetChildByPath(index int, path []int) int {

	var nextChild int = e.childMap[index][path[0]]

	for i := 1; i < len(path); i++ {

		nextChild = e.childMap[nextChild][path[i]]
	}
	return nextChild
}

// Identifiers

func (e *Expression) GetSymbolTypeByIndex(index int) SymbolType {

	return e.GetNodeByIndex(index).SymbolType
}

func (e *Expression) IsAtomic(index int) bool {

	symbolType := e.GetSymbolTypeByIndex(index)

	if symbolType == Variable || symbolType == Constant {

		return true

	} else {

		return false
	}
}

// Generating and Adding

func (e *Expression) GenerateId() int {

	var id int = 0

	for k := range e.treeMap {

		if k == id {

			id++

		} else {

			return id
		}
	}
	return id
}

func (e *Expression) AddToMap(node Symbol) int {

	// id := e.GenerateId()

	id := len(e.treeMap)

	// _, exists := e.reverseTree[node]

	e.treeMap[id] = node

	e.childMap[id] = make([]int, 0)

	return id
}

func (e *Expression) AddToMapWithAux(node Symbol, auxillaries []Symbol) int {

	// id := e.GenerateId()

	id := len(e.treeMap)

	// _, exists := e.reverseTree[node]

	e.treeMap[id] = node

	e.childMap[id] = make([]int, 0)

	e.auxMap[id] = auxillaries

	return id
}

func (e *Expression) SetRoot(node Symbol) int {

	if len(e.treeMap) == 0 {

		root := e.AddToMap(node)

		e.root = root

		return e.root

	} else {

		panic(errors.New("tree is not empty"))
	}
}

func (e *Expression) SetRootWithAux(node Symbol, auxillaries []Symbol) int {

	if len(e.treeMap) == 0 {

		root := e.AddToMapWithAux(node, auxillaries)

		e.root = root

		return e.root

	} else {

		panic(errors.New("tree is not empty"))
	}
}

func (e *Expression) SetRootByIndex(root int) {

	e.root = root
}

func (e *Expression) SetExpressionAsRoot(expression Expression) int {

	e.root = expression.root

	e.auxMap = expression.auxMap

	e.treeMap = expression.treeMap

	e.parentMap = expression.parentMap

	e.childMap = expression.childMap

	// e.reverseTree = expression.reverseTree

	return e.root
}

func (e *Expression) AppendAuxilliariesAt(index int, auxillaries []Symbol) {

	for i := 0; i < len(auxillaries); i++ {

		e.auxMap[index] = append(e.auxMap[index], auxillaries[i])
	}
}

func (e *Expression) InsertAuxilliariesAt(index int, auxillaries []Symbol) {

	for i := len(auxillaries) - 1; i >= 0; i-- {

		e.auxMap[index] = append(e.auxMap[index], Symbol{})

		copy(e.auxMap[index][1:], e.auxMap[index][0:])

		e.auxMap[index][0] = auxillaries[i]
	}
}

func (e *Expression) SetParent(parent int, child int) {

	e.parentMap[child] = parent

	e.childMap[parent] = append(e.childMap[parent], child)
}

func (e *Expression) AppendNode(parent int, child Symbol, childAux []Symbol) int {

	index := e.AddToMapWithAux(child, childAux)

	var childIndex int = len(e.treeMap) - 1

	e.parentMap[childIndex] = parent

	e.childMap[parent] = append(e.childMap[parent], childIndex)

	return index
}

func (e *Expression) AppendExpression(parent int, expression Expression, transferAux []Symbol, transferIndex int) int {

	transfer := expression.GetNodeByIndex(transferIndex)

	index := e.AddToMapWithAux(*transfer, transferAux)

	e.parentMap[index] = parent

	e.childMap[parent] = append(e.childMap[parent], index)

	for _, child := range expression.GetChildren(transferIndex) {

		childSign := expression.GetAuxilliariesByIndex(child)

		e.AppendExpression(index, expression, childSign, child)
	}
	return index
}

func (e *Expression) AppendSubtree(parent int, child int) int {

	return e.AppendExpression(parent, e.CopySubtree(child, 0, nil), e.GetAuxilliariesByIndex(0), 0)
}

func (e *Expression) AppendSubtreeFrom(parent int, child int, source Expression) int {

	return e.AppendExpression(parent, source.CopySubtree(child, 0, nil), source.GetAuxilliariesByIndex(0), 0)
}

func (e *Expression) AppendBulkNodes(parent int, children []int) {

	for _, child := range children {

		e.AppendExpression(parent, e.CopySubtree(child, 0, nil), e.GetAuxilliariesByIndex(0), 0)
	}
}

func (e *Expression) AppendBulkNodesFrom(parent int, children []int, source Expression) {

	for _, child := range children {

		e.AppendExpression(parent, source.CopySubtree(child, 0, nil), source.GetAuxilliariesByIndex(0), 0)
	}
}

func (e *Expression) AppendBulkExpressions(parent int, children []Expression) {

	for _, child := range children {

		e.AppendExpression(parent, child, child.GetAuxilliariesByIndex(0), 0)
	}
}

// Replacing and Removing

func (e *Expression) ReplaceNode(index int, symbol Symbol) {

	e.treeMap[index] = symbol
}

func (e *Expression) ReplaceNodeCascade(index int, expression Expression) {

	parent := e.GetParent(index)

	if parent == -1 {

		e.SetExpressionAsRoot(expression)

	} else if parent >= 0 && len(e.treeMap) != 0 {

		otherRoot := expression.GetRoot()

		e.treeMap[index] = *expression.GetNodeByIndex(otherRoot)

		for len(e.GetChildren(index)) != 0 {

			e.RemoveNode(e.GetChildAtBreadth(index, 0), true)
		}
		for _, otherChild := range expression.GetChildren(otherRoot) {

			otherAux := expression.GetAuxilliariesByIndex(otherChild)

			e.AppendExpression(index, expression, otherAux, otherChild)
		}
	}
}

func (e *Expression) RemoveNode(index int, startIndex bool) {

	for _, child := range e.GetChildren(index) {

		e.RemoveNode(child, false)
	}
	if startIndex && index != e.GetRoot() {

		parent := e.GetParent(index)

		if parent != -1 {

			e.childMap[parent] = append(e.childMap[parent][:index], e.childMap[parent][index+1:]...)
		}
	}
	delete(e.treeMap, index)

	delete(e.parentMap, index)

	delete(e.childMap, index)
}

// Copying

func (e *Expression) CopyTree() Expression {

	var copy Expression = Expression{}

	copy.root = e.root

	for key, value := range e.auxMap {

		copy.auxMap[key] = value
	}
	for key, value := range e.treeMap {

		copy.treeMap[key] = value
	}
	for key, value := range e.parentMap {

		copy.parentMap[key] = value
	}
	for key, value := range e.childMap {

		copy.childMap[key] = value
	}
	return copy
}

func (e *Expression) CopySubtree(parent int, copiedParent int, copiedExpression *Expression) Expression {

	if copiedExpression == nil {

		copiedExpression = &Expression{}

		copiedParent = 0

		copiedExpression.SetRootWithAux(e.GetNodeByIndex(parent).Copy(), e.GetAuxilliariesByIndex(parent))
	}
	for _, child := range e.childMap[parent] {

		copiedChild := e.GetNodeByIndex(child).Copy()

		sign := e.GetAuxilliariesByIndex((child))

		index := copiedExpression.AppendNode(copiedParent, copiedChild, sign)

		e.CopySubtree(child, index, copiedExpression)
	}
	return *copiedExpression
}

// Printers

func (e *Expression) buildString(index int) string {

	// if builder == nil {

	// 	builder = &strings.Builder{}
	// }

	symbol := e.GetNodeByIndex(index)

	parent := e.GetParent(index)

	if e.IsAtomic(index) {

		return symbol.CharacterValue

	} else {

		var operation string

		children := e.GetChildren(index)

		for i, child := range children {

			// substring := e.ToString(child, builder)

			if i == 0 {

				// builder.WriteString(substring)

				operation += e.buildString(child)

			} else {

				// builder.WriteString(symbol.characterValue)

				// builder.WriteString(substring)

				operation += symbol.CharacterValue + e.buildString(child)
			}
		}
		if e.GetAuxilliariesByIndex(index)[0].SymbolType == Subtraction {

			return "-(" + operation + ")"

		} else {

			if parent == -1 {

				return operation

			} else {

				return "(" + operation + ")"
			}
		}
	}
}

func (e *Expression) ToString() string {

	if len(e.treeMap) == 0 {

		return ""

	} else {

		return e.buildString(e.root)
	}
}
