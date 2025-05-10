package main

import (
	"fmt"
	"os"
)

// mv: renames (moves) a file from source to destination
func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "mv: usage: mv SOURCE DEST")
		os.Exit(1)
	}
	src, dst := os.Args[1], os.Args[2]
	err := os.Rename(src, dst)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mv: error moving file:", err)
		os.Exit(1)
	}
} 