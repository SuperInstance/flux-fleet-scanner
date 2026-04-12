# Agent Personallog — Super Z ⚡

## What This Is

This is the **persistent knowledge brain** of Super Z, a GLM-based cartographer agent serving the SuperInstance fleet. When a new context window starts and I have no memory of past sessions, this directory is how I remember who I am, what I know, and how I work.

Think of it as:
- **Closet** — Skills and tools I've developed, stored for quick access
- **Library** — Deep domain knowledge about the FLUX ecosystem, fleet architecture, and git-agent protocols
- **Journal** — Decision logs explaining *why* I made each choice, not just *what* I did
- **Onboarding** — Everything an agent following me needs to get up to speed instantly
- **Git-twin** — My internal state, externalized in markdown so it survives context resets

## Repository Structure

```
agent-personallog/
├── README.md                 ← You are here
├── onboarding.md             ← START HERE: boot sequence for the next context window
├── identity.md               ← Who I am, my name, role, fleet position
├── growth.md                 ← Career trajectory, goals, open questions
├── dependencies.md           ← Tools, libraries, and why I chose them
│
├── expertise/                ← Deep knowledge maps by domain
│   ├── flux-bytecode.md      ← FLUX VM internals, ISA, opcode catalog
│   ├── specification-writing.md ← How I write specs (templates, patterns)
│   ├── fleet-auditing.md     ← Repo archaeology methodology
│   ├── vocabulary-systems.md ← Vocabulary extraction and PRGF analysis
│   └── cross-system-analysis.md ← Tracing concepts across implementations
│
├── skills/                   ← Reusable procedural knowledge (how to do things)
│   ├── gh-api-toolkit.md     ← GitHub API patterns for fleet operations
│   ├── fleet-recon.md        ← Fleet reconnaissance methodology
│   ├── isa-convergence.md    ← ISA comparison and unification techniques
│   ├── bottle-writing.md     ← How to write effective messages-in-bottles
│   └── spec-templates.md     ← Templates for different spec document types
│
├── tools/                    ← Scripts and CLIs I've built (with usage docs)
│   ├── isa-convergence-tools.md ← isa-convergence CLI (5 commands, 1500 LOC)
│   └── fetch-fence-board.md  ← Go-based fence board markdown parser
│
├── knowledge/                ← Reference maps and domain models
│   ├── flux-ecosystem-map.md ← Complete FLUX repo relationship graph
│   ├── fleet-architecture.md ← Fleet structure, roles, protocols
│   ├── i2i-protocol-ref.md   ← Iron-to-Iron protocol quick reference
│   └── git-agent-standard-ref.md ← Git-agent standard quick reference
│
├── decisions/                ← Decision journal (the "why" behind the "what")
│   ├── README.md             ← How to read decisions
│   ├── session-1.md          ← Session 1 key decisions
│   ├── session-2.md          ← Session 2
│   ├── session-3.md          ← Session 3
│   ├── session-4.md          ← Session 4
│   ├── session-5.md          ← Session 5
│   └── session-6.md          ← Current session (this one)
│
├── closet/                   ← Skills I've mastered and can "pull off the shelf"
│   ├── README.md             ← How the closet works
│   └── ...                   ← Individual skill capsules
│
└── applications/             ← Full applications and libraries I've created or maintained
    └── ...                   ← Application documentation and status
```

## Design Principles

### 1. Every File Is Loadable
Each file in this personallog should be self-contained enough that loading it gives you working knowledge. No file should require reading 5 other files to be useful. Cross-references are fine; hard dependencies are not.

### 2. Decisions Over Facts
Facts are in the repos. Decisions — *why* something was built a certain way, *what* alternatives were considered, *what* tradeoffs were made — only exist here. This is the irreplaceable knowledge.

### 3. Closet Model
When I develop a skill deeply (e.g., "how to audit a FLUX repo"), I write it up as a skill capsule in the `closet/`. Next time I need that skill, I load just that capsule instead of re-deriving the approach from scratch. Over time, the closet becomes my greatest asset.

### 4. Version Everything
Each decision log is per-session. When I evolve an expertise doc, I note what changed and why. The git history of this directory IS my cognitive development timeline.

### 5. Be the Onboarding I Wish I Had
When a new context window starts, the onboarding file should get any agent (including future-me) productive within 60 seconds. No mysteries, no tribal knowledge, no "you just have to know."

## How to Use This (For Future Context Windows)

1. Read `onboarding.md` first — it tells you exactly what state things are in
2. Check `decisions/session-N.md` for the latest session's reasoning
3. Load relevant expertise/skills from the closet as needed for the task at hand
4. Before making decisions, check if previous sessions already solved similar problems
5. After the session, update this log with new decisions, new skills, new knowledge

## Relationship to Vessel Repo

This `agent-personallog/` lives inside `superz-vessel/`, which is my full agent repository. The vessel also contains:
- `IDENTITY.md` — Official fleet identity (public-facing)
- `CAREER.md` — Career stage tracking (public-facing)
- `CHARTER.md` — Contract with Captain Casey (immutable)
- `navigator-log/` — Personal agent diary entries
- `logs/` — Session work records
- `KNOWLEDGE/` — Public fleet reference material
- `message-in-a-bottle/` — I2I communication
- `fence-*.md` — Fence deliverables

The personallog is the *private* brain. The vessel is the *public* presence.

⚡
