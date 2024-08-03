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

func TestReturnStatements(t *testing.T) {
	cases := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}`, 10},
		{
			`
let f = fn(x) {
  return x;
  x + 10;
};
f(10);`,
			10,
		},
		{
			`
let f = fn(x) {
   let result = x + 10;
   return result;
   return 10;
};
f(10);`,
			20,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Return Statement Test Case %d", i), func(t *testing.T) {
			evaluated := testEval(c.input)
			testIntegerObject(t, evaluated, c.expected)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	cases := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}`, "unknown operator: BOOLEAN + BOOLEAN"},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Error Handling Test Case %d", i), func(t *testing.T) {
			evaluated := testEval(c.input)

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("Object is not an object error, got %T (%+v)",
					evaluated, evaluated)
				return
			}

			if errObj.Message != c.expectedMessage {
				t.Errorf("Incorrect error message, got %s want %s",
					errObj.Message, c.expectedMessage)
			}
		})
	}
}

func TestLetStatements(t *testing.T) {
	cases := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Let Statement Test Case %d", i), func(t *testing.T) {
			testIntegerObject(t, testEval(c.input), c.expected)
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("Object is not a function, got %T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("Incorrect amount of function parameters, got %d want %d",
			len(fn.Parameters), 1)
	}

	if fn.Parameters[0].String() != "x" {
		t.Errorf("Incorrect function parameter, got %s want %s",
			fn.Parameters[0].String(), "x")
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Errorf("Incorrect function body, got %s want %s",
			fn.Body.String(), expectedBody)
	}
}

func TestFunctionApplication(t *testing.T) {
	cases := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Function Application Test Case %d", i), func(t *testing.T) {
			testIntegerObject(t, testEval(c.input), c.expected)
		})
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) { x + y };
};

let addTwo = newAdder(2);
addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("Object is not a String object, got %T (%+v)",
			evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("Incorrect string object value, got %s want %s",
			str.Value, "Hello World!")
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("Object is not a String object, got %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("Incorrect value from string concatenation, got %s want %s",
			str.Value, "Hello World!")
	}
}

func TestBuiltinFunctions(t *testing.T) {
	cases := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to 'len' not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments, got 2 want 1"},
	}

	for _, c := range cases {
		eval := testEval(c.input)

		switch expected := c.expected.(type) {
		case int:
			testIntegerObject(t, eval, int64(expected))
		case string:
			err, ok := eval.(*object.Error)
			if !ok {
				t.Errorf("Object is not an error object, got %T (%+v)",
					err, err)
				continue
			}

			if err.Message != c.expected {
				t.Errorf(`Incorrect error message. 
Got: %q 
Want: %q`,
					err.Message, c.expected)
			}

		}
	}
}

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
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
