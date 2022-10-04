package scan

import (
	"fmt"

	"github.com/perlmonger42/go-lox/config"
	"github.com/perlmonger42/go-lox/lox"
)

func dumpTokens(text string) {
	config := config.New()
	lox := lox.New(config)
	scanner := New(lox, text)
	for _, token := range scanner.ScanTokens() {
		fmt.Printf("%s\n", token)
	}
}

func dumpTokensWithLineNumbers(text string) {
	config := config.New()
	lox := lox.New(config)
	scanner := New(lox, text)
	for _, token := range scanner.ScanTokens() {
		fmt.Printf("%s at line %d\n", token, token.Whence().Line())
	}
}

func ExampleParens() {
	dumpTokensWithLineNumbers("()")
	// Output:
	// LeftParen: `(` at line 1
	// RightParen: `)` at line 1
	// EOF at line 1
}

func ExampleUnexpectedInput() {
	dumpTokensWithLineNumbers("(~")
	// Output:
	// [line 1] Error at 'Other': Unexpected character ('~').
	// LeftParen: `(` at line 1
	// Other: `~` at line 1
	// EOF at line 1
}

func ExampleUnicode() {
	// "Ťėšťǐňġ" is entirely made of 2-byte UTF-8 encodings
	// "ṫẹṡṫịṅḡ" is entirely made of 3-byte UTF-8 encodings
	// "𝕠𝕟𝕖, 𝕥𝕨𝕠, 𝕥𝕙𝕣𝕖𝕖" has a bunch of 4-byte UTF-8 encodings
	dumpTokens("༺ Ťėšťǐňġ, ṫẹṡṫịṅḡ, 𝕠𝕟𝕖, 𝕥𝕨𝕠, 𝕥𝕙𝕣𝕖𝕖 ༻")

	// Output:
	// [line 1] Error at 'Other': Unexpected character ('༺').
	// [line 1] Error at 'Other': Unexpected character ('༻').
	// Other: `༺`
	// Identifier: `Ťėšťǐňġ`
	// Comma: `,`
	// Identifier: `ṫẹṡṫịṅḡ`
	// Comma: `,`
	// Identifier: `𝕠𝕟𝕖`
	// Comma: `,`
	// Identifier: `𝕥𝕨𝕠`
	// Comma: `,`
	// Identifier: `𝕥𝕙𝕣𝕖𝕖`
	// Other: `༻`
	// EOF
}

func ExampleOperators() {
	dumpTokens("!!====<<=>>=")
	// Output:
	// Bang: `!`
	// BangEqual: `!=`
	// EqualEqual: `==`
	// Equal: `=`
	// Less: `<`
	// LessEqual: `<=`
	// Greater: `>`
	// GreaterEqual: `>=`
	// EOF
}

func ExampleComment() {
	dumpTokensWithLineNumbers("=//stuff\tand nonsense\n=\n")
	// Output:
	// Equal: `=` at line 1
	// Equal: `=` at line 2
	// EOF at line 3
}

func Example_4_6() {
	dumpTokensWithLineNumbers(
		`// this is a comment
 		(( )){} // grouping stuff
 		!*+-/=<> <= == // operators`,
	)
	// Output:
	// LeftParen: `(` at line 2
	// LeftParen: `(` at line 2
	// RightParen: `)` at line 2
	// RightParen: `)` at line 2
	// LeftBrace: `{` at line 2
	// RightBrace: `}` at line 2
	// Bang: `!` at line 3
	// Star: `*` at line 3
	// Plus: `+` at line 3
	// Minus: `-` at line 3
	// Slash: `/` at line 3
	// Equal: `=` at line 3
	// Less: `<` at line 3
	// Greater: `>` at line 3
	// LessEqual: `<=` at line 3
	// EqualEqual: `==` at line 3
	// EOF at line 3
}

func ExampleStrings() {
	dumpTokensWithLineNumbers(
		`""
 		 "x"

 		 // missing close quote (unexpected newline)
 		 "unterminated
 		 "this is a string"
 		 "x\ny\"z"

 		 // missing close quote after escaping slash
 		 "x\

 		 // missing close quote (unexpected EOF)
 		 "`,
	)
	// Output:
	// [line 5] Error at 'InvalidString': Unterminated string literal
	// [line 10] Error at 'InvalidString': Unterminated string literal
	// [line 13] Error at 'InvalidString': Unterminated string literal
	// String: `""` = "" at line 1
	// String: `"x"` = "x" at line 2
	// InvalidString: `"unterminated` at line 5
	// String: `"this is a string"` = "this is a string" at line 6
	// String: `"x\ny\"z"` = "x\ny\"z" at line 7
	// InvalidString: `"x\` at line 10
	// InvalidString: `"` at line 13
	// EOF at line 13
}

func ExampleStringEscapeBeforeEof() {
	dumpTokensWithLineNumbers(`
       // missing close quote after escaping backslash (unexpected EOF)
 	  "abc\
 	`)
	// Output:
	// [line 3] Error at 'InvalidString': Unterminated string literal
	// InvalidString: `"abc\` at line 3
	// EOF at line 4
}

func ExampleNumbers() {
	dumpTokens(`
 		.     1
 		.2    3.  45
 		6.70 12.  .34
 		8.0625
 		700.00e+299
 		7e5000
 		1e  2E- 3e+
 	`)
	// Output:
	// [line 7] Error at 'InvalidNumber': Invalid number literal (7e5000): strconv.ParseFloat: parsing "7e5000": value out of range
	// [line 8] Error at 'InvalidNumber': Invalid number literal (1e): strconv.ParseFloat: parsing "1e": invalid syntax
	// [line 8] Error at 'InvalidNumber': Invalid number literal (2E-): strconv.ParseFloat: parsing "2E-": invalid syntax
	// [line 8] Error at 'InvalidNumber': Invalid number literal (3e+): strconv.ParseFloat: parsing "3e+": invalid syntax
	// Dot: `.`
	// Number: `1` = 1
	// Dot: `.`
	// Number: `2` = 2
	// Number: `3` = 3
	// Dot: `.`
	// Number: `45` = 45
	// Number: `6.70` = 6.7
	// Number: `12` = 12
	// Dot: `.`
	// Dot: `.`
	// Number: `34` = 34
	// Number: `8.0625` = 8.0625
	// Number: `700.00e+299` = 7e+301
	// InvalidNumber: `7e5000`
	// InvalidNumber: `1e`
	// InvalidNumber: `2E-`
	// InvalidNumber: `3e+`
	// EOF
}

func ExampleIdentiers() {
	dumpTokens(`
 		x     1y     π
 		_     _x     y_
 		shish-kebab
 		arabic٠١٢٣٤٥٦٧٨٩
 		Hello, 世界!
 	`)
	// Output:
	// Identifier: `x`
	// Number: `1` = 1
	// Identifier: `y`
	// Identifier: `π`
	// Identifier: `_`
	// Identifier: `_x`
	// Identifier: `y_`
	// Identifier: `shish`
	// Minus: `-`
	// Identifier: `kebab`
	// Identifier: `arabic٠١٢٣٤٥٦٧٨٩`
	// Identifier: `Hello`
	// Comma: `,`
	// Identifier: `世界`
	// Bang: `!`
	// EOF
}

func ExampleKeywords() {
	dumpTokens(`
         an
 		and class else false for fun if nil or
 		print return super this true var while
 		whiled
 	`)
	// Output:
	// Identifier: `an`
	// And: `and`
	// Class: `class`
	// Else: `else`
	// False: `false` = false
	// For: `for`
	// Fun: `fun`
	// If: `if`
	// Nil: `nil` = nil
	// Or: `or`
	// Print: `print`
	// Return: `return`
	// Super: `super`
	// This: `this`
	// True: `true` = true
	// Var: `var`
	// While: `while`
	// Identifier: `whiled`
	// EOF
}
