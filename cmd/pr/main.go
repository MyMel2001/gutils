package main

import (
	"bufio"
	"fmt"
	"os"
)

const linesPerPage = 60

// pr prints files with a header (filename and page number), paginated every 60 lines.
func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		printFile("-", os.Stdin)
	} else {
		for _, fname := range files {
			f, err := os.Open(fname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "pr: %s: %v\n", fname, err)
				continue
			}
			printFile(fname, f)
			f.Close()
		}
	}
}

// printFile prints the file with headers and pagination.
func printFile(name string, f *os.File) {
	scanner := bufio.NewScanner(f)
	lineCount := 0
	page := 1
	fmt.Printf("\n%s  Page %d\n\n", name, page)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		lineCount++
		if lineCount%linesPerPage == 0 {
			page++
			fmt.Printf("\n%s  Page %d\n\n", name, page)
		}
	}
}
