package parse

import (
	"fmt"

	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/token"
)

// ===== Node construction =====

func (p *Parser) newVarInitializedStatement(
	name token.T, init ast.Expr,
) *ast.VarInitialized {
	varStmt := &ast.VarInitialized{name, init}
	p.traceNode(varStmt)
	return varStmt
}

func (p *Parser) newVarUninitializedStatement(
	name token.T,
) *ast.VarUninitialized {
	varStmt := &ast.VarUninitialized{name}
	p.traceNode(varStmt)
	return varStmt
}

func (p *Parser) newFunction(
	name token.T, params []token.T, body []ast.Stmt,
) *ast.Function {
	function := &ast.Function{name, params, body}
	p.traceNode(function)
	return function
}

func (p *Parser) newNoopStatement() *ast.Noop {
	nullStmt := &ast.Noop{}
	p.traceNode(nullStmt)
	return nullStmt
}

func (p *Parser) newExpressionStatement(expr ast.Expr) *ast.Expression {
	expressionStmt := &ast.Expression{expr}
	p.traceNode(expressionStmt)
	return expressionStmt
}

func (p *Parser) newPrintStatement(tok token.T, expr ast.Expr) *ast.Print {
	printStmt := &ast.Print{tok, expr}
	p.traceNode(printStmt)
	return printStmt
}

func (p *Parser) newReturnStatement(tok token.T, expr ast.Expr) *ast.Return {
	returnStmt := &ast.Return{tok, expr}
	p.traceNode(returnStmt)
	return returnStmt
}

func (p *Parser) newPanicStatement(tok token.T, expr ast.Expr) *ast.Panic {
	panicStmt := &ast.Panic{tok, expr}
	p.traceNode(panicStmt)
	return panicStmt
}

func (p *Parser) newBlockStatement(tok token.T, body []ast.Stmt) *ast.Block {
	blockStmt := &ast.Block{tok, body}
	p.traceNode(blockStmt)
	return blockStmt
}

func (p *Parser) newIfStatement(cond ast.Expr, thenB, elseB ast.Stmt) *ast.If {
	ifStmt := &ast.If{cond, thenB, elseB}
	p.traceNode(ifStmt)
	return ifStmt
}

func (p *Parser) newWhileStatement(cond ast.Expr, body ast.Stmt) *ast.While {
	whileStmt := &ast.While{cond, body}
	p.traceNode(whileStmt)
	return whileStmt
}

func (p *Parser) newClass(
	name token.T, superclass *ast.Variable, methods []*ast.Function,
) *ast.Class {
	class := &ast.Class{name, superclass, methods}
	p.traceNode(class)
	return class
}

// ===== Parsing =====

func (p *Parser) program() (result []ast.Stmt) {
	var statements []ast.Stmt = []ast.Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() (result ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if perr, ok := r.(ParseError); ok {
				p.synchronize()
				msg := fmt.Sprintf(
					"parse error at %s; this code shouldn't be run",
					perr.Token.Whence(),
				)
				result = p.newPanicStatement(
					perr.Token,
					p.newLiteral(token.StringValue{msg}),
				)
			} else {
				panic(r)
			}
		}
	}()

	if p.match(token.Class) {
		return p.classDeclaration()
	}
	if p.match(token.Fun) {
		return p.function("function")
	}
	if p.match(token.Var) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) classDeclaration() ast.Stmt {
	var name token.T = p.consume(token.Identifier, "Expect class name.")

	var superclass *ast.Variable = nil
	if p.match(token.Less) {
		p.consume(token.Identifier, "Expect superclass name.")
		superclass = p.newVariable(p.previous())
	}

	var methods []*ast.Function

	p.consume(token.LeftBrace, "Expect '{' before class body.")
	for !p.check(token.RightBrace) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))
	}
	p.consume(token.RightBrace, "Expect '}' after class body.")

	return p.newClass(name, superclass, methods)
}

func (p *Parser) varDeclaration() ast.Stmt {
	var name token.T = p.consume(token.Identifier, "Expect variable name.")

	var stmt ast.Stmt
	if p.match(token.Equal) {
		stmt = p.newVarInitializedStatement(name, p.expression())
	} else {
		stmt = p.newVarUninitializedStatement(name)
	}

	p.consume(token.Semicolon, "Expect `;` after variable declaration.")
	return stmt
}

func (p *Parser) function(kind string) *ast.Function {
	var name token.T = p.consume(token.Identifier, "Expect "+kind+" name.")
	p.consume(token.LeftParen, "Expect `(` after "+kind+" name.")
	params := []token.T{}
	for !p.check(token.RightParen) {
		if len(params) >= 255 {
			p.Error(p.peek(), "Cannot have more than 255 parameters.")
		}
		params = append(params,
			p.consume(token.Identifier, "Expect parameter name."))
		if !p.match(token.Comma) {
			break
		}
	}
	p.consume(token.RightParen, "Expect `)` after parameters.")

	p.consume(token.LeftBrace, "Expect `{` before "+kind+" body.")
	return p.newFunction(name, params, p.block())
}

func (p *Parser) statement() ast.Stmt {
	if p.match(token.If) {
		return p.ifStatement()
	}
	if p.match(token.Print) {
		return p.printStatement()
	}
	if p.match(token.Return) {
		return p.returnStatement()
	}
	if p.match(token.While) {
		return p.whileStatement()
	}
	if p.match(token.For) {
		return p.forStatement()
	}
	if p.match(token.LeftBrace) {
		lbrace := p.previous()
		return p.newBlockStatement(lbrace, p.block())
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() ast.Stmt {
	p.consume(token.LeftParen, "Expect `(` after `if`.")
	var condition ast.Expr = p.expression()
	p.consume(token.RightParen, "Expect `)` after `if` condition.")

	var thenBranch ast.Stmt = p.statement()
	var elseBranch ast.Stmt = nil
	if p.match(token.Else) {
		elseBranch = p.statement()
	}

	return p.newIfStatement(condition, thenBranch, elseBranch)
}

func (p *Parser) whileStatement() ast.Stmt {
	p.consume(token.LeftParen, "Expect `(` after `while`.")
	var condition ast.Expr = p.expression()
	p.consume(token.RightParen, "Expect `)` after `while` condition.")

	var body ast.Stmt = p.statement()

	return p.newWhileStatement(condition, body)
}

func (p *Parser) forStatement() ast.Stmt {
	lParen := p.peek()
	p.consume(token.LeftParen, "Expect `(` after `for`.")
	var initializer ast.Stmt
	if p.match(token.Semicolon) {
		initializer = nil
	} else if p.match(token.Var) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition ast.Expr
	if !p.check(token.Semicolon) {
		condition = p.expression()
	}
	p.consume(token.Semicolon, "Expect `;` after `for` condition.")

	var increment ast.Stmt
	if !p.check(token.Semicolon) {
		increment = p.newExpressionStatement(p.expression())
	}
	p.consume(token.RightParen, "Expect `)` after `for` clauses.")

	var body ast.Stmt = p.statement()

	if increment != nil {
		body = p.newBlockStatement(lParen, []ast.Stmt{body, increment})
	}
	if condition == nil {
		condition = p.newLiteral(token.BooleanValue{true})
	}
	body = p.newWhileStatement(condition, body)
	if initializer != nil {
		body = p.newBlockStatement(lParen, []ast.Stmt{initializer, body})
	}

	return body
}

func (p *Parser) printStatement() ast.Stmt {
	var keyword token.T = p.previous()
	var value ast.Expr = p.expression()
	p.consume(token.Semicolon, "Expect `;` after value.")
	return p.newPrintStatement(keyword, value)
}

func (p *Parser) returnStatement() ast.Stmt {
	var keyword token.T = p.previous()
	var value ast.Expr
	if p.check(token.Semicolon) {
		p.advance()
	} else {
		value = p.expression()
		p.consume(token.Semicolon, "Expect expression or `;` after `return`.")
	}

	return p.newReturnStatement(keyword, value)
}

func (p *Parser) expressionStatement() ast.Stmt {
	var expr ast.Expr = p.expression()
	p.consume(token.Semicolon, "Expect `;` after expression statement.")
	return p.newExpressionStatement(expr)
}

func (p *Parser) block() []ast.Stmt {
	var statements []ast.Stmt = make([]ast.Stmt, 0, 5)

	for !p.check(token.RightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(token.RightBrace, "Expect `}` after block.")
	return statements
}
