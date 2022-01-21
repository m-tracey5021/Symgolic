package components

import (
	"errors"
	"fmt"
)

type Expression struct {
	root int

	auxMap map[int][]Symbol

	treeMap map[int]Symbol

	parentMap map[int]int

	childMap map[int][]int

	display string
}

// New

func NewEmptyExpression() Expression {

	var expression Expression = Expression{}

	expression.root = -1

	expression.auxMap = make(map[int][]Symbol)

	expression.treeMap = make(map[int]Symbol)

	expression.parentMap = make(map[int]int)

	expression.childMap = make(map[int][]int)

	return expression
}

func NewExpression(symbol Symbol) (int, Expression) {

	var expression Expression = Expression{}

	expression.auxMap = make(map[int][]Symbol)

	expression.treeMap = make(map[int]Symbol)

	expression.parentMap = make(map[int]int)

	expression.childMap = make(map[int][]int)

	root := expression.SetRoot(symbol)

	return root, expression
}

// Node Retrieval

func (e *Expression) GetNode(index int) *Symbol {

	node := e.treeMap[index]

	return &node
}

// Node Relationship Retrieval

func (e *Expression) GetRoot() int {

	return e.root
}

func (e *Expression) GetNodeByPath(path []int) int {

	root := e.GetRoot()

	return e.GetChildByPath(root, path)
}

func (e *Expression) GetAuxiliaries(index int) []Symbol {

	return e.auxMap[index]
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

	children := e.childMap[index]

	if len(children) == 0 || len(children) <= breadth {
		return -1
	}

	return e.childMap[index][breadth]
}

func (e *Expression) GetChildByPath(index int, path []int) int {

	nextChildren := e.childMap[index]

	childInRange := func() bool {

		if len(nextChildren) == 0 || len(nextChildren) <= path[0] {

			return false

		} else {

			return true
		}
	}

	if !childInRange() {

		return -1
	}

	nextChild := nextChildren[path[0]]

	for i := 1; i < len(path); i++ {

		nextChildren = e.childMap[nextChild]

		if !childInRange() {

			return -1
		}

		nextChild = nextChildren[path[i]]
	}
	return nextChild
}

func (e *Expression) GetSiblings(index int) []int {

	parent := e.GetParent(index)

	if parent == -1 {

		return make([]int, 0)

	} else {

		children := e.GetChildren(parent)

		siblings := make([]int, 0)

		for _, sibling := range children {

			if sibling != index {

				siblings = append(siblings, sibling)

				break
			}
		}
		return siblings
	}
}

func (e *Expression) GetSiblingsAndSelf(index int) []int {

	parent := e.GetParent(index)

	if parent == -1 {

		return make([]int, 0)

	} else {

		return e.GetChildren(parent)
	}
}

func (e *Expression) GetIndexAsChild(index int) int {

	for i, sibling := range e.GetSiblingsAndSelf(index) {

		if sibling == index {

			return i
		}
	}
	return -1
}

// Setters

func (e *Expression) SetRoot(node Symbol) int {

	if len(e.treeMap) == 0 {

		root := e.AddToMap(node)

		e.root = root

		e.updateDisplay()

		return e.root

	} else {

		panic(errors.New("tree is not empty"))
	}
}

func (e *Expression) SetRootWithAux(node Symbol, auxiliaries []Symbol) int {

	if len(e.treeMap) == 0 {

		root := e.AddToMapWithAux(node, auxiliaries)

		e.root = root

		e.updateDisplay()

		return e.root

	} else {

		panic(errors.New("tree is not empty"))
	}
}

func (e *Expression) SetRootByIndex(root int) {

	e.root = root

	e.updateDisplay()
}

func (e *Expression) SetExpressionAsRoot(expression Expression) int {

	e.root = expression.root

	e.auxMap = expression.auxMap

	e.treeMap = expression.treeMap

	e.parentMap = expression.parentMap

	e.childMap = expression.childMap

	e.updateDisplay()

	return e.root
}

func (e *Expression) SetParent(parent int, child int) {

	e.parentMap[child] = parent

	e.childMap[parent] = append(e.childMap[parent], child)

	e.updateDisplay()
}

func (e *Expression) SetNumericValue(index, value int) {

	e.GetNode(index).NumericValue = value

	e.updateDisplay()
}

func (e *Expression) SetAlphaValue(index int, value string) {

	e.GetNode(index).AlphaValue = value

	e.updateDisplay()
}

// Copying

func (e Expression) CopyTree() Expression {

	copy := NewEmptyExpression()

	copy.root = e.root

	for key, value := range e.auxMap {

		copy.auxMap[key] = value
	}
	for key, value := range e.treeMap {

		copy.treeMap[key] = value.Copy()
	}
	for key, value := range e.parentMap {

		copy.parentMap[key] = value
	}
	for key, value := range e.childMap {

		copiedChildren := make([]int, 0)

		for _, child := range value {

			copiedChildren = append(copiedChildren, child)
		}
		copy.childMap[key] = copiedChildren
	}
	copy.updateDisplay()

	return copy
}

func (e Expression) CopySubtree(index int) Expression {

	return e.CopySubtreeRecurse(index, -1, nil)
}

func (e Expression) CopySubtreeRecurse(parent int, copiedParent int, copiedExpression *Expression) Expression {

	if copiedExpression == nil {

		emptyExpression := NewEmptyExpression()

		copiedExpression = &emptyExpression

		copiedParent = 0

		copiedExpression.SetRootWithAux(e.GetNode(parent).Copy(), e.GetAuxiliaries(parent))
	}
	for _, child := range e.childMap[parent] {

		copiedChild := e.GetNode(child).Copy()

		sign := e.GetAuxiliaries(child)

		index := copiedExpression.AppendNodeWithAux(copiedParent, copiedChild, sign)

		e.CopySubtreeRecurse(child, index, copiedExpression)
	}
	return *copiedExpression
}

// Printers

func (e *Expression) updateDisplay() {

	e.display = e.ToString()
}

func (e *Expression) buildAuxiliaryString(index int) string {

	var auxiliaries string

	for _, aux := range e.GetAuxiliaries(index) {

		auxiliaries += aux.AlphaValue
	}
	return auxiliaries
}

func (e *Expression) buildString(index int) string {

	symbol := e.GetNode(index)

	parent := e.GetParent(index)

	if e.IsAtomic(index) {

		return e.buildAuxiliaryString(index) + symbol.AlphaValue

	} else if e.IsFunction(index) {

		var function string

		function += e.GetNode(index).AlphaValue

		function += "("

		for i, child := range e.GetChildren(index) {

			if i == 0 {

				function += e.buildString(child)

			} else {

				function += "," + e.buildString(child)
			}
		}
		function += ")"

		return function

	} else if e.IsNaryTuple(index) || e.IsSet(index) || e.IsVector(index) {

		bracketMap := map[SymbolType]string{

			NaryTuple: "()",

			Set: "{}",

			Vector: "[]",
		}
		var list string

		list += bracketMap[e.GetSymbolTypeByIndex(index)][:1]

		for i, child := range e.GetChildren(index) {

			if i == 0 {

				list += e.buildString(child)

			} else {

				list += "," + e.buildString(child)
			}
		}
		list += bracketMap[e.GetSymbolTypeByIndex(index)][1:]

		return e.buildAuxiliaryString(index) + list

	} else {

		var operation string

		children := e.GetChildren(index)

		for i, child := range children {

			if i == 0 {

				operation += e.buildString(child)

			} else {

				childAux := e.GetAuxiliaries(child)

				if e.IsSummation(index) && len(childAux) > 0 {

					if childAux[0].SymbolType == Subtraction {

						operation += e.buildString(child)
					}

				} else {

					operation += symbol.AlphaValue + e.buildString(child)
				}
			}
		}
		auxiliaries := e.GetAuxiliaries(index)

		if len(auxiliaries) > 0 {

			return e.buildAuxiliaryString(index) + "(" + operation + ")"

		} else {

			if parent == -1 ||
				e.IsEquality(parent) ||
				e.IsAssignment(parent) ||
				e.IsFunction(parent) ||
				e.IsNaryTuple(parent) ||
				e.IsSet(parent) ||
				e.IsVector(parent) {

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

func (e *Expression) PrintTree(index int, factor int, depth int) {

	for i := 0; i < (depth * factor); i++ {

		fmt.Print(" ")
	}
	fmt.Println(e.GetNode(index).AlphaValue)

	for _, child := range e.GetChildren(index) {

		e.PrintTree(child, factor, depth+1)
	}
}

// 3
