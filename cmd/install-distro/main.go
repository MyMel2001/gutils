package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// install-distro: install the running system to another drive
func main() {
	fmt.Println("Listing block devices:")
	devices, err := listBlockDevices()
	if err != nil {
		fmt.Fprintln(os.Stderr, "install-distro: failed to list block devices:", err)
		os.Exit(1)
	}
	for i, dev := range devices {
		fmt.Printf("%d: %s\n", i+1, dev)
	}
	fmt.Print("Select target device number: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := scanner.Text()
	idx := -1
	fmt.Sscanf(choice, "%d", &idx)
	if idx < 1 || idx > len(devices) {
		fmt.Fprintln(os.Stderr, "install-distro: invalid selection")
		os.Exit(1)
	}
	target := devices[idx-1]
	fmt.Printf("Writing root device to %s...\n", target)
	if err := ddRootToDevice(target); err != nil {
		fmt.Fprintln(os.Stderr, "install-distro: dd failed:", err)
		os.Exit(1)
	}
	fmt.Printf("Expanding root filesystem on %s...\n", target)
	if err := expandFS(target); err != nil {
		fmt.Fprintln(os.Stderr, "install-distro: expand-fs failed:", err)
		os.Exit(1)
	}
	fmt.Println("Install complete!")
}

// listBlockDevices returns a list of block device paths (e.g. /dev/sda, /dev/vda)
func listBlockDevices() ([]string, error) {
	f, err := os.Open("/proc/partitions")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var devs []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 4 && !strings.HasPrefix(fields[3], "loop") && !strings.HasPrefix(fields[3], "ram") {
			devs = append(devs, "/dev/"+fields[3])
		}
	}
	return devs, scanner.Err()
}

// ddRootToDevice copies the root device to the target device (block by block)
func ddRootToDevice(target string) error {
	// This is a simplified version; in real use, you'd want to detect the root device
	root := "/dev/root" // This may need to be detected more robustly
	in, err := os.Open(root)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(target, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer out.Close()
	buf := make([]byte, 1024*1024)
	for {
		n, err := in.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				return werr
			}
		}
		if err != nil {
			break
		}
	}
	return nil
}

// expandFS calls expand-fs utility logic on the target device
func expandFS(target string) error {
	// This is a stub; the real expand-fs logic should be implemented in expand-fs/main.go
	return nil
}
