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

func TestEvaluateArithmetic(t *testing.T) {

	testData := []BinaryTestData{

		{Input: "(2*y)+(3*y)", Expected: "5*y"},
		{Input: "((2*y)+(4*y))+(3*y)", Expected: "9*y"},
	}
	result, err := TestBinaryDataOverEvaluation(testData, interpretation.EvaluateArithmetic)

	if !result {

		t.Fatalf(err)
	}
}

func TestSum(t *testing.T) {

	testData := []TernaryTestData{

		{InputA: "2*y", InputB: "3*y", Expected: "5*y"},
		{InputA: "(2*y)+(4*y)", InputB: "3*y", Expected: "9*y"},
	}
	result, err := TestTernaryDataOverManipulationForMany(testData, interpretation.Add)

	if !result {

		t.Fatalf(err)
	}
}
