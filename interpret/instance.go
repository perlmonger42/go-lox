package interpret

import (
	"fmt"
	"strings"

	"github.com/perlmonger42/go-lox/token"
)

type LoxInstance struct {
	class  *LoxClass
	fields map[string]token.Value
}

var _ token.Object = &LoxInstance{}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class, fields: make(map[string]token.Value)}
}

func (i *LoxInstance) Get(name token.T) (token.Value, error) {
	if value, ok := i.fields[name.Lexeme()]; ok {
		return value, nil
	}

	if method, ok := i.class.FindMethod(name.Lexeme()); ok {
		bound_method := method.Bind(i)
		return token.ObjectValue{bound_method}, nil
	}

	return &token.NilValue{},
		&RuntimeError{name,
			fmt.Sprintf("Undefined property `%s`.", name.Lexeme())}
}

func (i *LoxInstance) Set(name token.T, val token.Value) {
	i.fields[name.Lexeme()] = val
}

func (i *LoxInstance) EqualsObject(o token.Object) bool {
	if it, ok := o.(*LoxInstance); ok {
		return i == it
	}
	return false
}

func (i *LoxInstance) String() string {
	values := []string{}
	for key, val := range i.fields {
		values = append(values, fmt.Sprintf("%s: %s", key, val.Show()))
	}
	return i.class.Name.Lexeme() + "{" + strings.Join(values, ", ") + "}"
}

func (i *LoxInstance) Show() string {
	return i.String()
}
