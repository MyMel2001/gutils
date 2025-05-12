package main

import (
	"fmt"
	"io"
	"os"
)

// cmp compares two files byte by byte and prints the first difference, or nothing if identical.
func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: cmp FILE1 FILE2")
		os.Exit(1)
	}
	f1, err1 := os.Open(os.Args[1])
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "cmp: %s: %v\n", os.Args[1], err1)
		os.Exit(1)
	}
	defer f1.Close()
	f2, err2 := os.Open(os.Args[2])
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "cmp: %s: %v\n", os.Args[2], err2)
		os.Exit(1)
	}
	defer f2.Close()
	buf1 := make([]byte, 4096)
	buf2 := make([]byte, 4096)
	offset := 0
	for {
		n1, e1 := f1.Read(buf1)
		n2, e2 := f2.Read(buf2)
		min := n1
		if n2 < min {
			min = n2
		}
		for i := 0; i < min; i++ {
			if buf1[i] != buf2[i] {
				fmt.Printf("cmp: %s %s differ: byte %d\n", os.Args[1], os.Args[2], offset+i+1)
				os.Exit(1)
			}
		}
		if n1 != n2 {
			fmt.Printf("cmp: EOF on %s at byte %d\n", os.Args[1], offset+min+1)
			os.Exit(1)
		}
		if e1 == io.EOF && e2 == io.EOF {
			break
		}
		if e1 != nil && e1 != io.EOF {
			fmt.Fprintf(os.Stderr, "cmp: %s: %v\n", os.Args[1], e1)
			os.Exit(1)
		}
		if e2 != nil && e2 != io.EOF {
			fmt.Fprintf(os.Stderr, "cmp: %s: %v\n", os.Args[2], e2)
			os.Exit(1)
		}
		offset += min
	}
	// No differences found
	os.Exit(0)
}
