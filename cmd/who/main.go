package main

import (
	"fmt"
	"os/user"
)

// who: lists logged-in users (minimal: just prints current user)
func main() {
	u, err := user.Current()
	if err != nil {
		fmt.Println("who: cannot get current user:", err)
		return
	}
	fmt.Println(u.Username)
}
