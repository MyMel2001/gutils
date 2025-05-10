package quantum

import (
	"fmt"
	"math"
	"math/cmplx"
)

// HostQuantumState represents a quantum state optimized for host execution
type HostQuantumState struct {
	amplitudes []Complex128
	numQubits  int
}

// NewHostQuantumState creates a new host-optimized quantum state
func NewHostQuantumState(numQubits int) *HostQuantumState {
	size := 1 << numQubits
	return &HostQuantumState{
		amplitudes: make([]Complex128, size),
		numQubits:  numQubits,
	}
}

// HostQuantumMachine represents a quantum computer optimized for host execution
type HostQuantumMachine struct {
	state       *HostQuantumState
	registers   [128]uint64
	quantumRegs [128]*HostQuantumState
	memory      []byte
}

// NewHostQuantumMachine creates a new host-optimized quantum machine
func NewHostQuantumMachine(numQubits int) *HostQuantumMachine {
	return &HostQuantumMachine{
		state:       NewHostQuantumState(numQubits),
		registers:   [128]uint64{},
		quantumRegs: [128]*HostQuantumState{},
		memory:      make([]byte, 1024*1024),
	}
}

// ExecuteQuantumRISCV executes a quantum RISC-V instruction on the host
func (m *HostQuantumMachine) ExecuteQuantumRISCV(inst RISCInstruction) error {
	switch inst.Opcode {
	case "qinit":
		// Initialize quantum register with |0⟩ state
		m.quantumRegs[inst.Rd] = NewHostQuantumState(1)
		m.quantumRegs[inst.Rd].amplitudes[0] = 1.0
	case "qapply":
		// Apply quantum gate using host-optimized operations
		if m.quantumRegs[inst.Rs1] == nil {
			return fmt.Errorf("quantum register x%d not initialized", inst.Rs1)
		}
		gateType := uint8(inst.Imm)
		m.applyHostGate(gateType, m.quantumRegs[inst.Rs1])
	case "qmeasure":
		// Measure quantum register using host-optimized measurement
		if m.quantumRegs[inst.Rs1] == nil {
			return fmt.Errorf("quantum register x%d not initialized", inst.Rs1)
		}
		result := m.measureHostState(m.quantumRegs[inst.Rs1])
		m.registers[inst.Rd] = result
	case "qentangle":
		// Entangle two quantum registers using host-optimized operations
		if m.quantumRegs[inst.Rs1] == nil || m.quantumRegs[inst.Rs2] == nil {
			return fmt.Errorf("quantum registers not initialized")
		}
		entangled := NewHostQuantumState(2)
		m.entangleHostStates(m.quantumRegs[inst.Rs1], m.quantumRegs[inst.Rs2], entangled)
		m.quantumRegs[inst.Rd] = entangled
	default:
		return fmt.Errorf("unknown quantum instruction: %s", inst.Opcode)
	}
	return nil
}

// applyHostGate applies a quantum gate using host-optimized operations
func (m *HostQuantumMachine) applyHostGate(gateType uint8, state *HostQuantumState) {
	switch gateType {
	case 0: // X gate
		state.amplitudes[0], state.amplitudes[1] = state.amplitudes[1], state.amplitudes[0]
	case 1: // Y gate
		state.amplitudes[0], state.amplitudes[1] = -1i*state.amplitudes[1], 1i*state.amplitudes[0]
	case 2: // Z gate
		state.amplitudes[1] = -state.amplitudes[1]
	case 3: // H gate
		invSqrt2 := complex(1.0/math.Sqrt2, 0)
		a, b := state.amplitudes[0], state.amplitudes[1]
		state.amplitudes[0] = invSqrt2 * (a + b)
		state.amplitudes[1] = invSqrt2 * (a - b)
	case 4: // S gate
		state.amplitudes[1] *= 1i
	case 5: // T gate
		state.amplitudes[1] *= cmplx.Exp(1i * math.Pi / 4)
	case 6: // CNOT gate
		// For 2-qubit states
		if state.numQubits == 2 {
			state.amplitudes[2], state.amplitudes[3] = state.amplitudes[3], state.amplitudes[2]
		}
	}
	m.normalizeHostState(state)
}

// measureHostState performs measurement using host-optimized operations
func (m *HostQuantumMachine) measureHostState(state *HostQuantumState) uint64 {
	// Calculate probabilities
	p0 := real(state.amplitudes[0] * cmplx.Conj(state.amplitudes[0]))
	p1 := real(state.amplitudes[1] * cmplx.Conj(state.amplitudes[1]))

	// Simple deterministic measurement (in a real implementation, this would be probabilistic)
	if p0 > p1 {
		return 0
	}
	return 1
}

// entangleHostStates entangles two quantum states using host-optimized operations
func (m *HostQuantumMachine) entangleHostStates(state1, state2, result *HostQuantumState) {
	// Create Bell state |Φ+⟩ = (|00⟩ + |11⟩)/√2
	result.amplitudes[0] = 1.0 / math.Sqrt2
	result.amplitudes[3] = 1.0 / math.Sqrt2
}

// normalizeHostState normalizes a quantum state using host-optimized operations
func (m *HostQuantumMachine) normalizeHostState(state *HostQuantumState) {
	var sum float64
	for _, amp := range state.amplitudes {
		sum += real(amp * cmplx.Conj(amp))
	}
	norm := 1.0 / math.Sqrt(sum)
	for i := range state.amplitudes {
		state.amplitudes[i] *= complex(norm, 0)
	}
}

// GetRegisters returns the current state of all registers
func (m *HostQuantumMachine) GetRegisters() [128]uint64 {
	return m.registers
}

// GetState returns the current quantum state
func (m *HostQuantumMachine) GetState() *HostQuantumState {
	return m.state
}

// SetRegister sets the value of a register
func (m *HostQuantumMachine) SetRegister(reg uint8, value uint64) {
	if reg == 0 {
		return // x0 is hardwired to zero
	}
	m.registers[reg] = value
}

// GetRegister gets the value of a register
func (m *HostQuantumMachine) GetRegister(reg uint8) uint64 {
	return m.registers[reg]
}

// LoadMemory loads a value from memory
func (m *HostQuantumMachine) LoadMemory(addr uint32, size uint8) (uint64, error) {
	switch size {
	case 1: // byte
		if addr >= uint32(len(m.memory)) {
			return 0, fmt.Errorf("memory access out of bounds: addr %d", addr)
		}
		return uint64(m.memory[addr]), nil
	case 2: // halfword
		if addr+1 >= uint32(len(m.memory)) {
			return 0, fmt.Errorf("memory access out of bounds: addr %d", addr)
		}
		return uint64(m.memory[addr]) | uint64(m.memory[addr+1])<<8, nil
	case 4: // word
		if addr+3 >= uint32(len(m.memory)) {
			return 0, fmt.Errorf("memory access out of bounds: addr %d", addr)
		}
		return uint64(m.memory[addr]) | uint64(m.memory[addr+1])<<8 | uint64(m.memory[addr+2])<<16 | uint64(m.memory[addr+3])<<24, nil
	default:
		return 0, fmt.Errorf("invalid memory access size: %d", size)
	}
}

// StoreMemory stores a value to memory
func (m *HostQuantumMachine) StoreMemory(addr uint32, value uint64, size uint8) error {
	switch size {
	case 1: // byte
		if addr >= uint32(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds: addr %d", addr)
		}
		m.memory[addr] = byte(value)
	case 2: // halfword
		if addr+1 >= uint32(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds: addr %d", addr)
		}
		m.memory[addr] = byte(value)
		m.memory[addr+1] = byte(value >> 8)
	case 4: // word
		if addr+3 >= uint32(len(m.memory)) {
			return fmt.Errorf("memory access out of bounds: addr %d", addr)
		}
		m.memory[addr] = byte(value)
		m.memory[addr+1] = byte(value >> 8)
		m.memory[addr+2] = byte(value >> 16)
		m.memory[addr+3] = byte(value >> 24)
	default:
		return fmt.Errorf("invalid memory access size: %d", size)
	}
	return nil
}
