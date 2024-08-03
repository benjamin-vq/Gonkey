package token

type TokenType int

type Token struct {
	Type    TokenType
	Literal string
}

var KEYWORDS = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent Figure out if the received identifier is a keyword or not
func LookupIdent(ident string) TokenType {

	if tt, valid := KEYWORDS[ident]; valid {
		return tt
	}

	return IDENT
}

func (tt TokenType) String() string {
	return tt.Literal()
}

func (tt TokenType) Literal() (lit string) {

	switch tt {
	case 0:
		lit = "" // EOF
	case 1:
		lit = "ILLEGAL"
	case 2:
		lit = "IDENT"
	case 3:
		lit = "INT"
	case 4:
		lit = "="
	case 5:
		lit = "+"
	case 6:
		lit = "-"
	case 7:
		lit = "!"
	case 8:
		lit = "*"
	case 9:
		lit = "/"
	case 10:
		lit = "<"
	case 11:
		lit = ">"
	case 12:
		lit = ","
	case 13:
		lit = ";"
	case 14:
		lit = "("
	case 15:
		lit = ")"
	case 16:
		lit = "{"
	case 17:
		lit = "}"
	case 18:
		lit = "=="
	case 19:
		lit = "!="
	case 20:
		lit = "FUNCTION"
	case 21:
		lit = "LET"
	case 22:
		lit = "TRUE"
	case 23:
		lit = "FALSE"
	case 24:
		lit = "IF"
	case 25:
		lit = "ELSE"
	case 26:
		lit = "RETURN"
	case 27:
		lit = "STRING"
	case 28:
		lit = "["
	case 29:
		lit = "]"
	}
	return lit
}

const (
	EOF TokenType = iota
	ILLEGAL

	IDENT
	INT

	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH

	LT
	GT

	COMMA
	SEMICOLON

	LPAREN
	RPAREN
	LBRACE
	RBRACE

	EQ
	NOT_EQ

	// Keywords
	FUNCTION
	LET
	TRUE
	FALSE
	IF
	ELSE
	RETURN

	STRING
	LBRACKET
	RBRACKET
)
