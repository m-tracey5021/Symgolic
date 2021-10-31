package tests

import (
	"symgolic/comparison"
	"symgolic/identities"
	"symgolic/parsing"
	"testing"
)

func TestAlgebraicIdentityA(t *testing.T) {

	original := parsing.ParseExpression("(a^2)+(2*a*b)+(b^2)")

	expected := parsing.ParseExpression("(a+b)^2")

	identityA := identities.NewAlgebraicIdentityA()

	change, result := identityA.Run(original.GetRoot(), &original)

	if !change {

		t.Fatalf("IdentityA should apply to this expression")

	} else {

		if !comparison.IsEqual(result, expected) {

			err := "Expected (a+b)^2 but instead got " + result.ToString()

			t.Fatalf(err)
		}
	}

}
