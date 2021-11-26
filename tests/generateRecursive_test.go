package tests

import (
	"symgolic/generic"
	"testing"
)

func TestGenerateRearrangedArrays(t *testing.T) {

	// arr := []int{1, 2, 3}

	// combinations := generation.GeneratePermutationsOfArray(arr)

	// fmt.Println(combinations)
}

func TestGenerateCombinationsByRow(t *testing.T) {

	matrix := [][]int{

		[]int{0, 7, 6},

		[]int{3, 5, 6},

		[]int{1, 2, 5},
	}
	expectedByRow := [][]int{

		[]int{0, 3, 1},
		[]int{0, 3, 2},
		[]int{0, 3, 5},
		[]int{0, 5, 1},
		[]int{0, 5, 2},
		[]int{0, 5, 5},
		[]int{0, 6, 1},
		[]int{0, 6, 2},
		[]int{0, 6, 5},
		[]int{7, 3, 1},
		[]int{7, 3, 2},
		[]int{7, 3, 5},
		[]int{7, 5, 1},
		[]int{7, 5, 2},
		[]int{7, 5, 5},
		[]int{7, 6, 1},
		[]int{7, 6, 2},
		[]int{7, 6, 5},
		[]int{6, 3, 1},
		[]int{6, 3, 2},
		[]int{6, 3, 5},
		[]int{6, 5, 1},
		[]int{6, 5, 2},
		[]int{6, 5, 5},
		[]int{6, 6, 1},
		[]int{6, 6, 2},
		[]int{6, 6, 5},
	}
	combinationsByRow := generic.GenerateCombinationsByRow(matrix)

	if len(combinationsByRow) != len(expectedByRow) {

		t.Fatalf("Combinations do not match expected length")

	} else {

		for i := 0; i < len(combinationsByRow); i++ {

			if len(combinationsByRow[i]) != len(expectedByRow[i]) {

				t.Fatalf("Inner combinations do not match expected length")
			}
			for j := 0; j < len(combinationsByRow[i]); j++ {

				if combinationsByRow[i][j] != expectedByRow[i][j] {

					t.Fatalf("Inner combinations do not match")
				}
			}
		}
	}
}
