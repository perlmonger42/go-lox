package interpret

import (
	"fmt"

	"github.com/perlmonger42/go-lox/token"
)

// ===== Environment and its constructors =====

type Environment interface {
	Define(name token.T, value Value) *RuntimeError
	Assign(name token.T, value Value) *RuntimeError
	AssignAt(distance int, name token.T, value Value) *RuntimeError
	GetLocal(name token.T) (Value, *RuntimeError)
	GetAt(distance int, name string) (Value, error)

	Dump(msg string)
}

func NewGlobalEnvironment() *globalEnv {
	return &globalEnv{
		Values: make(map[string]Value),
	}
}

func NewNestedEnvironment(outer Environment) *nestedEnv {
	v := &nestedEnv{
		globalEnv: *NewGlobalEnvironment(),
		Enclosing: outer,
	}
	return v
}

// ===== global Environment =====

type globalEnv struct {
	Values map[string]Value
}

func (v *globalEnv) Define(name token.T, value Value) *RuntimeError {
	v.Values[name.Lexeme()] = value
	// v.Dump("after defining global " + name.Lexeme() + " at " + name.Whence().String())
	return nil
}

func (v *globalEnv) Assign(name token.T, value Value) *RuntimeError {
	if _, ok := v.Values[name.Lexeme()]; ok {
		v.Values[name.Lexeme()] = value
		return nil
	}
	return &RuntimeError{
		name,
		fmt.Sprintf("Undefined variable '%s'.", name.Lexeme()),
	}
}

func (v *globalEnv) AssignAt(
	distance int,
	name token.T,
	value Value,
) *RuntimeError {
	if distance == 0 {
		return v.Assign(name, value)
	}
	panic(fmt.Sprintf("globalEnv.AssignAt called with distance %d > 0",
		distance))
}

func (v *globalEnv) GetLocal(name token.T) (Value, *RuntimeError) {
	if v, ok := v.Values[name.Lexeme()]; ok {
		return v, nil
	}
	// v.Dump(fmt.Sprintf("%q not found in global environment", name.Lexeme()))

	err := &RuntimeError{
		name,
		fmt.Sprintf("Undefined variable '%s'.", name.Lexeme()),
	}
	return token.NilValue{}, err
}

func (v *globalEnv) GetAt(distance int, name string) (Value, error) {
	if distance != 0 {
		panic(fmt.Sprintf(
			"[internal error] globalEnv.GetAt called with distance %d > 0",
			distance,
		))
	}

	if v, ok := v.Values[name]; ok {
		return v, nil
	}
	return &token.NilValue{}, fmt.Errorf("Undefined variable '%s'.", name)
}

// ===== nested environment =====

type nestedEnv struct {
	globalEnv
	Enclosing Environment
}

func (v *nestedEnv) Define(name token.T, value Value) (err *RuntimeError) {
	if _, ok := v.Values[name.Lexeme()]; ok {
		err = &RuntimeError{
			name,
			fmt.Sprintf("Variable '%s' redefined.", name.Lexeme()),
		}
	}

	v.Values[name.Lexeme()] = value
	// v.Dump("after defining local " + name.Lexeme() + " at " + name.Whence().String())
	return err
}

func (v *nestedEnv) Assign(name token.T, value Value) *RuntimeError {
	if _, ok := v.Values[name.Lexeme()]; ok {
		v.Values[name.Lexeme()] = value
		return nil
	}

	return v.Enclosing.Assign(name, value)
}

func (v *nestedEnv) AssignAt(
	distance int,
	name token.T,
	value Value,
) *RuntimeError {
	if distance <= 0 {
		return v.Assign(name, value)
	}
	return v.Enclosing.AssignAt(distance-1, name, value)
}

func (v *nestedEnv) GetAt(distance int, name string) (Value, error) {
	if distance > 0 {
		return v.Enclosing.GetAt(distance-1, name)
	}

	if v, ok := v.Values[name]; ok {
		return v, nil
	}
	return &token.NilValue{}, fmt.Errorf("Undefined variable '%s'.", name)
}

func (v *globalEnv) Dump(msg string) {

	if msg == "" {
		msg = "Global Environment"
	}
	fmt.Printf("===== %s =====\n", msg)
	for key, val := range v.Values {
		fmt.Printf("  %s: %s\n", key, val)
	}
	fmt.Print("==============\n\n\n")
}

func (v *nestedEnv) Dump(msg string) {

	if msg == "" {
		msg = "Nested Environment"
	}
	fmt.Printf("===== %s =====\n", msg)
	for key, val := range v.Values {
		fmt.Printf("  %s: %s\n", key, val)
	}
	v.Enclosing.Dump(msg)
	fmt.Print("==============\n")
}
