package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// wc: prints line, word, and byte count. Supports -l, -w, -c flags
func main() {
	countLines := false
	countWords := false
	countBytes := false
	files := []string{}

	// Parse flags
	for _, arg := range os.Args[1:] {
		if arg == "-l" {
			countLines = true
		} else if arg == "-w" {
			countWords = true
		} else if arg == "-c" {
			countBytes = true
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// Parse combined flags like -lw, -lwc
			for _, f := range arg[1:] {
				switch f {
				case 'l':
					countLines = true
				case 'w':
					countWords = true
				case 'c':
					countBytes = true
				default:
					fmt.Fprintf(os.Stderr, "wc: invalid option -- '%c'\n", f)
					os.Exit(1)
				}
			}
		} else {
			files = append(files, arg)
		}
	}

	// If no flags specified, show all counts
	if !countLines && !countWords && !countBytes {
		countLines = true
		countWords = true
		countBytes = true
	}

	if len(files) == 0 {
		l, w, b := countStream(os.Stdin)
		printCounts(l, w, b, "", countLines, countWords, countBytes)
		return
	}

	totalLines, totalWords, totalBytes := 0, 0, 0
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wc: cannot open file:", fname, err)
			continue
		}
		l, w, b := countStream(f)
		f.Close()
		printCounts(l, w, b, fname, countLines, countWords, countBytes)
		totalLines += l
		totalWords += w
		totalBytes += b
	}
	if len(files) > 1 {
		printCounts(totalLines, totalWords, totalBytes, "total", countLines, countWords, countBytes)
	}
}

// countStream counts lines, words, and bytes from a reader
func countStream(src io.Reader) (int, int, int) {
	scanner := bufio.NewScanner(src)
	lines, words, bytes := 0, 0, 0
	for scanner.Scan() {
		lines++
		line := scanner.Text()
		words += len(strings.Fields(line))
		bytes += len(line) + 1
	}
	return lines, words, bytes
}

// printCounts prints the requested counts
func printCounts(lines, words, bytes int, name string, countLines, countWords, countBytes bool) {
	if countLines {
		fmt.Printf("%7d", lines)
	}
	if countWords {
		fmt.Printf("%7d", words)
	}
	if countBytes {
		fmt.Printf("%7d", bytes)
	}
	if name != "" {
		fmt.Printf(" %s", name)
	}
	fmt.Println()
}
