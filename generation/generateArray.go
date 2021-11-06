package generation

import (
	"reflect"
)

func GenerateCombinations(terms interface{}) interface{} {

	termsType := reflect.TypeOf(terms).String()

	if termsType != reflect.Array.String() {

		return nil
	}

}

func GenerateCombinationsRecurse(terms, combinations [][]interface{}, accumulated []interface{}, limit int) [][]interface{} {

	last := len(terms) == 1

	n := len(terms[0])

	for i := 0; i < n; i++ {

		accumulated = append(accumulated, terms[0][i])

		item := accumulated

		if last {

			if len(item) == limit {

				combinations = append(combinations, item)
			}

		} else {

			combinations = GenerateCombinationsRecurse(terms[1:], combinations, item, limit)
		}
	}
	return combinations
}
