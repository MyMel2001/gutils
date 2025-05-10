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
			os.Exit(0)
		}
		if line == "pwd" {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, "pwd:", err)
			} else {
				fmt.Println(dir)
			}
			continue
		}
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "cd ") || line == "cd" {
			args := strings.Fields(line)
			dir := ""
			if len(args) < 2 {
				dir = os.Getenv("HOME")
			} else {
				dir = args[1]
			}
			if err := os.Chdir(dir); err != nil {
				fmt.Fprintln(os.Stderr, "cd:", err)
			}
			continue
		}
		cmd := exec.Command("/bin/sh", "-c", line)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
	}
} 