package lexer

import (
	"github.com/benja-vq/gonkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
}

let result = add(five, ten);
`

	cases := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{expectedType: token.LET, expectedLiteral: "let"},
		{expectedType: token.IDENT, expectedLiteral: "five"},
		{expectedType: token.ASSIGN, expectedLiteral: "="},
		{expectedType: token.INT, expectedLiteral: "5"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.LET, expectedLiteral: "let"},
		{expectedType: token.IDENT, expectedLiteral: "ten"},
		{expectedType: token.ASSIGN, expectedLiteral: "="},
		{expectedType: token.INT, expectedLiteral: "10"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.LET, expectedLiteral: "let"},
		{expectedType: token.IDENT, expectedLiteral: "add"},
		{expectedType: token.ASSIGN, expectedLiteral: "="},
		{expectedType: token.FUNCTION, expectedLiteral: "fn"},
		{expectedType: token.LPAREN, expectedLiteral: "("},
		{expectedType: token.IDENT, expectedLiteral: "x"},
		{expectedType: token.COMMA, expectedLiteral: ","},
		{expectedType: token.IDENT, expectedLiteral: "y"},
		{expectedType: token.RPAREN, expectedLiteral: ")"},
		{expectedType: token.LBRACE, expectedLiteral: "{"},
		{expectedType: token.IDENT, expectedLiteral: "x"},
		{expectedType: token.PLUS, expectedLiteral: "+"},
		{expectedType: token.IDENT, expectedLiteral: "y"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.RBRACE, expectedLiteral: "}"},
		{expectedType: token.LET, expectedLiteral: "let"},
		{expectedType: token.IDENT, expectedLiteral: "result"},
		{expectedType: token.ASSIGN, expectedLiteral: "="},
		{expectedType: token.IDENT, expectedLiteral: "add"},
		{expectedType: token.LPAREN, expectedLiteral: "("},
		{expectedType: token.IDENT, expectedLiteral: "five"},
		{expectedType: token.COMMA, expectedLiteral: ","},
		{expectedType: token.IDENT, expectedLiteral: "ten"},
		{expectedType: token.RPAREN, expectedLiteral: ")"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.EOF, expectedLiteral: ""},
	}

	lexer := NewLexer(input)

	read := []string{}
	for i, tt := range cases {
		tok := lexer.NextToken()

		read = append(read, tok.Literal)
		if tok.Type != tt.expectedType {
			t.Fatalf("Test case %d (token type) failed. got %q want %q. Position = %d - Read Position = %d - Consumed Input: %v",
				i, tok.Type.Literal(), tt.expectedType.Literal(), lexer.position, lexer.readPosition, read)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("Test case %d (literal) failed. got %q want %q. Position = %d - Read Position = %d - Consumed Input: %v",
				i, tok.Literal, tt.expectedLiteral, lexer.position, lexer.readPosition, read)
		}
	}
}
