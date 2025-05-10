package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

// chown: changes file owner and group
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "chown: usage: chown OWNER[:GROUP] FILE...")
		os.Exit(1)
	}
	ownerGroup := os.Args[1]
	parts := strings.SplitN(ownerGroup, ":", 2)
	uid, gid := -1, -1
	usr, err := user.Lookup(parts[0])
	if err == nil {
		uid64, _ := strconv.Atoi(usr.Uid)
		uid = uid64
	}
	if len(parts) == 2 {
		grp, err := user.LookupGroup(parts[1])
		if err == nil {
			gid64, _ := strconv.Atoi(grp.Gid)
			gid = gid64
		}
	}
	for _, fname := range os.Args[2:] {
		// Use -1 for unchanged uid/gid
		if err := os.Chown(fname, uid, gid); err != nil {
			fmt.Fprintln(os.Stderr, "chown: cannot change owner/group:", fname, err)
		}
	}
} 