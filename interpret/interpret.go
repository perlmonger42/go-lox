package interpret

import (
	"fmt"

	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/lox"
	"github.com/perlmonger42/go-lox/token"
)

type Value = token.Value

type T interface {
	InterpretStmts(stmts []ast.Stmt)
	InterpretExpr(expr ast.Expr) Value

	//InterpretStmts2(stmts []ast.Stmt) error
	//InterpretExprMaybe(expr ast.Expr) (Value, error)

	GetGlobalEnvironment() Environment
	GetCurrentEnvironment() Environment
	GetAt(distance int, name string) (token.Value, error)

	GetSuper(super *ast.Super) (distance int, superclass *LoxClass)
	GetThisAt(distance int) *LoxInstance

	Error(tok token.T, message string) RuntimeError

	Resolve(expr ast.Expr, name token.T, depth int)

	executeBlock(statements []ast.Stmt, newEnv Environment)

	printIndent()
	getLox() *lox.T
	indent() string
}

func New(lox *lox.T) T {
	i := &Interpreter{lox: lox, depth: 0}
	i.globals = NewGlobalEnvironment()
	i.environment = i.globals
	i.locals = make(map[ast.Expr]int)

	clock := token.New(token.Identifier, "clock", nil, token.NewPos(0))
	i.globals.Define(clock, token.ObjectValue{&ClockNative{}})

	str := token.New(token.Identifier, "str", nil, token.NewPos(0))
	i.globals.Define(str, token.ObjectValue{&StrNative{}})

	// i.globals.Dump("Interpreter Environment")
	return i
}

type Interpreter struct {
	lox         *lox.T
	depth       int
	globals     Environment
	environment Environment
	locals      map[ast.Expr]int
}

var _ T = &Interpreter{}

type RuntimeError struct {
	Token   token.T
	Message string
}

func (rte *RuntimeError) Error() string { return rte.Message }

func (i *Interpreter) GetGlobalEnvironment() Environment  { return i.globals }
func (i *Interpreter) GetCurrentEnvironment() Environment { return i.environment }
func (i *Interpreter) GetAt(distance int, name string) (Value, error) {
	return i.environment.GetAt(distance, name)
}
func (i *Interpreter) getLox() *lox.T { return i.lox }

func (i *Interpreter) Error(tok token.T, message string) RuntimeError {
	i.lox.Error(tok, message)
	return RuntimeError{tok, message}
}

func (i *Interpreter) Define(name token.T, value Value) {
	if i.getLox().Config.TraceEval {
		fmt.Printf("%sdefine %s <-- %s\n", i.indent(), name.Lexeme(), value)
	}
	if err := i.environment.Define(name, value); err != nil {
		i.Error(err.Token, err.Message)
	}
}

func (i *Interpreter) Assign(name token.T, value Value) {
	i.assignLocal(0, name, value)
}

func (i *Interpreter) GetSuper(
	super *ast.Super,
) (distance int, superclass *LoxClass) {
	if distance, ok := i.locals[super]; !ok {
		panic(i.Error(super.Keyword,
			"[internal error] `super` is not defined."))
	} else if value, err := i.GetAt(distance, "super"); err != nil {
		panic(i.Error(super.Keyword, fmt.Sprintf(
			"[internal error] `super` value could not be found: %s",
			err)))
	} else if object, ok := value.(token.ObjectValue); !ok {
		panic(i.Error(super.Keyword,
			"[internal error] `super` value is not an ObjectValue."))
	} else if superclass, ok := object.V.(*LoxClass); !ok {
		panic(i.Error(super.Keyword,
			"[internal error] `super` value is not a LoxClass."))
	} else {
		return distance, superclass
	}
}

func (i *Interpreter) GetThisAt(distance int) *LoxInstance {
	if value, err := i.GetAt(distance, "this"); err != nil {
		panic(fmt.Errorf(
			"[internal error] `this` value could not be found: %s", err))
	} else if object, ok := value.(token.ObjectValue); !ok {
		panic(fmt.Errorf(
			"[internal error] `this` value is not an ObjectValue."))
	} else if this, ok := object.V.(*LoxInstance); !ok {
		panic(fmt.Errorf(
			"[internal error] `this` value is not a LoxInstance."))
	} else {
		return this
	}
}

func (i *Interpreter) assignGlobal(name token.T, value Value) {
	if i.getLox().Config.TraceEval {
		fmt.Printf("%sassign %s <-- %s\n", i.indent(), name.Lexeme(), value)
	}
	if err := i.globals.Assign(name, value); err != nil {
		i.Error(err.Token, err.Message)
	}
}

func (i *Interpreter) assignLocal(depth int, name token.T, value Value) {
	if i.getLox().Config.TraceEval {
		fmt.Printf("%sassign %s <-- %s\n", i.indent(), name.Lexeme(), value)
	}
	if err := i.environment.AssignAt(depth, name, value); err != nil {
		i.Error(err.Token, err.Message)
	}
}

func (i *Interpreter) Resolve(expr ast.Expr, name token.T, depth int) {
	// fmt.Printf("resolving %s at %s with depth %d\n", expr, name.Whence(), depth)
	i.locals[expr] = depth
}

func (i *Interpreter) getFromGlobals(name token.T) Value {
	value, err := i.globals.GetLocal(name)
	if err != nil {
		i.Error(err.Token, err.Message)
	}
	if i.getLox().Config.TraceEval {
		fmt.Printf("%s%s <-- %s\n", i.indent(), value, name.Lexeme())
	}
	return value
}

func (i *Interpreter) getFromEnvironment(depth int, name token.T) Value {
	value, err := i.environment.GetAt(depth, name.Lexeme())
	if err != nil {
		i.Error(name, err.Error())
	}
	if i.getLox().Config.TraceEval {
		fmt.Printf("%s%s <-- %s\n", i.indent(), value, name.Lexeme())
	}
	return value
}

func (i *Interpreter) lookUpVariable(name token.T, expr ast.Expr) Value {
	if distance, ok := i.locals[expr]; ok {
		v := i.getFromEnvironment(distance, name)
		if i.lox.Config.TraceEval {

			fmt.Printf("fetching local %s at %s from depth %d\n",
				name.Lexeme(), name.Whence().String(), distance)
		}
		return v
	} else {
		v := i.getFromGlobals(name)
		if i.lox.Config.TraceEval {
			fmt.Printf("fetching global %s at %s\n",
				name.Lexeme(), name.Whence().String())
		}
		return v
	}
}

func isTruthy(val Value) bool {
	switch v := val.(type) {
	case token.NilValue:
		return false
	case token.BooleanValue:
		return v.V
	default:
		return true
	}
}
