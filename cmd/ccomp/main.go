package main

import (
	"fmt"
	"os"
	"strings"
)

// ccomp: checks for .c or .cpp files and prints a message about compiling (no real compilation)
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "ccomp: usage: ccomp SOURCE [-o OUTPUT] [EXTRA_ARGS...]")
		os.Exit(1)
	}
	source := os.Args[1]
	if !strings.HasSuffix(source, ".c") && !strings.HasSuffix(source, ".cpp") {
		fmt.Fprintln(os.Stderr, "ccomp: only .c or .cpp files supported")
		os.Exit(1)
	}
	fmt.Printf("Would compile %s (not implemented)\n", source)
}
