package main

import (
	"fmt"
	"os"
)

// stat: prints file information (size, mode, mod time, name)
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "stat: usage: stat FILE...")
		os.Exit(1)
	}
	for _, fname := range os.Args[1:] {
		info, err := os.Stat(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "stat: cannot stat:", fname, err)
			continue
		}
		fmt.Printf("  File: %s\n", fname)
		fmt.Printf("  Size: %d\tMode: %s\tModified: %s\n", info.Size(), info.Mode(), info.ModTime())
	}
}
