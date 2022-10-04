package ast

import (
	"github.com/perlmonger42/go-lox/token"
)

type Stmt interface {
	AsNode() Node // does nothing but prevent non-Nodes from looking like Nodes
	AsStmt() Stmt // does nothing but prevent non-Stmts from looking like Stmt

	Accept_Stmt(visitor Visitor_Stmt)
	Accept_Stmt_String(visitor Visitor_Stmt_String) string
	Accept_Stmt_Error(visitor Visitor_Stmt_Error) error
}

// A Visitor_Stmt is accepted by Stmt and has no return value
type Visitor_Stmt interface {
	Visit_NoopStmt(stmt *Noop)
	Visit_ExpressionStmt(stmt *Expression)
	Visit_PrintStmt(stmt *Print)
	Visit_ReturnStmt(stmt *Return)
	Visit_PanicStmt(stmt *Panic)
	Visit_VarInitializedStmt(stmt *VarInitialized)
	Visit_VarUninitializedStmt(stmt *VarUninitialized)
	Visit_FunctionStmt(stmt *Function)
	Visit_IfStmt(stmt *If)
	Visit_BlockStmt(stmt *Block)
	Visit_WhileStmt(stmt *While)
	Visit_ClassStmt(stmt *Class)
}

// A Visitor_Stmt_String is accepted by Stmt and returns string
type Visitor_Stmt_String interface {
	Visit_NoopStmt_String(stmt *Noop) string
	Visit_ExpressionStmt_String(stmt *Expression) string
	Visit_PrintStmt_String(stmt *Print) string
	Visit_ReturnStmt_String(stmt *Return) string
	Visit_PanicStmt_String(stmt *Panic) string
	Visit_VarInitializedStmt_String(stmt *VarInitialized) string
	Visit_VarUninitializedStmt_String(stmt *VarUninitialized) string
	Visit_FunctionStmt_String(stmt *Function) string
	Visit_IfStmt_String(stmt *If) string
	Visit_BlockStmt_String(stmt *Block) string
	Visit_WhileStmt_String(stmt *While) string
	Visit_ClassStmt_String(stmt *Class) string
}

// A Visitor_Stmt_Error is accepted by Stmt and returns error
type Visitor_Stmt_Error interface {
	Visit_NoopStmt_Error(stmt *Noop) error
	Visit_ExpressionStmt_Error(stmt *Expression) error
	Visit_PrintStmt_Error(stmt *Print) error
	Visit_ReturnStmt_Error(stmt *Return) error
	Visit_PanicStmt_Error(stmt *Panic) error
	Visit_VarInitializedStmt_Error(stmt *VarInitialized) error
	Visit_VarUninitializedStmt_Error(stmt *VarUninitialized) error
	Visit_FunctionStmt_Error(stmt *Function) error
	Visit_IfStmt_Error(stmt *If) error
	Visit_BlockStmt_Error(stmt *Block) error
	Visit_WhileStmt_Error(stmt *While) error
	Visit_ClassStmt_Error(stmt *Class) error
}

type Noop struct {
}

func (x *Noop) AsNode() Node { return x }
func (x *Noop) AsStmt() Stmt { return x }

func (x *Noop) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_NoopStmt(x)
}
func (x *Noop) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_NoopStmt_String(x)
}
func (x *Noop) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_NoopStmt_Error(x)
}

type Expression struct {
	Expression Expr
}

func (x *Expression) AsNode() Node { return x }
func (x *Expression) AsStmt() Stmt { return x }

func (x *Expression) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_ExpressionStmt(x)
}
func (x *Expression) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_ExpressionStmt_String(x)
}
func (x *Expression) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_ExpressionStmt_Error(x)
}

type Print struct {
	Keyword    token.T
	Expression Expr
}

func (x *Print) AsNode() Node { return x }
func (x *Print) AsStmt() Stmt { return x }

func (x *Print) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_PrintStmt(x)
}
func (x *Print) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_PrintStmt_String(x)
}
func (x *Print) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_PrintStmt_Error(x)
}

type Return struct {
	Keyword token.T
	Value   Expr
}

func (x *Return) AsNode() Node { return x }
func (x *Return) AsStmt() Stmt { return x }

func (x *Return) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_ReturnStmt(x)
}
func (x *Return) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_ReturnStmt_String(x)
}
func (x *Return) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_ReturnStmt_Error(x)
}

type Panic struct {
	Keyword    token.T
	Expression Expr
}

func (x *Panic) AsNode() Node { return x }
func (x *Panic) AsStmt() Stmt { return x }

func (x *Panic) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_PanicStmt(x)
}
func (x *Panic) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_PanicStmt_String(x)
}
func (x *Panic) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_PanicStmt_Error(x)
}

type VarInitialized struct {
	Name        token.T
	Initializer Expr
}

func (x *VarInitialized) AsNode() Node { return x }
func (x *VarInitialized) AsStmt() Stmt { return x }

func (x *VarInitialized) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_VarInitializedStmt(x)
}
func (x *VarInitialized) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_VarInitializedStmt_String(x)
}
func (x *VarInitialized) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_VarInitializedStmt_Error(x)
}

type VarUninitialized struct {
	Name token.T
}

func (x *VarUninitialized) AsNode() Node { return x }
func (x *VarUninitialized) AsStmt() Stmt { return x }

func (x *VarUninitialized) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_VarUninitializedStmt(x)
}
func (x *VarUninitialized) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_VarUninitializedStmt_String(x)
}
func (x *VarUninitialized) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_VarUninitializedStmt_Error(x)
}

type Function struct {
	Name   token.T
	Params []token.T
	Body   []Stmt
}

func (x *Function) AsNode() Node { return x }
func (x *Function) AsStmt() Stmt { return x }

func (x *Function) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_FunctionStmt(x)
}
func (x *Function) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_FunctionStmt_String(x)
}
func (x *Function) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_FunctionStmt_Error(x)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (x *If) AsNode() Node { return x }
func (x *If) AsStmt() Stmt { return x }

func (x *If) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_IfStmt(x)
}
func (x *If) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_IfStmt_String(x)
}
func (x *If) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_IfStmt_Error(x)
}

type Block struct {
	Token      token.T
	Statements []Stmt
}

func (x *Block) AsNode() Node { return x }
func (x *Block) AsStmt() Stmt { return x }

func (x *Block) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_BlockStmt(x)
}
func (x *Block) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_BlockStmt_String(x)
}
func (x *Block) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_BlockStmt_Error(x)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func (x *While) AsNode() Node { return x }
func (x *While) AsStmt() Stmt { return x }

func (x *While) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_WhileStmt(x)
}
func (x *While) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_WhileStmt_String(x)
}
func (x *While) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_WhileStmt_Error(x)
}

type Class struct {
	Name       token.T
	Superclass *Variable
	Methods    []*Function
}

func (x *Class) AsNode() Node { return x }
func (x *Class) AsStmt() Stmt { return x }

func (x *Class) Accept_Stmt(visitor Visitor_Stmt) {
	visitor.Visit_ClassStmt(x)
}
func (x *Class) Accept_Stmt_String(visitor Visitor_Stmt_String) string {
	return visitor.Visit_ClassStmt_String(x)
}
func (x *Class) Accept_Stmt_Error(visitor Visitor_Stmt_Error) error {
	return visitor.Visit_ClassStmt_Error(x)
}
