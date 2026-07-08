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

	"golang.org/x/term"
)

// susie: switch user in the current shell (like su) - no external dependencies
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

	// Get the user's shell from /etc/passwd
	shell := getShell(username)
	if shell == "" {
		shell = "/bin/highway"
	}

	// If not root, we need to authenticate via dosu first
	if os.Geteuid() != 0 {
		fmt.Print("Password: ")
		_, err := readPassword()
		if err != nil {
			fmt.Fprintln(os.Stderr, "susie: error reading password:", err)
			os.Exit(1)
		}
		fmt.Println()

		// Use dosu to escalate, then re-exec susie
		dosuPath, err := exec.LookPath("dosu")
		if err != nil {
			fmt.Fprintln(os.Stderr, "susie: dosu not found - need root privileges")
			os.Exit(1)
		}
		args := append([]string{"dosu", "susie", username}, os.Args[2:]...)
		if err := syscall.Exec(dosuPath, args, os.Environ()); err != nil {
			fmt.Fprintln(os.Stderr, "susie: failed to escalate:", err)
			os.Exit(1)
		}
		return
	}

	// We are root - set group and user ID
	if err := syscall.Setgid(gid); err != nil {
		fmt.Fprintln(os.Stderr, "susie: setgid failed:", err)
		os.Exit(1)
	}
	if err := syscall.Setuid(uid); err != nil {
		fmt.Fprintln(os.Stderr, "susie: setuid failed:", err)
		os.Exit(1)
	}

	// Change to home directory
	homeDir := u.HomeDir
	if homeDir != "" {
		os.Chdir(homeDir)
	}

	// Set HOME environment variable
	os.Setenv("HOME", homeDir)
	os.Setenv("USER", username)
	os.Setenv("LOGNAME", username)
	os.Setenv("SHELL", shell)

	// Start the user's shell
	cmd := exec.Command(shell)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "susie: failed to exec shell:", err)
		os.Exit(1)
	}
}

// getShell returns the shell for a user from /etc/passwd
func getShell(username string) string {
	f, err := os.Open("/etc/passwd")
	if err != nil {
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ":")
		if len(fields) >= 7 && fields[0] == username {
			return fields[6]
		}
	}
	return ""
}

// readPassword reads a line from stdin without echoing
func readPassword() (string, error) {
	fd := int(os.Stdin.Fd())
	pwBytes, err := term.ReadPassword(fd)
	fmt.Println()
	return strings.TrimSpace(string(pwBytes)), err
}
