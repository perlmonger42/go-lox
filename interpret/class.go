package interpret

import (
	"github.com/perlmonger42/go-lox/token"
)

type LoxClass struct {
	Name       token.T
	Superclass *LoxClass
	Methods    map[string]*LoxFunction
}

var _ token.Object = &LoxClass{}
var _ LoxCallable = &LoxClass{}

func NewLoxClass(
	name token.T,
	superclass *LoxClass,
	methods map[string]*LoxFunction,
) *LoxClass {
	return &LoxClass{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}
}

func (f *LoxClass) Arity() int {
	if initializer, ok := f.FindMethod("init"); ok {
		return initializer.Arity()
	}
	return 0
}

func (f *LoxClass) Call(i T, arguments []token.Value) (result token.Value) {
	var instance *LoxInstance = NewLoxInstance(f)
	if initializer, ok := f.FindMethod("init"); ok {
		initializer.Bind(instance).Call(i, arguments)
	}
	return token.ObjectValue{instance}
}

func (c *LoxClass) FindMethod(name string) (method *LoxFunction, ok bool) {
	method, ok = c.Methods[name]
	if !ok && c.Superclass != nil {
		method, ok = c.Superclass.FindMethod(name)
	}
	return
}

func (c *LoxClass) EqualsObject(o token.Object) bool {
	if lf, ok := o.(*LoxClass); ok {
		return c.Name == lf.Name
	}
	return false
}

func (c *LoxClass) String() string {
	return "class " + c.Name.Lexeme()
}

func (c *LoxClass) Show() string {
	return c.String()
}
