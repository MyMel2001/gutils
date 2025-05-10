package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// killall: sends SIGTERM (or SIGKILL with -9) to all processes with a given name (Linux only)
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "killall: usage: killall [-9] NAME")
		os.Exit(1)
	}
	sig := syscall.SIGTERM
	nameIdx := 1
	if os.Args[1] == "-9" {
		sig = syscall.SIGKILL
		nameIdx = 2
	}
	if len(os.Args) <= nameIdx {
		fmt.Fprintln(os.Stderr, "killall: missing process name")
		os.Exit(1)
	}
	name := os.Args[nameIdx]
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		fmt.Fprintln(os.Stderr, "killall: cannot read /proc:", err)
		os.Exit(1)
	}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(f.Name())
		if err != nil {
			continue
		}
		comm, err := ioutil.ReadFile("/proc/" + f.Name() + "/comm")
		if err != nil {
			continue
		}
		if strings.TrimSpace(string(comm)) == name {
			err = syscall.Kill(pid, sig)
			if err != nil {
				fmt.Fprintln(os.Stderr, "killall:", err)
			}
		}
	}
}
