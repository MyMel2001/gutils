package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// tr: translates characters in set1 to set2 for stdin or files
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "tr: usage: tr SET1 SET2 [FILE...]")
		os.Exit(1)
	}
	set1, set2 := os.Args[1], os.Args[2]
	files := os.Args[3:]
	if len(files) == 0 {
		trStream(os.Stdin, set1, set2)
		return
	}
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tr: cannot open file:", fname, err)
			continue
		}
		trStream(f, set1, set2)
		f.Close()
	}
}

func trStream(src *os.File, set1, set2 string) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Map(func(r rune) rune {
			idx := strings.IndexRune(set1, r)
			if idx >= 0 && idx < len(set2) {
				return rune(set2[idx])
			}
			return r
		}, line)
		fmt.Println(line)
	}
}
