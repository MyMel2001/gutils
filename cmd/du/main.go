package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// du: prints total size in bytes of files or directories
func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, path := range args {
		size, err := du(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "du: error:", path, err)
			continue
		}
		fmt.Printf("%d\t%s\n", size, path)
	}
}

// du returns the total size in bytes of the file or directory at path
func du(path string) (int64, error) {
	var total int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total, err
} 