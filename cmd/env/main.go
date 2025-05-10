package main

import (
	"fmt"
	"os"
)

// env: prints all environment variables
func main() {
	for _, e := range os.Environ() {
		fmt.Println(e)
	}
}
