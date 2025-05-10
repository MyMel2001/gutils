package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

// split: splits a file into pieces of N lines each (default 1000 lines)
func main() {
	lines := flag.Int("l", 1000, "lines per file")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "split: usage: split [-l LINES] FILE")
		os.Exit(1)
	}
	fname := flag.Arg(0)
	f, err := os.Open(fname)
	if err != nil {
		fmt.Fprintln(os.Stderr, "split: cannot open file:", fname, err)
		os.Exit(1)
	}
	defer f.Close()
	splitStream(f, *lines)
}

func splitStream(src *os.File, n int) {
	scanner := bufio.NewScanner(src)
	fileIdx, lineIdx := 0, 0
	var out *os.File
	for scanner.Scan() {
		if lineIdx == 0 {
			name := fmt.Sprintf("x%02d", fileIdx)
			out, _ = os.Create(name)
			fileIdx++
		}
		fmt.Fprintln(out, scanner.Text())
		lineIdx++
		if lineIdx >= n {
			out.Close()
			lineIdx = 0
		}
	}
	if out != nil {
		out.Close()
	}
}
