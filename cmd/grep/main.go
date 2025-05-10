package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// grep: prints lines matching a pattern
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "grep: usage: grep PATTERN [FILE...]")
		os.Exit(1)
	}
	pattern := os.Args[1]
	if len(os.Args) == 2 {
		grepStream(pattern, os.Stdin)
		return
	}
	for _, fname := range os.Args[2:] {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "grep: cannot open file:", fname, err)
			continue
		}
		grepStream(pattern, f)
		f.Close()
	}
}

// grepStream prints lines containing pattern from src
func grepStream(pattern string, src *os.File) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			fmt.Println(line)
		}
	}
} 