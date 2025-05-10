package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// highway: a minimal interactive shell
func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("highway$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "exit" || line == "quit" {
			break
		}
		if line == "" {
			continue
		}
		args := strings.Fields(line)
		cmd := exec.Command("/bin/sh", "-c", line)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
	}
} 