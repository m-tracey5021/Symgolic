package symbols

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

func NewExpression() Expression {

	var expression Expression = Expression{}

	expression.auxMap = make(map[int][]Symbol)

	expression.treeMap = make(map[int]Symbol)

	expression.parentMap = make(map[int]int)

	expression.childMap = make(map[int][]int)

	// result.reverseTree = make(map[Symbol]int)

	return expression
}

func NewExpressionWithRoot(symbol Symbol) (int, Expression) {

	var expression Expression = Expression{}

	expression.auxMap = make(map[int][]Symbol)

	expression.treeMap = make(map[int]Symbol)

	expression.parentMap = make(map[int]int)

	expression.childMap = make(map[int][]int)

	root := expression.SetRoot(symbol)

	return root, expression
}

// Retrieval

func (e *Expression) GetRoot() int {

	return e.root
}

func (e *Expression) GetNumericValueByIndex(index int) int {

	return e.treeMap[index].NumericValue
}

func (e *Expression) GetNumericValueByPath(path []int) int {

	return e.GetNumericValueByIndex(e.GetChildByPath(path))
}

func (e *Expression) GetNumericValueByPathFrom(index int, path []int) int {

	return e.GetNumericValueByIndex(e.GetChildByPathFrom(index, path))
}

func (e *Expression) GetAlphaValueByIndex(index int) string {

	return e.treeMap[index].AlphaValue
}

func (e *Expression) GetAuxilliariesByIndex(index int) []Symbol {

	return e.auxMap[index]
}

func (e *Expression) GetNodeByIndex(index int) *Symbol {

	node := e.treeMap[index]

	return &node
}

// func (e *Expression) GetIndexByNode(node Symbol) int {

// 	for key, value := range e.treeMap {

// 		if reflect.DeepEqual(node, value) {

// 			return key
// 		}
// 	}
// 	return -1
// 	// return e.reverseTree[node]
// }

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

	if len(children) == 0 {
		return -1
	}

	return e.childMap[index][breadth]
}

func (e *Expression) GetChildByPath(path []int) int {

	root := e.GetRoot()

	return e.GetChildByPathFrom(root, path)
}

func (e *Expression) GetChildByPathFrom(index int, path []int) int {

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

// Identifiers

func (e *Expression) GetSymbolTypeByIndex(index int) SymbolType {

	if index < 0 {

		return None

	} else {

		return e.GetNodeByIndex(index).SymbolType
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

		if !e.IsEquality(parent) {

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

		if e.IsEquality(parent) {

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

// Checking Validity

// func (e *Expression) GetDefinition(index int) int {

// 	if e.IsVariable(index) {

// 		matches := e.FindMatches(e.GetRoot(), e.GetAlphaValueByIndex(index), make([]int, 0))

// 		for _, match := range matches {

// 			parent := e.GetParent(match)

// 			if e.IsEquality(parent) {

// 				for _, child := range e.GetChildren(parent) {

// 					if e.GetAlphaValueByIndex(child) != e.GetAlphaValueByIndex(index) {

// 						return child
// 					}
// 				}
// 			}
// 		}
// 		return -1

// 	} else {

// 		return index
// 	}
// }

// func (e *Expression) FindInstancesOf(index int, compared string, instances []int) []int {

// 	if e.GetAlphaValueByIndex(index) == compared {

// 		instances = append(instances, index)
// 	}
// 	for _, child := range e.GetChildren(index) {

// 		instances = e.FindInstancesOf(child, compared, instances)
// 	}
// 	return instances
// }

// Generating and Adding

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

func (e *Expression) AddToMapWithAux(node Symbol, auxillaries []Symbol) int {

	id := e.GenerateId()

	e.treeMap[id] = node

	e.childMap[id] = make([]int, 0)

	e.auxMap[id] = auxillaries

	return id
}

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

func (e *Expression) SetRootWithAux(node Symbol, auxillaries []Symbol) int {

	if len(e.treeMap) == 0 {

		root := e.AddToMapWithAux(node, auxillaries)

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

func (e *Expression) AppendAuxilliariesAt(index int, auxillaries []Symbol) {

	for i := 0; i < len(auxillaries); i++ {

		e.auxMap[index] = append(e.auxMap[index], auxillaries[i])
	}
	e.updateDisplay()
}

func (e *Expression) InsertAuxilliariesAt(index int, auxillaries []Symbol) {

	for i := len(auxillaries) - 1; i >= 0; i-- {

		e.auxMap[index] = append(e.auxMap[index], Symbol{})

		copy(e.auxMap[index][1:], e.auxMap[index][0:])

		e.auxMap[index][0] = auxillaries[i]
	}
	e.updateDisplay()
}

func (e *Expression) SetParent(parent int, child int) {

	e.parentMap[child] = parent

	e.childMap[parent] = append(e.childMap[parent], child)

	e.updateDisplay()
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

func (e *Expression) AppendExpression(parent int, expression Expression, copy bool) {

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

	} else {

		e.AppendExpressionRecurse(parent, expression, root)
	}
	e.updateDisplay()
}

func (e *Expression) AppendExpressionRecurse(parent int, expression Expression, transferIndex int) {

	transfer := expression.GetNodeByIndex(transferIndex)

	transferAux := expression.GetAuxilliariesByIndex(transferIndex)

	index := e.AddToMapWithAux(*transfer, transferAux)

	e.parentMap[index] = parent

	e.childMap[parent] = append(e.childMap[parent], index)

	for _, child := range expression.GetChildren(transferIndex) {

		e.AppendExpressionRecurse(index, expression, child)
	}
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

// Replacing and Removing

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

			e.RemoveNode(index, true)

			e.AppendBulkSubtreesFrom(parent, expression.GetChildren(otherRoot), expression)

		} else {

			for len(e.GetChildren(index)) != 0 {

				e.RemoveNode(e.GetChildAtBreadth(index, 0), true)
			}
			e.treeMap[index] = *expression.GetNodeByIndex(otherRoot)

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

// Arithmetic

func (e *Expression) Negate() {

	root := e.GetRoot()

	negation := make([]Symbol, 0)

	negation = append(negation, Symbol{Subtraction, -1, "-"})

	e.InsertAuxilliariesAt(root, negation)
}

func (e *Expression) Subtract(other Expression) Expression {

	sub := NewExpression()

	root := sub.SetRoot(Symbol{Addition, -1, "+"})

	sub.AppendExpression(root, *e, true)

	lhs := other.CopyTree()

	lhs.Negate()

	sub.AppendExpression(root, lhs, false)

	return sub
}

func (e *Expression) Multiply(children []int) Expression {

	result := NewExpression()

	mul := Symbol{Multiplication, -1, "*"}

	root := result.SetRoot(mul)

	result.AppendBulkSubtreesFrom(root, children, *e)

	return result
}

// Copying

func (e *Expression) CopyTree() Expression {

	copy := NewExpression()

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
	copy.updateDisplay()

	return copy
}

func (e *Expression) CopySubtree(index int) Expression {

	return e.CopySubtreeRecurse(index, -1, nil)
}

func (e *Expression) CopySubtreeByPath(path []int) Expression {

	return e.CopySubtreeRecurse(e.GetChildByPath(path), -1, nil)
}

func (e *Expression) CopySubtreeRecurse(parent int, copiedParent int, copiedExpression *Expression) Expression {

	if copiedExpression == nil {

		newExpression := NewExpression()

		copiedExpression = &newExpression

		copiedParent = 0

		copiedExpression.SetRootWithAux(e.GetNodeByIndex(parent).Copy(), e.GetAuxilliariesByIndex(parent))
	}
	for _, child := range e.childMap[parent] {

		copiedChild := e.GetNodeByIndex(child).Copy()

		sign := e.GetAuxilliariesByIndex((child))

		index := copiedExpression.AppendNodeWithAux(copiedParent, copiedChild, sign)

		e.CopySubtreeRecurse(child, index, copiedExpression)
	}
	return *copiedExpression
}

// Printers

func (e *Expression) updateDisplay() {

	e.display = e.ToString()
}

func (e *Expression) buildString(index int) string {

	symbol := e.GetNodeByIndex(index)

	parent := e.GetParent(index)

	if e.IsAtomic(index) {

		return symbol.AlphaValue

	} else if e.IsFunction(index) {

		var function string

		function += e.GetAlphaValueByIndex(index)

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

		return list

	} else {

		var operation string

		children := e.GetChildren(index)

		for i, child := range children {

			if i == 0 {

				operation += e.buildString(child)

			} else {

				operation += symbol.AlphaValue + e.buildString(child)
			}
		}
		auxiliaries := e.GetAuxilliariesByIndex(index)

		if len(auxiliaries) > 0 {

			if e.GetAuxilliariesByIndex(index)[0].SymbolType == Subtraction {

				return "-(" + operation + ")"
			}
			return operation

		} else {

			if parent == -1 || e.IsEquality(parent) ||
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
	fmt.Println(e.GetNodeByIndex(index).AlphaValue)

	for _, child := range e.GetChildren(index) {

		e.PrintTree(child, factor, depth+1)
	}
}

// 3
