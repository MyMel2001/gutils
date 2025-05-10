package main

import (
	"fmt"
	"os"
)

// rm: removes files given as arguments
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "rm: usage: rm FILE...")
		os.Exit(1)
	}
	for _, fname := range os.Args[1:] {
		err := os.Remove(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "rm: cannot remove:", fname, err)
		}
	}
} 