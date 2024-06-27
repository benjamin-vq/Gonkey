package lexer

import "github.com/benja-vq/gonkey/token"

type Lexer struct {
	input        string
	position     int  // Current position in input
	readPosition int  // Current reading position in input
	char         byte // Character under examination. TODO: Support UTF-8
}

func NewLexer(input string) *Lexer {

	lexer := Lexer{
		input:        input,
		position:     0,
		readPosition: 0,
	}

	lexer.readChar()
	return &lexer
}

func (l *Lexer) NextToken() (tok token.Token) {

	l.skipWhitespace()

	switch l.char {
	// ASCII
	case 0:
		tok = newToken(token.EOF)
	case 33:
		if l.peekChar() == '=' {
			tok = newToken(token.NOT_EQ)
			l.readChar()
		} else {
			tok = newToken(token.BANG)
		}
	case 40:
		tok = newToken(token.LPAREN)
	case 41:
		tok = newToken(token.RPAREN)
	case 42:
		tok = newToken(token.ASTERISK)
	case 43:
		tok = newToken(token.PLUS)
	case 44:
		tok = newToken(token.COMMA)
	case 45:
		tok = newToken(token.MINUS)
	case 47:
		tok = newToken(token.SLASH)
	case 59:
		tok = newToken(token.SEMICOLON)
	case 60:
		tok = newToken(token.LT)
	case 61:
		if l.peekChar() == '=' {
			tok = newToken(token.EQ)
			l.readChar()
		} else {
			tok = newToken(token.ASSIGN)
		}
	case 62:
		tok = newToken(token.GT)
	case 123:
		tok = newToken(token.LBRACE)
	case 125:
		tok = newToken(token.RBRACE)
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL)
		}
	}

	l.readChar()
	return tok
}

// Read until we encounter a non-letter character, return the read chunk
// let myNumber = 5;
// ^~~
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.char) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// Read until we encounter a non-digit character, return the read number
// 5432 != 2345
// ^~~~
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.char) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func isLetter(char byte) bool {
	//             'a'           'z'           'A'           'Z'           '_'
	return char >= 97 && char <= 122 || char >= 65 && char <= 90 || char == 95
}

// TODO: Support floats
func isDigit(char byte) bool {
	//            '0'           '9'
	return char >= 48 && char <= 57
}

func (l *Lexer) skipWhitespace() (tok token.Token) {
	//           ' '            '\t'            '\n'            '\r'
	for l.char == 32 || l.char == 9 || l.char == 10 || l.char == 13 {
		l.readChar()
	}

	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition > len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func newToken(tokenType token.TokenType) token.Token {
	return token.Token{Type: tokenType, Literal: tokenType.Literal()}
}
