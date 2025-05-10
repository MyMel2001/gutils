package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// dhcp-get: requests a DHCP lease for the specified interface and prints lease info
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "dhcp-get: usage: dhcp-get INTERFACE")
		os.Exit(1)
	}
	iface := os.Args[1]

	// Create a UDP connection on port 68 (DHCP client port)
	laddr := &net.UDPAddr{IP: net.IPv4zero, Port: 68}
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: failed to listen on UDP port 68: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Build DHCP Discover packet for the interface
	discover, err := dhcpv4.NewDiscoveryForInterface(iface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: failed to build DHCPDISCOVER: %v\n", err)
		os.Exit(1)
	}

	// Send DHCP Discover to broadcast address
	bcast := &net.UDPAddr{IP: net.IPv4bcast, Port: 67}
	if _, err := conn.WriteTo(discover.ToBytes(), bcast); err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: failed to send DHCPDISCOVER: %v\n", err)
		os.Exit(1)
	}

	// Set a timeout for receiving the offer
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buf := make([]byte, 1500)
	n, _, err := conn.ReadFrom(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: failed to receive DHCPOFFER: %v\n", err)
		os.Exit(1)
	}

	offer, err := dhcpv4.FromBytes(buf[:n])
	if err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: failed to parse DHCPOFFER: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("DHCP lease offer for interface '%s':\n", iface)
	fmt.Printf("  IP Address: %s\n", offer.YourIPAddr)
	fmt.Printf("  Subnet Mask: %s\n", offer.SubnetMask())
	fmt.Printf("  Router: %v\n", offer.Router())
	fmt.Printf("  DNS: %v\n", offer.DNS())
	fmt.Printf("  Lease Time: %v\n", offer.IPAddressLeaseTime(0))
	fmt.Printf("  Server Identifier: %s\n", offer.ServerIdentifier())
}
