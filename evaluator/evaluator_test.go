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
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Eval Integer Expression Test Case %d", i), func(t *testing.T) {
			evaluated := testEval(c.input)
			testIntegerObject(t, evaluated, c.expected)
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
