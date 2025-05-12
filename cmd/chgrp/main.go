package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// chgrp changes the group ownership of files or directories specified as arguments.
func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: chgrp GROUP FILE...")
		os.Exit(1)
	}
	group := os.Args[1]
	grp, err := user.LookupGroup(group)
	if err != nil {
		gid, err2 := strconv.Atoi(group)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "chgrp: invalid group: %s\n", group)
			os.Exit(1)
		}
		changeGroup(gid, os.Args[2:])
		return
	}
	gid, err := strconv.Atoi(grp.Gid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chgrp: invalid group id: %s\n", grp.Gid)
		os.Exit(1)
	}
	changeGroup(gid, os.Args[2:])
}

// changeGroup changes the group of each file to gid.
func changeGroup(gid int, files []string) {
	status := 0
	for _, file := range files {
		if err := syscall.Chown(file, -1, gid); err != nil {
			fmt.Fprintf(os.Stderr, "chgrp: cannot change group of '%s': %v\n", file, err)
			status = 1
		}
	}
	os.Exit(status)
}
