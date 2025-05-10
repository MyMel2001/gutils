package main

import (
	"fmt"
	"time"
)

// date: prints the current date and time
func main() {
	fmt.Println(time.Now().Format(time.RFC1123))
}
