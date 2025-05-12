package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"
)

// shutdown: shuts down or reboots the system with optional delay, using syscalls only
func main() {
	// Parse command-line flags
	reboot := flag.Bool("r", false, "reboot instead of shutdown")
	delay := flag.Int("t", 0, "delay in seconds before shutdown/reboot")
	flag.Parse()

	// Wait for the specified delay, if any
	if *delay > 0 {
		fmt.Printf("Waiting %d seconds before %s...\n", *delay, action(*reboot))
		time.Sleep(time.Duration(*delay) * time.Second)
	}

	// Use syscall.Reboot to shutdown or reboot
	var how int
	if *reboot {
		how = syscall.LINUX_REBOOT_CMD_RESTART
	} else {
		how = syscall.LINUX_REBOOT_CMD_POWER_OFF
	}

	// Need to call syscall.Reboot with magic numbers
	err := syscall.Reboot(how)
	if err != nil {
		fmt.Fprintf(os.Stderr, "shutdown: failed to %s: %v\n", action(*reboot), err)
		os.Exit(1)
	}
}

// action returns the string "reboot" or "shutdown" based on the flag
func action(reboot bool) string {
	if reboot {
		return "reboot"
	}
	return "shutdown"
}
