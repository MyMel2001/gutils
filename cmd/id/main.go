package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"
)

// id: prints the current user's UID, GID, and groups
func main() {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "id: cannot get current user:", err)
		os.Exit(1)
	}
	fmt.Printf("uid=%s(%s) gid=%s(%s)", u.Uid, u.Username, u.Gid, u.Username)
	groups, err := u.GroupIds()
	if err == nil && len(groups) > 0 {
		fmt.Print(" groups=" + strings.Join(groups, ","))
	}
	fmt.Println()
}
