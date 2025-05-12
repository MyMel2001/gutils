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

// aliasMap stores user-defined aliases
var aliasMap = make(map[string]string)

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
	// Handle alias built-in
	if strings.HasPrefix(line, "alias") {
		handleAlias(line)
		return true
	}
	return false
}

// handleAlias processes the alias command
func handleAlias(line string) {
	args := strings.Fields(line)
	if len(args) == 1 {
		// Print all aliases
		for k, v := range aliasMap {
			fmt.Printf("alias %s='%s'\n", k, v)
		}
		return
	}
	for _, arg := range args[1:] {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			name := parts[0]
			val := strings.Trim(parts[1], "'\"")
			aliasMap[name] = val
		} else {
			// Print single alias
			if val, ok := aliasMap[arg]; ok {
				fmt.Printf("alias %s='%s'\n", arg, val)
			} else {
				fmt.Printf("alias: %s: not found\n", arg)
			}
		}
	}
}

// expandAlias replaces the first word with its alias if defined
func expandAlias(line string) string {
	fields := splitArgs(line)
	if len(fields) == 0 {
		return line
	}
	if val, ok := aliasMap[fields[0]]; ok {
		// Replace first word with alias value, append rest
		return val + " " + strings.Join(fields[1:], " ")
	}
	return line
}

// execLine executes a line as a command (supports |, &&, ;)
func execLine(line string) {
	if line == "" {
		return
	}
	// Expand alias before processing
	line = expandAlias(line)
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
	if len(cmds) == 0 {
		return
	}

	var input []byte

	for i, cmdStr := range cmds {
		args := splitArgs(strings.TrimSpace(cmdStr))
		if len(args) == 0 {
			continue
		}
		cmdPath, err := exec.LookPath(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "highway: command not found:", args[0])
			return
		}
		cmd := exec.Command(cmdPath, args[1:]...)

		// For the first command, use os.Stdin. For others, use the previous output.
		if i == 0 {
			cmd.Stdin = os.Stdin
		} else {
			cmd.Stdin = strings.NewReader(string(input))
		}

		// For the last command, output to os.Stdout. Otherwise, capture output.
		if i == len(cmds)-1 {
			cmd.Stdout = os.Stdout
		} else {
			var outBuf strings.Builder
			cmd.Stdout = &outBuf
			err = cmd.Run()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				return
			}
			input = []byte(outBuf.String())
			continue
		}

		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return
		}
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
