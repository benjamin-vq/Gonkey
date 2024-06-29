package parser

import (
	"github.com/benja-vq/gonkey/ast"
	"github.com/benja-vq/gonkey/lexer"
	"github.com/benja-vq/gonkey/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token
}

func NewParser(l *lexer.Lexer) *Parser {

	parser := Parser{
		l: l,
	}

	// Initialize currToken and peekToken
	parser.nextToken()
	parser.nextToken()

	return &parser
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
