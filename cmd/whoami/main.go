package main

import (
	"fmt"
	"os"
	"os/user"
)

// whoami: prints the current user name
func main() {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "whoami:", err)
		os.Exit(1)
	}
	fmt.Println(u.Username)
} 