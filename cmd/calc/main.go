package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
)

// calc: minimal calculator, supports +, -, *, /, and parentheses
func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		for _, expr := range args {
			calc(expr)
		}
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		calc(scanner.Text())
	}
}

func calc(expr string) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return
	}
	n, err := parser.ParseExpr(expr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "calc: parse error:", err)
		return
	}
	res, err := eval(n)
	if err != nil {
		fmt.Fprintln(os.Stderr, "calc: eval error:", err)
		return
	}
	fmt.Println(res)
}

func eval(n ast.Expr) (float64, error) {
	switch x := n.(type) {
	case *ast.BasicLit:
		return strconv.ParseFloat(x.Value, 64)
	case *ast.BinaryExpr:
		l, err := eval(x.X)
		if err != nil {
			return 0, err
		}
		r, err := eval(x.Y)
		if err != nil {
			return 0, err
		}
		switch x.Op {
		case token.ADD:
			return l + r, nil
		case token.SUB:
			return l - r, nil
		case token.MUL:
			return l * r, nil
		case token.QUO:
			return l / r, nil
		}
	}
	return 0, fmt.Errorf("unsupported expression")
}
