package main

import (
	"fmt"
	"os"
)

// yes: prints 'y' or a given string repeatedly until killed
func main() {
	str := "y"
	if len(os.Args) > 1 {
		str = os.Args[1]
	}
	for {
		fmt.Println(str)
	}
}
