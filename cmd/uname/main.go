package main

import (
	"fmt"
	"runtime"
)

// uname: prints the system name (kernel name)
func main() {
	fmt.Println(runtime.GOOS)
}
