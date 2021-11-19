package tests

import (
	"symgolic/comparison"
	"symgolic/parsing"
	"testing"
)

func TestInsertExpression(t *testing.T) {

	expression := parsing.ParseExpression("a+b+c")

	toInsert := parsing.ParseExpression("3/2")

	expression.InsertExpression(expression.GetRoot(), 1, toInsert)

	result := parsing.ParseExpression("a+(3/2)+b+c")

	if !comparison.IsEqual(result, expression) {

		err := "Expected " + result.ToString() + ", but got " + expression.ToString()

		t.Fatalf(err)
	}
}

func TestReplaceExpression(t *testing.T) {

	expression := parsing.ParseExpression("a+b+c")

	toReplace := parsing.ParseExpression("3/2")

	b := expression.GetNodeByPath([]int{1})

	expression.ReplaceNodeCascade(b, toReplace)

	result := parsing.ParseExpression("a+(3/2)+c")

	if !comparison.IsEqual(result, expression) {

		err := "Expected " + result.ToString() + ", but got " + expression.ToString()

		t.Fatalf(err)
	}

	expressionB := parsing.ParseExpression("a+b+c")

	toReplaceB := parsing.ParseExpression("d+e")

	bB := expressionB.GetNodeByPath([]int{1})

	expressionB.ReplaceNodeCascade(bB, toReplaceB)

	resultB := parsing.ParseExpression("a+d+e+c")

	if !comparison.IsEqual(resultB, expressionB) {

		err := "Expected " + resultB.ToString() + ", but got " + expressionB.ToString()

		t.Fatalf(err)
	}
}
