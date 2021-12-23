package tests

import (
	"fmt"
	"symgolic/evaluation"
	"symgolic/symbols"
	"testing"
)

func TestFindAllFactors(t *testing.T) {

	factorsShort := evaluation.FindFactors(100)

	fmt.Println(factorsShort)

	factors := evaluation.GeneratePossibleOperandCombinationsForValue(12, 4, symbols.Multiplication)

	fmt.Println(factors)
}

func TestFindRoot(t *testing.T) {

	rootsA := evaluation.FindRoots(25)

	rootsB := evaluation.FindRoots(256)

	fmt.Println(rootsA)

	fmt.Println(rootsB)
}
