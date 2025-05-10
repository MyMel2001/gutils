package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// realpath: prints the absolute path of a file or directory
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "realpath: usage: realpath FILE")
		os.Exit(1)
	}
	abs, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "realpath:", err)
		os.Exit(1)
	}
	fmt.Println(abs)
} 