package main

import (
	"fmt"
	"io"
	"os"
)

// cp: copies a file from source to destination
func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "cp: usage: cp SOURCE DEST")
		os.Exit(1)
	}
	src, dst := os.Args[1], os.Args[2]
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cp: cannot open source:", err)
		os.Exit(1)
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cp: cannot create destination:", err)
		os.Exit(1)
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cp: error copying:", err)
		os.Exit(1)
	}
} 