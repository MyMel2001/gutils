package main

import (
	"fmt"
	"os"
)

// stty is not implemented. This is a placeholder for POSIX compliance.
func main() {
	fmt.Fprintln(os.Stderr, "stty: not implemented")
	os.Exit(1)
}
