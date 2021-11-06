package tests

import (
	"fmt"
	"symgolic/generation"
	"testing"
)

func TestGenerateCombinations(t *testing.T) {

	terms := [][]int{

		[]int{0, 7, 6},

		[]int{3, 5, 6},

		[]int{1, 2, 5},
	}

	// combinations := make([][]int, 0)

	combinations := generation.GenerateCombinations(terms, make([][]interface{}, 0), make([]interface{}, 0), 3)

	fmt.Println(combinations)
}
