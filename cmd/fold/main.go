package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

// fold: wraps long lines to a specified width (default 80)
func main() {
	width := flag.Int("w", 80, "wrap width")
	flag.Parse()
	files := flag.Args()
	if len(files) == 0 {
		foldStream(os.Stdin, *width)
		return
	}
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "fold: cannot open file:", fname, err)
			continue
		}
		foldStream(f, *width)
		f.Close()
	}
}

func foldStream(src *os.File, width int) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		for len(line) > width {
			fmt.Println(line[:width])
			line = line[width:]
		}
		fmt.Println(line)
	}
}
