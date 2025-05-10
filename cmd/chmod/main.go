package main

import (
	"fmt"
	"os"
	"strconv"
)

// chmod: changes file permissions
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "chmod: usage: chmod MODE FILE...")
		os.Exit(1)
	}
	modeStr := os.Args[1]
	mode, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		fmt.Fprintln(os.Stderr, "chmod: invalid mode:", modeStr)
		os.Exit(1)
	}
	for _, fname := range os.Args[2:] {
		err := os.Chmod(fname, os.FileMode(mode))
		if err != nil {
			fmt.Fprintln(os.Stderr, "chmod: cannot change permissions:", fname, err)
		}
	}
} 