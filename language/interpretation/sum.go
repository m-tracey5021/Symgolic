package interpretation

import (
	. "symgolic/language/components"
)

func SSS(a, b ExpressionIndex, reverse bool) Expression { // adding two sums

	if reverse {

		return SSS(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.Expression.GetChildren(a.Index), a.Expression)

	result.AppendBulkSubtreesFrom(root, b.Expression.GetChildren(b.Index), b.Expression)

	return result
}

func SSM(a, b ExpressionIndex, reverse bool) Expression { // adding sum and multiplication

	if reverse {

		return SSM(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.Expression.GetChildren(a.Index), a.Expression)

	result.AppendSubtreeFrom(root, b.Index, b.Expression)

	return result
}

func SSD(a, b ExpressionIndex, reverse bool) Expression { // adding sum and division

	if reverse {

		return SSD(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.Expression.GetChildren(a.Index), a.Expression)

	result.AppendSubtreeFrom(root, b.Index, b.Expression)

	return result

}

func MSM(a, b ExpressionIndex, reverse bool) Expression { // adding multiplication and multiplication

	if reverse {

		return MSM(b, a, false)
	}

	isLikeTerm, coeff, terms := IsLikeTerm(a, b)

	if isLikeTerm {

		root, result := NewExpression(NewOperation(Multiplication))

		_, coeffExp := NewExpression(NewConstant(coeff))

		result.AppendExpression(root, coeffExp, false)

		result.AppendBulkExpressions(root, terms)

		return result

	} else {

		root, result := NewExpression(NewOperation(Addition))

		result.AppendSubtreeFrom(root, a.Index, a.Expression)

		result.AppendSubtreeFrom(root, b.Index, b.Expression)

		return result
	}
}

func MSD(a, b ExpressionIndex, reverse bool) Expression { // adding multiplication and division

	if reverse {

		return MSD(b, a, false)
	}
	root, result := NewExpression(NewOperation(Addition))

	result.AppendBulkSubtreesFrom(root, a.Expression.GetChildren(a.Index), a.Expression)

	result.AppendSubtreeFrom(root, b.Expression.GetRoot(), b.Expression)

	return result
}

func DSD(a, b ExpressionIndex, reverse bool) Expression { // adding division and division

	if reverse {

		return DSD(b, a, false)
	}
	childrenA := a.Expression.GetChildren(a.Index)

	childrenB := b.Expression.GetChildren(b.Index)

	root, result := NewExpression(NewOperation(Division))

	if IsEqualAt(a.At(childrenA[1]), b.At(childrenB[1])) {

		result.AppendExpression(root, Add(a.At(childrenA[0]), b.At(childrenB[0])), false)

		result.AppendSubtreeFrom(root, childrenA[1], a.Expression)

		// cancel if applicable

		return result

	} else {

		denomMul := Multiply(a.At(childrenA[1]), b.At(childrenB[1]))

		numAMul := Multiply(a.At(childrenA[0]), b.At(childrenB[1]))

		numBMul := Multiply(a.At(childrenA[1]), b.At(childrenB[0]))

		result.AppendExpression(root, Add(From(numAMul), From(numBMul)), false)

		result.AppendExpression(root, denomMul, false)

		// cancel if applicable

		return result
	}
}
