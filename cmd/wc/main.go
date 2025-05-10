package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// wc: prints line, word, and byte count
func main() {
	if len(os.Args) < 2 {
		wcStream(os.Stdin, "")
		return
	}
	for _, fname := range os.Args[1:] {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wc: cannot open file:", fname, err)
			continue
		}
		wcStream(f, fname)
		f.Close()
	}
}

// wcStream prints line, word, and byte count from src
func wcStream(src io.Reader, name string) {
	scanner := bufio.NewScanner(src)
	lines, words, bytes := 0, 0, 0
	for scanner.Scan() {
		lines++
		line := scanner.Text()
		words += len(strings.Fields(line))
		bytes += len(line) + 1 // +1 for newline
	}
	if name != "" {
		fmt.Printf("%d %d %d %s\n", lines, words, bytes, name)
	} else {
		fmt.Printf("%d %d %d\n", lines, words, bytes)
	}
} 