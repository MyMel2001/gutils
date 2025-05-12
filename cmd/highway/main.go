package main

import (
	"bufio"
	"fmt"
	"io"
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

// execLine executes a line as a command (supports |, &&, ;)
func execLine(line string) {
	if line == "" {
		return
	}
	// Handle ;
	seqCmds := splitByUnescaped(line, ';')
	for _, seqCmd := range seqCmds {
		seqCmd = strings.TrimSpace(seqCmd)
		if seqCmd == "" {
			continue
		}
		// Handle &&
		andCmds := splitByUnescaped(seqCmd, '&')
		if len(andCmds) > 1 {
			// Only treat as && if double ampersand
			var cmds []string
			for i := 0; i < len(andCmds); i++ {
				if i+1 < len(andCmds) && andCmds[i+1] == "" {
					cmds = append(cmds, strings.TrimSpace(andCmds[i]))
					i++ // skip next
				} else {
					cmds = append(cmds, strings.TrimSpace(andCmds[i]))
				}
			}
			execAndChain(cmds)
			continue
		}
		// Handle |
		pipeCmds := splitByUnescaped(seqCmd, '|')
		if len(pipeCmds) > 1 {
			execPipeline(pipeCmds)
			continue
		}
		// Single command
		execSingle(seqCmd)
	}
}

// splitByUnescaped splits a string by a delimiter, ignoring escaped delimiters
func splitByUnescaped(s string, delim rune) []string {
	var res []string
	var current strings.Builder
	inQuote := false
	quoteChar := byte(0)
	escape := false
	for _, c := range s {
		if escape {
			current.WriteRune(c)
			escape = false
			continue
		}
		if c == '\\' {
			escape = true
			continue
		}
		if inQuote {
			if byte(c) == quoteChar {
				inQuote = false
			} else {
				current.WriteRune(c)
			}
			continue
		}
		if c == '\'' || c == '"' {
			inQuote = true
			quoteChar = byte(c)
			continue
		}
		if c == delim && !inQuote {
			res = append(res, current.String())
			current.Reset()
			continue
		}
		current.WriteRune(c)
	}
	if current.Len() > 0 {
		res = append(res, current.String())
	}
	return res
}

// execAndChain executes commands chained with &&
func execAndChain(cmds []string) {
	for _, cmd := range cmds {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		args := splitArgs(cmd)
		if len(args) == 0 {
			continue
		}
		cmdPath, err := exec.LookPath(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "highway: command not found:", args[0])
			return
		}
		c := exec.Command(cmdPath, args[1:]...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return
		}
	}
}

// execPipeline executes a pipeline of commands separated by |
func execPipeline(cmds []string) {
	var commands []*exec.Cmd
	for _, cmdStr := range cmds {
		args := splitArgs(strings.TrimSpace(cmdStr))
		if len(args) == 0 {
			continue
		}
		cmdPath, err := exec.LookPath(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "highway: command not found:", args[0])
			return
		}
		commands = append(commands, exec.Command(cmdPath, args[1:]...))
	}
	if len(commands) == 0 {
		return
	}
	for i := 0; i < len(commands)-1; i++ {
		outPipe, inPipe := io.Pipe()
		commands[i].Stdout = inPipe
		commands[i+1].Stdin = outPipe
	}
	commands[0].Stdin = os.Stdin
	commands[len(commands)-1].Stdout = os.Stdout
	for _, cmd := range commands {
		cmd.Stderr = os.Stderr
	}
	// Start all but last
	for i := 0; i < len(commands)-1; i++ {
		if err := commands[i].Start(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return
		}
	}
	// Run last and wait for previous
	if err := commands[len(commands)-1].Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	for i := 0; i < len(commands)-1; i++ {
		commands[i].Wait()
	}
}

// execSingle executes a single command
func execSingle(cmdStr string) {
	args := splitArgs(cmdStr)
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
