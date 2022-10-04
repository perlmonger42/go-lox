package interpret

import (
	"fmt"

	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/config"
	"github.com/perlmonger42/go-lox/lox"
	"github.com/perlmonger42/go-lox/parse"
	"github.com/perlmonger42/go-lox/scan"
)

var traceEval bool
var traceTokens bool
var traceParsed bool

func eval(text string) {
	config := config.New()
	if traceEval {
		config.TraceEval = true
	}
	if traceTokens {
		config.TraceParseTokens = true
	}
	lox := lox.New(config)
	scanner := scan.New(lox, text)
	tokens := scanner.ScanTokens()
	if lox.HadError {
		return
	}
	parser := parse.New(lox, tokens)
	expr := parser.ParseExpr()
	if lox.HadError {
		return
	}
	if traceParsed {
		fmt.Printf("parsed: %s\n", ast.ToString(expr))
	}
	interpreter := New(lox)
	value := interpreter.InterpretExpr(expr)
	if lox.HadError {
		return
	}
	if value != nil {
		fmt.Println(value.Show())
	}
}

func ExampleZero() {
	eval("0")
	// Output:
	// 0
}

func ExampleParens() {
	eval("(1)")
	// Output:
	// 1
}

func ExampleTypeMismatch() {
	eval(`7+""`)
	eval(`nil-4`)
	// Output:
	// [line 1] Error at 'Plus': cannot apply Plus: `+` to types number and string (values 7 and "") (token.NumberValue and token.StringValue)
	// [line 1] Error at 'Minus': cannot apply Minus: `-` to types nil and number (values nil and 4) (token.NilValue and token.NumberValue)
}

func ExampleBinary() {
	eval("3+7")
	eval("true + nil")
	eval("6/false")
	eval("7/2")
	eval("6 / (7 - (4+3))")
	eval("-3 * 6 / (7 - (4+3))")
	// Output:
	// 10
	// [line 1] Error at 'Plus': cannot apply Plus: `+` to types boolean and nil (values true and nil) (token.BooleanValue and token.NilValue)
	// [line 1] Error at 'Slash': cannot apply Slash: `/` to types number and boolean (values 6 and false) (token.NumberValue and token.BooleanValue)
	// 3.5
	// +Inf
	// -Inf
}

func ExampleUnaryMinus() {
	eval("-8")
	eval(`-nil`)
	eval(`-false`)
	eval(`-true`)
	eval(`-"bar"`)
	// Output:
	// -8
	// [line 1] Error at 'Minus': cannot apply Minus: `-` to type nil (nil) (token.NilValue)
	// [line 1] Error at 'Minus': cannot apply Minus: `-` to type boolean (false) (token.BooleanValue)
	// [line 1] Error at 'Minus': cannot apply Minus: `-` to type boolean (true) (token.BooleanValue)
	// [line 1] Error at 'Minus': cannot apply Minus: `-` to type string ("bar") (token.StringValue)
}

func ExampleUnaryBang() {
	eval(`!nil`)
	eval("!false")
	eval("!true")
	eval(`!"hello, world!"`)
	eval(`!(430 +2)`)
	// Output:
	// true
	// true
	// false
	// [line 1] Error at 'Bang': cannot apply Bang: `!` to type string ("hello, world!") (token.StringValue)
	// [line 1] Error at 'Bang': cannot apply Bang: `!` to type number (432) (token.NumberValue)
}

func ExampleUnaryUnary() {
	eval("0 --9")
	eval("!!false")
	eval("!!true")
	// Output:
	// 9
	// false
	// true
}

func ExampleMinusOfMinusMinus() {
	eval("0 ---3")
	// Output:
	// -3
}

func ExampleExpr567() {
	eval("-5 * ( 6 * 7)")
	// Output:
	// -210
}

func ExampleExpr4567() {
	eval("4 - -5 * ( 6 * 7)")
	// Output:
	// 214
}

func ExampleComplex() {
	eval("1+2   * 3 / 4 - -5 * ( 6 * 7)")
	// Output:
	// 212.5
}

func ExampleStringConcat() {
	eval(`"a" + "\n" + "b"`)
	// Output:
	// "a\nb"
}

func ExampleComparison() {
	eval("7 >= 2+3")
	// Output:
	// true
}

func ExampleComparison2() {
	eval("7 >= 2+6")
	// Output:
	// false
}

func ExampleLess() {
	eval("1 < 2")
	eval("2 < 2")
	eval("3 < 2")
	eval(`"" < "a"`)
	eval(`"a" < "a"`)
	eval(`"ab" < "a"`)
	// Output:
	// true
	// false
	// false
	// true
	// false
	// false
}

func ExampleGreater() {
	eval("1 > 2")
	eval("2 > 2")
	eval("3 > 2")
	eval(`"" > "a"`)
	eval(`"a" > "a"`)
	eval(`"ab" > "a"`)
	// Output:
	// false
	// false
	// true
	// false
	// false
	// true
}

func ExampleLessEqual() {
	eval("1 <= 2")
	eval("2 <= 2")
	eval("3 <= 2")
	eval(`"" <= "a"`)
	eval(`"a" <= "a"`)
	eval(`"ab" <= "a"`)
	// Output:
	// true
	// true
	// false
	// true
	// true
	// false
}

func ExampleGreaterEqual() {
	eval("1 >= 2")
	eval("2 >= 2")
	eval("3 >= 2")
	eval(`"" >= "a"`)
	eval(`"a" >= "a"`)
	eval(`"ab" >= "a"`)
	// Output:
	// false
	// true
	// true
	// false
	// true
	// true
}

func ExampleEquality() {
	eval(`7 == 8`)
	eval(`4 == 4`)
	eval(`7 != 8`)
	eval(`4 != 4`)
	eval(`"foobie" == "foo"`)
	eval(`"foo" == "foo"`)
	eval(`"foobie" != "foo"`)
	eval(`"foo" != "foo"`)
	eval(`nil == nil`)
	eval(`nil != nil`)
	// Output:
	// false
	// true
	// true
	// false
	// false
	// true
	// true
	// false
	// true
	// false
}

func ExampleBooleanEquality() {
	eval(`false == false`)
	eval(`false == true`)
	eval(`true == false`)
	eval(`true == true`)
	eval(`false != false`)
	eval(`false != true`)
	eval(`true != false`)
	eval(`true != true`)
	// Output:
	// true
	// false
	// false
	// true
	// false
	// true
	// true
	// false
}

func ExampleTypeMismatchedEquality() {
	eval(`3 == "3"`)
	eval(`"0" == false`)
	eval(`false == nil`)
	eval(`nil == 0`)
	eval(`3 != "3"`)
	eval(`"0" != false`)
	eval(`false != nil`)
	eval(`nil != 0`)
	// Output:
	// false
	// false
	// false
	// false
	// true
	// true
	// true
	// true
}

func ExampleOr() {
	eval("nil or 1")
	eval(`false or 2`)
	eval(`true or 3`)
	eval(`"" or 4`)
	eval(`"0" or 5`)
	eval("0 or 6")
	eval("42 or 7")
	// Output:
	// 1
	// 2
	// true
	// ""
	// "0"
	// 0
	// 42
}

func ExampleAnd() {
	eval("nil and 1")
	eval(`false and 2`)
	eval(`true and 3`)
	eval(`"" and 4`)
	eval(`"0" and 5`)
	eval("0 and 6")
	eval("42 and 7")
	// Output:
	// nil
	// false
	// 3
	// 4
	// 5
	// 6
	// 7
}

func ExampleUnaryVsLowerPrecedence() {
	// unary - vs binary * - < != and or
	eval(`-7 * 5`)
	eval(`-7 - 4`)
	eval(`-7 < -8`)
	eval(`-7 != 7`)
	eval(`-7 and 7`)
	eval(`-7 or 7`)
	// Output:
	// -35
	// -11
	// false
	// true
	// 7
	// -7
}

func ExampleStarVsLowerPrecedence() {
	// binary '*' vs '-', '>=', '=='
	eval(`2  * 3  - 4`)
	eval(`2  - 3  * 4`)
	eval(`2  * 3 >= 5`)
	eval(`2 >= 3  * 5`)
	eval(`2  * 3 == 6`)
	eval(`2 == 3  * 6`)
	eval(`2  * 3 == 7`)
	eval(`2 == 3  * 7`)
	eval(`2  *  3 and 8`)
	eval(`2 and 3  *  8`)
	eval(`2  *  3 or  9`)
	eval(`2 or  3  *  9`)
	// Output:
	// 2
	// -10
	// true
	// false
	// true
	// false
	// false
	// false
	// 8
	// 24
	// 6
	// 2
}

func ExampleMinusVsLowerPrecedence() {
	// binary `-` vs binary `>`, `!=`
	eval(`9  -  8  >  7`)
	eval(`9  >  8  -  7`)
	eval(`9  -  8 !=  6`)
	eval(`9 !=  8  -  6`)
	eval(`9  -  8 and 5`)
	eval(`9 and 8  -  5`)
	eval(`9  -  8 or  4`)
	eval(`9 or  8  -  4`)
	// Output:
	// false
	// true
	// true
	// true
	// 5
	// 3
	// 1
	// 9
}

func ExampleCmpVsLowerPrecedence() {
	eval(`7 < 8 == (1==1)`)
	eval(`(1==0) == 7 < 8`)
	eval(`7 < 8 and 9`)
	eval(`7 and 8 < 9`)
	eval(`3 < 4 or 5`)
	eval(`5 < 4 or 6`)
	eval(`3 or 4 < 7`)
	eval(`3 < 2 or 8`)
	// Output:
	// true
	// false
	// 9
	// true
	// true
	// 6
	// 3
	// 8
}

func ExampleEqualityVsLowerPrecedence() {
	eval(`2 == 6 and nil `)
	eval(`3 == 7 and 1   `)
	eval(`4 != 8 and nil `)
	eval(`5 != 9 and 2   `)
	eval(`0 != 0  or nil `)
	eval(`0 != 0  or 1   `)
	eval(`1 == 1  or nil `)
	eval(`1 == 1  or 9   `)
	eval(`nil and 2 == 0`)
	eval(`nil and 3 == 3`)
	eval(`1   and 4 == 0`)
	eval(`2   and 5 == 5`)
	eval(`nil  or 0 != 0`)
	eval(`nil  or 0 == 0`)
	eval(`  3  or 1 != 1`)
	eval(`  4  or 1 == 1`)
	// Output:
	// false
	// false
	// nil
	// 2
	// nil
	// 1
	// true
	// true
	// nil
	// nil
	// false
	// true
	// false
	// true
	// 3
	// 4
}

func ExampleAndVsLowerPrecedence() {
	eval(`false and false  or   nil `) // F and F  or  F
	eval(`false and false  or   12  `) // F and F  or  T
	eval(`false and true   or   nil `) // F and T  or  F
	eval(`false and true   or   14  `) // F and T  or  T

	eval(`true  and false  or   nil `) // T and F  or  F
	eval(`true  and false  or   22  `) // T and F  or  T
	eval(`true  and  23    or  false`) // T and T  or  F
	eval(`true  and  24    or  true `) // T and T  or  T

	eval(`false or  nil   and  false`) // F  or  F and F
	eval(`false or  nil   and   true`) // F  or  F and T
	eval(`false or true   and   nil `) // F  or  T and F
	eval(`false or true   and   34  `) // F  or  T and T

	eval(` 41   or false  and  false`) // T  or  F and F
	eval(` 42   or false  and   true`) // T  or  T and F
	eval(` 43   or true   and  false`) // T  or  F and T
	eval(` 44   or true   and   true`) // T  or  T and T

	// Output:
	// nil
	// 12
	// nil
	// 14
	// nil
	// 22
	// 23
	// 24
	// nil
	// nil
	// nil
	// 34
	// 41
	// 42
	// 43
	// 44
}

func ExampleAssociativity() {
	eval(`7 + 4 - 100`)
	eval(`7 - 4 - 100`)
	eval(`7 * 4 / 100`)
	eval(`7 / 4 / 100`)
	// Output:
	// -89
	// -97
	// 0.28
	// 0.0175
}
