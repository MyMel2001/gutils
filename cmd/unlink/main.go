package main

import (
	"fmt"
	"os"
)

// unlink: removes a single file (like rm)
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "unlink: usage: unlink FILE")
		os.Exit(1)
	}
	if err := os.Remove(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, "unlink: cannot remove:", os.Args[1], err)
		os.Exit(1)
	}
}
