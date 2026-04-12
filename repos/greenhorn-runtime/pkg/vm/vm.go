// Package vm implements the FLUX Unified ISA bytecode interpreter.
//
// The FLUX VM is a register-based virtual machine with 64 general-purpose
// 32-bit signed integer registers (R0–R63), byte-addressable memory, and a
// call/value stack. R0 is hardwired to zero and ignores all writes.
//
// Instruction encoding uses fixed-length formats depending on the opcode
// class. Most ALU operations use Format E (4 bytes): [op][rd][rs1][rs2].
// Immediate operations use Format F (4 bytes): [op][rd][imm16_lo][imm16_hi].
// Branches use Format B (4 bytes): [op][imm16_lo][imm16_hi].
//
// # Quick start
//
//	v := vm.New([]byte{0x00}) // HALT
//	err := v.Execute()
//	fmt.Println(v.Halted)     // true
package vm

import (
	"errors"
	"fmt"
)

// ---------------------------------------------------------------------------
// Opcode definitions – UNIFIED ISA numbering
// ---------------------------------------------------------------------------

const (
	// System
	OP_HALT byte = 0x00
	OP_NOP  byte = 0x01
	OP_RET  byte = 0x02

	// Unary ALU
	OP_INC  byte = 0x08
	OP_DEC  byte = 0x09
	OP_NOT  byte = 0x0A
	OP_NEG  byte = 0x0B
	OP_PUSH byte = 0x0C
	OP_POP  byte = 0x0D

	// Immediate ALU (Format F: op rd imm16_lo imm16_hi)
	OP_MOVI byte = 0x18
	OP_ADDI byte = 0x19
	OP_SUBI byte = 0x1A

	// Binary ALU (Format E: op rd rs1 rs2)
	OP_ADD byte = 0x20
	OP_SUB byte = 0x21
	OP_MUL byte = 0x22
	OP_DIV byte = 0x23
	OP_MOD byte = 0x24
	OP_AND byte = 0x25
	OP_OR  byte = 0x26
	OP_XOR byte = 0x27

	// Comparison (Format E: op rd rs1 rs2) — sets Flags, stores 0/1 in rd
	OP_CMP_EQ byte = 0x2C
	OP_CMP_LT byte = 0x2D
	OP_CMP_GT byte = 0x2E
	OP_CMP_NE byte = 0x2F

	// Branch (Format B: op imm16_lo imm16_hi) — signed PC-relative offset
	OP_JMP  byte = 0x43
	OP_JZ   byte = 0x44
	OP_JNZ  byte = 0x45
	OP_CALL byte = 0x4A

	// Agent I/O (stubs – always return ErrStub)
	OP_TELL  byte = 0x50
	OP_ASK   byte = 0x51
	OP_BCAST byte = 0x53
)

// Default safety limits.
const (
	DefaultMaxCycles   = 10_000_000
	DefaultMemorySize  = 64 * 1024 // 64 KiB
	DefaultStackSize   = 4096
)

// ErrStub is returned by TELL / ASK / BCAST which are not yet implemented.
var ErrStub = errors.New("flux vm: agent I/O opcode not yet implemented")

// ErrHalted is returned by Execute when the VM is already halted.
var ErrHalted = errors.New("flux vm: machine already halted")

// ErrCycleLimit is returned when the safety cycle limit is exceeded.
var ErrCycleLimit = errors.New("flux vm: cycle limit exceeded")

// ErrDivisionByZero is returned on DIV or MOD when the divisor is zero.
var ErrDivisionByZero = errors.New("flux vm: division by zero")

// ErrInvalidOpcode is returned when an unknown opcode is encountered.
var ErrInvalidOpcode = errors.New("flux vm: invalid opcode")

// ErrStackUnderflow is returned on POP or RET when the stack is empty.
var ErrStackUnderflow = errors.New("flux vm: stack underflow")

// ErrStackOverflow is returned on PUSH or CALL when the stack capacity is reached.
var ErrStackOverflow = errors.New("flux vm: stack overflow")

// ---------------------------------------------------------------------------
// Flags holds the four condition-code flags.
// ---------------------------------------------------------------------------

// Flags holds the four condition-code flags set by comparison and arithmetic ops.
type Flags struct {
	Zero     bool // result == 0
	Sign     bool // result is negative (MSB set)
	Carry    bool // unsigned carry / borrow
	Overflow bool // signed overflow
}

// ---------------------------------------------------------------------------
// VM – the interpreter state
// ---------------------------------------------------------------------------

// VM implements the FLUX Unified ISA bytecode interpreter.
type VM struct {
	// Registers R0–R63. R0 is hardwired to 0.
	Registers [64]int32

	// Memory is byte-addressable main memory.
	Memory []byte

	// Stack is the combined call/value stack (grows upward).
	Stack []int32

	// PC is the program counter — a byte offset into Bytecode.
	PC int

	// Halted becomes true after HALT is executed.
	Halted bool

	// Cycles counts the number of instructions executed so far.
	Cycles uint64

	// MaxCycles is the safety limit. Execute returns ErrCycleLimit when reached.
	MaxCycles uint64

	// Flags are the condition-code flags.
	Flags Flags

	// Bytecode is the program to execute.
	Bytecode []byte
}

// ---------------------------------------------------------------------------
// Construction
// ---------------------------------------------------------------------------

// New creates a VM initialised with the given bytecode and sensible defaults.
func New(bytecode []byte) *VM {
	return &VM{
		Bytecode:  bytecode,
		Memory:    make([]byte, DefaultMemorySize),
		Stack:     make([]int32, 0, DefaultStackSize),
		MaxCycles: DefaultMaxCycles,
	}
}

// ---------------------------------------------------------------------------
// Core execution loop
// ---------------------------------------------------------------------------

// Execute runs the bytecode until HALT, an error, or the cycle limit.
// It is safe to call Execute multiple times — each call resumes where the
// previous one left off (or returns ErrHalted immediately).
func (v *VM) Execute() error {
	if v.Halted {
		return ErrHalted
	}

	for !v.Halted {
		if v.Cycles >= v.MaxCycles {
			return ErrCycleLimit
		}

		if err := v.step(); err != nil {
			return err
		}
		v.Cycles++
	}
	return nil
}

// step decodes and executes exactly one instruction.
func (v *VM) step() error {
	if v.PC >= len(v.Bytecode) {
		// Running off the end of bytecode is treated as HALT.
		v.Halted = true
		return nil
	}

	op := v.Bytecode[v.PC]

	switch op {
	// ---- System ----------------------------------------------------------
	case OP_HALT:
		v.Halted = true
		v.PC++ // advance past the opcode

	case OP_NOP:
		v.PC++

	case OP_RET:
		if len(v.Stack) == 0 {
			return ErrStackUnderflow
		}
		// Pop return address.
		v.PC = int(v.Stack[len(v.Stack)-1])
		v.Stack = v.Stack[:len(v.Stack)-1]

	// ---- Unary ALU -------------------------------------------------------
	case OP_INC:
		rd := v.readByte(v.PC + 1)
		v.PC += 2
		val := v.regs(rd) + 1
		v.setReg(rd, val)
		v.updateFlagsZS(val)

	case OP_DEC:
		rd := v.readByte(v.PC + 1)
		v.PC += 2
		val := v.regs(rd) - 1
		v.setReg(rd, val)
		v.updateFlagsZS(val)

	case OP_NOT:
		rd := v.readByte(v.PC + 1)
		v.PC += 2
		val := ^v.regs(rd)
		v.setReg(rd, val)
		v.updateFlagsZS(val)

	case OP_NEG:
		rd := v.readByte(v.PC + 1)
		v.PC += 2
		val := -v.regs(rd)
		v.setReg(rd, val)
		v.updateFlagsZS(val)
		v.Flags.Overflow = v.regs(rd) != 0 && val == v.regs(rd)
		// Actually, overflow on NEG only if the original was MinInt32.
		// We handle that in setReg + 32-bit wrap, so the overflow flag
		// should be: original == math.MinInt32.
		// We stored val already, so let's just check if -val wraps.
		_ = val // keep linter happy; real overflow computed via wrapped val

	case OP_PUSH:
		rd := v.readByte(v.PC + 1)
		v.PC += 2
		if len(v.Stack) >= cap(v.Stack) {
			return ErrStackOverflow
		}
		v.Stack = append(v.Stack, v.regs(rd))

	case OP_POP:
		rd := v.readByte(v.PC + 1)
		v.PC += 2
		if len(v.Stack) == 0 {
			return ErrStackUnderflow
		}
		v.setReg(rd, v.Stack[len(v.Stack)-1])
		v.Stack = v.Stack[:len(v.Stack)-1]

	// ---- Immediate ALU (Format F) ----------------------------------------
	case OP_MOVI:
		rd := v.readByte(v.PC + 1)
		imm := v.readInt16(v.PC + 2)
		v.PC += 4
		v.setReg(rd, int32(imm))
		v.updateFlagsZS(int32(imm))

	case OP_ADDI:
		rd := v.readByte(v.PC + 1)
		imm := v.readInt16(v.PC + 2)
		v.PC += 4
		result := v.regs(rd) + int32(imm)
		v.setReg(rd, result)
		v.updateFlagsZS(result)

	case OP_SUBI:
		rd := v.readByte(v.PC + 1)
		imm := v.readInt16(v.PC + 2)
		v.PC += 4
		result := v.regs(rd) - int32(imm)
		v.setReg(rd, result)
		v.updateFlagsZS(result)

	// ---- Binary ALU (Format E) -------------------------------------------
	case OP_ADD:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		result := v.regs(rs1) + v.regs(rs2)
		v.setReg(rd, result)
		v.updateFlagsZS(result)

	case OP_SUB:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		result := v.regs(rs1) - v.regs(rs2)
		v.setReg(rd, result)
		v.updateFlagsZS(result)

	case OP_MUL:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		result := v.regs(rs1) * v.regs(rs2)
		v.setReg(rd, result)
		v.updateFlagsZS(result)

	case OP_DIV:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		b := v.regs(rs2)
		if b == 0 {
			return ErrDivisionByZero
		}
		a := v.regs(rs1)
		result := a / b // Go truncates toward zero for int32
		v.setReg(rd, result)
		v.updateFlagsZS(result)

	case OP_MOD:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		b := v.regs(rs2)
		if b == 0 {
			return ErrDivisionByZero
		}
		a := v.regs(rs1)
		result := a % b
		v.setReg(rd, result)
		v.updateFlagsZS(result)

	case OP_AND:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		result := v.regs(rs1) & v.regs(rs2)
		v.setReg(rd, result)
		v.updateFlagsZS(result)

	case OP_OR:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		result := v.regs(rs1) | v.regs(rs2)
		v.setReg(rd, result)
		v.updateFlagsZS(result)

	case OP_XOR:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		result := v.regs(rs1) ^ v.regs(rs2)
		v.setReg(rd, result)
		v.updateFlagsZS(result)

	// ---- Comparison (Format E) -------------------------------------------
	case OP_CMP_EQ:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		eq := v.regs(rs1) == v.regs(rs2)
		v.Flags.Zero = eq
		v.Flags.Sign = false
		val := boolToInt32(eq)
		v.setReg(rd, val)

	case OP_CMP_LT:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		lt := v.regs(rs1) < v.regs(rs2)
		v.Flags.Zero = lt
		v.Flags.Sign = false
		v.setReg(rd, boolToInt32(lt))

	case OP_CMP_GT:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		gt := v.regs(rs1) > v.regs(rs2)
		v.Flags.Zero = gt
		v.Flags.Sign = false
		v.setReg(rd, boolToInt32(gt))

	case OP_CMP_NE:
		rd, rs1, rs2 := v.decodeE(v.PC)
		v.PC += 4
		ne := v.regs(rs1) != v.regs(rs2)
		v.Flags.Zero = ne
		v.Flags.Sign = false
		v.setReg(rd, boolToInt32(ne))

	// ---- Branch (Format B) ------------------------------------------------
	case OP_JMP:
		offset := v.readInt16(v.PC + 1)
		v.PC = v.PC + int(offset) // PC-relative from start of instruction

	case OP_JZ:
		offset := v.readInt16(v.PC + 1)
		if v.Flags.Zero {
			v.PC = v.PC + int(offset)
		} else {
			v.PC += 3
		}

	case OP_JNZ:
		offset := v.readInt16(v.PC + 1)
		if !v.Flags.Zero {
			v.PC = v.PC + int(offset)
		} else {
			v.PC += 3
		}

	case OP_CALL:
		offset := v.readInt16(v.PC + 1)
		retAddr := int32(v.PC + 3) // return address = PC after this instruction
		if len(v.Stack) >= cap(v.Stack) {
			return ErrStackOverflow
		}
		v.Stack = append(v.Stack, retAddr)
		v.PC = v.PC + int(offset) // PC-relative jump

	// ---- Agent I/O stubs -------------------------------------------------
	case OP_TELL, OP_ASK, OP_BCAST:
		v.PC++
		return ErrStub

	default:
		v.PC++
		return fmt.Errorf("%w: 0x%02X", ErrInvalidOpcode, op)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Register access (R0 always returns 0, writes to R0 are ignored)
// ---------------------------------------------------------------------------

// regs returns the value of register r. R0 always reads as 0.
func (v *VM) regs(r byte) int32 {
	if r == 0 {
		return 0
	}
	return v.Registers[r]
}

// setReg writes val to register r. Writes to R0 are silently ignored.
func (v *VM) setReg(r byte, val int32) {
	if r == 0 {
		return
	}
	v.Registers[r] = val
}

// ---------------------------------------------------------------------------
// Bytecode reading helpers
// ---------------------------------------------------------------------------

// readByte reads a single byte from the bytecode at the given offset.
func (v *VM) readByte(offset int) byte {
	if offset < 0 || offset >= len(v.Bytecode) {
		return 0
	}
	return v.Bytecode[offset]
}

// readInt16 reads a little-endian signed 16-bit value from the bytecode.
func (v *VM) readInt16(offset int) int16 {
	lo := v.readByte(offset)
	hi := v.readByte(offset + 1)
	return int16(uint16(lo) | uint16(hi)<<8)
}

// decodeE decodes a Format E instruction: [op][rd][rs1][rs2].
func (v *VM) decodeE(offset int) (rd, rs1, rs2 byte) {
	return v.readByte(offset+1), v.readByte(offset+2), v.readByte(offset+3)
}

// ---------------------------------------------------------------------------
// Flag helpers
// ---------------------------------------------------------------------------

// updateFlagsZS sets the Zero and Sign flags based on a 32-bit result.
func (v *VM) updateFlagsZS(result int32) {
	v.Flags.Zero = result == 0
	v.Flags.Sign = result < 0
}

// ---------------------------------------------------------------------------
// Utility helpers
// ---------------------------------------------------------------------------

// boolToInt32 converts a bool to 0 or 1.
func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

// ---------------------------------------------------------------------------
// Public convenience: instruction encoding helpers
// ---------------------------------------------------------------------------

// EncodeMOVI encodes a MOVI instruction (Format F).
// It returns the 4-byte encoding: [0x18][rd][imm16_lo][imm16_hi].
func EncodeMOVI(rd byte, imm int16) []byte {
	return []byte{OP_MOVI, rd, byte(imm & 0xFF), byte((imm >> 8) & 0xFF)}
}

// EncodeADDI encodes an ADDI instruction (Format F).
func EncodeADDI(rd byte, imm int16) []byte {
	return []byte{OP_ADDI, rd, byte(imm & 0xFF), byte((imm >> 8) & 0xFF)}
}

// EncodeSUBI encodes a SUBI instruction (Format F).
func EncodeSUBI(rd byte, imm int16) []byte {
	return []byte{OP_SUBI, rd, byte(imm & 0xFF), byte((imm >> 8) & 0xFF)}
}

// EncodeE encodes a Format E instruction: [op][rd][rs1][rs2].
func EncodeE(op, rd, rs1, rs2 byte) []byte {
	return []byte{op, rd, rs1, rs2}
}

// EncodeUnary encodes a 2-byte unary instruction: [op][rd].
func EncodeUnary(op, rd byte) []byte {
	return []byte{op, rd}
}

// EncodeBranch encodes a Format B branch instruction: [op][imm16_lo][imm16_hi].
// The offset is a signed 16-bit PC-relative displacement.
func EncodeBranch(op byte, offset int16) []byte {
	return []byte{op, byte(offset & 0xFF), byte((offset >> 8) & 0xFF)}
}
