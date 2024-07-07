package parser

import (
	"fmt"
	"github.com/benja-vq/gonkey/ast"
	"github.com/benja-vq/gonkey/config"
	"github.com/benja-vq/gonkey/lexer"
	"github.com/benja-vq/gonkey/token"
	"log"
	"strconv"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Precedences
const (
	LOWEST      = iota
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func NewParser(l *lexer.Lexer) *Parser {

	parser := Parser{
		l:      l,
		errors: []string{},
	}

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn, 6)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)

	parser.infixParseFns = make(map[token.TokenType]infixParseFn, 8)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)

	// Initialize currToken and peekToken
	parser.nextToken()
	parser.nextToken()

	return &parser
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) noPrefixParseFnError(tt token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function found for %s", tt)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(tt token.TokenType) {
	msg := fmt.Sprintf("Peeking returned an incorrect token, got %s want %s",
		p.peekToken.Type, tt)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {

	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}

}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENT) {
		log.Printf("Let statement expects next token to be an identifier, got %+v",
			p.peekToken)
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		log.Printf("Identifier expects an assignment, got %+v",
			p.peekToken)
		return nil
	}

	// TODO: Skipping expressions for now, just find a semicolon.
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	// TODO: Skipping expressions for now, just find a semicolon.
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	if config.Debug != nil && *config.Debug {
		defer untrace(trace("parseExpressionStatement"))
	}
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if config.Debug != nil && *config.Debug {
		defer untrace(trace("parseExpression"))
	}
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	leftExpression := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		p.nextToken()
		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	if config.Debug != nil && *config.Debug {
		defer untrace(trace("parseIntegerLiteral"))
	}
	lit := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("Could not parse %q as an integer",
			p.currToken.Literal))
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	if config.Debug != nil && *config.Debug {
		defer untrace(trace("parsePrefixExpression"))
	}
	prefixExp := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.nextToken()

	prefixExp.Right = p.parseExpression(PREFIX)

	return prefixExp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	if config.Debug != nil && *config.Debug {
		defer untrace(trace("parseInfixExpression"))
	}
	exp := &ast.InfixExpression{
		Token:    p.currToken,
		Left:     left,
		Operator: p.currToken.Literal,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.currToken,
		Value: p.currTokenIs(token.TRUE),
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) currTokenIs(tt token.TokenType) bool {
	return p.currToken.Type == tt
}

func (p *Parser) peekTokenIs(tt token.TokenType) bool {
	return p.peekToken.Type == tt
}

func (p *Parser) expectPeek(tt token.TokenType) bool {
	if p.peekTokenIs(tt) {
		p.nextToken()
		return true
	}

	p.peekError(tt)
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
}
