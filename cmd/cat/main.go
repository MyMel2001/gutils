package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// cat: prints file contents or stdin
func main() {
	if len(os.Args) < 2 {
		copyStream(os.Stdin, os.Stdout)
		return
	}
	for _, fname := range os.Args[1:] {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cat: cannot open file:", fname, err)
			continue
		}
		copyStream(f, os.Stdout)
		f.Close()
	}
}

// copyStream copies from src to dst
func copyStream(src io.Reader, dst io.Writer) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		dst.Write(append(scanner.Bytes(), '\n'))
	}
} 