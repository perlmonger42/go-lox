package interpret

import (
	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/token"
)

type LoxCallable interface {
	Callable
}

type LoxFunction struct {
	Declaration   *ast.Function
	Closure       Environment
	IsInitializer bool
}

var _ LoxCallable = &LoxFunction{}
var _ token.Object = &LoxFunction{}

func NewLoxFunction(
	declaration *ast.Function,
	closure Environment,
	isInitializer bool,
) *LoxFunction {
	return &LoxFunction{
		Declaration:   declaration,
		Closure:       closure,
		IsInitializer: isInitializer,
	}
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *LoxFunction) Call(i T, arguments []token.Value) (result token.Value) {
	defer func() {
		if r := recover(); r != nil {
			if pfr, ok := r.(PanicForReturn); ok {
				result = pfr.Result
			} else {
				panic(r)
			}
		}
		if f.IsInitializer {
			// instance initializer functions always return `this`
			// (even if it explicitly returns something else).
			if self, err := f.Closure.GetAt(0, "this"); err != nil {
				panic(err)
			} else {
				result = self
			}
		}
	}()

	environment := NewNestedEnvironment(f.Closure)
	for i, param := range f.Declaration.Params {
		environment.Define(param, arguments[i])
	}
	i.executeBlock(f.Declaration.Body, environment)
	return &token.NilValue{}
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	env := NewNestedEnvironment(f.Closure)
	this := &token.Token{
		Type_:   token.This,
		Whence_: f.Declaration.Name.Whence(),
		Lexeme_: "this",
	}
	env.Define(this, token.ObjectValue{instance})
	return NewLoxFunction(f.Declaration, env, f.IsInitializer)
	//bound := NewLoxFunction(f.Declaration, env)
	//return token.ObjectValue{bound}
}

func (f *LoxFunction) EqualsObject(o token.Object) bool {
	if lf, ok := o.(*LoxFunction); ok {
		return f.Declaration == lf.Declaration
	}
	return false
}

func (f *LoxFunction) String() string {
	return ast.StmtToString(f.Declaration)
}

func (f *LoxFunction) Show() string {
	return f.String()
}
