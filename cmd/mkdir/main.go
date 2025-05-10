package main

import (
	"fmt"
	"os"
)

// mkdir: creates directories given as arguments
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "mkdir: usage: mkdir DIR...")
		os.Exit(1)
	}
	for _, dir := range os.Args[1:] {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Fprintln(os.Stderr, "mkdir: cannot create directory:", dir, err)
		}
	}
} 