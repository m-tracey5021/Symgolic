package tests

import (
	"strconv"
	"symgolic/comparison"
	"symgolic/parsing"
	"testing"
)

type ComparisonTestData struct {
	first, second string

	expected bool
}

func TestIsEqual(t *testing.T) {

	data := []ComparisonTestData{
		{first: "a+b", second: "a+b", expected: true},
		{first: "a+b", second: "b+a", expected: true},
		{first: "a+(b*c)", second: "a+(b*c)", expected: true},
		{first: "a+(b*c)", second: "a+(c*b)", expected: true},
		{first: "(a*b)+(c*d)", second: "(a*b)+(c*d)", expected: true},
		{first: "(a*b)+(c*d)", second: "(d*c)+(b*a)", expected: true},
		{first: "(a*b)+(c*d)", second: "(a*c)+(b*d)", expected: false},
		{first: "a+b", second: "a+c", expected: false},
		{first: "a+(b*c)", second: "a+(b/c)", expected: false},
	}
	for _, input := range data {

		first := parsing.ParseExpression(input.first)

		second := parsing.ParseExpression(input.second)

		result := comparison.IsEqual(first, second)

		if result != input.expected {

			t.Fatalf("Expected equality of " + first.ToString() + " and " + second.ToString() + " to be: " + strconv.FormatBool(input.expected) + ", but got: " + strconv.FormatBool(result))
		}
	}
}

func TestIsEqualByBase(t *testing.T) {

	data := []ComparisonTestData{
		{first: "a^2", second: "a^3", expected: true},
		{first: "a^2", second: "b^3", expected: false},
		{first: "(a+b)^2", second: "(b+a)^3", expected: true},
		{first: "a+(b*c)", second: "a+(b*c)", expected: true},
		{first: "(a+(b*c))^(x+y)", second: "((b*c)+a)^(x+y)", expected: true},
		{first: "(a+(b*c))^(x+y)", second: "((a*b)+c)^(x+y)", expected: false},
	}
	for _, input := range data {

		first := parsing.ParseExpression(input.first)

		second := parsing.ParseExpression(input.second)

		result := comparison.IsEqualByBase(first, second)

		if result != input.expected {

			t.Fatalf("Expected equality of " + first.ToString() + " and " + second.ToString() + " to be: " + strconv.FormatBool(input.expected) + ", but got: " + strconv.FormatBool(result))
		}
	}
}
