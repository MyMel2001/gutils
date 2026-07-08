package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// grep: prints lines matching a pattern (supports regex and basic text search)
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "grep: usage: grep [-i] [-E] PATTERN [FILE...]")
		os.Exit(1)
	}

	ignoreCase := false
	extendedRegex := false
	pattern := ""
	fileStart := 1

	// Parse flags
	for fileStart < len(os.Args) {
		arg := os.Args[fileStart]
		if arg == "-i" {
			ignoreCase = true
			fileStart++
		} else if arg == "-E" {
			extendedRegex = true
			fileStart++
		} else if strings.HasPrefix(arg, "-") {
			fmt.Fprintf(os.Stderr, "grep: unknown flag: %s\n", arg)
			os.Exit(1)
		} else {
			pattern = arg
			fileStart++
			break
		}
	}

	if pattern == "" {
		fmt.Fprintln(os.Stderr, "grep: usage: grep [-i] [-E] PATTERN [FILE...]")
		os.Exit(1)
	}

	// Compile regex
	var re *regexp.Regexp
	var err error
	if extendedRegex {
		re, err = regexp.Compile(pattern)
	} else {
		// Treat as literal string, escape all regex metacharacters
		re, err = regexp.Compile(regexp.QuoteMeta(pattern))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "grep: invalid pattern: %v\n", err)
		os.Exit(1)
	}

	if ignoreCase {
		re = regexp.MustCompile("(?i)" + re.String())
	}

	if fileStart >= len(os.Args) {
		grepStream(re, os.Stdin, "")
		return
	}

	multipleFiles := len(os.Args) - fileStart > 1
	for _, fname := range os.Args[fileStart:] {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "grep: cannot open file:", fname, err)
			continue
		}
		grepStream(re, f, fname)
		if multipleFiles {
			// Print filename header
		}
		f.Close()
	}
}

// grepStream prints lines matching pattern from src
func grepStream(re *regexp.Regexp, src *os.File, filename string) {
	scanner := bufio.NewScanner(src)
	// Increase buffer for long lines
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			if filename != "" {
				fmt.Printf("%s:", filename)
			}
			fmt.Println(line)
		}
	}
}
