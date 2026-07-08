package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/vishvananda/netlink"
)

// dhcp-get: requests a DHCP lease for the specified interface and configures the IP
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "dhcp-get: usage: dhcp-get INTERFACE")
		os.Exit(1)
	}
	ifaceName := os.Args[1]

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: interface '%s' not found: %v\n", ifaceName, err)
		os.Exit(1)
	}

	// Create a UDP connection on port 68 (DHCP client port)
	laddr := &net.UDPAddr{IP: net.IPv4zero, Port: 68}
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: failed to listen on UDP port 68: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Build DHCP Discover packet for the interface
	discover, err := dhcpv4.NewDiscoveryForInterface(ifaceName)
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
	fmt.Println("Sent DHCPDISCOVER, waiting for OFFER...")

	// Set a timeout for receiving the offer
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
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

	fmt.Printf("Received DHCPOFFER:\n")
	fmt.Printf("  IP Address: %s\n", offer.YourIPAddr)
	fmt.Printf("  Subnet Mask: %s\n", offer.SubnetMask())
	fmt.Printf("  Router: %v\n", offer.Router())
	fmt.Printf("  DNS: %v\n", offer.DNS())
	fmt.Printf("  Lease Time: %v\n", offer.IPAddressLeaseTime(0))
	fmt.Printf("  Server Identifier: %s\n", offer.ServerIdentifier())

	// Build DHCP Request packet
	request, err := dhcpv4.NewRequestFromOffer(offer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: failed to build DHCPREQUEST: %v\n", err)
		os.Exit(1)
	}

	// Send DHCP Request
	if _, err := conn.WriteTo(request.ToBytes(), bcast); err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: failed to send DHCPREQUEST: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Sent DHCPREQUEST, waiting for ACK...")

	// Set a timeout for receiving the ACK
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	n, _, err = conn.ReadFrom(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: failed to receive DHCPACK: %v\n", err)
		os.Exit(1)
	}

	ack, err := dhcpv4.FromBytes(buf[:n])
	if err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: failed to parse DHCPACK: %v\n", err)
		os.Exit(1)
	}

	if ack.MessageType() != dhcpv4.MessageTypeAck {
		fmt.Fprintf(os.Stderr, "dhcp-get: expected DHCPACK but got %s\n", ack.MessageType())
		os.Exit(1)
	}

	fmt.Printf("Received DHCPACK!\n")
	fmt.Printf("  IP Address: %s\n", ack.YourIPAddr)
	fmt.Printf("  Subnet Mask: %s\n", ack.SubnetMask())
	fmt.Printf("  Router: %v\n", ack.Router())
	fmt.Printf("  DNS: %v\n", ack.DNS())
	fmt.Printf("  Lease Time: %v\n", ack.IPAddressLeaseTime(0))
	fmt.Printf("  Server Identifier: %s\n", ack.ServerIdentifier())

	// Configure the interface with the leased IP using netlink
	if err := configureInterface(iface, ack); err != nil {
		fmt.Fprintf(os.Stderr, "dhcp-get: failed to configure interface: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("DHCP lease obtained and interface configured successfully.")
}

// configureInterface sets the IP address, subnet mask, and default route on the interface
func configureInterface(iface *net.Interface, ack *dhcpv4.DHCPv4) error {
	ip := ack.YourIPAddr
	subnet := ack.SubnetMask()
	routers := ack.Router()

	// Get the link by name
	link, err := netlink.LinkByName(iface.Name)
	if err != nil {
		return fmt.Errorf("cannot find interface '%s': %w", iface.Name, err)
	}

	cidr := &net.IPNet{IP: ip, Mask: subnet}

	// Add the IP address to the interface
	addr := &netlink.Addr{
		IPNet: cidr,
		Label: iface.Name,
	}
	if err := netlink.AddrAdd(link, addr); err != nil {
		return fmt.Errorf("failed to add IP address %s: %w", cidr.String(), err)
	}
	fmt.Printf("Added IP address %s to %s\n", cidr.String(), iface.Name)

	// Set default route via router
	if len(routers) > 0 {
		route := &netlink.Route{
			Scope:     0, // RT_SCOPE_UNIVERSE
			LinkIndex: link.Attrs().Index,
			Gw:        routers[0],
			Dst:       nil,
		}
		if err := netlink.RouteAdd(route); err != nil {
			fmt.Fprintf(os.Stderr, "dhcp-get: warning: failed to set default route via %s: %v\n", routers[0], err)
		} else {
			fmt.Printf("Added default route via %s\n", routers[0])
		}
	}

	// Bring interface up
	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("failed to bring interface up: %w", err)
	}

	return nil
}
