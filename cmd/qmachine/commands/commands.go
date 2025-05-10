// Package commands provides command handlers for the QMachine REPL
package commands

import (
	"fmt"
	"strconv"
	"strings"

	"qmachine/help"
	"qmachine/quantum"
)

// Handler handles REPL command execution
type Handler struct {
	machine     *quantum.QuantumRISCVMachine
	hostMachine *quantum.HostQuantumMachine
	useHost     bool
}

// NewHandler creates a new command handler
func NewHandler(numQubits int) *Handler {
	return &Handler{
		machine:     quantum.NewQuantumRISCVMachine(numQubits),
		hostMachine: quantum.NewHostQuantumMachine(numQubits),
		useHost:     false,
	}
}

// ShowHelp displays all available commands and instructions
func (h *Handler) ShowHelp() {
	fmt.Println(help.GetBasicCommands())
	fmt.Println()
	fmt.Println(help.GetQuantumInstructions())
	fmt.Println()
	fmt.Println(help.GetRISCVInstructions())
}

// HandleGate processes quantum gate commands
func (h *Handler) HandleGate(args []string) error {
	if h.useHost {
		return fmt.Errorf("gate commands are exclusive to VM execution mode")
	}
	if len(args) < 2 {
		return fmt.Errorf("usage: gate <type> <target> [controls...]")
	}

	target, err := h.parseQubitIndex(args[1])
	if err != nil {
		return fmt.Errorf("invalid target qubit: %v", err)
	}

	controls, err := h.parseControlQubits(args[2:])
	if err != nil {
		return err
	}

	instruction, err := h.createGateInstruction(strings.ToUpper(args[0]), target, controls)
	if err != nil {
		return err
	}

	return h.machine.ExecuteRISCInstruction(fmt.Sprintf("qapply x%d, x%d, %d", instruction.Target, instruction.Controls[0], instruction.Opcode))
}

// HandleMeasure processes qubit measurement commands
func (h *Handler) HandleMeasure(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: measure <qubit>")
	}

	qubit, err := h.parseQubitIndex(args[0])
	if err != nil {
		return fmt.Errorf("invalid qubit index: %v", err)
	}

	// Convert uint8 to int for MeasureQubit
	result := h.machine.MeasureQubit(int(qubit))
	fmt.Printf("Measurement result: %d\n", result)
	return nil
}

// HandleState displays the current quantum state
func (h *Handler) HandleState() error {
	// Since GetQuantumState is not available, we'll show register state instead
	h.HandleRegisters()
	return nil
}

// HandleReset resets the quantum state
func (h *Handler) HandleReset() error {
	// Since Reset is not available, we'll recreate the machine
	h.machine = quantum.NewQuantumRISCVMachine(2000) // Using default size
	return nil
}

// HandleRISC processes RISC-V instructions
func (h *Handler) HandleRISC(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: riscv <instruction>")
	}

	instruction := strings.Join(args, " ")
	return h.machine.ExecuteRISCInstruction(instruction)
}

// HandleLoad loads a RISC-V program from a file
func (h *Handler) HandleLoad(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: load <file>")
	}

	return h.machine.LoadRISCProgram(args[0])
}

// HandleRun executes the loaded RISC-V program
func (h *Handler) HandleRun() error {
	return h.machine.ExecuteRISCProgram()
}

// HandleMode toggles between VM and host-native execution
func (h *Handler) HandleMode() {
	h.useHost = !h.useHost
	mode := "VM"
	if h.useHost {
		mode = "host-native"
	}
	fmt.Printf("Switched to %s execution mode\n", mode)
}

// HandleRegisters displays the current register state
func (h *Handler) HandleRegisters() {
	var registers [128]uint64
	if h.useHost {
		registers = h.hostMachine.GetRegisters()
	} else {
		registers = h.machine.GetRegisters()
	}

	fmt.Println("Register state:")
	for i, reg := range registers {
		fmt.Printf("  x%d: %d\n", i, reg)
	}
}

// Helper functions

func (h *Handler) parseQubitIndex(s string) (uint8, error) {
	index, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return uint8(index), nil
}

func (h *Handler) parseControlQubits(args []string) ([]uint8, error) {
	var controls []uint8
	for _, arg := range args {
		control, err := h.parseQubitIndex(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid control qubit: %v", err)
		}
		controls = append(controls, control)
	}
	return controls, nil
}

func (h *Handler) createGateInstruction(gateType string, target uint8, controls []uint8) (quantum.Instruction, error) {
	var opcode uint8
	switch gateType {
	case "X":
		opcode = 0x00
	case "Y":
		opcode = 0x01
	case "Z":
		opcode = 0x02
	case "H":
		opcode = 0x03
	case "S":
		opcode = 0x04
	case "T":
		opcode = 0x05
	case "CNOT":
		if len(controls) != 1 {
			return quantum.Instruction{}, fmt.Errorf("CNOT gate requires exactly one control qubit")
		}
		opcode = 0x06
	default:
		return quantum.Instruction{}, fmt.Errorf("unknown gate type: %s", gateType)
	}

	return quantum.Instruction{Opcode: opcode, Target: target, Controls: controls}, nil
}
