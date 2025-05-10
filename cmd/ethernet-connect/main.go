package main

import (
	"fmt"
	"os"
)

// ethernet-connect: lists ethernet interfaces and prints a message about connecting (no real connection)
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "ethernet-connect: usage: ethernet-connect INTERFACE")
		os.Exit(1)
	}
	iface := os.Args[1]
	// List likely ethernet interfaces (Linux: /sys/class/net/*, skip loopback and wireless)
	files, err := os.ReadDir("/sys/class/net")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ethernet-connect: cannot list interfaces:", err)
		os.Exit(1)
	}
	fmt.Println("Available ethernet interfaces:")
	any := false
	for _, f := range files {
		if f.Name() == "lo" {
			continue
		}
		if _, err := os.Stat("/sys/class/net/" + f.Name() + "/wireless"); err == nil {
			continue
		}
		fmt.Printf("  %s\n", f.Name())
		any = true
	}
	if !any {
		fmt.Println("  (none found)")
	}
	fmt.Printf("Would connect ethernet interface '%s' (not implemented)\n", iface)
}
