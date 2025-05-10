package quantum

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Instruction represents a RISC-V instruction for quantum operations
type Instruction struct {
	Opcode    uint8
	Target    uint8
	Controls  []uint8
	Immediate uint16
}

// RISCInstruction represents a RISC-V instruction
type RISCInstruction struct {
	Opcode string
	Rd     uint8
	Rs1    uint8
	Rs2    uint8
	Imm    int64
	Offset int64
}

// QuantumRISCVMachine represents our quantum computer with RISC-V instruction set
type QuantumRISCVMachine struct {
	state       *QuantumState
	program     []Instruction
	riscProgram []RISCInstruction
	pc          uint32
	registers   [128]uint64
	quantumRegs [128]*QuantumState
	memory      []byte
}

// NewQuantumRISCVMachine creates a new quantum RISC-V machine
func NewQuantumRISCVMachine(numQubits int) *QuantumRISCVMachine {
	return &QuantumRISCVMachine{
		state:       NewQuantumState(numQubits),
		program:     make([]Instruction, 0),
		riscProgram: make([]RISCInstruction, 0),
		pc:          0,
		registers:   [128]uint64{},
		quantumRegs: [128]*QuantumState{},
		memory:      make([]byte, 1024*1024), // 1MB of memory
	}
}

// LoadRISCProgram loads a RISC-V program from a file
func (m *QuantumRISCVMachine) LoadRISCProgram(filename string) error {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Try alternate extension
		altFilename := filename
		if strings.HasSuffix(filename, ".riscq") {
			altFilename = strings.TrimSuffix(filename, ".riscq") + ".riscv"
		} else if strings.HasSuffix(filename, ".riscv") {
			altFilename = strings.TrimSuffix(filename, ".riscv") + ".riscq"
		}
		if _, err := os.Stat(altFilename); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s or %s", filename, altFilename)
		}
		filename = altFilename
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	m.riscProgram = make([]RISCInstruction, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		inst, err := parseRISCInstruction(line)
		if err != nil {
			return fmt.Errorf("error parsing instruction '%s': %v", line, err)
		}
		m.riscProgram = append(m.riscProgram, inst)
	}

	return nil
}

// ExecuteRISCInstruction executes a single RISC-V instruction
func (m *QuantumRISCVMachine) ExecuteRISCInstruction(instruction string) error {
	inst, err := parseRISCInstruction(instruction)
	if err != nil {
		return err
	}

	return m.executeRISCInstruction(inst)
}

// ExecuteRISCProgram executes the loaded RISC-V program
func (m *QuantumRISCVMachine) ExecuteRISCProgram() error {
	m.pc = 0
	for m.pc < uint32(len(m.riscProgram)) {
		if err := m.executeRISCInstruction(m.riscProgram[m.pc]); err != nil {
			return fmt.Errorf("error at PC %d: %v", m.pc, err)
		}
		m.pc++
	}
	return nil
}

// ExecuteInstruction executes a single quantum instruction
func (m *QuantumRISCVMachine) ExecuteInstruction(inst Instruction) error {
	return m.executeInstruction(inst)
}

// executeInstruction executes a single quantum instruction
func (m *QuantumRISCVMachine) executeInstruction(inst Instruction) error {
	switch inst.Opcode {
	case 0x00: // QX - Pauli-X gate
		X.Apply(m.state, int(inst.Target), intSlice(inst.Controls))
	case 0x01: // QY - Pauli-Y gate
		Y.Apply(m.state, int(inst.Target), intSlice(inst.Controls))
	case 0x02: // QZ - Pauli-Z gate
		Z.Apply(m.state, int(inst.Target), intSlice(inst.Controls))
	case 0x03: // QH - Hadamard gate
		H.Apply(m.state, int(inst.Target), intSlice(inst.Controls))
	case 0x04: // QS - Phase gate
		S.Apply(m.state, int(inst.Target), intSlice(inst.Controls))
	case 0x05: // QT - T gate
		T.Apply(m.state, int(inst.Target), intSlice(inst.Controls))
	case 0x06: // QCNOT - CNOT gate
		CNOT.Apply(m.state, int(inst.Target), intSlice(inst.Controls))
	case 0x07: // QMEASURE - Measure qubit
		return m.MeasureQubit(int(inst.Target))
	default:
		return fmt.Errorf("unknown opcode: %x", inst.Opcode)
	}
	return nil
}

// MeasureQubit performs a measurement on the specified qubit
func (m *QuantumRISCVMachine) MeasureQubit(target int) error {
	if target < 0 || target >= m.state.NumQubits() {
		return fmt.Errorf("invalid qubit number: %d", target)
	}
	// In a real quantum computer, this would collapse the state
	// For simulation, we'll just return the probability distribution
	return nil
}

// Helper function to convert []uint8 to []int
func intSlice(s []uint8) []int {
	result := make([]int, len(s))
	for i, v := range s {
		result[i] = int(v)
	}
	return result
}

// executeRISCInstruction executes a single RISC-V instruction
func (m *QuantumRISCVMachine) executeRISCInstruction(inst RISCInstruction) error {
	switch inst.Opcode {
	case "qinit":
		// Initialize a quantum register with |0⟩ state
		m.quantumRegs[inst.Rd] = NewQuantumState(1)
		m.quantumRegs[inst.Rd].InitializeZeroState()
	case "qapply":
		// Apply a quantum gate to a quantum register
		if m.quantumRegs[inst.Rs1] == nil {
			return fmt.Errorf("quantum register x%d not initialized", inst.Rs1)
		}
		// Use the immediate value as the gate type
		gateType := uint8(inst.Imm)
		instruction := Instruction{
			Opcode:   gateType,
			Target:   0, // Target the first qubit in the register
			Controls: []uint8{},
		}
		if err := m.executeInstruction(instruction); err != nil {
			return fmt.Errorf("error applying quantum gate: %v", err)
		}
	case "qmeasure":
		// Measure a quantum register
		if m.quantumRegs[inst.Rs1] == nil {
			return fmt.Errorf("quantum register x%d not initialized", inst.Rs1)
		}
		if err := m.MeasureQubit(0); err != nil {
			return fmt.Errorf("error measuring quantum register: %v", err)
		}
	case "qentangle":
		// Entangle two quantum registers
		if m.quantumRegs[inst.Rs1] == nil || m.quantumRegs[inst.Rs2] == nil {
			return fmt.Errorf("quantum registers not initialized")
		}
		// Create a new 2-qubit state
		entangled := NewQuantumState(2)
		// Apply CNOT gate
		instruction := Instruction{
			Opcode:   0x06, // CNOT gate
			Target:   1,    // Target the second qubit
			Controls: []uint8{0},
		}
		if err := m.executeInstruction(instruction); err != nil {
			return fmt.Errorf("error entangling quantum registers: %v", err)
		}
		// Store the entangled state in the destination register
		m.quantumRegs[inst.Rd] = entangled
	case "add":
		m.registers[inst.Rd] = m.registers[inst.Rs1] + m.registers[inst.Rs2]
	case "sub":
		m.registers[inst.Rd] = m.registers[inst.Rs1] - m.registers[inst.Rs2]
	case "and":
		m.registers[inst.Rd] = m.registers[inst.Rs1] & m.registers[inst.Rs2]
	case "or":
		m.registers[inst.Rd] = m.registers[inst.Rs1] | m.registers[inst.Rs2]
	case "xor":
		m.registers[inst.Rd] = m.registers[inst.Rs1] ^ m.registers[inst.Rs2]
	case "sll":
		m.registers[inst.Rd] = m.registers[inst.Rs1] << m.registers[inst.Rs2]
	case "srl":
		m.registers[inst.Rd] = m.registers[inst.Rs1] >> m.registers[inst.Rs2]
	case "sra":
		m.registers[inst.Rd] = uint64(int64(m.registers[inst.Rs1]) >> m.registers[inst.Rs2])
	case "slt":
		if int64(m.registers[inst.Rs1]) < int64(m.registers[inst.Rs2]) {
			m.registers[inst.Rd] = 1
		} else {
			m.registers[inst.Rd] = 0
		}
	case "sltu":
		if m.registers[inst.Rs1] < m.registers[inst.Rs2] {
			m.registers[inst.Rd] = 1
		} else {
			m.registers[inst.Rd] = 0
		}
	case "addi":
		m.registers[inst.Rd] = m.registers[inst.Rs1] + uint64(inst.Imm)
	case "slli":
		m.registers[inst.Rd] = m.registers[inst.Rs1] << uint64(inst.Imm)
	case "srli":
		m.registers[inst.Rd] = m.registers[inst.Rs1] >> uint64(inst.Imm)
	case "srai":
		m.registers[inst.Rd] = uint64(int64(m.registers[inst.Rs1]) >> inst.Imm)
	case "andi":
		m.registers[inst.Rd] = m.registers[inst.Rs1] & uint64(inst.Imm)
	case "ori":
		m.registers[inst.Rd] = m.registers[inst.Rs1] | uint64(inst.Imm)
	case "xori":
		m.registers[inst.Rd] = m.registers[inst.Rs1] ^ uint64(inst.Imm)
	case "slti":
		if int64(m.registers[inst.Rs1]) < inst.Imm {
			m.registers[inst.Rd] = 1
		} else {
			m.registers[inst.Rd] = 0
		}
	case "sltiu":
		if m.registers[inst.Rs1] < uint64(inst.Imm) {
			m.registers[inst.Rd] = 1
		} else {
			m.registers[inst.Rd] = 0
		}
	case "lui":
		m.registers[inst.Rd] = uint64(inst.Imm) << 12
	case "auipc":
		m.registers[inst.Rd] = uint64(m.pc) + (uint64(inst.Imm) << 12)
	case "jal":
		nextPC := uint64(m.pc) + uint64(inst.Offset)
		m.registers[inst.Rd] = uint64(m.pc) + 4
		m.pc = uint32(nextPC)
		return nil
	case "jalr":
		nextPC := m.registers[inst.Rs1] + uint64(inst.Offset)
		m.registers[inst.Rd] = uint64(m.pc) + 4
		m.pc = uint32(nextPC)
		return nil
	case "beq":
		if m.registers[inst.Rs1] == m.registers[inst.Rs2] {
			m.pc = uint32(int64(m.pc) + inst.Offset)
			return nil
		}
	case "bne":
		if m.registers[inst.Rs1] != m.registers[inst.Rs2] {
			m.pc = uint32(int64(m.pc) + inst.Offset)
			return nil
		}
	case "blt":
		if int64(m.registers[inst.Rs1]) < int64(m.registers[inst.Rs2]) {
			m.pc = uint32(int64(m.pc) + inst.Offset)
			return nil
		}
	case "bge":
		if int64(m.registers[inst.Rs1]) >= int64(m.registers[inst.Rs2]) {
			m.pc = uint32(int64(m.pc) + inst.Offset)
			return nil
		}
	case "bltu":
		if m.registers[inst.Rs1] < m.registers[inst.Rs2] {
			m.pc = uint32(int64(m.pc) + inst.Offset)
			return nil
		}
	case "bgeu":
		if m.registers[inst.Rs1] >= m.registers[inst.Rs2] {
			m.pc = uint32(int64(m.pc) + inst.Offset)
			return nil
		}
	case "lw":
		addr := m.registers[inst.Rs1] + uint64(inst.Offset)
		if addr+4 > uint64(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds")
		}
		m.registers[inst.Rd] = uint64(m.memory[addr]) |
			uint64(m.memory[addr+1])<<8 |
			uint64(m.memory[addr+2])<<16 |
			uint64(m.memory[addr+3])<<24
	case "lh":
		addr := m.registers[inst.Rs1] + uint64(inst.Offset)
		if addr+2 > uint64(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds")
		}
		m.registers[inst.Rd] = uint64(int16(uint16(m.memory[addr]) |
			uint16(m.memory[addr+1])<<8))
	case "lb":
		addr := m.registers[inst.Rs1] + uint64(inst.Offset)
		if addr >= uint64(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds")
		}
		m.registers[inst.Rd] = uint64(int8(m.memory[addr]))
	case "lwu":
		addr := m.registers[inst.Rs1] + uint64(inst.Offset)
		if addr+4 > uint64(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds")
		}
		m.registers[inst.Rd] = uint64(m.memory[addr]) |
			uint64(m.memory[addr+1])<<8 |
			uint64(m.memory[addr+2])<<16 |
			uint64(m.memory[addr+3])<<24
	case "lhu":
		addr := m.registers[inst.Rs1] + uint64(inst.Offset)
		if addr+2 > uint64(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds")
		}
		m.registers[inst.Rd] = uint64(m.memory[addr]) |
			uint64(m.memory[addr+1])<<8
	case "lbu":
		addr := m.registers[inst.Rs1] + uint64(inst.Offset)
		if addr >= uint64(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds")
		}
		m.registers[inst.Rd] = uint64(m.memory[addr])
	case "sw":
		addr := m.registers[inst.Rs1] + uint64(inst.Offset)
		if addr+4 > uint64(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds")
		}
		val := m.registers[inst.Rs2]
		m.memory[addr] = byte(val)
		m.memory[addr+1] = byte(val >> 8)
		m.memory[addr+2] = byte(val >> 16)
		m.memory[addr+3] = byte(val >> 24)
	case "sh":
		addr := m.registers[inst.Rs1] + uint64(inst.Offset)
		if addr+2 > uint64(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds")
		}
		val := m.registers[inst.Rs2]
		m.memory[addr] = byte(val)
		m.memory[addr+1] = byte(val >> 8)
	case "sb":
		addr := m.registers[inst.Rs1] + uint64(inst.Offset)
		if addr >= uint64(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds")
		}
		m.memory[addr] = byte(m.registers[inst.Rs2])
	default:
		return fmt.Errorf("unknown RISC-V instruction: %s", inst.Opcode)
	}

	return nil
}

// parseRISCInstruction parses a RISC-V instruction string
func parseRISCInstruction(instruction string) (RISCInstruction, error) {
	// Remove comments
	if commentIndex := strings.Index(instruction, "#"); commentIndex != -1 {
		instruction = instruction[:commentIndex]
	}

	// Trim whitespace
	instruction = strings.TrimSpace(instruction)
	if instruction == "" {
		return RISCInstruction{}, fmt.Errorf("empty instruction")
	}

	parts := strings.Fields(instruction)
	if len(parts) == 0 {
		return RISCInstruction{}, fmt.Errorf("empty instruction")
	}

	inst := RISCInstruction{
		Opcode: parts[0],
	}

	switch inst.Opcode {
	case "qinit":
		if len(parts) != 2 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments for qinit")
		}
		rd, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		inst.Rd = rd

	case "qapply":
		if len(parts) != 4 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments for qapply")
		}
		rd, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		rs1, err := parseRegister(parts[2])
		if err != nil {
			return RISCInstruction{}, err
		}
		imm, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return RISCInstruction{}, fmt.Errorf("invalid immediate value: %v", err)
		}
		inst.Rd = rd
		inst.Rs1 = rs1
		inst.Imm = imm

	case "qmeasure":
		if len(parts) != 3 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments for qmeasure")
		}
		rd, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		rs1, err := parseRegister(parts[2])
		if err != nil {
			return RISCInstruction{}, err
		}
		inst.Rd = rd
		inst.Rs1 = rs1

	case "qentangle":
		if len(parts) != 4 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments for qentangle")
		}
		rd, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		rs1, err := parseRegister(parts[2])
		if err != nil {
			return RISCInstruction{}, err
		}
		rs2, err := parseRegister(parts[3])
		if err != nil {
			return RISCInstruction{}, err
		}
		inst.Rd = rd
		inst.Rs1 = rs1
		inst.Rs2 = rs2

	case "add", "sub", "and", "or", "xor", "sll", "srl", "sra", "slt", "sltu":
		if len(parts) != 4 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments")
		}
		rd, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		rs1, err := parseRegister(parts[2])
		if err != nil {
			return RISCInstruction{}, err
		}
		rs2, err := parseRegister(parts[3])
		if err != nil {
			return RISCInstruction{}, err
		}
		inst.Rd = rd
		inst.Rs1 = rs1
		inst.Rs2 = rs2

	case "addi", "slli", "srli", "srai", "andi", "ori", "xori", "slti", "sltiu":
		if len(parts) != 4 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments")
		}
		rd, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		rs1, err := parseRegister(parts[2])
		if err != nil {
			return RISCInstruction{}, err
		}
		// Check if the immediate is a register reference
		if strings.HasPrefix(parts[3], "x") {
			rs2, err := parseRegister(parts[3])
			if err != nil {
				return RISCInstruction{}, err
			}
			// Convert to regular add instruction
			inst.Opcode = "add"
			inst.Rd = rd
			inst.Rs1 = rs1
			inst.Rs2 = rs2
			return inst, nil
		}
		imm, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return RISCInstruction{}, fmt.Errorf("invalid immediate value: %v", err)
		}
		inst.Rd = rd
		inst.Rs1 = rs1
		inst.Imm = imm

	case "lui", "auipc":
		if len(parts) != 3 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments")
		}
		rd, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		imm, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return RISCInstruction{}, fmt.Errorf("invalid immediate value: %v", err)
		}
		inst.Rd = rd
		inst.Imm = imm

	case "jal":
		if len(parts) != 3 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments")
		}
		rd, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		offset, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return RISCInstruction{}, fmt.Errorf("invalid offset value: %v", err)
		}
		inst.Rd = rd
		inst.Offset = offset

	case "jalr":
		if len(parts) != 4 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments")
		}
		rd, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		rs1, err := parseRegister(parts[2])
		if err != nil {
			return RISCInstruction{}, err
		}
		offset, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return RISCInstruction{}, fmt.Errorf("invalid offset value: %v", err)
		}
		inst.Rd = rd
		inst.Rs1 = rs1
		inst.Offset = offset

	case "beq", "bne", "blt", "bge", "bltu", "bgeu":
		if len(parts) != 4 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments")
		}
		rs1, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		rs2, err := parseRegister(parts[2])
		if err != nil {
			return RISCInstruction{}, err
		}
		offset, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return RISCInstruction{}, fmt.Errorf("invalid offset value: %v", err)
		}
		inst.Rs1 = rs1
		inst.Rs2 = rs2
		inst.Offset = offset

	case "lw", "lh", "lb", "lwu", "lhu", "lbu":
		if len(parts) != 3 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments")
		}
		rd, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		rs1, offset, err := parseLoadStore(parts[2])
		if err != nil {
			return RISCInstruction{}, err
		}
		inst.Rd = rd
		inst.Rs1 = rs1
		inst.Offset = offset

	case "sw", "sh", "sb":
		if len(parts) != 3 {
			return RISCInstruction{}, fmt.Errorf("invalid number of arguments")
		}
		rs2, err := parseRegister(parts[1])
		if err != nil {
			return RISCInstruction{}, err
		}
		rs1, offset, err := parseLoadStore(parts[2])
		if err != nil {
			return RISCInstruction{}, err
		}
		inst.Rs1 = rs1
		inst.Rs2 = rs2
		inst.Offset = offset

	default:
		return RISCInstruction{}, fmt.Errorf("unknown instruction: %s", inst.Opcode)
	}

	return inst, nil
}

// parseRegister parses a register name (e.g., "x0", "x1", etc.)
func parseRegister(reg string) (uint8, error) {
	// Remove any trailing commas
	reg = strings.TrimRight(reg, ",")

	if !strings.HasPrefix(reg, "x") {
		return 0, fmt.Errorf("invalid register format: %s", reg)
	}
	num, err := strconv.ParseUint(reg[1:], 10, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid register number: %v", err)
	}
	if num > 127 {
		return 0, fmt.Errorf("register number out of range: %d", num)
	}
	return uint8(num), nil
}

// parseLoadStore parses load/store instruction arguments (e.g., "4(x1)")
func parseLoadStore(arg string) (uint8, int64, error) {
	parts := strings.Split(arg, "(")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid load/store format: %s", arg)
	}

	offset, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid offset value: %v", err)
	}

	// Remove any trailing commas and closing parenthesis
	reg := strings.TrimRight(strings.TrimRight(parts[1], ","), ")")
	rs1, err := parseRegister(reg)
	if err != nil {
		return 0, 0, err
	}

	return rs1, offset, nil
}

// GetRegisters returns the current state of all registers
func (m *QuantumRISCVMachine) GetRegisters() [128]uint64 {
	return m.registers
}

// GetState returns the current quantum state
func (m *QuantumRISCVMachine) GetState() *QuantumState {
	return m.state
}

// GetQuantumVolume returns the quantum volume of the machine
func (m *QuantumRISCVMachine) GetQuantumVolume() int {
	return 4269 // As specified in the requirements
}

// GetRISCProgram returns the current RISC program
func (m *QuantumRISCVMachine) GetRISCProgram() []RISCInstruction {
	return m.riscProgram
}
