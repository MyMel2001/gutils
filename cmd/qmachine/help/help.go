// Package help provides help text and documentation for the QMachine REPL
package help

// GetBasicCommands returns the basic command help text
func GetBasicCommands() string {
	return `Available commands:
  gate <type> <target> [controls...] - Apply a quantum gate
  measure <qubit>                    - Measure a qubit
  state                              - Show current quantum state
  reset                              - Reset quantum state
  riscv <instruction>                - Execute RISC-V instruction
  load <file>                        - Load RISC-V program from file
  run                                - Run loaded RISC-V program
  run-host                           - Run loaded program using host-native execution
  mode                               - Toggle between VM and host-native execution
  registers                          - Show RISC-V registers
  help                               - Show this help message
  exit                               - Exit REPL

Available gates: X, Y, Z, H, S, T, CNOT`
}

// GetQuantumInstructions returns help text for quantum RISC-V instructions
func GetQuantumInstructions() string {
	return `Custom Quantum RISC-V Instructions (Q-RISC-V Extensions):
  qinit rd                          - Initialize quantum register with |0⟩
  qapply rd, rs1, imm              - Apply quantum gate (imm: 0=X, 1=Y, 2=Z, 3=H, 4=S, 5=T, 6=CNOT)
  qmeasure rd, rs1                 - Measure quantum register
  qentangle rd, rs1, rs2          - Entangle two quantum registers`
}

// GetRISCVInstructions returns help text for standard RISC-V instructions
func GetRISCVInstructions() string {
	return `Standard RISC-V Instructions:
  add rd, rs1, rs2    - Add registers
  sub rd, rs1, rs2    - Subtract registers
  and rd, rs1, rs2    - AND registers
  or rd, rs1, rs2     - OR registers
  xor rd, rs1, rs2    - XOR registers
  sll rd, rs1, rs2    - Shift left logical
  srl rd, rs1, rs2    - Shift right logical
  sra rd, rs1, rs2    - Shift right arithmetic
  slt rd, rs1, rs2    - Set if less than
  sltu rd, rs1, rs2   - Set if less than unsigned
  addi rd, rs1, imm   - Add immediate
  slli rd, rs1, imm   - Shift left logical immediate
  srli rd, rs1, imm   - Shift right logical immediate
  srai rd, rs1, imm   - Shift right arithmetic immediate
  andi rd, rs1, imm   - AND immediate
  ori rd, rs1, imm    - OR immediate
  xori rd, rs1, imm   - XOR immediate
  slti rd, rs1, imm   - Set if less than immediate
  sltiu rd, rs1, imm  - Set if less than immediate unsigned
  lui rd, imm         - Load upper immediate
  auipc rd, imm       - Add upper immediate to PC
  jal rd, offset      - Jump and link
  jalr rd, rs1, offset - Jump and link register
  beq rs1, rs2, offset - Branch if equal
  bne rs1, rs2, offset - Branch if not equal
  blt rs1, rs2, offset - Branch if less than
  bge rs1, rs2, offset - Branch if greater or equal
  bltu rs1, rs2, offset - Branch if less than unsigned
  bgeu rs1, rs2, offset - Branch if greater or equal unsigned
  lw rd, offset(rs1)   - Load word
  lh rd, offset(rs1)   - Load halfword
  lb rd, offset(rs1)   - Load byte
  lwu rd, offset(rs1)  - Load word unsigned
  lhu rd, offset(rs1)  - Load halfword unsigned
  lbu rd, offset(rs1)  - Load byte unsigned
  sw rs2, offset(rs1)  - Store word
  sh rs2, offset(rs1)  - Store halfword
  sb rs2, offset(rs1)  - Store byte`
}
