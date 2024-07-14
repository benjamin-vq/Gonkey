package evaluator

import (
	"fmt"
	"github.com/benja-vq/gonkey/lexer"
	"github.com/benja-vq/gonkey/object"
	"github.com/benja-vq/gonkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	cases := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"27", 27},
		{"-5", -5},
		{"-27", -27},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Eval Integer Expression Test Case %d", i), func(t *testing.T) {
			evaluated := testEval(c.input)
			testIntegerObject(t, evaluated, c.expected)
		})
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Eval Boolean Expression Test Case %d", i), func(t *testing.T) {
			evaluated := testEval(c.input)
			testBooleanObject(t, evaluated, c.expected)
		})
	}
}

func TestBangOperator(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!27", false},
		{"!!true", true},
		{"!!false", false},
		{"!!27", true},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Bang Operator Test Case %d", i), func(t *testing.T) {
			evaluated := testEval(c.input)
			testBooleanObject(t, evaluated, c.expected)
		})
	}
}

func TestIfElseExpressions(t *testing.T) {
	cases := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("If Else Expression Test Case %d", i), func(t *testing.T) {
			evaluated := testEval(c.input)
			integer, ok := c.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Object is not an integer object, got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Incorrect object value, got %d want %d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Object is not a boolean object, got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("Incorrect object value, got %t want %t", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Object is not NULL, got %T (%+v)", obj, obj)
		return false
	}

	return true
}