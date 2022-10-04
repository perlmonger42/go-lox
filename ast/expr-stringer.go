package ast

import (
	"strings"
)

func ExprToString(expr Expr) string {
	return expr.Accept_Expr_String(&toStringVisitor{})
}

type toStringVisitor struct {
}

var _ Visitor_Expr_String = &toStringVisitor{}

func (x *toStringVisitor) ToString(expr Expr) string {
	return expr.Accept_Expr_String(x)
}

func (x *toStringVisitor) Visit_LogicalExpr_String(expr *Logical) string {
	return x.parenthesize(expr.Operator.Lexeme(), expr.Left, expr.Right)
}

func (x *toStringVisitor) Visit_BinaryExpr_String(expr *Binary) string {
	return x.parenthesize(expr.Operator.Lexeme(), expr.Left, expr.Right)
}

func (x *toStringVisitor) Visit_GroupingExpr_String(expr *Grouping) string {
	return x.parenthesize("group", expr.Expression)
}

func (x *toStringVisitor) Visit_LiteralExpr_String(expr *Literal) string {
	if expr.Value == nil {
		return "nil"
	}
	return expr.Value.Show()
}

func (x *toStringVisitor) Visit_CallExpr_String(expr *Call) string {
	args := []string{}
	for _, a := range expr.Arguments {
		args = append(args, ExprToString(a))
	}
	return ExprToString(expr.Callee) + "(" + strings.Join(args, ", ") + ")"
}

func (x *toStringVisitor) Visit_GetExpr_String(expr *Get) string {
	return "(" + ExprToString(expr.Object) + ")." + expr.Name.Lexeme()
}

func (x *toStringVisitor) Visit_SetExpr_String(expr *Set) string {
	return "(" + ExprToString(expr.Object) + ")." + expr.Name.Lexeme() +
		" = " + ExprToString(expr.Value)
}

func (x *toStringVisitor) Visit_UnaryExpr_String(expr *Unary) string {
	return x.parenthesize(expr.Operator.Lexeme(), expr.Right)
}

func (x *toStringVisitor) Visit_ThisExpr_String(expr *This) string {
	return "this"
}

func (x *toStringVisitor) Visit_SuperExpr_String(expr *Super) string {
	return "super." + expr.Method.Lexeme()
}

func (x *toStringVisitor) Visit_VariableExpr_String(expr *Variable) string {
	return expr.Name.Lexeme()
}

func (x *toStringVisitor) Visit_AssignExpr_String(expr *Assign) string {
	return expr.Name.Lexeme() + " = " + ExprToString(expr.Value) + ";"
}

func (x *toStringVisitor) parenthesize(name string, exprs ...Expr) string {
	var str strings.Builder

	str.WriteString("(")
	str.WriteString(name)
	for _, expr := range exprs {
		str.WriteString(" ")
		str.WriteString(expr.Accept_Expr_String(x))
	}
	str.WriteString(")")
	return str.String()
}
