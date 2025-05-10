package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// find: recursively lists files and directories from a given path (default .)
func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			fmt.Println(path)
		}
		return nil
	})
}
