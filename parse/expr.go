package parse

import (
	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/token"
)

// ===== Node construction =====

func (p *Parser) newLogical(
	op token.T, left ast.Expr, right ast.Expr,
) *ast.Logical {
	logical := &ast.Logical{op, left, right}
	p.traceNode(logical)
	return logical
}

func (p *Parser) newGrouping(expr ast.Expr) *ast.Grouping {
	group := &ast.Grouping{expr}
	p.traceNode(group)
	return group
}

func (p *Parser) newCall(
	callee ast.Expr,
	paren token.T,
	args []ast.Expr,
) *ast.Call {
	call := &ast.Call{callee, paren, args}
	p.traceNode(call)
	return call
}

func (p *Parser) newGet(expr ast.Expr, name token.T) *ast.Get {
	get := &ast.Get{expr, name}
	p.traceNode(get)
	return get
}

func (p *Parser) newSet(expr ast.Expr, name token.T, value ast.Expr) *ast.Set {
	set := &ast.Set{expr, name, value}
	p.traceNode(set)
	return set
}

func (p *Parser) newUnary(op token.T, right ast.Expr) *ast.Unary {
	unary := &ast.Unary{op, right}
	p.traceNode(unary)
	return unary
}

func (p *Parser) newBinary(
	op token.T, left ast.Expr, right ast.Expr,
) *ast.Binary {
	binary := &ast.Binary{op, left, right}
	p.traceNode(binary)
	return binary
}

func (p *Parser) newAssign(name token.T, value ast.Expr) *ast.Assign {
	assign := &ast.Assign{name, value}
	p.traceNode(assign)
	return assign
}

func (p *Parser) newVariable(tok token.T) *ast.Variable {
	variable := &ast.Variable{tok}
	p.traceNode(variable)
	return variable
}

func (p *Parser) newLiteral(value token.Value) *ast.Literal {
	literal := &ast.Literal{value}
	p.traceNode(literal)
	return literal
}

func (p *Parser) newThis(keyword token.T) *ast.This {
	this := &ast.This{keyword}
	p.traceNode(this)
	return this
}

func (p *Parser) newSuper(keyword, method token.T) *ast.Super {
	this := &ast.Super{keyword, method}
	p.traceNode(this)
	return this
}

// ===== Parsing =====

// expressionOnly parses an expression as the entire input.
// It is an error for there to be input following the expression.
func (p *Parser) expressionOnly() (result ast.Expr) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(ParseError); ok {
				result = nil
			} else {
				panic(r)
			}
		}
	}()

	result = p.expression()
	if p.peek().Type() != token.EOF {
		p.Error(p.peek(), "unexpected input after expression")
	}
	return result
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) assignment() ast.Expr {
	var expr ast.Expr = p.or()

	if p.match(token.Equal) {
		var equals token.T = p.previous()
		var value ast.Expr = p.assignment()

		if v, ok := expr.(*ast.Variable); ok {
			var name token.T = v.Name
			return p.newAssign(name, value)
		} else if get, ok := expr.(*ast.Get); ok {
			return p.newSet(get.Object, get.Name, value)
		}

		p.Error(equals, "Invalid assignment target.")
	}

	return expr
}

func (p *Parser) or() ast.Expr {
	var expr ast.Expr = p.and()

	for p.match(token.Or) {
		var operator token.T = p.previous()
		var right ast.Expr = p.and()
		expr = p.newLogical(operator, expr, right)
	}

	return expr
}

func (p *Parser) and() ast.Expr {
	var expr ast.Expr = p.equality()

	for p.match(token.And) {
		var operator token.T = p.previous()
		var right ast.Expr = p.equality()
		expr = p.newLogical(operator, expr, right)
	}

	return expr
}

func (p *Parser) equality() ast.Expr {
	var expr ast.Expr = p.comparison()

	for p.match(token.BangEqual, token.EqualEqual) {
		var operator token.T = p.previous()
		var right ast.Expr = p.comparison()
		expr = p.newBinary(operator, expr, right)
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	var expr ast.Expr = p.addition()

	for p.match(
		token.Greater, token.GreaterEqual, token.Less, token.LessEqual,
	) {
		var operator token.T = p.previous()
		var right ast.Expr = p.addition()
		expr = p.newBinary(operator, expr, right)
	}

	return expr
}

func (p *Parser) addition() ast.Expr {
	var expr ast.Expr = p.multiplication()

	for p.match(token.Minus, token.Plus) {
		var operator token.T = p.previous()
		var right ast.Expr = p.multiplication()
		expr = p.newBinary(operator, expr, right)
	}

	return expr
}

func (p *Parser) multiplication() ast.Expr {
	var expr ast.Expr = p.unary()

	for p.match(token.Slash, token.Star) {
		var operator token.T = p.previous()
		var right ast.Expr = p.unary()
		expr = p.newBinary(operator, expr, right)
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(token.Bang, token.Minus) {
		var operator token.T = p.previous()
		var right ast.Expr = p.unary()
		return p.newUnary(operator, right)
	}

	return p.call()
}

func (p *Parser) call() ast.Expr {
	var expr ast.Expr = p.primary()
	for {
		if p.match(token.LeftParen) {
			expr = p.finishCall(expr)
		} else if p.match(token.Dot) {
			name := p.consume(token.Identifier,
				"Expect property name after `.`.")
			expr = p.newGet(expr, name)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee ast.Expr) ast.Expr {
	var args []ast.Expr = []ast.Expr{}
	for !p.check(token.RightParen) {
		if len(args) >= 255 {
			p.Error(p.peek(), "Cannot have more than 255 arguments.")
		}
		args = append(args, p.expression())
		if !p.match(token.Comma) {
			break
		}
	}
	paren := p.consume(token.RightParen, "Expect `)` after arguments.")
	return p.newCall(callee, paren, args)
}

func (p *Parser) primary() ast.Expr {
	if p.match(token.False) {
		return p.newLiteral(token.BooleanValue{false})
	}
	if p.match(token.True) {
		return p.newLiteral(token.BooleanValue{true})
	}
	if p.match(token.Nil) {
		return p.newLiteral(token.NilValue{})
	}
	if p.match(token.This) {
		return p.newThis(p.previous())
	}
	if p.match(token.Super) {
		keyword := p.previous()
		p.consume(token.Dot, "Expect '.' after 'super'.")
		method := p.consume(token.Identifier,
			"Expect superclass method name after `super.`.")
		return p.newSuper(keyword, method)
	}

	if p.match(token.Number, token.String) {
		return p.newLiteral(p.previous().Literal())
	}

	if p.match(token.Identifier) {
		return p.newVariable(p.previous())
	}

	if p.match(token.LeftParen) {
		var expr ast.Expr = p.expression()
		p.consume(token.RightParen, "Expect ')' after expression.")
		return p.newGrouping(expr)
	}

	panic(p.Error(p.peek(), "Expect expression."))
}
