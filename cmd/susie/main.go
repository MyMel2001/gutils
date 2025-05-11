package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

// susie: switch user in the current shell (like su)
func main() {
	username := ""
	if len(os.Args) > 1 {
		username = os.Args[1]
	} else {
		fmt.Print("Username: ")
		reader := bufio.NewReader(os.Stdin)
		uname, _ := reader.ReadString('\n')
		username = strings.TrimSpace(uname)
	}
	if username == "" {
		fmt.Fprintln(os.Stderr, "susie: no username given")
		os.Exit(1)
	}

	u, err := user.Lookup(username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "susie: user '%s' not found\n", username)
		os.Exit(1)
	}
	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)

	fmt.Print("Password: ")
	_, err = readPassword()
	if err != nil {
		fmt.Fprintln(os.Stderr, "susie: error reading password:", err)
		os.Exit(1)
	}
	fmt.Println()

	// Try to use 'su' if not root
	if os.Geteuid() != 0 {
		cmd := exec.Command("su", username)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "susie: failed to switch user:", err)
			os.Exit(1)
		}
		return
	}

	// Set group and user ID
	if err := syscall.Setgid(gid); err != nil {
		fmt.Fprintln(os.Stderr, "susie: setgid failed:", err)
		os.Exit(1)
	}
	if err := syscall.Setuid(uid); err != nil {
		fmt.Fprintln(os.Stderr, "susie: setuid failed:", err)
		os.Exit(1)
	}

	shell := "/bin/sh"
	if uShell := u.Username; uShell != "" {
		// Try to get the shell from /etc/passwd
		if _, err := user.Lookup(username); err == nil {
			// Parse /etc/passwd for shell
			f, err := os.Open("/etc/passwd")
			if err == nil {
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					fields := strings.Split(scanner.Text(), ":")
					if len(fields) >= 7 && fields[0] == username {
						shell = fields[6]
						break
					}
				}
				f.Close()
			}
		}
	}
	cmd := exec.Command(shell)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "susie: failed to exec shell:", err)
		os.Exit(1)
	}
}

// readPassword reads a line from stdin without echoing
func readPassword() (string, error) {
	exec.Command("stty", "-echo").Run()
	reader := bufio.NewReader(os.Stdin)
	pw, err := reader.ReadString('\n')
	exec.Command("stty", "echo").Run()
	return strings.TrimSpace(pw), err
}
