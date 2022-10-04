package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bobappleyard/readline"

	"github.com/perlmonger42/go-lox/ast"
	"github.com/perlmonger42/go-lox/config"
	"github.com/perlmonger42/go-lox/interpret"
	"github.com/perlmonger42/go-lox/lox"
	"github.com/perlmonger42/go-lox/parse"
	"github.com/perlmonger42/go-lox/resolve"
	"github.com/perlmonger42/go-lox/scan"
	"github.com/perlmonger42/go-lox/token"
)

var (
	execute = flag.Bool("e", false, "execute arguments as a program")
	testing = flag.Bool("test", false, "execute Read Eval Read Compare Loop")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: go-lox [options] [file]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	os.Exit(64) // see "sysexits.h"
}

func main() {
	flag.Usage = usage
	flag.Parse()
	config := config.New()
	lox := lox.New(config)

	if *execute {
		text := strings.Join(flag.Args(), " ")
		runText(lox, text)
	} else if flag.NArg() == 0 {
		lox.Interactive = true
		runPrompt(lox)
	} else if flag.NArg() == 1 {
		runFile(lox, config, flag.Arg(0))
	} else {
		usage()
	}
	os.Exit(0)
}

func ConsoleReadline(config *config.T) (string, error) {
	prompt := config.Prompt
	if prompt == "" {
		prompt = "go-lox> "
	}
	if line, err := readline.String(prompt); err != nil {
		return "", err
	} else {
		readline.AddHistory(line)
		return line + "\n", nil
	}
}

func runFile(lox *lox.T, config *config.T, filename string) {
	if content, err := ioutil.ReadFile(filename); err != nil {
		fmt.Fprintf(os.Stderr, "error in go-lox: %s\n", err)
		os.Exit(66) // see "sysexits.h"
	} else {
		text := string(content) // convert []byte to string
		runText(lox, text)
	}
}

func runPrompt(lox *lox.T) {
	for {
		if line, err := ConsoleReadline(lox.Config); err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "go-lox: readline error: %s\n", err)
			os.Exit(66) // see "sysexits.h"
		} else {
			runText(lox, line)
		}
	}
}

func runText(lox *lox.T, text string) {
	lox.Config.TraceScanTokens = false
	lox.Config.TraceNodes = false
	lox.Config.TraceEval = false

	var scanner scan.T = scan.New(lox, text)
	var tokens []token.T = scanner.ScanTokens()
	if lox.HadError {
		return
	}
	if lox.Config.TraceScanTokens {
		for _, token := range tokens {
			fmt.Printf("NEXT is %s\n", token)
		}
	}

	var parser parse.T = parse.New(lox, tokens)
	var stmts []ast.Stmt = parser.Parse()
	if lox.Config.TraceParsed {
		for i, stmt := range stmts {
			fmt.Printf("%2d: %s\n", i, ast.ToString(stmt))
		}
	}
	if lox.HadError {
		return
	}

	var interpreter interpret.T = interpret.New(lox)

	var resolver *resolve.T = resolve.New(lox, interpreter)
	resolver.ResolveStmtList(stmts)
	if lox.HadError {
		return
	}

	fmt.Printf("running interpreter\n")
	//err := interpreter.InterpretStmts2(stmts)
	//if err != nil {
	//	fmt.Printf("error: (%s) %s\n", err, err)
	//}
	interpreter.InterpretStmts(stmts)
}
