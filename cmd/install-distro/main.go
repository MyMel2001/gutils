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
	if len(devices) == 0 {
		fmt.Fprintln(os.Stderr, "install-distro: no block devices found")
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
	fmt.Printf("WARNING: This will DESTROY all data on %s!\n", target)
	fmt.Print("Are you sure? (yes/no): ")
	scanner.Scan()
	confirm := strings.TrimSpace(scanner.Text())
	if confirm != "yes" {
		fmt.Println("Installation cancelled.")
		os.Exit(0)
	}

	fmt.Printf("Writing root filesystem to %s...\n", target)
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
			devName := fields[3]
			// Only show whole devices (no partition numbers), or show all
			// For safety, show all non-loop, non-ram devices
			devs = append(devs, "/dev/"+devName)
		}
	}
	return devs, scanner.Err()
}

// ddRootToDevice copies the root filesystem to the target device
func ddRootToDevice(target string) error {
	// Detect the root device by checking /proc/mounts
	root := detectRootDevice()
	if root == "" {
		return fmt.Errorf("could not detect root device")
	}
	fmt.Printf("Source root device: %s\n", root)
	fmt.Printf("Target device: %s\n", target)

	if root == target {
		return fmt.Errorf("source and target are the same device")
	}

	in, err := os.Open(root)
	if err != nil {
		return fmt.Errorf("cannot open source %s: %w", root, err)
	}
	defer in.Close()

	out, err := os.OpenFile(target, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("cannot open target %s: %w", target, err)
	}
	defer out.Close()

	buf := make([]byte, 1024*1024) // 1MB buffer
	totalWritten := int64(0)
	for {
		n, err := in.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				return fmt.Errorf("write error at offset %d: %w", totalWritten, werr)
			}
			totalWritten += int64(n)
		}
		if err != nil {
			break
		}
	}
	fmt.Printf("Written %d bytes to %s\n", totalWritten, target)
	return nil
}

// detectRootDevice finds the root device from /proc/mounts
func detectRootDevice() string {
	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "/" {
			dev := fields[0]
			// Resolve /dev/root or /dev/mmcblkXpY etc
			if dev == "rootfs" || dev == "/dev/root" {
				// Try to find the real device
				continue
			}
			return dev
		}
	}
	return ""
}

// expandFS calls expand-fs utility logic on the target device
func expandFS(target string) error {
	// Call our expand-fs utility
	return execCommand("expand-fs", target)
}

// execCommand runs a command using our own utilities
func execCommand(name string, args ...string) error {
	// Look for the command in PATH or current directory
	paths := []string{
		"./bin/" + name,
		"/bin/" + name,
		"/usr/bin/" + name,
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			// Use syscall.Exec to run it
			proc, err := os.StartProcess(p, append([]string{name}, args...), &os.ProcAttr{
				Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			})
			if err != nil {
				return err
			}
			state, err := proc.Wait()
			if err != nil {
				return err
			}
			if !state.Success() {
				return fmt.Errorf("%s exited with status %d", name, state.ExitCode())
			}
			return nil
		}
	}
	return fmt.Errorf("%s: command not found", name)
}
