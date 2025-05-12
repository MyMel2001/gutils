package main

import (
	"fmt"
	"os"
)

// rmdir removes empty directories specified in the command-line arguments.
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: rmdir DIR...")
		os.Exit(1)
	}
	status := 0
	for _, dir := range os.Args[1:] {
		err := os.Remove(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rmdir: failed to remove '%s': %v\n", dir, err)
			status = 1
		}
	}
	os.Exit(status)
}
