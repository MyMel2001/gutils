package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/mdlayher/wifi"
)

// wifi-connect: connects to a WiFi network using nl80211 directly (no external tools)
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "wifi-connect: usage: wifi-connect SSID [PASSWORD]")
		fmt.Fprintln(os.Stderr, "  If PASSWORD is omitted, connects to an open network")
		os.Exit(1)
	}
	ssid := os.Args[1]
	password := ""
	if len(os.Args) > 2 {
		password = os.Args[2]
	}

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
		fmt.Printf("  %s (MAC: %s)\n", ifi.Name, ifi.HardwareAddr)
	}
	// Use the first WiFi interface
	ifi := ifaces[0]

	// Scan for networks
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
		fmt.Fprintf(os.Stderr, "wifi-connect: SSID '%s' not found in scan results\n", ssid)
		fmt.Println("Available networks:")
		for _, ap := range aps {
			fmt.Printf("  %s (freq: %d MHz)\n", ap.SSID, ap.Frequency)
		}
		os.Exit(1)
	}

	fmt.Printf("Connecting to SSID '%s' on interface '%s'...\n", ssid, ifi.Name)

	// Disconnect any existing connection first
	if err := client.Disconnect(ifi); err != nil {
		// Ignore error if not connected
	}

	if password != "" {
		// Connect using WPA2-PSK via nl80211 (kernel handles the 4-way handshake)
		if err := client.ConnectWPAPSK(ifi, ssid, password); err != nil {
			fmt.Fprintln(os.Stderr, "wifi-connect: WPA connection failed:", err)
			os.Exit(1)
		}
	} else {
		// Connect to open network
		if err := client.Connect(ifi, ssid); err != nil {
			fmt.Fprintln(os.Stderr, "wifi-connect: connection failed:", err)
			os.Exit(1)
		}
	}

	// Wait for association to complete
	time.Sleep(2 * time.Second)

	// Verify connection by checking BSS
	bss, err := client.BSS(ifi)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wifi-connect: connected but could not verify BSS:", err)
		// Still try DHCP
	} else if bss.SSID == ssid {
		fmt.Printf("Associated with %s (BSSID: %s)\n", bss.SSID, bss.BSSID)
	} else {
		fmt.Fprintf(os.Stderr, "wifi-connect: connected to '%s' instead of '%s'\n", bss.SSID, ssid)
	}

	// Get DHCP lease using our dhcp-get utility
	fmt.Println("Requesting DHCP lease...")
	dhcpCmd := exec.Command("dhcp-get", ifi.Name)
	dhcpCmd.Stdout = os.Stdout
	dhcpCmd.Stderr = os.Stderr
	if err := dhcpCmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "wifi-connect: DHCP failed:", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to WiFi network.")
}
