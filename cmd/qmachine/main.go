package main

import (
	"flag"
	"fmt"
	"os"

	"qmachine/quantum"
	"qmachine/repl"
)

func main() {
	// Define command-line flags
	numQubits := flag.Int("qubits", 2000, "Number of qubits for the quantum computer")
	quantumFile := flag.String("quantum", "", "Path to quantum RISC-V file to execute")
	hostQuantumFile := flag.String("host-quantum", "", "Path to quantum RISC-V file to execute on host")
	flag.Parse()

	// Create the quantum computer REPL
	replInstance := repl.New(*numQubits)

	// Handle file execution modes
	if *hostQuantumFile != "" {
		fmt.Printf("Executing quantum RISC-V file on host: %s\n", *hostQuantumFile)
		if err := executeHostQuantumFile(*hostQuantumFile, *numQubits); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Quantum RISC-V program executed successfully using host-native execution")
		os.Exit(0)
	}

	if *quantumFile != "" {
		fmt.Printf("Executing quantum RISC-V file in VM mode: %s\n", *quantumFile)
		machine := quantum.NewQuantumRISCVMachine(*numQubits)

		// Load and execute the program
		if err := machine.LoadRISCProgram(*quantumFile); err != nil {
			fmt.Printf("Error loading quantum RISC-V program: %v\n", err)
			os.Exit(1)
		}

		// Print initial state
		fmt.Printf("\nInitial register state:\n")
		printRegisters(machine.GetRegisters())

		// Execute the program
		if err := machine.ExecuteRISCProgram(); err != nil {
			fmt.Printf("Error executing quantum RISC-V program: %v\n", err)
			os.Exit(1)
		}

		// Print final state
		fmt.Printf("\nFinal register state:\n")
		printRegisters(machine.GetRegisters())

		fmt.Println("\nQuantum RISC-V program executed successfully")
		os.Exit(0)
	}

	// Start interactive REPL mode
	replInstance.Start()
}

// executeHostQuantumFile executes a quantum RISC-V file using host-native execution
func executeHostQuantumFile(filename string, numQubits int) error {
	// Create a VM just to parse the program
	machine := quantum.NewQuantumRISCVMachine(numQubits)
	if err := machine.LoadRISCProgram(filename); err != nil {
		return fmt.Errorf("error loading quantum RISC-V program: %v", err)
	}

	// Create host machine for native execution
	hostMachine := quantum.NewHostQuantumMachine(numQubits)

	// Program counter for control flow
	pc := uint32(0)
	program := machine.GetRISCProgram()

	// Execute instructions until we reach the end of the program
	for pc < uint32(len(program)) {
		inst := program[pc]

		if isQuantumInstruction(inst.Opcode) {
			// Execute quantum instructions using host-native execution
			if err := hostMachine.ExecuteQuantumRISCV(inst); err != nil {
				return fmt.Errorf("error executing quantum instruction on host at PC %d: %v", pc, err)
			}
			pc++
		} else {
			// Execute classical RISC-V instructions
			switch inst.Opcode {
			case "add", "sub", "and", "or", "xor", "sll", "srl", "sra", "slt", "sltu":
				// R-type instructions
				rs1 := hostMachine.GetRegister(inst.Rs1)
				rs2 := hostMachine.GetRegister(inst.Rs2)
				var result uint64
				switch inst.Opcode {
				case "add":
					result = rs1 + rs2
				case "sub":
					result = rs1 - rs2
				case "and":
					result = rs1 & rs2
				case "or":
					result = rs1 | rs2
				case "xor":
					result = rs1 ^ rs2
				case "sll":
					result = rs1 << rs2
				case "srl":
					result = rs1 >> rs2
				case "sra":
					result = uint64(int64(rs1) >> rs2)
				case "slt":
					if int64(rs1) < int64(rs2) {
						result = 1
					}
				case "sltu":
					if rs1 < rs2 {
						result = 1
					}
				}
				hostMachine.SetRegister(inst.Rd, result)
				pc++

			case "addi", "slli", "srli", "srai", "andi", "ori", "xori", "slti", "sltiu":
				// I-type instructions
				rs1 := hostMachine.GetRegister(inst.Rs1)
				var result uint64
				switch inst.Opcode {
				case "addi":
					result = rs1 + uint64(inst.Imm)
				case "slli":
					result = rs1 << uint64(inst.Imm)
				case "srli":
					result = rs1 >> uint64(inst.Imm)
				case "srai":
					result = uint64(int64(rs1) >> uint64(inst.Imm))
				case "andi":
					result = rs1 & uint64(inst.Imm)
				case "ori":
					result = rs1 | uint64(inst.Imm)
				case "xori":
					result = rs1 ^ uint64(inst.Imm)
				case "slti":
					if int64(rs1) < inst.Imm {
						result = 1
					}
				case "sltiu":
					if rs1 < uint64(inst.Imm) {
						result = 1
					}
				}
				hostMachine.SetRegister(inst.Rd, result)
				pc++

			case "lui", "auipc":
				// U-type instructions
				switch inst.Opcode {
				case "lui":
					hostMachine.SetRegister(inst.Rd, uint64(inst.Imm<<12))
				case "auipc":
					hostMachine.SetRegister(inst.Rd, uint64(pc)+uint64(inst.Imm<<12))
				}
				pc++

			case "jal":
				// J-type instruction
				hostMachine.SetRegister(inst.Rd, uint64(pc+1))
				pc = uint32(int64(pc) + inst.Offset)

			case "jalr":
				// I-type jump instruction
				nextPc := uint32(int64(hostMachine.GetRegister(inst.Rs1)) + inst.Offset)
				hostMachine.SetRegister(inst.Rd, uint64(pc+1))
				pc = nextPc

			case "beq", "bne", "blt", "bge", "bltu", "bgeu":
				// B-type instructions
				rs1 := hostMachine.GetRegister(inst.Rs1)
				rs2 := hostMachine.GetRegister(inst.Rs2)
				var taken bool
				switch inst.Opcode {
				case "beq":
					taken = rs1 == rs2
				case "bne":
					taken = rs1 != rs2
				case "blt":
					taken = int64(rs1) < int64(rs2)
				case "bge":
					taken = int64(rs1) >= int64(rs2)
				case "bltu":
					taken = rs1 < rs2
				case "bgeu":
					taken = rs1 >= rs2
				}
				if taken {
					pc = uint32(int64(pc) + inst.Offset)
				} else {
					pc++
				}

			case "lw", "lh", "lb", "lwu", "lhu", "lbu":
				// Load instructions
				addr := uint32(int64(hostMachine.GetRegister(inst.Rs1)) + inst.Offset)
				var size uint8
				var signExtend bool
				switch inst.Opcode {
				case "lw":
					size = 4
					signExtend = true
				case "lh":
					size = 2
					signExtend = true
				case "lb":
					size = 1
					signExtend = true
				case "lwu":
					size = 4
					signExtend = false
				case "lhu":
					size = 2
					signExtend = false
				case "lbu":
					size = 1
					signExtend = false
				}
				val, err := hostMachine.LoadMemory(addr, size)
				if err != nil {
					return fmt.Errorf("error at PC %d: %v", pc, err)
				}
				if signExtend {
					switch size {
					case 1:
						val = uint64(int8(val))
					case 2:
						val = uint64(int16(val))
					case 4:
						val = uint64(int32(val))
					}
				}
				hostMachine.SetRegister(inst.Rd, val)
				pc++

			case "sw", "sh", "sb":
				// Store instructions
				addr := uint32(int64(hostMachine.GetRegister(inst.Rs1)) + inst.Offset)
				val := hostMachine.GetRegister(inst.Rs2)
				var size uint8
				switch inst.Opcode {
				case "sw":
					size = 4
				case "sh":
					size = 2
				case "sb":
					size = 1
				}
				if err := hostMachine.StoreMemory(addr, val, size); err != nil {
					return fmt.Errorf("error at PC %d: %v", pc, err)
				}
				pc++

			default:
				return fmt.Errorf("unknown instruction type at PC %d: %s", pc, inst.Opcode)
			}
		}
	}

	return nil
}

// printRegisters prints the current state of the registers
func printRegisters(registers [128]uint64) {
	for i, reg := range registers {
		if reg != 0 { // Only print non-zero registers to reduce noise
			fmt.Printf("  x%d: %d\n", i, reg)
		}
	}
}

// isQuantumInstruction checks if an instruction is a quantum instruction
func isQuantumInstruction(opcode string) bool {
	switch opcode {
	case "qinit", "qapply", "qmeasure", "qentangle":
		return true
	default:
		return false
	}
}
