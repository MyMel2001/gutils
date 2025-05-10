package main

import (
	"fmt"
	"os"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// ethernet-connect: brings up the specified ethernet interface using netlink
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "ethernet-connect: usage: ethernet-connect INTERFACE")
		os.Exit(1)
	}
	ifaceName := os.Args[1]

	links, err := netlink.LinkList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ethernet-connect: cannot list interfaces: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Available ethernet interfaces:")
	found := false
	for _, link := range links {
		attrs := link.Attrs()
		if attrs == nil || attrs.Name == "lo" || (attrs.Flags&unix.IFF_LOOPBACK != 0) {
			continue
		}
		if attrs.EncapType == "ether" {
			fmt.Printf("  %s\n", attrs.Name)
			if attrs.Name == ifaceName {
				found = true
			}
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "ethernet-connect: interface '%s' not found or not ethernet\n", ifaceName)
		os.Exit(1)
	}
	link, err := netlink.LinkByName(ifaceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ethernet-connect: failed to get interface '%s': %v\n", ifaceName, err)
		os.Exit(1)
	}
	if err := netlink.LinkSetUp(link); err != nil {
		fmt.Fprintf(os.Stderr, "ethernet-connect: failed to bring up '%s': %v\n", ifaceName, err)
		os.Exit(1)
	}
	fmt.Printf("Successfully brought up ethernet interface '%s'\n", ifaceName)
}
