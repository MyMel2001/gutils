package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"
)

// sort: multithreaded sort utility
func main() {
	lines := []string{}
	if len(os.Args) < 2 {
		readLines(os.Stdin, &lines)
	} else {
		for _, fname := range os.Args[1:] {
			f, err := os.Open(fname)
			if err != nil {
				fmt.Fprintln(os.Stderr, "sort: cannot open file:", fname, err)
				continue
			}
			readLines(f, &lines)
			f.Close()
		}
	}

	// Parallel sort: split into chunks, sort in goroutines, then merge
	workers := runtime.NumCPU()
	chunkSize := (len(lines) + workers - 1) / workers
	chunks := [][]string{}
	for i := 0; i < len(lines); i += chunkSize {
		end := i + chunkSize
		if end > len(lines) {
			end = len(lines)
		}
		chunks = append(chunks, lines[i:end])
	}

	var wg sync.WaitGroup
	for i := range chunks {
		wg.Add(1)
		go func(c *[]string) {
			sort.Strings(*c)
			wg.Done()
		}(&chunks[i])
	}
	wg.Wait()

	// Merge sorted chunks
	result := mergeChunks(chunks)
	for _, line := range result {
		fmt.Println(line)
	}
}

// readLines reads lines from r and appends to lines
func readLines(r *os.File, lines *[]string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		*lines = append(*lines, scanner.Text())
	}
}

// mergeChunks merges multiple sorted slices into one sorted slice
func mergeChunks(chunks [][]string) []string {
	result := []string{}
	idx := make([]int, len(chunks))
	for {
		minIdx := -1
		var minVal string
		for i, c := range chunks {
			if idx[i] < len(c) {
				if minIdx == -1 || c[idx[i]] < minVal {
					minIdx = i
					minVal = c[idx[i]]
				}
			}
		}
		if minIdx == -1 {
			break
		}
		result = append(result, minVal)
		idx[minIdx]++
	}
	return result
} 