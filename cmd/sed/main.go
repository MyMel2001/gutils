package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// sed supports basic 's/old/new/g' substitution for each line of input.
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: sed 's/old/new/g' [FILE...]")
		os.Exit(1)
	}
	cmd := os.Args[1]
	if !strings.HasPrefix(cmd, "s/") || !strings.HasSuffix(cmd, "/g") {
		fmt.Fprintln(os.Stderr, "Only 's/old/new/g' supported")
		os.Exit(1)
	}
	parts := strings.Split(cmd[2:len(cmd)-2], "/")
	if len(parts) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid substitution format")
		os.Exit(1)
	}
	old, new := parts[0], parts[1]
	re, err := regexp.Compile(regexp.QuoteMeta(old))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid regex: %v\n", err)
		os.Exit(1)
	}
	files := os.Args[2:]
	if len(files) == 0 {
		processSed(os.Stdin, re, new)
	} else {
		for _, fname := range files {
			f, err := os.Open(fname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "sed: %s: %v\n", fname, err)
				continue
			}
			processSed(f, re, new)
			f.Close()
		}
	}
}

// processSed reads lines from r, applies the substitution, and writes to stdout.
func processSed(r *os.File, re *regexp.Regexp, new string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(re.ReplaceAllString(line, new))
	}
}
