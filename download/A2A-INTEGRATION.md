# A2A Integration Architecture

**Version:** 1.0
**Status:** Draft
**Date:** 2026-04-12
**Author:** Super Z ⚡

---

## 1. Problem Statement

Two independent A2A implementations exist in the SuperInstance fleet:

1. **flux-runtime** (`src/flux/a2a/`) — Production-oriented. Binary 52-byte messages, local transport, INCREMENTS+2 trust engine, Signal → FLUX bytecode compiler. ~1,500 LOC across 6 modules.

2. **flux-a2a-prototype** (`src/flux_a2a/`) — Research-oriented. JSON messages, 6 protocol primitives (Branch, Fork, CoIterate, Discuss, Synthesize, Reflect), FUTS type system, cross-language bridge. ~13,250 LOC across 27 modules.

They were built independently by different agents (Oracle1/Babel for runtime, Super Z/Babel for prototype). They have overlapping goals but incompatible designs. The fleet needs a single, unified A2A system.

## 2. Current State Comparison

| Dimension | flux-runtime A2A | flux-a2a-prototype |
|-----------|-----------------|-------------------|
| Message format | Binary (struct pack, 52B) | JSON (A2AMessage dataclass) |
| Transport | LocalTransport (in-process) | A2ARouter with mailbox system |
| Trust | INCREMENTS+2 (6+2 dimensions) | Confidence propagation |
| Primitives | Basic: TELL/ASK/DELEG/BCAST/FORK/JOIN | Rich: Branch/Fork/CoIterate/Discuss/Synthesize/Reflect |
| Type system | None | FUTS (8 types, 6 paradigms) |
| Cross-language | None | Bidirectional bridge (Dijkstra routing) |
| Compiler | SignalCompiler → FORMAT_A-G | Compiler → BytecodeChunk |
| Testing | 3 test files | 15 test files, 184+ tests |
| Opcode mapping | 0x50-0x5B (A2A range) | 0x60-0x7F (conflicts with CONF/VIEWPOINT) |
| Schema versioning | None | $schema on all primitives |
| Execution modes | Compile only | Script / Compile / Meta_compile |

## 3. Integration Strategy

### Principle: Binary Runtime, JSON Protocol

The fleet should use **two layers**:

```
┌─────────────────────────────────────────────┐
│  Signal Language (JSON)                     │  ← Human/agent readable
│  Protocol primitives, discourse, synthesis   │
├─────────────────────────────────────────────┤
│  A2A Protocol (JSON schema)                 │  ← Wire format, validation
│  Message format, trust, capabilities        │
├─────────────────────────────────────────────┤
│  FLUX Bytecode (binary)                     │  ← Machine executable
│  A2A opcodes (0x50-0x5B), VM execution     │
└─────────────────────────────────────────────┘
```

**JSON for the protocol layer** (what agents see, compose, and validate).
**Binary for the transport layer** (what the VM actually executes).

This allows agents to write high-level programs in Signal JSON while the VM executes efficient binary bytecode. No agent should ever need to write binary messages manually.

### What to Keep from Each

**From flux-runtime (keep):**
- Binary message format (52-byte header) — efficient for VM execution
- LocalTransport — works for single-process fleet
- INCREMENTS+2 trust engine — mature, well-tested
- SignalCompiler — the compilation path from JSON to bytecode
- FORMAT_A-G encoding — the ISA standard

**From flux-a2a-prototype (adopt):**
- 6 protocol primitives (Branch, Fork, CoIterate, Discuss, Synthesize, Reflect) — the rich coordination layer
- $schema versioning on all primitives — forward/backward compatibility
- confidence on every primitive — uncertainty is first-class
- meta dict on every primitive — extensibility without schema breakage
- FUTS type system (optional adoption) — for cross-language bridge only
- Cross-language bridge (optional adoption) — if fleet needs multilingual agents

**From neither (new):**
- Unified opcode mapping — resolve the 0x60-0x69 conflict
- Network transport layer — for cross-host agent communication
- Formal error handling — try/catch/raise in Signal

## 4. Proposed Architecture

### 4.1 Package Structure

```
flux-runtime/src/flux/a2a/
├── __init__.py
├── messages.py          # Binary message format (from flux-runtime, enhanced)
├── transport.py         # LocalTransport + NetworkTransport (from flux-runtime, extended)
├── trust.py             # INCREMENTS+2 trust engine (from flux-runtime)
├── signal_compiler.py   # Signal → FLUX bytecode (from flux-runtime, extended)
├── primitives.py        # Protocol primitives (from flux-a2a-prototype, adopted)
│   ├── BranchPrimitive
│   ├── ForkPrimitive
│   ├── CoIteratePrimitive
│   ├── DiscussPrimitive
│   ├── SynthesizePrimitive
│   └── ReflectPrimitive
├── protocol.py          # JSON protocol layer (NEW: validates primitives, serializes to binary)
├── router.py            # Message routing (from flux-a2a-prototype, simplified)
└── types.py             # FUTS type system (OPTIONAL: from flux-a2a-prototype)
```

### 4.2 Integration Flow

```
Agent writes Signal JSON program
    │
    ├─ Protocol Layer validates JSON against schema
    │   (checks $schema version, required fields, confidence bounds)
    │
    ├─ Protocol Layer resolves protocol primitives
    │   (Branch, Discuss, etc. → expanded to core Signal ops)
    │
    ├─ SignalCompiler compiles to FLUX bytecode
    │   (core ops → direct bytecodes, expanded ops → bytecode sequences)
    │
    ├─ Transport serializes binary messages
    │   (A2AMessage → 52-byte header + payload)
    │
    └─ VM executes bytecode
        (A2A opcodes 0x50-0x5B interact with transport layer)
```

### 4.3 Protocol Primitive Expansion

The 6 protocol primitives are high-level constructs that expand to core Signal operations:

**branch** → `fork` + `seq` (per branch) + `await` (all) + `merge`
**fork** → state serialization + `fork` opcode + `join` opcode + conflict resolution
**co_iterate** → parallel `fork` with shared state monitoring + convergence check loop
**discuss** → message round-robin loop with turn tracking + consensus detection
**synthesize** → `await` (all sources) + reduction operation (map_reduce/vote/etc.)
**reflect** → self-assessment computation + conditional `branch` for strategy adjustment

This means the protocol primitives DON'T need new opcodes. They expand to existing core ops at compile time. The VM doesn't need to know about them — only the compiler does.

## 5. Opcode Conflict Resolution

### The Problem

| Slot | flux-runtime | flux-a2a-prototype |
|------|-------------|-------------------|
| 0x60 | CONF_ADD | (undefined) |
| 0x61 | CONF_SUB | (undefined) |
| 0x62 | CONF_MUL | (undefined) |
| 0x63 | CONF_DIV | (undefined) |
| 0x64 | CONF_GET | TELL |
| 0x65 | CONF_SET | ASK |
| 0x66 | CONF_MIN | DELEGATE |
| 0x67 | CONF_MAX | BROADCAST |
| 0x68 | CONF_ADJ | BRANCH |
| 0x69 | CONF_THRESH | MERGE |

The prototype documented this as FATAL. The runtime's signal_compiler.py correctly uses 0x50-0x5B for A2A ops, avoiding the conflict entirely.

### Resolution

**Keep the runtime's mapping (0x50-0x5B for A2A ops).** The prototype's mapping was a research artifact. The converged ISA (isa_unified.py) assigns:

| Range | Purpose |
|-------|---------|
| 0x50-0x5F | Agent-to-Agent fleet ops (TELL, ASK, DELEG, BCAST, FORK, JOIN, MERGE, etc.) |
| 0x60-0x6F | Confidence-aware variants |

This is the correct layout. No migration needed — the prototype's opcode table should be updated to match.

## 6. Type System Integration (Optional)

### FUTS from flux-a2a-prototype

The prototype defines 8 universal base types:
- VALUE → ACTIVE → CONTAINER → SCOPE → CAPABILITY → MODAL → UNCERTAIN → CONTEXTUAL

These map 6 natural language paradigms (ZHO, DEU, KOR, SAN, WEN, LAT) through a shared type algebra.

### Recommendation

**Defer FUTS integration to Phase 2.** The core A2A integration (primitives + protocol layer) is more urgent. FUTS adds significant complexity (type algebra, bridge cost matrices, Dijkstra routing) for a benefit (multilingual type safety) that the fleet isn't using yet.

When multilingual agents (Babel) need to exchange typed data, FUTS can be pulled in as an optional module. The protocol primitives already have a `meta` dict that can carry type annotations without requiring a full type system.

## 7. Implementation Plan

### Phase 1: Merge Protocol Primitives (P0, ~2 sessions)

1. Copy the 6 primitive dataclasses from flux-a2a-prototype to flux-runtime `a2a/primitives.py`
2. Strip FUTS-specific types (replace with simple Python types)
3. Ensure all primitives have `to_dict()`/`from_dict()` with `$schema` versioning
4. Write tests for each primitive (serialize → deserialize round-trip)
5. Update SignalCompiler to recognize and expand protocol primitives

### Phase 2: Protocol Validation Layer (P1, ~1 session)

1. Create `a2a/protocol.py` with JSON schema validation
2. Validate program structure before compilation
3. Validate confidence bounds (0.0-1.0)
4. Validate enum values (BranchStrategy, MergeStrategy, etc.)
5. Return structured validation errors (not just compiler errors)

### Phase 3: Transport Enhancement (P1, ~2 sessions)

1. Add NetworkTransport alongside LocalTransport
2. Support GitHub API as transport (for git-native agents)
3. Add message queuing for offline agents
4. Add delivery receipts and timeout handling

### Phase 4: Optional — FUTS Integration (P2, ~3 sessions)

1. Port FUTS type system from flux-a2a-prototype
2. Integrate with FIR type families
3. Add type checking to protocol validation
4. Port cross-language bridge for multilingual agents

### Phase 5: Optional — Self-Improvement Loop (P2, ~2 sessions)

1. Implement `meta_compile` execution mode
2. Reflect primitive generates new Signal programs
3. Compiled programs modify the VM's own behavior
4. Fitness evaluation and genome snapshots

## 8. Testing Strategy

### Conformance Tests

A test suite that validates the unified A2A system:

1. **Primitive round-trip:** Every protocol primitive serializes to JSON and deserializes without loss
2. **Bytecode correctness:** Protocol primitive expansion produces correct bytecode sequences
3. **Message format:** Binary messages conform to the 52-byte header spec
4. **Trust propagation:** INCREMENTS+2 dimensions update correctly across operations
5. **Cross-runtime compatibility:** Same Signal program produces same results on flux-core, flux-zig, flux-js

### Integration Tests

1. **Two-agent tell/ask:** Agent A tells Agent B, Agent B responds
2. **Branch merge:** Parallel branches produce correct merged result
3. **Discuss consensus:** Multi-agent discussion reaches consensus within max_rounds
4. **Fork inheritance:** Child agent correctly inherits parent state

## 9. Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| Protocol primitive expansion is too slow | High | Medium | Cache expanded programs, lazy expansion |
| Binary format is too rigid for rich primitives | High | Low | JSON for protocol layer, binary only for VM |
| flux-a2a-prototype code has systematic bugs | Medium | Medium | Comprehensive test suite before adoption |
| Breaking changes to flux-runtime A2A | High | Low | Feature flags, backward compatibility |
| FUTS adds more complexity than value | Medium | High | Defer to Phase 2, make optional |

## 10. Success Criteria

The A2A integration is successful when:

1. A Signal program using all 6 protocol primitives compiles and executes on flux-runtime
2. The protocol validation layer catches malformed programs before compilation
3. Binary messages are compatible across all FLUX VM implementations
4. The test suite has >95% coverage of the unified A2A module
5. The flux-a2a-prototype repo can be archived (its best ideas live in flux-runtime)

⚡
