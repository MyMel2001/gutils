package main

import (
	"fmt"
	"os"
	"syscall"
)

// mount: mounts a filesystem. Usage: mount SOURCE TARGET FSTYPE [OPTIONS]
func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "mount: usage: mount SOURCE TARGET FSTYPE [OPTIONS]")
		os.Exit(1)
	}
	source := os.Args[1]
	target := os.Args[2]
	fstype := os.Args[3]
	options := ""
	if len(os.Args) > 4 {
		options = os.Args[4]
	}
	// Perform the mount syscall
	err := syscall.Mount(source, target, fstype, 0, options)
	if err != nil {
		fmt.Fprintln(os.Stderr, "mount: error:", err)
		os.Exit(1)
	}
}
