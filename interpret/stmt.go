package interpret

import (
	"fmt"

	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/token"
)

var _ T = &Interpreter{}
var _ ast.Visitor_Stmt = &Interpreter{}

func (i *Interpreter) InterpretStmts(statements []ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if exception, ok := r.(RuntimeError); ok {
				fmt.Printf("runtime error: %s\n", exception)
			} else {
				panic(r)
			}
		}
	}()
	i.depth = 0

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) execute(stmt ast.Stmt) {
	stmt.Accept_Stmt(i)
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, newEnv Environment) {
	var previous Environment = i.environment
	i.environment = newEnv
	defer func() { i.environment = previous }()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) Visit_NoopStmt(stmt *ast.Noop) {
}

func (i *Interpreter) Visit_ExpressionStmt(stmt *ast.Expression) {
	i.evaluate(stmt.Expression)
}

func (i *Interpreter) Visit_PrintStmt(stmt *ast.Print) {
	var value Value = i.evaluate(stmt.Expression)
	fmt.Printf("%s\n", value.String())
}

type PanicForReturn struct {
	Token  token.T
	Result Value
}

func (pfr *PanicForReturn) Error() string {
	return fmt.Sprintf("returning %s from %s", pfr.Result, pfr.Result)
}

func (i *Interpreter) Visit_ReturnStmt(stmt *ast.Return) {
	var value Value
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	panic(PanicForReturn{Token: stmt.Keyword, Result: value})
}

func (i *Interpreter) Visit_PanicStmt(stmt *ast.Panic) {
	var value Value = i.evaluate(stmt.Expression)
	i.Error(stmt.Keyword, value.String())
}

func (i *Interpreter) Visit_BlockStmt(stmt *ast.Block) {
	env := NewNestedEnvironment(i.environment)
	i.executeBlock(stmt.Statements, env)
}

func (i *Interpreter) Visit_VarInitializedStmt(stmt *ast.VarInitialized) {
	var value Value = i.evaluate(stmt.Initializer)
	i.Define(stmt.Name, value)
}

func (i *Interpreter) Visit_VarUninitializedStmt(stmt *ast.VarUninitialized) {
	i.Define(stmt.Name, &token.NilValue{})
}

func (i *Interpreter) Visit_ClassStmt(stmt *ast.Class) {
	var superclass *LoxClass = nil
	if stmt.Superclass != nil {
		super1 := i.evaluate(stmt.Superclass)
		if super2, ok := super1.(token.ObjectValue); ok {
			if super3, ok := super2.V.(*LoxClass); ok {
				superclass = super3
			}
		}
		if superclass == nil {
			panic(i.Error(stmt.Superclass.Name,
				"Superclass must be a class."))
		}
	}

	i.Define(stmt.Name, token.NilValue{}) // make class name visible to methods

	saveEnvironment := i.environment
	if stmt.Superclass != nil {
		i.environment = NewNestedEnvironment(i.environment)
		super := &token.Token{Type_: token.Super, Lexeme_: "super"}
		i.Define(super, token.ObjectValue{superclass})
	}

	env := i.GetCurrentEnvironment()
	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.Methods {
		var function *LoxFunction = NewLoxFunction(
			method, env,
			method.Name.Lexeme() == "init",
		)
		methods[method.Name.Lexeme()] = function
	}

	class := NewLoxClass(stmt.Name, superclass, methods)

	if superclass != nil {
		i.environment = saveEnvironment
	}

	i.Assign(stmt.Name, token.ObjectValue{class})
}

func (i *Interpreter) Visit_FunctionStmt(stmt *ast.Function) {
	var function *LoxFunction = NewLoxFunction(
		stmt, i.GetCurrentEnvironment(), false,
	)
	i.environment.Define(stmt.Name, token.ObjectValue{function})
}

func (i *Interpreter) Visit_IfStmt(stmt *ast.If) {
	if isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
}

func (i *Interpreter) Visit_WhileStmt(stmt *ast.While) {
	for isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
}
