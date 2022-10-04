package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var outputDir = flag.String("d", ".", "output directory")

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 0 {
		usage()
	}
	defineAst(outputDir, exprDescription)
	defineAst(outputDir, stmtDescription)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: generate-ast [options]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	os.Exit(64) // see "sysexits.h"
}

type FieldDescription struct {
	FieldName string
	FieldType string
}

type Subclass struct {
	SubName string
	Fields  []FieldDescription
}

type Classes struct {
	BaseName     string
	VisitorTypes []string
	Imports      []string
	Subclasses   []Subclass
}

var exprDescription = &Classes{
	BaseName:     "Expr",
	VisitorTypes: []string{"token.Value", "string", "", "(token.Value, error)"},
	Imports:      []string{"github.com/perlmonger42/go-lox/token"},
	Subclasses: []Subclass{
		{"Grouping", []FieldDescription{
			{"Expression", "Expr"},
		}},
		{"This", []FieldDescription{
			{"Keyword", "token.T"},
		}},
		{"Super", []FieldDescription{
			{"Keyword", "token.T"},
			{"Method", "token.T"},
		}},
		{"Variable", []FieldDescription{
			{"Name", "token.T"},
		}},
		{"Literal", []FieldDescription{
			{"Value", "token.Value"},
		}},
		{"Call", []FieldDescription{
			{"Callee", "Expr"},
			{"Paren", "token.T"},
			{"Arguments", "[]Expr"},
		}},
		{"Get", []FieldDescription{
			{"Object", "Expr"},
			{"Name", "token.T"},
		}},
		{"Unary", []FieldDescription{
			{"Operator", "token.T"},
			{"Right", "Expr"},
		}},
		{"Binary", []FieldDescription{
			{"Operator", "token.T"},
			{"Left", "Expr"},
			{"Right", "Expr"},
		}},
		{"Logical", []FieldDescription{
			{"Operator", "token.T"},
			{"Left", "Expr"},
			{"Right", "Expr"},
		}},
		{"Set", []FieldDescription{
			{"Object", "Expr"},
			{"Name", "token.T"},
			{"Value", "Expr"},
		}},
		{"Assign", []FieldDescription{
			{"Name", "token.T"},
			{"Value", "Expr"},
		}},
	},
}

var stmtDescription = &Classes{
	BaseName:     "Stmt",
	VisitorTypes: []string{"", "string", "error"},
	Imports:      []string{"github.com/perlmonger42/go-lox/token"},
	Subclasses: []Subclass{
		{"Noop", []FieldDescription{}},
		{"Expression", []FieldDescription{
			{"Expression", "Expr"},
		}},
		{"Print", []FieldDescription{
			{"Keyword", "token.T"},
			{"Expression", "Expr"},
		}},
		{"Return", []FieldDescription{
			{"Keyword", "token.T"},
			{"Value", "Expr"},
		}},
		{"Panic", []FieldDescription{
			{"Keyword", "token.T"},
			{"Expression", "Expr"},
		}},
		{"VarInitialized", []FieldDescription{
			{"Name", "token.T"},
			{"Initializer", "Expr"},
		}},
		{"VarUninitialized", []FieldDescription{
			{"Name", "token.T"},
		}},
		{"Function", []FieldDescription{
			{"Name", "token.T"},
			{"Params", "[]token.T"},
			{"Body", "[]Stmt"},
		}},
		{"If", []FieldDescription{
			{"Condition", "Expr"},
			{"ThenBranch", "Stmt"},
			{"ElseBranch", "Stmt"},
		}},
		{"Block", []FieldDescription{
			{"Token", "token.T"},
			{"Statements", "[]Stmt"},
		}},
		{"While", []FieldDescription{
			{"Condition", "Expr"},
			{"Body", "Stmt"},
		}},
		{"Class", []FieldDescription{
			{"Name", "token.T"},
			{"Superclass", "*Variable"},
			{"Methods", "[]*Function"},
		}},
	},
}

var astSourceTemplate string = `
{{- $base := .BaseName}}
{{- $subclasses := .Subclasses}}
{{- $visitorReturnTypes := .VisitorTypes -}}
package ast

import ({{range .Imports -}}
	"{{.}}"
{{end}}
)

type {{$base}} interface {
	AsNode() Node // does nothing but prevent non-Nodes from looking like Nodes
	As{{$base}}() {{$base}} // does nothing but prevent non-{{$base}}s from looking like {{$base}}

{{range $visitorReturnTypes}}
{{- $type := .}}{{$Type := ""}}
{{- if $type}}{{$Type = title $type | printf "_%s" }}{{end}}
	Accept_{{$base}}{{$Type}}(visitor Visitor_{{$base}}{{$Type}}) {{$type}}
{{- end}}
}

{{range $visitorReturnTypes}}
{{- $type := .}}{{$Type := ""}}
{{- if $type}}{{$Type = title $type | printf "_%s" }}{{end}}
  // A Visitor_{{$base}}{{$Type}} is accepted by {{$base}} and {{if $type}}returns {{$type}}{{else}}has no return value{{end}}
  type Visitor_{{$base}}{{$Type}} interface {
    {{- range $subclasses}}
      Visit_{{.SubName}}{{$base}}{{$Type}}({{lc $base}} *{{.SubName}}) {{$type}}
    {{- end}}
  }
{{end}}

{{range $subclasses -}}
{{$subname := .SubName}}
  type {{$subname}} struct {
  {{- range .Fields}}
    {{.FieldName}} {{.FieldType}}
  {{- end}}
  }
  func (x *{{$subname}}) AsNode() Node { return x }
  func (x *{{$subname}}) As{{$base}}() {{$base}} { return x }
  {{range $visitorReturnTypes -}}
  {{- $type := .}}{{$Type := ""}}
  {{- if $type}}{{$Type = title $type | printf "_%s" }}{{end}}
    func (x *{{$subname}}) Accept_{{$base}}{{$Type}}(visitor Visitor_{{$base}}{{$Type}}) {{$type}} {
	  {{if $type}}return {{end}}visitor.Visit_{{$subname}}{{$base}}{{$Type}}(x)
    }
  {{- end}}
{{end}}
`

func defineAst(outputDir *string, classes *Classes) {
	lcName := strings.ToLower(classes.BaseName)
	path := fmt.Sprintf("%s/%s.go", *outputDir, lcName)
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open %q for writing\n", path)
		os.Exit(66) // see "sysexits.h"
	}
	defer f.Close()

	funcMap := template.FuncMap{
		// The name "lc" is what the function will be called in template text.
		"lc": strings.ToLower,
		"title": func(s string) string {
			if s == "(token.Value, error)" {
				return "MaybeValue"
			}
			s = strings.ReplaceAll(s, ".", "_")
			return strings.Title(s)
		},
	}
	t := template.Must(template.
		New("source").
		Funcs(funcMap).
		Parse(astSourceTemplate))

	err = t.Execute(f, classes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error evaluating template: %s\n", err)
		os.Exit(66) // see "sysexits.h"
	}
}
