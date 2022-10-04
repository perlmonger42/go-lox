// This lexer is based on Rob Pike's Ivy scanner, found at
// https://github.com/robpike/ivy/blob/master/scan/scan.go

package scan

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/perlmonger42/go-lox/lox"
	"github.com/perlmonger42/go-lox/token"
)

// A Scanner.T converts a string into an array of Tokens,
// via its ScanTokens method.
type T interface {
	ScanTokens() []token.T
}

type Value = token.Value

type Scanner struct {
	lox     *lox.T
	source  string
	tokens  []token.T
	start   int // index in source of first char of current token
	current int // index in source of char under read head
	line    int // line number of current character
}

func New(lox *lox.T, source string) T {
	return &Scanner{
		lox:    lox,
		source: source,
		tokens: []token.T{},
		line:   1,
	}
}

func (s *Scanner) Error(tok token.T, message string) {
	s.lox.Error(tok, message)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) newToken(typ token.Type, val Value) token.T {
	text := s.source[s.start:s.current]
	return token.New(typ, text, val, token.NewPos(s.line))
}

func (s *Scanner) addFullToken(tok token.T) token.T {
	s.tokens = append(s.tokens, tok)
	return tok
}

func (s *Scanner) addToken(typ token.Type) token.T {
	tok := s.newToken(typ, nil)
	if s.lox.Config.TraceScanTokens {
		fmt.Printf("token: %s\n", tok)
	}
	s.tokens = append(s.tokens, tok)
	return tok
}

func (s *Scanner) addTokenWithValue(typ token.Type, literal token.Value) token.T {
	tok := s.newToken(typ, literal)
	if s.lox.Config.TraceScanTokens {
		fmt.Printf("token: %s\n", tok)
	}
	s.tokens = append(s.tokens, tok)
	return tok
}

func (s *Scanner) ScanTokens() []token.T {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	pos := token.NewPos(s.line)
	s.tokens = append(s.tokens, token.New(token.EOF, "", nil, pos))
	return s.tokens
}

func (s *Scanner) scanToken() token.T {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LeftParen)
	case ')':
		s.addToken(token.RightParen)
	case '[':
		s.addToken(token.LeftBrack)
	case ']':
		s.addToken(token.RightBrack)
	case '{':
		s.addToken(token.LeftBrace)
	case '}':
		s.addToken(token.RightBrace)
	case ',':
		s.addToken(token.Comma)
	case '.':
		s.addToken(token.Dot)
	case '-':
		s.addToken(token.Minus)
	case '+':
		s.addToken(token.Plus)
	case ';':
		s.addToken(token.Semicolon)
	case '*':
		s.addToken(token.Star)
	case '!':
		s.addToken(s.match2nd('=', token.Bang, token.BangEqual))
	case '=':
		s.addToken(s.match2nd('=', token.Equal, token.EqualEqual))
	case '<':
		s.addToken(s.match2nd('=', token.Less, token.LessEqual))
	case '>':
		s.addToken(s.match2nd('=', token.Greater, token.GreaterEqual))
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for !s.isAtEnd() && s.peek() != '\n' {
				s.advance()
			}
		} else {
			s.addToken(token.Slash)
		}

	case '"':
		s.scanString()

	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		s.line++

	default:
		if isDigit(c) {
			s.scanNumber()
		} else if isAlpha(c) {
			s.scanIdentifier()
		} else {
			s.Error(
				s.addToken(token.Other),
				fmt.Sprintf("Unexpected character ('%c').", c),
			)
		}
	}
	return token.New(token.EOF, "", nil, token.NewPos(s.line))
}

var keywords map[string]token.Type = map[string]token.Type{
	"and":    token.And,
	"class":  token.Class,
	"else":   token.Else,
	"false":  token.False,
	"for":    token.For,
	"fun":    token.Fun,
	"if":     token.If,
	"nil":    token.Nil,
	"or":     token.Or,
	"print":  token.Print,
	"return": token.Return,
	"super":  token.Super,
	"this":   token.This,
	"true":   token.True,
	"var":    token.Var,
	"while":  token.While,
}

func (s *Scanner) scanIdentifier() {
	for isAlnum(s.peek()) {
		s.advance()
	}

	// See if the identifier is a reserved word.
	text := s.source[s.start:s.current]

	if t, ok := keywords[text]; ok {
		switch t {
		case token.True:
			s.addTokenWithValue(t, token.BooleanValue{true})
		case token.False:
			s.addTokenWithValue(t, token.BooleanValue{false})
		case token.Nil:
			s.addTokenWithValue(t, token.NilValue{})
		default:
			s.addToken(t)
		}
	} else {
		s.addToken(token.Identifier)
	}
}

// scanNumber parses a number. It isn't perfect - for instance, it accepts
// "1.2e" and ".3e-" - but when it's wrong, the input is invalid, and the
// scanner will notice because ParseFloat will notice.
func (s *Scanner) scanNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the ".".
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	// Look for exponent part.
	if s.match('e') || s.match('E') {
		if s.match('+') || s.match('-') {
		}
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	numStr := s.source[s.start:s.current]
	if f, err := strconv.ParseFloat(numStr, 64); err != nil {
		msg := fmt.Sprintf("Invalid number literal (%s): %s", numStr, err)
		s.Error(s.addToken(token.InvalidNumber), msg)
	} else {
		s.addTokenWithValue(token.Number, token.NumberValue{V: f})
	}
}
func (s *Scanner) scanString() {
	var r rune
	for {
		switch r = s.advance(); r {
		case '"':
			tok := s.newToken(token.String, nil)
			if value, ok := s.stringValue(tok); ok {
				s.addTokenWithValue(token.String, value)
			} else {
				s.addToken(token.InvalidString)
			}
			return
		case '\\':
			if r = s.advance(); r != eof && r != '\n' {
				continue
			}
			fallthrough
		case eof, '\n':
			if r == '\n' {
				s.current--
			}
			tok := s.addToken(token.InvalidString)
			s.Error(tok, "Unterminated string literal")
			s.stringValue(tok) // check for invalid escape sequences
			return
		}
	}
}

const eof = -1

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isAlnum(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func (s *Scanner) decodeNextRune() (rune, int) {
	r, w := utf8.DecodeRuneInString(s.source[s.current:])
	if w == 0 {
		r = eof
	}
	return r, w
}

func (s *Scanner) peek() rune {
	r, _ := s.decodeNextRune()
	return r
}

func (s *Scanner) peekNext() (r rune) {
	var w int
	r, w = utf8.DecodeRuneInString(s.source[s.current:])
	r, w = utf8.DecodeRuneInString(s.source[s.current+w:])
	if w == 0 {
		r = eof
	}
	return r
}

func (s *Scanner) advance() rune {
	r, w := s.decodeNextRune()
	s.current += w
	return r
}

func (s *Scanner) match(expected rune) bool {
	if s.peek() != expected {
		return false
	}
	s.advance()
	return true
}

func (s *Scanner) match2nd(r2 rune, t1 token.Type, t2 token.Type) token.Type {
	if s.match(r2) {
		return t2
	}
	return t1
}

func (s *Scanner) stringValue(tok token.T) (value token.StringValue, ok bool) {
	text := tok.Lexeme()
	if len(text) > 0 && text[0] == '"' {
		text = text[1:]
	}
	if len(text) > 0 && text[len(text)-1] == '"' {
		text = text[:len(text)-1]
	}

	str := ""
	ok = true
	escaped := false
	var r rune
	for _, r = range text {
		if escaped {
			r2 := r
			switch r {
			case '"':
				r2 = '"'
			case '\\':
				r2 = '\\'
			case 'a':
				r2 = '\a'
			case 'b':
				r2 = '\b'
			case 'f':
				r2 = '\f'
			case 'n':
				r2 = '\n'
			case 'r':
				r2 = '\r'
			case 't':
				r2 = '\t'
			case 'v':
				r2 = '\v'
			default:
				str = str + string('\\')
				r2 = r
				ok = false
				s.Error(tok,
					fmt.Sprintf("Invalid escape sequence in string (\\%c)", r))
			}
			str = str + string(r2)
			escaped = false
		} else if r == '\\' {
			escaped = true
		} else {
			str = str + string(r)
		}
	}
	if escaped {
		str = str + string('\\')
	}
	return token.StringValue{str}, ok
}
