package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// cat: prints file contents or stdin (standard Unix behavior)
func main() {
	if len(os.Args) < 2 {
		io.Copy(os.Stdout, os.Stdin)
		return
	}
	for _, fname := range os.Args[1:] {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cat: cannot open file:", fname, err)
			continue
		}
		io.Copy(os.Stdout, f)
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
