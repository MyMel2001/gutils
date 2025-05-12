package main

import (
	"fmt"
	"os"
	"syscall"
)

// mknod creates a special file (currently only supports FIFO pipes).
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: mknod NAME TYPE [MAJOR MINOR]")
		os.Exit(1)
	}
	name := os.Args[1]
	typeChar := os.Args[2]
	if typeChar == "p" {
		// Create a FIFO special file
		err := os.Mkdir(name, 0666)
		if err == nil {
			fmt.Fprintf(os.Stderr, "mknod: %s: is a directory, not a FIFO\n", name)
			os.Exit(1)
		}
		err = os.Remove(name)
		_ = err // ignore error
		err = syscall.Mkfifo(name, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mknod: %s: %v\n", name, err)
			os.Exit(1)
		}
		return
	}
	fmt.Fprintln(os.Stderr, "mknod: only FIFO (p) type is supported")
	os.Exit(1)
}
