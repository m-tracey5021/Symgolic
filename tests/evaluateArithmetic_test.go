package tests

import (
	"symgolic/evaluation"
	"symgolic/parsing"
	"testing"
)

type Target struct {
	Expression string

	Target int
}

func TestFindVariablesWhere(t *testing.T) {

	targets := map[Target][]map[string]int{

		Target{"(a+b)*c", 12}: []map[string]int{

			map[string]int{
				"a": 1,
				"b": 1,
				"c": 6,
			},
			map[string]int{
				"a": 1,
				"b": 2,
				"c": 4,
			},
			map[string]int{
				"a": 1,
				"b": 3,
				"c": 3,
			},
			map[string]int{
				"a": 2,
				"b": 2,
				"c": 3,
			},
		},
	}

	for input, output := range targets {

		parsed := parsing.ParseExpression(input.Expression)

		variableMaps := evaluation.FindVariablesWhere(parsed.GetRoot(), &parsed, input.Target)

		if len(variableMaps) != len(output) {

			err := "Differing number of variable maps produced"

			t.Fatalf(err)

		} else {

			for i := 0; i < len(output); i++ {

			}
		}

	}
}
