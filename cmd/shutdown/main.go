//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// shutdown: shuts down or reboots the system with optional delay, using syscalls only
func main() {
	// Parse command-line arguments manually (no flag package to keep it simple)
	reboot := false
	delay := 0
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-r":
			reboot = true
		case "-t":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &delay)
				i++
			}
		case "-h", "--help":
			fmt.Println("shutdown: usage: shutdown [-r] [-t seconds]")
			fmt.Println("  -r          reboot instead of shutdown")
			fmt.Println("  -t seconds  delay in seconds before shutdown/reboot")
			os.Exit(0)
		}
	}

	// Wait for the specified delay, if any
	if delay > 0 {
		action := "shutdown"
		if reboot {
			action = "reboot"
		}
		fmt.Printf("System will %s in %d seconds...\n", action, delay)
		for i := delay; i > 0; i-- {
			if i <= 10 || i%10 == 0 {
				fmt.Printf("%s in %d seconds...\n", action, i)
			}
			time.Sleep(1 * time.Second)
		}
	}

	// Sync filesystems before reboot/shutdown
	syscall.Sync()

	// Use syscall.Reboot to shutdown or reboot
	var how int
	if reboot {
		how = syscall.LINUX_REBOOT_CMD_RESTART
	} else {
		how = syscall.LINUX_REBOOT_CMD_POWER_OFF
	}

	err := syscall.Reboot(how)
	if err != nil {
		fmt.Fprintf(os.Stderr, "shutdown: failed: %v\n", err)
		os.Exit(1)
	}
}
