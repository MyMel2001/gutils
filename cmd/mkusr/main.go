package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// mkusr: create a new non-root user with a home directory and password (minimal, home-made)
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "mkusr: usage: mkusr USERNAME")
		os.Exit(1)
	}
	username := os.Args[1]

	if userExists(username) {
		fmt.Fprintf(os.Stderr, "mkusr: user '%s' already exists\n", username)
		os.Exit(1)
	}

	fmt.Printf("Set password for %s: ", username)
	pw, err := readPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, "mkusr: error reading password:", err)
		os.Exit(1)
	}
	fmt.Print("\nConfirm password: ")
	pw2, err := readPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, "mkusr: error reading password:", err)
		os.Exit(1)
	}
	fmt.Println()
	if pw != pw2 {
		fmt.Fprintln(os.Stderr, "mkusr: passwords do not match")
		os.Exit(1)
	}

	uid, gid := getNextUIDGID()
	home := "/home/" + username
	shell := "/bin/sh"
	passwdEntry := fmt.Sprintf("%s:x:%d:%d::%s:%s\n", username, uid, gid, home, shell)
	shadowEntry := fmt.Sprintf("%s:%s:19000:0:99999:7:::\n", username, hashPassword(pw))

	if err := appendToFile("/etc/passwd", passwdEntry); err != nil {
		fmt.Fprintln(os.Stderr, "mkusr: failed to update /etc/passwd:", err)
		os.Exit(1)
	}
	if err := appendToFile("/etc/shadow", shadowEntry); err != nil {
		fmt.Fprintln(os.Stderr, "mkusr: failed to update /etc/shadow:", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(home, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "mkusr: failed to create home directory:", err)
		os.Exit(1)
	}
	if err := os.Chown(home, uid, gid); err != nil {
		fmt.Fprintln(os.Stderr, "mkusr: failed to set home directory owner:", err)
		os.Exit(1)
	}
	fmt.Printf("mkusr: user '%s' created successfully\n", username)
}

// userExists checks if a user exists in /etc/passwd
func userExists(username string) bool {
	f, err := os.Open("/etc/passwd")
	if err != nil {
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), username+":") {
			return true
		}
	}
	return false
}

// getNextUIDGID finds the next available UID and GID (minimal, not robust)
func getNextUIDGID() (int, int) {
	maxUID := 1000
	f, err := os.Open("/etc/passwd")
	if err != nil {
		return 1001, 1001
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ":")
		if len(fields) > 2 {
			uid, _ := strconv.Atoi(fields[2])
			if uid > maxUID {
				maxUID = uid
			}
		}
	}
	return maxUID + 1, maxUID + 1
}

// hashPassword hashes the password using SHA-512 and base64 (not crypt-compatible, minimal)
func hashPassword(pw string) string {
	h := sha512.Sum512([]byte(pw))
	return base64.StdEncoding.EncodeToString(h[:])
}

// appendToFile appends a line to a file
func appendToFile(path, line string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line)
	return err
}

// readPassword reads a line from stdin (echoed, minimal)
func readPassword() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	pw, err := reader.ReadString('\n')
	return strings.TrimSpace(pw), err
}
