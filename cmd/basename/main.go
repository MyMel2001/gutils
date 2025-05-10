package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// basename: prints the final component of a path
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "basename: usage: basename PATH...")
		os.Exit(1)
	}
	for _, arg := range os.Args[1:] {
		fmt.Println(filepath.Base(arg))
	}
}
