package main

import (
	"fmt"
	"os"
)

// shortcut: creates hard or symbolic links (like ln)
func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "shortcut: usage: shortcut [-s] TARGET LINK_NAME")
		os.Exit(1)
	}
	if args[0] == "-s" {
		if len(args) != 3 {
			fmt.Fprintln(os.Stderr, "shortcut: usage: shortcut -s TARGET LINK_NAME")
			os.Exit(1)
		}
		target, linkName := args[1], args[2]
		if err := os.Symlink(target, linkName); err != nil {
			fmt.Fprintln(os.Stderr, "shortcut: cannot create symlink:", err)
			os.Exit(1)
		}
	} else {
		target, linkName := args[0], args[1]
		if err := os.Link(target, linkName); err != nil {
			fmt.Fprintln(os.Stderr, "shortcut: cannot create hard link:", err)
			os.Exit(1)
		}
	}
} 