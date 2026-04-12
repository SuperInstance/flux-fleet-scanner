# FLUX VM — Go Bytecode Interpreter

This package (`pkg/vm`) implements a register-based bytecode virtual machine that executes the **FLUX Unified ISA**. It is the Go counterpart of the Python-based VM in the `flux-runtime` project, intended for use in the `greenhorn-runtime` agent framework.

## Architecture

| Component     | Description                                  |
|---------------|----------------------------------------------|
| **Registers** | 64 general-purpose 32-bit signed integers (R0–R63) |
| **Memory**    | Byte-addressable linear memory (default 64 KiB) |
| **Stack**     | Combined call/value stack (grows upward, default cap 4096 entries) |
| **PC**        | Program counter — byte offset into the bytecode array |
| **Flags**     | Four condition-code flags: Zero, Sign, Carry, Overflow |

### Register Conventions

- **R0** is hardwired to zero. All writes to R0 are silently ignored, and it always reads as `0`.
- R1–R63 are general-purpose with no ABI-enforced roles (though callers may adopt their own conventions).
- Calling convention: `CALL` pushes the return address (PC after the CALL instruction) onto the stack. `RET` pops it and jumps.

### Memory Model

- Memory is a flat `[]byte` slice, addressable from `0` to `len(Memory)-1`.
- No memory instructions are implemented in this initial version; memory access can be added via load/store opcodes in a future extension.

## Supported Opcodes

### System (1 byte)

| Opcode | Hex  | Operands | Description        |
|--------|------|----------|--------------------|
| HALT   | 0x00 | —        | Stop execution     |
| NOP    | 0x01 | —        | No operation       |
| RET    | 0x02 | —        | Pop PC from stack  |

### Unary ALU (2 bytes: `[op][rd]`)

| Opcode | Hex  | Description                              |
|--------|------|------------------------------------------|
| INC    | 0x08 | rd ← rd + 1                              |
| DEC    | 0x09 | rd ← rd − 1                              |
| NOT    | 0x0A | rd ← bitwise NOT rd                      |
| NEG    | 0x0B | rd ← −rd                                 |
| PUSH   | 0x0C | Push rd onto stack                       |
| POP    | 0x0D | Pop stack top into rd                    |

### Immediate ALU — Format F (4 bytes: `[op][rd][imm16_lo][imm16_hi]`)

| Opcode | Hex  | Description                            |
|--------|------|----------------------------------------|
| MOVI   | 0x18 | rd ← sign-extended imm16               |
| ADDI   | 0x19 | rd ← rd + sign-extended imm16          |
| SUBI   | 0x1A | rd ← rd − sign-extended imm16          |

`imm16` is stored **little-endian** as two bytes.

### Binary ALU — Format E (4 bytes: `[op][rd][rs1][rs2]`)

| Opcode | Hex  | Description           |
|--------|------|-----------------------|
| ADD    | 0x20 | rd ← rs1 + rs2        |
| SUB    | 0x21 | rd ← rs1 − rs2        |
| MUL    | 0x22 | rd ← rs1 × rs2        |
| DIV    | 0x23 | rd ← rs1 / rs2 (trunc toward zero) |
| MOD    | 0x24 | rd ← rs1 % rs2        |
| AND    | 0x25 | rd ← rs1 & rs2        |
| OR     | 0x26 | rd ← rs1 \| rs2       |
| XOR    | 0x27 | rd ← rs1 ^ rs2        |

### Comparison — Format E (4 bytes: `[op][rd][rs1][rs2]`)

Compares rs1 and rs2, stores `1` in rd if the condition is true (else `0`), and updates `Flags.Zero`.

| Opcode | Hex  | Condition        |
|--------|------|------------------|
| CMP_EQ | 0x2C | rs1 == rs2       |
| CMP_LT | 0x2D | rs1 < rs2        |
| CMP_GT | 0x2E | rs1 > rs2        |
| CMP_NE | 0x2F | rs1 != rs2       |

### Branch — Format B (3 bytes: `[op][imm16_lo][imm16_hi]`)

Offset is a **signed 16-bit PC-relative** displacement from the start of the branch instruction.

| Opcode | Hex  | Description                                 |
|--------|------|---------------------------------------------|
| JMP    | 0x43 | PC ← PC + offset (unconditional)            |
| JZ     | 0x44 | PC ← PC + offset (if Flags.Zero)            |
| JNZ    | 0x45 | PC ← PC + offset (if !Flags.Zero)           |
| CALL   | 0x4A | Push PC+3, then PC ← PC + offset            |

### Agent I/O — Stubs (1 byte)

| Opcode | Hex  | Description                               |
|--------|------|-------------------------------------------|
| TELL   | 0x50 | Stub — always returns `ErrStub`           |
| ASK    | 0x51 | Stub — always returns `ErrStub`           |
| BCAST  | 0x53 | Stub — always returns `ErrStub`           |

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/greenhorn-runtime/pkg/vm"
)

func main() {
    // Build a simple program: load 10, add 20, halt
    prog := vm.Assemble(
        vm.EncodeMOVI(1, 10),
        vm.EncodeADDI(1, 20),
        []byte{vm.OP_HALT},
    )

    v := vm.New(prog)
    if err := v.Execute(); err != nil {
        fmt.Println("error:", err)
        return
    }

    fmt.Printf("R1 = %d\n", v.Registers[1]) // R1 = 30
}
```

## Encoding Helpers

The package provides helper functions for building bytecode:

| Function | Signature | Description |
|----------|-----------|-------------|
| `EncodeMOVI` | `(rd byte, imm int16) []byte` | Format F: load immediate |
| `EncodeADDI` | `(rd byte, imm int16) []byte` | Format F: add immediate |
| `EncodeSUBI` | `(rd byte, imm int16) []byte` | Format F: sub immediate |
| `EncodeE` | `(op, rd, rs1, rs2 byte) []byte` | Format E: binary ALU |
| `EncodeUnary` | `(op, rd byte) []byte` | Unary: INC/DEC/NOT/NEG/PUSH/POP |
| `EncodeBranch` | `(op byte, offset int16) []byte` | Format B: JMP/JZ/JNZ/CALL |

## Safety Features

- **Cycle limit**: `MaxCycles` (default 10 million) prevents infinite loops. Returns `ErrCycleLimit`.
- **Stack bounds**: Push beyond capacity returns `ErrStackOverflow`; pop from empty returns `ErrStackUnderflow`.
- **Division safety**: Division or modulo by zero returns `ErrDivisionByZero`.
- **Invalid opcodes**: Unknown opcodes return `ErrInvalidOpcode`.
- **Run-off-end**: If PC reaches the end of bytecode without HALT, execution stops gracefully (treated as HALT).

## Arithmetic Semantics

All arithmetic uses **32-bit signed integer wrapping** (two's complement), matching Go's `int32` behavior:
- Addition/subtraction/multiplication wrap naturally at the 32-bit boundary.
- Division truncates toward zero (Go semantics, matching C99/POSIX).
- Modulo result has the same sign as the dividend.

## Relationship to flux-runtime (Python)

This Go VM implements the same UNIFIED ISA as the Python `flux-runtime` VM. The two implementations share:

- The same opcode numbering and instruction encodings.
- The same register and memory model.
- The same comparison and flag semantics.

Differences:
- The Python VM is interpreted; this Go VM is also interpreted but leverages Go's type safety and performance.
- Agent I/O opcodes (TELL/ASK/BCAST) are stubs here; in flux-runtime they connect to the agent message bus.
- This VM provides encoding helpers and explicit error types for easier integration.

## Running Tests

```bash
cd /path/to/greenhorn-runtime
go test -v ./pkg/vm/
```

All 49 tests cover every opcode, edge cases (division by zero, stack underflow, R0 immutability), overflow wrapping, and the cycle safety limit.
