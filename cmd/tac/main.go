package main

import (
	"bufio"
	"fmt"
	"os"
)

// tac: prints lines of files or stdin in reverse order
func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		tacStream(os.Stdin)
		return
	}
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tac: cannot open file:", fname, err)
			continue
		}
		tacStream(f)
		f.Close()
	}
}

func tacStream(src *os.File) {
	scanner := bufio.NewScanner(src)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for i := len(lines) - 1; i >= 0; i-- {
		fmt.Println(lines[i])
	}
}
