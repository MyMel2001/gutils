// Package repl provides the REPL (Read-Eval-Print Loop) for QMachine
package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"qmachine/commands"
)

// REPL represents the quantum computer REPL
type REPL struct {
	handler *commands.Handler
	reader  *bufio.Reader
}

// New creates a new REPL instance
func New(numQubits int) *REPL {
	return &REPL{
		handler: commands.NewHandler(numQubits),
		reader:  bufio.NewReader(os.Stdin),
	}
}

// Start begins the REPL session
func (r *REPL) Start() {
	fmt.Printf("QMachine Quantum Computer Simulator\n")
	fmt.Printf("Type 'help' for available commands\n\n")

	for {
		fmt.Print("\nqmachine> ")
		input, err := r.reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		// Trim whitespace and newlines
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Split input into command and arguments
		parts := strings.Fields(input)
		command := parts[0]
		args := parts[1:]

		if err := r.processCommand(command, args); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

// processCommand handles the execution of REPL commands
func (r *REPL) processCommand(command string, args []string) error {
	switch command {
	case "exit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "help":
		r.handler.ShowHelp()
	case "gate":
		return r.handler.HandleGate(args)
	case "measure":
		return r.handler.HandleMeasure(args)
	case "state":
		return r.handler.HandleState()
	case "reset":
		return r.handler.HandleReset()
	case "riscv":
		return r.handler.HandleRISC(args)
	case "load":
		return r.handler.HandleLoad(args)
	case "run":
		return r.handler.HandleRun()
	case "run-host":
		r.handler.HandleMode()
		return r.handler.HandleRun()
	case "mode":
		r.handler.HandleMode()
	case "registers":
		r.handler.HandleRegisters()
	default:
		return fmt.Errorf("unknown command. Type 'help' for available commands")
	}
	return nil
}
