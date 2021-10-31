package identities

import (
	. "symgolic/symbols"
)

type Identity interface {
	Pass()

	Failure()

	Identify()

	Apply()

	Run()
}

type Stage struct {
	SymbolType SymbolType

	Recurse bool
}
