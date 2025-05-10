package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// lsblk: lists block devices (name and size in bytes)
func main() {
	blockDir := "/sys/block"
	entries, err := ioutil.ReadDir(blockDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "lsblk: cannot read /sys/block:", err)
		os.Exit(1)
	}
	for _, entry := range entries {
		name := entry.Name()
		sizePath := filepath.Join(blockDir, name, "size")
		sizeBytes := "?"
		if data, err := ioutil.ReadFile(sizePath); err == nil {
			// Size is in 512-byte sectors
			trimmed := strings.TrimSpace(string(data))
			if sectors, err := parseUint(trimmed); err == nil {
				size := sectors * 512
				sizeBytes = fmt.Sprintf("%d", size)
			}
		}
		fmt.Printf("%s\t%s\n", name, sizeBytes)
	}
}

// parseUint parses a string as uint64
func parseUint(s string) (uint64, error) {
	var n uint64
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
} 