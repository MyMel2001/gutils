package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// sha256sum: prints SHA-256 hash of files or stdin
func main() {
	if len(os.Args) < 2 {
		printHash(os.Stdin, "-")
		return
	}
	for _, fname := range os.Args[1:] {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "sha256sum: cannot open file:", fname, err)
			continue
		}
		printHash(f, fname)
		f.Close()
	}
}

// printHash prints the SHA-256 hash of src
func printHash(src io.Reader, name string) {
	h := sha256.New()
	if _, err := io.Copy(h, src); err != nil {
		fmt.Fprintln(os.Stderr, "sha256sum: error reading:", name, err)
		return
	}
	sum := hex.EncodeToString(h.Sum(nil))
	fmt.Printf("%s  %s\n", sum, name)
} 