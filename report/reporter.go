// Package report defines the error-reporting interface used by go-lox.
package report

import (
	"fmt"
	"os"

	"github.com/perlmonger42/go-lox/token"
)

type T interface {
	// Report generates an error report for a given location.
	Report(pos token.Pos, where string, message string)
}

// StderrReporter is a Reporter that simply prints messages to os.Stderr
type StderrReporter struct {
}

func (c *StderrReporter) Report(pos token.Pos, where string, message string) {
	pad := ""
	if where != "" {
		pad = " "
	}
	fmt.Fprintf(os.Stderr, "[%s] Error%s%s: %s\n", pos, pad, where, message)
}

func NewStderrReporter() T {
	return &StderrReporter{}
}

// StdoutReporter is a Reporter that simply prints messages to os.Stdout
type StdoutReporter struct {
}

func (c *StdoutReporter) Report(pos token.Pos, where string, message string) {
	pad := ""
	if where != "" {
		pad = " "
	}
	fmt.Fprintf(os.Stdout, "[%s] Error%s%s: %s\n", pos, pad, where, message)
}

func NewStdoutReporter() T {
	return &StdoutReporter{}
}
