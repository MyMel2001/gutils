package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// du: prints total size in bytes of files or directories
func main() {
	human := false
	args := os.Args[1:]
	filtered := []string{}
	for _, arg := range args {
		if arg == "-h" {
			human = true
		} else {
			filtered = append(filtered, arg)
		}
	}
	if len(filtered) == 0 {
		filtered = []string{"."}
	}
	for _, path := range filtered {
		size, err := du(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "du: error:", path, err)
			continue
		}
		if human {
			fmt.Printf("%s\t%s\n", humanSize(size), path)
		} else {
			fmt.Printf("%d\t%s\n", size, path)
		}
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

// humanSize returns a human-readable string for the given byte size
func humanSize(bytes int64) string {
	units := []string{"B", "K", "M", "G", "T", "P"}
	size := float64(bytes)
	unitIdx := 0
	for size >= 1024 && unitIdx < len(units)-1 {
		size /= 1024
		unitIdx++
	}
	if unitIdx == 0 {
		return fmt.Sprintf("%d", bytes)
	}
	return fmt.Sprintf("%.1f%s", size, units[unitIdx])
}
