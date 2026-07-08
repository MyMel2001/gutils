package main

import (
	"fmt"
	"os"
	"syscall"
)

// df: prints disk usage for the root filesystem
func main() {
	human := false
	for _, arg := range os.Args[1:] {
		if arg == "-h" {
			human = true
		}
	}

	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		fmt.Fprintln(os.Stderr, "df: error:", err)
		os.Exit(1)
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free

	if human {
		fmt.Printf("Filesystem  Size  Used Available\n")
		fmt.Printf("/ %s %s %s\n", humanSize(int64(total)), humanSize(int64(used)), humanSize(int64(free)))
	} else {
		fmt.Printf("Filesystem 1K-blocks Used Available\n")
		fmt.Printf("/ %d %d %d\n", total/1024, used/1024, free/1024)
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
