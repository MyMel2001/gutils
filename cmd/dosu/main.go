package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"golang.org/x/term"
)

// dosu: minimal su/sudo clone with password hash in /etc/dosu_passwd
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "dosu: usage: dosu COMMAND [ARGS...]")
		os.Exit(1)
	}

	// Read the stored hash from /etc/dosu_passwd
	hashBytes, err := os.ReadFile("/etc/dosu_passwd")
	if err != nil {
		fmt.Fprintln(os.Stderr, "dosu: cannot read /etc/dosu_passwd:", err)
		os.Exit(1)
	}
	storedHash := strings.TrimSpace(string(hashBytes))

	// Prompt for password
	fmt.Print("Password: ")
	pw, err := readPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, "dosu: error reading password:", err)
		os.Exit(1)
	}

	// Hash the entered password
	h := sha256.Sum256([]byte(pw))
	pwHash := hex.EncodeToString(h[:])

	// Compare hashes
	if pwHash != storedHash {
		fmt.Fprintln(os.Stderr, "dosu: authentication failed")
		os.Exit(1)
	}

	// Try to setuid(0) (must be setuid root)
	if err := syscall.Setuid(0); err != nil {
		fmt.Fprintln(os.Stderr, "dosu: failed to escalate privileges:", err)
		os.Exit(1)
	}

	// Execute the command as root
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "dosu: error running command:", err)
		os.Exit(1)
	}
}

// readPassword reads a line from stdin without echoing using golang.org/x/term
func readPassword() (string, error) {
	fd := int(os.Stdin.Fd())
	pwBytes, err := term.ReadPassword(fd)
	fmt.Println()
	return strings.TrimSpace(string(pwBytes)), err
} 