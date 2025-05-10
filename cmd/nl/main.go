package main

import (
	"bufio"
	"fmt"
	"os"
)

// nl: numbers lines of files or stdin
func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		nlStream(os.Stdin)
		return
	}
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "nl: cannot open file:", fname, err)
			continue
		}
		nlStream(f)
		f.Close()
	}
}

func nlStream(src *os.File) {
	scanner := bufio.NewScanner(src)
	n := 1
	for scanner.Scan() {
		fmt.Printf("%6d	%s\n", n, scanner.Text())
		n++
	}
}
