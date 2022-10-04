package parse

import (
	"fmt"

	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/lox"
	"github.com/perlmonger42/go-lox/token"
)

type T interface {
	Parse() []ast.Stmt // synonym for ParseProg
	ParseProg() []ast.Stmt
	ParseExpr() ast.Expr
}

func New(lox *lox.T, tokens []token.T) T {
	return &Parser{
		lox:     lox,
		tokens:  tokens,
		current: 0,
	}
}

type Parser struct {
	lox     *lox.T
	tokens  []token.T
	current int // index of current token in tokens[]
}

var _ T = &Parser{}

type ParseError struct {
	Token   token.T
	Message string
}

// Parse is a synonym for ParseProg.
// Sets p.lox.HadError if there was any parsing error.
func (p *Parser) Parse() (result []ast.Stmt) { return p.ParseProg() }

// ParseExpr parses a sequence of statements, consuming all the input.
// Sets p.lox.HadError if there was any parsing error.
func (p *Parser) ParseProg() (result []ast.Stmt) { return p.program() }

// ParseExpr parses an expression as the entire input.
// Sets p.lox.HadError if there was any parsing error.
// It is a parsing error for there to be input following the expression.
func (p *Parser) ParseExpr() (result ast.Expr) { return p.expressionOnly() }

func (p *Parser) Error(tok token.T, message string) ParseError {
	p.lox.Error(tok, message)
	return ParseError{tok, message}
}

func (p *Parser) traceToken() {
	if p.lox.Config.TraceParseTokens {
		tok := p.peek()
		fmt.Printf("[%s] consuming %s\n", tok.Whence(), tok)
	}
}

func (p *Parser) traceNode(node ast.Node) ast.Node {
	if p.lox.Config.TraceNodes {
		fmt.Printf("[%s] built %s\n", p.peek().Whence(), ast.ToString(node))
	}
	return node
}

// synchronize discards tokens until it finds what looks like
// a statement boundary
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type() == token.Semicolon {
			return
		}

		switch p.peek().Type() {
		case token.Class, token.Fun, token.Var, token.For,
			token.If, token.While, token.Print, token.Return:
			return
		}
		p.advance()
	}
}

func (p *Parser) match(types ...token.Type) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(typ token.Type, message string) token.T {
	if p.check(typ) {
		return p.advance()
	}

	tok := p.peek()
	panic(p.Error(tok, fmt.Sprintf("found %s; %s", tok.Type(), message)))
}

func (p *Parser) check(typ token.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type() == typ
}

func (p *Parser) advance() token.T {
	if !p.isAtEnd() {
		p.traceToken()
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type() == token.EOF
}

func (p *Parser) peek() token.T {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.T {
	return p.tokens[p.current-1]
}
