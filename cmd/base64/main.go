package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// base64: encodes or decodes files or stdin
func main() {
	args := os.Args[1:]
	decode := false
	files := args
	if len(args) > 0 && args[0] == "-d" {
		decode = true
		files = args[1:]
	}
	if len(files) == 0 {
		process(os.Stdin, decode)
		return
	}
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "base64: cannot open file:", fname, err)
			continue
		}
		process(f, decode)
		f.Close()
	}
}

// process encodes or decodes src to stdout
func process(src *os.File, decode bool) {
	if decode {
		dec := base64.NewDecoder(base64.StdEncoding, src)
		if _, err := io.Copy(os.Stdout, dec); err != nil {
			fmt.Fprintln(os.Stderr, "base64: decode error:", err)
		}
	} else {
		enc := base64.NewEncoder(base64.StdEncoding, os.Stdout)
		if _, err := io.Copy(enc, src); err != nil {
			fmt.Fprintln(os.Stderr, "base64: encode error:", err)
		}
		enc.Close()
	}
} 