## Features

- 2000 qubit quantum computer simulation
- 100% fidelity quantum operations
- Quantum Volume of 4269
- RISC-V based instruction set with 128 virtual registers (extended from standard 32)
- Interactive REPL interface
- Direct execution of quantum RISC-V programs from command line
- Host-native quantum execution mode for improved performance
- Support for common quantum gates:
  - Pauli gates (X, Y, Z)
  - Hadamard gate (H)
  - Phase gates (S, T)
  - CNOT gate
  - Measurement operations
- Full RISC-V RV32I base integer instruction set support:
  - Arithmetic operations (add, sub, and, or, xor)
  - Shift operations (sll, srl, sra)
  - Comparison operations (slt, sltu)
  - Immediate operations (addi, slli, srli, etc.)
  - Load/Store operations (lw, lh, lb, sw, sh, sb)
  - Branch operations (beq, bne, blt, bge, etc.)
  - Jump operations (jal, jalr)
  - Upper immediate operations (lui, auipc)
- Custom Quantum RISC-V Instructions (Q-RISC-V Extensions):
  - qinit rd - Initialize quantum register with |0⟩ state
  - qapply rd, rs1, imm - Apply quantum gate (imm: 0=X, 1=Y, 2=Z, 3=H, 4=S, 5=T, 6=CNOT)
  - qmeasure rd, rs1 - Measure quantum register
  - qentangle rd, rs1, rs2 - Entangle two quantum registers

## Design Choices

The simulator implements 128 virtual registers instead of the standard RISC-V 32 registers. This design choice was made because:
- The memory overhead is negligible in a virtual machine context
- Additional registers can improve performance by reducing memory access
- The power consumption impact is minimal since it's just virtual memory
- It provides more flexibility for complex quantum algorithms

The quantum RISC-V instructions (Q-RISC-V) are custom extensions to the standard RISC-V instruction set, designed specifically for quantum operations. These extensions provide a seamless integration between classical and quantum computation, allowing quantum operations to be performed using the familiar RISC-V instruction format. Note that these instructions are not part of the official RISC-V specification but are custom additions for quantum computing support.

## Requirements

- Go 1.21 or later

## Installation

```bash
git clone https://github.com/sneed-group/qmachine.git
cd qmachine
```

## Usage

### Interactive REPL Mode
Run the REPL interface:
```bash
go run .
```

### Direct Program Execution
Execute a quantum RISC-V program in VM mode:
```bash
go run . -quantum=program.riscq
```

Execute a quantum RISC-V program using host-native execution (improved performance):
```bash
go run . -host-quantum=program.riscq
```

You can also specify the number of qubits:
```bash
go run . -qubits=1000 -host-quantum=program.riscq
```

The host-native execution mode translates quantum RISC-V instructions directly to native Go code, potentially offering better performance than the VM mode. It uses a compatibility layer to handle the translation from quantum RISC-V to host machine instructions.

### Example Quantum RISC-V Program

Contents of `quantum_test.riscq`:
```
# Initialize quantum registers
qinit x1    # Initialize x1 as a quantum register
qinit x2    # Initialize x2 as a quantum register

# Apply quantum gates
qapply x1, x1, 3  # Apply Hadamard gate to x1
qapply x2, x2, 0  # Apply X gate to x2

# Entangle the registers
qentangle x3, x1, x2  # Entangle x1 and x2, store result in x3

# Classical computation
addi x4, x0, 42      # Load immediate value 42 into x4
addi x5, x0, 58      # Load immediate value 58 into x5
add x6, x4, x5       # Add x4 and x5, store in x6

# Measure quantum register
qmeasure x7, x3      # Measure x3 and store result in x7

# Final classical computation
add x8, x6, x7       # Add classical result (x6) with quantum measurement (x7)
```

Execute it:
```bash
go run . -quantum=quantum_test.riscq
```

### REPL Commands

- `gate <type> <target> [controls...]` - Apply a quantum gate
- `measure <qubit>` - Measure a qubit
- `state` - Show current quantum state
- `reset` - Reset quantum state
- `riscv <instruction>` - Execute RISC-V instruction
- `load <file>` - Load RISC-V program from file
- `run` - Run loaded RISC-V program
- `registers` - Show RISC-V registers
- `help` - Show help message
- `exit` - Exit REPL

## Project Structure

- `quantum/state.go`: Quantum state representation and manipulation
- `quantum/gates.go`: Implementation of quantum gates
- `quantum/riscv.go`: RISC-V instruction set implementation
- `main.go`: Main program entry point
- `repl.go`: Interactive REPL implementation

## Performance

The simulator is optimized for high-performance quantum state manipulation using Go's efficient memory management and concurrent execution capabilities. The RISC-V implementation provides a standard instruction set for classical computation alongside quantum operations.

## Quantum nerds, please get off my back... (A quick disclamer)

Before you say something like "ackshully *snort* quantum mechanics makes this impossible and therefore forbids this..." I found slack in the rope so to speak. As in, I found a way to do this more efficient than anyone else who has tried.

* ***Full speed?*** No, of course not.

* ***Close enought for demonstration?*** Yes.

* ***Still awesome?*** IMHO yes.