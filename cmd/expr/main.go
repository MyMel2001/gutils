package main

import (
	"fmt"
	"os"
	"strconv"
)

// expr evaluates basic integer arithmetic expressions (e.g., expr 1 + 2).
func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: expr INT OP INT")
		os.Exit(1)
	}
	a, err1 := strconv.Atoi(os.Args[1])
	b, err2 := strconv.Atoi(os.Args[3])
	if err1 != nil || err2 != nil {
		fmt.Fprintln(os.Stderr, "expr: arguments must be integers")
		os.Exit(1)
	}
	var result int
	switch os.Args[2] {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			fmt.Fprintln(os.Stderr, "expr: division by zero")
			os.Exit(1)
		}
		result = a / b
	case "%":
		if b == 0 {
			fmt.Fprintln(os.Stderr, "expr: division by zero")
			os.Exit(1)
		}
		result = a % b
	default:
		fmt.Fprintf(os.Stderr, "expr: unknown operator '%s'\n", os.Args[2])
		os.Exit(1)
	}
	fmt.Println(result)
}
