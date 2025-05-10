package main

import (
	"fmt"
	"os"
	"strconv"
)

// seq: prints a sequence of numbers from first to last (default step 1, default first=1)
func main() {
	args := os.Args[1:]
	first, step, last := 1, 1, 0
	if len(args) == 1 {
		last, _ = strconv.Atoi(args[0])
	} else if len(args) == 2 {
		first, _ = strconv.Atoi(args[0])
		last, _ = strconv.Atoi(args[1])
	} else if len(args) == 3 {
		first, _ = strconv.Atoi(args[0])
		step, _ = strconv.Atoi(args[1])
		last, _ = strconv.Atoi(args[2])
	} else {
		fmt.Fprintln(os.Stderr, "seq: usage: seq [FIRST] [STEP] LAST")
		os.Exit(1)
	}
	if step == 0 {
		fmt.Fprintln(os.Stderr, "seq: step cannot be zero")
		os.Exit(1)
	}
	for i := first; (step > 0 && i <= last) || (step < 0 && i >= last); i += step {
		fmt.Println(i)
	}
}
