package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// whereis: finds binaries in the current PATH
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "whereis: usage: whereis COMMAND...")
		os.Exit(1)
	}
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, ":")
	for _, cmd := range os.Args[1:] {
		found := false
		for _, dir := range paths {
			full := filepath.Join(dir, cmd)
			if fi, err := os.Stat(full); err == nil && !fi.IsDir() && (fi.Mode()&0111 != 0) {
				fmt.Println(full)
				found = true
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "whereis: %s not found\n", cmd)
		}
	}
} 