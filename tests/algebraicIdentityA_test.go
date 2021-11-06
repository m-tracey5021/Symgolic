package tests

import (
	"symgolic/comparison"
	"symgolic/identities"
	"symgolic/parsing"
	"testing"
)

func TestAlgebraicIdentityA(t *testing.T) {

	originals := map[string]string{
		"(a^2)+(2*a*b)+(b^2)": "(a+b)^2",
		"(a^2)+(6*a)+(3^2)":   "(a+3)^2",
		"(3^2)+(6*b)+(b^2)":   "(3+b)^2",
		"(2^2)+12+(3^2)":      "(2+3)^2",
		"(a+b)^2":             "(a^2)+(2*a*b)+(b^2)",
		"(2+(3*x))^2":         "(2^2)+(2*2*3*x)+((3*x)^2)",
		"a+b+c":               "a+b+c",
		"(1/2)+(3*x)":         "(1/2)+(3*x)",
		"(a^2)+(2*a*b)+(c^2)": "(a^2)+(2*a*b)+(c^2)",
	}

	for input, output := range originals {

		original := parsing.ParseExpression(input)

		expected := parsing.ParseExpression(output)

		identityA := identities.NewAlgebraicIdentityA(&original)

		_, result := identities.Run(original.GetRoot(), &original, &identityA)

		if !comparison.IsEqual(result, expected) {

			err := "Expected (a+b)^2 but instead got " + result.ToString()

			t.Fatalf(err)
		}
	}
}

func TestAlgebraicIdentityD(t *testing.T) {

	originals := map[string]string{
		"(x^2)+((a+b)*x)+(a*b)": "(x+a)*(x+b)",
		"(6^2)+12+2":            "(6+1)*(6+1)",
	}

	for input, output := range originals {

		original := parsing.ParseExpression(input)

		expected := parsing.ParseExpression(output)

		identityA := identities.NewAlgebraicIdentityA(&original)

		_, result := identities.Run(original.GetRoot(), &original, &identityA)

		if !comparison.IsEqual(result, expected) {

			err := "Expected (a+b)^2 but instead got " + result.ToString()

			t.Fatalf(err)
		}
	}
}
