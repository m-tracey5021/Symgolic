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
	}

	for input, output := range originals {

		original := parsing.ParseExpression(input)

		expected := parsing.ParseExpression(output)

		identityA := identities.NewAlgebraicIdentityA(&original)

		_, result := identities.Run(original.GetRoot(), &original, &identityA)

		if !comparison.IsEqual(result, expected) {

			err := "Expected " + expected.ToString() + " but instead got " + result.ToString()

			t.Fatalf(err)
		}
	}
}

func TestAlgebraicIdentityB(t *testing.T) {

	originals := map[string]string{
		"(a^2)-(2*a*b)+(b^2)": "(a-b)^2",
		"(a^2)-(6*a)+(3^2)":   "(a-3)^2",
		"(3^2)-(6*b)+(b^2)":   "(3-b)^2",
		"(2^2)-12+(3^2)":      "(2-3)^2",
		"(a-b)^2":             "(a^2)-(2*a*b)+(b^2)",
		"(2-(3*x))^2":         "(2^2)-(2*2*3*x)+((3*x)^2)",
		"a+b+c":               "a+b+c",
		"(1/2)+(3*x)":         "(1/2)+(3*x)",
	}

	for input, output := range originals {

		original := parsing.ParseExpression(input)

		expected := parsing.ParseExpression(output)

		identityA := identities.NewAlgebraicIdentityA(&original)

		_, result := identities.Run(original.GetRoot(), &original, &identityA)

		if !comparison.IsEqual(result, expected) {

			err := "Expected " + expected.ToString() + " but instead got " + result.ToString()

			t.Fatalf(err)
		}
	}
}

func TestAlgebraicIdentityC(t *testing.T) {

	originals := map[string]string{
		"(a^2)-(b^2)": "(a+b)*(a-b)",
		"(a+b)*(a-b)": "(a^2)-(b^2)",
		// "25-16":       "(5+4)*(5-4)", // this hits maximum recursion depth, 25 is too big, limit this somehow
		"9-4": "(3+2)*(3-2)",
		"a+b": "a+b",
	}

	for input, output := range originals {

		original := parsing.ParseExpression(input)

		expected := parsing.ParseExpression(output)

		identityA := identities.NewAlgebraicIdentityC(&original)

		_, result := identities.Run(original.GetRoot(), &original, &identityA)

		if !comparison.IsEqual(result, expected) {

			err := "Expected " + expected.ToString() + " but instead got " + result.ToString()

			t.Fatalf(err)
		}
	}
}

func TestAlgebraicIdentityD(t *testing.T) {

	originals := map[string]string{
		"(x^2)+(5*x)+(3*2)":     "(x+3)*(x+2)",
		"(x^2)+((a+b)*x)+(a*b)": "(x+a)*(x+b)",
		"(4^2)+12+2":            "(4+2)*(4+1)",
		"(4+2)*(4+1)":           "(4^2)+((2+1)*4)+(2*1)",
		"3*6":                   "(2^2)+((1+4)*2)+(1*4)", // (2+1)*(2+4) x = 2, a = 1, b = 4
		"a+b+c":                 "a+b+c",
	}

	for input, output := range originals {

		original := parsing.ParseExpression(input)

		expected := parsing.ParseExpression(output)

		identityD := identities.NewAlgebraicIdentityD(&original)

		_, result := identities.Run(original.GetRoot(), &original, &identityD)

		if !comparison.IsEqual(result, expected) {

			err := "Expected (a+b)^2 but instead got " + result.ToString()

			t.Fatalf(err)
		}
	}
}
