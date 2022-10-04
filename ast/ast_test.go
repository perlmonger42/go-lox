package ast

import (
	"fmt"
	"github.com/perlmonger42/go-lox/token"
)

func ExampleParens() {
	expression := &Binary{
		token.New(token.Star, "*", nil, token.NewPos(1)),
		&Unary{
			token.New(token.Minus, "-", nil, token.NewPos(1)),
			&Literal{&token.NumberValue{123.0}},
		},
		&Grouping{&Literal{&token.NumberValue{45.67}}},
	}

	fmt.Println(ToString(expression))
	// Output:
	// (* (- 123) (group 45.67))
}

func ExampleIf() {
	stmt := &If{
		Condition:  &Literal{&token.BooleanValue{true}},
		ThenBranch: &Expression{&Literal{&token.NumberValue{1.0}}},
		ElseBranch: nil,
	}

	fmt.Println(ToString(stmt))
	// Output:
	// if (true)
	//   1;
}

func ExampleIfElse() {
	stmt := &If{
		Condition:  &Literal{&token.BooleanValue{true}},
		ThenBranch: &Expression{&Literal{&token.NumberValue{1.0}}},
		ElseBranch: &Expression{&Literal{&token.NumberValue{2.0}}},
	}

	fmt.Println(ToString(stmt))
	// Output:
	// if (true)
	//   1;
	// else
	//   2;
}

func ExampleIfElseIf() {
	stmt := &If{
		Condition:  &Literal{&token.BooleanValue{true}},
		ThenBranch: &Expression{&Literal{&token.NumberValue{1.0}}},
		ElseBranch: &If{
			Condition:  &Literal{&token.BooleanValue{false}},
			ThenBranch: &Expression{&Literal{&token.NumberValue{2.0}}},
			ElseBranch: nil,
		},
	}

	fmt.Println(ToString(stmt))
	// Output:
	// if (true)
	//   1;
	// else if (false)
	//   2;
}

func ExampleBlockishIfElseIfElse() {
	stmt := &If{
		Condition: &Literal{&token.BooleanValue{true}},
		ThenBranch: &Block{
			&token.Token{Type_: token.LeftBrace},
			[]Stmt{&Expression{&Literal{&token.NumberValue{1.0}}}},
		},
		ElseBranch: &If{
			Condition: &Literal{&token.BooleanValue{false}},
			ThenBranch: &Block{
				&token.Token{Type_: token.LeftBrace},
				[]Stmt{&Expression{&Literal{&token.NumberValue{2.0}}}},
			},
			ElseBranch: &Block{
				&token.Token{Type_: token.LeftBrace},
				[]Stmt{&Expression{&Literal{&token.NumberValue{3.0}}}},
			},
		},
	}

	fmt.Println(ToString(stmt))
	// Output:
	// if (true) {
	//   1;
	// } else if (false) {
	//   2;
	// } else {
	//   3;
	// }
}

func ExampleWhile() {
	stmt := &While{
		Condition: &Literal{&token.BooleanValue{false}},
		Body:      &Expression{&Literal{&token.NumberValue{1.0}}},
	}

	fmt.Println(ToString(stmt))
	// Output:
	// while (false)
	//   1;
}

func ExampleBlockishWhile() {
	stmt := &While{
		Condition: &Literal{&token.BooleanValue{false}},
		Body: &Block{
			&token.Token{Type_: token.LeftBrace},
			[]Stmt{&Expression{&Literal{&token.NumberValue{1.0}}}},
		},
	}

	fmt.Println(ToString(stmt))
	// Output:
	// while (false) {
	//   1;
	// }
}
