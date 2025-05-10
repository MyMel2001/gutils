package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

// sum: prints a simple checksum and block count for files or stdin
func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		sumStream(os.Stdin, "-")
		return
	}
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, "sum: cannot open file:", fname, err)
			continue
		}
		sumStream(f, fname)
		f.Close()
	}
}

func sumStream(src *os.File, name string) {
	h := crc32.NewIEEE()
	blocks := 0
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			h.Write(buf[:n])
			blocks++
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "sum: read error:", err)
			return
		}
	}
	fmt.Printf("%d %d %s\n", h.Sum32(), blocks, name)
}
