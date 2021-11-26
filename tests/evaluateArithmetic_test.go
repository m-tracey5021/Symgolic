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
