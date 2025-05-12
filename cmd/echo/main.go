package main

import (
	"fmt"
	"os"
	"strings"
)

// echo: prints its arguments (standard Unix behavior)
func main() {
	args := os.Args[1:]
	fmt.Println(strings.Join(args, " "))
}
