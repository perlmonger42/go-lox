package ast

import (
	"github.com/perlmonger42/go-lox/token"
)

type Expr interface {
	AsNode() Node // does nothing but prevent non-Nodes from looking like Nodes
	AsExpr() Expr // does nothing but prevent non-Exprs from looking like Expr

	Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value
	Accept_Expr_String(visitor Visitor_Expr_String) string
	Accept_Expr(visitor Visitor_Expr)
	Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error)
}

// A Visitor_Expr_Token_Value is accepted by Expr and returns token.Value
type Visitor_Expr_Token_Value interface {
	Visit_GroupingExpr_Token_Value(expr *Grouping) token.Value
	Visit_ThisExpr_Token_Value(expr *This) token.Value
	Visit_SuperExpr_Token_Value(expr *Super) token.Value
	Visit_VariableExpr_Token_Value(expr *Variable) token.Value
	Visit_LiteralExpr_Token_Value(expr *Literal) token.Value
	Visit_CallExpr_Token_Value(expr *Call) token.Value
	Visit_GetExpr_Token_Value(expr *Get) token.Value
	Visit_UnaryExpr_Token_Value(expr *Unary) token.Value
	Visit_BinaryExpr_Token_Value(expr *Binary) token.Value
	Visit_LogicalExpr_Token_Value(expr *Logical) token.Value
	Visit_SetExpr_Token_Value(expr *Set) token.Value
	Visit_AssignExpr_Token_Value(expr *Assign) token.Value
}

// A Visitor_Expr_String is accepted by Expr and returns string
type Visitor_Expr_String interface {
	Visit_GroupingExpr_String(expr *Grouping) string
	Visit_ThisExpr_String(expr *This) string
	Visit_SuperExpr_String(expr *Super) string
	Visit_VariableExpr_String(expr *Variable) string
	Visit_LiteralExpr_String(expr *Literal) string
	Visit_CallExpr_String(expr *Call) string
	Visit_GetExpr_String(expr *Get) string
	Visit_UnaryExpr_String(expr *Unary) string
	Visit_BinaryExpr_String(expr *Binary) string
	Visit_LogicalExpr_String(expr *Logical) string
	Visit_SetExpr_String(expr *Set) string
	Visit_AssignExpr_String(expr *Assign) string
}

// A Visitor_Expr is accepted by Expr and has no return value
type Visitor_Expr interface {
	Visit_GroupingExpr(expr *Grouping)
	Visit_ThisExpr(expr *This)
	Visit_SuperExpr(expr *Super)
	Visit_VariableExpr(expr *Variable)
	Visit_LiteralExpr(expr *Literal)
	Visit_CallExpr(expr *Call)
	Visit_GetExpr(expr *Get)
	Visit_UnaryExpr(expr *Unary)
	Visit_BinaryExpr(expr *Binary)
	Visit_LogicalExpr(expr *Logical)
	Visit_SetExpr(expr *Set)
	Visit_AssignExpr(expr *Assign)
}

// A Visitor_Expr_MaybeValue is accepted by Expr and returns (token.Value, error)
type Visitor_Expr_MaybeValue interface {
	Visit_GroupingExpr_MaybeValue(expr *Grouping) (token.Value, error)
	Visit_ThisExpr_MaybeValue(expr *This) (token.Value, error)
	Visit_SuperExpr_MaybeValue(expr *Super) (token.Value, error)
	Visit_VariableExpr_MaybeValue(expr *Variable) (token.Value, error)
	Visit_LiteralExpr_MaybeValue(expr *Literal) (token.Value, error)
	Visit_CallExpr_MaybeValue(expr *Call) (token.Value, error)
	Visit_GetExpr_MaybeValue(expr *Get) (token.Value, error)
	Visit_UnaryExpr_MaybeValue(expr *Unary) (token.Value, error)
	Visit_BinaryExpr_MaybeValue(expr *Binary) (token.Value, error)
	Visit_LogicalExpr_MaybeValue(expr *Logical) (token.Value, error)
	Visit_SetExpr_MaybeValue(expr *Set) (token.Value, error)
	Visit_AssignExpr_MaybeValue(expr *Assign) (token.Value, error)
}

type Grouping struct {
	Expression Expr
}

func (x *Grouping) AsNode() Node { return x }
func (x *Grouping) AsExpr() Expr { return x }

func (x *Grouping) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_GroupingExpr_Token_Value(x)
}
func (x *Grouping) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_GroupingExpr_String(x)
}
func (x *Grouping) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_GroupingExpr(x)
}
func (x *Grouping) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_GroupingExpr_MaybeValue(x)
}

type This struct {
	Keyword token.T
}

func (x *This) AsNode() Node { return x }
func (x *This) AsExpr() Expr { return x }

func (x *This) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_ThisExpr_Token_Value(x)
}
func (x *This) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_ThisExpr_String(x)
}
func (x *This) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_ThisExpr(x)
}
func (x *This) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_ThisExpr_MaybeValue(x)
}

type Super struct {
	Keyword token.T
	Method  token.T
}

func (x *Super) AsNode() Node { return x }
func (x *Super) AsExpr() Expr { return x }

func (x *Super) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_SuperExpr_Token_Value(x)
}
func (x *Super) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_SuperExpr_String(x)
}
func (x *Super) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_SuperExpr(x)
}
func (x *Super) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_SuperExpr_MaybeValue(x)
}

type Variable struct {
	Name token.T
}

func (x *Variable) AsNode() Node { return x }
func (x *Variable) AsExpr() Expr { return x }

func (x *Variable) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_VariableExpr_Token_Value(x)
}
func (x *Variable) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_VariableExpr_String(x)
}
func (x *Variable) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_VariableExpr(x)
}
func (x *Variable) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_VariableExpr_MaybeValue(x)
}

type Literal struct {
	Value token.Value
}

func (x *Literal) AsNode() Node { return x }
func (x *Literal) AsExpr() Expr { return x }

func (x *Literal) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_LiteralExpr_Token_Value(x)
}
func (x *Literal) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_LiteralExpr_String(x)
}
func (x *Literal) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_LiteralExpr(x)
}
func (x *Literal) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_LiteralExpr_MaybeValue(x)
}

type Call struct {
	Callee    Expr
	Paren     token.T
	Arguments []Expr
}

func (x *Call) AsNode() Node { return x }
func (x *Call) AsExpr() Expr { return x }

func (x *Call) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_CallExpr_Token_Value(x)
}
func (x *Call) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_CallExpr_String(x)
}
func (x *Call) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_CallExpr(x)
}
func (x *Call) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_CallExpr_MaybeValue(x)
}

type Get struct {
	Object Expr
	Name   token.T
}

func (x *Get) AsNode() Node { return x }
func (x *Get) AsExpr() Expr { return x }

func (x *Get) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_GetExpr_Token_Value(x)
}
func (x *Get) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_GetExpr_String(x)
}
func (x *Get) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_GetExpr(x)
}
func (x *Get) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_GetExpr_MaybeValue(x)
}

type Unary struct {
	Operator token.T
	Right    Expr
}

func (x *Unary) AsNode() Node { return x }
func (x *Unary) AsExpr() Expr { return x }

func (x *Unary) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_UnaryExpr_Token_Value(x)
}
func (x *Unary) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_UnaryExpr_String(x)
}
func (x *Unary) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_UnaryExpr(x)
}
func (x *Unary) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_UnaryExpr_MaybeValue(x)
}

type Binary struct {
	Operator token.T
	Left     Expr
	Right    Expr
}

func (x *Binary) AsNode() Node { return x }
func (x *Binary) AsExpr() Expr { return x }

func (x *Binary) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_BinaryExpr_Token_Value(x)
}
func (x *Binary) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_BinaryExpr_String(x)
}
func (x *Binary) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_BinaryExpr(x)
}
func (x *Binary) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_BinaryExpr_MaybeValue(x)
}

type Logical struct {
	Operator token.T
	Left     Expr
	Right    Expr
}

func (x *Logical) AsNode() Node { return x }
func (x *Logical) AsExpr() Expr { return x }

func (x *Logical) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_LogicalExpr_Token_Value(x)
}
func (x *Logical) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_LogicalExpr_String(x)
}
func (x *Logical) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_LogicalExpr(x)
}
func (x *Logical) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_LogicalExpr_MaybeValue(x)
}

type Set struct {
	Object Expr
	Name   token.T
	Value  Expr
}

func (x *Set) AsNode() Node { return x }
func (x *Set) AsExpr() Expr { return x }

func (x *Set) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_SetExpr_Token_Value(x)
}
func (x *Set) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_SetExpr_String(x)
}
func (x *Set) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_SetExpr(x)
}
func (x *Set) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_SetExpr_MaybeValue(x)
}

type Assign struct {
	Name  token.T
	Value Expr
}

func (x *Assign) AsNode() Node { return x }
func (x *Assign) AsExpr() Expr { return x }

func (x *Assign) Accept_Expr_Token_Value(visitor Visitor_Expr_Token_Value) token.Value {
	return visitor.Visit_AssignExpr_Token_Value(x)
}
func (x *Assign) Accept_Expr_String(visitor Visitor_Expr_String) string {
	return visitor.Visit_AssignExpr_String(x)
}
func (x *Assign) Accept_Expr(visitor Visitor_Expr) {
	visitor.Visit_AssignExpr(x)
}
func (x *Assign) Accept_Expr_MaybeValue(visitor Visitor_Expr_MaybeValue) (token.Value, error) {
	return visitor.Visit_AssignExpr_MaybeValue(x)
}
