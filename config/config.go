package config

import (
	"github.com/perlmonger42/go-lox/report"
)

type T struct {
	Prompt           string
	Reporter         report.T
	TraceScanTokens  bool // print tokens after scanner creates them
	TraceParseTokens bool // print tokens as parser consumes them
	TraceNodes       bool // print AST nodes as they are built
	TraceParsed      bool // dump AST rendered as Lox
	TraceEval        bool // print intermediate values as executed
}

func New() *T {
	return &T{
		Prompt:   "> ",
		Reporter: report.NewStdoutReporter(),
	}
}
