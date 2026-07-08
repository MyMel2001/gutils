package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/term"
)

const (
	saltLen    = 16
	hashLen    = 32
	passwdFile = "/etc/dosu_passwd"
)

// dosu: minimal su/sudo clone with salted SHA-256 password hash in /etc/dosu_passwd
// Format of /etc/dosu_passwd: hex(salt) + ":" + hex(hash)
// where hash = SHA-256(salt + password)
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "dosu: usage: dosu COMMAND [ARGS...]")
		os.Exit(1)
	}

	// Check if we're being called to set a password
	if os.Args[1] == "--set-passwd" || os.Args[1] == "--passwd" {
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "dosu: usage: dosu --set-passwd USERNAME")
			os.Exit(1)
		}
		setPassword(os.Args[2])
		return
	}

	// Read the stored hash from /etc/dosu_passwd
	hashBytes, err := os.ReadFile(passwdFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dosu: cannot read", passwdFile+":", err)
		fmt.Fprintln(os.Stderr, "dosu: run 'dosu --set-passwd root' to set the password")
		os.Exit(1)
	}
	storedData := strings.TrimSpace(string(hashBytes))

	// Parse salt:hash format
	parts := strings.SplitN(storedData, ":", 2)
	var storedSalt, storedHash string
	if len(parts) == 2 {
		storedSalt = parts[0]
		storedHash = parts[1]
	} else {
		// Legacy format: just a hash (no salt) - still support it
		storedHash = parts[0]
	}

	// Prompt for password
	fmt.Print("Password: ")
	pw, err := readPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, "dosu: error reading password:", err)
		os.Exit(1)
	}

	// Verify password
	var pwHash string
	if storedSalt != "" {
		// Salted hash
		saltBytes, err := hex.DecodeString(storedSalt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dosu: invalid salt in", passwdFile)
			os.Exit(1)
		}
		h := sha256.Sum256(append(saltBytes, []byte(pw)...))
		pwHash = hex.EncodeToString(h[:])
	} else {
		// Legacy unsalted hash
		h := sha256.Sum256([]byte(pw))
		pwHash = hex.EncodeToString(h[:])
	}

	if pwHash != storedHash {
		fmt.Fprintln(os.Stderr, "dosu: authentication failed")
		os.Exit(1)
	}

	// Authentication successful - escalate privileges
	// First set the group to root (GID 0)
	if err := syscall.Setgid(0); err != nil {
		fmt.Fprintln(os.Stderr, "dosu: failed to set group to root:", err)
		os.Exit(1)
	}

	// Set supplementary groups to root's groups
	if err := syscall.Setgroups([]int{0}); err != nil {
		fmt.Fprintln(os.Stderr, "dosu: failed to set supplementary groups:", err)
		os.Exit(1)
	}

	// Set uid to root
	if err := syscall.Setuid(0); err != nil {
		fmt.Fprintln(os.Stderr, "dosu: failed to escalate privileges:", err)
		os.Exit(1)
	}

	// Execute the command as root
	cmdPath, err := exec.LookPath(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "dosu: command not found:", os.Args[1])
		os.Exit(1)
	}
	args := os.Args[1:]
	env := os.Environ()
	if err := syscall.Exec(cmdPath, args, env); err != nil {
		fmt.Fprintln(os.Stderr, "dosu: error running command:", err)
		os.Exit(1)
	}
}

// setPassword sets the password for a user by writing the salted hash to /etc/dosu_passwd
func setPassword(username string) {
	fmt.Printf("Setting password for %s\n", username)
	fmt.Print("New password: ")
	pw1, err := readPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, "dosu: error reading password:", err)
		os.Exit(1)
	}
	fmt.Println()

	fmt.Print("Confirm password: ")
	pw2, err := readPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, "dosu: error reading password:", err)
		os.Exit(1)
	}
	fmt.Println()

	if pw1 != pw2 {
		fmt.Fprintln(os.Stderr, "dosu: passwords do not match")
		os.Exit(1)
	}

	if len(pw1) < 4 {
		fmt.Fprintln(os.Stderr, "dosu: password too short (minimum 4 characters)")
		os.Exit(1)
	}

	// Generate random salt
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		fmt.Fprintln(os.Stderr, "dosu: failed to generate salt:", err)
		os.Exit(1)
	}

	// Compute salted hash
	h := sha256.Sum256(append(salt, []byte(pw1)...))
	saltHex := hex.EncodeToString(salt)
	hashHex := hex.EncodeToString(h[:])

	// Write to password file
	content := saltHex + ":" + hashHex + "\n"
	if err := os.WriteFile(passwdFile, []byte(content), 0600); err != nil {
		fmt.Fprintln(os.Stderr, "dosu: failed to write", passwdFile+":", err)
		os.Exit(1)
	}

	fmt.Printf("Password set for %s\n", username)
}

// readPassword reads a line from stdin without echoing using golang.org/x/term
func readPassword() (string, error) {
	fd := int(os.Stdin.Fd())
	pwBytes, err := term.ReadPassword(fd)
	fmt.Println()
	return strings.TrimSpace(string(pwBytes)), err
}
