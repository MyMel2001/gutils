package main

import (
	"fmt"
	"net"
	"os"
)

// ip: lists network interfaces and their addresses
func main() {
	ifs, err := net.Interfaces()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ip: error:", err)
		os.Exit(1)
	}
	for _, iface := range ifs {
		fmt.Printf("%s\n", iface.Name)
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ip: %s: %v\n", iface.Name, err)
			continue
		}
		for _, addr := range addrs {
			fmt.Printf("  %s\n", addr.String())
		}
	}
}
