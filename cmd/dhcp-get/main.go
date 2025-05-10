package main

import (
	"fmt"
	"os"
)

// dhcp-get: lists interfaces and prints a message about requesting DHCP (no real DHCP)
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "dhcp-get: usage: dhcp-get INTERFACE")
		os.Exit(1)
	}
	iface := os.Args[1]
	// List all interfaces (Linux: /sys/class/net/*)
	files, err := os.ReadDir("/sys/class/net")
	if err != nil {
		fmt.Fprintln(os.Stderr, "dhcp-get: cannot list interfaces:", err)
		os.Exit(1)
	}
	fmt.Println("Available interfaces:")
	for _, f := range files {
		fmt.Printf("  %s\n", f.Name())
	}
	fmt.Printf("Would request DHCP lease for interface '%s' (not implemented)\n", iface)
}
