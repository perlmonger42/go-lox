package ast

import (
	"strings"

	"github.com/perlmonger42/go-lox/token"
)

func StmtToString(stmt Stmt) string {
	return stmt.Accept_Stmt_String(&stmtToStringVisitor{})
}

type stmtToStringVisitor struct {
	level int
}

var _ Visitor_Stmt_String = &stmtToStringVisitor{}

func (x *stmtToStringVisitor) indent() {
	x.level++
}

func (x *stmtToStringVisitor) undent() {
	x.level--
}

func (x *stmtToStringVisitor) indentation() string {
	return strings.Join(make([]string, x.level+1), "  ")
}

func (x *stmtToStringVisitor) statementToString(stmt Stmt) string {
	return stmt.Accept_Stmt_String(x)
}

func (x *stmtToStringVisitor) idsToString(ids []token.T) string {
	params := []string{}
	for _, id := range ids {
		params = append(params, id.Lexeme())
	}
	return strings.Join(params, ", ")
}

// blockToString returns a string that begins with "{\n" and ends with "}".
// All lines but the first have their proper indentation.
func (x *stmtToStringVisitor) blockToString(stmts []Stmt) string {
	var body = "{\n"
	x.indent()
	for _, stmt := range stmts {
		body += x.statementToString(stmt)
	}
	x.undent()
	return body + x.indentation() + "}"
}

// ifStatementToString returns a string that begins with "if (" and ends with "\n".
// All lines but the first have their proper indentation.
func (x *stmtToStringVisitor) ifStatementToString(stmt *If) string {
	var s string = "if (" + ExprToString(stmt.Condition) + ")"

	var thenBranchIsABlock = false
	if block, ok := stmt.ThenBranch.(*Block); ok {
		thenBranchIsABlock = true
		s += " " + x.blockToString(block.Statements)
		if stmt.ElseBranch != nil {
			s += " "
		}
	} else {
		s += "\n"
		x.indent()
		s += x.statementToString(stmt.ThenBranch)
		x.undent()
	}

	if stmt.ElseBranch == nil {
		if thenBranchIsABlock {
			s += "\n"
		}
	} else if ifStmt, ok := stmt.ElseBranch.(*If); ok {
		s += "else " + x.ifStatementToString(ifStmt)
	} else if blk, ok := stmt.ElseBranch.(*Block); ok {
		s += "else " + x.Visit_BlockStmt_String(blk) + "\n"
	} else {
		s += "else\n"
		x.indent()
		s += x.statementToString(stmt.ElseBranch)
		x.undent()
	}
	return s
}

func (x *stmtToStringVisitor) Visit_NoopStmt_String(stmt *Noop) string {
	return x.indentation() + "/* continue */;\n"
}

func (x *stmtToStringVisitor) Visit_ExpressionStmt_String(stmt *Expression) string {
	return x.indentation() + ExprToString(stmt.Expression) + ";\n"
}

func (x *stmtToStringVisitor) Visit_PrintStmt_String(stmt *Print) string {
	return x.indentation() + "print " + ExprToString(stmt.Expression) + ";\n"
}

func (x *stmtToStringVisitor) Visit_ReturnStmt_String(stmt *Return) string {
	return x.indentation() + "return " + ExprToString(stmt.Value) + ";\n"
}

func (x *stmtToStringVisitor) Visit_PanicStmt_String(stmt *Panic) string {
	return x.indentation() + "panic " + ExprToString(stmt.Expression) + ";\n"
}

func (x *stmtToStringVisitor) Visit_VarInitializedStmt_String(stmt *VarInitialized) string {
	return x.indentation() + "var " + stmt.Name.Lexeme() + " = " +
		ExprToString(stmt.Initializer) + ";\n"
}

func (x *stmtToStringVisitor) Visit_VarUninitializedStmt_String(stmt *VarUninitialized) string {
	return x.indentation() + "var " + stmt.Name.Lexeme() + ";\n"
}

func (x *stmtToStringVisitor) Visit_FunctionStmt_String(stmt *Function) string {
	params := x.idsToString(stmt.Params)
	return x.indentation() + "fun (" + params + ") " +
		x.blockToString(stmt.Body) + "\n"
}

func (x *stmtToStringVisitor) Visit_BlockStmt_String(stmt *Block) string {
	return x.indentation() + x.blockToString(stmt.Statements) + "\n"
}

func (x *stmtToStringVisitor) Visit_IfStmt_String(stmt *If) string {
	return x.indentation() + x.ifStatementToString(stmt)
}

func (x *stmtToStringVisitor) Visit_WhileStmt_String(stmt *While) string {
	s := "while (" + ExprToString(stmt.Condition) + ")"
	if blkStmt, ok := stmt.Body.(*Block); ok {
		return s + " " + x.blockToString(blkStmt.Statements) + "\n"
	}
	x.indent()
	s += "\n" + x.statementToString(stmt.Body)
	x.undent()
	return s
}

func (x *stmtToStringVisitor) Visit_ClassStmt_String(stmt *Class) string {
	methods := []string{}
	x.indent()
	for _, method := range stmt.Methods {
		text := x.indentation() + method.Name.Lexeme() + " " +
			"(" + x.idsToString(method.Params) + ")" +
			x.blockToString(method.Body)
		methods = append(methods, text)
	}
	x.undent()

	return x.indentation() + "class " + "{\n" +
		strings.Join(methods, "") +
		x.indentation() + "}\n"
}
