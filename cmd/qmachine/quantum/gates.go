package quantum

import (
	"math"
	"math/cmplx"
)

// Gate represents a quantum gate operation
type Gate interface {
	Apply(state *QuantumState, target int, controls []int)
}

// SingleQubitGate represents a gate that operates on a single qubit
type SingleQubitGate struct {
	matrix [2][2]Complex128
}

// TwoQubitGate represents a gate that operates on two qubits
type TwoQubitGate struct {
	matrix [4][4]Complex128
}

// Common quantum gates
var (
	// Pauli gates
	X = &SingleQubitGate{
		matrix: [2][2]Complex128{
			{0, 1},
			{1, 0},
		},
	}
	Y = &SingleQubitGate{
		matrix: [2][2]Complex128{
			{0, -1i},
			{1i, 0},
		},
	}
	Z = &SingleQubitGate{
		matrix: [2][2]Complex128{
			{1, 0},
			{0, -1},
		},
	}

	// Hadamard gate
	H = &SingleQubitGate{
		matrix: [2][2]Complex128{
			{1 / math.Sqrt2, 1 / math.Sqrt2},
			{1 / math.Sqrt2, -1 / math.Sqrt2},
		},
	}

	// Phase gates
	S = &SingleQubitGate{
		matrix: [2][2]Complex128{
			{1, 0},
			{0, 1i},
		},
	}
	T = &SingleQubitGate{
		matrix: [2][2]Complex128{
			{1, 0},
			{0, cmplx.Exp(1i * math.Pi / 4)},
		},
	}

	// CNOT gate
	CNOT = &TwoQubitGate{
		matrix: [4][4]Complex128{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 1},
			{0, 0, 1, 0},
		},
	}
)

// Apply implements the Gate interface for SingleQubitGate
func (g *SingleQubitGate) Apply(state *QuantumState, target int, controls []int) {
	size := 1 << state.numQubits
	newAmplitudes := make([]Complex128, size)
	
	for i := 0; i < size; i++ {
		// Check if control conditions are met
		controlMet := true
		for _, control := range controls {
			if (i>>control)&1 == 0 {
				controlMet = false
				break
			}
		}
		
		if controlMet {
			// Apply gate to target qubit
			targetBit := (i >> target) & 1
			otherBits := i & ^(1 << target)
			
			for j := 0; j < 2; j++ {
				newIndex := otherBits | (j << target)
				newAmplitudes[newIndex] += state.amplitudes[i] * g.matrix[targetBit][j]
			}
		} else {
			newAmplitudes[i] = state.amplitudes[i]
		}
	}
	
	state.amplitudes = newAmplitudes
	state.Normalize()
}

// Apply implements the Gate interface for TwoQubitGate
func (g *TwoQubitGate) Apply(state *QuantumState, target int, controls []int) {
	if len(controls) != 1 {
		panic("TwoQubitGate requires exactly one control qubit")
	}
	
	size := 1 << state.numQubits
	newAmplitudes := make([]Complex128, size)
	
	for i := 0; i < size; i++ {
		control := controls[0]
		controlBit := (i >> control) & 1
		
		if controlBit == 1 {
			// Apply two-qubit gate
			targetBit := (i >> target) & 1
			otherBits := i & ^(1 << target)
			
			for j := 0; j < 2; j++ {
				newIndex := otherBits | (j << target)
				newAmplitudes[newIndex] += state.amplitudes[i] * g.matrix[targetBit][j]
			}
		} else {
			newAmplitudes[i] = state.amplitudes[i]
		}
	}
	
	state.amplitudes = newAmplitudes
	state.Normalize()
} 