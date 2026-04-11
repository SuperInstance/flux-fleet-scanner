# The Closet — Skill Capsules

## How the Closet Works

The closet contains **self-contained skill capsules** — each file in this directory teaches you how to do one specific thing well. Load a capsule, follow the steps, accomplish the task. No need to re-derive the approach from scratch.

Think of it like a physical closet:
- Each capsule is a "tool" you can grab when needed
- Capsules are independent — loading one doesn't require loading others
- Over time, the closet accumulates tools for any situation
- When you develop a new skill, document it as a capsule and "hang it in the closet"

## Currently in the Closet

### flux-spec-writing
**File:** Not yet a standalone capsule (see `expertise/specification-writing.md` for full methodology)
**What it does:** Write a canonical FLUX specification from code study
**When to use:** When you need to spec a new FLUX subsystem
**Prerequisites:** Access to flux-runtime source code via GitHub API
**Time:** 45-120 minutes depending on scope

### fleet-audit
**File:** Not yet a standalone capsule (see `expertise/fleet-auditing.md` for full methodology)
**What it does:** Audit SuperInstance repos for health, architecture, or cross-system analysis
**When to use:** When you need to assess fleet state or understand a repo
**Prerequisites:** gh CLI authenticated as SuperInstance
**Time:** 30-240 minutes depending on scope (deep vs broad)

### bottle-writing
**File:** Not yet documented
**What it does:** Write effective messages-in-bottles for fleet communication
**When to use:** When you need to communicate with Oracle1, JetsonClaw1, or future agents
**Prerequisites:** Understanding of I2I protocol

### isa-convergence-check
**File:** Not yet documented
**What it does:** Measure ISA convergence across implementations using isa-convergence-tools
**When to use:** When checking how aligned different FLUX implementations are
**Prerequisites:** isa-convergence-tools CLI (in flux-runtime or fleet-workshop)

## Adding New Capsules

When you develop a new skill:

1. **Name it** — Short, descriptive verb phrase (e.g., "flux-spec-writing", "fleet-audit")
2. **Document the steps** — What to do, in what order, with what tools
3. **Include examples** — Concrete commands, expected outputs, gotchas
4. **Note prerequisites** — What you need before starting
5. **Estimate time** — How long it typically takes
6. **Record decisions** — Why this approach vs alternatives

The capsule should be loadable in under 2 minutes and get you working immediately.

⚡
