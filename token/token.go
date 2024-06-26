package token

type TokenType int

type Token struct {
	Type    TokenType
	Literal string
}

var KEYWORDS = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// LookupIdent Figure out if the received identifier is a keyword or not
func LookupIdent(ident string) TokenType {

	if tt, valid := KEYWORDS[ident]; valid {
		return tt
	}

	return IDENT
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
		lit = "FUNCTION"
	case 19:
		lit = "LET"
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
	FUNCTION
	LET
)
