package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// xargs: runs a command with arguments from stdin (one per line, simple version)
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "xargs: usage: xargs COMMAND")
		os.Exit(1)
	}
	cmdName := os.Args[1]
	args := os.Args[2:]
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		cmd := exec.Command(cmdName, append(args, line)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
