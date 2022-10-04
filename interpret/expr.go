package interpret

import (
	"fmt"
	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/token"
)

var _ ast.Visitor_Expr_Token_Value = &Interpreter{}

func (i *Interpreter) InterpretExpr(expr ast.Expr) (result Value) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(RuntimeError); ok {
				result = nil
			} else {
				panic(r)
			}
		}
	}()
	i.depth = 0
	return i.evaluate(expr)
}

type Callable interface {
	Call(i T, arguments []token.Value) token.Value
	Arity() int
}
type TestCallable struct{}

func (t *TestCallable) Call(i T, arguments []token.Value) token.Value { return nil }
func (t *TestCallable) Arity() int                                    { return 0 }

var _ Callable = &TestCallable{}

func (i *Interpreter) indent() string {
	s := ""
	for n := 0; n < i.depth; n++ {
		s = s + "  "
	}
	return s
}

func (i *Interpreter) printIndent() {
	fmt.Print(i.indent())
}

func (i *Interpreter) evaluate(expr ast.Expr) Value {
	i.depth++
	value := expr.Accept_Expr_Token_Value(i)
	i.depth--

	if i.lox.Config.TraceEval {
		fmt.Printf("%s%s <-- %s\n", i.indent(), value, ast.ToString(expr))
	}
	return value
}

func (i *Interpreter) Visit_LiteralExpr_Token_Value(expr *ast.Literal) Value {
	return expr.Value
}

func (i *Interpreter) Visit_VariableExpr_Token_Value(expr *ast.Variable) Value {
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) Visit_ThisExpr_Token_Value(expr *ast.This) Value {
	return i.lookUpVariable(expr.Keyword, expr)
}

func (i *Interpreter) Visit_SuperExpr_Token_Value(expr *ast.Super) Value {
	distance, superclass := i.GetSuper(expr)

	// "this" is always bound just inside the environment that binds "super"
	object := i.GetThisAt(distance - 1)

	var method *LoxFunction
	var ok bool
	if method, ok = superclass.FindMethod(expr.Method.Lexeme()); !ok {
		panic(i.Error(expr.Method,
			fmt.Sprintf("Undefined property '%s'.",
				expr.Method.Lexeme())))
	} else {
		return token.ObjectValue{method.Bind(object)}
	}
}

func (i *Interpreter) Visit_GroupingExpr_Token_Value(expr *ast.Grouping) Value {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) Visit_CallExpr_Token_Value(expr *ast.Call) Value {
	var callee token.Value = i.evaluate(expr.Callee)
	arguments := []token.Value{}
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}
	function := i.GetCallable(expr.Paren, callee)
	if len(expr.Arguments) != function.Arity() {
		panic(i.Error(expr.Paren,
			fmt.Sprintf("expected %d arguments but got %d.",
				function.Arity(),
				len(expr.Arguments))))
	}
	return function.Call(i, arguments)
}

func (i *Interpreter) Visit_GetExpr_Token_Value(expr *ast.Get) Value {
	var lhs Value = i.evaluate(expr.Object)
	if obj, ok := lhs.(token.ObjectValue); ok {
		if instance, ok := obj.V.(*LoxInstance); ok {
			if value, err := instance.Get(expr.Name); err != nil {
				panic(i.Error(expr.Name, err.Error()))
			} else {
				return value
			}
		}
	}
	panic(i.Error(expr.Name, "Only instances have properties."))
}

func (i *Interpreter) Visit_SetExpr_Token_Value(expr *ast.Set) Value {
	var lhs Value = i.evaluate(expr.Object)
	if obj, ok := lhs.(token.ObjectValue); ok {
		if instance, ok := obj.V.(*LoxInstance); ok {
			var rhs Value = i.evaluate(expr.Value)
			instance.Set(expr.Name, rhs)
			return rhs
		}
	}
	panic(i.Error(expr.Name, "Only instances have fields."))
}

func (i *Interpreter) GetCallable(paren token.T, v token.Value) Callable {
	if obj, ok := v.(token.ObjectValue); ok {
		if callable, ok := obj.V.(Callable); ok {
			return callable
		} else {
			panic(i.Error(paren, fmt.Sprintf(
				"Can only call functions and classes (got %T; want Callable).",
				obj.V)))
		}
	} else {
		panic(i.Error(paren, fmt.Sprintf(
			"Can only call functions and classes (got %T; want ObjectValue).",
			v)))
	}
}

func (i *Interpreter) Visit_UnaryExpr_Token_Value(expr *ast.Unary) Value {
	var right Value = i.evaluate(expr.Right)
	switch expr.Operator.Type() {
	case token.Minus:
		switch r := right.(type) {
		case token.NumberValue:
			return token.NumberValue{-r.V}
		}
	case token.Bang:
		switch r := right.(type) {
		case token.BooleanValue:
			return token.BooleanValue{!r.V}
		case token.NilValue:
			return token.BooleanValue{true}
		}
	}
	panic(i.Error(
		expr.Operator,
		fmt.Sprintf("cannot apply %s to type %s (%v) (%T)",
			expr.Operator, right.TypeName(), right.Show(), right)),
	)

	// Unreachable.
	return nil
}

func (i *Interpreter) Visit_BinaryExpr_Token_Value(expr *ast.Binary) Value {
	var left Value = i.evaluate(expr.Left)
	var right Value = i.evaluate(expr.Right)
	switch expr.Operator.Type() {
	case token.Minus:
		switch l := left.(type) {
		case token.NumberValue:
			switch r := right.(type) {
			case token.NumberValue:
				return token.NumberValue{l.V - r.V}
			}
		}
	case token.Slash:
		switch l := left.(type) {
		case token.NumberValue:
			switch r := right.(type) {
			case token.NumberValue:
				return token.NumberValue{l.V / r.V}
			}
		}
	case token.Star:
		switch l := left.(type) {
		case token.NumberValue:
			switch r := right.(type) {
			case token.NumberValue:
				return token.NumberValue{l.V * r.V}
			}
		}
	case token.Plus:
		switch l := left.(type) {
		case token.NumberValue:
			switch r := right.(type) {
			case token.NumberValue:
				return token.NumberValue{l.V + r.V}
			}
		case token.StringValue:
			switch r := right.(type) {
			case token.StringValue:
				return token.StringValue{l.V + r.V}
			case token.NilValue:
				return token.StringValue{l.V + "{([<nil>])}"}
			}
		}
	case token.Greater:
		switch l := left.(type) {
		case token.NumberValue:
			switch r := right.(type) {
			case token.NumberValue:
				return token.BooleanValue{l.V > r.V}
			}
		case token.StringValue:
			switch r := right.(type) {
			case token.StringValue:
				return token.BooleanValue{l.V > r.V}
			}
		}
	case token.GreaterEqual:
		switch l := left.(type) {
		case token.NumberValue:
			switch r := right.(type) {
			case token.NumberValue:
				return token.BooleanValue{l.V >= r.V}
			}
		case token.StringValue:
			switch r := right.(type) {
			case token.StringValue:
				return token.BooleanValue{l.V >= r.V}
			}
		}
	case token.Less:
		switch l := left.(type) {
		case token.NumberValue:
			switch r := right.(type) {
			case token.NumberValue:
				return token.BooleanValue{l.V < r.V}
			}
		case token.StringValue:
			switch r := right.(type) {
			case token.StringValue:
				return token.BooleanValue{l.V < r.V}
			}
		}
	case token.LessEqual:
		switch l := left.(type) {
		case token.NumberValue:
			switch r := right.(type) {
			case token.NumberValue:
				return token.BooleanValue{l.V <= r.V}
			}
		case token.StringValue:
			switch r := right.(type) {
			case token.StringValue:
				return token.BooleanValue{l.V <= r.V}
			}
		}
	case token.BangEqual:
		return token.BooleanValue{V: !left.IsEqualTo(right)}
	case token.EqualEqual:
		return token.BooleanValue{V: left.IsEqualTo(right)}
	}
	panic(i.Error(
		expr.Operator,
		fmt.Sprintf(
			"cannot apply %s to types %s and %s (values %s and %s) (%T and %T)",
			expr.Operator, left.TypeName(), right.TypeName(),
			left.Show(), right.Show(),
			left, right)),
	)

	// Unreachable.
	return nil
}

func (i *Interpreter) Visit_LogicalExpr_Token_Value(expr *ast.Logical) Value {
	var left Value = i.evaluate(expr.Left)

	if expr.Operator.Type() == token.Or {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) Visit_AssignExpr_Token_Value(expr *ast.Assign) Value {
	var value Value = i.evaluate(expr.Value)

	if distance, ok := i.locals[expr]; ok {
		i.assignLocal(distance, expr.Name, value)
	} else {
		i.assignGlobal(expr.Name, value)
	}

	return value
}
