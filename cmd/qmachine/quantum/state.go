package quantum

import (
	"math"
	"math/cmplx"
)

// Complex128 represents a complex number with float64 precision
type Complex128 = complex128

// QuantumState represents the state of a quantum register
type QuantumState struct {
	amplitudes []Complex128
	numQubits  int
}

// NewQuantumState creates a new quantum state with the specified number of qubits
func NewQuantumState(numQubits int) *QuantumState {
	size := 1 << numQubits
	return &QuantumState{
		amplitudes: make([]Complex128, size),
		numQubits:  numQubits,
	}
}

// InitializeZeroState sets the quantum state to |0⟩^⊗n
func (qs *QuantumState) InitializeZeroState() {
	qs.amplitudes[0] = 1.0
}

// GetAmplitude returns the amplitude at the specified index
func (qs *QuantumState) GetAmplitude(index int) Complex128 {
	return qs.amplitudes[index]
}

// SetAmplitude sets the amplitude at the specified index
func (qs *QuantumState) SetAmplitude(index int, value Complex128) {
	qs.amplitudes[index] = value
}

// Normalize normalizes the quantum state
func (qs *QuantumState) Normalize() {
	var sum float64
	for _, amp := range qs.amplitudes {
		sum += real(amp*cmplx.Conj(amp))
	}
	norm := 1.0 / math.Sqrt(sum)
	for i := range qs.amplitudes {
		qs.amplitudes[i] *= complex(norm, 0)
	}
}

// NumQubits returns the number of qubits in the quantum state
func (qs *QuantumState) NumQubits() int {
	return qs.numQubits
}

// Clone creates a deep copy of the quantum state
func (qs *QuantumState) Clone() *QuantumState {
	clone := NewQuantumState(qs.numQubits)
	copy(clone.amplitudes, qs.amplitudes)
	return clone
} 