package main

import (
	"fmt"
	"os"
)

// ls: lists files in the current directory
func main() {
	human := false
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-h" {
			human = true
		}
	}

	dir, err := os.Open(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ls: cannot open directory:", err)
		os.Exit(1)
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ls: cannot read directory:", err)
		os.Exit(1)
	}

	for _, name := range files {
		if human {
			info, err := os.Stat(name)
			if err != nil {
				fmt.Println(name)
				continue
			}
			fmt.Printf("%s %s\n", humanSize(info.Size()), name)
		} else {
			fmt.Println(name)
		}
	}
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
