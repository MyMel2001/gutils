package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// echo: prints its arguments, or stdin if piped
func main() {
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped in
		io.Copy(os.Stdout, os.Stdin)
		return
	}
	args := os.Args[1:]
	fmt.Println(strings.Join(args, " "))
}
