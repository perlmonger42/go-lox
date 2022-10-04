package lox

import (
	"fmt"

	"github.com/perlmonger42/go-lox/config"
	"github.com/perlmonger42/go-lox/token"
)

type T struct {
	Config      *config.T
	Interactive bool
	HadError    bool
}

func New(config *config.T) *T {
	return &T{Config: config}
}

func (lox *T) Error(tok token.T, message string) {
	lox.HadError = true
	pos := tok.Whence()
	typ := tok.Type()
	if typ == token.EOF {
		lox.Config.Reporter.Report(pos, "at end", message)
	} else {
		lox.Config.Reporter.Report(pos, fmt.Sprintf("at '%s'", typ), message)
	}
}

func (l *T) Report(pos token.Pos, where string, message string) {
	l.Config.Reporter.Report(pos, where, message)
}
