package evaluator

import (
	"fmt"
	"github.com/benja-vq/gonkey/object"
	"testing"
)

func TestQuote(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{`quote(5)`, `5`},
		{`quote(5 + 8)`, `(5 + 8)`},
		{`quote(foobar)`, `foobar`},
		{`quote(foobar + barfoo)`, `(foobar + barfoo)`},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Quote Test Case %d", i), func(t *testing.T) {
			evaluated := testEval(c.input)
			quote, ok := evaluated.(*object.Quote)
			if !ok {
				t.Fatalf("Evaluated object is not a Quote object, got %T (%+v)",
					evaluated, evaluated)
			}

			if quote.Node == nil {
				t.Fatalf("quote.Node is nil")
			}

			if quote.Node.String() != c.expected {
				t.Errorf("Incorrect node string, got %s want %s",
					quote.Node.String(), c.expected)
			}
		})
	}
}

func TestQuoteUnquote(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{`quote(unquote(4))`, `4`},
		{`quote(unquote(4 + 4))`, `4`},
		{`quote(8 + unquote(4 + 4))`, `(8 + 8)`},
		{`quote(unquote(4 + 4) + 8)`, `(8 + 8)`},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Quote Unquote Test Case %d", i), func(t *testing.T) {
			evaluated := testEval(c.input)
			quote, ok := evaluated.(*object.Quote)
			if !ok {
				t.Fatalf("Evaluated object is not a quote object, got %T (%+v)",
					evaluated, evaluated)
			}

			if quote.Node == nil {
				t.Fatalf("quote.Node is nil")
			}

			if quote.Node.String() != c.expected {
				t.Errorf("Incorrect node string, got %s want %s",
					quote.Node.String(), c.expected)
			}
		})
	}
}
