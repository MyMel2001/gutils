package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

// kill: sends a signal to a process (default SIGTERM, -9 for SIGKILL)
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "kill: usage: kill [-9] PID...")
		os.Exit(1)
	}
	sig := syscall.SIGTERM
	start := 1
	if os.Args[1] == "-9" {
		sig = syscall.SIGKILL
		start = 2
	}
	for _, arg := range os.Args[start:] {
		pid, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "kill: invalid pid:", arg)
			continue
		}
		err = syscall.Kill(pid, sig)
		if err != nil {
			fmt.Fprintln(os.Stderr, "kill:", err)
		}
	}
}
