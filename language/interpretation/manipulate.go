package interpretation

import (
	. "symgolic/language/components"
)

type Manipulation func(ExpressionIndex) Expression

type ManipulationInPlace func(ExpressionIndex)

type ManipulationAgainst func(ExpressionIndex, ExpressionIndex) Expression

type ManipulationForMany func(...ExpressionIndex) Expression
