package resolve

import (
	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/lox"
	"github.com/perlmonger42/go-lox/token"
)

type Resolver interface {
	Resolve(expr ast.Expr, name token.T, depth int)
}

type FunctionType int

const (
	NOT_FUNCTION FunctionType = iota
	FUNCTION
	INITIALIZER // "init" method of a class
	METHOD
)

type ClassType int

const (
	NOT_CLASS ClassType = iota
	CLASS
	SUBCLASS
)

type T struct {
	lox             *lox.T
	resolver        Resolver
	scopes          []map[string]bool
	currentFunction FunctionType
	currentClass    ClassType
}

var _ ast.Visitor_Stmt = &T{}
var _ ast.Visitor_Expr = &T{}

func New(lox *lox.T, resolver Resolver) *T {
	return &T{
		lox:      lox,
		resolver: resolver,
	}
}

func (r *T) setCurrentClass(class ClassType) ClassType {
	was := r.currentClass
	r.currentClass = class
	return was
}

func (r *T) topScope() map[string]bool {
	if nScopes := len(r.scopes); nScopes > 0 {
		return r.scopes[nScopes-1]
	}
	return nil
}

func (r *T) topScopeFetch(name token.T) (defined bool, ok bool) {
	if s := r.topScope(); s != nil {
		defined, ok = s[name.Lexeme()]
		return
	}
	return false, false
}

func (r *T) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *T) declare(name token.T) {
	if s := r.topScope(); s != nil {
		s[name.Lexeme()] = false
	}
}

func (r *T) define(name token.T) {
	if s := r.topScope(); s != nil {
		s[name.Lexeme()] = true
	}
}

func (r *T) endScope() {
	r.scopes = r.scopes[0 : len(r.scopes)-1]
}

func (r *T) ResolveStmtList(statements []ast.Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *T) resolveStmt(statement ast.Stmt) {
	statement.Accept_Stmt(r)
}

func (r *T) resolveExpr(expr ast.Expr) {
	expr.Accept_Expr(r)
}

func (r *T) resolveLocal(expr ast.Expr, name token.T) {
	depth := 0
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme()]; ok {
			r.resolver.Resolve(expr, name, depth)
			return
		}
		depth++
	}

	// Not found. Assume it is global.
}

func (r *T) resolveFunction(function *ast.Function, ft FunctionType) {
	var enclosingFunction FunctionType = r.currentFunction
	r.currentFunction = ft
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStmtList(function.Body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *T) Visit_BlockStmt(stmt *ast.Block) {
	r.beginScope()
	r.ResolveStmtList(stmt.Statements)
	r.endScope()
}

func (r *T) Visit_NoopStmt(stmt *ast.Noop) {
}

func (r *T) Visit_ExpressionStmt(stmt *ast.Expression) {
	r.resolveExpr(stmt.Expression)
}

func (r *T) Visit_PrintStmt(stmt *ast.Print) {
	r.resolveExpr(stmt.Expression)
}

func (r *T) Visit_ReturnStmt(stmt *ast.Return) {
	if r.currentFunction == NOT_FUNCTION {
		r.lox.Error(stmt.Keyword, "Cannot return from top-level code.")
	}
	if stmt.Value != nil {
		if r.currentFunction == INITIALIZER {
			r.lox.Error(
				stmt.Keyword, "Cannot return a value from an initializer.",
			)
		}

		r.resolveExpr(stmt.Value)
	}
}

func (r *T) Visit_PanicStmt(stmt *ast.Panic) {
	r.resolveExpr(stmt.Expression)
}

func (r *T) Visit_VarInitializedStmt(stmt *ast.VarInitialized) {
	//r.declare(stmt.Name) // declare here to prevent var a=a+1;
	r.resolveExpr(stmt.Initializer)
	r.declare(stmt.Name) // declare here to allow var a=a+1;
	r.define(stmt.Name)
}

func (r *T) Visit_VarUninitializedStmt(stmt *ast.VarUninitialized) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
}

func (r *T) Visit_ClassStmt(stmt *ast.Class) {
	enclosingClass := r.setCurrentClass(CLASS)
	defer r.setCurrentClass(enclosingClass)

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.Superclass != nil {
		if stmt.Name.Lexeme() == stmt.Superclass.Name.Lexeme() {
			r.lox.Error(stmt.Superclass.Name,
				"A class can't inherit from itself.",
			)
		}
		r.resolveExpr(stmt.Superclass)
	}

	if stmt.Superclass != nil {
		r.setCurrentClass(SUBCLASS)
		r.beginScope()
		r.scopes[len(r.scopes)-1]["super"] = true
		defer r.endScope()
	}

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true
	defer r.endScope()

	for _, method := range stmt.Methods {
		var declaration FunctionType = METHOD
		if method.Name.Lexeme() == "init" {
			declaration = INITIALIZER
		}
		r.resolveFunction(method, declaration)
	}
}

func (r *T) Visit_FunctionStmt(stmt *ast.Function) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FUNCTION)
}

func (r *T) Visit_IfStmt(stmt *ast.If) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
}

func (r *T) Visit_WhileStmt(stmt *ast.While) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
}

func (r *T) Visit_GroupingExpr(expr *ast.Grouping) {
	r.resolveExpr(expr.Expression)
}

func (r *T) Visit_CallExpr(expr *ast.Call) {
	r.resolveExpr(expr.Callee)
	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}
}

func (r *T) Visit_GetExpr(expr *ast.Get) {
	r.resolveExpr(expr.Object)
}

func (r *T) Visit_SetExpr(expr *ast.Set) {
	r.resolveExpr(expr.Object)
	r.resolveExpr(expr.Value)
}

func (r *T) Visit_UnaryExpr(expr *ast.Unary) {
	r.resolveExpr(expr.Right)
}

func (r *T) Visit_BinaryExpr(expr *ast.Binary) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
}

func (r *T) Visit_LogicalExpr(expr *ast.Logical) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
}

func (r *T) Visit_VariableExpr(expr *ast.Variable) {
	if defined, ok := r.topScopeFetch(expr.Name); ok && !defined {
		r.lox.Error(expr.Name,
			"Cannot read local variable in its own initializer.",
		)
	}

	r.resolveLocal(expr, expr.Name)
}

func (r *T) Visit_ThisExpr(expr *ast.This) {
	if r.currentClass == NOT_CLASS {
		r.lox.Error(expr.Keyword, "Cannot use `this` outside of a class.")
		return
	}
	r.resolveLocal(expr, expr.Keyword)
}

func (r *T) Visit_SuperExpr(expr *ast.Super) {
	if r.currentClass == NOT_CLASS {
		r.lox.Error(expr.Keyword, "Can't use 'super' outside of a class.")
	} else if r.currentClass != SUBCLASS {
		r.lox.Error(expr.Keyword,
			"Can't use 'super' in a class with no superclass.")
	}
	r.resolveLocal(expr, expr.Keyword)
}

func (r *T) Visit_LiteralExpr(expr *ast.Literal) {
}

func (r *T) Visit_AssignExpr(expr *ast.Assign) {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
}
