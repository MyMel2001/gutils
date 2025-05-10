package main

import (
	"fmt"
	"os"
	"os/exec"
)

// dosu: runs a command as root using sudo
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "dosu: usage: dosu COMMAND [ARGS...]")
		os.Exit(1)
	}
	cmd := exec.Command("sudo", os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "dosu: error running command:", err)
		os.Exit(1)
	}
} 