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
	if len(os.Args) > 1 {
		scriptFile := os.Args[1]
		f, err := os.Open(scriptFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "highway: cannot open script:", err)
			os.Exit(1)
		}
		defer f.Close()
		execScript(f)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("highway$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}
		line = strings.TrimSpace(line)
		if handleBuiltins(line) {
			continue
		}
		execLine(line)
	}
}

// execScript reads and executes each line of a script file
func execScript(f *os.File) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if handleBuiltins(line) {
			continue
		}
		execLine(line)
	}
}

// handleBuiltins processes built-in commands, returns true if handled
func handleBuiltins(line string) bool {
	if line == "exit" || line == "quit" {
		os.Exit(0)
	}
	if line == "clear" {
		// Clear the terminal screen (ANSI escape code)
		fmt.Print("\033[2J\033[H")
		return true
	}
	if line == "pwd" {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "pwd:", err)
		} else {
			fmt.Println(dir)
		}
		return true
	}
	if line == "" {
		return true
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
		return true
	}
	return false
}

// execLine executes a line as a command (no external shell)
func execLine(line string) {
	if line == "" {
		return
	}
	args := splitArgs(line)
	if len(args) == 0 {
		return
	}
	cmdPath, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "highway: command not found:", args[0])
		return
	}
	cmd := exec.Command(cmdPath, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

// splitArgs splits a command line into arguments (basic, handles quotes)
func splitArgs(line string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := byte(0)
	escape := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		if escape {
			current.WriteByte(c)
			escape = false
			continue
		}
		if c == '\\' {
			escape = true
			continue
		}
		if inQuote {
			if c == quoteChar {
				inQuote = false
			} else {
				current.WriteByte(c)
			}
			continue
		}
		if c == '"' || c == '\'' {
			inQuote = true
			quoteChar = c
			continue
		}
		if c == ' ' || c == '\t' {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}
		current.WriteByte(c)
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
