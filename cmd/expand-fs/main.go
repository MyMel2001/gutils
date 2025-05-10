package main

import (
	"fmt"
	"os"
)

// expand-fs: expands the root filesystem on the given device (stub)
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "expand-fs: usage: expand-fs DEVICE")
		os.Exit(1)
	}
	device := os.Args[1]
	fmt.Printf("Would expand root filesystem on %s (not implemented)\n", device)
	// TODO: Implement real filesystem expansion logic
}
