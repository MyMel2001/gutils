package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"qmachine/quantum"
)

// REPL represents the quantum computer REPL
type REPL struct {
	machine     *quantum.QuantumRISCVMachine
	hostMachine *quantum.HostQuantumMachine
	reader      *bufio.Reader
	useHost     bool
}

// NewREPL creates a new REPL instance
func NewREPL(numQubits int) *REPL {
	return &REPL{
		machine:     quantum.NewQuantumRISCVMachine(numQubits),
		hostMachine: quantum.NewHostQuantumMachine(numQubits),
		reader:      bufio.NewReader(os.Stdin),
		useHost:     false,
	}
}

// Start begins the REPL session
func (r *REPL) Start() {
	fmt.Println("Welcome to QMachine REPL")
	fmt.Println("Available commands:")
	fmt.Println("  gate <type> <target> [controls...] - Apply a quantum gate")
	fmt.Println("  measure <qubit>                    - Measure a qubit")
	fmt.Println("  state                              - Show current quantum state")
	fmt.Println("  reset                              - Reset quantum state")
	fmt.Println("  riscv <instruction>                - Execute RISC-V instruction")
	fmt.Println("  load <file>                        - Load RISC-V program from file")
	fmt.Println("  run                                - Run loaded RISC-V program")
	fmt.Println("  run-host                           - Run loaded program using host-native execution")
	fmt.Println("  mode                               - Toggle between VM and host-native execution")
	fmt.Println("  registers                          - Show RISC-V registers")
	fmt.Println("  help                               - Show this help message")
	fmt.Println("  exit                               - Exit REPL")
	fmt.Println("\nAvailable gates: X, Y, Z, H, S, T, CNOT")
	fmt.Println("\nCustom Quantum RISC-V Instructions (Q-RISC-V Extensions):")
	fmt.Println("  qinit rd                          - Initialize quantum register with |0⟩")
	fmt.Println("  qapply rd, rs1, imm              - Apply quantum gate (imm: 0=X, 1=Y, 2=Z, 3=H, 4=S, 5=T, 6=CNOT)")
	fmt.Println("  qmeasure rd, rs1                 - Measure quantum register")
	fmt.Println("  qentangle rd, rs1, rs2          - Entangle two quantum registers")
	fmt.Println("\nStandard RISC-V Instructions:")
	fmt.Println("  add rd, rs1, rs2    - Add registers")
	fmt.Println("  sub rd, rs1, rs2    - Subtract registers")
	fmt.Println("  and rd, rs1, rs2    - AND registers")
	fmt.Println("  or rd, rs1, rs2     - OR registers")
	fmt.Println("  xor rd, rs1, rs2    - XOR registers")
	fmt.Println("  sll rd, rs1, rs2    - Shift left logical")
	fmt.Println("  srl rd, rs1, rs2    - Shift right logical")
	fmt.Println("  sra rd, rs1, rs2    - Shift right arithmetic")
	fmt.Println("  slt rd, rs1, rs2    - Set if less than")
	fmt.Println("  sltu rd, rs1, rs2   - Set if less than unsigned")
	fmt.Println("  addi rd, rs1, imm   - Add immediate")
	fmt.Println("  slli rd, rs1, imm   - Shift left logical immediate")
	fmt.Println("  srli rd, rs1, imm   - Shift right logical immediate")
	fmt.Println("  srai rd, rs1, imm   - Shift right arithmetic immediate")
	fmt.Println("  andi rd, rs1, imm   - AND immediate")
	fmt.Println("  ori rd, rs1, imm    - OR immediate")
	fmt.Println("  xori rd, rs1, imm   - XOR immediate")
	fmt.Println("  slti rd, rs1, imm   - Set if less than immediate")
	fmt.Println("  sltiu rd, rs1, imm  - Set if less than immediate unsigned")
	fmt.Println("  lui rd, imm         - Load upper immediate")
	fmt.Println("  auipc rd, imm       - Add upper immediate to PC")
	fmt.Println("  jal rd, offset      - Jump and link")
	fmt.Println("  jalr rd, rs1, offset - Jump and link register")
	fmt.Println("  beq rs1, rs2, offset - Branch if equal")
	fmt.Println("  bne rs1, rs2, offset - Branch if not equal")
	fmt.Println("  blt rs1, rs2, offset - Branch if less than")
	fmt.Println("  bge rs1, rs2, offset - Branch if greater or equal")
	fmt.Println("  bltu rs1, rs2, offset - Branch if less than unsigned")
	fmt.Println("  bgeu rs1, rs2, offset - Branch if greater or equal unsigned")
	fmt.Println("  lw rd, offset(rs1)   - Load word")
	fmt.Println("  lh rd, offset(rs1)   - Load halfword")
	fmt.Println("  lb rd, offset(rs1)   - Load byte")
	fmt.Println("  lwu rd, offset(rs1)  - Load word unsigned")
	fmt.Println("  lhu rd, offset(rs1)  - Load halfword unsigned")
	fmt.Println("  lbu rd, offset(rs1)  - Load byte unsigned")
	fmt.Println("  sw rs2, offset(rs1)  - Store word")
	fmt.Println("  sh rs2, offset(rs1)  - Store halfword")
	fmt.Println("  sb rs2, offset(rs1)  - Store byte")

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

		switch command {
		case "exit":
			fmt.Println("Goodbye!")
			return
		case "help":
			r.showHelp()
		case "gate":
			r.handleGateCommand(args)
		case "measure":
			r.handleMeasureCommand(args)
		case "state":
			r.handleStateCommand()
		case "reset":
			r.handleResetCommand()
		case "riscv":
			r.handleRISCCommand(args)
		case "load":
			r.handleLoadCommand(args)
		case "run":
			r.handleRunCommand()
		case "run-host":
			r.useHost = true
			r.handleRunCommand()
		case "mode":
			r.handleModeCommand()
		case "registers":
			r.handleRegistersCommand()
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func (r *REPL) showHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  gate <type> <target> [controls...] - Apply a quantum gate")
	fmt.Println("  measure <qubit>                    - Measure a qubit")
	fmt.Println("  state                              - Show current quantum state")
	fmt.Println("  reset                              - Reset quantum state")
	fmt.Println("  riscv <instruction>                - Execute RISC-V instruction")
	fmt.Println("  load <file>                        - Load RISC-V program from file")
	fmt.Println("  run                                - Run loaded RISC-V program")
	fmt.Println("  run-host                           - Run loaded program using host-native execution")
	fmt.Println("  mode                               - Toggle between VM and host-native execution")
	fmt.Println("  registers                          - Show RISC-V registers")
	fmt.Println("  help                               - Show this help message")
	fmt.Println("  exit                               - Exit REPL")
	fmt.Println("\nAvailable gates: X, Y, Z, H, S, T, CNOT")
	fmt.Println("\nCustom Quantum RISC-V Instructions (Q-RISC-V Extensions):")
	fmt.Println("  qinit rd                          - Initialize quantum register with |0⟩")
	fmt.Println("  qapply rd, rs1, imm              - Apply quantum gate (imm: 0=X, 1=Y, 2=Z, 3=H, 4=S, 5=T, 6=CNOT)")
	fmt.Println("  qmeasure rd, rs1                 - Measure quantum register")
	fmt.Println("  qentangle rd, rs1, rs2          - Entangle two quantum registers")
	fmt.Println("\nStandard RISC-V Instructions:")
	fmt.Println("  add rd, rs1, rs2    - Add registers")
	fmt.Println("  sub rd, rs1, rs2    - Subtract registers")
	fmt.Println("  and rd, rs1, rs2    - AND registers")
	fmt.Println("  or rd, rs1, rs2     - OR registers")
	fmt.Println("  xor rd, rs1, rs2    - XOR registers")
	fmt.Println("  sll rd, rs1, rs2    - Shift left logical")
	fmt.Println("  srl rd, rs1, rs2    - Shift right logical")
	fmt.Println("  sra rd, rs1, rs2    - Shift right arithmetic")
	fmt.Println("  slt rd, rs1, rs2    - Set if less than")
	fmt.Println("  sltu rd, rs1, rs2   - Set if less than unsigned")
	fmt.Println("  addi rd, rs1, imm   - Add immediate")
	fmt.Println("  slli rd, rs1, imm   - Shift left logical immediate")
	fmt.Println("  srli rd, rs1, imm   - Shift right logical immediate")
	fmt.Println("  srai rd, rs1, imm   - Shift right arithmetic immediate")
	fmt.Println("  andi rd, rs1, imm   - AND immediate")
	fmt.Println("  ori rd, rs1, imm    - OR immediate")
	fmt.Println("  xori rd, rs1, imm   - XOR immediate")
	fmt.Println("  slti rd, rs1, imm   - Set if less than immediate")
	fmt.Println("  sltiu rd, rs1, imm  - Set if less than immediate unsigned")
	fmt.Println("  lui rd, imm         - Load upper immediate")
	fmt.Println("  auipc rd, imm       - Add upper immediate to PC")
	fmt.Println("  jal rd, offset      - Jump and link")
	fmt.Println("  jalr rd, rs1, offset - Jump and link register")
	fmt.Println("  beq rs1, rs2, offset - Branch if equal")
	fmt.Println("  bne rs1, rs2, offset - Branch if not equal")
	fmt.Println("  blt rs1, rs2, offset - Branch if less than")
	fmt.Println("  bge rs1, rs2, offset - Branch if greater or equal")
	fmt.Println("  bltu rs1, rs2, offset - Branch if less than unsigned")
	fmt.Println("  bgeu rs1, rs2, offset - Branch if greater or equal unsigned")
	fmt.Println("  lw rd, offset(rs1)   - Load word")
	fmt.Println("  lh rd, offset(rs1)   - Load halfword")
	fmt.Println("  lb rd, offset(rs1)   - Load byte")
	fmt.Println("  lwu rd, offset(rs1)  - Load word unsigned")
	fmt.Println("  lhu rd, offset(rs1)  - Load halfword unsigned")
	fmt.Println("  lbu rd, offset(rs1)  - Load byte unsigned")
	fmt.Println("  sw rs2, offset(rs1)  - Store word")
	fmt.Println("  sh rs2, offset(rs1)  - Store halfword")
	fmt.Println("  sb rs2, offset(rs1)  - Store byte")
}

func (r *REPL) handleGateCommand(args []string) {
	if r.useHost {
		fmt.Println("Error: Gate commands are exclusive to VM execution mode.")
		return
	}
	if len(args) < 2 {
		fmt.Println("Usage: gate <type> <target> [controls...]")
		return
	}

	gateType := strings.ToUpper(args[0])
	target, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("Invalid target qubit: %v\n", err)
		return
	}

	var controls []uint8
	for _, arg := range args[2:] {
		control, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("Invalid control qubit: %v\n", err)
			return
		}
		controls = append(controls, uint8(control))
	}

	var instruction quantum.Instruction
	switch gateType {
	case "X":
		instruction = quantum.Instruction{Opcode: 0x00, Target: uint8(target), Controls: controls}
	case "Y":
		instruction = quantum.Instruction{Opcode: 0x01, Target: uint8(target), Controls: controls}
	case "Z":
		instruction = quantum.Instruction{Opcode: 0x02, Target: uint8(target), Controls: controls}
	case "H":
		instruction = quantum.Instruction{Opcode: 0x03, Target: uint8(target), Controls: controls}
	case "S":
		instruction = quantum.Instruction{Opcode: 0x04, Target: uint8(target), Controls: controls}
	case "T":
		instruction = quantum.Instruction{Opcode: 0x05, Target: uint8(target), Controls: controls}
	case "CNOT":
		if len(controls) != 1 {
			fmt.Println("CNOT gate requires exactly one control qubit")
			return
		}
		instruction = quantum.Instruction{Opcode: 0x06, Target: uint8(target), Controls: controls}
	default:
		fmt.Println("Unknown gate type. Available gates: X, Y, Z, H, S, T, CNOT")
		return
	}

	if err := r.machine.ExecuteInstruction(instruction); err != nil {
		fmt.Printf("Error applying gate: %v\n", err)
		return
	}

	fmt.Printf("Applied %s gate to qubit %d\n", gateType, target)
	if len(controls) > 0 {
		fmt.Printf("Control qubits: %v\n", controls)
	}
}

func (r *REPL) handleMeasureCommand(args []string) {
	if r.useHost {
		fmt.Println("Error: Measure commands are exclusive to VM execution mode.")
		return
	}
	if len(args) != 1 {
		fmt.Println("Usage: measure <qubit>")
		return
	}

	target, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Invalid qubit number: %v\n", err)
		return
	}

	if err := r.machine.MeasureQubit(target); err != nil {
		fmt.Printf("Error measuring qubit: %v\n", err)
		return
	}

	fmt.Printf("Measured qubit %d\n", target)
}

func (r *REPL) handleStateCommand() {
	if r.useHost {
		fmt.Println("Error: State commands are exclusive to VM execution mode.")
		return
	}
	state := r.machine.GetState()
	fmt.Printf("Current quantum state:\n")
	fmt.Printf("Number of qubits: %d\n", state.NumQubits())
	fmt.Printf("State vector size: %d\n", 1<<state.NumQubits())
	// Note: In a real implementation, you might want to show the actual state vector
	// but for a 2000-qubit system, this would be impractical to display
}

func (r *REPL) handleResetCommand() {
	if r.useHost {
		fmt.Println("Error: Reset commands are exclusive to VM execution mode.")
		return
	}
	r.machine = quantum.NewQuantumRISCVMachine(r.machine.GetState().NumQubits())
	fmt.Println("Quantum state reset to |0⟩^⊗n")
}

func (r *REPL) handleRISCCommand(args []string) {
	if r.useHost {
		fmt.Println("Error: Direct RISC-V command execution is exclusive to VM execution mode.")
		return
	}
	if len(args) == 0 {
		fmt.Println("Usage: riscv <instruction>")
		return
	}

	// Join all arguments to reconstruct the instruction
	instruction := strings.Join(args, " ")

	// Parse and execute the RISC-V instruction
	if err := r.machine.ExecuteRISCInstruction(instruction); err != nil {
		fmt.Printf("Error executing RISC-V instruction: %v\n", err)
		return
	}

	fmt.Printf("Executed RISC-V instruction: %s\n", instruction)
}

func (r *REPL) handleLoadCommand(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: load <file>")
		return
	}

	file := args[0]
	if err := r.machine.LoadRISCProgram(file); err != nil {
		fmt.Printf("Error loading RISC-V program: %v\n", err)
		return
	}

	fmt.Printf("Loaded RISC-V program from %s\n", file)
}

func (r *REPL) handleRunCommand() {
	if r.useHost {
		// Execute program using host-native execution
		for _, inst := range r.machine.GetRISCProgram() {
			if isQuantumInstruction(inst.Opcode) {
				if err := r.hostMachine.ExecuteQuantumRISCV(inst); err != nil {
					fmt.Printf("Error executing quantum instruction: %v\n", err)
					return
				}
			} else {
				// Execute classical RISC-V instructions using the standard machine
				var instruction string
				switch inst.Opcode {
				case "add", "sub", "and", "or", "xor", "sll", "srl", "sra", "slt", "sltu":
					instruction = fmt.Sprintf("%s x%d, x%d, x%d", inst.Opcode, inst.Rd, inst.Rs1, inst.Rs2)
				case "addi", "slli", "srli", "srai", "andi", "ori", "xori", "slti", "sltiu":
					instruction = fmt.Sprintf("%s x%d, x%d, %d", inst.Opcode, inst.Rd, inst.Rs1, inst.Imm)
				case "lui", "auipc":
					instruction = fmt.Sprintf("%s x%d, %d", inst.Opcode, inst.Rd, inst.Imm)
				case "jal":
					instruction = fmt.Sprintf("%s x%d, %d", inst.Opcode, inst.Rd, inst.Offset)
				case "jalr":
					instruction = fmt.Sprintf("%s x%d, x%d, %d", inst.Opcode, inst.Rd, inst.Rs1, inst.Offset)
				case "beq", "bne", "blt", "bge", "bltu", "bgeu":
					instruction = fmt.Sprintf("%s x%d, x%d, %d", inst.Opcode, inst.Rs1, inst.Rs2, inst.Offset)
				case "lw", "lh", "lb", "lwu", "lhu", "lbu":
					instruction = fmt.Sprintf("%s x%d, %d(x%d)", inst.Opcode, inst.Rd, inst.Offset, inst.Rs1)
				case "sw", "sh", "sb":
					instruction = fmt.Sprintf("%s x%d, %d(x%d)", inst.Opcode, inst.Rs2, inst.Offset, inst.Rs1)
				default:
					fmt.Printf("Unknown instruction type: %s\n", inst.Opcode)
					return
				}
				if err := r.machine.ExecuteRISCInstruction(instruction); err != nil {
					fmt.Printf("Error executing RISC-V instruction: %v\n", err)
					return
				}
			}
		}
		fmt.Println("RISC-V program executed successfully using host-native execution")
	} else {
		// Execute program using VM mode
		if err := r.machine.ExecuteRISCProgram(); err != nil {
			fmt.Printf("Error running RISC-V program: %v\n", err)
			return
		}
		fmt.Println("RISC-V program executed successfully using VM mode")
	}
}

func (r *REPL) handleModeCommand() {
	r.useHost = !r.useHost
	if r.useHost {
		fmt.Println("Switched to host-native execution mode")
	} else {
		fmt.Println("Switched to VM execution mode")
	}
}

func (r *REPL) handleRegistersCommand() {
	if r.useHost {
		fmt.Println("Error: Register inspection is exclusive to VM execution mode.")
		return
	}
	registers := r.machine.GetRegisters()
	fmt.Println("RISC-V Registers:")
	for i, reg := range registers {
		fmt.Printf("  x%d: %d\n", i, reg)
	}
}
