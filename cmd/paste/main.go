package main

import (
	"bufio"
	"fmt"
	"os"
)

// paste: merges lines of files side by side, separated by tabs
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "paste: usage: paste FILE1 FILE2")
		os.Exit(1)
	}
	f1, err1 := os.Open(os.Args[1])
	if err1 != nil {
		fmt.Fprintln(os.Stderr, "paste: cannot open file:", os.Args[1], err1)
		os.Exit(1)
	}
	defer f1.Close()
	f2, err2 := os.Open(os.Args[2])
	if err2 != nil {
		fmt.Fprintln(os.Stderr, "paste: cannot open file:", os.Args[2], err2)
		os.Exit(1)
	}
	defer f2.Close()
	scan1 := bufio.NewScanner(f1)
	scan2 := bufio.NewScanner(f2)
	for scan1.Scan() || scan2.Scan() {
		s1, s2 := "", ""
		if scan1.Scan() {
			s1 = scan1.Text()
		}
		if scan2.Scan() {
			s2 = scan2.Text()
		}
		fmt.Printf("%s\t%s\n", s1, s2)
	}
}
