package token

import "fmt"

type Pos interface {
	Line() int // just the line number, for now
	String() string
}

func NewPos(line int) Pos {
	return &Position{line}
}

type Position struct {
	line int
}

var _ Pos = &Position{}

func (p *Position) String() string {
	return fmt.Sprintf("line %d", p.Line())
}

func (p *Position) Line() int {
	return p.line
}
