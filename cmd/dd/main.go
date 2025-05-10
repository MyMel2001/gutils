package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// dd: copies data from if=INFILE to of=OUTFILE, with optional bs=BLOCKSIZE, count=N
func main() {
	infile, outfile := "", ""
	bs := 512
	count := -1
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "if=") {
			infile = arg[3:]
		} else if strings.HasPrefix(arg, "of=") {
			outfile = arg[3:]
		} else if strings.HasPrefix(arg, "bs=") {
			bsval, err := strconv.Atoi(arg[3:])
			if err == nil && bsval > 0 {
				bs = bsval
			}
		} else if strings.HasPrefix(arg, "count=") {
			cnt, err := strconv.Atoi(arg[6:])
			if err == nil && cnt >= 0 {
				count = cnt
			}
		}
	}
	if infile == "" || outfile == "" {
		fmt.Fprintln(os.Stderr, "dd: usage: dd if=INFILE of=OUTFILE [bs=BLOCKSIZE] [count=N]")
		os.Exit(1)
	}
	in, err := os.Open(infile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dd: cannot open input:", err)
		os.Exit(1)
	}
	defer in.Close()
	out, err := os.Create(outfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dd: cannot open output:", err)
		os.Exit(1)
	}
	defer out.Close()
	buf := make([]byte, bs)
	written := 0
	for count < 0 || written < count {
		n, err := in.Read(buf)
		if n > 0 {
			wn, werr := out.Write(buf[:n])
			if werr != nil || wn != n {
				fmt.Fprintln(os.Stderr, "dd: write error:", werr)
				os.Exit(1)
			}
			written++
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "dd: read error:", err)
			os.Exit(1)
		}
	}
} 