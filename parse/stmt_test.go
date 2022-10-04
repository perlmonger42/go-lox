package parse

import (
	"fmt"

	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/config"
	"github.com/perlmonger42/go-lox/lox"
	"github.com/perlmonger42/go-lox/scan"
)

func dumpProgram(text string) {
	config := config.New()
	lox := lox.New(config)
	scanner := scan.New(lox, text)
	tokens := scanner.ScanTokens()
	parser := New(lox, tokens)
	stmts := parser.ParseProg()
	for i, stmt := range stmts {
		fmt.Printf("%2d: %s", i+1, ast.ToString(stmt))
	}
}

func ExampleEmptyProgram() {
	dumpProgram("")
	// Output:
}

func ExampleMissingSemicolon() {
	dumpProgram("7")
	// Output:
	// [line 1] Error at end: found EOF; Expect `;` after expression statement.
	//  1: panic "parse error at line 1; this code shouldn't be run";
}

func ExampleExpressionStatement() {
	dumpProgram("40 + 2;")
	// Output:
	// 1: (+ 40 2);
}

func ExampleTwoStatements() {
	dumpProgram("42; print 6 * 9;")
	// Output:
	//  1: 42;
	//  2: print (* 6 9);
}

func ExampleSyncOnSemicolon() {
	dumpProgram("7+;\nprint 1;\n);\nprint 2;")
	// Output:
	// [line 1] Error at 'Semicolon': Expect expression.
	// [line 3] Error at 'RightParen': Expect expression.
	//  1: panic "parse error at line 1; this code shouldn't be run";
	//  2: print 1;
	//  3: panic "parse error at line 3; this code shouldn't be run";
	//  4: print 2;
}
