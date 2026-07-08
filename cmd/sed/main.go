package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// sed supports basic 's/old/new/g' substitution and in-place editing with -i
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: sed [-i] 's/old/new/g' [FILE...]")
		os.Exit(1)
	}

	inPlace := false
	cmdIndex := 1

	if os.Args[1] == "-i" {
		inPlace = true
		cmdIndex = 2
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: sed -i 's/old/new/g' [FILE...]")
			os.Exit(1)
		}
	}

	cmd := os.Args[cmdIndex]
	if !strings.HasPrefix(cmd, "s/") || !strings.HasSuffix(cmd, "/g") {
		fmt.Fprintln(os.Stderr, "Only 's/old/new/g' supported")
		os.Exit(1)
	}

	// Extract old and new patterns, handling escaped slashes
	inner := cmd[2 : len(cmd)-2]
	parts := splitSubst(inner)
	if len(parts) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid substitution format")
		os.Exit(1)
	}
	old, new := parts[0], parts[1]

	re, err := regexp.Compile(old)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid regex: %v\n", err)
		os.Exit(1)
	}

	files := os.Args[cmdIndex+1:]
	if len(files) == 0 {
		processSed(os.Stdin, re, new, "")
	} else {
		for _, fname := range files {
			f, err := os.Open(fname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "sed: %s: %v\n", fname, err)
				continue
			}
			if inPlace {
				// Read all lines, modify, write back
				var lines []string
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					lines = append(lines, re.ReplaceAllString(scanner.Text(), new))
				}
				f.Close()
				if err := os.WriteFile(fname, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
					fmt.Fprintf(os.Stderr, "sed: %s: write error: %v\n", fname, err)
				}
			} else {
				processSed(f, re, new, fname)
				f.Close()
			}
		}
	}
}

// splitSubst splits a substitution pattern by /, respecting escaped slashes
func splitSubst(s string) []string {
	var parts []string
	var current strings.Builder
	escape := false
	for _, c := range s {
		if escape {
			if c == '/' {
				current.WriteByte('/')
			} else {
				current.WriteByte('\\')
				current.WriteRune(c)
			}
			escape = false
			continue
		}
		if c == '\\' {
			escape = true
			continue
		}
		if c == '/' {
			parts = append(parts, current.String())
			current.Reset()
			continue
		}
		current.WriteRune(c)
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}

// processSed reads lines from r, applies the substitution, and writes to stdout.
func processSed(r *os.File, re *regexp.Regexp, new string, filename string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		result := re.ReplaceAllString(line, new)
		if filename != "" {
			// No filename prefix for sed by default
		}
		fmt.Println(result)
	}
}
