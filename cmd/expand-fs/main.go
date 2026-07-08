package main

import (
	"fmt"
	"os"
	"os/exec"
)

// expand-fs: expands the root filesystem on the given device
// Uses resize2fs or equivalent to expand the filesystem to fill the device
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "expand-fs: usage: expand-fs DEVICE")
		os.Exit(1)
	}
	device := os.Args[1]

	// Check if the device exists
	if _, err := os.Stat(device); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "expand-fs: device '%s' does not exist\n", device)
		os.Exit(1)
	}

	// Detect filesystem type
	fsType := detectFSType(device)
	fmt.Printf("expand-fs: detected filesystem type '%s' on %s\n", fsType, device)

	switch fsType {
	case "ext2", "ext3", "ext4":
		resizeExt(device)
	case "btrfs":
		resizeBtrfs(device)
	case "xfs":
		resizeXFS(device)
	case "f2fs":
		resizeF2FS(device)
	default:
		// Try resize2fs as a fallback
		fmt.Printf("expand-fs: unknown filesystem type '%s', trying resize2fs...\n", fsType)
		resizeExt(device)
	}
}

// detectFSType reads the superblock to determine the filesystem type
func detectFSType(device string) string {
	f, err := os.Open(device)
	if err != nil {
		return "unknown"
	}
	defer f.Close()

	// Read the superblock at offset 1024 (standard location)
	sb := make([]byte, 1024)
	if _, err := f.ReadAt(sb, 1024); err != nil {
		return "unknown"
	}

	// Check for ext2/3/4 magic (0xEF53) at offset 0x38
	if len(sb) > 0x3A && sb[0x38] == 0x53 && sb[0x39] == 0xEF {
		// Check revision level and feature flags to determine ext version
		return "ext4"
	}

	// Check for XFS magic at offset 0
	f.Seek(0, 0)
	magic := make([]byte, 4)
	f.Read(magic)
	if string(magic) == "XFSB" {
		return "xfs"
	}

	// Check for Btrfs magic at offset 0x40
	f.ReadAt(magic, 0x40)
	if string(magic) == "_BHR" {
		return "btrfs"
	}

	// Check for F2FS magic at offset 0
	f.Seek(0, 0)
	f2fsMagic := make([]byte, 4)
	f.Read(f2fsMagic)
	if string(f2fsMagic) == "\x10\x20\xF5\xF2" {
		return "f2fs"
	}

	return "unknown"
}

// resizeExt resizes an ext2/3/4 filesystem
func resizeExt(device string) {
	fmt.Printf("Resizing ext filesystem on %s...\n", device)
	// Try e2fsck first, then resize2fs
	exec.Command("e2fsck", "-f", "-y", device).Run()
	cmd := exec.Command("resize2fs", device)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "expand-fs: resize2fs failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("expand-fs: filesystem expanded successfully")
}

// resizeBtrfs resizes a btrfs filesystem
func resizeBtrfs(device string) {
	fmt.Printf("Resizing btrfs filesystem on %s...\n", device)
	cmd := exec.Command("btrfs", "filesystem", "resize", "max", device)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "expand-fs: btrfs resize failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("expand-fs: filesystem expanded successfully")
}

// resizeXFS resizes an XFS filesystem
func resizeXFS(device string) {
	fmt.Printf("Resizing XFS filesystem on %s...\n", device)
	cmd := exec.Command("xfs_growfs", device)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "expand-fs: xfs_growfs failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("expand-fs: filesystem expanded successfully")
}

// resizeF2FS resizes an F2FS filesystem
func resizeF2FS(device string) {
	fmt.Printf("Resizing F2FS filesystem on %s...\n", device)
	cmd := exec.Command("resize.f2fs", device)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "expand-fs: resize.f2fs failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("expand-fs: filesystem expanded successfully")
}
