package symbols

import (
	"errors"
)

type Expression struct {
	root int

	signMap map[int]bool

	treeMap map[int]Symbol

	parentMap map[int]int

	childMap map[int][]int

	reverseTree map[Symbol]int
}

// New

func NewExpression() Expression {

	var result Expression = Expression{}

	result.signMap = make(map[int]bool)

	result.treeMap = make(map[int]Symbol)

	result.parentMap = make(map[int]int)

	result.childMap = make(map[int][]int)

	result.reverseTree = make(map[Symbol]int)

	return result
}

// Retrieval

func (e *Expression) GetValuebyIndex(index int) int {

	return e.treeMap[index].NumericValue
}

func (e *Expression) GetSignByIndex(index int) bool {

	return e.signMap[index]
}

func (e *Expression) GetNodeByIndex(index int) *Symbol {

	node := e.treeMap[index]

	return &node
}

func (e *Expression) GetIndexByNode(node Symbol) int {

	return e.reverseTree[node]
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

func (e *Expression) AddToMap(node Symbol, sign bool) int {

	// id := e.GenerateId()

	id := len(e.treeMap)

	_, exists := e.reverseTree[node]

	if !exists {

		e.reverseTree[node] = id

		e.treeMap[id] = node

		e.childMap[id] = make([]int, 0)

		e.signMap[id] = sign

		return id

	} else {

		panic(errors.New("element already exists in map"))
	}
}

func (e *Expression) SetRoot(node Symbol, sign bool) int {

	if len(e.treeMap) == 0 {

		root := e.AddToMap(node, sign)

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

	e.treeMap = expression.treeMap

	e.parentMap = expression.parentMap

	e.childMap = expression.childMap

	e.reverseTree = expression.reverseTree

	e.root = expression.root

	return e.root
}

func (e *Expression) SetParent(parent int, child int) {

	e.parentMap[child] = parent

	e.childMap[parent] = append(e.childMap[parent], child)
}

func (e *Expression) AppendNode(parent int, child Symbol, childSign bool) int {

	index := e.AddToMap(child, childSign)

	var childIndex int = len(e.treeMap) - 1

	e.parentMap[childIndex] = parent

	e.childMap[parent] = append(e.childMap[parent], childIndex)

	return index
}

func (e *Expression) AppendExpression(parent int, expression Expression, transferSign bool, transferIndex int) int {

	transfer := expression.GetNodeByIndex(transferIndex)

	index := e.AddToMap(*transfer, transferSign)

	e.parentMap[index] = parent

	e.childMap[parent] = append(e.childMap[parent], index)

	for _, child := range expression.GetChildren(transferIndex) {

		childSign := expression.GetSignByIndex(child)

		e.AppendExpression(index, expression, childSign, child)
	}
	return index
}

func (e *Expression) AppendSubtree(parent int, child int) int {

	return e.AppendExpression(parent, e.CopySubtree(child, 0, nil), e.GetSignByIndex(0), 0)
}

func (e *Expression) AppendSubtreeFrom(parent int, child int, source Expression) int {

	return e.AppendExpression(parent, source.CopySubtree(child, 0, nil), source.GetSignByIndex(0), 0)
}

func (e *Expression) AppendBulkNodes(parent int, children []int) {

	for _, child := range children {

		e.AppendExpression(parent, e.CopySubtree(child, 0, nil), e.GetSignByIndex(0), 0)
	}
}

func (e *Expression) AppendBulkNodesFrom(parent int, children []int, source Expression) {

	for _, child := range children {

		e.AppendExpression(parent, source.CopySubtree(child, 0, nil), source.GetSignByIndex(0), 0)
	}
}

func (e *Expression) AppendBulkExpressions(parent int, children []Expression) {

	for _, child := range children {

		e.AppendExpression(parent, child, child.GetSignByIndex(0), 0)
	}
}

// Replacing and Removing

// Copying

func (e *Expression) CopyTree() Expression {

	var copy Expression = Expression{}

	for key, value := range e.treeMap {

		copy.treeMap[key] = value
	}
	for key, value := range e.parentMap {

		copy.parentMap[key] = value
	}
	for key, value := range e.childMap {

		copy.childMap[key] = value
	}
	for key, value := range e.reverseTree {

		copy.reverseTree[key] = value
	}
	copy.root = e.root

	return copy
}

func (e *Expression) CopySubtree(parent int, copiedParent int, copiedExpression *Expression) Expression {

	if copiedExpression == nil {

		copiedExpression = &Expression{}

		copiedParent = 0

		copiedExpression.SetRoot(e.GetNodeByIndex(parent).Copy(), e.GetSignByIndex(parent))
	}
	for _, child := range e.childMap[parent] {

		copiedChild := e.GetNodeByIndex(child).Copy()

		sign := e.GetSignByIndex((child))

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
		if !e.GetSignByIndex(index) {

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
