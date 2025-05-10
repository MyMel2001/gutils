package main

import (
	"bufio"
	"fmt"
	"os"
)

// uniq: removes adjacent duplicate lines from stdin or a file
func main() {
	var f *os.File
	var err error
	if len(os.Args) > 1 {
		f, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "uniq: cannot open file:", err)
			os.Exit(1)
		}
		defer f.Close()
	} else {
		f = os.Stdin
	}
	last := ""
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line != last {
			fmt.Println(line)
			last = line
		}
	}
}
