package tests

import (
	"fmt"
	"symgolic/language/components"
	"symgolic/language/interpretation"
	"testing"
)

func TestFindAllFactors(t *testing.T) {

	factorsShort := interpretation.FindFactors(100)

	fmt.Println(factorsShort)

	factors := interpretation.GeneratePossibleOperandCombinationsForValue(12, 4, components.Multiplication)

	fmt.Println(factors)
}

func TestFindRoot(t *testing.T) {

	rootsA := interpretation.FindRoots(25)

	rootsB := interpretation.FindRoots(256)

	fmt.Println(rootsA)

	fmt.Println(rootsB)
}
