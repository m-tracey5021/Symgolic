package components

func SSS(a, b Expression, reverse bool) Expression { // adding two sums

	if reverse {

		return SSS(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.GetChildren(a.GetRoot()), a)

	result.AppendBulkSubtreesFrom(root, b.GetChildren(b.GetRoot()), b)

	return result
}

func SSM(a, b Expression, reverse bool) Expression { // adding sum and multiplication

	if reverse {

		return SSM(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.GetChildren(a.GetRoot()), a)

	result.AppendSubtreeFrom(root, b.GetRoot(), b)

	return result
}

func SSD(a, b Expression, reverse bool) Expression { // adding sum and division

	if reverse {

		return SSD(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.GetChildren(a.GetRoot()), a)

	result.AppendSubtreeFrom(root, b.GetRoot(), b)

	return result

}

func MSM(a, b Expression, reverse bool) Expression { // adding multiplication and multiplication

	if reverse {

		return MSM(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	coeffA, termA := SeparateTerm(a.GetRoot(), a)

	coeffB, termB := SeparateTerm(b.GetRoot(), b)

	if termB == termA {

		resultCoeff := coeffA + coeffB

		return CreateLikeTerm(resultCoeff, termA)

	} else {

		result.AppendSubtreeFrom(root, a.GetRoot(), a)

		result.AppendSubtreeFrom(root, b.GetRoot(), b)

		return result
	}
}

func MSD(a, b Expression, reverse bool) Expression { // adding multiplication and division

	if reverse {

		return MSD(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.GetChildren(a.GetRoot()), a)

	result.AppendSubtreeFrom(root, b.GetRoot(), b)

	return result
}

func DSD(a, b Expression, reverse bool) Expression { // adding division and division

	if reverse {

		return DSD(b, a, false)
	}
	childrenA := a.GetChildren(a.GetRoot())

	childrenB := b.GetChildren(b.GetRoot())

	root, result := NewExpression(NewOperation(Division))

	if IsEqualAt(childrenA[1], childrenB[1], &a, &b) {

		result.AppendExpression(root, Add(ExpressionIndex{Expression: a, Index: childrenA[0]}, ExpressionIndex{Expression: b, Index: childrenB[0]}), false)

		result.AppendSubtreeFrom(root, childrenA[1], a)

		// cancel if applicable

		return result

	} else {

		denomMul := Multiply(ExpressionIndex{Expression: a, Index: childrenA[1]}, ExpressionIndex{Expression: b, Index: childrenB[1]})

		numAMul := Multiply(ExpressionIndex{Expression: a, Index: childrenA[0]}, ExpressionIndex{Expression: b, Index: childrenB[1]})

		numBMul := Multiply(ExpressionIndex{Expression: a, Index: childrenA[1]}, ExpressionIndex{Expression: b, Index: childrenB[0]})

		result.AppendExpression(root, Add(ExpressionIndex{Expression: numAMul, Index: numAMul.GetRoot()}, ExpressionIndex{Expression: numBMul, Index: numBMul.GetRoot()}), false)

		result.AppendExpression(root, denomMul, false)

		// cancel if applicable

		return result
	}
}
