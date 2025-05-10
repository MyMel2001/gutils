package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// ping: sends an ICMP echo request to a host. Usage: ping HOST
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "ping: usage: ping HOST")
		os.Exit(1)
	}
	host := os.Args[1]
	addr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ping: resolve error:", err)
		os.Exit(1)
	}
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ping: listen error:", err)
		os.Exit(1)
	}
	defer conn.Close()
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("SneedPing"),
		},
	}
	b, err := msg.Marshal(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ping: marshal error:", err)
		os.Exit(1)
	}
	start := time.Now()
	_, err = conn.WriteTo(b, &net.IPAddr{IP: addr.IP})
	if err != nil {
		fmt.Fprintln(os.Stderr, "ping: send error:", err)
		os.Exit(1)
	}
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	reply := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ping: timeout or error:", err)
		os.Exit(1)
	}
	dur := time.Since(start)
	parsed, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		fmt.Fprintln(os.Stderr, "ping: parse error:", err)
		os.Exit(1)
	}
	if parsed.Type == ipv4.ICMPTypeEchoReply {
		fmt.Printf("Reply from %s: time=%v\n", peer, dur)
	} else {
		fmt.Printf("Unexpected ICMP type: %v\n", parsed.Type)
	}
}
