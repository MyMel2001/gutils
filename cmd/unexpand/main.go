package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// unexpand: converts spaces to tabs (default 8 spaces per tab)
func main() {
	spaces := flag.Int("t", 8, "number of spaces per tab")
	flag.Parse()
	files := flag.Args()
	if len(files) == 0 {
		unexpandStream(os.Stdin, *spaces)
		return
	}
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unexpand: cannot open file:", fname, err)
			continue
		}
		unexpandStream(f, *spaces)
		f.Close()
	}
}

func unexpandStream(src *os.File, n int) {
	scanner := bufio.NewScanner(src)
	spacesStr := strings.Repeat(" ", n)
	for scanner.Scan() {
		fmt.Println(strings.ReplaceAll(scanner.Text(), spacesStr, "\t"))
	}
}
