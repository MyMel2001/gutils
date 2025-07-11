package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mdlayher/wifi"
)

// wifi-connect: connects to a WiFi network using WPA2-PSK
func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "wifi-connect: usage: wifi-connect SSID PASSWORD")
		os.Exit(1)
	}
	ssid := os.Args[1]
	password := os.Args[2]

	client, err := wifi.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "wifi-connect: cannot create wifi client:", err)
		os.Exit(1)
	}
	defer client.Close()

	ifaces, err := client.Interfaces()
	if err != nil || len(ifaces) == 0 {
		fmt.Fprintln(os.Stderr, "wifi-connect: no WiFi interfaces found:", err)
		os.Exit(1)
	}
	fmt.Println("Available WiFi interfaces:")
	for _, ifi := range ifaces {
		fmt.Printf("  %s\n", ifi.Name)
	}
	// Use the first WiFi interface
	ifi := ifaces[0]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Scan for networks
	if err := client.Scan(ctx, ifi); err != nil {
		fmt.Fprintln(os.Stderr, "wifi-connect: scan failed:", err)
		os.Exit(1)
	}
	aps, err := client.AccessPoints(ifi)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wifi-connect: cannot list access points:", err)
		os.Exit(1)
	}
	found := false
	for _, ap := range aps {
		if ap.SSID == ssid {
			found = true
			break
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "wifi-connect: SSID '%s' not found\n", ssid)
		os.Exit(1)
	}

	// Use github.com/mdlayher/wifi to connect
	fmt.Printf("Connecting to SSID '%s' on interface '%s'...\n", ssid, ifi.Name)
	if err := client.ConnectWPAPSK(ifi, ssid, password); err != nil {
		fmt.Fprintln(os.Stderr, "wifi-connect: failed to connect:", err)
		os.Exit(1)
	}

	bss, err := client.BSS(ifi)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wifi-connect: failed to get BSS information after connecting:", err)
		os.Exit(1)
	}

	if bss.Status == wifi.BSSStatusAssociated && bss.SSID == ssid {
		fmt.Println("Successfully connected to WiFi network.")
	} else {
		fmt.Fprintf(os.Stderr, "wifi-connect: connection to '%s' failed, status: %s\n", ssid, bss.Status)
		os.Exit(1)
	}
}
