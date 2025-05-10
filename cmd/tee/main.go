package main

import (
	"fmt"
	"io"
	"os"
)

// tee: writes stdin to stdout and to files given as arguments
func main() {
	if len(os.Args) < 2 {
		io.Copy(os.Stdout, os.Stdin)
		return
	}
	files := make([]*os.File, 0, len(os.Args)-1)
	for _, fname := range os.Args[1:] {
		f, err := os.Create(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tee: cannot open file:", fname, err)
			continue
		}
		files = append(files, f)
		defer f.Close()
	}
	mw := io.MultiWriter(append([]io.Writer{os.Stdout}, toWriters(files)...)...)
	io.Copy(mw, os.Stdin)
}

func toWriters(fs []*os.File) []io.Writer {
	ws := make([]io.Writer, len(fs))
	for i, f := range fs {
		ws[i] = f
	}
	return ws
}
