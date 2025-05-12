package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// cat: prints file contents or stdin, or just stdin if piped
func main() {
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped in, ignore arguments
		io.Copy(os.Stdout, os.Stdin)
		return
	}
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
