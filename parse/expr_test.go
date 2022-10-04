package parse

import (
	"fmt"

	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/config"
	"github.com/perlmonger42/go-lox/lox"
	"github.com/perlmonger42/go-lox/scan"
)

func dumpAst(text string) {
	config := config.New()
	lox := lox.New(config)
	scanner := scan.New(lox, text)
	tokens := scanner.ScanTokens()
	parser := New(lox, tokens)
	expr := parser.ParseExpr()

	if lox.HadError {
		return
	}
	fmt.Printf("%s\n", ast.ToString(expr))
}

func ExampleEmpty() {
	dumpAst("")
	// Output:
	// [line 1] Error at end: Expect expression.
}

func ExampleParens() {
	dumpAst("()")
	// Output:
	// [line 1] Error at 'RightParen': Expect expression.
}

func ExampleUnexpectedInput() {
	dumpAst("(~")
	// Output:
	// [line 1] Error at 'Other': Unexpected character ('~').
	// [line 1] Error at 'Other': Expect expression.
}

func ExampleNumber() {
	dumpAst("7")
	// Output:
	// 7
}

func Example4567() {
	dumpAst("4 - -5 * ( 6 * 7)")
	// Output:
	// (- 4 (* (- 5) (group (* 6 7))))
}
func ExampleComplex() {
	dumpAst("1+2   * 3 / 4 - -5 * ( 6 * 7)")
	// Output:
	// (- (+ 1 (/ (* 2 3) 4)) (* (- 5) (group (* 6 7))))
}

//func ExampleUnicode() {
//	// "Ťėšťǐňġ" is entirely made of 2-byte UTF-8 encodings
//	// "ṫẹṡṫịṅḡ" is entirely made of 3-byte UTF-8 encodings
//	// "𝕠𝕟𝕖, 𝕥𝕨𝕠, 𝕥𝕙𝕣𝕖𝕖" has a bunch of 4-byte UTF-8 encodings
//	//dumpAst("Ťėšťǐňġ + ṫẹṡṫịṅḡ * 𝕠𝕟𝕖 / 𝕥𝕨𝕠 - -𝕥𝕙𝕣𝕖𝕖")
//	// Output:
//}
