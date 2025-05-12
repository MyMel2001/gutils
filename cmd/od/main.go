package main

import (
	"fmt"
	"io"
	"os"
)

// od prints the octal dump of files or stdin, 16 bytes per line, with offsets.
func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		dump(os.Stdin)
	} else {
		for _, fname := range files {
			f, err := os.Open(fname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "od: %s: %v\n", fname, err)
				continue
			}
			dump(f)
			f.Close()
		}
	}
}

// dump prints the octal dump of the given reader.
func dump(r io.Reader) {
	buf := make([]byte, 16)
	offset := 0
	for {
		n, err := r.Read(buf)
		if n > 0 {
			fmt.Printf("%07o ", offset)
			for i := 0; i < n; i++ {
				fmt.Printf(" %03o", buf[i])
			}
			fmt.Println()
			offset += n
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "od: read error: %v\n", err)
			break
		}
	}
}
