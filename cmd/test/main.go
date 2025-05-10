package main

import (
	"os"
)

// test: checks if a file exists (returns 0 if exists, 1 if not)
func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	_, err := os.Stat(os.Args[1])
	if err == nil {
		os.Exit(0)
	}
	os.Exit(1)
}
