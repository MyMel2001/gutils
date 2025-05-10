package main

import (
	"fmt"
	"os"
	"time"
)

// touch: creates files or updates modification time
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "touch: usage: touch FILE...")
		os.Exit(1)
	}
	for _, fname := range os.Args[1:] {
		// Try to open the file, create if not exists
		f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "touch: cannot open file:", fname, err)
			continue
		}
		f.Close()
		// Update the modification and access time to now
		err = os.Chtimes(fname, time.Now(), time.Now())
		if err != nil {
			fmt.Fprintln(os.Stderr, "touch: cannot update times:", fname, err)
		}
	}
} 