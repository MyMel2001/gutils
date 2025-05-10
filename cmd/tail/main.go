package main

import (
	"bufio"
	"fmt"
	"os"
)

// tail: prints the last 10 lines of files or stdin
func main() {
	if len(os.Args) < 2 {
		tailStream(os.Stdin)
		return
	}
	for _, fname := range os.Args[1:] {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tail: cannot open file:", fname, err)
			continue
		}
		tailStream(f)
		f.Close()
	}
}

// tailStream prints the last 10 lines from src
func tailStream(src *os.File) {
	lines := make([]string, 0, 10)
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > 10 {
			lines = lines[1:]
		}
	}
	for _, line := range lines {
		fmt.Println(line)
	}
} 