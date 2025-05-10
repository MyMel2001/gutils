package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// sleep: sleeps for the given number of seconds
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "sleep: usage: sleep SECONDS")
		os.Exit(1)
	}
	secs, err := strconv.Atoi(os.Args[1])
	if err != nil || secs < 0 {
		fmt.Fprintln(os.Stderr, "sleep: invalid time:", os.Args[1])
		os.Exit(1)
	}
	time.Sleep(time.Duration(secs) * time.Second)
} 