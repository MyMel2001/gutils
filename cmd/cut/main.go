package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// cut: prints selected fields from each line, split by a delimiter (default tab)
func main() {
	delim := flag.String("d", "\t", "delimiter")
	field := flag.Int("f", 0, "field number (1-based)")
	flag.Parse()
	if *field <= 0 {
		fmt.Fprintln(os.Stderr, "cut: usage: cut -f N [-d DELIM] [FILE...]")
		os.Exit(1)
	}
	files := flag.Args()
	if len(files) == 0 {
		cutStream(os.Stdin, *delim, *field)
		return
	}
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cut: cannot open file:", fname, err)
			continue
		}
		cutStream(f, *delim, *field)
		f.Close()
	}
}

func cutStream(src *os.File, delim string, field int) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), delim)
		if field-1 < len(parts) {
			fmt.Println(parts[field-1])
		}
	}
}
