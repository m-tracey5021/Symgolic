package search

import (
	"symgolic/comparison"
	. "symgolic/symbols"
)

func SearchForInstancesOf(target, compared int, expression, other Expression, instances []int) []int {

	if comparison.IsEqual(target, compared, &expression, &other) {

		instances = append(instances, compared)
	}
	for _, child := range other.GetChildren(compared) {

		instances = SearchForInstancesOf(target, child, expression, other, instances)
	}
	return instances
}
