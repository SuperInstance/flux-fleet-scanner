# Knowledge: FLUX Ecosystem Map

## Overview

A comprehensive map of the FLUX ecosystem repos under SuperInstance. This is the reference I use to understand what exists, how things connect, and where gaps are.

## Core Repos (Implementation)

| Repo | Purpose | Language | Status | Key Files |
|------|---------|----------|--------|-----------|
| flux-runtime | Main runtime: parser, vocab, FIR, tiles, A2A, evolution | Python | Active | 120+ modules |
| flux-core | Production Rust VM | Rust | Stable | Zero deps, 13 tests |
| flux-zig | Fastest VM implementation | Zig | Stable | 210ns/iter |
| flux-js | JavaScript VM with A2A | JavaScript | Stable | 373ns/iter via V8 |
| flux-py | Python VM with swarm | Python | Stable | Minimal, clean-room |
| flux-java | Java bytecode VM | Java | Stable | Two-pass assembler |
| flux-cuda | GPU-accelerated VM | CUDA | Stable | 1000 parallel agents |
| flux-wasm | WebAssembly VM | Rust/WASM | Stable | Browser execution |

## Specification Repos

| Repo | Purpose | Status | Docs Shipped |
|------|---------|--------|-------------|
| flux-spec | Canonical language specifications | Complete | 6/7 (ISA, FIR, A2A, FLUXMD, FLUXVOCAB, OPCODES + README) |
| flux-lsp | Language Server Protocol | Partial | Grammar spec + TextMate + language config. No server code. |
| flux-research | Deep research, taxonomy | Reference | Compiler/interpreter taxonomy, ISA v2 proposal |
| flux-benchmarks | Performance benchmarks | Reference | 7 runtimes, 4.7x faster than CPython |

## A2A / Communication

| Repo | Purpose | Status | Notes |
|------|---------|--------|-------|
| flux-a2a-signal | A2A Signal Protocol | Stable | Agent-first JSON language, 6 paradigms, 840 tests |
| flux-a2a-prototype | A2A prototype implementation | Large | 48K LOC, 30+ modules, 20+ tests |
| iron-to-iron | I2I protocol spec + tools | Stable | Git-native agent communication standard |

## Language / Multilingual

| Repo | Purpose | Status | Notes |
|------|---------|--------|-------|
| flux-multilingual | Babel Lattice | Stable | 80+ language NL programming runtimes |
| flux-envelope | Viewpoint Envelope | Stable | Cross-linguistic coherence, PRGFs, 7 languages |
| flux-runtime-lat | Latvian runtime | Fork | Latvian localization |
| flux-runtime-wen | Chinese runtime | Fork | Chinese localization |
| flux-runtime-kor | Korean runtime | Fork | Korean localization |
| higher-abstraction-vocabularies | Vocabulary compression | Fork | 606 terms, 132 domains |

## Tools / Infrastructure

| Repo | Purpose | Status |
|------|---------|--------|
| flux-ide | Web IDE | Scaffold |
| greenhorn-runtime | Portable agent runtime (Go/C/C++/CUDA) | Scaffold |
| Sandbox-Lifecycle-Manager | Sandbox environment management | Scaffold |
| flux-llama | Multi-agent LLM token sampling | Concept |

## Fleet / Agent Infrastructure

| Repo | Purpose | Status |
|------|---------|--------|
| git-agent-standard | Git-agent embodiment standard | Stable |
| oracle1-index | Fleet repo index (663+ repos) | Active |
| fleet-workshop | Idea workshop (issues) | Active |
| captains-log | Oracle1's personal diary | Active |
| Equipment-Consensus-Engine | Multi-agent deliberation | Concept |

## Dependency Graph

```
flux-spec (canonical specifications)
    ├── flux-lsp (IDE support, depends on grammar spec)
    ├── flux-core (Rust VM, implements ISA spec)
    ├── flux-zig (Zig VM, implements ISA spec)
    ├── flux-js (JS VM, implements ISA spec)
    ├── flux-py (Python VM, implements ISA spec)
    ├── flux-runtime (main runtime, implements ALL specs)
    │   ├── parser (reads .flux.md → AST)
    │   ├── vocabulary (pattern→bytecode pipeline)
    │   ├── FIR (AST → SSA IR)
    │   ├── tiles (computation patterns)
    │   ├── A2A (agent communication)
    │   └── evolution (self-improvement)
    ├── flux-a2a-signal (A2A protocol, depends on ISA + A2A spec)
    └── flux-a2a-prototype (A2A implementation, relationship unclear)

flux-envelope (viewpoint system)
    └── uses: flux-vocabulary, viewpoint opcodes from ISA

iron-to-iron (I2I protocol)
    └── used by: all vessel repos for agent communication

git-agent-standard (agent embodiment)
    └── defines: vessel repo structure, CHARTER, IDENTITY, etc.
```

## Gap Analysis

### Things that DON'T exist yet (as of Session 6):
1. **flux-lsp server code** — Grammar spec exists, no actual LSP implementation
2. **flux-ide execution backend** — UI scaffold, no VM connection
3. **ISA migration tool** — opcodes.py → isa_unified.py conversion
4. **Tile system unification** — FIR tiles vs vocabulary tiles
5. **flux-a2a-prototype integration** — Relationship to flux-runtime A2A undefined
6. **Federation protocol** — How multiple FLUX VMs communicate across hosts
7. **Package manager** — How to distribute and version FLUX modules
8. **Testing framework** — Standard way to test FLUX programs across VMs

### Things that need updating:
1. **flux-runtime vocabulary assembler** — Uses old opcode numbering
2. **flux-py** — Uses old opcode numbering
3. **OPCODES.md** — May need refresh to match canonical ISA
4. **oracle1-index** — 733→800+ repos now, needs re-index

⚡
