package tests

import (
	"fmt"
	"symgolic/identities"
	"testing"
)

func TestGenerateCompatibleConstantMapsForValues(t *testing.T) {

	valuesA := map[int]string{

		6: "a+b",
		9: "a*b",
	}
	valuesB := map[int]string{

		12: "(a+b)*c", // a = 1, b = 3, c = 4
		3:  "a*b",
	}
	resultA := identities.GenerateCompatibleConstantMapsForValues(valuesA)

	resultB := identities.GenerateCompatibleConstantMapsForValues(valuesB)

	fmt.Println(resultA)

	fmt.Println(resultB)
}
