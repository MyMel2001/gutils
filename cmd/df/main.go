package main

import (
	"fmt"
	"os"
	"syscall"
)

// df: prints disk usage for the root filesystem
func main() {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		fmt.Fprintln(os.Stderr, "df: error:", err)
		os.Exit(1)
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free
	fmt.Printf("Filesystem 1K-blocks Used Available\n")
	fmt.Printf("/ %d %d %d\n", total/1024, used/1024, free/1024)
}
