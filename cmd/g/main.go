package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// gcomp: Go code evaluator using Yaegi, with REPL mode
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "gcomp: usage: gcomp SOURCE.go | - (eval), gcomp repl (REPL)")
		os.Exit(1)
	}

	if os.Args[1] == "repl" {
		repl()
		return
	}

	var code []byte
	var err error
	if os.Args[1] == "-" {
		code, err = io.ReadAll(os.Stdin)
	} else {
		code, err = os.ReadFile(os.Args[1])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "gcomp: error reading code:", err)
		os.Exit(1)
	}
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	_, err = i.Eval(string(code))
	if err != nil {
		fmt.Fprintln(os.Stderr, "gcomp: eval error:", err)
		os.Exit(1)
	}
}

// repl provides an interactive Go REPL using Yaegi
func repl() {
	fmt.Println("gcomp REPL (type Go code, end with Ctrl+D)")
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "quit" {
			break
		}
		_, err := i.Eval(line)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}
