package vm

import (
        "math"
        "testing"
)

// ---------------------------------------------------------------------------
// Helper: assemble bytecode from multiple encoded instructions
// ---------------------------------------------------------------------------

func assemble(parts ...[]byte) []byte {
        total := 0
        for _, p := range parts {
                total += len(p)
        }
        out := make([]byte, 0, total)
        for _, p := range parts {
                out = append(out, p...)
        }
        return out
}

// ---------------------------------------------------------------------------
// System instructions
// ---------------------------------------------------------------------------

func TestHALT_StopsExecution(t *testing.T) {
        prog := assemble(EncodeMOVI(1, 42), []byte{OP_HALT})
        v := New(prog)

        err := v.Execute()
        if err != nil {
                t.Fatalf("unexpected error: %v", err)
        }
        if !v.Halted {
                t.Fatal("expected Halted=true")
        }
        if v.Registers[1] != 42 {
                t.Fatalf("expected R1=42, got %d", v.Registers[1])
        }
}

func TestNOP_DoesNothing(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 10),
                []byte{OP_NOP},
                []byte{OP_NOP},
                EncodeADDI(1, 5),
                []byte{OP_HALT},
        )
        v := New(prog)

        err := v.Execute()
        if err != nil {
                t.Fatalf("unexpected error: %v", err)
        }
        if v.Registers[1] != 15 {
                t.Fatalf("expected R1=15, got %d", v.Registers[1])
        }
}

func TestExecute_AlreadyHalted(t *testing.T) {
        prog := []byte{OP_HALT}
        v := New(prog)
        v.Execute()

        err := v.Execute()
        if err != ErrHalted {
                t.Fatalf("expected ErrHalted, got %v", err)
        }
}

// ---------------------------------------------------------------------------
// MOVI / ADDI / SUBI
// ---------------------------------------------------------------------------

func TestMOVI_PositiveImmediate(t *testing.T) {
        prog := assemble(EncodeMOVI(5, 12345), []byte{OP_HALT})
        v := New(prog)
        v.Execute()

        if v.Registers[5] != 12345 {
                t.Fatalf("expected R5=12345, got %d", v.Registers[5])
        }
}

func TestMOVI_NegativeImmediate(t *testing.T) {
        prog := assemble(EncodeMOVI(3, -1000), []byte{OP_HALT})
        v := New(prog)
        v.Execute()

        if v.Registers[3] != -1000 {
                t.Fatalf("expected R3=-1000, got %d", v.Registers[3])
        }
}

func TestMOVI_ZeroImmediate(t *testing.T) {
        prog := assemble(EncodeMOVI(7, 0), []byte{OP_HALT})
        v := New(prog)
        v.Execute()

        if v.Registers[7] != 0 {
                t.Fatalf("expected R7=0, got %d", v.Registers[7])
        }
}

func TestADDI(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 100),
                EncodeADDI(1, 23),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[1] != 123 {
                t.Fatalf("expected R1=123, got %d", v.Registers[1])
        }
}

func TestSUBI(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 100),
                EncodeSUBI(1, 30),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[1] != 70 {
                t.Fatalf("expected R1=70, got %d", v.Registers[1])
        }
}

// ---------------------------------------------------------------------------
// Binary ALU: ADD, SUB, MUL, DIV, MOD
// ---------------------------------------------------------------------------

func TestADD(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 10),
                EncodeMOVI(2, 25),
                EncodeE(OP_ADD, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 35 {
                t.Fatalf("expected R3=35, got %d", v.Registers[3])
        }
}

func TestSUB(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 100),
                EncodeMOVI(2, 37),
                EncodeE(OP_SUB, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 63 {
                t.Fatalf("expected R3=63, got %d", v.Registers[3])
        }
}

func TestMUL(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 7),
                EncodeMOVI(2, 6),
                EncodeE(OP_MUL, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 42 {
                t.Fatalf("expected R3=42, got %d", v.Registers[3])
        }
}

func TestDIV(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 100),
                EncodeMOVI(2, 7),
                EncodeE(OP_DIV, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 14 {
                t.Fatalf("expected R3=14, got %d", v.Registers[3])
        }
}

func TestDIV_Negative(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, -100),
                EncodeMOVI(2, 7),
                EncodeE(OP_DIV, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        // Go truncates toward zero: -100/7 = -14
        if v.Registers[3] != -14 {
                t.Fatalf("expected R3=-14, got %d", v.Registers[3])
        }
}

func TestDIV_ByZero(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 100),
                EncodeMOVI(2, 0),
                EncodeE(OP_DIV, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)

        err := v.Execute()
        if err != ErrDivisionByZero {
                t.Fatalf("expected ErrDivisionByZero, got %v", err)
        }
}

func TestMOD(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 100),
                EncodeMOVI(2, 7),
                EncodeE(OP_MOD, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 2 {
                t.Fatalf("expected R3=2, got %d", v.Registers[3])
        }
}

func TestMOD_ByZero(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 100),
                EncodeMOVI(2, 0),
                EncodeE(OP_MOD, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)

        err := v.Execute()
        if err != ErrDivisionByZero {
                t.Fatalf("expected ErrDivisionByZero, got %v", err)
        }
}

// ---------------------------------------------------------------------------
// INC / DEC
// ---------------------------------------------------------------------------

func TestINC(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 41),
                EncodeUnary(OP_INC, 1),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[1] != 42 {
                t.Fatalf("expected R1=42, got %d", v.Registers[1])
        }
}

func TestDEC(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 43),
                EncodeUnary(OP_DEC, 1),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[1] != 42 {
                t.Fatalf("expected R1=42, got %d", v.Registers[1])
        }
}

// ---------------------------------------------------------------------------
// NEG / NOT
// ---------------------------------------------------------------------------

func TestNEG(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 50),
                EncodeUnary(OP_NEG, 1),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[1] != -50 {
                t.Fatalf("expected R1=-50, got %d", v.Registers[1])
        }
}

func TestNOT(t *testing.T) {
        prog := assemble(
                // Load 0x00FF into R1
                EncodeMOVI(1, 0x00FF),
                EncodeUnary(OP_NOT, 1),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        // NOT 0x00FF = 0xFFFFFF00 which is -256 as int32
        expected := int32(-256)
        if v.Registers[1] != expected {
                t.Fatalf("expected R1=%d (0x%08X), got %d (0x%08X)",
                        expected, uint32(expected), v.Registers[1], uint32(v.Registers[1]))
        }
}

// ---------------------------------------------------------------------------
// AND / OR / XOR
// ---------------------------------------------------------------------------

func TestAND(t *testing.T) {
        prog := assemble(
                // 0x0F0F & 0x00F0 = 0x000F... let's use values that fit in int16
                // Use 0x0F0F and 0x0F00
                EncodeMOVI(1, 0x0F0F),
                EncodeMOVI(2, 0x0F00),
                EncodeE(OP_AND, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 0x0F00 {
                t.Fatalf("expected R3=0x0F00, got 0x%04X", v.Registers[3])
        }
}

func TestOR(t *testing.T) {
        prog := assemble(
                // 0x00F0 | 0x000F = 0x00FF
                EncodeMOVI(1, 0x00F0),
                EncodeMOVI(2, 0x000F),
                EncodeE(OP_OR, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 0x00FF {
                t.Fatalf("expected R3=0x00FF, got 0x%04X", v.Registers[3])
        }
}

func TestXOR(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 0xFF),
                EncodeMOVI(2, 0xFF),
                EncodeE(OP_XOR, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 0 {
                t.Fatalf("expected R3=0, got %d", v.Registers[3])
        }
}

// ---------------------------------------------------------------------------
// CMP_EQ / CMP_LT / CMP_GT / CMP_NE
// ---------------------------------------------------------------------------

func TestCMP_EQ_True(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 42),
                EncodeMOVI(2, 42),
                EncodeE(OP_CMP_EQ, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 1 {
                t.Fatalf("expected R3=1, got %d", v.Registers[3])
        }
        if !v.Flags.Zero {
                t.Fatal("expected Flags.Zero=true")
        }
}

func TestCMP_EQ_False(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 10),
                EncodeMOVI(2, 20),
                EncodeE(OP_CMP_EQ, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 0 {
                t.Fatalf("expected R3=0, got %d", v.Registers[3])
        }
        if v.Flags.Zero {
                t.Fatal("expected Flags.Zero=false")
        }
}

func TestCMP_LT_True(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 5),
                EncodeMOVI(2, 10),
                EncodeE(OP_CMP_LT, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 1 {
                t.Fatalf("expected R3=1, got %d", v.Registers[3])
        }
        if !v.Flags.Zero {
                t.Fatal("expected Flags.Zero=true (comparison result is non-zero/true)")
        }
}

func TestCMP_LT_False(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 20),
                EncodeMOVI(2, 10),
                EncodeE(OP_CMP_LT, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 0 {
                t.Fatalf("expected R3=0, got %d", v.Registers[3])
        }
}

func TestCMP_GT_True(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 20),
                EncodeMOVI(2, 10),
                EncodeE(OP_CMP_GT, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 1 {
                t.Fatalf("expected R3=1, got %d", v.Registers[3])
        }
}

func TestCMP_NE_True(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 1),
                EncodeMOVI(2, 2),
                EncodeE(OP_CMP_NE, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 1 {
                t.Fatalf("expected R3=1, got %d", v.Registers[3])
        }
}

func TestCMP_NE_False(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 99),
                EncodeMOVI(2, 99),
                EncodeE(OP_CMP_NE, 3, 1, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 0 {
                t.Fatalf("expected R3=0, got %d", v.Registers[3])
        }
}

// ---------------------------------------------------------------------------
// JMP forward and backward
// ---------------------------------------------------------------------------

func TestJMP_Forward(t *testing.T) {
        // JMP +3 (skip over 2 instructions: MOVI sets R1=0, second MOVI sets R1=999)
        // Layout:
        //   0: JMP +7        -> jumps to offset 7 (HALT)
        //   3: MOVI R1 999   -> should be skipped
        //   7: HALT
        prog := assemble(
                EncodeBranch(OP_JMP, 7), // offset 0-2
                EncodeMOVI(1, 999),      // offset 3-6
                []byte{OP_HALT},         // offset 7
        )
        v := New(prog)
        v.Execute()

        if v.Registers[1] != 0 {
                t.Fatalf("expected R1=0 (skipped), got %d", v.Registers[1])
        }
}

func TestJMP_Backward(t *testing.T) {
        // Simple loop: count down from 3 to 0
        // Layout:
        //   0: MOVI R1 3        (4 bytes)
        //   4: MOVI R2 1        (4 bytes)
        //   8: DEC  R1           (2 bytes)
        //  10: JZ   +5 -> HALT  (3 bytes, offset 13)
        //  13: JMP  -5 -> DEC   (3 bytes, offset 8)
        //  16: HALT              (1 byte)
        prog := assemble(
                EncodeMOVI(1, 3),          // 0-3
                EncodeMOVI(2, 1),          // 4-7
                EncodeUnary(OP_DEC, 1),    // 8-9
                EncodeBranch(OP_JZ, 5),    // 10-12: jump to 15 (should be 16-1=15... let me recalc)
                []byte{OP_HALT},           // 13
        )
        // Wait, the offset for JZ should point to HALT at offset 13.
        // JZ is at offset 10, so offset = 13 - 10 = 3.
        // Let me fix:
        prog = assemble(
                EncodeMOVI(1, 3),          // 0-3
                EncodeUnary(OP_DEC, 1),    // 4-5
                EncodeBranch(OP_JZ, 4),    // 6-8: if zero, jump to 10 (HALT)
                EncodeBranch(OP_JMP, -6),  // 9-11: jump back to offset 3 (DEC)
                []byte{OP_HALT},           // 12
        )
        // Let me verify: DEC at 4, JZ at 6 (offset 3 bytes), JMP at 9 (offset 3 bytes), HALT at 12.
        // JZ at 6: if zero flag, PC = 6 + 4 = 10... that's wrong.
        // JZ at 6: want to jump to 12. offset = 12 - 6 = 6.
        // JMP at 9: want to jump to 4. offset = 4 - 9 = -5.
        prog = assemble(
                EncodeMOVI(1, 3),          // 0-3
                EncodeUnary(OP_DEC, 1),    // 4-5
                EncodeBranch(OP_JZ, 6),    // 6-8: if zero, jump to 12 (HALT). PC=6, 6+6=12. ✓
                EncodeBranch(OP_JMP, -5),  // 9-11: jump to 4 (DEC). PC=9, 9+(-5)=4. ✓
                []byte{OP_HALT},           // 12
        )
        v := New(prog)
        v.Execute()

        if v.Registers[1] != 0 {
                t.Fatalf("expected R1=0 after loop, got %d", v.Registers[1])
        }
}

// ---------------------------------------------------------------------------
// JZ / JNZ
// ---------------------------------------------------------------------------

func TestJZ_Taken(t *testing.T) {
        // Set up Flags.Zero = true via CMP_EQ with equal values, then JZ to skip.
        // Layout:
        //   0: MOVI R1 10        (4 bytes)
        //   4: MOVI R2 10        (4 bytes)
        //   8: CMP_EQ R3 R1 R2   (4 bytes) — sets Zero=true
        //  12: JZ   +3           (3 bytes) -> jump to 15 (HALT)
        //  15: MOVI R4 999       (4 bytes) — should be skipped
        //  19: HALT               (1 byte)
        prog := assemble(
                EncodeMOVI(1, 10),              // 0-3
                EncodeMOVI(2, 10),              // 4-7
                EncodeE(OP_CMP_EQ, 3, 1, 2),   // 8-11
                EncodeBranch(OP_JZ, 4),         // 12-14: jump to 16 (HALT). PC=12, 12+4=16. Wait...
                EncodeMOVI(4, 999),             // 15-18
                []byte{OP_HALT},                // 19
        )
        // JZ at offset 12, want to jump to 19. offset = 19 - 12 = 7.
        prog = assemble(
                EncodeMOVI(1, 10),              // 0-3
                EncodeMOVI(2, 10),              // 4-7
                EncodeE(OP_CMP_EQ, 3, 1, 2),   // 8-11
                EncodeBranch(OP_JZ, 7),         // 12-14: jump to 19. 12+7=19. ✓
                EncodeMOVI(4, 999),             // 15-18
                []byte{OP_HALT},                // 19
        )
        v := New(prog)
        v.Execute()

        if v.Registers[4] != 0 {
                t.Fatalf("expected R4=0 (skipped), got %d", v.Registers[4])
        }
}

func TestJZ_NotTaken(t *testing.T) {
        // CMP_EQ with different values -> Zero=false, JZ should fall through.
        prog := assemble(
                EncodeMOVI(1, 10),              // 0-3
                EncodeMOVI(2, 20),              // 4-7
                EncodeE(OP_CMP_EQ, 3, 1, 2),   // 8-11
                EncodeBranch(OP_JZ, 4),         // 12-14: should NOT jump
                EncodeMOVI(4, 888),             // 15-18
                []byte{OP_HALT},                // 19
        )
        v := New(prog)
        v.Execute()

        if v.Registers[4] != 888 {
                t.Fatalf("expected R4=888 (not skipped), got %d", v.Registers[4])
        }
}

func TestJNZ_Taken(t *testing.T) {
        // CMP_NE with different values -> Zero=true, JNZ should jump.
        prog := assemble(
                EncodeMOVI(1, 10),              // 0-3
                EncodeMOVI(2, 20),              // 4-7
                EncodeE(OP_CMP_NE, 3, 1, 2),   // 8-11: Zero=true (result 1)
                EncodeBranch(OP_JNZ, 7),        // 12-14: jump to 19. Zero=true, so JNZ checks !Zero=false... WAIT.
                // JNZ jumps when Zero is FALSE. CMP_NE sets Zero=true for inequality.
                // So we need a different approach: use CMP_EQ with unequal values -> Zero=false.
                []byte{OP_HALT},                // 15 -- unused
        )
        // Redo: CMP_EQ with unequal -> Zero=false. JNZ taken.
        prog = assemble(
                EncodeMOVI(1, 10),              // 0-3
                EncodeMOVI(2, 20),              // 4-7
                EncodeE(OP_CMP_EQ, 3, 1, 2),   // 8-11: Zero=false
                EncodeBranch(OP_JNZ, 7),        // 12-14: Zero=false, so jump to 19
                EncodeMOVI(4, 999),             // 15-18: should be skipped
                []byte{OP_HALT},                // 19
        )
        v := New(prog)
        v.Execute()

        if v.Registers[4] != 0 {
                t.Fatalf("expected R4=0 (skipped), got %d", v.Registers[4])
        }
}

func TestJNZ_NotTaken(t *testing.T) {
        // CMP_EQ with equal -> Zero=true. JNZ should NOT jump.
        prog := assemble(
                EncodeMOVI(1, 10),              // 0-3
                EncodeMOVI(2, 10),              // 4-7
                EncodeE(OP_CMP_EQ, 3, 1, 2),   // 8-11: Zero=true
                EncodeBranch(OP_JNZ, 7),        // 12-14: Zero=true, so NOT taken
                EncodeMOVI(4, 777),             // 15-18
                []byte{OP_HALT},                // 19
        )
        v := New(prog)
        v.Execute()

        if v.Registers[4] != 777 {
                t.Fatalf("expected R4=777 (not skipped), got %d", v.Registers[4])
        }
}

// ---------------------------------------------------------------------------
// CALL / RET
// ---------------------------------------------------------------------------

func TestCALL_PushesReturnAddress_RET_JumpsBack(t *testing.T) {
        // Layout:
        //   0: MOVI R1 10          (4 bytes)
        //   4: MOVI R2 20          (4 bytes)
        //   8: CALL +8             (3 bytes) -> jump to 16. Push return addr 11.
        //  11: MOVI R1 999         (4 bytes) — this runs after RET
        //  15: HALT                 (1 byte)
        //  16: ADD R3 R1 R2        (4 bytes) — subroutine: R3 = 10+20 = 30
        //  20: RET                  (1 byte)
        prog := assemble(
                EncodeMOVI(1, 10),         // 0-3
                EncodeMOVI(2, 20),         // 4-7
                EncodeBranch(OP_CALL, 8),  // 8-10: PC=8, jump to 8+8=16. Push ret=11.
                EncodeMOVI(1, 999),        // 11-14
                []byte{OP_HALT},           // 15
                EncodeE(OP_ADD, 3, 1, 2),  // 16-19
                []byte{OP_RET},            // 20
        )
        v := New(prog)
        v.Execute()

        // After CALL returns, R1 should be overwritten to 999.
        if v.Registers[1] != 999 {
                t.Fatalf("expected R1=999 after return, got %d", v.Registers[1])
        }
        // R3 should be 30 (computed in subroutine before RET).
        if v.Registers[3] != 30 {
                t.Fatalf("expected R3=30, got %d", v.Registers[3])
        }
        if !v.Halted {
                t.Fatal("expected Halted=true")
        }
}

func TestRET_EmptyStack(t *testing.T) {
        prog := []byte{OP_RET}
        v := New(prog)

        err := v.Execute()
        if err != ErrStackUnderflow {
                t.Fatalf("expected ErrStackUnderflow, got %v", err)
        }
}

// ---------------------------------------------------------------------------
// PUSH / POP
// ---------------------------------------------------------------------------

func TestPUSH_POP(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, 42),
                EncodeUnary(OP_PUSH, 1),
                EncodeMOVI(1, 0), // clobber R1
                EncodeUnary(OP_POP, 2),
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[2] != 42 {
                t.Fatalf("expected R2=42, got %d", v.Registers[2])
        }
}

func TestPOP_EmptyStack(t *testing.T) {
        prog := assemble(
                EncodeUnary(OP_POP, 1),
                []byte{OP_HALT},
        )
        v := New(prog)

        err := v.Execute()
        if err != ErrStackUnderflow {
                t.Fatalf("expected ErrStackUnderflow, got %v", err)
        }
}

// ---------------------------------------------------------------------------
// R0 immutability
// ---------------------------------------------------------------------------

func TestR0_Immutable(t *testing.T) {
        prog := assemble(
                EncodeMOVI(0, 42),          // write to R0 — should be ignored
                EncodeE(OP_ADD, 0, 1, 2),   // write to R0 — should be ignored
                EncodeUnary(OP_INC, 0),      // write to R0 — should be ignored
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[0] != 0 {
                t.Fatalf("expected R0=0, got %d", v.Registers[0])
        }
}

// ---------------------------------------------------------------------------
// TELL / ASK / BCAST stubs
// ---------------------------------------------------------------------------

func TestTELL_ReturnsError(t *testing.T) {
        prog := []byte{OP_TELL}
        v := New(prog)

        err := v.Execute()
        if err != ErrStub {
                t.Fatalf("expected ErrStub, got %v", err)
        }
}

func TestASK_ReturnsError(t *testing.T) {
        prog := []byte{OP_ASK}
        v := New(prog)

        err := v.Execute()
        if err != ErrStub {
                t.Fatalf("expected ErrStub, got %v", err)
        }
}

func TestBCAST_ReturnsError(t *testing.T) {
        prog := []byte{OP_BCAST}
        v := New(prog)

        err := v.Execute()
        if err != ErrStub {
                t.Fatalf("expected ErrStub, got %v", err)
        }
}

// ---------------------------------------------------------------------------
// Invalid opcode
// ---------------------------------------------------------------------------

func TestInvalidOpcode(t *testing.T) {
        prog := []byte{0xFF} // not a valid opcode
        v := New(prog)

        err := v.Execute()
        if err == nil {
                t.Fatal("expected error for invalid opcode")
        }
}

// ---------------------------------------------------------------------------
// Cycle limit
// ---------------------------------------------------------------------------

func TestCycleLimit_Enforced(t *testing.T) {
        // Infinite loop: JMP -2 (back to itself).
        // Layout:
        //   0: JMP -2 (back to offset -2... that's before the program)
        //   Let's do: JMP 0 at offset 0.
        // Actually JMP offset is from start of instruction.
        // JMP 0 at offset 0: PC = 0 + 0 = 0. Infinite loop.
        prog := EncodeBranch(OP_JMP, 0)
        v := New(prog)
        v.MaxCycles = 1000

        err := v.Execute()
        if err != ErrCycleLimit {
                t.Fatalf("expected ErrCycleLimit, got %v", err)
        }
        if v.Cycles != 1000 {
                t.Fatalf("expected 1000 cycles, got %d", v.Cycles)
        }
}

// ---------------------------------------------------------------------------
// Arithmetic overflow wraps
// ---------------------------------------------------------------------------

func TestArithmeticOverflow_Wraps(t *testing.T) {
        // Build MaxInt32 by shifting: MOVI 16384, then MUL by 16384 gives 2^28,
        // then MUL by 8 gives 2^31. But that's MinInt32... Let's use a simpler approach.
        // MaxInt32 = 2^31 - 1 = 0x7FFFFFFF
        // Use: load 32767 (0x7FFF), shift by adding to itself appropriately.
        // Simpler: just use two 16-bit halves.
        // R1 = 32767, R2 = R1 + R1 + 1 = 65535, R3 = R2 + R2 + 1 = 131071...
        // Actually the simplest: load 32767, store. Then: R1 * 65536 + 32767 = MaxInt32.
        // 32767 * 65536 = 2147418112. + 32767 = 2147450879... that's not MaxInt32.
        // MaxInt32 = 0x7FFFFFFF = 32767 * 65536 + 65535.
        // OR: just compute -(MinInt32 + 1) = -(-2^31 + 1) = 2^31 - 1.
        // Load -1, negate -> 1. Load 1, NEG -> -1. Hard to build MaxInt32.
        // Easiest: MOVI R1, -1 -> NOT R1 -> 0 = MaxUint32... no, NOT(-1) = 0.
        // MOVI R1, 1; NEG R1 -> -1; NOT R1 -> 0. That's not useful.
        // MOVI R1, 0; DEC R1 -> -1; NOT R1 -> 0.
        // Actually: MOVI R1, 0; NEG R1 -> 0 (can't neg 0). Hmm.
        // Let's try: MOVI R1, -1 -> R1 = -1 = 0xFFFFFFFF; DEC R1 -> -2 = 0xFFFFFFFE.
        // RSHIFT by 1? We don't have shift. OK, let me just use MOVI with a value
        // that we can compute. Actually the issue is that MOVI only takes int16.
        // Simplest approach: use arithmetic to build MaxInt32.
        // R1 = 32767 (0x7FFF)
        // R2 = R1 << 16 = R1 * 65536... but 65536 doesn't fit in int16.
        // Use: R2 = 1; repeated doubling: MUL R2 R2 R2 16 times? That's 16 instructions.
        // Or: just MOVI -1, then the result of adding 1 to it.
        // R1 = -1 (0xFFFFFFFF); R2 = 1; ADD R3 R1 R2 = 0 (wraps). That's not MaxInt32.
        // MaxInt32 = 0x7FFFFFFF. NOT(0x80000000). But we can't easily make 0x80000000.
        // Actually: MOVI R1, -1 (0xFFFFFFFF), NOT R1 -> 0. MOVI R1, -32768 (0xFFFF8000).
        // Hmm, let me take a completely different approach and use NEG on MinInt32.
        // To get MinInt32: MOVI -1; DEC; ... no. MOVI -32768, MUL 32768... 32768 doesn't fit.
        // OK let me just use ADD to overflow. Start with a large value using ADDI chaining.
        // R1 = 32767. ADDI R1, 32767 -> 65534. ADDI R1, 32767 -> 98301. Too slow.
        // Better: R1 = 32767, R2 = R1, MUL R2 R2 R1 -> 32767^2 = 1073676289.
        // Then ADDI 1073676289 + something... still slow.
        // The cleanest approach: build 0x7FFF0000 + 0xFFFF.
        // 0x7FFF = 32767. 0x7FFF0000 = 32767 * 65536.
        // We need 65536. We can get it: MOVI R1 1, repeatedly: ADD R2 R1 R1 (2), ADD R2 R2 R2 (4), ... (16).
        // That's 16 MUL/ADD ops but works.
        //
        // Actually, simplest of all: compute -2 via MOVI(-1) then DEC, then NOT(-2) = 1... no.
        // NOT(-2) = ~0xFFFFFFFE = 0x00000001 = 1. Not useful.
        // NOT(-1) = ~0xFFFFFFFF = 0. NOT(0) = -1.
        // NOT(1) = ~0x00000001 = -2.
        //
        // OK I'll use a different strategy. I'll construct the large value using
        // NEG. NEG(-2) = 2. Not helpful for MaxInt32.
        //
        // Let me just build it step by step. The simplest way to get 65536:
        // R1 = 256; R2 = R1 * R1 = 65536.
        // Then: MaxInt32 = 32767 * 65536 + 65535.
        // R1 = 256 (fits in int16)
        // R2 = R1 * R1 = 65536
        // R3 = 32767 (fits in int16)
        // R4 = R3 * R2 = 32767 * 65536 = 2147418112
        // R5 = R2 - 1 = 65535
        // R6 = R4 + R5 = 2147483647 = MaxInt32
        // R7 = 1
        // R8 = R6 + R7 = should wrap to MinInt32

        prog := assemble(
                EncodeMOVI(1, 256),                // R1 = 256
                EncodeE(OP_MUL, 2, 1, 1),         // R2 = 256*256 = 65536
                EncodeMOVI(3, 32767),              // R3 = 32767
                EncodeE(OP_MUL, 4, 3, 2),         // R4 = 32767 * 65536 = 2147418112
                EncodeMOVI(5, 1),                  // R5 = 1
                EncodeE(OP_SUB, 6, 2, 5),         // R6 = 65536 - 1 = 65535
                EncodeE(OP_ADD, 7, 4, 6),         // R7 = 2147418112 + 65535 = 2147483647 = MaxInt32
                EncodeMOVI(8, 1),                  // R8 = 1
                EncodeE(OP_ADD, 9, 7, 8),         // R9 = MaxInt32 + 1 = MinInt32 (overflow wraps)
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        // Verify we built MaxInt32 correctly.
        if v.Registers[7] != math.MaxInt32 {
                t.Fatalf("expected R7=MaxInt32 (%d), got %d", math.MaxInt32, v.Registers[7])
        }
        // MaxInt32 + 1 wraps to MinInt32.
        if v.Registers[9] != math.MinInt32 {
                t.Fatalf("expected R9=MinInt32 (%d), got %d", math.MinInt32, v.Registers[9])
        }
}

func TestArithmeticOverflow_SubWraps(t *testing.T) {
        // Build MinInt32 first using the same technique as above.
        // Then subtract 1 to get MaxInt32.
        prog := assemble(
                EncodeMOVI(1, 256),                // R1 = 256
                EncodeE(OP_MUL, 2, 1, 1),         // R2 = 65536
                EncodeMOVI(3, 32767),              // R3 = 32767
                EncodeE(OP_MUL, 4, 3, 2),         // R4 = 2147418112
                EncodeMOVI(5, 1),                  // R5 = 1
                EncodeE(OP_SUB, 6, 2, 5),         // R6 = 65535
                EncodeE(OP_ADD, 7, 4, 6),         // R7 = MaxInt32 = 2147483647
                EncodeE(OP_ADD, 8, 7, 5),         // R8 = MinInt32 = -2147483648
                EncodeE(OP_SUB, 9, 8, 5),         // R9 = MinInt32 - 1 should wrap to MaxInt32
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[9] != math.MaxInt32 {
                t.Fatalf("expected R9=MaxInt32 (%d), got %d", math.MaxInt32, v.Registers[9])
        }
}

func TestArithmeticOverflow_MulWraps(t *testing.T) {
        // 65536 * 65536 = 2^32 which wraps to 0 for int32
        prog := assemble(
                EncodeMOVI(1, 256),          // R1 = 256
                EncodeE(OP_MUL, 2, 1, 1),   // R2 = 256*256 = 65536
                EncodeE(OP_MUL, 3, 2, 2),   // R3 = 65536*65536 = 0 (overflow)
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[3] != 0 {
                t.Fatalf("expected R3=0 (overflow), got %d", v.Registers[3])
        }
}

// ---------------------------------------------------------------------------
// Flags sign flag
// ---------------------------------------------------------------------------

func TestFlags_SignAfterNegative(t *testing.T) {
        prog := assemble(
                EncodeMOVI(1, -1),
                EncodeUnary(OP_NEG, 1), // -(-1) = 1, Sign=false
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Flags.Sign {
                t.Fatal("expected Sign=false after NEG(-1)=1")
        }
}

// ---------------------------------------------------------------------------
// Running off end of bytecode treated as HALT
// ---------------------------------------------------------------------------

func TestRunoffEnd_Halts(t *testing.T) {
        prog := EncodeMOVI(1, 5) // 4 bytes, no HALT
        v := New(prog)
        v.Execute()

        if !v.Halted {
                t.Fatal("expected Halted=true after running off end")
        }
        if v.Registers[1] != 5 {
                t.Fatalf("expected R1=5, got %d", v.Registers[1])
        }
}

// ---------------------------------------------------------------------------
// Nested CALL/RET
// ---------------------------------------------------------------------------

func TestNestedCallRet(t *testing.T) {
        // Test two levels of CALL/RET.
        //   0: MOVI R1 1            (4 bytes)
        //   4: CALL subroutine1      (3 bytes) -> jump to 14
        //   7: MOVI R1 100           (4 bytes)
        //  11: HALT                   (1 byte)
        //  12: ADDI R1 10            (4 bytes) -- pad
        //
        // Actually let me lay this out precisely.
        //
        //   0: MOVI R1 1             (4) [0-3]
        //   4: CALL +10              (3) [4-6] -> PC=4+10=14. Push 7.
        //   7: MOVI R1 100           (4) [7-10]
        //  11: HALT                   (1) [11]
        //  12: ADDI R1 10            (4) [12-15] -- subroutine1 body
        //  16: CALL +8               (3) [16-18] -> PC=16+8=24. Push 19.
        //  19: MOVI R2 0             (4) [19-22] -- runs after inner RET
        //  23: HALT                   (1) [23]
        //  24: ADDI R1 100           (4) [24-27] -- subroutine2 body
        //  28: RET                    (1) [28]
        prog := assemble(
                EncodeMOVI(1, 1),           // [0-3]
                EncodeBranch(OP_CALL, 8),   // [4-6] -> jump to 12. Push 7.
                EncodeMOVI(1, 100),         // [7-10]
                []byte{OP_HALT},            // [11]
                EncodeADDI(1, 10),          // [12-15] subroutine1: R1 = 1+10 = 11
                EncodeBranch(OP_CALL, 8),   // [16-18] -> jump to 24. Push 19.
                EncodeMOVI(2, 0),           // [19-22]
                []byte{OP_HALT},            // [23]
                EncodeADDI(1, 100),         // [24-27] subroutine2: R1 = 11+100 = 111
                []byte{OP_RET},             // [28]
        )
        v := New(prog)
        v.Execute()

        // Trace:
        // 1. MOVI R1 1 -> R1=1
        // 2. CALL +8 -> push 7, PC=12
        // 3. ADDI R1 10 -> R1=11
        // 4. CALL +8 -> push 19, PC=24
        // 5. ADDI R1 100 -> R1=111
        // 6. RET -> pop 19, PC=19
        // 7. MOVI R2 0 -> R2=0
        // 8. HALT -> stopped
        if v.Registers[1] != 111 {
                t.Fatalf("expected R1=111, got %d", v.Registers[1])
        }
        if v.Registers[2] != 0 {
                t.Fatalf("expected R2=0, got %d", v.Registers[2])
        }
}

// ---------------------------------------------------------------------------
// Stack depth tracking
// ---------------------------------------------------------------------------

func TestStackDepth(t *testing.T) {
        // Push 5 values, then pop them all.
        prog := assemble(
                EncodeMOVI(1, 10),
                EncodeMOVI(2, 20),
                EncodeMOVI(3, 30),
                EncodeUnary(OP_PUSH, 1),
                EncodeUnary(OP_PUSH, 2),
                EncodeUnary(OP_PUSH, 3),
                EncodeUnary(OP_POP, 4),  // should get 30
                EncodeUnary(OP_POP, 5),  // should get 20
                EncodeUnary(OP_POP, 6),  // should get 10
                []byte{OP_HALT},
        )
        v := New(prog)
        v.Execute()

        if v.Registers[4] != 30 || v.Registers[5] != 20 || v.Registers[6] != 10 {
                t.Fatalf("expected R4=30 R5=20 R6=10, got R4=%d R5=%d R6=%d",
                        v.Registers[4], v.Registers[5], v.Registers[6])
        }
}
