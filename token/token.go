package token

import "fmt"

//go:generate stringer -type Type

// T represents a token or text string returned from the scanner.
type T interface {
	Type() Type     // The type of this item.
	Lexeme() string // The text of this item.
	Literal() Value // The value of the literal, if literal
	Whence() Pos    // The position at which this token appears
}

type Token struct {
	Type_    Type   // The type of this item.
	Lexeme_  string // The text of this item.
	Literal_ Value  // The value of the literal, if literal
	Whence_  Pos    // The position at which this token appears
}

func New(t Type, text string, literal Value, pos Pos) T {
	return &Token{t, text, literal, pos}
}

func (t *Token) Type() Type     { return t.Type_ }
func (t *Token) Lexeme() string { return t.Lexeme_ }
func (t *Token) Literal() Value { return t.Literal_ }
func (t *Token) Whence() Pos    { return t.Whence_ }

// Type identifies the type of lex items.
type Type int

const (
	EOF Type = iota // make EOF be the zero value

	// Punctuation and Operators
	LeftParen    // "("
	RightParen   // ")"
	LeftBrack    // "["
	RightBrack   // "]"
	LeftBrace    // "{"
	RightBrace   // "}"
	Comma        // ","
	Dot          // "."
	Minus        // "}"
	Plus         // "+"
	Star         // "*"
	Slash        // "/"
	Semicolon    // ";"
	Bang         // "!"
	BangEqual    // "!="
	Equal        // "="
	EqualEqual   // "=="
	Less         // "<"
	LessEqual    // "<="
	Greater      // ">"
	GreaterEqual // ">="

	// Keywords
	And    // "and"
	Class  // "class"
	Else   // "else"
	False  // "false"
	For    // "for"
	Fun    // "fun"
	If     // "if"
	Nil    // "nil"
	Or     // "or"
	Print  // "print"
	Return // "return"
	Super  // "super"
	This   // "this"
	True   // "true"
	Var    // "var"
	While  // "while"

	// Tokens with literal value
	String        // quoted string (includes quotes)
	InvalidString // unterminated, or has invalid escape sequence
	Number        // an integer or floating point number
	InvalidNumber // kinda close to being a number, but not quite right
	Identifier    // alphanumeric identifier

	Other // unrecognized character
)

func (i Token) String() string {
	if i.Type() == EOF {
		return "EOF"
	} else if i.Literal() != nil {
		return fmt.Sprintf("%s: %#q = %s",
			i.Type(), i.Lexeme(), i.Literal().Show())
	} else {
		return fmt.Sprintf("%s: %#q",
			i.Type(), i.Lexeme())
	}
}
