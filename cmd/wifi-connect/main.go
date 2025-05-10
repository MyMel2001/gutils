package main

import (
	"fmt"
	"os"
)

// wifi-connect: lists WiFi interfaces and prints a message about connecting (no real connection)
func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "wifi-connect: usage: wifi-connect SSID PASSWORD")
		os.Exit(1)
	}
	ssid := os.Args[1]
	password := os.Args[2]
	// List likely WiFi interfaces (Linux: /sys/class/net/*/wireless)
	files, err := os.ReadDir("/sys/class/net")
	if err != nil {
		fmt.Fprintln(os.Stderr, "wifi-connect: cannot list interfaces:", err)
		os.Exit(1)
	}
	fmt.Println("Available WiFi interfaces:")
	any := false
	for _, f := range files {
		path := "/sys/class/net/" + f.Name() + "/wireless"
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("  %s\n", f.Name())
			any = true
		}
	}
	if !any {
		fmt.Println("  (none found)")
	}
	fmt.Printf("Would connect to SSID '%s' with password '%s' (not implemented)\n", ssid, password)
}
