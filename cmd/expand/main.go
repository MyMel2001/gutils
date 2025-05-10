package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// expand: converts tabs to spaces (default 8 spaces per tab)
func main() {
	spaces := flag.Int("t", 8, "number of spaces per tab")
	flag.Parse()
	files := flag.Args()
	if len(files) == 0 {
		expandStream(os.Stdin, *spaces)
		return
	}
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "expand: cannot open file:", fname, err)
			continue
		}
		expandStream(f, *spaces)
		f.Close()
	}
}

func expandStream(src *os.File, n int) {
	scanner := bufio.NewScanner(src)
	tab := strings.Repeat(" ", n)
	for scanner.Scan() {
		fmt.Println(strings.ReplaceAll(scanner.Text(), "\t", tab))
	}
}
