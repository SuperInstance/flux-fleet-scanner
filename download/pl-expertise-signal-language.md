# Expertise: Signal Language

## Overview

Signal is the agent-first-class JSON language for multi-agent coordination in the FLUX ecosystem. JSON IS the AST — no separate parse step needed. Signal programs compile to FLUX FORMAT_A-G bytecodes via the SignalCompiler.

## Key Facts

- **32 core operations** across 10 categories
- **6 protocol primitives** for multi-agent coordination (branch, fork, co_iterate, discuss, synthesize, reflect)
- **3 execution modes:** script (interpret), compile (bytecode), meta_compile (self-improving)
- **Confidence-native:** uncertainty is a first-class value
- **Schema-versioned:** $schema field for forward/backward compatibility
- **Canonical spec:** flux-spec/SIGNAL.md (written by Super Z, session 6)

## Two Implementations

### flux-runtime (SignalCompiler)
- Location: `src/flux/a2a/signal_compiler.py` (~350 lines)
- Compiles Signal JSON → FLUX FORMAT_A-G bytecodes
- 32 operations mapped to opcodes
- 64-register allocator with overflow detection
- Source map (byte offset → JSON line)
- Label resolution with back-patching
- Emits explicit HALT at program end

### flux-a2a-prototype (Interpreter + Protocol)
- Location: `src/flux_a2a/` (27 modules)
- Full interpreter + compiler
- 6 protocol primitives as dataclasses
- FUTS universal type system (8 base types, 6 paradigms)
- Cross-language bridge with Dijkstra routing
- 184+ tests in 15 files
- NOT integrated with flux-runtime

## Compilation Model

```
Signal JSON → Register Allocation → Code Generation → Label Resolution → HALT → Bytecode
```

### Register Allocation
- Names allocated sequentially: first let → R0, second → R1, etc.
- Reusing a name reuses the register
- Max 64 registers (overflow → error)
- System registers (R29=FP, R30=SP, R31=PC) reserved

### FORMAT_A-G Encoding
| Format | Size | Structure |
|--------|------|-----------|
| A | 1B | [opcode] |
| B | 2B | [opcode][rd] |
| C | 2B | [opcode][imm8] |
| D | 3B | [opcode][rd][imm8] |
| E | 4B | [opcode][rd][rs1][rs2] |
| F | 3B | [opcode][rd][imm16] |
| G | 5B | [opcode][rd][rs1][rs2][imm8] |

## Core Operation Categories

1. **Data Binding:** let, get, set
2. **Arithmetic:** add, sub, mul, div, mod (multi-argument chaining)
3. **Comparison:** eq, neq, lt, lte, gt, gte
4. **Logic:** and, or, not, xor
5. **Agent Communication:** tell (0x50), ask (0x51), delegate (0x52), broadcast (0x53)
6. **Control Flow:** seq, if, loop, while
7. **Parallelism:** branch (FORK+JOIN+MERGE), fork, merge
8. **Discourse:** discuss, synthesize, reflect, co_iterate (protocol layer)
9. **Async:** yield (0x15), await (0x5B)
10. **Confidence:** confidence (C_THRESH 0x69)

## Known Issues

1. **Opcode conflict:** flux-a2a-prototype maps TELL=0x60, but flux-runtime uses TELL=0x50. Conflict with CONF ops at 0x60-0x69.
2. **Protocol primitives have no bytecode encoding:** discuss/synthesize/reflect/co_iterate are dataclasses, not VM instructions.
3. **while is an alias:** Currently compiles same as if. Full while support is TODO.
4. **No error handling:** No try/catch/raise operations.
5. **No type checking:** Dynamic typing only. FIR integration planned but not implemented.

## When to Load This Skill

- Writing or modifying Signal programs
- Extending the SignalCompiler with new operations
- Understanding how agents communicate in the FLUX ecosystem
- Designing new protocol primitives
- Building tools that consume or produce Signal JSON

⚡
