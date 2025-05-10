package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// ps: lists PID and command name for all user processes (Linux only)
func main() {
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ps: cannot read /proc:", err)
		os.Exit(1)
	}
	fmt.Printf("%5s %s\n", "PID", "CMD")
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(f.Name())
		if err != nil {
			continue
		}
		cmdline, err := ioutil.ReadFile("/proc/" + f.Name() + "/comm")
		if err != nil {
			continue
		}
		fmt.Printf("%5d %s", pid, cmdline)
	}
}
