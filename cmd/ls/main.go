package main

import (
	"fmt"
	"os"
)

// ls: lists files in the current directory
func main() {
	dir, err := os.Open(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ls: cannot open directory:", err)
		os.Exit(1)
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ls: cannot read directory:", err)
		os.Exit(1)
	}

	for _, name := range files {
		fmt.Println(name)
	}
} 