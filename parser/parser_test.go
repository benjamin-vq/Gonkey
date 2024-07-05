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
	checkParserErrors(t, p)

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

func TestReturnStatements(t *testing.T) {
	input := `
return -1;
return 5;
return 10;
return 993327;
`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("Incorrect amount of statements, got %d want 4",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("Statement is not a return statement, got %T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("Incorrect token literal, got %s want 'return'",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Not enough statements in program, got %d want 1", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program statement is not an expression, got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expression is not an identifier, got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("Incorrect identifier value, got %s want %s", ident.Value, "foobar")
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("Incorrect token literal for identifier, got %s want %s", ident.TokenLiteral(), "foobar")
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Not enough statements in program, got %d want 1", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program statement is not an expression, got %T", stmt)
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Statement is not an integer literal, got %T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("Incorrect literal value, got %d want 5", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("Incorrect token literal, got %s want %s", literal.TokenLiteral(), "5")
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) != 0 {
		t.Errorf("Parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("Parser error: %q", msg)
		}
		t.FailNow()
	}

	return
}
