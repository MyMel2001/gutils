package main

import (
	"fmt"
	"os"
)

// pwd: prints the current working directory
func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pwd:", err)
		os.Exit(1)
	}
	fmt.Println(dir)
} 