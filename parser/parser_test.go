package parser

import (
	"github.com/benja-vq/gonkey/ast"
	"github.com/benja-vq/gonkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 272727;
`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("Parsed program was nil")
	}

	if len(program.Statements) != 3 {
		t.Errorf("Incorrect amount of program statements, got %d want 3",
			len(program.Statements))
	}

	cases := []struct {
		expectedIdentifier string
	}{
		{expectedIdentifier: "x"},
		{expectedIdentifier: "y"},
		{expectedIdentifier: "foobar"},
	}

	for i, c := range cases {
		statement := program.Statements[i]

		if statement.TokenLiteral() != "let" {
			t.Errorf("Incorrect token literal for statement, got %s want let",
				statement.TokenLiteral())
		}

		// Type assertion from a Statement to a LetStatement
		let, ok := statement.(*ast.LetStatement)
		if !ok {
			t.Errorf("Statement is not a let statement, got %T", statement)
		}

		if let.Name.Value != c.expectedIdentifier {
			t.Errorf("Incorrect identifier, got %s want %s", let.Name.Value, c.expectedIdentifier)
		}

		if let.Name.TokenLiteral() != c.expectedIdentifier {
			t.Errorf("Incorrect token literal, got %s want %s", let.Name.TokenLiteral(), c.expectedIdentifier)
		}
	}
}
