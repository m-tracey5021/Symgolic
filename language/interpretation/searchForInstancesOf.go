package interpretation

import (
	. "symgolic/language/components"
)

func SearchForInstancesOf(a, b ExpressionIndex, instances []int) []int {

	if IsEqualAt(a, b) {

		instances = append(instances, b.Index)
	}
	for _, child := range b.Expression.GetChildren(b.Index) {

		instances = SearchForInstancesOf(a, b.At(child), instances)
	}
	return instances
}
