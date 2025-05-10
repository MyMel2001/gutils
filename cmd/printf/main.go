package main

import (
	"fmt"
	"os"
)

// printf: prints formatted output
func main() {
	if len(os.Args) < 2 {
		os.Exit(0)
	}
	format := os.Args[1]
	args := make([]interface{}, len(os.Args)-2)
	for i, v := range os.Args[2:] {
		args[i] = v
	}
	fmt.Printf(format, args...)
} 