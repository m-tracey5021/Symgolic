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

		12: "(a+b)*c", // a = 1, b = 2, c = 4
		2:  "a*b",
	}
	resultA := identities.GenerateCompatibleConstantMapsForValues(valuesA)

	resultB := identities.GenerateCompatibleConstantMapsForValues(valuesB)

	results := [][]map[string]int{resultA, resultB}

	expectedA := []map[string]int{

		map[string]int{
			"a": 3,
			"b": 3,
		},
	}
	expectedB := []map[string]int{

		map[string]int{
			"a": 2,
			"b": 1,
			"c": 4,
		},
		map[string]int{
			"a": 1,
			"b": 2,
			"c": 4,
		},
	}
	expectedResults := [][]map[string]int{expectedA, expectedB}

	for i := 0; i < len(expectedResults); i++ {

		expected := expectedResults[i]

		actual := results[i]

		if len(expected) == len(actual) {

			for j := 0; j < len(expected); j++ {

				if !identities.MappingsAreEqual(expected[j], actual[j]) {

					err := fmt.Sprintf("Mappings are not equal, expected %v, got %v", expected[j], actual[j])

					t.Fatalf(err)
				}
			}
		}
	}
}
