package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// mkusr: create a new non-root user with a home directory and password
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

	if len(pw) < 4 {
		fmt.Fprintln(os.Stderr, "mkusr: password too short (minimum 4 characters)")
		os.Exit(1)
	}

	uid, gid := getNextUIDGID()
	home := "/home/" + username
	shell := "/bin/highway"

	// Generate salted SHA-512 hash (like Linux shadow but simpler)
	hash := hashPassword(pw)

	passwdEntry := fmt.Sprintf("%s:x:%d:%d::%s:%s\n", username, uid, gid, home, shell)
	shadowEntry := fmt.Sprintf("%s:%s:19000:0:99999:7:::\n", username, hash)
	groupEntry := fmt.Sprintf("%s:x:%d:\n", username, gid)

	// Ensure /etc directory exists
	os.MkdirAll("/etc", 0755)

	if err := appendToFile("/etc/passwd", passwdEntry); err != nil {
		fmt.Fprintln(os.Stderr, "mkusr: failed to update /etc/passwd:", err)
		os.Exit(1)
	}
	if err := appendToFile("/etc/shadow", shadowEntry); err != nil {
		fmt.Fprintln(os.Stderr, "mkusr: failed to update /etc/shadow:", err)
		os.Exit(1)
	}
	if err := appendToFile("/etc/group", groupEntry); err != nil {
		fmt.Fprintln(os.Stderr, "mkusr: failed to update /etc/group:", err)
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
	fmt.Printf("mkusr: user '%s' created successfully (UID=%d, GID=%d)\n", username, uid, gid)
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

// getNextUIDGID finds the next available UID and GID
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

// hashPassword hashes the password using salted SHA-512
// Format: $6$<salt>$<hash> (compatible with Linux crypt format)
func hashPassword(pw string) string {
	// Generate 16-byte salt
	salt := make([]byte, 16)
	rand.Read(salt)
	saltStr := base64.StdEncoding.EncodeToString(salt)

	// Compute SHA-512 hash of salt+password
	h := sha512.Sum512(append([]byte(saltStr), []byte(pw)...))
	hashStr := base64.StdEncoding.EncodeToString(h[:])

	return "$6$" + saltStr + "$" + hashStr
}

// appendToFile appends a line to a file
func appendToFile(path, line string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line)
	return err
}

// readPassword reads a password without echoing
func readPassword() (string, error) {
	fd := int(os.Stdin.Fd())
	pwBytes, err := term.ReadPassword(fd)
	return strings.TrimSpace(string(pwBytes)), err
}
