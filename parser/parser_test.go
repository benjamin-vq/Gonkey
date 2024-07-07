package parser

import (
	"fmt"
	"github.com/benja-vq/gonkey/ast"
	"github.com/benja-vq/gonkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {

	cases := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Let Statement Test %d, Input: %s", i, c.input), func(t *testing.T) {
			l := lexer.NewLexer(c.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Errorf("Incorrect amount of program statements, got %d want %d",
					len(program.Statements), 1)
			}

			stmt := program.Statements[0]
			if !testLetStatement(t, stmt, c.expectedIdentifier) {
				return
			}

			let := stmt.(*ast.LetStatement)
			if let.Name.Value != c.expectedIdentifier {
				t.Errorf("Incorrect identifier, got %s want %s", let.Name.Value, c.expectedIdentifier)
			}

			if let.Name.TokenLiteral() != c.expectedIdentifier {
				t.Errorf("Incorrect token literal, got %s want %s", let.Name.TokenLiteral(), c.expectedIdentifier)
			}
		})
	}
}

func TestReturnStatements(t *testing.T) {

	cases := []struct {
		input         string
		expectedValue any
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Return Statement Test %d, Input %s", i, c.input), func(t *testing.T) {
			l := lexer.NewLexer(c.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("Incorrect amount of statements, got %d want %d",
					len(program.Statements), 1)
			}

			stmt := program.Statements[0]
			returnStmt, ok := stmt.(*ast.ReturnStatement)
			if !ok {
				t.Fatalf("Statement is not a return statement, got %T", stmt)
			}

			if returnStmt.TokenLiteral() != "return" {
				t.Fatalf("Incorrect return statement token literal, got %s want %s",
					returnStmt.TokenLiteral(), "return")
			}
		})
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

func TestParsingPrefixExpressions(t *testing.T) {
	cases := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Prefix Expression Test %d, Input %s", i, c.input), func(t *testing.T) {
			l := lexer.NewLexer(c.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("Incorrect amount of statements, got %d want %d",
					len(program.Statements), 1)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Statement is not an expression statement, got %T", stmt)
			}

			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("Expression is not a prefix expression, got %T", exp)
			}
			if exp.Operator != c.operator {
				t.Errorf("Incorrect expression operator, got %s want %s", exp.Operator, c.operator)
			}

			if !testLiteralExpression(t, exp.Right, c.value) {
				return
			}
		})
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	cases := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Infix Expression Test %d, Input %s", i, c.input), func(t *testing.T) {
			l := lexer.NewLexer(c.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("Incorrect amount of statements, got %d want %d",
					len(program.Statements), 1)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Statement is not an expression statement, got %T", stmt)
			}

			if !testInfixExpression(t, stmt.Expression, c.leftValue,
				c.operator, c.rightValue) {
				return
			}
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Operator Precedence Test Case %d", i), func(t *testing.T) {

			l := lexer.NewLexer(c.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			actual := program.String()
			if actual != c.expected {
				t.Errorf("Incorrect program, got %q want %q", actual, c.expected)
			}
		})
	}
}

func TestBooleanExpression(t *testing.T) {
	cases := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Boolean Expression Test %d", i), func(t *testing.T) {
			l := lexer.NewLexer(c.input)
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("Incorrect amount of program statements, got %d want %d",
					len(program.Statements), 1)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Program statement is not an expression statement, got %T",
					program.Statements[0])
			}

			boolean, ok := stmt.Expression.(*ast.BooleanLiteral)
			if !ok {
				t.Fatalf("Expression is not a boolean literal, got %T", stmt.Expression)
			}
			if boolean.Value != c.expectedBoolean {
				t.Errorf("Incorrect boolean value, got %t want %t", boolean.Value, c.expectedBoolean)
			}
		})
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {

	if stmt.TokenLiteral() != "let" {
		t.Errorf("Incorrect token literal, got %s want %s", stmt.TokenLiteral(), "let")
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement is not a let statement, got %T", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("Incorrect let statement name, got %s want %s", letStmt.Name.Value, name)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("Incorrect let statement name token literal, got %s want %s",
			letStmt.Name.TokenLiteral(), name)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {

	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("Unhandled expression type: %T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bl, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("Expression is not a boolean literal, got %T", exp)
		return false
	}

	if bl.Value != value {
		t.Errorf("Incorrect boolean literal value, got %t want %t", bl.Value, value)
		return false
	}

	if bl.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("Incorrect boolean token literal, got %s want %t", bl.TokenLiteral(), value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {

	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Expression is not an identifier, got %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("Incorrect identifier value, got %s want %s", ident.Value, value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("Incorrect identifier token literal, got %s want %s",
			ident.TokenLiteral(), value)
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {

	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Expression is not an infix expression, got %T", exp)
		return false
	}

	if !testLiteralExpression(t, infixExp.Left, left) {
		return false
	}

	if infixExp.Operator != operator {
		t.Errorf("Incorrect infix expression operator, got %s want %s", operator, infixExp.Operator)
		return false
	}

	if !testLiteralExpression(t, infixExp.Right, right) {
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expression is not an integer literal, got %T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("Incorrect integer expression value, got %d want %d", integer.Value, value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("Incorrect token literal, got %s want %d",
			integer.TokenLiteral(), value)
		return false
	}

	return true
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
