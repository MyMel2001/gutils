package main

import (
	"fmt"
	"os"
	"strconv"
)

// ps: lists PID and command name for all user processes (Linux only)
func main() {
	files, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ps: cannot read /proc:", err)
		os.Exit(1)
	}
	fmt.Printf("%-7s %-6s %s\n", "PID", "TTY", "CMD")
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(f.Name())
		if err != nil {
			continue
		}
		cmdline, err := os.ReadFile("/proc/" + f.Name() + "/comm")
		if err != nil {
			continue
		}
		// Get TTY
		tty := "?"
		if stat, err := os.Stat("/proc/" + f.Name() + "/fd/0"); err == nil {
			if link, err := os.Readlink("/proc/" + f.Name() + "/fd/0"); err == nil {
				if len(link) > 5 && link[:5] == "/dev/" {
					tty = link[5:]
				}
			}
			_ = stat
		}
		cmd := string(cmdline)
		cmd = cmd[:len(cmd)-1] // remove newline
		fmt.Printf("%-7d %-6s %s\n", pid, tty, cmd)
	}
}
