package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// dirname: prints the directory part of a path
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "dirname: usage: dirname PATH")
		os.Exit(1)
	}
	fmt.Println(filepath.Dir(os.Args[1]))
} 