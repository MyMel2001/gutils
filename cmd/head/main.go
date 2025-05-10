package main

import (
	"bufio"
	"fmt"
	"os"
)

// head: prints the first 10 lines of files or stdin
func main() {
	if len(os.Args) < 2 {
		headStream(os.Stdin)
		return
	}
	for _, fname := range os.Args[1:] {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "head: cannot open file:", fname, err)
			continue
		}
		headStream(f)
		f.Close()
	}
}

// headStream prints the first 10 lines from src
func headStream(src *os.File) {
	scanner := bufio.NewScanner(src)
	for i := 0; i < 10 && scanner.Scan(); i++ {
		fmt.Println(scanner.Text())
	}
} 